package warlock

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	AutoCastAbilities []*core.Spell
}

func (warlock *Warlock) petStatInheritance(ownerStats stats.Stats) stats.Stats {
	const petExpertiseScale = 1.53 * core.ExpertisePerQuarterPercentReduction
	const scalingFactor = 0.53153153153153 // TODO: changes from 80 (where it's 0.5) -> 85, clearly there's more to it..

	return stats.Stats{
		stats.Mana:  ownerStats[stats.Mana] / scalingFactor,
		stats.Armor: ownerStats[stats.Armor],

		stats.SpellPower:  ownerStats[stats.SpellPower] * scalingFactor,
		stats.AttackPower: ownerStats[stats.SpellPower] * scalingFactor * 2, // might also simply be pet SP * 2

		// almost certainly wrong, needs more testing
		stats.SpellCritPercent:    ownerStats[stats.SpellCritPercent],
		stats.PhysicalCritPercent: ownerStats[stats.SpellCritPercent],

		// pets inherit haste rating directly, evidenced by:
		// 1. haste staying the same if the warlock has windfury totem while the pet doesn't
		// 2. haste staying the same if warlock benefits from wrath of air (pet doesn't get this buff regardless)
		stats.HasteRating: ownerStats[stats.HasteRating],

		// unclear what exactly the scaling is here, but at hit cap they should definitely all be capped
		stats.HitRating:       ownerStats[stats.SpellHitPercent] * core.SpellHitRatingPerHitPercent,
		stats.ExpertiseRating: ownerStats[stats.SpellHitPercent] * petExpertiseScale,

		// for master demonologist
		stats.MasteryRating: ownerStats[stats.MasteryRating],
	}
}

func (warlock *Warlock) makePet(summonType proto.WarlockOptions_Summon, baseStats stats.Stats, meleeMod float64,
	powerModifier float64, aaOptions *core.AutoAttackOptions, statInheritance core.PetStatInheritance) *WarlockPet {

	name := proto.WarlockOptions_Summon_name[int32(summonType)]
	enabledOnStart := summonType == warlock.Options.Summon
	pet := &WarlockPet{
		Pet: core.NewPet(core.PetConfig{
			Name:            name,
			Owner:           &warlock.Character,
			BaseStats:       baseStats,
			StatInheritance: statInheritance,
			EnabledOnStart:  enabledOnStart,
			IsGuardian:      false,
		}),
	}
	if enabledOnStart {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			warlock.ActivePet = pet
		})
	}

	warlock.setPetOptions(pet, meleeMod, powerModifier, aaOptions)

	return pet
}

func (warlock *Warlock) setPetOptions(petAgent core.PetAgent, meleeMod float64, powerModifier float64,
	aaOptions *core.AutoAttackOptions) {

	pet := petAgent.GetPet()
	pet.EnableManaBarWithModifier(powerModifier)
	pet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	pet.AddStat(stats.AttackPower, -20)
	pet.AddStatDependency(stats.Agility, stats.PhysicalCritPercent,
		core.CritPerAgiMaxLevel[proto.Class_ClassPaladin])
	pet.AddStatDependency(stats.Intellect, stats.SpellCritPercent,
		core.CritPerIntMaxLevel[proto.Class_ClassPaladin])
	pet.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexPhysical] *= meleeMod
	if aaOptions != nil {
		pet.EnableAutoAttacks(petAgent, *aaOptions)
	}
	warlock.AddPet(petAgent)
}

func (warlock *Warlock) registerPets() {
	aaOptions := &core.AutoAttackOptions{
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
		stats.PhysicalCritPercent: 0.652,
		stats.SpellCritPercent:    3.3355,
	}
	impBaseStats := stats.Stats{
		stats.Strength:  429,
		stats.Agility:   309,
		stats.Stamina:   712,
		stats.Intellect: 391,
		stats.Spirit:    395,

		stats.Mana: 17415, // x / 0.16 = 2786, x / 0.02 = 348; x ~= 17415

		stats.PhysicalCritPercent: 0.652,
		stats.SpellCritPercent:    0.9075,
	}

	inheritance := warlock.petStatInheritance
	warlock.Felhunter = warlock.makePet(proto.WarlockOptions_Felhunter, baseStats, 0.8, 0.77, aaOptions, inheritance)
	warlock.Felguard = warlock.makePet(proto.WarlockOptions_Felguard, baseStats, 1.0, 0.77, aaOptions, inheritance)
	warlock.Imp = warlock.makePet(proto.WarlockOptions_Imp, impBaseStats, 1.0, 1.0, nil, inheritance)
	warlock.Succubus = warlock.makePet(proto.WarlockOptions_Succubus, baseStats, 1.025, 0.77, aaOptions, inheritance)
}

