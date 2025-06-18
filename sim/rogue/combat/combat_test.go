package combat

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterCombatRogue()
}

func TestCombat(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:         proto.Class_ClassRogue,
		Race:          proto.Race_RaceHuman,
		OtherRaces:    []proto.Race{proto.Race_RaceOrc},
		GearSet:       core.GetGearSet("../../../ui/rogue/combat/gear_sets", "preraid_combat"),
		OtherGearSets: []core.GearSetCombo{
			//core.GetGearSet("../../../ui/rogue/combat/gear_sets", "p3_combat"),
			//core.GetGearSet("../../../ui/rogue/combat/gear_sets", "p4_combat"),
		},
		Talents:     CombatTalents,
		Glyphs:      CombatGlyphs,
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Combat", SpecOptions: PlayerOptions},

		Rotation:       core.GetAplRotation("../../../ui/rogue/combat/apls", "combat"),
		OtherRotations: []core.RotationCombo{},
		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeSword,
			},
			HandTypes: []proto.HandType{
				proto.HandType_HandTypeMainHand,
				proto.HandType_HandTypeOffHand,
				proto.HandType_HandTypeOneHand,
			},
		},
	}))
}

var CombatTalents = "321233"

var CombatGlyphs = &proto.Glyphs{}

var PlayerOptions = &proto.Player_CombatRogue{
	CombatRogue: &proto.CombatRogue{
		Options: &proto.CombatRogue_Options{
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
