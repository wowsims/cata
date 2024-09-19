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

export const InfectedWoundsBuild = {
	name: 'Infected Wounds',
	data: SavedTalents.create({
		talentsString: '-2302322310310001220311-020331',
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

export const PRESET_BUILD_BOSS_DUMMY = PresetUtils.makePresetBuild('Single Target Dummy', {
	rotation: ROTATION_PRESET_SIMPLE,
	encounter: PresetUtils.makePresetEncounter(
		'Single Target Dummy',
		'http://localhost:5173/cata/druid/guardian/?i=rcmxe#eJzjEuVgdGDMYGxgZJzAyNjAxLiBifECE6MTpwCjBaMH4w1GRismAUYhBqkvjLOY2QJyEitTizjYBBiVeDmYDSQDmCKYEliBGp0YVjFzS3EKMoCBnsMJJpYLTOy3mDgFZ80EgZv2j5iaGCWYlOq4CquVchMz80qA2C0xtSgz1S2zKFXJqqSoNFUHLuOSmptflJiTWZWZlx6Un1gEk08Biodk5gLVm+goFZTmlAENqEpFiBSlAgVzgksScwtSU6Cm1gohXPGCKeUHE+NCZog7Ix26mGU5wcymKw5wRZpnz4AAl4MCm9JxJg9miIqGNAeo90QcJCFKT9pbQkQu2CumgcE1e0eI5jf2Rj1MBas+M1ZxByVmpiiEJBalp5YoREiwa91gZKAHCGhxoKZxDSnHkc3zsZg7xxFdDTg0GhZxOs5khIUlVA2LAwBf7n5L',
	),
});

export const PRESET_BUILD_MAGMAW = PresetUtils.makePresetBuild('Magmaw MT', {
	rotation: ROTATION_CLEAVE,
	encounter: PresetUtils.makePresetEncounter(
		'Magmaw MT',
		'http://localhost:5173/cata/druid/guardian/?i=rcmxe#eJzdVE9MHFUYn++9YZn5dtkOT6jLs+KwBEMn7WZZQCsxzkANWZO2AhqkjQcn7FIGl13KrhI5rVUjeCJeLERM6rEnw0Xl0EPTRJvYhF5a4NJGY2LS9tQL9WCceTOLu2Q10RgPfqd5v+/7/b6/GWxVwIJpKAOsAJQJfEVgi8CQqsEJSMM2wADRgEn8krpGQyM5+73svBLSIN6k0GTbCJkgbzW4xCHpCg1ztVkSlrC+J/IWadwlFeSe+SNZV2I0HmZ/QNiMIS4jiTXHG7EB6fGeIraiMhNCWXm8J8dVdOHjyURPkcmcDEisnR9BZmgKZXQXJC5Ivcku7hE2N6Cuvy/wX7pL2QQfxzQbxkgKGdGIi976mCA3YgqwpmtQBXKXidArfDKLzFS7XDUkfUWO7vvbR7LnYS/yATzGjDrKf8V234I9xC1MsgRGUxEmayQGAf+IwV1+dKYG3VeIBArCxyz+EiZ5ArV4FCMboIqUVLkYqlSxWYX5Gs/5Vex8AKKKYf4y9nG3gbiG0Q0IB0mFRlDJZg26rxIJVPxKOvgz+MQu0UQ2okzuF6x4L5EKuSKW8sUyZSpv9KatPOzeh397v4E9zZ/CpngY1Q0IIkWxV+/6A4/yiHjfXvJH2Mt7sJN1+O3XLKGOSGXuz/N+7GKdfscHZn+AVjvsaT6FE2z8bxySwburJb/ek/FJo9WdqbpZAeoeVp7n8E12ru5pADt0DWpv4x9kiwTZ/M48ptcHrd9HzfgS/Jiog/5ZHQemVtlZ5eCauPvLcDf/Wb/3/IVk9ghcpv7f4az1CX1W9b7KyVtW89g37R/dH902j04J+9Vse/fV9sMzp38y9VD8TkOailBp5YIV/F1arLa1Vc9umC/4yJbZ4ZNvm4M3f/DsgZkq07mfr5PF8Gn7/Ky9oKf69bQ+EQNjG6T/wkY+tP5NuXLmu2q9UyfWPx88GCOmUf5SHVwFf0A7ZhAjW5/CBfeoWl6ZnbNzWX0sa0+WnEJef92ZzRpn0oUFPVfIn9cXnFxOL027AbaT0Uv221m9VNADUjprZ3Qnr7+WnSzkM8WE3j04VcrOi3gn75QcO6f3JItHO/ycLVZqlcwtX6SLh4dy9qLjyp8s5Iul+XcmS/423vi/bENaMkbrb0NaemStBi3eq2xDspLwOzCwD1o=',
	),
});

export const PRESET_BUILD_NEF = PresetUtils.makePresetBuild('Nef Adds OT', {
	rotation: ROTATION_NEF,
	encounter: PresetUtils.makePresetEncounter(
		'Nef Adds OT',
		'http://localhost:3333/cata/druid/guardian/?i=rcmxe#eJzll19oW1Ucx3POuU1Pfu1ccqhrcmTrWXSjiywkabt11ZnbilJEpI6OMhTZXXJLM5LceHM7bccgdkh1Pii+6IYTN1/EJw2iW9W5dQP/oNA97I8FmexpsM0XEYoPm/dPksVWm6ooyP295J5ffr/v73c+v5PDDbRTJJCMxlAJodcQKmH0AUZzGA00+dFlhPqwHzEPP99yhHiHssqEqlOvH4V9lMRCQ9hMGCCnEXqftHBfwGNbVP4KS3O4eR5XPT8mr+JLECThFnbHBQHwcglwMBBuhiYgm+NFJnHc52HDfAcMskchFGmnEvNNSV4m0YVfJU7BC1I82lMEHglSL1s1i4BhfxNI9Op1iZsigLbCqnAL+MrIDKVvXiHc/PRPXyMsxRUYYTv/nmrvItW6pSVuFaEzZfTPWv9dEVPNbr18grBDiE8jOMD2w1N8F9wTCS0VMEMBx+PQx3qXK75MbmxR9fotmk3wys7ZJH8e8iz7n7VSmSF9q4J3O3tg+bp/OilLggdgNf34IBa3q4bYfj4BGsv9JdUN/N5FysAifkoYmUcebh/o7tgGCIbXQFsZLS1am+2/VLyrUXH7uD7OHwOZPQQ0YXnOT2P7dKL6Nkynczq7YE2kzRwm3VuJdbrrLtpqJ3+W2IO8D2I8Cv7wXdBaRj5bgdApr60qsdaZOp+TvcXJ/v4gYiHeDqvncStgmqpJS9aKAad22Kkr5lnizVY/9BVvzf32y6TmvtlZc996oYmZT84dw+7nm6CDrYXWMNRtt46iqcIr7prExZfuPFs97ubPwDDbsUJeEd5ZX+CTBQnaI3dbv4eZqmMJw438PluT/LFmrcPqft/ouYbTCxgdI869uks+RDb6rKfS0A9yoOPF60+eWHchuenZ7XDp3cLNZOi9kW8S5z78KSm84dLaQWKHenarcuVebpNDRw5b9nVym+OZS64fte1isv+7by27kUycIYWTN26jyXX9+UxOMdS0GNDyqhhRdD2j6SLRIwbFZhEXO4M0chl5/odWSn8p1y3loxNT/YtjHEIfHZcPo7O/dD796uyFZDVGfh0dRebJCDySVVOGrqXGDVU8rI3njciBJ8Zze1RdaKOi/suUUjSKwtBETkur2agY1gwlK1JWinguk82KPaooFnRVSQt1n5rPTghtn6lijKlCzdth5io9ritGRsubGcaYUISu5NNaLjNpDkgbHS2qRnS9019QTryDCzONZ5hw8QwtRJ82RtTlckSfNUbU7XJEnzdG1ONyRKcaI9rickRfNEa01eWITjdG1OtyRGcaI9rmZkTHcGHWQtSx/LttzOWMzq6AkZv/AFiMzq2AkZtfsH8D6cCcbQ==',
	),
});
