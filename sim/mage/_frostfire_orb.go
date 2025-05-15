package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (mage *Mage) registerFrostfireOrbSpell() {
	if mage.Talents.FrostfireOrb == 0 {
		return
	}

	mage.FrostfireOrb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 92283},
		SpellSchool: core.SpellSchoolFrost | core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

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

	FrostfireOrbFingerOfFrost *core.Aura

	TickCount int64
}

func (mage *Mage) NewFrostfireOrb() *FrostfireOrb {
	frostfireOrb := &FrostfireOrb{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Frostfire Orb",
			Owner:           &mage.Character,
			BaseStats:       frostfireOrbBaseStats,
			StatInheritance: createFrostfireOrbInheritance(),
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
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

	if success := spell.Cast(sim, ffo.mageOwner.CurrentTarget); !success {
		ffo.Disable(sim)
	}

}

var frostfireOrbBaseStats = stats.Stats{}

var createFrostfireOrbInheritance = func() func(stats.Stats) stats.Stats {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellHitPercent:  ownerStats[stats.SpellHitPercent],
			stats.SpellCritPercent: ownerStats[stats.SpellCritPercent],
			stats.SpellPower:       ownerStats[stats.SpellPower],
		}
	}
}

func (ffo *FrostfireOrb) registerFrostfireOrbTickSpell() {
	procChance := []float64{0, 0.07, 0.14, 0.20}[ffo.mageOwner.Talents.FingersOfFrost]

	ffo.FrostfireOrbTick = ffo.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 95969},
		SpellSchool:    core.SpellSchoolFrost | core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MageSpellFrostfireOrb,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		DamageMultiplier: 1,

		// should FFO benefit from meta gem?
		CritMultiplier:   ffo.DefaultCritMultiplier(),
		BonusCoefficient: 0.134,
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.278 * ffo.mageOwner.ClassSpellScaling
			randomTarget := sim.Encounter.TargetUnits[int(sim.Roll(0, float64(len(sim.Encounter.TargetUnits))))]
			spell.CalcAndDealDamage(sim, randomTarget, damage, spell.OutcomeMagicHitAndCrit)
			ffo.TickCount += 1
			if ffo.TickCount == 15 {
				ffo.TickCount = 0
			}
		},
	})

	ffo.FrostfireOrbFingerOfFrost = core.MakePermanent(ffo.RegisterAura(core.Aura{
		Label: "Frostfire Orb FoF",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == ffo.FrostfireOrbTick && sim.Proc(procChance, "FingersOfFrostProc") {
				ffo.mageOwner.FingersOfFrostAura.Activate(sim)
				ffo.mageOwner.FingersOfFrostAura.AddStack(sim)
			}
		},
	}))

}
