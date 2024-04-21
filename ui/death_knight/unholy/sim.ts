import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/other_inputs';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import {
	Debuffs,
	Faction,
	HandType,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../../core/proto/common';
import { Stats } from '../../core/proto_utils/stats';
import * as DeathKnightInputs from '../inputs';
import * as UnholyInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecUnholyDeathKnight, {
	cssClass: 'unholy-death-knight-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.DeathKnight),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatArmor,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatArmorPenetration,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatStrength,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
		Stat.StatExpertise,
	],
	defaults: {
		// Default equipped gear.
		gear: Presets.P1_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 3.22,
				[Stat.StatArmor]: 0.01,
				[Stat.StatAttackPower]: 1,
				[Stat.StatExpertise]: 1.13,
				[Stat.StatMeleeHaste]: 1.85,
				[Stat.StatMeleeHit]: 1.92,
				[Stat.StatMeleeCrit]: 0.76,
				[Stat.StatArmorPenetration]: 0.77,
				[Stat.StatSpellHit]: 0.8,
				[Stat.StatSpellCrit]: 0.34,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 3.1,
				[PseudoStat.PseudoStatOffHandDps]: 1.79,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.SingleTargetTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			devotionAura: true,
			bloodlust: true,
			markOfTheWild: true,
			icyTalons: true,
			leaderOfThePack: true,
			powerWordFortitude: true,
			hornOfWinter: true,
			abominationsMight: true,
			arcaneTactics: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({
		}),
		debuffs: Debuffs.create({
			sunderArmor: true,
			brittleBones: true,
			ebonPlaguebringer: true,
			shadowAndFlame: true
		}),
	},

	autoRotation: (player: Player<Spec.SpecUnholyDeathKnight>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 1) {
			return Presets.AOE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		} else {
			return Presets.SINGLE_TARGET_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellDamageDebuff
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.DamageReduction,
		BuffDebuffInputs.MeleeAttackSpeedDebuff,
		BuffDebuffInputs.BleedDebuff
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			// DeathKnightInputs.StartingRunicPower(),
			// DeathKnightInputs.PetUptime(),
			// UnholyInputs.SelfUnholyFrenzy,
			// UnholyInputs.UseAMSInput,
			// UnholyInputs.AvgAMSSuccessRateInput,
			// UnholyInputs.AvgAMSHitInput,
			// OtherInputs.TankAssignment,
			// OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
			OtherInputs.DarkIntentUptime
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.SingleTargetTalents,
			//Presets.AoeTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.SINGLE_TARGET_ROTATION_PRESET_DEFAULT,
			//Presets.AOE_ROTATION_PRESET_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P1_GEAR_PRESET,
		],
	},

	raidSimPresets: [
		{
			spec: Spec.SpecUnholyDeathKnight,
			talents: Presets.SingleTargetTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_GEAR_PRESET.gear,
					2: Presets.P1_GEAR_PRESET.gear,
					3: Presets.P1_GEAR_PRESET.gear,
					4: Presets.P1_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_GEAR_PRESET.gear,
					2: Presets.P1_GEAR_PRESET.gear,
					3: Presets.P1_GEAR_PRESET.gear,
					4: Presets.P1_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class UnholyDeathKnightSimUI extends IndividualSimUI<Spec.SpecUnholyDeathKnight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecUnholyDeathKnight>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
