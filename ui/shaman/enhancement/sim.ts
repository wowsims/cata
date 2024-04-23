import { FireElementalSection } from '../../core/components/fire_elemental_inputs.js';
import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl.js';
import { Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat, TristateEffect, UnitStats } from '../../core/proto/common.js';
import { ShamanImbue } from '../../core/proto/shaman.js';
import { Stats } from '../../core/proto_utils/stats.js';
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
			return new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(current).subtract(Stats.fromProto(previous)).getStat(Stat.StatAttackPower) * 0.55);
		}

		const base = statMod(playerStats.baseStats!);
		const gear = statMod(playerStats.gearStats!, playerStats.baseStats);
		const talents = statMod(playerStats.talentsStats!, playerStats.gearStats);
		const buffs = statMod(playerStats.buffsStats!, playerStats.talentsStats);
		const consumes = statMod(playerStats.consumesStats!, playerStats.buffsStats);
		const final = new Stats().withStat(Stat.StatSpellPower, Stats.fromProto(playerStats.finalStats).getStat(Stat.StatAttackPower) * 0.55);

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
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMastery,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMastery,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatIntellect]: 1.48,
				[Stat.StatAgility]: 1.59,
				[Stat.StatStrength]: 1.1,
				[Stat.StatSpellPower]: 1.13,
				[Stat.StatSpellHit]: 0, //default EP assumes cap
				[Stat.StatSpellCrit]: 0.91,
				[Stat.StatSpellHaste]: 0.37,
				[Stat.StatAttackPower]: 1.0,
				[Stat.StatMeleeHit]: 1.38,
				[Stat.StatMeleeCrit]: 0.81,
				[Stat.StatMeleeHaste]: 1.61, //haste is complicated
				[Stat.StatArmorPenetration]: 0.48,
				[Stat.StatExpertise]: 0, //default EP assumes cap
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 5.21,
				[PseudoStat.PseudoStatOffHandDps]: 2.21,
			},
		),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({
			vampiricTouch: true,
		}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [ShamanInputs.ShamanShieldInput(), ShamanInputs.ShamanImbueMH(), EnhancementInputs.ShamanImbueOH],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.ReplenishmentBuff, BuffDebuffInputs.MP5Buff, BuffDebuffInputs.SpellHasteBuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [EnhancementInputs.SyncTypeInput, OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	customSections: [ShamanInputs.TotemsSection, FireElementalSection],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_PRESET,
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
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDraenei,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
				},
			},
		},
	],
});

export class EnhancementShamanSimUI extends IndividualSimUI<Spec.SpecEnhancementShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecEnhancementShaman>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
