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
	MinEnergy         float64 // The minimum amount of energy needed to the AI casts a spell
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
	proto.WarlockOptions_Felguard: {
		stats.Health: 84606.8,
		stats.Armor:  12568,
	},
}

func (warlock *Warlock) SimplePetStatInheritanceWithScale(apScale float64) core.PetStatInheritance {
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

func ScaledAutoAttackConfig(swingSpeed float64) *core.AutoAttackOptions {
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
	name string,
	enabledOnStart bool,
	baseStats stats.Stats,
	aaOptions *core.AutoAttackOptions,
	statInheritance core.PetStatInheritance,
	isGuardian bool,
) *WarlockPet {
	pet := &WarlockPet{
		Pet: core.NewPet(core.PetConfig{
			Name:                            name,
			Owner:                           &warlock.Character,
			BaseStats:                       baseStats,
			StatInheritance:                 statInheritance,
			EnabledOnStart:                  enabledOnStart,
			IsGuardian:                      isGuardian,
			HasDynamicMeleeSpeedInheritance: true,
			HasDynamicCastSpeedInheritance:  true,
			HasResourceRegenInheritance:     true,
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
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Imp)]
	enabledOnStart := proto.WarlockOptions_Imp == warlock.Options.Summon
	return warlock.registerImpWithName(name, enabledOnStart, false)
}

func (warlock *Warlock) registerImpWithName(name string, enabledOnStart bool, isGuardian bool) *WarlockPet {
	pet := warlock.RegisterPet(proto.WarlockOptions_Imp, 0, 0, name, enabledOnStart, isGuardian)
	pet.registerFireboltSpell()
	return pet
}

func (warlock *Warlock) registerFelHunter() *WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Felhunter)]
	enabledOnStart := proto.WarlockOptions_Felhunter == warlock.Options.Summon
	return warlock.registerFelHunterWithName(name, enabledOnStart, false)
}

func (warlock *Warlock) registerFelHunterWithName(name string, enabledOnStart bool, isGuardian bool) *WarlockPet {
	pet := warlock.RegisterPet(proto.WarlockOptions_Felhunter, 2, 3.5, name, enabledOnStart, isGuardian)
	pet.registerShadowBiteSpell()
	return pet
}

func (warlock *Warlock) registerVoidWalker() *WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Voidwalker)]
	enabledOnStart := proto.WarlockOptions_Voidwalker == warlock.Options.Summon
	return warlock.registerVoidWalkerWithName(name, enabledOnStart, false)
}

func (warlock *Warlock) registerVoidWalkerWithName(name string, enabledOnStart bool, isGuardian bool) *WarlockPet {
	pet := warlock.RegisterPet(proto.WarlockOptions_Voidwalker, 2, 3.5, name, enabledOnStart, isGuardian)
	pet.registerTormentSpell()
	return pet
}

func (warlock *Warlock) registerSuccubus() *WarlockPet {
	name := proto.WarlockOptions_Summon_name[int32(proto.WarlockOptions_Succubus)]
	enabledOnStart := proto.WarlockOptions_Succubus == warlock.Options.Summon
	return warlock.registerSuccubusWithName(name, enabledOnStart, false)
}

func (warlock *Warlock) registerSuccubusWithName(name string, enabledOnStart bool, isGuardian bool) *WarlockPet {
	pet := warlock.RegisterPet(proto.WarlockOptions_Succubus, 3, 1.667, name, enabledOnStart, isGuardian)
	pet.registerLashOfPainSpell()
	return pet
}

func (warlock *Warlock) RegisterPet(
	t proto.WarlockOptions_Summon,
	swingSpeed float64,
	apScale float64,
	name string,
	enabledOnStart bool,
	isGuardian bool,
) *WarlockPet {
	baseStats, ok := petBaseStats[t]
	if !ok {
		panic("Undefined base stats for pet")
	}

	var attackOptions *core.AutoAttackOptions = nil
	if swingSpeed > 0 {
		attackOptions = ScaledAutoAttackConfig(swingSpeed)
	}

	inheritance := warlock.SimplePetStatInheritanceWithScale(apScale)
	return warlock.makePet(name, enabledOnStart, *baseStats, attackOptions, inheritance, isGuardian)
}

func (pet *WarlockPet) GetPet() *core.Pet {
	return &pet.Pet
}

func (pet *WarlockPet) Reset(_ *core.Simulation) {}

func (pet *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	waitUntil := time.Duration(1<<63 - 1)

	for _, spell := range pet.AutoCastAbilities {
		if spell.CanCast(sim, pet.CurrentTarget) && pet.CurrentEnergy() > pet.MinEnergy {
			spell.Cast(sim, pet.CurrentTarget)
			return
		}

		// calculate energy required
		cost := max(pet.MinEnergy, spell.Cost.GetCurrentCost())
		timeTillEnergy := max(0, (cost-pet.CurrentEnergy())/pet.EnergyRegenPerSecond())
		waitUntil = min(waitUntil, time.Duration(float64(time.Second)*timeTillEnergy))
	}

	// for now average the delay out to 100 ms so we don't need to roll random every time
	pet.WaitUntil(sim, sim.CurrentTime+waitUntil+time.Millisecond*100)
}

var petActionShadowBite = core.ActionID{SpellID: 54049}

func (pet *WarlockPet) registerShadowBiteSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       petActionShadowBite,
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

var petActionFireBolt = core.ActionID{SpellID: 3110}

func (pet *WarlockPet) registerFireboltSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       petActionFireBolt,
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

var petActionLashOfPain = core.ActionID{SpellID: 7814}

func (pet *WarlockPet) registerLashOfPainSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       petActionLashOfPain,
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

var petActionTorment = core.ActionID{SpellID: 3716}

func (pet *WarlockPet) registerTormentSpell() {
	pet.AutoCastAbilities = append(pet.AutoCastAbilities, pet.RegisterSpell(core.SpellConfig{
		ActionID:       petActionTorment,
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
