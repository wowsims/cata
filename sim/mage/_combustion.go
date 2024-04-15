package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {

	var combustionDotDamage float64
	var dotsToDuplicate []*core.Dot
	if mage.LivingBomb.Dot.Get(target).IsActive() {
		dotsToDuplicate = append(dotsToDuplicate, mage.LivingBomb.Dot)
	}
	if mage.Pyroblast.Dot.Get(target).IsActive() {
		dotsToDuplicate = append(dotsToDuplicate, mage.Pyroblast.Dot)
	}
	if mage.Ignite.Dot.Get(target).IsActive() {
		dotsToDuplicate = append(dotsToDuplicate, mage.Ignite.Dot)
	}
	for _, Dot := range dotsToDuplicate {
		combustionDotDamage += Dot.SnapshotBaseDamage
	}

	combustionAura := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Combustion",
			Tag: "FireMasteryDot"
			ActionID: actionID,
			Duration: 15 * time.Second,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {

			},
		})
	})

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 11129},
		SpellSchool: core.SpellSchoolFire,
		//ProcMask:    core.SpellFlagNoOnCastComplete,
		Flags: SpellFlagMage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			NumberOfTicks: 10
			TickLength: time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *Simulation, target *Unit, dot *Dot, isRollover bool) {
				dot.SnapshotBaseDamage = combustionDotDamage
				dot.Spell.bonusPeriodicDamageMultiplier = mage.GetFireMasteryBonusMultiplier()
			}
		}
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 1.113,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429*mage.ScalingBaseDamage,
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			spell.DealDamage(sim, result)
			combustionAura.Get(target).Apply(sim)
		},

	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Combustion,
		Type:  core.CooldownTypeDPS,
	})

}
