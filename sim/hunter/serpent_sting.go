package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	ImprovedSerpentStingMultiplier := 1.5
	IsSurvival := hunter.Spec == proto.Spec_SpecSurvivalHunter

	hunter.ImprovedSerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:                 core.ActionID{SpellID: 82834},
		SpellSchool:              core.SpellSchoolNature,
		ProcMask:                 core.ProcMaskDirect,
		ClassSpellMask:           HunterSpellSerpentSting,
		Flags:                    core.SpellFlagPassiveSpell,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           hunter.CritMultiplier(1, 0),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (hunter.GetBaseDamageFromCoeff(2.599999905) + 0.1599999964*spell.RangedAttackPower(target)) * 5
			dmg := baseDamage * core.TernaryFloat64(IsSurvival, ImprovedSerpentStingMultiplier, 1)
			spell.CalcAndDealDamage(sim, target, dmg*0.15, spell.OutcomeMeleeSpecialCritOnly)
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
				ActionID: core.ActionID{SpellID: 1978},
				Label:    "SerpentStingDot",
				Tag:      "SerpentSting",
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := hunter.GetBaseDamageFromCoeff(2.6) + 0.16*dot.Spell.RangedAttackPower(target)
				dot.Snapshot(target, baseDmg)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.Dot(target).Apply(sim)
					if IsSurvival {
						hunter.ImprovedSerpentSting.Cast(sim, target)
					}
				}
				spell.DealOutcome(sim, result)
			})
		},
	})
}
