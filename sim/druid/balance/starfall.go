package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

const (
	StarfallBonusCoeff = 0.364
	StarfallCoeff      = 0.58
	StarfallVariance   = 0.15
)

func (moonkin *BalanceDruid) registerStarfallSpell() {

	numberOfTicks := core.TernaryInt32(moonkin.Env.GetNumTargets() > 1, 20, 10)
	tickLength := time.Second

	starfallTickSpell := moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 50286},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: druid.DruidSpellStarfall,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   moonkin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: StarfallBonusCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := moonkin.CalcAndRollDamageRange(sim, StarfallCoeff, StarfallVariance)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	moonkin.Starfall = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48505},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellProc,
		Flags:       core.SpellFlagAPL,

		RelatedSelfBuff: moonkin.GetOrRegisterAura(core.Aura{
			Label:    "Starfall",
			ActionID: core.ActionID{SpellID: 48505},
			Duration: time.Second * 10,
		}),

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 32.6,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    moonkin.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Starfall",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, _ *core.Dot) {
				starfallTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)

			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})

	moonkin.AddEclipseCallback(func(eclipse Eclipse, gained bool, _ *core.Simulation) {
		if gained && eclipse == LunarEclipse {
			moonkin.Starfall.CD.Reset()
		}
	})
}
