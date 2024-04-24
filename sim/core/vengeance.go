package core

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core/stats"
)

type VengeanceTracker struct {
	eligibleDamage   float64
	apBonus          float64
	prevAPBonus      float64
	recentMaxAPBonus float64
	lastAttackedTime time.Duration // timestamp that the character was last attacked
}

const (
	VengeanceAPDecayRate     = 0.1 // AP bonus decays by 10% every 2 seconds, or 5% if the character has been hit in that time
	OutcomeVengeanceTriggers = OutcomeLanded
)

func Clamp(val float64, min float64, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

func UpdateVengeance(sim *Simulation, character *Character, tracker *VengeanceTracker, aura *Aura) {
	// Save the current AP bonus so we can apply the new buff correctly
	tracker.prevAPBonus = tracker.apBonus

	// If this character has been attacked in the last 2 seconds, apply half decay and add new damage to buff
	timeSinceLastHit := sim.CurrentTime - tracker.lastAttackedTime
	if timeSinceLastHit < time.Second*2 {

		// Decay existing bonus by half of the rate
		decay := VengeanceAPDecayRate / 2
		tracker.apBonus -= (decay * tracker.apBonus)

		// Add 5% of damage taken in the last 2 seconds
		tracker.apBonus += 0.05 * tracker.eligibleDamage

		// 4.3.0 change: the vengeance AP buff is always at least 33% of the incoming
		// damage if the tank has been hit in the last 2 seconds
		baseAPBonus := tracker.eligibleDamage / 3.0
		tracker.apBonus = math.Max(tracker.apBonus, baseAPBonus)
	} else {
		// No hits in the last 2 seconds - apply full decay
		tracker.apBonus -= (VengeanceAPDecayRate * tracker.recentMaxAPBonus)
	}

	// Vengeance tooltip is wrong in sake of simplicity as stated by blizzard
	// Actual formula used is Stamina + 10% of Base HP
	apBonusMax := character.GetStat(stats.Stamina) + 0.1*character.MaxHealth()
	tracker.apBonus = Clamp(tracker.apBonus, 0, apBonusMax)

	tracker.recentMaxAPBonus = math.Max(tracker.apBonus, tracker.recentMaxAPBonus)

	if sim.Log != nil {
		character.Log(sim, "Updated Vengeance for %s: Eligible Damage(%f) | AP Bonus(%f)", character.Name, tracker.eligibleDamage, tracker.apBonus)
	}

	tracker.eligibleDamage = 0

	// Update character stats
	character.AddStatDynamic(sim, stats.AttackPower, -tracker.prevAPBonus)
	character.AddStatDynamic(sim, stats.AttackPower, tracker.apBonus)
}

// To use: add a VengeanceTracker member to your spec-specific struct (e.g ProtWarrior, BloodDeathKnight, etc) then call this
// with your class's specific Vengeance spell ID
func ApplyVengeanceEffect(character *Character, tracker *VengeanceTracker, spellID int32) {

	// For sanity
	tracker.prevAPBonus = 0
	tracker.apBonus = 0
	tracker.eligibleDamage = 0
	tracker.recentMaxAPBonus = 0

	vengAura := MakePermanent(character.RegisterAura(Aura{
		Label:    "Vengeance",
		Duration: NeverExpires,
		ActionID: ActionID{SpellID: spellID}, // Different specs use different spell IDs even though the effect is the same
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Outcome.Matches(OutcomeVengeanceTriggers) {
				// Vengeance is based on the taken damage amount after mitigation
				// TODO: check how this treats dodge/parry/miss
				// https://worldofwarcraft.blizzard.com/en-us/news/1293873/tanking-with-a-vengeance seems to suggest a string of dodges will let it fall off
				// but simc's implementation retriggers vengeance on _any_ attack, even dodge/parry/miss.
				// I can't find any patch notes or other resources that support one or the other though
				tracker.lastAttackedTime = sim.CurrentTime
				tracker.eligibleDamage += result.Damage
			}
		},
	}))

	// Vengeance "ticks" every 2 seconds to update the AP buff
	character.RegisterResetEffect(func(sim *Simulation) {
		StartPeriodicAction(sim, PeriodicActionOptions{
			Period: time.Second * 2,
			OnAction: func(sim *Simulation) {
				UpdateVengeance(sim, character, tracker, vengAura)
			},
		})
	})
}
