package fire

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterFireMage()
}

func TestFire(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassMage,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceWorgen},
		GearSet:    core.GetGearSet("../../../ui/mage/fire/gear_sets", "p4_fire"),
		Talents:    FireTalents,
		Glyphs:     FireGlyphs,
		Consumes:   FullFireConsumes,

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
	Major1: int32(proto.MageMajorGlyph_GlyphOfEvocation),
	Major2: int32(proto.MageMajorGlyph_GlyphOfDragonSBreath),
	Major3: int32(proto.MageMajorGlyph_GlyphOfInvisibility),
}

var PlayerOptionsFire = &proto.Player_FireMage{
	FireMage: &proto.FireMage{
		Options: &proto.FireMage_Options{
			ClassOptions: &proto.MageOptions{},
		},
	},
}

var FullFireConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheFrostWyrm,
	Food:          proto.Food_FoodFirecrackerSalmon,
	DefaultPotion: proto.Potions_PotionOfSpeed,
}
