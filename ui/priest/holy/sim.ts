import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl.js';
import { Faction, PartyBuffs, Race, Spec, Stat } from '../../core/proto/common.js';
import { Stats } from '../../core/proto_utils/stats.js';
import * as PriestInputs from '../inputs';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHolyPriest, {
	cssClass: 'holy-priest-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Priest),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
		Stat.StatMastery,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMP5,
		Stat.StatMastery,
	],
	// modifyDisplayStats: (player: Player<Spec.SpecHolyPriest>) => {
	// 	let stats = new Stats();
	// 	stats = stats.addStat(Stat.StatSpellHit, player.getTalents().shadowFocus * 1 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

	// 	return {
	// 		talents: stats,
	// 	};
	// },

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.11,
			[Stat.StatSpirit]: 0.47,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellHit]: 0.87,
			[Stat.StatSpellCrit]: 0.74,
			[Stat.StatSpellHaste]: 1.65,
			[Stat.StatMP5]: 0.0,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PriestInputs.ArmorInput()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.MeleeHasteBuff,
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.ChannelClipDelay],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents, Presets.EnlightenmentTalents],
		rotations: [Presets.ROTATION_PRESET_DEFAULT, Presets.ROTATION_PRESET_AOE24, Presets.ROTATION_PRESET_AOE4PLUS],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P2_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecHolyPriest>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 4) {
			return Presets.ROTATION_PRESET_AOE4PLUS.rotation.rotation!;
		} else if (numTargets > 1) {
			return Presets.ROTATION_PRESET_AOE24.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecHolyPriest,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDwarf,
				[Faction.Horde]: Race.RaceUndead,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
			},
		},
	],
});

export class HolyPriestSimUI extends IndividualSimUI<Spec.SpecHolyPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHolyPriest>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
