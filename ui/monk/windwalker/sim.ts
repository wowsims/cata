import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, HandType, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { Sim } from '../../core/sim';
import { TypedEvent } from '../../core/typed_event';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecWindwalkerMonk, {
	cssClass: 'windwalker-monk-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Monk),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatExpertiseRating,
		Stat.StatMasteryRating,
		Stat.StatAttackPower,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps, PseudoStat.PseudoStatPhysicalHitPercent],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatStrength, Stat.StatAgility, Stat.StatAttackPower, Stat.StatExpertiseRating, Stat.StatMasteryRating],
		[
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatSpellHastePercent,
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),
	modifyDisplayStats: (player: Player<Spec.SpecWindwalkerMonk>) => {
		let stats = new Stats();

		TypedEvent.freezeAllAndDo(() => {
			if (player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item?.handType === HandType.HandTypeTwoHand) {
				stats = stats.addPseudoStat(PseudoStat.PseudoStatMeleeHastePercent, 40);
			}
		});

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.PREPATCH_DW_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.PREPATCH_DW_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			return expCap;
		})(),
		softCapBreakpoints: (() => {
			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0],
			});

			return [meleeHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			curseOfElements:true,
			physicalVulnerability: true,
			weakenedArmor: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.CritBuff, BuffDebuffInputs.MajorArmorDebuff, BuffDebuffInputs.SpellHasteBuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InFrontOfTarget, OtherInputs.InputDelay],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.PREPATCH_DW_EP_PRESET, Presets.PREPATCH_2H_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.PREPATCH_ROTATION_PRESET],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PREPATCH_DW_GEAR_PRESET, Presets.PREPATCH_2H_GEAR_PRESET],
	},

	autoRotation: (_: Player<Spec.SpecWindwalkerMonk>): APLRotation => {
		return Presets.PREPATCH_ROTATION_PRESET.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecWindwalkerMonk,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceAlliancePandaren,
				[Faction.Horde]: Race.RaceHordePandaren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					2: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					3: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					4: Presets.PREPATCH_DW_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					2: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					3: Presets.PREPATCH_DW_GEAR_PRESET.gear,
					4: Presets.PREPATCH_DW_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

const getActiveEPWeight = (player: Player<Spec.SpecWindwalkerMonk>, sim: Sim): Stats => {
	if (sim.getUseCustomEPValues()) {
		return player.getEpWeights();
	} else if (player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item?.handType === HandType.HandTypeTwoHand) {
		return Presets.PREPATCH_2H_EP_PRESET.epWeights;
	} else {
		return Presets.PREPATCH_DW_EP_PRESET.epWeights;
	}
};

export class WindwalkerMonkSimUI extends IndividualSimUI<Spec.SpecWindwalkerMonk> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWindwalkerMonk>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: (softCaps: StatCap[]) => {
					// Dynamic adjustments to the static Hit soft cap EP
					const meleeSoftCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent));
					if (meleeSoftCap) {
						const activeEPWeight = getActiveEPWeight(player, this.sim);
						const initialEP = activeEPWeight.getPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent);
						const mhWep = player.getEquippedItem(ItemSlot.ItemSlotMainHand);
						const ohWep = player.getEquippedItem(ItemSlot.ItemSlotOffHand);
						if (mhWep?.item.handType === HandType.HandTypeTwoHand || !ohWep) {
							meleeSoftCap.breakpoints = [meleeSoftCap.breakpoints[0]];
							meleeSoftCap.postCapEPs = [0];
						} else if (ohWep) {
							meleeSoftCap.postCapEPs = [initialEP / 5.5, 0];
						}
					}

					return softCaps;
				},
				getEPDefaults: (player: Player<Spec.SpecWindwalkerMonk>) => {
					return getActiveEPWeight(player, this.sim);
				},
			});
		});
	}
}
