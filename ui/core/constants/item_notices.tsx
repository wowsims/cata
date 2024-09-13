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

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
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
