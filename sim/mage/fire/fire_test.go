package fire

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterFireMage()
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:       proto.Class_ClassMage,
		Race:        proto.Race_RaceTroll,
		OtherRaces:  []proto.Race{proto.Race_RaceWorgen},
		GearSet:     core.GetGearSet("../../../ui/mage/fire/gear_sets", "p4_fire"),
		Talents:     FireTalents,
		Glyphs:      FireGlyphs,
		Consumables: FullArcaneConsumesSpec,

		SpecOptions: core.SpecOptionsCombo{Label: "Fire", SpecOptions: PlayerOptionsFire},
		Rotation:    core.GetAplRotation("../../../ui/mage/fire/apls", "fire"),

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

var FireTalents = "203-230330221120121213031-03"
var FireGlyphs = &proto.Glyphs{
	Prime1: int32(proto.MagePrimeGlyph_GlyphOfFireball),
	Prime2: int32(proto.MagePrimeGlyph_GlyphOfPyroblast),
	Prime3: int32(proto.MagePrimeGlyph_GlyphOfMoltenArmor),
	Major1: int32(proto.MageMajorGlyph_GlyphOfEvocation),
	Major2: int32(proto.MageMajorGlyph_GlyphOfDragonsBreath),
	Major3: int32(proto.MageMajorGlyph_GlyphOfInvisibility),
}

var PlayerOptionsFire = &proto.Player_FireMage{
	FireMage: &proto.FireMage{
		Options: &proto.FireMage_Options{
			ClassOptions: &proto.MageOptions{},
		},
	},
}
var FullArcaneConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58086, // Flask of the Draconic Mind
	FoodId:   62290, // Seafood Magnifique Feast
	PotId:    58091, // Volcanic Potion
	PrepotId: 58091, // Volcanic Potion
	TinkerId: 4179,  // Synapse Springs
}
