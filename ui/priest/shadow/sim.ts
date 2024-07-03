import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, PartyBuffs, Race, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import * as PriestInputs from '../inputs';
// import * as ShadowPriestInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecShadowPriest, {
	cssClass: 'shadow-priest-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Priest),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Some items may display and use stats a litle higher than their original value.'],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatSpellHit, Stat.StatSpellCrit, Stat.StatSpellHaste, Stat.StatMastery],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatIntellect,
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
		Stat.StatMastery,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		statCaps: (() => {
			return new Stats().withStat(Stat.StatSpellHit, 17 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			// Picked from Priest Discord
			// Sources:
			// https://docs.google.com/spreadsheets/d/17cJJUReg2uz-XxBB3oDWb1kCncdH_-X96mSb0HAu4Ko/edit?gid=0#gid=0
			// https://docs.google.com/spreadsheets/d/1WLOZ1YevGPw_WZs0JhGzVVy906W5y0i9UqHa3ejyBkE/htmlview?gid=16
			const breakpoints = [
				801, // 9-tick DP
				1066, // 7-tick SWP
				1280, // 6-tick VT
				1560, // Filler post 6-tick VT
				2033, // Semi BP Filler post 6-tick VT
				2400, // 10-tick DP
				2666, // Filler post 10-tick DP
				2933, // Filler post 10-tick DP
				3199, // 8-tick SWP
				3844, // 7-tick VT
				4395, // Semi BP - Filler post 11-tick DP
				5004, // Filler post 11-tick DP
				5337, // 9-tick SWP
				5607, // 12-tick DP
				5871, // Filler post 12-tick DP
				6135, // Filler post 12-tick DP
				6399, // 8-tick VT
				6524, // Semi BP - Filler post 8-tick VT
				6939, // Filler post 8-tick VT
				7209, // 13-tick DP
				7473, // 10-tick SWP
				7807, // Filler post 10-tick SWP
				8474, // Filler post 10-tick SWP
			];
			const hasteSoftCapConfig = {
				stat: Stat.StatSpellHaste,
				breakpoints,
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [
					0.71, // 9-tick DP
					1.03, // 7-tick SWP
					0.65, // 6-tick VT
					0.8, // Filler post 6-tick VT
					0.55, // Semi BP Filler post 6-tick VT
					0.51, // 10-tick DP
					0.52, // Filler post 10-tick DP
					0.66, // Filler post 10-tick DP
					0.7, // 8-tick SWP
					0.69, // 7-tick VT
					0.58, // Semi BP - Filler post 11-tick DP
					0.68, // Filler post 11-tick DP
					0.41, // 9-tick SWP
					0.48, // 12-tick DP
					0.19, // Filler post 12-tick DP
					1.11, // Filler post 12-tick DP
					0.42, // 8-tick VT
					0.48, // Semi BP - Filler post 8-tick VT
					0.62, // Filler post 8-tick VT
					0.33, // 13-tick DP
					0, // 10-tick SWP
					0.36, // Filler post 10-tick SWP
					1.17, // Filler post 10-tick SWP
				],
			};

			return [hasteSoftCapConfig];
		})(),
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
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.ManaBuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.DarkIntentUptime,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		rotations: [Presets.ROTATION_PRESET_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRE_RAID, Presets.P1_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecShadowPriest>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecShadowPriest,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					0: Presets.PRE_RAID.gear,
					1: Presets.P1_PRESET.gear,
				},
				[Faction.Horde]: {
					0: Presets.PRE_RAID.gear,
					1: Presets.P1_PRESET.gear,
				},
			},
		},
	],
});

export class ShadowPriestSimUI extends IndividualSimUI<Spec.SpecShadowPriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecShadowPriest>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
