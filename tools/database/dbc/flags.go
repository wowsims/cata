package dbc

import "github.com/wowsims/mop/sim/core/proto"

type ItemStaticFlags0 uint32

const (
	NO_PICKUP                         ItemStaticFlags0 = 0x00000001
	CONJURED                          ItemStaticFlags0 = 0x00000002
	HAS_LOOT_TABLE                    ItemStaticFlags0 = 0x00000004
	HEROIC_TOOLTIP                    ItemStaticFlags0 = 0x00000008
	DEPRECATED                        ItemStaticFlags0 = 0x00000010
	NO_USER_DESTROY                   ItemStaticFlags0 = 0x00000020
	PLAYER_CAST                       ItemStaticFlags0 = 0x00000040
	NO_EQUIP_COOLDOWN                 ItemStaticFlags0 = 0x00000080
	MULTI_LOOT_QUEST                  ItemStaticFlags0 = 0x00000100
	GIFT_WRAP                         ItemStaticFlags0 = 0x00000200
	USES_RESOURCES                    ItemStaticFlags0 = 0x00000400
	MULTI_DROP                        ItemStaticFlags0 = 0x00000800
	IN_GAME_REFUND                    ItemStaticFlags0 = 0x00001000
	PETITION                          ItemStaticFlags0 = 0x00002000
	HAS_TEXT                          ItemStaticFlags0 = 0x00004000
	NO_DISENCHANT                     ItemStaticFlags0 = 0x00008000
	REAL_DURATION                     ItemStaticFlags0 = 0x00010000
	NO_CREATOR                        ItemStaticFlags0 = 0x00020000
	PROSPECTABLE                      ItemStaticFlags0 = 0x00040000
	UNIQUE_EQUIPPABLE                 ItemStaticFlags0 = 0x00080000
	DISABLE_AURA_QUOTAS               ItemStaticFlags0 = 0x00100000
	IGNORE_DEFAULT_ARENA_RESTRICTIONS ItemStaticFlags0 = 0x00200000
	NO_DURABILITY_LOSS                ItemStaticFlags0 = 0x00400000
	USEABLE_WHILE_SHAPESHIFTED        ItemStaticFlags0 = 0x00800000
	HAS_QUEST_GLOW                    ItemStaticFlags0 = 0x01000000
	HIDE_UNUSABLE_RECIPE              ItemStaticFlags0 = 0x02000000
	NOT_USABLE_IN_ARENA               ItemStaticFlags0 = 0x04000000
	BOUND_TO_ACCOUNT                  ItemStaticFlags0 = 0x08000000
	NO_REAGENT_COST                   ItemStaticFlags0 = 0x10000000
	MILLABLE                          ItemStaticFlags0 = 0x20000000
	REPORT_TO_GUILD_CHAT              ItemStaticFlags0 = 0x40000000
	DONT_USE_DYNAMIC_DROP_CHANCE      ItemStaticFlags0 = 0x80000000
)

func (f ItemStaticFlags0) Has(flag ItemStaticFlags0) bool {
	return f&flag != 0
}

type SpellSchool int

const (
	PHYSICAL SpellSchool = 1 << iota // 0x1
	HOLY                             // 0x2
	FIRE                             // 0x4
	NATURE                           // 0x8
	FROST                            // 0x10
	SHADOW                           // 0x20
	ARCANE                           // 0x40

	SPELL_PENETRATION = FIRE | NATURE | FROST | SHADOW | ARCANE // 0x7E
)

func (f SpellSchool) Has(flag SpellSchool) bool {
	return f&flag != 0
}

type ItemStaticFlags1 uint32

