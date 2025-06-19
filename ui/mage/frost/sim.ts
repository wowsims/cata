import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLRotation } from '../../core/proto/apl';
import { Faction, IndividualBuffs, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { DefaultDebuffs, DefaultRaidBuffs, MAGE_BREAKPOINTS } from '../presets';
import * as FrostInputs from './inputs';
import * as Presets from './presets';

const hasteBreakpoints = MAGE_BREAKPOINTS.presets;

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFrostMage, {
	cssClass: 'frost-mage-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Mage),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpirit, Stat.StatSpellPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BIS.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.P1_EP_PRESET.epWeights,
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 15);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent, {
				breakpoints: [
					hasteBreakpoints.get('13-tick - Nether Tempest')!,
					hasteBreakpoints.get('5-tick - Living Bomb')!,
					hasteBreakpoints.get('14-tick - Nether Tempest')!,
					hasteBreakpoints.get('15-tick - Nether Tempest')!,
					hasteBreakpoints.get('16-tick - Nether Tempest')!,
					hasteBreakpoints.get('17-tick - Nether Tempest')!,
					hasteBreakpoints.get('6-tick - Living Bomb')!,
					hasteBreakpoints.get('18-tick - Nether Tempest')!,
					hasteBreakpoints.get('19-tick - Nether Tempest')!,
					hasteBreakpoints.get('7-tick - Living Bomb')!,
					hasteBreakpoints.get('20-tick - Nether Tempest')!,
					hasteBreakpoints.get('8-tick - Living Bomb')!,
					hasteBreakpoints.get('21-tick - Nether Tempest')!,
					hasteBreakpoints.get('22-tick - Nether Tempest')!,
				],
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0.45 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			const critSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellCritPercent, {
				breakpoints: [25],
				capType: StatCapType.TypeSoftCap,
				postCapEPs: [0.42 * Mechanics.CRIT_RATING_PER_CRIT_PERCENT],
			});

			return [critSoftCapConfig, hasteSoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultConsumables,
		// Default talents.
		talents: Presets.FrostDefaultTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultFrostOptions,
		other: Presets.OtherDefaults,
		// Default raid/party buffs settings.
		raidBuffs: DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({}),
		individualBuffs: IndividualBuffs.create({}),
		debuffs: DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
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
		epWeights: [Presets.P1_EP_PRESET],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.ROTATION_PRESET_DEFAULT, Presets.ROTATION_PRESET_AOE],
		// Preset talents that the user can quickly select.
		talents: [Presets.FrostDefaultTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.P1_PREBIS_RICH, Presets.P1_PREBIS_POOR, Presets.P1_BIS],
	},

	autoRotation: (player: Player<Spec.SpecFrostMage>): APLRotation => {
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets > 3) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFrostMage,
			talents: Presets.FrostDefaultTalents.data,
			specOptions: Presets.DefaultFrostOptions,
			consumables: Presets.DefaultConsumables,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceGnome,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PREBIS_RICH.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PREBIS_RICH.gear,
				},
			},
		},
	],
});

export class FrostMageSimUI extends IndividualSimUI<Spec.SpecFrostMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFrostMage>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				statSelectionPresets: [MAGE_BREAKPOINTS],
				enableBreakpointLimits: true,
			});
		});
	}
}