func (warlock *Warlock) registerPetAbilities() {
	warlock.Felhunter.registerShadowBiteSpell()

	warlock.Felguard.registerLegionStrikeSpell()
	warlock.Felguard.registerFelstormSpell()

	warlock.Imp.registerFireboltSpell(warlock)

	warlock.Succubus.registerLashOfPainSpell()
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {}

func petMasteryHelper(pet *core.Pet) {
	if pet.Owner.Spec == proto.Spec_SpecDemonologyWarlock {
		// The current code convention is to not include base Mastery points in the MasteryRating stat
		// value, only bonus Rating gained from gear / consumes. Therefore, we bake in the base points (8
		// for Warlock but not for all classes) to the damage multiplier calculation.
		petDamageMultiplier := func(masteryRating float64) float64 {
			return 1 + math.Floor(2.3*(8+core.MasteryRatingToMasteryPoints(masteryRating)))/100
		}

		// Set initial multiplier from base stats (should be 0 Mastery Rating at this point since
		// owner stats have not yet been inherited).
		pet.PseudoStats.DamageDealtMultiplier *= petDamageMultiplier(pet.GetStat(stats.MasteryRating))

		// Keep the multiplier updated when Mastery Rating changes.
		pet.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
			pet.PseudoStats.DamageDealtMultiplier *= petDamageMultiplier(newMasteryRating) / petDamageMultiplier(oldMasteryRating)
		})
	}
}

func (pet *WarlockPet) Initialize() {
	petMasteryHelper(&pet.Pet)
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
			BaseCostPercent: 3,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 6 * time.Second,
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
				// spell.RelatedDotSpell == nil to not double count spells with a separate dot component spell
				if spell.Dot(target) != nil && spell.RelatedDotSpell == nil && spell.Dot(target).IsActive() {
					activeDots++
				}
			}

			baseDamage *= 1 + 0.3*float64(activeDots)
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerFelstormSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 89751},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellFelGuardFelstorm,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 2},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 45 * time.Second,
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
			TickLength:    1 * time.Second,
			OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
				spell := dot.Spell
				baseDmg := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
				baseDmg += pet.Owner.CalcScalingSpellDmg(0.1155000031) + 0.231*spell.MeleeAttackPower()

				for _, target := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, target, baseDmg, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
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

	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30213},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellFelGuardLegionStrike,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 6},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 6 * time.Second,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			baseDmg += pet.Owner.CalcScalingSpellDmg(0.1439999938) + 0.264*spell.MeleeAttackPower()
			baseDmg /= float64(numberOfTargets)

			for _, target := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, target, baseDmg, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
			}
		},
	}))
}

func (pet *WarlockPet) registerFireboltSpell(warlock *Warlock) {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3110},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,
		MissileSpeed:   16,

		ManaCost: core.ManaCostOptions{BaseCostPercent: 2},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2500 * time.Millisecond,
			},
			IgnoreHaste: true,
			// Custom modify cast to not lower GCD
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.Unit.ApplyCastSpeedForSpell(spell.DefaultCast.CastTime, spell)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// seems to function similar to shadowbite, i.e. variance that only applies to the "base" damage, a
			// secondary "base" value of 182.5 (probably not entirely correct) and SP scaling via a secondary effect
			baseDamage := 182.5 + pet.Owner.CalcAndRollDamageRange(sim, 0.1230000034, 0.1099999994)
			baseDamage += 0.657 * spell.SpellPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if warlock.Talents.BurningEmbers > 0 && result.Landed() {
				dot := warlock.BurningEmbers.Dot(result.Target)
				if !dot.IsActive() {
					dot.SnapshotBaseDamage = 0.0 // ensure we don't use old dot data
				}
				dot.SnapshotBaseDamage += result.Damage * 0.25 * float64(warlock.Talents.BurningEmbers)
				dot.Apply(sim)
			}
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
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
			BaseCostPercent: 3,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    pet.NewTimer(),
				Duration: 12 * time.Second,
			},
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
