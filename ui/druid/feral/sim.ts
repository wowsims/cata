import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/other_inputs.js';
import { PhysicalDPSGemOptimizer } from '../../core/components/suggest_gems_action.js';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui.js';
import { Player } from '../../core/player.js';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLPrepullAction, APLRotation , APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { Cooldowns, Debuffs, Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, RaidBuffs, Spec, Stat, TristateEffect } from '../../core/proto/common.js';
import { FeralDruid_Rotation as DruidRotation } from '../../core/proto/druid.js';
import * as AplUtils from '../../core/proto_utils/apl_utils.js';
import { Gear } from '../../core/proto_utils/gear.js';
import { Stats } from '../../core/proto_utils/stats.js';
import * as FeralInputs from './inputs.js';
import * as Presets from './presets.js';

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
		Stat.StatMeleeHit,
		Stat.StatExpertise,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatExpertise,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatMastery,
		Stat.StatMana,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.PRERAID_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap(
			{
				[Stat.StatStrength]: 0.38,
				[Stat.StatAgility]: 1.0,
				[Stat.StatAttackPower]: 0.37,
				[Stat.StatMeleeHit]: 0.43,
				[Stat.StatExpertise]: 0.43,
				[Stat.StatMeleeCrit]: 0.40,
				[Stat.StatMeleeHaste]: 0.41,
				[Stat.StatMastery]: 0.58,
			},
			{
				[PseudoStat.PseudoStatMainHandDps]: 1.55,
			},
		),
		other: Presets.OtherDefaults,
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotationType: APLRotationType.TypeSimple,
		simpleRotation: Presets.DefaultRotation,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			markOfTheWild: true,
			strengthOfEarthTotem: true,
			abominationsMight: true,
			windfuryTotem: true,
			bloodlust: true,
			communion: true,
			arcaneBrilliance: true,
			manaSpringTotem: true,
		}),
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
		}),
		debuffs: Debuffs.create({
			savageCombat: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: FeralInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [BuffDebuffInputs.ManaBuff, BuffDebuffInputs.MP5Buff],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [FeralInputs.AssumeBleedActive, OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget, OtherInputs.DarkIntentUptime],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [Presets.StandardTalents, Presets.HybridTalents],
		rotations: [Presets.SIMPLE_ROTATION_DEFAULT, Presets.APL_ROTATION_DEFAULT, Presets.APL_ROTATION_AOE],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PRERAID_PRESET, Presets.P1_PRESET, Presets.P2_PRESET, Presets.P3_PRESET, Presets.P4_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecFeralDruid>): APLRotation => {
		return Presets.APL_ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecFeralDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		const [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const blockZerk = APLAction.fromJsonString(`{"condition":{"const":{"val":"false"}},"castSpell":{"spellId":{"spellId":50334}}}`);
		const doRotation = APLAction.fromJsonString(`{"catOptimalRotationAction":{"rotationType":${simple.rotationType},"manualParams":${simple.manualParams},"maintainFaerieFire":${simple.maintainFaerieFire},"allowAoeBerserk":${simple.allowAoeBerserk},"minRoarOffset":${simple.minRoarOffset.toFixed(2)},"ripLeeway":${simple.ripLeeway.toFixed(0)},"useRake":${simple.useRake},"useBite":${simple.useBite},"biteDuringExecute":${simple.biteDuringExecute},"biteTime":${simple.biteTime.toFixed(2)}}}`);

		actions.push(...([blockZerk, doRotation].filter(a => a) as Array<APLAction>));

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
			spec: Spec.SpecFeralDruid,
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
			otherDefaults: Presets.OtherDefaults,
		},
	],
});

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, SPEC_CONFIG);

		const _gemOptimizer = new FeralGemOptimizer(this);
	}
}

class FeralGemOptimizer extends PhysicalDPSGemOptimizer {
	constructor(simUI: IndividualSimUI<Spec.SpecFeralDruid>) {
		super(simUI, true, true, true, true);
	}

	calcCritCap(gear: Gear): Stats {
		const baseCritCapPercentage = 77.8; // includes 3% Crit debuff
		let agiProcs = 0;

		if (gear.hasRelic(47668)) {
			agiProcs += 200;
		}

		if (gear.hasRelic(50456)) {
			agiProcs += 44 * 5;
		}

		if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
			agiProcs += 510;
		}

		if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
			agiProcs += 450;
		}

		if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
			agiProcs += 300;
		}

		return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - (agiProcs * 1.1 * 1.06 * 1.02) / 83.33) * 45.91);
	}
}
