import { Player } from '../../core/player';
import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Glyphs, HandType, ItemSlot, Profession, PseudoStat, Spec, Stat } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { FuryWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';
import SMFFuryApl from './apls/smf.apl.json';
import TGFuryApl from './apls/tg.apl.json';
import P1FurySMFGear from './gear_sets/p1_fury_smf.gear.json';
import P1FuryTGGear from './gear_sets/p1_fury_tg.gear.json';
import P3FurySMFGear from './gear_sets/p3_fury_smf.gear.json';
import ItemSwapP3SMFGear from './gear_sets/p3_fury_smf_item_swap.gear.json';
import P3FuryTGGear from './gear_sets/p3_fury_tg.gear.json';
import ItemSwapP3TGGear from './gear_sets/p3_fury_tg_item_swap.gear.json';
import P4FurySMFGear from './gear_sets/p4_fury_smf.gear.json';
import ItemSwapP4SMFGear from './gear_sets/p4_fury_smf_item_swap.gear.json';
import P4FuryTGGear from './gear_sets/p4_fury_tg.gear.json';
import ItemSwapP4TGGear from './gear_sets/p4_fury_tg_item_swap.gear.json';
import PreraidFurySMFGear from './gear_sets/preraid_fury_smf.gear.json';
import PreraidFuryTGGear from './gear_sets/preraid_fury_tg.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

// Handlers for spec specific load checks
const FURY_SMF_PRESET_OPTIONS = {
	onLoad: (player: Player<Spec.SpecFuryWarrior>) => {
		PresetUtils.makeSpecChangeWarningToast(
			[
				{
					condition: (player: Player<Spec.SpecFuryWarrior>) =>
						player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeTwoHand,
					message: 'Check your gear: You have a two-handed weapon equipped, but the selected option is for one-handed weapons.',
				},
				// {
				// 	condition: (player: Player<Spec.SpecFuryWarrior>) =>
				// 		player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeTwoHand || !player.getTalents().singleMindedFury,
				// 	message: "Check your talents: You have selected a two-handed spec but don't have [Single-Minded Fury] talented.",
				// },
			],
			player,
		);
	},
};
const FURY_TG_PRESET_OPTIONS = {
	onLoad: (player: Player<any>) => {
		PresetUtils.makeSpecChangeWarningToast(
			[
				{
					condition: (player: Player<Spec.SpecFuryWarrior>) =>
						player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeOneHand,
					message: 'Check your gear: You have a one-handed weapon equipped, but the selected option is for two-handed weapons.',
				},
				// {
				// 	condition: (player: Player<Spec.SpecFuryWarrior>) => !player.getTalents().titansGrip,
				// 	message: "Check your talents: You have selected a one-handed spec but don't have [Titan's Grip] talented.",
				// },
			],
			player,
		);
	},
};

export const P3_PRERAID_FURY_SMF_PRESET = PresetUtils.makePresetGear('Preraid - SMF', PreraidFurySMFGear, FURY_SMF_PRESET_OPTIONS);
export const P3_PRERAID_FURY_TG_PRESET = PresetUtils.makePresetGear('Preraid - TG', PreraidFuryTGGear, FURY_TG_PRESET_OPTIONS);
export const P1_BIS_FURY_SMF_PRESET = PresetUtils.makePresetGear('P1 - SMF', P1FurySMFGear, FURY_SMF_PRESET_OPTIONS);
export const P1_BIS_FURY_TG_PRESET = PresetUtils.makePresetGear('P1 - TG', P1FuryTGGear, FURY_TG_PRESET_OPTIONS);
export const P3_BIS_FURY_SMF_PRESET = PresetUtils.makePresetGear('P3 - SMF', P3FurySMFGear, FURY_SMF_PRESET_OPTIONS);
export const P3_BIS_FURY_TG_PRESET = PresetUtils.makePresetGear('P3 - TG', P3FuryTGGear, FURY_TG_PRESET_OPTIONS);
export const P4_BIS_FURY_TG_PRESET = PresetUtils.makePresetGear('P4 - TG', P4FuryTGGear, FURY_TG_PRESET_OPTIONS);
export const P4_BIS_FURY_SMF_PRESET = PresetUtils.makePresetGear('P4 - SMF', P4FurySMFGear, FURY_SMF_PRESET_OPTIONS);

export const P3_ITEM_SWAP_TG = PresetUtils.makePresetItemSwapGear('P3 - TG', ItemSwapP3TGGear);
export const P3_ITEM_SWAP_SMF = PresetUtils.makePresetItemSwapGear('P3 - SMF', ItemSwapP3SMFGear);
export const P4_ITEM_SWAP_TG = PresetUtils.makePresetItemSwapGear('P4 - TG', ItemSwapP4TGGear);
export const P4_ITEM_SWAP_SMF = PresetUtils.makePresetItemSwapGear('P4 - SMF', ItemSwapP4SMFGear);

export const FURY_SMF_ROTATION = PresetUtils.makePresetAPLRotation('SMF', SMFFuryApl, FURY_SMF_PRESET_OPTIONS);
export const FURY_TG_ROTATION = PresetUtils.makePresetAPLRotation('TG', TGFuryApl, FURY_TG_PRESET_OPTIONS);

// Preset options for EP weights
export const P1_FURY_SMF_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 - SMF',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.2,
			[Stat.StatAgility]: 1.14,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.46,
			[Stat.StatHitRating]: 2.35,
			[Stat.StatCritRating]: 1.48,
			[Stat.StatHasteRating]: 1.05,
			[Stat.StatMasteryRating]: 0.95,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.15,
			[PseudoStat.PseudoStatOffHandDps]: 1.63,
		},
	),
	FURY_SMF_PRESET_OPTIONS,
);

