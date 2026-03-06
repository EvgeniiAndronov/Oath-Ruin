extends Node2D

class_name CombatSystem

@export var grid: BattleGrid
@export var selected_character: Character

var available_tiles := []
var highlight_color := Color(0, 1, 0, 0.5)

func _ready():
	if grid == null:
		grid = get_node_or_null("/root/BattleScene/BattleGrid")

func _input(event):
	if event is InputEventMouseButton and event.pressed:
		var world_pos = get_global_mouse_position()
		var grid_pos = grid.world_to_grid(world_pos) if grid else Vector2i(0, 0)
		
		if event.button_index == MOUSE_BUTTON_LEFT:
			_on_left_click(grid_pos)
		elif event.button_index == MOUSE_BUTTON_RIGHT:
			_on_right_click()

func _on_left_click(grid_pos: Vector2i):
	if selected_character and not selected_character.can_move:
		return
	
	if selected_character and grid_pos in available_tiles:
		selected_character.move_to(grid_pos)
		available_tiles.clear()
		_clear_highlights()
	else:
		select_character_at(grid_pos)

func _on_right_click():
	available_tiles.clear()
	selected_character = null
	_clear_highlights()

func select_character_at(grid_pos: Vector2i):
	var characters = get_tree().get_nodes_in_group("characters")
	for char in characters:
		if char.grid_position() == grid_pos:
			selected_character = char
			available_tiles = char.get_available_tiles()
			_highlight_tiles(available_tiles)
			return

func _highlight_tiles(tiles: Array[Vector2i]):
	_clear_highlights()
	for tile in tiles:
		_draw_highlight(tile, highlight_color)

func _draw_highlight(grid_pos: Vector2i, color: Color):
	if grid == null:
		return
	var world_pos = grid.grid_to_world(grid_pos)
	var rect = Rect2(world_pos - Vector2(grid.tile_size/2, grid.tile_size/2), Vector2(grid.tile_size, grid.tile_size))
	draw_rect(rect, color)

func _clear_highlights():
	queue_redraw()
