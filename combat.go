package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	critModifier         = 5     // Multiplier for crits. eg. If critModifier == 2, dmg done for crit is dmg*2.
	armourDeratingFactor = 0.999 // Multiplier applied to armour when hit. Must be <= 1. eg. If armourDebuff == 0.99
	// and an entity with 100 armour is hit, their armour is reduced to (100 * 0.99 = 99)
)

func combatCycle() {
	for !gameOverFlag {
		for _, e := range spawn_list {
			e.brain()
		}
		for _, e := range target_list {
			e.brain()
		}
		// TODO: Disabled until use case established.
		for _, e := range spawner_list {
			e.brain()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

/////////////////////
// SPAWN FUNCTIONS //
/////////////////////

// For each entity, find if their target is in range.
// For fixed locale (ie bases), find their nearest non-friendly neighbour.
func (e *spawn) findEnemy() *target {
	var foe *target

	// scan the map for the nearest available target.
	for _, x_array := range game_map[entity_layer][maxInt(0, e.loc.y-int(e.attack_range)):minInt(e.loc.y+int(e.attack_range), board_cells_High)] {
		for _, x_val := range x_array[maxInt(0, e.loc.x-int(e.attack_range)):minInt(e.loc.x+int(e.attack_range), board_cells_Wide)] {
			if x_val > 0 {
				// TODO: this need to iterate over a different array, that is an array of entities.
				// The game board entity layer array is purely for drawing. Is it still required?
				return foe
			}
		}
	}

	return foe
}

// Begin combat system.
// On game tick (or timer?), each entity attacks their chosen target.
func (self *spawn) attackEnemy(foe *target) {
	// If combat has not yet been initiated, do it now.
	if !self.inCombat {
		self.initiateCombat(foe)
	}
	// Damage is calculated as:
	now := time.Now().UnixMilli()
	// If sufficient time has passed, attack again.
	if int(now) > self.last_attack_time+int(1000/self.attacks_per_second) {
		// Update last attack time, regardless of whether the attack is a success or not. (You don't get a freebie for missing)
		self.last_attack_time = int(now)
		if self.successfullyHits(foe) {
			addScore(pointsSuccessfulHit)
			// TODO: assign points for successful hit.
			self.applyDamage(foe)
			if foe.alive == false {
				// TODO: assign points for successful kill.
				if len(self.target) > 1 {
					self.target = self.target[1:]
				} else {
					self.target = nil
				}
				self.inCombat = false
			}
		} else {
			console.console_add("Missed target")
		}
	}
}

// Calculates whether an attack is likely to succeed or not.
// Factors in % to hit.
func (s *spawn) successfullyHits(f *target) bool {
	if s.attack_success_pc >= rand.Float64() {
		return true
	}
	return false
}

// Calculates the damage to be applied to the target.
// Factors in enemy armour damage reduction.
// Calls out to assess if attack is a Crit, and modifies dmg accordingly.
func (s *spawn) calculateDamage(f *target) float64 {
	base_dmg := s.damage_per_attack * float64(s.attackCrits(f))

	// Damage tends toward zero at high armour levels ( > 100), but should be non-zero at all levels.
	return base_dmg - math.Pow(base_dmg, (f.armour/(f.armour+1)))
}

// Calculates whether an attack crits. Returns a multiplier to be applied to the
// attack.
func (s *spawn) attackCrits(f *target) int {
	if (1 - s.attack_success_pc) >= rand.Float64() {
		return critModifier
	}
	return 1
}

// Calculates the damage to be applied, and passes it to the target,
// who processes the Take Damage function.
func (s *spawn) applyDamage(f *target) {
	// Work out how much damage is to be applied.
	// Automatically factors in crits.
	f.takeDamage(s.calculateDamage(f))
}

// When being attacked, an entity my take damage. This function applies that
// damage, and marks the entity as no longer alive if appropriate. Also
// applies armour derating (damage reduces efficacy).
func (me *target) takeDamage(dmg float64) {
	// Apply armour derate
	me.armour = me.armour * armourDeratingFactor
	me.health = me.health - dmg
	console.console_add("Taking damage: " + fmt.Sprint(dmg))
	console.console_add("Health remaining: " + fmt.Sprint(me.health))
	dps[0] += dmg
	if me.health <= 0 {
		me.die()
	}
}

// Sends a signal to a target that they are being attacked.
func (me *spawn) initiateCombat(f *target) {
	me.inCombat = true
	f.receiveCombat(me)
}

// Recieves a signal that entity is under attack, and adds the attacker
// to the (end of the) target list.
func (me *target) receiveCombat(f *spawn) {
	me.inCombat = true
	me.target = append(me.target, f)
}

//////////////////////
// TARGET FUNCTIONS //
//////////////////////
// For each entity, find if their target is in range.
// For fixed locale (ie bases), find their nearest non-friendly neighbour.
func (e *target) findEnemy() *spawn {
	var foe *spawn

	// scan the map for the nearest available target.
	for _, x_array := range game_map[entity_layer][maxInt(0, e.loc.y-int(e.attack_range)):minInt(e.loc.y+int(e.attack_range), board_cells_High)] {
		for _, x_val := range x_array[maxInt(0, e.loc.x-int(e.attack_range)):minInt(e.loc.x+int(e.attack_range), board_cells_Wide)] {
			if x_val > 0 {
				// TODO: this need to iterate over a different array, that is an array of entities.
				// The game board entity layer array is purely for drawing. Is it still required?
				return foe
			}
		}
	}

	return foe
}

// Begin combat system.
// On game tick (or timer?), each entity attacks their chosen target.
func (self *target) attackEnemy(foe *spawn) {
	// If combat has not yet been initiated, do it now.
	if !self.inCombat {
		self.initiateCombat(foe)
	}
	// Damage is calculated as:
	now := time.Now().UnixMilli()
	// If sufficient time has passed, attack again.
	if int(now) > self.last_attack_time+int(1000/self.attacks_per_second) {
		console.console_add("Attacking.")
		// Update last attack time, regardless of whether the attack is a success or not. (You don't get a freebie for missing)
		self.last_attack_time = int(now)
		if self.successfullyHits(foe) {
			console.console_add("Successful hit.")
			addScore(pointsSuccessfulHit)
			// TODO: assign points for successful hit.
			self.applyDamage(foe)
			if foe.alive == false {
				// TODO: assign points for successful kill.
				if len(self.target) > 1 {
					self.target = self.target[1:]
				} else {
					self.target = nil
				}
				self.inCombat = false
			}
		} else {
			console.console_add("Missed target")
		}
	}
}

// Calculates whether an attack is likely to succeed or not.
// Factors in % to hit.
func (s *target) successfullyHits(f *spawn) bool {
	if s.attack_success_pc >= rand.Float64() {
		return true
	}
	return false
}

// Calculates the damage to be applied to the target.
// Factors in enemy armour damage reduction.
// Calls out to assess if attack is a Crit, and modifies dmg accordingly.
func (s *target) calculateDamage(f *spawn) float64 {
	base_dmg := s.damage_per_attack * float64(s.attackCrits(f))

	// Damage tends toward zero at high armour levels ( > 100), but should be non-zero at all levels.
	return base_dmg - math.Pow(base_dmg, (f.armour/(f.armour+1)))
}

// Calculates whether an attack crits. Returns a multiplier to be applied to the
// attack.
func (s *target) attackCrits(f *spawn) int {
	if (1 - s.attack_success_pc) >= rand.Float64() {
		return critModifier
	}
	return 1
}

// Calculates the damage to be applied, and passes it to the target,
// who processes the Take Damage function.
func (s *target) applyDamage(f *spawn) {
	// Work out how much damage is to be applied.
	// Automatically factors in crits.
	f.takeDamage(s.calculateDamage(f))
}

// When being attacked, an entity my take damage. This function applies that
// damage, and marks the entity as no longer alive if appropriate. Also
// applies armour derating (damage reduces efficacy).
func (me *spawn) takeDamage(dmg float64) {
	// Apply armour derate
	me.armour = me.armour * armourDeratingFactor
	me.health = me.health - dmg
	console.console_add("Taking damage: " + fmt.Sprint(dmg))
	console.console_add("Health remaining: " + fmt.Sprint(me.health))
	dps[0] += dmg
	if me.health <= 0 {
		me.die()
	}
}

// Sends a signal to a target that they are being attacked.
func (me *target) initiateCombat(f *spawn) {
	me.inCombat = true
	f.receiveCombat(me)
}

// Recieves a signal that entity is under attack, and adds the attacker
// to the (end of the) target list.
func (me *spawn) receiveCombat(f *target) {
	me.inCombat = true
	me.target = append(me.target, f)
}
