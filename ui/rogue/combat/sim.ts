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
				breakpoints: [16.90, 16.95, 16.96, 16.97, 16.98, 16.99, 17],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0, 0, 0, 0, 0, 0],
			});

			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0],
			});

			const hasteRatingSoftCapConfig = StatCap.fromStat(Stat.StatHasteRating, {
				breakpoints: [2070, 2150, 2250, 2350, 2450],
				capType: StatCapType.TypeSoftCap,
				// These are set by the active EP weight in the updateSoftCaps callback
				postCapEPs: [0, 0, 0, 0, 0],
			})

			return [meleeHitSoftCapConfig, spellHitSoftCapConfig, hasteRatingSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
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
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			mangle: true,
			faerieFire: true,
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
			OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotHands, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.CBAT_STANDARD_EP_PRESET, Presets.CBAT_4PT12_EP_PRESET, Presets.CBAT_T13_EP_PRESET, Presets.CBAT_NOKALED_EP_PRESET],
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
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

// Check if the player is wearing any combination of the legendary dagger stages in both MH and OH
const hasAnyLegendaryStage = (player: Player<Spec.SpecCombatRogue>): boolean => {
	const mhWepId = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.id;
	const ohWepId = player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.id;

	return (mhWepId == 77945 || mhWepId == 77947 || mhWepId == 77949) &&
		   (ohWepId == 77946 || ohWepId == 77948 || ohWepId == 77950);
}

// Check if the player is wearing the final legendary stage
const hasFinalLegendaryStage = (player: Player<Spec.SpecCombatRogue>): boolean => {
	const mhWepId = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.id;
	const ohWepId = player.getEquippedItem(ItemSlot.ItemSlotOffHand)?.id;

	return mhWepId == 77949 && ohWepId == 77950;
}

const getActiveEPWeight = (player: Player<Spec.SpecCombatRogue>, sim: Sim): Stats => {
	if (sim.getUseCustomEPValues()) {
		return player.getEpWeights();
	} else {
		const mhWepId = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.id;
		const playerGear = player.getGear();
		let avgIlvl = 0;
		playerGear.asArray().forEach(v => avgIlvl += v?.item.ilvl || 0);
		avgIlvl /= playerGear.asArray().length;
		if (mhWepId == 78472 || mhWepId == 77188 || mhWepId == 78481) { // No'Kaled MH
			return Presets.CBAT_NOKALED_EP_PRESET.epWeights;
		} else if (playerGear.getItemSetCount("Vestments of the Dark Phoenix") >= 4) {
			return Presets.CBAT_4PT12_EP_PRESET.epWeights;
		} else if (hasAnyLegendaryStage(player)) {
			return Presets.CBAT_T13_EP_PRESET.epWeights;
		}  else {
			return Presets.CBAT_STANDARD_EP_PRESET.epWeights;
		}
	}
}

export class CombatRogueSimUI extends IndividualSimUI<Spec.SpecCombatRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecCombatRogue>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateSoftCaps: (softCaps: StatCap[]) => {
					const activeEPWeight = getActiveEPWeight(player, this.sim);
					const hasteEP = activeEPWeight.getStat(Stat.StatHasteRating);
					const hasteSoftCap = softCaps.find(v => v.unitStat.equalsStat(Stat.StatHasteRating));
					const hasAnyLego = hasAnyLegendaryStage(player)
					const hasFinalLego = hasFinalLegendaryStage(player)
					if (hasteSoftCap) {
						// If wearing either Fear or Sleeper in MH, Haste EP is never overtaken by Mastery
						if (hasAnyLego && !hasFinalLego)
							hasteSoftCap.postCapEPs = [hasteEP, hasteEP, hasteEP, hasteEP, hasteEP]
						else
							hasteSoftCap.postCapEPs = [hasteEP - 0.1, hasteEP - 0.2, hasteEP - 0.3, hasteEP - 0.4, hasteEP - 0.5];
					}

					// Dynamic adjustments to the static Hit soft cap EP
					const meleeSoftCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent));
					const spellSoftCap = softCaps.find(v => v.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHitPercent));
					if (meleeSoftCap) {
						const initialEP = activeEPWeight.getPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent);
						meleeSoftCap.postCapEPs = [initialEP/2, 0];
					}
					if (spellSoftCap) {
						const initialEP = activeEPWeight.getPseudoStat(PseudoStat.PseudoStatSpellHitPercent);
						spellSoftCap.postCapEPs = [initialEP/1.25, initialEP/2, initialEP/4, initialEP/8, initialEP/16, initialEP/32, 0];
					}

					return softCaps
				},
				getEPDefaults: (player: Player<Spec.SpecCombatRogue>) => {
					return getActiveEPWeight(player, this.sim);
				},
				updateGearStatsModifier: (baseStats: Stats) => {
					// Human/Orc racials for MH. Maxing Expertise for OH is a DPS loss when the MH matches the racial.
					const mhWepType = player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.weaponType;
					if ((player.getRace() == Race.RaceHuman && (mhWepType == WeaponType.WeaponTypeSword || mhWepType ==  WeaponType.WeaponTypeMace) ||
						(player.getRace() == Race.RaceOrc && mhWepType == WeaponType.WeaponTypeAxe)))
					{
						return baseStats.addStat(Stat.StatExpertiseRating, 90);
					}
					return baseStats
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
