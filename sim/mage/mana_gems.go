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

	// Numbers may be different at 85, these are 80 values
	minManaGain := 3330.0
	maxManaGain := 3500.0

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = 3
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
