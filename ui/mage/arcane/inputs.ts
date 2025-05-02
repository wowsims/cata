import * as InputHelpers from '../../core/components/input_helpers';
import { Spec } from '../../core/proto/common';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecArcaneMage>({
	fieldName: 'focusMagicPercentUptime',
	label: 'Focus Magic Percent Uptime',
	labelTooltip: 'Percent of uptime for Focus Magic Buddy',
	extraCssClasses: ['within-raid-sim-hide'],
});

export const MageRotationConfig = {
	inputs: [],
};
