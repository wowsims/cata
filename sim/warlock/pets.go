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
		// based on testing for WotLK Classic the following is true:
		// - pets are meele hit capped if and only if the warlock has 210 (8%) spell hit rating or more
		//   - this is unaffected by suppression and by magic hit debuffs like FF
		// - pets gain expertise from 0% to 6.5% relative to the owners hit, reaching cap at 17% spell hit
		//   - this is also unaffected by suppression and by magic hit debuffs like FF
		//   - this is continious, i.e. not restricted to 0.25 intervals
		// - pets gain spell hit from 0% to 17% relative to the owners hit, reaching cap at 12% spell hit
		// spell hit rating is floor'd
		//   - affected by suppression and ff, but in weird ways:
		// 3/3 suppression => 262 hit  (9.99%) results in misses, 263 (10.03%) no misses
		// 2/3 suppression => 278 hit (10.60%) results in misses, 279 (10.64%) no misses
		// 1/3 suppression => 288 hit (10.98%) results in misses, 289 (11.02%) no misses
		// 0/3 suppression => 314 hit (11.97%) results in misses, 315 (12.01%) no misses
		// 3/3 suppression + FF => 209 hit (7.97%) results in misses, 210 (8.01%) no misses
		// 2/3 suppression + FF => 222 hit (8.46%) results in misses, 223 (8.50%) no misses
		//
		// the best approximation of this behaviour is that we scale the warlock's spell hit by `1/12*17` floor
		// the result and then add the hit percent from suppression/ff

		// does correctly not include ff/misery
		ownerHitChance := ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance

		// TODO: Account for sunfire/soulfrost
		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      ownerStats[stats.SpellPower] * 0.57,
			stats.SpellPower:       ownerStats[stats.SpellPower] * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			//With demonic tactics gone is there any crit inheritance?
			//stats.SpellCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			//stats.MeleeCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeHit: ownerHitChance * core.MeleeHitRatingPerHitChance,
			stats.SpellHit: math.Floor(ownerStats[stats.SpellHit] / 12.0 * 17.0),
			// TODO: revisit
			stats.Expertise: (ownerStats[stats.SpellHit] / core.SpellHitRatingPerHitChance) *
				PetExpertiseScale * core.ExpertisePerQuarterPercentReduction,

			// Resists, 40%
		}
	}
}

func (warlock *Warlock) makePet(summonType proto.WarlockOptions_Summon, baseStats stats.Stats, powerModifier float64,
	attackOptions *core.AutoAttackOptions, statInheritance core.PetStatInheritance) *WarlockPet {

	pet := &WarlockPet{
		Pet: core.NewPet(proto.WarlockOptions_Summon_name[int32(summonType)], &warlock.Character, baseStats,
			statInheritance, summonType == warlock.Options.Summon, false),
	}

	pet.EnableManaBarWithModifier(powerModifier)
	pet.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	pet.AddStat(stats.AttackPower, -20)
	pet.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritRatingPerCritChance*1/52.0833)
	if attackOptions != nil {
		pet.EnableAutoAttacks(pet, *attackOptions)
	}
	warlock.AddPet(pet)
	return pet
}

func (warlock *Warlock) registerPets() {
	autoAttackOptions := &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  88.8,
			BaseDamageMax:  133.3,
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

		stats.Mana:      1559,
		stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}
	impBaseStats := stats.Stats{
		stats.Strength:  429,
		stats.Agility:   309,
		stats.Stamina:   712,
		stats.Intellect: 391,
		stats.Spirit:    395,

		stats.Mana:      1559,
		stats.MeleeCrit: 3.2685 * core.CritRatingPerCritChance,
		stats.SpellCrit: 3.3355 * core.CritRatingPerCritChance,
	}

	inheritance := warlock.MakeStatInheritance()

	warlock.Felhunter = warlock.makePet(proto.WarlockOptions_Felhunter, baseStats, 0.77, autoAttackOptions, inheritance)
	warlock.Felguard = warlock.makePet(proto.WarlockOptions_Felguard, baseStats, 0.77, autoAttackOptions, inheritance)
	warlock.Imp = warlock.makePet(proto.WarlockOptions_Imp, impBaseStats, 1.0, nil, inheritance)
	warlock.Succubus = warlock.makePet(proto.WarlockOptions_Succubus, baseStats, 0.77, autoAttackOptions, inheritance)
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

func (pet *WarlockPet) Initialize() {}

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
			baseDamage := 166 + (1.228 * (0.5 * spell.SpellPower()))

			activeDots := 0

			for _, spell := range pet.Owner.Spellbook {
				if spell.ClassSpellMask&WarlockDoT > 0 && spell.Dot(target).IsActive() {
					activeDots++
				}
			}

			baseDamage *= 1 + 0.15*float64(activeDots)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
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

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.02,
		},
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
			// The .5 seems to be based on the spellpower of the owner. So dividing this by .15 ratio of imp to owner spell power.
			baseDamage := 124 + (0.657 * (0.5 / .15 * spell.SpellPower()))
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
