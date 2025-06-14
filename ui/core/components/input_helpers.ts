import { Player } from '../player.js';
import { Spec } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { ClassOptions, SpecOptions, SpecRotation } from '../proto_utils/utils.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { formatToNumber, randomUUID } from '../utils';
import { BooleanPickerConfig } from './pickers/boolean_picker.js';
import { EnumPickerConfig, EnumValueConfig } from './pickers/enum_picker.js';
import { IconEnumPickerConfig, IconEnumValueConfig } from './pickers/icon_enum_picker.jsx';
import { IconPickerConfig } from './pickers/icon_picker.jsx';
import { MultiIconPickerConfig } from './pickers/multi_icon_picker.js';
import { NumberPickerConfig } from './pickers/number_picker.js';

export function makeMultiIconInput<ModObject>(
	inputs: Array<IconPickerConfig<ModObject, any>>,
	label: string,
	categoryId?: ActionId,
): MultiIconPickerConfig<ModObject> {
	return {
		inputs: inputs,
		label: label,
		categoryId: categoryId,
		showWhen: p => inputs.filter(i => !i.showWhen || i.showWhen(p as ModObject)).length > 0,
	};
}

// Extend this to add player callbacks as optional config fields.
interface BasePlayerConfig<SpecType extends Spec, T> {
	getValue?: (player: Player<SpecType>) => T;
	setValue?: (eventID: EventID, player: Player<SpecType>, newVal: T) => void;
	changeEmitter?: (player: Player<SpecType>) => TypedEvent<any>;
	extraCssClasses?: Array<string>;
	showWhen?: (player: Player<SpecType>) => boolean;
}

/////////////////////////////////////////////////////////////////////////////////
//                                    BOOLEAN
/////////////////////////////////////////////////////////////////////////////////
export interface TypedBooleanPickerConfig<ModObject> extends BooleanPickerConfig<ModObject> {
	type: 'boolean';
}

interface WrappedBooleanInputConfig<SpecType extends Spec, ModObject> extends BooleanPickerConfig<ModObject> {
	getModObject: (player: Player<SpecType>) => ModObject;
}
function makeWrappedBooleanInput<SpecType extends Spec, ModObject>(
	config: WrappedBooleanInputConfig<SpecType, ModObject>,
): TypedBooleanPickerConfig<Player<SpecType>> {
	const getModObject = config.getModObject;
	return {
		id: config.id,
		type: 'boolean',
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		changedEvent: (player: Player<SpecType>) => config.changedEvent(getModObject(player)),
		getValue: (player: Player<SpecType>) => config.getValue(getModObject(player)),
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => config.setValue(eventID, getModObject(player), newValue),
		enableWhen: config.enableWhen ? (player: Player<SpecType>) => config.enableWhen!(getModObject(player)) : undefined,
		showWhen: config.showWhen ? (player: Player<SpecType>) => config.showWhen!(getModObject(player)) : undefined,
		extraCssClasses: config.extraCssClasses,
	};
}
export interface PlayerBooleanInputConfig<SpecType extends Spec, Message> extends BasePlayerConfig<SpecType, boolean> {
	fieldName: keyof Message;
	label: string;
	labelTooltip?: string;
	description?: string | Element;
	enableWhen?: (player: Player<SpecType>) => boolean;
	showWhen?: (player: Player<SpecType>) => boolean;
}
export function makeClassOptionsBooleanInput<SpecType extends Spec>(
	config: PlayerBooleanInputConfig<SpecType, ClassOptions<SpecType>>,
): TypedBooleanPickerConfig<Player<SpecType>> {
	return makeWrappedBooleanInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getClassOptions()[config.fieldName] as unknown as boolean),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: boolean) => {
				const newMessage = player.getClassOptions();
				(newMessage[config.fieldName] as unknown as boolean) = newVal;
				player.setClassOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	});
}
export function makeSpecOptionsBooleanInput<SpecType extends Spec>(
	config: PlayerBooleanInputConfig<SpecType, SpecOptions<SpecType>>,
): TypedBooleanPickerConfig<Player<SpecType>> {
	return makeWrappedBooleanInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSpecOptions()[config.fieldName] as unknown as boolean),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: boolean) => {
				const newMessage = player.getSpecOptions();
				(newMessage[config.fieldName] as unknown as boolean) = newVal;
				player.setSpecOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	});
}
export function makeRotationBooleanInput<SpecType extends Spec>(
	config: PlayerBooleanInputConfig<SpecType, SpecRotation<SpecType>>,
): TypedBooleanPickerConfig<Player<SpecType>> {
	return makeWrappedBooleanInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSimpleRotation()[config.fieldName] as unknown as boolean),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: boolean) => {
				const newMessage = player.getSimpleRotation();
				(newMessage[config.fieldName] as unknown as boolean) = newVal;
				player.setSimpleRotation(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.rotationChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	});
}

