import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Debuffs, IndividualBuffs, Profession, RaidBuffs, Stat } from '../../core/proto/common.js';
import { DisciplinePriest_Options as Options, PriestOptions_Armor } from '../../core/proto/priest.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import AOE24Apl from './apls/aoe_2_4.apl.json';
import AOE4PlusApl from './apls/aoe_4_plus.apl.json';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_AOE24 = PresetUtils.makePresetAPLRotation('AOE (2 to 4 targets)', AOE24Apl);
export const ROTATION_PRESET_AOE4PLUS = PresetUtils.makePresetAPLRotation('AOE (4+ targets)', AOE4PlusApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 0.11,
		[Stat.StatSpirit]: 0.47,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 0.87,
		[Stat.StatCritRating]: 0.74,
		[Stat.StatHasteRating]: 1.65,
		[Stat.StatMP5]: 0.0,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		// talentsString: '05032031--325023051223010323151301351',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfShadow,
		// 	major2: MajorGlyph.GlyphOfMindFlay,
		// 	major3: MajorGlyph.GlyphOfDispersion,
		// 	minor1: MinorGlyph.GlyphOfFortitude,
		// 	minor2: MinorGlyph.GlyphOfShadowProtection,
		// 	minor3: MinorGlyph.GlyphOfShadowfiend,
		// }),
	}),
};

export const EnlightenmentTalents = {
	name: 'Enlightenment',
	data: SavedTalents.create({
		// talentsString: '05032031303005022--3250230012230101231513011',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfShadow,
		// 	major2: MajorGlyph.GlyphOfMindFlay,
		// 	major3: MajorGlyph.GlyphOfShadowWordDeath,
		// 	minor1: MinorGlyph.GlyphOfFortitude,
		// 	minor2: MinorGlyph.GlyphOfShadowProtection,
		// 	minor3: MinorGlyph.GlyphOfShadowfiend,
		// }),
	}),
};

export const DefaultOptions = Options.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 123, // Flask of the Frost Wyrm (not found in list)
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
});
export const DefaultRaidBuffs = RaidBuffs.create({});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	// bloodFrenzy: true,
	// sunderArmor: true,
	// ebonPlaguebringer: true,
	// mangle: true,
	// criticalMass: true,
	// demoralizingShout: true,
	// frostFever: true,
});

export const OtherDefaults = {
	channelClipDelay: 100,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
