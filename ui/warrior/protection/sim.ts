import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation } from '../../core/proto/apl.js';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common.js';
import { ProtectionWarrior_Rotation as ProtectionWarriorRotation } from '../../core/proto/warrior.js';
import * as AplUtils from '../../core/proto_utils/apl_utils.js';
import { Stats } from '../../core/proto_utils/stats.js';
import * as ProtectionWarriorInputs from '../inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionWarrior, {
	cssClass: 'protection-warrior-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warrior),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatBlock,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMastery,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatBlock,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMastery,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BALANCED_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatArmor]: 0.174,
				[Stat.StatBonusArmor]: 0.155,
				[Stat.StatStamina]: 2.336,
				[Stat.StatStrength]: 1.555,
				[Stat.StatAgility]: 2.771,
				[Stat.StatAttackPower]: 0.32,
				[Stat.StatExpertise]: 1.44,
				[Stat.StatMeleeHit]: 1.432,
				[Stat.StatMeleeCrit]: 0.925,
				[Stat.StatMeleeHaste]: 0.431,
				[Stat.StatBlock]: 1.32,
				[Stat.StatDodge]: 2.606,
				[Stat.StatParry]: 2.649,
				// @todo: Calculate actual weights
				// This probably applies for all weights
				[Stat.StatMastery]: 0,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 6.081,
			},
		),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.StandardTalents.data,
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
			devotionAura: true,
			retributionAura: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			sunderArmor: true,
			mangle: true,
			vindication: true,
			bloodFrenzy: true,
			judgement: true,
			frostFever: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [ProtectionWarriorInputs.ShoutPicker(), ProtectionWarriorInputs.ShatteringThrow(),ProtectionWarriorInputs.Recklessness()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.StaminaBuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InspirationUptime,
			ProtectionWarriorInputs.StartingRage(),
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents, Presets.StandardTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_DEFAULT, Presets.ROTATION_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_BALANCED_PRESET,
			Presets.P1_BALANCED_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecProtectionWarrior>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecProtectionWarrior,
			talents: Presets.StandardTalents.data,
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
					1: Presets.P1_BALANCED_PRESET.gear,
					2: Presets.PRERAID_BALANCED_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_BALANCED_PRESET.gear,
					2: Presets.PRERAID_BALANCED_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class ProtectionWarriorSimUI extends IndividualSimUI<Spec.SpecProtectionWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
