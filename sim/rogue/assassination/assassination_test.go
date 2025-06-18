package assassination

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterAssassinationRogue()
}

func TestAssassination(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassRogue,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},
		GearSet:    core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "preraid_assassination"),

		OtherGearSets: []core.GearSetCombo{
			//core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p3_assassination"),
			//core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p4_assassination"),
		},

		Talents:     AssassinationTalents,
		Glyphs:      AssassinationGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Assassination", SpecOptions: PlayerOptionsAssassination},

		Rotation:       core.GetAplRotation("../../../ui/rogue/assassination/apls", "assassination"),
		OtherRotations: []core.RotationCombo{},

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
			},
		},

		// General practice is to not include stat weights in test suite configs to speed up test execution, but at least one spec should
		// include them so the core functionality is tested. Assassination Rogue was chosen because it was
		StatsToWeigh:    []proto.Stat{proto.Stat_StatCritRating},
		EPReferenceStat: proto.Stat_StatAgility,
	}))
}

var AssassinationTalents = "321232"

var AssassinationGlyphs = &proto.Glyphs{}

var PlayerOptionsAssassination = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				LethalPoison: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId: 76084, // Flask of Spring Blossoms
	PotId:   76089, // Virmen's Bite
}