const (
	HORDE_SPECIFIC                        ItemStaticFlags1 = 0x00000001
	ALLIANCE_SPECIFIC                     ItemStaticFlags1 = 0x00000002
	DONT_IGNORE_BUY_PRICE                 ItemStaticFlags1 = 0x00000004
	ONLY_CASTER_ROLL_NEED                 ItemStaticFlags1 = 0x00000008
	ONLY_NON_CASTER_ROLL_NEED             ItemStaticFlags1 = 0x00000010
	EVERYONE_CAN_ROLL_NEED                ItemStaticFlags1 = 0x00000020
	CANNOT_TRADE_BIND_ON_PICKUP           ItemStaticFlags1 = 0x00000040
	CAN_TRADE_BIND_ON_PICKUP              ItemStaticFlags1 = 0x00000080
	CAN_ONLY_ROLL_GREED                   ItemStaticFlags1 = 0x00000100
	CASTER_WEAPON                         ItemStaticFlags1 = 0x00000200
	DELETE_ON_LOGIN                       ItemStaticFlags1 = 0x00000400
	INTERNAL_ITEM                         ItemStaticFlags1 = 0x00000800
	NO_VENDOR_VALUE                       ItemStaticFlags1 = 0x00001000
	SHOW_BEFORE_DISCOVERED                ItemStaticFlags1 = 0x00002000
	OVERRIDE_GOLD_COST                    ItemStaticFlags1 = 0x00004000
	IGNORE_DEFAULT_RATED_BG_RESTRICTIONS  ItemStaticFlags1 = 0x00008000
	NOT_USABLE_IN_RATED_BG                ItemStaticFlags1 = 0x00010000
	BNET_ACCOUNT_TRADE_OK                 ItemStaticFlags1 = 0x00020000
	CONFIRM_BEFORE_USE                    ItemStaticFlags1 = 0x00040000
	REEVALUATE_BONDING_ON_TRANSFORM       ItemStaticFlags1 = 0x00080000
	NO_TRANSFORM_ON_CHARGE_DEPLETION      ItemStaticFlags1 = 0x00100000
	NO_ALTER_ITEM_VISUAL                  ItemStaticFlags1 = 0x00200000
	NO_SOURCE_FOR_ITEM_VISUAL             ItemStaticFlags1 = 0x00400000
	IGNORE_QUALITY_FOR_ITEM_VISUAL_SOURCE ItemStaticFlags1 = 0x00800000
	NO_DURABILITY                         ItemStaticFlags1 = 0x01000000
	ROLE_TANK                             ItemStaticFlags1 = 0x02000000
	ROLE_HEALER                           ItemStaticFlags1 = 0x04000000
	ROLE_DAMAGE                           ItemStaticFlags1 = 0x08000000
	CAN_DROP_IN_CHALLENGE_MODE            ItemStaticFlags1 = 0x10000000
	NEVER_STACK_IN_LOOT_UI                ItemStaticFlags1 = 0x20000000
	DISENCHANT_TO_LOOT_TABLE              ItemStaticFlags1 = 0x40000000
	CAN_BE_PLACED_IN_REAGENT_BANK         ItemStaticFlags1 = 0x80000000
)

func (f ItemStaticFlags1) Has(flag ItemStaticFlags1) bool {
	return f&flag != 0
}

type ItemStaticFlags2 uint32

