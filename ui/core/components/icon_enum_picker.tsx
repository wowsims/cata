import tippy from 'tippy.js';


import { ActionId } from '../proto_utils/action_id.js';
import { TypedEvent } from '../typed_event.js';
import { Input, InputConfig } from './input.js';

export enum IconEnumPickerDirection {
	Vertical = 'vertical',
	Horizontal = 'Horizontal',
}

export interface IconEnumValueConfig<ModObject, T> {
	value: T;
	// One of these should be set.
	// If actionId is set, shows the icon for that id.
	// If color is set, shows that color.
	// If iconUrl is set, shows that icon as grayscale
	actionId?: ActionId;
	color?: string;
	iconUrl?: string;
	// Text to be displayed on the icon.
	text?: string;
	// Hover tooltip.
	tooltip?: string;

	showWhen?: (obj: ModObject) => boolean;
}

export interface IconEnumPickerConfig<ModObject, T> extends InputConfig<ModObject, T> {
	numColumns?: number;
	values: Array<IconEnumValueConfig<ModObject, T>>;
	// Value that will be considered inactive.
	zeroValue: T;
	// Function for comparing two values.
	// Tooltip that will be shown whne hovering over the icon-picker-button
	tooltip?: string;
	// The direction the menu will open in relative to the root element
	direction?: IconEnumPickerDirection;
	equals: (a: T, b: T) => boolean;
	backupIconUrl?: (value: T) => ActionId;
	showWhen?: (obj: ModObject) => boolean;
}

// Icon-based UI for picking enum values.
// ModObject is the object being modified (Sim, Player, or Target).
export class IconEnumPicker<ModObject, T> extends Input<ModObject, T> {
	private readonly config: IconEnumPickerConfig<ModObject, T>;

	private currentValue: T;
	private storedValue: T | undefined;

	private readonly buttonElem: HTMLAnchorElement;
	private readonly buttonText: HTMLElement;

	constructor(parent: HTMLElement, modObj: ModObject, config: IconEnumPickerConfig<ModObject, T>) {
		super(parent, 'icon-enum-picker-root', modObj, config);
		this.rootElem.classList.add('icon-picker', (config.direction ?? 'vertical') === 'vertical' ? 'dropdown' : 'dropend');
		this.config = config;
		this.currentValue = this.config.zeroValue;

		this.config.changedEvent(this.modObject).on(_ => {
			if (this.showWhen()) {
				this.restoreValue();
				this.rootElem.classList.remove('hide');
			} else {
				this.storeValue();
				this.rootElem.classList.add('hide');
			}
		});

		if (config.tooltip) {
			tippy(this.rootElem, {
				content: config.tooltip,
			});
		}

		this.rootElem.appendChild(
			<>
				<a
					href="javascript:void(0)"
					className="icon-picker-button"
					attributes={{
						'aria-expanded': 'false',
						role: 'button',
					}}
					dataset={{
						bsToggle: 'dropdown',
						whtticon: 'false',
						disableWowheadTouchTooltip: 'true',
					}}></a>
				<ul className="dropdown-menu"></ul>
				<label className="form-label"></label>
			</>,
		);

		this.buttonElem = this.rootElem.querySelector('.icon-picker-button') as HTMLAnchorElement;
		this.buttonText = this.rootElem.querySelector('label') as HTMLElement;
		const dropdownMenu = this.rootElem.querySelector('.dropdown-menu') as HTMLElement;

		if (this.config.numColumns) dropdownMenu.style.gridTemplateColumns = `repeat(${this.config.numColumns}, 1fr)`;

		if (this.config.direction == IconEnumPickerDirection.Horizontal) dropdownMenu.style.gridAutoFlow = 'column';

		this.config.values.forEach((valueConfig, _) => {
			const optionContainer = document.createElement('li');
			optionContainer.classList.add('icon-dropdown-option', 'dropdown-option');
			dropdownMenu.appendChild(optionContainer);

			const option = document.createElement('a');
			option.classList.add('icon-picker-button');
			option.dataset.whtticon = 'false';
			option.dataset.disableWowheadTouchTooltip = 'true';
			optionContainer.appendChild(option);

			const updateOption = () => {
				this.setImage(option, valueConfig);
				if (this.showValueWhen(valueConfig)) {
					optionContainer.classList.remove('hide');
				} else {
					optionContainer.classList.add('hide');
					// Zero out the picker if the selected option is hidden
					if (this.currentValue == valueConfig.value) {
						this.setInputValue(this.config.zeroValue);
						this.inputChanged(TypedEvent.nextEventID());
					}
				}
			};

			this.config.changedEvent(this.modObject).on(updateOption);
			updateOption();

			option.addEventListener('click', event => {
				event.preventDefault();
				this.currentValue = valueConfig.value;
				this.storedValue = undefined;
				this.inputChanged(TypedEvent.nextEventID());
			});

			if (valueConfig.tooltip) {
				tippy(option, {
					content: valueConfig.tooltip,
				});
			}
		});

		this.init();
	}

