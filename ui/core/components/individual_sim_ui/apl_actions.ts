import { itemSwapEnabledSpecs } from '../../individual_sim_ui.js';
import { Player } from '../../player.js';
import {
	APLAction,
	APLActionActivateAllStatBuffProcAuras,
	APLActionActivateAura,
	APLActionActivateAuraWithStacks,
	APLActionAutocastOtherCooldowns,
	APLActionCancelAura,
	APLActionCastAllStatBuffCooldowns,
	APLActionCastFriendlySpell,
	APLActionCastSpell,
	APLActionCatOptimalRotationAction,
	APLActionChangeTarget,
	APLActionChannelSpell,
	APLActionCustomRotation,
	APLActionItemSwap,
	APLActionItemSwap_SwapSet as ItemSwapSet,
	APLActionMove,
	APLActionMoveDuration,
	APLActionMultidot,
	APLActionMultishield,
	APLActionResetSequence,
	APLActionSchedule,
	APLActionSequence,
	APLActionStrictMultidot,
	APLActionStrictSequence,
	APLActionTriggerICD,
	APLActionWait,
	APLActionWaitUntil,
	APLValue,
} from '../../proto/apl.js';
import { Spec } from '../../proto/common.js';
import { FeralDruid_Rotation_AplType } from '../../proto/druid.js';
import { EventID } from '../../typed_event.js';
import { randomUUID } from '../../utils';
import { Input, InputConfig } from '../input.js';
import { TextDropdownPicker } from '../pickers/dropdown_picker.jsx';
import { ListItemPickerConfig, ListPicker } from '../pickers/list_picker.jsx';
import * as AplHelpers from './apl_helpers.js';
import * as AplValues from './apl_values.js';

export interface APLActionPickerConfig extends InputConfig<Player<any>, APLAction> {}

export type APLActionKind = APLAction['action']['oneofKind'];
type APLActionImplStruct<F extends APLActionKind> = Extract<APLAction['action'], { oneofKind: F }>;
type APLActionImplTypesUnion = {
	[f in NonNullable<APLActionKind>]: f extends keyof APLActionImplStruct<f> ? APLActionImplStruct<f>[f] : never;
};
export type APLActionImplType = APLActionImplTypesUnion[NonNullable<APLActionKind>] | undefined;

export class APLActionPicker extends Input<Player<any>, APLAction> {
	private kindPicker: TextDropdownPicker<Player<any>, APLActionKind>;

	private readonly actionDiv: HTMLElement;
	private currentKind: APLActionKind;
	private actionPicker: Input<Player<any>, any> | null;

	private readonly conditionPicker: AplValues.APLValuePicker;

