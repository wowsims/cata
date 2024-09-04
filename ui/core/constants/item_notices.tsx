import { ItemNoticeData } from '../components/item_notice/item_notice';
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
	// Rogue Legendary Daggers (All Stages)
	[ // Fear
		77945,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
	[ // Vengeance
		77946,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
	[ // Sleeper
		77947,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
	[ // Dreamer
		77948,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
	[ // Golad
		77949,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
	[ // Tiriosh
		77950,
		{
			[Spec.SpecUnknown]: DTR_MISSING_IMPLEMENTATION_WARNING
		},
	],
]);
