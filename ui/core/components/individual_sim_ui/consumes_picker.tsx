import clsx from 'clsx';
import { element, fragment, ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { IconEnumPicker } from '../icon_enum_picker';
import { buildIconInput } from '../icon_inputs';
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
		const potionsRef = ref<HTMLDivElement>();

		const row = this.rootElem.appendChild(
			<ConsumeRow label="Potions">
				<div className="consumes-row-inputs consumes-potions"></div>
			</ConsumeRow>,
		);
		const potionsElem = potionsRef.value!;

		//makePrepopPotionsInput;
		const prePotOptions = ConsumablesInputs.makePrepopPotionsInput(relevantStatOptions(ConsumablesInputs.PRE_POTIONS_CONFIG, this.simUI), 'Prepop Potion');
		const prePotPicker = buildIconInput(potionsElem, this.simUI.player, prePotOptions);

		const potionsOptions = ConsumablesInputs.makePotionsInput(relevantStatOptions(ConsumablesInputs.POTIONS_CONFIG, this.simUI), 'Combat Potion');
		const potionsPicker = buildIconInput(potionsElem, this.simUI.player, potionsOptions);

		const conjuredOptions = ConsumablesInputs.makeConjuredInput(relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI));
		const conjuredPicker = buildIconInput(potionsElem, this.simUI.player, conjuredOptions);

		const events = TypedEvent.onAny([this.simUI.player.professionChangeEmitter]).on(() =>
			this.updateRow(row, [potionsPicker, conjuredPicker, prePotPicker]),
		);
		this.addOnDisposeCallback(() => events.dispose());
	}

	private buildElixirsPicker() {
		const flaskRef = ref<HTMLDivElement>();
		const battleElixirsRef = ref<HTMLDivElement>();
		const guardianElixirsRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<ConsumeRow label="Elixirs">
				<div className="consumes-row-inputs">
					<div ref={flaskRef} className="consumes-flasks"></div>
					<span className="elixir-space">or</span>
					<div ref={battleElixirsRef} className="consumes-battle-elixirs"></div>
					<div ref={guardianElixirsRef} className="consumes-guardian-elixirs"></div>
				</div>
			</ConsumeRow>,
		);
		const flasksElem = flaskRef.value!;
		const battleElixirsElem = battleElixirsRef.value!;
		const guardianElixirsElem = guardianElixirsRef.value!;

		const flasksOptions = ConsumablesInputs.makeFlasksInput(relevantStatOptions(ConsumablesInputs.FLASKS_CONFIG, this.simUI));
		buildIconInput(flasksElem, this.simUI.player, flasksOptions);
		const battleElixirOptions = ConsumablesInputs.makeBattleElixirsInput(relevantStatOptions(ConsumablesInputs.BATTLE_ELIXIRS_CONFIG, this.simUI));
		buildIconInput(battleElixirsElem, this.simUI.player, battleElixirOptions);
		const guardianElixirOptions = ConsumablesInputs.makeGuardianElixirsInput(relevantStatOptions(ConsumablesInputs.GUARDIAN_ELIXIRS_CONFIG, this.simUI));
		buildIconInput(guardianElixirsElem, this.simUI.player, guardianElixirOptions);
	}

	private buildFoodPicker() {
		const foodRef = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<ConsumeRow label="Food">
				<div ref={foodRef} className="consumes-row-inputs consumes-food"></div>
			</ConsumeRow>,
		);
		const foodsElem = foodRef.value!;

		const foodOptions = ConsumablesInputs.makeFoodInput(relevantStatOptions(ConsumablesInputs.FOOD_CONFIG, this.simUI));
		buildIconInput(foodsElem, this.simUI.player, foodOptions);
	}

	private buildEngPicker() {
		const engiConsumesRef = ref<HTMLDivElement>();
		const row = this.rootElem.appendChild(
			<ConsumeRow label="Engineering">
				<div className="consumes-row-inputs consumes-engi"></div>
			</ConsumeRow>,
		);
		const engiConsumesElem = engiConsumesRef.value!;

		const tinkerOptions = ConsumablesInputs.makeTinkerHandsInput(relevantStatOptions(ConsumablesInputs.TINKERS_HANDS_CONFIG, this.simUI), 'Gloves Tinkers');
		const tinkerPicker = buildIconInput(engiConsumesElem, this.simUI.player, tinkerOptions);

		const explosivePicker = buildIconInput(engiConsumesElem, this.simUI.player, ConsumablesInputs.ExplosiveBigDaddy);
		const highpoweredBoltGunPicker = buildIconInput(engiConsumesElem, this.simUI.player, ConsumablesInputs.HighpoweredBoltGun);

		const events = TypedEvent.onAny([this.simUI.player.professionChangeEmitter]).on(() =>
			this.updateRow(row, [explosivePicker, highpoweredBoltGunPicker, tinkerPicker]),
		);
		this.addOnDisposeCallback(() => events.dispose());
		this.updateRow(row, [explosivePicker, highpoweredBoltGunPicker, tinkerPicker]);
	}

	private buildPetPicker() {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
			const petConsumesRef = ref<HTMLDivElement>();
			this.rootElem.appendChild(
				<ConsumeRow label="Pet">
					<div ref={petConsumesRef} className="consumes-row-inputs consumes-pet"></div>
				</ConsumeRow>,
			);
			const petConsumesElem = petConsumesRef.value!;

			this.simUI.individualConfig.petConsumeInputs.map(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
		}
	}

	private updateRow(rowElem: Element, pickers: (IconPicker<Player<any>, any> | IconEnumPicker<Player<any>, any>)[]) {
		rowElem.classList[!!pickers.find(p => p?.showWhen()) ? 'remove' : 'add']('hide');
	}
}

const ConsumeRow = ({ label, children }: { label: string; children: JSX.Element }) => (
	<div className="consumes-row input-root input-inline">
		<label className="form-label">{label}</label>
		{children}
	</div>
);
