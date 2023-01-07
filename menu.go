package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Console struct {
	log          [30]string
	max_rows     int
	font_size    int
	print_height int
}

const (
	xMainMenuOffset = vp_width + (2 * scroll_button_offset)
	xButtonOffset   = xMainMenuOffset + spriteSize
	yButtonOffset   = 200
	yButtonSize     = 32
	yMainMenuOffset = 0 //spriteSize * cellsHigh
	xMainMenuSize   = 15 * spriteSize
	yMainMenuSize   = vp_height + (2 * scroll_button_offset)
	xConsole        = xMainMenuOffset + spriteSize
	yConsole        = yScreen - 450
	xInstructions   = xMainMenuOffset + spriteSize
	yInstructions   = 200
)

// var menu_layout = [][]*ebiten.Image{
// 	{img_button},
// 	{img_button},
// 	{img_button},
// 	{img_button, img_button},
// }

var instructions = []string{
	"R: randomise map",
	"O: optimise map",
	"P: place pathfinder",
	"U: place spawner",
	"X: place target",
	"C: Get current coords",
	" ",
	"TERRAIN:",
	"1: road terrain",
	"2: grassland terrain",
	"3: sand terrain",
	"4: forest",
	"5: water",
	"8: cliff (impassable)",
	"9: wall (impassable)",
	"Left+Click single cell",
	"Right+Click flood fill",
	" ",
	"Other:",
	"S: save map file",
	"L: load previous save",
}

func init_Menu() {

	console.max_rows = 10
	console.font_size = 14
	console.print_height = 14

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(console.font_size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Draw button helper function.
// Takes the x and y position of the buttons,
// within the menu (not screen pixels), 1 INDEX (not zero indexed!)
// ie. the first button on the left would be 1,1 (x,y).
// The first button on the second row would be (1,2).
// The second button on the third row would be (2,3).
// func draw_button(screen, button_img *ebiten.Image, x, y int) {
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(float64(xButtonOffset+(x-1)*(button_img.Bounds().Size().X)), float64(yButtonOffset+(y-1)*(button_img.Bounds().Size().Y)))
// 	//op.GeoM.Translate(float64(xButtonOffset), float64(yButtonOffset+(y-1*button_img.Bounds().Size().Y)))
// 	//fmt.Println("Button width: ", (x-1)*(button_img.Bounds().Size().X))
// 	screen.DrawImage(button_img, op)
// }

// func draw_buttons(screen *ebiten.Image, menu_array [][]*ebiten.Image) {
// 	for row_ind, row_buttons := range menu_array {
// 		fmt.Print("Row: ", row_ind)
// 		for col_ind, butts := range row_buttons {
// 			// if the column is > 0, shift this button right by the size of the previous button + some margin
// 			fmt.Println(" Col: ", col_ind)
// 			op := &ebiten.DrawImageOptions{}
// 			//op.GeoM.Translate(float64(xButtonOffset+(col_ind)*(butts.Bounds().Size().X)), float64(yButtonOffset+(row_ind)*(butts.Bounds().Size().Y)))
// 			op.GeoM.Translate(float64(xButtonOffset), float64(yButtonOffset))
// 			screen.DrawImage(butts, op)
// 		}
// 	}
// 	//op.GeoM.Translate(float64(xButtonOffset), float64(yButtonOffset+(y-1*button_img.Bounds().Size().Y)))
// 	//fmt.Println("Button width: ", (x-1)*(button_img.Bounds().Size().X))
// }

//func (screen *ebiten.Image) draw_Menu() {
func (g *Game) draw_Menu(screen *ebiten.Image) {
	// Draw border
	for x := xMainMenuOffset + spriteSize; x < (xScreen - spriteSize); x += spriteSize {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(0))
		screen.DrawImage(img_menu_border_top, op)
	}
	for x := xMainMenuOffset + spriteSize; x < (xScreen - spriteSize); x += spriteSize {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(yScreen-spriteSize))
		screen.DrawImage(img_menu_border_bottom, op)
	}
	for y := spriteSize; y < (yScreen - spriteSize); y += spriteSize {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xMainMenuOffset), float64(y))
		screen.DrawImage(img_menu_border_left, op)
	}
	for y := spriteSize; y < (yScreen - spriteSize); y += spriteSize {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(xScreen-spriteSize), float64(y))
		screen.DrawImage(img_menu_border_right, op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xMainMenuOffset), float64(0))
	screen.DrawImage(img_menu_tl, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xScreen-spriteSize), float64(0))
	screen.DrawImage(img_menu_tr, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xScreen-spriteSize), float64(yScreen-spriteSize))
	screen.DrawImage(img_menu_br, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xMainMenuOffset), float64(yScreen-spriteSize))
	screen.DrawImage(img_menu_bl, op)

	// Draw logo
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xMainMenuOffset+spriteSize), float64(spriteSize))
	screen.DrawImage(img_logo, op)

	// Draw menu
	// for i := 0; i < 5; i++ {
	// 	// op = &ebiten.DrawImageOptions{}
	// 	// op.GeoM.Translate(float64(xButtonOffset), float64(yButtonOffset+(i*yButtonSize)))
	// 	// screen.DrawImage(img_button, op)
	// 	draw_button(screen, img_button, 1, i)
	// }

	//draw_buttons(screen, menu_layout)

	// Draw instuctions:
	for i, inst := range instructions {
		text.Draw(screen, inst+string('\n'), mplusNormalFont, xInstructions, yInstructions+(i*console.print_height), color.White)
	}

	// Draw console
	for i, console_text := range console.log {
		text.Draw(screen, console_text+string('\n'), mplusNormalFont, xConsole, yConsole+(i*console.print_height), color.White)
	}

	// Draw minimap
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(mm_scaling), float64(mm_scaling))
	// TODO: set a better position. Make it so the position is fixed to the RHS, and shifts based on map size.
	op.GeoM.Translate(float64(xMainMenuOffset+spriteSize), float64(yScreen-(2*board_cells_High+50)))
	screen.DrawImage(img_minimap, op)
}

func (c *Console) console_add(s string) {
	for i := c.max_rows - 1; i > 0; i-- {
		c.log[i] = c.log[i-1]
	}

	c.log[0] = s
}

// load the assets

// add to slice

// make slice available to draw method

// take mouse inputs
