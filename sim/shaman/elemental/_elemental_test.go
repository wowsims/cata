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

		GearSet: core.GetGearSet("../../../ui/shaman/elemental/gear_sets", "p4.default"),
		Talents: TalentsTotemDuration,
		Glyphs:  StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "TalentsAoE",
				Talents: TalentsImprovedShields,
				Glyphs:  AoEGlyphs,
			},
			{
				Label:   "TalentsImprovedShields",
				Talents: TalentsImprovedShields,
				Glyphs:  AlternateGlyphs,
			},
		},
		Consumables: FullConsumesSpec,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsFireElemental},
		Rotation:    core.GetAplRotation("../../../ui/shaman/elemental/apls", "default"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/shaman/elemental/apls", "aoe"),
			core.GetAplRotation("../../../ui/shaman/elemental/apls", "unleash"),
		},
		ItemSwapSet: core.GetItemSwapGearSet("../../../ui/shaman/elemental/gear_sets", "p4_item_swap"),

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

var TalentsTotemDuration = "303202321223110132-201-20302"
var TalentsImprovedShields = "3032023212231101321-2030022"
var StandardGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfHealingStreamTotem),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfStoneclawTotem),
}
var AoEGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfChainLightning),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfStoneclawTotem),
}
var AlternateGlyphs = &proto.Glyphs{
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfHealingStreamTotem),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfStoneclawTotem),
}

var NoTotems = &proto.ShamanTotems{}
var TotemsBasic = &proto.ShamanTotems{
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_SearingTotem,
}

var TotemsFireElemental = &proto.ShamanTotems{
	Elements: &proto.TotemSet{
		Earth: proto.EarthTotem_TremorTotem,
		Air:   proto.AirTotem_WrathOfAirTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
	Ancestors: &proto.TotemSet{
		Earth: proto.EarthTotem_EarthElementalTotem,
		Fire:  proto.FireTotem_FireElementalTotem,
	},
	Spirits: &proto.TotemSet{
		Earth: proto.EarthTotem_TremorTotem,
		Air:   proto.AirTotem_WrathOfAirTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
	Earth: proto.EarthTotem_TremorTotem,
	Air:   proto.AirTotem_WrathOfAirTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_SearingTotem,
}

var PlayerOptionsFireElemental = &proto.Player_ElementalShaman{
	ElementalShaman: &proto.ElementalShaman{
		Options: &proto.ElementalShaman_Options{
			ClassOptions: &proto.ShamanOptions{
				Shield: proto.ShamanShield_LightningShield,
				Totems: TotemsFireElemental,
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
	TinkerId:   82174, // Synapse Springs
}
