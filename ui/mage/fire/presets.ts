import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, RaidBuffs, Spec, Stat } from '../../core/proto/common';
import {
	FireMage_Options as MageOptions,
	MageMajorGlyph as MajorGlyph,
	MageMinorGlyph as MinorGlyph,
	MagePrimeGlyph as PrimeGlyph,
} from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import FireApl from './apls/fire.apl.json';
//import FireAoeApl from './apls/fire_aoe.apl.json';
import P1FireBisGear from './gear_sets/p1_fire.gear.json';
import P1FirePrebisGear from './gear_sets/p1_fire_prebis_gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FIRE_P1_PRESET = PresetUtils.makePresetGear('Fire P1 Preset', P1FireBisGear, { talentTree: 1 });
export const FIRE_P1_PREBIS = PresetUtils.makePresetGear('Fire P1 Pre-raid', P1FirePrebisGear, { talentTree: 1 });

/* export const DefaultSimpleRotation = MageRotation.create({
	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,
}); */

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Fire', FireApl, { talentTree: 1 });

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Fire P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.33,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatSpellHit]: 1.09,
		[Stat.StatSpellCrit]: 0.62,
		[Stat.StatSpellHaste]: 0.82,
		[Stat.StatMastery]: 0.46,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		talentsString: '003-230330221120121213231-03',
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
	classOptions: {},
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
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const FIRE_BREAKPOINTS = new Map([
	[
		Stat.StatSpellHaste,
		// Picked from Mage Discord
		// Sources:
		// https://docs.google.com/spreadsheets/d/17cJJUReg2uz-XxBB3oDWb1kCncdH_-X96mSb0HAu4Ko/edit?gid=0#gid=0
		// https://docs.google.com/spreadsheets/d/1WLOZ1YevGPw_WZs0JhGzVVy906W5y0i9UqHa3ejyBkE/htmlview?gid=19
		new Map([
			['BL - 15-tick Combust', 1481],
			['5-tick LvB/Pyro', 1602],
			['12-tick Combust', 1922],
			['BL - 16-tick Combust', 2455],
			['BL - 7-tick LvB/Pyro', 3199],
			['13-tick Combust', 3212],
			['BL - 17-tick Combust', 3436],
			['14-tick Combust', 4488],
			['6-tick LvB/Pyro', 4805],
			['15-tick Combust', 5767],
			['16-tick Combust', 7033],
			['7-tick LvB/Pyro', 8000],
			['17-tick Combust', 8309],
			['18-tick Combust', 9602],
			['19-tick Combust', 10887],
			['8-tick LvB/Pyro', 11198],
			['20-tick Combust', 12182],
			['21-tick Combust', 13463],
			// ['9-tick LvB/Pyro', 14412],
			// ['22-tick Combust', 14704],
			// ['23-tick Combust', 16004],
			// ['24-tick Combust', 17290],
			// ['10-tick LvB/Pyro', 17600],
			// ['25-tick Combust', 18543],
			// ['26-tick Combust', 19821],
			// ['11-tick LvB/Pyro', 20820],
			// ['27-tick Combust', 21117],
			// ['28-tick Combust', 22424],
			// ['29-tick Combust', 23730],
			// ['12-tick LvB/Pyro', 24010],
		]),
	],
]);
