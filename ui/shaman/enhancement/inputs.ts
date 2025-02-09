import * as InputHelpers from '../../core/components/input_helpers.js';
import { Spec } from '../../core/proto/common.js';
import { ShamanImbue, ShamanSyncType } from '../../core/proto/shaman.js';
import { ActionId } from '../../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueOh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Off Hand Enchant' },
		{ actionId: ActionId.fromSpellId(8232), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(8024), value: ShamanImbue.FlametongueWeapon },
		{ actionId: ActionId.fromSpellId(8033), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman, ShamanSyncType>({
	fieldName: 'syncType',
	label: 'Sync/Stagger Setting',
	labelTooltip: `Choose your sync or stagger option Perfect
		<ul>
			<li><div>Auto: Will auto pick sync options based on your weapons attack speeds</div></li>
			<li><div>None: No Sync or Staggering, used for mismatched weapon speeds</div></li>
			<li><div>Perfect Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>
			<li><div>Delayed Offhand: Adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window</div></li>
		</ul>`,
	values: [
		{ name: 'Automatic', value: ShamanSyncType.Auto },
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});

export const UsePrepull2PT10 = InputHelpers.makeClassOptionsBooleanInput<Spec.SpecEnhancementShaman>({
	fieldName: 'usePrepullEnh2PT10',
	label: 'Use Tier 10 2PC effect during prepull',
	labelTooltip: 'Will use the Tier 10 2PC effect during the prepull.',
	getValue: player => player.getClassOptions().usePrepullEnh2PT10,
	setValue: (eventID, player, newValue: boolean) => {
		const newOptions = player.getClassOptions();
		newOptions.usePrepullEnh2PT10 = newValue;

		player.setClassOptions(eventID, newOptions);
	},
});
