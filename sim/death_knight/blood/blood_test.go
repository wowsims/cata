package blood

import (
	_ "github.com/wowsims/cata/sim/common" // imported to get item effects included.
)

func init() {
	RegisterBloodDeathKnight()
}

// func TestBlood(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:      proto.Class_ClassDeathKnight,
// 		Race:       proto.Race_RaceOrc,
// 		OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 		GearSet:     core.GetGearSet("../../../ui/DeathKnight/gear_sets", "p3_blood"),
// 		Talents:     BloodTalents,
// 		Glyphs:      BloodDefaultGlyphs,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsBlood},
// 		Rotation:    core.GetAplRotation("../../../ui/DeathKnight/apls", "blood_dps"),

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestUnholy(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:      proto.Class_ClassDeathKnight,
// 		Race:       proto.Race_RaceOrc,
// 		OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 		GearSet:     core.GetGearSet("../../../ui/DeathKnight/gear_sets", "p3_uh_dw"),
// 		Talents:     UnholyTalents,
// 		Glyphs:      UnholyDefaultGlyphs,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsUnholy},
// 		Rotation:    core.GetAplRotation("../../../ui/DeathKnight/apls", "uh_2h_ss"),

// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../../ui/DeathKnight/apls", "uh_dnd_aoe"),
// 			core.GetAplRotation("../../../ui/DeathKnight/apls", "unholy_dw_ss"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestFrost(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:      proto.Class_ClassDeathKnight,
// 		Race:       proto.Race_RaceOrc,
// 		OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 		GearSet:     core.GetGearSet("../../../ui/DeathKnight/gear_sets", "p3_frost"),
// 		Talents:     FrostTalents,
// 		Glyphs:      FrostDefaultGlyphs,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
// 		OtherSpecOptions: []core.SpecOptionsCombo{
// 			{Label: "Desync", SpecOptions: PlayerOptionsDesyncFrost},
// 		},
// 		Rotation: core.GetAplRotation("../../../ui/DeathKnight/apls", "frost_bl_pesti"),

// 		OtherRotations: []core.RotationCombo{
// 			core.GetAplRotation("../../../ui/DeathKnight/apls", "frost_uh_pesti"),
// 		},

// 		ItemFilter: ItemFilter,
// 	}))
// }

// func TestFrostUH(t *testing.T) {
// 	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator(core.CharacterSuiteConfig{
// 		Class:      proto.Class_ClassDeathKnight,
// 		Race:       proto.Race_RaceOrc,
// 		OtherRaces: []proto.Race{proto.Race_RaceHuman},

// 		GearSet:     core.GetGearSet("../../../ui/DeathKnight/gear_sets", "p3_frost"),
// 		Talents:     FrostUHTalents,
// 		Glyphs:      FrostUHDefaultGlyphs,
// 		Consumes:    FullConsumes,
// 		SpecOptions: core.SpecOptionsCombo{Label: "Basic", SpecOptions: PlayerOptionsFrost},
// 		Rotation:    core.GetAplRotation("../../../ui/DeathKnight/apls", "frost_uh_pesti"),

// 		ItemFilter: ItemFilter,
// 	}))
// }

// var BloodTalents = "2305120530003303231023001351--230220305003"
// var BloodDefaultGlyphs = &proto.Glyphs{
// 	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon),
// 	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfDeathStrike),
// 	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfDisease),
// 	// No interesting minor glyphs.
// }

// var FrostTalents = "23050005-32005350352203012300033101351"
// var FrostDefaultGlyphs = &proto.Glyphs{
// 	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfFrostStrike),
// 	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfObliterate),
// 	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfDisease),
// 	// No interesting minor glyphs.
// }

// var FrostUHTalents = "01-32002350342203012300033101351-230200305003"
// var FrostUHDefaultGlyphs = &proto.Glyphs{
// 	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfFrostStrike),
// 	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfObliterate),
// 	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfDisease),
// 	// No interesting minor glyphs.
// }

// var UnholyTalents = "-320043500002-2300303050032152000150013133051"
// var UnholyDefaultGlyphs = &proto.Glyphs{
// 	Major1: int32(proto.DeathKnightMajorGlyph_GlyphOfTheGhoul),
// 	Major2: int32(proto.DeathKnightMajorGlyph_GlyphOfDarkDeath),
// 	Major3: int32(proto.DeathKnightMajorGlyph_GlyphOfDeathAndDecay),
// 	// No interesting minor glyphs.
// }

// var PlayerOptionsBlood = &proto.Player_DeathKnight{
// 	DeathKnight: &proto.DeathKnight{
// 		Options: DeathKnightOptions,
// 	},
// }

// var PlayerOptionsUnholy = &proto.Player_DeathKnight{
// 	DeathKnight: &proto.DeathKnight{
// 		Options: DeathKnightOptions,
// 	},
// }

// var PlayerOptionsFrost = &proto.Player_DeathKnight{
// 	DeathKnight: &proto.DeathKnight{
// 		Options: DeathKnightOptions,
// 	},
// }

// var PlayerOptionsDesyncFrost = &proto.Player_DeathKnight{
// 	DeathKnight: &proto.DeathKnight{
// 		Options: DeathKnightOptions,
// 	},
// }

// var DeathKnightOptions = &proto.DeathKnight_Options{
// 	UnholyFrenzyTarget: &proto.UnitReference{Type: proto.UnitReference_Player, Index: 0},
// 	DrwPestiApply:      true,
// 	StartingRunicPower: 0,
// 	PetUptime:          1,
// }

// var FullConsumes = &proto.Consumes{
// 	Flask:         proto.Flask_FlaskOfEndlessRage,
// 	DefaultPotion: proto.Potions_PotionOfSpeed,
// 	PrepopPotion:  proto.Potions_PotionOfSpeed,
// 	Food:          proto.Food_FoodDragonfinFilet,
// }

// var ItemFilter = core.ItemFilter{
// 	ArmorType: proto.ArmorType_ArmorTypePlate,

// 	WeaponTypes: []proto.WeaponType{
// 		proto.WeaponType_WeaponTypeAxe,
// 		proto.WeaponType_WeaponTypeSword,
// 		proto.WeaponType_WeaponTypeMace,
// 	},
// }
