package core_test

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/death_knight/blood"
	"github.com/wowsims/cata/sim/druid/feral"
	"github.com/wowsims/cata/sim/hunter/marksmanship"
)

func getTestPlayerMM() *proto.Player {
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

func getTestPlayerBloodDk() *proto.Player {
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

func getTestPlayerFeralCat() *proto.Player {
	var StandardTalents = "-2320322312012121202301-020301"
	var StandardGlyphs = &proto.Glyphs{
		Prime1: int32(proto.DruidPrimeGlyph_GlyphOfRip),
		Prime2: int32(proto.DruidPrimeGlyph_GlyphOfBloodletting),
		Prime3: int32(proto.DruidPrimeGlyph_GlyphOfBerserk),
		Major1: int32(proto.DruidMajorGlyph_GlyphOfThorns),
		Major2: int32(proto.DruidMajorGlyph_GlyphOfFeralCharge),
		Major3: int32(proto.DruidMajorGlyph_GlyphOfRebirth),
	}

	var PlayerOptionsMonoCat = &proto.Player_FeralDruid{
		FeralDruid: &proto.FeralDruid{
			Options: &proto.FeralDruid_Options{
				AssumeBleedActive: true,
			},
		},
	}

	var FullConsumes = &proto.Consumes{
		Flask:         proto.Flask_FlaskOfTheWinds,
		Food:          proto.Food_FoodSkeweredEel,
		DefaultPotion: proto.Potions_PotionOfTheTolvir,
		PrepopPotion:  proto.Potions_PotionOfTheTolvir,
	}

	feral.RegisterFeralDruid()

	return &proto.Player{
		Race:           proto.Race_RaceTauren,
		Class:          proto.Class_ClassDruid,
		Equipment:      core.GetGearSet("../../ui/druid/feral/gear_sets", "preraid").GearSet,
		Rotation:       core.GetAplRotation("../../ui/druid/feral/apls", "default").Rotation,
		Consumes:       FullConsumes,
		Spec:           PlayerOptionsMonoCat,
		Glyphs:         StandardGlyphs,
		TalentsString:  StandardTalents,
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
			Iterations: 200,
			IsTest:     true,
			Debug:      false,
			RandomSeed: 123,
		},
	}
}

func compareValue(t *testing.T, loc string, vst reflect.Value, vmt reflect.Value) {
	switch vst.Kind() {
	case reflect.Pointer, reflect.Interface:
		if vst.IsNil() && vmt.IsNil() {
			break
		}
		if vst.IsNil() != vmt.IsNil() {
			t.Logf("%s: Expected %v but is %v in multi threaded result!", loc, vst.IsNil(), vmt.IsNil())
			t.Fail()
			break
		}
		compareValue(t, loc, vst.Elem(), vmt.Elem())
	case reflect.Struct:
		compareStruct(t, loc, vst, vmt)
	case reflect.Int32, reflect.Int, reflect.Int64:
		if vst.Int() != vmt.Int() {
			t.Logf("%s: Expected %d but is %d for multi threaded result!", loc, vst.Int(), vmt.Int())
			t.Fail()
		}
	case reflect.Float64:
		tolerance := 0.00001
		if strings.Contains(loc, "CastTimeMs") {
			tolerance = 2.2 // Castime is rounded in results and may be off 1ms per thread
		} else if strings.Contains(loc, "Resources") {
			tolerance = 0.001 // Seems to do some rounding at some point?
		}
		if math.Abs(vst.Float()-vmt.Float()) > tolerance {
			t.Logf("%s: Expected %f but is %f for multi threaded result!", loc, vst.Float(), vmt.Float())
			t.Fail()
		}
	case reflect.String:
		if vst.String() != vmt.String() {
			t.Logf("%s: Expected %s but is %s for multi threaded result!", loc, vst.String(), vmt.String())
			t.Fail()
		}
	case reflect.Bool:
		if vst.Bool() != vmt.Bool() {
			t.Logf("%s: Expected %t but is %t for multi threaded result!", loc, vst.Bool(), vmt.Bool())
			t.Fail()
		}
	case reflect.Slice, reflect.Array:
		if vst.Len() != vmt.Len() {
			t.Logf("%s: Expected length %d but is %d for multi threaded result!", loc, vst.Len(), vmt.Len())
			t.Fail()
			break
		}
		for i := 0; i < vst.Len(); i++ {
			compareValue(t, fmt.Sprintf("%s[%d]", loc, i), vst.Index(i), vmt.Index(i))
		}
	case reflect.Map:
		if vst.Len() != vmt.Len() {
			t.Logf("%s: Expected length %d but is %d for multi threaded result!", loc, vst.Len(), vmt.Len())
			t.Fail()
			break
		}
		for _, key := range vst.MapKeys() {
			mtVal := vmt.MapIndex(key)
			keyStr := ""
			switch key.Kind() {
			case reflect.Int32, reflect.Int, reflect.Int64:
				keyStr = fmt.Sprintf("%d", key.Int())
			default:
				keyStr = key.String()
			}
			if !mtVal.IsValid() {
				t.Logf("%s: Key %v not found in multi threaded result!", loc, keyStr)
				t.Fail()
				break
			}
			compareValue(t, fmt.Sprintf("%s[%s]", loc, keyStr), vst.MapIndex(key), mtVal)
		}
	default:
		t.Logf("%s: Has unhandled kind %s!", loc, vst.Kind().String())
		t.Fail()
	}
}

func checkActionMetrics(t *testing.T, loc string, st []*proto.ActionMetrics, mt []*proto.ActionMetrics) {
	actions := map[string]*proto.ActionMetrics{}

	for _, mtAction := range mt {
		_, exists := actions[mtAction.Id.String()]
		if exists {
			t.Logf("%s.Actions: %s exists multiple times in multi threaded results!", loc, mtAction.Id.String())
			t.Fail()
			continue
		}
		actions[mtAction.Id.String()] = mtAction
	}

	for _, stAction := range st {
		mtAction, exists := actions[stAction.Id.String()]
		if !exists {
			t.Logf("%s.Actions: %s does not exist in multi threaded results!", loc, mtAction.Id.String())
			t.Fail()
			continue
		}

		if stAction.IsMelee != mtAction.IsMelee {
			t.Logf("%s.Actions: %s expected IsMelee = %t but was %t in multi threaded results!", loc, stAction.Id.String(), stAction.IsMelee, mtAction.IsMelee)
			t.Fail()
			continue
		}

		compareValue(t, fmt.Sprintf("%s.Actions[%s]", loc, stAction.Id.String()), reflect.ValueOf(stAction.Targets), reflect.ValueOf(mtAction.Targets))
	}
}

func checkResourceMetrics(t *testing.T, loc string, st []*proto.ResourceMetrics, mt []*proto.ResourceMetrics) {
	resources := map[string]*proto.ResourceMetrics{}

	rkey := func(r *proto.ResourceMetrics) string {
		return fmt.Sprintf("%s-%d", r.Id.String(), r.Type)
	}

	for _, mtResource := range mt {
		key := rkey(mtResource)
		_, exists := resources[key]
		if exists {
			t.Logf("%s.Resources: %v exists multiple times in multi threaded results!", loc, key)
			t.Fail()
			continue
		}
		resources[key] = mtResource
	}

	for _, stResource := range st {
		stKey := rkey(stResource)
		mtResource, exists := resources[stKey]
		if !exists {
			t.Logf("%s.Resources: %s does not exist in multi threaded results!", loc, stKey)
			t.Fail()
			continue
		}

		compareValue(t, fmt.Sprintf("%s.Resources[%s]", loc, stKey), reflect.ValueOf(stResource), reflect.ValueOf(mtResource))
	}
}

func compareStruct(t *testing.T, loc string, vst reflect.Value, vmt reflect.Value) {
	for i := 0; i < vst.NumField(); i++ {
		fieldName := vst.Type().Field(i).Name
		fieldType := vst.Type().Field(i).Type.Name()

		if fieldType == "MessageState" {
			continue
		}

		stField := vst.Field(i)
		mtField := vmt.Field(i)

		if stField.Kind() == reflect.Ptr {
			if stField.IsNil() && mtField.IsNil() {
				continue
			} else if stField.IsNil() != mtField.IsNil() {
				t.Logf("%s.%s: Expected %v but is %v in multi threaded result!", loc, fieldName, stField.IsNil(), mtField.IsNil())
				t.Fail()
				continue
			}

			stField = stField.Elem()
			mtField = mtField.Elem()
		}

		if fieldName == "Actions" {
			checkActionMetrics(t, loc, stField.Interface().([]*proto.ActionMetrics), mtField.Interface().([]*proto.ActionMetrics))
			continue
		} else if fieldName == "Resources" {
			checkResourceMetrics(t, loc, stField.Interface().([]*proto.ResourceMetrics), mtField.Interface().([]*proto.ResourceMetrics))
			continue
		}

		compareValue(t, fmt.Sprintf("%s.%s", loc, fieldName), stField, mtField)
	}
}

func TestConcurrentRaidSim(t *testing.T) {
	testCases := []*proto.RaidSimRequest{
		makeTestCase(getTestPlayerMM()),
		makeTestCase(getTestPlayerBloodDk()),
		makeTestCase(getTestPlayerFeralCat()),
	}

	for i, rsr := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stRes := core.RunRaidSim(rsr)
			mtRes := core.RunConcurrentRaidSimSync(rsr)
			vst := reflect.ValueOf(stRes).Elem()
			vmt := reflect.ValueOf(mtRes).Elem()
			compareStruct(t, "RaidSimResult", vst, vmt)
		})
	}

	if t.Failed() {
		t.Log("A fail here means that either the combination of results is broken, or there's a state leak between iterations!")
	}
}
