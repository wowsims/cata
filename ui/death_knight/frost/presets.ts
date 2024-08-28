import * as Mechanics from '../../core/constants/mechanics';
import { Player } from '../../core/player';
import * as PresetUtils from '../../core/preset_utils';
import { makeSpecChangeWarningToast } from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, HandType, ItemSlot, Potions, Profession, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import TwoHandAPL from '../../death_knight/frost/apls/2h.apl.json';
import DualWieldAPL from '../../death_knight/frost/apls/dw.apl.json';
import MasterFrostAPL from '../../death_knight/frost/apls/masterfrost.apl.json';
import P12HGear from '../../death_knight/frost/gear_sets/p1.2h.gear.json';
import P1DWGear from '../../death_knight/frost/gear_sets/p1.dw.gear.json';
import P1MasterfrostGear from '../../death_knight/frost/gear_sets/p1.masterfrost.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Handlers for spec specific load checks
const DW_PRESET_OPTIONS = {
	onLoad: (player: Player<Spec.SpecFrostDeathKnight>) => {
		makeSpecChangeWarningToast(
			[
				{
					condition: (player: Player<Spec.SpecFrostDeathKnight>) =>
						player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeTwoHand,
					message: 'Check your gear: You have a two-handed weapon equipped, but the selected option is for dual wield.',
				},
				{
					condition: (player: Player<Spec.SpecFrostDeathKnight>) => !player.getTalents().threatOfThassarian,
					message: "Check your talents: You have selected a dual-wield spec but don't have [Threat Of Thassarian] talented.",
				},
			],
			player,
		);
	},
};

const TWOHAND_PRESET_OPTIONS = {
	onLoad: (player: Player<any>) => {
		makeSpecChangeWarningToast(
			[
				{
					condition: (player: Player<Spec.SpecFrostDeathKnight>) =>
						player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeOneHand,
					message: 'Check your gear: You have a one-handed weapon equipped, but the selected option is for dual wield',
				},
				{
					condition: (player: Player<Spec.SpecFrostDeathKnight>) => !player.getTalents().mightOfTheFrozenWastes,
					message: "Check your talents: You have selected a two-handed spec but don't have [Might of the Frozen Wastes] talented",
				},
			],
			player,
		);
	},
};

export const P1_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 Dual Wield', P1DWGear, DW_PRESET_OPTIONS);
export const P1_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 Two Hand', P12HGear, TWOHAND_PRESET_OPTIONS);
export const P1_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P1 Masterfrost', P1MasterfrostGear, DW_PRESET_OPTIONS);

export const DUAL_WIELD_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Dual Wield', DualWieldAPL, DW_PRESET_OPTIONS);
export const TWO_HAND_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Two Hand', TwoHandAPL, TWOHAND_PRESET_OPTIONS);
export const MASTERFROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Masterfrost', MasterFrostAPL, DW_PRESET_OPTIONS);

// Preset options for EP weights
export const P1_MASTERFROST_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Masterfrost',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.86,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.75,
			[Stat.StatHasteRating]: 1.38,
			[Stat.StatHitRating]: 1.08 + 0.59,
			[Stat.StatCritRating]: 0.64 + 0.43,
			[Stat.StatMasteryRating]: 1.41,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.5,
			[PseudoStat.PseudoStatOffHandDps]: 2.84,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 1.08 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.59 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
	DW_PRESET_OPTIONS,
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DualWieldTalents = {
	name: 'Dual Wield',
	data: SavedTalents.create({
		talentsString: '2032-20330022233112012301-003',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
	...DW_PRESET_OPTIONS,
};

export const TwoHandTalents = {
	name: 'Two Hand',
	data: SavedTalents.create({
		talentsString: '103-32030022233112012031-033',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
	...TWOHAND_PRESET_OPTIONS,
};

export const MasterfrostTalents = {
	name: 'Masterfrost',
	data: SavedTalents.create({
		talentsString: '2032-20330022233112012301-03',
		glyphs: Glyphs.create({
			prime1: DeathKnightPrimeGlyph.GlyphOfFrostStrike,
			prime2: DeathKnightPrimeGlyph.GlyphOfObliterate,
			prime3: DeathKnightPrimeGlyph.GlyphOfHowlingBlast,
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
			major2: DeathKnightMajorGlyph.GlyphOfBloodBoil,
			major3: DeathKnightMajorGlyph.GlyphOfDarkSuccor,
			minor1: DeathKnightMinorGlyph.GlyphOfDeathGate,
			minor2: DeathKnightMinorGlyph.GlyphOfPathOfFrost,
			minor3: DeathKnightMinorGlyph.GlyphOfHornOfWinter,
		}),
	}),
	...DW_PRESET_OPTIONS,
};

export const DefaultOptions = FrostDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
		petUptime: 1,
	},
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Enchanting,
	distanceFromTarget: 5,
};

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const PRESET_BUILD_DW = PresetUtils.makePresetBuild('Dual Wield', {
	gear: P1_DW_GEAR_PRESET,
	talents: DualWieldTalents,
	rotation: DUAL_WIELD_ROTATION_PRESET_DEFAULT,
});

export const PRESET_BUILD_2H = PresetUtils.makePresetBuild('Two Hand', {
	gear: P1_2H_GEAR_PRESET,
	talents: TwoHandTalents,
	rotation: TWO_HAND_ROTATION_PRESET_DEFAULT,
});

export const PRESET_BUILD_MASTERFROST = PresetUtils.makePresetBuild('Masterfrost', {
	gear: P1_MASTERFROST_GEAR_PRESET,
	talents: MasterfrostTalents,
	rotation: MASTERFROST_ROTATION_PRESET_DEFAULT,
});
