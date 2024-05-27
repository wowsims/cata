import { Dropdown } from 'bootstrap';
import clsx from 'clsx';
import tippy, { Instance as TippyInstance } from 'tippy.js';
import { element, fragment, ref } from 'tsx-vanilla';

import { TypedEvent } from '../typed_event.js';
import { existsInDOM } from '../utils';
import { Input, InputConfig } from './input.js';

export interface DropdownValueConfig<V> {
	value: V;
	submenu?: (string | V)[];
	headerText?: string;
	tooltip?: string;
	extraCssClasses?: string[];
}

export interface DropdownPickerConfig<ModObject, T, V = T> extends InputConfig<ModObject, T, V> {
	id: string;
	values: DropdownValueConfig<V>[];
	equals: (a: V | undefined, b: V | undefined) => boolean;
	setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<V>, isSelectButton?: boolean) => void;
	createMissingValue?: (val: V) => Promise<DropdownValueConfig<V>>;
	defaultLabel: string;
}

interface DropdownSubmenu<V> {
	path: (string | V)[];
	listElem: HTMLUListElement;
}

/** UI Input that uses a dropdown menu. */
export class DropdownPicker<ModObject, T, V = T> extends Input<ModObject, T, V> {
	private readonly config: DropdownPickerConfig<ModObject, T, V>;
	private valueConfigs: DropdownValueConfig<V>[];

	private readonly buttonElem: HTMLButtonElement;
	private readonly listElem: HTMLUListElement;
	private tooltip: TippyInstance | null = null;

	private currentSelection: DropdownValueConfig<V> | null;
	private submenus: Array<DropdownSubmenu<V>>;

	constructor(parent: HTMLElement, modObject: ModObject, config: DropdownPickerConfig<ModObject, T, V>) {
		super(parent, 'dropdown-picker-root', modObject, config);
		this.config = config;
		this.valueConfigs = this.config.values.filter(vc => !vc.headerText);
		this.currentSelection = null;
		this.submenus = [];

		this.rootElem.classList.add('dropdown');

		const buttonRef = ref<HTMLButtonElement>();
		const listRef = ref<HTMLUListElement>();
		this.rootElem.appendChild(
			<>
				<button
					ref={buttonRef}
					id={config.id}
					className="dropdown-picker-button btn dropdown-toggle open-on-click"
					dataset={{ bsToggle: 'dropdown' }}
					attributes={{ 'aria-expanded': false }}>
					{config.defaultLabel}
				</button>
				<ul ref={listRef} className="dropdown-picker-list dropdown-menu"></ul>
			</>,
		);

		this.buttonElem = buttonRef.value!;
		this.listElem = listRef.value!;

		this.buttonElem.addEventListener(
			'show.bs.dropdown',
			() => {
				this.buildDropdown(this.valueConfigs);
			},
			{ signal: this.signal },
		);
		this.buttonElem.addEventListener(
			'hidden.bs.dropdown',
			() => {
				this.clearDropdownInstances();
				this.listElem.replaceChildren(<></>);
			},
			{ signal: this.signal },
		);

		this.init();

		this.addOnDisposeCallback(() => {
			this.clearDropdownInstances();
			this.listElem.remove();
			Dropdown.getOrCreateInstance(this.buttonElem).dispose();
			this.buttonElem.remove();
		});
	}

	clearDropdownInstances() {
		this.listElem.querySelectorAll('[data-bs-toggle=dropdown]').forEach(elem => Dropdown.getOrCreateInstance(elem).dispose());
	}

	setOptions(newValueConfigs: DropdownValueConfig<V>[]) {
		const roomExistsInDOM = existsInDOM(this.rootElem);
		const listExistsInDOM = existsInDOM(this.listElem);
		const buttonExistsInDOM = existsInDOM(this.buttonElem);
		this.clearDropdownInstances();

		if (!roomExistsInDOM || !buttonExistsInDOM || !listExistsInDOM) {
			this.dispose();
			return;
		}

		this.valueConfigs = newValueConfigs.filter(vc => !vc.headerText);
		this.setInputValue(this.getSourceValue());
		return;
	}

