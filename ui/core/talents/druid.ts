import {DruidMajorGlyph, DruidMinorGlyph,DruidPrimeGlyph, DruidTalents } from '../proto/druid.js';
import { GlyphsConfig, } from './glyphs_picker.js';
import { newTalentsConfig,TalentsConfig } from './talents_picker.js';
import DruidTalentsJson from './trees/druid.json';

export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentsJson);

export const druidGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[DruidPrimeGlyph.GlyphOfBerserk]: {
			name: 'Glyph of Berserk',
			description: 'Increases the duration of Berserk by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_berserk.jpg',
		},
		[DruidPrimeGlyph.GlyphOfBloodletting]: {
			name: 'Glyph of Bloodletting',
			description: 'Each time you Shred or Mangle in Cat Form, the duration of your Rip on the target is extended by 2 sec, up to a maximum of 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_vampiricaura.jpg',
		},
		[DruidPrimeGlyph.GlyphOfInsectSwarm]: {
			name: 'Glyph of Insect Swarm',
			description: 'Increases the damage of your Insect Swarm ability by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_insectswarm.jpg',
		},
		[DruidPrimeGlyph.GlyphOfLacerate]: {
			name: 'Glyph of Lacerate',
			description: 'Increases the critical strike chance of your Lacerate ability by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_lacerate.jpg',
		},
		[DruidPrimeGlyph.GlyphOfLifebloom]: {
			name: 'Glyph of Lifebloom',
			description: 'Increases the critical effect chance of your Lifebloom by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_herb_felblossom.jpg',
		},
		[DruidPrimeGlyph.GlyphOfMangle]: {
			name: 'Glyph of Mangle',
			description: 'Increases the damage done by Mangle by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_mangle2.jpg',
		},
		[DruidPrimeGlyph.GlyphOfMoonfire]: {
			name: 'Glyph of Moonfire',
			description: 'Increases the periodic damage of your Moonfire ability by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg',
		},
		[DruidPrimeGlyph.GlyphOfRegrowth]: {
			name: 'Glyph of Regrowth',
			description: 'Your Regrowth heal-over-time will automatically refresh its duration on targets at or below 50% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_resistnature.jpg',
		},
		[DruidPrimeGlyph.GlyphOfRejuvenation]: {
			name: 'Glyph of Rejuvenation',
			description: 'Increases the healing done by your Rejuvenation by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_rejuvenation.jpg',
		},
		[DruidPrimeGlyph.GlyphOfRip]: {
			name: 'Glyph of Rip',
			description: 'Increases the periodic damage of your Rip by 15%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_ghoulfrenzy.jpg',
		},
		[DruidPrimeGlyph.GlyphOfSavageRoar]: {
			name: 'Glyph of Savage Roar',
			description: 'Your Savage Roar ability grants an additional 5% bonus damage done.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_skinteeth.jpg',
		},
		[DruidPrimeGlyph.GlyphOfStarfire]: {
			name: 'Glyph of Starfire',
			description: 'Your Starfire ability increases the duration of your Moonfire effect on the target by 3 sec, up to a maximum of 9 additional seconds. Only functions on the target with your most recently applied Moonfire.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_starfire.jpg',
		},
		[DruidPrimeGlyph.GlyphOfStarsurge]: {
			name: 'Glyph of Starsurge',
			description: 'When your Starsurge deals damage, the cooldown remaining on your Starfall is reduced by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcane03.jpg',
		},
		[DruidPrimeGlyph.GlyphOfSwiftmend]: {
			name: 'Glyph of Swiftmend',
			description: 'Your Swiftmend ability no longer consumes a Rejuvenation or Regrowth effect from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_relics_idolofrejuvenation.jpg',
		},
		[DruidPrimeGlyph.GlyphOfTigersFury]: {
			name: 'Glyph of Tiger\'s Fury',
			description: 'Reduces the cooldown of your Tiger\'s Fury ability by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mount_jungletiger.jpg',
		},
		[DruidPrimeGlyph.GlyphOfWrath]: {
			name: 'Glyph of Wrath',
			description: 'Increases the damage done by your Wrath by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_wrathv2.jpg',
		},
	},
	majorGlyphs: {
		[DruidMajorGlyph.GlyphOfBarkskin]: {
			name: 'Glyph of Barkskin',
			description: 'Reduces the chance you\'ll be critically hit by 25% while Barkskin is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg',
		},
		[DruidMajorGlyph.GlyphOfEntanglingRoots]: {
			name: 'Glyph of Entangling Roots',
			description: 'Reduces the cast time of your Entangling Roots by 0.2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stranglevines.jpg',
		},
		[DruidMajorGlyph.GlyphOfFaerieFire]: {
			name: 'Glyph of Faerie Fire',
			description: 'Increases the range of your Faerie Fire and Feral Faerie Fire abilities by 10 yds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_faeriefire.jpg',
		},
		[DruidMajorGlyph.GlyphOfFeralCharge]: {
			name: 'Glyph of Feral Charge',
			description: 'Reduces the cooldown of your Feral Charge (Cat) ability by 2 sec and the cooldown of your Feral Charge (Bear) ability by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_pet_bear.jpg',
		},
		[DruidMajorGlyph.GlyphOfFerociousBite]: {
			name: 'Glyph of Ferocious Bite',
			description: 'Your Ferocious Bite ability heals you for 1% of your maximum health for each 10 energy used.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg',
		},
		[DruidMajorGlyph.GlyphOfFocus]: {
			name: 'Glyph of Focus',
			description: 'Increases the damage done by Starfall by 10%, but decreases its radius by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcanepotency.jpg',
		},
		[DruidMajorGlyph.GlyphOfFrenziedRegeneration]: {
			name: 'Glyph of Frenzied Regeneration',
			description: 'While Frenzied Regeneration is active, healing effects on you are 30% more powerful but causes your Frenzied Regeneration to no longer convert rage into health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_bullrush.jpg',
		},
		[DruidMajorGlyph.GlyphOfHealingTouch]: {
			name: 'Glyph of Healing Touch',
			description: 'When you cast Healing Touch, the cooldown on your Nature\'s Swiftness is reduced by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingtouch.jpg',
		},
		[DruidMajorGlyph.GlyphOfHurricane]: {
			name: 'Glyph of Hurricane',
			description: 'Your Hurricane ability now also slows the movement speed of its victims by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_cyclone.jpg',
		},
		[DruidMajorGlyph.GlyphOfInnervate]: {
			name: 'Glyph of Innervate',
			description: 'When Innervate is cast on a friendly target other than the caster, the caster will gain 10% of maximum mana over 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[DruidMajorGlyph.GlyphOfMaul]: {
			name: 'Glyph of Maul',
			description: 'Your Maul ability now hits 1 additional target for 50% damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_maul.jpg',
		},
		[DruidMajorGlyph.GlyphOfMonsoon]: {
			name: 'Glyph of Monsoon',
			description: 'Reduces the cooldown of your Typhoon spell by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_riptide.jpg',
		},
		[DruidMajorGlyph.GlyphOfOmenOfClarity]: {
			name: 'Glyph of Omen of Clarity',
			description: 'Your Omen of Clarity talent has a 100% chance to be triggered by successfully casting Faerie Fire (Feral). Does not trigger on players or player-controlled pets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg',
		},
		[DruidMajorGlyph.GlyphOfPounce]: {
			name: 'Glyph of Pounce',
			description: 'Increases the range of your Pounce by 3 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_supriseattack.jpg',
		},
		[DruidMajorGlyph.GlyphOfRebirth]: {
			name: 'Glyph of Rebirth',
			description: 'Players resurrected by Rebirth are returned to life with 100% health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_reincarnation.jpg',
		},
		[DruidMajorGlyph.GlyphOfSolarBeam]: {
			name: 'Glyph of Solar Beam',
			description: 'Increases the duration of your Solar Beam silence effect by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_vehicle_sonicshockwave.jpg',
		},
		[DruidMajorGlyph.GlyphOfStarfall]: {
			name: 'Glyph of Starfall',
			description: 'Reduces the cooldown of Starfall by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_starfall.jpg',
		},
		[DruidMajorGlyph.GlyphOfThorns]: {
			name: 'Glyph of Thorns',
			description: 'Reduces the cooldown of your Thorns spell by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thorns.jpg',
		},
		[DruidMajorGlyph.GlyphOfWildGrowth]: {
			name: 'Glyph of Wild Growth',
			description: 'Wild Growth can affect 1 additional target, but its cooldown is increased by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_flourish.jpg',
		},
	},
	minorGlyphs: {
		[DruidMinorGlyph.GlyphOfAquaticForm]: {
			name: 'Glyph of Aquatic Form',
			description: 'Increases your swim speed by 50% while in Aquatic Form.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_aquaticform.jpg',
		},
		[DruidMinorGlyph.GlyphOfChallengingRoar]: {
			name: 'Glyph of Challenging Roar',
			description: 'Reduces the cooldown of your Challenging Roar ability by 30 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_challangingroar.jpg',
		},
		[DruidMinorGlyph.GlyphOfDash]: {
			name: 'Glyph of Dash',
			description: 'Reduces the cooldown of your Dash ability by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_dash.jpg',
		},
		[DruidMinorGlyph.GlyphOfMarkOfTheWild]: {
			name: 'Glyph of Mark of the Wild',
			description: 'Mana cost of your Mark of the Wild reduced by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_giftofthewild.jpg',
		},
		[DruidMinorGlyph.GlyphOfTheTreant]: {
			name: 'Glyph of the Treant',
			description: 'Your Tree of Life Form now resembles a Treant.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_treeoflife.jpg',
		},
		[DruidMinorGlyph.GlyphOfTyphoon]: {
			name: 'Glyph of Typhoon',
			description: 'Reduces the cost of your Typhoon spell by 8% and increases its radius by 10 yards, but it no longer knocks enemies back.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_typhoon.jpg',
		},
		[DruidMinorGlyph.GlyphOfUnburdenedRebirth]: {
			name: 'Glyph of Unburdened Rebirth',
			description: 'Your Rebirth spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispsplodegreen.jpg',
		},
	},
};
