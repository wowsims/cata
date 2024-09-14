import { TypedEvent } from '../../typed_event';
import { formatToNumber } from '../../utils';
import { Input, InputConfig } from '../input';

/**
 * Data for creating a number picker.
 */
export interface NumberPickerConfig<ModObject> extends InputConfig<ModObject, number> {
	id: string;
	// Whether the picker represents a float value. Default `false`
	float?: boolean;
	maxDecimalDigits?: number;
	// Whether to only allow positive values. Default `false`
	positive?: boolean;
	// Whether to show values of zero within the input. Default `true`
	showZeroes?: boolean;
}

// UI element for picking an arbitrary number field.
export class NumberPicker<ModObject> extends Input<ModObject, number> {
	private readonly inputElem: HTMLInputElement;
	private float: boolean;
	private positive: boolean;
	private showZeroes: boolean;
	private maxDecimalDigits: number;

	constructor(parent: HTMLElement | null, modObject: ModObject, config: NumberPickerConfig<ModObject>) {
		super(parent, 'number-picker-root', modObject, config);
		this.float = config.float ?? false;
		this.maxDecimalDigits = config.maxDecimalDigits ?? 2;
		this.positive = config.positive ?? false;
		this.showZeroes = config.showZeroes ?? true;

		this.inputElem = document.createElement('input');
		this.inputElem.id = config.id;
		this.inputElem.type = 'text';
		this.inputElem.classList.add('form-control', 'number-picker-input');

		if (this.positive) {
			this.inputElem.addEventListener(
				'change',
				() => {
					if (this.float) {
						this.inputElem.value = formatToNumber(Math.abs(Number(this.inputElem.value)), {
							minimumFractionDigits: 2,
							maximumFractionDigits: this.maxDecimalDigits,
						});
					} else {
						this.inputElem.value = Math.abs(parseInt(this.inputElem.value)).toString();
					}
				},
				{ signal: this.signal },
			);
		}

		this.inputElem.addEventListener(
			'change',
			() => {
				this.inputChanged(TypedEvent.nextEventID());
			},
			{ signal: this.signal },
		);

		this.inputElem.addEventListener(
			'input',
			() => {
				this.updateSize();
			},
			{ signal: this.signal },
		);

		this.rootElem.appendChild(this.inputElem);

		this.init();
		this.updateSize();
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): number {
		if (this.float) {
			return Number(this.inputElem.value || '') || 0;
		} else {
			return parseInt(this.inputElem.value || '') || 0;
		}
	}

	setInputValue(newValue: number) {
		if (newValue === 0 && !this.showZeroes) {
			this.inputElem.value = '';
			return;
		}

		if (this.float) {
			this.inputElem.value = formatToNumber(newValue, { minimumFractionDigits: 2, maximumFractionDigits: this.maxDecimalDigits });
		} else {
			this.inputElem.value = String(newValue);
		}
	}

	private updateSize() {
		const newSize = Math.max(3, this.inputElem.value.length);
		if (this.inputElem.size != newSize) this.inputElem.size = newSize;
	}
}
