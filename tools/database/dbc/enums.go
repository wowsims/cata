package dbc

import (
	"github.com/wowsims/mop/sim/core/proto"
)

const (
	ITEM_ENCHANTMENT_NONE             int = 0
	ITEM_ENCHANTMENT_COMBAT_SPELL     int = 1
	ITEM_ENCHANTMENT_DAMAGE           int = 2
	ITEM_ENCHANTMENT_EQUIP_SPELL      int = 3
	ITEM_ENCHANTMENT_RESISTANCE       int = 4
	ITEM_ENCHANTMENT_STAT             int = 5
	ITEM_ENCHANTMENT_TOTEM            int = 6
	ITEM_ENCHANTMENT_USE_SPELL        int = 7
	ITEM_ENCHANTMENT_PRISMATIC_SOCKET int = 8
	ITEM_ENCHANTMENT_RELIC_RANK       int = 9
	ITEM_ENCHANTMENT_APPLY_BONUS      int = 11
	ITEM_ENCHANTMENT_RELIC_EVIL       int = 12 // Scaling relic +ilevel, see enchant::initialize_relic
)

const (
	ITEM_SPELLTRIGGER_ON_USE          int = 0 // use after equip cooldown
	ITEM_SPELLTRIGGER_ON_EQUIP        int = 1
	ITEM_SPELLTRIGGER_CHANCE_ON_HIT   int = 2
	ITEM_SPELLTRIGGER_SOULSTONE       int = 4
	ITEM_SPELLTRIGGER_ON_NO_DELAY_USE int = 5 // no equip cooldown
	ITEM_SPELLTRIGGER_LEARN_SPELL_ID  int = 6
)

type SpellEffectType int

