import * as InputHelpers from '../../core/components/input_helpers.js';
import { Player } from '../../core/player.js';
import { Spec } from '../../core/proto/common.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const OkfUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecBalanceDruid>({
	fieldName: 'okfUptime',
	label: 'Owlkin Frenzy Uptime (%)',
	labelTooltip: 'Percentage of fight uptime for Owlkin Frenzy',
	percent: true,
});

export const StartInSolar = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecBalanceDruid>({
	fieldName: 'startInSolar',
	label: 'Start in Solar Eclipse',
	labelTooltip: 'Starts the fight in solar eclipse at full energy',
});

export const MasterySnapshot = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecBalanceDruid>({
	fieldName: 'masterySnapshot',
	label: 'Mastery snapshot amount (rating)',
	labelTooltip: 'Mastery amount to use when starting in Solar Eclipse',
	showWhen: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().startInSolar,
});
