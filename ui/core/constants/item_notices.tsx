import { ItemNoticeData, SetBonusNoticeData } from '../components/item_notice/item_notice';
import { Spec } from '../proto/common';

const MISSING_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">This item is not implemented!</p>
		<p>We are working hard on gathering all the old resources to allow for an initial implementation.</p>
		<p className="mb-0">Want to help out by providing additional information? Contact us on our Discord!</p>
	</>
);

const PROC_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">The proc rate for this item might not be correct!</p>
		<p>Current proc rate is 50%, based on old video data.</p>
		<p className="mb-0">Want to help out by providing additional information? Contact us on our Discord!</p>
	</>
);

const ITEM_DOESNT_EXIST_WARNING = <>This item didn't exist in the game and thus is not implemented!</>;
const ITEM_NO_ARMOR_WARNING = <>This item is missing Armor values.</>;

const NON_EXISTING_ITEMS = [
	// Rune of Zeth - 391
	69185,
	// Bladed Flamewrath Cover - 391
	71391,
	// Sleek Flamewrath Cloak - 391
	71388,
	// Cinch of the Flaming Ember - 391
	71399,
	// Fiery Quintessence - 391
	69198,
	// Rippling Flamewrath Drape - 391
	71389,
	// Ancient Petrified Seed - 391
	69199,
	// Essence of the Eternal Flame - 391
	69200,
	// Flowing Flamewrath Cape - 391
	71390,
	// Firemend Cinch - 391
	71397,
	// Stay of Execution - 391
	69184,
	// Durable Flamewrath Greatcloak - 391
	71392,
	// Embereye Belt - 391
	71393,
	// Flamebinding Girdle - 391
	71394,
	// Firearrow Belt - 391
	71396,
	// Firescar Sash - 391
	71395,
	// Girdle of the Indomitable Flame - 391
	71400,
	// Belt of the Seven Seals - 391
	71398,
];

