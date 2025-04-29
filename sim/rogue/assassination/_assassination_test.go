package assassination

import (
	"testing"

	_ "github.com/wowsims/mop/sim/common" // imported to get item effects included.
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func init() {
	RegisterAssassinationRogue()
}

func TestAssassination(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
		Class:      proto.Class_ClassRogue,
		Race:       proto.Race_RaceHuman,
		OtherRaces: []proto.Race{proto.Race_RaceOrc},
		GearSet:    core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p1_assassination"),

		OtherGearSets: []core.GearSetCombo{
			core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p3_assassination"),
			core.GetGearSet("../../../ui/rogue/assassination/gear_sets", "p4_assassination"),
		},

		Talents:     AssassinationTalents,
		Glyphs:      AssassinationGlyphs,
		Consumes:    FullConsumes,
		SpecOptions: core.SpecOptionsCombo{Label: "Assassination", SpecOptions: PlayerOptionsAssassinationDI},

		OtherSpecOptions: []core.SpecOptionsCombo{
			{Label: "MH Instant OH Deadly", SpecOptions: PlayerOptionsAssassinationID},
			{Label: "MH Instant OH Instant", SpecOptions: PlayerOptionsAssassinationII},
			{Label: "MH Deadly OH Deadly", SpecOptions: PlayerOptionsAssassinationDD},
		},

		Rotation:       core.GetAplRotation("../../../ui/rogue/assassination/apls", "mutilate"),
		OtherRotations: []core.RotationCombo{},

		ItemFilter: core.ItemFilter{
			ArmorType: proto.ArmorType_ArmorTypeLeather,

			RangedWeaponTypes: []proto.RangedWeaponType{
				proto.RangedWeaponType_RangedWeaponTypeBow,
				proto.RangedWeaponType_RangedWeaponTypeCrossbow,
				proto.RangedWeaponType_RangedWeaponTypeGun,
				proto.RangedWeaponType_RangedWeaponTypeThrown,
			},

			WeaponTypes: []proto.WeaponType{
				proto.WeaponType_WeaponTypeDagger,
			},
		},

		// General practice is to not include stat weights in test suite configs to speed up test execution, but at least one spec should
		// include them so the core functionality is tested. Assassination Rogue was chosen for this since it has appreciable EP
		// contributions from both physical and spell Crit, and therefore provides a good test case of school-specific EP consolidation for
		// Rating stats.
		StatsToWeigh:       []proto.Stat{proto.Stat_StatCritRating},
		PseudoStatsToWeigh: []proto.PseudoStat{proto.PseudoStat_PseudoStatPhysicalCritPercent, proto.PseudoStat_PseudoStatSpellCritPercent},
		EPReferenceStat:    proto.Stat_StatAgility,
	}))
}

var AssassinationTalents = "0333230013122110321-002-203003"

var AssassinationGlyphs = &proto.Glyphs{}

var PlayerOptionsAssassinationDI = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationID = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationDD = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_DeadlyPoison,
				OhImbue: proto.RogueOptions_DeadlyPoison,
				ThImbue: proto.RogueOptions_DeadlyPoison,
			},
		},
	},
}

var PlayerOptionsAssassinationII = &proto.Player_AssassinationRogue{
	AssassinationRogue: &proto.AssassinationRogue{
		Options: &proto.AssassinationRogue_Options{
			ClassOptions: &proto.RogueOptions{
				MhImbue: proto.RogueOptions_InstantPoison,
				OhImbue: proto.RogueOptions_InstantPoison,
				ThImbue: proto.RogueOptions_InstantPoison,
			},
		},
	},
}

var FullConsumes = &proto.Consumes{
	Flask:           proto.Flask_FlaskOfTheWinds,
	DefaultPotion:   proto.Potions_PotionOfTheTolvir,
	DefaultConjured: proto.Conjured_ConjuredRogueThistleTea,
}
