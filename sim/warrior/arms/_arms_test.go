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
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p4_arms_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p3_arms_bis"),
		},
		Talents:     ArmsTalents,
		Glyphs:      ArmsDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsArms},
		Rotation:    core.GetAplRotation("../../../ui/warrior/arms/apls", "arms"),

		ItemFilter: ItemFilter,
	}))
}

var ArmsTalents = "32120303120212312201-0322-3"
var ArmsDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfColossusSmash),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfShieldWall),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfRapidCharge),
	Minor1: int32(proto.WarriorMinorGlyph_GlyphOfBerserkerRage),
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
	FlaskId:  58088, // Flask of Titanic Strength
	FoodId:   62670, // Beer‑Basted Crocolisk
	PotId:    58146, // Golemblood Potion
	PrepotId: 58146, // Golemblood Potion

}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
}
