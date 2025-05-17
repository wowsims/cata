import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
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
	RangedWeaponType,
	RotationType,
	Spec,
	Stat,
} from '../../core/proto/common';
import { HunterStingType, SurvivalHunter_Rotation } from '../../core/proto/hunter';
import { StatCapType } from '../../core/proto/ui';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import * as HunterInputs from '../inputs';
import { sharedHunterDisplayStatsModifiers } from '../shared';
import * as SVInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecSurvivalHunter, {
	cssClass: 'survival-hunter-sim-ui',
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
	],
	epPseudoStats: [PseudoStat.PseudoStatRangedDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatRangedAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStamina, Stat.StatAgility, Stat.StatRangedAttackPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatPhysicalCritPercent, PseudoStat.PseudoStatRangedHastePercent],
	),
	modifyDisplayStats: (player: Player<Spec.SpecSurvivalHunter>) => {
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
		gear: Presets.SV_P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P4_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
		})(),
		softCapBreakpoints: (() => {
			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatRangedHastePercent, {
				breakpoints: [20],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0.89 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			return [hasteSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.SurvivalTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.SVDefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// faerieFire: true,
			// curseOfElements: true,
			// savageCombat: true,
			// bloodFrenzy: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: SVInputs.SVRotationConfig,
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.StaminaBuff, BuffDebuffInputs.SpellDamageDebuff, BuffDebuffInputs.MajorArmorDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			HunterInputs.PetUptime(),
			HunterInputs.AQTierPrepull(),
			HunterInputs.NaxxTierPrepull(),
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET, Presets.P3_EP_PRESET, Presets.P4_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.SurvivalTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SV, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		builds: [Presets.P2_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		gear: [Presets.SV_PRERAID_PRESET, Presets.SV_P1_PRESET, Presets.SV_P3_PRESET, Presets.SV_P4_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecSurvivalHunter>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 4) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_SV.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecSurvivalHunter>, simple: SurvivalHunter_Rotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action =>
				APLListItem.create({
					action: action,
				}),
			),
		});
	},

	raidSimPresets: [
		{
			spec: Spec.SpecSurvivalHunter,
			talents: Presets.SurvivalTalents.data,
			specOptions: Presets.SVDefaultOptions,

			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.SV_PRERAID_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.SV_PRERAID_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class SurvivalHunterSimUI extends IndividualSimUI<Spec.SpecSurvivalHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecSurvivalHunter>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				getEPDefaults: (player: Player<Spec.SpecSurvivalHunter>) => {
					if (player.getGear().getItemSetCount('Lightning-Charged Battlegear') >= 4) {
						return Presets.P1_EP_PRESET.epWeights;
					}
					if (player.getGear().getItemSetCount("Flamewaker's Battlegear") >= 4) {
						return Presets.P3_EP_PRESET.epWeights;
					}
					return Presets.P3_EP_PRESET.epWeights;
				},
			});
		});
	}
}
