package enhancement

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterEnhancementShaman()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceDwarf,
		OtherRaces: []proto.Race{proto.Race_RaceOrc, proto.Race_RaceTroll},

		// The above line is the actual line for the ring but it is causing an error in the test
		GearSet:     core.GetGearSet("../../../ui/shaman/enhancement/gear_sets", "p4.orc"),
		Talents:     StandardTalents,
		Glyphs:      StandardGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Standard", SpecOptions: PlayerOptionsStandard},
		Rotation:    core.GetAplRotation("../../../ui/shaman/enhancement/apls", "default"),

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

var StandardTalents = "3020023-2333310013003012321"
var StandardGlyphs = &proto.Glyphs{
	Prime1: int32(proto.ShamanPrimeGlyph_GlyphOfLavaLash),
	Prime2: int32(proto.ShamanPrimeGlyph_GlyphOfWindfuryWeapon),
	Prime3: int32(proto.ShamanPrimeGlyph_GlyphOfFeralSpirit),
}

var FullConsumes = &proto.Consumes{
	Flask:         proto.Flask_FlaskOfTheWinds,
	Food:          proto.Food_FoodGrilledDragon,
	DefaultPotion: proto.Potions_PotionOfTheTolvir,
	PrepopPotion:  proto.Potions_PotionOfTheTolvir,
}

var TotemsBasic = &proto.ShamanTotems{
	Elements: &proto.TotemSet{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
	Ancestors: &proto.TotemSet{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
	Spirits: &proto.TotemSet{
		Earth: proto.EarthTotem_StrengthOfEarthTotem,
		Air:   proto.AirTotem_WindfuryTotem,
		Water: proto.WaterTotem_ManaSpringTotem,
		Fire:  proto.FireTotem_SearingTotem,
	},
	Earth: proto.EarthTotem_StrengthOfEarthTotem,
	Air:   proto.AirTotem_WindfuryTotem,
	Water: proto.WaterTotem_ManaSpringTotem,
	Fire:  proto.FireTotem_SearingTotem,
}

var PlayerOptionsStandard = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{
				Shield:  proto.ShamanShield_LightningShield,
				Totems:  TotemsBasic,
				ImbueMh: proto.ShamanImbue_WindfuryWeapon,
			},
			ImbueOh: proto.ShamanImbue_FlametongueWeapon,
		},
	},
}
