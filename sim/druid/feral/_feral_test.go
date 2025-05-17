package feral

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFeralDruid()
}

var FeralItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypePolearm,
	},
	ArmorType:         proto.ArmorType_ArmorTypeLeather,
	RangedWeaponTypes: []proto.RangedWeaponType{},
}

func TestFeral(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassDruid,
		Race:        proto.Race_RaceWorgen,
		OtherRaces:  []proto.Race{proto.Race_RaceTroll},
		GearSet:     core.GetGearSet("../../../ui/druid/feral/gear_sets", "p4"),
		ItemSwapSet: core.GetItemSwapGearSet("../../../ui/druid/feral/gear_sets", "p4_item_swap"),

		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/druid/feral/gear_sets", "p3"),
		},

		OtherItemSwapSets: []core.ItemSwapSetCombo{
			{Label: "no_item_swap", ItemSwap: &proto.ItemSwap{}},
		},

		Talents:         StandardTalents,
		Glyphs:          StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{{Label: "HybridTalents", Talents: HybridTalents, Glyphs: HybridGlyphs}},
		Consumables:     FullConsumesSpec,
		SpecOptions:     core.SpecOptionsCombo{Label: "ExternalBleed", SpecOptions: PlayerOptionsMonoCat},
		Rotation:        core.GetAplRotation("../../../ui/druid/feral/apls", "default"),

		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/druid/feral/apls", "monocat"),
			core.GetAplRotation("../../../ui/druid/feral/apls", "aoe"),
		},

		StartingDistance: 25,
		ItemFilter:       FeralItemFilter,
	}))
}

// func TestFeralApl(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class: proto.Class_ClassDruid,
// 		Race:  proto.Race_RaceTauren,

// 		GearSet:     core.GetGearSet("../../../ui/feral_druid/gear_sets", "p3"),
// 		Talents:     StandardTalents,
// 		Glyphs:      StandardGlyphs,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsMonoCat},
// 		Rotation:    core.GetAplRotation("../../../ui/feral_druid/apls", "default"),
// 		ItemFilter:  FeralItemFilter,
// 	}))
// }

// func BenchmarkSimulate(b *testing.B) {
// 	rsr := &proto.RaidSimRequest{
// 		Raid: core.SinglePlayerRaidProto(
// 			&proto.Player{
// 				Race:      proto.Race_RaceTauren,
// 				Class:     proto.Class_ClassDruid,
// 				Equipment: core.GetGearSet("../../../ui/feral_druid/gear_sets", "p1").GearSet,
// 				Consumes:  FullConsumes,
// 				Spec:      PlayerOptionsMonoCat,
// 				Buffs:     core.FullIndividualBuffs,
// 				Glyphs:    StandardGlyphs,

// 				InFrontOfTarget: true,
// 			},
// 			core.FullPartyBuffs,
// 			core.FullRaidBuffs,
// 			core.FullDebuffs),
// 		Encounter: &proto.Encounter{
// 			Duration: 300,
// 			Targets: []*proto.Target{
// 				core.NewDefaultTarget(),
// 			},
// 		},
// 		SimOptions: core.AverageDefaultSimTestOptions,
// 	}

// 	core.RaidBenchmark(b, rsr)
// }

var StandardTalents = "-2320322312012121202301-020301"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfThorns),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfFeralCharge),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfRebirth),
}

var HybridTalents = "-2300322312310001220311-020331"
var HybridGlyphs = &proto.Glyphs{
	Major1: int32(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfMaul),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfRebirth),
}

var PlayerOptionsMonoCat = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			AssumeBleedActive: true,
		},
	},
}

var PlayerOptionsMonoCatNoBleed = &proto.Player_FeralDruid{
	FeralDruid: &proto.FeralDruid{
		Options: &proto.FeralDruid_Options{
			AssumeBleedActive: false,
		},
	},
}

// var PlayerOptionsFlowerCatAoe = &proto.Player_FeralDruid{
// 	FeralDruid: &proto.FeralDruid{
// 		Options: &proto.FeralDruid_Options{
// 			InnervateTarget:   &proto.UnitReference{}, // no Innervate
// 			AssumeBleedActive: false,
// 		},
// 		Rotation: &proto.FeralDruid_Rotation{
// 			RotationType:       proto.FeralDruid_Rotation_Aoe,
// 			BearWeaveType:      proto.FeralDruid_Rotation_None,
// 			UseRake:            true,
// 			UseBite:            true,
// 			MinCombosForRip:    5,
// 			MinCombosForBite:   5,
// 			BiteTime:           4.0,
// 			MaintainFaerieFire: true,
// 			BerserkBiteThresh:  25.0,
// 			BerserkFfThresh:    15.0,
// 			MaxFfDelay:         0.7,
// 			MinRoarOffset:      24.0,
// 			RipLeeway:          3,
// 			SnekWeave:          false,
// 			FlowerWeave:        true,
// 			RaidTargets:        30,
// 			PrePopOoc:          true,
// 		},
// 	},
// }

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:     58087, // Flask of the Winds
	FoodId:      62669, // Skewered Eel
	PotId:       58145, // Potion of the Tol'vir
	PrepotId:    58145, // Potion of the Tol'vir
	ExplosiveId: 89637, // Big Daddy Explosive
}