const (
	E_INSTAKILL                                                SpellEffectType = 1
	E_SCHOOL_DAMAGE                                            SpellEffectType = 2
	E_DUMMY                                                    SpellEffectType = 3
	E_PORTAL_TELEPORT                                          SpellEffectType = 4
	E_UNK_ITEM_MOD                                             SpellEffectType = 5
	E_APPLY_AURA                                               SpellEffectType = 6
	E_ENVIRONMENTAL_DAMAGE                                     SpellEffectType = 7
	E_POWER_DRAIN                                              SpellEffectType = 8
	E_HEALTH_LEECH                                             SpellEffectType = 9
	E_HEAL                                                     SpellEffectType = 10
	E_BIND                                                     SpellEffectType = 11
	E_PORTAL                                                   SpellEffectType = 12
	E_RITUAL_BASE                                              SpellEffectType = 13
	E_INCREASE_CURRENCY_CAP                                    SpellEffectType = 14
	E_RITUAL_ACTIVATE_PORTAL                                   SpellEffectType = 15
	E_QUEST_COMPLETE                                           SpellEffectType = 16
	E_WEAPON_DAMAGE_NOSCHOOL                                   SpellEffectType = 17
	E_RESURRECT                                                SpellEffectType = 18
	E_ADD_EXTRA_ATTACKS                                        SpellEffectType = 19
	E_DODGE                                                    SpellEffectType = 20
	E_EVADE                                                    SpellEffectType = 21
	E_PARRY                                                    SpellEffectType = 22
	E_BLOCK                                                    SpellEffectType = 23
	E_CREATE_ITEM                                              SpellEffectType = 24
	E_WEAPON                                                   SpellEffectType = 25
	E_DEFENSE                                                  SpellEffectType = 26
	E_PERSISTENT_AREA_AURA                                     SpellEffectType = 27
	E_SUMMON                                                   SpellEffectType = 28
	E_LEAP                                                     SpellEffectType = 29
	E_ENERGIZE                                                 SpellEffectType = 30
	E_WEAPON_PERCENT_DAMAGE                                    SpellEffectType = 31
	E_TRIGGER_MISSILE                                          SpellEffectType = 32
	E_OPEN_LOCK                                                SpellEffectType = 33
	E_SUMMON_CHANGE_ITEM                                       SpellEffectType = 34
	E_APPLY_AREA_AURA_PARTY                                    SpellEffectType = 35
	E_LEARN_SPELL                                              SpellEffectType = 36
	E_SPELL_DEFENSE                                            SpellEffectType = 37
	E_DISPEL                                                   SpellEffectType = 38
	E_LANGUAGE                                                 SpellEffectType = 39
	E_DUAL_WIELD                                               SpellEffectType = 40
	E_JUMP                                                     SpellEffectType = 41
	E_JUMP_DEST                                                SpellEffectType = 42
	E_TELEPORT_UNITS_FACE_CASTER                               SpellEffectType = 43
	E_SKILL_STEP                                               SpellEffectType = 44
	E_PLAY_MOVIE                                               SpellEffectType = 45
	E_SPAWN                                                    SpellEffectType = 46
	E_TRADE_SKILL                                              SpellEffectType = 47
	E_STEALTH                                                  SpellEffectType = 48
	E_DETECT                                                   SpellEffectType = 49
	E_TRANS_DOOR                                               SpellEffectType = 50
	E_FORCE_CRITICAL_HIT                                       SpellEffectType = 51
	E_SET_MAX_BATTLE_PET_COUNT                                 SpellEffectType = 52
	E_ENCHANT_ITEM                                             SpellEffectType = 53
	E_ENCHANT_ITEM_TEMPORARY                                   SpellEffectType = 54
	E_TAMECREATURE                                             SpellEffectType = 55
	E_SUMMON_PET                                               SpellEffectType = 56
	E_LEARN_PET_SPELL                                          SpellEffectType = 57
	E_WEAPON_DAMAGE                                            SpellEffectType = 58
	E_CREATE_RANDOM_ITEM                                       SpellEffectType = 59
	E_PROFICIENCY                                              SpellEffectType = 60
	E_SEND_EVENT                                               SpellEffectType = 61
	E_POWER_BURN                                               SpellEffectType = 62
	E_THREAT                                                   SpellEffectType = 63
	E_TRIGGER_SPELL                                            SpellEffectType = 64
	E_APPLY_AREA_AURA_RAID                                     SpellEffectType = 65
	E_RECHARGE_ITEM                                            SpellEffectType = 66
	E_HEAL_MAX_HEALTH                                          SpellEffectType = 67
	E_INTERRUPT_CAST                                           SpellEffectType = 68
	E_DISTRACT                                                 SpellEffectType = 69
	E_PULL                                                     SpellEffectType = 70
	E_PICKPOCKET                                               SpellEffectType = 71
	E_ADD_FARSIGHT                                             SpellEffectType = 72
	E_UNTRAIN_TALENTS                                          SpellEffectType = 73
	E_APPLY_GLYPH                                              SpellEffectType = 74
	E_HEAL_MECHANICAL                                          SpellEffectType = 75
	E_SUMMON_OBJECT_WILD                                       SpellEffectType = 76
	E_SCRIPT_EFFECT                                            SpellEffectType = 77
	E_ATTACK                                                   SpellEffectType = 78
	E_SANCTUARY                                                SpellEffectType = 79
	E_ADD_COMBO_POINTS                                         SpellEffectType = 80
	E_PUSH_ABILITY_TO_ACTION_BAR                               SpellEffectType = 81
	E_BIND_SIGHT                                               SpellEffectType = 82
	E_DUEL                                                     SpellEffectType = 83
	E_STUCK                                                    SpellEffectType = 84
	E_SUMMON_PLAYER                                            SpellEffectType = 85
	E_ACTIVATE_OBJECT                                          SpellEffectType = 86
	E_GAMEOBJECT_DAMAGE                                        SpellEffectType = 87
	E_GAMEOBJECT_REPAIR                                        SpellEffectType = 88
	E_GAMEOBJECT_SET_DESTRUCTION_STATE                         SpellEffectType = 89
	E_KILL_CREDIT                                              SpellEffectType = 90
	E_THREAT_ALL                                               SpellEffectType = 91
	E_ENCHANT_HELD_ITEM                                        SpellEffectType = 92
	E_FORCE_DESELECT                                           SpellEffectType = 93
	E_SELF_RESURRECT                                           SpellEffectType = 94
	E_SKINNING                                                 SpellEffectType = 95
	E_CHARGE                                                   SpellEffectType = 96
	E_CAST_BUTTON                                              SpellEffectType = 97
	E_KNOCK_BACK                                               SpellEffectType = 98
	E_DISENCHANT                                               SpellEffectType = 99
	E_INEBRIATE                                                SpellEffectType = 100
	E_FEED_PET                                                 SpellEffectType = 101
	E_DISMISS_PET                                              SpellEffectType = 102
	E_REPUTATION                                               SpellEffectType = 103
	E_SUMMON_OBJECT_SLOT1                                      SpellEffectType = 104
	E_SURVEY                                                   SpellEffectType = 105
	E_CHANGE_RAID_MARKER                                       SpellEffectType = 106
	E_SHOW_CORPSE_LOOT                                         SpellEffectType = 107
	E_DISPEL_MECHANIC                                          SpellEffectType = 108
	E_RESURRECT_PET                                            SpellEffectType = 109
	E_DESTROY_ALL_TOTEMS                                       SpellEffectType = 110
	E_DURABILITY_DAMAGE                                        SpellEffectType = 111
	E_ATTACK_ME                                                SpellEffectType = 114
	E_DURABILITY_DAMAGE_PCT                                    SpellEffectType = 115
	E_SKIN_PLAYER_CORPSE                                       SpellEffectType = 116
	E_SPIRIT_HEAL                                              SpellEffectType = 117
	E_SKILL                                                    SpellEffectType = 118
	E_APPLY_AREA_AURA_PET                                      SpellEffectType = 119
	E_TELEPORT_GRAVEYARD                                       SpellEffectType = 120
	E_NORMALIZED_WEAPON_DMG                                    SpellEffectType = 121
	E_SEND_TAXI                                                SpellEffectType = 123
	E_PULL_TOWARDS                                             SpellEffectType = 124
	E_MODIFY_THREAT_PERCENT                                    SpellEffectType = 125
	E_STEAL_BENEFICIAL_BUFF                                    SpellEffectType = 126
	E_PROSPECTING                                              SpellEffectType = 127
	E_APPLY_AREA_AURA_FRIEND                                   SpellEffectType = 128
	E_APPLY_AREA_AURA_ENEMY                                    SpellEffectType = 129
	E_REDIRECT_THREAT                                          SpellEffectType = 130
	E_PLAY_SOUND                                               SpellEffectType = 131
	E_PLAY_MUSIC                                               SpellEffectType = 132
	E_UNLEARN_SPECIALIZATION                                   SpellEffectType = 133
	E_KILL_CREDIT2                                             SpellEffectType = 134
	E_CALL_PET                                                 SpellEffectType = 135
	E_HEAL_PCT                                                 SpellEffectType = 136
	E_ENERGIZE_PCT                                             SpellEffectType = 137
	E_LEAP_BACK                                                SpellEffectType = 138
	E_CLEAR_QUEST                                              SpellEffectType = 139
	E_FORCE_CAST                                               SpellEffectType = 140
	E_FORCE_CAST_WITH_VALUE                                    SpellEffectType = 141
	E_TRIGGER_SPELL_WITH_VALUE                                 SpellEffectType = 142
	E_APPLY_AREA_AURA_OWNER                                    SpellEffectType = 143
	E_KNOCK_BACK_DEST                                          SpellEffectType = 144
	E_PULL_TOWARDS_DEST                                        SpellEffectType = 145
	E_ACTIVATE_RUNE                                            SpellEffectType = 146
	E_QUEST_FAIL                                               SpellEffectType = 147
	E_TRIGGER_MISSILE_SPELL_WITH_VALUE                         SpellEffectType = 148
	E_CHARGE_DEST                                              SpellEffectType = 149
	E_QUEST_START                                              SpellEffectType = 150
	E_TRIGGER_SPELL_2                                          SpellEffectType = 151
	E_SUMMON_RAF_FRIEND                                        SpellEffectType = 152
	E_CREATE_TAMED_PET                                         SpellEffectType = 153
	E_DISCOVER_TAXI                                            SpellEffectType = 154
	E_TITAN_GRIP                                               SpellEffectType = 155
	E_ENCHANT_ITEM_PRISMATIC                                   SpellEffectType = 156
	E_CREATE_LOOT                                              SpellEffectType = 157
	E_MILLING                                                  SpellEffectType = 158
	E_ALLOW_RENAME_PET                                         SpellEffectType = 159
	E_FORCE_CAST_2                                             SpellEffectType = 160
	E_TALENT_SPEC_COUNT                                        SpellEffectType = 161
	E_TALENT_SPEC_SELECT                                       SpellEffectType = 162
	E_OBLITERATE_ITEM                                          SpellEffectType = 163
	E_REMOVE_AURA                                              SpellEffectType = 164
	E_DAMAGE_FROM_MAX_HEALTH_PCT                               SpellEffectType = 165
	E_GIVE_CURRENCY                                            SpellEffectType = 166
	E_UPDATE_PLAYER_PHASE                                      SpellEffectType = 167
	E_ALLOW_CONTROL_PET                                        SpellEffectType = 168
	E_DESTROY_ITEM                                             SpellEffectType = 169
	E_UPDATE_ZONE_AURAS_AND_PHASES                             SpellEffectType = 170
	E_SUMMON_PERSONAL_GAMEOBJECT                               SpellEffectType = 171
	E_RESURRECT_WITH_AURA                                      SpellEffectType = 172
	E_UNLOCK_GUILD_VAULT_TAB                                   SpellEffectType = 173
	E_APPLY_AURA_ON_PET                                        SpellEffectType = 174
	E_SANCTUARY_2                                              SpellEffectType = 176
	E_CREATE_AREATRIGGER                                       SpellEffectType = 179
	E_UPDATE_AREATRIGGER                                       SpellEffectType = 180
	E_REMOVE_TALENT                                            SpellEffectType = 181
	E_DESPAWN_AREATRIGGER                                      SpellEffectType = 182
	E_REPUTATION_2                                             SpellEffectType = 184
	E_RANDOMIZE_ARCHAEOLOGY_DIGSITES                           SpellEffectType = 187
	E_LOOT                                                     SpellEffectType = 189
	E_TELEPORT_TO_DIGSITE                                      SpellEffectType = 191
	E_UNCAGE_BATTLEPET                                         SpellEffectType = 192
	E_START_PET_BATTLE                                         SpellEffectType = 193
	E_PLAY_SCENE                                               SpellEffectType = 198
	E_HEAL_BATTLEPET_PCT                                       SpellEffectType = 200
	E_ENABLE_BATTLE_PETS                                       SpellEffectType = 201
	E_APPLY_AURA_ON_UNKNOWN                                    SpellEffectType = 202 // originally "APPLY_AURA_ON_?"
	E_CHANGE_BATTLEPET_QUALITY                                 SpellEffectType = 204
	E_LAUNCH_QUEST_CHOICE                                      SpellEffectType = 205
	E_ALTER_ITEM                                               SpellEffectType = 206
	E_LAUNCH_QUEST_TASK                                        SpellEffectType = 207
	E_LEARN_GARRISON_BUILDING                                  SpellEffectType = 210
	E_LEARN_GARRISON_SPECIALIZATION                            SpellEffectType = 211
	E_CREATE_GARRISON                                          SpellEffectType = 214
	E_UPGRADE_CHARACTER_SPELLS                                 SpellEffectType = 215
	E_CREATE_SHIPMENT                                          SpellEffectType = 216
	E_UPGRADE_GARRISON                                         SpellEffectType = 217
	E_CREATE_CONVERSATION                                      SpellEffectType = 219
	E_ADD_GARRISON_FOLLOWER                                    SpellEffectType = 220
	E_CREATE_HEIRLOOM_ITEM                                     SpellEffectType = 222
	E_CHANGE_ITEM_BONUSES                                      SpellEffectType = 223
	E_ACTIVATE_GARRISON_BUILDING                               SpellEffectType = 224
	E_GRANT_BATTLEPET_LEVEL                                    SpellEffectType = 225
	E_TELEPORT_TO_LFG_DUNGEON                                  SpellEffectType = 227
	E_SET_FOLLOWER_QUALITY                                     SpellEffectType = 229
	E_INCREASE_FOLLOWER_ITEM_LEVEL                             SpellEffectType = 230
	E_INCREASE_FOLLOWER_EXPERIENCE                             SpellEffectType = 231
	E_REMOVE_PHASE                                             SpellEffectType = 232
	E_RANDOMIZE_FOLLOWER_ABILITIES                             SpellEffectType = 233
	E_GIVE_EXPERIENCE                                          SpellEffectType = 236
	E_GIVE_RESTED_EXPERIENCE_BONUS                             SpellEffectType = 237
	E_INCREASE_SKILL                                           SpellEffectType = 238
	E_END_GARRISON_BUILDING_CONSTRUCTION                       SpellEffectType = 239
	E_GIVE_ARTIFACT_POWER                                      SpellEffectType = 240
	E_GIVE_ARTIFACT_POWER_NO_BONUS                             SpellEffectType = 242
	E_APPLY_ENCHANT_ILLUSION                                   SpellEffectType = 243
	E_LEARN_FOLLOWER_ABILITY                                   SpellEffectType = 244
	E_UPGRADE_HEIRLOOM                                         SpellEffectType = 245
	E_FINISH_GARRISON_MISSION                                  SpellEffectType = 246
	E_ADD_GARRISON_MISSION                                     SpellEffectType = 247
	E_FINISH_SHIPMENT                                          SpellEffectType = 248
	E_FORCE_EQUIP_ITEM                                         SpellEffectType = 249
	E_TAKE_SCREENSHOT                                          SpellEffectType = 250
	E_SET_GARRISON_CACHE_SIZE                                  SpellEffectType = 251
	E_TELEPORT_UNITS                                           SpellEffectType = 252
	E_GIVE_HONOR                                               SpellEffectType = 253
	E_LEARN_TRANSMOG_SET                                       SpellEffectType = 255
	E_MODIFY_KEYSTONE                                          SpellEffectType = 258
	E_RESPEC_AZERITE_EMPOWERED_ITEM                            SpellEffectType = 259
	E_SUMMON_STABLED_PET                                       SpellEffectType = 260
	E_SCRAP_ITEM                                               SpellEffectType = 261
	E_REPAIR_ITEM                                              SpellEffectType = 263
	E_REMOVE_GEM                                               SpellEffectType = 264
	E_LEARN_AZERITE_ESSENCE_POWER                              SpellEffectType = 265
	E_APPLY_MOUNT_EQUIPMENT                                    SpellEffectType = 268
	E_UPGRADE_ITEM                                             SpellEffectType = 269
	E_APPLY_AREA_AURA_PARTY_NONRANDOM                          SpellEffectType = 271
	E_SET_COVENANT                                             SpellEffectType = 272
	E_CRAFT_RUNEFORGE_LEGENDARY                                SpellEffectType = 273
	E_LEARN_TRANSMOG_ILLUSION                                  SpellEffectType = 276
	E_SET_CHROMIE_TIME                                         SpellEffectType = 277
	E_LEARN_GARR_TALENT                                        SpellEffectType = 279
	E_LEARN_SOULBIND_CONDUIT                                   SpellEffectType = 281
	E_CONVERT_ITEMS_TO_CURRENCY                                SpellEffectType = 282
	E_COMPLETE_CAMPAIGN                                        SpellEffectType = 283
	E_SEND_CHAT_MESSAGE                                        SpellEffectType = 284
	E_MODIFY_KEYSTONE_2                                        SpellEffectType = 285
	E_GRANT_BATTLEPET_EXPERIENCE                               SpellEffectType = 286
	E_SET_GARRISON_FOLLOWER_LEVEL                              SpellEffectType = 287
	E_CRAFT_ITEM                                               SpellEffectType = 288
	E_MODIFY_AURA_STACKS                                       SpellEffectType = 289
	E_MODIFY_COOLDOWN                                          SpellEffectType = 290
	E_MODIFY_COOLDOWNS                                         SpellEffectType = 291
	E_MODIFY_COOLDOWNS_BY_CATEGORY                             SpellEffectType = 292
	E_MODIFY_CHARGES                                           SpellEffectType = 293
	E_CRAFT_LOOT                                               SpellEffectType = 294
	E_SALVAGE_ITEM                                             SpellEffectType = 295
	E_CRAFT_SALVAGE_ITEM                                       SpellEffectType = 296
	E_RECRAFT_ITEM                                             SpellEffectType = 297
	E_CANCEL_ALL_PRIVATE_CONVERSATIONS                         SpellEffectType = 298
	E_CRAFT_ENCHANT                                            SpellEffectType = 301
	E_GATHERING                                                SpellEffectType = 302
	E_CREATE_TRAIT_TREE_CONFIG                                 SpellEffectType = 303
	E_CHANGE_ACTIVE_COMBAT_TRAIT_CONFIG                        SpellEffectType = 304
	E_UPDATE_INTERACTIONS                                      SpellEffectType = 306
	E_CANCEL_PRELOAD_WORLD                                     SpellEffectType = 308
	E_PRELOAD_WORLD                                            SpellEffectType = 309
	E_ENSURE_WORLD_LOADED                                      SpellEffectType = 310
	E_CHANGE_ITEM_BONUSES_2                                    SpellEffectType = 311
	E_ADD_SOCKET_BONUS                                         SpellEffectType = 312
	E_LEARN_TRANSMOG_APPEARANCE_FROM_ITEM_MOD_APPEARANCE_GROUP SpellEffectType = 313
	E_KILL_CREDIT_LABEL_1                                      SpellEffectType = 314
	E_KILL_CREDIT_LABEL_2                                      SpellEffectType = 315
	E_UI_ACTION                                                SpellEffectType = 316
	E_LEARN_WARBAND_SCENE                                      SpellEffectType = 317
)

