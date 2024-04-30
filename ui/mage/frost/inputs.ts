import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

// export const WaterElementalDisobeyChance = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFrostMage>({
// 	fieldName: 'waterElementalDisobeyChance',
// 	percent: true,
// 	label: 'Water Ele Disobey %',
// 	labelTooltip: 'Percent of Water Elemental actions which will fail. This represents the Water Elemental moving around or standing still instead of casting.',
// 	changeEmitter: (player: Player<Spec.SpecFrostMage>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
// 	showWhen: (player: Player<Spec.SpecFrostMage>) => player.getTalents().summonWaterElemental,
// });

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                       FROST INPUTS
		// ********************************************************
		// InputHelpers.makeRotationBooleanInput<Spec.SpecFrostMage>({
		// 	fieldName: 'useIceLance',
		// 	label: 'Use Ice Lance',
		// 	labelTooltip: 'Casts Ice Lance at the end of Fingers of Frost, after using Deep Freeze.',
		// 	showWhen: (player: Player<Spec.SpecFrostMage>) => player.getTalentTree() == 2,
		// 	changeEmitter: (player: Player<Spec.SpecFrostMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
	],
};
