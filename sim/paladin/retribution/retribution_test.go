package retribution

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterRetributionPaladin()
}

func TestRetribution(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPaladin,
		Race:       proto.Race_RaceBloodElf,
		OtherRaces: []proto.Race{proto.Race_RaceHuman, proto.Race_RaceDraenei, proto.Race_RaceDwarf, proto.Race_RaceTauren},

		GearSet: core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "t11_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "preraid"),
			core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "t12_bis"),
			//core.GetGearSet("../../../ui/paladin/retribution/gear_sets", "t13_bis"),
		},
		Talents: StandardTalents,
		Glyphs:  StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "Crusader Strike Glyph",
				Talents: StandardTalents,
				Glyphs:  GlyphsWithCrusaderStrike,
			},
		},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: DefaultOptions},
		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "Snapshot Guardian", SpecOptions: OptionsWithSnapshotGuardian},
			{Label: "Seal of Righteousness", SpecOptions: OptionsWithSealOfRighteousness},
			{Label: "Seal of Justice", SpecOptions: OptionsWithSealOfJustice},
			{Label: "Seal of Insight", SpecOptions: OptionsWithSealOfInsight},
		},
		Rotation: core.GetAplRotation("../../../ui/paladin/retribution/apls", "default"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/paladin/retribution/apls", "apparatus"),
			//core.GetAplRotation("../../../ui/paladin/retribution/apls", "t13"),
			//core.GetAplRotation("../../../ui/paladin/retribution/apls", "t13-apparatus"),
		},

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
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeRelic,
			},
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:           proto.Race_RaceBloodElf,
				Class:          proto.Class_ClassPaladin,
				Equipment:      core.GetGearSet("../../../ui/retribution_paladin/gear_sets", "t11_bis").GearSet,
				Consumes:       FullConsumes,
				Spec:           DefaultOptions,
				Glyphs:         StandardGlyphs,
				TalentsString:  StandardTalents,
				Buffs:          core.FullIndividualBuffs,
				ReactionTimeMs: 100,
				Rotation:       core.GetAplRotation("../../../ui/paladin/retribution/apls", "default").Rotation,
			},
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
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
	Prime1: int32(proto.PaladinPrimeGlyph_GlyphOfTemplarSVerdict),
	Prime2: int32(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth),
	Prime3: int32(proto.PaladinPrimeGlyph_GlyphOfExorcism),
	Major1: int32(proto.PaladinMajorGlyph_GlyphOfHammerOfWrath),
	Major2: int32(proto.PaladinMajorGlyph_GlyphOfTheAsceticCrusader),
	Major3: int32(proto.PaladinMajorGlyph_GlyphOfConsecration),
	Minor1: int32(proto.PaladinMinorGlyph_GlyphOfBlessingOfMight),
	Minor2: int32(proto.PaladinMinorGlyph_GlyphOfTruth),
	Minor3: int32(proto.PaladinMinorGlyph_GlyphOfRighteousness),
}
var GlyphsWithCrusaderStrike = &proto.Glyphs{
	Prime1: int32(proto.PaladinPrimeGlyph_GlyphOfTemplarSVerdict),
	Prime2: int32(proto.PaladinPrimeGlyph_GlyphOfSealOfTruth),
	Prime3: int32(proto.PaladinPrimeGlyph_GlyphOfCrusaderStrike),
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
				Seal:             proto.PaladinSeal_Truth,
				Aura:             proto.PaladinAura_Retribution,
				SnapshotGuardian: false,
			},
		},
	},
}

var OptionsWithSnapshotGuardian = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal:             proto.PaladinSeal_Truth,
				Aura:             proto.PaladinAura_Retribution,
				SnapshotGuardian: true,
			},
		},
	},
}

var OptionsWithSealOfRighteousness = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal:             proto.PaladinSeal_Righteousness,
				Aura:             proto.PaladinAura_Retribution,
				SnapshotGuardian: false,
			},
		},
	},
}

var OptionsWithSealOfJustice = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal:             proto.PaladinSeal_Justice,
				Aura:             proto.PaladinAura_Retribution,
				SnapshotGuardian: false,
			},
		},
	},
}

var OptionsWithSealOfInsight = &proto.Player_RetributionPaladin{
	RetributionPaladin: &proto.RetributionPaladin{
		Options: &proto.RetributionPaladin_Options{
			ClassOptions: &proto.PaladinOptions{
				Seal:             proto.PaladinSeal_Insight,
				Aura:             proto.PaladinAura_Retribution,
				SnapshotGuardian: false,
			},
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTitanicStrength,
	DefaultPotion: proto.Potions_GolembloodPotion,
	PrepopPotion:  proto.Potions_GolembloodPotion,
	Food:          proto.Food_FoodBeerBasedCrocolisk,
}