const (
	DONT_DESTROY_ON_QUEST_ACCEPT                        ItemStaticFlags2 = 0x00000001
	CAN_BE_UPGRADED                                     ItemStaticFlags2 = 0x00000002
	UPGRADE_FROM_ITEM_OVERRIDES_DROP_UPGRADE            ItemStaticFlags2 = 0x00000004
	ALWAYS_FREE_FOR_ALL_IN_LOOT                         ItemStaticFlags2 = 0x00000008
	HIDE_ITEM_UPGRADES_IF_NOT_UPGRADED                  ItemStaticFlags2 = 0x00000010
	UPDATE_NPC_INTERACTIONS_WHEN_PICKED_UP              ItemStaticFlags2 = 0x00000020
	DOESNT_LEAVE_PROGRESSIVE_WIN_HISTORY                ItemStaticFlags2 = 0x00000040
	IGNORE_ITEM_HISTORY_TRACKER                         ItemStaticFlags2 = 0x00000080
	IGNORE_ITEM_LEVEL_CAP_IN_PVP                        ItemStaticFlags2 = 0x00000100
	DISPLAY_AS_HEIRLOOM                                 ItemStaticFlags2 = 0x00000200
	SKIP_USE_CHECK_ON_PICKUP                            ItemStaticFlags2 = 0x00000400
	NO_LOOT_OVERFLOW_MAIL                               ItemStaticFlags2 = 0x00000800
	DONT_DISPLAY_IN_GUILD_NEWS                          ItemStaticFlags2 = 0x00001000
	TRIAL_OF_THE_GLADIATOR_GEAR                         ItemStaticFlags2 = 0x00002000
	REQUIRES_STACK_CHANGE_LOG                           ItemStaticFlags2 = 0x00004000
	TOY                                                 ItemStaticFlags2 = 0x00008000
	SUPPRESS_NAME_SUFFIXES                              ItemStaticFlags2 = 0x00010000
	PUSH_LOOT                                           ItemStaticFlags2 = 0x00020000
	DONT_REPORT_LOOT_LOG_TO_PARTY                       ItemStaticFlags2 = 0x00040000
	ALWAYS_ALLOW_DUAL_WIELD                             ItemStaticFlags2 = 0x00080000
	OBLITERATABLE                                       ItemStaticFlags2 = 0x00100000
	ACTS_AS_TRANSMOG_HIDDEN_VISUAL_OPTION               ItemStaticFlags2 = 0x00200000
	EXPIRE_ON_WEEKLY_RESET                              ItemStaticFlags2 = 0x00400000
	DOESNT_SHOW_UP_IN_TRANSMOG_UI_UNTIL_COLLECTED       ItemStaticFlags2 = 0x00800000
	CAN_STORE_ENCHANTS                                  ItemStaticFlags2 = 0x01000000
	HIDE_QUEST_ITEM_FROM_OBJECT_TOOLTIP                 ItemStaticFlags2 = 0x02000000
	DO_NOT_TOAST                                        ItemStaticFlags2 = 0x04000000
	IGNORE_CREATION_CONTEXT_FOR_PROGRESSIVE_WIN_HISTORY ItemStaticFlags2 = 0x08000000
	FORCE_ALL_SPECS_FOR_ITEM_HISTORY                    ItemStaticFlags2 = 0x10000000
	SAVE_AFTER_CONSUME                                  ItemStaticFlags2 = 0x20000000
	LOOT_CONTAINER_SAVES_PLAYER_STATE                   ItemStaticFlags2 = 0x40000000
	NO_VOID_STORAGE                                     ItemStaticFlags2 = 0x80000000
)

func (f ItemStaticFlags2) Has(flag ItemStaticFlags2) bool {
	return f&flag != 0
}

type ItemStaticFlags3 uint32

