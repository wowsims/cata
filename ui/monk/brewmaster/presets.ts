import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common';
import { BrewmasterMonk_Options as BrewmasterMonkOptions, MonkMajorGlyph, MonkMinorGlyph, MonkStance } from '../../core/proto/monk';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1BISBalanced2HGear from './gear_sets/p1_bis_balanced_2h.gear.json';
import P1BISBalancedDWGear from './gear_sets/p1_bis_balanced_dw.gear.json';
import P1BISOffensive2HGear from './gear_sets/p1_bis_offensive_2h.gear.json';
import P1BISOffensiveDWGear from './gear_sets/p1_bis_offensive_dw.gear.json';
import P1PreBISPoorGear from './gear_sets/p1_prebis_poor.gear.json';
import P1PreBISRichGear from './gear_sets/p1_prebis_rich.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PREBIS_RICH_GEAR_PRESET = PresetUtils.makePresetGear('P1 - Pre-BIS 💰', P1PreBISRichGear);
export const P1_PREBIS_POOR_GEAR_PRESET = PresetUtils.makePresetGear('P1 - Pre-BIS 📉', P1PreBISPoorGear);

export const P1_BIS_BALANCED_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 - BIS DW (Balanced)', P1BISBalancedDWGear);
export const P1_BIS_BALANCED_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 - BIS 2H (Balanced)', P1BISBalanced2HGear);

export const P1_BIS_OFFENSIVE_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 - BIS DW (Offensive)', P1BISOffensiveDWGear);
export const P1_BIS_OFFENSIVE_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 - BIS 2H (Offensive)', P1BISOffensive2HGear);

export const ROTATION_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const PREPATCH_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 3.61,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 3.11,
			[Stat.StatHitRating]: 6.26,
			[Stat.StatHasteRating]: 3.57,
			[Stat.StatMasteryRating]: 1.60,
			[Stat.StatDodgeRating]: 0.24,
			[Stat.StatParryRating]: 0.36,
			[Stat.StatExpertiseRating]: 7.02,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 10.66,
			[PseudoStat.PseudoStatOffHandDps]: 5.28,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '213322',
		glyphs: Glyphs.create({
			major1: MonkMajorGlyph.GlyphOfFortifyingBrew,
			major2: MonkMajorGlyph.GlyphOfEnduringHealingSphere,
			major3: MonkMajorGlyph.GlyphOfFortuitousSpheres,
			minor1: MonkMinorGlyph.GlyphOfSpiritRoll,
			minor2: MonkMinorGlyph.GlyphOfJab,
			minor3: MonkMinorGlyph.GlyphOfWaterRoll,
		}),
	}),
};

export const DungeonTalents = {
	name: 'Raid',
	data: SavedTalents.create({
		talentsString: '213321',
		glyphs: Glyphs.create({
			major1: MonkMajorGlyph.GlyphOfFortifyingBrew,
			major2: MonkMajorGlyph.GlyphOfBreathOfFire,
			major3: MonkMajorGlyph.GlyphOfRapidRolling,
			minor1: MonkMinorGlyph.GlyphOfSpiritRoll,
			minor2: MonkMinorGlyph.GlyphOfJab,
			minor3: MonkMinorGlyph.GlyphOfWaterRoll,
		}),
	}),
};

export const DefaultOptions = BrewmasterMonkOptions.create({
	classOptions: {},
	stance: MonkStance.SturdyOx,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
