package arms

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterArmsWarrior()
}

func TestArms(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarrior,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceWorgen},
		StartingDistance: 9,

		GearSet: core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p1_arms_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p1_prebis_rich"),
			core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p1_prebis_poor"),
		},
		Talents:     ArmsTalents,
		Glyphs:      ArmsDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsArms},
		Rotation:    core.GetAplRotation("../../../ui/warrior/arms/apls", "arms"),

		ItemFilter: ItemFilter,
	}))
}

var ArmsTalents = "113332"
var ArmsDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfBullRush),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfUnendingRage),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfDeathFromAbove),
}

var PlayerOptionsArms = &proto.Player_ArmsWarrior{
	ArmsWarrior: &proto.ArmsWarrior{
		Options: &proto.ArmsWarrior_Options{
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

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
}
