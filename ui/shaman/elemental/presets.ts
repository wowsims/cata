import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Glyphs, Potions, Profession } from '../../core/proto/common.js';
import {
	AirTotem,
	EarthTotem,
	ElementalShaman_Options as ElementalShamanOptions,
	FireTotem,
	ShamanPrimeGlyph,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
	ShamanTotems,
	WaterTotem,
	TotemSet,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const TalentsTotemDuration = {
	name: 'Totem Duration',
	data: SavedTalents.create({
		talentsString: '303202321223110132-201-20302',
		glyphs: Glyphs.create({			
			prime1: ShamanPrimeGlyph.GlyphOfFlameShock,
			prime2: ShamanPrimeGlyph.GlyphOfLavaBurst,
			prime3: ShamanPrimeGlyph.GlyphOfLightningBolt,
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfThunder,
			major3: ShamanMajorGlyph.GlyphOfFireNova,
			minor1: ShamanMinorGlyph.GlyphOfThunderstorm,
			minor2: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfTheArcticWolf,
		}),
	}),
};

export const TalentsImprovedShields = {
	name: 'Improved Shields',
	data: SavedTalents.create({
		talentsString: '3032023212231101321-2030022',
		glyphs: Glyphs.create({			
			prime1: ShamanPrimeGlyph.GlyphOfFlameShock,
			prime2: ShamanPrimeGlyph.GlyphOfLavaBurst,
			prime3: ShamanPrimeGlyph.GlyphOfLightningBolt,
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfThunder,
			major3: ShamanMajorGlyph.GlyphOfFireNova,
			minor1: ShamanMinorGlyph.GlyphOfThunderstorm,
			minor2: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfTheArcticWolf,
		}),
	}),
};

export const DefaultOptions = ElementalShamanOptions.create({
	classOptions: {
		shield: ShamanShield.LightningShield,
		totems: ShamanTotems.create({
			elements: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			ancestors: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			spirits: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			earth: EarthTotem.StrengthOfEarthTotem,
			air: AirTotem.WrathOfAirTotem,
			fire: FireTotem.SearingTotem,
			water: WaterTotem.ManaSpringTotem,
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
	prepopPotion: Potions.VolcanicPotion,
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
});
