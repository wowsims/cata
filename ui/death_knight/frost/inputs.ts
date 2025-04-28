import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

// export const UseAMSInput = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFrostDeathKnight>({
// 	fieldName: 'useAms',
// 	label: 'Use AMS',
// 	labelTooltip: 'Use AMS around predicted damage for a RP gain.',
// 	showWhen: (player: Player<Spec.SpecFrostDeathKnight>) => player.getTalents().howlingBlast,
// 	changeEmitter: (player: Player<Spec.SpecFrostDeathKnight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
// });

// export const AvgAMSSuccessRateInput = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFrostDeathKnight>({
// 	fieldName: 'avgAmsSuccessRate',
// 	label: 'Avg AMS Success %',
// 	labelTooltip: 'Chance for damage to be taken during the 5 second window of AMS.',
// 	showWhen: (player: Player<Spec.SpecFrostDeathKnight>) => player.getSpecOptions().useAms == true && player.getTalents().howlingBlast,
// 	changeEmitter: (player: Player<Spec.SpecFrostDeathKnight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
// });

// export const AvgAMSHitInput = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFrostDeathKnight>({
// 	fieldName: 'avgAmsHit',
// 	label: 'Avg AMS Hit',
// 	labelTooltip: 'How much on average (+-10%) the character is hit for when AMS is successful.',
// 	showWhen: (player: Player<Spec.SpecFrostDeathKnight>) => player.getSpecOptions().useAms == true && player.getTalents().howlingBlast,
// 	changeEmitter: (player: Player<Spec.SpecFrostDeathKnight>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
// });
