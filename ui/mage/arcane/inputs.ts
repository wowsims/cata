import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { Spec } from '../../core/proto/common';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecArcaneMage>({
	fieldName: 'focusMagicPercentUptime',
	label: 'Focus Magic Percent Uptime',
	labelTooltip: 'Percent of uptime for Focus Magic Buddy',
	extraCssClasses: ['within-raid-sim-hide'],
});

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                      ARCANE INPUTS
		// ********************************************************
		// InputHelpers.makeRotationNumberInput<Spec.SpecArcaneMage>({
		// 	fieldName: 'only3ArcaneBlastStacksBelowManaPercent',
		// 	percent: true,
		// 	label: 'Stack Arcane Blast to 3 below mana %',
		// 	labelTooltip: 'When below this mana %, AM/ABarr will be used at 3 stacks of AB instead of 4.',
		// 	showWhen: (player: Player<Spec.SpecArcaneMage>) => player.getTalentTree() == 0,
		// 	changeEmitter: (player: Player<Spec.SpecArcaneMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
		// InputHelpers.makeRotationNumberInput<Spec.SpecArcaneMage>({
		// 	fieldName: 'blastWithoutMissileBarrageAboveManaPercent',
		// 	percent: true,
		// 	label: 'AB without Missile Barrage above mana %',
		// 	labelTooltip: 'When above this mana %, spam AB until a Missile Barrage proc occurs.',
		// 	showWhen: (player: Player<Spec.SpecArcaneMage>) => player.getTalentTree() == 0,
		// 	changeEmitter: (player: Player<Spec.SpecArcaneMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
		// InputHelpers.makeRotationNumberInput<Spec.SpecArcaneMage>({
		// 	fieldName: 'missileBarrageBelowManaPercent',
		// 	percent: true,
		// 	label: 'Use Missile Barrage ASAP below mana %',
		// 	labelTooltip: 'When below this mana %, use Missile Barrage proc as soon as possible. Can be useful to conserve mana.',
		// 	showWhen: (player: Player<Spec.SpecArcaneMage>) => player.getTalentTree() == 0,
		// 	changeEmitter: (player: Player<Spec.SpecArcaneMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
	// 	InputHelpers.makeRotationBooleanInput<Spec.SpecArcaneMage>({
	// 		fieldName: 'useArcaneBarrage',
	// 		label: 'Use Arcane Barrage',
	// 		labelTooltip: 'Includes Arcane Barrage in the rotation.',
	// 		enableWhen: (player: Player<Spec.SpecArcaneMage>) => player.getTalents().arcaneBarrage,
	// 		showWhen: (player: Player<Spec.SpecArcaneMage>) => player.getTalentTree() == 0,
	// 		changeEmitter: (player: Player<Spec.SpecArcaneMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
	// 	}),
	],
};
