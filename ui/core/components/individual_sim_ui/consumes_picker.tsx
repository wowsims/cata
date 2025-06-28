import { ref } from 'tsx-vanilla';

import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { Class, ConsumableType, Spec } from '../../proto/common';
import { Consumable } from '../../proto/db';
import { Database } from '../../proto_utils/database';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';
import { buildIconInput } from '../icon_inputs';
import * as ConsumablesInputs from '../inputs/consumables';
import { relevantStatOptions } from '../inputs/stat_options';
import { IconEnumPicker } from '../pickers/icon_enum_picker';
import { IconPicker } from '../pickers/icon_picker';
import { SettingsTab } from './settings_tab';

export class ConsumesPicker extends Component {
	protected settingsTab: SettingsTab;
	protected simUI: IndividualSimUI<any>;
	protected db: Database;

	constructor(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<any>, db: Database) {
		super(parentElem, 'consumes-picker-root');
		this.settingsTab = settingsTab;
		this.simUI = simUI;
		this.db = db;
	}

	private getConsumables(type: ConsumableType): Consumable[] {
		return this.db.getConsumablesByTypeAndStats(type, this.simUI.individualConfig.consumableStats ?? this.simUI.individualConfig.epStats);
	}

	public static create(parentElem: HTMLElement, settingsTab: SettingsTab, simUI: IndividualSimUI<any>): ConsumesPicker {
		const instance = new ConsumesPicker(parentElem, settingsTab, simUI, Database.getSync());
		instance.init();
		return instance;
	}

	private init(): void {
		this.buildPotionsPicker();
		this.buildElixirsPicker();
		this.buildFoodPicker();
		this.buildEngPicker();
		this.buildPetPicker();
	}

	private buildPotionsPicker(): void {
		const potionsRef = ref<HTMLDivElement>();

		const row = this.rootElem.appendChild(
			<ConsumeRow label="Potions">
				<div ref={potionsRef} className="picker-group icon-group consumes-row-inputs consumes-potions"></div>
			</ConsumeRow>,
		);
		const potionsElem = potionsRef.value!;

		let pots = this.getConsumables(ConsumableType.ConsumableTypePotion);
		if (this.simUI.player.getClass() !== Class.ClassWarrior && this.simUI.player.getSpec() !== Spec.SpecGuardianDruid) {
			pots = pots.filter(pot => pot.id !== 13442);
		}
		const prePotOptions = ConsumablesInputs.makeConsumableInput(pots, { consumesFieldName: 'prepotId' }, 'Prepop Potion');
		const potionsOptions = ConsumablesInputs.makeConsumableInput(pots, { consumesFieldName: 'potId' }, 'Combat Potion');

		const prePotPicker = buildIconInput(potionsElem, this.simUI.player, prePotOptions);

		const potionsPicker = buildIconInput(potionsElem, this.simUI.player, potionsOptions);

		const conjuredOptions = ConsumablesInputs.makeConjuredInput(relevantStatOptions(ConsumablesInputs.CONJURED_CONFIG, this.simUI));
		const conjuredPicker = buildIconInput(potionsElem, this.simUI.player, conjuredOptions);

		const events = TypedEvent.onAny([this.simUI.player.professionChangeEmitter]).on(() =>
			this.updateRow(row, [potionsPicker, conjuredPicker, prePotPicker]),
		);
		this.addOnDisposeCallback(() => events.dispose());

		this.updateRow(row, [potionsPicker, conjuredPicker, prePotPicker]);
	}

	private buildElixirsPicker(): void {
		const flaskRef = ref<HTMLDivElement>();
		const battleElixirsRef = ref<HTMLDivElement>();
		const guardianElixirsRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<ConsumeRow label="Elixirs">
				<div className="picker-group icon-group consumes-row-inputs">
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

		const flasks = this.getConsumables(ConsumableType.ConsumableTypeFlask);
		const simpleFlasksOptions = ConsumablesInputs.makeConsumableInput(flasks, { consumesFieldName: 'flaskId' }, '');
		buildIconInput(flasksElem, this.simUI.player, simpleFlasksOptions);

		const battleElixirs = this.getConsumables(ConsumableType.ConsumableTypeBattleElixir);
		const battleElixirOptions = ConsumablesInputs.makeConsumableInput(battleElixirs, { consumesFieldName: 'battleElixirId' }, '');

		const guardianElixirs = this.getConsumables(ConsumableType.ConsumableTypeGuardianElixir);
		const guardianElixirOptions = ConsumablesInputs.makeConsumableInput(guardianElixirs, { consumesFieldName: 'guardianElixirId' }, '');

		buildIconInput(battleElixirsElem, this.simUI.player, battleElixirOptions);

		buildIconInput(guardianElixirsElem, this.simUI.player, guardianElixirOptions);
	}

	private buildFoodPicker(): void {
		const foodRef = ref<HTMLDivElement>();
		this.rootElem.appendChild(
			<ConsumeRow label="Food">
				<div ref={foodRef} className="picker-group icon-group consumes-row-inputs consumes-food"></div>
			</ConsumeRow>,
		);
		const foodsElem = foodRef.value!;
		const foods = this.getConsumables(ConsumableType.ConsumableTypeFood);
		const foodsOptions = ConsumablesInputs.makeConsumableInput(foods, { consumesFieldName: 'foodId' }, '');
		buildIconInput(foodsElem, this.simUI.player, foodsOptions);
	}

	private buildEngPicker(): void {
		const engiConsumesRef = ref<HTMLDivElement>();
		const row = this.rootElem.appendChild(
			<ConsumeRow label="Engineering">
				<div ref={engiConsumesRef} className="picker-group icon-group consumes-row-inputs consumes-engi"></div>
			</ConsumeRow>,
		);
		const engiConsumesElem = engiConsumesRef.value!;

		const explosivesoptions = ConsumablesInputs.makeExplosivesInput(relevantStatOptions(ConsumablesInputs.EXPLOSIVE_CONFIG, this.simUI), 'Explosives');
		const explosivePicker = buildIconInput(engiConsumesElem, this.simUI.player, explosivesoptions);

		const events = this.simUI.player.professionChangeEmitter.on(() => this.updateRow(row, [explosivePicker]));
		this.addOnDisposeCallback(() => events.dispose());

		// Initial update of row based on current state.
		this.updateRow(row, [explosivePicker]);
	}

	private buildPetPicker(): void {
		if (this.simUI.individualConfig.petConsumeInputs?.length) {
			const petConsumesRef = ref<HTMLDivElement>();
			this.rootElem.appendChild(
				<ConsumeRow label="Pet">
					<div ref={petConsumesRef} className="picker-group icon-group consumes-row-inputs consumes-pet"></div>
				</ConsumeRow>,
			);
			const petConsumesElem = petConsumesRef.value!;

			// Create pickers for each pet consume input.
			this.simUI.individualConfig.petConsumeInputs.forEach(iconInput => buildIconInput(petConsumesElem, this.simUI.player, iconInput));
		}
	}

	private updateRow(rowElem: Element, pickers: (IconPicker<Player<any>, any> | IconEnumPicker<Player<any>, any>)[]) {
		rowElem.classList[!!pickers.find(p => p?.showWhen()) ? 'remove' : 'add']('hide');
	}
}

// A simple JSX stateless component for rows.
const ConsumeRow = ({ label, children }: { label: string; children: JSX.Element }) => (
	<div className="consumes-row input-root input-inline">
		<label className="form-label">{label}</label>
		{children}
	</div>
);
