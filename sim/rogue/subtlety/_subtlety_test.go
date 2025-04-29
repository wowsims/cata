package subtlety

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterSubtletyRogue()
}

func TestSubtlety(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassRogue,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},
		GearSet:    core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "p1_subtlety"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "p3_subtlety"),
			core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "p4_subtlety"),
		},
		Talents:     SubtletyTalents,
		Glyphs:      SubtletyGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Subtlety", SpecOptions: PlayerOptionsID},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Deadly OH Instant", SpecOptions: PlayerOptionsDI},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsDD},
		},
		Rotation:       core.GetAplRotation("../../../ui/rogue/subtlety/apls", "subtlety"),
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
			},
		},
	}))
}

var SubtletyTalents = "023003-002-0332031321310012321"

var SubtletyGlyphs = &proto.Glyphs{}

var PlayerOptionsDI = &proto.Player_SubtletyRogue{
	SubtletyRogue: &proto.SubtletyRogue{
		Options: &proto.SubtletyRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsID = &proto.Player_SubtletyRogue{
	SubtletyRogue: &proto.SubtletyRogue{
		Options: &proto.SubtletyRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsDD = &proto.Player_SubtletyRogue{
	SubtletyRogue: &proto.SubtletyRogue{
		Options: &proto.SubtletyRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsII = &proto.Player_SubtletyRogue{
	SubtletyRogue: &proto.SubtletyRogue{
		Options: &proto.SubtletyRogue_Options{
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
