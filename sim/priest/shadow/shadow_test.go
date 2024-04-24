package shadow

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get caster sets included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceNightElf, proto.Race_RaceDraenei},

		GearSet:  core.GetGearSet("../../../ui/priest/shadow/gear_sets", "p1"),
		Talents:  DefaultTalents,
		Glyphs:   DefaultGlyphs,
		Consumes: FullConsumes,

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

		Rotation: core.GetAplRotation("../../../ui/priest/shadow/apls", "default"),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeWand,
			},
		},
	}))
}

var DefaultTalents = "032212--322032210201222100231"
var DefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.PriestPrimeGlyph_GlyphOfMindFlay),
	Prime2: int32(proto.PriestPrimeGlyph_GlyphOfShadowWordPain),
	Prime3: int32(proto.PriestPrimeGlyph_GlyphOfMindFlay),
	Major1: 0,
	Major2: 0,
	Major3: 0,
	Minor1: 0,
	Minor2: 0,
	Minor3: 0,
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfTheDraconicMind,
	Food:            proto.Food_FoodSeafoodFeast,
	DefaultPotion:   proto.Potions_VolcanicPotion,
	PrepopPotion:    proto.Potions_VolcanicPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}

var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			ClassOptions: &proto.PriestOptions{
				Armor: proto.PriestOptions_InnerFire,
			},
		},
	},
}
