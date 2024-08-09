import * as PresetUtils from '../../core/preset_utils';
import { Consumes, Debuffs, Flask, Food, Glyphs, IndividualBuffs, Potions, Profession, RaidBuffs, Stat, TinkerHands } from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DestructionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
	WarlockPrimeGlyph as PrimeGlyph,
} from '../../core/proto/warlock';
import { Stats } from '../../core/proto_utils/stats';
import { WARLOCK_BREAKPOINTS } from '../presets';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P4WrathGear from './gear_sets/p4_wrath.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P4_WOTLK_PRESET = PresetUtils.makePresetGear('P4 Wrath', P4WrathGear, { tooltip: BIS_TOOLTIP });

export const APL_Default = PresetUtils.makePresetAPLRotation('Destro', DefaultApl);

// Preset options for EP weights
export const P1_EP_PRESET = PresetUtils.makePresetEpWeights(
	'P1',
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
		talentsString: '003-03202-3320202312201312211',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfConflagrate,
			prime2: PrimeGlyph.GlyphOfImmolate,
			prime3: PrimeGlyph.GlyphOfImp,
			major1: MajorGlyph.GlyphOfLifeTap,
			major2: MajorGlyph.GlyphOfSoulLink,
			major3: MajorGlyph.GlyphOfHealthstone,
			minor1: MinorGlyph.GlyphOfDrainSoul,
			minor2: MinorGlyph.GlyphOfRitualOfSouls,
			minor3: MinorGlyph.GlyphOfUnendingBreath,
		}),
	}),
};

export const DefaultOptions = WarlockOptions.create({
	classOptions: {
		summon: Summon.Imp,
		detonateSeed: false,
	},
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
	explosiveBigDaddy: false,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	markOfTheWild: true,
	icyTalons: true,
	moonkinForm: true,
	leaderOfThePack: true,
	powerWordFortitude: true,
	strengthOfEarthTotem: true,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: true,
	blessingOfKings: true,
	blessingOfMight: true,
	communion: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	vampiricTouch: true,
	darkIntent: true,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: false,
	demoralizingShout: true,
	frostFever: true,
});

export const OtherDefaults = {
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
	duration: 180,
	durationVariation: 30,
	darkIntentUptime: 90,
};

export const DESTRUCTION_BREAKPOINTS = WARLOCK_BREAKPOINTS;
