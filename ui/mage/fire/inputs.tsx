import * as InputHelpers from '../../core/components/input_helpers';
import { Spec } from '../../core/proto/common';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                       FIRE INPUTS
		// ********************************************************
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'combustThreshold',
			label: 'Combust Threshold - Bloodlust',
			labelTooltip: 'The value at which Combustion should be cast during Bloodlust',
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().combustThreshold,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'combustLastMomentLustPercentage',
			label: 'Combust Threshold - Last moment during Bloodlust',
			labelTooltip: 'The value at which Combustion should be cast when Bloodlust (+ Berserking) is about to run out.',
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().combustLastMomentLustPercentage,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'combustNoLustPercentage',
			label: 'Combust Threshold - Outside of Bloodlust',
			labelTooltip: 'The value at which Combustion should be cast when Bloodlust is not up.',
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().combustNoLustPercentage,
			positive: true,
		}),
	],
};
