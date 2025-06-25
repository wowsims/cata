package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	AstralStormBonusCoeff = 0.236
	AstralStormCoeff      = 0.199
)

func (moonkin *BalanceDruid) registerAstralStormSpell() {
	moonkin.AstralStormTickSpell = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 106998},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellProc,
		Flags:          core.SpellFlagAoE,
		ClassSpellMask: druid.DruidSpellAstralStorm,

		CritMultiplier:   moonkin.DefaultCritMultiplier(),
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: AstralStormBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			damage := moonkin.CalcScalingSpellDmg(AstralStormCoeff)

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	moonkin.AstralStorm = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 106996},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ClassSpellMask: druid.DruidSpellAstralStorm,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 50.3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Astral Storm (Aura)",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, _ *core.Dot) {
				moonkin.AstralStormTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
