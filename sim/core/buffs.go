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
		HornOfWinterAura(u, true)
	}
	if raidBuffs.TrueshotAura {
		TrueShotAura(u)
	}
	if raidBuffs.BattleShout {
		BattleShoutAura(u, true)
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
		MakePermanent(DarkIntentAura(u))
	}

	// +5% Spell Haste
	if raidBuffs.MoonkinAura {
		MoonkinAura(u)
	}
	if raidBuffs.MindQuickening {
		MindQuickeningAura(u)
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

	// Stamina & Strength/Agility secondary grouping
	applyStaminaBuffs(u, raidBuffs)

	registerManaTideTotemCD(agent, raidBuffs.ManaTideTotemCount)
	registerSkullBannerCD(agent, raidBuffs.SkullBannerCount)
	registerStormLashCD(agent, raidBuffs.StormlashTotemCount)

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

		// Other individual CDs
		registerUnholyFrenzyCD(agent, individual.UnholyFrenzyCount)
		registerTricksOfTheTradeCD(agent, individual.TricksOfTheTrade)
		registerHandOfSacrificeCD(agent, individual.HandOfSacrificeCount)
		registerPainSuppressionCD(agent, individual.PainSuppressionCount)
		registerGuardianSpiritCD(agent, individual.GuardianSpiritCount)
		registerRallyingCryCD(agent, individual.RallyingCryCount)
		registerShatteringThrowCD(agent, individual.ShatteringThrowCount)
	}
}

///////////////////////////////////////////////////////////////////////////
//							Strength, Agility, Intellect 5%
///////////////////////////////////////////////////////////////////////////

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
	return makeExclusiveAllStatPercentBuff(u, "Embrace of the Shale Spider", ActionID{SpellID: 90363}, 1.05)
}

///////////////////////////////////////////////////////////////////////////
//							Stamina
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/mop-classic/spell=21562/power-word-fortitude
func PowerWordFortitudeAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Power Word: Fortitude",
		ActionID{SpellID: 21562},
		[]StatConfig{
			{stats.Stamina, 1.1, true},
		},
	})
}

func QirajiFortitudeAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Qiraji Fortitude", ActionID{SpellID: 90364}, []StatConfig{{stats.Stamina, 1.1, true}}})
}
func CommandingShoutAura(unit *Unit, asExternal bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Commanding Shout",
		ActionID{SpellID: 469},
		[]StatConfig{
			{stats.Stamina, 1.1, true},
		}})
	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = time.Minute * 5
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
		CommandingShoutAura(u, true)
	}
}

//////// 3000 Mastery Rating

func RoarOfCourageAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Roar of Courage", ActionID{SpellID: 93435}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func SpiritBeastBlessingAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Spirit Beast Blessing", ActionID{SpellID: 128997}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func BlessingOfMightAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Blessing of Might", ActionID{SpellID: 19740}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}
func GraceOfAirAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Grace of Air", ActionID{SpellID: 116956}, []StatConfig{{stats.MasteryRating, 3000, false}}})
}

///////////////////////////////////////////////////////////////////////////
//							Attack Power
///////////////////////////////////////////////////////////////////////////

func TrueShotAura(unit *Unit) *Aura {
	return makeExclusiveBuff(unit, BuffConfig{
		"Trueshot Aura",
		ActionID{SpellID: 19506},
		[]StatConfig{
			{stats.AttackPower, 1.1, true},
			{stats.RangedAttackPower, 1.1, true},
		}})
}

func HornOfWinterAura(unit *Unit, asExternal bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Horn of Winter",
		ActionID{SpellID: 57330},
		[]StatConfig{
			{stats.AttackPower, 1.1, true},
			{stats.RangedAttackPower, 1.1, true},
		}})

	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = time.Minute * 5
	return baseAura
}

func BattleShoutAura(unit *Unit, asExternal bool) *Aura {
	baseAura := makeExclusiveBuff(unit, BuffConfig{
		"Battle Shout",
		ActionID{SpellID: 6673},
		[]StatConfig{
			{stats.AttackPower, 1.1, true},
			{stats.RangedAttackPower, 1.1, true},
		}})

	if asExternal {
		return baseAura
	}

	baseAura.OnReset = nil
	baseAura.Duration = time.Minute * 5
	return baseAura
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
	aura := makeExclusiveBuff(u, BuffConfig{"Unholy Aura", ActionID{SpellID: 55610}, nil})
	registerExclusiveMeleeHaste(aura, 1.10)
	return aura
}
func CacklingHowlAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Cackling Howl", ActionID{SpellID: 128432}, nil})
	registerExclusiveMeleeHaste(aura, 1.10)
	return aura
}
func SerpentsSwiftnessAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Serpent's Swiftness", ActionID{SpellID: 128433}, nil})
	registerExclusiveMeleeHaste(aura, 1.10)
	return aura
}
func SwiftbladesCunningAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Swiftblade's Cunning", ActionID{SpellID: 113742}, nil})
	registerExclusiveMeleeHaste(aura, 1.10)
	return aura
}
func UnleashedRageAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Unleashed Rage", ActionID{SpellID: 30809}, nil})
	registerExclusiveMeleeHaste(aura, 1.10)
	return aura
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
		ActionID{SpellID: 24907},
		[]StatConfig{}})

	registerExclusiveSpellHaste(aura, 0.05)
	return aura
}

