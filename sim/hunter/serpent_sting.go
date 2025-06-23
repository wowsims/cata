package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	IsSurvival := hunter.Spec == proto.Spec_SpecSurvivalHunter
	focusMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 118976})
	hunter.ImprovedSerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 82834},
		SpellSchool:              core.SpellSchoolNature,
		ProcMask:                 core.ProcMaskDirect,
		ClassSpellMask:           HunterSpellSerpentSting,
		Flags:                    core.SpellFlagPassiveSpell | core.SpellFlagRanged,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           hunter.CritMultiplier(1, 0),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (hunter.GetBaseDamageFromCoeff(2.599999905) + 0.1599999964*spell.RangedAttackPower()) * 5
			dmg := baseDamage
			spell.CalcAndDealDamage(sim, target, dmg*0.15, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	gainFocus := hunter.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{OtherID: 1978},
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagNone,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			hunter.AddFocus(sim, 3, focusMetrics)
		},
	})

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1978},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskProc,
		ClassSpellMask: HunterSpellSerpentSting,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       0,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplierAdditive: 1,

		// SS uses Spell Crit which is multiplied by toxicology
		CritMultiplier:   hunter.CritMultiplier(1, 0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentStingDot",
				Tag:   "Serpent Sting",
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := hunter.GetBaseDamageFromCoeff(2.6) + 0.16*dot.Spell.RangedAttackPower()
				dot.Snapshot(target, baseDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if hunter.Spec == proto.Spec_SpecSurvivalHunter {
					gainFocus.Cast(sim, target)
				}
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			if result.Landed() {
				if IsSurvival {
					hunter.ImprovedSerpentSting.Cast(sim, target)
				}
				spell.WaitTravelTime(sim, func(sim *core.Simulation) {
					spell.Dot(target).Apply(sim)
					spell.DealOutcome(sim, result)
				})

			}
		},
	})
}
