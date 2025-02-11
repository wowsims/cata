package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerShamanisticRageCD() {
	if !shaman.Talents.ShamanisticRage {
		return
	}

	actionID := core.ActionID{SpellID: 30823}
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.7
			shaman.PseudoStats.SpellCostPercentModifier -= 100
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.7
			shaman.PseudoStats.SpellCostPercentModifier += 100
		},
	})

	// 2pc T10 bonus aura
	dummySetAura := shaman.RegisterAura(core.Aura{
		Label:    "Frost Witch's Battlegear (2pc)",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Deactivate(sim)
		},
	})

	if shaman.usePrepullEnh_2PT10 && !shaman.CouldHaveSetBonus(ItemSetFrostWitchBattlegear, 2) {
		SharedEnhTier102PCAura(shaman, dummySetAura)
	}

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: SpellMaskShamanisticRage,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			if sim.CurrentTime < 0 {
				dummySetAura.Activate(sim)
			} else {
				dummySetAura.Deactivate(sim)
			}
			srAura.Activate(sim)
		},
		RelatedSelfBuff: srAura,
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
	})
}