/////////////////////////////////////////////////////////////////////////////////
//                                    NUMBER
/////////////////////////////////////////////////////////////////////////////////
export interface TypedNumberPickerConfig<ModObject> extends NumberPickerConfig<ModObject> {
	type: 'number';
}

interface WrappedNumberInputConfig<SpecType extends Spec, ModObject> extends NumberPickerConfig<ModObject> {
	getModObject: (player: Player<SpecType>) => ModObject;
}
function makeWrappedNumberInput<SpecType extends Spec, ModObject>(
	config: WrappedNumberInputConfig<SpecType, ModObject>,
): TypedNumberPickerConfig<Player<SpecType>> {
	const getModObject = config.getModObject;
	return {
		id: config.id,
		type: 'number',
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		float: config.float,
		showZeroes: config.showZeroes,
		maxDecimalDigits: config.maxDecimalDigits,
		positive: config.positive,
		changedEvent: (player: Player<SpecType>) => config.changedEvent(getModObject(player)),
		getValue: (player: Player<SpecType>) => config.getValue(getModObject(player)),
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: number) => config.setValue(eventID, getModObject(player), newValue),
		enableWhen: config.enableWhen ? (player: Player<SpecType>) => config.enableWhen!(getModObject(player)) : undefined,
		showWhen: config.showWhen ? (player: Player<SpecType>) => config.showWhen!(getModObject(player)) : undefined,
		extraCssClasses: config.extraCssClasses,
	};
}
export interface PlayerNumberInputConfig<SpecType extends Spec, Message>
	extends BasePlayerConfig<SpecType, number>,
		Pick<NumberPickerConfig<Player<SpecType>>, 'labelTooltip' | 'description' | 'showZeroes' | 'maxDecimalDigits' | 'float' | 'positive'> {
	fieldName: keyof Message;
	label: string;
	percent?: boolean;
	max?: number;
	enableWhen?: (player: Player<SpecType>) => boolean;
	showWhen?: (player: Player<SpecType>) => boolean;
	changeEmitter?: (player: Player<SpecType>) => TypedEvent<any>;
}

export const numberInputValueToPercentage = (value: number, config: PlayerNumberInputConfig<any, any>) =>
	Number(formatToNumber(value / 100, { maximumFractionDigits: config.maxDecimalDigits, useGrouping: false }));

export const numberInputPercentToValue = (value: number, config: PlayerNumberInputConfig<any, any>) =>
	Number(formatToNumber(value * 100, { maximumFractionDigits: config.maxDecimalDigits, useGrouping: false }));

