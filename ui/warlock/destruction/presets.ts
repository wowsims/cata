import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, RaidBuffs, Stat } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DestructionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
} from '../../core/proto/warlock';
import { Stats } from '../../core/proto_utils/stats';
import { WARLOCK_BREAKPOINTS } from '../presets';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1Gear);
export const DEFAULT_APL = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.25,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 0.87,
		[Stat.StatCritRating]: 0.48,
		[Stat.StatHasteRating]: 0.55,
		[Stat.StatMasteryRating]: 0.47,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DestructionTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '221231',
		glyphs: Glyphs.create({
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Imp,
		detonateSeed: false,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
	tinkerId: 82174, // Synapse Springs
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	blessingOfKings: true,
	leaderOfThePack: true,
	blessingOfMight: true,
	bloodlust: true,
	moonkinAura: true,
	skullBannerCount: 2,
	stormlashTotemCount: 4,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	curseOfElements: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};

export const DESTRUCTION_BREAKPOINTS = WARLOCK_BREAKPOINTS;
