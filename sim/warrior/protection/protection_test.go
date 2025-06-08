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

		GearSet: core.GetGearSet("../../../ui/warrior/protection/gear_sets", "p4_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/warrior/protection/gear_sets", "p3_bis"),
		},
		Talents:     DefaultTalents,
		Glyphs:      DefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.GetAplRotation("../../../ui/warrior/protection/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	HandTypes: []proto.HandType{
		proto.HandType_HandTypeMainHand,
		proto.HandType_HandTypeOneHand,
	},

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

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58085, // Flask of Steelskin
	FoodId:   62670, // Beer‑Basted Crocolisk
	PotId:    58090, // Earthen Potion
	PrepotId: 58090, // Earthen Potion
}
