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

				dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite}
				fmt.Println("Calc Time: ", sim.CurrentTime)
				for _, spell := range dotSpells {
					dots := spell.Dot(mage.CurrentTarget)
					if dots != nil && dots.IsActive() {
						normalizedDPS := 1000000000 * spell.Dot(mage.CurrentTarget).SnapshotBaseDamage / float64(spell.Dot(mage.CurrentTarget).TickPeriod())
						fmt.Println("Snapshot Spell:     ", spell.ActionID)
						fmt.Println("Snapshot Amt:       ", spell.Dot(mage.CurrentTarget).SnapshotBaseDamage)
						fmt.Println("Tick Period(s):     ", (spell.Dot(mage.CurrentTarget).TickPeriod()))
						fmt.Println("Normalized DPS:     ", normalizedDPS)
						combustionDotDamage += normalizedDPS
					}
				}
				dot.Snapshot(mage.CurrentTarget, combustionDotDamage)
				fmt.Println("Combustion base snapshot: ", dot.SnapshotBaseDamage)
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
