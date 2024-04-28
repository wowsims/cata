import * as OtherInputs from '../../core/components/other_inputs';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation } from '../../core/proto/apl';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common';
import { FireMage_Rotation, FireMage_Rotation_PrimaryFireSpell } from '../../core/proto/mage';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats } from '../../core/proto_utils/stats';
import * as MageInputs from '../inputs';
import * as FireInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFireMage, {
	cssClass: 'fire-mage-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Mage),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect,
			  Stat.StatSpirit,
			  Stat.StatSpellPower,
			  Stat.StatSpellHit,
			  Stat.StatSpellCrit,
			  Stat.StatSpellHaste,
			  Stat.StatMastery,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatMana,
		Stat.StatStamina,
		Stat.StatIntellect,
		Stat.StatSpirit,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMastery,
	],
	// modifyDisplayStats: (player: Player<Spec.SpecFireMage>) => {
	// 	let stats = new Stats();

	// 	if (player.getTalentTree() === 0) {
	// 		stats = stats.addStat(Stat.StatSpellHit, player.getTalents().arcaneFocus * 1 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
	// 	}

	// 	return {
	// 		talents: stats,
	// 	};
	// },

	defaults: {
		// Default equipped gear.
		gear: Presets.FIRE_P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.48,
			[Stat.StatSpirit]: 0.42,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellHit]: 0.38,
			[Stat.StatSpellCrit]: 0.58,
			[Stat.StatSpellHaste]: 0.94,
			[Stat.StatMastery]: 0.8
		}),
		// Default consumes settings.
		consumes: Presets.DefaultFireConsumes,
		// Default talents.
		talents: Presets.FireTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultFireOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,

		partyBuffs: PartyBuffs.create({
			manaTideTotems: 1,
		}),
		individualBuffs: IndividualBuffs.create({
			innervateCount: 0,
			vampiricTouch: true,
			focusMagic: true,
		}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [MageInputs.ArmorInput()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: FireInputs.MageRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		//Should add hymn of hope, revitalize, and
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.TankAssignment, OtherInputs.DarkIntentUptime],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset rotations that the user can quickly select.
		rotations: [Presets.FIRE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.FireTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.FIRE_P1_PRESET],
	},

	autoRotation: (player: Player<Spec.SpecFireMage>): APLRotation => {
		/*const numTargets = player.sim.encounter.targets.length;
 		if (numTargets > 3) {
			return Presets.FIRE_ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		} */
		return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	/* simpleRotation: (player: Player<Spec.SpecFireMage>, simple: FireMage_Rotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const prepullMirrorImage = APLPrepullAction.fromJsonString(
			`{"action":{"castSpell":{"spellId":{"spellId":55342}}},"doAtValue":{"const":{"val":"-2s"}}}`,
		);
		const berserking = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":26297}}}`,
		);
		const hyperspeedAcceleration = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"spellId":54758}}}`,
		);
		const combatPot = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"auraIsActive":{"auraId":{"spellId":12472}}}}},"castSpell":{"spellId":{"otherId":"OtherActionPotion"}}}`,
		);
		const evocation = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpLe","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"25%"}}}},"castSpell":{"spellId":{"spellId":12051}}}`,
		);

		const maintainImpScorch = APLAction.fromJsonString(
			`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":12873},"maxOverlap":{"const":{"val":"4s"}}}},"castSpell":{"spellId":{"spellId":42859}}}`,
		);
		const pyroWithHotStreak = APLAction.fromJsonString(
			`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44448}}},"castSpell":{"spellId":{"spellId":42891}}}`,
		);
		const livingBomb = APLAction.fromJsonString(
			`{"condition":{"and":{"vals":[{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"12s"}}}}]}},"multidot":{"spellId":{"spellId":55360},"maxDots":10,"maxOverlap":{"const":{"val":"0ms"}}}}`,
		);
		const cheekyFireBlastFinisher = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"spellCastTime":{"spellId":{"spellId":42859}}}}},"castSpell":{"spellId":{"spellId":42873}}}`,
		);
		const scorchFinisher = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpLe","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"4s"}}}},"castSpell":{"spellId":{"spellId":42859}}}`,
		);
		const fireball = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42833}}}`);
		const frostfireBolt = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":47610}}}`);
		const scorch = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42859}}}`);

		prepullActions.push(prepullMirrorImage);
		if (player.getTalents().improvedScorch > 0 && simple.maintainImprovedScorch) {
			actions.push(maintainImpScorch);
		}

		actions.push(
			...([
				berserking,
				hyperspeedAcceleration,
				combatPot,
				pyroWithHotStreak,
				livingBomb,
				cheekyFireBlastFinisher,
				scorchFinisher,
				evocation,
				simple.primaryFireSpell == FireMage_Rotation_PrimaryFireSpell.Fireball
					? fireball
					: simple.primaryFireSpell == FireMage_Rotation_PrimaryFireSpell.FrostfireBolt
					? frostfireBolt
					: scorch,
			].filter(a => a) as Array<APLAction>),
		);

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action =>
				APLListItem.create({
					action: action,
				}),
			),
		});
	}, */

	raidSimPresets: [
		{
			spec: Spec.SpecFireMage,
			talents: Presets.FireTalents.data,
			specOptions: Presets.DefaultFireOptions,
			consumes: Presets.DefaultFireConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.FIRE_P1_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.FIRE_P1_PRESET.gear,
				},
			},
		},
	],
});

export class FireMageSimUI extends IndividualSimUI<Spec.SpecFireMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFireMage>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