	constructor(parent: HTMLElement, player: Player<any>, config: APLActionPickerConfig) {
		super(parent, 'apl-action-picker-root', player, config);
		this.conditionPicker = new AplValues.APLValuePicker(this.rootElem, this.modObject, {
			label: 'If:',
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (_player: Player<any>) => this.getSourceValue()?.condition,
			setValue: (eventID: EventID, player: Player<any>, newValue: APLValue | undefined) => {
				const srcVal = this.getSourceValue();
				if (srcVal) {
					srcVal.condition = newValue;
					player.rotationChangeEmitter.emit(eventID);
				} else {
					this.setSourceValue(
						eventID,
						APLAction.create({
							condition: newValue,
						}),
					);
				}
			},
		});
		this.conditionPicker.rootElem.classList.add('apl-action-condition', 'apl-priority-list-only');

		this.actionDiv = document.createElement('div');
		this.actionDiv.classList.add('apl-action-picker-action');
		this.rootElem.appendChild(this.actionDiv);

		const isPrepull = this.rootElem.closest('.apl-prepull-action-picker') != null;

		const allActionKinds = (Object.keys(actionKindFactories) as Array<NonNullable<APLActionKind>>).filter(
			actionKind => actionKindFactories[actionKind].includeIf?.(player, isPrepull) ?? true,
		);

		this.kindPicker = new TextDropdownPicker(this.actionDiv, player, {
			id: randomUUID(),
			defaultLabel: 'Action',
			values: allActionKinds.map(actionKind => {
				const factory = actionKindFactories[actionKind];
				return {
					value: actionKind,
					label: factory.label,
					submenu: factory.submenu,
					tooltip: factory.fullDescription ? `<p>${factory.shortDescription}</p> ${factory.fullDescription}` : factory.shortDescription,
				};
			}),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (_player: Player<any>) => this.getSourceValue()?.action.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newKind: APLActionKind) => {
				const sourceValue = this.getSourceValue();
				const oldKind = sourceValue?.action.oneofKind;
				if (oldKind == newKind) {
					return;
				}

				if (newKind) {
					const factory = actionKindFactories[newKind];
					let newSourceValue = this.makeAPLAction(newKind, factory.newValue());
					if (sourceValue) {
						// Some pre-fill logic when swapping kinds.
						if (oldKind && this.actionPicker) {
							if (newKind == 'sequence') {
								if (sourceValue.action.oneofKind == 'strictSequence') {
									(newSourceValue.action as APLActionImplStruct<'sequence'>).sequence.actions = sourceValue.action.strictSequence.actions;
								} else {
									(newSourceValue.action as APLActionImplStruct<'sequence'>).sequence.actions = [
										this.makeAPLAction(oldKind, this.actionPicker.getInputValue()),
									];
								}
							} else if (newKind == 'strictSequence') {
								if (sourceValue.action.oneofKind == 'sequence') {
									(newSourceValue.action as APLActionImplStruct<'strictSequence'>).strictSequence.actions =
										sourceValue.action.sequence.actions;
								} else {
									(newSourceValue.action as APLActionImplStruct<'strictSequence'>).strictSequence.actions = [
										this.makeAPLAction(oldKind, this.actionPicker.getInputValue()),
									];
								}
							} else if (sourceValue.action.oneofKind == 'sequence' && sourceValue.action.sequence.actions?.[0]?.action.oneofKind == newKind) {
								newSourceValue = sourceValue.action.sequence.actions[0];
							} else if (
								sourceValue.action.oneofKind == 'strictSequence' &&
								sourceValue.action.strictSequence.actions?.[0]?.action.oneofKind == newKind
							) {
								newSourceValue = sourceValue.action.strictSequence.actions[0];
							}
						}
					}
					if (sourceValue) {
						sourceValue.action = newSourceValue.action;
					} else {
						this.setSourceValue(eventID, newSourceValue);
					}
				} else {
					sourceValue.action = {
						oneofKind: newKind,
					};
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentKind = undefined;
		this.actionPicker = null;

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

	getInputValue(): APLAction {
		const actionKind = this.kindPicker.getInputValue();
		return APLAction.create({
			condition: this.conditionPicker.getInputValue(),
			action: {
				oneofKind: actionKind,
				...(() => {
					const val: any = {};
					if (actionKind && this.actionPicker) {
						val[actionKind] = this.actionPicker.getInputValue();
					}
					return val;
				})(),
			},
		});
	}

	setInputValue(newValue: APLAction) {
		if (!newValue) {
			return;
		}

		this.conditionPicker.setInputValue(newValue.condition || APLValue.create({
			uuid: { value: randomUUID() }
		}));

		const newActionKind = newValue.action.oneofKind;
		this.updateActionPicker(newActionKind);

		if (newActionKind) {
			this.actionPicker!.setInputValue((newValue.action as any)[newActionKind]);
		}
	}

	private makeAPLAction<K extends NonNullable<APLActionKind>>(kind: K, implVal: APLActionImplTypesUnion[K]): APLAction {
		if (!kind) {
			return APLAction.create();
		}
		const obj: any = { oneofKind: kind };
		obj[kind] = implVal;
		return APLAction.create({ action: obj });
	}

	private updateActionPicker(newActionKind: APLActionKind) {
		const actionKind = this.currentKind;
		if (newActionKind == actionKind) {
			return;
		}
		this.currentKind = newActionKind;

		if (this.actionPicker) {
			this.actionPicker.rootElem.remove();
			this.actionPicker = null;
		}

		if (!newActionKind) {
			return;
		}

		this.kindPicker.setInputValue(newActionKind);

		const factory = actionKindFactories[newActionKind];
		this.actionPicker = factory.factory(this.actionDiv, this.modObject, {
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => (this.getSourceValue()?.action as any)?.[newActionKind] || factory.newValue(),
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				const sourceValue = this.getSourceValue();
				if (sourceValue) {
					(sourceValue?.action as any)[newActionKind] = newValue;
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});
		this.actionPicker.rootElem.classList.add('apl-action-' + newActionKind);
	}
}

type ActionKindConfig<T> = {
	label: string;
	submenu?: Array<string>;
	shortDescription: string;
	fullDescription?: string;
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean;
	newValue: () => T;
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>;
};

function itemSwapSetFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ItemSwapSet.Swap1,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				id: randomUUID(),
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: ItemSwapSet.Main, label: 'Main' },
					{ value: ItemSwapSet.Swap1, label: 'Swapped' },
				],
			}),
	};
}

function actionFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => APLValue.create({
			uuid: { value: randomUUID() }
		}) ,
		factory: (parent, player, config) => new APLActionPicker(parent, player, config),
	};
}

function actionListFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => [],
		factory: (parent, player, config) =>
			new ListPicker<Player<any>, APLAction>(parent, player, {
				...config,
				// Override setValue to replace undefined elements with default messages.
				setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLAction>) => {
					config.setValue(
						eventID,
						player,
						newValue.map(val => val || APLAction.create()),
					);
				},
				itemLabel: 'Action',
				newItem: APLAction.create,
				copyItem: (oldValue: APLAction) => (oldValue ? APLAction.clone(oldValue) : oldValue),
				newItemPicker: (
					parent: HTMLElement,
					listPicker: ListPicker<Player<any>, APLAction>,
					index: number,
					config: ListItemPickerConfig<Player<any>, APLAction>,
				) => new APLActionPicker(parent, player, config),
				allowedActions: ['create', 'delete', 'move'],
				actions: {
					create: {
						useIcon: true,
					},
				},
			}),
	};
}

function inputBuilder<T>(config: {
	label: string;
	submenu?: Array<string>;
	shortDescription: string;
	fullDescription?: string;
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean;
	newValue: () => T;
	fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, any>>;
}): ActionKindConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		includeIf: config.includeIf,
		newValue: config.newValue,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
	};
}

