import * as OtherInputs from '../../core/components/inputs/other_inputs';
import { ReforgeOptimizer } from '../../core/components/suggest_reforges_action';
import * as Mechanics from '../../core/constants/mechanics';
import { IndividualSimUI, registerSpecConfig } from '../../core/individual_sim_ui';
import { Player } from '../../core/player';
import { PlayerClasses } from '../../core/player_classes';
import { APLAction, APLListItem, APLRotation, APLRotation_Type, SimpleRotation } from '../../core/proto/apl';
import { Cooldowns, Faction, IndividualBuffs, ItemSlot, PartyBuffs, PseudoStat, Race, Spec, Stat } from '../../core/proto/common';
import { MageMajorGlyph } from '../../core/proto/mage';
import { StatCapType } from '../../core/proto/ui';
import { StatCap, Stats, UnitStat } from '../../core/proto_utils/stats';
import { formatToNumber } from '../../core/utils';
import { DefaultDebuffs, DefaultRaidBuffs } from '../presets';
import * as FireInputs from './inputs';
import * as Presets from './presets';

const combustBreakpoints = Presets.COMBUSTION_BREAKPOINT.presets;
const glyphedCombustBreakpoints = Presets.GLYPHED_COMBUSTION_BREAKPOINT.presets;

const relevantCombustionBreakpoints = [
	combustBreakpoints.get('12-tick - Combust')!,
	combustBreakpoints.get('13-tick - Combust')!,
	combustBreakpoints.get('14-tick - Combust')!,
	combustBreakpoints.get('15-tick - Combust')!,
	combustBreakpoints.get('16-tick - Combust')!,
	combustBreakpoints.get('17-tick - Combust')!,
	combustBreakpoints.get('18-tick - Combust')!,
	combustBreakpoints.get('19-tick - Combust')!,
	combustBreakpoints.get('20-tick - Combust')!,
	combustBreakpoints.get('21-tick - Combust')!,
];