func MindQuickeningAura(u *Unit) *Aura {
	aura := makeExclusiveBuff(u, BuffConfig{"Mind Quickening", ActionID{SpellID: 49868}, nil})
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
	return makeExclusiveBuff(u, BuffConfig{"Still Water", ActionID{SpellID: 126309},
		[]StatConfig{
			{stats.SpellPower, 1.10, true},
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false}}})
}
func ArcaneBrilliance(u *Unit) *Aura {
	// Mages: +10% Spell Power
	return makeExclusiveBuff(u, BuffConfig{"Arcane Brilliance", ActionID{SpellID: 1459},
		[]StatConfig{
			{stats.SpellPower, 1.10, true},
			{stats.PhysicalCritPercent, 5, false},
			{stats.SpellCritPercent, 5, false}}})
}
func BurningWrathAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Burning Wrath", ActionID{SpellID: 77747}, []StatConfig{{stats.SpellPower, 1.10, true}}})
}
func DarkIntentAura(u *Unit) *Aura {
	return makeExclusiveBuff(u, BuffConfig{"Dark Intent", ActionID{SpellID: 109773}, []StatConfig{{stats.SpellPower, 1.10, true}, {stats.Stamina, 1.10, true}}})
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

	//Todo: Only cancel the buffs that are supposed to be cancelled
	// Check beta when pets are better implemented?
	raidBuffs = &proto.RaidBuffs{}
	partyBuffs = &proto.PartyBuffs{}
	individualBuffs = &proto.IndividualBuffs{}

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

	aura := character.GetOrRegisterAura(Aura{
		Label:    "Bloodlust-" + actionID.String(),
		Tag:      BloodlustAuraTag,
		ActionID: actionID,
		Duration: BloodlustDuration,
		OnGain: func(aura *Aura, sim *Simulation) {
			aura.Unit.MultiplyAttackSpeed(sim, 1.3)
			aura.Unit.MultiplyResourceRegenSpeed(sim, 1.3)
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
			AddAura: func(sim *Simulation, character *Character) {
				tricksAura.Activate(sim)
			},
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
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.DamageDealtMultiplier, damageMult)

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
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.DamageTakenMultiplier, 0.7)
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
			AddAura: func(sim *Simulation, character *Character) {
				psAura.Activate(sim)
			},
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
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.DamageTakenMultiplier, 0.6)
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
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.HealingTakenMultiplier, 1.4)
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

const SkullBannerAuraTag = "SkullBanner"
const SkullBannerDuration = time.Second * 10
const SkullBannerCD = time.Minute * 3

func registerSkullBannerCD(agent Agent, numSkullBanners int32) {
	if numSkullBanners == 0 {
		return
	}

	sbAura := SkullBannerAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 114207, Tag: -1},
			AuraTag:          SkullBannerAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     SkullBannerDuration,
			AuraCD:           SkullBannerCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				sbAura.Activate(sim)
			},
		},
		numSkullBanners)
}

