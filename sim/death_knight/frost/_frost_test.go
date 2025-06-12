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

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/death_knight/frost/gear_sets", "p4.masterfrost"),
		Talents: MasterfrostTalents,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "TwoHand",
				Talents: TwoHandTalents,
				Glyphs:  FrostDefaultGlyphs,
			},
			{
				Label:   "DualWield",
				Talents: DualWieldTalents,
				Glyphs:  FrostDefaultGlyphs,
			},
		},
		Glyphs:      FrostDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/frost/apls", "masterfrost"),

		ItemFilter: ItemFilter,
	}))
}

var DualWieldTalents = "2032-20330022233112012301-003"
var TwoHandTalents = "103-32030022233112012031-033"
var MasterfrostTalents = "2032-30330012233112012301-03"

var FrostDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfPestilence),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfBloodBoil),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfDarkSuccor),
	// No interesting minor glyphs.
}

var PlayerOptionsFrost = &proto.Player_FrostDeathKnight{
	FrostDeathKnight: &proto.FrostDeathKnight{
		Options: &proto.FrostDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{
				PetUptime: 1.0,
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

	HandTypes: []proto.HandType{
		proto.HandType_HandTypeMainHand,
		proto.HandType_HandTypeOffHand,
		proto.HandType_HandTypeOneHand,
		proto.HandType_HandTypeTwoHand,
	},
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{},
}
