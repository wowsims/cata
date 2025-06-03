package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

const RipBaseNumTicks int32 = 8
const RipMaxNumTicks int32 = RipBaseNumTicks + 3

func (druid *Druid) registerRipSpell() {
	// Raw parameters from DB
	const coefficient = 0.10300000012
	const resourceCoefficient = 0.29199999571
	const attackPowerCoeff = 0.0484

	// Scaled parameters for spell code
	baseDamage := coefficient * druid.ClassSpellScaling // 112.7582
	comboPointCoeff := resourceCoefficient * druid.ClassSpellScaling // 319.664

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

		// https://www.wowhead.com/mop-classic/spell=137009/hotfix-passive
		DamageMultiplier: 1.2,

		BonusCritPercent: 0,
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
					return
				}

				cp := float64(druid.ComboPoints())
				ap := dot.Spell.MeleeAttackPower()
				dot.SnapshotPhysical(target, baseDamage + comboPointCoeff*cp + attackPowerCoeff*cp*ap)

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
				dot.Apply(sim)
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

	druid.Rip.ShortName = "Rip"
}

func (druid *Druid) CurrentRipCost() float64 {
	return druid.Rip.Cost.GetCurrentCost()
}

func (druid *Druid) ApplyBloodletting(target *core.Unit) {
	ripDot := druid.Rip.Dot(target)

	if ripDot.IsActive() && (ripDot.BaseTickCount < RipMaxNumTicks) {
		ripDot.BaseTickCount += 1
		ripDot.UpdateExpires(ripDot.ExpiresAt() + ripDot.BaseTickLength)
	}
}
