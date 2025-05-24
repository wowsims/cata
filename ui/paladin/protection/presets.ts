import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Profession, PseudoStat, Stat } from '../../core/proto/common.js';
import {
	PaladinMajorGlyph,
	PaladinSeal,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/P1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_GEAR_PRESET = PresetUtils.makePresetGear('P1', P1Gear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatArmor]: 0.04,
			[Stat.StatBonusArmor]: 0.04,
			[Stat.StatStamina]: 1.14,
			[Stat.StatMasteryRating]: 1.0,
			[Stat.StatStrength]: 0.5,
			[Stat.StatAgility]: 0,
			[Stat.StatAttackPower]: 0.15,
			[Stat.StatExpertiseRating]: 0.75,
			[Stat.StatHitRating]: 0.75,
			[Stat.StatCritRating]: 0.2,
			[Stat.StatHasteRating]: 0.3,
			[Stat.StatDodgeRating]: 0.6,
			[Stat.StatParryRating]: 0.6,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.33,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: {
			major3: PaladinMajorGlyph.GlyphOfFocusedShield,
		},
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	classOptions: {
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58085, // Flask of Steelskin
	foodId: 62663, // Lavascale Minestrone
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
	tinkerId: 82174, // Synapse Springs
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Leatherworking,
	iterationCount: 25000,
};
