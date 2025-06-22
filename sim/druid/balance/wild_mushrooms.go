package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	WildMushroomsBonusCoeff = 0.349
	WildMushroomsCoeff      = 0.295
	WildMushroomsVariance   = 0.19
)

func (moonkin *BalanceDruid) registerWildMushrooms() {

	wildMushroomsStackAura := moonkin.GetOrRegisterAura(core.Aura{
		Label:     "Wild Mushrooms (Tracker)",
		ActionID:  core.ActionID{SpellID: 88747},
		Duration:  core.NeverExpires,
		MaxStacks: 3,
	})

	moonkin.WildMushrooms = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 88747},
		Flags:    core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 11,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wildMushroomsStackAura.Activate(sim)
			wildMushroomsStackAura.AddStack(sim)
		},
	})

	wildMushroomsDamage := moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 78777},
		SpellSchool:      core.SpellSchoolNature,
		Flags:            core.SpellFlagAoE | core.SpellFlagPassiveSpell,
		ProcMask:         core.ProcMaskSpellDamage,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		ClassSpellMask:   druid.DruidSpellWildMushroomDetonate,
		CritMultiplier:   moonkin.DefaultCritMultiplier(),
		BonusCoefficient: WildMushroomsBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				damage := moonkin.CalcAndRollDamageRange(sim, WildMushroomsCoeff, WildMushroomsVariance)
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	moonkin.WildMushroomsDetonate = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88751},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: druid.DruidSpellWildMushroomDetonate,
		Flags:          core.SpellFlagAPL | core.SpellFlagPassiveSpell,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    moonkin.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   moonkin.DefaultCritMultiplier(),
		BonusCoefficient: WildMushroomsBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			for i := wildMushroomsStackAura.GetStacks(); i > 0; i-- {
				wildMushroomsDamage.Cast(sim, target)
				if !spell.ProcMask.Matches(core.ProcMaskSpellProc) {
					wildMushroomsStackAura.RemoveStack(sim)
				}
			}
		},
	})
}
