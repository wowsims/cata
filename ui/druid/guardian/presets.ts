import * as Mechanics from '../../core/constants/mechanics.js';
import * as PresetUtils from '../../core/preset_utils.js';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
	DruidPrimeGlyph,
	GuardianDruid_Options as DruidOptions,
	GuardianDruid_Rotation as DruidRotation,
} from '../../core/proto/druid.js';
import { SavedTalents } from '../../core/proto/ui.js';
// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.
import PreraidGear from './gear_sets/preraid.gear.json';
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-Raid BiS', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1/P2 BiS', P1Gear);
import P2Gear from './gear_sets/p2.gear.json';
export const P2_PRESET = PresetUtils.makePresetGear('P2', P2Gear);
import P3Gear from './gear_sets/p3.gear.json';
export const P3_PRESET = PresetUtils.makePresetGear('P3', P3Gear);
import P4Gear from './gear_sets/p4.gear.json';
export const P4_PRESET = PresetUtils.makePresetGear('P4', P4Gear);

export const DefaultSimpleRotation = DruidRotation.create({
	maintainFaerieFire: true,
	maintainDemoralizingRoar: true,
	demoTime: 4.0,
	pulverizeTime: 4.0,
	prepullStampede: true,
});

import { Stats } from '../../core/proto_utils/stats';
import CleaveApl from './apls/cleave.apl.json';
import DefaultApl from './apls/default.apl.json';
import NefApl from './apls/nef.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);
export const ROTATION_CLEAVE = PresetUtils.makePresetAPLRotation('2-Target Cleave', CleaveApl);
export const ROTATION_NEF = PresetUtils.makePresetAPLRotation('AoE (Nef Adds)', NefApl);

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecGuardianDruid, DefaultSimpleRotation);

// Preset options for EP weights
export const SURVIVAL_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Survival',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.04,
			[Stat.StatStamina]: 0.99,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 1.02,
			[Stat.StatBonusArmor]: 0.23,
			[Stat.StatDodgeRating]: 0.97,
			[Stat.StatMasteryRating]: 0.35,
			[Stat.StatStrength]: 0.11,
			[Stat.StatAttackPower]: 0.1,
			[Stat.StatHitRating]: 0.075 + 0.195,
			[Stat.StatExpertiseRating]: 0.15,
			[Stat.StatCritRating]: 0.11,
			[Stat.StatHasteRating]: 0.0,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.075 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.195 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

export const BALANCED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Balanced',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.02,
			[Stat.StatStamina]: 0.76,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.62,
			[Stat.StatBonusArmor]: 0.14,
			[Stat.StatDodgeRating]: 0.59,
			[Stat.StatMasteryRating]: 0.2,
			[Stat.StatStrength]: 0.21,
			[Stat.StatAttackPower]: 0.2,
			[Stat.StatHitRating]: 0.6,
			[Stat.StatExpertiseRating]: 0.93,
			[Stat.StatCritRating]: 0.25,
			[Stat.StatHasteRating]: 0.03,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.23,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.6 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-2300322312310001220311-020331',
		glyphs: Glyphs.create({
			prime1: DruidPrimeGlyph.GlyphOfMangle,
			prime2: DruidPrimeGlyph.GlyphOfLacerate,
			prime3: DruidPrimeGlyph.GlyphOfBerserk,
			major1: DruidMajorGlyph.GlyphOfFrenziedRegeneration,
			major2: DruidMajorGlyph.GlyphOfMaul,
			major3: DruidMajorGlyph.GlyphOfRebirth,
			minor1: DruidMinorGlyph.GlyphOfDash,
			minor2: DruidMinorGlyph.GlyphOfChallengingRoar,
			minor3: DruidMinorGlyph.GlyphOfUnburdenedRebirth,
		}),
	}),
};

export const DefaultOptions = DruidOptions.create({
	startingRage: 15,
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfSteelskin,
	food: Food.FoodSkeweredEel,
	prepopPotion: Potions.PotionOfTheTolvir,
	defaultPotion: Potions.PotionOfTheTolvir,
	defaultConjured: Conjured.ConjuredHealthstone,
	tinkerHands: TinkerHands.TinkerHandsSynapseSprings,
});

