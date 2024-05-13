package fury

import (
	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.

	"testing"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterFuryWarrior()
}

func TestFury(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		Talents:     FuryTalents,
		Glyphs:      FuryGlyphs,
		GearSet:     core.GetGearSet("../../../ui/warrior/fury/gear_sets", "p1_fury_smf"),
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFury},
		Rotation:    core.GetAplRotation("../../../ui/warrior/fury/apls", "fury"),

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypePlate,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
			},
		},
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

var FuryTalents = "302203-032222031301101223201"
var FuryGlyphs = &proto.Glyphs{
	Prime1: int32(proto.WarriorPrimeGlyph_GlyphOfBloodthirst),
	Prime2: int32(proto.WarriorPrimeGlyph_GlyphOfRagingBlow),
	Prime3: int32(proto.WarriorPrimeGlyph_GlyphOfSlam),
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
				StartingRage:       50,
				UseShatteringThrow: true,
				Shout:              proto.WarriorShout_WarriorShoutBattle,
			},
			UseRecklessness: true,
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTitanicStrength,
	DefaultPotion: proto.Potions_GolembloodPotion,
	PrepopPotion:  proto.Potions_GolembloodPotion,
	Food:          proto.Food_FoodBeerBasedCrocolisk,
}
