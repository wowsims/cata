package blood

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterBloodDeathKnight()
}

func TestBlood(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceWorgen,
		OtherRaces: []proto.Race{proto.Race_RaceGoblin},

		GearSet:  core.GetGearSet("../../../ui/death_knight/blood/gear_sets", "p1"),
		Talents:  BloodTalents,
		Glyphs:   BloodDefaultGlyphs,
		Consumes: FullConsumes,

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/blood/apls", "simple"),

		ItemFilter: ItemFilter,
	}))
}

var BloodTalents = "03323203132212311321--003"
var BloodDefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.DeathKnightPrimeGlyph_GlyphOfDeathStrike),
	Prime2: int32(proto.DeathKnightPrimeGlyph_GlyphOfHeartStrike),
	Prime3: int32(proto.DeathKnightPrimeGlyph_GlyphOfRuneStrike),
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfVampiricBlood),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfBoneShield),
}

var PlayerOptionsUnholy = &proto.Player_BloodDeathKnight{
	BloodDeathKnight: &proto.BloodDeathKnight{
		Options: &proto.BloodDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{},
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
