package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (mage *Mage) registerFlameOrbSpell() {

	mage.FlameOrb = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 82731},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, //tbd
		Flags:          SpellFlagMage | core.SpellFlagAPL,
		ClassSpellMask: MageSpellFlameOrb,

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

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.flameOrb.EnableWithTimeout(sim, mage.flameOrb, time.Second*15)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.FlameOrb,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerFlameOrbExplodeSpell() {

	mage.FlameOrbExplode = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 83619},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage | core.SpellFlagNoLogs,
		ClassSpellMask: MageSpellFlameOrb,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.193,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO implement proc chance when talents get fixed
			// procChance := []float64{0.0, 0.33, 0.66, 1.0}[mage.Talents.FirePower]

			// Debugging proc chance
			/* 			fmt.Println("---")
			   			fmt.Println(procChance)
			   			fmt.Println("Ignite Talents: ", mage.Talents.Ignite)
			   			fmt.Println("Fire Power Talents: ", mage.Talents.FirePower)
			   			fmt.Println("Blazing Speed Talents: ", mage.Talents.BlazingSpeed)
			   			fmt.Println("Impact Talents: ", mage.Talents.Impact) */

			damage := 1.318 * mage.ScalingBaseDamage
			// if sim.Proc(procChance, "FlameOrbExplosion") {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
				spell.SpellMetrics[target.UnitIndex].Hits++
			}
			// }
		},
	})
}

type FlameOrb struct {
	core.Pet

	mageOwner *Mage

	FlameOrbTick    *core.Spell
	FlameOrbExplode *core.Spell
	TickCount       int64
	TalentPoints    float64
}

func (mage *Mage) NewFlameOrb() *FlameOrb {
	flameOrb := &FlameOrb{
		Pet:          core.NewPet("Flame Orb", &mage.Character, flameOrbBaseStats, createFlameOrbInheritance(), false, true),
		mageOwner:    mage,
		TickCount:    0,
		TalentPoints: float64(mage.Talents.FirePower),
	}

	mage.AddPet(flameOrb)

	return flameOrb
}

func (fo *FlameOrb) GetPet() *core.Pet {
	return &fo.Pet
}

func (fo *FlameOrb) Initialize() {
	fo.registerFlameOrbTickSpell()
}

func (fo *FlameOrb) Reset(_ *core.Simulation) {
}

func (fo *FlameOrb) ExecuteCustomRotation(sim *core.Simulation) {

	// develop something where on timer expire, cast FlameOrbExplode
	spell := fo.FlameOrbTick
	if success := spell.Cast(sim, fo.CurrentTarget); !success {
		fo.Disable(sim)
	}
}

var flameOrbBaseStats = stats.Stats{}

var createFlameOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHit:   ownerStats[stats.SpellHit],
			stats.SpellCrit:  ownerStats[stats.SpellCrit],
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	}
}

func (fo *FlameOrb) registerFlameOrbTickSpell() {
	curTarget := fo.CurrentTarget

	fo.FlameOrbTick = fo.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 82739},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage | core.SpellFlagNoLogs,
		ClassSpellMask: MageSpellFlameOrb,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fo.mageOwner.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.134,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.278 * fo.mageOwner.ScalingBaseDamage
			spell.CalcAndDealDamage(sim, curTarget, damage, spell.OutcomeMagicHitAndCrit)
			curTarget = sim.Environment.NextTargetUnit(curTarget)

			fo.TickCount += 1
			//spell.SpellMetrics[target.UnitIndex].Casts--
			if fo.TickCount == 15 {
				fo.mageOwner.FlameOrbExplode.Cast(sim, fo.mageOwner.CurrentTarget)
				fo.TickCount = 0
			}
		},
	})
}
