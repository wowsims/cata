import * as InputHelpers from '../../core/components/input_helpers.js';
import { Spec } from '../../core/proto/common.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecGuardianDruid>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial Rage at the start of each iteration.',
});

export const GuardianDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationNumberInput<Spec.SpecGuardianDruid>({
			fieldName: 'maulRageThreshold',
			label: 'Maul Rage Threshold',
			labelTooltip: 'Queue Maul when Rage is above this value.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecGuardianDruid>({
			fieldName: 'pulverizeTime',
			label: 'Pulverize Refresh Leeway',
			labelTooltip: 'Refresh Pulverize when remaining duration is less than this value (in seconds).',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecGuardianDruid>({
			fieldName: 'maintainDemoralizingRoar',
			label: 'Maintain Demo Roar',
			labelTooltip: 'Keep Demoralizing Roar active on the primary target. If a stronger debuff is active, will not cast.',
		}),
	],
};
