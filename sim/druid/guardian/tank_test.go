package guardian

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterGuardianDruid()
}

func TestGuardian(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceTauren,

		GearSet: core.GetGearSet("../../../ui/druid/guardian/gear_sets", "preraid"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/druid/guardian/gear_sets", "p1"),
		},

		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsDefault},
		Rotation:    core.GetAplRotation("../../../ui/druid/guardian/apls", "default"),

		IsTank:          true,
		InFrontOfTarget: true,

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
				proto.WeaponType_WeaponTypePolearm,
			},
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeRelic,
			},
		},
	}))
}

// func BenchmarkSimulate(b *testing.B) {
// 	rsr := &proto.RaidSimRequest{
// 		Raid: core.SinglePlayerRaidProto(
// 			&proto.Player{
// 				Race:      proto.Race_RaceTauren,
// 				Class:     proto.Class_ClassDruid,
// 				Equipment: core.GetGearSet("../../../ui/feral_tank_druid/gear_sets", "p1").GearSet,
// 				Consumes:  FullConsumes,
// 				Spec:      PlayerOptionsDefault,
// 				Buffs:     core.FullIndividualBuffs,
//
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
//
// 	core.RaidBenchmark(b, rsr)
// }

var StandardTalents = "-2300322312310001220311-020331"
var StandardGlyphs = &proto.Glyphs{
	Prime1: int32(proto.DruidPrimeGlyph_GlyphOfMangle),
	Prime2: int32(proto.DruidPrimeGlyph_GlyphOfLacerate),
	Prime3: int32(proto.DruidPrimeGlyph_GlyphOfBerserk),
	Major1: int32(proto.DruidMajorGlyph_GlyphOfFrenziedRegeneration),
	Major2: int32(proto.DruidMajorGlyph_GlyphOfMaul),
	Major3: int32(proto.DruidMajorGlyph_GlyphOfRebirth),
}

var PlayerOptionsDefault = &proto.Player_GuardianDruid{
	GuardianDruid: &proto.GuardianDruid{
		Options: &proto.GuardianDruid_Options{
			StartingRage: 15,
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfSteelskin,
	Food:            proto.Food_FoodSkeweredEel,
	DefaultPotion:   proto.Potions_PotionOfTheTolvir,
	PrepopPotion:    proto.Potions_PotionOfTheTolvir,
	DefaultConjured: proto.Conjured_ConjuredHealthstone,
}
