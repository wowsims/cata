import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/other_inputs';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common';
import { Stats } from '../../core/proto_utils/stats';
import * as WarriorInputs from '../inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFuryWarrior, {
	cssClass: 'fury-warrior-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Warrior),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatArmor,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatExpertise,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatArmor,
	],
	// modifyDisplayStats: (player: Player<Spec.SpecFuryWarrior>) => {
	// 	//let stats = new Stats();
	// 	if (!player.getInFrontOfTarget()) {
	// 		// When behind target, dodge is the only outcome affected by Expertise.
	// 		//stats = stats.addStat(Stat.StatExpertise, player.getTalents().weaponMastery * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
	// 	}
	// 	return {
	// 	//	talents: stats,
	// 	};
	// },

	defaults: {
		// Default equipped gear.
		gear: Presets.P3_FURY_PRESET_ALLIANCE.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 2.72,
				[Stat.StatAgility]: 1.82,
				[Stat.StatAttackPower]: 1,
				[Stat.StatExpertise]: 2.55,
				[Stat.StatMeleeHit]: 0.79,
				[Stat.StatMeleeCrit]: 2.12,
				[Stat.StatMeleeHaste]: 1.72,
				[Stat.StatArmorPenetration]: 2.17,
				[Stat.StatArmor]: 0.03,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 6.29,
				[PseudoStat.PseudoStatOffHandDps]: 3.58,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.FuryTalents.data,
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
			bloodFrenzy: true,
			mangle: true,
			sunderArmor: true,
			curseOfElements: true,
			ebonPlaguebringer: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [WarriorInputs.ShoutPicker(), WarriorInputs.Recklessness(), WarriorInputs.ShatteringThrow()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		// just for Bryntroll
		BuffDebuffInputs.SpellDamageDebuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			WarriorInputs.StartingRage(),
			WarriorInputs.StanceSnapshot(),
			WarriorInputs.DisableExpertiseGemming(),
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.FuryTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_FURY, Presets.ROTATION_FURY_SUNDER],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_FURY_PRESET,
			Presets.P1_FURY_PRESET,
			Presets.P2_FURY_PRESET,
			Presets.P3_FURY_PRESET_ALLIANCE,
			Presets.P3_FURY_PRESET_HORDE,
			Presets.P4_FURY_PRESET_ALLIANCE,
			Presets.P4_FURY_PRESET_HORDE,
		],
	},

	autoRotation: (_player: Player<Spec.SpecFuryWarrior>): APLRotation => {
		return Presets.ROTATION_FURY.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFuryWarrior,
			talents: Presets.FuryTalents.data,
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
					1: Presets.P1_FURY_PRESET.gear,
					2: Presets.P2_FURY_PRESET.gear,
					3: Presets.P3_FURY_PRESET_ALLIANCE.gear,
					4: Presets.P4_FURY_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_FURY_PRESET.gear,
					2: Presets.P2_FURY_PRESET.gear,
					3: Presets.P3_FURY_PRESET_HORDE.gear,
					4: Presets.P4_FURY_PRESET_HORDE.gear,
				},
			},
		},
	],
});

export class FuryWarriorSimUI extends IndividualSimUI<Spec.SpecFuryWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFuryWarrior>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
