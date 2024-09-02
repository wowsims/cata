package core

import (
	"cmp"
	"slices"
	"strconv"
	"time"

	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type BuffConfig struct {
	Label    string
	ActionID ActionID
	Stats    []StatConfig
}

type StatConfig struct {
	Stat             stats.Stat
	Amount           float64
	IsMultiplicative bool
}

func makeExclusiveMultiplierBuff(aura *Aura, stat stats.Stat, value float64) {
	dep := aura.Unit.NewDynamicMultiplyStat(stat, value)
	aura.NewExclusiveEffect(stat.StatName()+"%Buff", false, ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *ExclusiveEffect, s *Simulation) {
			if ee.Aura.Unit.Env.MeasuringStats && ee.Aura.Unit.Env.State != Finalized {
				aura.Unit.StatDependencyManager.EnableDynamicStatDep(dep)
			} else {
				ee.Aura.Unit.EnableDynamicStatDep(s, dep)
			}
		},
		OnExpire: func(ee *ExclusiveEffect, s *Simulation) {
			if ee.Aura.Unit.Env.MeasuringStats {
				aura.Unit.StatDependencyManager.DisableDynamicStatDep(dep)
			} else {
				ee.Aura.Unit.DisableDynamicStatDep(s, dep)
			}
		},
	})
}

func makeExclusiveFlatStatBuff(aura *Aura, stat stats.Stat, value float64) {
	aura.NewExclusiveEffect(stat.StatName()+"Buff", false, ExclusiveEffect{
		Priority: value,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stat, value)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.AddStatDynamic(sim, stat, -value)
		},
	})
}

func registerExlusiveEffects(aura *Aura, config []StatConfig) {
	for _, statConfig := range config {
		if statConfig.IsMultiplicative {
			makeExclusiveMultiplierBuff(aura, statConfig.Stat, statConfig.Amount)
		} else {
			makeExclusiveFlatStatBuff(aura, statConfig.Stat, statConfig.Amount)
		}
	}
}

func makeExclusiveAllStatPercentBuff(unit *Unit, label string, actionID ActionID, value float64) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		label,
		actionID,
		[]StatConfig{
			{stats.Agility, value, true},
			{stats.Strength, value, true},
			{stats.Stamina, value, true},
			{stats.Intellect, value, true},
		}})
}

func makeExclusiveBuff(unit *Unit, config BuffConfig) *Aura {
	if config.Label == "" {
		panic("Buff without label.")
	}

	if ActionID.IsEmptyAction(config.ActionID) {
		panic("Buff without ActionID")
	}

	baseAura := MakePermanent(unit.GetOrRegisterAura(Aura{
		Label:      config.Label,
		ActionID:   config.ActionID,
		BuildPhase: CharacterBuildPhaseBuffs,
	}))

	registerExlusiveEffects(baseAura, config.Stats)
	return baseAura
}

// Applies buffs that affect individual players.
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, _ *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	character := agent.GetCharacter()

	// % Stats Buffs
	// https://www.wowhead.com/cata/spell=1126/mark-of-the-wild
	// https://www.wowhead.com/cata/spell=20217/blessing-of-kings
	// https://www.wowhead.com/cata/item=63140/drums-of-the-burning-wild
	if raidBuffs.BlessingOfKings {
		BlessingOfKingsAura(&character.Unit)
	}

	if raidBuffs.DrumsOfTheBurningWild {
		DrumsOfTheBurningWildAura(&character.Unit)
	}

	if raidBuffs.MarkOfTheWild {
		MarkOfTheWildAura(&character.Unit)
	}

	// Resistances
	if raidBuffs.ResistanceAura {
		ResistanceAura(&character.Unit)
	}

	if raidBuffs.ShadowProtection {
		ShadowProtectionAura(&character.Unit)
	}

	if raidBuffs.AspectOfTheWild {
		AspectOfTheWildAura(&character.Unit)
	}

	if raidBuffs.ElementalResistanceTotem {
		ElementalResistanceTotemAura(&character.Unit)
	}

	// Stamina
	applyStaminaBuffs(&character.Unit, raidBuffs)

	// Strength and Agility
	applyStrengthAgilityBuffs(&character.Unit, raidBuffs)

	// Mana
	if raidBuffs.ArcaneBrilliance {
		ArcaneBrilliance(&character.Unit)
	}

	if raidBuffs.FelIntelligence {
		FelIntelligence(&character.Unit)
	}

	// Melee haste
	if raidBuffs.WindfuryTotem {
		WindfuryTotem(&character.Unit)
	}

	if raidBuffs.IcyTalons {
		IcyTalons(&character.Unit)
	}

	if raidBuffs.HuntingParty {
		HuntingParty(&character.Unit)
	}

	// +Crit%
	if raidBuffs.LeaderOfThePack {
		LeaderOfThePack(&character.Unit)
	}

	if raidBuffs.Rampage {
		Rampage(&character.Unit)
	}

	if raidBuffs.ElementalOath {
		ElementalOath(&character.Unit)
	}

	if raidBuffs.HonorAmongThieves {
		HonorAmongThieves(&character.Unit)
	}

	if raidBuffs.TerrifyingRoar {
		TerrifyingRoar(&character.Unit)
	}

	if raidBuffs.FuriousHowl {
		FuriousHowl(&character.Unit)
	}

	// +% Attackpower
	if raidBuffs.AbominationsMight {
		AbominationsMightAura(&character.Unit)
	}

	if raidBuffs.UnleashedRage {
		UnleashedRageAura(&character.Unit)
	}

	if raidBuffs.TrueshotAura {
		TrueShotAura(&character.Unit)
	}

	if raidBuffs.BlessingOfMight {
		BlessingOfMightAura(&character.Unit)
	}

	// Spell Haste
	if raidBuffs.MoonkinForm {
		MoonkinAura(&character.Unit)
	}

	if raidBuffs.ShadowForm {
		ShadowFormAura(&character.Unit)
	}

	if raidBuffs.WrathOfAirTotem {
		WrathOfAirAura(&character.Unit)
	}

	// Spell Power
	if raidBuffs.FlametongueTotem {
		FlametongueTotem(&character.Unit)
	}

	// Arcane Brilliance already @Mana Buffs
	if raidBuffs.DemonicPact {
		DemonicPact(&character.Unit)
	}

	if raidBuffs.TotemicWrath {
		TotemicWrath(&character.Unit)
	}

	// +DMG%
	if raidBuffs.ArcaneTactics {
		ArcaneTactics(&character.Unit)
	}

	if raidBuffs.FerociousInspiration {
		FerociousInspiration(&character.Unit)
	}

	if raidBuffs.Communion {
		Communion(&character.Unit)
	}

	// MP5
	if raidBuffs.ManaSpringTotem {
		ManaSpringTotem(&character.Unit)
	}

	// Armor
	if raidBuffs.DevotionAura {
		DevotionAura(&character.Unit)
	}

	if raidBuffs.StoneskinTotem {
		StoneskinTotem(&character.Unit)
	}

	// Blessing Of Might @AttackPower&
	// Fel Inteligenc @Mana+

	var replenishmentActionID ActionID
	if individualBuffs.VampiricTouch {
		replenishmentActionID.SpellID = 34914
	} else if individualBuffs.SoulLeach {
		replenishmentActionID.SpellID = 30295
	} else if individualBuffs.Revitalize {
		replenishmentActionID.SpellID = 48544
	} else if individualBuffs.EnduringWinter {
		replenishmentActionID.SpellID = 44561
	} else if individualBuffs.Communion {
		replenishmentActionID.SpellID = 31876
	}

	if !(replenishmentActionID.IsEmptyAction()) {
		MakePermanent(replenishmentAura(&character.Unit, replenishmentActionID))
	}

	// 	character.AddStats(stats.Stats{
	// 		stats.Armor: GetTristateValueFloat(raidBuffs.DevotionAura, 1205, 1807.5),
	// 	})
	// }

	if raidBuffs.RetributionAura {
		RetributionAura(&character.Unit)
	}

	if len(character.Env.Raid.AllPlayerUnits) == 1 {
		if raidBuffs.Bloodlust {
			registerBloodlustCD(agent, 2825)
		} else if raidBuffs.Heroism {
			registerBloodlustCD(agent, 32182)
		} else if raidBuffs.TimeWarp {
			registerBloodlustCD(agent, 80353)
		}

		registerUnholyFrenzyCD(agent, individualBuffs.UnholyFrenzyCount)
		registerTricksOfTheTradeCD(agent, individualBuffs.TricksOfTheTrade)
		registerPowerInfusionCD(agent, individualBuffs.PowerInfusionCount)
		registerManaTideTotemCD(agent, raidBuffs.ManaTideTotemCount)
		registerInnervateCD(agent, individualBuffs.InnervateCount)
		registerDivineGuardianCD(agent, individualBuffs.DivineGuardianCount)
		registerHandOfSacrificeCD(agent, individualBuffs.HandOfSacrificeCount)
		registerPainSuppressionCD(agent, individualBuffs.PainSuppressionCount)
		registerGuardianSpiritCD(agent, individualBuffs.GuardianSpiritCount)

		if individualBuffs.FocusMagic {
			FocusMagicAura(nil, &character.Unit)
		}

		if individualBuffs.DarkIntent && character.Unit.Type == PlayerUnit {
			MakePermanent(DarkIntentAura(&character.Unit, character.Class == proto.Class_ClassWarlock))
		}
	}
}

