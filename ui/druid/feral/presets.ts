import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Spec, Stat } from '../../core/proto/common';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	FeralDruid_Options as FeralDruidOptions,
	FeralDruid_Rotation as FeralDruidRotation,
	FeralDruid_Rotation_AplType,
	FeralDruid_Rotation_BiteModeType,
} from '../../core/proto/druid';
import { SavedTalents } from '../../core/proto/ui';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);
import P4ItemSwapGear from './gear_sets/p4_item_swap.gear.json';
export const P4_ITEM_SWAP_PRESET = PresetUtils.makePresetItemSwapGear('P4', P4ItemSwapGear);

import DefaultApl from './apls/default.apl.json';
export const APL_ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL List View', DefaultApl);

import { Stats } from '../../core/proto_utils/stats';
import AoeApl from './apls/aoe.apl.json';
export const APL_ROTATION_AOE = PresetUtils.makePresetAPLRotation('APL AoE', AoeApl);
import TendonApl from './apls/tendon.apl.json';
export const APL_ROTATION_TENDON = PresetUtils.makePresetAPLRotation('Tendon APL', TendonApl);

// Preset options for EP weights
export const BEARWEAVE_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Bear-Weave',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 0.38,
			[Stat.StatAgility]: 1.0,
			[Stat.StatAttackPower]: 0.37,
			[Stat.StatHitRating]: 0.36,
			[Stat.StatExpertiseRating]: 0.34,
			[Stat.StatCritRating]: 0.32,
			[Stat.StatHasteRating]: 0.3,
			[Stat.StatMasteryRating]: 0.33,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.54,
		},
	),
);

export const MONOCAT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Mono-Cat',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 0.39,
			[Stat.StatAgility]: 1.0,
			[Stat.StatAttackPower]: 0.37,
			[Stat.StatHitRating]: 0.31,
			[Stat.StatExpertiseRating]: 0.31,
			[Stat.StatCritRating]: 0.31,
			[Stat.StatHasteRating]: 0.3,
			[Stat.StatMasteryRating]: 0.33,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.56,
		},
	),
);

export const DefaultRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.SingleTarget,
	bearWeave: true,
	minCombosForRip: 5,
	minCombosForBite: 5,
	useRake: true,
	useBite: true,
	mangleSpam: false,
	biteModeType: FeralDruid_Rotation_BiteModeType.Emperical,
	biteTime: 11.0,
	berserkBiteTime: 6.0,
	minRoarOffset: 31.0,
	ripLeeway: 1.0,
	maintainFaerieFire: true,
	snekWeave: true,
	manualParams: false,
	biteDuringExecute: true,
	allowAoeBerserk: false,
	meleeWeave: true,
	cancelPrimalMadness: false,
});

export const SIMPLE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('Single Target Default', Spec.SpecFeralDruid, DefaultRotation);

export const AoeRotation = FeralDruidRotation.create({
	rotationType: FeralDruid_Rotation_AplType.Aoe,
	bearWeave: true,
	maintainFaerieFire: false,
	snekWeave: true,
	allowAoeBerserk: false,
	cancelPrimalMadness: false,
});

export const AOE_ROTATION_DEFAULT = PresetUtils.makePresetSimpleRotation('AoE Default', Spec.SpecFeralDruid, AoeRotation);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Mono-Cat',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major3: DruidMajorGlyph.GlyphOfRebirth,
		}),
	}),
};

export const HybridTalents = {
	name: 'Hybrid',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			major2: DruidMajorGlyph.GlyphOfMaul,
			major3: DruidMajorGlyph.GlyphOfRebirth,
		}),
	}),
};

export const DefaultOptions = FeralDruidOptions.create({
	assumeBleedActive: true,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62669, // Skewered Eel
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir

	explosiveId: 89637, // Big Daddy Explosive
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	highHpThreshold: 0.8,
	iterationCount: 25000,
	profession1: Profession.Engineering,
	profession2: Profession.ProfessionUnknown,
};

export const PRESET_BUILD_DEFAULT = PresetUtils.makePresetBuild('Single Target Default', {
	rotation: SIMPLE_ROTATION_DEFAULT,
	encounter: PresetUtils.makePresetEncounter(
		'Single Target Default',
		'http://localhost:5173/mop/druid/feral/?i=rcmxe#eJzjEuNgzGBsYGScwMi4gpFxByNjAxPjBiZGJyYPRiEGqUNMs5jZAnISK1OLOLgFGJV4OZgMJAOYIpgqQBqcGLJYpJgUGE8wsdxiYnjE9ItRgknpKyPXJ8ZqpaTUxKLw1MSyVCWrkqLSVB2l3MTMvBIgdktMLcpMdcssQshk5jnn5yblF7vlFwVlFihZmeoolRanBiVmw5UAuU6ZJXBuEpAdkpkL5BsaAnmpRcWpRdlOcEEzVDMhOk0h2lxKizLz0l0rUpNLEeYVZRb4pKaWJ1YCDQTrDcpPLPJPSytOLVGyMgYKFeelZqP6JjUnNRVJpPYFU0ojMwMYWDoshLIiHbqYGZSOM3kwc4L5B4ocBCEyfg6Ss2aCwEl7S4jIBXvFNDC4Zu8IkXppb9TDVLDqM2MVd1BiZopCSGJRemqJQoQEu9YNRgZ6gIAWB2oa15ByHNk8H4u5cxzR1YBDo2ERp+NMRkgo3LSHqmFxAABYiZHH',
	),
});

export const PRESET_BUILD_TENDON = PresetUtils.makePresetBuild('Single Target Burst', {
	rotation: APL_ROTATION_TENDON,
	encounter: PresetUtils.makePresetEncounter(
		'Single Target Burst',
		'http://localhost:5173/mop/druid/feral/?i=rcmxe#eJzjEuZgzGBsYGScwMi4gpFxByNjAxOjE5MHoxCDVA/zLGa2gJzEytQiDm4BRiVuDiYDyQCmCpBaJ4YsFikmBcYTTCy3mBgeMR1jkmDmEubiyGLjYuFoms2sxM7FysWsa1oMF/z3gwUqaFjMJcLFLgVkcjzUUOLkAorqGugBlYpycUiBlM7rZEYSFtKW0uSSl5Ll4tjECNHDJajFz8EsxOTFIAU20dCwGKzvXyOrULxULFewVCCXoZA+kgZlLUWoBslNTGIcjEKcqxihNkGMMDIrRjcVJMS5Ca4MSAssPMYsJColjCbMsfMsoxAwNKwYIJIgl6cl5hSnwjwjJCIlhCwMctRbDSFeKe5JjBwSjBGMCcA4gJjwgimlkZkBDEQcFkJZkQ5dzAxKx5k8mDkhAsYOghDGB3vJWTNB4KS9JUTkgr1iGhhcs3eESL20N+phKlj1mbGKOygxM0UhJLEoPbVEIUKCXesGIwM9QECLAzWNa0g5jmyej8XcOY7oasCh0bCI03EmIyQUbtpD1bA4AADkI2mj',
	),
});
