import * as Mechanics from '../../core/constants/mechanics';
import { Player } from '../../core/player';
import * as PresetUtils from '../../core/preset_utils';
import { makeSpecChangeWarningToast } from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl';
import { ConsumesSpec, Glyphs, HandType, ItemSlot, Profession, PseudoStat, Spec, Stat } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import TwoHandAPL from '../../death_knight/frost/apls/2h.apl.json';
import DualWieldAPL from '../../death_knight/frost/apls/dw.apl.json';
import MasterFrostAPL from '../../death_knight/frost/apls/masterfrost.apl.json';
import P12HGear from '../../death_knight/frost/gear_sets/p1.2h.gear.json';
import P1DWGear from '../../death_knight/frost/gear_sets/p1.dw.gear.json';
import P1MasterfrostGear from '../../death_knight/frost/gear_sets/p1.masterfrost.gear.json';
import P3DWGear from '../../death_knight/frost/gear_sets/p3.dw.gear.json';
import P3MasterfrostGear from '../../death_knight/frost/gear_sets/p3.masterfrost.gear.json';
import P4MasterfrostGear from '../../death_knight/frost/gear_sets/p4.masterfrost.gear.json';
import PreBISGear from '../../death_knight/frost/gear_sets/prebis.gear.json';

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
				// {
				// 	condition: (player: Player<Spec.SpecFrostDeathKnight>) => !player.getTalents().threatOfThassarian,
				// 	message: "Check your talents: You have selected a dual-wield spec but don't have [Threat Of Thassarian] talented.",
				// },
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
				// {
				// 	condition: (player: Player<Spec.SpecFrostDeathKnight>) => !player.getTalents().mightOfTheFrozenWastes,
				// 	message: "Check your talents: You have selected a two-handed spec but don't have [Might of the Frozen Wastes] talented",
				// },
			],
			player,
		);
	},
};

export const P1_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 DW Obliterate', P1DWGear, DW_PRESET_OPTIONS);
export const P3_DW_GEAR_PRESET = PresetUtils.makePresetGear('P3 DW Obliterate', P3DWGear, DW_PRESET_OPTIONS);
export const P1_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 Two Hand', P12HGear, TWOHAND_PRESET_OPTIONS);
export const P1_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P1 Masterfrost', P1MasterfrostGear, DW_PRESET_OPTIONS);
export const P3_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P3 Masterfrost', P3MasterfrostGear, DW_PRESET_OPTIONS);
export const P4_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P4 Masterfrost', P4MasterfrostGear, DW_PRESET_OPTIONS);
export const PREBIS_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('Pre-bis Masterfrost', PreBISGear, DW_PRESET_OPTIONS);

export const DUAL_WIELD_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('DW Obliterate', DualWieldAPL, DW_PRESET_OPTIONS);
export const TWO_HAND_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Two Hand', TwoHandAPL, TWOHAND_PRESET_OPTIONS);
export const MASTERFROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Masterfrost', MasterFrostAPL, DW_PRESET_OPTIONS);

// Preset options for EP weights
export const P1_DUAL_WIELD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 DW Obliterate',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.92,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.56,
			[Stat.StatHasteRating]: 1.3,
			[Stat.StatHitRating]: 1.22,
			[Stat.StatCritRating]: 1.06,
			[Stat.StatMasteryRating]: 1.11,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.05,
			[PseudoStat.PseudoStatOffHandDps]: 3.85,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 146.53,
			[PseudoStat.PseudoStatSpellHitPercent]: 41.91,
		},
	),
	DW_PRESET_OPTIONS,
);

export const P3_DUAL_WIELD_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3 DW Obliterate',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.98,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.65,
			[Stat.StatHasteRating]: 1.8,
			[Stat.StatHitRating]: 1.29,
			[Stat.StatCritRating]: 1.24,
			[Stat.StatMasteryRating]: 1.23,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 6.23,
			[PseudoStat.PseudoStatOffHandDps]: 3.93,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 154.42,
			[PseudoStat.PseudoStatSpellHitPercent]: 41.91,
		},
	),
	DW_PRESET_OPTIONS,
);

export const P1_TWOHAND_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Two Hand',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.98,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.34,
			[Stat.StatHasteRating]: 1.94,
			[Stat.StatHitRating]: 1.61,
			[Stat.StatCritRating]: 1.24,
			[Stat.StatMasteryRating]: 1.26,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 10.08,
			[PseudoStat.PseudoStatOffHandDps]: 0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 1.61 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0,
		},
	),
	TWOHAND_PRESET_OPTIONS,
);

export const P1_MASTERFROST_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Masterfrost',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.86,
			[Stat.StatArmor]: 0.03,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 0.75,
			[Stat.StatHasteRating]: 1.38,
			[Stat.StatHitRating]: 1.67 + 1.4,
			[Stat.StatCritRating]: 0.64 + 0.43,
			[Stat.StatMasteryRating]: 1.41,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 4.5,
			[PseudoStat.PseudoStatOffHandDps]: 2.84,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 1.67 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 1.4 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
	DW_PRESET_OPTIONS,
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DualWieldTalents = {
	name: 'DW Obliterate',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfDeathGrip,
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
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
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
		talentsString: '',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfPestilence,
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
	profession2: Profession.Jewelcrafting,
	distanceFromTarget: 5,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58088, // Flask of Titanic Strength
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
});
export const PRESET_BUILD_DW = PresetUtils.makePresetBuild('P3 - DW Obliterate', {
	gear: P3_DW_GEAR_PRESET,
	talents: DualWieldTalents,
	rotationType: APLRotationType.TypeAuto,
	epWeights: P3_DUAL_WIELD_EP_PRESET,
});

export const PRESET_BUILD_2H = PresetUtils.makePresetBuild('P3 - Two Hand', {
	gear: P1_2H_GEAR_PRESET,
	talents: TwoHandTalents,
	rotationType: APLRotationType.TypeAuto,
	epWeights: P1_TWOHAND_EP_PRESET,
});

export const PRESET_BUILD_MASTERFROST = PresetUtils.makePresetBuild('P4 - Masterfrost', {
	gear: P4_MASTERFROST_GEAR_PRESET,
	talents: MasterfrostTalents,
	rotationType: APLRotationType.TypeAuto,
	epWeights: P1_MASTERFROST_EP_PRESET,
});

export const PRESET_BUILD_PREBIS = PresetUtils.makePresetBuild('P4 - Pre-bis', {
	gear: PREBIS_MASTERFROST_GEAR_PRESET,
	talents: MasterfrostTalents,
	rotationType: APLRotationType.TypeAuto,
	epWeights: P1_MASTERFROST_EP_PRESET,
});
