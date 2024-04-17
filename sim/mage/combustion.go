package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {

	var combustionDotDamage float64

	/* 	combustionAura := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Combustion",
			Tag:      "FireMasteryDot",
			ActionID: actionID,
			Duration: 15 * time.Second,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {

			},
		})
	}) */

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage,
		ClassSpellMask: MageSpellCombustion,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				combustionDotDamage = 0.0
				// Dots to snapshot onto combustion
				/* 				var dotsToDuplicate []*core.Dot
				   				if mage.LivingBomb.Dot(target).IsActive() {
				   					dotsToDuplicate = append(dotsToDuplicate, mage.LivingBomb.CurDot())
				   				}
				   				if mage.Pyroblast.Dot(target).IsActive() {
				   					dotsToDuplicate = append(dotsToDuplicate, mage.Pyroblast.Dot(target))
				   				}
				   				if mage.Ignite.Dot(target).IsActive() {
				   					dotsToDuplicate = append(dotsToDuplicate, mage.Ignite.Dot(target))
				   				}
				   				for _, Dot := range dotsToDuplicate {
				   					combustionDotDamage += Dot.SnapshotBaseDamage
				   				}
				   				dot.SnapshotBaseDamage = combustionDotDamage
				*/
				var dotSpells []*core.Spell
				dotSpells = append(dotSpells, mage.LivingBomb, mage.Ignite) //, mage.PyroblastDot)
				for _, spell := range dotSpells {
					dots := spell.Dot(mage.CurrentTarget)
					if dots != nil && dots.IsActive() {
						fmt.Println("Snapshot: ", spell.Dot(mage.CurrentTarget).SnapshotBaseDamage)
						combustionDotDamage += spell.Dot(mage.CurrentTarget).SnapshotBaseDamage
					}
				}
				dot.Snapshot(mage.CurrentTarget, combustionDotDamage)
				fmt.Println("Timer: ", sim.CurrentTime)
				fmt.Println("Combustion snapshot: ", combustionDotDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 1.113,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			spell.DealDamage(sim, result)
			spell.Dot(target).Apply(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Combustion,
		Type:  core.CooldownTypeDPS,
	})

}
