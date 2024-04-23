import * as OtherInputs from '../../core/components/other_inputs';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation } from '../../core/proto/apl';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common';
import { FrostMage_Rotation } from '../../core/proto/mage';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats } from '../../core/proto_utils/stats';
import * as MageInputs from '../inputs';
import * as FrostInputs from './inputs';
import * as Presets from './presets';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFrostMage, {
	cssClass: 'frost-mage-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Mage),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatSpellHit, Stat.StatSpellCrit, Stat.StatSpellHaste, Stat.StatMP5],
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
		Stat.StatMP5,
	],
	// modifyDisplayStats: (player: Player<Spec.SpecFrostMage>) => {
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
		gear: Presets.FROST_P3_PRESET_HORDE.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 0.48,
			[Stat.StatSpirit]: 0.42,
			[Stat.StatSpellPower]: 1,
			[Stat.StatSpellHit]: 0.38,
			[Stat.StatSpellCrit]: 0.58,
			[Stat.StatSpellHaste]: 0.94,
			[Stat.StatMP5]: 0.09,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultFrostConsumes,
		// Default talents.
		talents: Presets.FrostTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultFrostOptions,
		other: Presets.OtherDefaults,
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
		partyBuffs: PartyBuffs.create({
			manaTideTotems: 1,
		}),
		individualBuffs: IndividualBuffs.create({
			innervateCount: 0,
			vampiricTouch: true,
			focusMagic: true,
		}),
		debuffs: Debuffs.create({
			judgement: true,
			ebonPlaguebringer: true,
			shadowAndFlame: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [MageInputs.ArmorInput()],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: FrostInputs.MageRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		//Should add hymn of hope, revitalize, and
	],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			//FrostInputs.WaterElementalDisobeyChance,
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SIMPLE, Presets.FROST_ROTATION_PRESET_DEFAULT, Presets.FROST_ROTATION_PRESET_AOE],
		// Preset talents that the user can quickly select.
		talents: [Presets.FrostTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.FROST_P1_PRESET, Presets.FROST_P2_PRESET, Presets.FROST_P3_PRESET_ALLIANCE, Presets.FROST_P3_PRESET_HORDE],
	},

	autoRotation: (player: Player<Spec.SpecFrostMage>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 3) {
			return Presets.FROST_ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.FROST_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecFrostMage>, simple: FrostMage_Rotation, cooldowns: Cooldowns): APLRotation => {
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

		const deepFreeze = APLAction.fromJsonString(`{"condition":{"auraIsActive":{"auraId":{"spellId":44545}}},"castSpell":{"spellId":{"spellId":44572}}}`);
		const frostfireBoltWithBrainFreeze = APLAction.fromJsonString(
			`{"condition":{"auraIsActiveWithReactionTime":{"auraId":{"spellId":44549}}},"castSpell":{"spellId":{"spellId":47610}}}`,
		);
		const frostbolt = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":42842}}}`);
		const iceLance = APLAction.fromJsonString(
			`{"condition":{"cmp":{"op":"OpEq","lhs":{"auraNumStacks":{"auraId":{"spellId":44545}}},"rhs":{"const":{"val":"1"}}}},"castSpell":{"spellId":{"spellId":42914}}}`,
		);

		prepullActions.push(prepullMirrorImage);

		actions.push(
			...([
				berserking,
				hyperspeedAcceleration,
				combatPot,
				evocation,
				deepFreeze,
				frostfireBoltWithBrainFreeze,
				simple.useIceLance ? iceLance : null,
				frostbolt,
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
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFrostMage,
			talents: Presets.FrostTalents.data,
			specOptions: Presets.DefaultFrostOptions,
			consumes: Presets.DefaultFrostConsumes,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.FROST_P1_PRESET.gear,
					2: Presets.FROST_P2_PRESET.gear,
					3: Presets.FROST_P3_PRESET_ALLIANCE.gear,
				},
				[Faction.Horde]: {
					1: Presets.FROST_P1_PRESET.gear,
					2: Presets.FROST_P2_PRESET.gear,
					3: Presets.FROST_P3_PRESET_HORDE.gear,
				},
			},
		},
	],
});

export class FrostMageSimUI extends IndividualSimUI<Spec.SpecFrostMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFrostMage>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
