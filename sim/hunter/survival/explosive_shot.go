package survival

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *SurvivalHunter) registerExplosiveShotSpell() {
	actionID := core.ActionID{SpellID: 53301}
	minFlatDamage := 386.0
	maxFlatDamage := 464.0

	hunter.Hunter.ExplosiveShot = hunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Duration: time.Second * 6,
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfExplosiveShot), 6*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1,
		CritMultiplier: 1,//   hunter.critMultiplier(true, false, false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Explosive Shot - Dot",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(minFlatDamage, maxFlatDamage) + 0.273*dot.Spell.RangedAttackPower(target)
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeRangedHitAndCritSnapshot)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)

			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.TickOnce(sim)
			}
		},
	})
}
