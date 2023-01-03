package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	mm_scaling = 2
)

/*
"0: pathfinder"
"1: road terrain",
"2: grassland terrain",
"3: sand terrain",
"5: forest",
"9: water",
"98: cliff (impassable)",
"99: wall (impassable)",
*/

var mm_terrain map[int]color.RGBA
var mm_height, mm_width int
var img_vp_boundingbox *ebiten.Image

func initMinimap() {
	// Create the minimap variables
	mm_height = (len(game_map[terrain_layer]))
	mm_width = (len(game_map[terrain_layer][:]))
	img_minimap = ebiten.NewImage(mm_width, mm_height)
	// Including creating an image for the bounding box
	img_vp_boundingbox = ebiten.NewImage(vp_cells_wide, vp_cells_high)
	img_vp_boundingbox = create_vp_boundingbox(img_vp_boundingbox, vp_cells_wide, vp_cells_high)

	mm_terrain = make(map[int]color.RGBA)
	mm_terrain[0] = color.RGBA{20, 20, 20, 255}     // Gray, I think?
	mm_terrain[1] = color.RGBA{0, 0, 0, 255}        // Black, road
	mm_terrain[2] = color.RGBA{50, 250, 50, 255}    // Green, grassland
	mm_terrain[3] = color.RGBA{250, 250, 100, 255}  // Yellow? sand.
	mm_terrain[5] = color.RGBA{0, 180, 0, 255}      // Darker green? Forest.
	mm_terrain[9] = color.RGBA{0, 128, 250, 255}    // Blue. Water.
	mm_terrain[98] = color.RGBA{150, 75, 0, 255}    // Brownish? Cliff.
	mm_terrain[99] = color.RGBA{150, 150, 150, 255} // Gray. Wall.

	// Once initialised, spawn a go-routine that runs infinitely, updating the minimap periodically.
	go updateMinimap(&game_map)
}

func drawMinimap(b *board) {
	pixels := make([]uint8, 4*mm_width*mm_height)
	for y_ind, x_array := range b[terrain_layer] {
		for x_ind, x_val := range x_array {
			p := 4 * (x_ind + y_ind*mm_width)
			pixels[p] = mm_terrain[x_val].R
			pixels[p+1] = mm_terrain[x_val].G
			pixels[p+2] = mm_terrain[x_val].B
			pixels[p+3] = mm_terrain[x_val].A
		}
	}
	img_minimap.WritePixels(pixels)

}

func updateMinimap(b *board) {
	minimap_last_updated := time.Now().UnixMilli()
	drawMinimap(b)

	for {
		if minimap_last_updated < time.Now().UnixMilli()+10000 {
			drawMinimap(b)
			minimap_last_updated = time.Now().UnixMilli()
		}
		minimapDrawViewport()
		time.Sleep(100 * time.Millisecond)
	}
}

func minimapDrawViewport() {
	//ebitenutil.DrawRect(img_minimap, float64(viewport.vp_x_offset), float64(viewport.vp_y_offset), vp_cells_wide, vp_cells_high, color.White)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(viewport.vp_x_offset), float64(viewport.vp_y_offset))
	img_minimap.DrawImage(img_vp_boundingbox, op)
}

func create_vp_boundingbox(m *ebiten.Image, w, h int) *ebiten.Image {
	// Left border
	ebitenutil.DrawLine(m, float64(1), float64(0), float64(1), float64(h), color.White)
	// Right border
	ebitenutil.DrawLine(m, float64(w), float64(0), float64(w), float64(h), color.White)
	// Top border
	ebitenutil.DrawLine(m, float64(0), float64(0), float64(w), float64(0), color.White)
	// Bottom border
	ebitenutil.DrawLine(m, float64(0), float64(h-1), float64(w), float64(h-1), color.White)
	return m
}
