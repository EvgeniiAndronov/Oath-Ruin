extends Node

## Сетевой менеджер для подключения к серверу
## Singleton (autoload)

signal connected()
signal disconnected()
signal message_received(data: Dictionary)
signal error_occurred(error: String)

var socket = WebSocketPeer.new()
var server_url := "ws://localhost:8080/ws"
var is_connected := false

var player_id := ""
var player_name := ""

func _ready():
	set_process(false)  # Отключено по умолчанию

func _process(_delta):
	socket.poll()
	
	var state = socket.get_ready_state()
	
	match state:
		WebSocketPeer.STATE_OPEN:
			if not is_connected:
				is_connected = true
				connected.emit()
				print("✅ Connected to server")
		
		WebSocketPeer.STATE_CLOSED:
			if is_connected:
				is_connected = false
				disconnected.emit()
				print("❌ Disconnected from server")

func connect_to_server():
	if socket.get_ready_state() == WebSocketPeer.STATE_OPEN:
		return
	
	socket.connect_to_url(server_url)
	set_process(true)
	print("🔌 Connecting to ", server_url)

func disconnect_from_server():
	socket.close()
	set_process(false)
	is_connected = false

func send_message(msg_type: String, payload: Dictionary):
	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		error_occurred.emit("Not connected to server")
		return
	
	var message = {
		"type": msg_type,
		"payload": payload
	}
	
	var json = JSON.stringify(message)
	socket.send_text(json)
	print("📤 Sent: ", json)

func join_game(name: String, faction: String):
	player_name = name
	send_message("join_game", {
		"player_name": name,
		"faction": faction
	})

func move_character(character_id: String, x: int, y: int):
	send_message("move_character", {
		"character_id": character_id,
		"x": x,
		"y": y
	})

func use_ability(character_id: String, ability_id: String, target_x: int, target_y: int):
	send_message("use_ability", {
		"character_id": character_id,
		"ability_id": ability_id,
		"target_x": target_x,
		"target_y": target_y
	})

func end_turn():
	send_message("end_turn", {})

func get_pending_messages() -> Array:
	var messages := []
	while socket.get_available_packet_count() > 0:
		var packet = socket.get_packet()
		var json = packet.get_string_from_utf8()
		var data = JSON.parse_string(json)
		if data != null:
			messages.append(data)
			message_received.emit(data)
	return messages
