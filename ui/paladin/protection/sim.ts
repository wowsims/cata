import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation, APLRotation_Type } from '../../core/proto/apl.js';
import { Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, UnitStats } from '../../core/proto/common.js';
import { StatCapType } from '../../core/proto/ui.js';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats.js';
import * as PaladinInputs from '../inputs.js';
import * as Presets from './presets.js';

const getEPDefaults = (_: Player<Spec.SpecProtectionPaladin>) => {
	return Presets.P1_EP_PRESET.epWeights;
};

const SPEC_CONFIG = registerSpecConfig(Spec.SpecProtectionPaladin, {
	cssClass: 'protection-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	overwriteDisplayStats: (player: Player<Spec.SpecProtectionPaladin>) => {
		const playerStats = player.getCurrentStats();

		const statMod = (current: UnitStats, previous?: UnitStats) => {
			return new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(current).subtract(Stats.fromProto(previous)).getStat(Stat.StatAttackPower) * 0.5);
		};

		const base = statMod(playerStats.baseStats!);
		const gear = statMod(playerStats.gearStats!, playerStats.baseStats);
		const talents = statMod(playerStats.talentsStats!, playerStats.gearStats);
		const buffs = statMod(playerStats.buffsStats!, playerStats.talentsStats);
		const consumes = statMod(playerStats.consumesStats!, playerStats.buffsStats);
		const final = new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(playerStats.finalStats).getStat(Stat.StatAttackPower) * 0.5);

		return {
			base: base,
			gear: gear,
			talents: talents,
			buffs: buffs,
			consumes: consumes,
			final: final,
			stats: [Stat.StatSpellPower],
		};
	},

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatExpertiseRating,
		Stat.StatHasteRating,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodgeRating,
		Stat.StatParryRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatSpellHitPercent],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatArmor,
			Stat.StatBonusArmor,
			Stat.StatStamina,
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatMasteryRating,
		],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatBlockPercent,
			PseudoStat.PseudoStatDodgePercent,
			PseudoStat.PseudoStatParryPercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		// Values for now are pre-Cata initial WAG
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats();
		})(),
		softCapBreakpoints: (() => {
			const expertiseSoftCapConfig = StatCap.fromStat(Stat.StatExpertiseRating, {
				breakpoints: [7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION, 15 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [3.51, 0],
			});

			const hitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [7.5],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			});

			const spellHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHitPercent, {
				breakpoints: [15],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0],
			});

			return [expertiseSoftCapConfig, hitSoftCapConfig, spellHitSoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			blessingOfKings: true,
			blessingOfMight: true,
			bloodlust: true,
			elementalOath: true,
			powerWordFortitude: true,
			serpentsSwiftness: true,
			trueshotAura: true,
			skullBannerCount: 2,
			stormlashTotemCount: 4,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			curseOfElements: true,
			physicalVulnerability: true,
			weakenedArmor: true,
			weakenedBlows: true,
		}),
		rotationType: APLRotation_Type.TypeAuto,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PaladinInputs.StartingSealSelection()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.AbsorbFrac,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.APL_PRESET],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.P1_GEAR_PRESET],
		builds: [Presets.P1_BUILD_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecProtectionPaladin>): APLRotation => {
		return Presets.APL_PRESET.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecProtectionPaladin,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceBloodElf,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_GEAR_PRESET.gear,
				},
			},
		},
	],
});

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults,
				updateSoftCaps: softCaps => {
					softCaps[0].postCapEPs[0] = player.getEpWeights().getStat(Stat.StatExpertiseRating) * 0.9;
					return softCaps;
				},
			});
		});
	}
}
