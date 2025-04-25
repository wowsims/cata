import { RogueMajorGlyph, RogueMinorGlyph, RogueTalents } from '../proto/rogue.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import RogueTalentJson from './trees/rogue.json';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

export const rogueGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[RogueMajorGlyph.GlyphOfAmbush]: {
			name: 'Glyph of Ambush',
			description: 'Increases the range on Ambush by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_ambush.jpg',
		},
		[RogueMajorGlyph.GlyphOfBladeFlurry]: {
			name: 'Glyph of Blade Flurry',
			description: 'Reduces the penalty to Energy generation while Blade Flurry is active by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_punishingblow.jpg',
		},
		[RogueMajorGlyph.GlyphOfBlind]: {
			name: 'Glyph of Blind',
			description: 'Your Blind ability also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_mindsteal.jpg',
		},
		[RogueMajorGlyph.GlyphOfCloakOfShadows]: {
			name: 'Glyph of Cloak of Shadows',
			description: 'While Cloak of Shadows is active, you take 40% less physical damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_nethercloak.jpg',
		},
		[RogueMajorGlyph.GlyphOfCripplingPoison]: {
			name: 'Glyph of Crippling Poison',
			description: 'Increases the chance to inflict your target with Crippling Poison by an additional 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_poisonsting.jpg',
		},
		[RogueMajorGlyph.GlyphOfDeadlyThrow]: {
			name: 'Glyph of Deadly Throw',
			description: 'Increases the slowing effect on Deadly Throw by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_throwingknife_06.jpg',
		},
		[RogueMajorGlyph.GlyphOfEvasion]: {
			name: 'Glyph of Evasion',
			description: 'Increases the duration of Evasion by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowward.jpg',
		},
		[RogueMajorGlyph.GlyphOfExposeArmor]: {
			name: 'Glyph of Expose Armor',
			description: 'Increases the duration of Expose Armor by 12 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_warrior_riposte.jpg',
		},
		[RogueMajorGlyph.GlyphOfFanOfKnives]: {
			name: 'Glyph of Fan of Knives',
			description: 'Increases the radius of your Fan of Knives ability by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_fanofknives.jpg',
		},
		[RogueMajorGlyph.GlyphOfFeint]: {
			name: 'Glyph of Feint',
			description: 'Reduces the Energy cost of Feint by 20.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg',
		},
		[RogueMajorGlyph.GlyphOfGarrote]: {
			name: 'Glyph of Garrote',
			description: "Increases the duration of your Garrote ability's silence effect by 1.5 sec.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_garrote.jpg',
		},
		[RogueMajorGlyph.GlyphOfGouge]: {
			name: 'Glyph of Gouge',
			description: 'Your Gouge ability no longer requires that the target be facing you.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_gouge.jpg',
		},
		[RogueMajorGlyph.GlyphOfKick]: {
			name: 'Glyph of Kick',
			description:
				'Increases the cooldown of your Kick ability by 4 sec, but this cooldown is reduced by 6 sec when your Kick successfully interrupts a spell.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_kick.jpg',
		},
		[RogueMajorGlyph.GlyphOfPreparation]: {
			name: 'Glyph of Preparation',
			description: 'Your Preparation ability also instantly resets the cooldown of Kick, Dismantle, and Smoke Bomb.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_preparation.jpg',
		},
		[RogueMajorGlyph.GlyphOfSap]: {
			name: 'Glyph of Sap',
			description: 'Increases the duration of Sap against non-player targets by 80 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_sap.jpg',
		},
		[RogueMajorGlyph.GlyphOfSprint]: {
			name: 'Glyph of Sprint',
			description: 'Increases the movement speed of your Sprint ability by an additional 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMajorGlyph.GlyphOfTricksOfTheTrade]: {
			name: 'Glyph of Tricks of the Trade',
			description: "Removes the Energy cost of your Tricks of the Trade ability but reduces the recipient's damage bonus by 5%.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_tricksofthetrade.jpg',
		},
		[RogueMajorGlyph.GlyphOfVanish]: {
			name: 'Glyph of Vanish',
			description: 'Increases the duration of your Vanish effect by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_vanish.jpg',
		},
	},
	minorGlyphs: {
		[RogueMinorGlyph.GlyphOfBlurredSpeed]: {
			name: 'Glyph of Blurred Speed',
			description: 'Enables you to walk on water while your Sprint ability is active.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg',
		},
		[RogueMinorGlyph.GlyphOfDistract]: {
			name: 'Glyph of Distract',
			description: 'Increases the range of your Distract ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_rogue_distract.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickLock]: {
			name: 'Glyph of Pick Lock',
			description: 'Reduces the cast time of your Pick Lock ability by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_moonkey.jpg',
		},
		[RogueMinorGlyph.GlyphOfPickPocket]: {
			name: 'Glyph of Pick Pocket',
			description: 'Increases the range of your Pick Pocket ability by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_bag_11.jpg',
		},
		[RogueMinorGlyph.GlyphOfPoisons]: {
			name: 'Glyph of Poisons',
			description: 'You apply poisons to your weapons 50% faster.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/trade_brewpoison.jpg',
		},
		[RogueMinorGlyph.GlyphOfSafeFall]: {
			name: 'Glyph of Safe Fall',
			description: 'Increases the distance your Safe Fall ability allows you to fall without taking damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_feather_01.jpg',
		},
	},
};