export function makeClassOptionsNumberInput<SpecType extends Spec>(
	config: PlayerNumberInputConfig<SpecType, ClassOptions<SpecType>>,
): TypedNumberPickerConfig<Player<SpecType>> {
	const internalConfig = {
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		float: config.float,
		showZeroes: config.showZeroes,
		maxDecimalDigits: config.maxDecimalDigits,
		positive: config.positive,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getClassOptions()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getClassOptions();
				if (config?.max && newVal > config.max) {
					newVal = config.max;
				}
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setClassOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	};
	if (config.percent) {
		const getValue = internalConfig.getValue;
		internalConfig.getValue = (player: Player<SpecType>) => numberInputPercentToValue(getValue(player), config);
		const setValue = internalConfig.setValue;
		internalConfig.setValue = (eventID: EventID, player: Player<SpecType>, newVal: number) =>
			setValue(eventID, player, numberInputValueToPercentage(newVal, config));
	}
	return makeWrappedNumberInput<SpecType, Player<SpecType>>(internalConfig);
}
export function makeSpecOptionsNumberInput<SpecType extends Spec>(
	config: PlayerNumberInputConfig<SpecType, SpecOptions<SpecType>>,
): TypedNumberPickerConfig<Player<SpecType>> {
	const internalConfig = {
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		float: config.float,
		showZeroes: config.showZeroes,
		maxDecimalDigits: config.maxDecimalDigits,
		positive: config.positive,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSpecOptions()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getSpecOptions();
				if (config?.max && newVal > config.max) {
					newVal = config.max;
				}
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setSpecOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	};

	if (config.percent) {
		const getValue = internalConfig.getValue;
		internalConfig.getValue = (player: Player<SpecType>) => numberInputPercentToValue(getValue(player), config);
		const setValue = internalConfig.setValue;
		internalConfig.setValue = (eventID: EventID, player: Player<SpecType>, newVal: number) =>
			setValue(eventID, player, numberInputValueToPercentage(newVal, config));
	}
	return makeWrappedNumberInput<SpecType, Player<SpecType>>(internalConfig);
}
export function makeRotationNumberInput<SpecType extends Spec>(
	config: PlayerNumberInputConfig<SpecType, SpecRotation<SpecType>>,
): TypedNumberPickerConfig<Player<SpecType>> {
	const internalConfig = {
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		float: config.float,
		showZeroes: config.showZeroes,
		positive: config.positive,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSimpleRotation()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getSimpleRotation();
				if (config?.max && newVal > config.max) {
					newVal = config.max;
				}
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setSimpleRotation(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.rotationChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
		extraCssClasses: config.extraCssClasses,
	};
	if (config.percent) {
		const getValue = internalConfig.getValue;
		internalConfig.getValue = (player: Player<SpecType>) => numberInputPercentToValue(getValue(player), config);
		const setValue = internalConfig.setValue;
		internalConfig.setValue = (eventID: EventID, player: Player<SpecType>, newVal: number) =>
			setValue(eventID, player, numberInputValueToPercentage(newVal, config));
	}
	return makeWrappedNumberInput<SpecType, Player<SpecType>>(internalConfig);
}

/////////////////////////////////////////////////////////////////////////////////
//                                    ENUM
/////////////////////////////////////////////////////////////////////////////////
export interface TypedEnumPickerConfig<ModObject> extends EnumPickerConfig<ModObject> {
	type: 'enum';
}

interface WrappedEnumInputConfig<SpecType extends Spec, ModObject> extends EnumPickerConfig<ModObject> {
	getModObject: (player: Player<SpecType>) => ModObject;
}
function makeWrappedEnumInput<SpecType extends Spec, ModObject>(config: WrappedEnumInputConfig<SpecType, ModObject>): TypedEnumPickerConfig<Player<SpecType>> {
	const getModObject = config.getModObject;
	return {
		id: config.id,
		type: 'enum',
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		values: config.values,
		changedEvent: (player: Player<SpecType>) => config.changedEvent(getModObject(player)),
		getValue: (player: Player<SpecType>) => config.getValue(getModObject(player)),
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: number) => config.setValue(eventID, getModObject(player), newValue),
		enableWhen: config.enableWhen ? (player: Player<SpecType>) => config.enableWhen!(getModObject(player)) : undefined,
		showWhen: config.showWhen ? (player: Player<SpecType>) => config.showWhen!(getModObject(player)) : undefined,
	};
}

export interface PlayerEnumInputConfig<SpecType extends Spec, Message> {
	fieldName: keyof Message;
	label: string;
	labelTooltip?: string;
	description?: string | Element;
	values: Array<EnumValueConfig>;
	getValue?: (player: Player<SpecType>) => number;
	setValue?: (eventID: EventID, player: Player<SpecType>, newValue: number) => void;
	enableWhen?: (player: Player<SpecType>) => boolean;
	showWhen?: (player: Player<SpecType>) => boolean;
	changeEmitter?: (player: Player<SpecType>) => TypedEvent<any>;
}
// T is unused, but kept to have the same interface as the icon enum inputs.
export function makeClassOptionsEnumInput<SpecType extends Spec, _T>(
	config: PlayerEnumInputConfig<SpecType, ClassOptions<SpecType>>,
): TypedEnumPickerConfig<Player<SpecType>> {
	return makeWrappedEnumInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		values: config.values,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getClassOptions()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getClassOptions();
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setClassOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
	});
}
// T is unused, but kept to have the same interface as the icon enum inputs.
export function makeSpecOptionsEnumInput<SpecType extends Spec, _T>(
	config: PlayerEnumInputConfig<SpecType, SpecOptions<SpecType>>,
): TypedEnumPickerConfig<Player<SpecType>> {
	return makeWrappedEnumInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		values: config.values,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSpecOptions()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getSpecOptions();
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setSpecOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
	});
}
// T is unused, but kept to have the same interface as the icon enum inputs.
export function makeRotationEnumInput<SpecType extends Spec, _T>(
	config: PlayerEnumInputConfig<SpecType, SpecRotation<SpecType>>,
): TypedEnumPickerConfig<Player<SpecType>> {
	return makeWrappedEnumInput<SpecType, Player<SpecType>>({
		id: `${String(config.fieldName) || randomUUID()}`,
		label: config.label,
		labelTooltip: config.labelTooltip,
		description: config.description,
		values: config.values,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSimpleRotation()[config.fieldName] as unknown as number),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: number) => {
				const newMessage = player.getSimpleRotation();
				(newMessage[config.fieldName] as unknown as number) = newVal;
				player.setSimpleRotation(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.rotationChangeEmitter),
		enableWhen: config.enableWhen,
		showWhen: config.showWhen,
	});
}

