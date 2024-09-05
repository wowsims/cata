import { ItemNoticeData, SetBonusNoticeData } from '../components/item_notice/item_notice';
import { Spec } from '../proto/common';

const DTR_MISSING_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">This item is not implemented!</p>
		<p>We are working hard on gathering all the old resources to allow for an initial implementation.</p>
		<p className="mb-0">Want to help out by providing additional information? Contact us on our Discord!</p>
	</>
);

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
	// Dragonwrath, Tarecgosa's Rest
	[
		71086,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecBalanceDruid]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecArcaneMage]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecFireMage]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecFrostMage]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecShadowPriest]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecElementalShaman]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecAfflictionWarlock]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecDemonologyWarlock]: DTR_MISSING_IMPLEMENTATION_WARNING,
			[Spec.SpecDestructionWarlock]: DTR_MISSING_IMPLEMENTATION_WARNING,
		},
	],
]);

export const GENERIC_MISSING_SET_BONUS_NOTICE_DATA = new Map<number, string>([
	[2, "Not yet implemented"],
	[4, "Not yet implemented"],
]);

export const SET_BONUS_NOTICES = new Map<number, SetBonusNoticeData>([
	// Custom notices
	[
		1002, // Feral T12
		new Map<number, string>([
			[2, "Requires PTR testing for implementation"],
			[4, "Will be implemented in the next few days"],
		]),
	],
	[
		1008, // Warlock T12
		new Map<number, string>([
			[2, "Requires PTR testing to confirm exact pet behaviour & stats"],
			[4, "Exact proc behaviour may vary, needs PTR testing to confirm"],
		]),
	],

	// Generic "not yet implemented" notices
	[1058, null], // Feral T13
]);
