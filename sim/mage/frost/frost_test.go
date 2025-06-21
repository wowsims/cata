package frost

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFrostMage()
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/mage/frost/gear_sets", "p1_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/mage/frost/gear_sets", "p1_prebis_rich"),
		},
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumables: DefaultConsumables,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/mage/frost/apls", "frost"),

		ItemFilter: ItemFilter,
	}))
}

var FrostTalents = "111122"
var FrostDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfIcyVeins),
	Major2: int32(proto.MageMajorGlyph_GlyphOfSplittingIce),
}

var PlayerOptionsFrost = &proto.Player_FrostMage{
	FrostMage: &proto.FrostMage{
		Options: &proto.FrostMage_Options{
			ClassOptions: &proto.MageOptions{},
		},
	},
}

var DefaultConsumables = &proto.ConsumesSpec{
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
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
	},
}