	private buildDropdown(valueConfigs: DropdownValueConfig<V>[]) {
		this.listElem.replaceChildren();
		this.submenus = [];
		valueConfigs.forEach(valueConfig => {
			const containsSubmenuChildren = valueConfigs.some(vc => vc.submenu?.some(e => !(typeof e == 'string') && this.config.equals(e, valueConfig.value)));
			const buttonRef = ref<HTMLButtonElement>();
			const listItemRef = ref<HTMLLIElement>();
			const itemElem = (
				<li ref={listItemRef} className={clsx(valueConfig.extraCssClasses, valueConfig.headerText ? 'dropdown-picker-header' : 'dropdown-picker-item')}>
					{valueConfig.headerText && <h6 className="dropdown-header">{valueConfig.headerText}</h6>}
				</li>
			);

			if (!valueConfig.headerText) {
				const buttonElem = <button ref={buttonRef} className="dropdown-item" />;
				this.config.setOptionContent(buttonRef.value!, valueConfig);

				if (valueConfig.tooltip) {
					const tooltip = tippy(buttonRef.value!, {
						animation: false,
						theme: 'dropdown-tooltip',
						content: valueConfig.tooltip,
					});
					this.addOnDisposeCallback(() => tooltip?.destroy());
				}

				buttonRef.value!.addEventListener(
					'click',
					() => {
						this.updateValue(valueConfig);
						this.inputChanged(TypedEvent.nextEventID());
					},
					{ signal: this.signal },
				);
				this.addOnDisposeCallback(() => {
					buttonRef.value?.remove();
					itemElem.remove();
				});

				if (containsSubmenuChildren) {
					this.createSubmenu((valueConfig.submenu || []).concat([valueConfig.value]), buttonRef.value!, listItemRef.value!);
				} else {
					itemElem.appendChild(buttonElem);
				}
			}

			if (!containsSubmenuChildren) {
				if (valueConfig.submenu && valueConfig.submenu.length > 0) {
					this.createSubmenu(valueConfig.submenu);
				}
				const submenu = this.getSubmenu(valueConfig.submenu);
				if (submenu) {
					submenu.listElem.appendChild(itemElem);
				} else {
					this.listElem.appendChild(itemElem);
				}
			}
		});
	}

	private getSubmenu(path: (string | V)[] | undefined): DropdownSubmenu<V> | null {
		if (!path) {
			return null;
		}
		return this.submenus.find(submenu => this.equalPaths(submenu.path, path)) || null;
	}

	private createSubmenu(path: (string | V)[], buttonElem?: HTMLButtonElement, itemElem?: HTMLLIElement): DropdownSubmenu<V> {
		const submenu = this.getSubmenu(path);
		if (submenu) return submenu;

		let parent: DropdownSubmenu<V> | null = null;
		if (path.length > 1) parent = this.createSubmenu(path.slice(0, path.length - 1));

		if (!itemElem) itemElem = (<li className="dropdown-picker-item" />) as HTMLLIElement;

		if (!buttonElem) buttonElem = (<button className="dropdown-item" dataset={{ bsToggle: 'dropdown' }} attributes={{ 'aria-expanded': false }} />) as HTMLButtonElement;
		if (!buttonElem.childNodes.length) buttonElem.replaceChildren(path[path.length - 1] + ' \u00bb');

		const listRef = ref<HTMLUListElement>();

		itemElem.appendChild(
			<div className="dropend">
				{buttonElem}
				<ul ref={listRef} className="dropdown-submenu dropdown-menu"></ul>
			</div>,
		);

		(parent?.listElem || this.listElem).appendChild(itemElem);

		const newSubmenu = {
			path: path,
			listElem: listRef.value!,
		};
		this.submenus.push(newSubmenu);
		return newSubmenu;
	}

	private equalPaths(a: (string | V)[] | null | undefined, b: (string | V)[] | null | undefined): boolean {
		return (
			(a?.length || 0) == (b?.length || 0) &&
			(a || []).every((aVal, i) => (typeof aVal == 'string' ? aVal == (b![i] as string) : this.config.equals(aVal, b![i] as V)))
		);
	}

	getInputElem(): HTMLElement {
		return this.listElem;
	}

	getInputValue(): T {
		return this.valueToSource(this.currentSelection?.value as V);
	}

	setInputValue(newSrcValue: T) {
		const newValue = this.sourceToValue(newSrcValue);
		const newSelection = this.valueConfigs.find(v => this.config.equals(v.value, newValue))!;
		if (newSelection) {
			this.updateValue(newSelection);
		} else if (newValue == null) {
			this.updateValue(null);
		} else if (this.config.createMissingValue) {
			this.config.createMissingValue(newValue).then(newSelection => this.updateValue(newSelection));
		} else {
			this.updateValue(null);
		}
	}

	private updateValue(newValue: DropdownValueConfig<V> | null) {
		this.currentSelection = newValue;

		// Update button
		if (newValue) {
			this.buttonElem.innerHTML = '';
			this.config.setOptionContent(this.buttonElem, newValue, true);
		} else {
			this.buttonElem.textContent = this.config.defaultLabel;
		}
	}
}

export interface TextDropdownValueConfig<T> extends DropdownValueConfig<T> {
	label: string;
}

export interface TextDropdownPickerConfig<ModObject, T> extends Omit<DropdownPickerConfig<ModObject, T>, 'values' | 'setOptionContent'> {
	values: Array<TextDropdownValueConfig<T>>;
}

export class TextDropdownPicker<ModObject, T> extends DropdownPicker<ModObject, T> {
	constructor(parent: HTMLElement, modObject: ModObject, config: TextDropdownPickerConfig<ModObject, T>) {
		super(parent, modObject, {
			...config,
			setOptionContent: (button: HTMLButtonElement, valueConfig: DropdownValueConfig<T>) => {
				button.textContent = (valueConfig as TextDropdownValueConfig<T>).label;
			},
		});
	}
}
