package core

import (
	"time"

	googleProto "google.golang.org/protobuf/proto"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
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
			ee.Aura.Unit.EnableBuildPhaseStatDep(s, dep)
		},
		OnExpire: func(ee *ExclusiveEffect, s *Simulation) {
			ee.Aura.Unit.DisableBuildPhaseStatDep(s, dep)
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
func applyBuffEffects(agent Agent, raidBuffs *proto.RaidBuffs, _ *proto.PartyBuffs, individual *proto.IndividualBuffs) {
	char := agent.GetCharacter()
	u := &char.Unit

	// +10% Attack Power
	if raidBuffs.HornOfWinter {
		HornOfWinterAura(u, true, false)
	}
	if raidBuffs.TrueshotAura {
		TrueShotAura(u)
	}
	if raidBuffs.BattleShout {
		BattleShoutAura(u, true, false)
	}
	if raidBuffs.CommandingShout {
		CommandingShoutAura(u, true, false)
	}

	// +10% Melee and Ranged Attack Speed
	if raidBuffs.UnholyAura {
		UnholyAura(u)
	}
	if raidBuffs.CacklingHowl {
		CacklingHowlAura(u)
	}
	if raidBuffs.SerpentsSwiftness {
		SerpentsSwiftnessAura(u)
	}
	if raidBuffs.SwiftbladesCunning {
		SwiftbladesCunningAura(u)
	}
	if raidBuffs.UnleashedRage {
		UnleashedRageAura(u)
	}

	// +10% Spell Power
	if raidBuffs.StillWater {
		StillWaterAura(u)
	}
	if raidBuffs.ArcaneBrilliance {
		ArcaneBrilliance(u)
	}
	if raidBuffs.BurningWrath {
		BurningWrathAura(u)
	}
	if raidBuffs.DarkIntent {
		MakePermanent(DarkIntentAura(u, char.Class == proto.Class_ClassWarlock))
	}

	// +5% Spell Haste
	if raidBuffs.MoonkinAura {
		MoonkinAura(u)
	}
	if raidBuffs.MindQuickening {
		MindQuickeningAura(u)
	}
	if raidBuffs.ShadowForm {
		ShadowFormAura(u)
	}
	if raidBuffs.ElementalOath {
		ElementalOath(u)
	}

	// +5% Critical Strike Chance
	if raidBuffs.LeaderOfThePack {
		LeaderOfThePack(u)
	}
	if raidBuffs.TerrifyingRoar {
		TerrifyingRoar(u)
	}
	if raidBuffs.FuriousHowl {
		FuriousHowl(u)
	}
	if raidBuffs.LegacyOfTheWhiteTiger {
		LegacyOfTheWhiteTiger(u)
	}

	// +3000 Mastery Rating
	if raidBuffs.RoarOfCourage {
		RoarOfCourageAura(u)
	}
	if raidBuffs.SpiritBeastBlessing {
		SpiritBeastBlessingAura(u)
	}
	if raidBuffs.BlessingOfMight {
		BlessingOfMightAura(u)
	}
	if raidBuffs.GraceOfAir {
		GraceOfAirAura(u)
	}

	// +5% Strength, Agility, Intellect
	if raidBuffs.MarkOfTheWild {
		MarkOfTheWildAura(u)
	}
	if raidBuffs.EmbraceOfTheShaleSpider {
		EmbraceOfTheShaleSpiderAura(u)
	}
	if raidBuffs.LegacyOfTheEmperor {
		LegacyOfTheEmperorAura(u)
	}
	if raidBuffs.BlessingOfKings {
		BlessingOfKingsAura(u)
	}

	// +10% Stamina
	if raidBuffs.QirajiFortitude {
		QirajiFortitudeAura(u)
	}
	if raidBuffs.PowerWordFortitude {
		PowerWordFortitudeAura(u)
	}

	// Major Haste handled below
	// Mana Tidal totem count handled below

	// Stamina & Strength/Agility secondary grouping
	applyStaminaBuffs(u, raidBuffs)
	applyStrengthAgilityBuffs(u, raidBuffs)

	// Individual cooldowns and major buffs
	if len(char.Env.Raid.AllPlayerUnits)-char.Env.Raid.NumTargetDummies == 1 {
		// Major Haste
		if raidBuffs.Bloodlust {
			registerBloodlustCD(agent, 2825)
		}
		if raidBuffs.Heroism {
			registerBloodlustCD(agent, 32182)
		}
		if raidBuffs.TimeWarp {
			registerBloodlustCD(agent, 80353)
		}

		// Major Mana Replenishment
		registerManaTideTotemCD(agent, raidBuffs.ManaTideTotemCount)

		// Other individual CDs
		registerUnholyFrenzyCD(agent, individual.UnholyFrenzyCount)
		registerTricksOfTheTradeCD(agent, individual.TricksOfTheTrade)
		registerPowerInfusionCD(agent, individual.PowerInfusionCount)
		registerInnervateCD(agent, individual.InnervateCount)
		registerDivineGuardianCD(agent, individual.DivineGuardianCount)
		registerHandOfSacrificeCD(agent, individual.HandOfSacrificeCount)
		registerPainSuppressionCD(agent, individual.PainSuppressionCount)
		registerGuardianSpiritCD(agent, individual.GuardianSpiritCount)
		registerRallyingCryCD(agent, individual.RallyingCryCount)
		registerShatteringThrowCD(agent, individual.ShatteringThrowCount)

		if individual.FocusMagic {
			FocusMagicAura(nil, u)
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
	return aura
}

func LegacyOfTheEmperorAura(unit *Unit) *Aura {
	return makeExclusiveAllStatPercentBuff(unit, "Legacy of the Emperor", ActionID{SpellID: 115921}, 1.05)
}

func EmbraceOfTheShaleSpiderAura(u *Unit) *Aura {
	return makeExclusiveAllStatPercentBuff(u, "Embrace of the Shale Spider", ActionID{SpellID: 0}, 1.05)
}

///////////////////////////////////////////////////////////////////////////
//							Resistances
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/mop-classic/item=63140/drums-of-the-burning-wild
// https://www.wowhead.com/mop-classic/spell=1126/mark-of-the-wild
// https://www.wowhead.com/mop-classic/spell=20217/blessing-of-kings
// https://www.wowhead.com/mop-classic/spell=8184/elemental-resistance-totem
// https://www.wowhead.com/mop-classic/spell=19891/resistance-aura
// https://www.wowhead.com/mop-classic/spell=20043/aspect-of-the-wild
// https://www.wowhead.com/mop-classic/spell=27683/shadow-protection

///////////////////////////////////////////////////////////////////////////
//							Stamina
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/mop-classic/spell=21562/power-word-fortitude
func PowerWordFortitudeAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Power Word: Fortitude",
		ActionID{SpellID: 21562},
		[]StatConfig{
			{stats.Stamina, 585.0, false},
		},
	})
}

func QirajiFortitudeAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Qiraji Fortitude", ActionID{SpellID: 0}, []StatConfig{{stats.Stamina, 0.10, true}}})
}

// https://www.wowhead.com/mop-classic/spell=6307/blood-pact
func BloodPactAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Blood Pact",
		ActionID{SpellID: 6307},
		[]StatConfig{
			{stats.Stamina, 585.0, false},
		},
	})
}

// https://www.wowhead.com/mop-classic/spell=469/commanding-shout
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

func applyStaminaBuffs(u *Unit, raidBuffs *proto.RaidBuffs) {
	// +10% Stamina buffs
	if raidBuffs.PowerWordFortitude {
		PowerWordFortitudeAura(u)
	}
	if raidBuffs.QirajiFortitude {
		QirajiFortitudeAura(u)
	}
	if raidBuffs.CommandingShout {
		CommandingShoutAura(u, true, false)
	}
}

//////// 3000 Mastery Rating

func RoarOfCourageAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Roar of Courage", ActionID{SpellID: 0}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func SpiritBeastBlessingAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Spirit Beast Blessing", ActionID{SpellID: 0}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func BlessingOfMightAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Blessing of Might", ActionID{SpellID: 0}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func GraceOfAirAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Grace of Air", ActionID{SpellID: 0}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}

///////////////////////////////////////////////////////////////////////////
//							Strength and Agility
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/mop-classic/spell=8075/strength-of-earth-totem
func StrengthOfEarthTotemAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Strength of Earth Totem",
		ActionID{SpellID: 8075},
		[]StatConfig{
			{stats.Agility, 549.0, false},
			{stats.Strength, 549.0, false},
		}})
}

// https://www.wowhead.com/mop-classic/spell=57330/horn-of-winter
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

// https://www.wowhead.com/mop-classic/spell=6673/battle-shout
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

func applyStrengthAgilityBuffs(u *Unit, raidBuffs *proto.RaidBuffs) {
	// +5% Strength & Agility, Int buffs
	if raidBuffs.HornOfWinter {
		MakePermanent(HornOfWinterAura(u, true, false))
	}
	if raidBuffs.BattleShout {
		MakePermanent(BattleShoutAura(u, true, false))
	}
}

///////////////////////////////////////////////////////////////////////////
//							Attack Power
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/mop-classic/spell=30808/unleashed-rage
// https://www.wowhead.com/mop-classic/spell=19506/trueshot-aura
// https://www.wowhead.com/mop-classic/spell=53138/abominations-might
// https://www.wowhead.com/mop-classic/spell=19740/blessing-of-might

func TrueShotAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"True Shot Aura",
		ActionID{SpellID: 19506},
		[]StatConfig{
			{stats.AttackPower, 1.1, true},
			{stats.RangedAttackPower, 1.1, true},
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
func UnholyAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Unholy Aura", ActionID{SpellID: 0}, nil})
	registerExclusiveMeleeHaste(aura, 0.10)
	return aura
}
func CacklingHowlAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Cackling Howl", ActionID{SpellID: 0}, nil})
	registerExclusiveMeleeHaste(aura, 0.10)
	return aura
}
func SerpentsSwiftnessAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Serpent's Swiftness", ActionID{SpellID: 0}, nil})
	registerExclusiveMeleeHaste(aura, 0.10)
	return aura
}
func SwiftbladesCunningAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Swiftblade's Cunning", ActionID{SpellID: 0}, nil})
	registerExclusiveMeleeHaste(aura, 0.10)
	return aura
}
func UnleashedRageAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Unleashed Rage", ActionID{SpellID: 30808}, []StatConfig{{stats.AttackPower, 1.2, true}, {stats.RangedAttackPower, 1.1, true}}})
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

