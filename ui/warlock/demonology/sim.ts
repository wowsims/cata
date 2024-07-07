import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, ItemSlot, PartyBuffs, Race, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import * as WarlockInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecDemonologyWarlock, {
	cssClass: 'demonology-warlock-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warlock),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpellPower, Stat.StatSpellHit, Stat.StatSpellCrit, Stat.StatSpellHaste, Stat.StatMastery],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMastery,
		Stat.StatMP5,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge optimizer
		statCaps: (() => {
			return new Stats().withStat(Stat.StatSpellHit, 17 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			// Set up Mastery breakpoints for integer % damage increments.
			// These should be removed once the bugfix to make Mastery
			// continuous goes live!
			const masteryRatingBreakpoints: number[] = [];

			for (let masteryPercent = 19; masteryPercent <= 200; masteryPercent++) {
				masteryRatingBreakpoints.push(
					(masteryPercent / Mechanics.masteryPercentPerPoint.get(Spec.SpecDemonologyWarlock)!) * Mechanics.MASTERY_RATING_PER_MASTERY_POINT,
				);
			}

			const masterySoftCapConfig = {
				stat: Stat.StatMastery,
				breakpoints: masteryRatingBreakpoints,
				capType: StatCapType.TypeThreshold,
				postCapEPs: Array(masteryRatingBreakpoints.length).fill(0),
			};

			const hasteSoftCapConfig = {
				stat: Stat.StatSpellHaste,
				breakpoints: [16.65 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT, 25 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0.64, 0.61],
			};

			return [hasteSoftCapConfig, masterySoftCapConfig];
		})(),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,

		// Default talents.
		talents: Presets.DemonologyTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,

		// Default buffs and debuffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({}),

		individualBuffs: Presets.DefaultIndividualBuffs,

		debuffs: Presets.DefaultDebuffs,

		other: Presets.OtherDefaults,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [WarlockInputs.PetInput()],

	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.MajorArmorDebuff,
		BuffDebuffInputs.PhysicalDamageDebuff,
		BuffDebuffInputs.MeleeHasteBuff,
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.AttackPowerPercentBuff,
		BuffDebuffInputs.StrengthAndAgilityBuff,
		BuffDebuffInputs.StaminaBuff,
	],
	excludeBuffDebuffInputs: [],
	petConsumeInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarlockInputs.DetonateSeed(),
			WarlockInputs.PrepullMastery,
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.DarkIntentUptime,
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand, ItemSlot.ItemSlotRanged],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DemonologyTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.APL_Default],

		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P4_WOTLK_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecDemonologyWarlock>): APLRotation => {
		return Presets.APL_Default.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecDemonologyWarlock,
			talents: Presets.DemonologyTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PRERAID_PRESET.gear,
					2: Presets.P1_PRESET.gear,
					3: Presets.P4_WOTLK_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PRERAID_PRESET.gear,
					2: Presets.P1_PRESET.gear,
					3: Presets.P4_WOTLK_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class DemonologyWarlockSimUI extends IndividualSimUI<Spec.SpecDemonologyWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDemonologyWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
