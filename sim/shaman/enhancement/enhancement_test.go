package enhancement

import (
	"testing"

	"github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterEnhancementShaman()
	common.RegisterAllEffects()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassShaman,
		Race:       proto.Race_RaceDwarf,
		OtherRaces: []proto.Race{proto.Race_RaceOrc, proto.Race_RaceTroll, proto.Race_RaceDraenei, proto.Race_RaceAlliancePandaren},

		// The above line is the actual line for the ring but it is causing an error in the test
		GearSet: core.GetGearSet("../../../ui/shaman/enhancement/gear_sets", "preraid"),
		Talents: TalentsASEB,
		Glyphs:  StandardGlyphs,
		OtherTalentSets: []core.TalentsCombo{
			{
				Label:   "TalentsEchoUnleashed",
				Talents: TalentsEEUF,
				Glyphs:  StandardGlyphs,
			},
			{
				Label:   "TalentsEMPrimal",
				Talents: TalentsEMPE,
				Glyphs:  StandardGlyphs,
			},
		},
		Consumables: FullConsumesSpec,
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
	Major3: int32(proto.ShamanMajorGlyph_GlyphOfFireNova),
}

var FullConsumesSpec = &proto.ConsumesSpec{
	FlaskId:  58087, // Flask of the Winds
	FoodId:   62662, // Grilled Dragon
	PotId:    58145, // Potion of the Tol'vir
	PrepotId: 58145, // Potion of the Tol'vir
}

var PlayerOptionsStandard = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: &proto.EnhancementShaman_Options{
			ClassOptions: &proto.ShamanOptions{
				Shield:  proto.ShamanShield_LightningShield,
				ImbueMh: proto.ShamanImbue_WindfuryWeapon,
				FeleAutocast: &proto.FeleAutocastSettings{
					AutocastFireblast: true,
					AutocastFirenova:  true,
					AutocastImmolate:  true,
					AutocastEmpower:   false,
				},
			},
			ImbueOh: proto.ShamanImbue_FlametongueWeapon,
		},
	},
}
