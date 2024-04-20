package mage

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (mage *Mage) registerFrostfireOrbSpell() {

	mage.FrostfireOrb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 92283},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskEmpty,
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

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.frostfireOrb.EnableWithTimeout(sim, mage.frostfireOrb, time.Second*15)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.FrostfireOrb,
		Type:  core.CooldownTypeDPS,
	})
}

type FrostfireOrb struct {
	core.Pet

	mageOwner *Mage

	FrostfireOrbTick *core.Spell

	TickCount int64
}

func (mage *Mage) NewFrostfireOrb() *FrostfireOrb {
	frostfireOrb := &FrostfireOrb{
		Pet:       core.NewPet("Frostfire Orb", &mage.Character, frostfireOrbBaseStats, createFrostfireOrbInheritance(), false, true),
		mageOwner: mage,
		TickCount: 0,
	}

	mage.AddPet(frostfireOrb)

	return frostfireOrb
}

func (ffo *FrostfireOrb) GetPet() *core.Pet {
	return &ffo.Pet
}

func (ffo *FrostfireOrb) Initialize() {
	ffo.registerFrostfireOrbTickSpell()
}

func (ffo *FrostfireOrb) Reset(_ *core.Simulation) {
}

func (ffo *FrostfireOrb) ExecuteCustomRotation(sim *core.Simulation) {

	spell := ffo.FrostfireOrbTick
	fmt.Println(spell.CanCast(sim, ffo.CurrentTarget))

	if success := spell.Cast(sim, ffo.mageOwner.CurrentTarget); !success {
		ffo.Disable(sim)
	}

}

var frostfireOrbBaseStats = stats.Stats{}

var createFrostfireOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHit:   ownerStats[stats.SpellHit],
			stats.SpellCrit:  ownerStats[stats.SpellCrit],
			stats.SpellPower: ownerStats[stats.SpellPower],
		}
	}
}

func (ffo *FrostfireOrb) registerFrostfireOrbTickSpell() {

	ffo.FrostfireOrbTick = ffo.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 95969},
		SpellSchool:    core.SpellSchoolFrost,
		ProcMask:       core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:          SpellFlagMage | ArcaneMissileSpells | core.SpellFlagNoLogs,
		ClassSpellMask: MageSpellFrostfireOrb,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   ffo.mageOwner.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.134,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.278 * ffo.mageOwner.ScalingBaseDamage
			randomTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]
			spell.CalcAndDealDamage(sim, randomTarget, damage, spell.OutcomeMagicHitAndCrit)
			ffo.TickCount += 1
			if ffo.TickCount == 15 {
				ffo.TickCount = 0
			}
		},
	})
}
