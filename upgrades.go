package main

/// buy bone collector - repairs spawners
//var upgradeSpawnerSpeedIndex int
//var upgradeSpawn []

const tierBounds = 6 // there are 7 tiers

const (
	spawnerBaseCost = 10
)

type upgrade struct {
	field int
	tier  int
	cost  []int
	value []float64
}

// Spawn walking speed
// var speed = upgrade{
// 	state: 0,
// 	cost:  []int{0, 50, 100, 200, 500, 1000, 10000},
// 	value: []float64{10, 12, 15, 18, 22, 25, 30},
// }

// Spawn dmg per attack
var spawnDmg = upgrade{
	field: 0,
	tier:  0,
	cost:  []int{0, 500, 1000, 2000, 5000, 10000, 100000},
	value: []float64{10, 12.5, 15, 20, 25, 32.5, 45},
}

// Spawn attack speed (measured in attacks per second)
var spawnAttackRate = upgrade{
	field: 1,
	tier:  0,
	cost:  []int{0, 500, 1000, 2000, 5000, 10000, 100000},
	value: []float64{0.2, 0.5, 1, 1.25, 1.5, 2, 2.5},
}

// Spawn % to hit
// Maybe allow downgrade - lower toHit == higher crit %
var spawnToHit = upgrade{}

// Spawn max health
var spawnHealth = upgrade{}

// Spawn armour
var spawnArmour = upgrade{}

// SPAWNER upgrades

// Spawner rate of spawning
var spawnerRate = upgrade{
	field: 0,
	tier:  0,
	cost:  []int{0, 500, 1000, 2000, 5000, 10000, 100000},
	value: []float64{0.2, 0.5, 1, 1.25, 1.5, 2, 2.5},
}

// Spawner max health - NOT REQUIRED
//var spawnerHealth = upgrade{}

// Spawner max number of spawns before exhaustion
// effects all current and future spawners
var spawnerMaxOutput = upgrade{
	field: 1,
	tier:  0,
	cost:  []int{0, 5000, 10000, 20000, 50000, 100000, 1000000},
	value: []float64{10, 25, 50, 100, 150, 250, 500},
}

func upgradeSpawners(up *upgrade) {
	switch up.field {
	case 0: // Spawn rate
		if up.tier < tierBounds {
			up.tier++
			score = score - up.cost[up.tier]
			for _, s := range spawner_list {
				s.spawn_per_second = up.value[up.tier]
			}
		}
	case 1: // Max output
		if up.tier < tierBounds {
			up.tier++
			score = score - up.cost[up.tier]
			for _, s := range spawner_list {
				s.energy = int(spawnerMaxOutput.value[spawnerMaxOutput.tier])
			}
		}
	}

}