	/**
	 * Stores value of current input and hides the element for later
	 * restoration. Useful for events which trigger the element
	 * on and off.
	 */
	storeValue() {
		if (typeof this.storedValue !== 'undefined') return;

		this.storedValue = this.getInputValue();
		this.setInputValue(this.config.zeroValue);
		this.inputChanged(TypedEvent.nextEventID());
	}

	/**
	 * Restores value of current input and shows the element.
	 */
	restoreValue() {
		if (typeof this.storedValue === 'undefined') return;

		this.setInputValue(this.storedValue);
		this.inputChanged(TypedEvent.nextEventID());
		this.storedValue = undefined;
	}

	private setActionImage(elem: HTMLAnchorElement, actionId: ActionId) {
		actionId.fillAndSet(elem, true, true);
	}

	private setImage(elem: HTMLAnchorElement, valueConfig: IconEnumValueConfig<ModObject, T>) {
		if (valueConfig.showWhen && !valueConfig.showWhen(this.modObject)) {
			elem.removeAttribute('href');
			return;
		}

		const actionId = valueConfig.actionId;
		if (actionId) {
			this.setActionImage(elem, actionId);
			elem.style.filter = '';
		} else if (valueConfig.iconUrl) {
			elem.style.backgroundImage = `url(${valueConfig.iconUrl})`;
			elem.style.filter = 'grayscale(1)';
		} else {
			elem.style.backgroundImage = '';
			elem.style.filter = '';
			elem.style.backgroundColor = valueConfig.color!;
		}
	}

	update() {
		super.update();
		this.setActive(this.enabled && !this.config.equals(this.currentValue, this.config.zeroValue));
	}

	getInputElem(): HTMLElement {
		return this.buttonElem;
	}

	getInputValue(): T {
		return this.currentValue;
	}

	setInputValue(newValue: T) {
		this.currentValue = newValue;
		this.setActive(this.enabled && !this.config.equals(this.currentValue, this.config.zeroValue));

		this.buttonText.textContent = '';
		this.buttonText.style.display = 'none';

		const valueConfig = this.config.values.find(valueConfig => this.config.equals(valueConfig.value, this.currentValue))!;
		if (valueConfig) {
			this.setImage(this.buttonElem, valueConfig);
			if (valueConfig.text != undefined) {
				this.buttonText.style.display = 'block';
				this.buttonText.textContent = valueConfig.text;
			}
		} else if (this.config.backupIconUrl) {
			const backupId = this.config.backupIconUrl(this.currentValue);
			this.setActionImage(this.buttonElem, backupId);
			this.setActive(false);
		}
	}

	setActive(active: boolean) {
		if (active) this.buttonElem.classList.add('active');
		else this.buttonElem.classList.remove('active');
	}

	showWhen(): boolean {
		return (
			(!this.config.showWhen || this.config.showWhen(this.modObject)) &&
			!!this.config.values.find(
				valueConfig => valueConfig.actionId && valueConfig.actionId != null && (!valueConfig.showWhen || valueConfig.showWhen(this.modObject)),
			)
		);
	}

	showValueWhen(valueConfig: IconEnumValueConfig<ModObject, T>): boolean {
		return (!valueConfig.actionId || valueConfig.actionId != null) && (!valueConfig.showWhen || valueConfig.showWhen(this.modObject));
	}
}
