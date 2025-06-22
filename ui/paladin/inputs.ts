import * as InputHelpers from '../core/components/input_helpers';
import { Player } from '../core/player';
import { Spec } from '../core/proto/common';
import { PaladinSeal } from '../core/proto/paladin';
import { ActionId } from '../core/proto_utils/action_id';
import { PaladinSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingSealSelection = <SpecType extends PaladinSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, PaladinSeal>({
		fieldName: 'seal',
		values: [
			{ actionId: ActionId.fromSpellId(31801), value: PaladinSeal.Truth },
			{ actionId: ActionId.fromSpellId(20154), value: PaladinSeal.Righteousness },
			{ actionId: ActionId.fromSpellId(20165), value: PaladinSeal.Insight },
			{
				actionId: ActionId.fromSpellId(20164),
				value: PaladinSeal.Justice,
				showWhen: player => player.isSpec(Spec.SpecRetributionPaladin),
			},
		],
		changeEmitter: (player: Player<SpecType>) => player.changeEmitter,
	});

export const StartingHolyPower = <SpecType extends PaladinSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'startingHolyPower',
		label: 'Starting Holy Power',
		labelTooltip: "Initial Holy Power at the start of each iteration.",
	});
