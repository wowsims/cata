package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.2s GCD.
const PetGCD = time.Millisecond * 1200

const (
	Unknown PetAbilityType = iota

	// Beast Master Specific Buffs
	AncientHysteria         // Corehound: Burst Haste
	EmbraceOfTheShaleSpider // Shale Spider: Kings
	QirajiFortitude         // Silithid: Stamina

	// Enemy Debuffs
	Gore                // Boar: Phys Dmg Taken ↑
	Ravage              // Ravager: Phys Dmg Taken ↑
	StampedeDebuff      // Rhino: Phys Dmg Taken ↑
	AcidSpitDebuff      // Worm: Phys Dmg Taken ↑
	DemoralizingRoar    // Bear: Enemy DPS ↓
	DemoralizingScreech // Carrion Bird: Enemy DPS ↓
	FireBreathDebuff    // Dragonhawk: Magic Dmg Taken ↑
	LightningBreath     // Wind Serpent: Magic Dmg Taken ↑
	SporeCloud          // Spore Bat: Reduce Cast Speed
	TailSpin            // Fox: Reduce Cast Speed
	Trample             // Goat: Reduce Cast Speed
	LavaBreath          // Corehound: Exotic Cast Speed Debuff
	DustCloud           // Tallstrider: Reduce Armor
	TearArmor           // Raptor: Reduce Armor

	// Utility
	BurrowAttack        // Worm: AoE Damage
	FroststormBreathAoE // Chimera: AoE Damage

	MonstrousBite // Devilsaur: Reduce Healing
	SpiritMend    // Spirit Beast: Healing

	// Basic Attacks
	Bite  // FocusDump: Bite
	Claw  // FocusDump: Claw
	Smack // FocusDump: Smack
)

// These IDs are needed for certain talents.
const BiteSpellID = 17253
const ClawSpellID = 16827
const SmackSpellID = 49966

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {

	case Gore:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 35290, CD: time.Second * 10, School: core.SpellSchoolPhysical, DebuffAura: core.PhysVulnerabilityAura})
	case Ravage:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 50518, CD: time.Second * 10, School: core.SpellSchoolPhysical, DebuffAura: core.PhysVulnerabilityAura})
	case StampedeDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 57386, CD: time.Second * 10, School: core.SpellSchoolPhysical, DebuffAura: core.PhysVulnerabilityAura})
	case AcidSpitDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 55749, CD: time.Second * 10, School: core.SpellSchoolNature, DebuffAura: core.PhysVulnerabilityAura})
	case DemoralizingRoar:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 50256, CD: time.Second * 10, School: core.SpellSchoolNature, DebuffAura: core.WeakenedBlowsAura})
	case DemoralizingScreech:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 24423, CD: time.Second * 8, School: core.SpellSchoolNature, DebuffAura: core.WeakenedBlowsAura})
	case FireBreathDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 34889, CD: time.Second * 30, School: core.SpellSchoolFire, DebuffAura: core.FireBreathDebuff})
	case LightningBreath:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 24844, CD: time.Second * 30, School: core.SpellSchoolFire, DebuffAura: core.LightningBreathDebuff})
	case SporeCloud:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 50274, CD: time.Second * 8, School: core.SpellSchoolFire, DebuffAura: core.SporeCloud})
	case TailSpin:
		return hp.newTailSpin()
	case Trample:
		return hp.newTrample()
	case LavaBreath:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 58604, CD: time.Second * 8, School: core.SpellSchoolFire, DebuffAura: core.LavaBreathAura})
	case DustCloud:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 50285, CD: time.Second * 6, School: core.SpellSchoolNature, DebuffAura: core.WeakenedArmorAura})
	case TearArmor:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 50498, CD: time.Second * 25, School: core.SpellSchoolNature, DebuffAura: core.WeakenedArmorAura})
	case BurrowAttack:
		return hp.newBurrowAttack()
	case FroststormBreathAoE:
		return hp.newFrostStormBreath()
	case MonstrousBite:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 54680, CD: time.Second * 8, School: core.SpellSchoolNature, DebuffAura: core.MortalWoundsAura})
	case SpiritMend:
		return hp.newSpiritMend()

	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case Smack:
		return hp.newSmack()

	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

