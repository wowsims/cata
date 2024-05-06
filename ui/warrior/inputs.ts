import * as InputHelpers from '../core/components/input_helpers';
import { Spec } from '../core/proto/common';
import { WarriorShout } from '../core/proto/warrior';
import { ActionId } from '../core/proto_utils/action_id';
import { WarriorSpecs } from '../core/proto_utils/utils';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = <SpecType extends WarriorSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'startingRage',
		label: 'Starting Rage',
		labelTooltip: 'Initial rage at the start of each iteration.',
	});

export const ShoutPicker = <SpecType extends WarriorSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, WarriorShout>({
		fieldName: 'shout',
		values: [
			{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
			{ actionId: ActionId.fromSpellId(2048), value: WarriorShout.WarriorShoutBattle },
			{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
		],
	});

export const ShatteringThrow = <SpecType extends WarriorSpecs>() =>
	InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
		fieldName: 'useShatteringThrow',
		id: ActionId.fromSpellId(64382),
	});

// Arms/Fury only

export const Recklessness = <SpecType extends Spec.SpecArmsWarrior | Spec.SpecFuryWarrior | Spec.SpecProtectionWarrior>() =>
	InputHelpers.makeSpecOptionsBooleanIconInput<SpecType>({
		fieldName: 'useRecklessness',
		id: ActionId.fromSpellId(1719),
	});

export const StanceSnapshot = <SpecType extends Spec.SpecArmsWarrior | Spec.SpecFuryWarrior>() =>
	InputHelpers.makeSpecOptionsBooleanInput<SpecType>({
		fieldName: 'stanceSnapshot',
		label: 'Stance Snapshot',
		labelTooltip: 'Ability that is cast at the same time as stance swap will benefit from the bonus of the stance before the swap.',
	});

// Allows for auto gemming whilst ignoring expertise cap
// (Useful for Arms)
export const DisableExpertiseGemming = <SpecType extends Spec.SpecArmsWarrior | Spec.SpecFuryWarrior>() =>
	InputHelpers.makeSpecOptionsBooleanInput<SpecType>({
		fieldName: 'disableExpertiseGemming',
		label: 'Disable expertise gemming',
		labelTooltip: 'Disables auto gemming for expertise',
	});
