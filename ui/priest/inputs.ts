import * as InputHelpers from '../core/components/input_helpers';
import { PriestOptions_Armor } from '../core/proto/priest';
import { ActionId } from '../core/proto_utils/action_id';
import { PriestSpecs } from '../core/proto_utils/utils';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = <SpecType extends PriestSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, PriestOptions_Armor>({
		fieldName: 'armor',
		values: [
			{ value: PriestOptions_Armor.NoArmor, tooltip: 'No Inner Fire' },
			{ actionId: ActionId.fromSpellId(48168), value: PriestOptions_Armor.InnerFire },
		],
	});
