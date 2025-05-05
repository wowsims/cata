import * as PresetUtils from '../../core/preset_utils.js';
import { ConsumesSpec, Debuffs, Glyphs, Profession, RaidBuffs, Stat } from '../../core/proto/common.js';
import {
	HolyPaladin_Options as Paladin_Options,
	PaladinAura,
	PaladinMajorGlyph as MajorGlyph,
	PaladinMinorGlyph as MinorGlyph,
	PaladinSeal,
} from '../../core/proto/paladin.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('PreRaid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
// export const P2_PRESET = PresetUtils.makePresetGear('P2 Preset', P2Gear);
// export const P3_PRESET = PresetUtils.makePresetGear('P3 Preset', P3Gear);
// export const P4_PRESET = PresetUtils.makePresetGear('P4 Preset', P4Gear);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.375,
		[Stat.StatSpirit]: 1.125,
		[Stat.StatSpellPower]: 1,
		[Stat.StatCritRating]: 0.75,
		[Stat.StatHasteRating]: 0.85,
		[Stat.StatMasteryRating]: 0.5,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.

export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfDivinePlea,
			major2: MajorGlyph.GlyphOfDivinity,
		}),
	}),
};

export const DefaultOptions = Paladin_Options.create({
	classOptions: {
		aura: PaladinAura.Devotion,
		seal: PaladinSeal.Insight,
	},
});

export const DefaultRaidBuffs = RaidBuffs.create({});
export const DefaultConsumables = ConsumesSpec.create({
	flaskId: 58086, // Flask of the Draconic Mind
	foodId: 62290, // Seafood Magnifique Feast
	potId: 58091, // Volcanic Potion
});
export const DefaultDebuffs = Debuffs.create({
	// bloodFrenzy: true,
	// sunderArmor: true,
	// ebonPlaguebringer: true,
	// mangle: true,
	// criticalMass: true,
	// demoralizingShout: true,
	// frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 40,
	profession1: Profession.Engineering,
	profession2: Profession.Jewelcrafting,
};
