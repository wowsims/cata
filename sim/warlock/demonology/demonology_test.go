package demonology

import (
	"testing"

	"github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterDemonologyWarlock()
	common.RegisterAllEffects()
}

func TestDemonology(t *testing.T) {
	var defaultDemonologyWarlock = &proto.Player_DemonologyWarlock{
		DemonologyWarlock: &proto.DemonologyWarlock{
			Options: &proto.DemonologyWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon: proto.WarlockOptions_Felguard,
				},
			},
		},
	}

	var itemFilter = core.ItemFilter{
		WeaponTypes: []proto.WeaponType{
			proto.WeaponType_WeaponTypeSword,
			proto.WeaponType_WeaponTypeDagger,
			proto.WeaponType_WeaponTypeStaff,
		},
		HandTypes: []proto.HandType{
			proto.HandType_HandTypeOffHand,
		},
		ArmorType: proto.ArmorType_ArmorTypeCloth,
	}

	var fullConsumesSpec = &proto.ConsumesSpec{
		FlaskId:  76085, // Flask of the Warm Sun
		FoodId:   74650, // Mogu Fish Stew
		PotId:    76093, //Potion of the Jade Serpent
		PrepotId: 76093, // Potion of the Jade Serpent
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:            proto.Class_ClassWarlock,
		Race:             proto.Race_RaceOrc,
		OtherRaces:       []proto.Race{proto.Race_RaceTroll, proto.Race_RaceGoblin, proto.Race_RaceHuman},
		GearSet:          core.GetGearSet("../../../ui/warlock/demonology/gear_sets", "preraid"),
		Talents:          "231211",
		Glyphs:           &proto.Glyphs{},
		Consumables:      fullConsumesSpec,
		SpecOptions:      core.SpecOptionsCombo{Label: "Demonology Warlock", SpecOptions: defaultDemonologyWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/demonology/apls", "default"),
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
