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

var petBaseStats = map[proto.WarlockOptions_Summon]*stats.Stats{
	proto.WarlockOptions_Imp: {
		stats.Health: 48312.8,
		stats.Armor:  19680,
	},
	proto.WarlockOptions_Voidwalker: {
		stats.Health: 120900.8,
		stats.Armor:  19680,
	},
	proto.WarlockOptions_Succubus: {
		stats.Health: 84606.8,
		stats.Armor:  12568,
	},
	proto.WarlockOptions_Felhunter: {
		stats.Health: 84606.8,
		stats.Armor:  19680,
	},
}

func (warlock *Warlock) simplePetStatInheritanceWithScale(apScale float64) core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.SpellPower: ownerStats[stats.SpellPower], // All pets inherit spell 1:1
			stats.CritRating: ownerStats[stats.CritRating],

			// unclear what exactly the scaling is here, but at hit cap they should definitely all be capped
			stats.HitRating:       ownerStats[stats.HitRating],
			stats.ExpertiseRating: ownerStats[stats.HitRating] / core.SpellHitRatingPerHitPercent * core.ExpertisePerQuarterPercentReduction * 4, // 1% hit = 1% expertise

			stats.AttackPower: ownerStats[stats.SpellPower] * apScale,
		}
	}
}

func scaledAutoAttackConfig(swingSpeed float64) *core.AutoAttackOptions {
	return &core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  math.Floor(core.ClassBaseScaling[proto.Class_ClassWarlock]),
			BaseDamageMax:  math.Ceil(core.ClassBaseScaling[proto.Class_ClassWarlock]),
			SwingSpeed:     swingSpeed,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	}
}

func (warlock *Warlock) makePet(
	summonType proto.WarlockOptions_Summon,
	baseStats stats.Stats,
	aaOptions *core.AutoAttackOptions,
	statInheritance core.PetStatInheritance,
) *WarlockPet {
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

	warlock.setPetOptions(pet, aaOptions)

	return pet
}

func (warlock *Warlock) setPetOptions(petAgent core.PetAgent, aaOptions *core.AutoAttackOptions) {
	pet := petAgent.GetPet()
	if aaOptions != nil {
		pet.EnableAutoAttacks(petAgent, *aaOptions)
	}
	warlock.AddPet(petAgent)
}

func (warlock *Warlock) registerPets() {
	warlock.Imp = warlock.registerImp()
	warlock.Succubus = warlock.registerSuccubus()
	warlock.Felhunter = warlock.registerFelHunter()
	warlock.Voidwalker = warlock.registerVoidWalker()
}

func (warlock *Warlock) registerImp() *WarlockPet {
	return warlock.registerPet(proto.WarlockOptions_Imp, 0, 0)
}

func (warlock *Warlock) registerFelHunter() *WarlockPet {
	return warlock.registerPet(proto.WarlockOptions_Felhunter, 2, 3.5)
}

func (warlock *Warlock) registerVoidWalker() *WarlockPet {
	return warlock.registerPet(proto.WarlockOptions_Voidwalker, 2, 3.5)
}

func (warlock *Warlock) registerSuccubus() *WarlockPet {
	return warlock.registerPet(proto.WarlockOptions_Succubus, 3, 1.667)
}

func (warlock *Warlock) registerPet(t proto.WarlockOptions_Summon, swingSpeed float64, apScale float64) *WarlockPet {
	baseStats, ok := petBaseStats[t]
	if !ok {
		panic("Undefined base stats for pet")
	}

	var attackOptions *core.AutoAttackOptions = nil
	if swingSpeed > 0 {
		attackOptions = scaledAutoAttackConfig(swingSpeed)
	}

	inheritance := warlock.simplePetStatInheritanceWithScale(apScale)
	return warlock.makePet(t, *baseStats, attackOptions, inheritance)
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {}

// func petMasteryHelper(pet *core.Pet) {
// 	if pet.Owner.Spec == proto.Spec_SpecDemonologyWarlock {
// 		// The current code convention is to not include base Mastery points in the MasteryRating stat
// 		// value, only bonus Rating gained from gear / consumes. Therefore, we bake in the base points (8
// 		// for Warlock but not for all classes) to the damage multiplier calculation.
// 		petDamageMultiplier := func(masteryRating float64) float64 {
// 			return 1 + math.Floor(2.3*(8+core.MasteryRatingToMasteryPoints(masteryRating)))/100
// 		}

// 		// Set initial multiplier from base stats (should be 0 Mastery Rating at this point since
// 		// owner stats have not yet been inherited).
// 		pet.PseudoStats.DamageDealtMultiplier *= petDamageMultiplier(pet.GetStat(stats.MasteryRating))

// 		// Keep the multiplier updated when Mastery Rating changes.
// 		pet.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMasteryRating float64, newMasteryRating float64) {
// 			pet.PseudoStats.DamageDealtMultiplier *= petDamageMultiplier(newMasteryRating) / petDamageMultiplier(oldMasteryRating)
// 		})
// 	}
// }

func (pet *WarlockPet) Initialize() {
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
