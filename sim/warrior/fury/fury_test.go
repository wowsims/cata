package fury

import (
	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.

	"testing"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFuryWarrior()
}

func TestFury(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/warrior/fury/gear_sets", "p1_fury_tg"),

		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/fury/gear_sets", "p1_fury_smf"),
		},
		Talents: TGTalents,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "Single-Minded Fury",
				Talents: SMFTalents,
				Glyphs:  FuryGlyphs,
			},
		},
		Glyphs:           FuryGlyphs,
		Consumables:      FullConsumesSpec,
		SpecOptions:      core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFury},
		Rotation:         core.GetAplRotation("../../../ui/warrior/fury/apls", "default"),
		StartingDistance: 5,

		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
	},
}

var SMFTalents = "133333"
var TGTalents = "133133"
var FuryGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfBullRush),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfDeathFromAbove),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfUnendingRage),
}

var PlayerOptionsFury = &proto.Player_FuryWarrior{
	FuryWarrior: &proto.FuryWarrior{
		Options: &proto.FuryWarrior_Options{
			ClassOptions: &proto.WarriorOptions{
				StartingRage: 0,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76088, // Flask of Winter's Bite
	FoodId:   74646, // Black Pepper Ribs and Shrimp
	PotId:    76095, // Potion of Mogu Power
	PrepotId: 76095, // Potion of Mogu Power
}
