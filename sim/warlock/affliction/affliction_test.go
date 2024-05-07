package affliction

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterAfflictionWarlock()
}

func TestAffliction(t *testing.T) {
	var defaultAfflictionWarlock = &proto.Player_AfflictionWarlock{
		AfflictionWarlock: &proto.AfflictionWarlock{
			Options: &proto.AfflictionWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_Felhunter,
					DetonateSeed: false,
				},
			},
		},
	}

	var itemFilter = core.ItemFilter{
		WeaponTypes: []proto.WeaponType{
			proto.WeaponType_WeaponTypeSword,
			proto.WeaponType_WeaponTypeDagger,
			proto.WeaponType_WeaponTypeStaff,
		},
		HandTypes: []proto.HandType{
			proto.HandType_HandTypeOffHand,
		},
		ArmorType: proto.ArmorType_ArmorTypeCloth,
		RangedWeaponTypes: []proto.RangedWeaponType{
			proto.RangedWeaponType_RangedWeaponTypeWand,
		},
	}

	var fullConsumes = &proto.Consumes{
		Flask:             proto.Flask_FlaskOfTheDraconicMind,
		Food:              proto.Food_FoodSeveredSagefish,
		DefaultPotion:     proto.Potions_VolcanicPotion,
		ExplosiveBigDaddy: true,
		TinkerHands:       proto.TinkerHands_TinkerHandsSynapseSprings,
	}

	var afflictionTalents = "223222003013321321-03-33"
	var afflictionGlyphs = &proto.Glyphs{
		Prime1: int32(proto.WarlockPrimeGlyph_GlyphOfHaunt),
		Prime2: int32(proto.WarlockPrimeGlyph_GlyphOfUnstableAffliction),
		Prime3: int32(proto.WarlockPrimeGlyph_GlyphOfCorruption),
		Major2: int32(proto.WarlockMajorGlyph_GlyphOfShadowBolt),
		Major1: int32(proto.WarlockMajorGlyph_GlyphOfLifeTap),
		Major3: int32(proto.WarlockMajorGlyph_GlyphOfSoulSwap),
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarlock,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:          core.GetGearSet("../../../ui/warlock/affliction/gear_sets", "p1"),
		Talents:          afflictionTalents,
		Glyphs:           afflictionGlyphs,
		Consumes:         fullConsumes,
		SpecOptions:      core.SpecOptionsCombo{Label: "Affliction Warlock", SpecOptions: defaultAfflictionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/affliction/apls", "default"),
		OtherRotations:   []core.RotationCombo{},
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
