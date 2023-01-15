package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

type Button struct {
	enabled       bool // Is this button available for use?
	active        bool // Is the button the currently selected button?
	img           *ebiten.Image
	height, width int
	actionID      int // Used to switch-case the appropriate action from the menu click parser
	action        Action
}

// type Menu struct{
// 	button *Button
// }

type Action interface {
	getCost(int) int
	do()
}

type menuZoneDescriptor struct {
	x, y          float64 // Absolute x and y coords for the START of the menu zone.
	height, width int
}
type menuStructureDescriptor struct {
	loc      string
	actionID int
	action   Action
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

var buttonList []*Button
var img_btnPlaceSpawn, img_btnSpawnArmour, img_btnSpawnAttackSpeed, img_btnSpawnDmg, img_btnSpawnHP, img_btnSpawnToHit *ebiten.Image
var img_btnPlaceSpawner, img_btnSpawnerOutput, img_btnSpawnerRate *ebiten.Image
var img_btnFade *ebiten.Image

var menuZoneFile = []menuZoneDescriptor{
	{float64(xMainMenuOffset + spriteSize), float64(spriteSize), 128, 128}, // Logo
	{float64(xMainMenuOffset + spriteSize), float64(200), 32 * 8, 80},      // Buttons
	{float64(xMainMenuOffset + spriteSize), float64(540), 14 * 30, 80},     // Console
	{float64(xMainMenuOffset + spriteSize), float64(700), 32, 80},          // Minimap
}

var menuStructureFile = []menuStructureDescriptor{
	{"./assets/menu/place_spawner_button.png", 10, spawnerRate},
	//{"./assets/menu/place_spawn_button.png", 0},
	{"./assets/menu/spawn_armour_upgrade_button.png", 1, spawnArmour},
	{"./assets/menu/spawn_attack_speed_upgrade_button.png", 2, spawnAttackRate},
	{"./assets/menu/spawn_dmg_upgrade_button.png", 3, spawnDmg},
	{"./assets/menu/spawn_hp_upgrade_button.png", 4, spawnHealth},
	{"./assets/menu/spawn_tohit_upgrade_button.png", 5, spawnToHit},
	{"./assets/menu/spawner_output_upgrade_button.png", 11, spawnerMaxOutput},
	{"./assets/menu/spawner_rate_upgrade_button.png", 12, spawnerRate},
}

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

	// Call the fns to build an interable list of buttons, which are then drawn in the update method
	// Requires a pre-defined menu structure descriptor
	buttonList = buildButtonList(menuStructureFile)

	img_btnFade = ebiten.NewImage(80, 32)
	img_btnFade.Fill(&color.RGBA{10, 10, 10, 125})

}

func buildButtonList(msf []menuStructureDescriptor) []*Button {
	var bl []*Button

	for _, b := range msf {
		if err, btn := addBtn(b.loc, b.actionID, b.action); err == nil {
			bl = append(bl, btn)
		}
	}
	return bl
}

func addBtn(l string, actionid int, action Action) (error, *Button) {
	i, _, err := ebitenutil.NewImageFromFile(l)

	b := &Button{
		enabled:  false,
		active:   false,
		img:      i,
		height:   32,
		width:    80,
		actionID: actionid,
		action:   action,
	}
	return err, b
}

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
	// for i, inst := range instructions {
	// 	text.Draw(screen, inst+string('\n'), mplusNormalFont, xInstructions, yInstructions+(i*console.print_height), color.White)
	// }
	for i, btn := range buttonList {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(menuZoneFile[1].x, (menuZoneFile[1].y + float64(i*btn.height)))
		screen.DrawImage(btn.img, op)
		if !btn.enabled {
			screen.DrawImage(img_btnFade, op)
		}
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

//func (a Action)
