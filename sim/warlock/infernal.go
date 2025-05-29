package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var summonInfernalVariance = 0.12
var summonInfernalScale = 1.0
var summonInfernalCoefficient = 1.0

func (warlock *Warlock) registerSummonInfernal(timer *core.Timer) {
	summonInfernalAura := warlock.RegisterAura(core.Aura{
		Label:    "Summon Infernal",
		ActionID: core.ActionID{SpellID: 1122},
		Duration: time.Second * 60,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 1122},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonInfernal,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 25},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: 1500 * time.Millisecond,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: 10 * time.Minute,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warlock.DefaultCritMultiplier(),
		BonusCoefficient: summonInfernalCoefficient,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := warlock.CalcAndRollDamageRange(sim, 0.48500001431, 0.11999999732)
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			warlock.Infernal.EnableWithTimeout(sim, warlock.Infernal, spell.RelatedSelfBuff.Duration)
			// fake aura to show duration
			spell.RelatedSelfBuff.Activate(sim)
		},
		RelatedSelfBuff: summonInfernalAura,
	})
}

type InfernalPet struct {
	core.Pet
	immolationAura *core.Spell
}

func (warlock *Warlock) NewInfernalPet() *InfernalPet {
	baseStats := stats.Stats{
		stats.Health: 55740.8,
		stats.Armor:  19680,
	}

	inheritance := warlock.SimplePetStatInheritanceWithScale(0.25)
	attack := ScaledAutoAttackConfig(2)
	pet := &InfernalPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Infernal",
			Owner:                           &warlock.Character,
			BaseStats:                       baseStats,
			StatInheritance:                 inheritance,
			EnabledOnStart:                  false,
			IsGuardian:                      true,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
	}

	pet.Class = proto.Class_ClassWarlock
	pet.EnableAutoAttacks(pet, *attack)

	warlock.AddPet(pet)
	return pet
}

func (infernal *InfernalPet) GetPet() *core.Pet {
	return &infernal.Pet
}

func (infernal *InfernalPet) Initialize() {
	infernal.immolationAura = infernal.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20153},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAoE,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Immolation",
				ActionID: core.ActionID{SpellID: 19483},
			},
			NumberOfTicks:       31,
			TickLength:          2 * time.Second,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDmg := infernal.CalcScalingSpellDmg(0.1)
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}

func (infernal *InfernalPet) Reset(_ *core.Simulation) {
}

func (infernal *InfernalPet) ExecuteCustomRotation(sim *core.Simulation) {
	infernal.immolationAura.Cast(sim, nil)

	// wait till despawn, we only activate the aura once
	infernal.WaitUntil(sim, sim.CurrentTime+time.Minute)
}
