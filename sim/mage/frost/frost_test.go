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

		GearSet:     core.GetGearSet("../../../ui/mage/frost/gear_sets", "p1_frost_prebis"),
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/mage/frost/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: ItemFilter,
	}))
}

var FrostTalents = "111121"
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

var DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, // Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
	tinkerId: 82174, // Synapse Springs
});

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeCloth,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeStaff,
	},
}