export const OtherDefaults = {
	iterationCount: 50000,
};

export const PRESET_BUILD_SINGLE_TARGET = PresetUtils.makePresetBuild('Single Target', {
	rotation: ROTATION_NEF,
	encounter: PresetUtils.makePresetEncounter(
		'Single Target',
		'https://wowsims.github.io/cata/druid/guardian/#eJzlVF9sFEUY35lZrrtzV5iuBbcTc52uf3Jc0mZvD0ztg7utiTkT0YIG6+Omd6WL17vSO23o01k1bYmJhITY1ioibz4Q0xdpYzxArJhIPAgUWhPA/yYqJvJSiCHO7t4dV9OXRn0wftns7e/3/b7ffPPN5PBmCVigHxQAOAhAAYIPIChB0CUT0A4S4AoAHZAARaAlPI0C3Wl7f2pIChBAbwCskGIJqTeYJi9MoksltFhCOETuTCINHZpEkcMAbyJf3ULqD0xDE5yYARiTLxeQers5MuaVvy6qY0wLcLVbcAzgenL1N6QeCGtAcCEhX/+I1I+aNJHnBd9xaQ5VKzyXwmWkHuEuE1WXTeT3s0i92lztI0T+qHR17K9IJG+dQPw9y9+YTB9H6ihzeYHrjtyq7kULSkhv6oY9kM8p+inEK0D4n4aS+KecOoQu4X0UpHKDj9useRhuNeK6HjeMeIw/uq7HDEOPx2KtOv+Jx07DB6QzRUhm3hDVo99BdrIII6eKUD9XhO3Tl6C1cAEmPrkAz0KxBOuWYcX5uvkNvCiqSAsqdyncgANUxFBt0OrwBoxaYzlFpLBDUML0PqxEiYQUtAwE6mXj+oM0gEVpfhasmd9Wzk9eQ0oP3Y0TyuM4ZGAFEsjZ82MQ06gqAaX+NKghKa/EIO7lRCW0tzbF3TDclqOY47mboptRLPoo1mkbJtpGHJoFsidH0mig4jBfw/kOD/sOS68Az6GFNuN7liHxVFDqrS4juciTYCp5W3lnAikyrXN7lH6NVOk7L29Q4jSG71da/D5W7bJeC2J5FpTrvaU57flupCEPL46vxh9fE1fhSqv9tA/3KLvXMcYojdQ28OGKiO+NbuZzkecrxJpjdetcX7S276pd1FN+Y/lE3tzuwp9gcgWCo8i/VM9bB9BDsvtV0M9bDbtOhF/7eecVc2ufF7fNppeeDm/Zu+NbkwW0wwFPKBzcZ5WvZKPVND3lxufmIz5TMlv80kWz89wXbvxiGu+iwe/PwJHgDnvPgD3MjO0swXpUEP2X/5K6X7XWIy8kP1uP/sn2t2c619VPTXjTKrwnd04Bf4BLZtlLtA6Bffy6ND4xMGinU2xXyu7NO9kMe9YZSEWfSmSHWTqb2cOGnXSa5fu5wHaSLG+/kGL5LCsXJVJ2kjkZ9kyqN5tJ5tpYpLMvnxry9E7GyTt2msX03NYWf81Gy5iDgxOjaGRLV9oecbj9Y9lMLj/0Ym/eP63n/uunJYxHd/690xLGb1pT5RFcr5yWYOngT4i1Lxk=',
		{
			includeHealingModel: true,
		},
	),
});

