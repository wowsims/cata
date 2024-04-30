import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { IconEnumPicker } from '../icon_enum_picker';
import { buildIconInput } from '../icon_inputs.js';
import { IconPicker } from '../icon_picker';
import * as ConsumablesInputs from '../inputs/consumables';
import { relevantStatOptions } from '../inputs/stat_options';
import { SettingsTab } from './settings_tab';

export class ConsumesPicker extends Component {
	protected settingsTab: SettingsTab;
	protected simUI: IndividualSimUI<any>;

	constructor(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<any>) {
		super(parentElem, 'consumes-picker-root');
		this.settingsTab = settingsTab;
		this.simUI = simUI;

		this.buildPotionsPicker();
		this.buildElixirsPicker();
		this.buildFoodPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Potions</label>
				<div class="consumes-row-inputs consumes-potions"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const potionsElem = this.rootElem.querySelector('.consumes-potions') as HTMLElement;

		//makePrepopPotionsInput;
		const prePotOptions = ConsumablesInputs.makePrepopPotionsInput(relevantStatOptions(ConsumablesInputs.PRE_POTIONS_CONFIG, this.simUI), 'Prepop Potion');
		const prePotPicker = buildIconInput(potionsElem, this.simUI.player, prePotOptions);

		const potionsOptions = ConsumablesInputs.makePotionsInput(relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI), 'Combat Potion');
		const potionsPicker = buildIconInput(potionsElem, this.simUI.player, potionsOptions);

		const conjuredOptions = ConsumablesInputs.makeConjuredInput(relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI));
		const conjuredPicker = buildIconInput(potionsElem, this.simUI.player, conjuredOptions);

		TypedEvent.onAny([this.simUI.player.professionChangeEmitter]).on(() => {
			this.updateRow(row, [potionsPicker, conjuredPicker, prePotPicker]);
		});
	}

	private buildElixirsPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
      		<div class="consumes-row input-root input-inline">
				<label class="form-label">Elixirs</label>
				<div class="consumes-row-inputs">
					<div class="consumes-flasks"></div>
					<span class="elixir-space">or</span>
					<div class="consumes-battle-elixirs"></div>
					<div class="consumes-guardian-elixirs"></div>
				</div>
			</div>
    	`;

		const rowElem = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const flasksElem = this.rootElem.querySelector('.consumes-flasks') as HTMLElement;
		const battleElixirsElem = this.rootElem.querySelector('.consumes-battle-elixirs') as HTMLElement;
		const guardianElixirsElem = this.rootElem.querySelector('.consumes-guardian-elixirs') as HTMLElement;
		const flasksOptions = ConsumablesInputs.makeFlasksInput(relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI));
		buildIconInput(flasksElem, this.simUI.player, flasksOptions);
		const battleElixirOptions = ConsumablesInputs.makeBattleElixirsInput(relevantStatOptions(ConsumablesInputs.BATTLE_ELIXIRS_CONFIG, this.simUI));
		buildIconInput(battleElixirsElem, this.simUI.player, battleElixirOptions);
		const guardianElixirOptions = ConsumablesInputs.makeGuardianElixirsInput(relevantStatOptions(ConsumablesInputs.GUARDIAN_ELIXIRS_CONFIG, this.simUI));
		buildIconInput(guardianElixirsElem, this.simUI.player, guardianElixirOptions);
	}

	private buildFoodPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Food</label>
				<div class="consumes-row-inputs consumes-food"></div>
			</div>
    	`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const foodsElem = this.rootElem.querySelector('.consumes-food') as HTMLElement;

		const foodOptions = ConsumablesInputs.makeFoodInput(relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI));
		buildIconInput(foodsElem, this.simUI.player, foodOptions);
	}

	private buildEngPicker() {
		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="consumes-row input-root input-inline">
				<label class="form-label">Engineering</label>
				<div class="consumes-row-inputs consumes-engi"></div>
			</div>
		`;

		const row = this.rootElem.appendChild(fragment.children[0] as HTMLElement);
		const engiConsumesElem = this.rootElem.querySelector('.consumes-engi') as HTMLElement;

		const tinkerOptions = ConsumablesInputs.makeTinkerHandsInput(relevantStatOptions(ConsumablesInputs.TINKERS_HANDS_CONFIG, this.simUI), 'Gloves Tinkers');
		const tinkerPicker = buildIconInput(engiConsumesElem, this.simUI.player, tinkerOptions);

		const decoyPicker = buildIconInput(engiConsumesElem, this.simUI.player, ConsumablesInputs.ExplosiveDecoy);
		const sapperPicker = buildIconInput(engiConsumesElem, this.simUI.player, ConsumablesInputs.ThermalSapper);

		const explosiveOptions = ConsumablesInputs.makeExplosivesInput(relevantStatOptions(ConsumablesInputs.EXPLOSIVES_CONFIG, this.simUI), 'Explosives');
		const explosivePicker = buildIconInput(engiConsumesElem, this.simUI.player, explosiveOptions);

		TypedEvent.onAny([this.simUI.player.professionChangeEmitter]).on(() => this.updateRow(row, [decoyPicker, sapperPicker, explosivePicker, tinkerPicker]));
		this.updateRow(row, [decoyPicker, sapperPicker, explosivePicker, tinkerPicker]);
	}

	private buildPetPicker() {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
			const fragment = document.createElement('fragment');
			fragment.innerHTML = `
				<div class="consumes-row input-root input-inline">
					<label class="form-label">Pet</label>
					<div class="consumes-row-inputs consumes-pet"></div>
				</div>
			`;

			this.rootElem.appendChild(fragment.children[0] as HTMLElement);
			const petConsumesElem = this.rootElem.querySelector('.consumes-pet') as HTMLElement;

			this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
		}
	}

	private updateRow(rowElem: HTMLElement, pickers: (IconPicker<Player<any>, any> | IconEnumPicker<Player<any>, any>)[]) {
		if (!!pickers.find(p => p?.showWhen())) {
			rowElem.classList.remove('hide');
		} else {
			rowElem.classList.add('hide');
		}
	}
}
