// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.
import * as InputHelpers from '../../core/components/input_helpers.js';
import { Spec } from '../../core/proto/common.js';
import { WarriorSyncType } from "../../core/proto/warrior";


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

export const PrepullMastery = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFuryWarrior>({
	fieldName: 'prepullMastery',
	label: 'Prepull Mastery Rating',
	labelTooltip: 'Mastery rating in the prepull set equipped before entering combat. Only applies if value is greater than 0.',
});
