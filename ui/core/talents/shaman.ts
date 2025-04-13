import { ShamanMajorGlyph, ShamanMinorGlyph, ShamanPrimeGlyph, ShamanTalents } from '../proto/shaman.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import ShamanTalentJson from './trees/shaman.json';

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig(ShamanTalentJson);

export const shamanGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[ShamanPrimeGlyph.GlyphOfLavaBurst]: {
			name: "Glyph of Lava Burst",
			description: "Your Lava Burst spell deals $55454s1% more damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfShocking]: {
			name: "Glyph of Shocking",
			description: "Reduces your global cooldown when casting your shock spells by ${$m2/-1000}.1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfEarthlivingWeapon]: {
			name: "Glyph of Earthliving Weapon",
			description: "Increases the effectiveness of your Earthliving weapon's periodic healing by $55439s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfFireElementalTotem]: {
			name: "Glyph of Fire Elemental Totem",
			description: "Reduces the cooldown of your Fire Elemental Totem by ${$55455m1/-60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfFlameShock]: {
			name: "Glyph of Flame Shock",
			description: "Increases the duration of your Flame Shock by $55447s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfFlametongueWeapon]: {
			name: "Glyph of Flametongue Weapon",
			description: "Increases your spell critical strike chance by $55451s1% on each of your weapons with Flametongue Weapon active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfLightningBolt]: {
			name: "Glyph of Lightning Bolt",
			description: "Increases the damage dealt by Lightning Bolt by $55453s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfStormstrike]: {
			name: "Glyph of Stormstrike",
			description: "Increases the critical strike chance bonus from your Stormstrike ability by an additional $55446s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfLavaLash]: {
			name: "Glyph of Lava Lash",
			description: "Increases the damage dealt by your Lava Lash ability by $55444s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfWaterShield]: {
			name: "Glyph of Water Shield",
			description: "Increases the passive mana regeneration of your Water Shield spell by $55436s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfWindfuryWeapon]: {
			name: "Glyph of Windfury Weapon",
			description: "Increases the chance per swing for Windfury Weapon to trigger by $55445s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfFeralSpirit]: {
			name: "Glyph of Feral Spirit",
			description: "Your spirit wolves gain an additional $63271s1% of your attack power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfRiptide]: {
			name: "Glyph of Riptide",
			description: "Increases the duration of Riptide by $63273s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfEarthShield]: {
			name: "Glyph of Earth Shield",
			description: "Increases the amount healed by your Earth Shield by $63279s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
		[ShamanPrimeGlyph.GlyphOfUnleashedLightning]: {
			name: "Glyph of Unleashed Lightning",
			description: "Allows Lightning Bolt to be cast while moving.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primeshaman.jpg",
		},
	},
	majorGlyphs: {
		[ShamanMajorGlyph.GlyphOfChainHeal]: {
			name: "Glyph of Chain Heal",
			description: "Increases healing done by your Chain Heal spell to targets beyond the first by $55437s2%, but decreases the amount received by the initial target by $55437s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfChainLightning]: {
			name: "Glyph of Chain Lightning",
			description: "Your Chain Lightning spell now strikes $55449s1 additional targets, but deals $55449s2% less initial damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFireNova]: {
			name: "Glyph of Fire Nova",
			description: "Increases the radius of your Fire Nova spell by $55450s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHealingStreamTotem]: {
			name: "Glyph of Healing Stream Totem",
			description: "Your Healing Stream Totem also increases the Fire, Frost, and Nature resistance of party and raid members within $8185a1 yards by ${$cond($lte($PL,70),$PL,$cond($lte($PL,80),$PL+($PL-70)*5,$PL+($PL-70)*5+($PL-80)*7))}.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHealingWave]: {
			name: "Glyph of Healing Wave",
			description: "Your Healing Wave also heals you for $55440s1% of the healing effect when you heal someone else.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfTotemicRecall]: {
			name: "Glyph of Totemic Recall",
			description: "Causes your Totemic Recall ability to return an additional $55438s1% of the mana cost of any recalled totems.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfLightningShield]: {
			name: "Glyph of Lightning Shield",
			description: "Your Lightning Shield can no longer drop below 3 charges from dealing damage to attackers.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfGroundingTotem]: {
			name: "Glyph of Grounding Totem",
			description: "Instead of absorbing a spell, your Grounding Totem reflects the next harmful spell back at its caster, but the cooldown of your Grounding Totem is increased by ${$55441m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfFrostShock]: {
			name: "Glyph of Frost Shock",
			description: "Increases the duration of your Frost Shock by ${$55443m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfElementalMastery]: {
			name: "Glyph of Elemental Mastery",
			description: "While your Elemental Mastery ability is active, you take $55452s1% less damage from all sources.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfGhostWolf]: {
			name: "Glyph of Ghost Wolf",
			description: "Your Ghost Wolf form grants an additional $59289s1% movement speed.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfThunder]: {
			name: "Glyph of Thunder",
			description: "Reduces the cooldown on Thunderstorm by ${$63270m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfShamanisticRage]: {
			name: "Glyph of Shamanistic Rage",
			description: "Activating your Shamanistic Rage ability also cleanses you of all dispellable Magic debuffs.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfHex]: {
			name: "Glyph of Hex",
			description: "Reduces the cooldown of your Hex spell by ${$63291m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
		[ShamanMajorGlyph.GlyphOfStoneclawTotem]: {
			name: "Glyph of Stoneclaw Totem",
			description: "Your Stoneclaw Totem also places a damage absorb shield on you, equal to $63298s1 times the strength of the shield it places on your totems.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorshaman.jpg",
		},
	},
	minorGlyphs: {
		[ShamanMinorGlyph.GlyphOfWaterBreathing]: {
			name: "Glyph of Water Breathing",
			description: "Your Water Breathing spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
		[ShamanMinorGlyph.GlyphOfAstralRecall]: {
			name: "Glyph of Astral Recall",
			description: "Reduces the cooldown of your Astral Recall spell by ${$58058m1/-60000}.1 min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
		[ShamanMinorGlyph.DeprecatedGlyphOfTheArcticWolf]: {
			name: "Deprecated Glyph of the Arctic Wolf",
			description: "Alters the appearance of your Ghost Wolf transformation, causing it to resemble an arctic wolf.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_rune_05.jpg",
		},
		[ShamanMinorGlyph.GlyphOfRenewedLife]: {
			name: "Glyph of Renewed Life",
			description: "Your Reincarnation spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
		[ShamanMinorGlyph.GlyphOfTheArcticWolf]: {
			name: "Glyph of the Arctic Wolf",
			description: "Alters the appearance of your Ghost Wolf transformation, causing it to resemble an arctic wolf.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
		[ShamanMinorGlyph.GlyphOfWaterWalking]: {
			name: "Glyph of Water Walking",
			description: "Your Water Walking spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
		[ShamanMinorGlyph.GlyphOfThunderstorm]: {
			name: "Glyph of Thunderstorm",
			description: "Increases the mana you receive from your Thunderstorm spell by $62132s1%, but it no longer knocks enemies back or reduces their movement speed.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorshaman.jpg",
		},
	},
};
