import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Potions, PseudoStat, Stat } from '../../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinSeal,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
//import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('P1 PreRaid Preset', PreraidGear);
//export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatArmor]: 0.07,
			[Stat.StatBonusArmor]: 0.06,
			[Stat.StatStamina]: 1.14,
			[Stat.StatMastery]: 1.0,
			[Stat.StatStrength]: 1.0,
			[Stat.StatAgility]: 0,
			[Stat.StatAttackPower]: 0.26,
			[Stat.StatExpertise]: 0.5,
			[Stat.StatMeleeHit]: 0.5,
			[Stat.StatMeleeCrit]: 0.3,
			[Stat.StatMeleeHaste]: 0.17,
			[Stat.StatSpellPower]: 0,
			[Stat.StatDodge]: 0.6,
			[Stat.StatParry]: 0.6,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.33,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const GenericAoeTalents = {
	name: 'Baseline Example',
	data: SavedTalents.create({
		// talentsString: '-05005135200132311333312321-511302012003',
		// glyphs: {
		// 	major1: PaladinMajorGlyph.GlyphOfSealOfVengeance,
		// 	major2: PaladinMajorGlyph.GlyphOfRighteousDefense,
		// 	major3: PaladinMajorGlyph.GlyphOfDivinePlea,
		// 	minor1: PaladinMinorGlyph.GlyphOfSenseUndead,
		// 	minor2: PaladinMinorGlyph.GlyphOfLayOnHands,
		// 	minor3: PaladinMinorGlyph.GlyphOfBlessingOfKings,
		// },
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	classOptions: {
		aura: PaladinAura.Retribution,
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfStoneblood,
	food: Food.FoodDragonfinFilet,
	defaultPotion: Potions.EarthenPotion,
	prepopPotion: Potions.EarthenPotion,
});
