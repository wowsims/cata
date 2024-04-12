package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (mage *Mage) registerManaGemsCD() {

	var manaGain float64
	actionID := core.ActionID{ItemID: 36799} // this is the correct item ID for mana gem
	manaMetrics := mage.NewManaMetrics(actionID)
	hasT7_2pc := mage.HasSetBonus(ItemSetFrostfireGarb, 2)

	var T7GemAura *core.Aura
	if hasT7_2pc {
		T7GemAura = mage.NewTemporaryStatsAura("Improved Mana Gems T7", core.ActionID{SpellID: 61062}, stats.Stats{stats.SpellPower: 225}, 15*time.Second)
	}

	var serpentCoilAura *core.Aura
	if mage.HasTrinketEquipped(30720) {
		serpentCoilAura = mage.NewTemporaryStatsAura("Serpent-Coil Braid", core.ActionID{ItemID: 30720}, stats.Stats{stats.SpellPower: 225}, 15*time.Second)
	}

	//id for improved mana gem buff 83098
	//id for "Replenish Mana" is 5405
	//buff gives 2% of max mana as spell power for 15 seconds
	var ImprovedManaGemAura *core.Aura
	if mage.Talents.ImprovedManaGem > 0 {
		ImprovedManaGemAura = mage.NewTemporaryStatsAura("Improved Mana Gem",
			core.ActionID{SpellID: 83098},
			stats.Stats{stats.SpellPower: 0.01 * float64(mage.Talents.ImprovedManaGem) * mage.MaxMana()},
			15*time.Second)
	}

	// Multiply the mana restoral effect e.g. SCB restores +25% mana
	manaMultiplier := (1 +
		core.TernaryFloat64(serpentCoilAura != nil, 0.25, 0) +
		core.TernaryFloat64(hasT7_2pc, 0.25, 0))

	// Numbers may be different at 85, these are 80 values
	minManaGain := 3330.0 * manaMultiplier
	maxManaGain := 3500.0 * manaMultiplier

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = 3 // Now only have 1 gem, so 6 -> 3
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		// Don't use if we don't have any gems remaining!
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return remainingManaGems != 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {

			if remainingManaGems > 0 {
				manaGain = sim.Roll(minManaGain, maxManaGain)
			}

			if T7GemAura != nil {
				T7GemAura.Activate(sim)
			}
			if serpentCoilAura != nil {
				serpentCoilAura.Activate(sim)
			}

			if mage.Talents.ImprovedManaGem > 0 {
				ImprovedManaGemAura.Activate(sim)
			}

			mage.AddMana(sim, manaGain, manaMetrics)

			remainingManaGems--
			if remainingManaGems == 0 {
				mage.GetMajorCooldown(actionID).Disable()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Only pop if we have less than the max mana provided by the gem minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
			return character.MaxMana()-(character.CurrentMana()+totalRegen) >= maxManaGain
		},
	})
}
