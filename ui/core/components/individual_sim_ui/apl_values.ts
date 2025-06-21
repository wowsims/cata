import { Player } from '../../player.js';
import {
	APLValue,
	APLValueAllTrinketStatProcsActive,
	APLValueAnd,
	APLValueAnyTrinketStatProcsActive,
	APLValueAuraICDIsReadyWithReactionTime,
	APLValueAuraInternalCooldown,
	APLValueAuraIsActive,
	APLValueAuraIsActiveWithReactionTime,
	APLValueAuraIsInactiveWithReactionTime,
	APLValueAuraIsKnown,
	APLValueAuraNumStacks,
	APLValueAuraRemainingTime,
	APLValueAuraShouldRefresh,
	APLValueAutoTimeToNext,
	APLValueBossSpellIsCasting,
	APLValueBossSpellTimeToReady,
	APLValueCatExcessEnergy,
	APLValueCatNewSavageRoarDuration,
	APLValueChannelClipDelay,
	APLValueCompare,
	APLValueCompare_ComparisonOperator as ComparisonOperator,
	APLValueConst,
	APLValueCurrentComboPoints,
	APLValueCurrentEclipsePhase,
	APLValueCurrentEnergy,
	APLValueCurrentFocus,
	APLValueCurrentGenericResource,
	APLValueCurrentHealth,
	APLValueCurrentHealthPercent,
	APLValueCurrentLunarEnergy,
	APLValueCurrentMana,
	APLValueCurrentManaPercent,
	APLValueCurrentNonDeathRuneCount,
	APLValueCurrentRage,
	APLValueCurrentRuneActive,
	APLValueCurrentRuneCount,
	APLValueCurrentRuneDeath,
	APLValueCurrentRunicPower,
	APLValueCurrentSolarEnergy,
	APLValueCurrentTime,
	APLValueCurrentTimePercent,
	APLValueDotIsActive,
	APLValueDotPercentIncrease,
	APLValueDotRemainingTime,
	APLValueDotTickFrequency,
	APLValueEnergyRegenPerSecond,
	APLValueEnergyTimeToTarget,
	APLValueFocusRegenPerSecond,
	APLValueFocusTimeToTarget,
	APLValueFrontOfTarget,
	APLValueGCDIsReady,
	APLValueGCDTimeToReady,
	APLValueInputDelay,
	APLValueIsExecutePhase,
	APLValueIsExecutePhase_ExecutePhaseThreshold as ExecutePhaseThreshold,
	APLValueMageCurrentCombustionDotEstimate,
	APLValueMath,
	APLValueMath_MathOperator as MathOperator,
	APLValueMax,
	APLValueMaxComboPoints,
	APLValueMaxEnergy,
	APLValueMaxFocus,
	APLValueMaxHealth,
	APLValueMaxRage,
	APLValueMaxRunicPower,
	APLValueMin,
	APLValueMonkCurrentChi,
	APLValueMonkMaxChi,
	APLValueNextRuneCooldown,
	APLValueNot,
	APLValueNumberTargets,
	APLValueNumEquippedStatProcTrinkets,
	APLValueNumStatBuffCooldowns,
	APLValueOr,
	APLValueProtectionPaladinDamageTakenLastGlobal,
	APLValueRemainingTime,
	APLValueRemainingTimePercent,
	APLValueRuneCooldown,
	APLValueRuneSlotCooldown,
	APLValueSequenceIsComplete,
	APLValueSequenceIsReady,
	APLValueSequenceTimeToReady,
	APLValueShamanFireElementalDuration,
	APLValueSpellCanCast,
	APLValueSpellCastTime,
	APLValueSpellChanneledTicks,
	APLValueSpellCPM,
	APLValueSpellCurrentCost,
	APLValueSpellIsChanneling,
	APLValueSpellIsKnown,
	APLValueSpellIsReady,
	APLValueSpellNumCharges,
	APLValueSpellTimeToCharge,
	APLValueSpellTimeToReady,
	APLValueSpellTravelTime,
	APLValueTotemRemainingTime,
	APLValueTrinketProcsMaxRemainingICD,
	APLValueTrinketProcsMinRemainingTime,
	APLValueUnitIsMoving,
	APLValueWarlockHandOfGuldanInFlight,
	APLValueWarlockHauntInFlight,
} from '../../proto/apl.js';
import { Class, Spec } from '../../proto/common.js';
import { ShamanTotems_TotemType as TotemType } from '../../proto/shaman.js';
import SecondaryResource from '../../proto_utils/secondary_resource';
import { EventID } from '../../typed_event.js';
import { randomUUID } from '../../utils';
import { Input, InputConfig } from '../input.js';
import { TextDropdownPicker, TextDropdownValueConfig } from '../pickers/dropdown_picker.jsx';
import { ListItemPickerConfig, ListPicker } from '../pickers/list_picker.jsx';
import * as AplHelpers from './apl_helpers.js';

export interface APLValuePickerConfig extends InputConfig<Player<any>, APLValue | undefined> {}

type APLValue_Value = APLValue['value'];
export type APLValueKind = APLValue_Value['oneofKind'];
type ValidAPLValueKind = NonNullable<APLValueKind>;

export type APLValueImplStruct<F extends APLValueKind> = Extract<APLValue_Value, { oneofKind: F }>;

// Get the implementation type for a specific kind using infer
type APLValueImplFor<F extends ValidAPLValueKind> =
	APLValueImplStruct<F> extends { [K in F]: infer T }
		? T
		: never;

// Map all valid kinds to their implementation types
type APLValueImplMap = {
	[K in ValidAPLValueKind]: APLValueImplFor<K>;
};

export type APLValueImplType = APLValueImplMap[ValidAPLValueKind] | undefined;

export class APLValuePicker extends Input<Player<any>, APLValue | undefined> {
	private kindPicker: TextDropdownPicker<Player<any>, APLValueKind>;

	private currentKind: APLValueKind;
	private valuePicker: Input<Player<any>, any> | null;

