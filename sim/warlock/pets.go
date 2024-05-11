package warlock

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	AutoCastAbilities []*core.Spell
}

func (warlock *Warlock) MakeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		const petExpertiseScale = 1.53 * core.ExpertisePerQuarterPercentReduction
		const scalingFactor = 0.53153153153153 // TODO: changes from 80 (where it's 0.5) -> 85, clearly there's more to it..

		return stats.Stats{
			stats.Mana:  ownerStats[stats.Mana] / scalingFactor,
			stats.Armor: ownerStats[stats.Armor],

			stats.SpellPower:  ownerStats[stats.SpellPower] * scalingFactor,
			stats.AttackPower: ownerStats[stats.SpellPower] * scalingFactor * 2, // might also simply be pet SP * 2

			// almost certainly wrong, needs more testing
			stats.SpellCrit: ownerStats[stats.SpellCrit],
			stats.MeleeCrit: ownerStats[stats.SpellCrit],

			// pets inherit haste rating directly, evidenced by:
			// 1. haste staying the same if the warlock has windfury totem while the pet doesn't
			// 2. haste staying the same if warlock benefits from wrath of air (pet doesn't get this buff regardless)
			stats.MeleeHaste: ownerStats[stats.SpellHaste],

			// unclear what exactly the scaling is here, but at hit cap they should definitely all be capped
			stats.SpellHit:  ownerStats[stats.SpellHit],
			stats.MeleeHit:  ownerStats[stats.SpellHit],
			stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) * petExpertiseScale,

			// for master demonologist
			stats.Mastery: ownerStats[stats.Mastery],
		}
	}
}

func (warlock *Warlock) makePet(summonType proto.WarlockOptions_Summon, baseStats stats.Stats, meleeMod float64, powerModifier float64,
	attackOptions *core.AutoAttackOptions, statInheritance core.PetStatInheritance) *WarlockPet {

	pet := &WarlockPet{
		Pet: core.NewPet(proto.WarlockOptions_Summon_name[int32(summonType)], &warlock.Character, baseStats,
			statInheritance, summonType == warlock.Options.Summon, false),
	}

	pet.EnableManaBarWithModifier(powerModifier)
	pet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	pet.AddStat(stats.AttackPower, -20)
	pet.AddStatDependency(stats.Agility, stats.MeleeCrit,
		core.CritPerAgiMaxLevel[proto.Class_ClassPaladin]*core.CritRatingPerCritChance)
	pet.AddStatDependency(stats.Intellect, stats.SpellCrit,
		core.CritPerIntMaxLevel[proto.Class_ClassPaladin]*core.CritRatingPerCritChance)
	pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= meleeMod
	if attackOptions != nil {
		pet.EnableAutoAttacks(pet, *attackOptions)
	}
	warlock.AddPet(pet)
	return pet
}

func (warlock *Warlock) registerPets() {
	autoAttackOptions := &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  741.13,
			BaseDamageMax:  1111.62,
			SwingSpeed:     2,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}
	baseStats := stats.Stats{
		stats.Strength:  453,
		stats.Agility:   883,
		stats.Stamina:   353,
		stats.Intellect: 159,
		stats.Spirit:    225,

		stats.Mana: 23420, // x / 0.05 = 1171, x / 0.03 = 702; x ~= 23420

		// determined "base" crit chance with GetCritChanceFrom{Agility,Int}("pet") and then subtracted
		// the calculated amount of crit from agi/int.
		//
		// To calculate crit per agi/int the following two spells where used:
		// Rabies (3150) -5 int
		// Corrupted Agility (6817) -10 agi
		stats.MeleeCrit: 0.652 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}
	impBaseStats := stats.Stats{
		stats.Strength:  429,
		stats.Agility:   309,
		stats.Stamina:   712,
		stats.Intellect: 391,
		stats.Spirit:    395,

		stats.Mana: 17415, // x / 0.16 = 2786, x / 0.02 = 348; x ~= 17415

		stats.MeleeCrit: 0.652 * core.CritRatingPerCritChance,
		stats.SpellCrit: 0.9075 * core.CritRatingPerCritChance,
	}

	inheritance := warlock.MakeStatInheritance()

	warlock.Felhunter = warlock.makePet(proto.WarlockOptions_Felhunter, baseStats, 0.8, 0.77, autoAttackOptions, inheritance)
	warlock.Felguard = warlock.makePet(proto.WarlockOptions_Felguard, baseStats, 1.0, 0.77, autoAttackOptions, inheritance)
	warlock.Imp = warlock.makePet(proto.WarlockOptions_Imp, impBaseStats, 1.0, 1.0, nil, inheritance)
	// TODO: using the modifier for incubus for now, maybe the 1.025 from succubus is the correct one
	warlock.Succubus = warlock.makePet(proto.WarlockOptions_Succubus, baseStats, 1.05, 0.77, autoAttackOptions, inheritance)
}

