import * as InputHelpers from '../core/components/input_helpers';
import { Player } from '../core/player';
import { UnitReference, UnitReference_Type as UnitType } from '../core/proto/common';
import { ActionId } from '../core/proto_utils/action_id';
import { DruidSpecs } from '../core/proto_utils/utils';
import { EventID } from '../core/typed_event';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = <SpecType extends DruidSpecs>() =>
	InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
		fieldName: 'innervateTarget',
		id: ActionId.fromSpellId(29166),
		extraCssClasses: ['within-raid-sim-hide'],
		getValue: (player: Player<SpecType>) => player.getClassOptions().innervateTarget?.type == UnitType.Player,
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => {
			const newOptions = player.getClassOptions();
			newOptions.innervateTarget = UnitReference.create({
				type: newValue ? UnitType.Player : UnitType.Unknown,
				index: 0,
			});
			player.setClassOptions(eventID, newOptions);
		},
	});
