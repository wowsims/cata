package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}
	var combustionDotDamage float64

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, // application burst might have a proc that might necessitate separating
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

				dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.PyroblastDot}

				var spellDPS float64
				for _, spell := range dotSpells {
					dots := spell.Dot(target)
					// EJ states that combustion double dips on mastery for LB + Pyro, but not ignite
					// https://web.archive.org/web/20120207223126/http://elitistjerks.com/f75/t110187-cataclysm_mage_simulators_formulators/p3/
					if spell != mage.Ignite {
						if dots != nil && dots.IsActive() {
							spellDPS = spell.Dot(target).CalcSnapshotDamage(sim, target, dots.OutcomeTick).Damage / 3
							//fmt.Println(dots.Label, " snapped for : ", spellDPS)
						}
					} else {
						//This part is for ignite. The denominator will probably be variable if it works as intended in cata.
						if dots != nil && dots.IsActive() {
							spellDPS = spell.Dot(target).SnapshotBaseDamage / 2
							//fmt.Println(dots.Label, " snapped for : ", spellDPS)
						}
					}
					combustionDotDamage += spellDPS
				}
				dot.Snapshot(target, combustionDotDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         1.113,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.DealDamage(sim, result)
				spell.Dot(target).Apply(sim)
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.Combustion,
		Type:  core.CooldownTypeDPS,
	})
}
