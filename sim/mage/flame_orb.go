package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerFlameOrbSpell() {
	//Fire Power gives 33/66/100 percent chance to explode at end like Living Bomb
	orbExplosionChance := float64(mage.Talents.FirePower) / 3

	mage.FlameOrbTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 82739},
		SpellSchool: core.SpellSchoolFire,
		// no idea yet what it procs, likely nothing
		ProcMask:     core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:        SpellFlagMage | core.SpellFlagNoLogs,
		MissileSpeed: 20,

		DamageMultiplier: 1 + .01*float64(mage.Talents.TormentTheWeak),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := (0.278*mage.ScalingBaseDamage + 0.134*spell.SpellPower()) / float64(len(sim.Encounter.TargetUnits))
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	flameOrbExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 123}, //TODO find id, can't check log since no lvl 81 toon
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage,

		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower) +
			.05*float64(mage.Talents.CriticalMass),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.5*mage.ScalingBaseDamage + 0.515*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	mage.FlameOrb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 82731}, //82731 summons orb, 82739 is LIKELY the damaging ID
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.06,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute,
			},
		},

		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower) +
			.05*float64(mage.Talents.CriticalMass),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Flame Orb",
				Duration: time.Second * 15,
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if sim.RandomFloat(aura.Label) < orbExplosionChance {
						flameOrbExplosionSpell.Cast(sim, aura.Unit)
					}
				},
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for i := 0; i < 15; i++ {
				mage.FlameOrbTickSpell.Cast(sim, target)
			}
		},
	})
}