type PetDebuffSpellConfig struct {
	DebuffAura func(*core.Unit) *core.Aura
	SpellID    int32
	School     core.SpellSchool
	GCD        time.Duration
	CD         time.Duration

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) RegisterKillCommandSpell() *core.Spell {
	actionID := core.ActionID{SpellID: 34026}

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: HunterSpellKillCommand,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 0,
			},
		},
		DamageMultiplier: 1.5,
		CritMultiplier:   hp.CritMultiplier(1.0, 0.0),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.938*spell.RangedAttackPower() + 935
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}

func (hp *HunterPet) newPetDebuff(config PetDebuffSpellConfig) *core.Spell {
	auraArray := hp.NewEnemyAuraArray(config.DebuffAura)
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: config.SpellID},
		SpellSchool: config.School, // Adjust the spell school as needed
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,
		//ClassSpellMask: HunterPetSpellDebuff, // Define or adjust the class spell mask appropriately

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: config.CD,
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				auraArray.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuraArrays: auraArray.ToMap(),
	})
}

func (hp *HunterPet) newFocusDump(pat PetAbilityType, spellID int32) *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: HunterPetFocusDump,
		Flags:          core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Duration: time.Millisecond * 3200,
				Timer:    hp.NewTimer(),
			},
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         core.TernaryFloat64(hp.hunterOwner.Talents.BlinkStrikes, 1.5, 1),
		CritMultiplier:           hp.CritMultiplier(1.0, 0.0),
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hp.hunterOwner.CalcAndRollDamageRange(sim, 0.11400000006, 0.34999999404) + (spell.RangedAttackPower() * 0.168)

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
		flags = core.SpellFlagMeleeMetrics
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
				Duration: config.CD,
			},
		},
		ApplyEffects: applyEffects,
	})
}

func (hp *HunterPet) getFrostStormTickSpell() *core.Spell {
	config := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 92380},
		SpellSchool: core.SpellSchoolNature | core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		BonusCoefficient:         0.288,
		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           hp.DefaultCritMultiplier(),
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := 206 + (spell.MeleeAttackPower() * 0.40)
		spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
	}
	return hp.RegisterSpell(config)
}
func (hp *HunterPet) newFrostStormBreath() *core.Spell {
	frostStormTickSpell := hp.getFrostStormTickSpell()
	hp.frostStormBreath = hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 92380},
		SpellSchool: core.SpellSchoolNature | core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagChanneled | core.SpellFlagNoMetrics,
		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
		},

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		CritMultiplier:           hp.DefaultCritMultiplier(),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "FrostStormBreath-" + hp.Label,
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					frostStormTickSpell.Cast(sim, aoeTarget)
				}
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				frostStormTickSpell.SpellMetrics[target.UnitIndex].Casts += 1
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 206 + (spell.MeleeAttackPower() * 0.40)
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
		},
	})
	return hp.frostStormBreath
}

func (hp *HunterPet) newTailSpin() *core.Spell { panic("newTailSpin not implemented") }
func (hp *HunterPet) newTrample() *core.Spell  { panic("newTrample not implemented") }

func (hp *HunterPet) newBurrowAttack() *core.Spell { panic("newBurrowAttack not implemented") }
func (hp *HunterPet) newSpiritMend() *core.Spell   { panic("newSpiritMend not implemented") }

func (hp *HunterPet) registerRabidCD() {
	hunter := hp.hunterOwner
	if hunter.Options.PetSpec != proto.PetSpec_Ferocity {
		return
	}
	actionID := core.ActionID{SpellID: 53401}

	buffAura := hp.RegisterAura(core.Aura{
		Label:    "Rabid",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hp.MultiplyMeleeSpeed(sim, 1.7)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hp.MultiplyMeleeSpeed(sim, 1/1.7)
		},
	})

	rabidSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 90,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: rabidSpell,
		Type:  core.CooldownTypeDPS,
	})
}
