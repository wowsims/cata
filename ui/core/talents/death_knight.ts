import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightPrimeGlyph, DeathKnightTalents } from '../proto/death_knight.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DeathKnightTalentJson from './trees/death_knight.json';export const deathKnightTalentsConfig: TalentsConfig<DeathKnightTalents> = newTalentsConfig(DeathKnightTalentJson);

export const deathKnightGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[DeathKnightPrimeGlyph.GlyphOfHeartStrike]: {
			name: "Glyph of Heart Strike",
			description: "Increases the damage of your Heart Strike ability by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfDeathAndDecay]: {
			name: "Glyph of Death and Decay",
			description: "Increases the duration of your Death and Decay spell by $58629s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfFrostStrike]: {
			name: "Glyph of Frost Strike",
			description: "Reduces the cost of your Frost Strike by ${$58647m1/-10} Runic Power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfIcyTouch]: {
			name: "Glyph of Icy Touch",
			description: "Your Frost Fever disease deals $58631s1% additional damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfObliterate]: {
			name: "Glyph of Obliterate",
			description: "Increases the damage of your Obliterate ability by $m3%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfRaiseDead]: {
			name: "Glyph of Raise Dead",
			description: "Your Ghoul receives an additional $58686s1% of your Strength and $58686s1% of your Stamina.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfRuneStrike]: {
			name: "Glyph of Rune Strike",
			description: "Increases the critical strike chance of your Rune Strike by $58669s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfScourgeStrike]: {
			name: "Glyph of Scourge Strike",
			description: "Increases the Shadow damage portion of your Scourge Strike by $58642s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfDeathStrike]: {
			name: "Glyph of Death Strike",
			description: "Increases your Death Strike's damage by $59336s1% for every 5 Runic Power you currently have (up to a maximum of $59336s2%).  The Runic Power is not consumed by this effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.DeprecatedGlyphOfTheGhoul]: {
			name: "DEPRECATED Glyph of the Ghoul",
			description: "Your Ghoul receives an additional $58686s1% of your Strength and $58686s1% of your Stamina.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_inscription_minorglyph13.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfDeathCoil]: {
			name: "Glyph of Death Coil",
			description: "Increases the damage or healing done by Death Coil by $63333s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
		[DeathKnightPrimeGlyph.GlyphOfHowlingBlast]: {
			name: "Glyph of Howling Blast",
			description: "Your Howling Blast ability now infects your targets with Frost Fever.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primedeathknight.jpg",
		},
	},
	majorGlyphs: {
		[DeathKnightMajorGlyph.GlyphOfAntiMagicShell]: {
			name: "Glyph of Anti-Magic Shell",
			description: "Increases the duration of your Anti-Magic Shell by ${$58623m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfBoneShield]: {
			name: "Glyph of Bone Shield",
			description: "Increases your movement speed by $58673s1% while Bone Shield is active.  This does not stack with other movement-speed increasing effects.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfChainsOfIce]: {
			name: "Glyph of Chains of Ice",
			description: "Your Chains of Ice also causes $s1 Frost damage, with additional damage depending on your attack power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDeathGrip]: {
			name: "Glyph of Death Grip",
			description: "Increases the range of your Death Grip ability by $62259s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfPestilence]: {
			name: "Glyph of Pestilence",
			description: "Increases the radius of your Pestilence effect by $58657s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfStrangulate]: {
			name: "Glyph of Strangulate",
			description: "Increases the Silence duration of your Strangulate ability by ${$58618m1/1000} sec when used on a target who is casting a spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfPillarOfFrost]: {
			name: "Glyph of Pillar of Frost",
			description: "Empowers your Pillar of Frost, making you immune to all effects that cause loss of control of your character, but also freezing you in place while the ability is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfVampiricBlood]: {
			name: "Glyph of Vampiric Blood",
			description: "Increases the bonus healing received while your Vampiric Blood is active by an additional $58676s1%, but your Vampiric Blood no longer grants you health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfRuneTap]: {
			name: "Glyph of Rune Tap",
			description: "Your Rune Tap also heals your party for $59754s1% of their maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfBloodBoil]: {
			name: "Glyph of Blood Boil",
			description: "Increases the radius of your Blood Boil ability by $59332s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon]: {
			name: "Glyph of Dancing Rune Weapon",
			description: "Increases your threat generation by $63330s1% while your Dancing Rune Weapon is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfHungeringCold]: {
			name: "Glyph of Hungering Cold",
			description: "Your Hungering Cold ability no longer costs runic power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDarkSuccor]: {
			name: "Glyph of Dark Succor",
			description: "Your next Death Strike performed while in Frost or Unholy Presence, within $101568d after killing an enemy that yields experience or honor, will restore at least $101568s1% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majordeathknight.jpg",
		},
	},
	minorGlyphs: {
		[DeathKnightMinorGlyph.GlyphOfBloodTap]: {
			name: "Glyph of Blood Tap",
			description: "Your Blood Tap no longer causes damage to you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfDeathsEmbrace]: {
			name: "Glyph of Death's Embrace",
			description: "Your Death Coil refunds $s1 Runic Power when used to heal an allied minion.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfHornOfWinter]: {
			name: "Glyph of Horn of Winter",
			description: "Increases the duration of your Horn of Winter ability by ${$58680m1/60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfPathOfFrost]: {
			name: "Glyph of Path of Frost",
			description: "Your Path of Frost ability allows you to fall from a greater distance without suffering damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfResilientGrip]: {
			name: "Glyph of Resilient Grip",
			description: "When your Death Grip ability fails because its target is immune, its cooldown is reset.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfDeathGate]: {
			name: "Glyph of Death Gate",
			description: "Reduces the cast time of your Death Gate spell by $60200s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minordeathknight.jpg",
		},
	},
};
