package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	_ "image/png"
	"log"

	"golang.org/x/image/font"
)

type coords struct {
	x int
	y int
}

type node struct {
	loc      coords
	prev     *node
	terrain  int
	cost     float64 // cost to get to this point, including it's own terrain cost
	distance float64 // estimated distance to goal
	estimate float64 // cost + estimated distance to goal
}

// type spawn struct {
// 	name               string
// 	loc                coords
// 	mob_type           int
// 	alive              bool
// 	action             int  // Index of action currently being undertaken
// 	inCombat           bool // flag True when initiating combat. Reset to false when target dies.
// 	sprite_img         *ebiten.Image
// 	movement_speed     float64 // Should be between 0 (no movement) and ~ 50. Higher values may be OP.
// 	health             float64
// 	armour             float64 // Should be between 0 (no reduction) and 1000 (full reduction in damage taken)
// 	damage_per_attack  float64
// 	attacks_per_second float64
// 	attack_success_pc  float64 // Should be [0,1). Crit chance is calculated as the difference between 1 and attack_success_pc.
// 	attack_range       float64
// 	last_attack_time   int
// 	target             []*spawn
// 	path               []coords
// }

// type spawner struct {
// 	name     string
// 	loc      coords
// 	mob_type int
// 	alive    bool
// 	action   int // Index of action currently being undertaken
// 	//inCombat           bool // flag True when initiating combat. Reset to false when target dies.
// 	sprite_img *ebiten.Image
// 	//movement_speed     float64 // Should be between 0 (no movement) and ~ 50. Higher values may be OP.
// 	health float64
// 	armour float64 // Should be between 0 (no reduction) and 1000 (full reduction in damage taken)
// 	//damage_per_attack  float64
// 	//attacks_per_second float64
// 	//attack_success_pc  float64 // Should be [0,1). Crit chance is calculated as the difference between 1 and attack_success_pc.
// 	//attack_range       float64
// 	//last_attack_time   int
// 	target []*spawn
// 	path   []coords
// }

// type CombatEntity interface {
// 	// Choose a target
// 	findEnemy() *entity
// 	// Commence the combat sequence
// 	initiateCombat()
// 	receiveCombat()
// 	// Perform combat
// 	attackEnemy()
// 	takeDamage()
// }

type ThoughtfulEntity interface {
	brain()
	getLoc() (int, int)
	getSprite() *ebiten.Image
	//loc coords
	//sprite_img *ebiten.Image
}

// Game parameters
const (
	board_cells_High = 120 // Number of cells high.
	board_cells_Wide = 100 // Number of cells wide. Pixel count == cells * spriteSize
	spriteSize       = 16

	// used by Ebiten to set canvas size
	xScreen = (vp_cells_wide * spriteSize) + (2 * scroll_button_offset) + xMainMenuSize //1808 //xSize * spriteSize //
	yScreen = (vp_cells_high * spriteSize) + (2 * scroll_button_offset)                 //960  //ySize * spriteSize
)

// Game map board structure
// 3 dimensional array
// x,y and layers for terrain, entities
type board [2][board_cells_High][board_cells_Wide]int

// Terrain keys
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
	terrain_layer = 0
	entity_layer  = 1
)

// Entity actions
// TODO: This needs reworking to:
// Combat (offensive/defensive, doesn't matter),
// Move (to whatever target, including fleeing),
// Rest (do nothing, sleep),
// Action (eat, repair, gather resources - may need way to determin this is what they're doing - maybe separate numbers? Maybe tiers?)
// How to chain actions (ie. move to resource, then gather resource)? Maybe what is below works, because the target is defined, the move is implied.
const (
	// TODO: set a default value (0) - what do we want these to do by default? Attack enemy base?

	// 1. Self defence
	actionDefence = 1
	// 2. Feed
	actionFeed = 2
	// 3. Rest and recover
	actionRest = 3
	// 4. Local goal (attack enemy agents)
	actionAttackEnemy = 4
	// 5. Base repair
	actionRepair = 5
	// 6. Resource gathering
	actionGatherResources = 6
	// 7. Global goal (attack enemy base)
	actionAttackEnemyBase = 7
	// 8. Move
)

// Global variables
var game_map board
var img_minimap, img_scoreboard *ebiten.Image

