package balance

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get caster sets included. (we use spellfire here)
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterBalanceDruid()
}

func TestBalance(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassDruid,
		Race:  proto.Race_RaceNightElf,

		GearSet: core.GetGearSet("../../../ui/druid/balance/gear_sets", "preraid"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/druid/balance/gear_sets", "t14"),
		},
		Talents:        StandardTalents,
		Glyphs:         StandardGlyphs,
		Consumables:    FullConsumesSpec,
		SpecOptions:    core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBalance},
		Rotation:       core.GetAplRotation("../../../ui/druid/balance/apls", "standard"),
		OtherRotations: []core.RotationCombo{},
		ItemFilter:     ItemFilter,
	}))
}

var StandardTalents = "113221"
var StandardGlyphs = &proto.Glyphs{}

var PlayerOptionsBalance = &proto.Player_BalanceDruid{
	BalanceDruid: &proto.BalanceDruid{
		Options: &proto.BalanceDruid_Options{
			ClassOptions: &proto.DruidOptions{},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76085, // Flask of the Warm Sun
	FoodId:   74650, // Mogu Fish Stew
	PotId:    76093, // Potion of the Jade Serpent
	PrepotId: 76093, // Potion of the Jade Serpent
}

var ItemFilter = core.ItemFilter{
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
