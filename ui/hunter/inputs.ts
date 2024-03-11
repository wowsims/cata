import * as InputHelpers from '../core/components/input_helpers';
import { HunterOptions_Ammo as Ammo } from '../core/proto/hunter';
import { ActionId } from '../core/proto_utils/action_id';
import { HunterSpecs } from '../core/proto_utils/utils';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const WeaponAmmo = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Ammo>({
		fieldName: 'ammo',
		numColumns: 2,
		values: [
			{ value: Ammo.AmmoNone, tooltip: 'No Ammo' },
			{ actionId: ActionId.fromItemId(52021), value: Ammo.IcebladeArrow },
			{ actionId: ActionId.fromItemId(41165), value: Ammo.SaroniteRazorheads },
			{ actionId: ActionId.fromItemId(41586), value: Ammo.TerrorshaftArrow },
			{ actionId: ActionId.fromItemId(31737), value: Ammo.TimelessArrow },
			{ actionId: ActionId.fromItemId(34581), value: Ammo.MysteriousArrow },
			{ actionId: ActionId.fromItemId(33803), value: Ammo.AdamantiteStinger },
			{ actionId: ActionId.fromItemId(28056), value: Ammo.BlackflightArrow },
		],
	});

export const PetTypeInput = <SpecType extends HunterSpecs>() => makePetTypeInputConfig<SpecType>();

export const PetUptime = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'petUptime',
		label: 'Pet Uptime (%)',
		labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
		percent: true,
	});

export const UseHuntersMark = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
		fieldName: 'useHuntersMark',
		id: ActionId.fromSpellId(53338),
	});

export const TimeToTrapWeaveMs = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'timeToTrapWeaveMs',
		label: 'Weave Time',
		labelTooltip:
			'Amount of time for Explosive Trap, in milliseconds, between when you start moving towards the boss and when you re-engage your ranged autos.',
	});
