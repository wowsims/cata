package fire

import (
	"testing"

	"github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFireMage()
	common.RegisterAllEffects()
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/mage/fire/gear_sets", "p1_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/mage/fire/gear_sets", "p1_prebis"),
		},
		Talents:     FireTalents,
		Glyphs:      FireGlyphs,
		Consumables: FullFireConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
		Rotation:    core.GetAplRotation("../../../ui/mage/fire/apls", "fire"),

		ItemFilter: ItemFilter,
	}))
}

var FireTalents = "111122"
var FireGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfInfernoBlast),
	Major2: int32(proto.MageMajorGlyph_GlyphOfCombustion),
}

var PlayerOptionsFire = &proto.Player_FireMage{
	FireMage: &proto.FireMage{
		Options: &proto.FireMage_Options{
			ClassOptions: &proto.MageOptions{},
		},
	},
}
var FullFireConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76085, // Flask of the Warm Sun
	FoodId:   74650, // Mogu Fish Stew
	PotId:    76093, // Potion of the Jade Serpent
	PrepotId: 76093, // Potion of the Jade Serpent
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeCloth,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeStaff,
	},
	HandTypes: []proto.HandType{
		proto.HandType_HandTypeOffHand,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}
