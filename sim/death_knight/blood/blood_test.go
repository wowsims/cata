package blood

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBloodDeathKnight()
}

func TestBlood(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/death_knight/blood/gear_sets", "p1"),

		Talents:     BloodTalents,
		Glyphs:      BloodDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBlood},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/blood/apls", "defensive"),

		InFrontOfTarget: true,
		IsTank:          true,

		ItemFilter: ItemFilter,
	}))
}

var BloodTalents = "131131"
var BloodDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfFesteringBlood),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfRegenerativeMagic),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfOutbreak),
	Minor1: int32(proto.DeathKnightMinorGlyph_GlyphOfTheLongWinter),
}

var PlayerOptionsBlood = &proto.Player_BloodDeathKnight{
	BloodDeathKnight: &proto.BloodDeathKnight{
		Options: &proto.BloodDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76087, // Flask of the Earth
	FoodId:   74656, // Chun Tian Spring Rolls
	PotId:    76095, // Potion of Mogu Power
	PrepotId: 76095, // Potion of Mogu Power
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	HandTypes: []proto.HandType{
		proto.HandType_HandTypeTwoHand,
	},
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{},
}
