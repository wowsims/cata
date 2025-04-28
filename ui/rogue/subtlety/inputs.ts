import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const HonorAmongThievesCritRate = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecSubtletyRogue>({
	fieldName: 'honorAmongThievesCritRate',
	label: 'Honor of Thieves Crit Rate',
	labelTooltip: 'Number of crits other group members generate within 100 seconds',
	showWhen: (player: Player<Spec.SpecSubtletyRogue>) => false,
});
