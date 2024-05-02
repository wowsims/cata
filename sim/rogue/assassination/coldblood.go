package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/rogue"
)

func (sinRogue *AssassinationRogue) registerColdBloodCD() {
	if !sinRogue.Talents.ColdBlood {
		return
	}

	actionID := core.ActionID{SpellID: 14177}
	cbEnergyMetric := sinRogue.NewEnergyMetrics(actionID)

	coldBloodAura := sinRogue.RegisterAura(core.Aura{
		Label:    "Cold Blood",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range sinRogue.Spellbook {
				if spell.Flags.Matches(rogue.SpellFlagColdBlooded) {
					spell.BonusCritRating += 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range sinRogue.Spellbook {
				if spell.Flags.Matches(rogue.SpellFlagColdBlooded) {
					spell.BonusCritRating -= 100 * core.CritRatingPerCritChance
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// for Mutilate, the offhand hit comes first and is ignored, so the aura doesn't fade too early
			if spell.Flags.Matches(rogue.SpellFlagColdBlooded) && spell.ProcMask.Matches(core.ProcMaskMeleeMH|core.ProcMaskRangedSpecial) {
				aura.Deactivate(sim)
			}
		},
	})

	sinRogue.ColdBlood = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    sinRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			coldBloodAura.Activate(sim)
			sinRogue.AddEnergy(sim, 25, cbEnergyMetric)
		},
	})

	sinRogue.AddMajorCooldown(core.MajorCooldown{
		Spell: sinRogue.ColdBlood,
		Type:  core.CooldownTypeDPS,
	})
}
