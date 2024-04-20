package database

import (
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

// Note: EffectId AND SpellId are required for all enchants, because they are
// used by various importers/exporters. ItemId is optional.

var EnchantOverrides = []*proto.UIEnchant{
	// HANDS
	{EffectId: 4061, SpellId: 74132, Name: "Enchant Gloves - Mastery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mastery: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 4068, SpellId: 74198, Name: "Enchant Gloves - Haste", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 50, stats.SpellHaste: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 4075, SpellId: 74212, Name: "Enchant Gloves - Exceptional Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 4082, SpellId: 74220, Name: "Enchant Gloves - Greater Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 4106, SpellId: 74254, Name: "Enchant Gloves - Mighty Strength", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Strength: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 4107, SpellId: 74255, Name: "Enchant Gloves - Greater Mastery", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mastery: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	// Moved to engineering consumes
	//{EffectId: 4179, SpellId: 82175, Name: "Synapse Springs", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4180, SpellId: 82177, Name: "Quickflip Deflection Plates", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4181, SpellId: 82180, Name: "Tazik Shocker", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4182, SpellId: 82200, Name: "Spinal Healing Injector", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4183, SpellId: 82201, Name: "Z50 Mana Gulper", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},

	// FEET
	{EffectId: 4062, SpellId: 74189, Name: "Enchant Boots - Earthen Vitality", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4069, SpellId: 74199, Name: "Enchant Boots - Haste", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 50, stats.SpellHaste: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4076, SpellId: 74213, Name: "Enchant Boots - Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4092, SpellId: 74236, Name: "Enchant Boots - Precision", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 50, stats.SpellHit: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4094, SpellId: 74238, Name: "Enchant Boots - Mastery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mastery: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4105, SpellId: 74252, Name: "Enchant Boots - Assassin's Step", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4104, SpellId: 74253, Name: "Enchant Boots - Lavawalker", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Mastery: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},

	// CHEST
	{EffectId: 4063, SpellId: 74191, Name: "Enchant Chest - Mighty Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 15, stats.Strength: 15, stats.Agility: 15, stats.Intellect: 15, stats.Spirit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 4070, SpellId: 74200, Name: "Enchant Chest - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 4077, SpellId: 74214, Name: "Enchant Chest - Mighty Resilience", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Resilience: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 4088, SpellId: 74231, Name: "Enchant Chest - Exceptional Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 4102, SpellId: 74250, Name: "Enchant Chest - Peerless Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 20, stats.Strength: 20, stats.Agility: 20, stats.Intellect: 20, stats.Spirit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 4103, SpellId: 74251, Name: "Enchant Chest - Greater Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 75}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// CLOAK
	{EffectId: 4064, SpellId: 74192, Name: "Enchant Cloak - Greater Spell Piercing", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPenetration: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4072, SpellId: 74202, Name: "Enchant Cloak - Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4087, SpellId: 74230, Name: "Enchant Cloak - Critical Strike", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 50, stats.SpellCrit: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4090, SpellId: 74234, Name: "Enchant Cloak - Protection", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 250}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4096, SpellId: 74240, Name: "Enchant Cloak - Greater Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4100, SpellId: 74247, Name: "Enchant Cloak - Greater Critical Strike", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 65, stats.SpellCrit: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 4115, SpellId: 75172, Name: "Lightweave Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 4116, SpellId: 75175, Name: "Darkglow Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 4118, SpellId: 75178, Name: "Swordguard Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 1}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},

	// WRISTS
	{EffectId: 4065, SpellId: 74193, Name: "Enchant Bracer - Speed", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 50, stats.SpellHaste: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4071, SpellId: 74201, Name: "Enchant Bracer - Critical Strike", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 50, stats.SpellCrit: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4086, SpellId: 74229, Name: "Enchant Bracer - Superior Dodge", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Dodge: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4089, SpellId: 74232, Name: "Enchant Bracer - Precision", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 50, stats.SpellHit: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4093, SpellId: 74237, Name: "Enchant Bracer - Exceptional Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4095, SpellId: 74239, Name: "Enchant Bracer - Greater Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4101, SpellId: 74248, Name: "Enchant Bracer - Greater Critical Strike", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 65, stats.SpellCrit: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4108, SpellId: 74256, Name: "Enchant Bracer - Greater Speed", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 65, stats.SpellHaste: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4189, SpellId: 85007, Name: "Draconic Embossment - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 195}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 4190, SpellId: 85008, Name: "Draconic Embossment - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 4191, SpellId: 85009, Name: "Draconic Embossment - Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 4192, SpellId: 85010, Name: "Draconic Embossment - Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 4256, SpellId: 96261, Name: "Enchant Bracer - Major Strength", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Strength: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4257, SpellId: 96262, Name: "Enchant Bracer - Mighty Intellect", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 4258, SpellId: 96264, Name: "Enchant Bracer - Agility", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},

	// WEAPON
	{EffectId: 4066, SpellId: 74195, Name: "Enchant Weapon - Mending", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4067, SpellId: 74197, Name: "Enchant Weapon - Avalanche", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4073, SpellId: 74207, Name: "Enchant Shield - Protection", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Armor: 160}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 4074, SpellId: 74211, Name: "Enchant Weapon - Elemental Slayer", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4083, SpellId: 74223, Name: "Enchant Weapon - Hurricane", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4084, SpellId: 74225, Name: "Enchant Weapon - Heartsong", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4085, SpellId: 74226, Name: "Enchant Shield - Mastery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mastery: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 4091, SpellId: 74235, Name: "Enchant Off-Hand - Superior Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeOffHand},
	{EffectId: 4097, SpellId: 74242, Name: "Enchant Weapon - Power Torrent", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4098, SpellId: 74244, Name: "Enchant Weapon - Windwalk", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4099, SpellId: 74246, Name: "Enchant Weapon - Landslide", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4215, ItemId: 55055, SpellId: 92433, Name: "Elementium Shield Spike", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 4216, ItemId: 55056, SpellId: 92437, Name: "Pyrium Shield Spike", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 4217, ItemId: 55057, SpellId: 93448, Name: "Pyrium Weapon Chain", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHit: 40, stats.SpellHit: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 4227, SpellId: 95471, Name: "Enchant 2H Weapon - Mighty Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},

	// FINGER
	{EffectId: 4078, SpellId: 74215, Name: "Enchant Ring - Strength", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 4079, SpellId: 74216, Name: "Enchant Ring - Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 4080, SpellId: 74217, Name: "Enchant Ring - Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 4081, SpellId: 74218, Name: "Enchant Ring - Greater Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 60}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// LEGS
	{EffectId: 4109, ItemId: 54449, SpellId: 75149, Name: "Ghostly Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 55, stats.Spirit: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4110, ItemId: 54450, SpellId: 75150, Name: "Powerful Ghostly Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Intellect: 95, stats.Spirit: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4111, ItemId: 54447, SpellId: 75151, Name: "Enchanted Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 55, stats.Stamina: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4112, ItemId: 54448, SpellId: 75152, Name: "Powerful Enchanted Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Intellect: 95, stats.Stamina: 80}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4113, SpellId: 75154, Name: "Master's Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 95, stats.Stamina: 80}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 4114, SpellId: 75155, Name: "Sanctified Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 95, stats.Spirit: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 4122, ItemId: 56502, SpellId: 78169, Name: "Scorched Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeCrit: 45, stats.SpellCrit: 45, stats.AttackPower: 110, stats.RangedAttackPower: 110}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4124, ItemId: 56503, SpellId: 78170, Name: "Twilight Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 85, stats.Agility: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4126, ItemId: 56550, SpellId: 78171, Name: "Dragonscale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.MeleeCrit: 55, stats.SpellCrit: 55, stats.AttackPower: 190, stats.RangedAttackPower: 190}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4127, ItemId: 56551, SpellId: 78172, Name: "Charscale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 145, stats.Agility: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 4440, SpellId: 85067, Name: "Dragonbone Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 55, stats.SpellCrit: 55, stats.AttackPower: 190, stats.RangedAttackPower: 190}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 4439, SpellId: 85068, Name: "Charscale Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 145, stats.Agility: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},

	// Firelands
	// {EffectId: 4438, SpellId: 101600, Name: "Drakehide Leg Reinforcements", Phase: 3, Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 145, stats.Dodge: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	// {EffectId: 4270, ItemId: 71720, SpellId: 101598, Name: "Drakehide Leg Armor", phase: 3, Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 145, stats.Dodge: 55}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},

	// HEAD
	{EffectId: 4120, ItemId: 56477, SpellId: 78165, Name: "Savage Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 36}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeFeet, proto.ItemType_ItemTypeHands}},
	{EffectId: 4121, ItemId: 56517, SpellId: 78166, Name: "Heavy Savage Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 44}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeFeet, proto.ItemType_ItemTypeHands}},
	{EffectId: 4206, ItemId: 68764, SpellId: 86931, Name: "Arcanum of the Earthern Ring", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 90, stats.Dodge: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4207, ItemId: 68765, SpellId: 86932, Name: "Arcanum of Hyjal", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 60, stats.MeleeCrit: 35, stats.SpellCrit: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4208, ItemId: 68767, SpellId: 86933, Name: "Arcanum of the Highlands", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Strength: 60, stats.Mastery: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4209, ItemId: 68766, SpellId: 86934, Name: "Arcanum of Ramkahen", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 60, stats.MeleeHaste: 35, stats.SpellHaste: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4245, ItemId: 68770, SpellId: 96245, Name: "Arcanum of Vicious Intellect", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 60, stats.Resilience: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4246, ItemId: 68769, SpellId: 96246, Name: "Arcanum of Vicious Agility", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 60, stats.Resilience: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4247, ItemId: 68768, SpellId: 96247, Name: "Arcanum of Vicious Strength", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Strength: 60, stats.Resilience: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},

	// RANGED
	{EffectId: 4175, ItemId: 59594, SpellId: 81932, Name: "Gnomish X-Ray Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 4176, ItemId: 59595, SpellId: 81933, Name: "R19 Threatfinder", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHit: 88, stats.SpellHit: 88}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 4177, ItemId: 59596, SpellId: 81934, Name: "Safety Catch Removal Kit", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHaste: 88, stats.SpellHaste: 88}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 4267, ItemId: 70139, SpellId: 99623, Name: "Flintlocke's Woodchucker", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged, ClassAllowlist: []proto.Class{proto.Class_ClassHunter}},

	// WAIST
	// Moved to engineering consumes
	//{EffectId: 4187, SpellId: 84424, Name: "Invisibility Field", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4214, SpellId: 84425, Name: "Cardboard Assassin", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	//{EffectId: 4188, SpellId: 84427, Name: "Grounded Plasma Shield", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},

	// SHOULDERS
	{EffectId: 4193, SpellId: 86375, Name: "Swiftsteel Inscription", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 130, stats.Mastery: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 4194, SpellId: 86401, Name: "Lionsmane Inscription", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Strength: 130, stats.MeleeCrit: 25, stats.SpellCrit: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 4195, SpellId: 86402, Name: "Inscription of the Earth Prince", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 195, stats.Dodge: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 4196, SpellId: 86403, Name: "Felfire Inscription", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 130, stats.MeleeHaste: 25, stats.SpellHaste: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 4197, ItemId: 62321, SpellId: 86847, Name: "Inscription of Unbreakable Quartz", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 45, stats.Dodge: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4198, ItemId: 68717, SpellId: 86854, Name: "Greater Inscription of Unbreakable Quartz", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 75, stats.Dodge: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4199, ItemId: 62342, SpellId: 86898, Name: "Inscription of Charged Lodestone", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Intellect: 30, stats.MeleeHaste: 20, stats.SpellHaste: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4200, ItemId: 68715, SpellId: 86899, Name: "Greater Inscription of Charged Lodestone", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Intellect: 50, stats.MeleeHaste: 25, stats.SpellHaste: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4201, ItemId: 62344, SpellId: 86900, Name: "Inscription of Jagged Stone", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Strength: 30, stats.MeleeCrit: 20, stats.SpellCrit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4202, ItemId: 68716, SpellId: 86901, Name: "Greater Inscription of Jagged Stone", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Strength: 50, stats.MeleeCrit: 25, stats.SpellCrit: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4203, SpellId: 86906, Name: "Inscription of Shattered Crystal", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 30, stats.Mastery: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4204, ItemId: 68714, SpellId: 86907, Name: "Greater Inscription of Shattered Crystal", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Agility: 50, stats.Mastery: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4205, ItemId: 62347, SpellId: 86909, Name: "Inscription of Shattered Crystal", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Agility: 30, stats.Mastery: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4248, ItemId: 68772, SpellId: 96249, Name: "Greater Inscription of Vicious Intellect", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Intellect: 50, stats.Resilience: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4249, ItemId: 68773, SpellId: 96250, Name: "Greater Inscription of Vicious Strength", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Strength: 50, stats.Resilience: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 4250, ItemId: 68815, SpellId: 96251, Name: "Greater Inscription of Vicious Agility", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Agility: 50, stats.Resilience: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},

	////// WOTLK
	/////////////////////////////
	/////////////////////////////
	// Multi-slot
	{EffectId: 2988, ItemId: 29487, SpellId: 35419, Name: "Nature Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.NatureResistance: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 3329, ItemId: 38375, SpellId: 50906, Name: "Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},
	{EffectId: 3330, ItemId: 38376, SpellId: 50909, Name: "Heavy Borean Armor Kit", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, ExtraTypes: []proto.ItemType{proto.ItemType_ItemTypeChest, proto.ItemType_ItemTypeShoulder, proto.ItemType_ItemTypeLegs, proto.ItemType_ItemTypeHands, proto.ItemType_ItemTypeFeet}, EnchantType: proto.EnchantType_EnchantTypeKit},

	// Head
	{EffectId: 3795, ItemId: 44069, SpellId: 59777, Name: "Arcanum of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3796, ItemId: 44075, SpellId: 59784, Name: "Arcanum of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 29, stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3842, ItemId: 44875, SpellId: 61271, Name: "Arcanum of the Savage Gladiator", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.Resilience: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3812, ItemId: 44137, SpellId: 59944, Name: "Arcanum of the Frosty Soul", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.FrostResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3813, ItemId: 44138, SpellId: 59945, Name: "Arcanum of Toxic Warding", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.NatureResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3814, ItemId: 44139, SpellId: 59946, Name: "Arcanum of the Fleeing Shadow", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.ShadowResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3815, ItemId: 44140, SpellId: 59947, Name: "Arcanum of the Eclipsed Moon", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.ArcaneResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3816, ItemId: 44141, SpellId: 59948, Name: "Arcanum of the Flame's Soul", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 30, stats.FireResistance: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3819, ItemId: 44876, SpellId: 59960, Name: "Arcanum of Blissful Mending", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30, stats.MP5: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3820, ItemId: 44877, SpellId: 59970, Name: "Arcanum of Burning Mysteries", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30, stats.MeleeCrit: 20, stats.SpellCrit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3818, ItemId: 44878, SpellId: 59955, Name: "Arcanum of the Stalwart Protector", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 37, stats.Defense: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 3817, ItemId: 44879, SpellId: 59954, Name: "Arcanum of Torment", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50, stats.MeleeCrit: 20, stats.SpellCrit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead},
	{EffectId: 4222, SpellId: 67839, Name: "Mind Amplification Dish", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeHead, RequiredProfession: proto.Profession_Engineering},

	// Shoulder
	{EffectId: 2998, ItemId: 29187, SpellId: 35441, Name: "Inscription of Endurance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 7, stats.FireResistance: 7, stats.FrostResistance: 7, stats.NatureResistance: 7, stats.ShadowResistance: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3793, ItemId: 44067, SpellId: 59771, Name: "Inscription of Triumph", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3794, ItemId: 44068, SpellId: 59773, Name: "Inscription of Dominance", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 23, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3852, ItemId: 44957, SpellId: 62384, Name: "Greater Inscription of the Gladiator", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 30, stats.Resilience: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3806, ItemId: 44129, SpellId: 59927, Name: "Lesser Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.MeleeCrit: 10, stats.SpellCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3807, ItemId: 44130, SpellId: 59928, Name: "Lesser Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 18, stats.MP5: 5}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3875, ItemId: 44131, SpellId: 59929, Name: "Lesser Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 30, stats.RangedAttackPower: 30, stats.MeleeCrit: 10, stats.SpellCrit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3876, ItemId: 44132, SpellId: 59932, Name: "Lesser Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Dodge: 15, stats.Defense: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3810, ItemId: 44874, SpellId: 59937, Name: "Greater Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 24, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3809, ItemId: 44872, SpellId: 59936, Name: "Greater Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.SpellPower: 24, stats.MP5: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3808, ItemId: 44871, SpellId: 59934, Name: "Greater Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3811, ItemId: 44873, SpellId: 59941, Name: "Greater Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Dodge: 20, stats.Defense: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder},
	{EffectId: 3838, SpellId: 61120, Name: "Master's Inscription of the Storm", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 70, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3836, SpellId: 61118, Name: "Master's Inscription of the Crag", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 70, stats.MP5: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3835, SpellId: 61117, Name: "Master's Inscription of the Axe", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 120, stats.RangedAttackPower: 120, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},
	{EffectId: 3837, SpellId: 61119, Name: "Master's Inscription of the Pinnacle", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Dodge: 60, stats.Defense: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeShoulder, RequiredProfession: proto.Profession_Inscription},

	// Back
	{EffectId: 1262, ItemId: 37330, SpellId: 44596, Name: "Superior Arcane Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.ArcaneResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1354, ItemId: 37331, SpellId: 44556, Name: "Superior Fire Resistance", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.FireResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3230, SpellId: 44483, Name: "Superior Frost Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FrostResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1400, SpellId: 44494, Name: "Superior Nature Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.NatureResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1446, SpellId: 44590, Name: "Superior Shadow Resistance", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ShadowResistance: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1951, ItemId: 37347, SpellId: 44591, Name: "Titanweave", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Defense: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3256, ItemId: 37349, SpellId: 44631, Name: "Shadow Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Agility: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3294, ItemId: 44471, SpellId: 47672, Name: "Mighty Armor", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.BonusArmor: 225}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3831, ItemId: 44472, SpellId: 47898, Name: "Greater Speed", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHaste: 23, stats.SpellHaste: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3296, ItemId: 44488, SpellId: 47899, Name: "Wisdom", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Spirit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3243, SpellId: 44582, Name: "Spell Piercing", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPenetration: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3825, SpellId: 60609, Name: "Speed", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHaste: 15, stats.SpellHaste: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 983, SpellId: 44500, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 1099, SpellId: 60663, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack},
	{EffectId: 3605, SpellId: 55002, Name: "Flexweave Underlay", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3722, SpellId: 55642, Name: "Lightweave Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3728, SpellId: 55769, Name: "Darkglow Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3730, SpellId: 55777, Name: "Swordguard Embroidery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3859, SpellId: 63765, Name: "Springy Arachnoweave", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 27}.ToFloatArray(), Type: proto.ItemType_ItemTypeBack, RequiredProfession: proto.Profession_Engineering},

	// Chest
	{EffectId: 3245, ItemId: 37340, SpellId: 44588, Name: "Exceptional Resilience", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Resilience: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3252, SpellId: 44623, Name: "Super Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 8, stats.Strength: 8, stats.Agility: 8, stats.Intellect: 8, stats.Spirit: 8}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3832, ItemId: 44489, SpellId: 60692, Name: "Powerful Stats", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 10, stats.Strength: 10, stats.Agility: 10, stats.Intellect: 10, stats.Spirit: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3233, SpellId: 27958, Name: "Exceptional Mana", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Mana: 250}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3236, SpellId: 44492, Name: "Mighty Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 200}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 3297, SpellId: 47900, Name: "Super Health", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Health: 275}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 2381, SpellId: 44509, Name: "Greater Mana Restoration", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},
	{EffectId: 1953, SpellId: 47766, Name: "Greater Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeChest},

	// Wrist
	{EffectId: 3845, ItemId: 44484, SpellId: 44575, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2332, ItemId: 44498, SpellId: 60767, Name: "Superior Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3850, ItemId: 44944, SpellId: 62256, Name: "Major Stamina", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1119, SpellId: 44555, Name: "Exceptional Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1147, SpellId: 44593, Name: "Major Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3231, SpellId: 44598, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2661, SpellId: 44616, Name: "Greater Stats", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 6, stats.Strength: 6, stats.Agility: 6, stats.Intellect: 6, stats.Spirit: 6}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 2326, SpellId: 44635, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 1600, SpellId: 60616, Name: "Striking", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 38, stats.RangedAttackPower: 38}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist},
	{EffectId: 3756, SpellId: 57683, Name: "Fur Lining - Attack Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 130, stats.RangedAttackPower: 130}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3757, SpellId: 57690, Name: "Fur Lining - Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 102}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3758, SpellId: 57691, Name: "Fur Lining - Spell Power", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 76}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3759, SpellId: 57692, Name: "Fur Lining - Fire Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FireResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3760, SpellId: 57694, Name: "Fur Lining - Frost Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.FrostResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3761, SpellId: 57696, Name: "Fur Lining - Shadow Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ShadowResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3762, SpellId: 57699, Name: "Fur Lining - Nature Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.NatureResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3763, SpellId: 57701, Name: "Fur Lining - Arcane Resist", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.ArcaneResistance: 70}.ToFloatArray(), Type: proto.ItemType_ItemTypeWrist, RequiredProfession: proto.Profession_Leatherworking},

	// Hands
	{EffectId: 3253, ItemId: 44485, SpellId: 44625, Name: "Armsman", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Parry: 10}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 1603, SpellId: 60668, Name: "Crusher", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3246, SpellId: 44592, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 28}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3231, SpellId: 44484, Name: "Expertise", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Expertise: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3238, SpellId: 44506, Name: "Gatherer", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3829, SpellId: 44513, Name: "Greater Assult", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 35, stats.RangedAttackPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3222, SpellId: 44529, Name: "Major Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3234, SpellId: 44488, Name: "Precision", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 20, stats.SpellHit: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands},
	{EffectId: 3603, SpellId: 54998, Name: "Hand-Mounted Pyro Rocket", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3604, SpellId: 54999, Name: "Hyperspeed Accelerators", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3860, SpellId: 63770, Name: "Reticulated Armor Webbing", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.BonusArmor: 885}.ToFloatArray(), Type: proto.ItemType_ItemTypeHands, RequiredProfession: proto.Profession_Engineering},

	// Waist
	{EffectId: 3599, SpellId: 54736, Name: "Personal Electromagnetic Pulse Generator", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},
	{EffectId: 3601, SpellId: 54793, Name: "Frag Belt", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWaist, RequiredProfession: proto.Profession_Engineering},

	// Legs
	{EffectId: 3325, ItemId: 38371, SpellId: 50901, Name: "Jormungar Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 45, stats.Agility: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3326, ItemId: 38372, SpellId: 50902, Name: "Nerubian Leg Armor", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 55, stats.RangedAttackPower: 55, stats.MeleeCrit: 15, stats.SpellCrit: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3822, ItemId: 38373, SpellId: 60581, Name: "Frosthide Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 55, stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3823, ItemId: 38374, SpellId: 60582, Name: "Icescale Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3853, ItemId: 44963, SpellId: 62447, Name: "Earthen Leg Armor", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 28, stats.Resilience: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3718, ItemId: 41601, SpellId: 55630, Name: "Shining Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Spirit: 12, stats.SpellPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3719, ItemId: 41602, SpellId: 55631, Name: "Brilliant Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3720, ItemId: 41603, SpellId: 55632, Name: "Azure Spellthread", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.Stamina: 20, stats.SpellPower: 35}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3721, ItemId: 41604, SpellId: 55634, Name: "Sapphire Spellthread", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs},
	{EffectId: 3327, SpellId: 60583, Name: "Jormungar Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 55, stats.Agility: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3328, SpellId: 60584, Name: "Nerubian Leg Reinforcements", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 75, stats.RangedAttackPower: 75, stats.MeleeCrit: 22, stats.SpellCrit: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Leatherworking},
	{EffectId: 3873, SpellId: 56034, Name: "Master's Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 30, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},
	{EffectId: 3872, SpellId: 56039, Name: "Sanctified Spellthread", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 20, stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeLegs, RequiredProfession: proto.Profession_Tailoring},

	// Feet
	{EffectId: 1597, ItemId: 44490, SpellId: 60763, Name: "Greater Assault", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3232, ItemId: 44491, SpellId: 47901, Name: "Tuskarr's Vitality", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.Stamina: 15}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3824, SpellId: 60606, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 24, stats.RangedAttackPower: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 1075, SpellId: 44528, Name: "Greater Fortitude", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 22}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 1147, SpellId: 44508, Name: "Greater Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 18}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3244, SpellId: 44584, Name: "Greater Vitality", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MP5: 7}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 3826, SpellId: 60623, Name: "Icewalker", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeHit: 12, stats.SpellHit: 12, stats.MeleeCrit: 12, stats.SpellCrit: 12}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 983, SpellId: 44589, Name: "Superior Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 16}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet},
	{EffectId: 4223, SpellId: 55016, Name: "Nitro Boosts", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.MeleeCrit: 24, stats.SpellCrit: 24}.ToFloatArray(), Type: proto.ItemType_ItemTypeFeet, RequiredProfession: proto.Profession_Engineering},

	// Weapon
	{EffectId: 1103, SpellId: 44633, Name: "Exceptional Agility", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Agility: 26}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3844, SpellId: 44510, Name: "Exceptional Spirit", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Spirit: 45}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3251, ItemId: 37339, SpellId: 44621, Name: "Giant Slayer", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3239, ItemId: 37344, SpellId: 44524, Name: "Icebreaker", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3731, ItemId: 41976, SpellId: 55836, Name: "Titanium Weapon Chain", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{stats.MeleeHit: 28, stats.SpellHit: 28}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3833, ItemId: 44486, SpellId: 60707, Name: "Superior Potency", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 65, stats.RangedAttackPower: 65}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3834, ItemId: 44487, SpellId: 60714, Name: "Mighty Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 63}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3789, ItemId: 44492, SpellId: 59621, Name: "Berserking", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3241, ItemId: 44494, SpellId: 44576, Name: "Lifeward", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3790, ItemId: 44495, SpellId: 59625, Name: "Black Magic", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3788, ItemId: 44496, SpellId: 59619, Name: "Accuracy", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.MeleeHit: 25, stats.SpellHit: 25, stats.MeleeCrit: 25, stats.SpellCrit: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3830, SpellId: 44629, Name: "Exceptional Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 1606, SpellId: 60621, Name: "Greater Potency", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 50, stats.RangedAttackPower: 50}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3370, SpellId: 53343, Name: "Rune of Razorice", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3369, SpellId: 53341, Name: "Rune of Cinderglacier", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3366, SpellId: 53331, Name: "Rune of Lichbane", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3595, SpellId: 54447, Name: "Rune of Spellbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3594, SpellId: 54446, Name: "Rune of Swordbreaking", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3368, SpellId: 53344, Name: "Rune of the Fallen Crusader", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3870, ItemId: 46348, SpellId: 64579, Name: "Blood Draining", Quality: proto.ItemQuality_ItemQualityEpic, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon},
	{EffectId: 3883, SpellId: 70164, Name: "Rune of the Nerubian Carapace", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},

	// 2H Weapon
	{EffectId: 3247, ItemId: 44473, SpellId: 44595, Name: "Scourgebane", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3827, ItemId: 44483, SpellId: 60691, Name: "Massacre", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.AttackPower: 110, stats.RangedAttackPower: 110}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3828, SpellId: 44630, Name: "Greater Savagery", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 85, stats.RangedAttackPower: 85}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeTwoHand},
	{EffectId: 3854, ItemId: 45059, SpellId: 62948, Name: "Staff - Greater Spellpower", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.SpellPower: 81}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeStaff},
	{EffectId: 3367, SpellId: 53342, Name: "Rune of Spellshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3365, SpellId: 53323, Name: "Rune of Swordshattering", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},
	{EffectId: 3847, SpellId: 62158, Name: "Rune of the Stoneskin Gargoyle", Phase: 1, Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, ClassAllowlist: []proto.Class{proto.Class_ClassDeathKnight}},

	// Shield
	{EffectId: 1952, SpellId: 44489, Name: "Defense", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Defense: 20}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 1128, SpellId: 60653, Name: "Greater Intellect", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Intellect: 25}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 3748, ItemId: 42500, SpellId: 56353, Name: "Titanium Shield Spike", Quality: proto.ItemQuality_ItemQualityUncommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},
	{EffectId: 3849, ItemId: 44936, SpellId: 62201, Name: "Titanium Plating", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{stats.BlockValue: 81}.ToFloatArray(), Type: proto.ItemType_ItemTypeWeapon, EnchantType: proto.EnchantType_EnchantTypeShield},

	// Ring
	{EffectId: 3839, SpellId: 44645, Name: "Assault", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.AttackPower: 40, stats.RangedAttackPower: 40}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 3840, SpellId: 44636, Name: "Greater Spellpower", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.SpellPower: 23}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},
	{EffectId: 3791, SpellId: 59636, Name: "Stamina", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{stats.Stamina: 30}.ToFloatArray(), Type: proto.ItemType_ItemTypeFinger, RequiredProfession: proto.Profession_Enchanting},

	// Ranged
	{EffectId: 3607, ItemId: 41146, SpellId: 55076, Name: "Sun Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 3608, ItemId: 41167, SpellId: 55135, Name: "Heartseeker Scope", Quality: proto.ItemQuality_ItemQualityRare, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
	{EffectId: 3843, ItemId: 44739, SpellId: 61468, Name: "Diamond-cut Refractor Scope", Quality: proto.ItemQuality_ItemQualityCommon, Stats: stats.Stats{}.ToFloatArray(), Type: proto.ItemType_ItemTypeRanged},
}
