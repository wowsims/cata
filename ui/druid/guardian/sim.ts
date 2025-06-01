import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action.js';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation, APLRotation_Type as APLRotationType, SimpleRotation } from '../../core/proto/apl.js';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common.js';
import { GuardianDruid_Rotation as DruidRotation } from '../../core/proto/druid.js';
import { StatCapType } from '../../core/proto/ui';
import * as AplUtils from '../../core/proto_utils/apl_utils.js';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats.js';
import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecGuardianDruid, {
	cssClass: 'guardian-druid-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Druid),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],
	warnings: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodgeRating,
		Stat.StatMasteryRating,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatHitRating,
		Stat.StatExpertiseRating,
		Stat.StatCritRating,
		Stat.StatHasteRating,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps, PseudoStat.PseudoStatPhysicalHitPercent, PseudoStat.PseudoStatSpellHitPercent],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[
			Stat.StatHealth,
			Stat.StatStamina,
			Stat.StatAgility,
			Stat.StatArmor,
			Stat.StatBonusArmor,
			Stat.StatMasteryRating,
			Stat.StatStrength,
			Stat.StatAttackPower,
			Stat.StatExpertiseRating,
		],
		[
			PseudoStat.PseudoStatDodgePercent,
			PseudoStat.PseudoStatPhysicalHitPercent,
			PseudoStat.PseudoStatPhysicalCritPercent,
			PseudoStat.PseudoStatMeleeHastePercent,
			PseudoStat.PseudoStatSpellHitPercent,
			PseudoStat.PseudoStatSpellCritPercent,
		],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.SURVIVAL_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats().withStat(Stat.StatExpertiseRating, 15 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION).withPseudoStat(PseudoStat.PseudoStatPhysicalHitPercent, 7.5).withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 15);
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultSimpleRotation,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			// ebonPlaguebringer: true,
			// criticalMass: true,
			// bloodFrenzy: true,
			// frostFever: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.GuardianDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellDamageDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.AbsorbFrac,
			OtherInputs.BurstWindow,
			OtherInputs.HpPercentForDefensives,
			DruidInputs.StartingRage,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.SURVIVAL_EP_PRESET, Presets.BALANCED_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents, Presets.InfectedWoundsBuild],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SIMPLE, Presets.ROTATION_DEFAULT, Presets.ROTATION_CLEAVE, Presets.ROTATION_NEF],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
		builds: [
			//Presets.PRESET_BUILD_BOSS_DUMMY,
		],
	},

	autoRotation: (_player: Player<Spec.SpecGuardianDruid>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	// simpleRotation: (player: Player<Spec.SpecGuardianDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
	// 	// const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

	// 	// const stampedeSpellId = player.getTalents().stampede == 1 ? 81016 : 81017;
	// 	// const preStampede = APLPrepullAction.fromJsonString(
	// 	// 	`{"action":{"activateAura":{"auraId":{"spellId":${stampedeSpellId.toFixed(0)}}}},"doAtValue":{"const":{"val":"-0.1s"}}}`,
	// 	// );

	// 	// const emergencySpellId = player.getTalents().pulverize ? 80313 : 33745;
	// 	// const emergencyPulverize = APLAction.fromJsonString(
	// 	// 	`{"condition":{"and":{"vals":[{"dotIsActive":{"targetUnit":{"type":"Target"},"spellId":{"spellId":33745}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"targetUnit":{"type":"Target"},"spellId":{"spellId":33745}}},"rhs":{"const":{"val":"${simple.pulverizeTime.toFixed(
	// 	// 		1,
	// 	// 	)}s"}}}}]}},"castSpell":{"spellId":{"spellId":${emergencySpellId.toFixed(0)}},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const faerieFireMaintain = APLAction.fromJsonString(
	// 	// 	`{"condition":{"or":{"vals":[{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":770}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":770}}},"rhs":{"const":{"val":"6s"}}}}]}},"castSpell":{"spellId":{"spellId":16857},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const demoRoar = APLAction.fromJsonString(
	// 	// 	`{"condition":{"auraShouldRefresh":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":99},"maxOverlap":{"const":{"val":"${simple.demoTime.toFixed(
	// 	// 		1,
	// 	// 	)}s"}}}},"castSpell":{"spellId":{"spellId":99},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const berserk = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":50334}}}`);
	// 	// const enrage = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":5229}}}`);
	// 	// const synapseSprings = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":82174}}}`);
	// 	// const lacerateForProcs = APLAction.fromJsonString(
	// 	// 	`{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"targetUnit":{"type":"Target"},"spellId":{"spellId":33745}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":50334}}}}},{"cmp":{"op":"OpGt","lhs":{"spellTimeToReady":{"spellId":{"spellId":50334}}},"rhs":{"const":{"val":"3s"}}}}]}},"castSpell":{"spellId":{"spellId":33745},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const mangle = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":33878},"target":{"type":"Target"}}}`);
	// 	// const thrash = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":77758},"target":{"type":"Target"}}}`);
	// 	// const faerieFireFiller = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":16857},"target":{"type":"Target"}}}`);
	// 	// const pulverize = APLAction.fromJsonString(
	// 	// 	`{"condition":{"and":{"vals":[{"dotIsActive":{"targetUnit":{"type":"Target"},"spellId":{"spellId":33745}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":80951}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":80951}}},"rhs":{"const":{"val":"${simple.pulverizeTime.toFixed(
	// 	// 		1,
	// 	// 	)}s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":80313},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const lacerateBuild = APLAction.fromJsonString(
	// 	// 	`{"condition":{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"sourceUnit":{"type":"Target"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},"castSpell":{"spellId":{"spellId":33745},"target":{"type":"Target"}}}`,
	// 	// );
	// 	// const maul = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":6807},"target":{"type":"Target"}}}`);

	// 	// prepullActions.push(...([simple.prepullStampede ? preStampede : null].filter(a => a) as Array<APLPrepullAction>));

	// 	// actions.push(
	// 	// 	...([
	// 	// 		emergencyPulverize,
	// 	// 		simple.maintainFaerieFire ? faerieFireMaintain : null,
	// 	// 		simple.maintainDemoralizingRoar ? demoRoar : null,
	// 	// 		berserk,
	// 	// 		enrage,
	// 	// 		synapseSprings,
	// 	// 		lacerateForProcs,
	// 	// 		mangle,
	// 	// 		thrash,
	// 	// 		faerieFireFiller,
	// 	// 		pulverize,
	// 	// 		lacerateBuild,
	// 	// 		maul,
	// 	// 	].filter(a => a) as Array<APLAction>),
	// 	// );

	// 	// return APLRotation.create({
	// 	// 	simple: SimpleRotation.create({
	// 	// 		cooldowns: Cooldowns.create({
	// 	// 			hpPercentForDefensives: cooldowns.hpPercentForDefensives,
	// 	// 		}),
	// 	// 	}),
	// 	// 	prepullActions: prepullActions,
	// 	// 	priorityList: actions.map(action =>
	// 	// 		APLListItem.create({
	// 	// 			action: action,
	// 	// 		}),
	// 	// 	),
	// 	// });
	// },

	raidSimPresets: [
		{
			spec: Spec.SpecGuardianDruid,
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

export class GuardianDruidSimUI extends IndividualSimUI<Spec.SpecGuardianDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecGuardianDruid>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this);
		});
	}
}
