package blood

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBloodDeathKnight()
}

func TestBlood(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassDeathKnight,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet: core.GetGearSet("../../../ui/death_knight/blood/gear_sets", "p3-balanced"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/death_knight/blood/gear_sets", "p1"),
		},

		Talents:     BloodTalents,
		Glyphs:      BloodDefaultGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBlood},
		Rotation:    core.GetAplRotation("../../../ui/death_knight/blood/apls", "simple"),

		InFrontOfTarget: true,
		IsTank:          true,

		ItemFilter: ItemFilter,
	}))
}

var BloodTalents = "02323203102122111321-3-033"
var BloodDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfAntiMagicShell),
	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon),
	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfBoneShield),
	Minor1: int32(proto.DeathKnightMinorGlyph_GlyphOfDeathGate),
	Minor2: int32(proto.DeathKnightMinorGlyph_GlyphOfPathOfFrost),
	Minor3: int32(proto.DeathKnightMinorGlyph_GlyphOfHornOfWinter),
}

var PlayerOptionsBlood = &proto.Player_BloodDeathKnight{
	BloodDeathKnight: &proto.BloodDeathKnight{
		Options: &proto.BloodDeathKnight_Options{
			ClassOptions: &proto.DeathKnightOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58085, // Flask of Steelskin
	FoodId:   62670, // Beer‑Basted Crocolisk
	PotId:    58146, // Golemblood Potion
	PrepotId: 58146, // Golemblood Potion
	TinkerId: 82174, // Synapse Springs
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	HandTypes: []proto.HandType{
		proto.HandType_HandTypeTwoHand,
	},
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
	RangedWeaponTypes: []proto.RangedWeaponType{},
}
