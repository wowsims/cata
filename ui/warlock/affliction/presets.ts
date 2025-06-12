import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, RaidBuffs, Stat } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	AfflictionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
} from '../../core/proto/warlock';
import { Stats } from '../../core/proto_utils/stats';
import { WARLOCK_BREAKPOINTS } from '../presets';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import ItemSwapP4 from './gear_sets/p4_item_swap.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4 - BIS', P4Gear);

export const P4_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P4 - Item Swap', ItemSwapP4);

export const APL_Default = PresetUtils.makePresetAPLRotation('Affliction', DefaultApl);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.26,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 0.93,
		[Stat.StatCritRating]: 0.52,
		[Stat.StatHasteRating]: 0.58,
		[Stat.StatMasteryRating]: 0.38,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const AfflictionTalents = {
	name: 'Affliction',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			minor3: MinorGlyph.GlyphOfUnendingBreath,
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Felhunter,
		detonateSeed: false,
	},
});

export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
	prepotId: 58091, // Volcanic Potion
});

export const DefaultRaidBuffs = RaidBuffs.create({});

export const DefaultIndividualBuffs = IndividualBuffs.create({});

export const DefaultDebuffs = Debuffs.create({
	// bloodFrenzy: true,
	// sunderArmor: true,
	// ebonPlaguebringer: true,
	// mangle: true,
	// criticalMass: false,
	// demoralizingShout: true,
	// frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};

export const AFFLICTION_BREAKPOINTS = WARLOCK_BREAKPOINTS;