func DarkIntentAura(unit *Unit, isWarlock bool) *Aura {
	var dotDmgEffect *ExclusiveEffect
	procAura := unit.RegisterAura(Aura{
		Label:     "Dark Intent Proc",
		ActionID:  ActionID{SpellID: 85759},
		Duration:  7 * time.Second,
		MaxStacks: 3,
		OnStacksChange: func(aura *Aura, sim *Simulation, oldStacks, newStacks int32) {
			dotDmgEffect.SetPriority(sim, TernaryFloat64(isWarlock, 0.03, 0.01)*float64(newStacks))
		},
	})
	dotDmgEffect = procAura.NewExclusiveEffect("DarkIntent", false, ExclusiveEffect{
		Priority: 0,
		OnGain: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.PseudoStats.DotDamageMultiplierAdditive += ee.Priority
		},
		OnExpire: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.PseudoStats.DotDamageMultiplierAdditive -= ee.Priority
		},
	})

	// proc this based on the uptime configuration
	// We assume lock precasts dot so first tick might happen after 2 seconds already
	ApplyFixedUptimeAura(procAura, unit.DarkIntentUptimePercent, time.Second*2, time.Second*2)

	// var periodicHandler OnPeriodicDamage
	// if selfBuff {
	// 	periodicHandler = func(_ *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
	// 		if result.Outcome.Matches(OutcomeCrit) && spell.SchoolIndex > stats.SchoolIndexPhysical {
	// 			procAura.Activate(sim)
	// 			procAura.AddStack(sim)
	// 		}
	// 	}
	// }

	return unit.RegisterAura(Aura{
		Label:    "Dark Intent",
		ActionID: ActionID{SpellID: 85767},
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyCastSpeed(1.03)
			aura.Unit.MultiplyAttackSpeed(sim, 1.03)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / 1.03)
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.03)
		},
		// OnPeriodicDamageDealt: periodicHandler,
		BuildPhase: CharacterBuildPhaseBuffs,
	})
}

func StoneskinTotem(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Stoneskin Totem",
		ActionID{SpellID: 8071},
		[]StatConfig{
			{stats.Armor, 4075, false},
		},
	})
}

func DevotionAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Devotion Aura",
		ActionID{SpellID: 465},
		[]StatConfig{
			{stats.Armor, 4075, false},
		},
	})
}

func BlessingOfKingsAura(unit *Unit) *Aura {
	return makeExclusiveAllStatPercentBuff(unit, "Blessing of Kings", ActionID{SpellID: 20217}, 1.05)
}

func MarkOfTheWildAura(unit *Unit) *Aura {
	aura := makeExclusiveAllStatPercentBuff(unit, "Mark of the Wild", ActionID{SpellID: 1126}, 1.05)
	registerExlusiveEffects(aura, []StatConfig{
		{stats.FireResistance, 97, false},
		{stats.FrostResistance, 97, false},
		{stats.ShadowResistance, 97, false},
		{stats.NatureResistance, 97, false},
	})
	return aura
}

func DrumsOfTheBurningWildAura(unit *Unit) *Aura {
	aura := makeExclusiveAllStatPercentBuff(unit, "Drums of the burning Wild", ActionID{ItemID: 63140}, 1.04)
	registerExlusiveEffects(aura, []StatConfig{
		{stats.FireResistance, 78, false},
		{stats.FrostResistance, 78, false},
		{stats.ShadowResistance, 78, false},
		{stats.NatureResistance, 78, false},
	})
	return aura
}

///////////////////////////////////////////////////////////////////////////
//							Resistances
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/cata/item=63140/drums-of-the-burning-wild
// https://www.wowhead.com/cata/spell=1126/mark-of-the-wild
// https://www.wowhead.com/cata/spell=20217/blessing-of-kings
// https://www.wowhead.com/cata/spell=8184/elemental-resistance-totem
// https://www.wowhead.com/cata/spell=19891/resistance-aura
// https://www.wowhead.com/cata/spell=20043/aspect-of-the-wild
// https://www.wowhead.com/cata/spell=27683/shadow-protection

func ElementalResistanceTotemAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Elemental Resistance Totem",
		ActionID{SpellID: 8184},
		[]StatConfig{
			{stats.FireResistance, 195, false},
			{stats.FrostResistance, 195, false},
			{stats.NatureResistance, 195, false},
		},
	})
}

func ResistanceAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Resistance Aura",
		ActionID{SpellID: 19891},
		[]StatConfig{
			{stats.FireResistance, 195, false},
			{stats.FrostResistance, 195, false},
			{stats.ShadowResistance, 195, false},
		},
	})
}

func ShadowProtectionAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Shadow Protection",
		ActionID{SpellID: 27683},
		[]StatConfig{
			{stats.ShadowResistance, 195, false},
		},
	})
}

func AspectOfTheWildAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Aspect of the Wild",
		ActionID{SpellID: 20043},
		[]StatConfig{
			{stats.NatureResistance, 195, false},
		},
	})
}

///////////////////////////////////////////////////////////////////////////
//							Stamina
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/cata/spell=21562/power-word-fortitude
func PowerWordFortitudeAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Power Word: Fortitude",
		ActionID{SpellID: 21562},
		[]StatConfig{
			{stats.Stamina, 585.0, false},
		},
	})
}

// https://www.wowhead.com/cata/spell=6307/blood-pact
func BloodPactAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Blood Pact",
		ActionID{SpellID: 6307},
		[]StatConfig{
			{stats.Stamina, 585.0, false},
		},
	})
}

// https://www.wowhead.com/cata/spell=469/commanding-shout
func CommandingShoutAura(unit *Unit, asExternal bool, withGlyph bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Commanding Shout",
		ActionID{SpellID: 469},
		[]StatConfig{
			{stats.Stamina, 585.0, false},
		},
	})

	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = TernaryDuration(withGlyph, time.Minute*4, time.Minute*2)
	return baseAura
}

func applyStaminaBuffs(unit *Unit, raidBuffs *proto.RaidBuffs) {
	if raidBuffs.PowerWordFortitude {
		PowerWordFortitudeAura(unit)
	}

	if raidBuffs.CommandingShout {
		CommandingShoutAura(unit, true, false)
	}

	if raidBuffs.BloodPact {
		BloodPactAura(unit)
	}
}

///////////////////////////////////////////////////////////////////////////
//							Strength and Agility
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/cata/spell=8075/strength-of-earth-totem
func StrengthOfEarthTotemAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Strength of Earth Totem",
		ActionID{SpellID: 8075},
		[]StatConfig{
			{stats.Agility, 549.0, false},
			{stats.Strength, 549.0, false},
		}})
}
func RoarOfCourageAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Roar of Courage",
		ActionID{SpellID: 93435},
		[]StatConfig{
			{stats.Agility, 549.0, false},
			{stats.Strength, 549.0, false},
		}})
}

// https://www.wowhead.com/cata/spell=57330/horn-of-winter
func HornOfWinterAura(unit *Unit, asExternal bool, withGlyph bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Horn of Winter",
		ActionID{SpellID: 57330},
		[]StatConfig{
			{stats.Agility, 549.0, false},
			{stats.Strength, 549.0, false},
		}})

	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = TernaryDuration(withGlyph, time.Minute*3, time.Minute*2)
	return baseAura
}

// https://www.wowhead.com/cata/spell=6673/battle-shout
func BattleShoutAura(unit *Unit, asExternal bool, withGlyph bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Battle Shout",
		ActionID{SpellID: 6673},
		[]StatConfig{
			{stats.Agility, 549.0, false},
			{stats.Strength, 549.0, false},
		}})

	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = TernaryDuration(withGlyph, time.Minute*4, time.Minute*2)
	return baseAura
}

func applyStrengthAgilityBuffs(unit *Unit, raidBuffs *proto.RaidBuffs) {
	if raidBuffs.StrengthOfEarthTotem {
		MakePermanent(StrengthOfEarthTotemAura(unit))
	}

	if raidBuffs.HornOfWinter {
		MakePermanent(HornOfWinterAura(unit, true, false))
	}

	if raidBuffs.BattleShout {
		MakePermanent(BattleShoutAura(unit, true, false))
	}
}

///////////////////////////////////////////////////////////////////////////
//							Attack Power
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/cata/spell=30808/unleashed-rage
// https://www.wowhead.com/cata/spell=19506/trueshot-aura
// https://www.wowhead.com/cata/spell=53138/abominations-might
// https://www.wowhead.com/cata/spell=19740/blessing-of-might

func UnleashedRageAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Unleashed Rage",
		ActionID{SpellID: 30808},
		[]StatConfig{
			{stats.AttackPower, 1.2, true},
			{stats.RangedAttackPower, 1.1, true},
		}})
}

func TrueShotAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"True Shot Aura",
		ActionID{SpellID: 19506},
		[]StatConfig{
			{stats.AttackPower, 1.2, true},
			{stats.RangedAttackPower, 1.1, true},
		}})
}

func AbominationsMightAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Abominations Might",
		ActionID{SpellID: 53138},
		[]StatConfig{
			{stats.AttackPower, 1.2, true},
			{stats.RangedAttackPower, 1.1, true},
		}})
}

func BlessingOfMightAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Blessing of Might",
		ActionID{SpellID: 19740},
		[]StatConfig{
			{stats.AttackPower, 1.2, true},
			{stats.RangedAttackPower, 1.1, true},
			{stats.MP5, 326, false},
		}})
}

///////////////////////////////////////////////////////////////////////////
//							Mp5
///////////////////////////////////////////////////////////////////////////

func FelIntelligence(unit *Unit) *Aura {
	if !unit.HasManaBar() {
		return nil
	}
	return makeExclusiveBuff(unit, BuffConfig{
		"Fel Intelligence",
		ActionID{SpellID: 54424},
		[]StatConfig{
			{stats.Mana, 2126, false},
			{stats.MP5, 326, false},
		}})
}

func ArcaneBrilliance(unit *Unit) *Aura {
	if !unit.HasManaBar() {
		return nil
	}
	return makeExclusiveBuff(unit, BuffConfig{
		"Arcane Brilliance",
		ActionID{SpellID: 1459},
		[]StatConfig{
			{stats.Mana, 2126, false},
			{stats.SpellPower, 1.06, true},
		}})
}

func ManaSpringTotem(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Mana Spring Totem",
		ActionID{SpellID: 5675},
		[]StatConfig{
			{stats.MP5, 326, false},
		}})
}