func LegacyOfTheWhiteTiger(unit *Unit) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Legacy of the White Tiger",
		ActionID{SpellID: 116781},
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
func MindQuickeningAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Mind Quickening", ActionID{SpellID: 0}, nil})
	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}
func ShadowFormAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Shadow Form", ActionID{SpellID: 15473}, nil})
	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}
func ElementalOath(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Elemental Oath", ActionID{SpellID: 51470}, nil})
	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}

// /////////////////////////////////////////////////////////////////////////
//
//	Spell Power
//
// /////////////////////////////////////////////////////////////////////////

func StillWaterAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Still Water", ActionID{SpellID: 0}, []StatConfig{{stats.SpellPower, 0.10, true}}})
}
func ArcaneBrilliance(u *Unit) *Aura {
	// Mages: +10% Spell Power
	return makeExclusiveBuff(u, BuffConfig{"Arcane Brilliance", ActionID{SpellID: 1459}, []StatConfig{{stats.SpellPower, 0.10, true}}})
}
func BurningWrathAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Burning Wrath", ActionID{SpellID: 0}, []StatConfig{{stats.SpellPower, 0.10, true}}})
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
	// Str/Agi
	raidBuffs.HornOfWinter = false
	raidBuffs.BattleShout = false
	// Crit%
	raidBuffs.LeaderOfThePack = false
	raidBuffs.ElementalOath = false
	raidBuffs.TerrifyingRoar = false
	raidBuffs.FuriousHowl = false
	raidBuffs.LegacyOfTheWhiteTiger = false
	// AP%
	raidBuffs.TrueshotAura = false
	raidBuffs.UnleashedRage = false
	raidBuffs.BlessingOfMight = false
	// SP%
	raidBuffs.ArcaneBrilliance = false
	// +5% Spell haste
	raidBuffs.ShadowForm = false
	// Mana
	// +Armor
	// 10% Haste
	// raidBuffs.HuntingParty = false
	// raidBuffs.IcyTalons = false
	// raidBuffs.WindfuryTotem = false
	// +3% All Damage
	// +Spell Resistances
	// +5% Base Stats and Spell Resistances
	raidBuffs.MarkOfTheWild = false
	raidBuffs.BlessingOfKings = false
	raidBuffs.LegacyOfTheEmperor = false

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
	individualBuffs.RallyingCryCount = 0

	if !petAgent.GetPet().enabledOnStart {
		// What do we do with permanent pets that are not enabled at start?
	}

	applyBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
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
			return !character.HasActiveAura(SatedAuraLabel)
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

var RallyingCryAuraTag = "RallyingCry"

const RallyingCryDuration = time.Second * 10
const RallyingCryCD = time.Minute * 3

func registerRallyingCryCD(agent Agent, numRallyingCries int32) {
	if numRallyingCries == 0 {
		return
	}

	buffAura := RallyingCryAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 97462, Tag: -1},
			AuraTag:          RallyingCryAuraTag,
			CooldownPriority: CooldownPriorityLow,
			AuraDuration:     RallyingCryDuration,
			AuraCD:           RallyingCryCD,
			Type:             CooldownTypeSurvival,

			ShouldActivate: func(_ *Simulation, _ *Character) bool {
				return true
			},

			AddAura: func(sim *Simulation, _ *Character) {
				buffAura.Activate(sim)
			},
		},
		numRallyingCries,
	)
}

func RallyingCryAura(character *Character, actionTag int32) *Aura {
	actionID := ActionID{SpellID: 97462, Tag: actionTag}
	healthMetrics := character.NewHealthMetrics(actionID)

	var bonusHealth float64

	return character.GetOrRegisterAura(Aura{
		Label:    "RallyingCry-" + actionID.String(),
		Tag:      RallyingCryAuraTag,
		ActionID: actionID,
		Duration: RallyingCryDuration,

		OnGain: func(_ *Aura, sim *Simulation) {
			bonusHealth = character.MaxHealth() * 0.2
			character.UpdateMaxHealth(sim, bonusHealth, healthMetrics)
		},

		OnExpire: func(_ *Aura, sim *Simulation) {
			character.UpdateMaxHealth(sim, -bonusHealth, healthMetrics)
		},
	})
}

const ShatteringThrowCD = time.Minute * 5

func registerShatteringThrowCD(agent Agent, numShatteringThrows int32) {
	if numShatteringThrows == 0 {
		return
	}

	stAura := ShatteringThrowAura(agent.GetCharacter().Env.Encounter.TargetUnits[0], -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 64382, Tag: -1},
			AuraTag:          ShatteringThrowAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     ShatteringThrowDuration,
			AuraCD:           ShatteringThrowCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				stAura.Activate(sim)
			},
		},
		numShatteringThrows)
}