export const P1_FURY_TG_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 - TG',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.21,
			[Stat.StatAgility]: 1.23,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.59,
			[Stat.StatHitRating]: 2.56,
			[Stat.StatCritRating]: 1.59,
			[Stat.StatHasteRating]: 1.15,
			[Stat.StatMasteryRating]: 1.31,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.77,
			[PseudoStat.PseudoStatOffHandDps]: 1.6,
		},
	),
	FURY_TG_PRESET_OPTIONS,
);

export const P3_FURY_SMF_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3 - SMF',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.2,
			[Stat.StatAgility]: 1.07,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.57,
			[Stat.StatHitRating]: 2.46,
			[Stat.StatCritRating]: 1.39,
			[Stat.StatHasteRating]: 1.19,
			[Stat.StatMasteryRating]: 1.02,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.49,
			[PseudoStat.PseudoStatOffHandDps]: 1.71,
		},
	),
	FURY_SMF_PRESET_OPTIONS,
);

export const P3_FURY_TG_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P3 - TG',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.21,
			[Stat.StatAgility]: 1.14,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertiseRating]: 1.71,
			[Stat.StatHitRating]: 2.66,
			[Stat.StatCritRating]: 1.48,
			[Stat.StatHasteRating]: 1.19,
			[Stat.StatMasteryRating]: 1.17,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.77,
			[PseudoStat.PseudoStatOffHandDps]: 1.6,
		},
	),
	FURY_TG_PRESET_OPTIONS,
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const FurySMFTalents = {
	name: 'SMF',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
		}),
	}),
	...FURY_SMF_PRESET_OPTIONS,
};

export const FuryTGTalents = {
	name: 'TG',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
		}),
	}),
	...FURY_TG_PRESET_OPTIONS,
};

export const DefaultOptions = WarriorOptions.create({
	classOptions: {
		startingRage: 0,
	},
	syncType: 0,
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58088, // Flask of Titanic Strength
	foodId: 62670, // Beer-Basted Crocolisk
	potId: 58146, // Golemblood Potion
	prepotId: 58146, // Golemblood Potion
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};

export const P1_PRESET_BUILD_SMF = PresetUtils.makePresetBuild('P1 - SMF', {
	gear: P1_BIS_FURY_SMF_PRESET,
	talents: FurySMFTalents,
	rotation: FURY_SMF_ROTATION,
	epWeights: P1_FURY_SMF_EP_PRESET,
});

export const P1_PRESET_BUILD_TG = PresetUtils.makePresetBuild('P1 - TG', {
	gear: P1_BIS_FURY_TG_PRESET,
	talents: FuryTGTalents,
	rotation: FURY_TG_ROTATION,
	epWeights: P1_FURY_TG_EP_PRESET,
});

export const P3_PRESET_BUILD_SMF = PresetUtils.makePresetBuild('P3 - SMF', {
	gear: P3_BIS_FURY_SMF_PRESET,
	talents: FurySMFTalents,
	rotation: FURY_SMF_ROTATION,
	epWeights: P3_FURY_SMF_EP_PRESET,
});

export const P3_PRESET_BUILD_TG = PresetUtils.makePresetBuild('P3 - TG', {
	gear: P3_BIS_FURY_TG_PRESET,
	talents: FuryTGTalents,
	rotation: FURY_TG_ROTATION,
	epWeights: P3_FURY_TG_EP_PRESET,
});

export const P4_PRESET_BUILD_TG = PresetUtils.makePresetBuild('P4 - TG', {
	gear: P4_BIS_FURY_TG_PRESET,
	talents: FuryTGTalents,
	rotation: FURY_TG_ROTATION,
	epWeights: P3_FURY_TG_EP_PRESET,
});

export const P4_PRESET_BUILD_SMF = PresetUtils.makePresetBuild('P4 - SMF', {
	gear: P4_BIS_FURY_SMF_PRESET,
	talents: FurySMFTalents,
	rotation: FURY_SMF_ROTATION,
	epWeights: P3_FURY_SMF_EP_PRESET,
});
