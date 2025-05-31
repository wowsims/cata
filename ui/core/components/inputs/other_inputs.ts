import { Player } from '../../player.js';
import { UnitReference } from '../../proto/common.js';
import { emptyUnitReference } from '../../proto_utils/utils.js';
import { Sim } from '../../sim.js';
import { EventID } from '../../typed_event.js';
import { BooleanPicker } from '../pickers/boolean_picker.js';
import { EnumPicker } from '../pickers/enum_picker.js';

export function makeShow1hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	parent.classList.remove('hide');
	return new BooleanPicker<Sim>(parent, sim, {
		id: 'show-1h-weapons-selector',
		extraCssClasses: ['show-1h-weapons-selector', 'mb-0'],
		label: '1H',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().oneHandedWeapons,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.oneHandedWeapons = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShow2hWeaponsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	parent.classList.remove('hide');
	return new BooleanPicker<Sim>(parent, sim, {
		id: 'show-2h-weapons-selector',
		extraCssClasses: ['show-2h-weapons-selector', 'mb-0'],
		label: '2H',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().twoHandedWeapons,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.twoHandedWeapons = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShowMatchingGemsSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		id: 'show-matching-gems-selector',
		extraCssClasses: ['show-matching-gems-selector', 'input-inline', 'mb-0'],
		label: 'Match Socket',
		inline: true,
		changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
		getValue: (sim: Sim) => sim.getFilters().matchingGemsOnly,
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			const filters = sim.getFilters();
			filters.matchingGemsOnly = newValue;
			sim.setFilters(eventID, filters);
		},
	});
}

export function makeShowEPValuesSelector(parent: HTMLElement, sim: Sim): BooleanPicker<Sim> {
	return new BooleanPicker<Sim>(parent, sim, {
		id: 'show-ep-values-selector',
		extraCssClasses: ['show-ep-values-selector', 'input-inline', 'mb-0'],
		label: 'Show EP',
		inline: true,
		changedEvent: (sim: Sim) => sim.showEPValuesChangeEmitter,
		getValue: (sim: Sim) => sim.getShowEPValues(),
		setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
			sim.setShowEPValues(eventID, newValue);
		},
	});
}

export function makePhaseSelector(parent: HTMLElement, sim: Sim): EnumPicker<Sim> {
	return new EnumPicker<Sim>(parent, sim, {
		id: 'phase-selector',
		extraCssClasses: ['phase-selector'],
		values: [
			{ name: 'Phase 1 (Tier 14)', value: 1 },
			{ name: 'Phase 2 (Tier 15)', value: 2 },
			{ name: 'Phase 3 (Tier 16)', value: 3 },
		],
		changedEvent: (sim: Sim) => sim.phaseChangeEmitter,
		getValue: (sim: Sim) => sim.getPhase(),
		setValue: (eventID: EventID, sim: Sim, newValue: number) => {
			sim.setPhase(eventID, newValue);
		},
	});
}

export const InputDelay = {
	id: 'input-delay',
	type: 'number' as const,
	label: 'Input Delay',
	labelTooltip:
		"Player input delay, in milliseconds. Specifies the maximum delay on actions that cannot be spell queued, such as spell casts that are waiting on resource gains or waiting for a cooldown to expire. Also used with certain APL values (such as 'Aura Is Active With Reaction Time'). Roughly models the sum of reaction time + server latency.",
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getReactionTime(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setReactionTime(eventID, newValue);
	},
};

export const ChallengeMode = {
	id: 'challenge-mode',
	type: 'boolean' as const,
	label: 'Challenge Mode',
	labelTooltip: 'Enables Challenge Mode',
	changedEvent: (player: Player<any>) => player.challengeModeChangeEmitter,
	getValue: (player: Player<any>) => player.getChallengeModeEnabled(),
	setValue: (eventID: EventID, player: Player<any>, value: boolean) => {
		player.setChallengeModeEnabled(eventID, value);
	},
};

export const ChannelClipDelay = {
	id: 'channel-clip-delay',
	type: 'number' as const,
	label: 'Channel Clip Delay',
	labelTooltip:
		'Clip delay following channeled spells, in milliseconds. This delay occurs following any full or partial channel ending after the GCD becomes available, due to the player not being able to queue the next spell.',
	changedEvent: (player: Player<any>) => player.miscOptionsChangeEmitter,
	getValue: (player: Player<any>) => player.getChannelClipDelay(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setChannelClipDelay(eventID, newValue);
	},
};

export const InFrontOfTarget = {
	id: 'in-front-of-target',
	type: 'boolean' as const,
	label: 'In Front of Target',
	labelTooltip: 'Stand in front of the target, causing Blocks and Parries to be included in the attack table.',
	changedEvent: (player: Player<any>) => player.inFrontOfTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getInFrontOfTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: boolean) => {
		player.setInFrontOfTarget(eventID, newValue);
	},
};

export const DistanceFromTarget = {
	id: 'distance-from-target',
	type: 'number' as const,
	label: 'Distance From Target',
	labelTooltip: 'Distance from targets, in yards. Used to calculate travel time for certain spells.',
	changedEvent: (player: Player<any>) => player.distanceFromTargetChangeEmitter,
	getValue: (player: Player<any>) => player.getDistanceFromTarget(),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		player.setDistanceFromTarget(eventID, newValue);
	},
};

