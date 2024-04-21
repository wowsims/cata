package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (druid *Druid) registerRipSpell() {
	ripBaseNumTicks := int32(8)

	comboPointCoeff := 161.0

	glyphMulti := core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfRip), 1.15, 1.0)

	druid.Rip = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 1079},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		BonusCritRating:  0,
		DamageMultiplier: glyphMulti * druid.RazorClawsMultiplier(druid.GetStat(stats.Mastery)),
		CritMultiplier:   druid.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label: "Rip",
			}),
			NumberOfTicks: ripBaseNumTicks,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				cp := float64(druid.ComboPoints())
				ap := dot.Spell.MeleeAttackPower()

				dot.SnapshotBaseDamage = 56 + comboPointCoeff*cp + 0.0207*cp*ap

				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.NumberOfTicks = ripBaseNumTicks
				dot.Apply(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	druid.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		druid.Rip.DamageMultiplier *= druid.RazorClawsMultiplier(newMastery) / druid.RazorClawsMultiplier(oldMastery)
	})
}

func (druid *Druid) MaxRipTicks() int32 {
	base := int32(8)
	shredGlyphBonus := core.TernaryInt32(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfBloodletting), 3, 0)
	return base + shredGlyphBonus
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.ApplyCostModifiers(druid.Rip.DefaultCast.Cost)
}

func (druid *Druid) ApplyBloodletting(target *core.Unit) {
	ripDot := druid.Rip.Dot(target)

	if ripDot.IsActive() && (ripDot.NumberOfTicks < 11) {
		ripDot.NumberOfTicks += 1
		ripDot.RecomputeAuraDuration()
		ripDot.UpdateExpires(ripDot.ExpiresAt() + time.Second*2)
	}
}
