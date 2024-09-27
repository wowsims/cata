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
import { formatToNumber } from '../../core/utils';
import { sharedMageDisplayStatsModifiers } from '../shared';
import * as FireInputs from './inputs';
import * as Presets from './presets';

const hasteBreakpoints = Presets.FIRE_BREAKPOINTS.find(entry => entry.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent))!.presets!;

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFireMage, {
	cssClass: 'fire-mage-sim-ui',
	cssScheme: PlayerClasses.getCssClass(PlayerClasses.Mage),
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [],

	// All stats for which EP should be calculated.
	epStats: [Stat.StatIntellect, Stat.StatSpellPower, Stat.StatHitRating, Stat.StatCritRating, Stat.StatHasteRating, Stat.StatMasteryRating],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatSpellPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: UnitStat.createDisplayStatArray(
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpellPower, Stat.StatMasteryRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),
	modifyDisplayStats: (player: Player<Spec.SpecFireMage>) => {
		return sharedMageDisplayStatsModifiers(player);
	},

	defaults: {
		// Default equipped gear.
		gear: Presets.FIRE_P1_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.DEFAULT_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 17);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			// Set up Mastery breakpoints for integer % damage increments.
			// These should be removed once the bugfix to make Mastery
			// continuous goes live!
			const masterySoftCapConfig = StatCap.fromStat(Stat.StatMasteryRating, {
				breakpoints: [(23 / Mechanics.masteryPercentPerPoint.get(Spec.SpecFireMage)!) * Mechanics.MASTERY_RATING_PER_MASTERY_POINT],
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0],
			});

			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent, {
				breakpoints: [
					hasteBreakpoints.get('5-tick - LvB/Pyro')!,
					hasteBreakpoints.get('12-tick - Combust')!,
					hasteBreakpoints.get('13-tick - Combust')!,
					hasteBreakpoints.get('14-tick - Combust')!,
					hasteBreakpoints.get('6-tick - LvB/Pyro')!,
					hasteBreakpoints.get('15-tick - Combust')!,
					hasteBreakpoints.get('16-tick - Combust')!,
					hasteBreakpoints.get('7-tick - LvB/Pyro')!,
					hasteBreakpoints.get('17-tick - Combust')!,
					hasteBreakpoints.get('18-tick - Combust')!,
					hasteBreakpoints.get('19-tick - Combust')!,
					hasteBreakpoints.get('8-tick - LvB/Pyro')!,
					hasteBreakpoints.get('20-tick - Combust')!,
					hasteBreakpoints.get('21-tick - Combust')!,
				],
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0.61 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			return [masterySoftCapConfig, hasteSoftCapConfig];
		})(),
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
		}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [],
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
		epWeights: [Presets.DEFAULT_EP_PRESET],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.FIRE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.FireTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.FIRE_P1_PRESET, Presets.FIRE_P1_PREBIS, Presets.FIRE_P3_PRESET],
	},

	autoRotation: (_player): APLRotation => {
		/*const numTargets = player.sim.encounter.targets.length;
 		if (numTargets > 3) {
			return Presets.FIRE_ROTATION_PRESET_AOE.rotation.rotation!;
		} else {
			return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		} */
		return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

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
					2: Presets.FIRE_P1_PREBIS.gear,
				},
				[Faction.Horde]: {
					1: Presets.FIRE_P1_PRESET.gear,
					2: Presets.FIRE_P1_PREBIS.gear,
				},
			},
		},
	],
});

export class FireMageSimUI extends IndividualSimUI<Spec.SpecFireMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFireMage>) {
		super(parentElem, player, SPEC_CONFIG);

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				statSelectionPresets: Presets.FIRE_BREAKPOINTS,
				updateSoftCaps: softCaps => {
					const raidBuffs = player.getRaid()?.getBuffs();
					const hasBL = !!(raidBuffs?.bloodlust || raidBuffs?.timeWarp || raidBuffs?.heroism);
					const hasPI = !!player.getBuffs().powerInfusionCount;
					const hasBerserking = player.getRace() === Race.RaceTroll;

					const modifyHaste = (oldHastePercent: number, modifier: number) =>
						Number(formatToNumber(((oldHastePercent / 100 + 1) / modifier - 1) * 100, { maximumFractionDigits: 5 }));

					this.individualConfig.defaults.softCapBreakpoints!.forEach(softCap => {
						const softCapToModify = softCaps.find(sc => sc.unitStat.equals(softCap.unitStat));
						if (softCap.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent) && softCapToModify) {
							const adjustedHastedBreakpoints = new Set([...softCap.breakpoints]);
							const hasCloseMatchingValue = (value: number) =>
								[...adjustedHastedBreakpoints.values()].find(bp => bp.toFixed(2) === value.toFixed(2));

							// LvB/Pyro are not worth adjusting for
							const excludedHasteBreakpoints = [
								hasteBreakpoints.get('5-tick - LvB/Pyro')!,
								hasteBreakpoints.get('6-tick - LvB/Pyro')!,
								hasteBreakpoints.get('7-tick - LvB/Pyro')!,
								hasteBreakpoints.get('8-tick - LvB/Pyro')!,
							];
							softCap.breakpoints.forEach(breakpoint => {
								const isExcludedFromPiZerk = excludedHasteBreakpoints.includes(breakpoint);
								if (hasBL) {
									const blBreakpoint = modifyHaste(breakpoint, 1.3);

									if (blBreakpoint > 0) {
										if (!hasCloseMatchingValue(blBreakpoint)) adjustedHastedBreakpoints.add(blBreakpoint);
										if (hasBerserking) {
											const berserkingBreakpoint = modifyHaste(blBreakpoint, 1.2);
											if (berserkingBreakpoint > 0 && !hasCloseMatchingValue(blBreakpoint)) {
												adjustedHastedBreakpoints.add(berserkingBreakpoint);
											}
										}
									}
								}
								if (hasPI && !isExcludedFromPiZerk) {
									const piBreakpoint = modifyHaste(breakpoint, 1.2);
									if (piBreakpoint > 0) {
										if (!hasCloseMatchingValue(piBreakpoint)) adjustedHastedBreakpoints.add(piBreakpoint);
										if (hasBerserking) {
											const berserkingBreakpoint = modifyHaste(piBreakpoint, 1.2);
											if (berserkingBreakpoint > 0 && !hasCloseMatchingValue(piBreakpoint)) {
												adjustedHastedBreakpoints.add(berserkingBreakpoint);
											}
										}
									}
								}
							});
							softCapToModify.breakpoints = [...adjustedHastedBreakpoints].sort((a, b) => a - b);
						}
					});
					return softCaps;
				},
				additionalSoftCapTooltipInformation: {
					[Stat.StatHasteRating]: () => {
						const raidBuffs = player.getRaid()?.getBuffs();
						const hasBL = !!(raidBuffs?.bloodlust || raidBuffs?.timeWarp || raidBuffs?.heroism);
						const hasPI = !!player.getBuffs().powerInfusionCount;
						const hasBerserking = player.getRace() === Race.RaceTroll;

						return (
							<>
								{(hasBL || hasPI || hasBerserking) && (
									<>
										<p className="mb-0">Additional breakpoints have been created using the following cooldowns:</p>
										<ul className="mb-0">
											{hasBL && <li>Bloodlust</li>}
											{hasPI && <li>Power Infusion</li>}
											{hasBerserking && <li>Berserking</li>}
										</ul>
									</>
								)}
							</>
						);
					},
				},
			});
		});
	}
}
