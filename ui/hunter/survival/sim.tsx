import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs';
import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLRotation } from '../../core/proto/apl';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, RotationType, Spec, Stat } from '../../core/proto/common';
import { HunterStingType, SurvivalHunter_Rotation } from '../../core/proto/hunter';
import { StatCapType } from '../../core/proto/ui';
import * as AplUtils from '../../core/proto_utils/apl_utils';
import { Stats } from '../../core/proto_utils/stats';
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
	epStats: [Stat.StatStamina, Stat.StatAgility, Stat.StatRangedAttackPower, Stat.StatMeleeHit, Stat.StatMeleeCrit, Stat.StatMeleeHaste, Stat.StatMastery],
	epPseudoStats: [PseudoStat.PseudoStatRangedDps],
	// Reference stat against which to calculate EP.
	epReferenceStat: Stat.StatRangedAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatAgility,
		Stat.StatRangedAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
	],
	modifyDisplayStats: (player: Player<Spec.SpecSurvivalHunter>) => {
		return sharedHunterDisplayStatsModifiers(player);
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.SV_P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withStat(Stat.StatMeleeHit, 8 * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE);
			return hitCap;
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			const hasteSoftCapConfig = {
				stat: Stat.StatMeleeHaste,
				breakpoints: [2650],
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0.87],
			};

			return [hasteSoftCapConfig];
		})(),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.SurvivalTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.SVDefaultOptions,
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
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({
			vampiricTouch: true,
		}),
		debuffs: Debuffs.create({
			sunderArmor: true,
			curseOfElements: true,
			savageCombat: true,
			bloodFrenzy: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [HunterInputs.PetTypeInput()],
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
			SVInputs.SniperTrainingUptime,
			OtherInputs.InputDelay,
			OtherInputs.DistanceFromTarget,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
			OtherInputs.DarkIntentUptime,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		// Preset talents that the user can quickly select.
		talents: [Presets.SurvivalTalents],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_SIMPLE_DEFAULT, Presets.ROTATION_PRESET_SV, Presets.ROTATION_PRESET_SV_ADVANCED, Presets.ROTATION_PRESET_AOE],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.SV_PRERAID_PRESET, Presets.SV_P1_PRESET],
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
		const blackArrow = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":63672}}}`);
		const explosiveShot4 = APLAction.fromJsonString(
			`{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":60053}}}}},"castSpell":{"spellId":{"spellId":60053}}}`,
		);
		const explosiveShot3 = APLAction.fromJsonString(
			`{"condition":{"dotIsActive":{"spellId":{"spellId":60053}}},"castSpell":{"spellId":{"spellId":60052}}}`,
		);
		//const arcaneShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49045}}}`);

		if (simple.type == RotationType.Aoe) {
			actions.push(
				...([
					simple.sting == HunterStingType.ScorpidSting ? scorpidSting : null,
					simple.sting == HunterStingType.SerpentSting ? serpentSting : null,
					simple.trapWeave ? trapWeave : null,
					volley,
				].filter(a => a) as Array<APLAction>),
			);
		} else {
			// SV
			actions.push(
				...([
					killShot,
					explosiveShot4,
					simple.allowExplosiveShotDownrank ? explosiveShot3 : null,
					simple.trapWeave ? trapWeave : null,
					simple.sting == HunterStingType.ScorpidSting ? scorpidSting : null,
					simple.sting == HunterStingType.SerpentSting ? serpentSting : null,
					blackArrow,
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
			spec: Spec.SpecSurvivalHunter,
			talents: Presets.SurvivalTalents.data,
			specOptions: Presets.SVDefaultOptions,

			consumes: Presets.DefaultConsumes,
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
				updateSoftCaps: softCaps => {
					const hasT114PC = player.getGear().getItemSetCount('Lightning-Charged Battlegear') >= 4;
					this.individualConfig.defaults.softCapBreakpoints!.forEach(softCap => {
						const softCapToModify = softCaps.findIndex(sc => sc.stat === softCap.stat);
						// Remove the threshold if 4-set T11 is not found
						if (!hasT114PC && softCap.stat === Stat.StatMeleeHaste && softCapToModify !== -1) {
							softCaps.splice(softCapToModify, 1);
						}
					});
					return softCaps;
				},
				additionalSoftCapTooltipInformation: {
					[Stat.StatMeleeHaste]: () => {
						const hasT114PC = player.getGear().getItemSetCount('Lightning-Charged Battlegear') >= 4;
						return <>{hasT114PC && <p className="mb-0">T11 4-set was found, added haste threshold.</p>}</>;
					},
				},
			});
		});
	}
}
