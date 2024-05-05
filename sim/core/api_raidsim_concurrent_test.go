package core_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/death_knight/blood"
	"github.com/wowsims/cata/sim/hunter/marksmanship"
)

// MM hunter
func getTestPlayer1() *proto.Player {
	var FullConsumes = &proto.Consumes{
		Flask:         proto.Flask_FlaskOfTheWinds,
		DefaultPotion: proto.Potions_PotionOfTheTolvir,
	}
	var MMTalents = "032002-2302320032120231221-03"

	var MMGlyphs = &proto.Glyphs{
		Prime1: int32(proto.HunterPrimeGlyph_GlyphOfArcaneShot),
		Prime2: int32(proto.HunterPrimeGlyph_GlyphOfRapidFire),
		Prime3: int32(proto.HunterPrimeGlyph_HunterPrimeGlyphNone),
	}
	var FerocityTalents = &proto.HunterPetTalents{
		SerpentSwiftness: 2,
		Dive:             true,
		SpikedCollar:     3,
		Bloodthirsty:     1,
		CullingTheHerd:   3,
		SpidersBite:      3,
		Rabid:            true,
		CallOfTheWild:    true,
		SharkAttack:      2,
	}

	var PlayerOptionsBasic = &proto.Player_MarksmanshipHunter{
		MarksmanshipHunter: &proto.MarksmanshipHunter{
			Options: &proto.MarksmanshipHunter_Options{
				ClassOptions: &proto.HunterOptions{
					PetType:           proto.HunterOptions_Wolf,
					PetTalents:        FerocityTalents,
					PetUptime:         0.9,
					TimeToTrapWeaveMs: 0,
				},
			},
		},
	}

	marksmanship.RegisterMarksmanshipHunter()

	return &proto.Player{
		Race:           proto.Race_RaceOrc,
		Class:          proto.Class_ClassHunter,
		Equipment:      core.GetGearSet("../../ui/hunter/marksmanship/gear_sets", "preraid_mm").GearSet,
		Rotation:       core.GetAplRotation("../../ui/hunter/marksmanship/apls", "mm").Rotation,
		Consumes:       FullConsumes,
		Spec:           PlayerOptionsBasic,
		Glyphs:         MMGlyphs,
		TalentsString:  MMTalents,
		Buffs:          core.FullIndividualBuffs,
		ReactionTimeMs: 100,
	}
}

func getTestPlayer2() *proto.Player {
	var BloodTalents = "03323203132212311321--003"
	var BloodDefaultGlyphs = &proto.Glyphs{
		Prime1: int32(proto.DeathKnightPrimeGlyph_GlyphOfDeathStrike),
		Prime2: int32(proto.DeathKnightPrimeGlyph_GlyphOfHeartStrike),
		Prime3: int32(proto.DeathKnightPrimeGlyph_GlyphOfRuneStrike),
		Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfVampiricBlood),
		Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon),
		Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfBoneShield),
	}

	var PlayerOptionsUnholy = &proto.Player_BloodDeathKnight{
		BloodDeathKnight: &proto.BloodDeathKnight{
			Options: &proto.BloodDeathKnight_Options{
				ClassOptions: &proto.DeathKnightOptions{},
			},
		},
	}

	var FullConsumes = &proto.Consumes{
		Flask:         proto.Flask_FlaskOfTitanicStrength,
		DefaultPotion: proto.Potions_GolembloodPotion,
		PrepopPotion:  proto.Potions_GolembloodPotion,
		Food:          proto.Food_FoodBeerBasedCrocolisk,
	}

	blood.RegisterBloodDeathKnight()

	return &proto.Player{
		Race:           proto.Race_RaceWorgen,
		Class:          proto.Class_ClassDeathKnight,
		Equipment:      core.GetGearSet("../../ui/death_knight/blood/gear_sets", "p1").GearSet,
		Rotation:       core.GetAplRotation("../../ui/death_knight/blood/apls", "p1").Rotation,
		Consumes:       FullConsumes,
		Spec:           PlayerOptionsUnholy,
		Glyphs:         BloodDefaultGlyphs,
		TalentsString:  BloodTalents,
		Buffs:          core.FullIndividualBuffs,
		ReactionTimeMs: 100,
	}
}

func makeTestCase(player *proto.Player) *proto.RaidSimRequest {
	return &proto.RaidSimRequest{
		Raid: core.SinglePlayerRaidProto(
			player,
			core.FullPartyBuffs,
			core.FullRaidBuffs,
			core.FullDebuffs),
		Encounter: &proto.Encounter{
			Duration: 300,
			Targets: []*proto.Target{
				core.NewDefaultTarget(),
			},
		},
		SimOptions: &proto.SimOptions{
			Iterations: 1000,
			IsTest:     true,
			Debug:      false,
			RandomSeed: 123,
		},
	}
}

func TestConcurrentRaidSim(t *testing.T) {
	testCases := []*proto.RaidSimRequest{
		makeTestCase(getTestPlayer1()),
		makeTestCase(getTestPlayer2()),
	}

	for i, rsr := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defaultResult := core.RunRaidSim(rsr)
			concurrentResult := core.RunConcurrentRaidSimSync(rsr)

			dpsDefault := defaultResult.RaidMetrics.Dps.Avg
			dpsConcurrent := concurrentResult.RaidMetrics.Dps.Avg

			diff := math.Abs(dpsDefault - dpsConcurrent)

			if diff > 0.0001 {
				t.Logf("DPS expected %0.03f but was %0.03f!", dpsDefault, dpsConcurrent)
				t.Fail()
			}
		})
	}
}
