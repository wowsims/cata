import { ContentBlock } from '../core/components/content_block';
import { Input } from '../core/components/input';
import * as InputHelpers from '../core/components/input_helpers';
import { IconEnumPicker } from '../core/components/pickers/icon_enum_picker';
import { IndividualSimUI } from '../core/individual_sim_ui';
import { Player } from '../core/player';
import { AirTotem, CallTotem, EarthTotem, FireTotem, ShamanImbue, ShamanShield, ShamanTotems, TotemSet, WaterTotem } from '../core/proto/shaman';
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

	contentBlock.bodyElement.appendChild(totemDropdownGroup);

	const _callTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['call-totem-picker'],
		numColumns: 1,
		values: [
			{ actionId: ActionId.fromSpellId(66842), value: CallTotem.Elements },
			{ actionId: ActionId.fromSpellId(66843), value: CallTotem.Ancestors },
			{ actionId: ActionId.fromSpellId(66844), value: CallTotem.Spirits },
		],
		equals: (a: CallTotem, b: CallTotem) => a == b,
		zeroValue: CallTotem.NoCall,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().classOptions?.call || CallTotem.Elements,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			switch (newValue) {
				case CallTotem.Elements: {
					_earthTotemPicker.setInputValue(newOptions.classOptions?.totems?.elements?.earth || EarthTotem.NoEarthTotem);
					_waterTotemPicker.setInputValue(newOptions.classOptions?.totems?.elements?.water || WaterTotem.NoWaterTotem);
					_fireTotemPicker.setInputValue(newOptions.classOptions?.totems?.elements?.fire || FireTotem.NoFireTotem);
					_airTotemPicker.setInputValue(newOptions.classOptions?.totems?.elements?.air || AirTotem.NoAirTotem);
					break;
				}
				case CallTotem.Ancestors: {
					_earthTotemPicker.setInputValue(newOptions.classOptions?.totems?.ancestors?.earth || EarthTotem.NoEarthTotem);
					_waterTotemPicker.setInputValue(newOptions.classOptions?.totems?.ancestors?.water || WaterTotem.NoWaterTotem);
					_fireTotemPicker.setInputValue(newOptions.classOptions?.totems?.ancestors?.fire || FireTotem.NoFireTotem);
					_airTotemPicker.setInputValue(newOptions.classOptions?.totems?.ancestors?.air || AirTotem.NoAirTotem);
					break;
				}
				case CallTotem.Spirits: {
					_earthTotemPicker.setInputValue(newOptions.classOptions?.totems?.spirits?.earth || EarthTotem.NoEarthTotem);
					_waterTotemPicker.setInputValue(newOptions.classOptions?.totems?.spirits?.water || WaterTotem.NoWaterTotem);
					_fireTotemPicker.setInputValue(newOptions.classOptions?.totems?.spirits?.fire || FireTotem.NoFireTotem);
					_airTotemPicker.setInputValue(newOptions.classOptions?.totems?.spirits?.air || AirTotem.NoAirTotem);
					break;
				}
			}
			newOptions.classOptions!.call = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const _earthTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: ['earth-totem-picker'],
		numColumns: 1,
		values: [
			{ color: '#ffdfba', value: EarthTotem.NoEarthTotem },
			{ actionId: ActionId.fromSpellId(8075), value: EarthTotem.StrengthOfEarthTotem },
			{ actionId: ActionId.fromSpellId(8071), value: EarthTotem.StoneskinTotem },
			{ actionId: ActionId.fromSpellId(8143), value: EarthTotem.TremorTotem },
			{ actionId: ActionId.fromSpellId(2062), value: EarthTotem.EarthElementalTotem },
		],
		equals: (a: EarthTotem, b: EarthTotem) => a == b,
		zeroValue: EarthTotem.NoEarthTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => {
			const newOptions = player.getSpecOptions();
			let value = EarthTotem.NoEarthTotem;
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					value = player.getSpecOptions().classOptions?.totems?.elements?.earth || EarthTotem.NoEarthTotem;
					break;
				}
				case CallTotem.Ancestors: {
					value = player.getSpecOptions().classOptions?.totems?.ancestors?.earth || EarthTotem.NoEarthTotem;
					break;
				}
				case CallTotem.Spirits: {
					value = player.getSpecOptions().classOptions?.totems?.spirits?.earth || EarthTotem.NoEarthTotem;
					break;
				}
			}
			return value;
		},
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					newOptions.classOptions!.totems.earth = newValue;
					if (!newOptions.classOptions?.totems.elements) newOptions.classOptions.totems.elements = TotemSet.create();
					newOptions.classOptions.totems.elements!.earth = newValue;
					break;
				}
				case CallTotem.Ancestors: {
					if (!newOptions.classOptions?.totems.ancestors) newOptions.classOptions.totems.ancestors = TotemSet.create();
					newOptions.classOptions.totems.ancestors!.earth = newValue;
					break;
				}
				case CallTotem.Spirits: {
					if (!newOptions.classOptions?.totems.spirits) newOptions.classOptions.totems.spirits = TotemSet.create();
					newOptions.classOptions.totems.spirits!.earth = newValue;
					break;
				}
			}
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
		getValue: (player: Player<ShamanSpecs>) => {
			const newOptions = player.getSpecOptions();
			let value = WaterTotem.NoWaterTotem;
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					value = player.getSpecOptions().classOptions?.totems?.elements?.water || WaterTotem.NoWaterTotem;
					break;
				}
				case CallTotem.Ancestors: {
					value = player.getSpecOptions().classOptions?.totems?.ancestors?.water || WaterTotem.NoWaterTotem;
					break;
				}
				case CallTotem.Spirits: {
					value = player.getSpecOptions().classOptions?.totems?.spirits?.water || WaterTotem.NoWaterTotem;
					break;
				}
			}
			return value;
		},
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					newOptions.classOptions!.totems.water = newValue;
					if (!newOptions.classOptions?.totems.elements) newOptions.classOptions.totems.elements = TotemSet.create();
					newOptions.classOptions.totems.elements!.water = newValue;
					break;
				}
				case CallTotem.Ancestors: {
					if (!newOptions.classOptions?.totems.ancestors) newOptions.classOptions.totems.ancestors = TotemSet.create();
					newOptions.classOptions.totems.ancestors!.water = newValue;
					break;
				}
				case CallTotem.Spirits: {
					if (!newOptions.classOptions?.totems.spirits) newOptions.classOptions.totems.spirits = TotemSet.create();
					newOptions.classOptions.totems.spirits!.water = newValue;
					break;
				}
			}
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
			{ actionId: ActionId.fromSpellId(2894), value: FireTotem.FireElementalTotem },
		],
		equals: (a: FireTotem, b: FireTotem) => a == b,
		zeroValue: FireTotem.NoFireTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => {
			const newOptions = player.getSpecOptions();
			let value = FireTotem.NoFireTotem;
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					value = player.getSpecOptions().classOptions?.totems?.elements?.fire || FireTotem.NoFireTotem;
					break;
				}
				case CallTotem.Ancestors: {
					value = player.getSpecOptions().classOptions?.totems?.ancestors?.fire || FireTotem.NoFireTotem;
					break;
				}
				case CallTotem.Spirits: {
					value = player.getSpecOptions().classOptions?.totems?.spirits?.fire || FireTotem.NoFireTotem;
					break;
				}
			}
			return value;
		},
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					newOptions.classOptions!.totems.fire = newValue;
					if (!newOptions.classOptions?.totems.elements) newOptions.classOptions.totems.elements = TotemSet.create();
					newOptions.classOptions.totems.elements!.fire = newValue;
					break;
				}
				case CallTotem.Ancestors: {
					if (!newOptions.classOptions?.totems.ancestors) newOptions.classOptions.totems.ancestors = TotemSet.create();
					newOptions.classOptions.totems.ancestors!.fire = newValue;
					break;
				}
				case CallTotem.Spirits: {
					if (!newOptions.classOptions?.totems.spirits) newOptions.classOptions.totems.spirits = TotemSet.create();
					newOptions.classOptions.totems.spirits!.fire = newValue;
					break;
				}
			}
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
		getValue: (player: Player<ShamanSpecs>) => {
			const newOptions = player.getSpecOptions();
			let value = AirTotem.NoAirTotem;
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					value = player.getSpecOptions().classOptions?.totems?.elements?.air || AirTotem.NoAirTotem;
					break;
				}
				case CallTotem.Ancestors: {
					value = player.getSpecOptions().classOptions?.totems?.ancestors?.air || AirTotem.NoAirTotem;
					break;
				}
				case CallTotem.Spirits: {
					value = player.getSpecOptions().classOptions?.totems?.spirits?.air || AirTotem.NoAirTotem;
					break;
				}
			}
			return value;
		},
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.classOptions?.totems) newOptions.classOptions!.totems = ShamanTotems.create();
			switch (newOptions.classOptions?.call) {
				case CallTotem.Elements: {
					newOptions.classOptions!.totems.air = newValue;
					if (!newOptions.classOptions?.totems.elements) newOptions.classOptions.totems.elements = TotemSet.create();
					newOptions.classOptions.totems.elements!.air = newValue;
					break;
				}
				case CallTotem.Ancestors: {
					if (!newOptions.classOptions?.totems.ancestors) newOptions.classOptions.totems.ancestors = TotemSet.create();
					newOptions.classOptions.totems.ancestors!.air = newValue;
					break;
				}
				case CallTotem.Spirits: {
					if (!newOptions.classOptions?.totems.spirits) newOptions.classOptions.totems.spirits = TotemSet.create();
					newOptions.classOptions.totems.spirits!.air = newValue;
					break;
				}
			}
			player.setSpecOptions(eventID, newOptions);
		},
	});

	return contentBlock;
}
