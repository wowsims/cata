package frost

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterFrostDeathKnight()
}

func TestFrost(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/death_knight/frost/gear_sets", "p12h"),
		Talents:     FrostTalents,
		Glyphs:      FrostDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/frost/apls", "2hst"),

		ItemFilter: ItemFilter,
	}))
}

var FrostTalents = "103-32030022233112012031-033"
var FrostDefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.DeathKnightPrimeGlyph_GlyphOfFrostStrike),
	Prime2: int32(proto.DeathKnightPrimeGlyph_GlyphOfObliterate),
	Prime3: int32(proto.DeathKnightPrimeGlyph_GlyphOfHowlingBlast),
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

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTitanicStrength,
	DefaultPotion: proto.Potions_GolembloodPotion,
	PrepopPotion:  proto.Potions_GolembloodPotion,
	Food:          proto.Food_FoodBeerBasedCrocolisk,
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
}
