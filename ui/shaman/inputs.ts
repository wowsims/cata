import { ContentBlock } from '../core/components/content_block';
import { IndividualSimUI } from '../core/individual_sim_ui';
import { Input } from '../core/components/input';
import * as InputHelpers from '../core/components/input_helpers';
import { ShamanImbue, ShamanShield} from '../core/proto/shaman';
import { ActionId } from '../core/proto_utils/action_id';
import { ShamanSpecs } from '../core/proto_utils/utils';
import { Player } from '../core/player';
import { EventID } from '../core/typed_event';
import { buildIconInput } from '../core/components/icon_inputs';

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

	const feleAbilities = Input.newGroupContainer();
	feleAbilities.classList.add('totem-dropdowns-container', 'icon-group');

	contentBlock.bodyElement.appendChild(feleAbilities);

	const _fireBlastPicker = <SpecType extends ShamanSpecs>() =>
		InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
			fieldName: 'feleAutocast',
			id: ActionId.fromSpellId(57984),
			getValue: (player: Player<SpecType>) => player.getClassOptions().feleAutocast!.autocastFireblast,
			setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => {
				const newOptions = player.getClassOptions();
				newOptions.feleAutocast!.autocastFireblast = newValue
				player.setClassOptions(eventID, newOptions);
			},
		});

		const _fireNovaPicker = <SpecType extends ShamanSpecs>() =>
		InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
			fieldName: 'feleAutocast',
			id: ActionId.fromSpellId(117588),
			getValue: (player: Player<SpecType>) => player.getClassOptions().feleAutocast!.autocastFirenova,
			setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => {
				const newOptions = player.getClassOptions();
				newOptions.feleAutocast!.autocastFirenova = newValue
				player.setClassOptions(eventID, newOptions);
			},
		});

		const _ImmolationPicker = <SpecType extends ShamanSpecs>() =>
		InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
			fieldName: 'feleAutocast',
			id: ActionId.fromSpellId(118297),
			getValue: (player: Player<SpecType>) => player.getClassOptions().feleAutocast!.autocastImmolate,
			setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => {
				const newOptions = player.getClassOptions();
				newOptions.feleAutocast!.autocastImmolate = newValue
				player.setClassOptions(eventID, newOptions);
			},
		});

		const _EmpowerPicker = <SpecType extends ShamanSpecs>() =>
		InputHelpers.makeClassOptionsBooleanIconInput<SpecType>({
			fieldName: 'feleAutocast',
			id: ActionId.fromSpellId(118350),
			getValue: (player: Player<SpecType>) => player.getClassOptions().feleAutocast!.autocastEmpower,
			setValue: (eventID: EventID, player: Player<SpecType>, newValue: boolean) => {
				const newOptions = player.getClassOptions();
				newOptions.feleAutocast!.autocastEmpower = newValue
				player.setClassOptions(eventID, newOptions);
			},
		});

	buildIconInput(feleAbilities, simUI.player, _fireBlastPicker())
	buildIconInput(feleAbilities, simUI.player, _fireNovaPicker())
	buildIconInput(feleAbilities, simUI.player, _ImmolationPicker())
	buildIconInput(feleAbilities, simUI.player, _EmpowerPicker())

	return contentBlock;
}