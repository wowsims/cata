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
			stats.Stamina:     ownerStats[stats.Stamina] * 1.0 / 3.0,
			stats.SpellPower:  ownerStats[stats.SpellPower], // All pets inherit spell 1:1
			stats.CritRating:  ownerStats[stats.CritRating],
			stats.HasteRating: ownerStats[stats.HasteRating],

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
			Name:                            name,
			Owner:                           &warlock.Character,
			BaseStats:                       baseStats,
			StatInheritance:                 statInheritance,
			EnabledOnStart:                  enabledOnStart,
			IsGuardian:                      false,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
		}),
	}

	// set pet class for proper scaling values
	pet.Class = pet.Owner.Class
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

	pet.EnableEnergyBar(core.EnergyBarOptions{
		MaxEnergy: 200,
		UnitClass: proto.Class_ClassWarlock,
	})

	warlock.AddPet(petAgent)
}

func (warlock *Warlock) registerPets() {
	warlock.Imp = warlock.registerImp()
	warlock.Succubus = warlock.registerSuccubus()
	warlock.Felhunter = warlock.registerFelHunter()
	warlock.Voidwalker = warlock.registerVoidWalker()
}

func (warlock *Warlock) registerImp() *WarlockPet {
	pet := warlock.registerPet(proto.WarlockOptions_Imp, 0, 0)
	pet.registerFireboltSpell()
	return pet
}

func (warlock *Warlock) registerFelHunter() *WarlockPet {
	pet := warlock.registerPet(proto.WarlockOptions_Felhunter, 2, 3.5)
	pet.registerShadowBiteSpell()
	return pet
}

func (warlock *Warlock) registerVoidWalker() *WarlockPet {
	pet := warlock.registerPet(proto.WarlockOptions_Voidwalker, 2, 3.5)
	pet.registerTormentSpell()
	return pet
}

func (warlock *Warlock) registerSuccubus() *WarlockPet {
	pet := warlock.registerPet(proto.WarlockOptions_Succubus, 3, 1.667)
	pet.registerLashOfPainSpell()
	return pet
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
func (pet *WarlockPet) Initialize() {
}

func (pet *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	waitUntil := time.Duration(1<<63 - 1)

	for _, spell := range pet.AutoCastAbilities {
		if spell.CanCast(sim, pet.CurrentTarget) {
			spell.Cast(sim, pet.CurrentTarget)
			return
		}

		// calculate energy required
		timeTillEnergy := max(0, (spell.Cost.GetCurrentCost()-pet.CurrentEnergy())/pet.EnergyRegenPerSecond())
		waitUntil = min(waitUntil, time.Duration(float64(time.Second)*timeTillEnergy))
	}

	pet.WaitUntil(sim, sim.CurrentTime+waitUntil)
}

func (pet *WarlockPet) registerShadowBiteSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 54049},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellFelHunterShadowBite,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.38,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, pet.CalcScalingSpellDmg(0.38), spell.OutcomeMagicHitAndCrit)
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

func (pet *WarlockPet) registerFireboltSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3110},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellImpFireBolt,
		MissileSpeed:   16,

		EnergyCost: core.EnergyCostOptions{
			Cost: 40,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second * 1,
				CastTime: time.Millisecond * 1750,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.907,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, pet.CalcScalingSpellDmg(0.907), spell.OutcomeMagicHitAndCrit)
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
		EnergyCost: core.EnergyCostOptions{
			Cost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.907,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, pet.CalcScalingSpellDmg(0.907), spell.OutcomeMagicHitAndCrit)
		},
	}))
}

func (pet *WarlockPet) registerTormentSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 3716},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: WarlockSpellVoidwalkerTorment,
		EnergyCost: core.EnergyCostOptions{
			Cost: 50,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.3,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, pet.CalcScalingSpellDmg(0.3), spell.OutcomeMagicHitAndCrit)
		},
	}))
}
