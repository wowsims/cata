import * as PresetUtils from '../../core/preset_utils';
import { ConsumesSpec, Debuffs, Glyphs, IndividualBuffs, Profession, PseudoStat, RaidBuffs, Stat } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DemonologyWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
} from '../../core/proto/warlock';
import { Stats, UnitStat } from '../../core/proto_utils/stats';
import { WARLOCK_BREAKPOINTS } from '../presets';
import IncinerateAPL from './apls/incinerate.apl.json';
import ShadowBoltAPL from './apls/shadow-bolt.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P3Gear from './gear_sets/p3.gear.json';
import ItemSwapP3 from './gear_sets/p3_item_swap.gear.json';
import P4Gear from './gear_sets/p4.gear.json';
import ItemSwapP4 from './gear_sets/p4_item_swap.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 - BIS', P1Gear);
export const P3_PRESET = PresetUtils.makePresetGear('P3 - BIS', P3Gear);
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);

export const P3_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P3 - Mastery', ItemSwapP3);
export const P4_ITEM_SWAP = PresetUtils.makePresetItemSwapGear('P4 - Mastery', ItemSwapP4);

export const APL_ShadowBolt = PresetUtils.makePresetAPLRotation('Shadow Bolt', ShadowBoltAPL);
export const APL_Incinerate = PresetUtils.makePresetAPLRotation('Incinerate', IncinerateAPL);

// Preset options for EP weights
export const DEFAULT_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.27,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 0.92,
		[Stat.StatCritRating]: 0.51,
		[Stat.StatHasteRating]: 2.75,
		[Stat.StatMasteryRating]: 0.57,
	}),
);

export const Mastery_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Mastery',
	Stats.fromMap({
		[Stat.StatIntellect]: 1.27,
		[Stat.StatSpellPower]: 1.0,
		[Stat.StatHitRating]: 0.92,
		[Stat.StatCritRating]: 0.51,
		[Stat.StatHasteRating]: 2.75,
		[Stat.StatMasteryRating]: 0.76,
	}),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wotlk.wowhead.com/talent-calc and copy the numbers in the url.

export const DemonologyTalentsShadowBolt = {
	name: 'Shadow bolt',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major2: MajorGlyph.GlyphOfLifeTap,
			minor3: MinorGlyph.GlyphOfUnendingBreath,
		}),
	}),
};
export const DemonologyTalentsIncinerate = {
	name: 'Incinerate',
	data: SavedTalents.create({
		talentsString: '',
		glyphs: Glyphs.create({
			major1: MajorGlyph.GlyphOfSoulstone,
			major2: MajorGlyph.GlyphOfLifeTap,
			minor3: MinorGlyph.GlyphOfUnendingBreath,
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Felguard,
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

export const PRESET_BUILD_SHADOWBOLT = PresetUtils.makePresetBuild('Shadow Bolt', {
	talents: DemonologyTalentsShadowBolt,
	rotation: APL_ShadowBolt,
});

export const PRESET_BUILD_INCINERATE = PresetUtils.makePresetBuild('Incinerate', {
	talents: DemonologyTalentsIncinerate,
	rotation: APL_Incinerate,
});

export const DEMONOLOGY_BREAKPOINTS = [
	{
		unitStat: UnitStat.fromPseudoStat(PseudoStat.PseudoStatSpellHastePercent),
		presets: new Map([
			...[...WARLOCK_BREAKPOINTS.find(entry => entry.unitStat.equalsPseudoStat(PseudoStat.PseudoStatSpellHastePercent))!.presets!],
			['16-tick - Immo Aura', 3.35918],
			['8-tick - Immolate:Inferno', 7.12373],
			['17-tick - Immo Aura', 9.95053],
			['18-tick - Immo Aura', 16.61809],
			['9-tick - Immolate:Inferno', 21.43291],
			['19-tick - Immo Aura', 23.38064],
			['20-tick - Immo Aura', 29.95453],
			['10-tick - Immolate:Inferno', 35.71591],
			['21-tick - Immo Aura', 36.70542],
			['22-tick - Immo Aura', 43.3692],
			['11-tick - Immolate:Inferno', 50.03752],
			['23-tick - Immo Aura', 50.03753],
			['24-tick - Immo Aura', 56.6171],
			['25-tick - Immo Aura', 63.26533],
			['12-tick - Immolate:Inferno', 64.24857],
			['26-tick - Immo Aura', 69.92356],
			['27-tick - Immo Aura', 76.52254],
			['13-tick - Immolate:Inferno', 78.6246],
			['28-tick - Immo Aura', 83.31809],
			['29-tick - Immo Aura', 89.93356],
			['14-tick - Immolate:Inferno', 92.86404],
			['30-tick - Immo Aura', 96.65687],
			['31-tick - Immo Aura', 103.45884],
			['15-tick - Immolate:Inferno', 107.11082],
			['32-tick - Immo Aura', 109.86363],
			// ['33-tick - Immo Aura', 116.68477],
			// ['16-tick - Immolate:Inferno', 121.48396],
			// ['34-tick - Immo Aura', 123.46374],
			// ['35-tick - Immo Aura', 130.14965],
			// ['17-tick - Immolate:Inferno', 135.7564],
			// ['36-tick - Immo Aura', 136.68645],
			// ['37-tick - Immo Aura', 143.60542],
			// ['38-tick - Immo Aura', 149.68795],
			// ['18-tick - Immolate:Inferno', 149.8959],
			// ['39-tick - Immo Aura', 156.73948],
			// ['40-tick - Immo Aura', 163.50468],
			// ['19-tick - Immolate:Inferno', 164.20082],
			// ['41-tick - Immo Aura', 169.90561],
			// ['42-tick - Immo Aura', 176.62525],
			// ['20-tick - Immolate:Inferno', 178.68094],
			// ['43-tick - Immo Aura', 183.68802],
			// ['21-tick - Immolate:Inferno', 192.8258],
			// ['22-tick - Immolate:Inferno', 207.21969],
			// ['23-tick - Immolate:Inferno', 221.37122],
			// ['24-tick - Immolate:Inferno', 235.75829],
		]),
	},
];
