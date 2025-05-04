package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.2s GCD.
const PetGCD = time.Millisecond * 1200

const (
	Unknown PetAbilityType = iota

	// Raid/Party-Wide Buffs
	RoarOfCourage       // Cat: Increase Mastery
	SpiritBeastBlessing // Spirit Beast: Increase Mastery
	CacklingHowl        // Hyena: Increase Haste
	SerpentsSwiftness   // Serpent: Increase Haste
	BellowingRoar       // Hydra (retired): Increase Crit
	FuriousHowl         // Wolf: Increase Crit
	TerrifyingRoar      // Devilsaur: Increase Crit
	FearlessRoar        // Quilen: Increase Crit
	StillWater          // Water Strider: Spell Power + Crit

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
	EternalGuardian     // Quilen: Battle Res
	SonicBlast          // Bat: Stun
	Sting               // Wasp: Stun
	WebWrap             // Shale Spider: Stun
	ParalyzingQuill     // Porcupine: Incapacitate
	NetherShock         // Nether Ray: Interrupt
	Pummel              // Gorilla: Interrupt
	SerenityDust        // Moth: Interrupt
	Clench              // Scorpid: Disarm
	Snatch              // BirdOfPrey: Disarm
	BadManners          // Monkey: CC
	PetrifyingGaze      // Basilisk: CC
	Lullaby             // Crane: CC
	AnkleCrack          // Crocodile: Slow
	TimeWarp            // Warpstalker: Slow
	FrostBreathSlow     // Chimera: Slow
	Pin                 // Crab: Root
	LockJaw             // Dog: Root
	Web                 // Spider: Root
	VenomWebSpray       // Silithid: Root
	MonstrousBite       // Devilsaur: Reduce Healing
	HornToss            // Rhino: Knockback
	SpiritMend          // Spirit Beast: Healing

	// Hunter-Specific Utility
	ShellShield    // Turtle: Damage Reduction
	HardenCarapace // Beetle: Damage Reduction
	SurfaceTrot    // Water Strider: Water Walking

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
	case RoarOfCourage:
		return hp.newRoarOfCourage()
	case SpiritBeastBlessing:
		return hp.newSpiritBeastBlessing()
	case CacklingHowl:
		return hp.newCacklingHowl()
	case SerpentsSwiftness:
		return hp.newSerpentsSwiftness()
	case BellowingRoar:
		return hp.newBellowingRoar()
	case FuriousHowl:
		return hp.newFuriousHowl()
	case TerrifyingRoar:
		return hp.newTerrifyingRoar()
	case FearlessRoar:
		return hp.newFearlessRoar()
	case StillWater:
		return hp.newStillWater()

	case AncientHysteria:
		return hp.newAncientHysteria()
	case EmbraceOfTheShaleSpider:
		return hp.newEmbraceOfTheShaleSpider()
	case QirajiFortitude:
		return hp.newQirajiFortitude()
	case Gore:
		//return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 35290, CD: time.Second * 25, School: core.SpellSchoolPhysical, DebuffAura: core.GoreAura})
	case Ravage:
		//return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 35290, CD: time.Second * 25, School: core.SpellSchoolPhysical, DebuffAura: core.RavageAura})
	case StampedeDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 35290, CD: time.Second * 10, School: core.SpellSchoolPhysical, DebuffAura: core.StampedeAura})
	case AcidSpitDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 55749, CD: time.Second * 10, School: core.SpellSchoolNature, DebuffAura: core.AcidSpitAura})
	case DemoralizingRoar:
		return hp.newDemoralizingRoar()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
	case FireBreathDebuff:
		return hp.newPetDebuff(PetDebuffSpellConfig{SpellID: 24844, CD: time.Second * 30, School: core.SpellSchoolFire, DebuffAura: core.FireBreathDebuff})
	case LightningBreath:
		return hp.newLightningBreath()
	case SporeCloud:
		return hp.newSporeCloud()
	case TailSpin:
		return hp.newTailSpin()
	case Trample:
		return hp.newTrample()
	case LavaBreath:
		return hp.newLavaBreath()
	case DustCloud:
		return hp.newDustCloud()
	case TearArmor:
		return hp.newTearArmor()

	case BurrowAttack:
		return hp.newBurrowAttack()
	case FroststormBreathAoE:
		return hp.newFrostStormBreath()
	case EternalGuardian:
		return hp.newEternalGuardian()
	case SonicBlast:
		return hp.newSonicBlast()
	case Sting:
		return hp.newSting()
	case WebWrap:
		return hp.newWebWrap()
	case ParalyzingQuill:
		return hp.newParalyzingQuill()
	case NetherShock:
		return hp.newNetherShock()
	case Pummel:
		return hp.newPummel()
	case SerenityDust:
		return hp.newSerenityDust()
	case Clench:
		return hp.newClench()
	case Snatch:
		return hp.newSnatch()
	case BadManners:
		return hp.newBadManners()
	case PetrifyingGaze:
		return hp.newPetrifyingGaze()
	case Lullaby:
		return hp.newLullaby()
	case AnkleCrack:
		return hp.newAnkleCrack()
	case TimeWarp:
		return hp.newTimeWarp()
	case FrostBreathSlow:
		return hp.newFrostBreathSlow()
	case Pin:
		return hp.newPin()
	case LockJaw:
		return hp.newLockJaw()
	case Web:
		return hp.newWeb()
	case VenomWebSpray:
		return hp.newVenomWebSpray()
	case MonstrousBite:
		return hp.newMonstrousBite()
	case HornToss:
		return hp.newHornToss()
	case SpiritMend:
		return hp.newSpiritMend()

	case ShellShield:
		return hp.newShellShield()
	case HardenCarapace:
		return hp.newHardenCarapace()
	case SurfaceTrot:
		return hp.newSurfaceTrot()

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
	return nil
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
			baseDamage := 0.938*spell.RangedAttackPower(target) + 935
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
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

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
		DamageMultiplier:         1,
		CritMultiplier:           hp.CritMultiplier(1.0, 0.0),
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
func (hp *HunterPet) newDemoralizingRoar() *core.Spell {
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
func (hp *HunterPet) newRoarOfCourage() *core.Spell { panic("newRoarOfCourage not implemented") }
func (hp *HunterPet) newSpiritBeastBlessing() *core.Spell {
	panic("newSpiritBeastBlessing not implemented")
}
func (hp *HunterPet) newCacklingHowl() *core.Spell { panic("newCacklingHowl not implemented") }
func (hp *HunterPet) newSerpentsSwiftness() *core.Spell {
	panic("newSerpentsSwiftness not implemented")
}
func (hp *HunterPet) newBellowingRoar() *core.Spell  { panic("newBellowingRoar not implemented") }
func (hp *HunterPet) newFuriousHowl() *core.Spell    { panic("newFuriousHowl not implemented") }
func (hp *HunterPet) newTerrifyingRoar() *core.Spell { panic("newTerrifyingRoar not implemented") }
func (hp *HunterPet) newFearlessRoar() *core.Spell   { panic("newFearlessRoar not implemented") }
func (hp *HunterPet) newStillWater() *core.Spell     { panic("newStillWater not implemented") }

func (hp *HunterPet) newAncientHysteria() *core.Spell { panic("newAncientHysteria not implemented") }
func (hp *HunterPet) newEmbraceOfTheShaleSpider() *core.Spell {
	panic("newEmbraceOfTheShaleSpider not implemented")
}
func (hp *HunterPet) newQirajiFortitude() *core.Spell { panic("newQirajiFortitude not implemented") }

func (hp *HunterPet) newLightningBreath() *core.Spell { panic("newLightningBreath not implemented") }
func (hp *HunterPet) newSporeCloud() *core.Spell      { panic("newSporeCloud not implemented") }
func (hp *HunterPet) newTailSpin() *core.Spell        { panic("newTailSpin not implemented") }
func (hp *HunterPet) newTrample() *core.Spell         { panic("newTrample not implemented") }
func (hp *HunterPet) newLavaBreath() *core.Spell      { panic("newLavaBreath not implemented") }
func (hp *HunterPet) newDustCloud() *core.Spell       { panic("newDustCloud not implemented") }
func (hp *HunterPet) newTearArmor() *core.Spell       { panic("newTearArmor not implemented") }

func (hp *HunterPet) newBurrowAttack() *core.Spell    { panic("newBurrowAttack not implemented") }
func (hp *HunterPet) newEternalGuardian() *core.Spell { panic("newEternalGuardian not implemented") }
func (hp *HunterPet) newSonicBlast() *core.Spell      { panic("newSonicBlast not implemented") }
func (hp *HunterPet) newSting() *core.Spell           { panic("newSting not implemented") }
func (hp *HunterPet) newWebWrap() *core.Spell         { panic("newWebWrap not implemented") }
func (hp *HunterPet) newParalyzingQuill() *core.Spell { panic("newParalyzingQuill not implemented") }
func (hp *HunterPet) newNetherShock() *core.Spell     { panic("newNetherShock not implemented") }
func (hp *HunterPet) newPummel() *core.Spell          { panic("newPummel not implemented") }
func (hp *HunterPet) newSerenityDust() *core.Spell    { panic("newSerenityDust not implemented") }
func (hp *HunterPet) newClench() *core.Spell          { panic("newClench not implemented") }
func (hp *HunterPet) newSnatch() *core.Spell          { panic("newSnatch not implemented") }
func (hp *HunterPet) newBadManners() *core.Spell      { panic("newBadManners not implemented") }
func (hp *HunterPet) newPetrifyingGaze() *core.Spell  { panic("newPetrifyingGaze not implemented") }
func (hp *HunterPet) newLullaby() *core.Spell         { panic("newLullaby not implemented") }
func (hp *HunterPet) newAnkleCrack() *core.Spell      { panic("newAnkleCrack not implemented") }
func (hp *HunterPet) newTimeWarp() *core.Spell        { panic("newTimeWarp not implemented") }
func (hp *HunterPet) newFrostBreathSlow() *core.Spell { panic("newFrostBreathSlow not implemented") }
func (hp *HunterPet) newPin() *core.Spell             { panic("newPin not implemented") }
func (hp *HunterPet) newLockJaw() *core.Spell         { panic("newLockJaw not implemented") }
func (hp *HunterPet) newWeb() *core.Spell             { panic("newWeb not implemented") }
func (hp *HunterPet) newVenomWebSpray() *core.Spell   { panic("newVenomWebSpray not implemented") }
func (hp *HunterPet) newMonstrousBite() *core.Spell   { panic("newMonstrousBite not implemented") }
func (hp *HunterPet) newHornToss() *core.Spell        { panic("newHornToss not implemented") }
func (hp *HunterPet) newSpiritMend() *core.Spell      { panic("newSpiritMend not implemented") }

func (hp *HunterPet) newShellShield() *core.Spell    { panic("newShellShield not implemented") }
func (hp *HunterPet) newHardenCarapace() *core.Spell { panic("newHardenCarapace not implemented") }
func (hp *HunterPet) newSurfaceTrot() *core.Spell    { panic("newSurfaceTrot not implemented") }
