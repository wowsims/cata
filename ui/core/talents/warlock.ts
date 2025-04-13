import { WarlockMajorGlyph, WarlockMinorGlyph, WarlockPrimeGlyph, WarlockTalents } from '../proto/warlock.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarlockTalentJson from './trees/warlock.json';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);

export const warlockGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[WarlockPrimeGlyph.GlyphOfIncinerate]: {
			name: 'Glyph of Incinerate',
			description: 'Increases the damage done by Incinerate by $56242m1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfConflagrate]: {
			name: 'Glyph of Conflagrate',
			description: 'Reduces the cooldown of your Conflagrate by $/1000;S1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfCorruption]: {
			name: 'Glyph of Corruption',
			description:
				'Your Corruption spell has a $s1% chance to cause you to enter a Shadow Trance state after damaging the opponent.  The Shadow Trance state reduces the casting time of your next Shadow Bolt spell by $17941s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfBaneOfAgony]: {
			name: 'Glyph of Bane of Agony',
			description: 'Increases the duration of your Bane of Agony by ${$56241m1/1000} sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfFelguard]: {
			name: 'Glyph of Felguard',
			description: "Increases the damage done by your Felguard's Legion Strike by $s1%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfImmolate]: {
			name: 'Glyph of Immolate',
			description: 'Increases the periodic damage of your Immolate by $56228s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfImp]: {
			name: 'Glyph of Imp',
			description: "Increases the damage done by your Imp's Firebolt spell by $56248s1%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfShadowburn]: {
			name: 'Glyph of Shadowburn',
			description:
				"If your Shadowburn fails to kill the target at or below $s1% health, your Shadowburn's cooldown is instantly reset. This effect has a $91001d cooldown.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfUnstableAffliction]: {
			name: 'Glyph of Unstable Affliction',
			description: 'Decreases the casting time of your Unstable Affliction by ${$56233m1/-1000}.1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfHaunt]: {
			name: 'Glyph of Haunt',
			description: 'The bonus damage granted by your Haunt spell is increased by an additional $63302s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfMetamorphosis]: {
			name: 'Glyph of Metamorphosis',
			description: 'Increases the duration of your Metamorphosis by ${$63303m1/1000} sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfChaosBolt]: {
			name: 'Glyph of Chaos Bolt',
			description: 'Reduces the cooldown on Chaos Bolt by ${$63304m1/-1000} sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
		[WarlockPrimeGlyph.GlyphOfLashOfPain]: {
			name: 'Glyph of Lash of Pain',
			description: "Increases the damage done by your Succubus' Lash of Pain by $s1%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primewarlock.jpg',
		},
	},
	majorGlyphs: {
		[WarlockMajorGlyph.GlyphOfDeathCoilWl]: {
			name: 'Glyph of Death Coil',
			description: 'Increases the duration of your Death Coil by ${$56232m1/1000}.1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFear]: {
			name: 'Glyph of Fear',
			description: 'Your Fear causes the target to tremble in place instead of fleeing in fear, but now causes Fear to have a $s1 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfFelhunter]: {
			name: 'Glyph of Felhunter',
			description: 'When your Felhunter uses Devour Magic, you will also be healed for that amount.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHealthstone]: {
			name: 'Glyph of Healthstone',
			description: 'You receive $56224s1% more healing from using a healthstone.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfHowlOfTerror]: {
			name: 'Glyph of Howl of Terror',
			description: 'Reduces the cooldown on your Howl of Terror spell by ${$56217m1/-1000} sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulSwap]: {
			name: 'Glyph of Soul Swap',
			description:
				'Your Soul Swap leaves your damage-over-time spells behind on the target you Soul Swapped from, but gives Soul Swap a $s1 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowBolt]: {
			name: 'Glyph of Shadow Bolt',
			description: 'Reduces the mana cost of your Shadow Bolt by $56240s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulstone]: {
			name: 'Glyph of Soulstone',
			description: 'Increases the amount of health you gain from resurrecting via a Soulstone by an additional $56231s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSeduction]: {
			name: 'Glyph of Seduction',
			description: "Your Succubus's Seduction ability also removes all damage over time effects from the target.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfVoidwalker]: {
			name: 'Glyph of Voidwalker',
			description: "Increases your Voidwalker's total health by $56247s1%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfDemonicCircle]: {
			name: 'Glyph of Demonic Circle',
			description: 'Reduces the cooldown on Demonic Circle by ${$63309m1/-1000} sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfShadowflame]: {
			name: 'Glyph of Shadowflame',
			description: 'Your Shadowflame also applies a $63310s1% movement speed slow to its victims.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfLifeTap]: {
			name: 'Glyph of Life Tap',
			description: 'Reduces the global cooldown of your Life Tap by 0.5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
		[WarlockMajorGlyph.GlyphOfSoulLink]: {
			name: 'Glyph of Soul Link',
			description: 'Increases the percentage of damage shared via your Soul Link by an additional $63312s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorwarlock.jpg',
		},
	},
	minorGlyphs: {
		[WarlockMinorGlyph.GlyphOfHealthFunnel]: {
			name: 'Glyph of Health Funnel',
			description: 'Reduces the pushback suffered from damaging attacks while channeling your Health Funnel spell by $56238s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfUnendingBreath]: {
			name: 'Glyph of Unending Breath',
			description: 'Increases the swim speed of targets affected by your Unending Breath spell by $58079s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfDrainSoul]: {
			name: 'Glyph of Drain Soul',
			description: 'Your Drain Soul restores $58068s1% of your total mana after you kill a target that yields experience or honor.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfEyeOfKilrogg]: {
			name: 'Glyph of Eye of Kilrogg',
			description: 'Increases the movement speed of your Eye of Kilrogg by $s1% and allows it to fly in areas where flying mounts are enabled.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfCurseOfExhaustion]: {
			name: 'Glyph of Curse of Exhaustion',
			description: 'Increases the range of your Curse of Exhaustion spell by $58080s1 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfSubjugateDemon]: {
			name: 'Glyph of Subjugate Demon',
			description: 'Reduces the cast time of your Subjugate Demon spell by $58107s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
		[WarlockMinorGlyph.GlyphOfRitualOfSouls]: {
			name: 'Glyph of Ritual of Souls',
			description: 'Reduces the mana cost of your Ritual of Souls spell by $58094s1%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorwarlock.jpg',
		},
	},
};
