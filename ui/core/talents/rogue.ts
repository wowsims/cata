import { RogueMajorGlyph, RogueMinorGlyph, RogueTalents } from '../proto/rogue.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import RogueTalentJson from './trees/rogue.json';export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

export const rogueGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[RogueMajorGlyph.GlyphOfShadowWalk]: {
			name: "Glyph of Shadow Walk",
			description: "Your Shadow Walk ability also increases your stealth detection while active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_envelopingshadows.jpg",
		},
		[RogueMajorGlyph.GlyphOfAmbush]: {
			name: "Glyph of Ambush",
			description: "Increases the range of Ambush by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_ambush.jpg",
		},
		[RogueMajorGlyph.GlyphOfBladeFlurry]: {
			name: "Glyph of Blade Flurry",
			description: "Your attacks have a 30% higher chance of applying Non-Lethal poisons while Blade Flurry is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_punishingblow.jpg",
		},
		[RogueMajorGlyph.GlyphOfSharpKnives]: {
			name: "Glyph of Sharp Knives",
			description: "Your Fan of Kinves also damages the armor of its victims, applying 1 application of the Weakened Armor effect to each target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_stone_sharpeningstone_05.jpg",
		},
		[RogueMajorGlyph.GlyphOfRecuperate]: {
			name: "Glyph of Recuperate",
			description: "Increases the healing of your Recuperate ability by an additional 1.0% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_recuperate.jpg",
		},
		[RogueMajorGlyph.GlyphOfEvasion]: {
			name: "Glyph of Evasion",
			description: "Increases the duration of Evasion by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowward.jpg",
		},
		[RogueMajorGlyph.GlyphOfRecovery]: {
			name: "Glyph of Recovery",
			description: "While Recuperate is active, you receive 20% increased healing from other sources.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sturdyrecuperate.jpg",
		},
		[RogueMajorGlyph.GlyphOfExposeArmor]: {
			name: "Glyph of Expose Armor",
			description: "Your Expose Armor ability causes three applications of Weakened Armor.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_riposte.jpg",
		},
		[RogueMajorGlyph.GlyphOfFeint]: {
			name: "Glyph of Feint",
			description: "Increases the duration of Feint by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_feint.jpg",
		},
		[RogueMajorGlyph.GlyphOfGarrote]: {
			name: "Glyph of Garrote",
			description: "Increases the duration of your Garrote ability\'s silence effect by 1.0 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_garrote.jpg",
		},
		[RogueMajorGlyph.GlyphOfGouge]: {
			name: "Glyph of Gouge",
			description: "Your Gouge ability no longer requires that the target be facing you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_gouge.jpg",
		},
		[RogueMajorGlyph.GlyphOfSmokeBomb]: {
			name: "Glyph of Smoke Bomb",
			description: "Increases the duration of your Smoke Bomb by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_smoke.jpg",
		},
		[RogueMajorGlyph.GlyphOfCheapShot]: {
			name: "Glyph of Cheap Shot",
			description: "Increases the duration of your Cheap Shot by 0.5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_cheapshot.jpg",
		},
		[RogueMajorGlyph.GlyphOfHemorraghingVeins]: {
			name: "Glyph of Hemorraghing Veins",
			description: "Your Sanguinary Veins ability now also increases damage done to targets affected by your Hemorrhage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofsacrifice.jpg",
		},
		[RogueMajorGlyph.GlyphOfKick]: {
			name: "Glyph of Kick",
			description: "Increases the cooldown of your Kick ability by 4 sec, but this cooldown is reduced by 6 sec when your Kick successfully interrupts a spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_kick.jpg",
		},
		[RogueMajorGlyph.GlyphOfRedirect]: {
			name: "Glyph of Redirect",
			description: "Reduces the cooldown of Redirect by 50 seconds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_redirect.jpg",
		},
		[RogueMajorGlyph.GlyphOfShiv]: {
			name: "Glyph of Shiv",
			description: "Reduces the cooldown of your Shiv ability by 3 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_throwingknife_04.jpg",
		},
		[RogueMajorGlyph.GlyphOfSprint]: {
			name: "Glyph of Sprint",
			description: "Increases the movement speed of your Sprint ability by an additional 30%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg",
		},
		[RogueMajorGlyph.GlyphOfVendetta]: {
			name: "Glyph of Vendetta",
			description: "Reduces the damage bonus of your Vendetta ability by 5% but increases its duration by 10 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_deadliness.jpg",
		},
		[RogueMajorGlyph.GlyphOfStealth]: {
			name: "Glyph of Stealth",
			description: "Reduces the cooldown of your Stealth ability by 4 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_stealth.jpg",
		},
		[RogueMajorGlyph.GlyphOfDeadlyMomentum]: {
			name: "Glyph of Deadly Momentum",
			description: "When you land a killing blow on an opponent that yields experience or honor, your Slice and Dice and Recuperate abilities are refreshed to their original duration.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_deadlymomentum.jpg",
		},
		[RogueMajorGlyph.GlyphOfCloakOfShadows]: {
			name: "Glyph of Cloak of Shadows",
			description: "While Cloak of Shadows is active, you take 40% less physical damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_nethercloak.jpg",
		},
		[RogueMajorGlyph.GlyphOfVanish]: {
			name: "Glyph of Vanish",
			description: "Increases the duration of your Vanish effect by 2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_vanish.jpg",
		},
		[RogueMajorGlyph.GlyphOfBlind]: {
			name: "Glyph of Blind",
			description: "Your Blind ability also removes all damage over time effects from the target that would cause Blind to break early.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_mindsteal.jpg",
		},
	},
	minorGlyphs: {
		[RogueMinorGlyph.GlyphOfDecoy]: {
			name: "Glyph of Decoy",
			description: "When you Vanish, you leave behind a brief illusion that very closely resembles you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_vanish.jpg",
		},
		[RogueMinorGlyph.GlyphOfDetection]: {
			name: "Glyph of Detection",
			description: "Teaches you the ability Detection.\u000D\u000A\u000D\u000A Focus intently on trying to detect certain creatures.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_spy.jpg",
		},
		[RogueMinorGlyph.GlyphOfHemorrhage]: {
			name: "Glyph of Hemorrhage",
			description: "Your Hemorrhage ability only causes lingering damage over time to targets that were already afflicted by a Bleed effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg",
		},
		[RogueMinorGlyph.GlyphOfPickPocket]: {
			name: "Glyph of Pick Pocket",
			description: "Increases the range of your Pick Pocket ability by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_bag_11.jpg",
		},
		[RogueMinorGlyph.GlyphOfDistract]: {
			name: "Glyph of Distract",
			description: "Increases the range of your Distract ability by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_distract.jpg",
		},
		[RogueMinorGlyph.GlyphOfPickLock]: {
			name: "Glyph of Pick Lock",
			description: "Reduces the cast time of your Pick Lock ability by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_moonkey.jpg",
		},
		[RogueMinorGlyph.GlyphOfSafeFall]: {
			name: "Glyph of Safe Fall",
			description: "Increases the distance your Safe Fall ability allows you to fall without taking damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_feather_01.jpg",
		},
		[RogueMinorGlyph.GlyphOfBlurredSpeed]: {
			name: "Glyph of Blurred Speed",
			description: "Enables you to walk on water while your Sprint ability is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg",
		},
		[RogueMinorGlyph.GlyphOfPoisons]: {
			name: "Glyph of Poisons",
			description: "You apply poisons to your weapons 50% faster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/trade_brewpoison.jpg",
		},
		[RogueMinorGlyph.GlyphOfKillingSpree]: {
			name: "Glyph of Killing Spree",
			description: "Your Killing Spree returns you to your starting location when the effect ends.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_murderspree.jpg",
		},
		[RogueMinorGlyph.GlyphOfTricksOfTheTrade]: {
			name: "Glyph of Tricks of the Trade",
			description: "Your Tricks of the Trade ability no longer costs Energy, but also no longer increases the damage dealt by the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_tricksofthetrade.jpg",
		},
		[RogueMinorGlyph.GlyphOfDisguise]: {
			name: "Glyph of Disguise",
			description: "When you Pick Pocket a humanoid enemy, you also copy their appearance for 0ms. Your disguise will unravel upon entering combat.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_disguise.jpg",
		},
		[RogueMinorGlyph.GlyphOfHeadhunting]: {
			name: "Glyph of Headhunting",
			description: "Your Throw and Deadly Throw abilities will now throw axes regardless of your currently equipped weapon.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_throwingaxe_03.jpg",
		},
		[RogueMinorGlyph.GlyphOfImprovedDistraction]: {
			name: "Glyph of Improved Distraction",
			description: "Distract now summons a decoy at the target location.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_distract.jpg",
		},
	},
};
