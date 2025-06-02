import * as Mechanics from '../../core/constants/mechanics.js';
import * as PresetUtils from '../../core/preset_utils.js';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl.js';
import { ConsumesSpec, Glyphs, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import { PaladinMajorGlyph, PaladinMinorGlyph, PaladinSeal, ProtectionPaladin_Options as ProtectionPaladinOptions } from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1_Gear from './gear_sets/p1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1_Gear);

export const APL_PRESET = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.14,
			[Stat.StatAgility]: 0.13,
			[Stat.StatStamina]: 1.34,
			[Stat.StatHitRating]: 3.69,
			[Stat.StatCritRating]: 1.42,
			[Stat.StatHasteRating]: 3.0,
			[Stat.StatExpertiseRating]: 3.51,
			[Stat.StatDodgeRating]: 0.98,
			[Stat.StatParryRating]: 0.97,
			[Stat.StatMasteryRating]: 1.68,
			[Stat.StatArmor]: 1.0,
			[Stat.StatBonusArmor]: 0.89,
			[Stat.StatAttackPower]: 1.0,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.99,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 2.619 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 1.067 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '112222',
		glyphs: Glyphs.create({
			major1: PaladinMajorGlyph.GlyphOfFocusedShield,
			major2: PaladinMajorGlyph.GlyphOfTheAlabasterShield,
			major3: PaladinMajorGlyph.GlyphOfDivineProtection,

			minor1: PaladinMinorGlyph.GlyphOfFocusedWrath,
		}),
	}),
};

export const P1_BUILD_PRESET = PresetUtils.makePresetBuild('P1', {
	gear: P1_GEAR_PRESET,
	epWeights: P1_EP_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
});

export const DefaultOptions = ProtectionPaladinOptions.create({
	classOptions: {
		seal: PaladinSeal.Insight,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76087, // Flask of the Earth
	foodId: 74656, // Chun Tian Spring Rolls
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
	tinkerId: 126734, // Synapse Springs Mark II
});

export const OtherDefaults = {
	profession1: Profession.Blacksmithing,
	profession2: Profession.Enchanting,
	distanceFromTarget: 5,
	iterationCount: 25000,
};