/////////////////////////////////////////////////////////////////////////////////
//                                  ICON
/////////////////////////////////////////////////////////////////////////////////
export interface TypedIconPickerConfig<ModObject, T> extends IconPickerConfig<ModObject, T> {
	type: 'icon';
}

interface WrappedIconInputConfig<SpecType extends Spec, ModObject, T> extends IconPickerConfig<ModObject, T> {
	getModObject: (player: Player<SpecType>) => ModObject;
}
function makeWrappedIconInput<SpecType extends Spec, ModObject, T>(
	config: WrappedIconInputConfig<SpecType, ModObject, T>,
): TypedIconPickerConfig<Player<SpecType>, T> {
	const getModObject = config.getModObject;
	return {
		type: 'icon',
		actionId: config.actionId,
		states: config.states,
		changedEvent: (player: Player<SpecType>) => config.changedEvent(getModObject(player)),
		showWhen: (player: Player<SpecType>) => !config.showWhen || (config.showWhen(getModObject(player)) as any),
		getValue: (player: Player<SpecType>) => config.getValue(getModObject(player)),
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: T) => config.setValue(eventID, getModObject(player), newValue),
		extraCssClasses: config.extraCssClasses,
	};
}

interface WrappedTypedInputConfig<Message, ModObject, T> {
	getModObject: (player: Player<any>) => ModObject;
	getValue: (modObj: ModObject) => Message;
	setValue: (eventID: EventID, modObj: ModObject, messageVal: Message) => void;
	changeEmitter: (modObj: ModObject) => TypedEvent<any>;
	extraCssClasses?: Array<string>;

	showWhen?: (obj: ModObject) => boolean;
	getFieldValue?: (modObj: ModObject) => T;
	setFieldValue?: (eventID: EventID, modObj: ModObject, newValue: T) => void;
}

export function makeBooleanIconInput<SpecType extends Spec, Message, ModObject>(
	config: WrappedTypedInputConfig<Message, ModObject, boolean>,
	actionId: ActionId,
	fieldName: keyof Message,
	value?: number,
): TypedIconPickerConfig<Player<SpecType>, boolean> {
	return makeWrappedIconInput<SpecType, ModObject, boolean>({
		getModObject: config.getModObject,
		actionId: actionId,
		states: 2,
		changedEvent: config.changeEmitter,
		showWhen: config.showWhen,
		getValue:
			config.getFieldValue ||
			((modObj: ModObject) =>
				value ? (config.getValue(modObj)[fieldName] as unknown as number) == value : (config.getValue(modObj)[fieldName] as unknown as boolean)),
		setValue:
			config.setFieldValue ||
			((eventID: EventID, modObj: ModObject, newValue: boolean) => {
				const newMessage = config.getValue(modObj);
				if (value) {
					if (newValue) {
						(newMessage[fieldName] as unknown as number) = value;
					} else if ((newMessage[fieldName] as unknown as number) == value) {
						(newMessage[fieldName] as unknown as number) = 0;
					}
				} else {
					(newMessage[fieldName] as unknown as boolean) = newValue;
				}
				config.setValue(eventID, modObj, newMessage);
			}),
		extraCssClasses: config.extraCssClasses,
	});
}

