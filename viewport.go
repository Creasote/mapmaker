package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Viewport is designed to take the full game board, and only display a sub-image of it.
// It can scroll around the game board.
// It can generate thumbnails to display in a mini-map.

const (
	vp_width  = vp_cells_wide * spriteSize // The number of pixels wide to be displayed on the screen (640px).
	vp_height = vp_cells_high * spriteSize // The number of pixels high to be displayed (480px).
	//width      int // Width of the viewport (see vp_width)
	//height     int // Height of the viewport (see vp_height)
	vp_cells_wide        = 90
	vp_cells_high        = 60
	scroll_button_offset = 15
	vp_x_max             = vp_width + scroll_button_offset
	vp_y_max             = vp_height + scroll_button_offset
)

type Viewport struct {
	vp_x_offset int // The number of CELLS the viewport is shifted horizontally from the origin. A value of 10 means
	// the viewport will start displaying the gameboard from the 10th horizontal cell.
	vp_y_offset int // The number of CELLS the viewport is shifted vertically from the origin. A value of 10 means
	// the viewport will start displaying the gameboard from the 10th vertical cell.

}

// Viewport helper functions

// func (g *Game) portGameBoard(b *board) [][][]int {
// 	return b[:][g.viewport.vp_y_offset : g.viewport.vp_y_offset+g.viewport.cells_high][g.viewport.vp_x_offset : g.viewport.vp_x_offset+g.viewport.cells_wide]
// }

// Function scales the array sizes for RANGE calls to only output the
// func viewportScale() (int, int, int, int) {
// 	x_from := viewport.vp_x_offset
// 	x_to := viewport.vp_x_offset + vp_cells_wide
// 	y_from := viewport.vp_y_offset
// 	y_to := viewport.vp_y_offset + vp_cells_high
// 	return x_from, x_to, y_from, y_to
// }

func draw_ViewportMap(screen *ebiten.Image) {
	for y_ind, y := range game_map[terrain_layer][viewport.vp_y_offset : viewport.vp_y_offset+vp_cells_high] {
		for x_ind, terr := range y[viewport.vp_x_offset : viewport.vp_x_offset+vp_cells_wide] {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(((float64(x_ind) * spriteSize) + scroll_button_offset), ((float64(y_ind) * spriteSize) + scroll_button_offset))
			screen.DrawImage(terrain_map[terr], op)
		}
	}
}

func draw_ViewportEntities(screen *ebiten.Image) {
	for _, ent := range entity_list {
		if tf, modified_x, modified_y := isInViewport(ent.loc.x, ent.loc.y); tf == true {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate((float64(modified_x) * spriteSize), (float64(modified_y) * spriteSize))
			screen.DrawImage(ent.sprite_img, op)
			// If the entity has a path, draw the waypoints.
			if len(ent.path) > 0 {
				for _, waypoint := range ent.path {
					if tf, wp_modfied_x, wp_modified_y := isInViewport(waypoint.x, waypoint.y); tf == true {
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate((float64(wp_modfied_x) * spriteSize), (float64(wp_modified_y) * spriteSize))
						screen.DrawImage(img_path, op)
					}
				}
			}
		}
	}
}

// Recieve the entities x and y coords, and determine whether this falls within the displayable viewport.
// If it does, modify the x and y coords for screen drawing.
func isInViewport(x, y int) (bool, int, int) {
	if x > viewport.vp_x_offset && x < viewport.vp_x_offset+vp_cells_wide {
		if y > viewport.vp_y_offset && y < viewport.vp_y_offset+vp_cells_high {
			// The +1 modifier is a magic number to account for the size of the scroll offset.
			// However, x,y are in cells, whereas the scrollbar offset const is in pixels.
			// 1 cell == 16 pixels.
			return true, x + 1 - viewport.vp_x_offset, y + 1 - viewport.vp_y_offset
		}
	}
	return false, 0, 0
}

// Take the x,y location of the click, and normalises it to the underlying array,
// including scaling, and shifting the viewport.
func viewportClick(x, y int) (int, int) {
	x = ((x - scroll_button_offset) / spriteSize) + viewport.vp_x_offset
	y = ((y - scroll_button_offset) / spriteSize) + viewport.vp_y_offset
	// Ensure value is valid.
	x, y = validateViewportClickOffsets(x, y)

	return x, y
}

// Check to see if the mouse cursor is over the viewport scroll controls.
// Returns a boolean true/false, as well as a modifier pair (x,y)
// Negative values indicate the cursor is over the top (y) or left (x) scroll button.
// Positive valuees indicate the cursor is over the bottom (y) or right (x) scroll button.
// These values can be applied directly to the viewport offset value.
func viewportInScroll(x, y int) (bool, int, int) {
	if x < scroll_button_offset {
		return true, -1, 0
	}
	if y < scroll_button_offset {
		return true, 0, -1
	}
	if x > vp_x_max && x < vp_x_max+scroll_button_offset {
		return true, 1, 0
	}
	if y > vp_y_max && y < vp_y_max+scroll_button_offset {
		return true, 0, 1
	}
	return false, 0, 0
}

// Takes in the viewport position modifiers, validates, and then applies.
// Validation requires that the offset cannot be negative, and
// cannot be greater than that which takes the far side of the viewport past the game map bounds.
func updateViewportScroll(x, y int) {
	viewport.vp_x_offset, viewport.vp_y_offset = validateViewportOffsets(viewport.vp_x_offset+x, viewport.vp_y_offset+y)

}

// Validates the offsets generated for viewport operations.
// Value must be between cell zero and (the right/down most cell visible on the viewport)
func validateViewportOffsets(x, y int) (int, int) {
	x = minInt(x, board_cells_Wide-1-vp_cells_wide)
	x = maxInt(x, 0)

	y = minInt(y, board_cells_High-1-vp_cells_high)
	y = maxInt(y, 0)

	return x, y
}

// Validates the offsets applied to mouse inputs in the viewport.
// Values must be between zero (the first row/column on the board)
// to the maximum board size.
func validateViewportClickOffsets(x, y int) (int, int) {
	x = minInt(x, board_cells_Wide-1)
	x = maxInt(x, 0)

	y = minInt(y, board_cells_High-1)
	y = maxInt(y, 0)

	return x, y
}
