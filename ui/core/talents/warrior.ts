import { WarriorMajorGlyph, WarriorMinorGlyph, WarriorTalents } from '../proto/warrior.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarriorTalentJson from './trees/warrior.json';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);

export const warriorGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarriorMajorGlyph.GlyphOfCleaving]: {
			name: 'Glyph of Cleaving',
			description: 'Increases the number of targets your Cleave hits by 1.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_cleave.jpg',
		},
		[WarriorMajorGlyph.GlyphOfColossusSmash]: {
			name: 'Glyph of Colossus Smash',
			description: 'Your Colossus Smash also applies the Sunder Armor effect to your target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_colossussmash.jpg',
		},
		[WarriorMajorGlyph.GlyphOfDeathWish]: {
			name: 'Glyph of Death Wish',
			description: 'Death Wish no longer increases damage taken.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathpact.jpg',
		},
		[WarriorMajorGlyph.GlyphOfHeroicThrow]: {
			name: 'Glyph of Heroic Throw',
			description: 'Your Heroic Throw applies a stack of Sunder Armor.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_axe_66.jpg',
		},
		[WarriorMajorGlyph.GlyphOfIntercept]: {
			name: 'Glyph of Intercept',
			description: 'Increases the duration of your Intercept stun by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[WarriorMajorGlyph.GlyphOfIntervene]: {
			name: 'Glyph of Intervene',
			description: 'Increases the number of attacks you intercept for your Intervene target by 1.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_victoryrush.jpg',
		},
		[WarriorMajorGlyph.GlyphOfLongCharge]: {
			name: 'Glyph of Long Charge',
			description: 'Increases the range of your Charge ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg',
		},
		[WarriorMajorGlyph.GlyphOfPiercingHowl]: {
			name: 'Glyph of Piercing Howl',
			description: 'Increases the radius of Piercing Howl by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathscream.jpg',
		},
		[WarriorMajorGlyph.GlyphOfRapidCharge]: {
			name: 'Glyph of Rapid Charge',
			description: 'Reduces the cooldown of your Charge ability by 1 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg',
		},
		[WarriorMajorGlyph.GlyphOfResonatingPower]: {
			name: 'Glyph of Resonating Power',
			description: 'Reduces the rage cost of your Thunder Clap ability by 5.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg',
		},
		[WarriorMajorGlyph.GlyphOfShieldWall]: {
			name: 'Glyph of Shield Wall',
			description: 'Shield Wall now reduces damage taken by an additional 20%, but its cooldown is increased by 2 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldwall.jpg',
		},
		[WarriorMajorGlyph.GlyphOfShockwave]: {
			name: 'Glyph of Shockwave',
			description: 'Reduces the cooldown on Shockwave by 3 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shockwave.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSpellReflection]: {
			name: 'Glyph of Spell Reflection',
			description: 'Reduces the cooldown on Spell Reflection by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldreflection.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSunderArmor]: {
			name: 'Glyph of Sunder Armor',
			description: 'When you use Sunder Armor or Devastate, a second nearby target also receives Sunder Armor.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_sunder.jpg',
		},
		[WarriorMajorGlyph.GlyphOfSweepingStrikes]: {
			name: 'Glyph of Sweeping Strikes',
			description: 'Reduces the rage cost of your Sweeping Strikes ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_slicedice.jpg',
		},
		[WarriorMajorGlyph.GlyphOfThunderClap]: {
			name: 'Glyph of Thunder Clap',
			description: 'Increases the radius of your Thunder Clap ability by 2 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg',
		},
		[WarriorMajorGlyph.GlyphOfVictoryRush]: {
			name: 'Glyph of Victory Rush',
			description: 'Increases the total healing provided by your Victory Rush by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg',
		},
	},
	minorGlyphs: {
		[WarriorMinorGlyph.GlyphOfBattle]: {
			name: 'Glyph of Battle',
			description: 'Increases the duration by 2 min and area of effect by 50% of your Battle Shout.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_battleshout.jpg',
		},
		[WarriorMinorGlyph.GlyphOfBerserkerRage]: {
			name: 'Glyph of Berserker Rage',
			description: 'Berserker Rage generates 5 Rage when used.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_ancestralguardian.jpg',
		},
		[WarriorMinorGlyph.GlyphOfBloodyHealing]: {
			name: 'Glyph of Bloody Healing',
			description: 'Increases the healing you receive from Bloodthirst by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_bloodlust.jpg',
		},
		[WarriorMinorGlyph.GlyphOfCommand]: {
			name: 'Glyph of Command',
			description: 'Increases the duration by 2 min and area of effect by 50% of your Commanding Shout.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_rallyingcry.jpg',
		},
		[WarriorMinorGlyph.GlyphOfDemoralizingShout]: {
			name: 'Glyph of Demoralizing Shout',
			description: 'Increases the duration by 15 sec and area of effect by 50% of your Demoralizing Shout.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_warcry.jpg',
		},
		[WarriorMinorGlyph.GlyphOfEnduringVictory]: {
			name: 'Glyph of Enduring Victory',
			description: 'Increases the window of opportunity in which you can use Victory Rush by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg',
		},
		[WarriorMinorGlyph.GlyphOfFuriousSundering]: {
			name: 'Glyph of Furious Sundering',
			description: 'Reduces the cost of Sunder Armor by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_sunder.jpg',
		},
		[WarriorMinorGlyph.GlyphOfIntimidatingShout]: {
			name: 'Glyph of Intimidating Shout',
			description: 'All targets of your Intimidating Shout now tremble in place instead of fleeing in fear.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_golemthunderclap.jpg',
		},
	},
};
