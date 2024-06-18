import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Stat, TinkerHands } from '../../core/proto/common.js';
import { BloodDeathKnight_Options, DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefensiveBloodApl from './apls/defensive.apl.json';
import SimpleBloodApl from './apls/simple.apl.json';
import P1BloodGear from './gear_sets/p1.gear.json';
import PreRaidBloodGear from './gear_sets/preraid.gear.json';

export const PRERAID_BLOOD_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreRaidBloodGear);
export const P1_BLOOD_PRESET = PresetUtils.makePresetGear('P1', P1BloodGear);

export const BLOOD_SIMPLE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Simple', SimpleBloodApl);
export const BLOOD_DEFENSIVE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Defensive', DefensiveBloodApl);

// Preset options for EP weights
export const P1_BLOOD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.45,
			[Stat.StatAgility]: 1.2,
			[Stat.StatStamina]: 3,
			[Stat.StatAttackPower]: 1,
			[Stat.StatMeleeHit]: 6,
			[Stat.StatMeleeCrit]: 1.65,
			[Stat.StatMeleeHaste]: 1.58,
			[Stat.StatExpertise]: 5,
			[Stat.StatArmor]: 1,
			[Stat.StatDodge]: 2.5,
			[Stat.StatParry]: 2.44,
			[Stat.StatBonusArmor]: 1,
			[Stat.StatMastery]: 7,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 12.29,
			[PseudoStat.PseudoStatOffHandDps]: 0.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const BloodTalents = {
	name: 'Blood',
	data: SavedTalents.create({
		talentsString: '03323200132222311321-2-003',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfDeathStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfHeartStrike,
			prime3: DeathKnightPrimeGlyph.GlyphOfRuneStrike,
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon,
			major3: DeathKnightMajorGlyph.GlyphOfBoneShield,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
};

export const DefaultOptions = BloodDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	distanceFromTarget: 5,
};
