package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type spawn struct {
	name               string
	loc                coords
	mob_type           int
	alive              bool
	action             int  // Index of action currently being undertaken
	inCombat           bool // flag True when initiating combat. Reset to false when target dies.
	sprite_img         *ebiten.Image
	movement_speed     float64 // Should be between 0 (no movement) and ~ 50. Higher values may be OP.
	health             float64
	energy             int     //
	armour             float64 // Should be between 0 (no reduction) and 1000 (full reduction in damage taken)
	damage_per_attack  float64
	attacks_per_second float64
	attack_success_pc  float64 // Should be [0,1). Crit chance is calculated as the difference between 1 and attack_success_pc.
	attack_range       float64
	last_attack_time   int
	target             []*target
	path               []coords
}

type spawner struct {
	name             string
	loc              coords
	mob_type         int
	alive            bool
	action           int // Index of action currently being undertaken
	sprite_img       *ebiten.Image
	health           float64
	energy           int     // Expends 1 energy per spawn
	armour           float64 // Should be between 0 (no reduction) and 1000 (full reduction in damage taken)
	spawn_per_second float64
	last_spawn_time  int
	target           []*target
	path             []coords
}

type target struct {
	name               string
	loc                coords
	mob_type           int
	alive              bool
	action             int  // Index of action currently being undertaken
	inCombat           bool // flag True when initiating combat. Reset to false when target dies.
	sprite_img         *ebiten.Image
	health             float64
	armour             float64 // Should be between 0 (no reduction) and 1000 (full reduction in damage taken)
	damage_per_attack  float64
	attacks_per_second float64
	attack_success_pc  float64 // Should be [0,1). Crit chance is calculated as the difference between 1 and attack_success_pc.
	attack_range       float64
	last_attack_time   int
	target             []*spawn
	//path               []coords
}

// Adds a generic player entity. Use as a template for various entity types.
func placeSpawn(x, y int) {

	// Set up the target list, pre-populate the goal entity.
	t := []*target{}
	t = append(t, target_list[0])

	spawn_list = append(spawn_list, &spawn{
		name:               "Clicker",
		loc:                coords{x, y},
		mob_type:           0,
		action:             actionAttackEnemyBase,
		inCombat:           false,
		alive:              true,
		sprite_img:         img_player,
		movement_speed:     50,
		health:             spawnHealth.value[spawnHealth.tier],
		energy:             100,
		armour:             spawnArmour.value[spawnArmour.tier],
		damage_per_attack:  spawnDmg.value[spawnDmg.tier],
		attacks_per_second: spawnAttackRate.value[spawnAttackRate.tier],
		attack_success_pc:  spawnToHit.value[spawnToHit.tier],
		attack_range:       2,
		last_attack_time:   0,
		target:             t,
		path:               []coords{},
	})
	spawnCount++
	addScore(pointsManualSpawn)
}

// Place a spawner that regularly spawns new mobs.
func placeSpawner(x, y int) {
	// Pay for your goods first
	if score > power(spawnerBaseCost, len(spawner_list)) { //(spawnerBaseCost * len(spawner_list)) {
		score = score - power(spawnerBaseCost, len(spawner_list)) //(spawnerBaseCost * len(spawner_list))
		// Generate a new entity for the spawner for drawing.
		t := []*target{}
		t = append(t, target_list[0])

		spawner_list = append(spawner_list, &spawner{
			name:             "Spawner",
			loc:              coords{x, y},
			mob_type:         0,
			alive:            true,
			action:           actionAttackEnemyBase,
			sprite_img:       img_spawner,
			health:           1,
			energy:           int(spawnerMaxOutput.value[spawnerMaxOutput.tier]),
			armour:           0,
			spawn_per_second: spawnerRate.value[spawnerRate.tier],
			target:           t,
			path:             []coords{},
		})
		spawner_list[0].pathfind(&game_map)
	}
}

// Sets the Goal location.
func setGoal(x, y int) {
	target_list[0].loc = coords{x, y}
	game_map[entity_layer][y][x] = 1 // Sets goal, able to be retrieved on load.
}