const NO_ARMOR_VALUE_ITEMS = [
	// Sleek Flamewrath Cloak - 391
	71388,
	// Rippling Flamewrath Drape - 391
	71389,
	// Flowing Flamewrath Cape - 391
	71390,
	// Bladed Flamewrath Cover - 391
	71391,
	// Embereye Belt - 391
	71393,
	// Flamebinding Girdle - 391
	71394,
	// Firescar Sash - 391
	71395,
	// Firearrow Belt - 391
	71396,
	// Firemend Cinch - 391
	71397,
	// Belt of the Seven Seals - 391
	71398,
	// Cinch of the Flaming Ember - 391
	71399,
	// Girdle of the Indomitable Flame - 391
	71400,
	// Cinderweb Leggings - 391
	71402,
	// Flickering Shoulders - 391
	71403,
	// Arachnaflame Treads - 391
	71404,
	// Carapace of Imbibed Flame - 391
	71405,
	// Robes of Smoldering Devastation - 391
	71407,
	// Ward of the Red Widow - 391
	71408,
	// Cindersilk Gloves - 391
	71410,
	// Cowl of the Clicking Menace - 391
	71411,
	// Thoracic Flame Kilt - 391
	71412,
	// Spaulders of Manifold Eyes - 391
	71413,
	// Dreadfire Drape - 391
	71415,
	// Hood of Rampant Disdain - 391
	71416,
	// Flaming Core Chestguard - 391
	71417,
	// Earthcrack Bracers - 391
	71418,
	// Fireskin Gauntlets - 391
	71419,
	// Cracked Obsidian Stompers - 391
	71420,
	// Flickering Cowl - 391
	71421,
	// Incendic Chestguard - 391
	71424,
	// Lava Line Wristbands - 391
	71425,
	// Grips of the Raging Giant - 391
	71426,
	// Flickering Wristbands - 391
	71428,
	// Moltenfeather Leggings - 391
	71429,
	// Greathelm of the Voracious Maw - 391
	71430,
	// Lavaworm Legplates - 391
	71431,
	// Spaulders of Recurring Flame - 391
	71432,
	// Wings of Flame - 391
	71434,
	// Leggings of Billowing Fire - 391
	71435,
	// Phoenix-Down Treads - 391
	71436,
	// Clawshaper Gauntlets - 391
	71437,
	// Craterflame Spaulders - 391
	71438,
	// Clutch of the Firemother - 391
	71439,
	// Gloves of Dissolving Smoke - 391
	71440,
	// Scalp of the Bandit Prince - 391
	71442,
	// Uncrushable Belt of Fury - 391
	71443,
	// Legplates of Frenzied Devotion - 391
	71444,
	// Coalwalker Sandals - 391
	71447,
	// Flickering Shoulderpads - 391
	71450,
	// Treads of Implicit Obedience - 391
	71451,
	// Bracers of the Dread Hunter - 391
	71452,
	// Legplates of Absolute Control - 391
	71453,
	// Breastplate of the Incendiary Soul - 391
	71455,
	// Shoulderpads of the Forgotten Gate - 391
	71456,
	// Decimation Treads - 391
	71457,
	// Flickering Handguards - 391
	71458,
	// Helm of Blazing Glory - 391
	71459,
	// Shard of Torment - 391
	71460,
	// Mantle of Closed Doors - 391
	71461,
	// Glowing Wing Bracers - 391
	71463,
	// Gatekeeper's Embrace - 391
	71464,
	// Casque of Flame - 391
	71465,
	// Sandals of Leaping Coals - 391
	71467,
	// Grips of Unerring Precision - 391
	71468,
	// Breastplate of Shifting Visions - 391
	71469,
	// Bracers of the Fiery Path - 391
	71470,
	// Wristwraps of Arrogant Doom - 391
	71471,
	// Firecat Leggings - 391
	71474,
	// Treads of the Penitent Man - 391
	71475,
	// Elementium Deathplate Breastplate - 391
	71476,
	// Elementium Deathplate Gauntlets - 391
	71477,
	// Elementium Deathplate Helmet - 391
	71478,
	// Elementium Deathplate Greaves - 391
	71479,
	// Elementium Deathplate Pauldrons - 391
	71480,
	// Elementium Deathplate Chestguard - 391
	71481,
	// Elementium Deathplate Handguards - 391
	71482,
	// Elementium Deathplate Faceguard - 391
	71483,
	// Elementium Deathplate Legguards - 391
	71484,
	// Elementium Deathplate Shoulderguards - 391
	71485,
	// Obsidian Arborweave Raiment - 391
	71486,
	// Obsidian Arborweave Grips - 391
	71487,
	// Obsidian Arborweave Headpiece - 391
	71488,
	// Obsidian Arborweave Spaulders - 391
	71490,
	// Obsidian Arborweave Handwraps - 391
	71491,
	// Obsidian Arborweave Helm - 391
	71492,
	// Obsidian Arborweave Legwraps - 391
	71493,
	// Obsidian Arborweave Mantle - 391
	71495,
	// Obsidian Arborweave Gloves - 391
	71496,
	// Obsidian Arborweave Cover - 391
	71497,
	// Obsidian Arborweave Leggings - 391
	71498,
	// Obsidian Arborweave Vestment - 391
	71499,
	// Obsidian Arborweave Shoulderwraps - 391
	71500,
	// Flamewaker's Tunic - 391
	71501,
	// Flamewaker's Gloves - 391
	71502,
	// Flamewaker's Headguard - 391
	71503,
	// Flamewaker's Legguards - 391
	71504,
	// Flamewaker's Spaulders - 391
	71505,
	// Firehawk Gloves - 391
	71507,
	// Firehawk Hood - 391
	71508,
	// Firehawk Leggings - 391
	71509,
	// Firehawk Robes - 391
	71510,
	// Firehawk Mantle - 391
	71511,
	// Immolation Battleplate - 391
	71512,
	// Immolation Gauntlets - 391
	71513,
	// Immolation Helmet - 391
	71514,
	// Immolation Legplates - 391
	71515,
	// Immolation Pauldrons - 391
	71516,
	// Immolation Breastplate - 391
	71517,
	// Immolation Gloves - 391
	71518,
	// Immolation Headguard - 391
	71519,
	// Immolation Greaves - 391
	71520,
	// Immolation Mantle - 391
	71521,
	// Immolation Chestguard - 391
	71522,
	// Immolation Handguards - 391
	71523,
	// Immolation Faceguard - 391
	71524,
	// Immolation Legguards - 391
	71525,
	// Immolation Shoulderguards - 391
	71526,
	// Handwraps of the Cleansing Flame - 391
	71527,
	// Cowl of the Cleansing Flame - 391
	71528,
	// Legwraps of the Cleansing Flame - 391
	71529,
	// Robes of the Cleansing Flame - 391
	71530,
	// Mantle of the Cleansing Flame - 391
	71531,
	// Gloves of the Cleansing Flame - 391
	71532,
	// Hood of the Cleansing Flame - 391
	71533,
	// Leggings of the Cleansing Flame - 391
	71534,
	// Vestment of the Cleansing Flame - 391
	71535,
	// Shoulderwraps of the Cleansing Flame - 391
	71536,
	// Dark Phoenix Tunic - 391
	71537,
	// Dark Phoenix Gloves - 391
	71538,
	// Dark Phoenix Helmet - 391
	71539,
	// Dark Phoenix Legguards - 391
	71540,
	// Dark Phoenix Spaulders - 391
	71541,
	// Erupting Volcanic Tunic - 391
	71542,
	// Erupting Volcanic Handwraps - 391
	71543,
	// Erupting Volcanic Faceguard - 391
	71544,
	// Erupting Volcanic Legwraps - 391
	71545,
	// Erupting Volcanic Mantle - 391
	71546,
	// Erupting Volcanic Cuirass - 391
	71547,
	// Erupting Volcanic Grips - 391
	71548,
	// Erupting Volcanic Helmet - 391
	71549,
	// Erupting Volcanic Legguards - 391
	71550,
	// Erupting Volcanic Spaulders - 391
	71551,
	// Erupting Volcanic Hauberk - 391
	71552,
	// Erupting Volcanic Gloves - 391
	71553,
	// Erupting Volcanic Headpiece - 391
	71554,
	// Erupting Volcanic Kilt - 391
	71555,
	// Erupting Volcanic Shoulderwraps - 391
	71556,
	// Balespider's Handwraps - 391
	71594,
	// Balespider's Hood - 391
	71595,
	// Balespider's Leggings - 391
	71596,
	// Balespider's Robes - 391
	71597,
	// Balespider's Mantle - 391
	71598,
	// Helmet of the Molten Giant - 391
	71599,
	// Battleplate of the Molten Giant - 391
	71600,
	// Gauntlets of the Molten Giant - 391
	71601,
	// Legplates of the Molten Giant - 391
	71602,
	// Pauldrons of the Molten Giant - 391
	71603,
	// Chestguard of the Molten Giant - 391
	71604,
	// Handguards of the Molten Giant - 391
	71605,
	// Faceguard of the Molten Giant - 391
	71606,
	// Legguards of the Molten Giant - 391
	71607,
	// Shoulderguards of the Molten Giant - 391
	71608,
	// Majordomo's Chain of Office - 397
	71613,
	// Fingers of Incineration - 397
	71614,
	// Crown of Flame - 397
	71616,
].filter(id => !NON_EXISTING_ITEMS.find(i => id === i));

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
	...NON_EXISTING_ITEMS.map((itemID): [number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: ITEM_DOESNT_EXIST_WARNING,
		},
	]),
	...NO_ARMOR_VALUE_ITEMS.map((itemID): [number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: ITEM_NO_ARMOR_WARNING,
		},
	]),
	// Dragonwrath, Tarecgosa's Rest
	[
		71086,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecBalanceDruid]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecArcaneMage]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecFireMage]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecFrostMage]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecShadowPriest]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecElementalShaman]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecAfflictionWarlock]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecDemonologyWarlock]: MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecDestructionWarlock]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	// Rogue Legendary Daggers (All Stages)
	[
		// Fear
		77945,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// Vengeance
		77946,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// Sleeper
		77947,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// Dreamer
		77948,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// Golad
		77949,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// Tiriosh
		77950,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING,
		},
	],
	[
		// VPLC - Normal
		68925,
		{ [Spec.SpecUnknown]: PROC_IMPLEMENTATION_WARNING },
	],
	[
		// VPLC - Heroic
		69110,
		{ [Spec.SpecUnknown]: PROC_IMPLEMENTATION_WARNING },
	],
]);

export const GENERIC_MISSING_SET_BONUS_NOTICE_DATA = new Map<number, string>([
	[2, 'Not yet implemented'],
	[4, 'Not yet implemented'],
]);

export const SET_BONUS_NOTICES = new Map<number, SetBonusNoticeData>([
	// Custom notices
	[
		1002, // Feral T12
		new Map<number, string>([
			[2, 'Not implemented, requires PTR testing!'],
			[4, 'Implemented and working for both cat and bear'],
		]),
	],
	[
		1008, // Warlock T12
		new Map<number, string>([
			[2, 'Requires PTR testing to confirm exact pet behaviour & stats'],
			[4, 'Exact proc behaviour may vary, needs PTR testing to confirm'],
		]),
	],

	// Generic "not yet implemented" notices
	[1058, null], // Feral T13
]);