type ItemSubClass int

// Define each item subclass as a bit flag (only those with a name).
const (
	OneHandedAxes    ItemSubClass = 1 << 0  // 1    from "One-Handed Axes" (SubClassID 0)
	TwoHandedAxes    ItemSubClass = 1 << 1  // 2    from "Two-Handed Axes" (SubClassID 1)
	Bows             ItemSubClass = 1 << 2  // 4    from "Bows" (SubClassID 2)
	Guns             ItemSubClass = 1 << 3  // 8    from "Guns" (SubClassID 3)
	OneHandedMaces   ItemSubClass = 1 << 4  // 16   from "One-Handed Maces" (SubClassID 4)
	TwoHandedMaces   ItemSubClass = 1 << 5  // 32   from "Two-Handed Maces" (SubClassID 5)
	Polearms         ItemSubClass = 1 << 6  // 64   from "Polearms" (SubClassID 6)
	OneHandedSwords  ItemSubClass = 1 << 7  // 128  from "One-Handed Swords" (SubClassID 7)
	TwoHandedSwords  ItemSubClass = 1 << 8  // 256  from "Two-Handed Swords" (SubClassID 8)
	Staves           ItemSubClass = 1 << 10 // 1024 from "Staves" (SubClassID 10)
	OneHandedExotics ItemSubClass = 1 << 11 // 2048 from "One-Handed Exotics" (SubClassID 11)
	TwoHandedExotics ItemSubClass = 1 << 12 // 4096 from "Two-Handed Exotics" (SubClassID 12)
	FistWeapons      ItemSubClass = 1 << 13 // 8192 from "Fist Weapons" (SubClassID 13)
	Daggers          ItemSubClass = 1 << 15 // 32768 from "Daggers" (SubClassID 15)
)

type ItemQuality int

const (
	JUNK      ItemQuality = 0
	COMMON    ItemQuality = 1
	UNCOMMON  ItemQuality = 2
	RARE      ItemQuality = 3
	EPIC      ItemQuality = 4
	LEGENDARY ItemQuality = 5
	ARTIFACT  ItemQuality = 6
	HEIRLOOM  ItemQuality = 7
)

func (raw ItemQuality) ToProto() proto.ItemQuality {
	switch raw {
	case JUNK:
		return proto.ItemQuality_ItemQualityJunk
	case COMMON:
		return proto.ItemQuality_ItemQualityCommon
	case UNCOMMON:
		return proto.ItemQuality_ItemQualityUncommon
	case RARE:
		return proto.ItemQuality_ItemQualityRare
	case EPIC:
		return proto.ItemQuality_ItemQualityEpic
	case LEGENDARY:
		return proto.ItemQuality_ItemQualityLegendary
	case ARTIFACT:
		return proto.ItemQuality_ItemQualityArtifact
	case HEIRLOOM:
		return proto.ItemQuality_ItemQualityHeirloom
	}
	return proto.ItemQuality_ItemQualityUncommon
}

const (
	ITEM_CLASS_CONSUMABLE = iota
	ITEM_CLASS_CONTAINER
	ITEM_CLASS_WEAPON
	ITEM_CLASS_GEM
	ITEM_CLASS_ARMOR
	ITEM_CLASS_REAGENT
	ITEM_CLASS_PROJECTILE
	ITEM_CLASS_TRADE_GOODS
	ITEM_CLASS_GENERIC
	ITEM_CLASS_RECIPE
	ITEM_CLASS_MONEY
	ITEM_CLASS_QUIVER
	ITEM_CLASS_QUEST
	ITEM_CLASS_KEY
	ITEM_CLASS_PERMANENT
	ITEM_CLASS_MISC
	ITEM_CLASS_GLYPH
)

const (
	ITEM_SUBCLASS_WEAPON_AXE = iota
	ITEM_SUBCLASS_WEAPON_AXE2
	ITEM_SUBCLASS_WEAPON_BOW
	ITEM_SUBCLASS_WEAPON_GUN
	ITEM_SUBCLASS_WEAPON_MACE
	ITEM_SUBCLASS_WEAPON_MACE2
	ITEM_SUBCLASS_WEAPON_POLEARM
	ITEM_SUBCLASS_WEAPON_SWORD
	ITEM_SUBCLASS_WEAPON_SWORD2
	ITEM_SUBCLASS_WEAPON_WARGLAIVE
	ITEM_SUBCLASS_WEAPON_STAFF
	ITEM_SUBCLASS_WEAPON_EXOTIC
	ITEM_SUBCLASS_WEAPON_EXOTIC2
	ITEM_SUBCLASS_WEAPON_FIST
	ITEM_SUBCLASS_WEAPON_MISC
	ITEM_SUBCLASS_WEAPON_DAGGER
	ITEM_SUBCLASS_WEAPON_THROWN
	ITEM_SUBCLASS_WEAPON_SPEAR
	ITEM_SUBCLASS_WEAPON_CROSSBOW
	ITEM_SUBCLASS_WEAPON_WAND
	ITEM_SUBCLASS_WEAPON_FISHING_POLE
)

const ITEM_SUBCLASS_WEAPON_INVALID = 31

const (
	ITEM_SUBCLASS_ARMOR_MISC = iota
	ITEM_SUBCLASS_ARMOR_CLOTH
	ITEM_SUBCLASS_ARMOR_LEATHER
	ITEM_SUBCLASS_ARMOR_MAIL
	ITEM_SUBCLASS_ARMOR_PLATE
	ITEM_SUBCLASS_ARMOR_COSMETIC
	ITEM_SUBCLASS_ARMOR_SHIELD
	ITEM_SUBCLASS_ARMOR_LIBRAM
	ITEM_SUBCLASS_ARMOR_IDOL
	ITEM_SUBCLASS_ARMOR_TOTEM
	ITEM_SUBCLASS_ARMOR_SIGIL
	ITEM_SUBCLASS_ARMOR_RELIC
)

const (
	ITEM_SUBCLASS_CONSUMABLE = iota
	ITEM_SUBCLASS_POTION
	ITEM_SUBCLASS_ELIXIR
	ITEM_SUBCLASS_FLASK
	ITEM_SUBCLASS_SCROLL
	ITEM_SUBCLASS_FOOD
	ITEM_SUBCLASS_ITEM_ENHANCEMENT
	ITEM_SUBCLASS_BANDAGE
	ITEM_SUBCLASS_CONSUMABLE_OTHER
)

const (
	INVTYPE_NON_EQUIP = iota
	INVTYPE_HEAD
	INVTYPE_NECK
	INVTYPE_SHOULDERS
	INVTYPE_BODY
	INVTYPE_CHEST
	INVTYPE_WAIST
	INVTYPE_LEGS
	INVTYPE_FEET
	INVTYPE_WRISTS
	INVTYPE_HANDS
	INVTYPE_FINGER
	INVTYPE_TRINKET
	INVTYPE_WEAPON
	INVTYPE_SHIELD
	INVTYPE_RANGED
	INVTYPE_CLOAK
	INVTYPE_2HWEAPON
	INVTYPE_BAG
	INVTYPE_TABARD
	INVTYPE_ROBE
	INVTYPE_WEAPONMAINHAND
	INVTYPE_WEAPONOFFHAND
	INVTYPE_HOLDABLE
	INVTYPE_AMMO
	INVTYPE_THROWN
	INVTYPE_RANGEDRIGHT
	INVTYPE_QUIVER
	INVTYPE_RELIC
	INVTYPE_MAX
)