// /////////////////////////////////////////////////////////////////////////
//
//	Melee Haste
//
// /////////////////////////////////////////////////////////////////////////
func registerExclusiveMeleeHaste(aura *Aura, value float64) {
	aura.NewExclusiveEffect("AttackSpeed%", false, ExclusiveEffect{
		OnGain: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(s, value)
		},
		OnExpire: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.MultiplyAttackSpeed(s, 1/value)
		},
	})
}

func WindfuryTotem(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Windfury Totem",
		ActionID{SpellID: 8512},
		[]StatConfig{}})

	registerExclusiveMeleeHaste(baseAura, 1.1)
	return baseAura
}

func IcyTalons(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Icy Talons",
		ActionID{SpellID: 55610},
		[]StatConfig{}})

	registerExclusiveMeleeHaste(baseAura, 1.1)
	return baseAura
}

func HuntingParty(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Hunting Party",
		ActionID{SpellID: 53290},
		[]StatConfig{}})

	registerExclusiveMeleeHaste(baseAura, 1.1)
	return baseAura
}

// /////////////////////////////////////////////////////////////////////////
//
//	+Crit %
//
// /////////////////////////////////////////////////////////////////////////

func LeaderOfThePack(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Leader Of The Pack",
		ActionID{SpellID: 17007},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

func ElementalOath(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Elemental Oath",
		ActionID{SpellID: 51470},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

func HonorAmongThieves(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Honor Among Thieves",
		ActionID{SpellID: 51701},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

func Rampage(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Rampage",
		ActionID{SpellID: 29801},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

func TerrifyingRoar(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Terrifying Roar",
		ActionID{SpellID: 90309},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

func FuriousHowl(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Furious Howl",
		ActionID{SpellID: 24604},
		[]StatConfig{
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false},
		}})

	return baseAura
}

// /////////////////////////////////////////////////////////////////////////
//
//	Spell Haste
//
// /////////////////////////////////////////////////////////////////////////
// Builds an ExclusiveEffect representing a SpellHaste bonus multiplier
// spellHastePercent should be given as the percent value i.E. 0.05 for +5%
func registerExclusiveSpellHaste(aura *Aura, spellHastePercent float64) {
	aura.NewExclusiveEffect("SpellHaste%Buff", false, ExclusiveEffect{
		Priority: spellHastePercent,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(1 + ee.Priority)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(1 / (1 + ee.Priority))
		},
	})
}

func MoonkinAura(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{
		"Moonkin Aura",
		ActionID{SpellID: 24858},
		[]StatConfig{}})

	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}

func WrathOfAirAura(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{
		"Wrath of Air",
		ActionID{SpellID: 3738},
		[]StatConfig{}})

	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}

func ShadowFormAura(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{
		"Shadow Form",
		ActionID{SpellID: 15473},
		[]StatConfig{}})

	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}

// /////////////////////////////////////////////////////////////////////////
//
//	Spell Power
//
// /////////////////////////////////////////////////////////////////////////
func FlametongueTotem(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Flametongue Totem",
		ActionID{SpellID: 8227},
		[]StatConfig{
			{stats.SpellPower, 1.06, true},
		}})
}

func DemonicPact(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Demonic Pact",
		ActionID{SpellID: 53646},
		[]StatConfig{
			{stats.SpellPower, 1.1, true},
		}})
}

func TotemicWrath(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Totemic Wrath",
		ActionID{SpellID: 77746},
		[]StatConfig{
			{stats.SpellPower, 1.1, true},
		}})
}

// /////////////////////////////////////////////////////////////////////////
//
//	Damage Done%
//
// /////////////////////////////////////////////////////////////////////////
func registerExclusiveDamageDone(aura *Aura, damageDoneMod float64) {
	aura.NewExclusiveEffect("DamageDone%Buff", false, ExclusiveEffect{
		Priority: damageDoneMod,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.DamageDealtMultiplier *= (1 + ee.Priority)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.PseudoStats.DamageDealtMultiplier /= (1 + ee.Priority)
		},
	})
}

func ArcaneTactics(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{"Arcane Tactics", ActionID{SpellID: 82930}, []StatConfig{}})
	registerExclusiveDamageDone(aura, 0.03)
	return aura
}

func FerociousInspiration(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{"Ferocious Inspiration", ActionID{SpellID: 34460}, []StatConfig{}})
	registerExclusiveDamageDone(aura, 0.03)
	return aura
}

func Communion(unit *Unit) *Aura {
	aura := makeExclusiveBuff(unit, BuffConfig{"Communion", ActionID{SpellID: 31876}, []StatConfig{}})
	registerExclusiveDamageDone(aura, 0.03)
	return aura
}

/////////////
/// OLD /////
////////////

// Applies buffs to pets.
func applyPetBuffEffects(petAgent PetAgent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) {
	// Summoned pets, like Mage Water Elemental, aren't around to receive raid buffs.
	if petAgent.GetPet().IsGuardian() {
		return
	}

	raidBuffs = googleProto.Clone(raidBuffs).(*proto.RaidBuffs)
	partyBuffs = googleProto.Clone(partyBuffs).(*proto.PartyBuffs)
	individualBuffs = googleProto.Clone(individualBuffs).(*proto.IndividualBuffs)

	// Remove buffs that do not apply to pets
	// Or those that will be applied from the player (BL)
	raidBuffs.Bloodlust = false
	raidBuffs.Heroism = false
	raidBuffs.TimeWarp = false
	// Stam
	raidBuffs.PowerWordFortitude = false
	raidBuffs.CommandingShout = false
	raidBuffs.BloodPact = false // does apply to the imp itself, but not to any other pet
	// Str/Agi
	raidBuffs.StrengthOfEarthTotem = false
	raidBuffs.HornOfWinter = false
	raidBuffs.BattleShout = false
	// Crit%
	raidBuffs.LeaderOfThePack = false
	raidBuffs.HonorAmongThieves = false
	raidBuffs.ElementalOath = false
	raidBuffs.Rampage = false
	raidBuffs.TerrifyingRoar = false
	raidBuffs.FuriousHowl = false
	// AP%
	raidBuffs.TrueshotAura = false
	raidBuffs.UnleashedRage = false
	raidBuffs.AbominationsMight = false
	raidBuffs.BlessingOfMight = false
	// SP%
	raidBuffs.ArcaneBrilliance = false
	raidBuffs.DemonicPact = false
	raidBuffs.TotemicWrath = false
	raidBuffs.FlametongueTotem = false
	// +5% Spell haste
	raidBuffs.MoonkinForm = false
	raidBuffs.ShadowForm = false
	raidBuffs.WrathOfAirTotem = false
	// Mana
	raidBuffs.FelIntelligence = false // does apply to the fel hunter itself, but not to any other pet
	// +Armor
	raidBuffs.DevotionAura = false
	raidBuffs.StoneskinTotem = false
	// 10% Haste
	// raidBuffs.HuntingParty = false
	// raidBuffs.IcyTalons = false
	// raidBuffs.WindfuryTotem = false
	// +3% All Damage
	raidBuffs.ArcaneTactics = false
	raidBuffs.FerociousInspiration = false
	raidBuffs.Communion = false
	// +Spell Resistances
	raidBuffs.ElementalResistanceTotem = false
	raidBuffs.ResistanceAura = false
	raidBuffs.ShadowProtection = false
	raidBuffs.AspectOfTheWild = false
	// +5% Base Stats and Spell Resistances
	raidBuffs.MarkOfTheWild = false
	raidBuffs.BlessingOfKings = false
	raidBuffs.DrumsOfTheBurningWild = false

	individualBuffs.HymnOfHopeCount = 0
	individualBuffs.InnervateCount = 0
	individualBuffs.PowerInfusionCount = 0
	individualBuffs.DivineGuardianCount = 0
	individualBuffs.GuardianSpiritCount = 0
	individualBuffs.HandOfSacrificeCount = 0
	individualBuffs.PainSuppressionCount = 0
	individualBuffs.PowerInfusionCount = 0
	individualBuffs.TricksOfTheTrade = proto.TristateEffect_TristateEffectMissing
	individualBuffs.UnholyFrenzyCount = 0

	if !petAgent.GetPet().enabledOnStart {
		// What do we do with permanent pets that are not enabled at start?
	}

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
}

