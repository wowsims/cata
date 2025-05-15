package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	noxiousStingsMultiplier := 1 + 0.05*float64(hunter.Talents.NoxiousStings)

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
			baseDamage := (460 * 5) + 0.40*spell.RangedAttackPower()
			dmg := baseDamage * (float64(hunter.Talents.ImprovedSerpentSting) * 0.15)
			spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMeleeSpecialCritOnly)
		},
	})

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1978},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskProc,
		ClassSpellMask: HunterSpellSerpentSting,
		Flags:          core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{
			Cost: 25,
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
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= noxiousStingsMultiplier

				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if sim.GetRemainingDuration() > time.Millisecond*1 {
						core.StartDelayedAction(sim, core.DelayedActionOptions{
							Priority: core.ActionPriorityDOT - 1,
							DoAt:     sim.CurrentTime + time.Millisecond*1,
							OnAction: func(sim *core.Simulation) {
								hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= noxiousStingsMultiplier
							},
						})
					} else {
						hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= noxiousStingsMultiplier
					}
				},
			},

			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := 460 + 0.08*dot.Spell.RangedAttackPower()
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
					if hunter.Talents.ImprovedSerpentSting > 0 {
						hunter.ImprovedSerpentSting.Cast(sim, target)
					}
				}
				spell.DealOutcome(sim, result)
			})
		},
	})
}
