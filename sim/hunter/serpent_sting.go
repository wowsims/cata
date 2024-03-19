package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	noxiousStingsMultiplier := 1 + 0.05*float64(hunter.Talents.NoxiousStings)

	impSSCritChance := float64(hunter.Talents.ImprovedSerpentSting) * 5
	impSSCritChance += core.TernaryFloat64(hunter.HasSetBonus(ItemSetLightningChargedBattleGear, 2), 5, 0)

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
		BonusCritRating: impSSCritChance + impSSCritChance + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSerpentSting), 6, 0)*core.CritRatingPerCritChance,

		DamageMultiplierAdditive: 1 + 0.15*float64(hunter.Talents.ImprovedSerpentSting),
		// according to in-game testing (which happens to match the wowhead 60% mortal shots flag on wowhead)
		// serpent-sting gets 60% crit modifier instead of 30% crit modifier from mortal shots
		CritMultiplier:   hunter.SpellCritMultiplier(1, float64(hunter.Talents.Toxicology)*0.5),
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
				dot.SnapshotBaseDamage = (460 + 0.40*dot.Spell.RangedAttackPower(target)) / 5
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