export const PRESET_BUILD_AOE = PresetUtils.makePresetBuild('AoE (Nef Adds)', {
	rotation: ROTATION_NEF,
	encounter: PresetUtils.makePresetEncounter(
		'AoE (Nef Adds)',
		'https://wowsims.github.io/cata/druid/guardian/#eJztmG1sFEUYx3dmtte9h2K3k0L3xkiHFchxps3eHYVSxVswGGKMQQIhRGM4ettQc72rd1e0EJJaQipIwJAQSgUJ8MVojNr4UipSKATUQFJJyksTRRujGF6+GJPCB3F2rz1O0B5gNF68ySa3+995/s8z83t2L1moUBBHJlqN2hB6A6E2jD7AaACjBUUqOo9QHVYRldhJtYu4FkfDrVZCcamIDSOgat8A0a5x3X2ikwwOkPMDBIrVo1tk7w4EperGNln7ketkUyfx7kYA6uUhot2o9B5wIrfIWgfXXds7iTi8WxFMVN+/QLT9U3UkeTsQqOrHIv6QR5fFfcmeUKpeOEh019Z0RIft0naOaHuFy6aMS6n6rnD5tlIfE4rVrp+IXVCxemqXU5msvtljF9pxidi1gLqzh2jt3HaUoETde52kg3cg3a0Qw7MYi13xvYdhBEn/00GNe42okxaQIwi9QyYwd1laqjZ78ZSqQNAwgoFA0C8OwzD8gYAR9PurDPET9PfjacrxPqzu3iZr+37A/Egf9h7tw8bpPlzbNYjNE2fwomNn8BdYHsDFQ3jM+bvQML4GGtEn0FsSlIGLyYC1Mr0YioBU+ZMwCZQXXSAr10dk3Q1CrjKq/UkqM1wn0aVsCSyiT4LHV6HI1N0uu6isjNyQmQIixl9dkwTm0xQXndgveg+rRcJo+LLMhDegOTBRnwDubmTbd14kTPzaHUbrWRiW02X351p7m2vWpW1uJ1F6u9HfK/0PSYSbU3p3D6GbERPPxHq6Dp5jK+BBn+dOAzEVsN8PdbR2vOTjxBq3Zc9eoiiCja6crmWvQIxG/7VSRhkqu0e3dx59dPy8f0nKtmBlUKp8sgHzm2MD0XWsFeK06Z5cp7OHb3MG6lMVQskQkpjT57OM6aDpk6G8G92ZNMP2H0oezJXcaden2VNg0sdBCdjK1x3Y6U6UXYYQ090ZhMm+cgEz/eDaslPdrKTjdvAXmT7G6sBg1aDqD0BJN3I7DkRpdzmuMi3pzdLS0bPT0Rc2IOphFVA6hEsAK/UZa9m+osAUZ9rhi6KXWLFdj/K6KyO/tYlk5KvejPzbq0VUnKVfPfQRNhMq6UNQokPWcrN2UbiwUTljcfa1W+d2jSvZC7CULrnL/fIxb3aCT0dkqPBNsp+H3jHhjj2cwaY5nuTPPTMVjq13Z80lHBnBaB9Jv25XmJvJDLd91rb4G7OscuPlZ3umDIZmvjQPzh1ovhryvL38q8DxD6+FuEv/ucqZKK20zNGXdbnp6dpljy9Dc9PKQGhqgzPOhuafPmWPK6HA96T54JWbaO2U+bHGpnDKivAF8ZjFl4cTicZ4ggdq+CJexf18mab4/lv/1W2Rk+Y9TDf3tLbPv99c6R38aL+5Cx371fv8tv7B0JiXuR3tQaJvyhZGrfpUIl7fkrL4E/GWWMq3/pmWplVWgscbePbN+nAyleSpOG+KR6xoNV8aT4WjvN4O4S83RqN8lcWTzQkrHOHWGisWbeXxNcIltdriVsyZJq4iLYlwqjEeExGp1TzME+FYJN7UuFYAjDc0JK1U9dR0fZppM+7NzThQYCzlKeNyh/FnuRkHC4ylvGZ8KDfjWQXGUl4z/jw345oCYymvGR/OzXh2gbGU14z7cjOeU2As5TXjI7kZ1xYYS3nN+GhuxnMLjKX8ZTxMmvttxpXjfwMxCpClvIZ87C4gF7502SOPIR+/C8iFT132yFPIvwOEaV7o',
		{
			includeHealingModel: true,
		},
	),
});
