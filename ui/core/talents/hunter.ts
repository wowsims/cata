import { HunterTalents, HunterMajorGlyph, HunterMinorGlyph, HunterPetTalents, HunterPrimeGlyph } from '../proto/hunter.js';

import { GlyphsConfig } from './glyphs_picker.js';
import { TalentsConfig, newTalentsConfig } from './talents_picker.js';

import HunterTalentJson from './trees/hunter.json';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);

export const hunterGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[HunterPrimeGlyph.GlyphOfAimedShot]: {
			name: 'Glyph of Aimed Shot',
			description: 'When you critically hit with Aimed Shot, you instantly gain 5 Focus.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_07.jpg',
		},
		[HunterPrimeGlyph.GlyphOfArcaneShot]: {
			name: 'Glyph of Arcane Shot',
			description: 'Your Arcane Shot deals 12% more damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_impalingbolt.jpg',
		},
		[HunterPrimeGlyph.GlyphOfChimeraShot]: {
			name: 'Glyph of Chimera Shot',
			description: 'Reduces the cooldown of Chimera Shot by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_chimerashot2.jpg',
		},
		[HunterPrimeGlyph.GlyphOfDazzledPrey]: {
			name: 'Glyph of Dazzled Prey',
			description: 'Your Steady Shot and Cobra Shot abilities generate an additional 2 Focus on targets afflicted by a daze effect.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_cheatdeath.jpg',
		},
		[HunterPrimeGlyph.GlyphOfExplosiveShot]: {
			name: 'Glyph of Explosive Shot',
			description: 'Increases the critical strike chance of Explosive Shot by 6%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_explosiveshot.jpg',
		},
		[HunterPrimeGlyph.GlyphOfKillCommand]: {
			name: 'Glyph of Kill Command',
			description: 'Reduces the Focus cost of your Kill Command by 3.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_killcommand.jpg',
		},
		[HunterPrimeGlyph.GlyphOfKillShot]: {
			name: 'Glyph of Kill Shot',
			description: 'If the damage from your Kill Shot fails to kill a target at or below 20% health, your Kill Shot\'s cooldown is instantly reset. This effect has a 6 sec cooldown.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_assassinate2.jpg',
		},
		[HunterPrimeGlyph.GlyphOfRapidFire]: {
			name: 'Glyph of Rapid Fire',
			description: 'Increases the haste from Rapid Fire by an additional 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_runningshot.jpg',
		},
		[HunterPrimeGlyph.GlyphOfSerpentSting]: {
			name: 'Glyph of Serpent Sting',
			description: 'Increases the periodic critical strike chance of your Serpent Sting by 6%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_quickshot.jpg',
		},
		[HunterPrimeGlyph.GlyphOfSteadyShot]: {
			name: 'Glyph of Steady Shot',
			description: 'Increases the damage dealt by Steady Shot by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_steadyshot.jpg',
		},
	},
	majorGlyphs: {
		[HunterMajorGlyph.GlyphOfBestialWrath]: {
			name: 'Glyph of Bestial Wrath',
			description: 'Decreases the cooldown of Bestial Wrath by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_ferociousbite.jpg',
		},
		[HunterMajorGlyph.GlyphOfConcussiveShot]: {
			name: 'Glyph of Concussive Shot',
			description: 'Your Concussive Shot also limits the maximum run speed of your target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_stun.jpg',
		},
		[HunterMajorGlyph.GlyphOfDeterrence]: {
			name: 'Glyph of Deterrence',
			description: 'Decreases the cooldown of Deterrence by 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_whirlwind.jpg',
		},
		[HunterMajorGlyph.GlyphOfDisengage]: {
			name: 'Glyph of Disengage',
			description: 'Decreases the cooldown of Disengage by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[HunterMajorGlyph.GlyphOfFreezingTrap]: {
			name: 'Glyph of Freezing Trap',
			description: 'When your Freezing Trap breaks, the victim\'s movement speed is reduced by 70% for 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[HunterMajorGlyph.GlyphOfIceTrap]: {
			name: 'Glyph of Ice Trap',
			description: 'Increases the radius of the effect from your Ice Trap by 2 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg',
		},
		[HunterMajorGlyph.GlyphOfImmolationTrap]: {
			name: 'Glyph of Immolation Trap',
			description: 'Decreases the duration of the effect from your Immolation Trap by 6 sec, but damage while active is increased by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_fire_flameshock.jpg',
		},
		[HunterMajorGlyph.GlyphOfMasterSCall]: {
			name: 'Glyph of Master\'s Call',
			description: 'Increases the duration of your Master\'s Call by 4 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_masterscall.jpg',
		},
		[HunterMajorGlyph.GlyphOfMending]: {
			name: 'Glyph of Mending',
			description: 'Increases the total amount of healing done by your Mend Pet ability by 60%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_mendpet.jpg',
		},
		[HunterMajorGlyph.GlyphOfMisdirection]: {
			name: 'Glyph of Misdirection',
			description: 'When you use Misdirection on your pet, the cooldown on your Misdirection is reset.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_misdirection.jpg',
		},
		[HunterMajorGlyph.GlyphOfRaptorStrike]: {
			name: 'Glyph of Raptor Strike',
			description: 'Reduces damage taken by 20% for 5 sec after using Raptor Strike.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_meleedamage.jpg',
		},
		[HunterMajorGlyph.GlyphOfScatterShot]: {
			name: 'Glyph of Scatter Shot',
			description: 'Increases the range of Scatter Shot by 3 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_golemstormbolt.jpg',
		},
		[HunterMajorGlyph.GlyphOfSilencingShot]: {
			name: 'Glyph of Silencing Shot',
			description: 'When you successfully silence an enemy\'s spell cast with Silencing Shot, you instantly gain 10 Focus.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_theblackarrow.jpg',
		},
		[HunterMajorGlyph.GlyphOfSnakeTrap]: {
			name: 'Glyph of Snake Trap',
			description: 'Snakes from your Snake Trap take 90% reduced damage from area of effect spells.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_snaketrap.jpg',
		},
		[HunterMajorGlyph.GlyphOfTrapLauncher]: {
			name: 'Glyph of Trap Launcher',
			description: 'Reduces the Focus cost of Trap Launcher by 10.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_traplauncher.jpg',
		},
		[HunterMajorGlyph.GlyphOfWyvernSting]: {
			name: 'Glyph of Wyvern Sting',
			description: 'Decreases the cooldown of your Wyvern Sting by 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_spear_02.jpg',
		},
	},
	minorGlyphs: {
		[HunterMinorGlyph.GlyphOfAspectOfThePack]: {
			name: 'Glyph of Aspect of the Pack',
			description: 'Increases the range of your Aspect of the Pack ability by 15 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mount_whitetiger.jpg',
		},
		[HunterMinorGlyph.GlyphOfFeignDeath]: {
			name: 'Glyph of Feign Death',
			description: 'Reduces the cooldown of your Feign Death spell by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feigndeath.jpg',
		},
		[HunterMinorGlyph.GlyphOfLesserProportion]: {
			name: 'Glyph of Lesser Proportion',
			description: 'Slightly reduces the size of your pet.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_bestialdiscipline.jpg',
		},
		[HunterMinorGlyph.GlyphOfRevivePet]: {
			name: 'Glyph of Revive Pet',
			description: 'Reduces the pushback suffered from damaging attacks while casting Revive Pet by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_beastsoothe.jpg',
		},
		[HunterMinorGlyph.GlyphOfScareBeast]: {
			name: 'Glyph of Scare Beast',
			description: 'Reduces the pushback suffered from damaging attacks while casting Scare Beast by 75%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_cower.jpg',
		},
	},
};