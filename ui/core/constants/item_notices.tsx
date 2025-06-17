import { ItemNoticeData, SetBonusNoticeData } from '../components/item_notice/item_notice';
import { Spec } from '../proto/common';
import { MISSING_ITEM_EFFECTS } from './missing_effects_auto_gen';

const WantToHelpMessage = () => <p className="mb-0">Want to help out by providing additional information? Contact us on our Discord!</p>;

export const MISSING_RANDOM_SUFFIX_WARNING = <p className="mb-0">Please select a random suffix</p>;

const MISSING_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">This item is not implemented!</p>
		<p>We are working hard on gathering all the old resources to allow for an initial implementation.</p>
		<WantToHelpMessage />
	</>
);

// const DTR_FIRST_IMPLEMENTATION_WARNING = (
// 	<>
// 		<p className="fw-bold">This item was implemented based on the first round of testing on PTR.</p>
// 		<p>Results may change as we get more logs and reports on interactions.</p>
// 		<WantToHelpMessage />
// 	</>
// );

const TENTATIVE_IMPLEMENTATION_WARNING = (
	<>
		<p>
			This item <span className="fw-bold">is</span> implemented, but detailed proc behavior will be confirmed on PTR.
		</p>
		<WantToHelpMessage />
	</>
);

const WILL_NOT_BE_IMPLEMENTED_WARNING = <>The equip/use effect on this item will not be implemented!</>;

const WILL_NOT_BE_IMPLEMENTED_ITEMS = [
	// Eye of Blazing Power - Normal, Heroic
	68983, 69149,
	// Windward Heart - LFR, Normal, Heroic
	77981, 77209, 78001,
	// Heart of Unliving - LFR, Normal, Heroic
	77976, 77199, 77996,
	// Maw of the Dragonlord - LFR, Normal, Heroic
	78485, 77196, 78476,
];

const TENTATIVE_IMPLEMENTATION_ITEMS = [
	// Veil of Lies
	72900,
	// Arrow of Time
	72897,
	// Rosary of Light
	72901,
	// Varo'then's Brooch
	72899,
];

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
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
	...MISSING_ITEM_EFFECTS.map((itemID):[number, ItemNoticeData] => [
		itemID,
		{
			[Spec.SpecUnknown]: MISSING_IMPLEMENTATION_WARNING
		}
	])
]);

export const GENERIC_MISSING_SET_BONUS_NOTICE_DATA = new Map<number, string>([
	[2, 'Not yet implemented'],
	[4, 'Not yet implemented'],
]);

export const SET_BONUS_NOTICES = new Map<number, SetBonusNoticeData>([
	// Custom notices

	// Generic "not yet implemented" notices
	[928, null], // Resto Druid T11
	[933, null], // Holy Paladin T11
	[935, null], // Healing Priest T11
	[938, null], // Resto Shaman T11

	[1004, null], // Resto Druid T12
	[1009, null], // Healing Priest T12
	[1011, null], // Holy Paladin T12
	[1014, null], // Resto Shaman T12

	[1056, null], // Blood DK T13
	[1060, null], // Resto Druid T13
	[1066, null], // Healing Priest T13
	[
		1069, // Resto Shaman T13
		new Map<number, string>([[2, 'Not implemented']]),
	],
]);
