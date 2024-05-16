import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import { TankGemOptimizer } from '../../core/components/suggest_gems_action.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation , APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import {
	Class,
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../../core/proto/common.js';
import {
	GuardianDruid_Rotation as DruidRotation,
} from '../../core/proto/druid.js';
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
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
	],
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
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatHealth]: 0.20,
			[Stat.StatStamina]: 3.35,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 1.45,
			[Stat.StatBonusArmor]: 0.21,
			[Stat.StatDodge]: 0.33,
			[Stat.StatMastery]: 0.34,
			[Stat.StatStrength]: 0.13,
			[Stat.StatAttackPower]: 0.05,
			[Stat.StatMeleeHit]: 0.11,
			[Stat.StatExpertise]: 0.22,
			[Stat.StatMeleeCrit]: 0.08,
			[Stat.StatMeleeHaste]: 0.06,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
		}),
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
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
		}),
		debuffs: Debuffs.create({
			ebonPlaguebringer: true,
			criticalMass: true,
			bloodFrenzy: true,
			frostFever: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.GuardianDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellCritDebuff,
		BuffDebuffInputs.SpellDamageDebuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
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
		talents: [
			Presets.StandardTalents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_SIMPLE,
			Presets.ROTATION_DEFAULT,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
		],
	},

	autoRotation: (_player: Player<Spec.SpecGuardianDruid>): APLRotation => {
		return Presets.ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecGuardianDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const emergencyLacerate = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}`);
		const demoRoar = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":48560},"maxOverlap":{"const":{"val":"1.5s"}}}},"castSpell":{"spellId":{"spellId":48560}}}`);
		const mangle = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":48564}}}`);
		const delayFaerieFireForMangle = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48564}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":48564}}},"rhs":{"const":{"val":"1.0s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48564}}}}}`);
		const faerieFire = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":16857}}}`);
		const delayFillersForMangle = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":48564}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":48564}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":48564}}}}}`);
		const lacerate = APLAction.fromJsonString(`{"condition":{"or":{"vals":[{"cmp":{"op":"OpLt","lhs":{"auraNumStacks":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":48568}}},"rhs":{"const":{"val":"5"}}}},{"cmp":{"op":"OpLe","lhs":{"dotRemainingTime":{"spellId":{"spellId":48568}}},"rhs":{"const":{"val":"${simple.pulverizeTime.toFixed(1)}s"}}}}]}},"castSpell":{"spellId":{"spellId":48568}}}`);
		const swipe = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"${(simple.maulRageThreshold + 15).toFixed(0)}"}}}},"castSpell":{"spellId":{"spellId":48562}}}`);
		const queueMaul = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGe","lhs":{"currentRage":{}},"rhs":{"const":{"val":"${simple.maulRageThreshold.toFixed(0)}"}}}},"castSpell":{"spellId":{"spellId":48480,"tag":1}}}`);
		const waitForFaerieFire = APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"gcdIsReady":{}},{"not":{"val":{"spellIsReady":{"spellId":{"spellId":16857}}}}},{"cmp":{"op":"OpLt","lhs":{"spellTimeToReady":{"spellId":{"spellId":16857}}},"rhs":{"const":{"val":"1.5s"}}}}]}},"wait":{"duration":{"spellTimeToReady":{"spellId":{"spellId":16857}}}}}`);

		actions.push(...[
			emergencyLacerate,
			simple.maintainDemoralizingRoar ? demoRoar : null,
			mangle,
			delayFaerieFireForMangle,
			faerieFire,
			delayFillersForMangle,
			lacerate,
			swipe,
			queueMaul,
			waitForFaerieFire,
		].filter(a => a) as Array<APLAction>)

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action => APLListItem.create({
				action: action,
			}))
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
		//const _gemOptimizer = new TankGemOptimizer(this);
	}
}
