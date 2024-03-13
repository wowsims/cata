import * as InputHelpers from '../../core/components/input_helpers.js';
import { Player } from '../../core/player.js';
import { RotationType, Spec } from '../../core/proto/common.js';
import { HunterStingType } from '../../core/proto/hunter';
import { TypedEvent } from '../../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SniperTrainingUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecSurvivalHunter>({
	fieldName: 'sniperTrainingUptime',
	label: 'ST Uptime (%)',
	labelTooltip: 'Uptime for the Sniper Training talent, as a percent of the fight duration.',
	percent: true,
	showWhen: (player: Player<Spec.SpecSurvivalHunter>) => player.getTalents().sniperTraining > 0,
	changeEmitter: (player: Player<Spec.SpecSurvivalHunter>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const SVRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecSurvivalHunter, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: RotationType.SingleTarget },
				{ name: 'AOE', value: RotationType.Aoe },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecSurvivalHunter, HunterStingType>({
			fieldName: 'sting',
			label: 'Sting',
			labelTooltip: 'Maintains the selected Sting on the primary target.',
			values: [
				{ name: 'None', value: HunterStingType.NoSting },
				{ name: 'Scorpid Sting', value: HunterStingType.ScorpidSting },
				{ name: 'Serpent Sting', value: HunterStingType.SerpentSting },
			],
			showWhen: (player: Player<Spec.SpecSurvivalHunter>) => player.getSimpleRotation().type == RotationType.SingleTarget,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSurvivalHunter>({
			fieldName: 'trapWeave',
			label: 'Trap Weave',
			labelTooltip: 'Uses Explosive Trap at appropriate times. Note that selecting this will disable Black Arrow because they share a CD.',
		}),
		// InputHelpers.makeRotationBooleanInput<Spec.SpecSurvivalHunter>({
		// 	fieldName: 'allowExplosiveShotDownrank',
		// 	label: 'Allow ES Downrank',
		// 	labelTooltip: 'Weaves Explosive Shot Rank 3 during LNL procs. This works because the rank 3 and rank 4 dots can stack.',
		// 	showWhen: (player: Player<Spec.SpecSurvivalHunter>) =>
		// 		player.getSimpleRotation().type != RotationType.Custom && player.getTalents().explosiveShot && player.getTalents().lockAndLoad > 0,
		// 	changeEmitter: (player: Player<Spec.SpecSurvivalHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSurvivalHunter>({
			fieldName: 'multiDotSerpentSting',
			label: 'Multi-Dot Serpent Sting',
			labelTooltip: 'Casts Serpent Sting on multiple targets',
			changeEmitter: (player: Player<Spec.SpecSurvivalHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecSurvivalHunter>({
			fieldName: 'viperStartManaPercent',
			label: 'Viper Start Mana %',
			labelTooltip: 'Switch to Aspect of the Viper when mana goes below this amount.',
			percent: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecSurvivalHunter>({
			fieldName: 'viperStopManaPercent',
			label: 'Viper Stop Mana %',
			labelTooltip: 'Switch back to Aspect of the Hawk when mana goes above this amount.',
			percent: true,
		}),
	],
};
