import * as Mechanics from '../../core/constants/mechanics';
import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Debuffs, Flask, Food, Glyphs, Potions, Profession, PseudoStat, RaidBuffs, Stat, TinkerHands } from '../../core/proto/common.js';
import {
	AirTotem,
	CallTotem,
	EarthTotem,
	EnhancementShaman_Options as EnhancementShamanOptions,
	FireTotem,
	ShamanImbue,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanPrimeGlyph,
	ShamanShield,
	ShamanSyncType,
	ShamanTotems,
	TotemSet,
	WaterTotem,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import { Stats } from '../../core/proto_utils/stats';
import DefaultApl from './apls/default.apl.json';
import P1NonOrcGear from './gear_sets/p1.non-orc.gear.json';
import P1OrcGear from './gear_sets/p1.orc.gear.json';
import P3NonOrcGear from './gear_sets/p3.non-orc.gear.json';
import P3OrcGear from './gear_sets/p3.orc.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-raid', PreraidGear);

export const P1_ORC_PRESET = PresetUtils.makePresetGear('P1 - Orc', P1OrcGear);
export const P1_NON_ORC_PRESET = PresetUtils.makePresetGear('P1 - Non-Orc', P1NonOrcGear);

export const P3_ORC_PRESET = PresetUtils.makePresetGear('P3 - Orc', P3OrcGear);
export const P3_NON_ORC_PRESET = PresetUtils.makePresetGear('P3 - Non-Orc', P3NonOrcGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Preset options for EP weights
export const P3_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Default',
	Stats.fromMap(
		{
			[Stat.StatIntellect]: 0.07,
			[Stat.StatAgility]: 2.47,
			[Stat.StatSpellPower]: 0,
			[Stat.StatHitRating]: 0.89 + 0.6,
			[Stat.StatCritRating]: 0.26 + 0.58,
			[Stat.StatHasteRating]: 0.22 + 0.44,
			[Stat.StatAttackPower]: 1.0,
			[Stat.StatExpertiseRating]: 1.3,
			[Stat.StatMasteryRating]: 1.21,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 3.05,
			[PseudoStat.PseudoStatOffHandDps]: 2.56,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.89 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.6 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '3022003-2333310013003012321',
		glyphs: Glyphs.create({
			prime1: ShamanPrimeGlyph.GlyphOfLavaLash,
			prime2: ShamanPrimeGlyph.GlyphOfStormstrike,
			prime3: ShamanPrimeGlyph.GlyphOfFeralSpirit,
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfHealingStreamTotem,
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
		call: CallTotem.Elements,
		totems: ShamanTotems.create({
			elements: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			ancestors: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			spirits: TotemSet.create({
				earth: EarthTotem.StrengthOfEarthTotem,
				air: AirTotem.WrathOfAirTotem,
				fire: FireTotem.SearingTotem,
				water: WaterTotem.ManaSpringTotem,
			}),
			earth: EarthTotem.StrengthOfEarthTotem,
			air: AirTotem.WrathOfAirTotem,
			fire: FireTotem.SearingTotem,
			water: WaterTotem.ManaSpringTotem,
		}),
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
});

export const OtherDefaults = {
	distanceFromTarget: 5,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfTheTolvir,
	prepopPotion: Potions.PotionOfTheTolvir,
	flask: Flask.FlaskOfTheWinds,
	food: Food.FoodSeafoodFeast,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
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
	faerieFire: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
});
