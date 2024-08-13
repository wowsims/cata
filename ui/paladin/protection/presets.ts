import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Flask, Food, Potions, PseudoStat, Stat } from '../../core/proto/common.js';
import {
	PaladinAura as PaladinAura,
	PaladinMajorGlyph,
	PaladinMinorGlyph,
	PaladinPrimeGlyph,
	PaladinSeal,
	ProtectionPaladin_Options as ProtectionPaladinOptions,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
//import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';
import T11Gear from './gear_sets/T11.gear.json';
import T11CTCGear from './gear_sets/T11CTC.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('P1 PreRaid Preset', PreraidGear);
export const T11_PRESET = PresetUtils.makePresetGear('T11 Balanced Preset', T11Gear);
export const T11CTC_PRESET = PresetUtils.makePresetGear('T11 CTC Preset', T11CTCGear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap(
		{
			[Stat.StatArmor]: 0.04,
			[Stat.StatBonusArmor]: 0.04,
			[Stat.StatStamina]: 1.14,
			[Stat.StatMastery]: 1.0,
			[Stat.StatStrength]: 0.5,
			[Stat.StatAgility]: 0,
			[Stat.StatAttackPower]: 0.15,
			[Stat.StatExpertise]: 0.75,
			[Stat.StatMeleeHit]: 0.75,
			[Stat.StatMeleeCrit]: 0.2,
			[Stat.StatMeleeHaste]: 0.3,
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
		"talentsString": "-32023013122121101231-032032",
		"glyphs": {
		  "prime1": PaladinPrimeGlyph.GlyphOfShieldOfTheRighteous,
		  "prime2": PaladinPrimeGlyph.GlyphOfCrusaderStrike,
		  "prime3": PaladinPrimeGlyph.GlyphOfSealOfTruth,
		  "major1": PaladinMajorGlyph.GlyphOfTheAsceticCrusader,
		  "major2": PaladinMajorGlyph.GlyphOfLayOnHands,
		  "major3": PaladinMajorGlyph.GlyphOfHolyWrath,
		  "minor1": PaladinMinorGlyph.GlyphOfTruth,
		  "minor2": PaladinMinorGlyph.GlyphOfBlessingOfMight,
		  "minor3": PaladinMinorGlyph.GlyphOfInsight,
		},
	}),
};

export const DefaultOptions = ProtectionPaladinOptions.create({
	classOptions: {
		aura: PaladinAura.Retribution,
		seal: PaladinSeal.Truth,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	food: Food.FoodLavascaleMinestrone,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
});
