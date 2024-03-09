import { WarlockTalents, WarlockMajorGlyph, WarlockMinorGlyph } from '../proto/warlock.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import WarlockTalentJson from './trees/warlock.json';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);


export const warlockGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[WarlockPrimeGlyph.GlyphOfBaneOfAgony]: {
			name: 'Glyph of Bane of Agony',
			description: 'Increases the duration of your Bane of Agony by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curseofsargeras.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfChaosBolt]: {
			name: 'Glyph of Chaos Bolt',
			description: 'Reduces the cooldown on Chaos Bolt by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_chaosbolt.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfConflagrate]: {
			name: 'Glyph of Conflagrate',
			description: 'Reduces the cooldown of your Conflagrate by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfCorruption]: {
			name: 'Glyph of Corruption',
			description: 'Your Corruption spell has a 4% chance to cause you to enter a Shadow Trance state after damaging the opponent. The Shadow Trance state reduces the casting time of your next Shadow Bolt spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_abominationexplosion.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfFelguard]: {
			name: 'Glyph of Felguard',
			description: 'Increases the damage done by your Felguard\'s Legion Strike by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfHaunt]: {
			name: 'Glyph of Haunt',
			description: 'The bonus damage granted by your Haunt spell is increased by an additional 3%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_haunt.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfImmolate]: {
			name: 'Glyph of Immolate',
			description: 'Increases the periodic damage of your Immolate by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_immolation.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfImp]: {
			name: 'Glyph of Imp',
			description: 'Increases the damage done by your Imp\'s Firebolt spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonimp.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfIncinerate]: {
			name: 'Glyph of Incinerate',
			description: 'Increases the damage done by Incinerate by 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_burnout.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfLashOfPain]: {
			name: 'Glyph of Lash of Pain',
			description: 'Increases the damage done by your Succubus\' Lash of Pain by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_curse.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfMetamorphosis]: {
			name: 'Glyph of Metamorphosis',
			description: 'Increases the duration of your Metamorphosis by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonform.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfShadowburn]: {
			name: 'Glyph of Shadowburn',
			description: 'If your Shadowburn fails to kill the target at or below 20% health, your Shadowburn\'s cooldown is instantly reset. This effect has a 6 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_scourgebuild.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfUnstableAffliction]: {
			name: 'Glyph of Unstable Affliction',
			description: 'Decreases the casting time of your Unstable Affliction by 0.2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unstableaffliction_3.jpg',
		},
	},
	majorGlyphs: {
		[WarlockMajorGlyph.GlyphOfDeathCoil]: {
			name: 'Glyph of Death Coil',
			description: 'Increases the duration of your Death Coil by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDemonicCircle]: {
			name: 'Glyph of Demonic Circle',
			description: 'Reduces the cooldown on Demonic Circle by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demoniccircleteleport.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFear]: {
			name: 'Glyph of Fear',
			description: 'Your Fear causes the target to tremble in place instead of fleeing in fear, but now causes Fear to have a 5 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_possession.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelhunter]: {
			name: 'Glyph of Felhunter',
			description: 'When your Felhunter uses Devour Magic, you will also be healed for that amount.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthstone]: {
			name: 'Glyph of Healthstone',
			description: 'You receive 30% more healing from using a healthstone.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_stone_04.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHowlOfTerror]: {
			name: 'Glyph of Howl of Terror',
			description: 'Reduces the cooldown on your Howl of Terror spell by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathscream.jpg',
		},
		[WarlockMajorGlyph.GlyphOfLifeTap]: {
			name: 'Glyph of Life Tap',
			description: 'Reduces the global cooldown of your Life Tap by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_burningspirit.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSeduction]: {
			name: 'Glyph of Seduction',
			description: 'Your Succubus\'s Seduction ability also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_mindsteal.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowBolt]: {
			name: 'Glyph of Shadow Bolt',
			description: 'Reduces the mana cost of your Shadow Bolt by 15%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowbolt.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowflame]: {
			name: 'Glyph of Shadowflame',
			description: 'Your Shadowflame also applies a 70% movement speed slow to its victims.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_shadowflame.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulLink]: {
			name: 'Glyph of Soul Link',
			description: 'Increases the percentage of damage shared via your Soul Link by an additional 5%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_gathershadows.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulSwap]: {
			name: 'Glyph of Soul Swap',
			description: 'Your Soul Swap leaves your damage-over-time spells behind on the target you Soul Swapped from, but gives Soul Swap a 30 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_soulswap.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulstone]: {
			name: 'Glyph of Soulstone',
			description: 'Increases the amount of health you gain from resurrecting via a Soulstone by an additional 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_orb_04.jpg',
		},
		[WarlockMajorGlyph.GlyphOfVoidwalker]: {
			name: 'Glyph of Voidwalker',
			description: 'Increases your Voidwalker\'s total health by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg',
		},
	},
	minorGlyphs: {
		[WarlockMinorGlyph.GlyphOfCurseOfExhaustion]: {
			name: 'Glyph of Curse of Exhaustion',
			description: 'Increases the range of your Curse of Exhaustion spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_grimward.jpg',
		},
		[WarlockMinorGlyph.GlyphOfDrainSoul]: {
			name: 'Glyph of Drain Soul',
			description: 'Your Drain Soul restores 10% of your total mana after you kill a target that yields experience or honor.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_haunting.jpg',
		},
		[WarlockMinorGlyph.GlyphOfEyeOfKilrogg]: {
			name: 'Glyph of Eye of Kilrogg',
			description: 'Increases the movement speed of your Eye of Kilrogg by 50% and allows it to fly in areas where flying mounts are enabled.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_evileye.jpg',
		},
		[WarlockMinorGlyph.GlyphOfHealthFunnel]: {
			name: 'Glyph of Health Funnel',
			description: 'Reduces the pushback suffered from damaging attacks while channeling your Health Funnel spell by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
		[WarlockMinorGlyph.GlyphOfRitualOfSouls]: {
			name: 'Glyph of Ritual of Souls',
			description: 'Reduces the mana cost of your Ritual of Souls spell by 70%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadesofdarkness.jpg',
		},
		[WarlockMinorGlyph.GlyphOfSubjugateDemon]: {
			name: 'Glyph of Subjugate Demon',
			description: 'Reduces the cast time of your Subjugate Demon spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_enslavedemon.jpg',
		},
		[WarlockMinorGlyph.GlyphOfUnendingBreath]: {
			name: 'Glyph of Unending Breath',
			description: 'Increases the swim speed of targets affected by your Unending Breath spell by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonbreath.jpg',
		},
	},
};