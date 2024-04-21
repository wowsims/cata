import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/other_inputs';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import {
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
	WeaponType,
} from '../../core/proto/common';
import { RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { Stats } from '../../core/proto_utils/stats';
import * as RogueInputs from '../inputs';
// import * as CombatInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecCombatRogue, {
	cssClass: 'combat-rogue-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Rogue),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Rotations are not fully optimized, especially for non-standard setups.'],
	// warnings: [
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.sim.encounter.changeEmitter,
	// 			getContent: () => {
	// 				let hasNoArmor = false;
	// 				for (const target of simUI.sim.encounter.targets) {
	// 					if (new Stats(target.stats).getStat(Stat.StatArmor) <= 0) {
	// 						hasNoArmor = true;
	// 						break;
	// 					}
	// 				}
	// 				if (hasNoArmor) {
	// 					return 'One or more targets have no armor. Check advanced encounter settings.';
	// 				} else {
	// 					return '';
	// 				}
	// 			},
	// 		};
	// 	},
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.player.changeEmitter,
	// 			getContent: () => {
	// 				if (
	// 					simUI.player.getTalents().mutilate &&
	// 					(simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType != WeaponType.WeaponTypeDagger ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType != WeaponType.WeaponTypeDagger)
	// 				) {
	// 					return '"Mutilate" talent selected, but daggers not equipped in both hands.';
	// 				} else {
	// 					return '';
	// 				}
	// 			},
	// 		};
	// 	},
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.player.changeEmitter,
	// 			getContent: () => {
	// 				if (simUI.player.getTalents().hackAndSlash) {
	// 					if (
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeSword ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeAxe ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeSword ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeAxe
	// 					) {
	// 						return '';
	// 					} else {
	// 						return '"Hack and Slash" talent selected, but swords or axes not equipped.';
	// 					}
	// 				} else {
	// 					return '';
	// 				}
	// 			},
	// 		};
	// 	},
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.player.changeEmitter,
	// 			getContent: () => {
	// 				if (simUI.player.getTalents().closeQuartersCombat) {
	// 					if (
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeFist ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeDagger ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeFist ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeDagger
	// 					) {
	// 						return '';
	// 					} else {
	// 						return '"Close Quarters Combat" talent selected, but fists or daggers not equipped.';
	// 					}
	// 				} else {
	// 					return '';
	// 				}
	// 			},
	// 		};
	// 	},
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.player.changeEmitter,
	// 			getContent: () => {
	// 				if (simUI.player.getTalents().maceSpecialization) {
	// 					if (
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType == WeaponType.WeaponTypeMace ||
	// 						simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponType == WeaponType.WeaponTypeMace
	// 					) {
	// 						return '';
	// 					} else {
	// 						return '"Mace Specialization" talent selected, but maces not equipped.';
	// 					}
	// 				} else {
	// 					return '';
	// 				}
	// 			},
	// 		};
	// 	},
	// 	(simUI: IndividualSimUI<Spec.SpecCombatRogue>) => {
	// 		return {
	// 			updateOn: simUI.player.changeEmitter,
	// 			getContent: () => {
	// 				const mhWeaponSpeed = simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponSpeed;
	// 				const ohWeaponSpeed = simUI.player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed;
	// 				const mhImbue = simUI.player.getSpecOptions().classOptions!.mhImbue;
	// 				const ohImbue = simUI.player.getSpecOptions().classOptions!.ohImbue;
	// 				if (
	// 					typeof mhWeaponSpeed == 'undefined' ||
	// 					typeof ohWeaponSpeed == 'undefined' ||
	// 					!simUI.player.getSpecOptions().classOptions!.applyPoisonsManually
	// 				) {
	// 					return '';
	// 				}
	// 				if (mhWeaponSpeed < ohWeaponSpeed && ohImbue == RogueOptions_PoisonImbue.DeadlyPoison) {
	// 					return 'Deadly poison applied to slower (off hand) weapon.';
	// 				}
	// 				if (ohWeaponSpeed < mhWeaponSpeed && mhImbue == RogueOptions_PoisonImbue.DeadlyPoison) {
	// 					return 'Deadly poison applied to slower (main hand) weapon.';
	// 				}
	// 				return '';
	// 			},
	// 		};
	// 	},
	// ],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
		Stat.StatExpertise,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatSpellHit,
		Stat.StatMeleeCrit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
		Stat.StatExpertise,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET_COMBAT.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatAgility]: 1.86,
				[Stat.StatStrength]: 1.14,
				[Stat.StatAttackPower]: 1,
				[Stat.StatSpellCrit]: 0.28,
				[Stat.StatSpellHit]: 0.08,
				[Stat.StatMeleeHit]: 1.39,
				[Stat.StatMeleeCrit]: 1.32,
				[Stat.StatMeleeHaste]: 1.48,
				[Stat.StatMastery]: 0.84,
				[Stat.StatExpertise]: 0.98,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 2.94,
				[PseudoStat.PseudoStatOffHandDps]: 2.45,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.CombatTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			bloodlust: true,
			markOfTheWild: true,
			icyTalons: true,
			moonkinForm: true,
			leaderOfThePack: true,
			powerWordFortitude: true,
			strengthOfEarthTotem: true,
			trueshotAura: true,
			wrathOfAirTotem: true,
			demonicPact: true,
			blessingOfKings: true,
			blessingOfMight: true,
			communion: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({
		}),
		debuffs: Debuffs.create({
			mangle: true,
			sunderArmor: true,
			shadowAndFlame: true,
			earthAndMoon: true,
			bloodFrenzy: true,
		}),
	},

	playerInputs: {
		inputs: [RogueInputs.ApplyPoisonsManually()],
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [RogueInputs.MainHandImbue(), RogueInputs.OffHandImbue(), RogueInputs.ThrownImbue()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.CritBuff,
		BuffDebuffInputs.SpellCritDebuff,
		BuffDebuffInputs.SpellDamageDebuff,
		BuffDebuffInputs.MajorArmorDebuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			// RogueInputs.StartingOverkillDuration(),
			// RogueInputs.VanishBreakTime(),
			// RogueInputs.AssumeBleedActive(),
			// OtherInputs.TankAssignment,
			// OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.CombatTalents],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_COMBAT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_PRESET_COMBAT,
		],
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
			consumes: Presets.DefaultConsumes,
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
		},
	],
});

export class CombatRogueSimUI extends IndividualSimUI<Spec.SpecCombatRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecCombatRogue>) {
		super(parentElem, player, SPEC_CONFIG);
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