export interface PlayerBooleanIconInputConfig<SpecType extends Spec, Message, T> extends BasePlayerConfig<SpecType, T> {
	fieldName: keyof Message;
	id: ActionId;
	value?: number;
}
export function makeClassOptionsBooleanIconInput<SpecType extends Spec>(
	config: PlayerBooleanIconInputConfig<SpecType, ClassOptions<SpecType>, boolean>,
): TypedIconPickerConfig<Player<SpecType>, boolean> {
	return makeBooleanIconInput<SpecType, ClassOptions<SpecType>, Player<SpecType>>(
		{
			getModObject: (player: Player<SpecType>) => player,
			getValue: (player: Player<SpecType>) => player.getClassOptions(),
			setValue: (eventID: EventID, player: Player<SpecType>, newVal: ClassOptions<SpecType>) => player.setClassOptions(eventID, newVal),
			changeEmitter: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
			extraCssClasses: config.extraCssClasses,
			getFieldValue: config.getValue,
			setFieldValue: config.setValue,
			showWhen: config.showWhen,
		},
		config.id,
		config.fieldName,
		config.value,
	);
}
export function makeSpecOptionsBooleanIconInput<SpecType extends Spec>(
	config: PlayerBooleanIconInputConfig<SpecType, SpecOptions<SpecType>, boolean>,
): TypedIconPickerConfig<Player<SpecType>, boolean> {
	return makeBooleanIconInput<SpecType, SpecOptions<SpecType>, Player<SpecType>>(
		{
			getModObject: (player: Player<SpecType>) => player,
			getValue: (player: Player<SpecType>) => player.getSpecOptions(),
			setValue: (eventID: EventID, player: Player<SpecType>, newVal: SpecOptions<SpecType>) => player.setSpecOptions(eventID, newVal),
			changeEmitter: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
			extraCssClasses: config.extraCssClasses,
			getFieldValue: config.getValue,
			setFieldValue: config.setValue,
			showWhen: config.showWhen,
		},
		config.id,
		config.fieldName,
		config.value,
	);
}

function makeNumberIconInput<SpecType extends Spec, Message, ModObject>(
	config: WrappedTypedInputConfig<Message, ModObject, number>,
	actionId: ActionId,
	fieldName: keyof Message,
	multiplier?: number,
): TypedIconPickerConfig<Player<SpecType>, number> {
	return makeWrappedIconInput<SpecType, ModObject, number>({
		getModObject: config.getModObject,
		actionId: actionId,
		states: 0, // Must be assigned externally.
		changedEvent: config.changeEmitter,
		getValue: (modObj: ModObject) => config.getValue(modObj)[fieldName] as unknown as number,
		setValue: (eventID: EventID, modObj: ModObject, newValue: number) => {
			const newMessage = config.getValue(modObj);
			if (multiplier) {
				const sign = newValue - (newMessage[fieldName] as unknown as number);
				newValue += (multiplier - 1) * sign;
			}
			if (newValue < 0) {
				newValue = 0;
			}
			(newMessage[fieldName] as unknown as number) = newValue;
			config.setValue(eventID, modObj, newMessage);
		},
	});
}
export function makeTristateIconInput<SpecType extends Spec, Message, ModObject>(
	config: WrappedTypedInputConfig<Message, ModObject, number>,
	id: ActionId,
	impId: ActionId,
	fieldName: keyof Message,
): TypedIconPickerConfig<Player<SpecType>, number> {
	const input = makeNumberIconInput<SpecType, Message, ModObject>(config, id, fieldName);
	input.states = 3;
	input.improvedId = impId;
	return input;
}
export function makeQuadstateIconInput<SpecType extends Spec, Message, ModObject>(
	config: WrappedTypedInputConfig<Message, ModObject, number>,
	id: ActionId,
	impId: ActionId,
	impId2: ActionId,
	fieldName: keyof Message,
): TypedIconPickerConfig<Player<SpecType>, number> {
	const input = makeNumberIconInput<SpecType, Message, ModObject>(config, id, fieldName);
	input.states = 4;
	input.improvedId = impId;
	input.improvedId2 = impId2;
	return input;
}
export function makeMultistateIconInput<SpecType extends Spec, Message, ModObject>(
	config: WrappedTypedInputConfig<Message, ModObject, number>,
	id: ActionId,
	numStates: number,
	fieldName: keyof Message,
	multiplier?: number,
): TypedIconPickerConfig<Player<SpecType>, number> {
	const input = makeNumberIconInput<SpecType, Message, ModObject>(config, id, fieldName, multiplier);
	input.states = numStates;
	return input;
}

export interface TypedIconEnumPickerConfig<ModObject, T> extends IconEnumPickerConfig<ModObject, T> {
	type: 'iconEnum';
}