func SkullBannerAura(character *Character, actionTag int32) *Aura {
	return character.GetOrRegisterAura(Aura{
		Label:    "Skull Banner",
		Tag:      SkullBannerAuraTag,
		ActionID: ActionID{SpellID: 114206, Tag: actionTag},
		Duration: SkullBannerDuration,
	}).AttachMultiplicativePseudoStatBuff(&character.PseudoStats.CritDamageMultiplier, 1.2)
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

func ManaTideTotemAura(character *Character, actionTag int32) *Aura {
	actionID := ManaTideTotemActionID.WithTag(actionTag)
	dep := character.NewDynamicMultiplyStat(stats.Spirit, 2)
	return character.GetOrRegisterAura(Aura{
		Label:    "ManaTideTotem-" + actionID.String(),
		Tag:      ManaTideTotemAuraTag,
		ActionID: actionID,
		Duration: ManaTideTotemDuration,
	}).AttachStatDependency(dep)
}

const StormLashAuraTag = "StormLash"
const StormLashDuration = time.Second * 10
const StormLashCD = time.Minute * 5

func registerStormLashCD(agent Agent, numStormLashes int32) {
	if numStormLashes == 0 {
		return
	}

	sbAura := StormLashAura(agent.GetCharacter(), -1)

	registerExternalConsecutiveCDApproximation(
		agent,
		externalConsecutiveCDApproximation{
			ActionID:         ActionID{SpellID: 120668, Tag: -1},
			AuraTag:          StormLashAuraTag,
			CooldownPriority: CooldownPriorityDefault,
			AuraDuration:     StormLashDuration,
			AuraCD:           StormLashCD,
			Type:             CooldownTypeDPS,

			ShouldActivate: func(sim *Simulation, character *Character) bool {
				return true
			},
			AddAura: func(sim *Simulation, character *Character) {
				sbAura.Activate(sim)
			},
		},
		numStormLashes)
}

var StormLashSpellExceptions = map[int32]float64{
	1120:   2.0, // Drain Soul
	45284:  2.0, // Lightning Bolt
	51505:  2.0, // Lava Burst
	103103: 2.0, // Malefic Grasp
	15407:  1.0, // Mind Flay
	129197: 1.0, // Mind Flay - Insanity
}

// Source: https://www.wowhead.com/mop-classic/spell=120668/stormlash-totem#comments
func StormLashAura(character *Character, actionTag int32) *Aura {
	for _, pet := range character.Pets {
		if !pet.IsGuardian() {
			StormLashAura(&pet.Character, actionTag)
		}
	}

	damage := 0.0

	stormlashSpell := character.RegisterSpell(SpellConfig{
		ActionID:    ActionID{SpellID: 120687, Tag: actionTag},
		Flags:       SpellFlagNoOnCastComplete | SpellFlagPassiveSpell,
		SpellSchool: SpellSchoolNature,
		ProcMask:    ProcMaskEmpty,

		DamageMultiplier: 1,
		CritMultiplier:   character.DefaultCritMultiplier(),

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
		},
	})

	getStormLashSpellOverride := func(spell *Spell) float64 {
		return StormLashSpellExceptions[spell.ActionID.SpellID]
	}

	handler := func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
		if !aura.Icd.IsReady(sim) || !result.Landed() || result.Damage <= 0 || !spell.ProcMask.Matches(ProcMaskDirect|ProcMaskSpecial) || !sim.Proc(0.5, "Stormlash") {
			return
		}

		baseMultiplierExtension := getStormLashSpellOverride(spell)
		ap := Ternary(spell.IsMelee(), stormlashSpell.MeleeAttackPower(), stormlashSpell.RangedAttackPower())
		sp := stormlashSpell.SpellPower()
		scaledAP := ap * 0.2
		scaledSP := sp * 0.3

		baseDamage := max(scaledAP, scaledSP)
		baseMultiplier := 2.0
		speedMultiplier := 1.0
		if baseMultiplierExtension != 0 {
			baseMultiplier = baseMultiplier * baseMultiplierExtension
		}
		if spell.Unit.Type == PetUnit {
			baseMultiplier *= 0.2
		}

		if spell.ProcMask.Matches(ProcMaskWhiteHit) {
			swingSpeed := 0.0
			baseMultiplier *= 0.4

			if spell.IsRanged() {
				ranged := spell.Unit.AutoAttacks.Ranged()
				if ranged != nil {
					swingSpeed = ranged.SwingSpeed
				}
			} else if spell.IsMH() {
				mh := spell.Unit.AutoAttacks.MH()
				if mh != nil {
					swingSpeed = mh.SwingSpeed
				}
			} else {
				baseMultiplier /= 2
				oh := spell.Unit.AutoAttacks.OH()
				if oh != nil {
					swingSpeed = oh.SwingSpeed
				}
			}

			speedMultiplier = swingSpeed / 2.6
		} else {
			speedMultiplier = max(spell.DefaultCast.CastTime.Seconds(), 1.5) / 1.5
		}

		avg := baseDamage * baseMultiplier * speedMultiplier
		min, max := ApplyVarianceMinMax(avg, 0.30)
		damage = sim.RollWithLabel(min, max, StormLashAuraTag)

		if sim.Log != nil {
			var chosenStat = Ternary(scaledAP > scaledSP, stats.AttackPower, stats.SpellPower)
			var statValue = Ternary(chosenStat == stats.AttackPower, ap, sp)

			character.Log(sim, "[DEBUG] Damage portion for Stormlash procced by %s: Stat=%s, BaseStatValue=%0.2f, BaseDamage=%0.2f, BaseMultiplier=%0.2f, SpeedMultiplier=%0.2f, PreOutcomeDamageAvg=%0.2f, PreOutcomeDamageMin=%0.2f, PreOutcomeDamageMax=%0.2f, PreOutcomeDamageActual=%0.2f",
				spell.ActionID, chosenStat.StatName(), statValue, baseDamage, baseMultiplier, speedMultiplier, avg, min, max, damage)
		}
		stormlashSpell.Cast(sim, result.Target)
		aura.Icd.Use(sim)
	}

	return character.RegisterAura(Aura{
		Label:    "Stormlash Totem",
		Tag:      StormLashAuraTag,
		ActionID: ActionID{SpellID: 120668, Tag: actionTag},
		Duration: StormLashDuration,
		Icd: &Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Millisecond * 70,
		},
		OnGain: func(aura *Aura, sim *Simulation) {
			for _, pet := range character.Pets {
				if pet.IsEnabled() && !pet.IsGuardian() {
					pet.GetAura(aura.Label).Activate(sim)
				}
			}
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			handler(aura, sim, spell, result)
		},
		OnPeriodicDamageDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			isValidDot := getStormLashSpellOverride(spell) != 0
			if isValidDot {
				handler(aura, sim, spell, result)
			}
		},
	})
}
