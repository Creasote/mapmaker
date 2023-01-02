package main

import (
	"fmt"

	//"math/rand"
	"time"
	//"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	_ "image/png"
	"log"

	//"github.com/hajimehoshi/bitmapfont/v2"
	//"github.com/hajimehoshi/bitmapfont"

	"golang.org/x/image/font"
	//"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/text"
)

type coords struct {
	x int
	y int
}

type node struct {
	loc      coords
	prev     *node
	terrain  int
	cost     float32 // cost to get to this point, including it's own terrain cost
	distance float32 // estimated distance to goal
	estimate float32 // cost + estimated distance to goal
}

type entity struct {
	name               string
	loc                coords
	mob_type           int
	sprite_img         *ebiten.Image
	movement_speed     float32
	health             float32
	armour             float32
	damage_per_attack  float32
	attacks_per_second float32
	attack_success_pc  float32
	attack_range       float32
	last_attack_time   int
	target             *entity
	path               []coords
}

const (
	terrain_layer = 0
	entity_layer  = 1
)

type board [2][board_cells_High][board_cells_Wide]int

const (
	road                 = 1
	grassland            = 2
	sand                 = 3
	forest               = 5
	water                = 9
	impassable_threshold = 90
	cliff                = 98
	wall                 = 99
)

var terrain_list = []int{road, grassland, sand, forest, water, cliff, wall}
var terrain_map = map[int]*ebiten.Image{}

const (
	//screenWidth  = 100 // Maximum # of cells, for a total of 1600px wide
	//screenHeight = 60  // Maximum # of cells, for a total of 960px high
	board_cells_High = 120
	board_cells_Wide = 100 // Number of cells wide, for a total of 20x16 = 320px
	spriteSize       = 16
	//xViewport        = board_cells_Wide * spriteSize
	//yViewport        = board_cells_High * spriteSize
	xScreen = (vp_cells_wide * spriteSize) + (2 * scroll_button_offset) + xMainMenuSize //1808 //xSize * spriteSize // used by Ebiten to set canvas size
	yScreen = (vp_cells_high * spriteSize) + (2 * scroll_button_offset)                 //960  //ySize * spriteSize
)

var game_map board

var START = coords{0, 0}
var GOAL = coords{18, 6}

var entity_list []*entity

var console Console
var mplusNormalFont font.Face

var viewport = Viewport{
	vp_x_offset: 0,
	vp_y_offset: 0,
}

type Game struct {
	keylist    []ebiten.Key
	flood_mode bool // true places tiles in flood mode,
	// false places a single tile at the cursor position
	object_value int /* records what type of tile or object is placed on mouse click:
	"0: pathfinder"
	"1: road terrain",
	"2: grassland terrain",
	"3: sand terrain",
	"5: forest",
	"9: water",
	"98: cliff (impassable)",
	"99: wall (impassable)",
	*/
	scroll_state       bool // Goes True when the cursor is detected in the scroll zone. Flips back to False when it shifts outside.
	scroll_state_since time.Time
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return xScreen, yScreen
}

func (g *Game) Update() error {
	// Take mouse inputs
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.flood_mode = false
		g.parse_mouseclick(ebiten.CursorPosition())
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.flood_mode = true
		g.parse_mouseclick(ebiten.CursorPosition())
	}

	if tf, x_vp_modifier, y_vp_modifier := viewportInScroll(ebiten.CursorPosition()); tf == true {
		updateViewportScroll(x_vp_modifier, y_vp_modifier)
		//viewport.vp_x_offset = +x_vp_modifier
		//viewport.vp_y_offset += y_vp_modifier
	}

	// TODO: should keyboard inputs be processed first? It's possible someone could change the input type and click in the same tick.
	// Take keyboard inputs
	g.keylist = inpututil.AppendPressedKeys(g.keylist[:0])
	g.parse_keyboard()
	//parse_keyboard(&g.keylist)

	// Update pathing for entities
	for _, e := range entity_list {
		if len(e.path) == 0 {
			e.pathfind(&game_map)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw map terrain
	//for y_ind, y := range game_map[terrain_layer]  {
	//for y_ind, y := range game_map[terrain_layer][viewport.vp_y_offset : viewport.vp_y_offset+vp_cells_high] {
	// for y_ind, y := range game_map[terrain_layer][viewport.vp_y_offset : viewport.vp_y_offset+vp_cells_high] {
	// 	for x_ind, terr := range y[viewport.vp_x_offset : viewport.vp_x_offset+vp_cells_wide] {
	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate((float64(x_ind) * spriteSize), (float64(y_ind) * spriteSize))
	// 		screen.DrawImage(terrain_map[terr], op)
	// 	}
	// }
	draw_ViewportMap(screen)

	// TODO: Pass entity drawing via a viewport modifier
	// Draw each of the entities
	draw_ViewportEntities(screen)
	// for _, ent := range entity_list {
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate((float64(ent.loc.x) * spriteSize), (float64(ent.loc.y) * spriteSize))
	// 	screen.DrawImage(ent.sprite_img, op)
	// 	// If the entity has a path, draw the waypoints.
	// 	if len(ent.path) > 0 {
	// 		for _, waypoint := range ent.path {
	// 			op := &ebiten.DrawImageOptions{}
	// 			op.GeoM.Translate((float64(waypoint.x) * spriteSize), (float64(waypoint.y) * spriteSize))
	// 			screen.DrawImage(img_path, op)
	// 		}
	// 	}
	// }

	// Draw Menu
	draw_Menu(screen)
	//screen.draw_Menu()

	// Print processing rate for performance monitoring
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())) //ebitenutil.DebugPrint(screen, "This is NOT a test.")

}

