package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) registerManaGemsCD() {

	var manaGain float64
	actionID := core.ActionID{ItemID: 36799} // this is the correct item ID for mana gem
	manaMetrics := mage.NewManaMetrics(actionID)

	//buff gives points% of max mana as spell power for 15 seconds
	var improvedManaGemAura *core.Aura
	if mage.Talents.ImprovedManaGem > 0 {
		spBonusMod := mage.AddDynamicMod(core.SpellModConfig{
			Kind:      core.SpellMod_BonusSpellPower_Flat,
			ClassMask: MageSpellsAll,
		})

		improvedManaGemAura = mage.GetOrRegisterAura(core.Aura{
			Label:    "Improved Mana Gem",
			ActionID: core.ActionID{SpellID: []int32{0, 31584, 31585}[mage.Talents.ImprovedManaGem]},
			Duration: time.Second * 15,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				spBonusMod.UpdateFloatValue(0.01 * float64(mage.Talents.ImprovedManaGem) * mage.MaxMana())
				spBonusMod.Activate()
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				spBonusMod.Deactivate()
			},
		})
	}

	minManaGain := 11801.0
	maxManaGain := 12405.0

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
			if mage.Talents.ImprovedManaGem > 0 {
				improvedManaGemAura.Activate(sim)
			}

			manaGain = sim.Roll(minManaGain, maxManaGain)
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
			totalRegen := character.ManaRegenPerSecondWhileCombat() * 5
			return character.MaxMana()-(character.CurrentMana()+totalRegen) >= maxManaGain
		},
	})
}