const relevantGlyphedCombustionBreakpoints = [
	glyphedCombustBreakpoints.get('21-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('22-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('23-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('24-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('25-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('26-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('27-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('28-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('29-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('30-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('31-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('32-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('33-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('34-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('35-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('36-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('37-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('38-tick - Combust (Glyph)')!,
	glyphedCombustBreakpoints.get('39-tick - Combust (Glyph)')!,
];

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
		[Stat.StatHealth, Stat.StatMana, Stat.StatStamina, Stat.StatIntellect, Stat.StatSpellPower, Stat.StatExpertiseRating],
		[PseudoStat.PseudoStatSpellHitPercent, PseudoStat.PseudoStatSpellCritPercent, PseudoStat.PseudoStatSpellHastePercent],
	),

	defaults: {
		// Default equipped gear.
		gear: Presets.P1_BIS_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Presets.DEFAULT_EP_PRESET.epWeights,
		// Default stat caps for the Reforge Optimizer
		statCaps: (() => {
			return new Stats().withPseudoStat(PseudoStat.PseudoStatSpellHitPercent, 15);
		})(),
		// Default soft caps for the Reforge optimizer
		softCapBreakpoints: (() => {
			const hasteSoftCapConfig = StatCap.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent, {
				breakpoints: relevantCombustionBreakpoints,
				capType: StatCapType.TypeThreshold,
				postCapEPs: [0.61 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT],
			});

			return [hasteSoftCapConfig];
		})(),
		// Default consumes settings.
		consumables: Presets.DefaultFireConsumables,
		// Default rotation settings.
		rotationType: APLRotation_Type.TypeSimple,
		simpleRotation: Presets.P1TrollDefaultSimpleRotation,
		// Default talents.
		talents: Presets.FireTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultFireOptions,
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
	rotationInputs: FireInputs.MageRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [],
	excludeBuffDebuffInputs: [],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [OtherInputs.InputDelay, OtherInputs.DistanceFromTarget, OtherInputs.TankAssignment],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotHands, ItemSlot.ItemSlotTrinket1, ItemSlot.ItemSlotTrinket2],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: true,
	},

	presets: {
		epWeights: [Presets.DEFAULT_EP_PRESET],
		// Preset rotations that the user can quickly select.
		rotations: [Presets.P1_SIMPLE_ROTATION_DEFAULT, Presets.P1_SIMPLE_ROTATION_NO_TROLL, Presets.FIRE_ROTATION_PRESET_DEFAULT],
		// Preset talents that the user can quickly select.
		talents: [Presets.FireTalents],
		// Preset gear configurations that the user can quickly select.
		gear: [Presets.PREBIS_PRESET, Presets.P1_BIS_PRESET],
		builds: [Presets.P1_PRESET_BUILD, Presets.P1_PRESET_BUILD_NO_TROLL],
	},

	autoRotation: (): APLRotation => {
		// return Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
		return Presets.P1_SIMPLE_ROTATION_DEFAULT.rotation.rotation!;
	},

	// simpleRotation: (player, simple): APLRotation => {
	// 	const rotation = Presets.FIRE_ROTATION_PRESET_DEFAULT.rotation.rotation!;
	// 	const { combustThreshold, combustLastMomentLustPercentage, combustNoLustPercentage } = simple;

	// 	const maxCombustDuringLust = APLAction.fromJsonString(
	// 		`{"condition":{"and":{"vals":[{"or":{"vals":[{"and":{"vals":[{"auraIsKnown":{"auraId":{"spellId":26297}}},{"auraIsActive":{"auraId":{"spellId":26297}}},{"cmp":{"op":"OpLe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"17s"}}}}]}},{"and":{"vals":[{"not":{"val":{"auraIsKnown":{"auraId":{"spellId":26297}}}}},{"cmp":{"op":"OpGt","lhs":{"auraRemainingTime":{"auraId":{"spellId":2825,"tag":-1}}},"rhs":{"const":{"val":"2s"}}}}]}}]}},{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":44457}}},{"cmp":{"op":"OpGt","lhs":{"mageCurrentCombustionDotEstimate":{}},"rhs":{"const":{"val":"${combustThreshold}"}}}}]}},"castSpell":{"spellId":{"spellId":11129}}}`,
	// 	);
	// 	const lastMomentCombustDuringLust = APLAction.fromJsonString(
	// 		`{"condition":{"and":{"vals":[{"or":{"vals":[{"and":{"vals":[{"auraIsKnown":{"auraId":{"spellId":26297}}},{"auraIsActive":{"auraId":{"spellId":26297}}},{"cmp":{"op":"OpLt","lhs":{"auraRemainingTime":{"auraId":{"spellId":26297}}},"rhs":{"const":{"val":"2.5s"}}}},{"cmp":{"op":"OpLe","lhs":{"currentTime":{}},"rhs":{"const":{"val":"17s"}}}}]}},{"and":{"vals":[{"not":{"val":{"auraIsKnown":{"auraId":{"spellId":26297}}}}},{"auraIsActive":{"auraId":{"spellId":2825,"tag":-1}}},{"cmp":{"op":"OpLe","lhs":{"auraRemainingTime":{"auraId":{"spellId":2825,"tag":-1}}},"rhs":{"const":{"val":"2s"}}}}]}}]}},{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":44457}}},{"cmp":{"op":"OpGt","lhs":{"mageCurrentCombustionDotEstimate":{}},"rhs":{"const":{"val":"${combustLastMomentLustPercentage}"}}}}]}},"castSpell":{"spellId":{"spellId":11129}}}`,
	// 	);
	// 	const combustOutsideOfLustAndBerserking = APLAction.fromJsonString(
	// 		`{"condition":{"and":{"vals":[{"or":{"vals":[{"and":{"vals":[{"auraIsKnown":{"auraId":{"spellId":26297}}},{"cmp":{"op":"OpGt","lhs":{"currentTime":{}},"rhs":{"const":{"val":"17s"}}}}]}},{"and":{"vals":[{"not":{"val":{"auraIsKnown":{"auraId":{"spellId":26297}}}}},{"not":{"val":{"auraIsActive":{"auraId":{"spellId":2825,"tag":-1}}}}}]}}]}},{"auraIsActive":{"sourceUnit":{"type":"CurrentTarget"},"auraId":{"spellId":44457}}},{"cmp":{"op":"OpGt","lhs":{"mageCurrentCombustionDotEstimate":{}},"rhs":{"const":{"val":"${combustNoLustPercentage}"}}}}]}},"castSpell":{"spellId":{"spellId":11129}}}`,
	// 	);
	// 	const lastMomentCombustBeforeEncounter = APLAction.fromJsonString(
	// 		`{"condition":{"and":{"vals":[{"cmp":{"op":"OpLt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"15s"}}}},{"cmp":{"op":"OpGt","lhs":{"mageCurrentCombustionDotEstimate":{}},"rhs":{"const":{"val":"${combustLastMomentLustPercentage}"}}}}]}},"castSpell":{"spellId":{"spellId":11129}}}`,
	// 	);

	// 	const modifiedSimpleRotation = rotation;

	// 	modifiedSimpleRotation.priorityList[5] = APLListItem.create({
	// 		action: maxCombustDuringLust,
	// 	});
	// 	modifiedSimpleRotation.priorityList[6] = APLListItem.create({
	// 		action: lastMomentCombustDuringLust,
	// 	});
	// 	modifiedSimpleRotation.priorityList[7] = APLListItem.create({
	// 		action: combustOutsideOfLustAndBerserking,
	// 	});
	// 	modifiedSimpleRotation.priorityList[8] = APLListItem.create({
	// 		action: lastMomentCombustBeforeEncounter,
	// 	});

	// 	return APLRotation.create({
	// 		simple: SimpleRotation.create({
	// 			cooldowns: Cooldowns.create(),
	// 		}),
	// 		prepullActions: modifiedSimpleRotation.prepullActions,
	// 		priorityList: modifiedSimpleRotation.priorityList,
	// 	});
	// },
	// Hide all the MCDs since the simeple rotation handles them.
	hiddenMCDs: [
		// Berserking
		26297,
		// Bloodlust
		2825,
		// Evocation
		12051,
		// Flame Orb
		82731,
		// Mana Gem
		36799,
		// Mirror Image
		55342,
		// Synapse Springs
		82174,
		// Volcanic Potion
		58091,
	],

	raidSimPresets: [
		{
			spec: Spec.SpecFireMage,
			talents: Presets.FireTalents.data,
			specOptions: Presets.DefaultFireOptions,
			consumables: Presets.DefaultFireConsumables,
			otherDefaults: Presets.OtherDefaults,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceWorgen,
				[Faction.Horde]: Race.RaceTroll,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_BIS_PRESET.gear,
					2: Presets.PREBIS_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_BIS_PRESET.gear,
					2: Presets.PREBIS_PRESET.gear,
				},
			},
		},
	],
});

export class FireMageSimUI extends IndividualSimUI<Spec.SpecFireMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFireMage>) {
		super(parentElem, player, SPEC_CONFIG);

		const statSelectionPresets = [
			{
				unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
				presets: new Map([...Presets.COMBUSTION_BREAKPOINT.presets, ...Presets.GLYPHED_COMBUSTION_BREAKPOINT.presets].sort((a, b) => a[1] - b[1])),
			},
		];

		player.sim.waitForInit().then(() => {
			new ReforgeOptimizer(this, {
				statSelectionPresets: statSelectionPresets,
				enableBreakpointLimits: true,
				updateSoftCaps: softCaps => {
					const raidBuffs = player.getRaid()?.getBuffs();
					const hasBL = !!raidBuffs?.bloodlust;
					const hasBerserking = player.getRace() === Race.RaceTroll;
					const hasGlyphOfCombustion = !!player.getMajorGlyphs().find(glyph => glyph === MageMajorGlyph.GlyphOfCombustion);

					const modifyHaste = (oldHastePercent: number, modifier: number) =>
						Number(formatToNumber(((oldHastePercent / 100 + 1) / modifier - 1) * 100, { maximumFractionDigits: 5 }));

					this.individualConfig.defaults.softCapBreakpoints!.forEach(softCap => {
						const softCapToModify = softCaps.find(sc => sc.unitStat.equals(softCap.unitStat));
						if (softCap.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent) && softCapToModify) {
							if (hasGlyphOfCombustion) softCap.breakpoints = relevantGlyphedCombustionBreakpoints;

							const adjustedHastedBreakpoints = new Set([...softCap.breakpoints]);
							const hasCloseMatchingValue = (value: number) =>
								[...adjustedHastedBreakpoints.values()].find(bp => bp.toFixed(2) === value.toFixed(2));

							softCap.breakpoints.forEach(breakpoint => {
								if (hasBL) {
									const blBreakpoint = modifyHaste(breakpoint, 1.3);

									if (blBreakpoint > 0) {
										if (!hasCloseMatchingValue(blBreakpoint)) adjustedHastedBreakpoints.add(blBreakpoint);
										if (hasBerserking) {
											const berserkingBreakpoint = modifyHaste(blBreakpoint, 1.2);
											if (berserkingBreakpoint > 0 && !hasCloseMatchingValue(berserkingBreakpoint)) {
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
						const hasBL = !!raidBuffs?.bloodlust;
						const hasBerserking = player.getRace() === Race.RaceTroll;

						return (
							<>
								{(hasBL || hasBerserking) && (
									<>
										<p className="mb-0">Additional breakpoints have been created using the following cooldowns:</p>
										<ul className="mb-0">
											{hasBL && <li>Bloodlust</li>}
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
