import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { ArmsWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorPrimeGlyph } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';
import ArmsApl from './apls/arms.apl.json';
import P1ArmsBisGear from './gear_sets/p1_arms_bis.gear.json';
import P1ArmsRealisticBisGear from './gear_sets/p1_arms_realistic_bis.gear.json';
import PreraidArmsGear from './gear_sets/preraid_arms.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_ARMS_PRESET = PresetUtils.makePresetGear('Preraid Arms', PreraidArmsGear);
export const P1_ARMS_BIS_PRESET = PresetUtils.makePresetGear('P1 Arms - BIS', P1ArmsBisGear);
export const P1_ARMS_REALISTIC_PRESET = PresetUtils.makePresetGear('P1 Arms - Realistic', P1ArmsRealisticBisGear);

export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Arms', ArmsApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.21,
			[Stat.StatAgility]: 1.12,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 1.75,
			[Stat.StatMeleeHit]: 2.77,
			[Stat.StatMeleeCrit]: 1.45,
			[Stat.StatMeleeHaste]: 0.68,
			[Stat.StatMastery]: 0.89,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 9.22,
			[PseudoStat.PseudoStatOffHandDps]: 0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const ArmsTalents = {
	name: 'Arms',
	data: SavedTalents.create({
		talentsString: '30220303120212312211-0322-3',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfMortalStrike,
			prime2: WarriorPrimeGlyph.GlyphOfOverpower,
			prime3: WarriorPrimeGlyph.GlyphOfSlam,
			major1: WarriorMajorGlyph.GlyphOfCleaving,
			major2: WarriorMajorGlyph.GlyphOfSweepingStrikes,
			major3: WarriorMajorGlyph.GlyphOfThunderClap,
			minor1: WarriorMinorGlyph.GlyphOfBerserkerRage,
			minor2: WarriorMinorGlyph.GlyphOfCommand,
			minor3: WarriorMinorGlyph.GlyphOfBattle,
		}),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		startingRage: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};
