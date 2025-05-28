import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
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
		Stat.StatIntellect,
		Stat.StatAgility,
		Stat.StatRangedAttackPower,
		Stat.StatHitRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMP5,
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
		gear: Presets.BM_P3_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P3_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
			return hitCap;
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.BeastMasteryTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.BMDefaultOptions,
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
	playerIconInputs: [HunterInputs.PetTypeInput()], //[HunterInputs.PetTypeInput()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: BMInputs.BMRotationConfig,
	petConsumeInputs: [],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.StaminaBuff, BuffDebuffInputs.SpellDamageDebuff],
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
		epWeights: [Presets.P1_EP_PRESET, Presets.P3_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.BeastMasteryTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_BM, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.BM_P3_PRESET, Presets.BM_PRERAID_PRESET, Presets.BM_P1_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecBeastMasteryHunter>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 4) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_BM.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecBeastMasteryHunter>, simple: BeastMasteryHunter_Rotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const combatPot = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}`,
		);

		const serpentSting = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"6s"}}}},"multidot":{"spellId":{"spellId":49001},"maxDots":${
				simple.multiDotSerpentSting ? 3 : 1
			},"maxOverlap":{"const":{"val":"0ms"}}}}`,
		);
		const scorpidSting = APLAction.fromJsonString(
			`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":3043},"maxOverlap":{"const":{"val":"0ms"}}}},"castSpell":{"spellId":{"spellId":3043}}}`,
		);
		const trapWeave = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49067}}}`,
		);
		const volley = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":58434}}}`);
		const killShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":61006}}}`);
		const aimedShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49050}}}`);
		const multiShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49048}}}`);
		const steadyShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49052}}}`);

		if (simple.type == RotationType.Aoe) {
			actions.push(
				...([
					combatPot,
					simple.sting == HunterStingType.ScorpidSting ? scorpidSting : null,
					simple.sting == HunterStingType.SerpentSting ? serpentSting : null,
					simple.trapWeave ? trapWeave : null,
					volley,
				].filter(a => a) as Array<APLAction>),
			);
		} else {
			actions.push(
				...([
					combatPot,
					killShot,
					simple.trapWeave ? trapWeave : null,
					simple.sting == HunterStingType.ScorpidSting ? scorpidSting : null,
					simple.sting == HunterStingType.SerpentSting ? serpentSting : null,
					aimedShot,
					multiShot,
					steadyShot,
				].filter(a => a) as Array<APLAction>),
			);
		}

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
			spec: Spec.SpecBeastMasteryHunter,
			talents: Presets.BeastMasteryTalents.data,
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
					1: Presets.BM_PRERAID_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.BM_PRERAID_PRESET.gear,
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
					if (player.getGear().getItemSetCount('Lightning-Charged Battlegear') >= 4) {
						return Presets.P1_EP_PRESET.epWeights;
					}
					if (player.getGear().getItemSetCount("Flamewaker's Battlegear") >= 4) {
						return Presets.P3_EP_PRESET.epWeights;
					}
					return Presets.P1_EP_PRESET.epWeights;
				},
			});
		});
	}
}
