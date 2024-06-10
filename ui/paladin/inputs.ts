import * as InputHelpers from '../core/components/input_helpers';
import { Player } from '../core/player';
import { PaladinAura, PaladinSeal } from '../core/proto/paladin';
import { ActionId } from '../core/proto_utils/action_id';
import { PaladinSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AuraSelection = <SpecType extends PaladinSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, PaladinAura>({
		fieldName: 'aura',
		values: [
			{ actionId: ActionId.fromSpellId(465), value: PaladinAura.Devotion },
			{ actionId: ActionId.fromSpellId(7294), value: PaladinAura.Retribution },
			{ actionId: ActionId.fromSpellId(19891), value: PaladinAura.Resistance },
		],
	});

export const StartingSealSelection = <SpecType extends PaladinSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, PaladinSeal>({
		fieldName: 'seal',
		values: [
			{ actionId: ActionId.fromSpellId(31801), value: PaladinSeal.Truth },
			{ actionId: ActionId.fromSpellId(20154), value: PaladinSeal.Righteousness },
			{ actionId: ActionId.fromSpellId(20165), value: PaladinSeal.Insight },
			{ actionId: ActionId.fromSpellId(20164), value: PaladinSeal.Justice },
		],
		changeEmitter: (player: Player<SpecType>) => player.changeEmitter,
	});
