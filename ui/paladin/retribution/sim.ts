import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation, APLRotation_Type } from '../../core/proto/apl.js';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common.js';
import { Stats, UnitStat } from '../../core/proto_utils/stats.js';
import { TypedEvent } from '../../core/typed_event.js';
import * as PaladinInputs from '../inputs.js';
import * as RetributionInputs from './inputs.js';
import * as Presets from './presets.js';

const isGlyphOfSealOfTruthActive = (player: Player<Spec.SpecRetributionPaladin>): boolean => {
	// const currentSeal = player.getSpecOptions().classOptions?.seal;
	// return (
	// 	player.getPrimeGlyps().includes(PaladinPrimeGlyph.GlyphOfSealOfTruth) &&
	// 	(currentSeal === PaladinSeal.Truth || currentSeal === PaladinSeal.Righteousness)
	// );
	return false;
};

const modifyDisplayStats = (player: Player<Spec.SpecRetributionPaladin>) => {
	let stats = new Stats();

	TypedEvent.freezeAllAndDo(() => {
		if (isGlyphOfSealOfTruthActive(player)) {
			stats = stats.addStat(Stat.StatExpertiseRating, 2.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
		}
	});

	return {
		talents: stats,
	};
};

const getStatCaps = () => {
	const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
	const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

	return hitCap.add(expCap);
};

const updateGearStatsModifier = (player: Player<Spec.SpecRetributionPaladin>) => (baseStats: Stats) => {
	if (isGlyphOfSealOfTruthActive(player)) {
		return baseStats.addStat(Stat.StatExpertiseRating, 2.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
	} else {
		return baseStats;
	}
};

const getEPDefaults = (player: Player<Spec.SpecRetributionPaladin>) => {
	let hasP3Setup = false;
	let hasP4Setup = false;

	const items = player.getGear().getEquippedItems();

	for (const item of items) {
		const phase = item?.item.phase || 0;
		if (phase > 3) {
			hasP4Setup = true;
		} else if (phase > 2) {
			hasP3Setup = true;
		}
	}

	return hasP4Setup ? Presets.P4_EP_PRESET.epWeights : hasP3Setup ? Presets.P3_EP_PRESET.epWeights : Presets.P2_EP_PRESET.epWeights;
};

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRetributionPaladin, {
	cssClass: 'retribution-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatExpertiseRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatIntellect,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatSpellPower,
			Stat.StatMana,
			Stat.StatHealth,
			Stat.StatMasteryRating,
		],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatSpellHastePercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatSpellHitPercent,
		],
	),
	modifyDisplayStats,

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_BIS_RET_PRESET.gear,
		// Default item swap set.
		itemSwap: Presets.ITEM_SWAP_4P_T11.itemSwap,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P4_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: getStatCaps(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// faerieFire: true,
			// bloodFrenzy: true,
			// mangle: true,
			// ebonPlaguebringer: true,
			// criticalMass: true,
		}),
		rotationType: APLRotation_Type.TypeAuto,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PaladinInputs.AuraSelection(), PaladinInputs.StartingSealSelection()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [RetributionInputs.StartingHolyPower(), OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P2_EP_PRESET, Presets.P3_EP_PRESET, Presets.P4_EP_PRESET],
		rotations: [Presets.ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.P2_BIS_RET_PRESET, Presets.P3_BIS_RET_PRESET, Presets.P4_BIS_RET_PRESET, Presets.PRERAID_RET_PRESET],
		itemSwaps: [Presets.ITEM_SWAP_4P_T11],
		builds: [Presets.P2_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
	},

	autoRotation: (_: Player<Spec.SpecRetributionPaladin>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRetributionPaladin,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceBloodElf,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P2_BIS_RET_PRESET.gear,
					2: Presets.P2_BIS_RET_PRESET.gear,
					3: Presets.P3_BIS_RET_PRESET.gear,
					4: Presets.P4_BIS_RET_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P2_BIS_RET_PRESET.gear,
					2: Presets.P2_BIS_RET_PRESET.gear,
					3: Presets.P3_BIS_RET_PRESET.gear,
					4: Presets.P4_BIS_RET_PRESET.gear,
				},
			},
		},
	],
});

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateGearStatsModifier: updateGearStatsModifier(player),
				getEPDefaults,
			});
		});
	}
}
