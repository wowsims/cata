import * as Mechanics from '../../core/constants/mechanics.js';
import * as PresetUtils from '../../core/preset_utils.js';
import { Conjured, Consumes, Flask, Food, Glyphs, Potions, Profession, PseudoStat, Spec, Stat, TinkerHands } from '../../core/proto/common';
import {
	DruidMajorGlyph,
	DruidMinorGlyph,
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
import BalerocMTApl from './apls/balerocMT.apl.json';
import BalerocOTApl from './apls/balerocOT.apl.json';
import BlackhornOTApl from './apls/blackhorn.apl.json';
export const ROTATION_DEFAULT = PresetUtils.makePresetAPLRotation('APL Default', DefaultApl);
export const ROTATION_CLEAVE = PresetUtils.makePresetAPLRotation('2-Target Cleave', CleaveApl);
export const ROTATION_NEF = PresetUtils.makePresetAPLRotation('AoE (Nef Adds)', NefApl);
export const ROTATION_BETH = PresetUtils.makePresetAPLRotation("Beth'tilac Phase 2", BethApl);
export const ROTATION_BALEROC_MT = PresetUtils.makePresetAPLRotation("Baleroc MT", BalerocMTApl);
export const ROTATION_BALEROC_OT = PresetUtils.makePresetAPLRotation("Baleroc OT", BalerocOTApl);
export const ROTATION_BLACKHORN_OT = PresetUtils.makePresetAPLRotation("Blackhorn OT", BlackhornOTApl);

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecGuardianDruid, DefaultSimpleRotation);

// Preset options for EP weights
export const SURVIVAL_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Survival',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.05,
			[Stat.StatStamina]: 1.08,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.93,
			[Stat.StatBonusArmor]: 0.21,
			[Stat.StatDodgeRating]: 0.85,
			[Stat.StatMasteryRating]: 0.31,
			[Stat.StatStrength]: 0.08,
			[Stat.StatAttackPower]: 0.08,
			[Stat.StatHitRating]: 0.22,
			[Stat.StatExpertiseRating]: 0.37,
			[Stat.StatCritRating]: 0.28,
			[Stat.StatHasteRating]: 0.06,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.0,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.185 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.035 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

