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
		Duration: 60 * time.Second,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 18540},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonDoomguard,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 25},
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
	baseStats := stats.Stats{
		stats.Health: 84606.8,
		stats.Armor:  19680,
	}

	pet := &DoomguardPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Doomguard",
			Owner:                           &warlock.Character,
			BaseStats:                       baseStats,
			NonHitExpStatInheritance:        warlock.SimplePetStatInheritanceWithScale(0),
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
			HasResourceRegenInheritance:     false,
		}),
	}

	pet.Class = proto.Class_ClassWarlock
	pet.EnableEnergyBar(core.EnergyBarOptions{
		MaxEnergy: 100,
		UnitClass: proto.Class_ClassWarlock,
	})

	warlock.AddPet(pet)
	return pet
}

func (doomguard *DoomguardPet) GetPet() *core.Pet {
	return &doomguard.Pet
}

func (pet *DoomguardPet) Initialize() {
	pet.Pet.Initialize()
	pet.registerDoomBolt()
}

func (pet *DoomguardPet) Reset(_ *core.Simulation) {}

func (pet *DoomguardPet) ExecuteCustomRotation(sim *core.Simulation) {
	if pet.DoomBolt.CanCast(sim, pet.CurrentTarget) {
		pet.DoomBolt.Cast(sim, pet.CurrentTarget)

		// calculate energy required
		timeTillEnergy := max(0, (pet.DoomBolt.Cost.GetCurrentCost()-pet.CurrentEnergy())/pet.EnergyRegenPerSecond())
		delay := min(0, time.Duration(float64(time.Second)*timeTillEnergy))
		pet.WaitUntil(sim, sim.CurrentTime+delay)
		return
	}
}

func (pet *DoomguardPet) registerDoomBolt() {
	doomBoltExecuteMod := pet.AddDynamicMod(core.SpellModConfig{
		ClassMask:  WarlockSpellDoomguardDoomBolt,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2,
	})

	pet.DoomBolt = pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 85692},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellDoomguardDoomBolt,
		MissileSpeed:   20,

		EnergyCost: core.EnergyCostOptions{
			Cost: 35,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 3000 * time.Millisecond,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.9,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, pet.CalcAndRollDamageRange(sim, 0.9, 0.1), spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	pet.RegisterResetEffect(func(sim *core.Simulation) {
		doomBoltExecuteMod.Deactivate()
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 20 {
				doomBoltExecuteMod.Activate()
			}
		})
	})
}
