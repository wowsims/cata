import * as PresetUtils from '../../core/preset_utils.js';
import { Conjured, Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, RaidBuffs, Spec } from '../../core/proto/common.js';
import {
	HolyPaladin_Options as Paladin_Options,
	PaladinAura,
	PaladinSeal,
	PaladinPrimeGlyph as PrimeGlyph,
	PaladinMajorGlyph as MajorGlyph,
	PaladinMinorGlyph as MinorGlyph,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
// export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
// export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
// export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '03331001221131312301-3-032002',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfHolyShock,
			prime2: PrimeGlyph.GlyphOfSealOfInsight,
			prime3: PrimeGlyph.GlyphOfDivineFavor,
			major1: MajorGlyph.GlyphOfDivinePlea,
			major2: MajorGlyph.GlyphOfDivinity,
			major3: MajorGlyph.GlyphOfTheAsceticCrusader,
			minor1: MinorGlyph.GlyphOfInsight,
			minor2: MinorGlyph.GlyphOfBlessingOfKings,
			minor3: MinorGlyph.GlyphOfBlessingOfMight,
		}),
	}),
};

export const DefaultOptions = Paladin_Options.create({
	classOptions: {
		aura: PaladinAura.Devotion,
		seal: PaladinSeal.Insight,
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

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.VolcanicPotion,
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 40,
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
};
