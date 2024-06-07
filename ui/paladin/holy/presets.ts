import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Potions } from '../../core/proto/common.js';
import { HolyPaladin_Options as HolyPaladinOptions, PaladinAura, PaladinMajorGlyph, PaladinMinorGlyph } from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		// talentsString: '50350151020013053100515221-50023131203',
		// glyphs: {
		// 	major1: PaladinMajorGlyph.GlyphOfHolyLight,
		// 	major2: PaladinMajorGlyph.GlyphOfSealOfWisdom,
		// 	major3: PaladinMajorGlyph.GlyphOfBeaconOfLight,
		// 	minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
		// 	minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
		// 	minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings,
		// },
	}),
};

export const DefaultOptions = HolyPaladinOptions.create({
	classOptions: {
		aura: PaladinAura.DevotionAura,
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.MythicalManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});