const (
	IMMEDIATELY_TRIGGER_ON_USE_BINDING_EFFECTS    ItemStaticFlags3 = 0x00000001
	ALWAYS_DISPLAY_ITEM_LEVEL_IN_TOOLTIP          ItemStaticFlags3 = 0x00000002
	DISPLAY_RANDOM_ADDITIONAL_STATS_IN_TOOLTIP    ItemStaticFlags3 = 0x00000004
	ACTIVATE_ON_EQUIP_EFFECTS_WHEN_TRANSMOGRIFIED ItemStaticFlags3 = 0x00000008
	ENFORCE_TRANSMOG_WITH_CHILD_ITEM              ItemStaticFlags3 = 0x00000010
	SCRAPABLE                                     ItemStaticFlags3 = 0x00000020
	BYPASS_REP_REQUIREMENTS_FOR_TRANSMOG          ItemStaticFlags3 = 0x00000040
	DISPLAY_ONLY_ON_DEFINED_RACES                 ItemStaticFlags3 = 0x00000080
	REGULATED_COMMODITY                           ItemStaticFlags3 = 0x00000100
	CREATE_LOOT_IMMEDIATELY                       ItemStaticFlags3 = 0x00000200
	GENERATE_LOOT_SPEC_ITEM                       ItemStaticFlags3 = 0x00000400
	HIDDEN_IN_REWARD_SUMMARIES                    ItemStaticFlags3 = 0x00000800
	DISALLOW_WHILE_LEVEL_LINKED                   ItemStaticFlags3 = 0x00001000
	DISALLOW_ENCHANT                              ItemStaticFlags3 = 0x00002000
	SQUISH_USING_ITEM_LEVEL_AS_PLAYER_LEVEL       ItemStaticFlags3 = 0x00004000
	ALWAYS_SHOW_SELL_PRICE_IN_TOOLTIP             ItemStaticFlags3 = 0x00008000
	COSMETIC_ITEM                                 ItemStaticFlags3 = 0x00010000
	NO_SPELL_EFFECT_TOOLTIP_PREFIXES              ItemStaticFlags3 = 0x00020000
	IGNORE_COSMETIC_COLLECTION_BEHAVIOR           ItemStaticFlags3 = 0x00040000
	NPC_ONLY                                      ItemStaticFlags3 = 0x00080000
	NOT_RESTORABLE                                ItemStaticFlags3 = 0x00100000
	DONT_DISPLAY_AS_CRAFTING_REAGENT              ItemStaticFlags3 = 0x00200000
	DISPLAY_REAGENT_QUALITY_AS_CRAFTED_QUALITY    ItemStaticFlags3 = 0x00400000
	NO_SALVAGE                                    ItemStaticFlags3 = 0x00800000
	RECRAFTABLE                                   ItemStaticFlags3 = 0x01000000
	CC_TRINKET                                    ItemStaticFlags3 = 0x02000000
	KEEP_THROUGH_FACTION_CHANGE                   ItemStaticFlags3 = 0x04000000
	NOT_MULTICRAFTABLE                            ItemStaticFlags3 = 0x08000000
	DONT_REPORT_LOOT_LOG_TO_SELF                  ItemStaticFlags3 = 0x10000000
	SEND_TELEMETRY_ON_USE                         ItemStaticFlags3 = 0x20000000
)

func (f ItemStaticFlags3) Has(flag ItemStaticFlags3) bool {
	return f&flag != 0
}

type InventoryTypeFlag uint32

const (
	HEAD             InventoryTypeFlag = 0x2
	NECK             InventoryTypeFlag = 0x4
	SHOULDER         InventoryTypeFlag = 0x8
	BODY             InventoryTypeFlag = 0x10
	CHEST            InventoryTypeFlag = 0x20
	WAIST            InventoryTypeFlag = 0x40
	LEGS             InventoryTypeFlag = 0x80
	FEET             InventoryTypeFlag = 0x100
	WRIST            InventoryTypeFlag = 0x200
	HAND             InventoryTypeFlag = 0x400
	FINGER           InventoryTypeFlag = 0x800
	TRINKET          InventoryTypeFlag = 0x1000
	MAIN_HAND        InventoryTypeFlag = 0x2000
	OFF_HAND         InventoryTypeFlag = 0x4000
	RANGED           InventoryTypeFlag = 0x8000
	CLOAK            InventoryTypeFlag = 0x10000
	TWO_H_WEAPON     InventoryTypeFlag = 0x20000
	BAG              InventoryTypeFlag = 0x40000
	TABARD           InventoryTypeFlag = 0x80000
	ROBE             InventoryTypeFlag = 0x100000
	WEAPON_MAIN_HAND InventoryTypeFlag = 0x200000
	WEAPON_OFF_HAND  InventoryTypeFlag = 0x400000
	HOLDABLE         InventoryTypeFlag = 0x800000
	AMMO             InventoryTypeFlag = 0x1000000
	THROWN           InventoryTypeFlag = 0x2000000
	RANGED_RIGHT     InventoryTypeFlag = 0x4000000
	QUIVER           InventoryTypeFlag = 0x8000000
	RELIC            InventoryTypeFlag = 0x10000000
)