	constructor(parent: HTMLElement, player: Player<any>, config: APLValuePickerConfig) {
		super(parent, 'apl-value-picker-root', player, config);

		const isPrepull = this.rootElem.closest('.apl-prepull-action-picker') != null;

		const allValueKinds = (Object.keys(valueKindFactories) as ValidAPLValueKind[]).filter(
			(valueKind): valueKind is ValidAPLValueKind => (!!valueKind && valueKindFactories[valueKind].includeIf?.(player, isPrepull)) ?? true,
		);

		if (this.rootElem.parentElement!.classList.contains('list-picker-item')) {
			const itemHeaderElem = ListPicker.getItemHeaderElem(this) || this.rootElem;
			ListPicker.makeListItemValidations(
				itemHeaderElem,
				player,
				player => player.getCurrentStats().rotationStats?.uuidValidations?.find(v => v.uuid?.value === this.rootElem.id)?.validations || [],
			);
		}

		this.kindPicker = new TextDropdownPicker(this.rootElem, player, {
			defaultLabel: 'No Condition',
			id: randomUUID(),
			values: [
				{
					value: undefined,
					label: '<None>',
				} as TextDropdownValueConfig<APLValueKind>,
			].concat(
				allValueKinds.map(kind => {
					const factory = valueKindFactories[kind];
					const resolveString = factory.dynamicStringResolver || ((value: string) => value);
					return {
						value: kind,
						label: resolveString(factory.label, player),
						submenu: factory.submenu,
						tooltip: factory.fullDescription
							? `<p>${resolveString(factory.shortDescription, player)}</p> ${resolveString(factory.fullDescription, player)}`
							: resolveString(factory.shortDescription, player),
					};
				}),
			),
			equals: (a, b) => a == b,
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: (_player: Player<any>) => this.getSourceValue()?.value.oneofKind,
			setValue: (eventID: EventID, player: Player<any>, newKind: APLValueKind) => {
				const sourceValue = this.getSourceValue();
				const oldKind = sourceValue?.value.oneofKind;
				if (oldKind == newKind) {
					return;
				}

				if (newKind) {
					const factory = valueKindFactories[newKind];
					let newSourceValue = this.makeAPLValue(newKind, factory.newValue());
					if (sourceValue) {
						// Some pre-fill logic when swapping kinds.
						if (oldKind && this.valuePicker) {
							if (newKind == 'not') {
								(newSourceValue.value as APLValueImplStruct<'not'>).not.val = this.makeAPLValue(oldKind, this.valuePicker.getInputValue());
							} else if (sourceValue.value.oneofKind == 'not' && sourceValue.value.not.val?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.not.val;
							} else if (newKind == 'and') {
								if (sourceValue.value.oneofKind == 'or') {
									(newSourceValue.value as APLValueImplStruct<'and'>).and.vals = sourceValue.value.or.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'and'>).and.vals = [
										this.makeAPLValue(oldKind, this.valuePicker.getInputValue()),
									];
								}
							} else if (newKind == 'or') {
								if (sourceValue.value.oneofKind == 'and') {
									(newSourceValue.value as APLValueImplStruct<'or'>).or.vals = sourceValue.value.and.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'or'>).or.vals = [this.makeAPLValue(oldKind, this.valuePicker.getInputValue())];
								}
							} else if (newKind == 'min') {
								if (sourceValue.value.oneofKind == 'max') {
									(newSourceValue.value as APLValueImplStruct<'min'>).min.vals = sourceValue.value.max.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'min'>).min.vals = [
										this.makeAPLValue(oldKind, this.valuePicker.getInputValue()),
									];
								}
							} else if (newKind == 'max') {
								if (sourceValue.value.oneofKind == 'min') {
									(newSourceValue.value as APLValueImplStruct<'max'>).max.vals = sourceValue.value.min.vals;
								} else {
									(newSourceValue.value as APLValueImplStruct<'max'>).max.vals = [
										this.makeAPLValue(oldKind, this.valuePicker.getInputValue()),
									];
								}
							} else if (sourceValue.value.oneofKind == 'and' && sourceValue.value.and.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.and.vals[0];
							} else if (sourceValue.value.oneofKind == 'or' && sourceValue.value.or.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.or.vals[0];
							} else if (sourceValue.value.oneofKind == 'min' && sourceValue.value.min.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.min.vals[0];
							} else if (sourceValue.value.oneofKind == 'max' && sourceValue.value.max.vals?.[0]?.value.oneofKind == newKind) {
								newSourceValue = sourceValue.value.max.vals[0];
							} else if (newKind == 'cmp') {
								(newSourceValue.value as APLValueImplStruct<'cmp'>).cmp.lhs = this.makeAPLValue(oldKind, this.valuePicker.getInputValue());
							}
						}
					}
					if (sourceValue) {
						sourceValue.value = newSourceValue.value;
					} else {
						this.setSourceValue(eventID, newSourceValue);
					}
				} else {
					this.setSourceValue(eventID, undefined);
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});

		this.currentKind = undefined;
		this.valuePicker = null;

		this.init();
	}

	getInputElem(): HTMLElement | null {
		return this.rootElem;
	}

	getInputValue(): APLValue | undefined {
		const kind = this.kindPicker.getInputValue();
		if (!kind) {
			return undefined;
		} else {
			return APLValue.create({
				value: {
					oneofKind: kind,
					...(() => {
						const val: any = {};
						if (kind && this.valuePicker) {
							val[kind] = this.valuePicker.getInputValue();
						}
						return val;
					})(),
				},
				uuid: { value: randomUUID() },
			});
		}
	}

	setInputValue(newValue: APLValue | undefined) {
		const newKind = newValue?.value.oneofKind;
		this.updateValuePicker(newKind);

		if (newKind && newValue) {
			this.valuePicker!.setInputValue((newValue.value as any)[newKind]);
		}

		if (newValue) {
			if (!newValue.uuid || newValue.uuid.value == '') {
				newValue.uuid = {
					value: randomUUID(),
				};
			}
			this.rootElem.id = newValue.uuid!.value;
		}
	}

	private makeAPLValue<K extends ValidAPLValueKind>(kind: K, implVal: APLValueImplMap[K]): APLValue {
		if (!kind) {
			return APLValue.create({
				uuid: { value: randomUUID() },
			});
		}
		const obj: any = { oneofKind: kind };
		obj[kind] = implVal;
		return APLValue.create({
			value: obj,
			uuid: { value: randomUUID() },
		});
	}

	private updateValuePicker(newKind: APLValueKind) {
		const oldKind = this.currentKind;
		if (newKind == oldKind) {
			return;
		}
		this.currentKind = newKind;

		if (this.valuePicker) {
			this.valuePicker.rootElem.remove();
			this.valuePicker = null;
		}

		if (!newKind) {
			return;
		}

		this.kindPicker.setInputValue(newKind);

		const factory = valueKindFactories[newKind];
		this.valuePicker = factory.factory(this.rootElem, this.modObject, {
			id: randomUUID(),
			changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
			getValue: () => {
				const sourceVal = this.getSourceValue();
				return sourceVal ? (sourceVal.value as any)[newKind] || factory.newValue() : factory.newValue();
			},
			setValue: (eventID: EventID, player: Player<any>, newValue: any) => {
				const sourceVal = this.getSourceValue();
				if (sourceVal) {
					(sourceVal.value as any)[newKind] = newValue;
				}
				player.rotationChangeEmitter.emit(eventID);
			},
		});
	}
}

