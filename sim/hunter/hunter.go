package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

var TalentTreeSizes = [3]int{19, 19, 20}

const ThoridalTheStarsFuryItemID = 34334

type Hunter struct {
	core.Character

	Talents             *proto.HunterTalents
	Options             *proto.HunterOptions
	BeastMasteryOptions *proto.BeastMasteryHunter_Options
	MarksmanshipOptions *proto.MarksmanshipHunter_Options
	SurvivalOptions     *proto.SurvivalHunter_Options

	pet *HunterPet

	// The most recent time at which moving could have started, for trap weaving.
	mayMoveAt time.Duration

	AspectOfTheHawk *core.Spell
	AspectOfTheFox      *core.Spell

	FireTrapTimer *core.Timer

	// Hunter spells
	KillCommand     *core.Spell
	ArcaneShot      *core.Spell
	ExplosiveTrap   *core.Spell
	KillShot        *core.Spell
	RapidFire       *core.Spell
	MultiShot       *core.Spell
	RaptorStrike    *core.Spell
	SerpentSting    *core.Spell
	SteadyShot      *core.Spell
	ScorpidSting    *core.Spell
	SilencingShot   *core.Spell
	// BM only spells

	// MM only spells
	AimedShot       *core.Spell
	ChimeraShot *core.Spell

	// Survival only spells
	ExplosiveShot *core.Spell
	BlackArrow *core.Spell
	CobraShot *core.Spell



	// Fake spells to encapsulate weaving logic.
	TrapWeaveSpell *core.Spell

	AspectOfTheHawkAura *core.Aura
	ImprovedSteadyShotAura    *core.Aura
	LockAndLoadAura           *core.Aura
	RapidFireAura             *core.Aura
	ScorpidStingAuras         core.AuraArray
	KillingStreakCounterAura 	*core.Aura
	KillingStreakAura 	*core.Aura
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}
func (hunter *Hunter) HasPrimeGlyph(glyph proto.HunterPrimeGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}
func (hunter *Hunter) HasMajorGlyph(glyph proto.HunterMajorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}
func (hunter *Hunter) HasMinorGlyph(glyph proto.HunterMinorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	// if hunter.Talents.TrueshotAura {
	// 	raidBuffs.TrueshotAura = true
	// }
	// if hunter.Talents.FerociousInspiration == 3 && hunter.pet != nil {
	// 	raidBuffs.FerociousInspiration = true
	// }
}
func (hunter *Hunter) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (hunter *Hunter) CritMultiplier(isRanged bool, isMFDSpell bool, doubleDipMS bool) float64 {
	primaryModifier := 1.0
	secondaryModifier := 0.0

	if isRanged {
		//secondaryModifier += mortalShotsFactor * float64(hunter.Talents.MortalShots)
		if isMFDSpell {
			//secondaryModifier += 0.02 * float64(hunter.Talents.MarkedForDeath)
		}
	}

	return hunter.MeleeCritMultiplier(primaryModifier, secondaryModifier)
}

func (hunter *Hunter) Initialize() {
	hunter.AutoAttacks.MHConfig().CritMultiplier = 1 //hunter.critMultiplier(false, false, false)
	hunter.AutoAttacks.OHConfig().CritMultiplier = 1 //hunter.critMultiplier(false, false, false)
	hunter.AutoAttacks.RangedConfig().CritMultiplier = 1// hunter.critMultiplier(false, false, false)
	hunter.FireTrapTimer = hunter.NewTimer()
	hunter.registerSteadyShotSpell() // Todo: Probably add a counter so we can check for imp ss?
	hunter.registerArcaneShotSpell()
	hunter.registerKillShotSpell()
	hunter.registerAspectOfTheHawkSpell()
	hunter.registerSerpentStingSpell()
	hunter.registerMultiShotSpell(hunter.NewTimer())
	hunter.registerKillCommandSpell()
	hunter.registerExplosiveTrapSpell(hunter.FireTrapTimer)
	hunter.registerCobraShotSpell()
}

func (hunter *Hunter) Reset(_ *core.Simulation) {
	hunter.mayMoveAt = 0
}

func NewHunter(character *core.Character, options *proto.Player, hunterOptions *proto.HunterOptions) *Hunter {
	hunter := &Hunter{
		Character: *character,
		Talents:   &proto.HunterTalents{},
		Options:   hunterOptions,
	}
	core.FillTalentsProto(hunter.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	// Todo: Verify that is is actually 4 focus per second
	hunter.EnableFocusBar(100 + (float64(hunter.Talents.KindredSpirits) * 5), 4.0, true)

	hunter.PseudoStats.CanParry = true


	// Passive bonus (used to be from quiver).
	hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
	rangedWeapon := hunter.WeaponFromRanged(0)

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		// We don't know crit multiplier until later when we see the target so just
		// use 0 for now.
		MainHand:        hunter.WeaponFromMainHand(0),
		OffHand:         hunter.WeaponFromOffHand(0),
		Ranged:          rangedWeapon,
		//ReplaceMHSwing:  hunter.TryRaptorStrike, //Todo: Might be weaving
		AutoSwingRanged: true,
	})

	hunter.AutoAttacks.RangedConfig().ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower(target)) +
			spell.BonusWeaponDamage()
		spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
	}

	hunter.pet = hunter.NewHunterPet()

	hunter.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.RangedAttackPower, 1)
	hunter.AddStatDependency(stats.Agility, stats.MeleeCrit, core.CritPerAgiMaxLevel[character.Class]*core.CritRatingPerCritChance)

	return hunter
}

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
