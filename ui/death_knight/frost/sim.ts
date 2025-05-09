import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as DeathKnightInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFrostDeathKnight, {
	cssClass: 'frost-death-knight-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.DeathKnight),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatArmor,
		Stat.StatAttackPower,
		Stat.StatExpertiseRating,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatPhysicalHitPercent,
		PseudoStat.PseudoStatSpellHitPercent,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatArmor, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),
	defaults: {
		// Default equipped gear.
		gear: Presets.P4_MASTERFROST_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_MASTERFROST_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats();
		})(),
		softCapBreakpoints: (() => {
			const physicalHitPercentSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [1 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT, 0],
			});

			const spellHitPercentSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHitPercent, {
				breakpoints: [17],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			});

			const expertiseRatingSoftCapConfig = StatCap.fromStat(Stat.StatExpertiseRating, {
				breakpoints: [6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			});

			return [physicalHitPercentSoftCapConfig, spellHitPercentSoftCapConfig, expertiseRatingSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.MasterfrostTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({}),
	},

	autoRotation: (player: Player<Spec.SpecFrostDeathKnight>): APLRotation => {
		const talents = player.getTalents();

		// if (talents.mightOfTheFrozenWastes > 0) {
		// 	return Presets.TWO_HAND_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		// } else if (talents.epidemic > 0) {
		// 	return Presets.DUAL_WIELD_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		// }

		return Presets.MASTERFROST_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellDamageDebuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.DamageReduction, BuffDebuffInputs.CastSpeedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DeathKnightInputs.StartingRunicPower(),
			// DeathKnightInputs.PetUptime(),
			// FrostInputs.UseAMSInput,
			// FrostInputs.AvgAMSSuccessRateInput,
			// FrostInputs.AvgAMSHitInput,
			// OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
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
		ItemSlot.ItemSlotRanged,
	],
	encounterPicker: {
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_MASTERFROST_EP_PRESET, Presets.P1_TWOHAND_EP_PRESET, Presets.P1_DUAL_WIELD_EP_PRESET, Presets.P3_DUAL_WIELD_EP_PRESET],
		talents: [Presets.DualWieldTalents, Presets.TwoHandTalents, Presets.MasterfrostTalents],
		rotations: [Presets.DUAL_WIELD_ROTATION_PRESET_DEFAULT, Presets.TWO_HAND_ROTATION_PRESET_DEFAULT, Presets.MASTERFROST_ROTATION_PRESET_DEFAULT],
		gear: [
			Presets.P1_DW_GEAR_PRESET,
			Presets.P3_DW_GEAR_PRESET,
			Presets.P1_2H_GEAR_PRESET,
			Presets.PREBIS_MASTERFROST_GEAR_PRESET,
			Presets.P1_MASTERFROST_GEAR_PRESET,
			Presets.P3_MASTERFROST_GEAR_PRESET,
			Presets.P4_MASTERFROST_GEAR_PRESET,
		],
		builds: [Presets.PRESET_BUILD_PREBIS, Presets.PRESET_BUILD_DW, Presets.PRESET_BUILD_2H, Presets.PRESET_BUILD_MASTERFROST],
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFrostDeathKnight,
			talents: Presets.MasterfrostTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_MASTERFROST_GEAR_PRESET.gear,
					2: Presets.P1_MASTERFROST_GEAR_PRESET.gear,
					3: Presets.P3_MASTERFROST_GEAR_PRESET.gear,
					4: Presets.P3_MASTERFROST_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_MASTERFROST_GEAR_PRESET.gear,
					2: Presets.P1_MASTERFROST_GEAR_PRESET.gear,
					3: Presets.P3_MASTERFROST_GEAR_PRESET.gear,
					4: Presets.P3_MASTERFROST_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class FrostDeathKnightSimUI extends IndividualSimUI<Spec.SpecFrostDeathKnight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFrostDeathKnight>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: (softCaps: StatCap[]) => {
					const talents = player.getTalents();
					// if (talents.mightOfTheFrozenWastes > 0 || talents.epidemic > 0) {
					// 	const physicalHitCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent));
					// 	if (physicalHitCap) {
					// 		physicalHitCap.breakpoints = [8];
					// 		physicalHitCap.postCapEPs = [0];
					// 	}

					// 	const spellHitCap = softCaps.findIndex(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHitPercent));
					// 	if (spellHitCap > -1) {
					// 		softCaps.splice(spellHitCap, 1);
					// 	}
					// }
					return softCaps;
				},
				// getEPDefaults: (player: Player<Spec.SpecFrostDeathKnight>) => {
				// 	const talents = player.getTalents();
				// 	return talents.mightOfTheFrozenWastes > 0
				// 		? Presets.P1_TWOHAND_EP_PRESET.epWeights
				// 		: talents.epidemic > 0
				// 		? Presets.P3_DUAL_WIELD_EP_PRESET.epWeights
				// 		: Presets.P1_MASTERFROST_EP_PRESET.epWeights;
				// },
			});
		});
	}
}
