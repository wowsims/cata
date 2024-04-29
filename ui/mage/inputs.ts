import * as InputHelpers from '../core/components/input_helpers';
import { MageOptions_ArmorType } from '../core/proto/mage';
import { ActionId } from '../core/proto_utils/action_id';
import { MageSpecs } from '../core/proto_utils/utils';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = <SpecType extends MageSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, MageOptions_ArmorType>({
		fieldName: 'armor',
		values: [
			{ value: MageOptions_ArmorType.NoArmor, tooltip: 'No Armor' },
			{ actionId: ActionId.fromSpellId(6117), value: MageOptions_ArmorType.MageArmor },
			{ actionId: ActionId.fromSpellId(30482), value: MageOptions_ArmorType.MoltenArmor },
			{ actionId: ActionId.fromSpellId(7302), value: MageOptions_ArmorType.FrostArmor },
		],
	});
