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
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman},

		GearSet:     core.GetGearSet("../../../ui/paladin/protection/gear_sets", "T12"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: DefaultOptions},
		Rotation:    core.GetAplRotation("../../../ui/paladin/protection/apls", "default"),

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
				Equipment:      core.GetGearSet("../../../ui/paladin/protection/gear_sets", "T12").GearSet,
				Consumables:    FullConsumesSpec,
				Spec:           DefaultOptions,
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

var StandardTalents = "-32023013122121101231-032032"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfTheAsceticCrusader),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfLayOnHands),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfFocusedShield),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfTruth),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfMight),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfInsight),
}

var DefaultOptions = &proto.Player_ProtectionPaladin{
	ProtectionPaladin: &proto.ProtectionPaladin{
		Options: &proto.ProtectionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal: proto.PaladinSeal_Truth,
				Aura: proto.PaladinAura_Retribution,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58085, // Flask of Steelskin
	FoodId:   62663, // Lavascale Minestrone
	PotId:    58146, // Golemblood Potion
	PrepotId: 58146, // Golemblood Potion
	TinkerId: 82174, // Synapse Springs
}
