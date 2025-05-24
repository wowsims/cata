import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation, APLRotation_Type } from '../../core/proto/apl.js';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common.js';
import { Stats, UnitStat } from '../../core/proto_utils/stats.js';
import * as PaladinInputs from '../inputs.js';
import * as RetributionInputs from './inputs.js';
import * as Presets from './presets.js';

const getStatCaps = () => {
	const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 7.5);
	const expCap = new Stats().withStat(Stat.StatExpertiseRating, 7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

	return hitCap.add(expCap);
};

const getEPDefaults = (_: Player<Spec.SpecRetributionPaladin>) => {
	return Presets.P1_EP_PRESET.epWeights;
};

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRetributionPaladin, {
	cssClass: 'retribution-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatExpertiseRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatStrength,
			Stat.StatAgility,
			Stat.StatIntellect,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
			Stat.StatSpellPower,
			Stat.StatMana,
			Stat.StatHealth,
			Stat.StatMasteryRating,
		],
		[
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatSpellHastePercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatSpellHitPercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: getStatCaps(),
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
			elementalOath: true,
			powerWordFortitude: true,
			serpentsSwiftness: true,
			timeWarp: true,
			trueshotAura: true,
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
		inputs: [RetributionInputs.StartingHolyPower(), OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		rotations: [Presets.APL_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.P1_GEAR_PRESET],
		itemSwaps: [],
		builds: [Presets.P1_BUILD_PRESET],
	},

	autoRotation: (_: Player<Spec.SpecRetributionPaladin>): APLRotation => {
		return Presets.APL_PRESET.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRetributionPaladin,
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

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults,
			});
		});
	}
}
