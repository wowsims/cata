import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import * as DeathKnightInputs from '../inputs';
import * as Presets from './presets';

const getEPDefaults = (player: Player<Spec.SpecUnholyDeathKnight>) => {
	let hasP3Setup = false;

	const items = player.getGear().getEquippedItems();

	for (const item of items) {
		const phase = item?.item.phase || 0;
		if (phase > 2) {
			hasP3Setup = true;
			break;
		}
	}

	return hasP3Setup ? Presets.P3_UNHOLY_EP_PRESET.epWeights : Presets.P2_UNHOLY_EP_PRESET.epWeights;
};

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
		Stat.StatExpertiseRating,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatArmor, Stat.StatStrength, Stat.StatAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
		],
	),
	defaults: {
		// Default equipped gear.
		gear: Presets.P4_BIS_GEAR_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P3_UNHOLY_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return hitCap.add(expCap);
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({}),
	},

	autoRotation: (_: Player<Spec.SpecUnholyDeathKnight>): APLRotation => Presets.DEFAULT_ROTATION_PRESET.rotation.rotation!,

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellDamageDebuff],
	excludeBuffDebuffInputs: [BuffDebuffInputs.DamageReduction, BuffDebuffInputs.CastSpeedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DeathKnightInputs.StartingRunicPower(),
			// DeathKnightInputs.PetUptime(),
			// UnholyInputs.SelfUnholyFrenzy,
			// UnholyInputs.UseAMSInput,
			// UnholyInputs.AvgAMSSuccessRateInput,
			// UnholyInputs.AvgAMSHitInput,
			// OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			OtherInputs.InputDelay,
		],
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotNeck,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotBack,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotWrist,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotWaist,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotFeet,
		ItemSlot.ItemSlotFinger1,
		ItemSlot.ItemSlotFinger2,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P2_UNHOLY_EP_PRESET, Presets.P3_UNHOLY_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.DefaultTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.DEFAULT_ROTATION_PRESET],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PREBIS_GEAR_PRESET, Presets.P2_BIS_GEAR_PRESET, Presets.P3_BIS_GEAR_PRESET, Presets.P4_BIS_GEAR_PRESET],
		builds: [Presets.PREBIS_PRESET, Presets.P2_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
	},

	raidSimPresets: [
		{
			spec: Spec.SpecUnholyDeathKnight,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P2_BIS_GEAR_PRESET.gear,
					2: Presets.P2_BIS_GEAR_PRESET.gear,
					3: Presets.P3_BIS_GEAR_PRESET.gear,
					4: Presets.P3_BIS_GEAR_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P2_BIS_GEAR_PRESET.gear,
					2: Presets.P2_BIS_GEAR_PRESET.gear,
					3: Presets.P3_BIS_GEAR_PRESET.gear,
					4: Presets.P3_BIS_GEAR_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class UnholyDeathKnightSimUI extends IndividualSimUI<Spec.SpecUnholyDeathKnight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecUnholyDeathKnight>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults,
			});
		});
	}
}
