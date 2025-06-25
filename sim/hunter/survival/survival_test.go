package survival

import (
	"testing"

	"github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterSurvivalHunter()
	common.RegisterAllEffects()
}

func TestSurvival(t *testing.T) {
	var talentSets []core.TalentsCombo
	talentSets = core.GenerateTalentVariationsForRows(SurvivalTalents, SurvivalDefaultGlyphs, []int{4, 5})

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet: core.GetGearSet("../../../ui/hunter/presets", "p1"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/hunter/presets", "preraid"),
			core.GetGearSet("../../../ui/hunter/presets", "preraid_celestial"),
		},
		Talents:         SurvivalTalents,
		OtherTalentSets: talentSets,
		Glyphs:          SurvivalDefaultGlyphs,
		Consumables:     FullConsumesSpec,
		SpecOptions:     core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:        core.GetAplRotation("../../../ui/hunter/survival/apls", "sv"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/hunter/survival/apls", "aoe"),
		},
		StartingDistance: 5.1,
		ItemFilter:       ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeBow,
		proto.RangedWeaponType_RangedWeaponTypeCrossbow,
		proto.RangedWeaponType_RangedWeaponTypeGun,
	},
}

func BenchmarkSimulate(b *testing.B) {
	rsr := &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			&proto.Player{
				Race:           proto.Race_RaceOrc,
				Class:          proto.Class_ClassHunter,
				Equipment:      core.GetGearSet("../../../ui/hunter/presets", "p1").GearSet,
				Consumables:    FullConsumesSpec,
				Spec:           PlayerOptionsBasic,
				Glyphs:         SurvivalDefaultGlyphs,
				TalentsString:  SurvivalTalents,
				Buffs:          core.FullIndividualBuffs,
				ReactionTimeMs: 100,
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

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76084, // Flask of Spring Blossoms
	FoodId:   74648, // Sea Mist Rice Noodles
	PotId:    76089, // Virmen's Bite
	PrepotId: 76089, // Virmen's Bite
}

var SurvivalTalents = "312111"
var SurvivalDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfLiberation),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfAnimalBond),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfDeterrence),
}

var PlayerOptionsBasic = &proto.Player_SurvivalHunter{
	SurvivalHunter: &proto.SurvivalHunter{
		Options: &proto.SurvivalHunter_Options{
			ClassOptions: &proto.HunterOptions{
				PetType:   proto.HunterOptions_Wolf,
				PetUptime: 1,
			},
		},
	},
}
