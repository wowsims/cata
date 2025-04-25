package marksmanship

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterMarksmanshipHunter()
}

func TestMM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GetGearSet("../../../ui/hunter/marksmanship/gear_sets", "preraid_mm"),
		Talents:     MMTalents,
		Glyphs:      MMGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.GetAplRotation("../../../ui/hunter/marksmanship/apls", "mm"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/hunter/marksmanship/apls", "mm"),
		},

		ItemFilter:       ItemFilter,
		StartingDistance: 5.1,
	}))
}

var ItemFilter = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypePolearm,
		proto.WeaponType_WeaponTypeStaff,
		proto.WeaponType_WeaponTypeSword,
	},
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
				Equipment:      core.GetGearSet("../../../ui/hunter/marksmanship/gear_sets", "preraid_mm").GearSet,
				Consumes:       FullConsumes,
				Spec:           PlayerOptionsBasic,
				Glyphs:         MMGlyphs,
				TalentsString:  MMTalents,
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

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheWinds,
	DefaultPotion: proto.Potions_PotionOfTheTolvir,
}
var MMTalents = "032002-2302320032120231221-03"

var MMGlyphs = &proto.Glyphs{}
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
