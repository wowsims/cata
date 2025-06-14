import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl.js';
import { Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat, UnitStats } from '../../core/proto/common.js';
import { Stats, UnitStat } from '../../core/proto_utils/stats.js';
import * as ShamanInputs from '../inputs.js';
import * as EnhancementInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecEnhancementShaman, {
	cssClass: 'enhancement-shaman-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Shaman),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	overwriteDisplayStats: (player: Player<Spec.SpecEnhancementShaman>) => {
		const playerStats = player.getCurrentStats();

		const statMod = (current: UnitStats, previous?: UnitStats) => {
			return new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(current).subtract(Stats.fromProto(previous)).getStat(Stat.StatAttackPower) * 0.65);
		};

		const base = statMod(playerStats.baseStats!);
		const gear = statMod(playerStats.gearStats!, playerStats.baseStats);
		const talents = statMod(playerStats.talentsStats!, playerStats.gearStats);
		const buffs = statMod(playerStats.buffsStats!, playerStats.talentsStats);
		const consumes = statMod(playerStats.consumesStats!, playerStats.buffsStats);
		const final = new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(playerStats.finalStats).getStat(Stat.StatAttackPower) * 0.65);

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
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatExpertiseRating,
		Stat.StatSpellPower,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
		PseudoStat.PseudoStatPhysicalHitPercent,
		PseudoStat.PseudoStatSpellHitPercent,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatStamina,
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatIntellect,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatSpellPower,
			Stat.StatMasteryRating,
		],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatSpellHastePercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_ORC_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P3_EP_PRESET.epWeights,
		// Default stat caps for the Reforge optimizer
		statCaps: (() => {
			const physHitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 7.5);
			const spellHitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 15);
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return physHitCap.add(spellHitCap.add(expCap));
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [ShamanInputs.ShamanShieldInput(), ShamanInputs.ShamanImbueMH(), EnhancementInputs.ShamanImbueOH],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [EnhancementInputs.SyncTypeInput, OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotBack,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
	],
	customSections: [ShamanInputs.TotemsSection],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P3_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_ORC_PRESET,
		],
	},

	autoRotation: (player: Player<Spec.SpecEnhancementShaman>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecEnhancementShaman,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Alliance]: Race.RaceDraenei,
				[Faction.Horde]: Race.RaceOrc,
				[Faction.Unknown]: Race.RaceUnknown,
			},
			defaultGear: {
				[Faction.Alliance]: {
					1: Presets.P1_ORC_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_ORC_PRESET.gear,
				},
				[Faction.Unknown]: {},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class EnhancementShamanSimUI extends IndividualSimUI<Spec.SpecEnhancementShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecEnhancementShaman>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