// Entity brain is responsible for deciding what the entity will do.
// It is responsible for measuring entity vitals, and determining a course of action.
// Inputs can include: health, hunger, target information
// Actions can include: attacking, fleeing, feeding
func (me *spawn) brain() {
	//console.console_add("Pathfinder brain initiated.")
	// Step 1: Decide what to do
	// Step 2: Set target
	// Step 3: Navigate to target
	// Step 4: Do action
	// Repeat.

	// If no target, select one. Priority chain:
	// 1. Entity attacking me.
	// 2. Food / energy source.
	// 3. Reset / recovery.
	// 4. Other imperatives (resource gathering, repairing)
	// 5. Global goal (enemy base)

	// TODO: implement targeting chain logic.
	// me.setTarget()

	// Action chain:
	// 1. Self defence
	// 2. Feed
	// 3. Recover
	// 4. Local goal (attack enemy agents)
	// 5. Base repair
	// 6. Resource gathering
	// 7. Global goal (attack enemy base)
	if !(me.energy > 0) {
		me.health--
	}
	if !(me.health > 0) {
		me.die()
	}
	switch me.action {
	case actionDefence:
		if len(me.target) > 0 {
			me.attackEnemy(me.target[0])
		}
	case actionRest:
		// TODO: add some regeneration

	case actionAttackEnemyBase:
		if estimate_distance(me.loc, me.target[0].loc) <= me.attack_range {
			// attack
			me.attackEnemy(me.target[0])
		} else if len(me.path) == 0 {
			me.pathfind(&game_map)
		} else {
			// Too far from target, and path is set, so get moving.
			me.move_entity()
		}
	}
}

func (me *spawner) brain() {
	//TODO: what actions does a spawner have? Just spawning?
	if !(me.energy > 0) {
		me.health--
	}
	if !(me.health > 0) {
		me.die()
	}
	switch me.action {
	case actionAttackEnemyBase:
		now := time.Now().UnixMilli()
		if now > int64(me.last_spawn_time)+(int64(1000/me.spawn_per_second)) {
			me.last_spawn_time = int(now)
			spawn_list = append(spawn_list, &spawn{
				name:               "Spawned",
				loc:                me.loc,
				mob_type:           0,
				action:             actionAttackEnemyBase,
				inCombat:           false,
				alive:              true,
				sprite_img:         img_player,
				movement_speed:     50,
				health:             spawnHealth.value[spawnHealth.tier],
				energy:             100,
				armour:             spawnArmour.value[spawnArmour.tier],
				damage_per_attack:  spawnDmg.value[spawnDmg.tier],
				attacks_per_second: spawnAttackRate.value[spawnAttackRate.tier],
				attack_success_pc:  spawnToHit.value[spawnToHit.tier],
				attack_range:       2,
				last_attack_time:   0,
				target:             me.target,
				path:               me.path,
			})
			spawnCount++
			addScore(pointsAutoSpawn)
			me.energy--
		}
	default:
		// ruh-roh
	}

}

func (me *target) brain() {
	//console.console_add("Pathfinder brain initiated.")
	// Step 1: Decide what to do
	// Step 2: Set target
	// Step 3: Navigate to target
	// Step 4: Do action
	// Repeat.

	// If no target, select one. Priority chain:
	// 1. Entity attacking me.
	// 2. Food / energy source.
	// 3. Reset / recovery.
	// 4. Other imperatives (resource gathering, repairing)
	// 5. Global goal (enemy base)

	// TODO: implement targeting chain logic.
	// me.setTarget()

	// Action chain:
	// 1. Self defence
	// 2. Feed
	// 3. Recover
	// 4. Local goal (attack enemy agents)
	// 5. Base repair
	// 6. Resource gathering
	// 7. Global goal (attack enemy base)

	// Attack enemy base
	// me.setTarget(enemy_base)
	//for me.alive {
	switch me.action {
	case actionDefence:
		if len(me.target) > 0 {
			me.attackEnemy(me.target[0])
		}
	case actionRest:
		// TODO: add some regeneration
		// time.Sleep(1000 * time.Millisecond)

		// case actionAttackEnemyBase:
		// 	if estimate_distance(me.loc, me.target[0].loc) <= me.attack_range {
		// 		// attack
		// 		me.attackEnemy(me.target[0])
		// 	} else if len(me.path) == 0 {
		// 		// TODO: investigate why pathing glitches when using go-routines.
		// 		me.pathfind(&game_map)
		// 	} else {
		// 		// Too far from target, and path is set, so get moving.
		// 		me.move_entity()
		// 	}
	}
}

func (me *spawn) move_entity() {
	//for len(me.path) > 0 {
	if len(me.path) > 0 {
		// TODO: remove loop, change to "while me.task == <some movement interface>"
		//if len(me.path) > 0 {
		// Set the location to the next waypoint, and remove that waypoint from the path list.
		// TODO: Include calculation that slows down based on terrain.
		me.loc, me.path = me.path[len(me.path)-1], me.path[:len(me.path)-1]
		//}

		// TODO: This is an arbitrary delay. Review. Change to move "units" per time, which are "spent" moving over terrain.
		//time.Sleep(time.Duration(10000/me.movement_speed) * time.Millisecond)
	}
}

func (me *target) die() {
	me.alive = false
	me.target = nil
	me.action = actionRest
	gameOverFlag = true
}

func (me *spawn) die() {
	me.alive = false
	me.target = nil
	me.action = actionRest
}

func (me *spawner) die() {
	me.alive = false
	me.action = actionRest
}
