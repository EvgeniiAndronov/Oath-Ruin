package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Типы сообщений
type MessageType string

const (
	MsgJoinGame      MessageType = "join_game"
	MsgMoveCharacter MessageType = "move_character"
	MsgUseAbility    MessageType = "use_ability"
	MsgEndTurn       MessageType = "end_turn"
	MsgGameState     MessageType = "game_state"
	MsgError         MessageType = "error"
)

// Сообщения клиента
type ClientMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Сообщения сервера
type ServerMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// MovePayload - перемещение персонажа
type MovePayload struct {
	PlayerID string `json:"player_id"`
	CharacterID string `json:"character_id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
}

// JoinPayload - вход в игру
type JoinPayload struct {
	PlayerName string `json:"player_name"`
	Faction    string `json:"faction"` // "oath" или "ruin"
}

// Игрок
type Player struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Faction  string `json:"faction"`
	Conn     *websocket.Conn
}

// Состояние игры
type GameState struct {
	Players    map[string]*Player `json:"-"`
	GridWidth  int                `json:"grid_width"`
	GridHeight int                `json:"grid_height"`
	Characters []Character        `json:"characters"`
	CurrentTurn string           `json:"current_turn"`
	mu         sync.RWMutex
}

type Character struct {
	ID        string `json:"id"`
	PlayerID  string `json:"player_id"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	HP        int    `json:"hp"`
	MaxHP     int    `json:"max_hp"`
	CanMove   bool   `json:"can_move"`
	CanAct    bool   `json:"can_act"`
}

// Глобальное состояние (для прототипа)
var gameState = &GameState{
	Players:    make(map[string]*Player),
	Characters: make([]Character, 0),
	GridWidth:  16,
	GridHeight: 16,
	CurrentTurn: "",
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Для разработки, в продакшене ограничить
	},
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/health", healthCheck)

	log.Println("🎮 Oath & Ruin Server starting on :8080...")
	log.Printf("📡 WebSocket: ws://localhost:8080/ws")
	log.Printf("❤️  Health: http://localhost:8080/health")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"game":   "Oath & Ruin",
	})
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	player := &Player{
		ID:   generateID(),
		Conn: conn,
	}
	gameState.Players[player.ID] = player

	log.Printf("👤 Player joined: %s", player.ID)

	// Отправить приветствие
	sendMessage(conn, ServerMessage{
		Type: MsgGameState,
		Payload: gameState,
	})

	// Чтение сообщений
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("👤 Player disconnected: %s", player.ID)
			delete(gameState.Players, player.ID)
			break
		}

		handleMessage(player, message)
	}
}

func handleMessage(player *Player, data []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		sendError(player.Conn, "Invalid message format")
		return
	}

	log.Printf("📨 Message from %s: type=%s", player.ID, msg.Type)

	switch msg.Type {
	case MsgJoinGame:
		handleJoinGame(player, msg.Payload)
	case MsgMoveCharacter:
		handleMove(player, msg.Payload)
	case MsgUseAbility:
		handleAbility(player, msg.Payload)
	case MsgEndTurn:
		handleEndTurn(player)
	default:
		sendError(player.Conn, "Unknown message type")
	}
}

func handleJoinGame(player *Player, payload json.RawMessage) {
	var data JoinPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		sendError(player.Conn, "Invalid join payload")
		return
	}

	player.Name = data.PlayerName
	player.Faction = data.Faction

	// Создать тестового персонажа
	character := Character{
		ID:        generateID(),
		PlayerID:  player.ID,
		Name:      data.PlayerName,
		Class:     "warrior",
		X:         1,
		Y:         1,
		HP:        100,
		MaxHP:     100,
		CanMove:   true,
		CanAct:    true,
	}
	gameState.Characters = append(gameState.Characters, character)

	// Рассылка всем игрокам
	broadcast(ServerMessage{
		Type: MsgGameState,
		Payload: gameState,
	})
}

func handleMove(player *Player, payload json.RawMessage) {
	var data MovePayload
	if err := json.Unmarshal(payload, &data); err != nil {
		sendError(player.Conn, "Invalid move payload")
		return
	}

	// Найти персонажа
	for i, char := range gameState.Characters {
		if char.ID == data.CharacterID && char.PlayerID == player.ID {
			// Проверка: может ли двигаться
			if !char.CanMove {
				sendError(player.Conn, "Character already moved")
				return
			}

			// Проверка: валидная ли клетка (упрощённо)
			if data.X < 0 || data.X >= gameState.GridWidth || 
			   data.Y < 0 || data.Y >= gameState.GridHeight {
				sendError(player.Conn, "Invalid position")
				return
			}

			// Обновить позицию
			gameState.Characters[i].X = data.X
			gameState.Characters[i].Y = data.Y
			gameState.Characters[i].CanMove = false

			log.Printf("🚶 Character %s moved to (%d, %d)", data.CharacterID, data.X, data.Y)

			// Рассылка всем
			broadcast(ServerMessage{
				Type: MsgGameState,
				Payload: gameState,
			})
			return
		}
	}

	sendError(player.Conn, "Character not found")
}

func handleAbility(player *Player, payload json.RawMessage) {
	// TODO: Реализовать использование способностей
	log.Println("⚔️ Ability usage not implemented yet")
}

func handleEndTurn(player *Player) {
	// Восстановить действия персонажей игрока
	for i := range gameState.Characters {
		if gameState.Characters[i].PlayerID == player.ID {
			gameState.Characters[i].CanMove = true
			gameState.Characters[i].CanAct = true
		}
	}

	// Передать ход следующему игроку
	// TODO: Реализовать логику очереди ходов

	broadcast(ServerMessage{
		Type: MsgGameState,
		Payload: gameState,
	})
}

func sendMessage(conn *websocket.Conn, msg ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Marshal error:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("Send error:", err)
	}
}

func sendError(conn *websocket.Conn, message string) {
	sendMessage(conn, ServerMessage{
		Type: MsgError,
		Payload: map[string]string{"error": message},
	})
}

func broadcast(msg ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Broadcast marshal error:", err)
		return
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()

	for _, player := range gameState.Players {
		if err := player.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Send to %s error: %v", player.ID, err)
		}
	}
}

func generateID() string {
	// Упрощённая генерация ID (для прототипа)
	return fmt.Sprintf("id_%d", len(gameState.Players)+len(gameState.Characters))
}
