import * as PresetUtils from '../../core/preset_utils';
import {
	Consumes,
	Debuffs,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	Potions,
	Profession,
	RaidBuffs,
	TinkerHands,
} from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	AfflictionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Summon as Summon,
	WarlockPrimeGlyph as PrimeGlyph,
} from '../../core/proto/warlock';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P4WrathGear from './gear_sets/p4_wrath.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear,);
export const P4_WOTLK_PRESET = PresetUtils.makePresetGear('P4 Wrath', P4WrathGear, { tooltip: BIS_TOOLTIP });

export const APL_Default = PresetUtils.makePresetAPLRotation('Affliction', DefaultApl);

export const AfflictionTalents = {
	name: 'Affliction',
	data: SavedTalents.create({
		talentsString: '223222003013321321-03-33',
		glyphs: Glyphs.create({
			prime1: PrimeGlyph.GlyphOfHaunt,
			prime2: PrimeGlyph.GlyphOfUnstableAffliction,
			prime3: PrimeGlyph.GlyphOfCorruption,
			major1: MajorGlyph.GlyphOfShadowBolt,
			major2: MajorGlyph.GlyphOfLifeTap,
			major3: MajorGlyph.GlyphOfSoulSwap,
			minor1: MinorGlyph.GlyphOfDrainSoul,
			minor2: MinorGlyph.GlyphOfRitualOfSouls,
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
	judgement: true,
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
