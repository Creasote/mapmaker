package main

import (
	"time"
)

// Adds a generic player entity. Use as a template for various entity types.
func placeEntity(x, y int) {

	// Set up the target list, pre-populate the goal entity.
	t := []*entity{}
	t = append(t, entity_list[0])

	entity_list = append(entity_list, &entity{
		name:               "Clicker",
		loc:                coords{x, y},
		mob_type:           0,
		action:             actionAttackEnemyBase,
		inCombat:           false,
		alive:              true,
		sprite_img:         img_player,
		movement_speed:     50,
		health:             100,
		armour:             10,
		damage_per_attack:  10,
		attacks_per_second: 1,
		attack_success_pc:  0.75,
		attack_range:       2,
		last_attack_time:   0,
		target:             t,
		path:               []coords{},
	})
}

// Place a spawner that regularly spawns new mobs.
func placeSpawner(x, y int) {
	// Generate a new entity for the spawner for drawing.
	t := []*entity{}
	t = append(t, entity_list[0])
	entity_list = append(entity_list, &entity{
		name:               "Spawner",
		loc:                coords{x, y},
		mob_type:           0,
		alive:              true,
		action:             actionDefence,
		inCombat:           false,
		sprite_img:         img_spawner,
		movement_speed:     0,
		health:             1,
		armour:             0,
		damage_per_attack:  0,
		attacks_per_second: 0,
		attack_success_pc:  0,
		attack_range:       0,
		last_attack_time:   0,
		target:             nil,
		path:               []coords{},
	})

	// Generate a pre-canned path for all spawned mobs
	//precannedPath := pathfind(&game_map)
	for {
		// t := []*entity{}
		// t = append(t, entity_list[0])
		// Generate a new spawned entity every given period.
		entity_list = append(entity_list, &entity{
			name:               "Spawned",
			loc:                coords{x, y},
			mob_type:           0,
			action:             actionAttackEnemyBase,
			inCombat:           false,
			alive:              true,
			sprite_img:         img_player,
			movement_speed:     50,
			health:             100,
			armour:             10,
			damage_per_attack:  10,
			attacks_per_second: 1,
			attack_success_pc:  0.75,
			attack_range:       2,
			last_attack_time:   0,
			target:             t,
			path:               []coords{},
		})

		time.Sleep(5000 * time.Millisecond)
	}
}

// Sets the Goal location.
func setGoal(x, y int) {
	entity_list[0].loc = coords{x, y}
	game_map[entity_layer][y][x] = 1 // Sets goal, able to be retrieved on load.
}

// Entity brain is responsible for deciding what the entity will do.
// It is responsible for measuring entity vitals, and determining a course of action.
// Inputs can include: health, hunger, target information
// Actions can include: attacking, fleeing, feeding
func (me *entity) brain() {
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

	case actionAttackEnemyBase:
		if estimate_distance(me.loc, entity_list[0].loc) <= me.attack_range {
			// attack
			me.attackEnemy(me.target[0])
		} else if len(me.path) == 0 {
			// TODO: investigate why pathing glitches when using go-routines.
			me.pathfind(&game_map)
		} else {
			// Too far from target, and path is set, so get moving.
			me.move_entity()
		}
	}
	//}
	//time.Sleep(1000 * time.Millisecond)
}

func (me *entity) move_entity() {
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

func (me *entity) die() {
	me.alive = false
	me.target = nil
	me.action = actionRest
}

// ENTITY TYPES

//

// create sprite - first one in the array is the GOAL
// TODO: see about moving goal / enemy to separate array?
// var goal_entity = entity{
// 	name:               "Goal",
// 	loc:                GOAL,
// 	mob_type:           0,
// 	action:             actionDefence,
// 	inCombat:           false,
// 	alive:              true,
// 	sprite_img:         img_goal,
// 	movement_speed:     0,
// 	health:             100000,
// 	armour:             10000,
// 	damage_per_attack:  20,
// 	attacks_per_second: 5,
// 	attack_success_pc:  0.95,
// 	attack_range:       5,
// 	last_attack_time:   0,
// 	target:             nil,
// 	path:               []coords{},
// }

var entityPathfinder = entity{
	name: "Spawned Zombie",
	//loc:                coords{x, y},
	mob_type:           0,
	action:             actionAttackEnemyBase,
	inCombat:           false,
	alive:              true,
	sprite_img:         img_player,
	movement_speed:     50,
	health:             100,
	armour:             10,
	damage_per_attack:  10,
	attacks_per_second: 1,
	attack_success_pc:  0.75,
	attack_range:       2,
	last_attack_time:   0,
	//target:             t,
	//path:               []coords{},
}

//entity_list = append(entity_list, &goal_entity)
