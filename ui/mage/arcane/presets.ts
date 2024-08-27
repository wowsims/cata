import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, RaidBuffs, Stat, TinkerHands, UnitReference } from '../../core/proto/common';
import { ArcaneMage_Options as MageOptions, MageMajorGlyph as MajorGlyph, MagePrimeGlyph as PrimeGlyph } from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import ArcaneApl from './apls/arcane.apl.json';
import P1ArcaneBisGear from './gear_sets/p1_arcane.gear.json';
import P1ArcanePrebisGear from './gear_sets/P1_arcane_prebis_gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
export const ARCANE_P1_PRESET = PresetUtils.makePresetGear('Arcane P1 Preset', P1ArcaneBisGear, { talentTree: 0 });
export const ARCANE_P1_PREBIS_PRESET = PresetUtils.makePresetGear('Arcane P1 Preraid Preset', P1ArcanePrebisGear, { talentTree: 0 });
/* export const DefaultSimpleRotation = MageRotation.create({
	only3ArcaneBlastStacksBelowManaPercent: 0.15,
	blastWithoutMissileBarrageAboveManaPercent: 0.2,
	missileBarrageBelowManaPercent: 0,
	useArcaneBarrage: false,
}); */

//export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecArcaneMage, DefaultSimpleRotation);
export const ARCANE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Arcane', ArcaneApl, { talentTree: 0 });
//export const ARCANE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Arcane AOE', ArcaneAoeApl, { talentTree: 0 });

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Arcane P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.76,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 1.31,
		[Stat.StatCritRating]: 0.53,
		[Stat.StatHasteRating]: 0.91,
		[Stat.StatMasteryRating]: 0.57,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const ArcaneTalents = {
	name: 'Arcane',
	data: SavedTalents.create({
		talentsString: '303302221230122210121-23-03',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfArcaneMissiles,
			prime2: PrimeGlyph.GlyphOfArcaneBlast,
			prime3: PrimeGlyph.GlyphOfMageArmor,
			major1: MajorGlyph.GlyphOfEvocation,
			major2: MajorGlyph.GlyphOfArcanePower,
			major3: MajorGlyph.GlyphOfManaShield,
		}),
	}),
};

export const DefaultArcaneOptions = MageOptions.create({
	classOptions: {},
	focusMagicPercentUptime: 99,
	focusMagicTarget: UnitReference.create(),
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

export const DefaultArcaneConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
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
