package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (warlock *Warlock) registerSummonDoomguard(timer *core.Timer) {
	summonDoomguardAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Doomguard",
		ActionID: core.ActionID{SpellID: 18540},
		Duration: time.Duration(45+10*warlock.Talents.AncientGrimoire) * time.Second,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 18540},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonDoomguard,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 80},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{GCD: core.GCDDefault},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: 10 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.Doomguard.EnableWithTimeout(sim, warlock.Doomguard, spell.RelatedSelfBuff.Duration)
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: summonDoomguardAura,
	})
}

type DoomguardPet struct {
	core.Pet

	DoomBolt *core.Spell
}

func (warlock *Warlock) NewDoomguardPet() *DoomguardPet {
	// probably wrong, but nobody is ever going to test this
	baseStats := stats.Stats{
		stats.Strength:            453,
		stats.Agility:             883,
		stats.Stamina:             353,
		stats.Intellect:           159,
		stats.Spirit:              225,
		stats.Mana:                23420,
		stats.PhysicalCritPercent: 0.652,
		stats.SpellCritPercent:    3.3355,
	}

	pet := &DoomguardPet{
		Pet: core.NewPet(core.PetConfig{
			Name:            "Doomguard",
			Owner:           &warlock.Character,
			BaseStats:       baseStats,
			StatInheritance: warlock.petStatInheritance,
			EnabledOnStart:  false,
			IsGuardian:      true,
		}),
	}
	warlock.setPetOptions(pet, 1.0, 0.77, nil)

	return pet
}

func (doomguard *DoomguardPet) GetPet() *core.Pet {
	return &doomguard.Pet
}

func (pet *DoomguardPet) Initialize() {
	pet.registerDoomBolt()
	petMasteryHelper(&pet.Pet)
}

func (pet *DoomguardPet) Reset(_ *core.Simulation) {}

func (pet *DoomguardPet) ExecuteCustomRotation(sim *core.Simulation) {
	if pet.DoomBolt.CanCast(sim, pet.CurrentTarget) {
		pet.DoomBolt.Cast(sim, pet.CurrentTarget)
		minDelay := 150.0
		maxDelay := 750.0
		delayRange := maxDelay - minDelay
		// ~150-750ms delay between casts
		// Research: https://docs.google.com/spreadsheets/d/e/2PACX-1vSaFavGbmrd0l3r7XsPWivap9wMjeaRB6Sl5ieg_GpJ8AfdWkzdG3o2czJ60WHFIZwK0QK5yWF22p8D/pubchart?oid=586881278&format=interactive
		randomDelay := time.Duration(minDelay+delayRange*sim.RandomFloat("Doomguard Cast Delay")) * time.Millisecond
		pet.WaitUntil(sim, pet.NextGCDAt()+randomDelay)
		return
	}
}

func (pet *DoomguardPet) registerDoomBolt() {
	pet.DoomBolt = pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 85692},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellDoomguardDoomBolt,
		MissileSpeed:   20,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 0},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 3000 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,
		BonusCoefficient:         1.28550004959,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// seems weird to have class = warrior here, but seems to fit better then paladin
			baseDmg := core.GetClassSpellScalingCoefficient(proto.Class_ClassWarrior) * 1.05599999428
			baseDmg = sim.Roll(core.ApplyVarianceMinMax(baseDmg, 0.1099999994))
			result := spell.CalcDamage(sim, target, baseDmg, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
