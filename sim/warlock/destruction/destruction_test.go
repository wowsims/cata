package destruction

import (
	"testing"

	"github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterDestructionWarlock()
	common.RegisterAllEffects()
}

func TestDestruction(t *testing.T) {
	var defaultDestructionWarlock = &proto.Player_DestructionWarlock{
		DestructionWarlock: &proto.DestructionWarlock{
			Options: &proto.DestructionWarlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:       proto.WarlockOptions_Imp,
					DetonateSeed: false,
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
		RangedWeaponTypes: []proto.RangedWeaponType{
			proto.RangedWeaponType_RangedWeaponTypeWand,
		},
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
		GearSet:          core.GetGearSet("../../../ui/warlock/destruction/gear_sets", "p1-prebis"),
		Talents:          "221211",
		Glyphs:           &proto.Glyphs{},
		Consumables:      fullConsumesSpec,
		SpecOptions:      core.SpecOptionsCombo{Label: "Destruction Warlock", SpecOptions: defaultDestructionWarlock},
		OtherSpecOptions: []core.SpecOptionsCombo{},
		Rotation:         core.GetAplRotation("../../../ui/warlock/destruction/apls", "default"),
		ItemFilter:       itemFilter,
		StartingDistance: 25,
	}))
}
