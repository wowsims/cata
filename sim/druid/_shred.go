package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (druid *Druid) registerShredSpell() {
	flatDamageBonus := 302.0 / 5.4
	hasBloodletting := druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfBloodletting)
	rendAndTearMod := []float64{1.0, 1.07, 1.13, 1.2}[druid.Talents.RendAndTear]

	druid.Shred = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5221},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: DruidSpellShred,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   40,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return !druid.PseudoStats.InFrontOfTarget && !druid.CannotShredTarget
		},

		DamageMultiplier: 5.4,
		CritMultiplier:   druid.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := flatDamageBonus +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier += .3
			}

			if druid.AssumeBleedActive || druid.Rip.Dot(target).IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
				modifier *= rendAndTearMod
			}
			baseDamage *= modifier

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				if hasBloodletting {
					druid.ApplyBloodletting(target)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := flatDamageBonus + spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())

			modifier := 1.0
			if druid.BleedCategories.Get(target).AnyActive() {
				modifier += .3
			}
			if druid.AssumeBleedActive || druid.Rip.Dot(target).IsActive() || druid.Rake.Dot(target).IsActive() || druid.Lacerate.Dot(target).IsActive() {
				modifier *= rendAndTearMod
			}
			baseDamage *= modifier
			baseres := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)

			attackTable := spell.Unit.AttackTables[target.UnitIndex]
			critChance := spell.PhysicalCritChance(attackTable)
			critMod := (critChance * (spell.CritMultiplier - 1))

			baseres.Damage *= (1 + critMod)

			return baseres
		},
	})
}

func (druid *Druid) CanShred() bool {
	return !druid.PseudoStats.InFrontOfTarget && !druid.CannotShredTarget && druid.CurrentEnergy() >= druid.CurrentShredCost()
}

func (druid *Druid) CurrentShredCost() float64 {
	return druid.Shred.Cost.GetCurrentCost()
}
