import { Player } from '../../core/player';
import * as PresetUtils from '../../core/preset_utils';
import { makeSpecChangeWarningToast } from '../../core/preset_utils';
import { APLRotation_Type as APLRotationType } from '../../core/proto/apl';
import { ConsumesSpec, Glyphs, HandType, ItemSlot, Profession, PseudoStat, Spec, Stat } from '../../core/proto/common';
import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, FrostDeathKnight_Options } from '../../core/proto/death_knight';
import { SavedTalents } from '../../core/proto/ui';
import { Stats } from '../../core/proto_utils/stats';
import TwoHandAPL from '../../death_knight/frost/apls/2h.apl.json';
// import DualWieldAPL from '../../death_knight/frost/apls/dw.apl.json';
import MasterFrostAPL from '../../death_knight/frost/apls/masterfrost.apl.json';
import P12HGear from '../../death_knight/frost/gear_sets/p1.2h.gear.json';
// import P1DWGear from '../../death_knight/frost/gear_sets/p1.dw.gear.json';
import P1MasterfrostGear from '../../death_knight/frost/gear_sets/p1.masterfrost.gear.json';
// import PreBISGear from '../../death_knight/frost/gear_sets/prebis.gear.json';

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
			],
			player,
		);
	},
};

// export const P1_DW_GEAR_PRESET = PresetUtils.makePresetGear('P1 DW Obliterate', P1DWGear, DW_PRESET_OPTIONS);
export const P1_2H_GEAR_PRESET = PresetUtils.makePresetGear('P1 Two Hand', P12HGear, TWOHAND_PRESET_OPTIONS);
export const P1_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('P1 Masterfrost', P1MasterfrostGear, DW_PRESET_OPTIONS);
// export const PREBIS_MASTERFROST_GEAR_PRESET = PresetUtils.makePresetGear('Pre-bis Masterfrost', PreBISGear, DW_PRESET_OPTIONS);

// export const DUAL_WIELD_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('DW Obliterate', DualWieldAPL, DW_PRESET_OPTIONS);
export const TWO_HAND_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Two Hand', TwoHandAPL, TWOHAND_PRESET_OPTIONS);
export const MASTERFROST_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Masterfrost', MasterFrostAPL, DW_PRESET_OPTIONS);

// // Preset options for EP weights
// export const P1_DUAL_WIELD_EP_PRESET = PresetUtils.makePresetEpWeights(
// 	'P1 DW Obliterate',
// 	Stats.fromMap(
// 		{
// 			[Stat.StatStrength]: 2.92,
// 			[Stat.StatArmor]: 0.03,
// 			[Stat.StatAttackPower]: 1,
// 			[Stat.StatExpertiseRating]: 0.56,
// 			[Stat.StatHasteRating]: 1.3,
// 			[Stat.StatHitRating]: 1.22,
// 			[Stat.StatCritRating]: 1.06,
// 			[Stat.StatMasteryRating]: 1.11,
// 		},
// 		{
// 			[PseudoStat.PseudoStatMainHandDps]: 6.05,
// 			[PseudoStat.PseudoStatOffHandDps]: 3.85,
// 			[PseudoStat.PseudoStatPhysicalHitPercent]: 146.53,
// 			[PseudoStat.PseudoStatSpellHitPercent]: 41.91,
// 		},
// 	),
// 	DW_PRESET_OPTIONS,
// );

export const P1_TWOHAND_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Two Hand',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 0.7,
			[Stat.StatExpertiseRating]: 0.7,
			[Stat.StatHasteRating]: 0.68,
			[Stat.StatMasteryRating]: 0.64,
			[Stat.StatCritRating]: 0.63,
			[Stat.StatAttackPower]: 0.37,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 2.61,
			[PseudoStat.PseudoStatOffHandDps]: 0,
		},
	),
	TWOHAND_PRESET_OPTIONS,
);

export const P1_MASTERFROST_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 Masterfrost',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 1.0,
			[Stat.StatHitRating]: 0.54,
			[Stat.StatExpertiseRating]: 0.54,
			[Stat.StatCritRating]: 0.51,
			[Stat.StatMasteryRating]: 0.51,
			[Stat.StatAttackPower]: 0.38,
			[Stat.StatHasteRating]: 0.37,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 1.06,
			[PseudoStat.PseudoStatOffHandDps]: 0.5,
		},
	),
	DW_PRESET_OPTIONS,
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DefaultTalents = {
	name: 'Default',
	data: SavedTalents.create({
		talentsString: '221111',
		glyphs: Glyphs.create({
			major1: DeathKnightMajorGlyph.GlyphOfAntiMagicShell,
			major2: DeathKnightMajorGlyph.GlyphOfPestilence,
			major3: DeathKnightMajorGlyph.GlyphOfLoudHorn,
			minor1: DeathKnightMinorGlyph.GlyphOfArmyOfTheDead,
			minor2: DeathKnightMinorGlyph.GlyphOfTranquilGrip,
			minor3: DeathKnightMinorGlyph.GlyphOfDeathGate,
		}),
	}),
	...DW_PRESET_OPTIONS,
};

export const DefaultOptions = FrostDeathKnight_Options.create({
	classOptions: {
		startingRunicPower: 0,
	},
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76088, // Flask of Winter's Bite
	foodId: 74646, // Black Pepper Ribs and Shrimp
	potId: 76095, // Potion of Mogu Power
	prepotId: 76095, // Potion of Mogu Power
});

// export const PRESET_BUILD_DW = PresetUtils.makePresetBuild('P1 - DW Obliterate', {
// 	gear: P1_DW_GEAR_PRESET,
// 	talents: DefaultTalents,
// 	rotationType: APLRotationType.TypeAPL,
// 	rotation: DUAL_WIELD_ROTATION_PRESET_DEFAULT,
// 	epWeights: P1_DUAL_WIELD_EP_PRESET,
// });

export const PRESET_BUILD_2H = PresetUtils.makePresetBuild('P1 - Two Hand', {
	gear: P1_2H_GEAR_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
	rotation: TWO_HAND_ROTATION_PRESET_DEFAULT,
	epWeights: P1_TWOHAND_EP_PRESET,
});

export const PRESET_BUILD_MASTERFROST = PresetUtils.makePresetBuild('P1 - Masterfrost', {
	gear: P1_MASTERFROST_GEAR_PRESET,
	talents: DefaultTalents,
	rotationType: APLRotationType.TypeAuto,
	rotation: MASTERFROST_ROTATION_PRESET_DEFAULT,
	epWeights: P1_MASTERFROST_EP_PRESET,
});

// export const PRESET_BUILD_PREBIS = PresetUtils.makePresetBuild('P1 - Pre-bis', {
// 	gear: PREBIS_MASTERFROST_GEAR_PRESET,
// 	talents: DefaultTalents,
// 	rotationType: APLRotationType.TypeAPL,
// 	rotation: MASTERFROST_ROTATION_PRESET_DEFAULT,
// 	epWeights: P1_MASTERFROST_EP_PRESET,
// });