func InspirationAura(unit *Unit, points int32) *Aura {
	multiplier := 1 - []float64{0, .03, .07, .10}[points]

	return unit.GetOrRegisterAura(Aura{
		Label:    "Inspiration",
		ActionID: ActionID{SpellID: 15357},
		Duration: time.Second * 15,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= multiplier
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= multiplier
		},
	})
}

func ApplyInspiration(unit *Unit, uptime float64) {
	if uptime <= 0 {
		return
	}
	uptime = min(1, uptime)

	inspirationAura := InspirationAura(unit, 3)

	ApplyFixedUptimeAura(inspirationAura, uptime, time.Millisecond*2500, 1)
}

func RetributionAura(unit *Unit) *Aura {
	actionID := ActionID{SpellID: 7294}

	baseDamage := 116.0

	procSpell := unit.RegisterSpell(SpellConfig{
		ActionID:    actionID,
		SpellSchool: SpellSchoolHoly,
		ProcMask:    ProcMaskEmpty,
		Flags:       SpellFlagBinary,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHit)
		},
	})

	return unit.RegisterAura(Aura{
		Label:    "Retribution Aura",
		ActionID: actionID,
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if result.Landed() && spell.SpellSchool.Matches(SpellSchoolPhysical) {
				procSpell.Cast(sim, spell.Unit)
			}
		},
	})
}

// Used for approximating cooldowns applied by other players to you, such as
// bloodlust, innervate, power infusion, etc. This is specifically for buffs
// which can be consecutively applied multiple times to a single player.
type externalConsecutiveCDApproximation struct {
	ActionID         ActionID
	AuraTag          string
	CooldownPriority int32
	Type             CooldownType
	AuraDuration     time.Duration
	AuraCD           time.Duration

	// Callback for extra activation conditions.
	ShouldActivate CooldownActivationCondition

	// Applies the buff.
	AddAura CooldownActivation
}

// numSources is the number of other players assigned to apply the buff to this player.
// E.g. the number of other shaman in the group using bloodlust.
func registerExternalConsecutiveCDApproximation(agent Agent, config externalConsecutiveCDApproximation, numSources int32) {
	if numSources == 0 {
		panic("Need at least 1 source!")
	}
	character := agent.GetCharacter()

	var nextExternalIndex int

	externalTimers := make([]*Timer, numSources)
	for i := 0; i < int(numSources); i++ {
		externalTimers[i] = character.NewTimer()
	}
	sharedTimer := character.NewTimer()

	spell := character.RegisterSpell(SpellConfig{
		ActionID: config.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    sharedTimer,
				Duration: config.AuraDuration, // Assumes that multiple buffs are different sources.
			},
		},
		ExtraCastCondition: func(sim *Simulation, target *Unit) bool {
			if !externalTimers[nextExternalIndex].IsReady(sim) {
				return false
			}

			if character.HasActiveAuraWithTag(config.AuraTag) {
				return false
			}

			return true
		},

		ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
			config.AddAura(sim, character)
			externalTimers[nextExternalIndex].Set(sim.CurrentTime + config.AuraCD)

			nextExternalIndex = (nextExternalIndex + 1) % len(externalTimers)

			if externalTimers[nextExternalIndex].IsReady(sim) {
				sharedTimer.Set(sim.CurrentTime + config.AuraDuration)
			} else {
				sharedTimer.Set(sim.CurrentTime + externalTimers[nextExternalIndex].TimeToReady(sim))
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: config.CooldownPriority,
		Type:     config.Type,

		ShouldActivate: config.ShouldActivate,
	})
}

var BloodlustActionID = ActionID{SpellID: 2825}

const SatedAuraLabel = "Sated"
const BloodlustAuraTag = "Bloodlust"
const BloodlustDuration = time.Second * 40
const BloodlustCD = time.Minute * 10

func registerBloodlustCD(agent Agent, spellID int32) {
	character := agent.GetCharacter()
	BloodlustActionID.SpellID = spellID
	bloodlustAura := BloodlustAura(character, -1)

	spell := character.RegisterSpell(SpellConfig{
		ActionID: bloodlustAura.ActionID,
		Flags:    SpellFlagNoOnCastComplete | SpellFlagNoMetrics | SpellFlagNoLogs,

		Cast: CastConfig{
			CD: Cooldown{
				Timer:    character.NewTimer(),
				Duration: BloodlustCD,
			},
		},

		ApplyEffects: func(sim *Simulation, target *Unit, _ *Spell) {
			if !target.HasActiveAura(SatedAuraLabel) {
				bloodlustAura.Activate(sim)
			}
		},
	})

	character.AddMajorCooldown(MajorCooldown{
		Spell:    spell,
		Priority: CooldownPriorityBloodlust,
		Type:     CooldownTypeDPS,
		ShouldActivate: func(sim *Simulation, character *Character) bool {
			// Haste portion doesn't stack with Power Infusion, so prefer to wait.
			return !character.HasActiveAuraWithTag(PowerInfusionAuraTag) && !character.HasActiveAura(SatedAuraLabel)
		},
	})
}

