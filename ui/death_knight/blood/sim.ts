import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
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
		Stat.StatExpertiseRating,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatHealth,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodgeRating,
		Stat.StatParryRating,
		Stat.StatNatureResistance,
		Stat.StatShadowResistance,
		Stat.StatFrostResistance,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatArmor,
			Stat.StatStamina,
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatMasteryRating,
		],
		[
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatDodgePercent,
			PseudoStat.PseudoStatParryPercent,
		],
	),
	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BLOOD_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_BLOOD_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return hitCap.add(expCap);
		})(),
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

	modifyDisplayStats: (player: Player<Spec.SpecBloodDeathKnight>) => {
		// Blood Presence is a combat buff but we want to include its bonus in the stats
		const currentStats = player.getCurrentStats()
		if (currentStats.finalStats) {
			const bonusStamina = currentStats.finalStats.stats[Stat.StatStamina]*0.08
			const bonusHealth = bonusStamina*14
			const stats = new Stats().addStat(Stat.StatHealth, bonusHealth).addStat(Stat.StatStamina, bonusStamina)
			return {
				buffs: stats,
			};
		}
		return {}
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
		epWeights: [Presets.P1_BLOOD_EP_PRESET, Presets.P3_BLOOD_EP_PRESET],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT, Presets.BLOOD_DEFENSIVE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.BloodTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_BLOOD_PRESET, Presets.P1_BLOOD_PRESET, Presets.P3_BLOOD_BALANCED_PRESET, Presets.P3_BLOOD_DEFENSIVE_PRESET, Presets.P3_BLOOD_OFFENSIVE_PRESET],
		builds: [Presets.P1_PRESET, Presets.P3_PRESET],
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
				[Faction.Horde]: Race.RaceOrc,
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
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecFuryWarrior>) => {
					const hasP3Setup = player
						.getGear()
						.getEquippedItems()
						.some(item => (item?.item.phase || 0) >= 3);

					return hasP3Setup ? Presets.P3_BLOOD_EP_PRESET.epWeights : Presets.P1_BLOOD_EP_PRESET.epWeights;
				},
			});
		});
	}
}
