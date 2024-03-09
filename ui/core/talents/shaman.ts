import { ShamanTalents, ShamanMajorGlyph, ShamanMinorGlyph, ShamanPrimeGlyph } from '../proto/shaman.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import ShamanTalentJson from './trees/shaman.json';

export const shamanTalentsConfig: TalentsConfig<ShamanTalents> = newTalentsConfig(ShamanTalentJson);


export const shamanGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[ShamanPrimeGlyph.GlyphOfEarthShield]: {
			name: 'Glyph of Earth Shield',
			description: 'Increases the amount healed by your Earth Shield by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_skinofearth.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfEarthlivingWeapon]: {
			name: 'Glyph of Earthliving Weapon',
			description: 'Increases the effectiveness of your Earthliving weapon\'s periodic healing by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_unleashweapon_life.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfFeralSpirit]: {
			name: 'Glyph of Feral Spirit',
			description: 'Your spirit wolves gain an additional 30% of your attack power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_feralspirit.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfFireElementalTotem]: {
			name: 'Glyph of Fire Elemental Totem',
			description: 'Reduces the cooldown of your Fire Elemental Totem by 5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_elemental_totem.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfFlameShock]: {
			name: 'Glyph of Flame Shock',
			description: 'Increases the duration of your Flame Shock by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfFlametongueWeapon]: {
			name: 'Glyph of Flametongue Weapon',
			description: 'Increases your spell critical strike chance by 2% on each of your weapons with Flametongue Weapon active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_unleashweapon_flame.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfLavaBurst]: {
			name: 'Glyph of Lava Burst',
			description: 'Your Lava Burst spell deals 10% more damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_lavaburst.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfLavaLash]: {
			name: 'Glyph of Lava Lash',
			description: 'Increases the damage dealt by your Lava Lash ability by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_lavalash.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfLightningBolt]: {
			name: 'Glyph of Lightning Bolt',
			description: 'Increases the damage dealt by Lightning Bolt by 4%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfRiptide]: {
			name: 'Glyph of Riptide',
			description: 'Increases the duration of Riptide by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_riptide.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfShocking]: {
			name: 'Glyph of Shocking',
			description: 'Reduces your global cooldown when casting your shock spells by 0.0 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_earthshock.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfStormstrike]: {
			name: 'Glyph of Stormstrike',
			description: 'Increases the critical strike chance bonus from your Stormstrike ability by an additional 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_stormstrike.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfUnleashedLightning]: {
			name: 'Glyph of Unleashed Lightning',
			description: 'Allows Lightning Bolt to be cast while moving.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfWaterShield]: {
			name: 'Glyph of Water Shield',
			description: 'Increases the passive mana regeneration of your Water Shield spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_shaman_watershield.jpg',
		},
		[ShamanPrimeGlyph.GlyphOfWindfuryWeapon]: {
			name: 'Glyph of Windfury Weapon',
			description: 'Increases the chance per swing for Windfury Weapon to trigger by 2%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_unleashweapon_wind.jpg',
		},
	},
	majorGlyphs: {
		[ShamanMajorGlyph.GlyphOfChainHeal]: {
			name: 'Glyph of Chain Heal',
			description: 'Increases healing done by your Chain Heal spell to targets beyond the first by 15%, but decreases the amount received by the initial target by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg',
		},
		[ShamanMajorGlyph.GlyphOfChainLightning]: {
			name: 'Glyph of Chain Lightning',
			description: 'Your Chain Lightning spell now strikes 2 additional targets, but deals 10% less initial damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_chainlightning.jpg',
		},
		[ShamanMajorGlyph.GlyphOfElementalMastery]: {
			name: 'Glyph of Elemental Mastery',
			description: 'While your Elemental Mastery ability is active, you take 20% less damage from all sources.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispheal.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFireNova]: {
			name: 'Glyph of Fire Nova',
			description: 'Increases the radius of your Fire Nova spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_firenova.jpg',
		},
		[ShamanMajorGlyph.GlyphOfFrostShock]: {
			name: 'Glyph of Frost Shock',
			description: 'Increases the duration of your Frost Shock by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostshock.jpg',
		},
		[ShamanMajorGlyph.GlyphOfGhostWolf]: {
			name: 'Glyph of Ghost Wolf',
			description: 'Your Ghost Wolf form grants an additional 5% movement speed.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_spiritwolf.jpg',
		},
		[ShamanMajorGlyph.GlyphOfGroundingTotem]: {
			name: 'Glyph of Grounding Totem',
			description: 'Instead of absorbing a spell, your Grounding Totem reflects the next harmful spell back at its caster, but the cooldown of your Grounding Totem is increased by 35 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_groundingtotem.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHealingStreamTotem]: {
			name: 'Glyph of Healing Stream Totem',
			description: 'Your Healing Stream Totem also increases the Fire, Frost, and Nature resistance of party and raid members within 30 yards by 85.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_04.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHealingWave]: {
			name: 'Glyph of Healing Wave',
			description: 'Your Healing Wave also heals you for 20% of the healing effect when you heal someone else.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_healingwavegreater.jpg',
		},
		[ShamanMajorGlyph.GlyphOfHex]: {
			name: 'Glyph of Hex',
			description: 'Reduces the cooldown of your Hex spell by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_hex.jpg',
		},
		[ShamanMajorGlyph.GlyphOfLightningShield]: {
			name: 'Glyph of Lightning Shield',
			description: 'Your Lightning Shield can no longer drop below 3 charges from dealing damage to attackers.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightningshield.jpg',
		},
		[ShamanMajorGlyph.GlyphOfShamanisticRage]: {
			name: 'Glyph of Shamanistic Rage',
			description: 'Activating your Shamanistic Rage ability also cleanses you of all dispellable Magic debuffs.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_shamanrage.jpg',
		},
		[ShamanMajorGlyph.GlyphOfStoneclawTotem]: {
			name: 'Glyph of Stoneclaw Totem',
			description: 'Your Stoneclaw Totem also places a damage absorb shield on you, equal to 4 times the strength of the shield it places on your totems.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_stoneclawtotem.jpg',
		},
		[ShamanMajorGlyph.GlyphOfThunder]: {
			name: 'Glyph of Thunder',
			description: 'Reduces the cooldown on Thunderstorm by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg',
		},
		[ShamanMajorGlyph.GlyphOfTotemicRecall]: {
			name: 'Glyph of Totemic Recall',
			description: 'Causes your Totemic Recall ability to return an additional 50% of the mana cost of any recalled totems.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_totemrecall.jpg',
		},
	},
	minorGlyphs: {
		[ShamanMinorGlyph.GlyphOfAstralRecall]: {
			name: 'Glyph of Astral Recall',
			description: 'Reduces the cooldown of your Astral Recall spell by 7.5 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_astralrecal.jpg',
		},
		[ShamanMinorGlyph.GlyphOfRenewedLife]: {
			name: 'Glyph of Renewed Life',
			description: 'Your Reincarnation spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_improvedreincarnation.jpg',
		},
		[ShamanMinorGlyph.GlyphOfTheArcticWolf]: {
			name: 'Glyph of the Arctic Wolf',
			description: 'Alters the appearance of your Ghost Wolf transformation, causing it to resemble an arctic wolf.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mount_whitedirewolf.jpg',
		},
		[ShamanMinorGlyph.GlyphOfThunderstorm]: {
			name: 'Glyph of Thunderstorm',
			description: 'Increases the mana you receive from your Thunderstorm spell by 2%, but it no longer knocks enemies back or reduces their movement speed.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shaman_thunderstorm.jpg',
		},
		[ShamanMinorGlyph.GlyphOfWaterBreathing]: {
			name: 'Glyph of Water Breathing',
			description: 'Your Water Breathing spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg',
		},
		[ShamanMinorGlyph.GlyphOfWaterWalking]: {
			name: 'Glyph of Water Walking',
			description: 'Your Water Walking spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_windwalkon.jpg',
		},
	},
};