func BloodlustAura(character *Character, actionTag int32) *Aura {
	actionID := BloodlustActionID.WithTag(actionTag)

	sated := character.GetOrRegisterAura(Aura{
		Label:    SatedAuraLabel,
		ActionID: ActionID{SpellID: 57724},
		Duration: time.Minute * 10,
	})

	for _, pet := range character.Pets {
		if !pet.IsGuardian() {
			BloodlustAura(&pet.Character, actionTag)
		}
	}

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
			aura.Unit.MultiplyResourceRegenSpeed(sim, 1.3)
			for _, pet := range character.Pets {
				if pet.IsEnabled() && !pet.IsGuardian() {
					pet.GetAura(aura.Label).Activate(sim)
				}
			}

			sated.Activate(sim)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.3)
			aura.Unit.MultiplyResourceRegenSpeed(sim, 1/1.3)
		},
	})
	multiplyCastSpeedEffect(aura, 1.3)
	return aura
}

var PowerInfusionActionID = ActionID{SpellID: 10060}
var PowerInfusionAuraTag = "PowerInfusion"

const PowerInfusionDuration = time.Second * 15
const PowerInfusionCD = time.Minute * 2

func registerPowerInfusionCD(agent Agent, numPowerInfusions int32) {
	if numPowerInfusions == 0 {
		return
	}

	piAura := PowerInfusionAura(&agent.GetCharacter().Unit, -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         PowerInfusionActionID.WithTag(-1),
			AuraTag:          PowerInfusionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PowerInfusionDuration,
			AuraCD:           PowerInfusionCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Haste portion doesn't stack with Bloodlust, so prefer to wait.
				return !character.HasActiveAuraWithTag(BloodlustAuraTag)
			},
			AddAura: func(sim *Simulation, character *Character) { piAura.Activate(sim) },
		},
		numPowerInfusions)
}

func PowerInfusionAura(character *Unit, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 10060, Tag: actionTag}
	aura := character.GetOrRegisterAura(Aura{
		Label:    "PowerInfusion-" + actionID.String(),
		Tag:      PowerInfusionAuraTag,
		ActionID: actionID,
		Duration: PowerInfusionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				character.PseudoStats.CostMultiplier -= 0.2
			}
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			if character.HasManaBar() {
				character.PseudoStats.CostMultiplier += 0.2
			}
		},
	})
	multiplyCastSpeedEffect(aura, 1.2)
	return aura
}

func multiplyCastSpeedEffect(aura *Aura, multiplier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("MultiplyCastSpeed", false, ExclusiveEffect{
		Priority: multiplier,
		OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(multiplier)
		},
		OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
			ee.Aura.Unit.MultiplyCastSpeed(1 / multiplier)
		},
	})
}

var TricksOfTheTradeAuraTag = "TricksOfTheTrade"

func registerTricksOfTheTradeCD(agent Agent, tristateConfig proto.TristateEffect) {
	if tristateConfig == proto.TristateEffect_TristateEffectMissing {
		return
	}

	// TristateEffectRegular is interpreted as the Glyphed version (since it
	// is weaker).
	damageMultiplier := GetTristateValueFloat(tristateConfig, 1.1, 1.15)
	unit := &agent.GetCharacter().Unit
	tricksAura := TricksOfTheTradeAura(unit, -1, damageMultiplier)

	// Add a small offset to the tooltip CD to account for input delays
	// between the Rogue pressing Tricks and hitting a target.
	effectiveCD := time.Second*30 + unit.ReactionTime

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 57933, Tag: -1},
			AuraTag:          TricksOfTheTradeAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     tricksAura.Duration,
			AuraCD:           effectiveCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return !character.GetExclusiveEffectCategory("PercentDamageModifier").AnyActive()
			},
			AddAura: func(sim *Simulation, character *Character) { tricksAura.Activate(sim) },
		},
		1)
}

func TricksOfTheTradeAura(character *Unit, actionTag int32, damageMult float64) *Aura {
	actionID := ActionID{SpellID: 57933, Tag: actionTag}

	aura := character.GetOrRegisterAura(Aura{
		Label:    "TricksOfTheTrade-" + actionID.String(),
		Tag:      TricksOfTheTradeAuraTag,
		ActionID: actionID,
		Duration: time.Second * 6,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageDealtMultiplier *= damageMult
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageDealtMultiplier /= damageMult
		},
	})

	RegisterPercentDamageModifierEffect(aura, damageMult)
	return aura
}

var UnholyFrenzyAuraTag = "UnholyFrenzy"

const UnholyFrenzyDuration = time.Second * 30
const UnholyFrenzyCD = time.Minute * 3

func registerUnholyFrenzyCD(agent Agent, numUnholyFrenzy int32) {
	if numUnholyFrenzy == 0 {
		return
	}

	ufAura := UnholyFrenzyAura(&agent.GetCharacter().Unit, -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 49016, Tag: -1},
			AuraTag:          UnholyFrenzyAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     UnholyFrenzyDuration,
			AuraCD:           UnholyFrenzyCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return !character.GetExclusiveEffectCategory("PercentDamageModifier").AnyActive()
			},
			AddAura: func(sim *Simulation, character *Character) { ufAura.Activate(sim) },
		},
		numUnholyFrenzy)
}

func UnholyFrenzyAura(character *Unit, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 49016, Tag: actionTag}

	aura := character.GetOrRegisterAura(Aura{
		Label:    "UnholyFrenzy-" + actionID.String(),
		Tag:      UnholyFrenzyAuraTag,
		ActionID: actionID,
		Duration: UnholyFrenzyDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.2)
			aura.Unit.MultiplyResourceRegenSpeed(sim, 1.2)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1/1.2)
			aura.Unit.MultiplyResourceRegenSpeed(sim, 1/1.2)
		},
	})
	return aura
}

func RegisterPercentDamageModifierEffect(aura *Aura, percentDamageModifier float64) *ExclusiveEffect {
	return aura.NewExclusiveEffect("PercentDamageModifier", true, ExclusiveEffect{
		Priority: percentDamageModifier,
	})
}

var DivineGuardianAuraTag = "DivineGuardian"

const DivineGuardianDuration = time.Second * 6
const DivineGuardianCD = time.Minute * 2

func registerDivineGuardianCD(agent Agent, numDivineGuardians int32) {
	if numDivineGuardians == 0 {
		return
	}

	dgAura := DivineGuardianAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 70940, Tag: -1},
			AuraTag:          DivineGuardianAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     DivineGuardianDuration,
			AuraCD:           DivineGuardianCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { dgAura.Activate(sim) },
		},
		numDivineGuardians)
}

func DivineGuardianAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 53530, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "DivineGuardian-" + actionID.String(),
		Tag:      DivineGuardianAuraTag,
		ActionID: actionID,
		Duration: DivineGuardianDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})
}

var HandOfSacrificeAuraTag = "HandOfSacrifice"