export const TankAssignment = {
	id: 'tank-assignment',
	type: 'enum' as const,
	extraCssClasses: ['tank-selector', 'threat-metrics', 'within-raid-sim-hide'],
	label: 'Tank Assignment',
	labelTooltip:
		'Determines which mobs will be tanked. Most mobs default to targeting the Main Tank, but in preset multi-target encounters this is not always true.',
	values: [
		{ name: 'None', value: -1 },
		{ name: 'Main Tank', value: 0 },
		{ name: 'Tank 2', value: 1 },
		{ name: 'Tank 3', value: 2 },
		{ name: 'Tank 4', value: 3 },
	],
	changedEvent: (player: Player<any>) => player.getRaid()!.tanksChangeEmitter,
	getValue: (player: Player<any>) => (player.getRaid()?.getTanks() || []).findIndex(tank => UnitReference.equals(tank, player.makeUnitReference())),
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const newTanks = [];
		if (newValue != -1) {
			for (let i = 0; i < newValue; i++) {
				newTanks.push(emptyUnitReference());
			}
			newTanks.push(player.makeUnitReference());
		}
		player.getRaid()!.setTanks(eventID, newTanks);
	},
};

export const IncomingHps = {
	id: 'incoming-hps',
	type: 'number' as const,
	label: 'Incoming HPS',
	labelTooltip: `
		<p>Average amount of healing received per second. Used for calculating chance of death.</p>
		<p>If set to 0, defaults to 25% of the primary target's base DPS.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().hps,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.hps = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HealingCadence = {
	id: 'healing-cadence',
	type: 'number' as const,
	float: true,
	label: 'Healing Cadence',
	labelTooltip: `
		<p>How often the incoming heal 'ticks', in seconds. Generally, longer durations favor Effective Hit Points (EHP) for minimizing Chance of Death, while shorter durations favor avoidance.</p>
		<p>Example: if Incoming HPS is set to 1000 and this is set to 1s, then every 1s a heal will be received for 1000. If this is instead set to 2s, then every 2s a heal will be recieved for 2000.</p>
		<p>If set to 0, default values for Healing Cadence and Cadence +/- are inferred from boss damage parameters.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().cadenceSeconds,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.cadenceSeconds = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HealingCadenceVariation = {
	id: 'healing-cadence-variation',
	type: 'number' as const,
	float: true,
	label: 'Cadence +/-',
	labelTooltip: `
		<p>Magnitude of random variation in healing intervals, in seconds.</p>
		<p>Example: if Healing Cadence is set to 1s with 0.5s variation, then the interval between successive heals will vary uniformly between 0.5 and 1.5s. If the variation is instead set to 2s, then 50% of healing intervals will fall between 0s and 1s, and the other 50% will fall between 1s and 3s.</p>
		<p>The amount of healing per 'tick' is automatically scaled up or down based on the randomized time since the last tick, so as to keep HPS constant.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().cadenceVariation,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.cadenceVariation = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const AbsorbFrac = {
	id: 'healing-model-absorb-frac',
	type: 'number' as const,
	float: true,
	label: 'Absorb %',
	labelTooltip: `
		<p>% of each incoming heal 'tick' to model as an absorb shield rather than as a direct heal.</p>
	`,
	changedEvent: (player: Player<any>) => player.healingModelChangeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().absorbFrac * 100,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.absorbFrac = newValue / 100;
		player.setHealingModel(eventID, healingModel);
	},
};

export const BurstWindow = {
	id: 'burst-window',
	type: 'number' as const,
	float: false,
	label: 'TMI Burst Window',
	labelTooltip: `
		<p>Size in whole seconds of the burst window for calculating TMI. It is important to use a consistent setting when comparing this metric.</p>
		<p>Default is 6 seconds. If set to 0, TMI calculations are disabled.</p>
	`,
	changedEvent: (player: Player<any>) => player.getRaid()!.changeEmitter,
	getValue: (player: Player<any>) => player.getHealingModel().burstWindow,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const healingModel = player.getHealingModel();
		healingModel.burstWindow = newValue;
		player.setHealingModel(eventID, healingModel);
	},
	enableWhen: (player: Player<any>) => (player.getRaid()?.getTanks() || []).find(tank => UnitReference.equals(tank, player.makeUnitReference())) != null,
};

export const HpPercentForDefensives = {
	id: 'hp-percent-for-defensives',
	type: 'number' as const,
	float: true,
	label: 'HP % for Defensive CDs',
	labelTooltip: `
		<p>% of Maximum Health, below which defensive cooldowns are allowed to be used.</p>
		<p>If set to 0, this restriction is disabled.</p>
	`,
	changedEvent: (player: Player<any>) => player.rotationChangeEmitter,
	getValue: (player: Player<any>) => player.getSimpleCooldowns().hpPercentForDefensives * 100,
	setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
		const cooldowns = player.getSimpleCooldowns();
		cooldowns.hpPercentForDefensives = newValue / 100;
		player.setSimpleCooldowns(eventID, cooldowns);
	},
};
