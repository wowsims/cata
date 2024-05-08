import * as InputHelpers from '../core/components/input_helpers';
import { HunterOptions_Ammo as Ammo } from '../core/proto/hunter';
import { ActionId } from '../core/proto_utils/action_id';
import { HunterSpecs } from '../core/proto_utils/utils';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const PetTypeInput = <SpecType extends HunterSpecs>() => makePetTypeInputConfig<SpecType>();

export const PetUptime = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'petUptime',
		label: 'Pet Uptime (%)',
		labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
		percent: true,
	});
