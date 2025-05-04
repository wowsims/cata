import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, PseudoStat, Stat } from '../../core/proto/common';
import { AssassinationRogue_Options as RogueOptions, RogueMajorGlyph, RogueOptions_PoisonOptions } from '../../core/proto/rogue';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import AssassinationApl from './apls/assassination.apl.json';
import P1AssassinationGear from './gear_sets/p1_assassination.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET_ASSASSINATION = PresetUtils.makePresetGear('P1 Assassination', P1AssassinationGear);

export const ROTATION_PRESET_ASSASSINATION = PresetUtils.makePresetAPLRotation('Assassination', AssassinationApl);

// Preset options for EP weights
export const ASN_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Asn',
	Stats.fromMap(
		{
			[Stat.StatAgility]: 2.64,
			[Stat.StatStrength]: 1.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatCritRating]: 1.18,
			[Stat.StatHitRating]: 2.62,
			[Stat.StatHasteRating]: 1.35,
			[Stat.StatMasteryRating]: 1.45,
			[Stat.StatExpertiseRating]: 1.2,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.0,
			[PseudoStat.PseudoStatOffHandDps]: 0.97,
			[PseudoStat.PseudoStatSpellHitPercent]: 130.5,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 162.0,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const AssassinationTalentsDefault = {
	name: 'Assassination',
	data: SavedTalents.create({
		talentsString: '300003',
		glyphs: Glyphs.create({
			major1: RogueMajorGlyph.GlyphOfVendetta,
		}),
	}),
};

export const DefaultOptions = RogueOptions.create({
	classOptions: {
		lethalPoison: RogueOptions_PoisonOptions.DeadlyPoison,
		applyPoisonsManually: false,
		startingOverkillDuration: 20,
		vanishBreakTime: 0.1,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58087, // Flask of the Winds
	foodId: 62669, // Skewered Eel
	potId: 58145, // Potion of the Tol'vir
	prepotId: 58145, // Potion of the Tol'vir
});

export const OtherDefaults = {
	distanceFromTarget: 5,
};
