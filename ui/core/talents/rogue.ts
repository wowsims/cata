import { RogueMajorGlyph, RogueMinorGlyph, RoguePrimeGlyph, RogueTalents } from '../proto/rogue.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import RogueTalentJson from './trees/rogue.json';export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

export const rogueGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[RoguePrimeGlyph.GlyphOfAdrenalineRush]: {
			name: "Glyph of Adrenaline Rush",
			description: "Increases the duration of Adrenaline Rush by ${$56808m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfBackstab]: {
			name: "Glyph of Backstab",
			description: "Your Backstab critical strikes grant you $56800s1 Energy.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfEviscerate]: {
			name: "Glyph of Eviscerate",
			description: "Increases the critical strike chance of Eviscerate by $56802s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfRevealingStrike]: {
			name: "Glyph of Revealing Strike",
			description: "Increases Revealing Strike's bonus effectiveness to your finishing moves by an additional $56814s1%",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfHemorrhage]: {
			name: "Glyph of Hemorrhage",
			description: "Your Hemorrhage ability also causes the target to bleed, dealing $56807s1% of the direct strike's damage over $89775d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfRupture]: {
			name: "Glyph of Rupture",
			description: "Increases the duration of Rupture by ${$56801m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfSinisterStrike]: {
			name: "Glyph of Sinister Strike",
			description: "Your Sinister Strikes have a $h% chance to add an additional combo point.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfSliceAndDice]: {
			name: "Glyph of Slice and Dice",
			description: "Increases the duration of Slice and Dice by ${$56810m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfVendetta]: {
			name: "Glyph of Vendetta",
			description: "Increases the duration of your Vendetta ability by $63249s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfKillingSpree]: {
			name: "Glyph of Killing Spree",
			description: "Increases the bonus to your damage while Killing Spree is active by an additional $63252s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfShadowDance]: {
			name: "Glyph of Shadow Dance",
			description: "Increases the duration of Shadow Dance by ${$63253m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfMutilate]: {
			name: "Glyph of Mutilate",
			description: "Reduces the cost of Mutilate by $63268s1 Energy.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
		[RoguePrimeGlyph.GlyphOfStabbing]: {
			name: "Glyph of Stabbing",
			description: "Your Backstab critical strikes grant you $56800s1 Energy.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_primerogue.jpg",
		},
	},
	majorGlyphs: {
		[RogueMajorGlyph.GlyphOfAmbush]: {
			name: "Glyph of Ambush",
			description: "Increases the range on Ambush by $56813s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfBladeFlurry]: {
			name: "Glyph of Blade Flurry",
			description: "Reduces the penalty to Energy generation while Blade Flurry is active by $56818s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfCripplingPoison]: {
			name: "Glyph of Crippling Poison",
			description: "Increases the chance to inflict your target with Crippling Poison by an additional $56820s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfDeadlyThrow]: {
			name: "Glyph of Deadly Throw",
			description: "Increases the slowing effect on Deadly Throw by $56806s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfEvasion]: {
			name: "Glyph of Evasion",
			description: "Increases the duration of Evasion by ${$56799m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfExposeArmor]: {
			name: "Glyph of Expose Armor",
			description: "Increases the duration of Expose Armor by ${$56803m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfFeint]: {
			name: "Glyph of Feint",
			description: "Reduces the Energy cost of Feint by $56804s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfGarrote]: {
			name: "Glyph of Garrote",
			description: "Increases the duration of your Garrote ability's silence effect by ${$56812m1/1000}.1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfGouge]: {
			name: "Glyph of Gouge",
			description: "Your Gouge ability no longer requires that the target be facing you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfPreparation]: {
			name: "Glyph of Preparation",
			description: "Your Preparation ability also instantly resets the cooldown of Kick, Dismantle, and Smoke Bomb.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfSap]: {
			name: "Glyph of Sap",
			description: "Increases the duration of Sap against non-player targets by ${$56798m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfKick]: {
			name: "Glyph of Kick",
			description: "Increases the cooldown of your Kick ability by ${$56805m1/1000} sec, but this cooldown is reduced by ${$56805m2/1000} sec when your Kick successfully interrupts a spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfSprint]: {
			name: "Glyph of Sprint",
			description: "Increases the movement speed of your Sprint ability by an additional $56811s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfFanOfKnives]: {
			name: "Glyph of Fan of Knives",
			description: "Increases the radius of your Fan of Knives ability by $63254s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfTricksOfTheTrade]: {
			name: "Glyph of Tricks of the Trade",
			description: "Removes the Energy cost of your Tricks of the Trade ability but reduces the recipient's damage bonus by $63256s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfCloakOfShadows]: {
			name: "Glyph of Cloak of Shadows",
			description: "While Cloak of Shadows is active, you take $63269s1% less physical damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfVanish]: {
			name: "Glyph of Vanish",
			description: "Increases the duration of your Vanish effect by ${$89758m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
		[RogueMajorGlyph.GlyphOfBlind]: {
			name: "Glyph of Blind",
			description: "Your Blind ability also removes all damage over time effects from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_majorrogue.jpg",
		},
	},
	minorGlyphs: {
		[RogueMinorGlyph.GlyphOfPickPocket]: {
			name: "Glyph of Pick Pocket",
			description: "Increases the range of your Pick Pocket ability by $58017s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
		[RogueMinorGlyph.GlyphOfDistract]: {
			name: "Glyph of Distract",
			description: "Increases the range of your Distract ability by $58032s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
		[RogueMinorGlyph.GlyphOfPickLock]: {
			name: "Glyph of Pick Lock",
			description: "Reduces the cast time of your Pick Lock ability by $58027s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
		[RogueMinorGlyph.GlyphOfSafeFall]: {
			name: "Glyph of Safe Fall",
			description: "Increases the distance your Safe Fall ability allows you to fall without taking damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
		[RogueMinorGlyph.GlyphOfBlurredSpeed]: {
			name: "Glyph of Blurred Speed",
			description: "Enables you to walk on water while your Sprint ability is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
		[RogueMinorGlyph.GlyphOfPoisons]: {
			name: "Glyph of Poisons",
			description: "You apply poisons to your weapons $58038s1% faster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_glyph_minorrogue.jpg",
		},
	},
};
