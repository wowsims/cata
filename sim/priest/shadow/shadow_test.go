package shadow

import (
	"testing"

	"github.com/wowsims/mop/sim/common" // imported to get caster sets included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterShadowPriest()
	common.RegisterAllEffects()
}

func TestShadow(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassPriest,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceNightElf, proto.Race_RaceDraenei},

		GearSet: core.GetGearSet("../../../ui/priest/shadow/gear_sets", "pre_raid"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/priest/shadow/gear_sets", "p1"),
		},
		Talents:     DefaultTalents,
		Glyphs:      &proto.Glyphs{},
		Consumables: FullConsumesSpec,

		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBasic},

		Rotation: core.GetAplRotation("../../../ui/priest/shadow/apls", "default"),

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType: proto.ArmorType_ArmorTypeCloth,
		},
	}))
}

var DefaultTalents = "223113"

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76085, // Flask of the Warm Sun
	FoodId:   74650, // Mogu Fish Stew
	PotId:    76093, //Potion of the Jade Serpent
	PrepotId: 76093, // Potion of the Jade Serpent

}
var PlayerOptionsBasic = &proto.Player_ShadowPriest{
	ShadowPriest: &proto.ShadowPriest{
		Options: &proto.ShadowPriest_Options{
			ClassOptions: &proto.PriestOptions{
				Armor: proto.PriestOptions_InnerFire,
			},
		},
	},
}
