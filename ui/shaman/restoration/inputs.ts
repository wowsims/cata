import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const TriggerEarthShield = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRestorationShaman>({
	fieldName: 'earthShieldPPM',
	label: 'Earth Shield PPM',
	labelTooltip: 'How many times Earth Shield should be triggered per minute.',
	showWhen: (player: Player<Spec.SpecRestorationShaman>) => player.getTalents().earthShield,
	changeEmitter: (player: Player<Spec.SpecRestorationShaman>) =>
		TypedEvent.onAny([player.specOptionsChangeEmitter, player.rotationChangeEmitter, player.talentsChangeEmitter]),
});
