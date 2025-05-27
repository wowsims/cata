import * as Mechanics from '../../core/constants/mechanics.js';
import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Spec, Stat } from '../../core/proto/common';
import { DruidMajorGlyph, GuardianDruid_Options as DruidOptions, GuardianDruid_Rotation as DruidRotation } from '../../core/proto/druid.js';
import { SavedTalents } from '../../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1/P2', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);

export const DefaultSimpleRotation = DruidRotation.create({
	maintainFaerieFire: true,
	maintainDemoralizingRoar: true,
	demoTime: 4.0,
	pulverizeTime: 4.0,
	prepullStampede: true,
});

import { Stats } from '../../core/proto_utils/stats';
import BalerocMTApl from './apls/balerocMT.apl.json';
import BalerocOTApl from './apls/balerocOT.apl.json';
import BethApl from './apls/bethtilac.apl.json';
import BlackhornOTApl from './apls/blackhorn.apl.json';
import CleaveApl from './apls/cleave.apl.json';
import DefaultApl from './apls/default.apl.json';
import NefApl from './apls/nef.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);
export const ROTATION_CLEAVE = PresetUtils.makePresetAPLRotation('2-Target Cleave', CleaveApl);
export const ROTATION_NEF = PresetUtils.makePresetAPLRotation('AoE (Nef Adds)', NefApl);
export const ROTATION_BETH = PresetUtils.makePresetAPLRotation("Beth'tilac Phase 2", BethApl);
export const ROTATION_BALEROC_MT = PresetUtils.makePresetAPLRotation('Baleroc MT', BalerocMTApl);
export const ROTATION_BALEROC_OT = PresetUtils.makePresetAPLRotation('Baleroc OT', BalerocOTApl);
export const ROTATION_BLACKHORN_OT = PresetUtils.makePresetAPLRotation('Blackhorn OT', BlackhornOTApl);

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecGuardianDruid, DefaultSimpleRotation);

// Preset options for EP weights
export const SURVIVAL_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Survival',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.05,
			[Stat.StatStamina]: 1.08,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.93,
			[Stat.StatBonusArmor]: 0.21,
			[Stat.StatDodgeRating]: 0.85,
			[Stat.StatMasteryRating]: 0.31,
			[Stat.StatStrength]: 0.08,
			[Stat.StatAttackPower]: 0.08,
			[Stat.StatHitRating]: 0.22,
			[Stat.StatExpertiseRating]: 0.37,
			[Stat.StatCritRating]: 0.28,
			[Stat.StatHasteRating]: 0.06,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.185 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.035 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

export const BALANCED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Balanced',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.04,
			[Stat.StatStamina]: 0.88,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.66,
			[Stat.StatBonusArmor]: 0.15,
			[Stat.StatDodgeRating]: 0.6,
			[Stat.StatMasteryRating]: 0.22,
			[Stat.StatStrength]: 0.16,
			[Stat.StatAttackPower]: 0.15,
			[Stat.StatHitRating]: 0.61,
			[Stat.StatExpertiseRating]: 1.07,
			[Stat.StatCritRating]: 0.36,
			[Stat.StatHasteRating]: 0.1,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.5,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.535 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.075 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			major2: DruidMajorGlyph.GlyphOfMaul,
			major3: DruidMajorGlyph.GlyphOfRebirth,
		}),
	}),
};

export const InfectedWoundsBuild = {
	name: 'Infected Wounds',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			major2: DruidMajorGlyph.GlyphOfMaul,
			major3: DruidMajorGlyph.GlyphOfRebirth,
		}),
	}),
};

export const DefaultOptions = DruidOptions.create({
	startingRage: 10,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58085, // Flask of Steelskin
	foodId: 62669, // Skewered Eel
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
	conjuredId: 5512, // Conjured Healthstone
	tinkerId: 82174, // Synapse Springs
});
export const OtherDefaults = {
	iterationCount: 50000,
	profession1: Profession.Engineering,
	profession2: Profession.ProfessionUnknown,
};

//export const PRESET_BUILD_BOSS_DUMMY = PresetUtils.makePresetBuild('Single Target Dummy', {
//	rotation: ROTATION_PRESET_SIMPLE,
//	encounter: PresetUtils.makePresetEncounter(
//		'Single Target Dummy',
//		'http://localhost:5173/mop/druid/guardian/?i=rcmxe#eJzjEuVgdGDMYGxgZJzAyNjAxLiBifECE6MTpwCjBaMH4w1GRismAUYhBqkvjLOY2QJyEitTizjYBBiVeDmYDSQDmCKYEliBGp0YVjFzS3EKMoCBnsMJJpYLTOy3mDgFZ80EgZv2j5iaGCWYlOq4CquVchMz80qA2C0xtSgz1S2zKFXJqqSoNFUHLuOSmptflJiTWZWZlx6Un1gEk08Biodk5gLVm+goFZTmlAENqEpFiBSlAgVzgksScwtSU6Cm1gohXPGCKeUHE+NCZog7Ix26mGU5wcymKw5wRZpnz4AAl4MCm9JxJg9miIqGNAeo90QcJCFKT9pbQkQu2CumgcE1e0eI5jf2Rj1MBas+M1ZxByVmpiiEJBalp5YoREiwa91gZKAHCGhxoKZxDSnHkc3zsZg7xxFdDTg0GhZxOs5khIUlVA2LAwBf7n5L',
//	),
//});
