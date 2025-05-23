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
import * as DemonologyInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecDemonologyWarlock, {
	cssClass: 'demonology-warlock-sim-ui',
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
		gear: Presets.P3_PRESET.gear,
		itemSwap: Presets.P3_ITEM_SWAP.itemSwap,

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
			const masteryPercentPerPoint = Mechanics.masteryPercentPerPoint.get(Spec.SpecDemonologyWarlock)!;
			for (let masteryPercent = 19; masteryPercent <= 200; masteryPercent++) {
				masteryRatingBreakpoints.push((masteryPercent / masteryPercentPerPoint) * Mechanics.MASTERY_RATING_PER_MASTERY_POINT);
			}

			const masterySoftCapConfig = StatCap.fromStat(Stat.StatMasteryRating, {
				breakpoints: masteryRatingBreakpoints,
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0],
			});

			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent, {
				breakpoints: [16.65, 25],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0.64 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT, 0.61 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			return [hasteSoftCapConfig, masterySoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,

		// Default talents.
		talents: Presets.DemonologyTalentsIncinerate.data,
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
			DemonologyInputs.AssumePrepullMasteryElixir,
		],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotNeck,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotBack,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotWrist,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotWaist,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotFeet,
		ItemSlot.ItemSlotFinger1,
		ItemSlot.ItemSlotFinger2,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.DEFAULT_EP_PRESET, Presets.Mastery_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DemonologyTalentsShadowBolt, Presets.DemonologyTalentsIncinerate],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.APL_ShadowBolt, Presets.APL_Incinerate],

		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		itemSwaps: [Presets.P3_ITEM_SWAP, Presets.P4_ITEM_SWAP],

		builds: [Presets.PRESET_BUILD_SHADOWBOLT, Presets.PRESET_BUILD_INCINERATE],
	},

	autoRotation: (_player: Player<Spec.SpecDemonologyWarlock>): APLRotation => {
		return Presets.APL_Incinerate.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecDemonologyWarlock,
			talents: Presets.DemonologyTalentsIncinerate.data,
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

export class DemonologyWarlockSimUI extends IndividualSimUI<Spec.SpecDemonologyWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDemonologyWarlock>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				statSelectionPresets: Presets.DEMONOLOGY_BREAKPOINTS,
				getEPDefaults: (player: Player<Spec.SpecFuryWarrior>) => {
					const playerWeights = player.getEpWeights();
					const defaultWeights = Presets.DEFAULT_EP_PRESET.epWeights;

					if (this.individualConfig.presets.epWeights.some(preset => playerWeights.equals(preset.epWeights))) return playerWeights;

					return defaultWeights;
				},
			});
		});
	}
}
