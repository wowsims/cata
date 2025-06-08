import { DruidMajorGlyph, DruidMinorGlyph, DruidPrimeGlyph, DruidTalents } from '../proto/druid.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DruidTalentJson from './trees/druid.json';export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentJson);

export const druidGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[DruidPrimeGlyph.GlyphOfMangle]: {
			name: "Glyph of Mangle",
			description: "Increases the damage done by Mangle by $54813s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_mangle2.jpg",
		},
		[DruidPrimeGlyph.GlyphOfBloodletting]: {
			name: "Glyph of Bloodletting",
			description: "Each time you Shred or Mangle in Cat Form, the duration of your Rip on the target is extended by $54815s1 sec, up to a maximum of $54815s2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_vampiricaura.jpg",
		},
		[DruidPrimeGlyph.GlyphOfRip]: {
			name: "Glyph of Rip",
			description: "Increases the periodic damage of your Rip by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_ghoulfrenzy.jpg",
		},
		[DruidPrimeGlyph.GlyphOfSwiftmend]: {
			name: "Glyph of Swiftmend",
			description: "Your Swiftmend ability no longer consumes a Rejuvenation or Regrowth effect from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_idolofrejuvenation.jpg",
		},
		[DruidPrimeGlyph.GlyphOfRegrowth]: {
			name: "Glyph of Regrowth",
			description: "Your Regrowth heal-over-time will automatically refresh its duration on targets at or below $s1% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_resistnature.jpg",
		},
		[DruidPrimeGlyph.GlyphOfRejuvenation]: {
			name: "Glyph of Rejuvenation",
			description: "Increases the healing done by your Rejuvenation by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_rejuvenation.jpg",
		},
		[DruidPrimeGlyph.GlyphOfLifebloom]: {
			name: "Glyph of Lifebloom",
			description: "Increases the critical effect chance of your Lifebloom by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_herb_felblossom.jpg",
		},
		[DruidPrimeGlyph.GlyphOfStarfire]: {
			name: "Glyph of Starfire",
			description: "Your Starfire ability increases the duration of your Moonfire effect on the target by $54845s1 sec, up to a maximum of $54845s2 additional seconds.  Only functions on the target with your most recently applied Moonfire.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_starfire.jpg",
		},
		[DruidPrimeGlyph.GlyphOfInsectSwarm]: {
			name: "Glyph of Insect Swarm",
			description: "Increases the damage of your Insect Swarm ability by $s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_insectswarm.jpg",
		},
		[DruidPrimeGlyph.GlyphOfWrath]: {
			name: "Glyph of Wrath",
			description: "Increases the damage done by your Wrath by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_wrathv2.jpg",
		},
		[DruidPrimeGlyph.GlyphOfMoonfire]: {
			name: "Glyph of Moonfire",
			description: "Increases the periodic damage of your Moonfire ability by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg",
		},
		[DruidPrimeGlyph.GlyphOfBerserk]: {
			name: "Glyph of Berserk",
			description: "Increases the duration of Berserk by ${$62969m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_berserk.jpg",
		},
		[DruidPrimeGlyph.GlyphOfStarsurge]: {
			name: "Glyph of Starsurge",
			description: "When your Starsurge deals damage, the cooldown remaining on your Starfall is reduced by $s1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcane03.jpg",
		},
		[DruidPrimeGlyph.GlyphOfSavageRoar]: {
			name: "Glyph of Savage Roar",
			description: "Your Savage Roar ability grants an additional $63055s1% bonus damage done.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_skinteeth.jpg",
		},
		[DruidPrimeGlyph.GlyphOfLacerate]: {
			name: "Glyph of Lacerate",
			description: "Increases the critical strike chance of your Lacerate ability by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_lacerate.jpg",
		},
		[DruidPrimeGlyph.GlyphOfTigersFury]: {
			name: "Glyph of Tiger's Fury",
			description: "Reduces the cooldown of your Tiger's Fury ability by ${$m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mount_jungletiger.jpg",
		},
	},
	majorGlyphs: {
		[DruidMajorGlyph.GlyphOfFrenziedRegeneration]: {
			name: "Glyph of Frenzied Regeneration",
			description: "While Frenzied Regeneration is active, healing effects on you are $54810s1% more powerful but causes your Frenzied Regeneration to no longer convert rage into health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_bullrush.jpg",
		},
		[DruidMajorGlyph.GlyphOfMaul]: {
			name: "Glyph of Maul",
			description: "Your Maul ability now hits $s1 additional target for $s3% damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_maul.jpg",
		},
		[DruidMajorGlyph.GlyphOfSolarBeam]: {
			name: "Glyph of Solar Beam",
			description: "Increases the duration of your Solar Beam silence effect by $/1000;S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_vehicle_sonicshockwave.jpg",
		},
		[DruidMajorGlyph.GlyphOfPounce]: {
			name: "Glyph of Pounce",
			description: "Increases the range of your Pounce by $s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_supriseattack.jpg",
		},
		[DruidMajorGlyph.GlyphOfInnervate]: {
			name: "Glyph of Innervate",
			description: "When Innervate is cast on a friendly target other than the caster, the caster will gain ${$54833m1*10}% of maximum mana over $54833d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg",
		},
		[DruidMajorGlyph.GlyphOfRebirth]: {
			name: "Glyph of Rebirth",
			description: "Players resurrected by Rebirth are returned to life with $54733s2% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg",
		},
		[DruidMajorGlyph.GlyphOfHealingTouch]: {
			name: "Glyph of Healing Touch",
			description: "When you cast Healing Touch, the cooldown on your Nature's Swiftness is reduced by $s1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg",
		},
		[DruidMajorGlyph.GlyphOfHurricane]: {
			name: "Glyph of Hurricane",
			description: "Your Hurricane ability now also slows the movement speed of its victims by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg",
		},
		[DruidMajorGlyph.GlyphOfStarfall]: {
			name: "Glyph of Starfall",
			description: "Reduces the cooldown of Starfall by ${$54828m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_starfall.jpg",
		},
		[DruidMajorGlyph.GlyphOfEntanglingRoots]: {
			name: "Glyph of Entangling Roots",
			description: "Reduces the cast time of your Entangling Roots by 0.2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_stranglevines.jpg",
		},
		[DruidMajorGlyph.GlyphOfThorns]: {
			name: "Glyph of Thorns",
			description: "Reduces the cooldown of your Thorns spell by $/1000;S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_thorns.jpg",
		},
		[DruidMajorGlyph.GlyphOfFocus]: {
			name: "Glyph of Focus",
			description: "Increases the damage done by Starfall by $62080s2%, but decreases its radius by $62080s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcanepotency.jpg",
		},
		[DruidMajorGlyph.GlyphOfWildGrowth]: {
			name: "Glyph of Wild Growth",
			description: "Wild Growth can affect $62970s1 additional target, but its cooldown is increased by ${$62970m2/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_flourish.jpg",
		},
		[DruidMajorGlyph.GlyphOfMonsoon]: {
			name: "Glyph of Monsoon",
			description: "Reduces the cooldown of your Typhoon spell by ${$63056m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_riptide.jpg",
		},
		[DruidMajorGlyph.GlyphOfBarkskin]: {
			name: "Glyph of Barkskin",
			description: "Reduces the chance you'll be critically hit by $63057s1% while Barkskin is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg",
		},
		[DruidMajorGlyph.GlyphOfFerociousBite]: {
			name: "Glyph of Ferocious Bite",
			description: "Your Ferocious Bite ability heals you for $s1% of your maximum health for each $s2 energy used.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg",
		},
		[DruidMajorGlyph.GlyphOfFaerieFire]: {
			name: "Glyph of Faerie Fire",
			description: "Increases the range of your Faerie Fire and Feral Faerie Fire abilities by $s1 yds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_faeriefire.jpg",
		},
		[DruidMajorGlyph.GlyphOfFeralCharge]: {
			name: "Glyph of Feral Charge",
			description: "Reduces the cooldown of your Feral Charge (Cat) ability by ${$m1/-1000} sec and the cooldown of your Feral Charge (Bear) ability by ${$m2/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_pet_bear.jpg",
		},
	},
	minorGlyphs: {
		[DruidMinorGlyph.GlyphOfAquaticForm]: {
			name: "Glyph of Aquatic Form",
			description: "Increases your swim speed by $57856s1% while in Aquatic Form.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_aquaticform.jpg",
		},
		[DruidMinorGlyph.GlyphOfUnburdenedRebirth]: {
			name: "Glyph of Unburdened Rebirth",
			description: "Your Rebirth spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispsplodegreen.jpg",
		},
		[DruidMinorGlyph.GlyphOfChallengingRoar]: {
			name: "Glyph of Challenging Roar",
			description: "Reduces the cooldown of your Challenging Roar ability by ${$57858m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_challangingroar.jpg",
		},
		[DruidMinorGlyph.GlyphOfMarkOfTheWild]: {
			name: "Glyph of Mark of the Wild",
			description: "Mana cost of your Mark of the Wild reduced by $57855s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_giftofthewild.jpg",
		},
		[DruidMinorGlyph.GlyphOfDash]: {
			name: "Glyph of Dash",
			description: "Reduces the cooldown of your Dash ability by $59219s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_dash.jpg",
		},
		[DruidMinorGlyph.GlyphOfTyphoon]: {
			name: "Glyph of Typhoon",
			description: "Reduces the cost of your Typhoon spell by $62135s1% and increases its radius by $62135s2 yards, but it no longer knocks enemies back.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_typhoon.jpg",
		},
		[DruidMinorGlyph.GlyphOfTheTreant]: {
			name: "Glyph of the Treant",
			description: "Your Tree of Life Form now resembles a Treant.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_treeoflife.jpg",
		},
	},
};
