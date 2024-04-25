package combat

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterCombatRogue()
}

func TestCombat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassRogue,
		Race:        proto.Race_RaceHuman,
		OtherRaces:  []proto.Race{proto.Race_RaceOrc},
		GearSet:     core.GetGearSet("../../../ui/rogue/combat/gear_sets", "p1_combat_test"),
		Talents:     SubtletyTalents,
		Glyphs:      SubtletyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Combat", SpecOptions: PlayerOptionsID},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Deadly OH Instant", SpecOptions: PlayerOptionsDI},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsDD},
		},
		Rotation:       core.GetAplRotation("../../../ui/rogue/combat/apls", "combat"),
		OtherRotations: []core.RotationCombo{},
		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeBow,
				proto.RangedWeaponType_RangedWeaponTypeCrossbow,
				proto.RangedWeaponType_RangedWeaponTypeGun,
				proto.RangedWeaponType_RangedWeaponTypeThrown,
			},
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeSword,
			},
		},
	}))
}

var SubtletyTalents = "0322-2332030310230012321-003"

var SubtletyGlyphs = &proto.Glyphs{
	Prime1: int32(proto.RoguePrimeGlyph_GlyphOfSinisterStrike),
	Prime2: int32(proto.RoguePrimeGlyph_GlyphOfSliceAndDice),
	Prime3: int32(proto.RoguePrimeGlyph_GlyphOfAdrenalineRush),
}

var PlayerOptionsDI = &proto.Player_CombatRogue{
	CombatRogue: &proto.CombatRogue{
		Options: &proto.CombatRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsID = &proto.Player_CombatRogue{
	CombatRogue: &proto.CombatRogue{
		Options: &proto.CombatRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsDD = &proto.Player_CombatRogue{
	CombatRogue: &proto.CombatRogue{
		Options: &proto.CombatRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsII = &proto.Player_CombatRogue{
	CombatRogue: &proto.CombatRogue{
		Options: &proto.CombatRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_InstantPoison,
			},
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfTheWinds,
	DefaultPotion:   proto.Potions_PotionOfTheTolvir,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
}
