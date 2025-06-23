import * as InputHelpers from '../core/components/input_helpers';
import { DeathKnightSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = <SpecType extends DeathKnightSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'startingRunicPower',
		label: 'Starting Runic Power',
		labelTooltip: 'Initial RP at the start of each iteration.',
	});