type ValueKindConfig<T> = {
	label: string;
	submenu?: Array<string>;
	shortDescription: string;
	fullDescription?: string;
	newValue: () => T;
	includeIf?: (player: Player<any>, isPrepull: boolean) => boolean;
	factory: (parent: HTMLElement, player: Player<any>, config: InputConfig<Player<any>, T>) => Input<Player<any>, T>;
	dynamicStringResolver?: (value: string, player: Player<any>) => string;
};

function comparisonOperatorFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ComparisonOperator.OpEq,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				id: randomUUID(),
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: ComparisonOperator.OpEq, label: '==' },
					{ value: ComparisonOperator.OpNe, label: '!=' },
					{ value: ComparisonOperator.OpGe, label: '>=' },
					{ value: ComparisonOperator.OpGt, label: '>' },
					{ value: ComparisonOperator.OpLe, label: '<=' },
					{ value: ComparisonOperator.OpLt, label: '<' },
				],
			}),
	};
}

function mathOperatorFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => MathOperator.OpAdd,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				id: randomUUID(),
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: MathOperator.OpAdd, label: '+' },
					{ value: MathOperator.OpSub, label: '-' },
					{ value: MathOperator.OpMul, label: '*' },
					{ value: MathOperator.OpDiv, label: '/' },
				],
			}),
	};
}

function executePhaseThresholdFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => ExecutePhaseThreshold.E20,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				id: randomUUID(),
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: ExecutePhaseThreshold.E20, label: '20%' },
					{ value: ExecutePhaseThreshold.E25, label: '25%' },
					{ value: ExecutePhaseThreshold.E35, label: '35%' },
					{ value: ExecutePhaseThreshold.E45, label: '45%' },
					{ value: ExecutePhaseThreshold.E90, label: '90%' },
				],
			}),
	};
}

function totemTypeFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => TotemType.Water,
		factory: (parent, player, config) =>
			new TextDropdownPicker(parent, player, {
				id: randomUUID(),
				...config,
				defaultLabel: 'None',
				equals: (a, b) => a == b,
				values: [
					{ value: TotemType.Earth, label: 'Earth' },
					{ value: TotemType.Air, label: 'Air' },
					{ value: TotemType.Fire, label: 'Fire' },
					{ value: TotemType.Water, label: 'Water' },
				],
			}),
	};
}

export function valueFieldConfig(
	field: string,
	options?: Partial<AplHelpers.APLPickerBuilderFieldConfig<any, any>>,
): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () =>
			APLValue.create({
				uuid: { value: randomUUID() },
			}),
		factory: (parent, player, config) => new APLValuePicker(parent, player, config),
		...(options || {}),
	};
}

export function valueListFieldConfig(field: string): AplHelpers.APLPickerBuilderFieldConfig<any, any> {
	return {
		field: field,
		newValue: () => [],
		factory: (parent, player, config) =>
			new ListPicker<Player<any>, APLValue | undefined>(parent, player, {
				...config,
				// Override setValue to replace undefined elements with default messages.
				setValue: (eventID: EventID, player: Player<any>, newValue: Array<APLValue | undefined>) => {
					config.setValue(
						eventID,
						player,
						newValue.map(val => {
							return (
								val ||
								APLValue.create({
									uuid: { value: randomUUID() },
								})
							);
						}),
					);
				},
				itemLabel: 'Value',
				newItem: () => {
					return APLValue.create({
						uuid: { value: randomUUID() },
					});
				},
				copyItem: (oldValue: APLValue | undefined) => (oldValue ? APLValue.clone(oldValue) : oldValue),
				newItemPicker: (
					parent: HTMLElement,
					listPicker: ListPicker<Player<any>, APLValue | undefined>,
					index: number,
					config: ListItemPickerConfig<Player<any>, APLValue | undefined>,
				) => new APLValuePicker(parent, player, config),
				allowedActions: ['copy', 'create', 'delete', 'move'],
				actions: {
					create: {
						useIcon: true,
					},
				},
			}),
	};
}

function inputBuilder<T extends APLValueImplType>(
	config: {
		fields: Array<AplHelpers.APLPickerBuilderFieldConfig<T, keyof T>>;
	} & Omit<ValueKindConfig<T>, 'factory'>,
): ValueKindConfig<T> {
	return {
		label: config.label,
		submenu: config.submenu,
		shortDescription: config.shortDescription,
		fullDescription: config.fullDescription,
		newValue: config.newValue,
		includeIf: config.includeIf,
		factory: AplHelpers.aplInputBuilder(config.newValue, config.fields),
		dynamicStringResolver: config.dynamicStringResolver,
	};
}

