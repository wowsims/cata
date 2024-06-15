import * as PresetUtils from '../../core/preset_utils.js';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PartyBuffs,
	Potions,
	Profession,
	RaidBuffs,
	Stat,
	UnitReference,
} from '../../core/proto/common.js';
import { BalanceDruid_Options as BalanceDruidOptions, DruidMajorGlyph, DruidMinorGlyph, DruidPrimeGlyph } from '../../core/proto/druid.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import PreraidGear from './gear_sets/preraid.gear.json';
import T11Gear from './gear_sets/t11.gear.json';

export const PreraidPresetGear = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const T11PresetGear = PresetUtils.makePresetGear('T11', T11Gear);

export const PresetRotationDefault = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 0.43,
		[Stat.StatSpirit]: 0.34,
		[Stat.StatSpellPower]: 1,
		[Stat.StatSpellCrit]: 0.82,
		[Stat.StatSpellHaste]: 0.8,
		[Stat.StatMastery]: 0.0,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '33230221123212111001-01-020331',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfInsectSwarm,
			prime2: DruidPrimeGlyph.GlyphOfMoonfire,
			prime3: DruidPrimeGlyph.GlyphOfWrath,
			major1: DruidMajorGlyph.GlyphOfStarfall,
			major2: DruidMajorGlyph.GlyphOfRebirth,
			major3: DruidMajorGlyph.GlyphOfFocus,
			minor1: DruidMinorGlyph.GlyphOfTyphoon,
			minor2: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
			minor3: DruidMinorGlyph.GlyphOfMarkOfTheWild,
		}),
	}),
};

export const DefaultOptions = BalanceDruidOptions.create({
	classOptions: {
		innervateTarget: UnitReference.create(),
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
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

export const DefaultPartyBuffs = PartyBuffs.create({});

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
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	darkIntentUptime: 100,
};
