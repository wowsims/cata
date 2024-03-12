import * as InputHelpers from '../../core/components/input_helpers';
import { Player } from '../../core/player';
import { RotationType, Spec } from '../../core/proto/common';
import { HunterStingType } from '../../core/proto/hunter';
import { TypedEvent } from '../../core/typed_event';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MMRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecMarksmanshipHunter, RotationType>({
			fieldName: 'type',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: RotationType.SingleTarget },
				{ name: 'AOE', value: RotationType.Aoe },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecMarksmanshipHunter, HunterStingType>({
			fieldName: 'sting',
			label: 'Sting',
			labelTooltip: 'Maintains the selected Sting on the primary target.',
			values: [
				{ name: 'None', value: HunterStingType.NoSting },
				{ name: 'Scorpid Sting', value: HunterStingType.ScorpidSting },
				{ name: 'Serpent Sting', value: HunterStingType.SerpentSting },
			],
			showWhen: (player: Player<Spec.SpecMarksmanshipHunter>) => player.getSimpleRotation().type == RotationType.SingleTarget,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMarksmanshipHunter>({
			fieldName: 'trapWeave',
			label: 'Trap Weave',
			labelTooltip: 'Uses Explosive Trap at appropriate times. Note that selecting this will disable Black Arrow because they share a CD.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMarksmanshipHunter>({
			fieldName: 'multiDotSerpentSting',
			label: 'Multi-Dot Serpent Sting',
			labelTooltip: 'Casts Serpent Sting on multiple targets',
			changeEmitter: (player: Player<Spec.SpecMarksmanshipHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMarksmanshipHunter>({
			fieldName: 'viperStartManaPercent',
			label: 'Viper Start Mana %',
			labelTooltip: 'Switch to Aspect of the Viper when mana goes below this amount.',
			percent: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMarksmanshipHunter>({
			fieldName: 'viperStopManaPercent',
			label: 'Viper Stop Mana %',
			labelTooltip: 'Switch back to Aspect of the Hawk when mana goes above this amount.',
			percent: true,
		}),
	],
};
