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

		GearSet: core.GetGearSet("../../../ui/priest/shadow/gear_sets", "p4"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/priest/shadow/gear_sets", "p3"),
		},
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumables: FullConsumesSpec,

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

		Rotation: core.GetAplRotation("../../../ui/priest/shadow/apls", "default"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/priest/shadow/apls", "default"),
		},

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
	Prime3: int32(proto.PriestPrimeGlyph_GlyphOfShadowWordDeath),
	Major1: int32(proto.PriestMajorGlyph_GlyphOfFade),
	Major2: int32(proto.PriestMajorGlyph_GlyphOfInnerFire),
	Major3: int32(proto.PriestMajorGlyph_GlyphOfSpiritTap),
	Minor1: int32(proto.PriestMinorGlyph_GlyphOfFading),
	Minor2: int32(proto.PriestMinorGlyph_GlyphOfFortitude),
	Minor3: int32(proto.PriestMinorGlyph_GlyphOfShadowfiend),
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58086, // Flask of the Draconic Mind
	FoodId:   62290, // Seafood Magnifique Feast
	PotId:    58091, // Volcanic Potion
	PrepotId: 58091, // Volcanic Potion
	TinkerId: 82174, // Synapse Springs
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
