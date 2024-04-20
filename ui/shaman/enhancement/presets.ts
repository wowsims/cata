import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Debuffs, Faction, Flask, Food, Glyphs, Potions, RaidBuffs, TinkerHands, TristateEffect } from '../../core/proto/common.js';
import {
	AirTotem,
	EarthTotem,
	EnhancementShaman_Options as EnhancementShamanOptions,
	FireTotem,
	ShamanImbue,
	ShamanPrimeGlyph,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
	ShamanSyncType,
	ShamanTotems,
	WaterTotem,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';
import PrepatchGear from './gear_sets/prepatch.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PREPATCH_PRESET = PresetUtils.makePresetGear('Prepatch Preset', PrepatchGear);
export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);


export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '3020023-2333310013003012321',
		glyphs: Glyphs.create({
			prime1: ShamanPrimeGlyph.GlyphOfLavaLash,
			prime2: ShamanPrimeGlyph.GlyphOfFeralSpirit,
			prime3: ShamanPrimeGlyph.GlyphOfWindfuryWeapon,
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfChainLightning,
			major3: ShamanMajorGlyph.GlyphOfFireNova,
			minor1: ShamanMinorGlyph.GlyphOfWaterWalking,
			minor2: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfTheArcticWolf,
		}),
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shield: ShamanShield.LightningShield,
		totems: ShamanTotems.create({
			earth: EarthTotem.StoneskinTotem,
			fire: FireTotem.SearingTotem,
			water: WaterTotem.ManaSpringTotem,
			air: AirTotem.WindfuryTotem,
			useFireElemental: true
		}),
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodGrilledDragon,
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

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
	judgement: true,
});
