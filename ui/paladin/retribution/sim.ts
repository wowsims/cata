import * as BuffDebuffInputs from '../../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../../core/components/inputs/other_inputs.js';
import {ReforgeOptimizer} from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics.js';
import {IndividualSimUI, registerSpecConfig} from '../../core/individual_sim_ui.js';
import {Player} from '../../core/player.js';
import {PlayerClasses} from '../../core/player_classes';
import {APLRotation, APLRotation_Type} from '../../core/proto/apl.js';
import {
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat
} from '../../core/proto/common.js';
import {PaladinPrimeGlyph, PaladinSeal} from '../../core/proto/paladin';
import {Stats} from '../../core/proto_utils/stats.js';
import {TypedEvent} from '../../core/typed_event.js';
import * as PaladinInputs from '../inputs.js';
import * as Presets from './presets.js';
import * as RetributionInputs from './inputs.js';

const isGlyphOfSealOfTruthActive = (player: Player<Spec.SpecRetributionPaladin>): boolean => {
	const currentSeal = player.getSpecOptions().classOptions?.seal;
	return (
		player.getPrimeGlyps().includes(PaladinPrimeGlyph.GlyphOfSealOfTruth) &&
		(currentSeal === PaladinSeal.Truth || currentSeal === PaladinSeal.Righteousness)
	);
};

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRetributionPaladin, {
	cssClass: 'retribution-paladin-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Paladin),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatMastery
	],
	epPseudoStats: [PseudoStat.PseudoStatMainHandDps],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatMP5,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatSpellHaste,
		Stat.StatSpellPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatMana,
		Stat.StatHealth,
		Stat.StatMastery,
	],
	modifyDisplayStats: (player: Player<Spec.SpecRetributionPaladin>) => {
		let stats = new Stats();

		TypedEvent.freezeAllAndDo(() => {
			if (isGlyphOfSealOfTruthActive(player)) {
				stats = stats.addStat(Stat.StatExpertise, 2.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
			}
		});

		return {
			talents: stats,
		};
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BIS_RET_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			const hitCap = new Stats().withStat(Stat.StatMeleeHit, 8 * Mechanics.MELEE_HIT_RATING_PER_HIT_CHANCE);
			const expCap = new Stats().withStat(Stat.StatExpertise, 6.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);

			return hitCap.add(expCap);
		})(),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.RetTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
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
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: Debuffs.create({
			exposeArmor: true,
			bloodFrenzy: true,
			mangle: true,
			ebonPlaguebringer: true,
			criticalMass: true,
		}),
		rotationType: APLRotation_Type.TypeSimple,
		simpleRotation: Presets.ROTATION_PRESET_DEFAULT,
	},

	playerInputs: {
		inputs: [RetributionInputs.SnapshotGuardian()],
	},
	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [PaladinInputs.AuraSelection(), PaladinInputs.StartingSealSelection()],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.SpellDamageDebuff,
		BuffDebuffInputs.ManaBuff,
		BuffDebuffInputs.SpellHasteBuff,
		BuffDebuffInputs.PowerInfusion,
	],
	excludeBuffDebuffInputs: [BuffDebuffInputs.BleedDebuff],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.TankAssignment, OtherInputs.InFrontOfTarget],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		epWeights: [Presets.P1_EP_PRESET],
		rotations: [Presets.ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.RetTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.P1_BIS_RET_PRESET, Presets.PRERAID_RET_PRESET, Presets.P1_NONHC_RET_PRESET],
	},

	autoRotation: (_player: Player<Spec.SpecRetributionPaladin>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (_player: Player<Spec.SpecRetributionPaladin>): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRetributionPaladin,
			talents: Presets.RetTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceBloodElf,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.PRERAID_RET_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.PRERAID_RET_PRESET.gear,
				},
			},
		},
	],
});

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				updateGearStatsModifier: (baseStats: Stats) => {
					if (isGlyphOfSealOfTruthActive(player)) {
						return baseStats.addStat(Stat.StatExpertise, 2.5 * 4 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
					} else {
						return baseStats;
					}
				},
			});
		});
	}
}
