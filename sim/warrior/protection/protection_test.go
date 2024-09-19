package protection

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterProtectionWarrior()
}

func TestProtectionWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet: core.GetGearSet("../../../ui/warrior/protection/gear_sets", "p1_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/protection/gear_sets", "p3_bis"),
			core.GetGearSet("../../../ui/warrior/protection/gear_sets", "preraid"),
		},
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.GetAplRotation("../../../ui/warrior/protection/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

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
		proto.WeaponType_WeaponTypeShield,
	},
}

var DefaultTalents = "320003-002-33213201121210212031"
var DefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.WarriorPrimeGlyph_GlyphOfRevenge),
	Prime2: int32(proto.WarriorPrimeGlyph_GlyphOfShieldSlam),
	Prime3: int32(proto.WarriorPrimeGlyph_GlyphOfDevastate),
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfShieldWall),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfShockwave),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfThunderClap),
	Minor1: int32(proto.WarriorMinorGlyph_GlyphOfDemoralizingShout),
	Minor2: int32(proto.WarriorMinorGlyph_GlyphOfBattle),
	Minor3: int32(proto.WarriorMinorGlyph_GlyphOfCommand),
}

var PlayerOptionsBasic = &proto.Player_ProtectionWarrior{
	ProtectionWarrior: &proto.ProtectionWarrior{
		Options: &proto.ProtectionWarrior_Options{
			ClassOptions: &proto.WarriorOptions{
				StartingRage: 0,
			},
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfSteelskin,
	DefaultPotion: proto.Potions_EarthenPotion,
	PrepopPotion:  proto.Potions_EarthenPotion,
	Food:          proto.Food_FoodBeerBasedCrocolisk,
}
