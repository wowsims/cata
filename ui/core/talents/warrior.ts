import { WarriorMajorGlyph, WarriorMinorGlyph, WarriorPrimeGlyph, WarriorTalents } from '../proto/warrior.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarriorTalentJson from './trees/warrior.json';export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);

export const warriorGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[WarriorPrimeGlyph.GlyphOfDevastate]: {
			name: "Glyph of Devastate",
			description: "Increases the critical strike chance of Devastate by $58388s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_sword_11.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfBloodthirst]: {
			name: "Glyph of Bloodthirst",
			description: "Increases the damage of Bloodthirst by $58367s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_bloodlust.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfMortalStrike]: {
			name: "Glyph of Mortal Strike",
			description: "Increases the damage of Mortal Strike by $58368s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_savageblow.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfOverpower]: {
			name: "Glyph of Overpower",
			description: "Increases the damage of Overpower by $58386s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_meleedamage.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfSlam]: {
			name: "Glyph of Slam",
			description: "Increases the critical strike chance of Slam by $58385s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_decisivestrike.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfRevenge]: {
			name: "Glyph of Revenge",
			description: "Increases the damage of Revenge by $58364s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_revenge.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfShieldSlam]: {
			name: "Glyph of Shield Slam",
			description: "Increases the damage of Shield Slam by $58375s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_shield_05.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfRagingBlow]: {
			name: "Glyph of Raging Blow",
			description: "Increases the critical strike chance of Raging Blow by $58370s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_swiftstrike.jpg",
		},
		[WarriorPrimeGlyph.GlyphOfBladestorm]: {
			name: "Glyph of Bladestorm",
			description: "Reduces the cooldown on Bladestorm by ${$63324m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_bladestorm.jpg",
		},
	},
	majorGlyphs: {
		[WarriorMajorGlyph.GlyphOfLongCharge]: {
			name: "Glyph of Long Charge",
			description: "Increases the range of your Charge ability by $58097s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg",
		},
		[WarriorMajorGlyph.GlyphOfThunderClap]: {
			name: "Glyph of Thunder Clap",
			description: "Increases the radius of your Thunder Clap ability by $58098s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg",
		},
		[WarriorMajorGlyph.GlyphOfRapidCharge]: {
			name: "Glyph of Rapid Charge",
			description: "Reduces the cooldown of your Charge ability by $/1000;58355s1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg",
		},
		[WarriorMajorGlyph.GlyphOfCleaving]: {
			name: "Glyph of Cleaving",
			description: "Increases the number of targets your Cleave hits by 1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_cleave.jpg",
		},
		[WarriorMajorGlyph.GlyphOfPiercingHowl]: {
			name: "Glyph of Piercing Howl",
			description: "Increases the radius of Piercing Howl by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathscream.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHeroicThrow]: {
			name: "Glyph of Heroic Throw",
			description: "Your Heroic Throw applies a stack of Sunder Armor.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_axe_66.jpg",
		},
		[WarriorMajorGlyph.GlyphOfIntervene]: {
			name: "Glyph of Intervene",
			description: "Increases the number of attacks you intercept for your Intervene target by $58377s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_victoryrush.jpg",
		},
		[WarriorMajorGlyph.GlyphOfSunderArmor]: {
			name: "Glyph of Sunder Armor",
			description: "When you use Sunder Armor or Devastate, a second nearby target also receives Sunder Armor.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_sunder.jpg",
		},
		[WarriorMajorGlyph.GlyphOfSweepingStrikes]: {
			name: "Glyph of Sweeping Strikes",
			description: "Reduces the rage cost of your Sweeping Strikes ability by $58384s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_slicedice.jpg",
		},
		[WarriorMajorGlyph.GlyphOfResonatingPower]: {
			name: "Glyph of Resonating Power",
			description: "Reduces the rage cost of your Thunder Clap ability by ${$58356m1/-10}.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg",
		},
		[WarriorMajorGlyph.GlyphOfVictoryRush]: {
			name: "Glyph of Victory Rush",
			description: "Increases the total healing provided by your Victory Rush by $58382s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg",
		},
		[WarriorMajorGlyph.GlyphOfShockwave]: {
			name: "Glyph of Shockwave",
			description: "Reduces the cooldown on Shockwave by ${$63325m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shockwave.jpg",
		},
		[WarriorMajorGlyph.GlyphOfSpellReflection]: {
			name: "Glyph of Spell Reflection",
			description: "Reduces the cooldown on Spell Reflection by ${$63328m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldreflection.jpg",
		},
		[WarriorMajorGlyph.GlyphOfShieldWall]: {
			name: "Glyph of Shield Wall",
			description: "Shield Wall now reduces damage taken by an additional $63329m2%, but its cooldown is increased by ${$63329m1/60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldwall.jpg",
		},
		[WarriorMajorGlyph.GlyphOfColossusSmash]: {
			name: "Glyph of Colossus Smash",
			description: "Your Colossus Smash also applies the Sunder Armor effect to your target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_colossussmash.jpg",
		},
		[WarriorMajorGlyph.GlyphOfIntercept]: {
			name: "Glyph of Intercept",
			description: "Increases the duration of your Intercept stun by ${$94372m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_sprint.jpg",
		},
		[WarriorMajorGlyph.GlyphOfDeathWish]: {
			name: "Glyph of Death Wish",
			description: "Death Wish no longer increases damage taken.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathpact.jpg",
		},
	},
	minorGlyphs: {
		[WarriorMinorGlyph.GlyphOfBattle]: {
			name: "Glyph of Battle",
			description: "Increases the duration by ${$58095m1/60000} min and area of effect by $58095s2% of your Battle Shout.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_battleshout.jpg",
		},
		[WarriorMinorGlyph.GlyphOfBerserkerRage]: {
			name: "Glyph of Berserker Rage",
			description: "Berserker Rage generates $/10;23690s1 Rage when used.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_ancestralguardian.jpg",
		},
		[WarriorMinorGlyph.GlyphOfDemoralizingShout]: {
			name: "Glyph of Demoralizing Shout",
			description: "Increases the duration by ${$58099m1/1000} sec and area of effect by $58099s2% of your Demoralizing Shout.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_warcry.jpg",
		},
		[WarriorMinorGlyph.GlyphOfEnduringVictory]: {
			name: "Glyph of Enduring Victory",
			description: "Increases the window of opportunity in which you can use Victory Rush by ${$58104m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg",
		},
		[WarriorMinorGlyph.GlyphOfBloodyHealing]: {
			name: "Glyph of Bloody Healing",
			description: "Increases the healing you receive from Bloodthirst by $58369s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_bloodlust.jpg",
		},
		[WarriorMinorGlyph.GlyphOfFuriousSundering]: {
			name: "Glyph of Furious Sundering",
			description: "Reduces the cost of Sunder Armor by $63326s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_sunder.jpg",
		},
		[WarriorMinorGlyph.GlyphOfIntimidatingShout]: {
			name: "Glyph of Intimidating Shout",
			description: "All targets of your Intimidating Shout now tremble in place instead of fleeing in fear.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_golemthunderclap.jpg",
		},
		[WarriorMinorGlyph.GlyphOfCommand]: {
			name: "Glyph of Command",
			description: "Increases the duration by ${$68164m1/60000} min and area of effect by $68164s2% of your Commanding Shout.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_rallyingcry.jpg",
		},
	},
};
