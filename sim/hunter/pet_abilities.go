package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1200

// Todo: Pet not done
// Apply Wild Hunt
const (
	Unknown PetAbilityType = iota
	AcidSpit
	Bite
	Claw
	DemoralizingScreech
	FireBreath
	FuriousHowl
	FroststormBreath
	Gore
	LavaBreath
	LightningBreath
	MonstrousBite
	NetherShock
	Pin
	PoisonSpit
	Rake
	Ravage
	SavageRend
	ScorpidPoison
	Smack
	Snatch
	SonicBlast
	SpiritStrike
	SporeCloud
	Stampede
	Sting
	Swipe
	TendonRip
	VenomWebSpray
)

// These IDs are needed for certain talents.
const BiteSpellID = 17253
const ClawSpellID = 16827
const SmackSpellID = 52476

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case AcidSpit:
		return hp.newAcidSpit()
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
	case FireBreath:
		return hp.newFireBreath()
	case FroststormBreath:
		return hp.newFroststormBreath()
	case FuriousHowl:
		return hp.newFuriousHowl()
	case Gore:
		return hp.newGore()
	case LavaBreath:
		return hp.newLavaBreath()
	case LightningBreath:
		return hp.newLightningBreath()
	case MonstrousBite:
		return hp.newMonstrousBite()
	case NetherShock:
		return hp.newNetherShock()
	case Pin:
		return hp.newPin()
	case Ravage:
		return hp.newRavage()
	case Smack:
		return hp.newSmack()
	case Snatch:
		return hp.newSnatch()
	case SonicBlast:
		return hp.newSonicBlast()
	case SpiritStrike:
		return hp.newSpiritStrike()
	case SporeCloud:
		return hp.newSporeCloud()
	case Stampede:
		return hp.newStampede()
	case Sting:
		return hp.newSting()
	case Swipe:
		return hp.newSwipe()
	case TendonRip:
		return hp.newTendonRip()
	case VenomWebSpray:
		return hp.newVenomWebSpray()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (hp *HunterPet) newFocusDump(pat PetAbilityType, spellID int32) *core.Spell {

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: HunterPetFocusDump,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Duration: time.Millisecond * 3320,
				Timer:    hp.NewTimer(),
			},
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(132, 188) + (spell.MeleeAttackPower() * 0.2)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newBite() *core.Spell {
	return hp.newFocusDump(Bite, BiteSpellID)
}
func (hp *HunterPet) newClaw() *core.Spell {
	return hp.newFocusDump(Claw, ClawSpellID)
}
func (hp *HunterPet) newSmack() *core.Spell {
	return hp.newFocusDump(Smack, SmackSpellID)
}

type PetSpecialAbilityConfig struct {
	Type    PetAbilityType
	SpellID int32
	School  core.SpellSchool
	GCD     time.Duration
	CD      time.Duration

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) newSpecialAbility(config PetSpecialAbilityConfig) *core.Spell {
	var flags core.SpellFlag
	var applyEffects core.ApplySpellResults
	var procMask core.ProcMask
	onSpellHitDealt := config.OnSpellHitDealt
	if config.School == core.SpellSchoolPhysical {
		flags = core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage
		procMask = core.ProcMaskSpellDamage
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}

		}
	} else {
		procMask = core.ProcMaskMeleeMHSpecial
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	}

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: config.SpellID},
		SpellSchool: config.School,
		ProcMask:    procMask,
		Flags:       flags,

		DamageMultiplier: 1, //* hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: config.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(config.CD),
			},
		},
		ApplyEffects: applyEffects,
	})
}

func (hp *HunterPet) newAcidSpit() *core.Spell {
	acidSpitAuras := hp.NewEnemyAuraArray(core.AcidSpitAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: AcidSpit,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55754,
		School:  core.SpellSchoolNature,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				aura := acidSpitAuras.Get(result.Target)
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newDemoralizingScreech() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.DemoralizingScreechAura)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: DemoralizingScreech,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55487,
		School:  core.SpellSchoolPhysical,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					debuffs.Get(aoeTarget).Activate(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newFireBreath() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: FireBreath,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55485,
		School:  core.SpellSchoolFire,

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newFroststormBreath() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: FroststormBreath,

		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 55492,
		School:  core.SpellSchoolFrost,
	})
}

func (hp *HunterPet) newFuriousHowl() *core.Spell {
	actionID := core.ActionID{SpellID: 24604}

	howlSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 45),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			//petAura.Activate(sim)
			//ownerAura.Activate(sim)
		},
	})

	hp.hunterOwner.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL | core.SpellFlagMCD,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return howlSpell.CanCast(sim, target)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			howlSpell.Cast(sim, target)
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: howlSpell,
		Type:  core.CooldownTypeDPS,
	})

	return nil
}

func (hp *HunterPet) newGore() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: Gore,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 35295,
		School:  core.SpellSchoolPhysical,
	})
}

func (hp *HunterPet) newLavaBreath() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: LavaBreath,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 58611,
		School:  core.SpellSchoolFire,
	})
}

func (hp *HunterPet) newLightningBreath() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: LightningBreath,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 25012,
		School:  core.SpellSchoolNature,
	})
}

func (hp *HunterPet) newMonstrousBite() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: MonstrousBite,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55499,
		School:  core.SpellSchoolPhysical,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {

			}
		},
	})
}

func (hp *HunterPet) newNetherShock() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: NetherShock,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 53589,
		School:  core.SpellSchoolShadow,
	})
}

func (hp *HunterPet) newPin() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53548},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},

		// /DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		},
	})
}

func (hp *HunterPet) newRavage() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Ravage,
		CD:      time.Second * 40,
		SpellID: 53562,
		School:  core.SpellSchoolPhysical,
	})
}

func (hp *HunterPet) newSnatch() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Snatch,
		CD:      time.Second * 60,
		SpellID: 53543,
		School:  core.SpellSchoolPhysical,
	})
}

func (hp *HunterPet) newSonicBlast() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SonicBlast,
		CD:      time.Second * 60,
		SpellID: 53568,
		School:  core.SpellSchoolNature,
	})
}

func (hp *HunterPet) newSpiritStrike() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: SpiritStrike,

		GCD:     0,
		CD:      time.Second * 10,
		SpellID: 61198,
		School:  core.SpellSchoolArcane,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newSporeCloud() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.SporeCloudAura)
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53598},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
			},
		},

		// /DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
			for _, target := range spell.Unit.Env.Encounter.TargetUnits {
				debuffs.Get(target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newStampede() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.StampedeAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Stampede,
		CD:      time.Second * 60,
		SpellID: 57386,
		School:  core.SpellSchoolPhysical,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSting() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.StingAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Sting,
		GCD:     PetGCD,
		CD:      time.Second * 6,
		SpellID: 56631,
		School:  core.SpellSchoolNature,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSwipe() *core.Spell {
	// TODO: This is frontal cone, but might be more realistic as single-target
	// since pets are hard to control.
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Swipe,
		GCD:     PetGCD,
		CD:      time.Second * 5,
		SpellID: 53533,
		School:  core.SpellSchoolPhysical,
	})
}

func (hp *HunterPet) newTendonRip() *core.Spell {
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    TendonRip,
		CD:      time.Second * 20,
		SpellID: 53575,
		School:  core.SpellSchoolPhysical,
	})
}

func (hp *HunterPet) newVenomWebSpray() *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55509},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},

		// /DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		},
	})
}