var img_road, img_grassland, img_sand, img_forest, img_water, img_wall, img_cliff *ebiten.Image
var img_player, img_path, img_goal *ebiten.Image
var img_menu_border_top, img_menu_border_left, img_menu_border_right, img_menu_border_bottom *ebiten.Image
var img_menu_tl, img_menu_tr, img_menu_bl, img_menu_br *ebiten.Image
var img_logo, img_button, img_button_pressed *ebiten.Image

func init() {
	console.console_add("Beginning init...")
	var err error
	img_road, _, err = ebitenutil.NewImageFromFile("./assets/road.png")
	terrain_map[1] = img_road
	img_grassland, _, err = ebitenutil.NewImageFromFile("./assets/grass.png")
	terrain_map[2] = img_grassland
	img_sand, _, err = ebitenutil.NewImageFromFile("./assets/sand.png")
	terrain_map[3] = img_sand
	img_forest, _, err = ebitenutil.NewImageFromFile("./assets/forest.png")
	terrain_map[5] = img_forest
	img_water, _, err = ebitenutil.NewImageFromFile("./assets/water.png")
	terrain_map[9] = img_water
	img_wall, _, err = ebitenutil.NewImageFromFile("./assets/wall.png")
	terrain_map[99] = img_wall
	img_cliff, _, err = ebitenutil.NewImageFromFile("./assets/cliff.png")
	terrain_map[98] = img_cliff
	img_player, _, err = ebitenutil.NewImageFromFile("./assets/bee.png")
	img_goal, _, err = ebitenutil.NewImageFromFile("./assets/goal.png")
	img_path, _, err = ebitenutil.NewImageFromFile("./assets/path.png")
	img_menu_border_top, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_top.png")
	img_menu_border_left, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_left.png")
	img_menu_border_right, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_right.png")
	img_menu_border_bottom, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_bottom.png")
	img_menu_tl, _, err = ebitenutil.NewImageFromFile("./assets/menu/tl.png")
	img_menu_tr, _, err = ebitenutil.NewImageFromFile("./assets/menu/tr.png")
	img_menu_bl, _, err = ebitenutil.NewImageFromFile("./assets/menu/bl.png")
	img_menu_br, _, err = ebitenutil.NewImageFromFile("./assets/menu/br.png")
	img_logo, _, err = ebitenutil.NewImageFromFile("./assets/menu/logo.png")
	img_button, _, err = ebitenutil.NewImageFromFile("./assets/menu/save.png")
	img_button_pressed, _, err = ebitenutil.NewImageFromFile("./assets/menu/save_depressed.png")

	// REMOVE THIS LATER
	terrain_map[0] = img_menu_border_top
	// PLACEHOLDER TO REMOVE ABOVE.

	if err != nil {
		log.Fatal(err)
	}
	console.console_add("Images successfully loaded.")

	//var terrain_list = []int{road, grassland, sand, forest, water, cliff, wall}
	// var terrain_icons = []*ebiten.Image{
	// 	img_road,
	// 	img_grassland,
	// 	img_sand,
	// 	img_forest,
	// 	img_water,
	// 	img_cliff,
	// 	img_wall,
	// }

	// create sprite

	goal_entity := entity{
		name:               "Goal",
		loc:                GOAL,
		mob_type:           0,
		sprite_img:         img_goal,
		movement_speed:     0,
		health:             0,
		armour:             0,
		damage_per_attack:  0,
		attacks_per_second: 0,
		attack_success_pc:  0,
		attack_range:       0,
		last_attack_time:   0,
		target:             nil,
		path:               []coords{},
	}
	entity_list = append(entity_list, &goal_entity)

	console.console_add("Initialising menu...")
	init_Menu()
	console.console_add("Init complete.")
}

func (e *entity) move_entity() {
	for {
		if len(e.path) > 0 {
			// Set the location to the next waypoint, and remove that waypoint from the path list.
			e.loc, e.path = e.path[len(e.path)-1], e.path[:len(e.path)-1]
		}
		//time.Sleep(500 * time.Millisecond)
		//sleep_time := 10000 / e.movement_speed
		time.Sleep(time.Duration(10000/e.movement_speed) * time.Millisecond)
	}
}

func main() {
	// Create a background go-routine that polls for user movement. Maybe set variable sleep time based on movement speed?
	for _, char := range entity_list {
		go char.move_entity()
	}

	ebiten.SetWindowSize(xScreen, yScreen)
	ebiten.SetWindowTitle("Ebiten Test")

	console.console_add("Entering game loop.")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}

}
