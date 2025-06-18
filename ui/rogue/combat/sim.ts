import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, WeaponType } from '../../core/proto/common';
import { RogueOptions_PoisonOptions } from '../../core/proto/rogue';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { Sim } from '../../core/sim';
import * as RogueInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecCombatRogue, {
	cssClass: 'combat-rogue-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Rogue),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
		Stat.StatExpertiseRating,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
	],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_MSV_GEARSET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.CBAT_STANDARD_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			return expCap;
		})(),
		softCapBreakpoints: (() => {
			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [7.5, 26.5],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0.3 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT, 0],
			});

			return [meleeHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.CombatTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			blessingOfKings: true,
			trueshotAura: true,
			swiftbladesCunning: true,
			legacyOfTheWhiteTiger: true,
			blessingOfMight: true,
			bloodlust: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			weakenedArmor: true,
			physicalVulnerability: true,
			masterPoisoner: true,
		}),
	},

	playerInputs: {
		inputs: [RogueInputs.ApplyPoisonsManually()],
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [RogueInputs.LethalPoison()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.CritBuff, 
		BuffDebuffInputs.AttackPowerBuff,
		BuffDebuffInputs.MasteryBuff,
		BuffDebuffInputs.StatsBuff,
		BuffDebuffInputs.AttackSpeedBuff,
		
		BuffDebuffInputs.MajorHasteBuff,
		BuffDebuffInputs.StormLashTotem,
		BuffDebuffInputs.Skullbanner,
		BuffDebuffInputs.ShatteringThrow,
		BuffDebuffInputs.TricksOfTheTrade,

		BuffDebuffInputs.SpellDamageDebuff, 
		BuffDebuffInputs.MajorArmorDebuff, 
		BuffDebuffInputs.PhysicalDamageDebuff
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
			RogueInputs.StartingComboPoints(),
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
		epWeights: [Presets.CBAT_STANDARD_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.CombatTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_COMBAT],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_GEARSET, Presets.P1_MSV_GEARSET, Presets.P1_T14_GEARSET],
	},

	autoRotation: (player: Player<Spec.SpecCombatRogue>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 2) {
			return Presets.ROTATION_PRESET_COMBAT.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_COMBAT.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecCombatRogue,
			talents: Presets.CombatTalents.data,
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
					1: Presets.P1_MSV_GEARSET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_MSV_GEARSET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class CombatRogueSimUI extends IndividualSimUI<Spec.SpecCombatRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecCombatRogue>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateGearStatsModifier: (baseStats: Stats) => {
					// Human/Orc racials for MH. Maxing Expertise for OH is a DPS loss when the MH matches the racial.
					const mhWepType = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType;
					const playerRace = player.getRace();
					if (
						(playerRace == Race.RaceHuman && (mhWepType == WeaponType.WeaponTypeSword || mhWepType == WeaponType.WeaponTypeMace)) ||
						(playerRace == Race.RaceOrc && (mhWepType == WeaponType.WeaponTypeAxe || mhWepType == WeaponType.WeaponTypeFist)) ||
						(playerRace == Race.RaceGnome && (mhWepType == WeaponType.WeaponTypeDagger || mhWepType == WeaponType.WeaponTypeSword))
					) {
						return baseStats.addStat(Stat.StatExpertiseRating, 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
					}
					return baseStats;
				},
			});
		});

		this.player.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			const encounter = this.sim.encounter;
			if (!options.classOptions!.applyPoisonsManually) {
				options.classOptions!.lethalPoison = RogueOptions_PoisonOptions.DeadlyPoison
			}
			this.player.setSpecOptions(c, options);
		});
		this.sim.encounter.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			const encounter = this.sim.encounter;
			if (!options.classOptions!.applyPoisonsManually) {
				options.classOptions!.lethalPoison = RogueOptions_PoisonOptions.DeadlyPoison
			}
			this.player.setSpecOptions(c, options);
		});
	}
}
