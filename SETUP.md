# 🚀 Настройка проекта Oath & Ruin

## 📁 Шаг 1: Перенос файлов

Скопируй всё из папки `game_setup` в свой репозиторий:

```bash
# Перейди в папку проекта
cd /Users/evgenii/Develop/game_Oath_\&_Ruin

# Скопируй файлы (если ещё не скопировал)
cp -R /Users/evgenii/Desktop/game_setup/* .

# Проверь структуру
ls -la
```

**Ожидаемая структура:**
```
oath-and-ruin/
├── .gitignore
├── README.md
├── client/
│   ├── project.godot
│   ├── scenes/
│   │   ├── main_menu.tscn
│   │   └── battle.tscn
│   └── scripts/
│       ├── autoload/
│       ├── components/
│       ├── systems/
│       └── ui/
├── server/
│   ├── go.mod
│   └── cmd/
│       └── server/
│           └── main.go
└── docs/
    ├── vision.md
    └── architecture/
        └── adr-001-monorepo.md
```

---

## 🎮 Шаг 2: Настройка Godot

1. **Открой Godot 4.x**
2. **Импортируй проект:**
   - Click "Import"
   - Выбери `client/project.godot`
   - Click "Import & Edit"

3. **Проверь сцены:**
   - Открой `scenes/main_menu.tscn`
   - Открой `scenes/battle.tscn`

4. **Настрой autoload (сеть):**
   - Project → Project Settings → Autoload
   - Добавь `scripts/autoload/network.gd` как `Network`

5. **Запусти:**
   - Click F5 или Play button
   - Должно открыться главное меню

---

## 🔧 Шаг 3: Настройка Go сервера

```bash
# Перейди в папку сервера
cd server

# Инициализируй модуль (если go.mod ещё не создан)
go mod init oath-and-ruin/server

# Установи зависимости
go mod tidy

# Запусти сервер
go run cmd/server/main.go
```

**Ожидаемый вывод:**
```
🎮 Oath & Ruin Server starting on :8080...
📡 WebSocket: ws://localhost:8080/ws
❤️  Health: http://localhost:8080/health
```

**Проверь:**
- Открой браузер: http://localhost:8080/health
- Должно вернуться: `{"status":"ok","game":"Oath & Ruin"}`

---

## 🔗 Шаг 4: Подключение клиента к серверу

1. **В Godot:**
   - Открой `scripts/autoload/network.gd`
   - Проверь `server_url := "ws://localhost:8080/ws"`

2. **Запусти сервер** (если ещё не запущен):
   ```bash
   cd server
   go run cmd/server/main.go
   ```

3. **Запусти клиент в Godot:**
   - Click Play (F5)
   - Click "Играть"
   - Откроется сцена боя

4. **Проверь подключение:**
   - В консоли сервера должно быть: `👤 Player joined: id_...`
   - В консоли Godot: `✅ Connected to server`

---

## 🎯 Шаг 5: Тестирование

### Тест 1: Сетка и перемещение
1. Запусти сцену `battle.tscn`
2. Кликни на синего персонажа (Character1)
3. Зелёным подсветятся доступные клетки
4. Кликни на клетку для перемещения
5. Персонаж должен переместиться с анимацией

### Тест 2: Сеть
1. Запусти сервер
2. Запусти клиент
3. В консоли сервера проверь сообщения
4. Открой `http://localhost:8080/health` в браузере

---

## 🐛 Возможные проблемы

### Godot не видит скрипты
**Решение:** Перезапусти Godot, проверь пути в `.tscn` файлах

### Сервер не запускается
```bash
# Проверь версию Go
go version  # Должна быть 1.21+

# Переустанови зависимости
cd server
rm -rf go.sum
go mod tidy
```

### WebSocket не подключается
- Проверь, что сервер запущен
- Проверь URL в `network.gd`
- Отключи фаервол на время теста

---

## 📚 Следующие шаги

1. **Изучи код:**
   - `client/scripts/systems/grid.gd` — сетка
   - `client/scripts/components/character.gd` — персонаж
   - `server/cmd/server/main.go` — сервер

2. **Поэкспериментируй:**
   - Измени размер сетки
   - Поменяй скорость перемещения
   - Добавь препятствия

3. **Документируй:**
   - Запиши вопросы
   - Отметь что непонятно

---

## ✅ Чеклист готовности

- [ ] Файлы скопированы в репозиторий
- [ ] Godot проект импортирован
- [ ] Сцены открываются без ошибок
- [ ] Сервер запускается
- [ ] Health check работает
- [ ] Персонаж перемещается по клику
- [ ] Сетка отрисовывается

**Как всё заработает — приступим к первому спринту!** 🎮
