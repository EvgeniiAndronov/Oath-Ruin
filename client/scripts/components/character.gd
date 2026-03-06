extends CharacterBody2D

class_name Character

@export var movement_range: int = 5
@export var grid: BattleGrid
@export var character_id: String = ""
@export var player_id: String = ""
@export var character_name: String = ""
@export var character_class: String = "warrior"

var hp: int = 100
var max_hp: int = 100
var can_move: bool = true
var can_act: bool = true

func _ready():
	if grid == null:
		grid = get_node_or_null("/root/BattleScene/BattleGrid")

func get_available_tiles() -> Array:
	if grid == null:
		return []
	
	var tiles := []
	var visited := {}
	var queue = []
	queue.append([grid_position(), 0])
	
	while not queue.is_empty():
		var current = queue.pop_front()
		var pos = current[0]
		var cost = current[1]
		
		if visited.has(pos) or cost > movement_range:
			continue
		
		visited[pos] = true
		
		if pos != grid_position() and grid.is_walkable(pos):
			tiles.append(pos)
		
		for dir in [Vector2i.UP, Vector2i.DOWN, Vector2i.LEFT, Vector2i.RIGHT]:
			var next = pos + dir
			if not visited.has(next):
				queue.append([next, cost + 1])
	
	return tiles

func grid_position() -> Vector2i:
	if grid == null:
		return Vector2i(0, 0)
	return grid.world_to_grid(position)

func move_to(new_pos: Vector2i):
	if grid == null:
		return
	
	# Освободить старую клетку
	grid.set_occupied(grid_position(), false)
	
	# Анимация перемещения
	var tween = create_tween()
	tween.tween_property(self, "position", grid.grid_to_world(new_pos), 0.3)
	
	# Занять новую клетку
	await tween.finished
	grid.set_occupied(new_pos, true)
	
	can_move = false

func reset_turn():
	can_move = true
	can_act = true
