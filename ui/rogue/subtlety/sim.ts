import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as RogueInputs from '../inputs';
import * as SubInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecSubtletyRogue, {
	cssClass: 'subtlety-rogue-sim-ui',
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
		gear: Presets.P4_PRESET_SUB.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Stat caps for reforge optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			return expCap;
		})(),
		softCapBreakpoints: (() => {
			const spellHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHitPercent, {
				breakpoints: [17],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			});

			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [98.02, 0],
			});

			return [meleeHitSoftCapConfig, spellHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.SubtletyTalents.data,
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
			RogueInputs.AssumeBleedActive(),
			SubInputs.HonorAmongThievesCritRate,
			// OtherInputs.TankAssignment,
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
		epWeights: [Presets.P1_EP_PRESET, Presets.P4_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.SubtletyTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SUBTLETY, Presets.ROTATION_PRESET_SUBTLETY_MASTERY],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET_SUB, Presets.P1_PRESET_SUB, Presets.P3_PRESET_SUB, Presets.P4_PRESET_SUB],
	},

	autoRotation: (player: Player<Spec.SpecSubtletyRogue>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 5) {
			return Presets.ROTATION_PRESET_SUBTLETY.rotation.rotation!;
		} else {
			// TODO: Need a sub rotation here
			return Presets.ROTATION_PRESET_SUBTLETY.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecSubtletyRogue,
			talents: Presets.SubtletyTalents.data,
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
					1: Presets.P1_PRESET_SUB.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET_SUB.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class SubtletyRogueSimUI extends IndividualSimUI<Spec.SpecSubtletyRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecSubtletyRogue>) {
		super(parentElem, player, SPEC_CONFIG);

		// Auto Reforging
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecSubtletyRogue>) => {
					if (player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.id == 77949) {
						return Presets.P4_EP_PRESET.epWeights;
					} else {
						return Presets.P1_EP_PRESET.epWeights;
					}
				},
			});
		});

		this.player.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			if (!options.classOptions!.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return;
				}
				if (mhWeaponSpeed <= ohWeaponSpeed) {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.DeadlyPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
				} else {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.DeadlyPoison;
				}
			}
			this.player.setSpecOptions(c, options);
		});
		this.sim.encounter.changeEmitter.on(c => {
			const options = this.player.getSpecOptions();
			if (!options.classOptions!.applyPoisonsManually) {
				const mhWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
				const ohWeaponSpeed = this.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
				if (typeof mhWeaponSpeed == 'undefined' || typeof ohWeaponSpeed == 'undefined') {
					return;
				}
				if (mhWeaponSpeed <= ohWeaponSpeed) {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.DeadlyPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.InstantPoison;
				} else {
					options.classOptions!.mhImbue = RogueOptions_PoisonImbue.InstantPoison;
					options.classOptions!.ohImbue = RogueOptions_PoisonImbue.DeadlyPoison;
				}
			}
			this.player.setSpecOptions(c, options);
		});
	}
}
