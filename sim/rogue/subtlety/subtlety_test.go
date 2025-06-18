package subtlety

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterSubtletyRogue()
}

func TestSubtlety(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:         proto.Class_ClassRogue,
		Race:          proto.Race_RaceHuman,
		OtherRaces:    []proto.Race{proto.Race_RaceOrc},
		GearSet:       core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "preraid_subtlety"),
		OtherGearSets: []core.GearSetCombo{
			//core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "p3_subtlety"),
			//core.GetGearSet("../../../ui/rogue/subtlety/gear_sets", "p4_subtlety"),
		},
		Talents:        SubtletyTalents,
		Glyphs:         SubtletyGlyphs,
		Consumables:    FullConsumesSpec,
		SpecOptions:    core.SpecOptionsCombo{Label: "Subtlety", SpecOptions: PlayerOptions},
		Rotation:       core.GetAplRotation("../../../ui/rogue/subtlety/apls", "subtlety"),
		OtherRotations: []core.RotationCombo{},
		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
			},
		},
	}))
}

var SubtletyTalents = "321233"

var SubtletyGlyphs = &proto.Glyphs{}

var PlayerOptions = &proto.Player_SubtletyRogue{
	SubtletyRogue: &proto.SubtletyRogue{
		Options: &proto.SubtletyRogue_Options{
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
