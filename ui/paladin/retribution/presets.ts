import * as PresetUtils from '../../core/preset_utils.js';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions } from '../../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinJudgement as PaladinJudgement,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	RetributionPaladin_Options as RetributionPaladinOptions,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3MaceGear from './gear_sets/p3_mace.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import P5Gear from './gear_sets/p5.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Mace Preset', P3MaceGear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);
export const P5_PRESET = PresetUtils.makePresetGear('P5 Preset', P5Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const AuraMasteryTalents = {
	name: 'Aura Mastery',
	data: SavedTalents.create({
		// talentsString: '050501-05-05232051203331302133231331',
		// glyphs: Glyphs.create({
		// 	major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
		// 	major2: PaladinMajorGlyph.GlyphOfJudgement,
		// 	major3: PaladinMajorGlyph.GlyphOfReckoning,
		// 	minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
		// 	minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
		// 	minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings,
		// }),
	}),
};

export const DivineSacTalents = {
	name: 'Divine Sacrifice & Guardian',
	data: SavedTalents.create({
		// talentsString: '03-453201002-05222051203331302133201331',
		// glyphs: Glyphs.create({
		// 	major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
		// 	major2: PaladinMajorGlyph.GlyphOfJudgement,
		// 	major3: PaladinMajorGlyph.GlyphOfReckoning,
		// 	minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
		// 	minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
		// 	minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings,
		// }),
	}),
};

export const DefaultOptions = RetributionPaladinOptions.create({
	classOptions: {
		aura: PaladinAura.RetributionAura,
		judgement: PaladinJudgement.JudgementOfWisdom,
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	defaultConjured: Conjured.ConjuredDarkRune,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodDragonfinFilet,
});
