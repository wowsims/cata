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
export const PRERAID_PRESET = PresetUtils.makePresetGear('Pre-Raid', PreraidGear);
import P1Gear from './gear_sets/p1.gear.json';
export const P1_PRESET = PresetUtils.makePresetGear('P1/P2', P1Gear);
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
import BethApl from './apls/bethtilac.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);
export const ROTATION_CLEAVE = PresetUtils.makePresetAPLRotation('2-Target Cleave', CleaveApl);
export const ROTATION_NEF = PresetUtils.makePresetAPLRotation('AoE (Nef Adds)', NefApl);
export const ROTATION_BETH = PresetUtils.makePresetAPLRotation("Beth'tilac Phase 2", BethApl);

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecGuardianDruid, DefaultSimpleRotation);

// Preset options for EP weights
export const SURVIVAL_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Survival',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.06,
			[Stat.StatStamina]: 1.18,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 1.44,
			[Stat.StatBonusArmor]: 0.33,
			[Stat.StatDodgeRating]: 0.69,
			[Stat.StatMasteryRating]: 0.31,
			[Stat.StatStrength]: 0.10,
			[Stat.StatAttackPower]: 0.09,
			[Stat.StatHitRating]: 0.23,
			[Stat.StatExpertiseRating]: 0.46,
			[Stat.StatCritRating]: 0.41,
			[Stat.StatHasteRating]: 0.06,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.23 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

export const BALANCED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Balanced',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.04,
			[Stat.StatStamina]: 1.0,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 1.11,
			[Stat.StatBonusArmor]: 0.25,
			[Stat.StatDodgeRating]: 0.52,
			[Stat.StatMasteryRating]: 0.23,
			[Stat.StatStrength]: 0.16,
			[Stat.StatAttackPower]: 0.15,
			[Stat.StatHitRating]: 0.53,
			[Stat.StatExpertiseRating]: 0.99,
			[Stat.StatCritRating]: 0.44,
			[Stat.StatHasteRating]: 0.11,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.45,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.495 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.035 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
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
			prime2: DruidPrimeGlyph.GlyphOfRip,
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

export const PRESET_BUILD_BETHTILAC = PresetUtils.makePresetBuild("Beth'tilac Phase 2", {
	rotation: ROTATION_BETH,
	encounter: PresetUtils.makePresetEncounter(
		"Beth'tilac Phase 2",
		'http://localhost:5173/cata/druid/guardian/?i=rcmxe#eJzVVF1oHFUY3Xtndjtzk2x3b43dvbV6M5Y2LknYzZ9pUGcSSojVSi2VxkjRye7d7iS7M+vObENXCtEiVn1QAlIbilh8KaUPNQ/+5EGLglZRqIg1lmpBREQrCL7EH9Q7d3bTjdUHH/1gGL7znXvu+Q7MoHYFGKAA5gB4HoA5CM5CcAGCUTUGhsA4WAZgGMYADpF5dUGK7C6ah1hFaYkBrUWR0snd8JEwPzYaOQfAeQBPSy1EjYdEbTHOQ/kSVON9oj7Vv4YvKAmoPQvR5GNaybRsjz9jJqtYbMyqMG3Yq1RZ1+pkBys5FbNo1Sz7wB7HrDTmOY7vtUqc39+llavFg1ygxurIYXwGoASSlROLAMcCIzuM4D1poBifvHhFwooAjKkAWeLcAAk9IJDYZ783OPVTsZPvNZC5XQbahNYrrx2B9M9GNQQM00AaiionBptndSP7AyPnHjZWM0FxFCEygom4tg6FkdSdcVE7UqYj3NYvK7KmIg53p3syLpYJHA7hzWRTLYna+tM9abeLZtL9/I0VEqmvfCO5oYZRJNM/eA33F24nG2pxFB5qoi810ZtkxO51+sD2a7AfQAe5pbaZ0/uGfHwDiV8fQyfZWtuCopnMYGCw93afiknsulAmyQTaicdRay/CMAa5o0+egoikEgrAbe+AJpDwYBDoQzeliCLj1unmkYK4Od89QRx582fZn2GD3IXSpAfFtChqXQSqOCApT0SEPtdYasL89RAcDBS+OAKEwq2kA7VfgliwoJJtukjxe0FCJIjypaclrJJ1vk/lx85V+I/Hw/heshPdgYcDJ2s2bdNakLoI6ufRxlS7Esbq6QYQuOoLXHG+uDBKWkV/8eja/q0r8pq+scU0KaBJPPEfMu4iqWZnr6/IKJnayCNTlxrAv2SeIp1CWfpn5TV7tBH+n+BhHRvw2+/gr2AFgpNS8Jk8aDwjbVXFp/b9ZSO+542bn/zh/mX9tkfvRJ+/Uv5STx7N7r88Jn2j04j2oTwuCWro5axR//H8pCcXjvv1gb49QC7oHXlRF/WRjz/y66re+5xcnl+QautHmVfY5llFM0t7B+g4nUiA1DII/Q9rLve+0dSejR8+NvJ3jkjk1LvyyHGg73vo1bevfqU3OMY8+A3g6N12tljNMTpWYXbtUOpbsMvJsSL1CoxOVSs2LRdMl1EnLxBmZ52q7bEKnbW8goBcz8zO8H82nXJcl+bMknmAUbNU7qH3OR7jFNOjezlvn5VzZre59B6L02atYpHajkenGC3597EcZQeZTS3/HsulTtmzHN661OXDrMdyXdS17CyjluejputWS/yQkDdt7i/f7Zn2TKDMVcsVxgHu1bfmOdR1zBnhN8emqvl8TwL8BQ0Xc7A=',
	),
});
