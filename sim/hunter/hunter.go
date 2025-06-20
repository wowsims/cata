package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const ThoridalTheStarsFuryItemID = 34334

type Hunter struct {
	core.Character

	ClassSpellScaling float64

	Talents             *proto.HunterTalents
	Options             *proto.HunterOptions
	BeastMasteryOptions *proto.BeastMasteryHunter_Options
	MarksmanshipOptions *proto.MarksmanshipHunter_Options
	SurvivalOptions     *proto.SurvivalHunter_Options

	Pet          *HunterPet
	StampedePet  []*HunterPet
	DireBeastPet *HunterPet

	// The most recent time at which moving could have started, for trap weaving.
	mayMoveAt time.Duration

	AspectOfTheHawk *core.Spell

	// Hunter spells
	SerpentSting         *core.Spell
	ExplosiveTrap        *core.Spell
	ExplosiveShot        *core.Spell
	ImprovedSerpentSting *core.Spell

	// Fake spells to encapsulate weaving logic.
	HuntersMarkSpell *core.Spell
}

func (hunter *Hunter) GetCharacter() *core.Character {
	return &hunter.Character
}

func (hunter *Hunter) GetHunter() *Hunter {
	return hunter
}

func NewHunter(character *core.Character, options *proto.Player, hunterOptions *proto.HunterOptions) *Hunter {
	hunter := &Hunter{
		Character:         *character,
		Talents:           &proto.HunterTalents{},
		Options:           hunterOptions,
		ClassSpellScaling: core.GetClassSpellScalingCoefficient(proto.Class_ClassHunter),
	}

	core.FillTalentsProto(hunter.Talents.ProtoReflect(), options.TalentsString)
	focusPerSecond := 5.0

	// TODO: Fix this to work with the new talent system.
	// hunter.EnableFocusBar(100+(float64(hunter.Talents.KindredSpirits)*5), focusPerSecond, true, nil)
	hunter.EnableFocusBar(100, focusPerSecond, true, nil)

	hunter.PseudoStats.CanParry = true

	// Passive bonus (used to be from quiver).
	//hunter.PseudoStats.RangedSpeedMultiplier *= 1.15
	rangedWeapon := hunter.WeaponFromRanged(0)

	hunter.EnableAutoAttacks(hunter, core.AutoAttackOptions{
		Ranged: rangedWeapon,
		//ReplaceMHSwing:  hunter.TryRaptorStrike, //Todo: Might be weaving
		AutoSwingRanged: true,
		AutoSwingMelee:  false,
	})

	hunter.AutoAttacks.RangedConfig().ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := hunter.RangedWeaponDamage(sim, spell.RangedAttackPower())

		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

		spell.WaitTravelTime(sim, func(sim *core.Simulation) {
			spell.DealDamage(sim, result)
		})
	}

	hunter.AddStatDependencies()

	hunter.Pet = hunter.NewHunterPet()
	hunter.StampedePet = make([]*HunterPet, 4)
	for index := range 4 {
		hunter.StampedePet[index] = hunter.NewStampedePet(index)
	}

	hunter.DireBeastPet = hunter.NewDireBeastPet()
	return hunter
}

func (hunter *Hunter) Initialize() {
	hunter.AutoAttacks.RangedConfig().CritMultiplier = hunter.DefaultCritMultiplier()
	// hunter.addBloodthirstyGloves()
	// Add Stampede pets

	// Add Dire Beast pet
	// hunter.ApplyGlyphs()

	hunter.RegisterSpells()

}

func (hunter *Hunter) GetBaseDamageFromCoeff(coeff float64) float64 {
	return coeff * hunter.ClassSpellScaling
}

func (hunter *Hunter) ApplyTalents() {
	hunter.applyThrillOfTheHunt()
	hunter.ApplyHotfixes()

	if hunter.Pet != nil {

		hunter.Pet.ApplyTalents()
	}

	hunter.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeMail, 86538)
}

func (hunter *Hunter) RegisterSpells() {
	hunter.registerSteadyShotSpell()
	hunter.registerArcaneShotSpell()
	hunter.registerKillShotSpell()
	hunter.registerHawkSpell()
	hunter.RegisterLynxRushSpell()
	hunter.registerSerpentStingSpell()
	hunter.registerMultiShotSpell()
	hunter.registerKillCommandSpell()
	hunter.registerExplosiveTrapSpell()
	hunter.registerCobraShotSpell()
	hunter.registerRapidFireCD()
	hunter.registerSilencingShotSpell()
	hunter.registerHuntersMarkSpell()
	hunter.registerAMOCSpell()
	hunter.registerBarrageSpell()
	hunter.registerGlaiveTossSpell()
	hunter.registerBarrageSpell()
	hunter.registerFervorSpell()
	hunter.RegisterDireBeastSpell()
	hunter.RegisterStampedeSpell()
	hunter.registerPowerShotSpell()
}

