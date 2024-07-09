import { Player } from '../../core/player';
import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Flask, Food, Glyphs, HandType, ItemSlot, Potions, Profession, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import { FuryWarrior_Options as WarriorOptions, WarriorMajorGlyph, WarriorMinorGlyph, WarriorPrimeGlyph } from '../../core/proto/warrior';
import { Stats } from '../../core/proto_utils/stats';
import FuryApl from './apls/fury.apl.json';
import P1FurySMFGear from './gear_sets/p1_fury_smf.gear.json';
import P1FuryTGGear from './gear_sets/p1_fury_tg.gear.json';
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
				{
					condition: (player: Player<Spec.SpecFuryWarrior>) =>
						player.getEquippedItem(ItemSlot.ItemSlotMainHand)?.item.handType === HandType.HandTypeTwoHand || !player.getTalents().singleMindedFury,
					message: "Check your talents: You have selected a two-handed spec but don't have [Single-Minded Fury] talented.",
				},
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
				{
					condition: (player: Player<Spec.SpecFuryWarrior>) => !player.getTalents().titansGrip,
					message: "Check your talents: You have selected a one-handed spec but don't have [Titan's Grip] talented.",
				},
			],
			player,
		);
	},
};

export const PRERAID_FURY_SMF_PRESET = PresetUtils.makePresetGear('Preraid Fury SMF', PreraidFurySMFGear, FURY_SMF_PRESET_OPTIONS);
export const PRERAID_FURY_TG_PRESET = PresetUtils.makePresetGear('Preraid Fury TG', PreraidFuryTGGear, FURY_TG_PRESET_OPTIONS);
export const P1_FURY_SMF_PRESET = PresetUtils.makePresetGear('P1 Fury SMF', P1FurySMFGear, FURY_SMF_PRESET_OPTIONS);
export const P1_FURY_TG_PRESET = PresetUtils.makePresetGear('P1 Fury TG', P1FuryTGGear, FURY_TG_PRESET_OPTIONS);

export const ROTATION_FURY = PresetUtils.makePresetAPLRotation('Fury', FuryApl);

// Preset options for EP weights
export const P1_FURY_TG_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 - TG',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.21,
			[Stat.StatAgility]: 1.23,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 1.78,
			[Stat.StatMeleeHit]: 2.61,
			[Stat.StatMeleeCrit]: 1.59,
			[Stat.StatMeleeHaste]: 1.07,
			[Stat.StatMastery]: 1.22,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.34,
			[PseudoStat.PseudoStatOffHandDps]: 1.77,
		},
	),
	FURY_TG_PRESET_OPTIONS,
);

export const P1_FURY_SMF_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1 - SMF',
	Stats.fromMap(
		{
			[Stat.StatStrength]: 2.2,
			[Stat.StatAgility]: 1.14,
			[Stat.StatAttackPower]: 1,
			[Stat.StatExpertise]: 1.61,
			[Stat.StatMeleeHit]: 2.36,
			[Stat.StatMeleeCrit]: 1.47,
			[Stat.StatMeleeHaste]: 0.92,
			[Stat.StatMastery]: 0.99,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.23,
			[PseudoStat.PseudoStatOffHandDps]: 1.7,
		},
	),
	FURY_SMF_PRESET_OPTIONS,
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.

export const FurySMFTalents = {
	name: 'Fury SMF',
	data: SavedTalents.create({
		talentsString: '302203-032222031301101223201',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfBloodthirst,
			prime2: WarriorPrimeGlyph.GlyphOfRagingBlow,
			prime3: WarriorPrimeGlyph.GlyphOfSlam,
			major1: WarriorMajorGlyph.GlyphOfCleaving,
			major2: WarriorMajorGlyph.GlyphOfDeathWish,
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfBattle,
			minor3: WarriorMinorGlyph.GlyphOfBerserkerRage,
		}),
	}),
	...FURY_SMF_PRESET_OPTIONS,
};

export const FuryTGTalents = {
	name: 'Fury TG',
	data: SavedTalents.create({
		talentsString: '302203-03222203130110122321',
		glyphs: Glyphs.create({
			prime1: WarriorPrimeGlyph.GlyphOfBloodthirst,
			prime2: WarriorPrimeGlyph.GlyphOfRagingBlow,
			prime3: WarriorPrimeGlyph.GlyphOfSlam,
			major1: WarriorMajorGlyph.GlyphOfCleaving,
			major2: WarriorMajorGlyph.GlyphOfDeathWish,
			major3: WarriorMajorGlyph.GlyphOfColossusSmash,
			minor1: WarriorMinorGlyph.GlyphOfCommand,
			minor2: WarriorMinorGlyph.GlyphOfBattle,
			minor3: WarriorMinorGlyph.GlyphOfBerserkerRage,
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

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTitanicStrength,
	food: Food.FoodBeerBasedCrocolisk,
	defaultPotion: Potions.GolembloodPotion,
	prepopPotion: Potions.GolembloodPotion,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	profession1: Profession.Engineering,
	profession2: Profession.Blacksmithing,
	distanceFromTarget: 5,
};
