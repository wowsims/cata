package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerLightningShieldSpell() {
	if shaman.SelfBuffs.Shield != proto.ShamanShield_LightningShield {
		return
	}

	actionID := core.ActionID{SpellID: 324}

	shaman.LightningShieldDamage = shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 26364},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		ClassSpellMask:   SpellMaskLightningShield,
		DamageMultiplier: 1,
		ThreatMultiplier: 1, //fix when spirit weapons is fixed
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 0.38800001144,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.CalcScalingSpellDmg(0.56499999762)
			result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)
			spell.DealDamage(sim, result)
		},
	})

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 3500,
	}

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 60,
		MaxStacks: 7,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 1)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}
			if !icd.IsReady(sim) {
				return
			}
			icd.Use(sim)

			aura.RemoveStack(sim)
			shaman.LightningShieldDamage.Cast(sim, spell.Unit)
		},
	})

	shaman.LightningShield = shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			shaman.LightningShieldAura.Activate(sim)
		},
	})
}
