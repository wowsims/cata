package protection

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterProtectionPaladin()
}

func TestProtection(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPaladin,
		Race:  proto.Race_RaceBloodElf,

		GearSet:     core.GetGearSet("../../../ui/paladin/protection/gear_sets", "p1"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Seal of Insight", SpecOptions: SealOfInsight},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Seal of Righteousness", SpecOptions: SealOfRighteousness},
			{Label: "Seal of Truth", SpecOptions: SealOfTruth},
		},
		Rotation: core.GetAplRotation("../../../ui/paladin/protection/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeShield,
			},
			HandTypes: []proto.HandType{
				proto.HandType_HandTypeMainHand,
				proto.HandType_HandTypeOneHand,
				proto.HandType_HandTypeOffHand,
			},
			ArmorType:         proto.ArmorType_ArmorTypePlate,
			RangedWeaponTypes: []proto.RangedWeaponType{},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:           proto.Race_RaceBloodElf,
				Class:          proto.Class_ClassPaladin,
				Equipment:      core.GetGearSet("../../../ui/paladin/protection/gear_sets", "p1").GearSet,
				Consumables:    FullConsumesSpec,
				Spec:           SealOfInsight,
				Glyphs:         StandardGlyphs,
				TalentsString:  StandardTalents,
				Buffs:          core.FullIndividualBuffs,
				ReactionTimeMs: 100,
				Rotation:       core.GetAplRotation("../../../ui/paladin/protection/apls", "default").Rotation,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration:          300,
			DurationVariation: 30,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: core.AverageDefaultSimTestOptions,
	}

	core.RaidBenchmark(b, rsr)
}

var StandardTalents = "112222"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfFocusedShield),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfTheAlabasterShield),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfDivineProtection),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfFocusedWrath),
}

var SealOfInsight = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: &proto.ProtectionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Insight,
			},
		},
	},
}

var SealOfRighteousness = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: &proto.ProtectionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Righteousness,
			},
		},
	},
}

var SealOfTruth = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: &proto.ProtectionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Truth,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76087, // Flask of the Earth
	FoodId:   74656, // Chun Tian Spring Rolls
	PotId:    76095, // Potion of Mogu Power
	PrepotId: 76095, // Potion of Mogu Power
}
