package frost

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFrostDeathKnight()
}

func TestFrostMasterfrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet:     core.GetGearSet("../../../ui/death_knight/frost/gear_sets", "p1.masterfrost"),
		Talents:     DefaultTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/frost/apls", "masterfrost"),

		ItemFilter: ItemFilterMasterfrost,
	}))
}

func TestFrostTwoHand(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet:     core.GetGearSet("../../../ui/death_knight/frost/gear_sets", "p1.2h"),
		Talents:     DefaultTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/frost/apls", "2h"),

		ItemFilter: ItemFilterTwoHand,
	}))
}

var DefaultTalents = "221111"

var FrostDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfPestilence),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfLoudHorn),
	// No interesting minor glyphs.
}

var PlayerOptionsFrost = &proto.Player_FrostDeathKnight{
	FrostDeathKnight: &proto.FrostDeathKnight{
		Options: &proto.FrostDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76088, // Flask of Winter's Bite
	FoodId:   74646, // Black Pepper Ribs and Shrimp
	PotId:    76095, // Potion of Mogu Power
	PrepotId: 76095, // Potion of Mogu Power
}

var ItemFilterMasterfrost = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	HandTypes: []proto.HandType{
		proto.HandType_HandTypeMainHand,
		proto.HandType_HandTypeOffHand,
		proto.HandType_HandTypeOneHand,
	},
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{},
}

var ItemFilterTwoHand = core.ItemFilter{
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
