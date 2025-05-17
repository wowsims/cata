package balance

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceNightElf,

		GearSet: core.GetGearSet("../../../ui/druid/balance/gear_sets", "t13"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/druid/balance/gear_sets", "t12"),
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBalance},
		Rotation:    core.GetAplRotation("../../../ui/druid/balance/apls", "t13"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/druid/balance/apls", "t12"),
		},
		ItemFilter: ItemFilter,
	}))
}

var StandardTalents = "33230221123212111001-01-020331"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfStarfall),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfRebirth),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfMonsoon),
	Minor1: int32(proto.DruidMinorGlyph_GlyphOfTyphoon),
	Minor2: int32(proto.DruidMinorGlyph_GlyphOfUnburdenedRebirth),
	Minor3: int32(proto.DruidMinorGlyph_GlyphOfMarkOfTheWild),
}

var PlayerOptionsBalance = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options: &proto.BalanceDruid_Options{
			ClassOptions: &proto.DruidOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58086, // Flask of the Draconic Mind
	FoodId:   62671, // Severed Sagefish Head
	PotId:    58091, // Volcanic Potion
	PrepotId: 58091, // Volcanic Potion
}
var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypePolearm,
	},
	ArmorType:         proto.ArmorType_ArmorTypeLeather,
	RangedWeaponTypes: []proto.RangedWeaponType{},
}
