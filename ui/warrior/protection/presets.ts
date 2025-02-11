import * as PresetUtils from '../../core/preset_utils.js';
import {
	BattleElixir,
	Consumes,
	Food,
	Glyphs,
	GuardianElixir,
	Potions,
	Profession,
	PseudoStat,
	Stat,
	TinkerHands,
} from '../../core/proto/common.js';
import { SavedTalents } from '../../core/proto/ui.js';
import {
	ProtectionWarrior_Options as ProtectionWarriorOptions,
	WarriorMajorGlyph,
	WarriorMinorGlyph,
	WarriorPrimeGlyph,
} from '../../core/proto/warrior.js';
import { Stats } from '../../core/proto_utils/stats';
import ItemSwapP4Gear from '../arms/gear_sets/p4_arms_item_swap.gear.json';
import DefaultApl from './apls/default.apl.json';
import P1BISGear from './gear_sets/p1_bis.gear.json';
import P3BISGear from './gear_sets/p3_bis.gear.json';
import P4BISGear from './gear_sets/p4_bis.gear.json';
import PreraidBISGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_BALANCED_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidBISGear);
export const P1_BALANCED_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1BISGear);
export const P3_BALANCED_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3BISGear);
export const P4_BALANCED_PRESET = PresetUtils.makePresetGear('P4 - BIS', P4BISGear);

export const P4_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P4', ItemSwapP4Gear);

export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('Default APL', DefaultApl);

// Preset options for EP weights
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatArmor]: 2.155,
			[Stat.StatBonusArmor]: 2.155,
			[Stat.StatStamina]: 12.442,
			[Stat.StatStrength]: 1.4,
			[Stat.StatAgility]: 0.26,
			[Stat.StatAttackPower]: 0.196,
			[Stat.StatExpertiseRating]: 0.863,
			[Stat.StatHitRating]: 0.736,
			[Stat.StatCritRating]: 0.336,
			[Stat.StatHasteRating]: 0.048,
			[Stat.StatDodgeRating]: 4.801,
			[Stat.StatParryRating]: 4.801,
			[Stat.StatMasteryRating]: 7.415,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.081,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '320003-002-33213201121210212031',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfRevenge,
			prime2: WarriorPrimeGlyph.GlyphOfShieldSlam,
			prime3: WarriorPrimeGlyph.GlyphOfDevastate,
			major1: WarriorMajorGlyph.GlyphOfShieldWall,
			major2: WarriorMajorGlyph.GlyphOfShockwave,
			major3: WarriorMajorGlyph.GlyphOfThunderClap,
			minor1: WarriorMinorGlyph.GlyphOfBattle,
			minor2: WarriorMinorGlyph.GlyphOfBerserkerRage,
			minor3: WarriorMinorGlyph.GlyphOfDemoralizingShout,
		}),
	}),
};

export const DefaultOptions = ProtectionWarriorOptions.create({
	classOptions: {
		startingRage: 0,
	},
});

export const DefaultConsumes = Consumes.create({
	// flask: Flask.FlaskOfSteelskin,
	battleElixir: BattleElixir.ElixirOfTheMaster,
	guardianElixir: GuardianElixir.ElixirOfDeepEarth,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.EarthenPotion,
	prepopPotion: Potions.EarthenPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
	explosiveBigDaddy: true,
});

export const OtherDefaults = {
	profession1: Profession.Leatherworking,
	profession2: Profession.Inscription,
};
