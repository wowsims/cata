package beast_mastery

import (
	"testing"

	"github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBeastMasteryHunter()
	common.RegisterAllEffects()
}

func TestBeastMastery(t *testing.T) {
	var talentSets []core.TalentsCombo
	talentSets = core.GenerateTalentVariationsForRows(BeastMasteryTalents, BeastMasteryDefaultGlyphs, []int{4, 5})

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet: core.GetGearSet("../../../ui/hunter/presets", "p1"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/hunter/presets", "preraid"),
			core.GetGearSet("../../../ui/hunter/presets", "preraid_celestial"),
		},
		Talents:         BeastMasteryTalents,
		OtherTalentSets: talentSets,
		Glyphs:          BeastMasteryDefaultGlyphs,
		Consumables:     FullConsumesSpec,
		SpecOptions:     core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:        core.GetAplRotation("../../../ui/hunter/beast_mastery/apls", "bm"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/hunter/beast_mastery/apls", "aoe"),
		},

		ItemFilter:       ItemFilter,
		StartingDistance: 5.1,
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
				Glyphs:         BeastMasteryDefaultGlyphs,
				TalentsString:  BeastMasteryTalents,
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

var BeastMasteryTalents = "312111"
var BeastMasteryDefaultGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfPathfinding),
	Major2: int32(proto.HunterMajorGlyph_GlyphOfAnimalBond),
	Major3: int32(proto.HunterMajorGlyph_GlyphOfDeterrence),
}

var PlayerOptionsBasic = &proto.Player_BeastMasteryHunter{
	BeastMasteryHunter: &proto.BeastMasteryHunter{
		Options: &proto.BeastMasteryHunter_Options{
			ClassOptions: &proto.HunterOptions{
				PetType:   proto.HunterOptions_Wolf,
				PetUptime: 1,
			},
		},
	},
}
