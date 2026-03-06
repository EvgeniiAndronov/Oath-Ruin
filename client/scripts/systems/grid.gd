extends Node2D

class_name BattleGrid

@export var tile_size: int = 64
@export var grid_width: int = 16
@export var grid_height: int = 16

var grid := {}  # {Vector2i: TileData}

func _ready():
	_generate_grid()

func _generate_grid():
	for x in range(grid_width):
		for y in range(grid_height):
			var pos = Vector2i(x, y)
			grid[pos] = {
				"walkable": true,
				"occupied": false,
				"terrain": "grass"
			}

func _draw():
	for x in range(grid_width):
		for y in range(grid_height):
			var world_pos = grid_to_world(Vector2i(x, y))
			var color = Color.WHITE
			if not grid[Vector2i(x, y)]["walkable"]:
				color = Color.RED
			draw_rect(
				Rect2(world_pos - Vector2(tile_size/2, tile_size/2), Vector2(tile_size, tile_size)),
				color,
				false,
				1
			)

func world_to_grid(world_pos: Vector2) -> Vector2i:
	return Vector2i(
		floor(world_pos.x / tile_size),
		floor(world_pos.y / tile_size)
	)

func grid_to_world(grid_pos: Vector2i) -> Vector2:
	return Vector2(
		grid_pos.x * tile_size + tile_size / 2,
		grid_pos.y * tile_size + tile_size / 2
	)

func is_valid(pos: Vector2i) -> bool:
	return pos.x >= 0 and pos.x < grid_width and \
		   pos.y >= 0 and pos.y < grid_height

func is_walkable(pos: Vector2i) -> bool:
	return is_valid(pos) and grid[pos]["walkable"] and not grid[pos]["occupied"]

func set_occupied(pos: Vector2i, occupied: bool):
	if is_valid(pos):
		grid[pos]["occupied"] = occupied
