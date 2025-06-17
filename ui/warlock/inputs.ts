import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { WarlockOptions_Summon as Summon } from '../core/proto/warlock.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { WarlockSpecs } from '../core/proto_utils/utils';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const PetInput = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, Summon>({
		fieldName: 'summon',
		values: [
			{ value: Summon.NoSummon, tooltip: 'No Pet' },
			{ actionId: ActionId.fromSpellId(691), value: Summon.Felhunter },
			{
				actionId: ActionId.fromSpellId(30146),
				value: Summon.Felguard,
				showWhen: (player: Player<SpecType>) => player.getSpec() == Spec.SpecDemonologyWarlock,
			},
			{ actionId: ActionId.fromSpellId(688), value: Summon.Imp },
			{ actionId: ActionId.fromSpellId(712), value: Summon.Succubus },
			{ actionId: ActionId.fromSpellId(697), value: Summon.Voidwalker },
		],
		changeEmitter: (player: Player<SpecType>) => player.changeEmitter,
	});

export const DetonateSeed = <SpecType extends WarlockSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'detonateSeed',
		label: 'Detonate Seed on Cast',
		labelTooltip: 'Simulates raid doing damage to targets such that seed detonates immediately on cast.',
	});

