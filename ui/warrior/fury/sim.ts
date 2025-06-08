import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Debuffs, Faction, HandType, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as WarriorInputs from '../inputs';
import * as FuryInputs from './inputs';
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
		Stat.StatExpertiseRating,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatOffHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatStrength, Stat.StatAgility, Stat.StatAttackPower, Stat.StatExpertiseRating, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatPhysicalCritPercent, PseudoStat.PseudoStatMeleeHastePercent],
	),
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
		gear: Presets.P3_BIS_FURY_TG_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P3_FURY_TG_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return expCap;
		})(),
		softCapBreakpoints: (() => {
			const meleeHitSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, {
				breakpoints: [8, 27],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [1.23 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT, 0],
			});

			return [meleeHitSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.

		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.FuryTGTalents.data,
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
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			bloodFrenzy: true,
			mangle: true,
			faerieFire: true,
			curseOfElements: true,
			ebonPlaguebringer: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		// just for Bryntroll
		BuffDebuffInputs.SpellDamageDebuff,
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			FuryInputs.SyncTypeInput,
			WarriorInputs.StartingRage(),
			WarriorInputs.StanceSnapshot(),
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			FuryInputs.AssumePrepullMasteryElixir,
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
		ItemSlot.ItemSlotRanged,
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.P1_FURY_SMF_EP_PRESET, Presets.P1_FURY_TG_EP_PRESET, Presets.P3_FURY_SMF_EP_PRESET, Presets.P3_FURY_TG_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.FurySMFTalents, Presets.FuryTGTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.FURY_SMF_ROTATION, Presets.FURY_TG_ROTATION],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.P3_PRERAID_FURY_SMF_PRESET,
			Presets.P3_PRERAID_FURY_TG_PRESET,
			Presets.P1_BIS_FURY_SMF_PRESET,
			Presets.P1_BIS_FURY_TG_PRESET,
			Presets.P3_BIS_FURY_SMF_PRESET,
			Presets.P3_BIS_FURY_TG_PRESET,
			Presets.P4_BIS_FURY_TG_PRESET,
			Presets.P4_BIS_FURY_SMF_PRESET,
		],
		itemSwaps: [Presets.P3_ITEM_SWAP_SMF, Presets.P3_ITEM_SWAP_TG, Presets.P4_ITEM_SWAP_TG, Presets.P4_ITEM_SWAP_SMF],
		builds: [
			Presets.P1_PRESET_BUILD_SMF,
			Presets.P1_PRESET_BUILD_TG,
			Presets.P3_PRESET_BUILD_SMF,
			Presets.P3_PRESET_BUILD_TG,
			Presets.P4_PRESET_BUILD_TG,
			Presets.P4_PRESET_BUILD_SMF,
		],
	},

	autoRotation: (_player: Player<Spec.SpecFuryWarrior>): APLRotation => {
		return Presets.FURY_TG_ROTATION.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFuryWarrior,
			talents: Presets.FurySMFTalents.data,
			specOptions: Presets.DefaultOptions,

			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P3_BIS_FURY_SMF_PRESET.gear,
					2: Presets.P3_BIS_FURY_TG_PRESET.gear,
					3: Presets.P3_PRERAID_FURY_SMF_PRESET.gear,
					4: Presets.P3_PRERAID_FURY_TG_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P3_BIS_FURY_SMF_PRESET.gear,
					2: Presets.P3_BIS_FURY_TG_PRESET.gear,
					3: Presets.P3_PRERAID_FURY_SMF_PRESET.gear,
					4: Presets.P3_PRERAID_FURY_TG_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class FuryWarriorSimUI extends IndividualSimUI<Spec.SpecFuryWarrior> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFuryWarrior>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecFuryWarrior>) => {
					const hasP3Setup = player
						.getGear()
						.getEquippedItems()
						.some(item => (item?.item.phase || 0) >= 3);

					if (player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeOneHand || !player.getTalents().titansGrip) {
						return hasP3Setup ? Presets.P3_FURY_SMF_EP_PRESET.epWeights : Presets.P1_FURY_SMF_EP_PRESET.epWeights;
					}
					return hasP3Setup ? Presets.P3_FURY_TG_EP_PRESET.epWeights : Presets.P1_FURY_TG_EP_PRESET.epWeights;
				},
			});
		});
	}
}
