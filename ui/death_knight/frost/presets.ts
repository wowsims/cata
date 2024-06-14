import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Explosive, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import TwoHandAPL from '../../death_knight/frost/apls/2h.apl.json';
import DualWieldAPL from '../../death_knight/frost/apls/dw.apl.json';
import MasterFrostAPL from '../../death_knight/frost/apls/masterfrost.apl.json';
import P12HGear from '../../death_knight/frost/gear_sets/p1.2h.gear.json';
import P1DWGear from '../../death_knight/frost/gear_sets/p1.dw.gear.json';
import P1MasterfrostGear from '../../death_knight/frost/gear_sets/p1.masterfrost.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 Dual Wield', P1DWGear);
export const P1_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 Two Hand', P12HGear);
export const P1_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P1 Masterfrost', P1MasterfrostGear);

export const DUAL_WIELD_ROTATION_RESET_DEFAULT = PresetUtils.makePresetAPLRotation('Dual Wield', DualWieldAPL);
export const TWO_HAND_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Two Hand', TwoHandAPL);
export const MASTERFROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Masterfrost', MasterFrostAPL);

// Preset options for EP weights
export const P1_MASTERFROST_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Masterfrost',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.86,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 0.75,
			[Stat.StatMeleeHaste]: 1.38,
			[Stat.StatMeleeHit]: 1.08,
			[Stat.StatMeleeCrit]: 0.64,
			[Stat.StatSpellHit]: 0.59,
			[Stat.StatSpellCrit]: 0.43,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.5,
			[PseudoStat.PseudoStatOffHandDps]: 2.84,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DualWieldTalents = {
	name: 'Dual Wield',
	data: SavedTalents.create({
		talentsString: '2032-20330022233112012301-003',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const TwoHandTalents = {
	name: 'Two Hand',
	data: SavedTalents.create({
		talentsString: '103-32030022233112012031-033',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const MasterfrostTalents = {
	name: 'Masterfrost',
	data: SavedTalents.create({
		talentsString: '2032-20330022233112012301-03',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const DefaultOptions = FrostDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
		petUptime: 1,
	},
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Enchanting,
	distanceFromTarget: 5,
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});
