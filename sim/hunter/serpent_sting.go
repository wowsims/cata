package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	noxiousStingsMultiplier := 1 + 0.01*float64(hunter.Talents.NoxiousStings)

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1978},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
		},

		// Need to specially apply LethalShots here, because this spell uses an empty proc mask
		//TOdo: add glyph of serpent sting modifier *6percent crit chance
		BonusCritRating: 1 * core.CritRatingPerCritChance * (float64(hunter.Talents.ImprovedSerpentSting)),

		DamageMultiplierAdditive: 1 +
			0.15*float64(hunter.Talents.ImprovedSerpentSting),
		// according to in-game testing (which happens to match the wowhead 60% mortal shots flag on wowhead)
		// serpent-sting gets 60% crit modifier instead of 30% crit modifier from mortal shots
		CritMultiplier:   hunter.MeleeCritMultiplier(1, 1 + (float64(hunter.Talents.Toxicology)*0.5)),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentSting",
				Tag:   "SerpentSting",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= noxiousStingsMultiplier
					
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= noxiousStingsMultiplier
				},
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 242 + 0.04*dot.Spell.RangedAttackPower(target)
				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