const valueKindFactories: { [f in ValidAPLValueKind]: ValueKindConfig<APLValueImplMap[f]> } = {
	// Operators
	const: inputBuilder({
		label: 'Const',
		shortDescription: 'A fixed value.',
		fullDescription: `
		<p>
			Examples:
			<ul>
				<li><b>Number:</b> '123', '0.5', '-10'</li>
				<li><b>Time:</b> '100ms', '5s', '3m'</li>
				<li><b>Percentage:</b> '30%'</li>
			</ul>
		</p>
		`,
		newValue: APLValueConst.create,
		fields: [AplHelpers.stringFieldConfig('val')],
	}),
	cmp: inputBuilder({
		label: 'Compare',
		submenu: ['Logic'],
		shortDescription: 'Compares two values.',
		newValue: APLValueCompare.create,
		fields: [valueFieldConfig('lhs'), comparisonOperatorFieldConfig('op'), valueFieldConfig('rhs')],
	}),
	math: inputBuilder({
		label: 'Math',
		submenu: ['Logic'],
		shortDescription: 'Do basic math on two values.',
		newValue: APLValueMath.create,
		fields: [valueFieldConfig('lhs'), mathOperatorFieldConfig('op'), valueFieldConfig('rhs')],
	}),
	max: inputBuilder({
		label: 'Max',
		submenu: ['Logic'],
		shortDescription: 'Returns the largest value among the subvalues.',
		newValue: APLValueMax.create,
		fields: [valueListFieldConfig('vals')],
	}),
	min: inputBuilder({
		label: 'Min',
		submenu: ['Logic'],
		shortDescription: 'Returns the smallest value among the subvalues.',
		newValue: APLValueMin.create,
		fields: [valueListFieldConfig('vals')],
	}),
	and: inputBuilder({
		label: 'All of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if all of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueAnd.create,
		fields: [valueListFieldConfig('vals')],
	}),
	or: inputBuilder({
		label: 'Any of',
		submenu: ['Logic'],
		shortDescription: 'Returns <b>True</b> if any of the sub-values are <b>True</b>, otherwise <b>False</b>',
		newValue: APLValueOr.create,
		fields: [valueListFieldConfig('vals')],
	}),
	not: inputBuilder({
		label: 'Not',
		submenu: ['Logic'],
		shortDescription: 'Returns the opposite of the inner value, i.e. <b>True</b> if the value is <b>False</b> and vice-versa.',
		newValue: APLValueNot.create,
		fields: [valueFieldConfig('val')],
	}),

	// Encounter
	currentTime: inputBuilder({
		label: 'Current Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration.',
		newValue: APLValueCurrentTime.create,
		fields: [],
	}),
	currentTimePercent: inputBuilder({
		label: 'Current Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the current sim iteration, as a percentage.',
		newValue: APLValueCurrentTimePercent.create,
		fields: [],
	}),
	remainingTime: inputBuilder({
		label: 'Remaining Time',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration.',
		newValue: APLValueRemainingTime.create,
		fields: [],
	}),
	remainingTimePercent: inputBuilder({
		label: 'Remaining Time (%)',
		submenu: ['Encounter'],
		shortDescription: 'Elapsed time of the remaining sim iteration, as a percentage.',
		newValue: APLValueRemainingTimePercent.create,
		fields: [],
	}),
	isExecutePhase: inputBuilder({
		label: 'Is Execute Phase',
		submenu: ['Encounter'],
		shortDescription:
			"<b>True</b> if the encounter is in Execute Phase, meaning the target's health is less than the given threshold, otherwise <b>False</b>.",
		newValue: APLValueIsExecutePhase.create,
		fields: [executePhaseThresholdFieldConfig('threshold')],
	}),
	numberTargets: inputBuilder({
		label: 'Number of Targets',
		submenu: ['Encounter'],
		shortDescription: 'Count of targets in the current encounter',
		newValue: APLValueNumberTargets.create,
		fields: [],
	}),
	frontOfTarget: inputBuilder({
		label: 'Front of Target',
		submenu: ['Encounter'],
		shortDescription: '<b>True</b> if facing from of target',
		newValue: APLValueFrontOfTarget.create,
		fields: [],
	}),

	// Boss
	bossSpellIsCasting: inputBuilder({
		label: 'Spell is Casting',
		submenu: ['Boss'],
		shortDescription: '',
		newValue: APLValueBossSpellIsCasting.create,
		fields: [
			AplHelpers.unitFieldConfig('targetUnit', 'targets'),
			AplHelpers.actionIdFieldConfig('spellId', 'non_instant_spells', 'targetUnit', 'currentTarget'),
		],
	}),
	bossSpellTimeToReady: inputBuilder({
		label: 'Spell Time to Ready',
		submenu: ['Boss'],
		shortDescription: '',
		newValue: APLValueBossSpellTimeToReady.create,
		fields: [AplHelpers.unitFieldConfig('targetUnit', 'targets'), AplHelpers.actionIdFieldConfig('spellId', 'spells', 'targetUnit', 'currentTarget')],
	}),

	// Unit
	unitIsMoving: inputBuilder({
		label: 'Is moving',
		submenu: ['Unit'],
		shortDescription: '',
		newValue: APLValueUnitIsMoving.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources')],
	}),

	// Resources
	currentHealth: inputBuilder({
		label: 'Current Health',
		submenu: ['Resources', 'Health'],
		shortDescription: 'Amount of currently available Health.',
		newValue: APLValueCurrentHealth.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources')],
	}),
	currentHealthPercent: inputBuilder({
		label: 'Current Health (%)',
		submenu: ['Resources', 'Health'],
		shortDescription: 'Amount of currently available Health, as a percentage.',
		newValue: APLValueCurrentHealthPercent.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources')],
	}),
	maxHealth: inputBuilder({
		label: 'Max Health',
		submenu: ['Resources', 'Health'],
		shortDescription: 'Amount of currently available maximum Health.',
		newValue: APLValueMaxHealth.create,
		fields: [],
	}),
	currentMana: inputBuilder({
		label: 'Current Mana',
		submenu: ['Resources', 'Mana'],
		shortDescription: 'Amount of currently available Mana.',
		newValue: APLValueCurrentMana.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			return clss !== Class.ClassDeathKnight && clss !== Class.ClassHunter && clss !== Class.ClassRogue && clss !== Class.ClassWarrior;
		},
		fields: [],
	}),
	currentManaPercent: inputBuilder({
		label: 'Current Mana (%)',
		submenu: ['Resources', 'Mana'],
		shortDescription: 'Amount of currently available Mana, as a percentage.',
		newValue: APLValueCurrentManaPercent.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			return clss !== Class.ClassDeathKnight && clss !== Class.ClassHunter && clss !== Class.ClassRogue && clss !== Class.ClassWarrior;
		},
		fields: [],
	}),
	currentRage: inputBuilder({
		label: 'Current Rage',
		submenu: ['Resources', 'Rage'],
		shortDescription: 'Amount of currently available Rage.',
		newValue: APLValueCurrentRage.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassWarrior;
		},
		fields: [],
	}),
	maxRage: inputBuilder({
		label: 'Max Rage',
		submenu: ['Resources', 'Rage'],
		shortDescription: 'Amount of maximum available Rage.',
		newValue: APLValueMaxRage.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassWarrior;
		},
		fields: [],
	}),
	currentFocus: inputBuilder({
		label: 'Current Focus',
		submenu: ['Resources', 'Focus'],
		shortDescription: 'Amount of currently available Focus.',
		newValue: APLValueCurrentFocus.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassHunter,
		fields: [],
	}),
	maxFocus: inputBuilder({
		label: 'Max Focus',
		submenu: ['Resources', 'Focus'],
		shortDescription: 'Amount of maximum available Focus.',
		newValue: APLValueMaxFocus.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassHunter,
		fields: [],
	}),
	focusRegenPerSecond: inputBuilder({
		label: 'Focus Regen Per Second',
		submenu: ['Resources', 'Focus'],
		shortDescription: 'Focus regen per second.',
		newValue: APLValueFocusRegenPerSecond.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassHunter,
		fields: [],
	}),
	focusTimeToTarget: inputBuilder({
		label: 'Estimated Time To Target Focus',
		submenu: ['Resources', 'Focus'],
		shortDescription: 'Estimated time until target Focus is reached, will return 0 if at or above target.',
		newValue: APLValueFocusTimeToTarget.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassHunter,
		fields: [valueFieldConfig('targetFocus')],
	}),
	currentEnergy: inputBuilder({
		label: 'Current Energy',
		submenu: ['Resources', 'Energy'],
		shortDescription: 'Amount of currently available Energy.',
		newValue: APLValueCurrentEnergy.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue || clss === Class.ClassMonk;
		},
		fields: [],
	}),
	maxEnergy: inputBuilder({
		label: 'Max Energy',
		submenu: ['Resources', 'Energy'],
		shortDescription: 'Amount of maximum available Energy.',
		newValue: APLValueMaxEnergy.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue || clss === Class.ClassMonk;
		},
		fields: [],
	}),
	energyRegenPerSecond: inputBuilder({
		label: 'Energy Regen Per Second',
		submenu: ['Resources', 'Energy'],
		shortDescription: 'Energy regen per second.',
		newValue: APLValueEnergyRegenPerSecond.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue || clss === Class.ClassMonk;
		},
		fields: [],
	}),
	energyTimeToTarget: inputBuilder({
		label: 'Estimated Time To Target Energy',
		submenu: ['Resources', 'Energy'],
		shortDescription: 'Estimated time until target Energy is reached, will return 0 if at or above target.',
		newValue: APLValueEnergyTimeToTarget.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue || clss === Class.ClassMonk;
		},
		fields: [valueFieldConfig('targetEnergy')],
	}),
	currentComboPoints: inputBuilder({
		label: 'Current Combo Points',
		submenu: ['Resources', 'Combo Points'],
		shortDescription: 'Amount of currently available Combo Points.',
		newValue: APLValueCurrentComboPoints.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue;
		},
		fields: [],
	}),
	maxComboPoints: inputBuilder({
		label: 'Max Combo Points',
		submenu: ['Resources', 'Combo Points'],
		shortDescription: 'Amount of maximum available Combo Points.',
		newValue: APLValueMaxComboPoints.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return spec === Spec.SpecFeralDruid || spec === Spec.SpecGuardianDruid || clss === Class.ClassRogue;
		},
		fields: [],
	}),
	monkCurrentChi: inputBuilder({
		label: 'Current Chi',
		submenu: ['Resources', 'Chi'],
		shortDescription: 'Amount of currently available Chi.',
		newValue: APLValueMonkCurrentChi.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() === Class.ClassMonk,
		fields: [],
	}),
	monkMaxChi: inputBuilder({
		label: 'Max Chi',
		submenu: ['Resources', 'Chi'],
		shortDescription: 'Amount of maximum available Chi.',
		newValue: APLValueMonkMaxChi.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() === Class.ClassMonk,
		fields: [],
	}),
	currentRunicPower: inputBuilder({
		label: 'Current Runic Power',
		submenu: ['Resources', 'Runic Power'],
		shortDescription: 'Amount of currently available Runic Power.',
		newValue: APLValueCurrentRunicPower.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [],
	}),
	maxRunicPower: inputBuilder({
		label: 'Max Runic Power',
		submenu: ['Resources', 'Runic Power'],
		shortDescription: 'Amount of maximum available Runic Power.',
		newValue: APLValueMaxRunicPower.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [],
	}),
	currentLunarEnergy: inputBuilder({
		label: 'Solar Energy',
		submenu: ['Resources', 'Eclipse'],
		shortDescription: 'Amount of currently available Solar Energy.',
		newValue: APLValueCurrentSolarEnergy.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecBalanceDruid,
		fields: [],
	}),
	currentSolarEnergy: inputBuilder({
		label: 'Lunar Energy',
		submenu: ['Resources', 'Eclipse'],
		shortDescription: 'Amount of currently available Lunar Energy',
		newValue: APLValueCurrentLunarEnergy.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecBalanceDruid,
		fields: [],
	}),
	druidCurrentEclipsePhase: inputBuilder({
		label: 'Current Eclipse Phase',
		submenu: ['Resources', 'Eclipse'],
		shortDescription: 'The eclipse phase the druid currently is in.',
		newValue: APLValueCurrentEclipsePhase.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecBalanceDruid,
		fields: [AplHelpers.eclipseTypeFieldConfig('eclipsePhase')],
	}),
	currentGenericResource: inputBuilder({
		label: '{GENERIC_RESOURCE}',
		submenu: ['Resources'],
		shortDescription: 'Amount of currently available {GENERIC_RESOURCE}.',
		newValue: APLValueCurrentGenericResource.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => SecondaryResource.hasSecondaryResource(player.getSpec()),
		fields: [],
		dynamicStringResolver: (value: string, player: Player<any>) => player.secondaryResource?.replaceResourceName(value) || '',
	}),

	// Resources Rune
	currentRuneCount: inputBuilder({
		label: 'Num Runes',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of currently available Runes of certain type including Death.',
		newValue: APLValueCurrentRuneCount.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeTypeFieldConfig('runeType', true)],
	}),
	currentNonDeathRuneCount: inputBuilder({
		label: 'Num Non Death Runes',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of currently available Runes of certain type ignoring Death',
		newValue: APLValueCurrentNonDeathRuneCount.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeTypeFieldConfig('runeType', false)],
	}),
	currentRuneActive: inputBuilder({
		label: 'Rune Is Ready',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Is the rune of a certain slot currently available.',
		newValue: APLValueCurrentRuneActive.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeSlotFieldConfig('runeSlot')],
	}),
	currentRuneDeath: inputBuilder({
		label: 'Rune Is Death',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Is the rune of a certain slot currently converted to Death.',
		newValue: APLValueCurrentRuneDeath.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeSlotFieldConfig('runeSlot')],
	}),
	runeCooldown: inputBuilder({
		label: 'Rune Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a rune of certain type is ready to use.<br><b>NOTE:</b> Returns 0 if there is a rune available',
		newValue: APLValueRuneCooldown.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeTypeFieldConfig('runeType', false)],
	}),
	nextRuneCooldown: inputBuilder({
		label: 'Next Rune Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a 2nd rune of certain type is ready to use.<br><b>NOTE:</b> Returns 0 if there are 2 runes available',
		newValue: APLValueNextRuneCooldown.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeTypeFieldConfig('runeType', false)],
	}),
	runeSlotCooldown: inputBuilder({
		label: 'Rune Slot Cooldown',
		submenu: ['Resources', 'Runes'],
		shortDescription: 'Amount of time until a rune of certain slot is ready to use.<br><b>NOTE:</b> Returns 0 if rune is ready',
		newValue: APLValueRuneSlotCooldown.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassDeathKnight,
		fields: [AplHelpers.runeSlotFieldConfig('runeSlot')],
	}),

	// GCD
	gcdIsReady: inputBuilder({
		label: 'GCD Is Ready',
		submenu: ['GCD'],
		shortDescription: '<b>True</b> if the GCD is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueGCDIsReady.create,
		fields: [],
	}),
	gcdTimeToReady: inputBuilder({
		label: 'GCD Time To Ready',
		submenu: ['GCD'],
		shortDescription: 'Amount of time remaining before the GCD comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueGCDTimeToReady.create,
		fields: [],
	}),

	// Auto attacks
	autoTimeToNext: inputBuilder({
		label: 'Time To Next Auto',
		submenu: ['Auto'],
		shortDescription: 'Amount of time remaining before the next Main-hand or Off-hand melee attack, or <b>0</b> if autoattacks are not engaged.',
		newValue: APLValueAutoTimeToNext.create,
		includeIf(player: Player<any>, _isPrepull: boolean) {
			const clss = player.getClass();
			const spec = player.getSpec();
			return (
				clss !== Class.ClassHunter &&
				clss !== Class.ClassMage &&
				clss !== Class.ClassPriest &&
				clss !== Class.ClassWarlock &&
				spec !== Spec.SpecBalanceDruid &&
				spec !== Spec.SpecElementalShaman
			);
		},
		fields: [],
	}),

	// Spells
	spellIsKnown: inputBuilder({
		label: 'Spell Known',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if the spell is currently known, otherwise <b>False</b>.',
		newValue: APLValueSpellIsKnown.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellCurrentCost: inputBuilder({
		label: 'Current Cost',
		submenu: ['Spell'],
		shortDescription: 'Returns current resource cost of spell',
		newValue: APLValueSpellCurrentCost.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellCanCast: inputBuilder({
		label: 'Can Cast',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if all requirements for casting the spell are currently met, otherwise <b>False</b>.',
		fullDescription: `
			<p>The <b>Cast Spell</b> action does not need to be conditioned on this, because it applies this check automatically.</p>
		`,
		newValue: APLValueSpellCanCast.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellIsReady: inputBuilder({
		label: 'Is Ready',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if the spell is not on cooldown, otherwise <b>False</b>.',
		newValue: APLValueSpellIsReady.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellTimeToReady: inputBuilder({
		label: 'Time To Ready',
		submenu: ['Spell'],
		shortDescription: 'Amount of time remaining before the spell comes off cooldown, or <b>0</b> if it is not on cooldown.',
		newValue: APLValueSpellTimeToReady.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellCastTime: inputBuilder({
		label: 'Cast Time',
		submenu: ['Spell'],
		shortDescription: 'Amount of time to cast the spell including any haste and spell cast time adjustments.',
		newValue: APLValueSpellCastTime.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellTravelTime: inputBuilder({
		label: 'Travel Time',
		submenu: ['Spell'],
		shortDescription: 'Amount of time for the spell to travel to the target.',
		newValue: APLValueSpellTravelTime.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellCpm: inputBuilder({
		label: 'CPM',
		submenu: ['Spell'],
		shortDescription: 'Casts Per Minute for the spell so far in the current iteration.',
		newValue: APLValueSpellCPM.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellIsChanneling: inputBuilder({
		label: 'Is Channeling',
		submenu: ['Spell'],
		shortDescription: '<b>True</b> if this spell is currently being channeled, otherwise <b>False</b>.',
		newValue: APLValueSpellIsChanneling.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', '')],
	}),
	spellChanneledTicks: inputBuilder({
		label: 'Channeled Ticks',
		submenu: ['Spell'],
		shortDescription: 'The number of completed ticks in the current channel of this spell, or <b>0</b> if the spell is not being channeled.',
		newValue: APLValueSpellChanneledTicks.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'channel_spells', '')],
	}),
	spellNumCharges: inputBuilder({
		label: 'Number of Charges',
		submenu: ['Spell'],
		shortDescription: 'The number of charges that are currently available for the spell.',
		newValue: APLValueSpellNumCharges.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	spellTimeToCharge: inputBuilder({
		label: 'Time to next Charge',
		submenu: ['Spell'],
		shortDescription: 'The time until the next charge is available. 0 if spell has all charges avaialable.',
		newValue: APLValueSpellTimeToCharge.create,
		fields: [AplHelpers.actionIdFieldConfig('spellId', 'castable_spells', '')],
	}),
	channelClipDelay: inputBuilder({
		label: 'Channel Clip Delay',
		submenu: ['Spell'],
		shortDescription: 'The amount of time specified by the <b>Channel Clip Delay</b> setting.',
		newValue: APLValueChannelClipDelay.create,
		fields: [],
	}),
	inputDelay: inputBuilder({
		label: 'Input Delay',
		submenu: ['Spell'],
		shortDescription: 'The amount of time specified by the <b>Input Dleay</b> setting.',
		newValue: APLValueInputDelay.create,
		fields: [],
	}),

	// Auras
	auraIsKnown: inputBuilder({
		label: 'Aura Known',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura is currently known, otherwise <b>False</b>.',
		newValue: APLValueAuraIsKnown.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit')],
	}),
	auraIsActive: inputBuilder({
		label: 'Aura Active',
		submenu: ['Aura'],
		shortDescription: '<b>True</b> if the aura is currently active, otherwise <b>False</b>.',
		newValue: APLValueAuraIsActive.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit')],
	}),
	auraIsActiveWithReactionTime: inputBuilder({
		label: 'Aura Active (with Reaction Time)',
		submenu: ['Aura'],
		shortDescription:
			'<b>True</b> if the aura is currently active AND it has been active for at least as long as the player reaction time (configured in Settings), otherwise <b>False</b>.',
		newValue: APLValueAuraIsActiveWithReactionTime.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit')],
	}),
	auraIsInactiveWithReactionTime: inputBuilder({
		label: 'Aura Inactive (with Reaction Time)',
		submenu: ['Aura'],
		shortDescription:
			'<b>True</b> if the aura is not currently active AND it has been inactive for at least as long as the player reaction time (configured in Settings), otherwise <b>False</b>.',
		newValue: APLValueAuraIsInactiveWithReactionTime.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit')],
	}),
	auraRemainingTime: inputBuilder({
		label: 'Aura Remaining Time',
		submenu: ['Aura'],
		shortDescription: 'Time remaining before this aura will expire, or 0 if the aura is not currently active.',
		newValue: APLValueAuraRemainingTime.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'auras', 'sourceUnit')],
	}),
	auraNumStacks: inputBuilder({
		label: 'Aura Num Stacks',
		submenu: ['Aura'],
		shortDescription: 'Number of stacks of the aura.',
		newValue: APLValueAuraNumStacks.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'stackable_auras', 'sourceUnit')],
	}),
	auraInternalCooldown: inputBuilder({
		label: 'Aura Remaining ICD',
		submenu: ['Aura'],
		shortDescription: "Time remaining before this aura's internal cooldown will be ready, or <b>0</b> if the ICD is ready now.",
		newValue: APLValueAuraInternalCooldown.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'icd_auras', 'sourceUnit')],
	}),
	auraIcdIsReadyWithReactionTime: inputBuilder({
		label: 'Aura ICD Is Ready (with Reaction Time)',
		submenu: ['Aura'],
		shortDescription:
			"<b>True</b> if the aura's ICD is currently ready OR it was put on CD recently, within the player's reaction time (configured in Settings), otherwise <b>False</b>.",
		newValue: APLValueAuraICDIsReadyWithReactionTime.create,
		fields: [AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources'), AplHelpers.actionIdFieldConfig('auraId', 'icd_auras', 'sourceUnit')],
	}),
	auraShouldRefresh: inputBuilder({
		label: 'Should Refresh Aura',
		submenu: ['Aura'],
		shortDescription: 'Whether this aura should be refreshed, e.g. for the purpose of maintaining a debuff.',
		fullDescription: `
		<p>This condition checks not only the specified aura but also any other auras on the same unit, including auras applied by other raid members, which apply the same debuff category.</p>
		<p>For example, 'Should Refresh Debuff(Sunder Armor)' will return <b>False</b> if the unit has an active Expose Armor aura.</p>
		`,
		newValue: () =>
			APLValueAuraShouldRefresh.create({
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
			AplHelpers.unitFieldConfig('sourceUnit', 'aura_sources_targets_first'),
			AplHelpers.actionIdFieldConfig('auraId', 'exclusive_effect_auras', 'sourceUnit', 'currentTarget'),
			valueFieldConfig('maxOverlap', {
				label: 'Overlap',
				labelTooltip: 'Maximum amount of time before the aura expires when it may be refreshed.',
			}),
		],
	}),

	// Aura Sets
	allTrinketStatProcsActive: inputBuilder({
		label: 'All Item Proc Buffs Active',
		submenu: ['Aura Sets'],
		shortDescription: '<b>True</b> if all item/enchant procs that buff the specified stat type(s) are currently active, otherwise <b>False</b>.',
		fullDescription: `
		<p>For stacking proc buffs, this condition also checks that the buff has been stacked to its maximum possible strength.</p>
		`,
		newValue: () =>
			APLValueAllTrinketStatProcsActive.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
			AplHelpers.minIcdInput,
		],
	}),
	anyTrinketStatProcsActive: inputBuilder({
		label: 'Any Item Proc Buff Active',
		submenu: ['Aura Sets'],
		shortDescription: '<b>True</b> if any item/enchant procs that buff the specified stat type(s) are currently active, otherwise <b>False</b>.',
		fullDescription: `
		<p>For stacking proc buffs, this condition also checks that the buff has been stacked to its maximum possible strength.</p>
		`,
		newValue: () =>
			APLValueAnyTrinketStatProcsActive.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
			AplHelpers.minIcdInput,
		],
	}),
	trinketProcsMinRemainingTime: inputBuilder({
		label: 'Item Procs Min Remaining Time',
		submenu: ['Aura Sets'],
		shortDescription:
			'Shortest remaining duration on any active item/enchant procs that buff the specified stat type(s), or infinity if none are currently active.',
		newValue: () =>
			APLValueTrinketProcsMinRemainingTime.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
			AplHelpers.minIcdInput,
		],
	}),
	trinketProcsMaxRemainingIcd: inputBuilder({
		label: 'Item Procs Max Remaining ICD',
		submenu: ['Aura Sets'],
		shortDescription: 'Longest remaining ICD on any inactive item/enchant procs that buff the specified stat type(s), or 0 if all are currently active.',
		newValue: () =>
			APLValueTrinketProcsMaxRemainingICD.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
			AplHelpers.minIcdInput,
		],
	}),
	numEquippedStatProcTrinkets: inputBuilder({
		label: 'Num Equipped Stat Proc Effects',
		submenu: ['Aura Sets'],
		shortDescription: 'Number of equipped passive item/enchant effects that buff the specified stat type(s) when they proc.',
		newValue: () =>
			APLValueNumEquippedStatProcTrinkets.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [
			AplHelpers.statTypeFieldConfig('statType1'),
			AplHelpers.statTypeFieldConfig('statType2'),
			AplHelpers.statTypeFieldConfig('statType3'),
			AplHelpers.minIcdInput,
		],
	}),
	numStatBuffCooldowns: inputBuilder({
		label: 'Num Stat Buff Cooldowns',
		submenu: ['Aura Sets'],
		shortDescription: 'Number of registered Major Cooldowns that buff the specified stat type(s) when they are cast.',
		fullDescription: `
		<p>Both manually casted cooldowns as well as cooldowns controlled by "Cast All Stat Buff Cooldowns" and "Autocast Other Cooldowns" actions are included in the total count returned by this value.</p>
		`,
		newValue: () =>
			APLValueNumStatBuffCooldowns.create({
				statType1: -1,
				statType2: -1,
				statType3: -1,
			}),
		fields: [AplHelpers.statTypeFieldConfig('statType1'), AplHelpers.statTypeFieldConfig('statType2'), AplHelpers.statTypeFieldConfig('statType3')],
	}),

	// DoT
	dotIsActive: inputBuilder({
		label: 'Dot Is Active',
		submenu: ['DoT'],
		shortDescription: '<b>True</b> if the specified dot is currently ticking, otherwise <b>False</b>.',
		newValue: APLValueDotIsActive.create,
		fields: [AplHelpers.unitFieldConfig('targetUnit', 'targets'), AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', '')],
	}),
	dotRemainingTime: inputBuilder({
		label: 'Dot Remaining Time',
		submenu: ['DoT'],
		shortDescription: 'Time remaining before the last tick of this DoT will occur, or 0 if the DoT is not currently ticking.',
		newValue: APLValueDotRemainingTime.create,
		fields: [AplHelpers.unitFieldConfig('targetUnit', 'targets'), AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', '')],
	}),
	dotTickFrequency: inputBuilder({
		label: 'Dot Tick Frequency',
		submenu: ['DoT'],
		shortDescription: 'The time between each tick.',
		newValue: APLValueDotTickFrequency.create,
		fields: [AplHelpers.unitFieldConfig('targetUnit', 'targets'), AplHelpers.actionIdFieldConfig('spellId', 'dot_spells', '')],
	}),
	dotPercentIncrease: inputBuilder({
		label: 'Dot Damage Increase %',
		submenu: ['DoT'],
		shortDescription: 'How much stronger a new DoT would be compared to the old.',
		newValue: APLValueDotPercentIncrease.create,
		fields: [AplHelpers.unitFieldConfig('targetUnit', 'targets'), AplHelpers.actionIdFieldConfig('spellId', 'expected_dot_spells', '')],
	}),
	sequenceIsComplete: inputBuilder({
		label: 'Sequence Is Complete',
		submenu: ['Sequence'],
		shortDescription: '<b>True</b> if there are no more subactions left to execute in the sequence, otherwise <b>False</b>.',
		newValue: APLValueSequenceIsComplete.create,
		fields: [AplHelpers.stringFieldConfig('sequenceName')],
	}),
	sequenceIsReady: inputBuilder({
		label: 'Sequence Is Ready',
		submenu: ['Sequence'],
		shortDescription: '<b>True</b> if the next subaction in the sequence is ready to be executed, otherwise <b>False</b>.',
		newValue: APLValueSequenceIsReady.create,
		fields: [AplHelpers.stringFieldConfig('sequenceName')],
	}),
	sequenceTimeToReady: inputBuilder({
		label: 'Sequence Time To Ready',
		submenu: ['Sequence'],
		shortDescription: 'Returns the amount of time remaining until the next subaction in the sequence will be ready.',
		newValue: APLValueSequenceTimeToReady.create,
		fields: [AplHelpers.stringFieldConfig('sequenceName')],
	}),

	// Class/spec specific values
	totemRemainingTime: inputBuilder({
		label: 'Totem Remaining Time',
		submenu: ['Shaman'],
		shortDescription: 'Returns the amount of time remaining until the totem will expire.',
		newValue: APLValueTotemRemainingTime.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassShaman,
		fields: [totemTypeFieldConfig('totemType')],
	}),
	shamanFireElementalDuration: inputBuilder({
		label: 'Fire Elemental Total Duration',
		submenu: ['Shaman'],
		shortDescription: 'Returns the duration of Fire Elemental depending on if Totemic Focus is talented or not.',
		newValue: APLValueShamanFireElementalDuration.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getClass() == Class.ClassShaman,
		fields: [],
	}),
	catExcessEnergy: inputBuilder({
		label: 'Excess Energy',
		submenu: ['Feral Druid'],
		shortDescription: 'Returns the amount of excess energy available, after subtracting energy that will be needed to maintain DoTs.',
		newValue: APLValueCatExcessEnergy.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecFeralDruid,
		fields: [],
	}),
	catNewSavageRoarDuration: inputBuilder({
		label: 'New Savage Roar Duration',
		submenu: ['Feral Druid'],
		shortDescription: 'Returns duration of savage roar based on current combo points',
		newValue: APLValueCatNewSavageRoarDuration.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecFeralDruid,
		fields: [],
	}),
	warlockHandOfGuldanInFlight: inputBuilder({
		label: 'Hand of Guldan in Flight',
		submenu: ['Warlock'],
		shortDescription: 'Returns <b>True</b> if the impact of Hand of Guldan currenty is in flight.',
		newValue: APLValueWarlockHandOfGuldanInFlight.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecDemonologyWarlock,
		fields: [],
	}),
	warlockHauntInFlight: inputBuilder({
		label: 'Haunt In Flight',
		submenu: ['Warlock'],
		shortDescription: 'Returns <b>True</b> if Haunt currently is in flight.',
		newValue: APLValueWarlockHauntInFlight.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecAfflictionWarlock,
		fields: [],
	}),
	mageCurrentCombustionDotEstimate: inputBuilder({
		label: 'Combustion Dot Value',
		submenu: ['Mage'],
		shortDescription: 'Returns the current estimated size of your Combustion Dot.',
		newValue: APLValueMageCurrentCombustionDotEstimate.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() == Spec.SpecFireMage,
		fields: [],
	}),
	brewmasterMonkCurrentStaggerPercent: inputBuilder({
		label: 'Current Stagger (%)',
		submenu: ['Tank'],
		shortDescription: 'Amount of current Stagger, as a percentage.',
		newValue: APLValueMonkCurrentChi.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() === Spec.SpecBrewmasterMonk,
		fields: [],
	}),
	protectionPaladinDamageTakenLastGlobal: inputBuilder({
		label: 'Damage Taken Last Global',
		submenu: ['Tank'],
		shortDescription: 'Amount of damage taken in the last 1.5s.',
		newValue: APLValueProtectionPaladinDamageTakenLastGlobal.create,
		includeIf: (player: Player<any>, _isPrepull: boolean) => player.getSpec() === Spec.SpecProtectionPaladin,
		fields: [],
	}),
};
