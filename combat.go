package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	critModifier         = 5    // Multiplier for crits. eg. If critModifier == 2, dmg done for crit is dmg*2.
	armourDeratingFactor = 0.99 // Multiplier applied to armour when hit. Must be <= 1. eg. If armourDebuff == 0.99
	// and an entity with 100 armour is hit, their armour is reduced to (100 * 0.99 = 99)
	maxArmour = 1000
)

// For each entity, find if their target is in range.
// For fixed locale (ie bases), find their nearest non-friendly neighbour.
func (e *entity) findEnemy() *entity {
	var foe *entity

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
func (self *entity) attackEnemy(foe *entity) {
	// Damage is calculated as:
	now := time.Now().UnixMilli()
	// If sufficient time has passed, attack again.
	if int(now) > self.last_attack_time+int(1000/self.attacks_per_second) {
		// Update last attack time, regardless of whether the attack is a success or not. (You don't get a freebie for missing)
		self.last_attack_time = int(now)
		if self.successfullyHits(foe) {
			// TODO: assign points for successful hit.
			self.applyDamage(foe)
			if foe.alive == false {
				// TODO: assign points for successful kill.
				self.target = nil
			}
		}
	}
}

// Calculates whether an attack is likely to succeed or not.
// Factors in % to hit.
func (s *entity) successfullyHits(f *entity) bool {
	if s.attack_success_pc >= rand.Float32() {
		return true
	}
	return false
}

// Calculates the damage to be applied to the target.
// Factors in enemy armour damage reduction.
// Calls out to assess if attack is a Crit, and modifies dmg accordingly.
func (s *entity) calculateDamage(f *entity) float32 {
	base_dmg := s.damage_per_attack * float32(s.attackCrits(f))
	armourDebuff := (maxArmour - f.armour) / maxArmour

	// calculate armour debuff. At armour > 1000, zero damage is inflicted.
	armourDebuff = float32(math.Max(float64(armourDebuff), 0))

	return base_dmg * armourDebuff
}

// Calculates whether an attack crits. Returns a multiplier to be applied to the
// attack.
func (s *entity) attackCrits(f *entity) int {
	if (1 - s.attack_success_pc) >= rand.Float32() {
		return critModifier
	}
	return 1
}

// Calculates the damage to be applied, and passes it to the target,
// who processes the Take Damage function.
func (s *entity) applyDamage(f *entity) {
	// Work out how much damage is to be applied.
	// Automatically factors in crits.
	f.takeDamage(s.calculateDamage(f))
}

// When being attacked, an entity my take damage. This function applies that
// damage, and marks the entity as no longer alive if appropriate. Also
// applies armour derating (damage reduces efficacy).
func (me *entity) takeDamage(dmg float32) {
	// Apply armour derate
	me.armour = me.armour * armourDeratingFactor
	me.health = me.health - dmg
	if me.health <= 0 {
		me.alive = false
	}
}
