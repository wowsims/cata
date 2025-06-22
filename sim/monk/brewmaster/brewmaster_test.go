package brewmaster

import (
	"testing"

	"github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBrewmasterMonk()
	common.RegisterAllEffects()
}

func TestBrewmaster(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMonk,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/monk/brewmaster/gear_sets", "p1_bis_balanced_2h"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/monk/brewmaster/gear_sets", "p1_bis_balanced_dw"),
			core.GetGearSet("../../../ui/monk/brewmaster/gear_sets", "p1_prebis_rich"),
			core.GetGearSet("../../../ui/monk/brewmaster/gear_sets", "p1_prebis_poor"),
		},
		Talents: BrewmasterDefaultTalents,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "Dungeon",
				Talents: BrewmasterDungeonTalents,
				Glyphs:  BrewmasterDefaultGlyphs,
			},
		},
		Glyphs:      BrewmasterDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBrewmaster},
		Rotation:    core.GetAplRotation("../../../ui/monk/brewmaster/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: ItemFilter,
	}))
}

var BrewmasterDefaultTalents = "213322"
var BrewmasterDungeonTalents = "213321"
var BrewmasterDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.MonkMajorGlyph_GlyphOfFortifyingBrew),
	Major2: int32(proto.MonkMajorGlyph_GlyphOfEnduringHealingSphere),
	Major3: int32(proto.MonkMajorGlyph_GlyphOfFortuitousSpheres),
}

var PlayerOptionsBrewmaster = &proto.Player_BrewmasterMonk{
	BrewmasterMonk: &proto.BrewmasterMonk{
		Options: &proto.BrewmasterMonk_Options{
			ClassOptions: &proto.MonkOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76084, // Flask of Spring Blossoms
	FoodId:   74648, // Sea Mist Rice Noodles
	PotId:    76089, // Virmen's Bite
	PrepotId: 76089, // Virmen's Bite
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeLeather,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeFist,
	},
}