// EffectAuraType defines the custom type for aura effects.
type EffectAuraType int

// Enum constants defined using the A_ naming convention.
// (Keys that were commented out in the JavaScript object are left commented here.)
const (
	A_NONE                       EffectAuraType = 0
	A_BIND_SIGHT                 EffectAuraType = 1
	A_MOD_POSSESS                EffectAuraType = 2
	A_PERIODIC_DAMAGE            EffectAuraType = 3
	A_DUMMY                      EffectAuraType = 4
	A_MOD_CONFUSE                EffectAuraType = 5
	A_MOD_CHARM                  EffectAuraType = 6
	A_MOD_FEAR                   EffectAuraType = 7
	A_PERIODIC_HEAL              EffectAuraType = 8
	A_MOD_ATTACKSPEED            EffectAuraType = 9
	A_MOD_THREAT                 EffectAuraType = 10
	A_MOD_TAUNT                  EffectAuraType = 11
	A_MOD_STUN                   EffectAuraType = 12
	A_MOD_DAMAGE_DONE            EffectAuraType = 13
	A_MOD_DAMAGE_TAKEN           EffectAuraType = 14
	A_DAMAGE_SHIELD              EffectAuraType = 15
	A_MOD_STEALTH                EffectAuraType = 16
	A_MOD_STEALTH_DETECT         EffectAuraType = 17
	A_MOD_INVISIBILITY           EffectAuraType = 18
	A_MOD_INVISIBILITY_DETECT    EffectAuraType = 19
	A_OBS_MOD_HEALTH             EffectAuraType = 20
	A_OBS_MOD_POWER              EffectAuraType = 21
	A_MOD_RESISTANCE             EffectAuraType = 22
	A_PERIODIC_TRIGGER_SPELL     EffectAuraType = 23
	A_PERIODIC_ENERGIZE          EffectAuraType = 24
	A_MOD_PACIFY                 EffectAuraType = 25
	A_MOD_ROOT                   EffectAuraType = 26
	A_MOD_SILENCE                EffectAuraType = 27
	A_REFLECT_SPELLS             EffectAuraType = 28
	A_MOD_STAT                   EffectAuraType = 29
	A_MOD_SKILL                  EffectAuraType = 30
	A_MOD_INCREASE_SPEED         EffectAuraType = 31
	A_MOD_INCREASE_MOUNTED_SPEED EffectAuraType = 32
	A_MOD_DECREASE_SPEED         EffectAuraType = 33
	A_MOD_INCREASE_HEALTH        EffectAuraType = 34
	A_MOD_INCREASE_ENERGY        EffectAuraType = 35
	A_MOD_SHAPESHIFT             EffectAuraType = 36
	A_EFFECT_IMMUNITY            EffectAuraType = 37
	A_STATE_IMMUNITY             EffectAuraType = 38
	A_SCHOOL_IMMUNITY            EffectAuraType = 39
	A_DAMAGE_IMMUNITY            EffectAuraType = 40
	A_DISPEL_IMMUNITY            EffectAuraType = 41
	A_PROC_TRIGGER_SPELL         EffectAuraType = 42
	A_PROC_TRIGGER_DAMAGE        EffectAuraType = 43
	A_TRACK_CREATURES            EffectAuraType = 44
	A_TRACK_RESOURCES            EffectAuraType = 45
	// 46 is commented out in the JS mapping
	A_MOD_PARRY_PERCENT                  EffectAuraType = 47
	A_PERIODIC_TRIGGER_SPELL_FROM_CLIENT EffectAuraType = 48
	A_MOD_DODGE_PERCENT                  EffectAuraType = 49
	A_MOD_CRITICAL_HEALING_AMOUNT        EffectAuraType = 50
	A_MOD_BLOCK_PERCENT                  EffectAuraType = 51
	A_MOD_WEAPON_CRIT_PERCENT            EffectAuraType = 52
	A_PERIODIC_LEECH                     EffectAuraType = 53
	A_MOD_HIT_CHANCE                     EffectAuraType = 54
	A_MOD_SPELL_HIT_CHANCE               EffectAuraType = 55
	A_TRANSFORM                          EffectAuraType = 56
	A_MOD_SPELL_CRIT_CHANCE              EffectAuraType = 57
	A_MOD_INCREASE_SWIM_SPEED            EffectAuraType = 58
	A_MOD_DAMAGE_DONE_CREATURE           EffectAuraType = 59
	A_MOD_PACIFY_SILENCE                 EffectAuraType = 60
	A_MOD_SCALE                          EffectAuraType = 61
	A_PERIODIC_HEALTH_FUNNEL             EffectAuraType = 62
	A_MOD_ADDITIONAL_POWER_COST          EffectAuraType = 63
	A_PERIODIC_MANA_LEECH                EffectAuraType = 64
	A_MOD_CASTING_SPEED_NOT_STACK        EffectAuraType = 65
	A_FEIGN_DEATH                        EffectAuraType = 66
	A_MOD_DISARM                         EffectAuraType = 67
	A_MOD_STALKED                        EffectAuraType = 68
	A_SCHOOL_ABSORB                      EffectAuraType = 69
	A_PERIODIC_WEAPON_PERCENT_DAMAGE     EffectAuraType = 70
	A_STORE_TELEPORT_RETURN_POINT        EffectAuraType = 71
	A_MOD_POWER_COST_SCHOOL_PCT          EffectAuraType = 72
	A_MOD_POWER_COST_SCHOOL              EffectAuraType = 73
	A_REFLECT_SPELLS_SCHOOL              EffectAuraType = 74
	A_MOD_LANGUAGE                       EffectAuraType = 75
	A_FAR_SIGHT                          EffectAuraType = 76
	A_MECHANIC_IMMUNITY                  EffectAuraType = 77
	A_MOUNTED                            EffectAuraType = 78
	A_MOD_DAMAGE_PERCENT_DONE            EffectAuraType = 79
	A_MOD_PERCENT_STAT                   EffectAuraType = 80
	A_SPLIT_DAMAGE_PCT                   EffectAuraType = 81
	A_WATER_BREATHING                    EffectAuraType = 82
	A_MOD_BASE_RESISTANCE                EffectAuraType = 83
	A_MOD_REGEN                          EffectAuraType = 84
	A_MOD_POWER_REGEN                    EffectAuraType = 85
	A_CHANNEL_DEATH_ITEM                 EffectAuraType = 86
	A_MOD_DAMAGE_PERCENT_TAKEN           EffectAuraType = 87
	A_MOD_HEALTH_REGEN_PERCENT           EffectAuraType = 88
	A_PERIODIC_DAMAGE_PERCENT            EffectAuraType = 89
	// 90 is commented out
	A_MOD_DETECT_RANGE                   EffectAuraType = 91
	A_PREVENTS_FLEEING                   EffectAuraType = 92
	A_MOD_UNATTACKABLE                   EffectAuraType = 93
	A_INTERRUPT_REGEN                    EffectAuraType = 94
	A_GHOST                              EffectAuraType = 95
	A_SPELL_MAGNET                       EffectAuraType = 96
	A_MANA_SHIELD                        EffectAuraType = 97
	A_MOD_SKILL_TALENT                   EffectAuraType = 98
	A_MOD_ATTACK_POWER                   EffectAuraType = 99
	A_AURAS_VISIBLE                      EffectAuraType = 100
	A_MOD_RESISTANCE_PCT                 EffectAuraType = 101
	A_MOD_MELEE_ATTACK_POWER_VERSUS      EffectAuraType = 102
	A_MOD_TOTAL_THREAT                   EffectAuraType = 103
	A_WATER_WALK                         EffectAuraType = 104
	A_FEATHER_FALL                       EffectAuraType = 105
	A_HOVER                              EffectAuraType = 106
	A_ADD_FLAT_MODIFIER                  EffectAuraType = 107
	A_ADD_PCT_MODIFIER                   EffectAuraType = 108
	A_ADD_TARGET_TRIGGER                 EffectAuraType = 109
	A_MOD_POWER_REGEN_PERCENT            EffectAuraType = 110
	A_INTERCEPT_MELEE_RANGED_ATTACKS     EffectAuraType = 111
	A_OVERRIDE_CLASS_SCRIPTS             EffectAuraType = 112
	A_MOD_RANGED_DAMAGE_TAKEN            EffectAuraType = 113
	A_MOD_RANGED_DAMAGE_TAKEN_PCT        EffectAuraType = 114
	A_MOD_HEALING                        EffectAuraType = 115
	A_MOD_REGEN_DURING_COMBAT            EffectAuraType = 116
	A_MOD_MECHANIC_RESISTANCE            EffectAuraType = 117
	A_MOD_HEALING_PCT                    EffectAuraType = 118
	A_PVP_TALENTS                        EffectAuraType = 119
	A_UNTRACKABLE                        EffectAuraType = 120
	A_EMPATHY                            EffectAuraType = 121
	A_MOD_OFFHAND_DAMAGE_PCT             EffectAuraType = 122
	A_MOD_TARGET_RESISTANCE              EffectAuraType = 123
	A_MOD_RANGED_ATTACK_POWER            EffectAuraType = 124
	A_MOD_MELEE_DAMAGE_TAKEN             EffectAuraType = 125
	A_MOD_MELEE_DAMAGE_TAKEN_PCT         EffectAuraType = 126
	A_RANGED_ATTACK_POWER_ATTACKER_BONUS EffectAuraType = 127
	A_MOD_FIXATE                         EffectAuraType = 128
	A_MOD_SPEED_ALWAYS                   EffectAuraType = 129
	A_MOD_MOUNTED_SPEED_ALWAYS           EffectAuraType = 130
	A_MOD_RANGED_ATTACK_POWER_VERSUS     EffectAuraType = 131
	A_MOD_INCREASE_ENERGY_PERCENT        EffectAuraType = 132
	A_MOD_INCREASE_HEALTH_PERCENT        EffectAuraType = 133
	A_MOD_MANA_REGEN_INTERRUPT           EffectAuraType = 134
	A_MOD_HEALING_DONE                   EffectAuraType = 135
	A_MOD_HEALING_DONE_PERCENT           EffectAuraType = 136
	A_MOD_TOTAL_STAT_PERCENTAGE          EffectAuraType = 137
	A_MOD_MELEE_HASTE                    EffectAuraType = 138
	A_FORCE_REACTION                     EffectAuraType = 139
	A_MOD_RANGED_HASTE                   EffectAuraType = 140
	// 141 is commented out
	A_MOD_BASE_RESISTANCE_PCT          EffectAuraType = 142
	A_MOD_RECOVERY_RATE_BY_SPELL_LABEL EffectAuraType = 143
	A_SAFE_FALL                        EffectAuraType = 144
	A_MOD_INCREASE_HEALTH_PERCENT2     EffectAuraType = 145
	A_ALLOW_TAME_PET_TYPE              EffectAuraType = 146
	A_MECHANIC_IMMUNITY_MASK           EffectAuraType = 147
	A_MOD_CHARGE_RECOVERY_RATE         EffectAuraType = 148
	A_REDUCE_PUSHBACK                  EffectAuraType = 149
	A_MOD_SHIELD_BLOCKVALUE_PCT        EffectAuraType = 150
	A_TRACK_STEALTHED                  EffectAuraType = 151
	A_MOD_DETECTED_RANGE               EffectAuraType = 152
	A_MOD_AUTOATTACK_RANGE             EffectAuraType = 153
	A_MOD_STEALTH_LEVEL                EffectAuraType = 154
	A_MOD_WATER_BREATHING              EffectAuraType = 155
	A_MOD_REPUTATION_GAIN              EffectAuraType = 156
	A_PET_DAMAGE_MULTI                 EffectAuraType = 157
	A_ALLOW_TALENT_SWAPPING            EffectAuraType = 158
	A_NO_PVP_CREDIT                    EffectAuraType = 159
	// 160 is commented out
	A_MOD_HEALTH_REGEN_IN_COMBAT        EffectAuraType = 161
	A_POWER_BURN                        EffectAuraType = 162
	A_MOD_CRIT_DAMAGE_BONUS             EffectAuraType = 163
	A_FORCE_BREATH_BAR                  EffectAuraType = 164
	A_MELEE_ATTACK_POWER_ATTACKER_BONUS EffectAuraType = 165
	A_MOD_ATTACK_POWER_PCT              EffectAuraType = 166
	A_MOD_RANGED_ATTACK_POWER_PCT       EffectAuraType = 167
	A_MOD_DAMAGE_DONE_VERSUS            EffectAuraType = 168
	A_SET_FFA_PVP                       EffectAuraType = 169
	A_DETECT_AMORE                      EffectAuraType = 170
	A_MOD_SPEED_NOT_STACK               EffectAuraType = 171
	A_MOD_MOUNTED_SPEED_NOT_STACK       EffectAuraType = 172
	// 173 is commented out
	A_MOD_SPELL_DAMAGE_OF_STAT_PERCENT            EffectAuraType = 174
	A_MOD_SPELL_HEALING_OF_STAT_PERCENT           EffectAuraType = 175
	A_SPIRIT_OF_REDEMPTION                        EffectAuraType = 176
	A_AOE_CHARM                                   EffectAuraType = 177
	A_MOD_MAX_POWER_PCT                           EffectAuraType = 178
	A_MOD_POWER_DISPLAY                           EffectAuraType = 179
	A_MOD_FLAT_SPELL_DAMAGE_VERSUS                EffectAuraType = 180
	A_MOD_SPELL_CURRENCY_REAGENTS_COUNT_PCT       EffectAuraType = 181
	A_SUPPRESS_ITEM_PASSIVE_EFFECT_BY_SPELL_LABEL EffectAuraType = 182
	A_MOD_CRIT_CHANCE_VERSUS_TARGET_HEALTH        EffectAuraType = 183
	A_MOD_ATTACKER_MELEE_HIT_CHANCE               EffectAuraType = 184
	A_MOD_ATTACKER_RANGED_HIT_CHANCE              EffectAuraType = 185
	A_MOD_ATTACKER_SPELL_HIT_CHANCE               EffectAuraType = 186
	A_MOD_ATTACKER_MELEE_CRIT_CHANCE              EffectAuraType = 187
	A_MOD_UI_HEALING_RANGE                        EffectAuraType = 188
	A_MOD_RATING                                  EffectAuraType = 189
	A_MOD_FACTION_REPUTATION_GAIN                 EffectAuraType = 190
	A_USE_NORMAL_MOVEMENT_SPEED                   EffectAuraType = 191
	A_MOD_MELEE_RANGED_HASTE                      EffectAuraType = 192
	A_MELEE_SLOW                                  EffectAuraType = 193
	A_MOD_TARGET_ABSORB_SCHOOL                    EffectAuraType = 194
	A_LEARN_SPELL                                 EffectAuraType = 195
	A_MOD_COOLDOWN                                EffectAuraType = 196
	A_MOD_ATTACKER_SPELL_AND_WEAPON_CRIT_CHANCE   EffectAuraType = 197
	A_MOD_COMBAT_RATING_FROM_COMBAT_RATING        EffectAuraType = 198
	// 199 is commented out
	A_MOD_XP_PCT                        EffectAuraType = 200
	A_FLY                               EffectAuraType = 201
	A_IGNORE_COMBAT_RESULT              EffectAuraType = 202
	A_PREVENT_INTERRUPT                 EffectAuraType = 203
	A_PREVENT_CORPSE_RELEASE            EffectAuraType = 204
	A_MOD_CHARGE_COOLDOWN               EffectAuraType = 205
	A_MOD_INCREASE_VEHICLE_FLIGHT_SPEED EffectAuraType = 206
	A_MOD_INCREASE_MOUNTED_FLIGHT_SPEED EffectAuraType = 207
	A_MOD_INCREASE_FLIGHT_SPEED         EffectAuraType = 208
	A_MOD_MOUNTED_FLIGHT_SPEED_ALWAYS   EffectAuraType = 209
	A_MOD_VEHICLE_SPEED_ALWAYS          EffectAuraType = 210
	A_MOD_FLIGHT_SPEED_NOT_STACK        EffectAuraType = 211
	A_MOD_HONOR_GAIN_PCT                EffectAuraType = 212
	A_MOD_RAGE_FROM_DAMAGE_DEALT        EffectAuraType = 213
	// 214 is commented out
	A_ARENA_PREPARATION                               EffectAuraType = 215
	A_HASTE_SPELLS                                    EffectAuraType = 216
	A_MOD_MELEE_HASTE_2                               EffectAuraType = 217
	A_ADD_PCT_MODIFIER_BY_SPELL_LABEL                 EffectAuraType = 218
	A_ADD_FLAT_MODIFIER_BY_SPELL_LABEL                EffectAuraType = 219
	A_MOD_ABILITY_SCHOOL_MASK                         EffectAuraType = 220
	A_MOD_DETAUNT                                     EffectAuraType = 221
	A_REMOVE_TRANSMOG_COST                            EffectAuraType = 222
	A_REMOVE_BARBER_SHOP_COST                         EffectAuraType = 223
	A_LEARN_TALENT                                    EffectAuraType = 224
	A_MOD_VISIBILITY_RANGE                            EffectAuraType = 225
	A_PERIODIC_DUMMY                                  EffectAuraType = 226
	A_PERIODIC_TRIGGER_SPELL_WITH_VALUE               EffectAuraType = 227
	A_DETECT_STEALTH                                  EffectAuraType = 228
	A_MOD_AOE_DAMAGE_AVOIDANCE                        EffectAuraType = 229
	A_MOD_MAX_HEALTH                                  EffectAuraType = 230
	A_PROC_TRIGGER_SPELL_WITH_VALUE                   EffectAuraType = 231
	A_MECHANIC_DURATION_MOD                           EffectAuraType = 232
	A_CHANGE_MODEL_FOR_ALL_HUMANOIDS                  EffectAuraType = 233
	A_MECHANIC_DURATION_MOD_NOT_STACK                 EffectAuraType = 234
	A_MOD_HOVER_NO_HEIGHT_OFFSET                      EffectAuraType = 235
	A_CONTROL_VEHICLE                                 EffectAuraType = 236
	A_237                                             EffectAuraType = 237
	A_238                                             EffectAuraType = 238
	A_MOD_SCALE_2                                     EffectAuraType = 239
	A_MOD_EXPERTISE                                   EffectAuraType = 240
	A_FORCE_MOVE_FORWARD                              EffectAuraType = 241
	A_MOD_SPELL_DAMAGE_FROM_HEALING                   EffectAuraType = 242
	A_MOD_FACTION                                     EffectAuraType = 243
	A_COMPREHEND_LANGUAGE                             EffectAuraType = 244
	A_MOD_AURA_DURATION_BY_DISPEL                     EffectAuraType = 245
	A_MOD_AURA_DURATION_BY_DISPEL_NOT_STACK           EffectAuraType = 246
	A_CLONE_CASTER                                    EffectAuraType = 247
	A_MOD_COMBAT_RESULT_CHANCE                        EffectAuraType = 248
	A_MOD_DAMAGE_PERCENT_DONE_BY_TARGET_AURA_MECHANIC EffectAuraType = 249
	A_MOD_INCREASE_HEALTH_2                           EffectAuraType = 250
	A_MOD_ENEMY_DODGE                                 EffectAuraType = 251
	A_MOD_SPEED_SLOW_ALL                              EffectAuraType = 252
	A_MOD_BLOCK_CRIT_CHANCE                           EffectAuraType = 253
	A_MOD_DISARM_OFFHAND                              EffectAuraType = 254
	A_MOD_MECHANIC_DAMAGE_TAKEN_PERCENT               EffectAuraType = 255
	A_NO_REAGENT_USE                                  EffectAuraType = 256
	A_MOD_TARGET_RESIST_BY_SPELL_CLASS                EffectAuraType = 257
	A_OVERRIDE_SUMMONED_OBJECT                        EffectAuraType = 258
	A_MOD_HOT_PCT                                     EffectAuraType = 259
	A_SCREEN_EFFECT                                   EffectAuraType = 260
	A_PHASE                                           EffectAuraType = 261
	A_ABILITY_IGNORE_AURASTATE                        EffectAuraType = 262
	A_DISABLE_CASTING_EXCEPT_ABILITIES                EffectAuraType = 263
	A_DISABLE_ATTACKING_EXCEPT_ABILITIES              EffectAuraType = 264
	// 265 is commented out
	A_SET_VIGNETTE                       EffectAuraType = 266
	A_MOD_IMMUNE_AURA_APPLY_SCHOOL       EffectAuraType = 267
	A_MOD_ARMOR_PCT_FROM_STAT            EffectAuraType = 268
	A_MOD_IGNORE_TARGET_RESIST           EffectAuraType = 269
	A_MOD_SCHOOL_MASK_DAMAGE_FROM_CASTER EffectAuraType = 270
	A_MOD_SPELL_DAMAGE_FROM_CASTER       EffectAuraType = 271
	A_MOD_BLOCK_VALUE_PCT                EffectAuraType = 272
	A_X_RAY                              EffectAuraType = 273
	A_MOD_BLOCK_VALUE_FLAT               EffectAuraType = 274
	A_MOD_IGNORE_SHAPESHIFT              EffectAuraType = 275
	A_MOD_DAMAGE_DONE_FOR_MECHANIC       EffectAuraType = 276
	// 277 is commented out
	A_MOD_DISARM_RANGED                    EffectAuraType = 278
	A_INITIALIZE_IMAGES                    EffectAuraType = 279
	A_SPELL_AURA_MOD_ARMOR_PENETRATION_PCT EffectAuraType = 280
	A_PROVIDE_SPELL_FOCUS                  EffectAuraType = 281
	A_MOD_BASE_HEALTH_PCT                  EffectAuraType = 282
	A_MOD_HEALING_RECEIVED                 EffectAuraType = 283
	A_LINKED                               EffectAuraType = 284
	A_LINKED_2                             EffectAuraType = 285
	A_MOD_RECOVERY_RATE                    EffectAuraType = 286
	A_DEFLECT_SPELLS                       EffectAuraType = 287
	A_IGNORE_HIT_DIRECTION                 EffectAuraType = 288
	A_PREVENT_DURABILITY_LOSS              EffectAuraType = 289
	A_MOD_CRIT_PCT                         EffectAuraType = 290
	A_MOD_XP_QUEST_PCT                     EffectAuraType = 291
	A_OPEN_STABLE                          EffectAuraType = 292
	A_OVERRIDE_SPELLS                      EffectAuraType = 293
	A_PREVENT_REGENERATE_POWER             EffectAuraType = 294
	A_MOD_PERIODIC_DAMAGE_TAKEN            EffectAuraType = 295
	A_SET_VEHICLE_ID                       EffectAuraType = 296
	A_MOD_ROOT_DISABLE_GRAVITY             EffectAuraType = 297
	A_MOD_STUN_DISABLE_GRAVITY             EffectAuraType = 298
	// 299 is commented out
	A_SHARE_DAMAGE_PCT   EffectAuraType = 300
	A_SCHOOL_HEAL_ABSORB EffectAuraType = 301
	// 302 is commented out
	A_MOD_DAMAGE_DONE_VERSUS_AURASTATE          EffectAuraType = 303
	A_MOD_DRUNK                                 EffectAuraType = 304
	A_MOD_MINIMUM_SPEED                         EffectAuraType = 305
	A_MOD_CRIT_CHANCE_FOR_CASTER                EffectAuraType = 306
	A_CAST_WHILE_WALKING_BY_SPELL_LABEL         EffectAuraType = 307
	A_MOD_CRIT_CHANCE_FOR_CASTER_WITH_ABILITIES EffectAuraType = 308
	A_MOD_RESILIENCE                            EffectAuraType = 309
	A_MOD_CREATURE_AOE_DAMAGE_AVOIDANCE         EffectAuraType = 310
	A_IGNORE_COMBAT                             EffectAuraType = 311
	A_ANIM_REPLACEMENT_SET                      EffectAuraType = 312
	// 313 is commented out
	A_PREVENT_RESURRECTION   EffectAuraType = 314
	A_UNDERWATER_WALKING     EffectAuraType = 315
	A_SCHOOL_ABSORB_OVERKILL EffectAuraType = 316
	A_MOD_SPELL_POWER_PCT    EffectAuraType = 317
	A_MASTERY                EffectAuraType = 318
	A_MOD_MELEE_HASTE_3      EffectAuraType = 319
	A_MOD_RANGED_HASTE_2     EffectAuraType = 320
	A_MOD_NO_ACTIONS         EffectAuraType = 321
	A_INTERFERE_TARGETTING   EffectAuraType = 322
	// 323 is commented out
	A_OVERRIDE_UNLOCKED_AZERITE_ESSENCE_RANK EffectAuraType = 324
	A_LEARN_PVP_TALENT                       EffectAuraType = 325
	A_PHASE_GROUP                            EffectAuraType = 326
	A_PHASE_ALWAYS_VISIBLE                   EffectAuraType = 327
	A_TRIGGER_SPELL_ON_POWER_PCT             EffectAuraType = 328
	A_MOD_POWER_GAIN_PCT                     EffectAuraType = 329
	A_CAST_WHILE_WALKING                     EffectAuraType = 330
	A_FORCE_WEATHER                          EffectAuraType = 331
	A_OVERRIDE_ACTIONBAR_SPELLS              EffectAuraType = 332
	A_OVERRIDE_ACTIONBAR_SPELLS_TRIGGERED    EffectAuraType = 333
	A_MOD_AUTOATTACK_CRIT_CHANCE             EffectAuraType = 334
	// 335 is commented out
	A_MOUNT_RESTRICTIONS                     EffectAuraType = 336
	A_MOD_VENDOR_ITEMS_PRICES                EffectAuraType = 337
	A_MOD_DURABILITY_LOSS                    EffectAuraType = 338
	A_MOD_CRIT_CHANCE_FOR_CASTER_PET         EffectAuraType = 339
	A_MOD_RESURRECTED_HEALTH_BY_GUILD_MEMBER EffectAuraType = 340
	A_MOD_SPELL_CATEGORY_COOLDOWN            EffectAuraType = 341
	A_MOD_MELEE_RANGED_HASTE_2               EffectAuraType = 342
	A_MOD_MELEE_DAMAGE_FROM_CASTER           EffectAuraType = 343
	A_MOD_AUTOATTACK_DAMAGE                  EffectAuraType = 344
	A_BYPASS_ARMOR_FOR_CASTER                EffectAuraType = 345
	A_ENABLE_ALT_POWER                       EffectAuraType = 346
	A_MOD_SPELL_COOLDOWN_BY_HASTE            EffectAuraType = 347
	A_MOD_MONEY_GAIN                         EffectAuraType = 348
	A_MOD_CURRENCY_GAIN                      EffectAuraType = 349
	A_350                                    EffectAuraType = 350
	// 351,352 are commented out
	A_MOD_CAMOUFLAGE EffectAuraType = 353
	// 354 is commented out
	A_MOD_CASTING_SPEED                    EffectAuraType = 355
	A_PROVIDE_TOTEM_CATEGORY               EffectAuraType = 356
	A_ENABLE_BOSS1_UNIT_FRAME              EffectAuraType = 357
	A_WORGEN_ALTERED_FORM                  EffectAuraType = 358
	A_MOD_HEALING_DONE_VERSUS_AURASTATE    EffectAuraType = 359
	A_PROC_TRIGGER_SPELL_COPY              EffectAuraType = 360
	A_OVERRIDE_AUTOATTACK_WITH_MELEE_SPELL EffectAuraType = 361
	// 362 is commented out
	A_MOD_NEXT_SPELL EffectAuraType = 363
	// 364 is commented out
	A_MAX_FAR_CLIP_PLANE                    EffectAuraType = 365
	A_OVERRIDE_SPELL_POWER_BY_AP_PCT        EffectAuraType = 366
	A_OVERRIDE_AUTOATTACK_WITH_RANGED_SPELL EffectAuraType = 367
	// 368 is commented out
	A_ENABLE_POWER_BAR_TIMER    EffectAuraType = 369
	A_SPELL_OVERRIDE_NAME_GROUP EffectAuraType = 370
	// 371 is commented out
	A_OVERRIDE_MOUNT_FROM_SET         EffectAuraType = 372
	A_MOD_SPEED_NO_CONTROL            EffectAuraType = 373
	A_MOD_FALL_DAMAGE_PCT             EffectAuraType = 374
	A_HIDE_MODEL_AND_EQUIPEMENT_SLOTS EffectAuraType = 375
	A_MOD_CURRENCY_GAIN_FROM_SOURCE   EffectAuraType = 376
	A_CAST_WHILE_WALKING_ALL          EffectAuraType = 377
	A_MOD_POSSESS_PET                 EffectAuraType = 378
	A_MOD_MANA_REGEN_PCT              EffectAuraType = 379
	// 380 is commented out
	A_MOD_DAMAGE_TAKEN_FROM_CASTER_PET EffectAuraType = 381
	A_MOD_PET_STAT_PCT                 EffectAuraType = 382
	A_IGNORE_SPELL_COOLDOWN            EffectAuraType = 383
	// 384,385,386,387 are commented out
	A_MOD_TAXI_FLIGHT_SPEED EffectAuraType = 388
	// 389,390,391,392 are commented out
	A_BLOCK_SPELLS_IN_FRONT                  EffectAuraType = 393
	A_SHOW_CONFIRMATION_PROMPT               EffectAuraType = 394
	A_AREA_TRIGGER                           EffectAuraType = 395
	A_TRIGGER_SPELL_ON_POWER_AMOUNT          EffectAuraType = 396
	A_BATTLEGROUND_PLAYER_POSITION_FACTIONAL EffectAuraType = 397
	A_BATTLEGROUND_PLAYER_POSITION           EffectAuraType = 398
	A_MOD_TIME_RATE                          EffectAuraType = 399
	A_MOD_SKILL_2                            EffectAuraType = 400
	// 401 is commented out
	A_MOD_OVERRIDE_POWER_DISPLAY      EffectAuraType = 402
	A_OVERRIDE_SPELL_VISUAL           EffectAuraType = 403
	A_OVERRIDE_ATTACK_POWER_BY_SP_PCT EffectAuraType = 404
	A_MOD_RATING_PCT                  EffectAuraType = 405
	A_KEYBOUND_OVERRIDE               EffectAuraType = 406
	A_MOD_FEAR_2                      EffectAuraType = 407
	A_SET_ACTION_BUTTON_SPELL_COUNT   EffectAuraType = 408
	A_CAN_TURN_WHILE_FALLING          EffectAuraType = 409
	// 410 is commented out
	A_MOD_MAX_CHARGES EffectAuraType = 411
	// 412 is commented out
	A_MOD_RANGED_ATTACK_DEFLECT_CHANCE        EffectAuraType = 413
	A_MOD_RANGED_ATTACK_BLOCK_CHANCE_IN_FRONT EffectAuraType = 414
	// 415 is commented out
	A_MOD_COOLDOWN_BY_HASTE_REGEN        EffectAuraType = 416
	A_MOD_GLOBAL_COOLDOWN_BY_HASTE_REGEN EffectAuraType = 417
	A_MOD_MAX_POWER                      EffectAuraType = 418
	A_MOD_BASE_MANA_PCT                  EffectAuraType = 419
	A_MOD_BATTLE_PET_XP_PCT              EffectAuraType = 420
	A_MOD_ABSORB_EFFECTS_DONE_PCT        EffectAuraType = 421
	A_MOD_ABSORB_EFFECTS_TAKEN_PCT       EffectAuraType = 422
	A_MOD_MANA_COST_PCT                  EffectAuraType = 423
	A_CASTER_IGNORE_LOS                  EffectAuraType = 424
	// 425,426 are commented out
	A_SCALE_PLAYER_LEVEL         EffectAuraType = 427
	A_LINKED_SUMMON              EffectAuraType = 428
	A_MOD_SUMMON_DAMAGE          EffectAuraType = 429
	A_PLAY_SCENE                 EffectAuraType = 430
	A_MOD_OVERRIDE_ZONE_PVP_TYPE EffectAuraType = 431
	// 432,433,434,435 are commented out
	A_MOD_ENVIRONMENTAL_DAMAGE_TAKEN EffectAuraType = 436
	A_MOD_MINIMUM_SPEED_RATE         EffectAuraType = 437
	A_PRELOAD_PHASE                  EffectAuraType = 438
	// 439 is commented out
	A_MOD_MULTISTRIKE_DAMAGE EffectAuraType = 440
	A_MOD_MULTISTRIKE_CHANCE EffectAuraType = 441
	A_MOD_READINESS          EffectAuraType = 442
	A_MOD_LEECH              EffectAuraType = 443
	// 444,445 are commented out
	A_SPELL_AURA_ADVANCED_FLYING EffectAuraType = 446
	A_MOD_XP_FROM_CREATURE_TYPE  EffectAuraType = 447
	// 448 is commented out (Related to PvP rules)
	// 449,450 are commented out
	A_OVERRIDE_PET_SPECS EffectAuraType = 451
	// 452 is commented out
	A_CHARGE_RECOVERY_MOD                     EffectAuraType = 453
	A_CHARGE_RECOVERY_MULTIPLIER              EffectAuraType = 454
	A_MOD_ROOT_2                              EffectAuraType = 455
	A_CHARGE_RECOVERY_AFFECTED_BY_HASTE       EffectAuraType = 456
	A_CHARGE_RECOVERY_AFFECTED_BY_HASTE_REGEN EffectAuraType = 457
	A_IGNORE_DUAL_WIELD_HIT_PENALTY           EffectAuraType = 458
	A_IGNORE_MOVEMENT_FORCES                  EffectAuraType = 459
	A_RESET_COOLDOWNS_ON_DUEL_START           EffectAuraType = 460
	// 461 is commented out
	A_MOD_HEALING_AND_ABSORB_FROM_CASTER       EffectAuraType = 462
	A_CONVERT_CRIT_RATING_PCT_TO_PARRY_RATING  EffectAuraType = 463
	A_MOD_ATTACK_POWER_OF_BONUS_ARMOR          EffectAuraType = 464
	A_MOD_BONUS_ARMOR                          EffectAuraType = 465
	A_MOD_BONUS_ARMOR_PCT                      EffectAuraType = 466
	A_MOD_STAT_BONUS_PCT                       EffectAuraType = 467
	A_TRIGGER_SPELL_ON_HEALTH_BELOW_PCT        EffectAuraType = 468
	A_SHOW_CONFIRMATION_PROMPT_WITH_DIFFICULTY EffectAuraType = 469
	A_MOD_AURA_TIME_RATE_BY_SPELL_LABEL        EffectAuraType = 470
	A_MOD_VERSATILITY                          EffectAuraType = 471
	// 472 is commented out
	A_PREVENT_DURABILITY_LOSS_FROM_COMBAT   EffectAuraType = 473
	A_REPLACE_ITEM_BONUS_TREE               EffectAuraType = 474
	A_ALLOW_USING_GAMEOBJECTS_WHILE_MOUNTED EffectAuraType = 475
	A_MOD_CURRENCY_GAIN_LOOTED_PCT          EffectAuraType = 476
	// 477,478,479 are commented out
	A_MOD_ARTIFACT_ITEM_LEVEL EffectAuraType = 480
	A_CONVERT_CONSUMED_RUNE   EffectAuraType = 481
	// 482 is commented out
	A_SUPPRESS_TRANSFORMS          EffectAuraType = 483
	A_ALLOW_INTERRUPT_SPELL        EffectAuraType = 484
	A_MOD_MOVEMENT_FORCE_MAGNITUDE EffectAuraType = 485
	// 486 is commented out
	A_COSMETIC_MOUNTED EffectAuraType = 487
	// 488 is commented out
	A_MOD_ALTERNATIVE_DEFAULT_LANGUAGE EffectAuraType = 489
	// 490,491 are commented out
	A_MOD_RESTED_XP_CONSUMPTION             EffectAuraType = 492
	A_MOD_RESTED_XP_CONSUMPTION_DUP         EffectAuraType = 493 // duplicate string?
	A_SET_POWER_POINT_CHARGE                EffectAuraType = 494
	A_TRIGGER_SPELL_ON_EXPIRE               EffectAuraType = 495
	A_ALLOW_CHANGING_EQUIPMENT_IN_TORGHAST  EffectAuraType = 496
	A_MOD_ANIMA_GAIN                        EffectAuraType = 497
	A_CURRENCY_LOSS_PCT_ON_DEATH            EffectAuraType = 498
	A_MOD_RESTED_XP_CONSUMPTION_2           EffectAuraType = 499 // differentiate duplicate
	A_IGNORE_SPELL_CHARGE_COOLDOWN          EffectAuraType = 500
	A_MOD_CRITICAL_DAMAGE_TAKEN_FROM_CASTER EffectAuraType = 501
	A_MOD_VERSATILITY_DAMAGE_DONE_BENEFIT   EffectAuraType = 502
	A_MOD_VERSATILITY_HEALING_DONE_BENEFIT  EffectAuraType = 503
	A_MOD_HEALING_TAKEN_FROM_CASTER         EffectAuraType = 504
	A_MOD_PLAYER_CHOICE_REROLLS             EffectAuraType = 505
	A_DISABLE_INERTIA                       EffectAuraType = 506
	A_MOD_DAMAGE_TAKEN_FROM_CASTER_BY_LABEL EffectAuraType = 507
	// 508,509 are commented out
	A_MODIFIED_RAID_INSTANCE  EffectAuraType = 510
	A_APPLY_PROFESSION_EFFECT EffectAuraType = 511
	A_CONVERT_RUNE            EffectAuraType = 512
	// 513-518 are commented out
	A_MOD_COOLDOWN_RECOVERY_RATE_ALL EffectAuraType = 519
	// 520-524 are commented out
	A_DISPLAY_PROFESSION_EQUIPMENT EffectAuraType = 525
	// 526,527 are commented out
	A_ALLOW_BLOCKING_SPELLS  EffectAuraType = 528
	A_MOD_SPELL_BLOCK_CHANCE EffectAuraType = 529
	// 530-535 are commented out
	A_IGNORE_SPELL_CREATURE_TYPE_REQUIREMENTS EffectAuraType = 536
	// 537 is commented out
	A_MOD_FAKE_INEBRIATION_MOVEMENT_ONLY  EffectAuraType = 538
	A_ALLOW_MOUNT_IN_COMBAT               EffectAuraType = 539
	A_MOD_SUPPORT_STAT                    EffectAuraType = 540
	A_MOD_REQUIRED_MOUNT_CAPABILITY_FLAGS EffectAuraType = 541
	// 542-546 are commented out
	A_MOD_CRIT_PERCENT_VERSUS EffectAuraType = 547
	A_MOD_RUNE_REGEN_SPEED    EffectAuraType = 548
	// 549,550 are commented out
	A_EXTRA_ATTACKS                EffectAuraType = 551
	A_MOD_SPELL_CRIT_CHANCE_SCHOOL EffectAuraType = 552
	A_MOD_POWER_COST_SCHOOL2       EffectAuraType = 553
	// 554,555 are commented out
	A_MOD_MELEE_DAMAGE_TAKEN2            EffectAuraType = 556
	A_MOD_RANGED_HASTE_QUIVER            EffectAuraType = 557
	A_MOD_RESISTANCE_EXCLUSIVE           EffectAuraType = 558
	A_MOD_PET_TALENT_POINTS              EffectAuraType = 559
	A_RETAIN_COMBO_POINTS                EffectAuraType = 560
	A_MOD_SHIELD_BLOCKVALUE_PCT2         EffectAuraType = 561
	A_SPLIT_DAMAGE_FLAT                  EffectAuraType = 562
	A_PET_DAMAGE_MULTI2                  EffectAuraType = 563
	A_MOD_SHIELD_BLOCKVALUE              EffectAuraType = 564
	A_SPELL_AURA_MOD_AOE_AVOIDANCE       EffectAuraType = 565
	A_MELEE_ATTACK_POWER_ATTACKER_BONUS2 EffectAuraType = 566
	// 567,568 are commented out
	A_MOD_ATTACKER_SPELL_CRIT_CHANCE EffectAuraType = 569
	// 570 is commented out
	A_MOD_RESISTANCE_OF_STAT_PERCENT   EffectAuraType = 571
	A_MOD_CRITICAL_THREAT              EffectAuraType = 572
	A_MOD_ATTACKER_RANGED_CRIT_CHANCE  EffectAuraType = 573
	A_MOD_TARGET_ABILITY_ABSORB_SCHOOL EffectAuraType = 574
	// 575,576 are commented out
	A_MOD_ATTACKER_MELEE_CRIT_DAMAGE  EffectAuraType = 577
	A_MOD_ATTACKER_RANGED_CRIT_DAMAGE EffectAuraType = 578
	A_MOD_SCHOOL_CRIT_DMG_TAKEN       EffectAuraType = 579
	// 580-582 are commented out
	A_MOD_RATING_FROM_STAT EffectAuraType = 583
	// 584 is commented out
	A_RAID_PROC_FROM_CHARGE EffectAuraType = 585
	// 586,587 are commented out
	A_MOD_DISPEL_RESIST                 EffectAuraType = 588
	A_MOD_SPELL_DAMAGE_OF_ATTACK_POWER  EffectAuraType = 589
	A_MOD_SPELL_HEALING_OF_ATTACK_POWER EffectAuraType = 590
	A_MOD_SCALE_3                       EffectAuraType = 591
	// 592 is commented out
	A_MOD_COMBAT_RESULT_CHANCE2         EffectAuraType = 593
	A_MOD_TARGET_RESIST_BY_SPELL_CLASS2 EffectAuraType = 594
	// 595-598 are commented out
	A_MOD_IGNORE_TARGET_RESIST2      EffectAuraType = 599
	A_SCHOOL_MASK_DAMAGE_FROM_CASTER EffectAuraType = 600
	A_IGNORE_MELEE_RESET             EffectAuraType = 601
	// 602,603 are commented out
	A_MOD_HONOR_GAIN_PCT2 EffectAuraType = 604
	// 605 is commented out
	A_MOD_BASE_HEALTH_PCT2      EffectAuraType = 606
	A_MOD_ATTACK_POWER_OF_ARMOR EffectAuraType = 607
	A_ABILITY_PERIODIC_CRIT     EffectAuraType = 608
	// 609-614 are commented out
	A_MOD_RANGED_HASTE_3 EffectAuraType = 615
	// 616-618 are commented out
	A_MOD_BLIND                  EffectAuraType = 619
	A_MOD_VENDOR_ITEMS_PRICES2   EffectAuraType = 620
	A_INCREASE_SKILL_GAIN_CHANCE EffectAuraType = 621
	// 622 is commented out
	A_MOD_GATHERING_ITEMS_GAINED_PERCENT EffectAuraType = 623
	A_MOD_DAMAGE_FROM_MANA               EffectAuraType = 624
)

