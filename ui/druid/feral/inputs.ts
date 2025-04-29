import * as InputHelpers from '../../core/components/input_helpers.js';
import { Player } from '../../core/player.js';
import { APLRotation_Type } from '../../core/proto/apl.js';
import { Spec } from '../../core/proto/common.js';
import { FeralDruid_Rotation_AplType as AplType, FeralDruid_Rotation_BiteModeType as BiteModeType } from '../../core/proto/druid.js';
import { TypedEvent } from '../../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: "Assume bleed always exists for 'Rend and Tear' calculations. Otherwise will only calculate based on own rip/rake/lacerate.",
	extraCssClasses: ['within-raid-sim-hide'],
});

export const CannotShredTarget = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'cannotShredTarget',
	label: 'Cannot Shred Target',
	labelTooltip: 'Alternative to "In Front of Target" for modeling bosses that do not Parry or Block, but which you still cannot Shred.',
});

function ShouldShowAdvParamST(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.SingleTarget;
}

function ShouldShowAdvParamAoe(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.Aoe;
}

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, AplType>({
			fieldName: 'rotationType',
			label: 'Type',
			values: [
				{ name: 'Single Target', value: AplType.SingleTarget },
				{ name: 'AOE', value: AplType.Aoe },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'Maintain Faerie Fire debuff. Overwrites any external Sunder effects specified in settings.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'meleeWeave',
			label: 'Enable leave-weaving',
			labelTooltip: 'Weave out of melee range for Stampede procs',
			// showWhen: (player: Player<Spec.SpecFeralDruid>) =>
			// 	player.getSimpleRotation().rotationType == AplType.SingleTarget &&
			// 	player.getTalents().stampede > 0 &&
			// 	!player.getSpecOptions().cannotShredTarget &&
			// 	!player.getInFrontOfTarget(),
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) =>
				TypedEvent.onAny([
					player.rotationChangeEmitter,
					player.talentsChangeEmitter,
					player.specOptionsChangeEmitter,
					player.inFrontOfTargetChangeEmitter,
				]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'bearWeave',
			label: 'Enable bear-weaving',
			labelTooltip: 'Weave into Bear Form while pooling Energy',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'snekWeave',
			label: 'Use Albino Snake',
			labelTooltip: 'Reset swing timer at the end of bear-weaves using Albino Snake pet',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().bearWeave,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'allowAoeBerserk',
			label: 'Allow AoE Berserk',
			labelTooltip: 'Allow Berserk usage in AoE rotation',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().rotationType == AplType.Aoe,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'prepullTranquility',
			label: 'Enable pre-pull Tranquility',
			labelTooltip: 'Swap in configured healing trinkets before the pull and proc them using Tranquility',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.shouldEnableTargetDummies(),
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.itemSwapSettings.changeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'manualParams',
			label: 'Manual Advanced Parameters',
			labelTooltip: 'Manually specify advanced parameters, otherwise will use preset defaults',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().rotationType == AplType.SingleTarget,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'minRoarOffset',
			label: 'Roar Offset',
			labelTooltip: 'Targeted offset in Rip/Roar timings',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'ripLeeway',
			label: 'Rip Leeway',
			labelTooltip: 'Rip leeway when determining roar clips',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useRake',
			label: 'Use Rake',
			labelTooltip: 'Use rake during rotation',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: 'Bite during rotation',
			labelTooltip: 'Use bite during rotation rather than just at end',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: 'Bite Time',
			labelTooltip: 'Min seconds on Rip/Roar to bite',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Emperical,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'berserkBiteTime',
			label: 'Bite Time during Berserk',
			labelTooltip: 'More aggressive threshold when Berserk is active',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Emperical,
		}),
		// InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
		// 	fieldName: 'biteDuringExecute',
		// 	label: 'Bite during Execute phase',
		// 	labelTooltip: 'Bite aggressively during Execute phase',
		// 	showWhen: (player: Player<Spec.SpecFeralDruid>) =>
		// 		ShouldShowAdvParamST(player) && player.getTalents().bloodInTheWater > 0,
		// 	changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
		// InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
		// 	fieldName: 'cancelPrimalMadness',
		// 	label: 'Enable Primal Madness cancellation',
		// 	labelTooltip: 'Click off Primal Madness buff when doing so will result in net Energy gains',
		// 	showWhen: (player: Player<Spec.SpecFeralDruid>) =>
		// 		(ShouldShowAdvParamST(player) || (player.getSimpleRotation().rotationType == AplType.Aoe)) && (player.getTalents().primalMadness > 0),
		// 	changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		// }),
		// Can be uncommented if/when analytical bite mode is added
		//InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BiteModeType>({
		//	fieldName: 'biteModeType',
		//	label: 'Bite Mode',
		//	labelTooltip: 'Underlying "Bite logic" to use',
		//	values: [
		//		{ name: 'Emperical', value: BiteModeType.Emperical },
		//	],
		//	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().useBite == true
		//}),
	],
};
