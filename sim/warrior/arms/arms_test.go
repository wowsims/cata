package arms

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterArmsWarrior()
}

func TestArms(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassWarrior,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},

		GearSet:     core.GetGearSet("../../../ui/warrior/arms/gear_sets", "p1_arms_bis"),
		Talents:     ArmsTalents,
		Glyphs:      ArmsDefaultGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsArms},
		Rotation:    core.GetAplRotation("../../../ui/warrior/arms/apls", "arms"),

		ItemFilter: ItemFilter,
	}))
}

var ArmsTalents = "32120303120212312201-0322-3"
var ArmsDefaultGlyphs = &proto.Glyphs{
	Prime1: int32(proto.WarriorPrimeGlyph_GlyphOfMortalStrike),
	Prime2: int32(proto.WarriorPrimeGlyph_GlyphOfOverpower),
	Prime3: int32(proto.WarriorPrimeGlyph_GlyphOfSlam),
	Major1: int32(proto.WarriorMajorGlyph_GlyphOfColossusSmash),
	Major2: int32(proto.WarriorMajorGlyph_GlyphOfShieldWall),
	Major3: int32(proto.WarriorMajorGlyph_GlyphOfRapidCharge),
	Minor1: int32(proto.WarriorMinorGlyph_GlyphOfBerserkerRage),
}

var PlayerOptionsArms = &proto.Player_ArmsWarrior{
	ArmsWarrior: &proto.ArmsWarrior{
		Options: &proto.ArmsWarrior_Options{
			ClassOptions: &proto.WarriorOptions{
				StartingRage: 0,
			},
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTitanicStrength,
	DefaultPotion: proto.Potions_GolembloodPotion,
	PrepopPotion:  proto.Potions_GolembloodPotion,
	Food:          proto.Food_FoodBeerBasedCrocolisk,
	TinkerHands:   proto.TinkerHands_TinkerHandsSynapseSprings,
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
	},
}
