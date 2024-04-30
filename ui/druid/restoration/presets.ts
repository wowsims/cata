import * as PresetUtils from '../../core/preset_utils';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	RaidBuffs,
	TristateEffect,
	UnitReference,
} from '../../core/proto/common';
import { DruidMajorGlyph, DruidMinorGlyph, RestorationDruid_Options as RestorationDruidOptions } from '../../core/proto/druid';
import { SavedTalents } from '../../core/proto/ui';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const CelestialFocusTalents = {
	name: 'Celestial Focus',
	data: SavedTalents.create({
		// talentsString: '05320031103--230023312131502331050313051',
		// glyphs: Glyphs.create({
		// 	major1: DruidMajorGlyph.GlyphOfWildGrowth,
		// 	major2: DruidMajorGlyph.GlyphOfSwiftmend,
		// 	major3: DruidMajorGlyph.GlyphOfNourish,
		// 	minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		// 	minor3: DruidMinorGlyph.GlyphOfTheWild,
		// 	minor1: DruidMinorGlyph.GlyphOfDash,
		// }),
	}),
};
export const ThiccRestoTalents = {
	name: 'Thicc Resto',
	data: SavedTalents.create({
		// talentsString: '05320001--230023312331502531053313051',
		// glyphs: Glyphs.create({
		// 	major1: DruidMajorGlyph.GlyphOfWildGrowth,
		// 	major2: DruidMajorGlyph.GlyphOfSwiftmend,
		// 	major3: DruidMajorGlyph.GlyphOfNourish,
		// 	minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		// 	minor3: DruidMinorGlyph.GlyphOfTheWild,
		// 	minor1: DruidMinorGlyph.GlyphOfDash,
		// }),
	}),
};

export const DefaultOptions = RestorationDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.MythicalManaPotion,
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
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

export const DefaultIndividualBuffs = IndividualBuffs.create({
	vampiricTouch: true,
});

export const DefaultPartyBuffs = PartyBuffs.create({
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
	judgement: true,
});

export const OtherDefaults = {
	distanceFromTarget: 18,
};
