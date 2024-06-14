import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Faction, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { FuryWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorPrimeGlyph, WarriorShout } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';
import FuryApl from './apls/fury.apl.json';
import P1FurySMFGear from './gear_sets/p1_fury_smf.gear.json';
import P1FuryTGGear from './gear_sets/p1_fury_tg.gear.json';
import PreraidFurySMFGear from './gear_sets/preraid_fury_smf.gear.json';
import PreraidFuryTGGear from './gear_sets/preraid_fury_tg.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_FURY_SMF_PRESET = PresetUtils.makePresetGear('Preraid Fury SMF', PreraidFurySMFGear);
export const PRERAID_FURY_TG_PRESET = PresetUtils.makePresetGear('Preraid Fury TG', PreraidFuryTGGear);
export const P1_FURY_SMF_PRESET = PresetUtils.makePresetGear('P1 Fury SMF', P1FurySMFGear);
export const P1_FURY_TG_PRESET = PresetUtils.makePresetGear('P1 Fury TG', P1FuryTGGear);

export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('Fury', FuryApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.2,
			[Stat.StatAgility]: 1.18,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 1.6,
			[Stat.StatMeleeHit]: 1.24,
			[Stat.StatMeleeCrit]: 1.54,
			[Stat.StatMeleeHaste]: 1.01,
			[Stat.StatMastery]: 0.92,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.63,
			[PseudoStat.PseudoStatOffHandDps]: 1.4,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const FurySMFTalents = {
	name: 'Fury SMF',
	data: SavedTalents.create({
		talentsString: '302203-032222031301101223201',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfBloodthirst,
			prime2: WarriorPrimeGlyph.GlyphOfRagingBlow,
			prime3: WarriorPrimeGlyph.GlyphOfSlam,
			major1: WarriorMajorGlyph.GlyphOfCleaving,
			major2: WarriorMajorGlyph.GlyphOfDeathWish,
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfBattle,
			minor3: WarriorMinorGlyph.GlyphOfBerserkerRage,
		}),
	}),
};

export const FuryTGTalents = {
	name: 'Fury TG',
	data: SavedTalents.create({
		talentsString: '302203-03222203130110122321',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfBloodthirst,
			prime2: WarriorPrimeGlyph.GlyphOfRagingBlow,
			prime3: WarriorPrimeGlyph.GlyphOfSlam,
			major1: WarriorMajorGlyph.GlyphOfCleaving,
			major2: WarriorMajorGlyph.GlyphOfDeathWish,
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfBattle,
			minor3: WarriorMinorGlyph.GlyphOfBerserkerRage,
		}),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		startingRage: 0,
		useShatteringThrow: true,
		shout: WarriorShout.WarriorShoutCommanding,
	},
	useRecklessness: true,
	disableExpertiseGemming: false,
	syncType: 0,
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
