import { ItemNoticeData, SetBonusNoticeData } from '../components/item_notice/item_notice';
import { Spec } from '../proto/common';

const WantToHelpMessage = () => <p className="mb-0">Want to help out by providing additional information? Contact us on our Discord!</p>;

const MISSING_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">This item is not implemented!</p>
		<p>We are working hard on gathering all the old resources to allow for an initial implementation.</p>
		<WantToHelpMessage />
	</>
);

const VPLC_IMPLEMENTATION_WARNING = (
	<>
		<p>
			Current proc rate is <span className="fw-bold">50%</span> confirmed on PTR.
		</p>
		<p>Scales with: 3% Damage buff & 8% Spell Damage debuff.</p>
		<WantToHelpMessage />
	</>
);

const TENTATIVE_IMPLEMENTATION_WARNING = (
	<>
		<p>
			This item <span className="fw-bold">is</span> implemented, but detailed proc behavior will be confirmed on PTR.
		</p>
		<WantToHelpMessage />
	</>
);

const ITEM_DOESNT_EXIST_WARNING = <>
	<p>
		This item never existed in the original game, therefore any effects or procs it might have are not implemented.
	</p>
	<p>
		Once we get a clear indication from Blizzard whether they decide to include it or not, we will either implement it and remove this notice or remove the item entirely.
	</p>
	<WantToHelpMessage />
</>;

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
	// Gigantiform Bracers, - 391
	71584,
	// Bracers of Forked Lightning, - 391
	71570,
	// Bracers of Misting Ash, - 391
	71582,
	// Flamebinder Bracers, - 391
	71569,
	// Firesoul Wristguards, - 391
	71572,
	// Emberflame Bracers, - 391
	71571,
	// Bracers of Regal Force, - 391
	71585,
	// Smolderskull Bindings, - 391
	71581,
	// Bracers of Imperious Truths - 391
	71583,
	// Crystalline Brimstone Ring - 391
	71574,
	// Soothing Brimstone Circle - 391
	71578,
	// Splintered Brimstone Seal - 391
	71566,
	// Serrated Brimstone Signet - 391
	71588,
	// Deflecting Brimstone Band - 391
	71591,
	// Necklace of Smoke Signals - 391
	71565,
	// Stoneheart Necklace - 391
	71589,
	// Stoneheart Choker - 391
	71586,
	// Amulet of Burning Brilliance - 391
	71573,
	// Firemind Pendant - 391
	71576,
	// Rotting Skull - 384, 410
	77987, 78007,
	// Kiroptyric Sigil - 384, 410
	77984, 78004,
	// Fire of the Deep - 384, 410
	77988, 78008,
	// Reflection of the Light - 384, 410
	77986, 78006,
	// Bottled Wishes - 384, 410
	77985, 78005,
];

const WILL_NOT_BE_IMPLEMENTED_WARNING = <>The equip/use effect on this item is will not be implemented!</>;

const WILL_NOT_BE_IMPLEMENTED_ITEMS = [
	// Eye of Blazing Power - Normal, Heroic
	68983, 69149,
	// Windward Heart - LFR, Normal, Heroic
	77981, 77209, 78001,
	// Heart of Unliving - LFR, Normal, Heroic
	77976, 77199, 77996,
	// Seal of the Seven Signs - LFR, Normal, Heroic
	77969, 77204, 77989,

	// Maw of the Dragonlord - LFR, Normal, Heroic
	78485, 77196, 78476,
];

const TENTATIVE_IMPLEMENTATION_ITEMS = [
	// Vial of Shadows - LFR, Normal, Heroic
	77979, 77207, 77999,
	// Bone-Link Fetish - LFR, Normal, Heroic
	77982, 77210, 78002,
	// Cunning of the Cruel - LFR, Normal, Heroic
	77980, 77208, 78000,
	// Indomitable Pride - LFR, Normal, Heroic
	77983, 77211, 78003,
	// Creche of the Final Dragon - LFR, Normal, Heroic
	77972, 77205, 77992,
	// Insignia of the Corrupted Mind - LFR, Normal, Heroic
	77971, 77203, 77991,
	// Soulshifter Vortex - LFR, Normal, Heroic
	77970, 77206, 77990,
	// Starcatcher Compass - LFR, Normal, Heroic
	77973, 77202, 77993,
	// Eye of Unmaking - LFR, Normal, Heroic
	77977, 77200, 77997,
	// Resolve of Undying - LFR, Normal, Heroic
	77978, 77201, 77998,
	// Will of Unbinding - LFR, Normal, Heroic
	77975, 77198, 77995,
	// Wrath of Unchaining - LFR, Normal, Heroic
	77974, 77197, 77994,

	// Kiril, Fury of Beasts - LFR, Normal, Heroic
	78482, 77194, 78473,

	// Ti'tahk, the Steps of Time - LFR, Normal, Heroic
	78486, 77190, 78477,

	// Veil of Lies
	72900,
	// Foul Gift of the Demon Lord
	72898,
	// Arrow of Time
	72897,
	// Rosary of Light
	72901,
	// Varo'then's Brooch
	72899,
];

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
	...NON_EXISTING_ITEMS.map((itemID): [number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: ITEM_DOESNT_EXIST_WARNING,
		},
	]),
	...WILL_NOT_BE_IMPLEMENTED_ITEMS.map((itemID): [number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: WILL_NOT_BE_IMPLEMENTED_WARNING,
		},
	]),
	...TENTATIVE_IMPLEMENTATION_ITEMS.map((itemID): [number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: TENTATIVE_IMPLEMENTATION_WARNING,
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
		{
			[Spec.SpecUnknown]: VPLC_IMPLEMENTATION_WARNING,
			[Spec.SpecArcaneMage]: false,
			[Spec.SpecFireMage]: false,
			[Spec.SpecDemonologyWarlock]: false,
			[Spec.SpecAfflictionWarlock]: false,
			[Spec.SpecDestructionWarlock]: false,
			[Spec.SpecShadowPriest]: false,
		},
	],
	[
		// VPLC - Heroic
		69110,
		{
			[Spec.SpecUnknown]: VPLC_IMPLEMENTATION_WARNING,
			[Spec.SpecArcaneMage]: false,
			[Spec.SpecFireMage]: false,
			[Spec.SpecDemonologyWarlock]: false,
			[Spec.SpecAfflictionWarlock]: false,
			[Spec.SpecDestructionWarlock]: false,
			[Spec.SpecShadowPriest]: false,
		},
	],
]);

export const GENERIC_MISSING_SET_BONUS_NOTICE_DATA = new Map<number, string>([
	[2, 'Not yet implemented'],
	[4, 'Not yet implemented'],
]);

export const IGNITE_BONUS_WARNING = new Map<number, string>([[2, 'Tentative implementation, requires PTR testing to validate munching / roll-over model']]);

export const SET_BONUS_NOTICES = new Map<number, SetBonusNoticeData>([
	// Custom notices
	[1002, IGNITE_BONUS_WARNING], // Feral T12
	[1006, IGNITE_BONUS_WARNING], // Rogue T12
	[1012, IGNITE_BONUS_WARNING], // Ret T12
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
