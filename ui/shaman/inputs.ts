import { ContentBlock } from '../core/components/content_block';
import { IconEnumPicker } from '../core/components/icon_enum_picker';
import { IconPicker } from '../core/components/icon_picker';
import { Input } from '../core/components/input';
import * as InputHelpers from '../core/components/input_helpers';
import { IndividualSimUI } from '../core/individual_sim_ui';
import { Player } from '../core/player';
import { Spec } from '../core/proto/common';
import { AirTotem, EarthTotem, FireTotem, ShamanImbue, ShamanShield, ShamanTotems, WaterTotem } from '../core/proto/shaman';
import { ActionId } from '../core/proto_utils/action_id';
import { ShamanSpecs } from '../core/proto_utils/utils';
import { EventID } from '../core/typed_event';

// Configuration for class-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanShieldInput = <SpecType extends ShamanSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, ShamanShield>({
		fieldName: 'shield',
		values: [
			{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
			{ actionId: ActionId.fromSpellId(52127), value: ShamanShield.WaterShield },
			{ actionId: ActionId.fromSpellId(324), value: ShamanShield.LightningShield },
		],
	});

export const ShamanImbueMH = <SpecType extends ShamanSpecs>() =>
	InputHelpers.makeClassOptionsEnumIconInput<SpecType, ShamanImbue>({
		fieldName: 'imbueMh',
		values: [
			{ value: ShamanImbue.NoImbue, tooltip: 'No Main Hand Enchant' },
			{ actionId: ActionId.fromSpellId(8232), value: ShamanImbue.WindfuryWeapon },
			{ actionId: ActionId.fromSpellId(8024), value: ShamanImbue.FlametongueWeapon },
			{ actionId: ActionId.fromSpellId(8033), value: ShamanImbue.FrostbrandWeapon },
		],
	});

export function TotemsSection(parentElem: HTMLElement, simUI: IndividualSimUI<any>): ContentBlock {
	const contentBlock = new ContentBlock(parentElem, 'totems-settings', {
		header: { title: 'Totems' },
	});

	const totemDropdownGroup = Input.newGroupContainer();
	totemDropdownGroup.classList.add('totem-dropdowns-container', 'icon-group');

	const fireElementalContainer = document.createElement('div');
	fireElementalContainer.classList.add('fire-elemental-input-container');

	contentBlock.bodyElement.appendChild(totemDropdownGroup);
	contentBlock.bodyElement.appendChild(fireElementalContainer);

	const _earthTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['earth-totem-picker'],
		numColumns: 1,
		values: [
			{ color: '#ffdfba', value: EarthTotem.NoEarthTotem },
			{ actionId: ActionId.fromSpellId(8075), value: EarthTotem.StrengthOfEarthTotem },
			{ actionId: ActionId.fromSpellId(8071), value: EarthTotem.StoneskinTotem },
			{ actionId: ActionId.fromSpellId(8143), value: EarthTotem.TremorTotem },
			{ actionId: ActionId.fromSpellId(2062), value: EarthTotem.EarthElementalTotem},
		],
		equals: (a: EarthTotem, b: EarthTotem) => a == b,
		zeroValue: EarthTotem.NoEarthTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.totems?.earth || EarthTotem.NoEarthTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			newOptions.classOptions!.totems.earth = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const _waterTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['water-totem-picker'],
		numColumns: 1,
		values: [
			{ color: '#bae1ff', value: WaterTotem.NoWaterTotem },
			{ actionId: ActionId.fromSpellId(5675), value: WaterTotem.ManaSpringTotem },
			{ actionId: ActionId.fromSpellId(5394), value: WaterTotem.HealingStreamTotem },
		],
		equals: (a: WaterTotem, b: WaterTotem) => a == b,
		zeroValue: WaterTotem.NoWaterTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.totems?.water || WaterTotem.NoWaterTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			newOptions.classOptions!.totems.water = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const _fireTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['fire-totem-picker'],
		numColumns: 1,
		values: [
			{ color: '#ffb3ba', value: FireTotem.NoFireTotem },
			{ actionId: ActionId.fromSpellId(78770), value: FireTotem.MagmaTotem },
			{ actionId: ActionId.fromSpellId(3599), value: FireTotem.SearingTotem },
			{ actionId: ActionId.fromSpellId(8227), value: FireTotem.FlametongueTotem },
			{ actionId: ActionId.fromSpellId(2894), value: FireTotem.FireElementalTotem},
		],
		equals: (a: FireTotem, b: FireTotem) => a == b,
		zeroValue: FireTotem.NoFireTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.totems?.fire || FireTotem.NoFireTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			newOptions.classOptions!.totems.fire = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const _airTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['air-totem-picker'],
		numColumns: 1,
		values: [
			{ color: '#baffc9', value: AirTotem.NoAirTotem },
			{ actionId: ActionId.fromSpellId(8512), value: AirTotem.WindfuryTotem },
			{ actionId: ActionId.fromSpellId(3738), value: AirTotem.WrathOfAirTotem },
		],
		equals: (a: AirTotem, b: AirTotem) => a == b,
		zeroValue: AirTotem.NoAirTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.totems?.air || AirTotem.NoAirTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			newOptions.classOptions!.totems.air = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	// Enchancement Shaman uses the Fire Elemental Inputs with custom inputs.
	if (simUI.player.getSpec() != Spec.SpecEnhancementShaman) {
		const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>(
			{
				getModObject: (player: Player<ShamanSpecs>) => player,
				getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.totems || ShamanTotems.create(),
				setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
					const newOptions = player.getSpecOptions();
					newOptions.classOptions!.totems = newVal;
					player.setSpecOptions(eventID, newOptions);
				},
				changeEmitter: (player: Player<any>) => player.specOptionsChangeEmitter,
			},
			ActionId.fromSpellId(2894),
			'useFireElemental',
		);

		new IconPicker(fireElementalContainer, simUI.player, fireElementalBooleanIconInput);
	}

	return contentBlock;
}
