import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action.js';
import * as Mechanics from '../../core/constants/mechanics.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation, APLRotation_Type as APLRotationType, SimpleRotation } from '../../core/proto/apl.js';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat } from '../../core/proto/common.js';
import { GuardianDruid_Rotation as DruidRotation } from '../../core/proto/druid.js';
import * as AplUtils from '../../core/proto_utils/apl_utils.js';
import { Stats } from '../../core/proto_utils/stats.js';
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
		Stat.StatDodge,
		Stat.StatMastery,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatExpertise,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatNatureResistance,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAgility,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatArmor,
		Stat.StatBonusArmor,
		Stat.StatDodge,
		Stat.StatMastery,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatExpertise,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatNatureResistance,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatHealth]: 0.04,
				[Stat.StatStamina]: 0.99,
				[Stat.StatAgility]: 1.0,
				[Stat.StatArmor]: 1.02,
				[Stat.StatBonusArmor]: 0.23,
				[Stat.StatDodge]: 0.97,
				[Stat.StatMastery]: 0.35,
				[Stat.StatStrength]: 0.11,
				[Stat.StatAttackPower]: 0.1,
				[Stat.StatMeleeHit]: 0.075,
				[Stat.StatSpellHit]: 0.195,
				[Stat.StatExpertise]: 0.15,
				[Stat.StatMeleeCrit]: 0.11,
				[Stat.StatMeleeHaste]: 0.0,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 0.0,
			},
		),
		// For breakpoints add additional entries to the array.
		// Used for Reforge Optimizer
		statCaps: (() => {
			const meleeHitCap = new Stats().withStat(Stat.StatMeleeHit, 5 * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE);
			const spellHitCap = new Stats().withStat(Stat.StatSpellHit, 4 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
			const expCap = new Stats().withStat(Stat.StatExpertise, 5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return [meleeHitCap.add(spellHitCap).add(expCap)];
		})(),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultSimpleRotation,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			powerWordFortitude: true,
			markOfTheWild: true,
			bloodlust: true,
			strengthOfEarthTotem: true,
			abominationsMight: true,
			windfuryTotem: true,
			communion: true,
			devotionAura: true,
		}),
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			ebonPlaguebringer: true,
			criticalMass: true,
			bloodFrenzy: true,
			frostFever: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.GuardianDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.SpellCritDebuff, BuffDebuffInputs.SpellDamageDebuff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.TankAssignment,
			OtherInputs.IncomingHps,
			OtherInputs.HealingCadence,
			OtherInputs.HealingCadenceVariation,
			OtherInputs.BurstWindow,
			OtherInputs.InspirationUptime,
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
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SIMPLE, Presets.ROTATION_DEFAULT],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecGuardianDruid>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecGuardianDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const stampedeSpellId = player.getTalents().stampede == 1 ? 81016 : 81017;
		const preStampede = APLPrepullAction.fromJsonString(
			`{"action":{"activateAura":{"auraId":{"spellId":${stampedeSpellId.toFixed(0)}}}},"doAtValue":{"const":{"val":"-0.1s"}}}`,
		);

		const emergencySpellId = player.getTalents().pulverize ? 80313 : 33745;
		const emergencyPulverize = APLAction.fromJsonString(
			`{"condition":{"and":{"vals":[{"dotIsActive":{"spellId":{"spellId":33745}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":33745}}},"rhs":{"const":{"val":"${simple.pulverizeTime.toFixed(
				1,
			)}s"}}}}]}},"castSpell":{"spellId":{"spellId":${emergencySpellId.toFixed(0)}}}}`,
		);
		const faerieFireMaintain = APLAction.fromJsonString(
			`{"condition":{"or":{"vals":[{"not":{"val":{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":770}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":770}}},"rhs":{"const":{"val":"6s"}}}}]}},"castSpell":{"spellId":{"spellId":16857}}}`,
		);
		const demoRoar = APLAction.fromJsonString(
			`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":99},"maxOverlap":{"const":{"val":"${simple.demoTime.toFixed(
				1,
			)}s"}}}},"castSpell":{"spellId":{"spellId":99}}}`,
		);
		const berserk = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":50334}}}`);
		const enrage = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":5229}}}`);
		const synapseSprings = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":82174}}}`);
		const lacerateForProcs = APLAction.fromJsonString(
			`{"condition":{"and":{"vals":[{"not":{"val":{"dotIsActive":{"spellId":{"spellId":33745}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":50334}}}}}]}},"castSpell":{"spellId":{"spellId":33745}}}`,
		);
		const mangle = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":33878}}}`);
		const thrash = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":77758}}}`);
		const faerieFireFiller = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":16857}}}`);
		const pulverize = APLAction.fromJsonString(
			`{"condition":{"and":{"vals":[{"dotIsActive":{"spellId":{"spellId":33745}}},{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},{"or":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":80951}}}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":80951}}},"rhs":{"const":{"val":"${simple.pulverizeTime.toFixed(
				1,
			)}s"}}}}]}}]}},"castSpell":{"spellId":{"spellId":80313}}}`,
		);
		const lacerateBuild = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":33745}}},"rhs":{"const":{"val":"3"}}}},"castSpell":{"spellId":{"spellId":33745}}}`,
		);
		const maul = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":6807}}}`);

		prepullActions.push(...([simple.prepullStampede ? preStampede : null].filter(a => a) as Array<APLPrepullAction>));

		actions.push(
			...([
				emergencyPulverize,
				simple.maintainFaerieFire ? faerieFireMaintain : null,
				simple.maintainDemoralizingRoar ? demoRoar : null,
				berserk,
				enrage,
				synapseSprings,
				lacerateForProcs,
				mangle,
				thrash,
				faerieFireFiller,
				pulverize,
				lacerateBuild,
				maul,
			].filter(a => a) as Array<APLAction>),
		);

		return APLRotation.create({
			simple: SimpleRotation.create({
				cooldowns: Cooldowns.create({
					hpPercentForDefensives: cooldowns.hpPercentForDefensives,
				}),
			}),
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
			spec: Spec.SpecGuardianDruid,
			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
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
