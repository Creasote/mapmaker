package main

/// buy bone collector - repairs spawners
//var upgradeSpawnerSpeedIndex int
//var upgradeSpawn []

const tierBounds = 6 // there are 7 tiers

const (
	spawnerBaseCost = 100
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
var spawnToHit = upgrade{
	field: 2,
	tier:  0,
	cost:  []int{500, 500, 500, 500, 500, 500, 500},
	value: []float64{0.75, 0.8, 0.85, 0.95, 0.99, 0.25, 0.50},
}

// Spawn max health
var spawnHealth = upgrade{
	field: 3,
	tier:  0,
	cost:  []int{0, 5000, 10000, 20000, 50000, 100000, 1000000},
	value: []float64{100, 125, 150, 200, 250, 500, 750},
}

// Spawn armour
var spawnArmour = upgrade{
	field: 4,
	tier:  0,
	cost:  []int{0, 5000, 10000, 20000, 50000, 100000, 1000000},
	value: []float64{10, 12.5, 16, 20, 25, 32.5, 40},
}

// func upgradeSpawns(up *upgrade){
// 	switch
// }
// SPAWNER upgrades

// Spawner rate of spawning (measured as spawns per second)
var spawnerRate = upgrade{
	field: 10,
	tier:  0,
	cost:  []int{0, 500, 1000, 2000, 5000, 10000, 100000},
	value: []float64{0.2, 0.5, 1, 1.25, 1.5, 2, 2.5},
}

// Spawner max number of spawns before exhaustion
// effects all current and future spawners
var spawnerMaxOutput = upgrade{
	field: 11,
	tier:  0,
	cost:  []int{0, 5000, 10000, 20000, 50000, 100000, 1000000},
	value: []float64{100000000, 25, 50, 100, 150, 250, 500},
}

func doUpgrade(up *upgrade) {
	if up.tier < tierBounds && score > up.cost[up.tier+1] {
		up.tier++
		score = score - up.cost[up.tier]
		switch up.field {
		case 0: // Spawn Damage
			for _, s := range spawn_list {
				s.damage_per_attack = up.value[up.tier]
			}
		case 1: // Spawn attack rate
			for _, s := range spawn_list {
				s.attacks_per_second = up.value[up.tier]
			}
		case 2: // Spawn To-hit
			for _, s := range spawn_list {
				s.attack_success_pc = up.value[up.tier]
			}
		case 3: // Spawn health
			for _, s := range spawn_list {
				s.health = up.value[up.tier]
			}
		case 4: // Spawn armour
			for _, s := range spawn_list {
				s.armour = up.value[up.tier]
			}
		case 10: // Spawn rate
			for _, s := range spawner_list {
				s.spawn_per_second = up.value[up.tier]
			}

		case 11: // Max output
			score = score - up.cost[up.tier]
			for _, s := range spawner_list {
				s.energy = int(spawnerMaxOutput.value[spawnerMaxOutput.tier])
			}
		}
	}

}

func (up upgrade) getCost(offset int) int {
	return up.cost[up.tier+offset]
}

func (up upgrade) do() {
	doUpgrade(&up)
}