func (i InventoryTypeFlag) Has(flag InventoryTypeFlag) bool {
	return i&flag != 0
}

type GemType int

const (
	Meta   GemType = 0x1
	Red    GemType = 0x2
	Yellow GemType = 0x4
	Blue   GemType = 0x8
	// Combined colors:
	Orange     GemType = Red | Yellow        // 0x6
	Purple     GemType = Red | Blue          // 0xa
	Green      GemType = Yellow | Blue       // 0xc
	Prismatic  GemType = Red | Yellow | Blue // 0xe
	ShaTouched GemType = 0x10
	Cogwheel   GemType = 0x20
)

func (gem GemType) ToProto() proto.GemColor {
	switch gem {
	case Meta:
		return proto.GemColor_GemColorMeta
	case Red:
		return proto.GemColor_GemColorRed
	case Yellow:
		return proto.GemColor_GemColorYellow
	case Blue:
		return proto.GemColor_GemColorBlue
	case Orange:
		return proto.GemColor_GemColorOrange
	case Purple:
		return proto.GemColor_GemColorPurple
	case Green:
		return proto.GemColor_GemColorGreen
	case Prismatic:
		return proto.GemColor_GemColorPrismatic
	case ShaTouched:
		return proto.GemColor_GemColorShaTouched
	case Cogwheel:
		return proto.GemColor_GemColorCogwheel
	default:
		return proto.GemColor_GemColorUnknown
	}
}

func (i GemType) Has(flag GemType) bool {
	return i&flag != 0
}

type ClassMask uint32

const (
	WARRIOR      ClassMask = 0x1
	PALADIN      ClassMask = 0x2
	HUNTER       ClassMask = 0x4
	ROGUE        ClassMask = 0x8
	PRIEST       ClassMask = 0x10
	DEATH_KNIGHT ClassMask = 0x20
	SHAMAN       ClassMask = 0x40
	MAGE         ClassMask = 0x80
	WARLOCK      ClassMask = 0x100
	MONK         ClassMask = 0x200
	DRUID        ClassMask = 0x400
	DEMON_HUNTER ClassMask = 0x800
	EVOKER       ClassMask = 0x1000
)

func (c ClassMask) Has(flag ClassMask) bool {
	return c&flag != 0
}

// --- Weapon constants (EquippedItemClass = 2) ---
const (
	ITEM_SUBCLASS_BIT_WEAPON_NONE         = 0x00000000
	ITEM_SUBCLASS_BIT_WEAPON_1H_AXE       = 0x00000001
	ITEM_SUBCLASS_BIT_WEAPON_2H_AXE       = 0x00000002
	ITEM_SUBCLASS_BIT_WEAPON_BOW          = 0x00000004
	ITEM_SUBCLASS_BIT_WEAPON_GUNS         = 0x00000008
	ITEM_SUBCLASS_BIT_WEAPON_MACE_1H      = 0x00000010
	ITEM_SUBCLASS_BIT_WEAPON_MACE_2H      = 0x00000020
	ITEM_SUBCLASS_BIT_WEAPON_POLEARM      = 0x00000040
	ITEM_SUBCLASS_BIT_WEAPON_SWORD_1H     = 0x00000080
	ITEM_SUBCLASS_BIT_WEAPON_SWORD_2H     = 0x00000100
	ITEM_SUBCLASS_BIT_WEAPON_OBSOLETE     = 0x00000200
	ITEM_SUBCLASS_BIT_WEAPON_STAFF        = 0x00000400
	ITEM_SUBCLASS_BIT_WEAPON_1H_EXOTIC    = 0x00000800
	ITEM_SUBCLASS_BIT_WEAPON_2H_EXOTIC    = 0x00001000
	ITEM_SUBCLASS_BIT_WEAPON_FIST         = 0x00002000
	ITEM_SUBCLASS_BIT_WEAPON_MISC         = 0x00004000
	ITEM_SUBCLASS_BIT_WEAPON_DAGGERS      = 0x00008000
	ITEM_SUBCLASS_BIT_WEAPON_THROWN       = 0x00010000
	ITEM_SUBCLASS_BIT_WEAPON_SPEAR        = 0x00020000
	ITEM_SUBCLASS_BIT_WEAPON_CROSSBOW     = 0x00040000
	ITEM_SUBCLASS_BIT_WEAPON_WAND         = 0x00080000
	ITEM_SUBCLASS_BIT_WEAPON_FISHING_POLE = 0x00100000
)

