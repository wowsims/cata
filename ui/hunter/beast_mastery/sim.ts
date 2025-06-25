import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLRotation } from '../../core/proto/apl';
import {
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	RotationType,
	Spec,
	Stat,
} from '../../core/proto/common';
import { BeastMasteryHunter_Rotation, HunterStingType } from '../../core/proto/hunter';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import * as HunterInputs from '../inputs';
import { sharedHunterDisplayStatsModifiers } from '../shared';
import * as BMInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecBeastMasteryHunter, {
	cssClass: 'beast-mastery-hunter-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Hunter),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],
	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatRangedAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
		Stat.StatExpertiseRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatRangedDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatRangedAttackPower, Stat.StatMasteryRating, Stat.StatExpertiseRating],
		[PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatPhysicalCritPercent, PseudoStat.PseudoStatRangedHastePercent],
	),
	modifyDisplayStats: (player: Player<Spec.SpecBeastMasteryHunter>) => {
		return sharedHunterDisplayStatsModifiers(player);
	},
	itemSwapSlots: [
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotTrinket1,
		ItemSlot.ItemSlotTrinket2,
		ItemSlot.ItemSlotFinger1,
		ItemSlot.ItemSlotFinger2,
	],
	defaults: {
		// Default equipped gear.
		gear: Presets.P1_PRESET_GEAR.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats()
				.withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 7.5)
				.withStat(Stat.StatExpertiseRating, 7.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
		})(),

		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.DefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.BMDefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			blessingOfKings: true,
			trueshotAura: true,
			leaderOfThePack: true,
			blessingOfMight: true,
			commandingShout: true,
			unholyAura: true,
			bloodlust: true,
			skullBannerCount: 2,
			stormlashTotemCount: 4,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			weakenedArmor: true,
			physicalVulnerability: true,
			curseOfElements: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [HunterInputs.PetTypeInput()], //[HunterInputs.PetTypeInput()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: BMInputs.BMRotationConfig,
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.StaminaBuff, BuffDebuffInputs.SpellDamageDebuff, BuffDebuffInputs.MajorArmorDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [HunterInputs.PetUptime(), HunterInputs.GlaiveTossChance(), OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.TankAssignment],
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
		rotations: [Presets.ROTATION_PRESET_BM, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		builds: [Presets.PRERAID_PRESET, Presets.PRERAID_PRESET_CELESTIAL, Presets.P1_PRESET],
		gear: [Presets.PRERAID_PRESET_GEAR, Presets.PRERAID_CELESTIAL_PRESET_GEAR, Presets.P1_PRESET_GEAR],
	},

	autoRotation: (player: Player<Spec.SpecBeastMasteryHunter>): APLRotation => {
		return Presets.ROTATION_PRESET_BM.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecBeastMasteryHunter,
			talents: Presets.DefaultTalents.data,
			specOptions: Presets.BMDefaultOptions,

			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PRERAID_CELESTIAL_PRESET_GEAR.gear,
				},
				[Faction.Horde]: {
					1: Presets.PRERAID_CELESTIAL_PRESET_GEAR.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class BeastMasteryHunterSimUI extends IndividualSimUI<Spec.SpecBeastMasteryHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecBeastMasteryHunter>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecFuryWarrior>) => {
					return Presets.P1_EP_PRESET.epWeights;
				},
			});
		});
	}
}
