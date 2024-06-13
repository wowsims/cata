import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common';
import { Stats } from '../../core/proto_utils/stats';
import * as DeathKnightInputs from '../inputs';
import * as BloodInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBloodDeathKnight, {
	cssClass: 'blood-death-knight-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.DeathKnight),
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
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatSpellHit,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMastery,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatDodge,
		Stat.StatParry,
		Stat.StatMastery,
	],
	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BLOOD_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatArmor]: 0.05,
				[Stat.StatBonusArmor]: 0.03,
				[Stat.StatStamina]: 1,
				[Stat.StatStrength]: 0.33,
				[Stat.StatAgility]: 0.6,
				[Stat.StatAttackPower]: 0.06,
				[Stat.StatExpertise]: 0.67,
				[Stat.StatMeleeHit]: 0.67,
				[Stat.StatMeleeCrit]: 0.28,
				[Stat.StatMeleeHaste]: 0.21,
				[Stat.StatDodge]: 0.7,
				[Stat.StatParry]: 0.58,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 3.1,
				[PseudoStat.PseudoStatOffHandDps]: 0.0,
			},
		),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.BloodTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			retributionAura: true,
			powerWordFortitude: true,
			markOfTheWild: true,
			icyTalons: true,
			hornOfWinter: true,
			abominationsMight: true,
			leaderOfThePack: true,
			bloodlust: true,
			arcaneTactics: true,
			devotionAura: true,
			resistanceAura: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			bloodFrenzy: true,
			sunderArmor: true,
			ebonPlaguebringer: true,
			criticalMass: true,
			vindication: true,
			frostFever: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: BloodInputs.BloodDeathKnightRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellDamageDebuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.SpellHasteBuff, BuffDebuffInputs.BleedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.InspirationUptime,
			OtherInputs.InFrontOfTarget,
			DeathKnightInputs.StartingRunicPower(),
			OtherInputs.DarkIntentUptime,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset rotations that the user can quickly select.
		rotations: [Presets.BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT, Presets.BLOOD_DEFENSIVE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.BloodTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_BLOOD_PRESET, Presets.P1_BLOOD_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecBloodDeathKnight>): APLRotation => {
		return Presets.BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBloodDeathKnight,
			talents: Presets.BloodTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_BLOOD_PRESET.gear,
					2: Presets.P1_BLOOD_PRESET.gear,
					3: Presets.P1_BLOOD_PRESET.gear,
					4: Presets.P1_BLOOD_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_BLOOD_PRESET.gear,
					2: Presets.P1_BLOOD_PRESET.gear,
					3: Presets.P1_BLOOD_PRESET.gear,
					4: Presets.P1_BLOOD_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class BloodDeathKnightSimUI extends IndividualSimUI<Spec.SpecBloodDeathKnight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBloodDeathKnight>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
