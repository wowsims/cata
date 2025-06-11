package windwalker

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterWindwalkerMonk()
}

func TestWindwalker(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMonk,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/monk/windwalker/gear_sets", "p1_prebis_dw"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/monk/windwalker/gear_sets", "p1_prebis_2h"),
			core.GetGearSet("../../../ui/monk/windwalker/gear_sets", "p1_bis_dw"),
			core.GetGearSet("../../../ui/monk/windwalker/gear_sets", "p1_bis_2h"),
		},
		Talents:     WindwalkerTalents,
		Glyphs:      WindwalkerDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsWindwalker},
		Rotation:    core.GetAplRotation("../../../ui/monk/windwalker/apls", "default"),

		ItemFilter: ItemFilter,
	}))
}

var WindwalkerTalents = "213322"
var WindwalkerDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.MonkMajorGlyph_GlyphOfSpinningCraneKick),
	Major2: int32(proto.MonkMajorGlyph_GlyphOfTouchOfKarma),
	Minor1: int32(proto.MonkMinorGlyph_GlyphOfBlackoutKick),
}

var PlayerOptionsWindwalker = &proto.Player_WindwalkerMonk{
	WindwalkerMonk: &proto.WindwalkerMonk{
		Options: &proto.WindwalkerMonk_Options{
			ClassOptions: &proto.MonkOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76084,  // Flask of Spring Blossoms
	FoodId:   74648,  // Sea Mist Rice Noodles
	PotId:    76089,  // Virmen's Bite
	PrepotId: 76089,  // Virmen's Bite
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeLeather,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeFist,
	},
}
