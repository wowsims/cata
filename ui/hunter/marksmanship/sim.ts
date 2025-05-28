import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
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
import { HunterStingType, MarksmanshipHunter_Rotation } from '../../core/proto/hunter';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import * as HunterInputs from '../inputs';
import { sharedHunterDisplayStatsModifiers } from '../shared';
import * as MMInputs from './inputs';
import * as Presets from './presets';
import { MM_T12_PRESET, P3_EP_PRESET } from './presets';
const SPEC_CONFIG = registerSpecConfig(Spec.SpecMarksmanshipHunter, {
	cssClass: 'marksmanship-hunter-sim-ui',
	cssScheme: 'hunter',
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
	modifyDisplayStats: (player: Player<Spec.SpecMarksmanshipHunter>) => {
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
		gear: Presets.MM_T12_PRESET.gear,
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
		talents: Presets.MarksmanTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.MMDefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// sunderArmor: true,
			// faerieFire: true,
			// curseOfElements: true,
			// savageCombat: true,
			// bloodFrenzy: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [HunterInputs.PetTypeInput()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: MMInputs.MMRotationConfig,
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
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET, P3_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.MarksmanTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_MM, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		gear: [MM_T12_PRESET, Presets.MM_PRERAID_PRESET, Presets.MM_P1_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecMarksmanshipHunter>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 4) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_MM.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecMarksmanshipHunter>, simple: MarksmanshipHunter_Rotation, cooldowns: Cooldowns): APLRotation => {
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
		const silencingShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":34490}}}`);
		const chimeraShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":53209}}}`);
		//const arcaneShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49045}}}`);

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
			// MM
			actions.push(
				...([
					combatPot,
					silencingShot,
					killShot,
					simple.sting == HunterStingType.ScorpidSting ? scorpidSting : null,
					simple.sting == HunterStingType.SerpentSting ? serpentSting : null,
					simple.trapWeave ? trapWeave : null,
					chimeraShot,
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
			spec: Spec.SpecMarksmanshipHunter,
			talents: Presets.MarksmanTalents.data,
			specOptions: Presets.MMDefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.MM_PRERAID_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.MM_PRERAID_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class MarksmanshipHunterSimUI extends IndividualSimUI<Spec.SpecMarksmanshipHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecMarksmanshipHunter>) {
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
