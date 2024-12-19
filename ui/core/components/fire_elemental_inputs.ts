import * as InputHelpers from '../components/input_helpers.js';
import { IndividualSimUI } from '../individual_sim_ui';
import { Player } from '../player';
import { ShamanTotems } from '../proto/shaman';
import { ActionId } from '../proto_utils/action_id.js';
import { ShamanSpecs } from '../proto_utils/utils';
import { EventID } from '../typed_event';
import { ContentBlock } from './content_block';
import { Input } from './input';
import { BooleanPicker } from './pickers/boolean_picker.js';
import { IconPicker } from './pickers/icon_picker.jsx';
import { NumberPicker } from './pickers/number_picker.js';

export function FireElementalSection(parentElem: HTMLElement, simUI: IndividualSimUI<any>): ContentBlock {
	const contentBlock = new ContentBlock(parentElem, 'fire-elemental-settings', {
		header: { title: 'Fire Elemental' },
	});

	const fireElementalIconContainer = Input.newGroupContainer();
	fireElementalIconContainer.classList.add('fire-elemental-icon-container');

	contentBlock.bodyElement.appendChild(fireElementalIconContainer);

	const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>(
		{
			getModObject: (player: Player<ShamanSpecs>) => player,
			getValue: (player: Player<ShamanSpecs>) => player.getClassOptions().totems || ShamanTotems.create(),
			setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
				const newOptions = player.getClassOptions();
				newOptions.totems = newVal;
				player.setClassOptions(eventID, newOptions);

				// Hacky fix ItemSwapping is in the Rotation proto, this will let the Rotation know to update showWhen
				// TODO move the ItemSwap enabled to a spec option and have the ItemSwap proto be apart of player.
				player.rotationChangeEmitter.emit(eventID);
			},
			changeEmitter: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		},
		ActionId.fromSpellId(2894),
		'useFireElemental',
	);

	new IconPicker(fireElementalIconContainer, simUI.player, fireElementalBooleanIconInput);

	new BooleanPicker(contentBlock.bodyElement, simUI.player, {
		id: 'fire-elemental-use-tier-ten',
		label: 'Use Tier 10 (4pc)',
		labelTooltip: 'Will use Tier 10 (4pc) to snapshot Fire Elemental.',
		inline: true,
		getValue: (player: Player<ShamanSpecs>) => player.getClassOptions().totems?.enhTierTenBonus || false,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: boolean) => {
			const newOptions = player.getClassOptions();

			if (newOptions.totems) {
				newOptions.totems.enhTierTenBonus = newVal;
			}

			player.setClassOptions(eventID, newOptions);
		},
		changedEvent: (player: Player<ShamanSpecs>) => player.currentStatsEmitter,
		showWhen: (player: Player<ShamanSpecs>) => {
			const hasBonus = player.getCurrentStats().sets.includes("Frost Witch's Battlegear (4pc)");
			return hasBonus;
		},
	});

	return contentBlock;
}
