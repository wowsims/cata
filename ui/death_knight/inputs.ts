import * as InputHelpers from '../core/components/input_helpers';
import { Player } from '../core/player';
import { DeathKnightSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRunicPower = <SpecType extends DeathKnightSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'startingRunicPower',
		label: 'Starting Runic Power',
		labelTooltip: 'Initial RP at the start of each iteration.',
	});

// export const PetUptime = <SpecType extends DeathKnightSpecs>() =>
// 	InputHelpers.makeClassOptionsNumberInput<SpecType>({
// 		fieldName: 'petUptime',
// 		label: 'Ghoul Uptime (%)',
// 		labelTooltip: 'Percent of the fight duration for which your ghoul will be on target.',
// 		percent: true,
// 		showWhen: (player: Player<SpecType>) => player.getTalents().masterOfGhouls,
// 	});