type ConsumableClass int

const (
	EXPLOSIVES_AND_DEVICES ConsumableClass = iota
	POTION
	ELIXIR
	FLASK
	SCROLL
	FOOD
	ITEM_ENHANCEMENT
	BANDAGE
	OTHER
)

type RPPMModifierType int

const (
	RPPMModifierHaste     RPPMModifierType = iota + 1 // 1
	RPPMModifierCrit                                  // 2
	RPPMModifierClass                                 // 3
	RPPMModifierSpec                                  // 4
	RPPMModifierRace                                  // 5
	RPPMModifierIlevel                                // 6
	RPPMModifierUnkAdjust                             // 7
)

const (
	PROC_FLAG_NONE int = 0

	PROC_FLAG_HEARTBEAT int = 0x00000001 // 00 Heartbeat
	PROC_FLAG_KILL      int = 0x00000002 // 01 Kill target (in most cases need XP/Honor reward)

	PROC_FLAG_DEAL_MELEE_SWING int = 0x00000004 // 02 Done melee auto attack
	PROC_FLAG_TAKE_MELEE_SWING int = 0x00000008 // 03 Taken melee auto attack

	PROC_FLAG_DEAL_MELEE_ABILITY int = 0x00000010 // 04 Done attack by Spell that has dmg class melee
	PROC_FLAG_TAKE_MELEE_ABILITY int = 0x00000020 // 05 Taken attack by Spell that has dmg class melee

	PROC_FLAG_DEAL_RANGED_ATTACK int = 0x00000040 // 06 Done ranged auto attack
	PROC_FLAG_TAKE_RANGED_ATTACK int = 0x00000080 // 07 Taken ranged auto attack

	PROC_FLAG_DEAL_RANGED_ABILITY int = 0x00000100 // 08 Done attack by Spell that has dmg class ranged
	PROC_FLAG_TAKE_RANGED_ABILITY int = 0x00000200 // 09 Taken attack by Spell that has dmg class ranged

	PROC_FLAG_DEAL_HELPFUL_ABILITY int = 0x00000400 // 10 Done positive spell that has dmg class none
	PROC_FLAG_TAKE_HELPFUL_ABILITY int = 0x00000800 // 11 Taken positive spell that has dmg class none

	PROC_FLAG_DEAL_HARMFUL_ABILITY int = 0x00001000 // 12 Done negative spell that has dmg class none
	PROC_FLAG_TAKE_HARMFUL_ABILITY int = 0x00002000 // 13 Taken negative spell that has dmg class none

	PROC_FLAG_DEAL_HELPFUL_SPELL int = 0x00004000 // 14 Done positive spell that has dmg class magic
	PROC_FLAG_TAKE_HELPFUL_SPELL int = 0x00008000 // 15 Taken positive spell that has dmg class magic

	PROC_FLAG_DEAL_HARMFUL_SPELL int = 0x00010000 // 16 Done negative spell that has dmg class magic
	PROC_FLAG_TAKE_HARMFUL_SPELL int = 0x00020000 // 17 Taken negative spell that has dmg class magic

	PROC_FLAG_DEAL_HARMFUL_PERIODIC int = 0x00040000 // 18 Successful do periodic (damage)
	PROC_FLAG_TAKE_HARMFUL_PERIODIC int = 0x00080000 // 19 Taken spell periodic (damage)

	PROC_FLAG_TAKE_ANY_DAMAGE int = 0x00100000 // 20 Taken any damage

	PROC_FLAG_DEAL_HELPFUL_PERIODIC int = 0x00200000 // 21

	PROC_FLAG_MAIN_HAND_WEAPON_SWING int = 0x00400000 // 22 Done main-hand melee attacks (spell and autoattack)
	PROC_FLAG_OFF_HAND_WEAPON_SWING  int = 0x00800000 // 23 Done off-hand melee attacks (spell and autoattack)

	PROC_FLAG_ANY_DIRECT_TAKEN int = PROC_FLAG_TAKE_MELEE_SWING |
		PROC_FLAG_TAKE_MELEE_ABILITY |
		PROC_FLAG_TAKE_RANGED_ABILITY |
		PROC_FLAG_TAKE_HARMFUL_ABILITY |
		PROC_FLAG_TAKE_RANGED_ATTACK |
		PROC_FLAG_TAKE_HELPFUL_SPELL |
		PROC_FLAG_TAKE_HELPFUL_ABILITY |
		PROC_FLAG_TAKE_ANY_DAMAGE |
		PROC_FLAG_TAKE_HARMFUL_SPELL
	PROC_FLAG_ANY_DIRECT_DEALT int = PROC_FLAG_DEAL_MELEE_SWING |
		PROC_FLAG_DEAL_MELEE_ABILITY |
		PROC_FLAG_DEAL_RANGED_ATTACK |
		PROC_FLAG_DEAL_RANGED_ABILITY |
		PROC_FLAG_DEAL_HARMFUL_ABILITY |
		PROC_FLAG_DEAL_HARMFUL_SPELL
	PROC_FLAG_ANY_HEAL int = PROC_FLAG_DEAL_HELPFUL_PERIODIC |
		PROC_FLAG_DEAL_HELPFUL_ABILITY |
		PROC_FLAG_DEAL_HELPFUL_SPELL
)

const (
	ATTR_EX_3_CAN_PROC_FROM_PROCS int = 0x4000000
)

const (
	NORMAL_DUNGEON = 1 << iota
	HEROIC_DUNGEON
	NORMAL_RAID_10_MAN
	NORMAL_RAID_25_MAN
	HEROIC_RAID_10_MAN
	HEROIC_RAID_25_MAN
	LOOKING_FOR_RAID
	CHALLENGE_MODE
	NORMAL_RAID_40_MAN
)
