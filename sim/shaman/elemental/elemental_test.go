package elemental

import (
	"testing"

	_ "github.com/wowsims/cata/sim/common"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func init() {
	RegisterElementalShaman()
}

func TestElemental(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceTroll,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},

		GearSet: core.GetGearSet("../../../ui/shaman/elemental/gear_sets", "p1"),
		Talents: TalentsTotemDuration,
		Glyphs:  StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "TalentsImprovedShields",
				Talents: TalentsImprovedShields,
				Glyphs:  AlternateGlyphs,
			},
		},
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsFireElemental},
		Rotation:    core.GetAplRotation("../../../ui/shaman/elemental/apls", "default"),
		OtherRotations: []core.RotationCombo{
			core.GetAplRotation("../../../ui/shaman/elemental/apls", "aoe"),
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
			ArmorType: proto.ArmorType_ArmorTypeMail,
			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeRelic,
			},
		},
	}))
}

var TalentsTotemDuration = "303202321223110132-201-20302"
var TalentsImprovedShields = "3032023212231101321-2030022"
var StandardGlyphs = &proto.Glyphs{
	Prime1: int32(proto.ShamanPrimeGlyph_GlyphOfFlameShock),
	Prime2: int32(proto.ShamanPrimeGlyph_GlyphOfLavaBurst),
	Prime3: int32(proto.ShamanPrimeGlyph_GlyphOfLightningBolt),
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfThunder),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfFireNova),
}
var AlternateGlyphs = &proto.Glyphs{
	Prime1: int32(proto.FireTotem_FireElementalTotem),
	Prime2: int32(proto.ShamanImbue_FlametongueWeapon),
	Prime3: int32(proto.ShamanPrimeGlyph_GlyphOfUnleashedLightning),
	Major1: int32(proto.ShamanMajorGlyph_GlyphOfLightningShield),
	Major2: int32(proto.ShamanMajorGlyph_GlyphOfThunder),
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfChainLightning),
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

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfTheDraconicMind,
	Food:            proto.Food_FoodSeveredSagefish,
	DefaultPotion:   proto.Potions_VolcanicPotion,
	PrepopPotion:    proto.Potions_VolcanicPotion,
	DefaultConjured: proto.Conjured_ConjuredDarkRune,
}
