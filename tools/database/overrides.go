package database

import (
	"regexp"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

var OtherItemIdsToFetch = []string{
	// Hallow's End Ilvl bumped rings
	"211817",
	"211844",
	"211847",
	"211850",
	"211851",
}
var ConsumableOverrides = []*proto.Consumable{
	{Id: 62290, BuffsMainStat: true, Stats: stats.Stats{stats.Stamina: 90}.ToProtoArray()},
	{Id: 62649, BuffsMainStat: true, Stats: stats.Stats{stats.Stamina: 90}.ToProtoArray()},
}
var ItemOverrides = []*proto.UIItem{
	// Boosted 359 green weapon damage stats are way off
	// min and max show as double of their actual values in tooltips...
	// some scaling is happening that wowhead can't pick up

	// Balance T9 "of Conquest" Alliance set
	{Id: 48158, SetName: "Malfurion's Regalia"},
	{Id: 48159, SetName: "Malfurion's Regalia"},
	{Id: 48160, SetName: "Malfurion's Regalia"},
	{Id: 48161, SetName: "Malfurion's Regalia"},
	{Id: 48162, SetName: "Malfurion's Regalia"},

	// Death Knight T9 "of Conquest" Horde set
	{Id: 48501, SetName: "Koltira's Battlegear"},
	{Id: 48502, SetName: "Koltira's Battlegear"},
	{Id: 48503, SetName: "Koltira's Battlegear"},
	{Id: 48504, SetName: "Koltira's Battlegear"},
	{Id: 48505, SetName: "Koltira's Battlegear"},

	// Death Knight T9 "of Conquest" Tank Horde set
	{Id: 48558, SetName: "Koltira's Plate"},
	{Id: 48559, SetName: "Koltira's Plate"},
	{Id: 48560, SetName: "Koltira's Plate"},
	{Id: 48561, SetName: "Koltira's Plate"},
	{Id: 48562, SetName: "Koltira's Plate"},

	// Valorous T8 Sets
	{Id: 45375, Phase: 2},
	{Id: 45381, Phase: 2},
	{Id: 45382, Phase: 2},
	{Id: 45376, Phase: 2},
	{Id: 45370, Phase: 2},
	{Id: 45371, Phase: 2},
	{Id: 45383, Phase: 2},
	{Id: 45372, Phase: 2},
	{Id: 45377, Phase: 2},
	{Id: 45384, Phase: 2},
	{Id: 45379, Phase: 2},
	{Id: 45385, Phase: 2},
	{Id: 45380, Phase: 2},
	{Id: 45373, Phase: 2},
	{Id: 45374, Phase: 2},
	{Id: 45391, Phase: 2},
	{Id: 45386, Phase: 2},
	{Id: 45340, Phase: 2},
	{Id: 45335, Phase: 2},
	{Id: 45336, Phase: 2},
	{Id: 45341, Phase: 2},
	{Id: 45337, Phase: 2},
	{Id: 45342, Phase: 2},
	{Id: 45338, Phase: 2},
	{Id: 45343, Phase: 2},
	{Id: 45339, Phase: 2},
	{Id: 45344, Phase: 2},
	{Id: 45419, Phase: 2},
	{Id: 45417, Phase: 2},
	{Id: 45420, Phase: 2},
	{Id: 45421, Phase: 2},
	{Id: 45422, Phase: 2},
	{Id: 45387, Phase: 2},
	{Id: 45392, Phase: 2},
	{Id: 46131, Phase: 2},
	{Id: 45365, Phase: 2},
	{Id: 45367, Phase: 2},
	{Id: 45369, Phase: 2},
	{Id: 45368, Phase: 2},
	{Id: 45388, Phase: 2},
	{Id: 45393, Phase: 2},
	{Id: 46313, Phase: 2},
	{Id: 45351, Phase: 2},
	{Id: 45355, Phase: 2},
	{Id: 45345, Phase: 2},
	{Id: 45356, Phase: 2},
	{Id: 45346, Phase: 2},
	{Id: 45347, Phase: 2},
	{Id: 45357, Phase: 2},
	{Id: 45352, Phase: 2},
	{Id: 45358, Phase: 2},
	{Id: 45348, Phase: 2},
	{Id: 45359, Phase: 2},
	{Id: 45349, Phase: 2},
	{Id: 45353, Phase: 2},
	{Id: 45354, Phase: 2},
	{Id: 45394, Phase: 2},
	{Id: 45395, Phase: 2},
	{Id: 45389, Phase: 2},
	{Id: 45360, Phase: 2},
	{Id: 45361, Phase: 2},
	{Id: 45362, Phase: 2},
	{Id: 45363, Phase: 2},
	{Id: 45364, Phase: 2},
	{Id: 45390, Phase: 2},
	{Id: 45429, Phase: 2},
	{Id: 45424, Phase: 2},
	{Id: 45430, Phase: 2},
	{Id: 45425, Phase: 2},
	{Id: 45426, Phase: 2},
	{Id: 45431, Phase: 2},
	{Id: 45427, Phase: 2},
	{Id: 45432, Phase: 2},
	{Id: 45428, Phase: 2},
	{Id: 45433, Phase: 2},
	{Id: 45396, Phase: 2},
	{Id: 45397, Phase: 2},
	{Id: 45398, Phase: 2},
	{Id: 45399, Phase: 2},
	{Id: 45400, Phase: 2},
	{Id: 45413, Phase: 2},
	{Id: 45412, Phase: 2},
	{Id: 45406, Phase: 2},
	{Id: 45414, Phase: 2},
	{Id: 45401, Phase: 2},
	{Id: 45411, Phase: 2},
	{Id: 45402, Phase: 2},
	{Id: 45408, Phase: 2},
	{Id: 45409, Phase: 2},
	{Id: 45403, Phase: 2},
	{Id: 45415, Phase: 2},
	{Id: 45410, Phase: 2},
	{Id: 45404, Phase: 2},
	{Id: 45405, Phase: 2},
	{Id: 45416, Phase: 2},

	// Other items Wowhead has the wrong phase listed for
	// Ick's loot table from Pit of Saron
	{Id: 49812, Phase: 4},
	{Id: 49808, Phase: 4},
	{Id: 49811, Phase: 4},
	{Id: 49807, Phase: 4},
	{Id: 49810, Phase: 4},
	{Id: 49809, Phase: 4},

	// Drape of Icy Intent
	{Id: 45461, Phase: 2},

	// Cata pre-patch event items
	{Id: 53492, Phase: 5},

	// Heirloom Dwarven Handcannon, Wowhead partially glitchs out and shows us some other lvl calc for this
	{Id: 44093, Stats: stats.Stats{stats.CritRating: 30, stats.ResilienceRating: 13, stats.AttackPower: 34}.ToProtoArray()},

	// Dungeon items and quest rewards from patch 4.3
	{Id: 72880, Phase: 4}, // Alurmi's Ring
	{Id: 72852, Phase: 4}, // Archivist's Gloves
	{Id: 72874, Phase: 4}, // Boots of the Forked Road
	{Id: 72879, Phase: 4}, // Boots of the Treacherous Path
	{Id: 72873, Phase: 4}, // Bronze Blaster
	{Id: 72877, Phase: 4}, // Chain of the Demon Hunter
	{Id: 72882, Phase: 4}, // Chronicler's Chestguard
	{Id: 72887, Phase: 4}, // Cinch of the World Shaman
	{Id: 76152, Phase: 4}, // Cowl of Destiny
	{Id: 72871, Phase: 4}, // Crescent Wand
	{Id: 72878, Phase: 4}, // Demonic Skull
	{Id: 72883, Phase: 4}, // Historian's Sash
	{Id: 72876, Phase: 4}, // Ironfeather Longbow
	{Id: 66540, Phase: 4}, // Miniature Winter Veil Tree
	{Id: 72888, Phase: 4}, // Ring of the Loyal Companion
	{Id: 72858, Phase: 4}, // Safeguard Gloves
	{Id: 76153, Phase: 4}, // Signet of the Twilight Prophet
	{Id: 76155, Phase: 4}, // Thorns of the Dying Day
	{Id: 72886, Phase: 4}, // Thrall's Gratitude
	{Id: 72872, Phase: 4}, // Time Strand Gauntlets
	{Id: 72875, Phase: 4}, // Time Twister's Gauntlets
	{Id: 72881, Phase: 4}, // Treads of the Past
	{Id: 72884, Phase: 4}, // Writhing Wand
}

// Keep these sorted by item ID.
var ItemAllowList = map[int32]struct{}{
	2140: {},
	//Shaman Dungeon Set 3 Tidefury
	27510: {}, // Tidefury Gauntlets
	27802: {}, // Tidefury Shoulderguards
	27909: {}, // Tidefury Kilt
	28231: {}, // Tidefury Chestpiece
	28349: {}, // Tidefury Helm

	29309: {}, // Band of the Eternal Restorer

	31026: {}, // Slayer's Handguards
	31027: {}, // Slayer's Helm
	31028: {}, // Slayer's Chestguard
	31029: {}, // Slayer's Legguards
	31030: {}, // Slayer's Shoulderpads
	34448: {}, // Slayer's Bracers
	34558: {}, // Slayer's Belt
	34575: {}, // Slayer's Boots

	34677: {}, // Shattered Sun Pendant of Restoration

	45703: {}, // Spark of Hope
}

// Keep these sorted by item ID.
var ItemDenyList = map[int32]struct{}{
	17782: {}, // talisman of the binding shard
	17783: {}, // talisman of the binding fragment
	17802: {}, // Deprecated version of Thunderfury
	18582: {},
	18583: {},
	18584: {},
	24265: {},
	32384: {},
	32421: {},
	32422: {},
	33482: {},
	33350: {},
	34576: {}, // Battlemaster's Cruelty
	34577: {}, // Battlemaster's Depreavity
	34578: {}, // Battlemaster's Determination
	34579: {}, // Battlemaster's Audacity
	34580: {}, // Battlemaster's Perseverence

	38694: {}, // "Family" Shoulderpads heirloom
	45084: {}, // 'Book of Crafting Secrets' heirloom

	// '10 man' onyxia head rewards
	49312: {},
	49313: {},
	49314: {},

	50251: {}, // 'one hand shadows edge'
	53500: {}, // Tectonic Plate

	// Old Valentine's day event rewards
	51804: {},
	51805: {},
	51806: {},
	51807: {},
	51808: {},
	68172: {},
	68173: {},
	68174: {},
	68175: {},
	68176: {},

	48880: {}, // DK's Tier 9 Duplicates
	48881: {}, // DK's Tier 9 Duplicates
	48882: {}, // DK's Tier 9 Duplicates
	48883: {}, // DK's Tier 9 Duplicates
	48884: {}, // DK's Tier 9 Duplicates
	48885: {}, // DK's Tier 9 Duplicates
	48886: {}, // DK's Tier 9 Duplicates
	48887: {}, // DK's Tier 9 Duplicates
	48888: {}, // DK's Tier 9 Duplicates
	48889: {}, // DK's Tier 9 Duplicates
	48890: {}, // DK's Tier 9 Duplicates
	48891: {}, // DK's Tier 9 Duplicates
	48892: {}, // DK's Tier 9 Duplicates
	48893: {}, // DK's Tier 9 Duplicates
	48894: {}, // DK's Tier 9 Duplicates
	48895: {}, // DK's Tier 9 Duplicates
	48896: {}, // DK's Tier 9 Duplicates
	48897: {}, // DK's Tier 9 Duplicates
	48898: {}, // DK's Tier 9 Duplicates
	48899: {}, // DK's Tier 9 Duplicates
	68710: {}, // Stump of Time Duplicate (Not available ingame)
	68711: {}, // Mandala of Stirring Patterns Duplicate
	68712: {}, // Impatience of Youth Duplicate
	68713: {}, // Mirror of Broken Images Duplicate
	65104: {}, // DONTUSEUnheeded Warning
	65015: {}, // DONTUSEFury of Angerforge

	232548: {}, // The Horseman's Sinister Saber - 353
	232544: {}, // The Horseman's Horrific Helmet - 353
	232536: {}, // Band of Ghoulish Glee - 353
	232537: {}, // The Horseman's Signet - 353
	232540: {}, // Wicked Witch's Ring - 353
	232538: {}, // Seal of the Petrified Pumpkin - 353

	71331: {}, // Direbrew's Bloodied Shanker - 365
	71332: {}, // Tremendous Tankard O' Terror - 365
	71333: {}, // Bubblier Brightbrew Charm - 365
	71334: {}, // Bitterer Balebrew Charm - 365
	71335: {}, // Coren's Chilled Chromium Coaster - 365
	71336: {}, // Petrified Pickled Egg - 365
	71337: {}, // Mithril Stopwatch - 365
	71338: {}, // Brawler's Trophy - 365

	// T11 BoE items which are not available in the game
	65005: {}, // Claws of Agony - 372
	65006: {}, // Claws of Torment - 372
	65008: {}, // Shadowforge's Lightbound Smock - 372
	65009: {}, // Hide of Chromaggus - 372
	65010: {}, // Ironstar's Impenetrable Cover - 372
	65011: {}, // Corehammer's Riveted Girdle - 372
	65012: {}, // Treads of Savage Beatings - 372
	65013: {}, // Maldo's Sword Cane - 372
	65014: {}, // Maimgor's Bite - 372
	65016: {}, // Theresa's Booklight - 372
	65097: {}, // Bracers of the Dark Pool - 372
	65098: {}, // Crossfire Carbine - 372
	65099: {}, // Tsanga's Helm - 372
	65100: {}, // Phase-Twister Leggings - 372
	65101: {}, // Heaving Plates of Protection - 372
	65102: {}, // Chelley's Staff of Dark Mending - 372
	65103: {}, // Soul Blade - 372

	// Firelands "upgraded" items which are not available in the game
	69184: {}, // Stay of Execution - 391
	69185: {}, // Rune of Zeth - 391
	69198: {}, // Fiery Quintessence - 391
	69199: {}, // Ancient Petrified Seed - 391
	69200: {}, // Essence of the Eternal Flame - 391
	71388: {}, // Sleek Flamewrath Cloak - 391
	71389: {}, // Rippling Flamewrath Drape - 391
	71390: {}, // Flowing Flamewrath Cape - 391
	71391: {}, // Bladed Flamewrath Cover - 391
	71392: {}, // Durable Flamewrath Greatcloak - 391
	71393: {}, // Embereye Belt - 391
	71394: {}, // Flamebinding Girdle - 391
	71395: {}, // Firescar Sash - 391
	71396: {}, // Firearrow Belt - 391
	71397: {}, // Firemend Cinch - 391
	71398: {}, // Belt of the Seven Seals - 391
	71399: {}, // Cinch of the Flaming Ember - 391
	71400: {}, // Girdle of the Indomitable Flame - 391
	71565: {}, // Necklace of Smoke Signals - 391
	71566: {}, // Splintered Brimstone Seal - 391
	71569: {}, // Flamebinder Bracers, - 391
	71570: {}, // Bracers of Forked Lightning, - 391
	71571: {}, // Emberflame Bracers, - 391
	71572: {}, // Firesoul Wristguards, - 391
	71573: {}, // Amulet of Burning Brilliance - 391
	71574: {}, // Crystalline Brimstone Ring - 391
	71576: {}, // Firemind Pendant - 391
	71578: {}, // Soothing Brimstone Circle - 391
	71581: {}, // Smolderskull Bindings, - 391
	71582: {}, // Bracers of Misting Ash, - 391
	71583: {}, // Bracers of Imperious Truths - 391
	71584: {}, // Gigantiform Bracers, - 391
	71585: {}, // Bracers of Regal Force, - 391
	71586: {}, // Stoneheart Choker - 391
	71588: {}, // Serrated Brimstone Signet - 391
	71589: {}, // Stoneheart Necklace - 391
	71591: {}, // Deflecting Brimstone Band - 391

	// Patch 4.3 items not available in the game
	78610: {}, // Arrowflick Gauntlets - 384
	78527: {}, // Arrowflick Gauntlets - 410
	78601: {}, // Band of Reconstruction - 384
	78523: {}, // Band of Reconstruction - 410
	78587: {}, // Batwing Cloak - 384
	78509: {}, // Batwing Cloak - 410
	78640: {}, // Belt of Hidden Keys - 384
	78565: {}, // Belt of Hidden Keys - 410
	78641: {}, // Belt of Universal Curing - 384
	78566: {}, // Belt of Universal Curing - 410
	78591: {}, // Bladeshatter Treads - 384
	78511: {}, // Bladeshatter Treads - 410
	78644: {}, // Blinding Girdle of Truth - 384
	78563: {}, // Blinding Girdle of Truth - 410
	78583: {}, // Bones of the Damned - 384
	78499: {}, // Bones of the Damned - 410
	78596: {}, // Boneshard Boots - 384
	78512: {}, // Boneshard Boots - 410
	78592: {}, // Boots of Fungoid Growth - 384
	78517: {}, // Boots of Fungoid Growth - 410
	77985: {}, // Bottled Wishes - 384
	78005: {}, // Bottled Wishes - 410
	78654: {}, // Bracers of Manifold Pockets - 384
	78574: {}, // Bracers of Manifold Pockets - 410
	78655: {}, // Bracers of the Black Dream - 384
	78577: {}, // Bracers of the Black Dream - 410
	78651: {}, // Bracers of the Spectral Wolf - 384
	78572: {}, // Bracers of the Spectral Wolf - 410
	78650: {}, // Bracers of Unrelenting Excellence - 384
	78570: {}, // Bracers of Unrelenting Excellence - 410
	78622: {}, // Cameo of Terrible Memories - 384
	78546: {}, // Cameo of Terrible Memories - 410
	78584: {}, // Chestplate of the Unshakable Titan - 384
	78500: {}, // Chestplate of the Unshakable Titan - 410
	78656: {}, // Chronoboost Bracers - 384
	78576: {}, // Chronoboost Bracers - 410
	78608: {}, // Clockwinder's Immaculate Gloves - 384
	78532: {}, // Clockwinder's Immaculate Gloves - 410
	78642: {}, // Cord of Dragon Sinew - 384
	78561: {}, // Cord of Dragon Sinew - 410
	78638: {}, // Darting Chakram - 384
	78558: {}, // Darting Chakram - 410
	78582: {}, // Decaying Herbalist's Robes - 384
	78505: {}, // Decaying Herbalist's Robes - 410
	78645: {}, // Demonbone Waistguard - 384
	78564: {}, // Demonbone Waistguard - 410
	78653: {}, // Dragonbelly Bracers - 384
	78571: {}, // Dragonbelly Bracers - 410
	78579: {}, // Dragonflayer Vest - 384
	78501: {}, // Dragonflayer Vest - 410
	78586: {}, // Dreamcrusher Drape - 384
	78506: {}, // Dreamcrusher Drape - 410
	78599: {}, // Emergency Descent Loop - 384
	78524: {}, // Emergency Descent Loop - 410
	77988: {}, // Fire of the Deep - 384
	78008: {}, // Fire of the Deep - 410
	78648: {}, // Flashing Bracers of Warmth - 384
	78573: {}, // Flashing Bracers of Warmth - 410
	78646: {}, // Forgesmelter Waistplate - 384
	78560: {}, // Forgesmelter Waistplate - 410
	78604: {}, // Fungus-Born Gloves - 384
	78531: {}, // Fungus-Born Gloves - 410
	78606: {}, // Gauntlets of Feathery Blows - 384
	78526: {}, // Gauntlets of Feathery Blows - 410
	78580: {}, // Ghostworld Chestguard - 384
	78502: {}, // Ghostworld Chestguard - 410
	78643: {}, // Girdle of Shamanic Fury - 384
	78562: {}, // Girdle of Shamanic Fury - 410
	78612: {}, // Gleaming Grips of Mending - 384
	78529: {}, // Gleaming Grips of Mending - 410
	78611: {}, // Gloves of Ghostly Dreams - 384
	78528: {}, // Gloves of Ghostly Dreams - 410
	78621: {}, // Glowing Wings of Hope - 384
	78538: {}, // Glowing Wings of Hope - 410
	78605: {}, // Grimfist Crushers - 384
	78525: {}, // Grimfist Crushers - 410
	78623: {}, // Guardspike Choker - 384
	78544: {}, // Guardspike Choker - 410
	78629: {}, // Gutripper Shard - 384
	78550: {}, // Gutripper Shard - 410
	78649: {}, // Heartcrusher Wristplates - 384
	78569: {}, // Heartcrusher Wristplates - 410
	78618: {}, // Helmet of Perpetual Rebirth - 384
	78540: {}, // Helmet of Perpetual Rebirth - 410
	78616: {}, // Hood of Hidden Flesh - 384
	78541: {}, // Hood of Hidden Flesh - 410
	78627: {}, // Hungermouth Wand - 384
	78548: {}, // Hungermouth Wand - 410
	78589: {}, // Indefatigable Greatcloak - 384
	78507: {}, // Indefatigable Greatcloak - 410
	78615: {}, // Jaw of Repudiation - 384
	78535: {}, // Jaw of Repudiation - 410
	78597: {}, // Kavan's Forsaken Treads - 384
	78518: {}, // Kavan's Forsaken Treads - 410
	77984: {}, // Kiroptyric Sigil - 384
	78004: {}, // Kiroptyric Sigil - 410
	78590: {}, // Kneebreaker Boots - 384
	78515: {}, // Kneebreaker Boots - 410
	78609: {}, // Lightfinger Handwraps - 384
	78530: {}, // Lightfinger Handwraps - 410
	78631: {}, // Lightning Spirit in a Bottle - 384
	78552: {}, // Lightning Spirit in a Bottle - 410
	78635: {}, // Lightwarper Vestments - 384
	78556: {}, // Lightwarper Vestments - 410
	78652: {}, // Luminescent Bracers - 384
	78575: {}, // Luminescent Bracers - 410
	78630: {}, // Mindbender Lens - 384
	78553: {}, // Mindbender Lens - 410
	78588: {}, // Nanoprecise Cape - 384
	78510: {}, // Nanoprecise Cape - 410
	78625: {}, // Necklace of Black Dragon's Teeth - 384
	78543: {}, // Necklace of Black Dragon's Teeth - 410
	78617: {}, // Nocturnal Gaze - 384
	78539: {}, // Nocturnal Gaze - 410
	78624: {}, // Opal of the Secret Order - 384
	78547: {}, // Opal of the Secret Order - 410
	77986: {}, // Reflection of the Light - 384
	78006: {}, // Reflection of the Light - 410
	78603: {}, // Ring of Torn Flesh - 384
	78520: {}, // Ring of Torn Flesh - 410
	78633: {}, // Ripfang Relic - 384
	78554: {}, // Ripfang Relic - 410
	78634: {}, // Robes of Searing Shadow - 384
	78555: {}, // Robes of Searing Shadow - 410
	78594: {}, // Rooftop Griptoes - 384
	78516: {}, // Rooftop Griptoes - 410
	77987: {}, // Rotting Skull - 384
	78007: {}, // Rotting Skull - 410
	78595: {}, // Sabatons of the Graceful Spirit - 384
	78513: {}, // Sabatons of the Graceful Spirit - 410
	78628: {}, // Scintillating Rods - 384
	78549: {}, // Scintillating Rods - 410
	78600: {}, // Seal of the Grand Architect - 384
	78522: {}, // Seal of the Grand Architect - 410
	78581: {}, // Shadowbinder Chestguard - 384
	78504: {}, // Shadowbinder Chestguard - 410
	78578: {}, // Shining Carapace of Glory - 384
	78503: {}, // Shining Carapace of Glory - 410
	78602: {}, // Signet of the Resolute - 384
	78521: {}, // Signet of the Resolute - 410
	78593: {}, // Silver Sabatons of Fury - 384
	78514: {}, // Silver Sabatons of Fury - 410
	78620: {}, // Soulgaze Cowl - 384
	78542: {}, // Soulgaze Cowl - 410
	78598: {}, // Splinterfoot Sandals - 384
	78519: {}, // Splinterfoot Sandals - 410
	78632: {}, // Stoutheart Talisman - 384
	78551: {}, // Stoutheart Talisman - 410
	78639: {}, // Tentacular Belt - 384
	78567: {}, // Tentacular Belt - 410
	78607: {}, // The Hands of Gilly - 384
	78533: {}, // The Hands of Gilly - 410
	78626: {}, // Threadlinked Chain - 384
	78545: {}, // Threadlinked Chain - 410
	78636: {}, // Unexpected Backup - 384
	78557: {}, // Unexpected Backup - 410
	78647: {}, // Vestal's Irrepressible Girdle - 384
	78568: {}, // Vestal's Irrepressible Girdle - 410
	78614: {}, // Visage of Petrification - 384
	78534: {}, // Visage of Petrification - 410
	78637: {}, // Windslicer Boomerang - 384
	78559: {}, // Windslicer Boomerang - 410
	78613: {}, // Wolfdream Circlet - 384
	78537: {}, // Wolfdream Circlet - 410
	78585: {}, // Woundlicker Cover - 384
	78508: {}, // Woundlicker Cover - 410
	78619: {}, // Zeherah's Dragonskull Crown - 384
	78536: {}, // Zeherah's Dragonskull Crown - 410
}

// Item icons to include in the DB, so they don't need to be separately loaded in the UI.
var ExtraItemIcons = []int32{
	// Pet foods
	33874,
	43005,

	// Demonic Rune
	12662,

	// Food IDs
	27655,
	27657,
	27658,
	27664,
	33052,
	33825,
	33872,
	34753,
	34754,
	34756,
	34758,
	34767,
	34769,
	42994,
	42995,
	42996,
	42998,
	42999,
	43000,
	43015,
	62290,
	62649,
	62671,
	62670,
	62661,
	62665,
	62668,
	62669,
	62664,
	62666,
	62667,
	62662,
	62663,

	// Flask IDs
	13512,
	22851,
	22853,
	22854,
	22861,
	22866,
	33208,
	40079,
	44939,
	46376,
	46377,
	46378,
	46379,

	// Elixer IDs
	40072,
	40078,
	40097,
	40109,
	44328,
	44332,

	// Elixer IDs
	13452,
	13454,
	22824,
	22827,
	22831,
	22833,
	22834,
	22835,
	22840,
	28103,
	28104,
	31679,
	32062,
	32067,
	32068,
	39666,
	40068,
	40070,
	40073,
	40076,
	44325,
	44327,
	44329,
	44330,
	44331,

	// Potions / In Battle Consumes
	13442,
	20520,
	22105,
	22788,
	22828,
	22832,
	22837,
	22838,
	22839,
	22849,
	31677,
	33447,
	33448,
	36892,
	40093,
	40211,
	40212,
	40536,
	40771,
	41119,
	41166,
	42545,
	42641,

	// Poisons
	43231,
	43233,
	43235,

	// Thistle Tea
	7676,

	// Scrolls
	37094,
	43466,
	43464,
	37092,
	37098,
	43468,

	// Drums
	49633,
	49634,
}

// Item Ids of consumables to allow
var ConsumableAllowList = []int32{
	//Fortune Cookie and Feast
	62649,
	62290,
	//Migty Rage Potion
	13442,
	// Dark Rune
	20520,
	46376, // Flask of the Frost Wyrm
	45568, // Firecracker Salmon
	54221, // Potion of Speed
}
var ConsumableDenyList = []int32{
	57099,
}

// Raid buffs / debuffs
var SharedSpellsIcons = []int32{
	// Revitalize, Rejuv, WG
	48545,
	26982,
	53251,

	// Registered CD's
	49016,
	57933,
	64382,
	10060,
	16190,
	29166,
	53530,
	33206,
	2825,
	54758,

	// Raid Buffs
	43002,
	57567,
	54038,

	48470,
	17051,

	25898,
	25899,

	48942,
	20140,
	8071,
	16293,

	48161,
	14767,

	8075,
	52456,
	57623,

	48073,

	48934,
	20045,
	47436,

	53138,
	30808,
	19506,

	31869,
	31583,
	34460,

	57472,
	50720,

	53648,

	47440,
	12861,
	47982,
	18696,

	48938,
	20245,
	5675,
	16206,

	17007,
	34300,
	29801,

	55610,
	8512,
	29193,

	48160,
	31878,
	53292,
	54118,
	44561,

	24907,
	48396,
	51470,

	3738,
	47240,
	57722,
	8227,

	54043,
	48170,
	31025,
	31035,
	6562,
	31033,
	53307,
	16840,
	54648,

	// Raid Debuffs
	8647,
	47467,
	55749,

	770,
	33602,
	702,
	18180,
	56631,
	53598,

	26016,
	47437,
	12879,
	48560,
	16862,
	55487,

	48566,
	46855,
	57386,

	30706,
	20337,
	58410,

	47502,
	12666,
	55095,
	51456,
	53696,
	48485,

	3043,
	29859,
	58413,
	65855,

	17800,
	17803,
	12873,
	28593,

	33198,
	51161,
	48511,
	1490,

	20271,
	53408,

	11374,
	15235,

	27013,

	58749,
	49071,

	30708,
}

// If any of these match the item name, don't include it.
var DenyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile(`30 Epic`),
	regexp.MustCompile(`130 Epic`),
	regexp.MustCompile(`63 Blue`),
	regexp.MustCompile(`63 Green`),
	regexp.MustCompile(`66 Epic`),
	regexp.MustCompile(`90 Epic`),
	regexp.MustCompile(`90 Green`),
	regexp.MustCompile(`Boots 1`),
	regexp.MustCompile(`Boots 2`),
	regexp.MustCompile(`Boots 3`),
	regexp.MustCompile(`Bracer 1`),
	regexp.MustCompile(`Bracer 2`),
	regexp.MustCompile(`Bracer 3`),
	regexp.MustCompile(`DB\d`),
	regexp.MustCompile(`DEPRECATED`),
	regexp.MustCompile(`OLD`),
	regexp.MustCompile(`Deprecated`),
	regexp.MustCompile(`Deprecated: Keanna`),
	regexp.MustCompile(`Indalamar`),
	regexp.MustCompile(`Monster -`),
	regexp.MustCompile(`NEW`),
	regexp.MustCompile(`PH`),
	regexp.MustCompile(`QR XXXX`),
	regexp.MustCompile(`TEST`),
	regexp.MustCompile(`Test`),
	regexp.MustCompile(`Enchant Template`),
	regexp.MustCompile(`Arcane Amalgamation`),
	regexp.MustCompile(`Deleted`),
	regexp.MustCompile(`DELETED`),
	regexp.MustCompile(`zOLD`),
	regexp.MustCompile(`Archaic Spell`),
	regexp.MustCompile(`Well Repaired`),
}

// Allows manual overriding for Gem fields in case WowHead is wrong.
var GemOverrides = []*proto.UIGem{
	{Id: 33131, Stats: stats.Stats{stats.AttackPower: 32, stats.RangedAttackPower: 32}.ToProtoArray()},
}
var GemAllowList = map[int32]struct{}{
	//22459: {}, // Void Sphere
	//36766: {}, // Bright Dragon's Eye
	//36767: {}, // Solid Dragon's Eye
}
var GemDenyList = map[int32]struct{}{
	// pvp non-unique gems not in game currently.
	32735: {},
	34142: {}, // Infinite Sphere
	34143: {}, // Chromatic Sphere
	35489: {},
	37430: {}, // Solid Sky Sapphire (Unused)
	38545: {},
	38546: {},
	38547: {},
	38548: {},
	38549: {},
	38550: {},
	63696: {},
	63697: {},
}
