package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (mage *Mage) registerFlameOrbSpell() {

	// Frostfire Orb talent converts Flame Orb spell into Frostfire Orb spell.
	// Don't allow a user access to both at the same time.
	if mage.Talents.FrostfireOrb != 0 {
		return
	}

	flameOrb := mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 82731},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty, //tbd
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellFlameOrb,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 6,
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
		Spell: flameOrb,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerFlameOrbExplodeSpell() {

	mage.FlameOrbExplode = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 83619},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFlameOrb,
		Flags:          core.SpellFlagAoE,

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultCritMultiplier(),
		BonusCoefficient: 0.193,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 1.318 * mage.ClassSpellScaling

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)
			}

		},
	})
}

type FlameOrb struct {
	core.Pet

	mageOwner *Mage

	FlameOrbTick *core.Spell

	TickCount int64
}

func (mage *Mage) NewFlameOrb() *FlameOrb {
	flameOrb := &FlameOrb{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Flame Orb",
			Owner:           &mage.Character,
			BaseStats:       flameOrbBaseStats,
			StatInheritance: createFlameOrbInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
		mageOwner: mage,
		TickCount: 0,
	}

	flameOrb.Pet.OnPetEnable = flameOrb.enable

	mage.AddPet(flameOrb)

	return flameOrb
}

func (fo *FlameOrb) enable(sim *core.Simulation) {

	fo.PseudoStats.DamageDealtMultiplier = fo.Owner.PseudoStats.DamageDealtMultiplier
	fo.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] = fo.Owner.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire]

	fo.EnableDynamicStats(func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	})
}

func (fo *FlameOrb) GetPet() *core.Pet {
	return &fo.Pet
}

func (fo *FlameOrb) Initialize() {
	fo.registerFlameOrbTickSpell()
}

func (fo *FlameOrb) Reset(_ *core.Simulation) {
	fo.TickCount = 0
}

func (fo *FlameOrb) ExecuteCustomRotation(sim *core.Simulation) {
	spell := fo.FlameOrbTick
	if success := spell.Cast(sim, fo.CurrentTarget); !success {
		fo.Disable(sim)
	}
}

var flameOrbBaseStats = stats.Stats{}

var createFlameOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			stats.SpellPower:       ownerStats[stats.SpellPower],
		}
	}
}

func (fo *FlameOrb) registerFlameOrbTickSpell() {
	fo.FlameOrbTick = fo.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 82739},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFlameOrb,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   fo.mageOwner.DefaultCritMultiplier(),
		BonusCoefficient: 0.134,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			damage := fo.mageOwner.CalcAndRollDamageRange(sim, 0.278, 0.25)
			randomTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]
			spell.CalcAndDealDamage(sim, randomTarget, damage, spell.OutcomeMagicHitAndCrit)

			fo.TickCount += 1
			if fo.TickCount == 15 {
				procChance := []float64{0.0, 0.33, 0.66, 1.0}[fo.mageOwner.Talents.FirePower]
				if sim.Proc(procChance, "FlameOrbExplosion") {
					fo.mageOwner.FlameOrbExplode.Cast(sim, fo.mageOwner.CurrentTarget)
				}
				fo.TickCount = 0
			}
		},
	})
}
