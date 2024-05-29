import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/other_inputs.js';
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
} from '../../core/proto/common';
import { RogueOptions_PoisonImbue } from '../../core/proto/rogue';
import { Stats } from '../../core/proto_utils/stats';
import * as RogueInputs from '../inputs';
// import * as SinInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecAssassinationRogue, {
	cssClass: 'assassination-rogue-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Rogue),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: ['Rotations are not fully optimized, especially for non-standard setups.'],

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
		gear: Presets.P1_PRESET_ASSASSINATION.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatAgility]: 2.58,
				[Stat.StatStrength]: 1.05,
				[Stat.StatAttackPower]: 1,
				[Stat.StatSpellCrit]: 0.26,
				[Stat.StatSpellHit]: 1.31,
				[Stat.StatMeleeHit]: 0.70,
				[Stat.StatMeleeCrit]: 0.62,
				[Stat.StatMeleeHaste]: 1.10,
				[Stat.StatMastery]: 1.23,
				[Stat.StatExpertise]: 1.04,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 2.49,
				[PseudoStat.PseudoStatOffHandDps]: 1.0,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.AssassinationTalentsDefault.data,
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
			earthAndMoon: true,
			bloodFrenzy: true,
			shadowAndFlame: true,
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
		talents: [Presets.AssassinationTalentsDefault],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_MUTILATE,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_PRESET_ASSASSINATION,
			Presets.P1_PRESET_ASN_EXPERTISE
		],
	},

	autoRotation: (player: Player<Spec.SpecAssassinationRogue>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 5) {
			return Presets.ROTATION_PRESET_MUTILATE.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_MUTILATE.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecAssassinationRogue,
			talents: Presets.AssassinationTalentsDefault.data,
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
					1: Presets.P1_PRESET_ASSASSINATION.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET_ASSASSINATION.gear,
				},
			},
		},
	],
});

export class AssassinationRogueSimUI extends IndividualSimUI<Spec.SpecAssassinationRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecAssassinationRogue>) {
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
