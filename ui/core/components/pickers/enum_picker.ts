import { TypedEvent } from '../../typed_event.js';
import { Input, InputConfig } from '../input.js';

export interface EnumValueConfig {
	name: string;
	value: number;
	tooltip?: string;
}

export interface EnumPickerConfig<ModObject> extends InputConfig<ModObject, number> {
	id: string;
	values: Array<EnumValueConfig>;
}

export class EnumPicker<ModObject> extends Input<ModObject, number> {
	private readonly selectElem: HTMLSelectElement;

	constructor(parent: HTMLElement | null, modObject: ModObject, config: EnumPickerConfig<ModObject>) {
		super(parent, 'enum-picker-root', modObject, config);

		this.selectElem = document.createElement('select');
		this.selectElem.id = config.id;
		this.selectElem.classList.add('enum-picker-selector', 'form-select');

		config.values.forEach(value => {
			const option = document.createElement('option');
			option.value = String(value.value);
			option.textContent = value.name;
			this.selectElem.appendChild(option);

			if (value.tooltip) {
				option.title = value.tooltip;
			}
		});
		this.rootElem.appendChild(this.selectElem);

		this.init();

		this.selectElem.addEventListener(
			'change',
			() => {
				this.inputChanged(TypedEvent.nextEventID());
			},
			{ signal: this.signal },
		);
	}

	getInputElem(): HTMLElement {
		return this.selectElem;
	}

	getInputValue(): number {
		return Number(this.selectElem.value);
	}

	setInputValue(newValue: number) {
		this.selectElem.value = String(newValue);
	}
}
