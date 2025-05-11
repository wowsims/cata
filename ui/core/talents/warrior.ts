import { WarriorMajorGlyph, WarriorMinorGlyph, WarriorTalents } from '../proto/warrior.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarriorTalentJson from './trees/warrior.json';export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);

export const warriorGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[WarriorMajorGlyph.GlyphOfLongCharge]: {
			name: "Glyph of Long Charge",
			description: "Increases the range of your Charge ability by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_charge.jpg",
		},
		[WarriorMajorGlyph.GlyphOfUnendingRage]: {
			name: "Glyph of Unending Rage",
			description: "Increases your maximum Rage by 20.0.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_intensifyrage.jpg",
		},
		[WarriorMajorGlyph.GlyphOfEnragedSpeed]: {
			name: "Glyph of Enraged Speed",
			description: "While Enraged, you move 20% faster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_racial_bloodrage.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHinderingStrikes]: {
			name: "Glyph of Hindering Strikes",
			description: "Your Heroic Strike and Cleave now also reduce the target\'s movement speed by 50% for 8s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_cleave.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHeavyRepercussions]: {
			name: "Glyph of Heavy Repercussions",
			description: "While your Shield Block is active, your Shield Slam hits for an additional 50% damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_defend.jpg",
		},
		[WarriorMajorGlyph.GlyphOfBloodthirst]: {
			name: "Glyph of Bloodthirst",
			description: "Increases the healing of Bloodthirst by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_bloodlust.jpg",
		},
		[WarriorMajorGlyph.GlyphOfRudeInterruption]: {
			name: "Glyph of Rude Interruption",
			description: "Successfully interrupting a spell with Pummel increases your damage by 6% for 20s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_commandingshout.jpg",
		},
		[WarriorMajorGlyph.GlyphOfGagOrder]: {
			name: "Glyph of Gag Order",
			description: "Your Pummel and Heroic Throw also silence the target for 3s. Does not work against players.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_axe_66.jpg",
		},
		[WarriorMajorGlyph.GlyphOfBlitz]: {
			name: "Glyph of Blitz",
			description: "Your Charge also roots and snares an additional 2 nearby targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_victoryrush.jpg",
		},
		[WarriorMajorGlyph.GlyphOfMortalStrike]: {
			name: "Glyph of Mortal Strike",
			description: "When your Mortal Strike is affecting a target, healing effects on you are increased by 10%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_savageblow.jpg",
		},
		[WarriorMajorGlyph.GlyphOfDieByTheSword]: {
			name: "Glyph of Die by the Sword",
			description: "While Die by the Sword is active, using Overpower increases its duration by 1 sec and using Wild Strike increases its duration by 0.5 sec per use.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_challange.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHamstring]: {
			name: "Glyph of Hamstring",
			description: "When you spend Rage to apply Hamstring, the Rage cost of your next Hamstring is reduced by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_shockwave.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHoldTheLine]: {
			name: "Glyph of Hold the Line",
			description: "Improves the damage of your next Revenge by 50% following a successful parry.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_revenge.jpg",
		},
		[WarriorMajorGlyph.GlyphOfShieldSlam]: {
			name: "Glyph of Shield Slam",
			description: "Your Shield Slam now dispels 1 magical effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_shield_05.jpg",
		},
		[WarriorMajorGlyph.GlyphOfHoarseVoice]: {
			name: "Glyph of Hoarse Voice",
			description: "Reduces the cooldown and Rage generation of your Battle and Commanding Shout by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_battleshout.jpg",
		},
		[WarriorMajorGlyph.GlyphOfSweepingStrikes]: {
			name: "Glyph of Sweeping Strikes",
			description: "When you hit a target with Sweeping Strikes, you gain 1.0 Rage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_slicedice.jpg",
		},
		[WarriorMajorGlyph.GlyphOfResonatingPower]: {
			name: "Glyph of Resonating Power",
			description: "Increases the damage and cooldown of Thunder Clap by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_thunderclap.jpg",
		},
		[WarriorMajorGlyph.GlyphOfVictoryRush]: {
			name: "Glyph of Victory Rush",
			description: "Increases the total healing provided by your Victory Rush by 50%.\u000D\u000A This glyph has no effect if combined with the Impending Victory talent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg",
		},
		[WarriorMajorGlyph.GlyphOfRagingWind]: {
			name: "Glyph of Raging Wind",
			description: "Your Raging Blow hits increase the damage of your next Whirlwind by 10%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_swiftstrike.jpg",
		},
		[WarriorMajorGlyph.GlyphOfWhirlwind]: {
			name: "Glyph of Whirlwind",
			description: "Increases the radius of Whirlwind by 4 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_whirlwind.jpg",
		},
		[WarriorMajorGlyph.GlyphOfDeathFromAbove]: {
			name: "Glyph of Death From Above",
			description: "Reduces the cooldown on Heroic Leap by 15 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_heroicleap.jpg",
		},
		[WarriorMajorGlyph.GlyphOfVictoriousThrow]: {
			name: "Glyph of Victorious Throw",
			description: "Increases the range of Victory Rush and Impending Victory by 15 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_warrior_wildstrike.jpg",
		},
		[WarriorMajorGlyph.GlyphOfSpellReflection]: {
			name: "Glyph of Spell Reflection",
			description: "Reduces the cooldown on Spell Reflection by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldreflection.jpg",
		},
		[WarriorMajorGlyph.GlyphOfShieldWall]: {
			name: "Glyph of Shield Wall",
			description: "Shield Wall now reduces damage taken by an additional 20%, but its cooldown is increased by 2 min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldwall.jpg",
		},
		[WarriorMajorGlyph.GlyphOfColossusSmash]: {
			name: "Glyph of Colossus Smash",
			description: "Your Colossus Smash also applies the Sunder Armor effect to your target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_colossussmash.jpg",
		},
		[WarriorMajorGlyph.GlyphOfBullRush]: {
			name: "Glyph of Bull Rush",
			description: "Your Charge generates 15 additional Rage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_character_tauren_male.jpg",
		},
		[WarriorMajorGlyph.GlyphOfRecklessness]: {
			name: "Glyph of Recklessness",
			description: "Decreases the critical chance of Recklessness by 12% but increases its duration by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_criticalstrike.jpg",
		},
		[WarriorMajorGlyph.GlyphOfIncite]: {
			name: "Glyph of Incite",
			description: "Using Demoralizing Shout makes your next 3 Heroic Strike or Cleave abilities free.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_incite.jpg",
		},
		[WarriorMajorGlyph.GlyphOfImpalingThrows]: {
			name: "Glyph of Impaling Throws",
			description: "Heroic Throw now leaves an axe in the target, which can be retrieved by moving within 5 yards of the target to finish the cooldown of Heroic Throw. This effect will only occur when Heroic Throw is cast from more than 10 yards away from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_axe_05.jpg",
		},
		[WarriorMajorGlyph.GlyphOfTheExecutor]: {
			name: "Glyph of the Executor",
			description: "Killing an enemy with Execute grants you 30 rage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/warrior_talent_icon_lambstotheslaughter.jpg",
		},
	},
	minorGlyphs: {
		[WarriorMinorGlyph.GlyphOfMysticShout]: {
			name: "Glyph of Mystic Shout",
			description: "Your Battle Shout and Commanding Shout cause you to hover in the air for 1s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_battleshout.jpg",
		},
		[WarriorMinorGlyph.GlyphOfBloodcurdlingShout]: {
			name: "Glyph of Bloodcurdling Shout",
			description: "Your Battle Shout and Commanding Shout terrify small animals.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_furiousresolve.jpg",
		},
		[WarriorMinorGlyph.GlyphOfGushingWound]: {
			name: "Glyph of Gushing Wound",
			description: "Your Deep Wounds are even bloodier than normal.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_backstab.jpg",
		},
		[WarriorMinorGlyph.GlyphOfMightyVictory]: {
			name: "Glyph of Mighty Victory",
			description: "When your Victory Rush or Impending Victory heal you, you grow slightly larger.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_devastate.jpg",
		},
		[WarriorMinorGlyph.GlyphOfBloodyHealing]: {
			name: "Glyph of Bloody Healing",
			description: "Increases the healing you receive from bandages by 20% while your Deep Wounds is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_backstab.jpg",
		},
		[WarriorMinorGlyph.GlyphOfIntimidatingShout]: {
			name: "Glyph of Intimidating Shout",
			description: "All targets of your Intimidating Shout now tremble in place instead of fleeing in fear.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_golemthunderclap.jpg",
		},
		[WarriorMinorGlyph.GlyphOfThunderStrike]: {
			name: "Glyph of Thunder Strike",
			description: "Your Thunder Clap visual includes a lightning strike.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_thunderclap.jpg",
		},
		[WarriorMinorGlyph.GlyphOfCrowFeast]: {
			name: "Glyph of Crow Feast",
			description: "Your Execute critical strikes summon a flock of carrion birds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_sword_48.jpg",
		},
		[WarriorMinorGlyph.GlyphOfBurningAnger]: {
			name: "Glyph of Burning Anger",
			description: "You get so angry when Enraged that you catch on fire.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unholyfrenzy.jpg",
		},
		[WarriorMinorGlyph.GlyphOfTheBlazingTrail]: {
			name: "Glyph of the Blazing Trail",
			description: "Your Charge leaves a trail of fire in its wake. If you\'re going to Charge why not do it with some style?",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_burningspeed.jpg",
		},
		[WarriorMinorGlyph.GlyphOfTheRagingWhirlwind]: {
			name: "Glyph of the Raging Whirlwind",
			description: "Whirlwind gives you 15 rage over 6s, but for that time you no longer generate rage from autoattacks.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_unleashedrage.jpg",
		},
		[WarriorMinorGlyph.GlyphOfTheSubtleDefender]: {
			name: "Glyph of the Subtle Defender",
			description: "Removes the threat generation bonus from Defensive Stance.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_shieldguard.jpg",
		},
		[WarriorMinorGlyph.GlyphOfTheWatchfulEye]: {
			name: "Glyph of the Watchful Eye",
			description: "Intervene will now target the party or raid member with the lowest health within 25 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_bloodyeye.jpg",
		},
		[WarriorMinorGlyph.GlyphOfTheWeaponmaster]: {
			name: "Glyph of the Weaponmaster",
			description: "Your Shout abilities cause the appearance of your weapon to change to that of a random weapon from your primary bag for 0ms.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_warrior_weaponmastery.jpg",
		},
	},
};
