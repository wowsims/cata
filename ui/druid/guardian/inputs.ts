import * as InputHelpers from '../../core/components/input_helpers.js';
import { Player } from '../../core/player.js';
import { Spec } from '../../core/proto/common.js';
import { TypedEvent } from '../../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecGuardianDruid>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial Rage at the start of each iteration.',
});

export const GuardianDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecGuardianDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'Maintain Faerie Fire debuff. Overwrites any external Sunder effects specified in settings.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecGuardianDruid>({
			fieldName: 'maintainDemoralizingRoar',
			label: 'Maintain Demo Roar',
			labelTooltip: 'Keep Demoralizing Roar active on the primary target. If a stronger debuff is active, will not cast.',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecGuardianDruid>({
			fieldName: 'demoTime',
			label: 'Demo Roar refresh leeway',
			labelTooltip:
				'Refresh Demoralizing Roar when remaining duration is less than this value (in seconds). Larger values provide safety buffer against misses, but at the cost of lower DPS.',
			showWhen: (player: Player<Spec.SpecGuardianDruid>) => player.getSimpleRotation().maintainDemoralizingRoar,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecGuardianDruid>({
			fieldName: 'pulverizeTime',
			label: 'Pulverize refresh leeway',
			labelTooltip:
				'Refresh Pulverize when remaining duration is less than this value (in seconds). Note that Mangle, Thrash, and Faerie Fire usage on cooldown takes priority over this rule, unless Lacerate itself is about to fall off.',
		}),
		// InputHelpers.makeRotationBooleanInput<Spec.SpecGuardianDruid>({
		// 	fieldName: 'prepullStampede',
		// 	label: 'Assume pre-pull Stampede',
		// 	labelTooltip: 'Activate Stampede Haste buff at the start of each pull. Models the effects of initiating the pull with Feral Charge.',
		// 	showWhen: (player: Player<Spec.SpecGuardianDruid>) =>
		// 		player.getTalents().stampede > 0,
		// 	changeEmitter: (player: Player<Spec.SpecGuardianDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
	],
};
