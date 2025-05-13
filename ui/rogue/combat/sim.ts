import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, WeaponType } from '../../core/proto/common';
import { RogueOptions_PoisonImbue } from '../../core/proto/rogue';
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
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
		Stat.StatExpertiseRating,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatPhysicalHitPercent,
		PseudoStat.PseudoStatSpellHitPercent,
	],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET_COMBAT.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.CBAT_STANDARD_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			return expCap;
		})(),
		softCapBreakpoints: (() => {
			// Running just under spell cap is typically preferrable to being over.
			const spellHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHitPercent, {
				breakpoints: [16.95, 16.96, 16.97, 16.98, 16.99, 17],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0, 0, 0, 0, 0],
			});

			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0],
			});

			return [meleeHitSoftCapConfig, spellHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.CombatTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// mangle: true,
			// faerieFire: true,
			// shadowAndFlame: true,
			// earthAndMoon: true,
			// bloodFrenzy: true,
		}),
	},

	playerInputs: {
		inputs: [RogueInputs.ApplyPoisonsManually()],
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [RogueInputs.MainHandImbue(), RogueInputs.OffHandImbue(), RogueInputs.ThrownImbue()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.CritBuff, BuffDebuffInputs.SpellDamageDebuff, BuffDebuffInputs.MajorArmorDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			// RogueInputs.StartingOverkillDuration(),
			// RogueInputs.VanishBreakTime(),
			// RogueInputs.AssumeBleedActive(),
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
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.CBAT_STANDARD_EP_PRESET, Presets.CBAT_4PT12_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.CombatTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_COMBAT],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET_COMBAT, Presets.P1_PRESET_COMBAT, Presets.P3_PRESET_COMBAT, Presets.P4_PRESET_COMBAT],
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
					1: Presets.P1_PRESET_COMBAT.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET_COMBAT.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

const getActiveEPWeight = (player: Player<Spec.SpecCombatRogue>, sim: Sim): Stats => {
	if (sim.getUseCustomEPValues()) {
		return player.getEpWeights();
	} else {
		const playerGear = player.getGear();
		if (playerGear.getItemSetCount('Vestments of the Dark Phoenix') >= 4) {
			return Presets.CBAT_4PT12_EP_PRESET.epWeights;
		} else {
			return Presets.CBAT_STANDARD_EP_PRESET.epWeights;
		}
	}
};

export class CombatRogueSimUI extends IndividualSimUI<Spec.SpecCombatRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecCombatRogue>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: (softCaps: StatCap[]) => {
					const activeEPWeight = getActiveEPWeight(player, this.sim);

					// Dynamic adjustments to the static Hit soft cap EP
					const meleeSoftCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent));
					if (meleeSoftCap) {
						const initialEP = activeEPWeight.getPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent);
						meleeSoftCap.postCapEPs = [initialEP * 0.6, 0];
					}

					return softCaps;
				},
				getEPDefaults: (player: Player<Spec.SpecCombatRogue>) => {
					return getActiveEPWeight(player, this.sim);
				},
				updateGearStatsModifier: (baseStats: Stats) => {
					// Human/Orc racials for MH. Maxing Expertise for OH is a DPS loss when the MH matches the racial.
					const mhWepType = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType;
					const playerRace = player.getRace();
					if (
						(playerRace == Race.RaceHuman && (mhWepType == WeaponType.WeaponTypeSword || mhWepType == WeaponType.WeaponTypeMace)) ||
						(playerRace == Race.RaceOrc && (mhWepType == WeaponType.WeaponTypeAxe || mhWepType == WeaponType.WeaponTypeFist)) ||
						(playerRace == Race.RaceGnome && (mhWepType == WeaponType.WeaponTypeDagger || mhWepType == WeaponType.WeaponTypeSword))
					) {
						return baseStats.addStat(Stat.StatExpertiseRating, 90);
					}
					return baseStats;
				},
			});
		});

		this.player.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			const encounter = this.sim.encounter;
			if (!options.classOptions!.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return;
				}
				if (encounter.targets.length > 3) {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
				} else {
					if (mhWeaponSpeed <= ohWeaponSpeed) {
						options.classOptions!.mhImbue = RogueOptions_PoisonImbue.DeadlyPoison;
						options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
					} else {
						options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
						options.classOptions!.ohImbue = RogueOptions_PoisonImbue.DeadlyPoison;
					}
				}
			}
			this.player.setSpecOptions(c, options);
		});
		this.sim.encounter.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			const encounter = this.sim.encounter;
			if (!options.classOptions!.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return;
				}
				if (encounter.targets.length > 3) {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
				} else {
					if (mhWeaponSpeed <= ohWeaponSpeed) {
						options.classOptions!.mhImbue = RogueOptions_PoisonImbue.DeadlyPoison;
						options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
					} else {
						options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
						options.classOptions!.ohImbue = RogueOptions_PoisonImbue.DeadlyPoison;
					}
				}
			}
			this.player.setSpecOptions(c, options);
		});
	}
}
