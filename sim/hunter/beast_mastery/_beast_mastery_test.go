package beast_mastery

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBeastMasteryHunter()
}

func TestBM(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassHunter,
		Race:       proto.Race_RaceOrc,
		OtherRaces: []proto.Race{proto.Race_RaceDwarf},

		GearSet:     core.GetGearSet("../../../ui/hunter/beast_mastery/gear_sets", "preraid_bm"),
		Talents:     BMTalents,
		Glyphs:      BMGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},
		Rotation:    core.GetAplRotation("../../../ui/hunter/beast_mastery/apls", "bm"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/hunter/beast_mastery/apls", "bm_advanced"),
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
				Equipment:      core.GetGearSet("../../../ui/hunter/beast_mastery/gear_sets", "preraid_bm").GearSet,
				Consumes:       FullConsumes,
				Spec:           PlayerOptionsBasic,
				Glyphs:         BMGlyphs,
				TalentsString:  BMTalents,
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

var BMTalents = "2330230311320112121-2302-03"
var BMGlyphs = &proto.Glyphs{
	Major1: int32(proto.HunterMajorGlyph_GlyphOfBestialWrath),
}
var FerocityTalents = &proto.HunterPetTalents{
	SerpentSwiftness: 2,
	Dash:             true,
	SpikedCollar:     3,
	Bloodthirsty:     1,
	CullingTheHerd:   3,
	SpidersBite:      3,
	Rabid:            true,
	CallOfTheWild:    true,
	SharkAttack:      2,
}

var PlayerOptionsBasic = &proto.Player_BeastMasteryHunter{
	BeastMasteryHunter: &proto.BeastMasteryHunter{
		Options: &proto.BeastMasteryHunter_Options{
			ClassOptions: &proto.HunterOptions{
				PetType:    proto.HunterOptions_Wolf,
				PetTalents: FerocityTalents,
				PetUptime:  1,
			},
		},
	},
}
