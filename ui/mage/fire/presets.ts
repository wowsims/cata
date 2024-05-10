import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, RaidBuffs, Spec } from '../../core/proto/common';
import {
	FireMage_Options as MageOptions,
	MageMajorGlyph as MajorGlyph,
	MageMinorGlyph as MinorGlyph,
	MagePrimeGlyph as PrimeGlyph,
} from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import FireApl from './apls/fire.apl.json';
//import FireAoeApl from './apls/fire_aoe.apl.json';
import P1FireGear from './gear_sets/p1_fire.gear.json';


// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FIRE_P1_PRESET = PresetUtils.makePresetGear('Fire P1 Preset', P1FireGear, { talentTree: 1 });

/* export const DefaultSimpleRotation = MageRotation.create({
	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,
}); */

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Fire', FireApl, { talentTree: 1 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '203-230330221120121213031-03',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfFireball,
			prime2: PrimeGlyph.GlyphOfPyroblast,
			prime3: PrimeGlyph.GlyphOfMoltenArmor,
			major1: MajorGlyph.GlyphOfEvocation,
			major2: MajorGlyph.GlyphOfDragonSBreath,
			major3: MajorGlyph.GlyphOfInvisibility,
			minor1: MinorGlyph.GlyphOfMirrorImage,
			minor2: MinorGlyph.GlyphOfArmors,
			minor3: MinorGlyph.GlyphOfTheMonkey,
		}),
	}),
};

export const DefaultFireOptions = MageOptions.create({
	classOptions: {
	},
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	markOfTheWild: true,
	icyTalons: true,
	moonkinForm: true,
	leaderOfThePack: true,
	powerWordFortitude: true,
	strengthOfEarthTotem: true,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: true,
	blessingOfKings: true,
	blessingOfMight: true,
	communion: true,
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
});

export const DefaultDebuffs = Debuffs.create({
    ebonPlaguebringer: true,
    shadowAndFlame: true,
    judgement: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
