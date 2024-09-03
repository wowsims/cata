import { ItemNoticeData } from '../components/item_notice/item_notice';
import { Spec } from '../proto/common';

const DTR_PARTIAL_IMPLEMENTATION_WARNING = (
	<>
		<p className="fw-bold">This item is not implemented!</p>
		<p>
			We are working hard on gathering all the old resources to allow for an initial implementation.
		</p>
		<p className="mb-0">
			Want to help out by providing additional information? Contact us on our Discord!
		</p>
	</>
);

export const ITEM_NOTICES = new Map<number, ItemNoticeData>([
	// Dragonwrath, Tarecgosa's Rest
	[
		71086,
		{
			[Spec.SpecUnknown]: (
				<>
					<p>This item is unsupported for this spec.</p>
				</>
			),
			[Spec.SpecBalanceDruid]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecArcaneMage]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecFireMage]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecFrostMage]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecShadowPriest]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecElementalShaman]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecAfflictionWarlock]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecDemonologyWarlock]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
			[Spec.SpecDestructionWarlock]: DTR_PARTIAL_IMPLEMENTATION_WARNING,
		},
	],
]);