const actionKindFactories: { [f in NonNullable<APLActionKind>]: ActionKindConfig<APLActionImplTypesUnion[f]> } = {
	['castSpell']: inputBuilder({
		label: 'Cast',
		shortDescription: 'Casts the spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		newValue: APLActionCastSpell.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', ''), AplHelpers.unitFieldConfig('target', 'targets')],
	}),
	['castFriendlySpell']: inputBuilder({
		label: 'Cast at Player',
		shortDescription: 'Casts a friendly spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		newValue: APLActionCastFriendlySpell.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'friendly_spells', ''), AplHelpers.unitFieldConfig('target', 'players')],
		includeIf: (player: Player<any>, _isPrepull: boolean) => (player.getRaid()!.size() > 1) || player.shouldEnableTargetDummies(),
	}),
	['multidot']: inputBuilder({
		label: 'Multi Dot',
		submenu: ['Casting'],
		shortDescription: 'Keeps a DoT active on multiple targets by casting the specified spell.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () =>
			APLActionMultidot.create({
				maxDots: 3,
				maxOverlap: {
					value: {
						oneofKind: 'const',
						const: {
							val: '0ms',
						},
					},
				},
			}),
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_dot_spells', ''),
			AplHelpers.numberFieldConfig('maxDots', false, {
				label: 'Max Dots',
				labelTooltip: 'Maximum number of DoTs to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a DoT expires when it may be refreshed.',
			}),
		],
	}),
	['strictMultidot']: inputBuilder({
		label: 'Strict Multi Dot',
		submenu: ['Casting'],
		shortDescription: 'Like a regular <b>Multi Dot</b>, except all Dots are applied immediately after each other. Keeps a DoT active on multiple targets by casting the specified spell. Will take Cast Time/GCD into account when refreshing subsequent DoTs.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () =>
			APLActionStrictMultidot.create({
				maxDots: 3,
				maxOverlap: {
					value: {
						oneofKind: 'const',
						const: {
							val: '0ms',
						},
					},
				},
			}),
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'castable_dot_spells', ''),
			AplHelpers.numberFieldConfig('maxDots', false, {
				label: 'Max Dots',
				labelTooltip: 'Maximum number of DoTs to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a DoT expires when it may be refreshed.',
			}),
		],
	}),
	['multishield']: inputBuilder({
		label: 'Multi Shield',
		submenu: ['Casting'],
		shortDescription: 'Keeps a Shield active on multiple targets by casting the specified spell.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull && player.getSpec().isHealingSpec,
		newValue: () =>
			APLActionMultishield.create({
				maxShields: 3,
				maxOverlap: {
					value: {
						oneofKind: 'const',
						const: {
							val: '0ms',
						},
					},
				},
			}),
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'shield_spells', ''),
			AplHelpers.numberFieldConfig('maxShields', false, {
				label: 'Max Shields',
				labelTooltip: 'Maximum number of Shields to simultaneously apply.',
			}),
			AplValues.valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before a Shield expires when it may be refreshed.',
			}),
		],
	}),
	['channelSpell']: inputBuilder({
		label: 'Channel',
		submenu: ['Casting'],
		shortDescription: 'Channels the spell if possible, i.e. resource/cooldown/GCD/etc requirements are all met.',
		fullDescription: `
			<p>The difference between channeling a spell vs casting the spell is that channels can be interrupted. If the <b>Interrupt If</b> parameter is empty, this action is equivalent to <b>Cast</b>.</p>
			<p>The channel will be interrupted only if all of the following are true:</p>
			<ul>
				<li>Immediately following a tick of the channel</li>
				<li>The <b>Interrupt If</b> condition evaluates to <b>True</b></li>
				<li>Another action in the APL list is available</li>
			</ul>
			<p>Note that if you simply want to allow other actions to interrupt the channel, set <b>Interrupt If</b> to <b>True</b>.</p>
		`,
		newValue: () =>
			APLActionChannelSpell.create({
				interruptIf: {
					value: {
						oneofKind: 'gcdIsReady',
						gcdIsReady: {},
					},
				},
			}),
		fields: [
			AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', ''),
			AplHelpers.unitFieldConfig('target', 'targets'),
			AplValues.valueFieldConfig('interruptIf', {
				label: 'Interrupt If',
				labelTooltip: 'Condition which must be true to allow the channel to be interrupted.',
			}),
			AplHelpers.booleanFieldConfig('allowRecast', 'Recast', {
				labelTooltip: 'If checked, interrupts of this channel will recast the spell.',
			}),
		],
	}),
	['castAllStatBuffCooldowns']: inputBuilder({
		label: 'Cast All Stat Buff Cooldowns',
		submenu: ['Casting'],
		shortDescription: 'Casts all cooldowns that buff the specified stat type(s).',
		fullDescription: `
			<ul>
				<li>Does not cast cooldowns which are already controlled by other actions in the priority list.</li>
				<li>By default, this action will cast such cooldowns greedily as they become available. However, when embedded in a sequence, the action will only fire when ALL cooldowns matching the specified buff type(s) are ready.</li>
			</ul>
		`,
		newValue: () =>
			APLActionCastAllStatBuffCooldowns.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
		],
	}),
	['autocastOtherCooldowns']: inputBuilder({
		label: 'Autocast Other Cooldowns',
		submenu: ['Casting'],
		shortDescription: 'Auto-casts cooldowns as soon as they are ready.',
		fullDescription: `
			<ul>
				<li>Does not auto-cast cooldowns which are already controlled by other actions in the priority list.</li>
				<li>Cooldowns are usually cast immediately upon becoming ready, but there are some basic smart checks in place, e.g. don't use Mana CDs when near full mana.</li>
			</ul>
		`,
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: APLActionAutocastOtherCooldowns.create,
		fields: [],
	}),
	['wait']: inputBuilder({
		label: 'Wait',
		submenu: ['Timing'],
		shortDescription: 'Pauses all APL actions for a specified amount of time.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () =>
			APLActionWait.create({
				duration: {
					value: {
						oneofKind: 'const',
						const: {
							val: '1000ms',
						},
					},
				},
			}),
		fields: [AplValues.valueFieldConfig('duration')],
	}),
	['waitUntil']: inputBuilder({
		label: 'Wait Until',
		submenu: ['Timing'],
		shortDescription: 'Pauses all APL actions until the specified condition is <b>True</b>.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () => APLActionWaitUntil.create(),
		fields: [AplValues.valueFieldConfig('condition')],
	}),
	['schedule']: inputBuilder({
		label: 'Scheduled Action',
		submenu: ['Timing'],
		shortDescription: 'Executes the inner action once at each specified timing.',
		includeIf: (player: Player<any>, isPrepull: boolean) => !isPrepull,
		newValue: () =>
			APLActionSchedule.create({
				schedule: '0s, 60s',
				innerAction: {
					action: { oneofKind: 'castSpell', castSpell: {} },
				},
			}),
		fields: [
			AplHelpers.stringFieldConfig('schedule', {
				label: 'Do At',
				labelTooltip: 'Comma-separated list of timings. The inner action will be performed once at each timing.',
			}),
			actionFieldConfig('innerAction'),
		],
	}),
	['sequence']: inputBuilder({
		label: 'Sequence',
		submenu: ['Sequences'],
		shortDescription: 'A list of sub-actions to execute in the specified order.',
		fullDescription: `
			<p>Once one of the sub-actions has been performed, the next sub-action will not necessarily be immediately executed next. The system will restart at the beginning of the whole actions list (not the sequence). If the sequence is executed again, it will perform the next sub-action.</p>
			<p>When all actions have been performed, the sequence does NOT automatically reset; instead, it will be skipped from now on. Use the <b>Reset Sequence</b> action to reset it, if desired.</p>
		`,
		includeIf: (_, isPrepull: boolean) => !isPrepull,
		newValue: APLActionSequence.create,
		fields: [AplHelpers.stringFieldConfig('name'), actionListFieldConfig('actions')],
	}),
	['resetSequence']: inputBuilder({
		label: 'Reset Sequence',
		submenu: ['Sequences'],
		shortDescription: 'Restarts a sequence, so that the next time it executes it will perform its first sub-action.',
		fullDescription: `
			<p>Use the <b>name</b> field to refer to the sequence to be reset. The desired sequence must have the same (non-empty) value for its <b>name</b>.</p>
		`,
		includeIf: (_, isPrepull: boolean) => !isPrepull,
		newValue: APLActionResetSequence.create,
		fields: [AplHelpers.stringFieldConfig('sequenceName')],
	}),
	['strictSequence']: inputBuilder({
		label: 'Strict Sequence',
		submenu: ['Sequences'],
		shortDescription:
			'Like a regular <b>Sequence</b>, except all sub-actions are executed immediately after each other and the sequence resets automatically upon completion.',
		fullDescription: `
			<p>Strict Sequences do not begin unless ALL sub-actions are ready.</p>
		`,
		includeIf: (_, isPrepull: boolean) => !isPrepull,
		newValue: APLActionStrictSequence.create,
		fields: [actionListFieldConfig('actions')],
	}),
	['changeTarget']: inputBuilder({
		label: 'Change Target',
		submenu: ['Misc'],
		shortDescription: 'Sets the current target, which is the target of auto attacks and most casts by default.',
		newValue: () => APLActionChangeTarget.create(),
		fields: [AplHelpers.unitFieldConfig('newTarget', 'targets')],
	}),
	['activateAura']: inputBuilder({
		label: 'Activate Aura',
		submenu: ['Misc'],
		shortDescription: 'Activates an aura',
		includeIf: (_, isPrepull: boolean) => isPrepull,
		newValue: () => APLActionActivateAura.create(),
		fields: [AplHelpers.actionIdFieldConfig('auraId', 'auras')],
	}),
	['activateAuraWithStacks']: inputBuilder({
		label: 'Activate Aura With Stacks',
		submenu: ['Misc'],
		shortDescription: 'Activates an aura with the specified number of stacks',
		includeIf: (_, isPrepull: boolean) => isPrepull,
		newValue: () => APLActionActivateAuraWithStacks.create({
			numStacks: 1,
		}),
		fields: [AplHelpers.actionIdFieldConfig('auraId', 'stackable_auras'), AplHelpers.numberFieldConfig('numStacks', false, {
			label: 'stacks',
			labelTooltip: 'Desired number of initial aura stacks.',
		})],
	}),
	['activateAllStatBuffProcAuras']: inputBuilder({
		label: 'Activate All Stat Buff Proc Auras',
		submenu: ['Misc'],
		shortDescription: 'Activates all item/enchant proc auras that buff the specified stat type(s) using the specified item set.',
		includeIf: (_, isPrepull: boolean) => isPrepull,
		newValue: () =>
			APLActionActivateAllStatBuffProcAuras.create({
				swapSet: ItemSwapSet.Main,
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			itemSwapSetFieldConfig('swapSet'),
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
		],
	}),
	['cancelAura']: inputBuilder({
		label: 'Cancel Aura',
		submenu: ['Misc'],
		shortDescription: 'Deactivates an aura, equivalent to /cancelaura.',
		newValue: () => APLActionCancelAura.create(),
		fields: [AplHelpers.actionIdFieldConfig('auraId', 'auras')],
	}),
	['triggerIcd']: inputBuilder({
		label: 'Trigger ICD',
		submenu: ['Misc'],
		shortDescription: "Triggers an aura's ICD, putting it on cooldown. Example usage would be to desync an ICD cooldown before combat starts.",
		includeIf: (_, isPrepull: boolean) => isPrepull,
		newValue: () => APLActionTriggerICD.create(),
		fields: [AplHelpers.actionIdFieldConfig('auraId', 'icd_auras')],
	}),
	['itemSwap']: inputBuilder({
		label: 'Item Swap',
		submenu: ['Misc'],
		shortDescription: 'Swaps items, using the swap set specified in Settings.',
		includeIf: (player: Player<any>, _isPrepull: boolean) => itemSwapEnabledSpecs.includes(player.getSpec()),
		newValue: () => APLActionItemSwap.create(),
		fields: [itemSwapSetFieldConfig('swapSet')],
	}),
	['move']: inputBuilder({
		label: 'Move',
		submenu: ['Misc'],
		shortDescription: 'Starts a move to the desired range from target.',
		newValue: () => APLActionMove.create(),
		fields: [
			AplValues.valueFieldConfig('rangeFromTarget', {
				label: 'to Range',
				labelTooltip: 'Desired range from target.',
			}),
		],
	}),
	['moveDuration']: inputBuilder({
		label: 'Move duration',
		submenu: ['Misc'],
		shortDescription: 'The characters moves for the given duration.',
		newValue: () => APLActionMoveDuration.create(),
		fields: [
			AplValues.valueFieldConfig('duration', {
				label: 'Duration',
				labelTooltip: 'Amount of time the character should move.',
			}),
		],
	}),
	['customRotation']: inputBuilder({
		label: 'Custom Rotation',
		//submenu: ['Misc'],
		shortDescription: 'INTERNAL ONLY',
		includeIf: (_player: Player<any>, _isPrepull: boolean) => false, // Never show this, because its internal only.
		newValue: () => APLActionCustomRotation.create(),
		fields: [],
	}),

	// Class/spec specific actions
	['catOptimalRotationAction']: inputBuilder({
		label: 'Optimal Rotation Action',
		submenu: ['Feral Druid'],
		shortDescription: 'Executes optimized Feral DPS rotation using hardcoded algorithm.',
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecFeralDruid,
		newValue: () =>
			APLActionCatOptimalRotationAction.create({
				rotationType: FeralDruid_Rotation_AplType.SingleTarget,
				maintainFaerieFire: true,
				manualParams: true,
				minRoarOffset: 31.0,
				ripLeeway: 1,
				useRake: true,
				useBite: true,
				biteTime: 11.0,
				berserkBiteTime: 6.0,
				biteDuringExecute: true,
				allowAoeBerserk: false,
				meleeWeave: true,
				bearWeave: true,
				snekWeave: true,
				cancelPrimalMadness: false,
			}),
		fields: [
			AplHelpers.rotationTypeFieldConfig('rotationType'),
			AplHelpers.booleanFieldConfig('maintainFaerieFire', 'Maintain Faerie Fire', {
				labelTooltip: 'Maintain Faerie Fire debuff. Overwrites any external Sunder effects specified in settings.',
			}),
			AplHelpers.booleanFieldConfig('meleeWeave', 'Enable leave-weaving', {
				labelTooltip: 'Weave out of melee range for Stampede procs. Ignored for AoE rotation or if Stampede is not talented.',
			}),
			AplHelpers.booleanFieldConfig('bearWeave', 'Enable bear-weaving', {
				labelTooltip: 'Weave into Bear Form while pooling Energy. Ignored for AoE rotation.',
			}),
			AplHelpers.booleanFieldConfig('snekWeave', 'Use Albino Snake', {
				labelTooltip: 'Reset swing timer at the end of bear-weaves using Albino Snake pet. Ignored if not bear-weaving.',
			}),
			AplHelpers.booleanFieldConfig('allowAoeBerserk', 'Allow AoE Berserk', {
				labelTooltip: 'Allow Berserk usage in AoE rotation. Ignored for single target rotation.',
			}),
			AplHelpers.booleanFieldConfig('manualParams', 'Manual Advanced Parameters', {
				labelTooltip: 'Manually specify advanced parameters, otherwise will use preset defaults.',
			}),
			AplHelpers.numberFieldConfig('minRoarOffset', true, {
				label: 'Roar Offset',
				labelTooltip: 'Targeted offset in Rip/Roar timings. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.numberFieldConfig('ripLeeway', false, {
				label: 'Rip Leeway',
				labelTooltip: 'Rip leeway when optimizing Roar clips. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.booleanFieldConfig('useRake', 'Use Rake', {
				labelTooltip: 'Use Rake during rotation. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.booleanFieldConfig('useBite', 'Bite during rotation', {
				labelTooltip:
					'Use Bite during rotation rather than exclusively at end of fight. Ignored for AOE rotation or if not using manual advanced parameters.',
			}),
			AplHelpers.numberFieldConfig('biteTime', true, {
				label: 'Bite Time',
				labelTooltip: 'Min seconds remaining on Rip/Roar to allow a Bite. Ignored if not Biting during rotation.',
			}),
			AplHelpers.numberFieldConfig('berserkBiteTime', true, {
				label: 'Bite Time during Berserk',
				labelTooltip: 'More aggressive threshold when Berserk is active.',
			}),
			AplHelpers.booleanFieldConfig('biteDuringExecute', 'Bite during Execute phase', {
				labelTooltip:
					'Bite aggressively during Execute phase. Ignored if Blood in the Water is not talented, or if not using manual advanced parameters.',
			}),
			AplHelpers.booleanFieldConfig('cancelPrimalMadness', 'Enable Primal Madness cancellation', {
				labelTooltip:
					'Click off Primal Madness buff when doing so will result in net Energy gains. Ignored if Primal Madness is not talented, or if not using manual advanced parameters.',
			}),
		],
	}),
};
