package unholy

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterUnholyDeathKnight()
}

func TestUnholy(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceWorgen,
		OtherRaces: []proto.Race{proto.Race_RaceGoblin},

		GearSet:     core.GetGearSet("../../../ui/death_knight/unholy/gear_sets", "p3.bis"),
		Talents:     UnholyTalents,
		Glyphs:      UnholyDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/unholy/apls", "st"),

		ItemFilter: ItemFilter,
	}))
}

var UnholyTalents = "203-1-13300321230331121231"
var UnholyDefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.DeathKnightPrimeGlyph_GlyphOfDeathCoil),
	Prime2: int32(proto.DeathKnightPrimeGlyph_GlyphOfScourgeStrike),
	Prime3: int32(proto.DeathKnightPrimeGlyph_GlyphOfRaiseDead),
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfPestilence),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfBloodBoil),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell),
	Minor1: int32(proto.DeathKnightMinorGlyph_GlyphOfDeathSEmbrace),
	Minor2: int32(proto.DeathKnightMinorGlyph_GlyphOfHornOfWinter),
}

var PlayerOptionsUnholy = &proto.Player_UnholyDeathKnight{
	UnholyDeathKnight: &proto.UnholyDeathKnight{
		Options: &proto.UnholyDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{
				PetUptime:          1.0,
				StartingRunicPower: 100,
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
