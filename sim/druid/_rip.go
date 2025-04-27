package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const RipBaseNumTicks = int32(8)

func (druid *Druid) registerRipSpell() {
	baseDamage := 56.0
	comboPointCoeff := 161.0
	attackPowerCoeff := 0.0207
	glyphMulti := core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfRip), 1.15, 1.0)

	// Blood in the Water refreshes use the CP value from the last "raw" Rip cast, so we need to store that here.
	var comboPointSnapshot int32

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

		BonusCritPercent: 0,
		DamageMultiplier: glyphMulti * druid.RazorClawsMultiplier(druid.GetStat(stats.MasteryRating)),
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		MaxRange:         core.MaxMeleeRange,

		Dot: core.DotConfig{
			Aura: druid.applyRendAndTear(core.Aura{
				Label: "Rip",
			}),
			NumberOfTicks: RipBaseNumTicks,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				if isRollover {
					panic("Rip cannot roll-over snapshots!")
				}

				cp := float64(comboPointSnapshot)
				ap := dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage = baseDamage + comboPointCoeff*cp + attackPowerCoeff*cp*ap
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)

				// Store snapshot power parameters for later use.
				druid.UpdateBleedPower(druid.Rip, sim, target, true, true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHitNoHitCounter)
			if result.Landed() {
				dot := spell.Dot(target)
				dot.BaseTickCount = RipBaseNumTicks
				comboPointSnapshot = druid.ComboPoints()
				dot.Apply(sim)
				druid.ApplyFeral4pT12(sim)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			cp := 5.0 // Hard-code this so that snapshotting calculations can be performed at any CP value.
			ap := spell.MeleeAttackPower()
			baseTickDamage := baseDamage + comboPointCoeff*cp + attackPowerCoeff*cp*ap
			result := spell.CalcPeriodicDamage(sim, target, baseTickDamage, spell.OutcomeExpectedMagicAlwaysHit)
			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := critChance * (spell.CritMultiplier - 1)
			result.Damage *= 1 + critMod
			return result
		},
	})

	druid.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		druid.Rip.DamageMultiplier *= druid.RazorClawsMultiplier(newMastery) / druid.RazorClawsMultiplier(oldMastery)
	})

	druid.Rip.ShortName = "Rip"
}

func (druid *Druid) MaxRipTicks() int32 {
	return RipBaseNumTicks + core.TernaryInt32(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfBloodletting), 3, 0)
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.Cost.GetCurrentCost()
}

func (druid *Druid) ApplyBloodletting(target *core.Unit) {
	ripDot := druid.Rip.Dot(target)

	if ripDot.IsActive() && (ripDot.BaseTickCount < RipBaseNumTicks+3) {
		ripDot.BaseTickCount += 1
		ripDot.UpdateExpires(ripDot.ExpiresAt() + time.Second*2)
	}
}
