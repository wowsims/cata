import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation, APLRotation_Type as APLRotationType } from '../../core/proto/apl';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import { FeralDruid_Rotation as DruidRotation, FeralDruid_Rotation_AplType as FeralRotationType } from '../../core/proto/druid';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import { TypedEvent } from '../../core/typed_event';
import * as FeralInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralDruid, {
	cssClass: 'feral-druid-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Druid),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatExpertiseRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
		Stat.StatMasteryRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatStrength, Stat.StatAgility, Stat.StatAttackPower, Stat.StatExpertiseRating, Stat.StatMasteryRating, Stat.StatMana],
		[PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatPhysicalCritPercent, PseudoStat.PseudoStatMeleeHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.BEARWEAVE_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 8);
			const expCap = new Stats().withStat(Stat.StatExpertiseRating, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return hitCap.add(expCap);
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultRotation,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: FeralInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			FeralInputs.AssumeBleedActive,
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			FeralInputs.CannotShredTarget,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotHands, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.BEARWEAVE_EP_PRESET, Presets.MONOCAT_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents, Presets.HybridTalents],
		rotations: [Presets.SIMPLE_ROTATION_DEFAULT, Presets.AOE_ROTATION_DEFAULT, Presets.APL_ROTATION_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		itemSwaps: [Presets.P4_ITEM_SWAP_PRESET],
		builds: [Presets.PRESET_BUILD_DEFAULT, Presets.PRESET_BUILD_TENDON],
	},

	autoRotation: (_player: Player<Spec.SpecFeralDruid>): APLRotation => {
		return Presets.APL_ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecFeralDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const synapseSprings = APLAction.fromJsonString(
			`{"condition":{"or":{"vals":[{"auraIsActive":{"auraId":{"spellId":5217}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"11s"}}}}]}},"castSpell":{"spellId":{"spellId":82174}}}`,
		);
		const potion = APLAction.fromJsonString(
			`{"condition":{"or":{"vals":[{"and":{"vals":[{"auraIsActive":{"auraId":{"spellId":5217}}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"math":{"op":"OpAdd","lhs":{"spellTimeToReady":{"spellId":{"spellId":50334}}},"rhs":{"const":{"val":"26s"}}}}}}]}},{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"26s"}}}},{"auraIsActive":{"auraId":{"spellId":50334}}}]}},"castSpell":{"spellId":{"itemId":58145}}}`,
		);
		const trollRacial = APLAction.fromJsonString(`{"condition":{"auraIsActive":{"auraId":{"spellId":50334}}},"castSpell":{"spellId":{"spellId":26297}}}`);
		const blockZerk = APLAction.fromJsonString(`{"condition":{"const":{"val":"false"}},"castSpell":{"spellId":{"spellId":50334}}}`);
		const blockEnrage = APLAction.fromJsonString(`{"condition":{"const":{"val":"false"}},"castSpell":{"spellId":{"spellId":5229}}}`);
		const doRotation = APLAction.fromJsonString(
			`{"catOptimalRotationAction":{"rotationType":${simple.rotationType},"manualParams":${simple.manualParams},"maintainFaerieFire":${
				simple.maintainFaerieFire
			},"allowAoeBerserk":${simple.allowAoeBerserk},"meleeWeave":${simple.meleeWeave},"bearWeave":${simple.bearWeave},"snekWeave":${
				simple.snekWeave
			},"minRoarOffset":${simple.minRoarOffset.toFixed(2)},"ripLeeway":${simple.ripLeeway.toFixed(0)},"useRake":${simple.useRake},"useBite":${
				simple.useBite
			},"biteDuringExecute":${simple.biteDuringExecute},"biteTime":${simple.biteTime.toFixed(2)},"berserkBiteTime":${simple.berserkBiteTime.toFixed(
				2,
			)},"cancelPrimalMadness":${simple.cancelPrimalMadness}}}`,
		);
		const autocasts = APLAction.fromJsonString(`{"autocastOtherCooldowns":{}}`);

		const singleTarget = simple.rotationType == FeralRotationType.SingleTarget;
		actions.push(
			...([
				singleTarget ? synapseSprings : null,
				singleTarget ? potion : null,
				singleTarget ? trollRacial : null,
				blockZerk,
				blockEnrage,
				doRotation,
				autocasts,
			].filter(a => a) as Array<APLAction>),
		);

		if (simple.prepullTranquility && player.shouldEnableTargetDummies()) {
			player.getRaid()?.setTargetDummies(TypedEvent.nextEventID(), 4);

			const trinketSwap = APLPrepullAction.fromJsonString(`{"action":{"itemSwap":{"swapSet":"Swap1"}},"doAtValue":{"const":{"val":"-125s"}}}`);
			const tranq = APLPrepullAction.fromJsonString(
				`{"action":{"channelSpell":{"spellId":{"spellId":740},"interruptIf":{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"-2s"}}}}}},"doAtValue":{"const":{"val":"-5.5s"}}}`,
			);
			const swapBack = APLPrepullAction.fromJsonString(`{"action":{"itemSwap":{"swapSet":"Main"}},"doAtValue":{"const":{"val":"-1.5s"}}}`);
			const shiftCat = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":768}}},"doAtValue":{"const":{"val":"-1.5s"}}}`);

			prepullActions.push(...([trinketSwap, tranq, swapBack, shiftCat].filter(a => a) as Array<APLPrepullAction>));
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

	hiddenMCDs: [50334, 5229],

	raidSimPresets: [
		{
			spec: Spec.SpecFeralDruid,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumables: Presets.DefaultConsumables,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
			},
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, SPEC_CONFIG);
		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
