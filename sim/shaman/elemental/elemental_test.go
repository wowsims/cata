package elemental

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterElementalShaman()
}

func TestElemental(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/shaman/elemental/gear_sets", "preraid"),
		Talents: TalentsASEB,
		Glyphs:  StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "TalentsEchoUnleashed",
				Talents: TalentsEEUF,
				Glyphs:  AoEGlyphs,
			},
			{
				Label:   "TalentsEMPrimal",
				Talents: TalentsEMPE,
				Glyphs:  StandardGlyphs,
			},
		},
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsFireElemental},
		Rotation:    core.GetAplRotation("../../../ui/shaman/elemental/apls", "default"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/shaman/elemental/apls", "aoe"),
			core.GetAplRotation("../../../ui/shaman/elemental/apls", "unleash"),
		},

		ItemFilter: core.ItemFilter{
			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeAxe,
				proto.WeaponType_WeaponTypeDagger,
				proto.WeaponType_WeaponTypeFist,
				proto.WeaponType_WeaponTypeMace,
				proto.WeaponType_WeaponTypeOffHand,
				proto.WeaponType_WeaponTypeShield,
				proto.WeaponType_WeaponTypeStaff,
			},
			ArmorType:         proto.ArmorType_ArmorTypeMail,
			RangedWeaponTypes: []proto.RangedWeaponType{},
		},
	}))
}

var TalentsEMUF = "313131"
var TalentsEMPE = "313132"
var TalentsEMEB = "313133"

var TalentsASUF = "313231"
var TalentsASPE = "313232"
var TalentsASEB = "313233"

var TalentsEEUF = "313331"
var TalentsEEPE = "313332"
var TalentsEEEB = "313333"

var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfHealingStreamTotem),
}
var AoEGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfChainLightning),
}

var NoTotems = &proto.ShamanTotems{}

var PlayerOptionsFireElemental = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			ClassOptions: &proto.ShamanOptions{
				Shield: proto.ShamanShield_LightningShield,
				FeleAutocast: &proto.FeleAutocastSettings{
					AutocastFireblast: true,
					AutocastFirenova:  true,
					AutocastImmolate:  true,
					AutocastEmpower:   false,
				},
			},
		},
	},
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:    58086, // Flask of the Draconic Mind
	FoodId:     62671, // Severed Sagefish Head
	PotId:      58091, // Volcanic Potion
	PrepotId:   58091, // Volcanic Potion
	ConjuredId: 20520, // Dark Rune
}