func (hunter *Hunter) AddStatDependencies() {
	hunter.AddStatDependency(stats.Agility, stats.AttackPower, 2)
	hunter.AddStatDependency(stats.Agility, stats.RangedAttackPower, 2)
	hunter.AddStatDependency(stats.Agility, stats.PhysicalCritPercent, core.CritPerAgiMaxLevel[hunter.Class])
}

func (hunter *Hunter) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.TrueshotAura = true

	// if hunter.Talents.FerociousInspiration && hunter.Options.PetType != proto.HunterOptions_PetNone {
	// 	raidBuffs.FerociousInspiration = true
	// }

	if hunter.Options.PetType == proto.HunterOptions_CoreHound {
		raidBuffs.Bloodlust = true
	}
	switch hunter.Options.PetType {
	case proto.HunterOptions_CoreHound:
		raidBuffs.Bloodlust = true

	case proto.HunterOptions_ShaleSpider:
		raidBuffs.EmbraceOfTheShaleSpider = true

	case proto.HunterOptions_Wolf:
		raidBuffs.FuriousHowl = true
	case proto.HunterOptions_Devilsaur:
		raidBuffs.TerrifyingRoar = true
	case proto.HunterOptions_WaterStrider:
		raidBuffs.StillWater = true
	case proto.HunterOptions_Hyena:
		raidBuffs.CacklingHowl = true
	case proto.HunterOptions_Serpent:
		raidBuffs.SerpentsSwiftness = true
	case proto.HunterOptions_SporeBat:
		raidBuffs.MindQuickening = true
	case proto.HunterOptions_Cat:
		raidBuffs.RoarOfCourage = true
	case proto.HunterOptions_SpiritBeast:
		raidBuffs.SpiritBeastBlessing = true
	}
	// if hunter.Options.PetType == proto.HunterOptions_ShaleSpider {
	// 	raidBuffs.BlessingOfKings = true
	// }

	// if hunter.Options.PetType == proto.HunterOptions_Wolf || hunter.Options.PetType == proto.HunterOptions_Devilsaur {
	// 	raidBuffs.FuriousHowl = true
	// }

	// TODO: Fix this to work with the new talent system.
	//
	//	if hunter.Talents.HuntingParty {
	//		raidBuffs.HuntingParty = true
	//	}
}

func (hunter *Hunter) AddPartyBuffs(_ *proto.PartyBuffs) {
}

func (hunter *Hunter) HasMajorGlyph(glyph proto.HunterMajorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}
func (hunter *Hunter) HasMinorGlyph(glyph proto.HunterMinorGlyph) bool {
	return hunter.HasGlyph(int32(glyph))
}

func (hunter *Hunter) Reset(_ *core.Simulation) {
	hunter.mayMoveAt = 0
}

const (
	HunterSpellFlagsNone int64 = 0
	SpellMaskSpellRanged int64 = 1 << iota
	HunterSpellAutoShot
	HunterSpellSteadyShot
	HunterSpellCobraShot
	HunterSpellArcaneShot
	HunterSpellKillCommand
	HunterSpellChimeraShot
	HunterSpellExplosiveShot
	HunterSpellExplosiveTrap
	HunterSpellBlackArrow
	HunterSpellMultiShot
	HunterSpellAimedShot
	HunterSpellSerpentSting
	HunterSpellKillShot
	HunterSpellRapidFire
	HunterSpellBestialWrath
	HunterPetFocusDump
	HunterPetDamage
	HunterSpellsTierTwelve = HunterSpellArcaneShot | HunterSpellKillCommand | HunterSpellChimeraShot | HunterSpellExplosiveShot |
		HunterSpellMultiShot | HunterSpellAimedShot
	HunterSpellsAll = HunterSpellSteadyShot | HunterSpellCobraShot |
		HunterSpellArcaneShot | HunterSpellKillCommand | HunterSpellChimeraShot | HunterSpellExplosiveShot |
		HunterSpellExplosiveTrap | HunterSpellBlackArrow | HunterSpellMultiShot | HunterSpellAimedShot |
		HunterSpellSerpentSting | HunterSpellKillShot | HunterSpellRapidFire | HunterSpellBestialWrath
)

// Agent is a generic way to access underlying hunter on any of the agents.
type HunterAgent interface {
	GetHunter() *Hunter
}
