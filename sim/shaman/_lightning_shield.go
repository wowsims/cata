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

	procSpell := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 26364},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		ClassSpellMask:   SpellMaskLightningShield,
		DamageMultiplier: 1,
		ThreatMultiplier: 1, //fix when spirit weapons is fixed
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		BonusCoefficient: 0.267,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.ClassSpellScaling * 0.38899999857
			result := shaman.calcDamageStormstrikeCritChance(sim, target, baseDamage, spell)
			spell.DealDamage(sim, result)
		},
	})

	icd := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 3500,
	}

	icdStaticShock := core.Cooldown{
		Timer:    shaman.NewTimer(),
		Duration: time.Millisecond * 500,
	}

	shaman.LightningShieldAura = shaman.RegisterAura(core.Aura{
		Label:     "Lightning Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 10,
		MaxStacks: 9,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.SetStacks(sim, 3)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !icdStaticShock.IsReady(sim) || !aura.IsActive() {
				return
			}
			if shaman.Talents.StaticShock > 0 && spell == shaman.LavaLash || spell == shaman.Stormstrike || spell == shaman.PrimalStrike {
				if sim.RandomFloat("Static Shock") < 0.15*float64(shaman.Talents.StaticShock) {
					procSpell.Cast(sim, result.Target)
				}
			}
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
			procSpell.Cast(sim, spell.Unit)
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