const HandOfSacrificeDuration = time.Millisecond * 10500 // subtract Divine Shield GCD
const HandOfSacrificeCD = time.Minute * 5                // use Divine Shield CD here

func registerHandOfSacrificeCD(agent Agent, numSacs int32) {
	if numSacs == 0 {
		return
	}

	hosAura := HandOfSacrificeAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 6940, Tag: -1},
			AuraTag:          HandOfSacrificeAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     HandOfSacrificeDuration,
			AuraCD:           HandOfSacrificeCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				hosAura.Activate(sim)
			},
		},
		numSacs)
}

func HandOfSacrificeAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 6940, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "HandOfSacrifice-" + actionID.String(),
		Tag:      HandOfSacrificeAuraTag,
		ActionID: actionID,
		Duration: HandOfSacrificeDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.7
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.7
		},
	})
}

var PainSuppressionAuraTag = "PainSuppression"

const PainSuppressionDuration = time.Second * 8
const PainSuppressionCD = time.Minute * 3

func registerPainSuppressionCD(agent Agent, numPainSuppressions int32) {
	if numPainSuppressions == 0 {
		return
	}

	psAura := PainSuppressionAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 33206, Tag: -1},
			AuraTag:          PainSuppressionAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     PainSuppressionDuration,
			AuraCD:           PainSuppressionCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) { psAura.Activate(sim) },
		},
		numPainSuppressions)
}

func PainSuppressionAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 33206, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "PainSuppression-" + actionID.String(),
		Tag:      PainSuppressionAuraTag,
		ActionID: actionID,
		Duration: PainSuppressionDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier *= 0.6
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.DamageTakenMultiplier /= 0.6
		},
	})
}

var GuardianSpiritAuraTag = "GuardianSpirit"

const GuardianSpiritDuration = time.Second * 10
const GuardianSpiritCD = time.Minute * 3

func registerGuardianSpiritCD(agent Agent, numGuardianSpirits int32) {
	if numGuardianSpirits == 0 {
		return
	}

	character := agent.GetCharacter()
	gsAura := GuardianSpiritAura(character, -1)
	healthMetrics := character.NewHealthMetrics(ActionID{SpellID: 47788})

	character.AddDynamicDamageTakenModifier(func(sim *Simulation, _ *Spell, result *SpellResult) {
		if (result.Damage >= character.CurrentHealth()) && gsAura.IsActive() {
			result.Damage = character.CurrentHealth()
			character.GainHealth(sim, 0.5*character.MaxHealth(), healthMetrics)
			gsAura.Deactivate(sim)
		}
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 47788, Tag: -1},
			AuraTag:          GuardianSpiritAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     GuardianSpiritDuration,
			AuraCD:           GuardianSpiritCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				gsAura.Activate(sim)
			},
		},
		numGuardianSpirits)
}

func GuardianSpiritAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 47788, Tag: actionTag}

	return character.GetOrRegisterAura(Aura{
		Label:    "GuardianSpirit-" + actionID.String(),
		Tag:      GuardianSpiritAuraTag,
		ActionID: actionID,
		Duration: GuardianSpiritDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.HealingTakenMultiplier *= 1.4
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			character.PseudoStats.HealingTakenMultiplier /= 1.4
		},
	})
}

const ShatteringThrowCD = time.Minute * 5

// func registerShatteringThrowCD(agent Agent, numShatteringThrows int32) {
// 	if numShatteringThrows == 0 {
// 		return
// 	}

// 	stAura := ShatteringThrowAura(agent.GetCharacter().Env.Encounter.TargetUnits[0])

// 	registerExternalConsecutiveCDApproximation(
// 		agent,
// 		externalConsecutiveCDApproximation{
// 			ActionID:         ActionID{SpellID: 64382, Tag: -1},
// 			AuraTag:          ShatteringThrowAuraTag,
// 			CooldownPriority: CooldownPriorityDefault,
// 			AuraDuration:     ShatteringThrowDuration,
// 			AuraCD:           ShatteringThrowCD,
// 			Type:             CooldownTypeDPS,

// 			ShouldActivate: func(sim *Simulation, unit *Unit) bool {
// 				return true
// 			},
// 			AddAura: func(sim *Simulation, unit *Unit) { stAura.Activate(sim) },
// 		},
// 		numShatteringThrows)
// }

var InnervateAuraTag = "Innervate"

const InnervateDuration = time.Second * 10
const InnervateCD = time.Minute * 3

func InnervateManaThreshold(character *Character) float64 {
	if character.Class == proto.Class_ClassMage {
		// Mages burn mana really fast so they need a higher threshold.
		return character.MaxMana() * 0.7
	} else {
		return character.MaxMana() * 0.45
	}
}

func registerInnervateCD(agent Agent, numInnervates int32) {
	if numInnervates == 0 {
		return
	}

	innervateThreshold := 0.0
	var innervateAura *Aura

	character := agent.GetCharacter()
	character.Env.RegisterPostFinalizeEffect(func() {
		innervateThreshold = InnervateManaThreshold(character)
		innervateAura = InnervateAura(character, -1, 0.05)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 29166, Tag: -1},
			AuraTag:          InnervateAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     InnervateDuration,
			AuraCD:           InnervateCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// Only cast innervate when very low on mana, to make sure all other mana CDs are prioritized.
				return character.CurrentMana() <= innervateThreshold
			},
			AddAura: func(sim *Simulation, character *Character) {
				innervateAura.Activate(sim)
			},
		},
		numInnervates)
}

func InnervateAura(character *Character, actionTag int32, reg float64) *Aura {
	actionID := ActionID{SpellID: 29166, Tag: actionTag}
	manaMetrics := character.NewManaMetrics(actionID)
	return character.GetOrRegisterAura(Aura{
		Label:    "Innervate-" + actionID.String(),
		Tag:      InnervateAuraTag,
		ActionID: actionID,
		Duration: InnervateDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			manaPerTick := aura.Unit.MaxMana() * reg / 10.0
			StartPeriodicAction(sim, PeriodicActionOptions{
				Period:   InnervateDuration / 10,
				NumTicks: 10,
				OnAction: func(sim *Simulation) {
					character.AddMana(sim, manaPerTick, manaMetrics)
				},
			})
		},
	})
}

var ManaTideTotemActionID = ActionID{SpellID: 16190}
var ManaTideTotemAuraTag = "ManaTideTotem"

const ManaTideTotemDuration = time.Second * 12
const ManaTideTotemCD = time.Minute * 5