var START = coords{0, 0}
var GOAL = coords{18, 6}

var spawn_list []*spawn
var spawner_list []*spawner
var target_list []*target
var entity_list []*ThoughtfulEntity

var console Console
var mplusNormalFont font.Face

var viewport = Viewport{
	vp_x_offset: 0,
	vp_y_offset: 0,
}

var score int                                        // Total score
var sps = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}     // Records the scores per second. Store live score in sps[10], avg calculated over [0:9]
var dps = []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // Records entries for the last 10 seconds. Store live score in dps[10], avg calculated over [0:9]
var spawnCount int                                   // records total spawns

type Game struct {
	keylist    []ebiten.Key
	flood_mode bool // true places tiles in flood mode,
	// false places a single tile at the cursor position
	paint bool // when true, it indicates cell-by-cell painting is active.
	// When active, any cells the mouse passes over will be "painted" with the object
	// per game.object_value (only applies to terrain, not entities).
	// Activated by holding down mouse button in single paint mode. Deactivated
	// simply by releasing mouse button.
	object_value int /* records what type of tile or object is placed on mouse click:
	-2: spawner
	-1: goal
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
	tick               int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return xScreen, yScreen
}

func (g *Game) Update() error {
	// Take mouse inputs
	// Get the cursor position once to re-use during update:
	cursor_x, cursor_y := ebiten.CursorPosition()

	// Left-Click will place either an entity or a terrain.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.flood_mode = false
		g.parse_mouseclick(cursor_x, cursor_y)
	}

	// Holding the left button for more than 30 ticks (~0.5 seconds)
	// will put it in to paint mode. Drag around to pain terrain (only).
	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 30 {
		if g.object_value > 0 {
			g.paint = true
		}
	}

	// To exit terrain paint mode, simply release the mouse button.
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.paint = false
	}

	// Right-click will flood fill terrain.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.flood_mode = true
		g.parse_mouseclick(cursor_x, cursor_y)
	}

	// If paint mode is enabled, apply a cell-by-cell paint.
	if g.paint {
		cursor_x, cursor_y = viewportClick(cursor_x, cursor_y)
		g.place_terrain(&game_map, cursor_x, cursor_y)
	}

	if tf, x_vp_modifier, y_vp_modifier := viewportInScroll(cursor_x, cursor_y); tf == true {
		updateViewportScroll(x_vp_modifier, y_vp_modifier)
	}

	// TODO: should keyboard inputs be processed first? It's possible someone could change the input type and click in the same tick.
	// Take keyboard inputs
	g.keylist = inpututil.AppendPressedKeys(g.keylist[:0])
	g.parse_keyboard()
	//parse_keyboard(&g.keylist)

	// Update pathing for entities
	// TODO: Use same iteration to remove dead enemies.
	// for _, e := range entity_list {
	// 	if len(e.path) == 0 {
	// 		// TODO: investigate why pathing glitches when using go-routines.
	// 		e.pathfind(&game_map)
	// 	}
	// }

	// Tick related updates
	// g.tick++
	// if g.tick > 60 {
	// 	g.tick = 0
	// 	createMinimap(&game_map)
	// }

	for ind, ent := range spawn_list {
		// ent.brain()
		if !ent.alive {
			if ind < len(spawn_list)-1 {
				// entity is not the last in the array
				arrayEnd := spawn_list[ind+1:]
				spawn_list = spawn_list[:ind]
				spawn_list = append(spawn_list, arrayEnd...)

			} else {
				// entity is the last one in the array
				spawn_list = spawn_list[:ind]
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw map terrain
	draw_ViewportMap(screen)

	// Draw each of the entities
	drawViewportSpawners(screen, spawner_list)
	drawViewportTarget(screen, target_list)
	drawViewportSpawns(screen, spawn_list)

	// Draw Menu
	g.draw_Menu(screen)
	//screen.draw_Menu()

	// Draw Scoreboard
	drawScoreboard(screen)

	// Print processing rate for performance monitoring
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())) //ebitenutil.DebugPrint(screen, "This is NOT a test.")

}

var img_road, img_grassland, img_sand, img_forest, img_water, img_wall, img_cliff *ebiten.Image
var img_player, img_path, img_goal, img_spawner *ebiten.Image
var img_menu_border_top, img_menu_border_left, img_menu_border_right, img_menu_border_bottom *ebiten.Image
var img_menu_tl, img_menu_tr, img_menu_bl, img_menu_br *ebiten.Image
var img_logo, img_button, img_button_pressed *ebiten.Image

func init() {
	console.console_add("Beginning init...")
	var err error
	img_road, _, err = ebitenutil.NewImageFromFile("./assets/road.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[1] = img_road
	img_grassland, _, err = ebitenutil.NewImageFromFile("./assets/grass.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[2] = img_grassland
	img_sand, _, err = ebitenutil.NewImageFromFile("./assets/sand.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[3] = img_sand
	img_forest, _, err = ebitenutil.NewImageFromFile("./assets/forest.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[5] = img_forest
	img_water, _, err = ebitenutil.NewImageFromFile("./assets/water.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[9] = img_water
	img_wall, _, err = ebitenutil.NewImageFromFile("./assets/wall.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[99] = img_wall
	img_cliff, _, err = ebitenutil.NewImageFromFile("./assets/cliff.png")
	if err != nil {
		log.Fatal(err)
	}
	terrain_map[98] = img_cliff
	img_player, _, err = ebitenutil.NewImageFromFile("./assets/bee.png")
	if err != nil {
		log.Fatal(err)
	}
	img_goal, _, err = ebitenutil.NewImageFromFile("./assets/goal.png")
	if err != nil {
		log.Fatal(err)
	}
	img_spawner, _, err = ebitenutil.NewImageFromFile("./assets/spawner.png")
	if err != nil {
		log.Fatal(err)
	}
	img_path, _, err = ebitenutil.NewImageFromFile("./assets/path.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_border_top, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_top.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_border_left, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_left.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_border_right, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_right.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_border_bottom, _, err = ebitenutil.NewImageFromFile("./assets/menu/border_bottom.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_tl, _, err = ebitenutil.NewImageFromFile("./assets/menu/tl.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_tr, _, err = ebitenutil.NewImageFromFile("./assets/menu/tr.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_bl, _, err = ebitenutil.NewImageFromFile("./assets/menu/bl.png")
	if err != nil {
		log.Fatal(err)
	}
	img_menu_br, _, err = ebitenutil.NewImageFromFile("./assets/menu/br.png")
	if err != nil {
		log.Fatal(err)
	}
	img_logo, _, err = ebitenutil.NewImageFromFile("./assets/menu/logo.png")
	if err != nil {
		log.Fatal(err)
	}
	img_button, _, err = ebitenutil.NewImageFromFile("./assets/menu/save.png")
	if err != nil {
		log.Fatal(err)
	}
	img_button_pressed, _, err = ebitenutil.NewImageFromFile("./assets/menu/save_depressed.png")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: REMOVE THIS LATER - placeholder for when no map is loaded
	terrain_map[0] = img_menu_border_top

	console.console_add("Images successfully loaded.")

	// // create sprite - first one in the array is the GOAL
	// // TODO: see about moving goal / enemy to separate array?
	goal_entity := target{
		name:       "Goal",
		loc:        GOAL,
		mob_type:   0,
		action:     actionDefence,
		inCombat:   false,
		alive:      true,
		sprite_img: img_goal,
		//movement_speed:     0,
		health:             10000000,
		armour:             100000,
		damage_per_attack:  20,
		attacks_per_second: 5,
		attack_success_pc:  0.95,
		attack_range:       5,
		last_attack_time:   0,
		target:             nil,
		//path:               []coords{},
	}
	target_list = append(target_list, &goal_entity)
	//go entity_list[0].brain()
	go combatCycle()

	console.console_add("Initialising menu...")
	init_Menu()

	// Init minimap
	initMinimap()

	// Init scoreboard (all scoreboard updating happens from here)
	initScoreboard()

	console.console_add("Init complete.")
}

func main() {
	// Create a background go-routine that polls for user movement. Maybe set variable sleep time based on movement speed?
	// for _, char := range entity_list {
	// 	go char.move_entity()
	// }

	ebiten.SetWindowSize(xScreen, yScreen)
	ebiten.SetWindowTitle("Ebiten Test")

	console.console_add("Entering game loop.")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}

}