interface WrappedEnumIconInputConfig<SpecType extends Spec, ModObject, T> extends IconEnumPickerConfig<ModObject, T> {
	getModObject: (player: Player<SpecType>) => ModObject;
}
function makeWrappedEnumIconInput<SpecType extends Spec, ModObject, T>(
	config: WrappedEnumIconInputConfig<SpecType, ModObject, T>,
): TypedIconEnumPickerConfig<Player<SpecType>, T> {
	const getModObject = config.getModObject;
	return {
		type: 'iconEnum',
		numColumns: config.numColumns,
		values: config.values.map(value => {
			if (value.showWhen) {
				const showWhen = value.showWhen;
				value.showWhen = ((player: Player<SpecType>) => showWhen(getModObject(player))) as any;
			}
			return value as unknown as IconEnumValueConfig<Player<SpecType>, T>;
		}),
		equals: config.equals,
		showWhen: (player: Player<SpecType>): boolean => !config.showWhen || (config.showWhen(getModObject(player)) as any),
		zeroValue: config.zeroValue,
		changedEvent: (player: Player<SpecType>) => config.changedEvent(getModObject(player)),
		getValue: (player: Player<SpecType>) => config.getValue(getModObject(player)),
		setValue: (eventID: EventID, player: Player<SpecType>, newValue: T) => config.setValue(eventID, getModObject(player), newValue),
		extraCssClasses: config.extraCssClasses,
	};
}

export interface PlayerEnumIconInputConfig<SpecType extends Spec, Message, T> extends BasePlayerConfig<SpecType, T> {
	fieldName: keyof Message;
	values: Array<IconEnumValueConfig<Player<SpecType>, T>>;
	numColumns?: number;
}
export function makeClassOptionsEnumIconInput<SpecType extends Spec, T>(
	config: PlayerEnumIconInputConfig<SpecType, ClassOptions<SpecType>, T>,
): TypedIconEnumPickerConfig<Player<SpecType>, T> {
	return makeWrappedEnumIconInput<SpecType, Player<SpecType>, T>({
		numColumns: config.numColumns || 1,
		values: config.values,
		equals: (a: T, b: T) => a == b,
		showWhen: config.showWhen,
		zeroValue: 0 as unknown as T,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getClassOptions()[config.fieldName] as unknown as T),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: T) => {
				const newMessage = player.getClassOptions();
				(newMessage[config.fieldName] as unknown as T) = newVal;
				player.setClassOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		extraCssClasses: config.extraCssClasses,
	});
}
export function makeSpecOptionsEnumIconInput<SpecType extends Spec, T>(
	config: PlayerEnumIconInputConfig<SpecType, SpecOptions<SpecType>, T>,
): TypedIconEnumPickerConfig<Player<SpecType>, T> {
	return makeWrappedEnumIconInput<SpecType, Player<SpecType>, T>({
		numColumns: config.numColumns || 1,
		values: config.values,
		equals: (a: T, b: T) => a == b,
		showWhen: config.showWhen,
		zeroValue: 0 as unknown as T,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSpecOptions()[config.fieldName] as unknown as T),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: T) => {
				const newMessage = player.getSpecOptions();
				(newMessage[config.fieldName] as unknown as T) = newVal;
				player.setSpecOptions(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.specOptionsChangeEmitter),
		extraCssClasses: config.extraCssClasses,
	});
}
export function makeRotationEnumIconInput<SpecType extends Spec, T>(
	config: PlayerEnumIconInputConfig<SpecType, SpecRotation<SpecType>, T>,
): TypedIconEnumPickerConfig<Player<SpecType>, T> {
	return makeWrappedEnumIconInput<SpecType, Player<SpecType>, T>({
		numColumns: config.numColumns || 1,
		values: config.values,
		equals: (a: T, b: T) => a == b,
		showWhen: config.showWhen,
		zeroValue: 0 as unknown as T,
		getModObject: (player: Player<SpecType>) => player,
		getValue: config.getValue || ((player: Player<SpecType>) => player.getSimpleRotation()[config.fieldName] as unknown as T),
		setValue:
			config.setValue ||
			((eventID: EventID, player: Player<SpecType>, newVal: T) => {
				const newMessage = player.getSimpleRotation();
				(newMessage[config.fieldName] as unknown as T) = newVal;
				player.setSimpleRotation(eventID, newMessage);
			}),
		changedEvent: config.changeEmitter || ((player: Player<SpecType>) => player.rotationChangeEmitter),
		extraCssClasses: config.extraCssClasses,
	});
}
