import clsx from 'clsx';
import { element } from 'tsx-vanilla';

import { UnitReference } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { DropdownPicker, DropdownPickerConfig, DropdownValueConfig } from './dropdown_picker.jsx';

export interface UnitValue {
	value: UnitReference | undefined;
	text?: string;
	iconUrl?: string | ActionId;
	color?: string;
}

export interface UnitValueConfig extends DropdownValueConfig<UnitValue> {}
export interface UnitPickerConfig<ModObject>
	extends Omit<DropdownPickerConfig<ModObject, UnitReference | undefined, UnitValue>, 'equals' | 'setOptionContent' | 'defaultLabel'> {
	hideLabelWhenDefaultSelected?: boolean;
}

export class UnitPicker<ModObject> extends DropdownPicker<ModObject, UnitReference | undefined, UnitValue> {
	constructor(parent: HTMLElement, modObject: ModObject, config: UnitPickerConfig<ModObject>) {
		super(parent, modObject, {
			...config,
			equals: (a, b) => UnitReference.equals(a?.value || UnitReference.create(), b?.value || UnitReference.create()),
			defaultLabel: 'Unit',
			setOptionContent: (button, valueConfig, isSelectButton) => {
				const unitConfig = valueConfig.value;

				button.className = button.className.replace(/text-[\w]*/, '');
				if (unitConfig.color) {
					button.classList.add(`text-${unitConfig.color}`);
				}

				if (unitConfig.iconUrl) {
					let icon: HTMLElement | HTMLImageElement | null = null;
					if (unitConfig.iconUrl instanceof ActionId) {
						icon = (<img className="unit-picker-item-icon" />) as HTMLImageElement;
						unitConfig.iconUrl.fill().then(filledId => {
							if (icon) (icon as HTMLImageElement).src! = filledId.iconUrl;
						});
					} else if (unitConfig.iconUrl.startsWith('fa-')) {
						icon = (<i className={clsx('fa', unitConfig.iconUrl, 'unit-picker-item-icon')} />) as HTMLElement;
					} else {
						icon = (<img className="unit-picker-item-icon" src={unitConfig.iconUrl} />) as HTMLImageElement;
					}
					button.appendChild(icon);
				}

				const hideLabel = config.hideLabelWhenDefaultSelected && isSelectButton && !unitConfig.value;
				if (unitConfig.text && !hideLabel) {
					button.insertAdjacentText('beforeend', unitConfig.text);
				}
			},
		});
		this.rootElem.classList.add('unit-picker-root');
	}
}
