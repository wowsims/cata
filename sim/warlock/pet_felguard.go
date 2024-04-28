package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warlock *Warlock) registerSummonFelguardSpell() {
	warlock.SummonFelguard = warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30146},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSummonFelguard,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.8,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//warlock.ChangeActivePet(sim, warlock.Imp.WarlockPet)
			warlock.ChangeActivePet(sim, PetFelguard)
		},
	})
}

type FelguardPet struct {
	core.Pet

	LegionStrike *core.Spell
	Felstorm     *core.Spell
}

func (warlock *Warlock) NewFelguardPet() *FelguardPet {
	baseStats := stats.Stats{
		stats.Strength:  314,
		stats.Agility:   90,
		stats.Stamina:   328,
		stats.Intellect: 150,
		stats.Spirit:    209,
		stats.Mana:      1559,
		stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}

	autoAttackOptions := &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  88.8,
			BaseDamageMax:  133.3,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}

	felguard := &FelguardPet{
		Pet: core.NewPet(PetFelguard, &warlock.Character, baseStats, warlock.MakeStatInheritance(), false, false),
	}

	felguard.EnableManaBarWithModifier(0.77)
	felguard.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	felguard.AddStat(stats.AttackPower, -20)
	felguard.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)
	felguard.EnableAutoAttacks(felguard, *autoAttackOptions)
	core.ApplyPetConsumeEffects(&warlock.Character, warlock.Consumes)
	warlock.AddPet(felguard)
	return felguard
}

func (felguard *FelguardPet) GetPet() *core.Pet {
	return &felguard.Pet
}

func (felguard *FelguardPet) Reset(_ *core.Simulation) {
}

func (felguard *FelguardPet) Initialize() {
	felguard.registerFelstormSpell()
	felguard.registerLegionStrikeSpell()
}

// TODO: script
func (felguard *FelguardPet) ExecuteCustomRotation(sim *core.Simulation) {
	if felguard.Felstorm.CanCast(sim, felguard.CurrentTarget) {
		felguard.Felstorm.Cast(sim, felguard.CurrentTarget)
	} else if felguard.LegionStrike.CanCast(sim, felguard.CurrentTarget) {
		felguard.LegionStrike.Cast(sim, felguard.CurrentTarget)
	}
}

func (felguard *FelguardPet) registerFelstormSpell() {
	numHits := felguard.Env.GetNumTargets()
	results := make([]*core.SpellResult, numHits)

	felguard.Felstorm = felguard.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 89751},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,
		ClassSpellMask: WarlockSpellFelGuardFelstorm,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.02,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    felguard.NewTimer(),
				Duration: time.Second * 45,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.0,
		CritMultiplier:   2,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Felstorm",
			},
			NumberOfTicks: 6,
			TickLength:    time.Second * 1,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				target := felguard.CurrentTarget
				spell := dot.Spell
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					//TODO: Does this scale with melee attack power as well??
					baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
					// This is the formula from the tooltip but... it multiplies by .5 and then by 2???? so it's just spell power???
					baseDamage += (dot.Spell.SpellPower()*0.50)*2*0.33 + 130
					results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
				curTarget = target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					spell.DealDamage(sim, results[hitIndex])
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)
		},
	})
}

func (felguard *FelguardPet) registerLegionStrikeSpell() {
	numberOfTargets := felguard.Env.GetNumTargets()
	results := make([]*core.SpellResult, numberOfTargets)

	felguard.LegionStrike = felguard.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30213},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL | core.SpellFlagIncludeTargetBonusDamage,
		ClassSpellMask: WarlockSpellFelGuardLegionStrike,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: 1,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    felguard.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			//TODO: Does this scale with melee attack power as well??
			baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			// This is the formula from the tooltip but... it multiplies by .5 and then by 2???? so it's just spell power???
			baseDamage += ((spell.SpellPower()*0.50)*2*0.264 + 139) / float64(numberOfTargets)

			curTarget := target
			for hitIndex := int32(0); hitIndex < numberOfTargets; hitIndex++ {
				results[hitIndex] = spell.CalcDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			curTarget = target
			for hitIndex := int32(0); hitIndex < numberOfTargets; hitIndex++ {
				spell.DealDamage(sim, results[hitIndex])
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}
