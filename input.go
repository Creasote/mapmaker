package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Work out what zone a mouse click occurs, and pass on to relevant handler.
func (g *Game) parse_mouseclick(x, y int) {
	if x < scroll_button_offset {
		// TODO: do we need this click action to move viewport? IF so, add to move right/down, too.
		viewport.vp_x_offset = maxInt(0, viewport.vp_x_offset-1)
	} else if y < scroll_button_offset {
		viewport.vp_y_offset = maxInt(0, viewport.vp_y_offset-1)
	} else if x < vp_x_max && y < vp_y_max {
		// The click was in the game map. Normalise the coordinates to match the
		// underlying grid structure, including removing the scroll button spacing.
		x, y = viewportClick(x, y)

		// Object Value (-1) is the goal (X)
		if g.object_value == -1 {
			set_goal(x, y)
			// Object Value (0) is a pathfinder entity.
		} else if g.object_value == 0 {
			place_entity(x, y)
			// All others are terrain types.
		} else {
			// Set game map array to value.
			g.place_terrain(&game_map, x, y)
		}

	}
	// else mouseclick was in the menu.
	// TODO: put menu mouseclick handling here.
}

// Take keyboard input, and pass to relevant handler.
func (g *Game) parse_keyboard() {
	for i, k := range g.keylist {
		switch k {
		case ebiten.KeyS:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				fmt.Println("Saving...")
				console.console_add("Saving map file...")
				err := save_map(game_map)
				if err != nil {
					fmt.Println("Error saving, map not exported")
					console.console_add("Error saving, map not exported!")
				}
			}
		case ebiten.KeyL:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				console.console_add("Loading map file from disk...")
				restored_map, err := load_map()
				if err != nil {
					console.console_add("Did not successfully load map file from disk.")
				}

				game_map = restored_map
				for y_ind, x_axis := range game_map[entity_layer] {
					for x_ind, cell_val := range x_axis {
						// Goal location is demarked as a 1 on the entity layer.
						// TODO: remove magic number.
						if cell_val == 1 {
							entity_list[0].loc = coords{x_ind, y_ind}
						}
					}
				}
				drawMinimap(&game_map)
			}

			// TODO: Add Diamond Square map generation

			// Create random noise map
		case ebiten.KeyR:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				game_map = randomise_board(&game_map)
			}

			// Optimise map
		case ebiten.KeyO:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				optimise_board(&game_map)
			}

			// Display current coords on mouse click
		case ebiten.KeyC:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				x, y := viewportClick(ebiten.CursorPosition())
				str := fmt.Sprint(x) + "." + fmt.Sprint(y)
				console.console_add("Grid coords: " + str)
			}

			// Pathfinder
		case ebiten.KeyP:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 0
				console.console_add("Placing PATHFINDER units")
			}

			// Set Goal location
		case ebiten.KeyX:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = -1
				console.console_add("Updating Goal (X)")
			}

			// Road
		case ebiten.Key1:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 1
				console.console_add("Placing ROAD terrain")
			}

			// Grassland
		case ebiten.Key2:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 2
				console.console_add("Placing GRASS terrain")
			}
			// Sand
		case ebiten.Key3:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 3
				console.console_add("Placing SAND terrain")
			}

			// Forest
		case ebiten.Key4:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 5
				console.console_add("Placing FOREST terrain")
			}

			// Water
		case ebiten.Key5:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 9
				console.console_add("Placing WATER terrain")
			}

			// Cliff
		case ebiten.Key8:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 98
				console.console_add("Placing CLIFF obstacle")
			}

			// Wall
		case ebiten.Key9:
			// Required debounce
			if inpututil.IsKeyJustPressed(k) {
				g.object_value = 99
				console.console_add("Placing WALL obstacle")
			}
		}

		// Remove the key we just tested against. Avoids reprosessing
		// and also clears out un-allocted keys.
		if len(g.keylist) > 1 {
			g.keylist = append((g.keylist)[:i], (g.keylist)[i+1:]...)
		} else {
			g.keylist = (g.keylist)[:0]
		}
	}
}
