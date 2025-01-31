// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

import * as InputHelpers from '../../core/components/input_helpers.js';
import { Profession, Spec, Stat } from '../../core/proto/common.js';
import { WarriorSyncType } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecFuryWarrior, WarriorSyncType>({
	fieldName: 'syncType',
	label: 'Swing Sync Setting',
	labelTooltip: `Choose your sync option Perfect
		<ul>
			<li><div>None: No Sync used for mismatched weapon speeds</div></li>
			<li><div>Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>
		</ul>`,
	values: [
		{ name: 'None', value: WarriorSyncType.WarriorNoSync },
		{ name: 'Sync', value: WarriorSyncType.WarriorSyncMainhandOffhandSwings },
	],
});

export const AssumePrepullMasteryElixir = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFuryWarrior>({
	fieldName: 'useItemSwapBonusStats',
	label: 'Assume Prepull Mastery Elixir',
	labelTooltip: 'Enabling this assumes you are using a Mastery Elixir during prepull.',
	getValue: player => player.getSpecOptions().useItemSwapBonusStats,
	setValue: (eventID, player, newVal) => {
		const newMessage = player.getSpecOptions();
		newMessage.useItemSwapBonusStats = newVal;

		const bonusStats = newVal ? new Stats().withStat(Stat.StatMasteryRating, 225 + (player.hasProfession(Profession.Alchemy) ? 40 : 0)) : new Stats();
		player.itemSwapSettings.setBonusStats(eventID, bonusStats);
		player.setSpecOptions(eventID, newMessage);
	},
});
