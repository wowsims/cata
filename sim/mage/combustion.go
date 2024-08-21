package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerCombustionSpell() {
	if !mage.Talents.Combustion {
		return
	}

	mage.Combustion = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage, // need to check proc mask for impact damage
		ClassSpellMask: MageSpellCombustionApplication,
		Flags:          core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultMageCritMultiplier(),
		BonusCoefficient:         1.113,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.429 * mage.ClassSpellScaling
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.DealDamage(sim, result)
				spell.RelatedDotSpell.Cast(sim, target)
			}
		},
	})

	dotBase := map[int64]float64{
		MageSpellLivingBombDot: 0.25 * mage.ClassSpellScaling,
		MageSpellPyroblastDot:  0.175 * mage.ClassSpellScaling,
	}

	mage.Combustion.RelatedDotSpell = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 11129}.WithTag(1),
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		ClassSpellMask: MageSpellCombustion,
		Flags:          core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultMageCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Combustion Dot",
			},
			NumberOfTicks:       10,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				combustionDotDamage := 0.0
				dotSpells := []*core.Spell{mage.LivingBomb, mage.Ignite, mage.Pyroblast.RelatedDotSpell}
				for _, spell := range dotSpells {
					dot := spell.Dot(target)
					if dot.IsActive() {
						if spell.ClassSpellMask&(MageSpellLivingBombDot|MageSpellPyroblastDot) != 0 {
							dps := dotBase[spell.ClassSpellMask] + dot.BonusCoefficient*dot.Spell.SpellPower()
							dps *= spell.DamageMultiplier * spell.DamageMultiplierAdditive
							combustionDotDamage += dps / dot.BaseTickLength.Seconds()
						} else {
							combustionDotDamage += dot.SnapshotBaseDamage / 2
						}
					}
				}
				dot.Snapshot(target, combustionDotDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}
