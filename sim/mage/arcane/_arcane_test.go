package arcane

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterArcaneMage()
}

func TestArcane(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class: proto.Class_ClassMage,
		Race:  proto.Race_RaceTroll,

		GearSet:     core.GetGearSet("../../../ui/mage/arcane/gear_sets", "p3"),
		Talents:     ArcaneTalents,
		Glyphs:      ArcaneGlyphs,
		Consumables: FullArcaneConsumesSpec,

		SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},
		Rotation:    core.GetAplRotation("../../../ui/mage/arcane/apls", "arcane"),
		// OtherRotations: []core.RotationCombo{
		// 	core.GetAplRotation("../../ui/mage/apls", "arcane_aoe"),
		// },

		ItemFilter: ItemFilter,
	}))
}

var ItemFilter = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var ArcaneTalents = "303322021230122210121-23-03"
var ArcaneGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfEvocation),
	Major2: int32(proto.MageMajorGlyph_GlyphOfArcanePower),
	Major3: int32(proto.MageMajorGlyph_GlyphOfManaShield),
}

var PlayerOptionsArcane = &proto.Player_ArcaneMage{
	ArcaneMage: &proto.ArcaneMage{
		Options: &proto.ArcaneMage_Options{
			ClassOptions:            &proto.MageOptions{},
			FocusMagicPercentUptime: 90,
		},
	},
}

var FullArcaneConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58086, // Flask of the Draconic Mind
	FoodId:   62290, // Seafood Magnifique Feast
	PotId:    58091, // Volcanic Potion
	PrepotId: 58091, // Volcanic Potion

}