func registerManaTideTotemCD(agent Agent, numManaTideTotems int32) {
	if numManaTideTotems == 0 {
		return
	}

	initialDelay := time.Duration(0)
	var mttAura *Aura

	character := agent.GetCharacter()
	mttAura = ManaTideTotemAura(character, -1)

	character.Env.RegisterPostFinalizeEffect(func() {
		// Use first MTT at 60s, or halfway through the fight, whichever comes first.
		initialDelay = min(character.Env.BaseDuration/2, time.Second*60)
	})

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ManaTideTotemActionID.WithTag(-1),
			AuraTag:          ManaTideTotemAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ManaTideTotemDuration,
			AuraCD:           ManaTideTotemCD,
			Type:             CooldownTypeMana,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				// A normal resto shaman would wait to use MTT.
				return sim.CurrentTime >= initialDelay
			},
			AddAura: func(sim *Simulation, character *Character) {
				mttAura.Activate(sim)
			},
		},
		numManaTideTotems)
}

// TODO: Should this be a raid aura on every character available?
func ManaTideTotemAura(character *Character, actionTag int32) *Aura {
	actionID := ManaTideTotemActionID.WithTag(actionTag)
	dep := character.NewDynamicMultiplyStat(stats.Spirit, 2)
	return character.GetOrRegisterAura(Aura{
		Label:    "ManaTideTotem-" + actionID.String(),
		Tag:      ManaTideTotemAuraTag,
		ActionID: actionID,
		Duration: ManaTideTotemDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, dep)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, dep)
		},
	})
}

const ReplenishmentAuraDuration = time.Second * 15

// Creates the actual replenishment aura for a unit.
func replenishmentAura(unit *Unit, _ ActionID) *Aura {
	if unit.ReplenishmentAura != nil {
		return unit.ReplenishmentAura
	}

	replenishmentDep := unit.NewDynamicStatDependency(stats.Mana, stats.MP5, 0.005)

	unit.ReplenishmentAura = unit.RegisterAura(Aura{
		Label:    "Replenishment",
		ActionID: ActionID{SpellID: 57669},
		Duration: ReplenishmentAuraDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, replenishmentDep)
		},
		OnExpire: func(aura *Aura, sim *Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, replenishmentDep)
		},
	})

	return unit.ReplenishmentAura
}

type ReplenishmentSource int

// Returns a new aura whose activation will give the Replenishment buff to 10 party/raid members.
func (raid *Raid) NewReplenishmentSource(actionID ActionID) ReplenishmentSource {
	newReplSource := ReplenishmentSource(len(raid.curReplenishmentUnits))
	raid.curReplenishmentUnits = append(raid.curReplenishmentUnits, []*Unit{})

	if raid.replenishmentUnits != nil {
		return newReplSource
	}

	// Get the list of all eligible units (party/raid members + their pets, but no guardians).
	var manaUsers []*Unit
	for _, party := range raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			if character.HasManaBar() {
				manaUsers = append(manaUsers, &character.Unit)
			}
		}
		for _, petAgent := range party.Pets {
			pet := petAgent.GetPet()
			if pet.HasManaBar() && !pet.IsGuardian() {
				manaUsers = append(manaUsers, &pet.Unit)
			}
		}
	}
	raid.replenishmentUnits = manaUsers

	// Initialize replenishment aura for all applicable units.
	for _, unit := range raid.replenishmentUnits {
		replenishmentAura(unit, actionID)
	}

	return newReplSource
}

func (raid *Raid) resetReplenishment(_ *Simulation) {
	raid.leftoverReplenishmentUnits = raid.replenishmentUnits
	for i := 0; i < len(raid.curReplenishmentUnits); i++ {
		raid.curReplenishmentUnits[i] = nil
	}
}

func (raid *Raid) ProcReplenishment(sim *Simulation, src ReplenishmentSource) {
	if sim.GetRemainingDuration() <= 0 {
		return
	}
	// If the raid is fully covered by one or more replenishment sources, we can
	// skip the mana sorting.
	if len(raid.curReplenishmentUnits)*10 >= len(raid.replenishmentUnits) {
		if len(raid.curReplenishmentUnits[src]) == 0 {
			if len(raid.leftoverReplenishmentUnits) > 10 {
				raid.curReplenishmentUnits[src] = raid.leftoverReplenishmentUnits[:10]
				raid.leftoverReplenishmentUnits = raid.leftoverReplenishmentUnits[10:]
			} else {
				raid.curReplenishmentUnits[src] = raid.leftoverReplenishmentUnits
				raid.leftoverReplenishmentUnits = nil
			}
		}
		for _, unit := range raid.curReplenishmentUnits[src] {
			unit.ReplenishmentAura.Activate(sim)
		}
		return
	}

	eligible := append(raid.curReplenishmentUnits[src], raid.leftoverReplenishmentUnits...)
	slices.SortFunc(eligible, func(v1, v2 *Unit) int {
		return cmp.Compare(v1.CurrentManaPercent(), v2.CurrentManaPercent())
	})
	raid.curReplenishmentUnits[src] = eligible[:10]
	raid.leftoverReplenishmentUnits = eligible[10:]
	for _, unit := range raid.curReplenishmentUnits[src] {
		unit.ReplenishmentAura.Activate(sim)
	}
	for _, unit := range raid.leftoverReplenishmentUnits {
		unit.ReplenishmentAura.Deactivate(sim)
	}
}

func FocusMagicAura(caster *Unit, target *Unit) (*Aura, *Aura) {
	actionID := ActionID{SpellID: 54648}

	var casterAura *Aura
	var onHitCallback OnSpellHit
	casterIndex := -1
	if caster != nil {
		casterIndex = int(caster.Index)
		casterAura = caster.GetOrRegisterAura(Aura{
			Label:    "Focus Magic",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCritPercent: 3,
				})
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCritPercent: -3,
				})
			},
		})

		onHitCallback = func(_ *Aura, sim *Simulation, _ *Spell, result *SpellResult) {
			if result.DidCrit() {
				casterAura.Activate(sim)
			}
		}
	}

	var aura *Aura
	if target != nil {
		aura = target.GetOrRegisterAura(Aura{
			Label:      "Focus Magic" + strconv.Itoa(casterIndex),
			ActionID:   actionID.WithTag(int32(casterIndex)),
			Duration:   NeverExpires,
			BuildPhase: CharacterBuildPhaseBuffs,
			OnReset: func(aura *Aura, sim *Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: onHitCallback,
		})
		aura.NewExclusiveEffect("FocusMagic", true, ExclusiveEffect{
			OnGain: func(ee *ExclusiveEffect, sim *Simulation) {
				ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCritPercent: 3,
				})
			},
			OnExpire: func(ee *ExclusiveEffect, sim *Simulation) {
				ee.Aura.Unit.AddStatsDynamic(sim, stats.Stats{
					stats.SpellCritPercent: -3,
				})
			},
		})
	}

	return casterAura, aura
}
