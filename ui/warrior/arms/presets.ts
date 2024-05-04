import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Faction, Flask, Food, Glyphs, Potions, Profession, Race, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { ArmsWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorPrimeGlyph, WarriorShout } from '../../core/proto/warrior';
import ArmsApl from './apls/arms.apl.json';
import P1ArmsBisGear from './gear_sets/p1_arms_bis.gear.json';
import PreraidArmsGear from './gear_sets/preraid_arms.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_ARMS_PRESET = PresetUtils.makePresetGear('Preraid Arms', PreraidArmsGear);
export const P1_ARMS_PRESET = PresetUtils.makePresetGear('P1 Arms', P1ArmsBisGear);

export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Arms', ArmsApl);

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
		useShatteringThrow: true,
		shout: WarriorShout.WarriorShoutCommanding,
	},
	useRecklessness: true,
	disableExpertiseGemming: false,
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
