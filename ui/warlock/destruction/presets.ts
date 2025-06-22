import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, PseudoStat, RaidBuffs, Stat } from '../../core/proto/common';
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
import P1PreBisGear from './gear_sets/p1-prebis.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const P1_PREBIS_PRESET = PresetUtils.makePresetGear('P1 - Pre-BIS', P1PreBisGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1Gear);
export const DEFAULT_APL = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.24,
		[Stat.StatSpellPower]: 1,
		[Stat.StatHitRating]: 0.83,
		[Stat.StatCritRating]: 0.59,
		[Stat.StatHasteRating]: 0.57,
		[Stat.StatMasteryRating]: 0.61,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DestructionTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		talentsString: '221211',
		glyphs: Glyphs.create({}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Imp,
		detonateSeed: false,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 76085, // Flask of the Warm Sun
	foodId: 74650, // Mogu Fish Stew
	potId: 76093, //Potion of the Jade Serpent
	prepotId: 76093, // Potion of the Jade Serpent
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
	weakenedArmor: true,
	physicalVulnerability: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};

export const DESTRUCTION_BREAKPOINTS = WARLOCK_BREAKPOINTS;