export const BALANCED_EP_PRESET = PresetUtils.makePresetEpWeights(
	'Balanced',
	Stats.fromMap(
		{
			[Stat.StatHealth]: 0.04,
			[Stat.StatStamina]: 0.88,
			[Stat.StatAgility]: 1.0,
			[Stat.StatArmor]: 0.66,
			[Stat.StatBonusArmor]: 0.15,
			[Stat.StatDodgeRating]: 0.60,
			[Stat.StatMasteryRating]: 0.22,
			[Stat.StatStrength]: 0.16,
			[Stat.StatAttackPower]: 0.15,
			[Stat.StatHitRating]: 0.61,
			[Stat.StatExpertiseRating]: 1.07,
			[Stat.StatCritRating]: 0.36,
			[Stat.StatHasteRating]: 0.10,
		},
		{
			[PseudoStat.PseudoStatMainHandDps]: 0.50,
			[PseudoStat.PseudoStatPhysicalHitPercent]: 0.535 * Mechanics.PHYSICAL_HIT_RATING_PER_HIT_PERCENT,
			[PseudoStat.PseudoStatSpellHitPercent]: 0.075 * Mechanics.SPELL_HIT_RATING_PER_HIT_PERCENT,
		},
	),
);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/mop-classic/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '-2300322312310001220311-020331',
		glyphs: Glyphs.create({
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
	profession1: Profession.Engineering,
	profession2: Profession.ProfessionUnknown,
};

export const PRESET_BUILD_BOSS_DUMMY = PresetUtils.makePresetBuild('Single Target Dummy', {
	rotation: ROTATION_PRESET_SIMPLE,
	encounter: PresetUtils.makePresetEncounter(
		'Single Target Dummy',
		'http://localhost:5173/mop/druid/guardian/?i=rcmxe#eJzjEuVgdGDMYGxgZJzAyNjAxLiBifECE6MTpwCjBaMH4w1GRismAUYhBqkvjLOY2QJyEitTizjYBBiVeDmYDSQDmCKYEliBGp0YVjFzS3EKMoCBnsMJJpYLTOy3mDgFZ80EgZv2j5iaGCWYlOq4CquVchMz80qA2C0xtSgz1S2zKFXJqqSoNFUHLuOSmptflJiTWZWZlx6Un1gEk08Biodk5gLVm+goFZTmlAENqEpFiBSlAgVzgksScwtSU6Cm1gohXPGCKeUHE+NCZog7Ix26mGU5wcymKw5wRZpnz4AAl4MCm9JxJg9miIqGNAeo90QcJCFKT9pbQkQu2CumgcE1e0eI5jf2Rj1MBas+M1ZxByVmpiiEJBalp5YoREiwa91gZKAHCGhxoKZxDSnHkc3zsZg7xxFdDTg0GhZxOs5khIUlVA2LAwBf7n5L',
	),
});

export const PRESET_BUILD_MAGMAW = PresetUtils.makePresetBuild('Magmaw MT', {
	rotation: ROTATION_CLEAVE,
	encounter: PresetUtils.makePresetEncounter(
		'Magmaw MT',
		'http://localhost:5173/mop/druid/guardian/?i=rcmxe#eJzdVE9MHFUYn++9YZn5dtkOT6jLs+KwBEMn7WZZQCsxzkANWZO2AhqkjQcn7FIGl13KrhI5rVUjeCJeLERM6rEnw0Xl0EPTRJvYhF5a4NJGY2LS9tQL9WCceTOLu2Q10RgPfqd5v+/7/b6/GWxVwIJpKAOsAJQJfEVgi8CQqsEJSMM2wADRgEn8krpGQyM5+73svBLSIN6k0GTbCJkgbzW4xCHpCg1ztVkSlrC+J/IWadwlFeSe+SNZV2I0HmZ/QNiMIS4jiTXHG7EB6fGeIraiMhNCWXm8J8dVdOHjyURPkcmcDEisnR9BZmgKZXQXJC5Ivcku7hE2N6Cuvy/wX7pL2QQfxzQbxkgKGdGIi976mCA3YgqwpmtQBXKXidArfDKLzFS7XDUkfUWO7vvbR7LnYS/yATzGjDrKf8V234I9xC1MsgRGUxEmayQGAf+IwV1+dKYG3VeIBArCxyz+EiZ5ArV4FCMboIqUVLkYqlSxWYX5Gs/5Vex8AKKKYf4y9nG3gbiG0Q0IB0mFRlDJZg26rxIJVPxKOvgz+MQu0UQ2okzuF6x4L5EKuSKW8sUyZSpv9KatPOzeh397v4E9zZ/CpngY1Q0IIkWxV+/6A4/yiHjfXvJH2Mt7sJN1+O3XLKGOSGXuz/N+7GKdfscHZn+AVjvsaT6FE2z8bxySwburJb/ek/FJo9WdqbpZAeoeVp7n8E12ru5pADt0DWpv4x9kiwTZ/M48ptcHrd9HzfgS/Jiog/5ZHQemVtlZ5eCauPvLcDf/Wb/3/IVk9ghcpv7f4az1CX1W9b7KyVtW89g37R/dH902j04J+9Vse/fV9sMzp38y9VD8TkOailBp5YIV/F1arLa1Vc9umC/4yJbZ4ZNvm4M3f/DsgZkq07mfr5PF8Gn7/Ky9oKf69bQ+EQNjG6T/wkY+tP5NuXLmu2q9UyfWPx88GCOmUf5SHVwFf0A7ZhAjW5/CBfeoWl6ZnbNzWX0sa0+WnEJef92ZzRpn0oUFPVfIn9cXnFxOL027AbaT0Uv221m9VNADUjprZ3Qnr7+WnSzkM8WE3j04VcrOi3gn75QcO6f3JItHO/ycLVZqlcwtX6SLh4dy9qLjyp8s5Iul+XcmS/423vi/bENaMkbrb0NaemStBi3eq2xDspLwOzCwD1o=',
	),
});

export const PRESET_BUILD_NEF = PresetUtils.makePresetBuild('Nef Adds OT', {
	rotation: ROTATION_NEF,
	encounter: PresetUtils.makePresetEncounter(
		'Nef Adds OT',
		'http://localhost:3333/mop/druid/guardian/?i=rcmxe#eJzll19oW1Ucx3POuU1Pfu1ccqhrcmTrWXSjiywkabt11ZnbilJEpI6OMhTZXXJLM5LceHM7bccgdkh1Pii+6IYTN1/EJw2iW9W5dQP/oNA97I8FmexpsM0XEYoPm/dPksVWm6ooyP295J5ffr/v73c+v5PDDbRTJJCMxlAJodcQKmH0AUZzGA00+dFlhPqwHzEPP99yhHiHssqEqlOvH4V9lMRCQ9hMGCCnEXqftHBfwGNbVP4KS3O4eR5XPT8mr+JLECThFnbHBQHwcglwMBBuhiYgm+NFJnHc52HDfAcMskchFGmnEvNNSV4m0YVfJU7BC1I82lMEHglSL1s1i4BhfxNI9Op1iZsigLbCqnAL+MrIDKVvXiHc/PRPXyMsxRUYYTv/nmrvItW6pSVuFaEzZfTPWv9dEVPNbr18grBDiE8jOMD2w1N8F9wTCS0VMEMBx+PQx3qXK75MbmxR9fotmk3wys7ZJH8e8iz7n7VSmSF9q4J3O3tg+bp/OilLggdgNf34IBa3q4bYfj4BGsv9JdUN/N5FysAifkoYmUcebh/o7tgGCIbXQFsZLS1am+2/VLyrUXH7uD7OHwOZPQQ0YXnOT2P7dKL6Nkynczq7YE2kzRwm3VuJdbrrLtpqJ3+W2IO8D2I8Cv7wXdBaRj5bgdApr60qsdaZOp+TvcXJ/v4gYiHeDqvncStgmqpJS9aKAad22Kkr5lnizVY/9BVvzf32y6TmvtlZc996oYmZT84dw+7nm6CDrYXWMNRtt46iqcIr7prExZfuPFs97ubPwDDbsUJeEd5ZX+CTBQnaI3dbv4eZqmMJw438PluT/LFmrcPqft/ouYbTCxgdI869uks+RDb6rKfS0A9yoOPF60+eWHchuenZ7XDp3cLNZOi9kW8S5z78KSm84dLaQWKHenarcuVebpNDRw5b9nVym+OZS64fte1isv+7by27kUycIYWTN26jyXX9+UxOMdS0GNDyqhhRdD2j6SLRIwbFZhEXO4M0chl5/odWSn8p1y3loxNT/YtjHEIfHZcPo7O/dD796uyFZDVGfh0dRebJCDySVVOGrqXGDVU8rI3njciBJ8Zze1RdaKOi/suUUjSKwtBETkur2agY1gwlK1JWinguk82KPaooFnRVSQt1n5rPTghtn6lijKlCzdth5io9ritGRsubGcaYUISu5NNaLjNpDkgbHS2qRnS9019QTryDCzONZ5hw8QwtRJ82RtTlckSfNUbU7XJEnzdG1ONyRKcaI9rickRfNEa01eWITjdG1OtyRGcaI9rmZkTHcGHWQtSx/LttzOWMzq6AkZv/AFiMzq2AkZtfsH8D6cCcbQ==',
	),
});

export const PRESET_BUILD_BETHTILAC = PresetUtils.makePresetBuild("Beth'tilac Phase 2", {
	rotation: ROTATION_BETH,
	encounter: PresetUtils.makePresetEncounter(
		"Beth'tilac Phase 2",
		'http://localhost:5173/mop/druid/guardian/?i=rcmxe#eJzVVF1oHFUY3Xtndjtzk2x3b43dvbV6M5Y2LknYzZ9pUGcSSojVSi2VxkjRye7d7iS7M+vObENXCtEiVn1QAlIbilh8KaUPNQ/+5EGLglZRqIg1lmpBREQrCL7EH9Q7d3bTjdUHH/1gGL7znXvu+Q7MoHYFGKAA5gB4HoA5CM5CcAGCUTUGhsA4WAZgGMYADpF5dUGK7C6ah1hFaYkBrUWR0snd8JEwPzYaOQfAeQBPSy1EjYdEbTHOQ/kSVON9oj7Vv4YvKAmoPQvR5GNaybRsjz9jJqtYbMyqMG3Yq1RZ1+pkBys5FbNo1Sz7wB7HrDTmOY7vtUqc39+llavFg1ygxurIYXwGoASSlROLAMcCIzuM4D1poBifvHhFwooAjKkAWeLcAAk9IJDYZ783OPVTsZPvNZC5XQbahNYrrx2B9M9GNQQM00AaiionBptndSP7AyPnHjZWM0FxFCEygom4tg6FkdSdcVE7UqYj3NYvK7KmIg53p3syLpYJHA7hzWRTLYna+tM9abeLZtL9/I0VEqmvfCO5oYZRJNM/eA33F24nG2pxFB5qoi810ZtkxO51+sD2a7AfQAe5pbaZ0/uGfHwDiV8fQyfZWtuCopnMYGCw93afiknsulAmyQTaicdRay/CMAa5o0+egoikEgrAbe+AJpDwYBDoQzeliCLj1unmkYK4Od89QRx582fZn2GD3IXSpAfFtChqXQSqOCApT0SEPtdYasL89RAcDBS+OAKEwq2kA7VfgliwoJJtukjxe0FCJIjypaclrJJ1vk/lx85V+I/Hw/heshPdgYcDJ2s2bdNakLoI6ufRxlS7Esbq6QYQuOoLXHG+uDBKWkV/8eja/q0r8pq+scU0KaBJPPEfMu4iqWZnr6/IKJnayCNTlxrAv2SeIp1CWfpn5TV7tBH+n+BhHRvw2+/gr2AFgpNS8Jk8aDwjbVXFp/b9ZSO+542bn/zh/mX9tkfvRJ+/Uv5STx7N7r88Jn2j04j2oTwuCWro5axR//H8pCcXjvv1gb49QC7oHXlRF/WRjz/y66re+5xcnl+QautHmVfY5llFM0t7B+g4nUiA1DII/Q9rLve+0dSejR8+NvJ3jkjk1LvyyHGg73vo1bevfqU3OMY8+A3g6N12tljNMTpWYXbtUOpbsMvJsSL1CoxOVSs2LRdMl1EnLxBmZ52q7bEKnbW8goBcz8zO8H82nXJcl+bMknmAUbNU7qH3OR7jFNOjezlvn5VzZre59B6L02atYpHajkenGC3597EcZQeZTS3/HsulTtmzHN661OXDrMdyXdS17CyjluejputWS/yQkDdt7i/f7Zn2TKDMVcsVxgHu1bfmOdR1zBnhN8emqvl8TwL8BQ0Xc7A=',
	),
});

export const PRESET_BUILD_BALEROC_MT = PresetUtils.makePresetBuild("Baleroc MT", {
	rotation: ROTATION_BALEROC_MT,
	encounter: PresetUtils.makePresetEncounter(
		"Baleroc MT",
		'http://localhost:5173/mop/druid/guardian/?i=cmxe#eJzVUr9PFEEUZvbQ3IkmB5oIJJIHlV4QkYhRY7LLxeCRQLzIxWDn3O7s3eRmZ8/Z2bscFbEyxsLQQaNURit7Y60mklAZbKwoKLQ12vl2lkPA+Af4is3szPd+fN/3TpzNEiAOqZMVQp4RsmKRNxbZskgxlyfXSIlsE3LDypOBnuGH65njZUE7TGX78mTsVDYzOVS2lqwHxzCv2PM60zec6+8xMeF8sHp3Le+nRTYy6dV950nmUs4c5787/etrSfyyL6SHL/bl9PDKHkrxuzYcH3ueK2XSnMmWs1f7jDOUQj/a19ObLXvUN/HZntn8lMQ3e2on2/yxkVk+WaSCqdCFqWkowdJgb2Gb9PyHseK9dw78nn9UfjlzFGPk2Hl7emaNXOm/evPd3a92F+OUyCRZJU/JwEiFygZEbdoEP1Rwi7k8oJqHEoqCeqxQn/PBrTO3wbxx0HUG1TCKoM2FgCoDjcnMg2oHTJkpUBQxCoFUGvQC5TJ982LFZe2vBlhKemE7mhjEeR5nsmQA5iTXnAq4U0kgywxCH26LUHUg0tRtRIVtC4eSoby4jFaOd1ubmbjJ1cw0b8Z4E0vNhfmNmtjb5zivqQNuiG/AI6CuGwexwDRvAuapqiGDFhUxwychwrbJxml0CDVGFQShYkBrNcWiiLeY6EC7zgU20MkIAWM6YZokTU1PNqBUxjNi66HwjMjJi89VpP+lxjhUYw0BbTAQ3GdpQ4/7PndjofdrLFS6qs5JnykZHpEUKnVkx2UzNjR5TWIdD7hvfDOmI0NTOqJVgexHu/uxSjYJenGu68U9rvG72KSqgR4HzcRoVnhBZg/RCbiMtfGrK3+6MxHfkxfqjIpkYFpjqL2x7I/6iMUxAyo7h/qlrgNqRFtMYSZunm4zJvdQs4IGDBTzE42ZASZNE45dQiNIyEU+hUXNqGcWCdsdbLI/ReKFIVdYROtxE5KffckVC3CjPbw+QLLbJe/8Bt31bfM='
	),
});

export const PRESET_BUILD_BALEROC_OT = PresetUtils.makePresetBuild("Baleroc OT", {
	rotation: ROTATION_BALEROC_OT,
	encounter: PresetUtils.makePresetEncounter(
		"Baleroc OT",
		'http://localhost:5173/mop/druid/guardian/?i=rcmxe#eJzVUk9oHFUY3zcz2c6+JGUztSQZsL7sQeKShG3aFA3F3QSJG0kxmBDizZeZN7uvOzuzzMxmSU6xHhqLhxIQbRG1J6EnCYIY0YMXFVpILyXtwYr00IMFQZD0In7vTSbdre3Bo++wO9/3/b4/v9/34X4dEVRCVbSB0GWENhT0pYJ2FTSdyaKXURntITSZmlSyyEiZl9NX1fS8S9dYoHdnUa5XVwqD88qy8k4XpE6nrqvdZqYvJd9Y6SdF21XUu0ri+bX4m/Jh14Ca6zYeu3AfTpsaVgb6ckdwF1ZHT4b4ONbPp7GmP9rXchkM7tHC2MnQ0ExlMmUMmv24N9eNM9tIYL66mDbF/8f31GeFPtlGceiShrGSVQxAb8ahHQgtm0u4bMzgnnFsQBS8ty4q2MwP6Mjo/QG1OU2YD6NTMqYZPefbQ1ANK6dDE4P9zZ+aiBgl81VcMMdwNncU92yjjISr+oV0UmGnzRdXOBNXuPMekhWGzBfwsbtKVqIU3TpsowtLQrCpSyqfvq8aGfOImFF/OHzo/vvdLmPOfAOfNSbjOTpYtqkF+bg/f1zvMjLXE0fc7lQ8E+Blw6Nmj7Rvb3ba39/TOuyEQ9V08LKx9B/0zZvD7ZN9va/JyTQjs5M4nqq3yBN11afX7WDRa8KpglQfTQjzgWLvK+iaGl/l26VL6osZ+Tn3sNR39Yp4j4ovxR93ioMx7EGRpHOfZcpqDC2slg6u+rnSYAz9ufhK7NktDjny3S5O3bwh3u/F8ft6469r6nrPNHVZ4FtkfIKUyfKAlt9Dqf/h27B/LLWZwxfmv5h6EiPluP/tsakr6HTfmbPfvfVLMcGUyqiAttAHyDixSL0aCVu0QRw/IK8xi9dpxH2PTLvUZvnqrEOsKrNqzB4hUZWRFT8MSYu7LllhJIJkZpOVNSLLjJOAAiYAIPUk+hzlXhyzmwH3Kv9qAKU822+FYwMwz6YKZ0pmPR5x6pI3FwVknRHfIa+7frBGwohatTC/p8BQnu+NrsMqR5LWciYucyMmmzea4Gl6EXelGTagt8NhXlmHWD7ECA8JtaxmvelCmj1G5mhQAQar1G0yCLmu35LZME3kkwqjAan7ASO0UglYGPJV5q6RVpW70CASI9QZiwRTkTQ+UaiR8jx8A7bqu7YUWUQcHoTRs9QYISvNiNRpjRGXOyxuaHPH4VbTjQ5rnFtMVJ31HBZ4/hOSksUqsONeoylp8ooHdWzCHbk3uXRgKEuHdMUF9kPxefxR3EI3Eezi+WQXSzyC34UGDWqw43pDLJrlP0czHXTq3GtGcl+J/PHNhPxAXlJl1BUD0woD7eXKHqsPWBizTr21jn7x1gloRFdZAJlweVGLMe8ANePSOiMBc4TGTAJFU8ExIXSitIUs4JNfiBi15SFBu/Ymh1OIXUhy+QVYPVyCMA4lD1gdLtoGdxvJpEu29A+pJMjo'
	),
});

export const PRESET_BUILD_BLACKHORN_OT = PresetUtils.makePresetBuild("Blackhorn OT", {
	rotation: ROTATION_BLACKHORN_OT,
	encounter: PresetUtils.makePresetEncounter(
		"Blackhorn OT",
		'http://localhost:5173/mop/druid/guardian/?i=cmxe#eJzdUs9rE0EUzpu1MWmCTUPRpgedFmJtQImJEQ3iboLSnCS1QuvNSXaSHZPshp1Jgznl4sWLUvHQHkQvHjyUXoRKL95swUJO2quehOhRKJ6cyQ+t4j+g32F4881733u8+UZP+CAOBljQBngE0EawiaCDIOsPwSXIwT5A2pNGIQh7pvi65s1XyT3q+gIhmAn4tHgkj+6MyLqs/w3ADqBNL3qpBab8454ezhk76EgHaZ/RdzhA8Fzr07eNB9ppv4qeBT8Z4xvzH1sTGx/0uVP3uwtbJ7/pke3ZscCrF1917J3pHs1pvVRPmxkD2Qkjsr6msKtf7jMdfbrUw3s9s/dO4YueeDpS725prckl4tYIF9TF2SopVizHtXEihXN4edIb2wfPP4i2+dY4dDUi5Enmz5zeaq4Ew5k16G/rtT7IOdBXYckH4fAtYlcwb5I6Zrbczwqpxi5er5I6pyYWrEbxGWZjTouObfI5XKCiSakkWK1RJULlDOv5dF/5qrEKZal8PGOa2KRnSVGwFSKYY/f0YtnFfu2Q+L0BEbhpsaKFhYNNxkmhSvG848pkMsvlq5C/N+y0LDstyE5j2Yb8zrpFOMW5PI7GklHslHDB4VzdXVojzGZ2WSrL0QdqmPFD01Gzr2nclJp3fSg8eoO6JcyFfIpdWBQutcvCUqrCooouVpRg3mlKSw3YDK/TolCqhUaplEbxaBqlomntfFwdKkqoKKGipIqSqWjiIarvbWut4HCqgSm1/8OU1xKZx383JWbHfppy95cp4/ADEKDaLg==',
	),
});
