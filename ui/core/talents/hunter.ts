import { HunterMajorGlyph, HunterMinorGlyph, HunterPrimeGlyph, HunterTalents } from '../proto/hunter.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import HunterTalentJson from './trees/hunter.json';export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);

export const hunterGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[HunterPrimeGlyph.GlyphOfAimedShot]: {
			name: "Glyph of Aimed Shot",
			description: "When you critically hit with Aimed Shot, you instantly gain $s1 Focus.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_spear_07.jpg",
		},
		[HunterPrimeGlyph.GlyphOfArcaneShot]: {
			name: "Glyph of Arcane Shot",
			description: "Your Arcane Shot deals $s1% more damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_impalingbolt.jpg",
		},
		[HunterPrimeGlyph.GlyphOfTheDazzledPrey]: {
			name: "Glyph of the Dazzled Prey",
			description: "Your Steady Shot and Cobra Shot abilities generate an additional $s1 Focus on targets afflicted by a daze effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_cheatdeath.jpg",
		},
		[HunterPrimeGlyph.GlyphOfRapidFire]: {
			name: "Glyph of Rapid Fire",
			description: "Increases the haste from Rapid Fire by an additional $56828s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_runningshot.jpg",
		},
		[HunterPrimeGlyph.GlyphOfSerpentSting]: {
			name: "Glyph of Serpent Sting",
			description: "Increases the periodic critical strike chance of your Serpent Sting by $m1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_quickshot.jpg",
		},
		[HunterPrimeGlyph.GlyphOfSteadyShot]: {
			name: "Glyph of Steady Shot",
			description: "Increases the damage dealt by Steady Shot by $56826s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_steadyshot.jpg",
		},
		[HunterPrimeGlyph.GlyphOfKillCommand]: {
			name: "Glyph of Kill Command",
			description: "Reduces the Focus cost of your Kill Command by $s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_killcommand.jpg",
		},
		[HunterPrimeGlyph.GlyphOfChimeraShot]: {
			name: "Glyph of Chimera Shot",
			description: "Reduces the cooldown of Chimera Shot by ${$63065m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_chimerashot2.jpg",
		},
		[HunterPrimeGlyph.GlyphOfExplosiveShot]: {
			name: "Glyph of Explosive Shot",
			description: "Increases the critical strike chance of Explosive Shot by $63066s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_explosiveshot.jpg",
		},
		[HunterPrimeGlyph.GlyphOfKillShot]: {
			name: "Glyph of Kill Shot",
			description: "If the damage from your Kill Shot fails to kill a target at or below $s1% health, your Kill Shot's cooldown is instantly reset. This effect has a $90967d cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_assassinate2.jpg",
		},
	},
	majorGlyphs: {
		[HunterMajorGlyph.GlyphOfTrapLauncher]: {
			name: "Glyph of Trap Launcher",
			description: "Reduces the Focus cost of Trap Launcher by $s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_traplauncher.jpg",
		},
		[HunterMajorGlyph.GlyphOfMending]: {
			name: "Glyph of Mending",
			description: "Increases the total amount of healing done by your Mend Pet ability by $m2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_mendpet.jpg",
		},
		[HunterMajorGlyph.GlyphOfConcussiveShot]: {
			name: "Glyph of Concussive Shot",
			description: "Your Concussive Shot also limits the maximum run speed of your target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_stun.jpg",
		},
		[HunterMajorGlyph.GlyphOfBestialWrath]: {
			name: "Glyph of Bestial Wrath",
			description: "Decreases the cooldown of Bestial Wrath by ${$56830m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg",
		},
		[HunterMajorGlyph.GlyphOfDeterrence]: {
			name: "Glyph of Deterrence",
			description: "Decreases the cooldown of Deterrence by ${$56850m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_whirlwind.jpg",
		},
		[HunterMajorGlyph.GlyphOfDisengage]: {
			name: "Glyph of Disengage",
			description: "Decreases the cooldown of Disengage by ${$56844m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg",
		},
		[HunterMajorGlyph.GlyphOfFreezingTrap]: {
			name: "Glyph of Freezing Trap",
			description: "When your Freezing Trap breaks, the victim's movement speed is reduced by $61394s1% for $61394d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg",
		},
		[HunterMajorGlyph.GlyphOfIceTrap]: {
			name: "Glyph of Ice Trap",
			description: "Increases the radius of the effect from your Ice Trap by $56847s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg",
		},
		[HunterMajorGlyph.GlyphOfMisdirection]: {
			name: "Glyph of Misdirection",
			description: "When you use Misdirection on your pet, the cooldown on your Misdirection is reset.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_misdirection.jpg",
		},
		[HunterMajorGlyph.GlyphOfImmolationTrap]: {
			name: "Glyph of Immolation Trap",
			description: "Decreases the duration of the effect from your Immolation Trap by 6 sec, but damage while active is increased by $56846s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg",
		},
		[HunterMajorGlyph.GlyphOfSilencingShot]: {
			name: "Glyph of Silencing Shot",
			description: "When you successfully silence an enemy's spell cast with Silencing Shot, you instantly gain $s1 Focus.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_theblackarrow.jpg",
		},
		[HunterMajorGlyph.GlyphOfSnakeTrap]: {
			name: "Glyph of Snake Trap",
			description: "Snakes from your Snake Trap take $56849s1% reduced damage from area of effect spells.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_snaketrap.jpg",
		},
		[HunterMajorGlyph.GlyphOfWyvernSting]: {
			name: "Glyph of Wyvern Sting",
			description: "Decreases the cooldown of your Wyvern Sting by ${$m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_spear_02.jpg",
		},
		[HunterMajorGlyph.GlyphOfMastersCall]: {
			name: "Glyph of Master's Call",
			description: "Increases the duration of your Master's Call by $/1000;S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_masterscall.jpg",
		},
		[HunterMajorGlyph.GlyphOfScatterShot]: {
			name: "Glyph of Scatter Shot",
			description: "Increases the range of Scatter Shot by $63069s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_golemstormbolt.jpg",
		},
		[HunterMajorGlyph.GlyphOfRaptorStrike]: {
			name: "Glyph of Raptor Strike",
			description: "Reduces damage taken by $63087s1% for $63087d after using Raptor Strike.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_meleedamage.jpg",
		},
	},
	minorGlyphs: {
		[HunterMinorGlyph.GlyphOfRevivePet]: {
			name: "Glyph of Revive Pet",
			description: "Reduces the pushback suffered from damaging attacks while casting Revive Pet by $57866s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_beastsoothe.jpg",
		},
		[HunterMinorGlyph.GlyphOfLesserProportion]: {
			name: "Glyph of Lesser Proportion",
			description: "Slightly reduces the size of your pet.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_bestialdiscipline.jpg",
		},
		[HunterMinorGlyph.GlyphOfFeignDeath]: {
			name: "Glyph of Feign Death",
			description: "Reduces the cooldown of your Feign Death spell by ${$57903m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feigndeath.jpg",
		},
		[HunterMinorGlyph.GlyphOfAspectOfThePack]: {
			name: "Glyph of Aspect of the Pack",
			description: "Increases the range of your Aspect of the Pack ability by $57904s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mount_whitetiger.jpg",
		},
		[HunterMinorGlyph.GlyphOfScareBeast]: {
			name: "Glyph of Scare Beast",
			description: "Reduces the pushback suffered from damaging attacks while casting Scare Beast by $57902s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_druid_cower.jpg",
		},
	},
};
