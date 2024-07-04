import * as PresetUtils from '../../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	Stat,
	TinkerHands,
	TristateEffect,
} from '../../core/proto/common.js';
import {
	PriestMajorGlyph as MajorGlyph,
	PriestMinorGlyph as MinorGlyph,
	PriestOptions_Armor,
	PriestPrimeGlyph as PrimeGlyph,
	ShadowPriest_Options as Options,
} from '../../core/proto/priest.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreRaidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const PRE_RAID = PresetUtils.makePresetGear('Pre Raid', PreRaidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.0,
		[Stat.StatSpirit]: 0.53,
		[Stat.StatSpellPower]: 0.82,
		[Stat.StatSpellHit]: 0.52,
		[Stat.StatSpellCrit]: 0.42,
		[Stat.StatSpellHaste]: 0.76,
		[Stat.StatMastery]: 0.46,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://www.wowhead.com/cata/talent-calc/priest and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '032212--322032210201222100231',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfShadowWordPain,
			prime2: PrimeGlyph.GlyphOfMindFlay,
			prime3: PrimeGlyph.GlyphOfShadowWordDeath,
			major1: MajorGlyph.GlyphOfFade,
			major2: MajorGlyph.GlyphOfInnerFire,
			major3: MajorGlyph.GlyphOfSpiritTap,
			minor1: MinorGlyph.GlyphOfFading,
			minor2: MinorGlyph.GlyphOfFortitude,
			minor3: MinorGlyph.GlyphOfShadowfiend,
		}),
	}),
};

export const DefaultOptions = Options.create({
	classOptions: {
		armor: PriestOptions_Armor.InnerFire,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
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
	darkIntent: true,
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
	channelClipDelay: 40,
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	darkIntentUptime: 100,
};
