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

		GearSet:     core.GetGearSet("../../../ui/warrior/fury/gear_sets", "p4_fury_tg"),
		ItemSwapSet: core.GetItemSwapGearSet("../../../ui/warrior/fury/gear_sets", "p4_fury_tg_item_swap"),

		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/fury/gear_sets", "p4_fury_smf"),
		},
		OtherItemSwapSets: []core.ItemSwapSetCombo{
			core.GetItemSwapGearSet("../../../ui/warrior/fury/gear_sets", "p4_fury_smf_item_swap"),
		},
		Talents: SMFTalents,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "Titan's Grip",
				Talents: TGTalents,
				Glyphs:  FuryGlyphs,
			},
		},
		Glyphs:      FuryGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFury},
		Rotation:    core.GetAplRotation("../../../ui/warrior/fury/apls", "tg"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/warrior/fury/apls", "smf"),
		},

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

var SMFTalents = "302003-032222031301101223201-2"
var TGTalents = "302003-03222203130110122321-2"
var FuryGlyphs = &proto.Glyphs{
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfColossusSmash),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfCleaving),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfDeathWish),
	Minor1: int32(proto.WarriorMinorGlyph_GlyphOfBerserkerRage),
	Minor2: int32(proto.WarriorMinorGlyph_GlyphOfBattle),
	Minor3: int32(proto.WarriorMinorGlyph_GlyphOfCommand),
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
	FlaskId:  58088, // Flask of Titanic Strength
	FoodId:   62670, // Beer‑Basted Crocolisk
	PotId:    58146, // Golemblood Potion
	PrepotId: 58146, // Golemblood Potion

}
