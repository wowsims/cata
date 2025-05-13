import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as WarlockInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecAfflictionWarlock, {
	cssClass: 'affliction-warlock-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warlock),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpellPower, Stat.StatMasteryRating, Stat.StatMP5],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET.gear,

		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.DEFAULT_EP_PRESET.epWeights,
		// Default stat caps for the Reforge optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 17);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			// Set up Mastery breakpoints for integer % damage increments.
			// These should be removed once the bugfix to make Mastery
			// continuous goes live!
			const masteryRatingBreakpoints = [];
			const masteryPercentPerPoint = Mechanics.masteryPercentPerPoint.get(Spec.SpecAfflictionWarlock)!;

			for (let masteryPercent = 14; masteryPercent <= 200; masteryPercent++) {
				masteryRatingBreakpoints.push((masteryPercent / masteryPercentPerPoint) * Mechanics.MASTERY_RATING_PER_MASTERY_POINT);
			}

			const masterySoftCapConfig = StatCap.fromStat(Stat.StatMasteryRating, {
				breakpoints: masteryRatingBreakpoints,
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0],
			});

			return [masterySoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,

		// Default talents.
		talents: Presets.AfflictionTalents.data,
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
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	petConsumeInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarlockInputs.DetonateSeed(),
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.ChannelClipDelay,
		],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.DEFAULT_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.AfflictionTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.APL_Default],

		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		itemSwaps: [Presets.P4_ITEM_SWAP],
	},

	autoRotation: (_player: Player<Spec.SpecAfflictionWarlock>): APLRotation => {
		return Presets.APL_Default.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecAfflictionWarlock,
			talents: Presets.AfflictionTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
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
					3: Presets.P3_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PRERAID_PRESET.gear,
					2: Presets.P1_PRESET.gear,
					3: Presets.P3_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class AfflictionWarlockSimUI extends IndividualSimUI<Spec.SpecAfflictionWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecAfflictionWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				statSelectionPresets: Presets.AFFLICTION_BREAKPOINTS,
			});
		});
	}
}
