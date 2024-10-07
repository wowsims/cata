import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands, UnitReference } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph, UnholyDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import SingleTargetApl from '../../death_knight/unholy/apls/st.apl.json';
import P1BISGear from '../../death_knight/unholy/gear_sets/p1.bis.gear.json';
import P1PreBISGear from '../../death_knight/unholy/gear_sets/p1.prebis.gear.json';
import P1RealisticBISGear from '../../death_knight/unholy/gear_sets/p1.realistic-bis.gear.json';
import P3BISGear from '../../death_knight/unholy/gear_sets/p3.bis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const P1_PREBIS_GEAR_PRESET = PresetUtils.makePresetGear('P1 - Preraid', P1PreBISGear);
export const P1_REALISTIC_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P1 - Realistic BIS', P1RealisticBISGear);
export const P1_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);
export const P3_BIS_GEAR_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3BISGear);

export const SINGLE_TARGET_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Single Target', SingleTargetApl);

export const AOE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('AoE', SingleTargetApl);

// Preset options for EP weights
export const P1_UNHOLY_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 4.49,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.94,
			[Stat.StatHasteRating]: 2.4,
			[Stat.StatHitRating]: 2.6,
			[Stat.StatCritRating]: (1.43 + 0.69),
			[Stat.StatMasteryRating]: 1.65,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.13,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const SingleTargetTalents = {
	name: 'Single Target',
	data: SavedTalents.create({
		talentsString: '203-1-13300321230331121231',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfDeathCoil,
			prime2: DeathKnightPrimeGlyph.GlyphOfScourgeStrike,
			prime3: DeathKnightPrimeGlyph.GlyphOfRaiseDead,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathSEmbrace,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const AoeTalents = {
	name: 'AOE',
	data: SavedTalents.create({
		// talentsString: '-320050500002-2302303050032052000150013133151',
		// glyphs: Glyphs.create({
		// 	major1: DeathKnightMajorGlyph.GlyphOfTheGhoul,
		// 	major2: DeathKnightMajorGlyph.GlyphOfIcyTouch,
		// 	major3: DeathKnightMajorGlyph.GlyphOfDeathAndDecay,
		// 	minor1: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		// 	minor2: DeathKnightMinorGlyph.GlyphOfPestilence,
		// 	minor3: DeathKnightMinorGlyph.GlyphOfRaiseDead,
		// }),
	}),
};

export const DefaultOptions = UnholyDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 80,
		petUptime: 1,
	},
	unholyFrenzyTarget: UnitReference.create(),
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
	distanceFromTarget: 5,
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});