// --- Armor constants (EquippedItemClass = 4) ---
const (
	ITEM_SUBCLASS_BIT_ARMOR_NONE    = 0x00000000
	ITEM_SUBCLASS_BIT_ARMOR_MISC    = 0x00000001
	ITEM_SUBCLASS_BIT_ARMOR_CLOTH   = 0x00000002
	ITEM_SUBCLASS_BIT_ARMOR_LEATHER = 0x00000004
	ITEM_SUBCLASS_BIT_ARMOR_MAIL    = 0x00000008
	ITEM_SUBCLASS_BIT_ARMOR_PLATE   = 0x00000010
	ITEM_SUBCLASS_BIT_ARMOR_BUCKLER = 0x00000020 // (Unused)
	ITEM_SUBCLASS_BIT_ARMOR_SHIELD  = 0x00000040
	ITEM_SUBCLASS_BIT_ARMOR_LIBRAM  = 0x00000080
	ITEM_SUBCLASS_BIT_ARMOR_IDOL    = 0x00000100
	ITEM_SUBCLASS_BIT_ARMOR_TOTEM   = 0x00000200
	ITEM_SUBCLASS_BIT_ARMOR_SIGIL   = 0x00000400
)
const (
	rangedMask  = ITEM_SUBCLASS_BIT_WEAPON_BOW | ITEM_SUBCLASS_BIT_WEAPON_GUNS | ITEM_SUBCLASS_BIT_WEAPON_CROSSBOW
	twoHandMask = ITEM_SUBCLASS_BIT_WEAPON_2H_AXE | ITEM_SUBCLASS_BIT_WEAPON_MACE_2H |
		ITEM_SUBCLASS_BIT_WEAPON_SWORD_2H | ITEM_SUBCLASS_BIT_WEAPON_2H_EXOTIC |
		ITEM_SUBCLASS_BIT_WEAPON_SPEAR | ITEM_SUBCLASS_BIT_WEAPON_POLEARM
)

const (
	OffHandValue = 65
	ShieldValue1 = 96
	ShieldValue2 = 64
)

type RatingModType uint

const (
	RATING_MOD_DODGE        = 0x00000004
	RATING_MOD_PARRY        = 0x00000008
	RATING_MOD_HIT_MELEE    = 0x00000020
	RATING_MOD_HIT_RANGED   = 0x00000040
	RATING_MOD_HIT_SPELL    = 0x00000080
	RATING_MOD_CRIT_MELEE   = 0x00000100
	RATING_MOD_CRIT_RANGED  = 0x00000200
	RATING_MOD_CRIT_SPELL   = 0x00000400
	RATING_MOD_MULTISTRIKE  = 0x00000800
	RATING_MOD_READINESS    = 0x00001000
	RATING_MOD_SPEED        = 0x00002000
	RATING_MOD_RESILIENCE   = 0x00008000
	RATING_MOD_LEECH        = 0x00010000
	RATING_MOD_HASTE_MELEE  = 0x00020000
	RATING_MOD_HASTE_RANGED = 0x00040000
	RATING_MOD_HASTE_SPELL  = 0x00080000
	RATING_MOD_AVOIDANCE    = 0x00100000
	RATING_MOD_EXPERTISE    = 0x00800000
	RATING_MOD_MASTERY      = 0x02000000
	RATING_MOD_PVP_POWER    = 0x04000000

	RATING_MOD_VERS_DAMAGE = 0x10000000
	RATING_MOD_VERS_HEAL   = 0x20000000
	RATING_MOD_VERS_MITIG  = 0x40000000
)
