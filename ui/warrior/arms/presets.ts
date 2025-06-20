import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { ArmsWarrior_Options as WarriorOptions, WarriorMajorGlyph } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';
import ArmsApl from './apls/arms.apl.json';
import P1ArmsBisGear from './gear_sets/p1_arms_bis.gear.json';
import P1PreBisGearPoor from './gear_sets/p1_prebis_poor.gear.json';
import P1PreBisGearRich from './gear_sets/p1_prebis_rich.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PREBIS_ARMS_RICH_PRESET = PresetUtils.makePresetGear('P1 - Pre-BIS 💰', P1PreBisGearRich);
export const P1_PREBIS_ARMS_POOR_PRESET = PresetUtils.makePresetGear('P1 - Pre-BIS 📉', P1PreBisGearPoor);
export const P1_ARMS_BIS_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1ArmsBisGear);

export const ROTATION_ARMS = PresetUtils.makePresetAPLRotation('Default', ArmsApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1,
			[Stat.StatAttackPower]: 0.45,
			[Stat.StatExpertiseRating]: 1.20,
			[Stat.StatHitRating]: 1.40,
			[Stat.StatCritRating]: 0.59,
			[Stat.StatHasteRating]: 0.32,
			[Stat.StatMasteryRating]: 0.39,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.71,
			[PseudoStat.PseudoStatOffHandDps]: 0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const ArmsTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '113332',
		glyphs: Glyphs.create({
			major1: WarriorMajorGlyph.GlyphOfBullRush,
			major2: WarriorMajorGlyph.GlyphOfUnendingRage,
			major3: WarriorMajorGlyph.GlyphOfDeathFromAbove,
		}),
	}),
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		startingRage: 0,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76088, // Flask of Winter's Bite
	foodId: 74646, // Black Pepper Ribs and Shrimp
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 9,
};
