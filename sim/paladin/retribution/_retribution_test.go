package retribution

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassPaladin,
		Race:  proto.Race_RaceBloodElf,

		GearSet: core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p4_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p2_bis"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p2_with_apparatus"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p2_with_double_passive"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p3_bis"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p3_with_double_passive"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p3_with_on_use"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p4_with_apparatus"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p4_with_on_use"),
		},
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: DefaultOptions},
		Rotation:    core.GetAplRotation("../../../ui/paladin/retribution/apls", "default"),
		ItemSwapSet: core.GetItemSwapGearSet("../../../ui/paladin/retribution/gear_sets", "item_swap_4p_t11"),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeSword,
				proto.WeaponType_WeaponTypePolearm,
				proto.WeaponType_WeaponTypeMace,
			},
			HandTypes: []proto.HandType{
				proto.HandType_HandTypeTwoHand,
			},
			ArmorType: proto.ArmorType_ArmorTypePlate,
			RangedWeaponTypes: []proto.RangedWeaponType{},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:           proto.Race_RaceBloodElf,
				Class:          proto.Class_ClassPaladin,
				Equipment:      core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "p4_bis").GearSet,
				Consumables:    FullConsumesSpec,
				Spec:           DefaultOptions,
				Glyphs:         StandardGlyphs,
				TalentsString:  StandardTalents,
				Buffs:          core.FullIndividualBuffs,
				ReactionTimeMs: 100,
				Rotation:       core.GetAplRotation("../../../ui/paladin/retribution/apls", "default").Rotation,
				ItemSwap:       core.GetItemSwapGearSet("../../../ui/paladin/retribution/gear_sets", "item_swap_4p_t11").ItemSwap,
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

var StandardTalents = "203002-02-23203213211113002311"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfTheAsceticCrusader),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfConsecration),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfMight),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfTruth),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfRighteousness),
}

var DefaultOptions = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Truth,
				Aura: proto.PaladinAura_Retribution,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58088, // Flask of Titanic Strength
	FoodId:   62670, // Beer‑Basted Crocolisk
	PotId:    58146, // Golemblood Potion
	PrepotId: 58146, // Golemblood Potion
	TinkerId: 82174, // Synapse Springs
}
