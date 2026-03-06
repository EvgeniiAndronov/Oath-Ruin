# ⚔️ Oath & Ruin

Пошаговая тактическая PvP/PvE игра в стиле Divinity: Original Sin 2 с элементами roguelike.

## 🎮 Описание

**Oath & Ruin** — это многопользовательская тактическая RPG, где две фракции сражаются за контроль над Разломами — источниками магической энергии.

- **PvP 3v3** — командные бои с тактической глубиной
- **PvE** — сюжетная кампания и процедурные подземелья
- **5 классов** — воин, маг, лучник, жрец, убийца
- **Система комбо** — сочетания стихий и состояний

## 🏗️ Архитектура

```
oath-and-ruin/
├── client/          # Godot 4.x клиент
├── server/          # Go сервер
├── docs/            # Документация
└── README.md
```

## 🚀 Быстрый старт

### Клиент (Godot)

1. Установи [Godot 4.x](https://godotengine.org/download)
2. Открой `client/project.godot`
3. Запусти сцену `scenes/main_menu.tscn`

### Сервер (Go)

```bash
cd server
go mod tidy
go run cmd/server/main.go
```

## 📚 Документация

- [Vision](docs/vision.md)
- [Архитектура](docs/architecture/overview.md)
- [API](docs/api/server-api.md)

## 🛠️ Стек

| Компонент | Технология |
|-----------|------------|
| Клиент | Godot 4.x (GDScript) |
| Сервер | Go 1.21+ |
| Протокол | WebSocket + JSON |
| База данных | PostgreSQL |

## 📋 Статус разработки

- [ ] Прототип сетки и перемещения
- [ ] Боевая система (способности)
- [ ] Сетевой код (клиент-сервер)
- [ ] Matchmaking
- [ ] PvE режим

## 👥 Команда

- **PM/Architect/Dev:** Evgenii Andronov
- **Art:** TBA

## 📄 Лицензия

Proprietary (все права защищены)

---

**Присоединяйся к разработке!** 🎮
