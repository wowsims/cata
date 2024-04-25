package assassination

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterAssassinationRogue()
}

func TestAssassination(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassRogue,
		Race:        proto.Race_RaceHuman,
		OtherRaces:  []proto.Race{proto.Race_RaceOrc},
		GearSet:     core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p1_assassination_test"),
		Talents:     AssassinationTalents,
		Glyphs:      AssassinationGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Assassination", SpecOptions: PlayerOptionsAssassinationDI},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsAssassinationID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsAssassinationII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsAssassinationDD},
		},
		Rotation:       core.GetAplRotation("../../../ui/rogue/assassination/apls", "mutilate"),
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

var AssassinationTalents = "0333230013122110321-002-203003"

var AssassinationGlyphs = &proto.Glyphs{
	Prime1: int32(proto.RoguePrimeGlyph_GlyphOfBackstab),
	Prime2: int32(proto.RoguePrimeGlyph_GlyphOfRupture),
	Prime3: int32(proto.RoguePrimeGlyph_GlyphOfMutilate),
}

var PlayerOptionsAssassinationDI = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationID = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationDD = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationII = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
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
