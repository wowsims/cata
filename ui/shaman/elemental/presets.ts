import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Faction, Flask, Food, Glyphs, Potions, Profession } from '../../core/proto/common.js';
import {
	AirTotem,
	EarthTotem,
	ElementalShaman_Options as ElementalShamanOptions,
	FireTotem,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
	ShamanTotems,
	WaterTotem,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import AdvancedApl from './apls/advanced.apl.json';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P2Gear from './gear_sets/p2.gear.json';
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
import P3HordeGear from './gear_sets/p3_horde.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
export const P3_PRESET_ALLI = PresetUtils.makePresetGear('P3 Preset [A]', P3AllianceGear, { faction: Faction.Alliance });
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3 Preset [H]', P3HordeGear, { faction: Faction.Horde });
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);
export const ROTATION_PRESET_ADVANCED = PresetUtils.makePresetAPLRotation('Advanced', AdvancedApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		// talentsString: '3232023212231101321-203002',
		// glyphs: Glyphs.create({			
		// 	prime1: ShamanPrimeGlyph.ShamanPrimeGlyphNone,
		// 	prime2: ShamanPrimeGlyph.ShamanPrimeGlyphNone,
		// 	prime3: ShamanPrimeGlyph.ShamanPrimeGlyphNone,
		// 	major1: ShamanMajorGlyph.ShamanMajorGlyphNone,
		// 	major2: ShamanMajorGlyph.ShamanMajorGlyphNone,
		// 	major3: ShamanMajorGlyph.ShamanMajorGlyphNone,
		// 	minor1: ShamanMinorGlyph.ShamanMinorGlyphNone,
		// 	minor2: ShamanMinorGlyph.ShamanMinorGlyphNone,
		// 	minor3: ShamanMinorGlyph.ShamanMinorGlyphNone,
		// }),
	}),
};

export const DefaultOptions = ElementalShamanOptions.create({
	classOptions: {
		shield: ShamanShield.WaterShield,
		totems: ShamanTotems.create({
			earth: EarthTotem.StrengthOfEarthTotem,
			air: AirTotem.WrathOfAirTotem,
			fire: FireTotem.FlametongueTotem,
			water: WaterTotem.ManaSpringTotem,
			useFireElemental: true,
		}),
	},
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.VolcanicPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
});
