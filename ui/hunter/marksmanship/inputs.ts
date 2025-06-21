import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { RotationType, Spec } from '../../core/proto/common';
import { HunterStingType } from '../../core/proto/hunter';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MMRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecMarksmanshipHunter, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: RotationType.SingleTarget },
				{ name: 'AOE', value: RotationType.Aoe },
			],
		}),
	],
};