func (warlock *Warlock) registerPetAbilities() {
	warlock.Felhunter.registerShadowBiteSpell()

	warlock.Felguard.registerLegionStrikeSpell()
	warlock.Felguard.registerFelstormSpell()

	warlock.Imp.registerFireboltSpell()

	warlock.Succubus.registerLashOfPainSpell()
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {}

func (pet *WarlockPet) Initialize() {
	if pet.Owner.Spec == proto.Spec_SpecDemonologyWarlock {
		masteryBonus := func(mastery float64) float64 {
			return 1 + math.Floor(2.3*core.MasteryRatingToMasteryPoints(mastery))/100.0
		}

		pet.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
			pet.PseudoStats.DamageDealtMultiplier /= masteryBonus(oldMastery)
			pet.PseudoStats.DamageDealtMultiplier *= masteryBonus(newMastery)
		})

		// unfortunately mastery is only the *bonus* mastery and not the base value, so we add
		// this manually here such that AddOnMasteryStatChanged() can calculate the correct values
		// while still having 0 mastery = 0% dmg at the start
		pet.AddStats(stats.Stats{stats.Mastery: 8 * core.MasteryRatingPerMasteryPoint})
		pet.PseudoStats.DamageDealtMultiplier *= masteryBonus(8 * core.MasteryRatingPerMasteryPoint)
	}
}

func (pet *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	waitUntil := time.Duration(1<<63 - 1)

	for _, spell := range pet.AutoCastAbilities {
		if spell.CanCast(sim, pet.CurrentTarget) {
			spell.Cast(sim, pet.CurrentTarget)
			return
		}

		waitUntil = min(waitUntil, max(sim.CurrentTime, spell.CD.ReadyAt()))
	}

	pet.WaitUntil(sim, waitUntil)
}

func (pet *WarlockPet) registerShadowBiteSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54049},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellFelHunterShadowBite,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.03,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// shadowbite is a weird spell that seems to get it's SP scaling via a secondary effect,
			// so even though it has variance that only applies to the "base" damage
			// the second "base" value of 182.5 is probably not correct
			baseDamage := 182.5 + pet.Owner.CalcAndRollDamageRange(sim, 0.12600000203, 0.34999999404)
			baseDamage += 1.228 * spell.SpellPower()

			activeDots := 0
			for _, spell := range pet.Owner.Spellbook {
				if spell.ClassSpellMask&WarlockDoT > 0 && spell.Dot(target).IsActive() {
					activeDots++
				}
			}

			baseDamage *= 1 + 0.3*float64(activeDots)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerFelstormSpell() {
	numHits := pet.Env.GetNumTargets()
	results := make([]*core.SpellResult, numHits)

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
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
				Timer:    pet.NewTimer(),
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
				target := pet.CurrentTarget
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
	}))
}

func (pet *WarlockPet) registerLegionStrikeSpell() {
	numberOfTargets := pet.Env.GetNumTargets()
	results := make([]*core.SpellResult, numberOfTargets)

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
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
				Timer:    pet.NewTimer(),
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
	}))
}

func (pet *WarlockPet) registerFireboltSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3110},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,

		ManaCost: core.ManaCostOptions{BaseCost: 0.02},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// seems to function similar to shadowbite, i.e. variance that only applies to the "base" damage, a
			// secondary "base" value of 182.5 (probably not entirely correct) and SP scaling via a secondary effect
			baseDamage := 182.5 + pet.Owner.CalcAndRollDamageRange(sim, 0.1230000034, 0.1099999994)
			baseDamage += 1.228 * spell.SpellPower()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerLashOfPainSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 7814},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellSuccubusLashOfPain,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 187 + (0.612 * (0.5 * spell.SpellPower()))
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}))
}
