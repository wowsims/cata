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
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/mage/arcane/gear_sets", "p1_bis"),
		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/mage/arcane/gear_sets", "prebis"),
		},
		Talents:     ArcaneTalents,
		Glyphs:      ArcaneGlyphs,
		Consumables: FullArcaneConsumesSpec,

		SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptionsArcane},
		Rotation:    core.GetAplRotation("../../../ui/mage/arcane/apls", "default"),

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

var ArcaneTalents = "311122"
var ArcaneGlyphs = &proto.Glyphs{
	Major1: int32(proto.MageMajorGlyph_GlyphOfArcanePower),
	Major2: int32(proto.MageMajorGlyph_GlyphOfRapidDisplacement),
	Major3: int32(proto.MageMajorGlyph_GlyphOfEvocation),
	Minor1: int32(proto.MageMinorGlyph_GlyphOfMomentum),
	Minor2: int32(proto.MageMinorGlyph_GlyphOfRapidTeleportation),
	Minor3: int32(proto.MageMinorGlyph_GlyphOfMirrorImage),
}

var PlayerOptionsArcane = &proto.Player_ArcaneMage{
	ArcaneMage: &proto.ArcaneMage{
		Options: &proto.ArcaneMage_Options{
			ClassOptions: &proto.MageOptions{},
		},
	},
}

var FullArcaneConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  76085, // Flask of the Warm Sun
	FoodId:   74650, // Mogu Fish Stew
	PotId:    76093, // Potion of the Jade Serpent
	PrepotId: 76093, // Potion of the Jade Serpent
}
