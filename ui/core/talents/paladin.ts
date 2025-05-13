import { PaladinMajorGlyph, PaladinMinorGlyph, PaladinTalents } from '../proto/paladin.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PaladinTalentJson from './trees/paladin.json';export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

export const paladinGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PaladinMajorGlyph.GlyphOfDoubleJeopardy]: {
			name: "Glyph of Double Jeopardy",
			description: "Judging a target increases the damage of your next Judgment by 20%, but only if used on a different second target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousfury.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDevotionAura]: {
			name: "Glyph of Devotion Aura",
			description: "Devotion Aura no longer affects party or raid members, but the cooldown is reduced by 60 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_auramastery.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHolyWrath]: {
			name: "Glyph of Holy Wrath",
			description: "Your Holy Wrath now also stuns Elementals, Dragonkin, and Aberrations.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_weaponmastery.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivineProtection]: {
			name: "Glyph of Divine Protection",
			description: "Reduces the magical damage reduction of your Divine Protection to 20% but adds 20% physical damage reduction.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_divineprotection.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTemplarsVerdict]: {
			name: "Glyph of Templar's Verdict",
			description: "You take 10% less damage for 6s after dealing damage with Templar\'s Verdict or Exorcism.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_paladin_templarsverdict.jpg",
		},
		[PaladinMajorGlyph.GlyphOfAvengingWrath]: {
			name: "Glyph of Avenging Wrath",
			description: "While Avenging Wrath is active, you are healed for 1% of your maximum health every 0.0 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_avenginewrath.jpg",
		},
		[PaladinMajorGlyph.GlyphOfConsecration]: {
			name: "Glyph of Consecration",
			description: "You can now target Consecration anywhere within 25 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg",
		},
		[PaladinMajorGlyph.GlyphOfFocusedShield]: {
			name: "Glyph of Focused Shield",
			description: "Your Avenger\'s Shield hits 2 fewer targets, but for 30% more damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfBurdenOfGuilt]: {
			name: "Glyph of Burden of Guilt",
			description: "Your Judgment hits fill your target with doubt and remorse, reducing movement speed by 50% for 2s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_bg_topdps.jpg",
		},
		[PaladinMajorGlyph.GlyphOfBlindingLight]: {
			name: "Glyph of Blinding Light",
			description: "Your Blinding Light now knocks down targets for 0ms instead of Blinding them.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_blindinglight.jpg",
		},
		[PaladinMajorGlyph.GlyphOfFinalWrath]: {
			name: "Glyph of Final Wrath",
			description: "Your Holy Wrath does an additional 50% damage to targets with less than 20% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_blessedresillience.jpg",
		},
		[PaladinMajorGlyph.GlyphOfWordOfGlory]: {
			name: "Glyph of Word of Glory",
			description: "Increases your damage by 3% per Holy Power spent after you cast Word of Glory or Eternal Flame on a friendly target. Lasts 6s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_helmet_96.jpg",
		},
		[PaladinMajorGlyph.GlyphOfIllumination]: {
			name: "Glyph of Illumination",
			description: "Your Holy Shock criticals grant 1% mana return, but Holy Insight returns 10% less mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_pureofheart.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHarshWords]: {
			name: "Glyph of Harsh Words",
			description: "Your Word of Glory can now also be used on enemy targets, causing Holy damage approximately equal to the amount it would have healed.\u000D\u000A\u000D\u000A Does not work with Eternal Flame.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_helmet_96.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivinity]: {
			name: "Glyph of Divinity",
			description: "Increases the cooldown of your Lay on Hands by 0 min but causes it to give you 10% of your maximum mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg",
		},
		[PaladinMajorGlyph.GlyphOfLightOfDawn]: {
			name: "Glyph of Light of Dawn",
			description: "Light of Dawn affects 2 fewer targets, but heals each target for 25% more.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_paladin_lightofdawn.jpg",
		},
		[PaladinMajorGlyph.GlyphOfBlessedLife]: {
			name: "Glyph of Blessed Life",
			description: "You have a 50% chance to gain a charge of Holy Power whenever you are affected by a Stun, Fear or Immobilize effect.\u000D\u000A\u000D\u000A This effect cannot occur more than once every 20 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingaura.jpg",
		},
		[PaladinMajorGlyph.GlyphOfFlashOfLight]: {
			name: "Glyph of Flash of Light",
			description: "When you Flash of Light a target, it increases your next heal done to that target within 7s by 10%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_flashheal.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDenounce]: {
			name: "Glyph of Denounce",
			description: "Your Holy Shocks reduce the cast time of your next Denounce by 0.5 sec. This effect stacks up to 3 times.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_purifyingpower.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDazingShield]: {
			name: "Glyph of Dazing Shield",
			description: "Your Avenger\'s Shield now also dazes targets for 10s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfImmediateTruth]: {
			name: "Glyph of Immediate Truth",
			description: "Increases the instant damage done by Seal of Truth by 40%, but decreases the damage done by Censure by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg",
		},
		[PaladinMajorGlyph.GlyphOfBeaconOfLight]: {
			name: "Glyph of Beacon of Light",
			description: "Removes the global cooldown on Beacon of Light.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_beaconoflight.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHammerOfTheRighteous]: {
			name: "Glyph of Hammer of the Righteous",
			description: "The physical damage reduction caused by Hammer of the Righteous now lasts 50% longer.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_hammeroftherighteous.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivineStorm]: {
			name: "Glyph of Divine Storm",
			description: "Your Divine Storm also heals you for 5% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_divinestorm.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTheAlabasterShield]: {
			name: "Glyph of the Alabaster Shield",
			description: "Your successful blocks increase the damage of your next Shield of the Righteous by 10%. Stacks up to 1 times.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/trade_archaeology_stoneshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivinePlea]: {
			name: "Glyph of Divine Plea",
			description: "Divine Plea returns 50% less mana but has a 50% shorter cooldown.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_aspiration.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHolyShock]: {
			name: "Glyph of Holy Shock",
			description: "Decreases the healing of Holy Shock by 50% but increases its damage by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_searinglight.jpg",
		},
		[PaladinMajorGlyph.GlyphOfInquisition]: {
			name: "Glyph of Inquisition",
			description: "When you land a killing blow on an opponent that yields experience or honor, the duration of your Inquisition is increased by 30 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_paladin_inquisition.jpg",
		},
		[PaladinMajorGlyph.GlyphOfProtectorOfTheInnocent]: {
			name: "Glyph of Protector of the Innocent",
			description: "When you use Word of Glory to heal another target, it also heals you for 20% of the amount.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_protectoroftheinnocent.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTheBattleHealer]: {
			name: "Glyph of the Battle Healer",
			description: "Melee attacks from Seal of Insight heal the most wounded member of your raid or party for 30% of the normal heal instead of you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_vindication.jpg",
		},
		[PaladinMajorGlyph.GlyphOfMassExorcism]: {
			name: "Glyph of Mass Exorcism",
			description: "Reduces the range of Exorcism to melee range, but causes 25% damage to all enemies within 8 yards of the primary target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivineShield]: {
			name: "Glyph of Divine Shield",
			description: "Removing harmful effects with Divine Shield heals you for 10% for each effect removed. This heal cannot exceed 50% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_divineshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHandOfSacrifice]: {
			name: "Glyph of Hand of Sacrifice",
			description: "Hand of Sacrifice no longer redirects damage to the Paladin.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofsacrifice.jpg",
		},
	},
	minorGlyphs: {
		[PaladinMinorGlyph.GlyphOfTheLuminousCharger]: {
			name: "Glyph of the Luminous Charger",
			description: "Your Paladin class mounts glow with holy light.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusaderaura.jpg",
		},
		[PaladinMinorGlyph.GlyphOfTheMountedKing]: {
			name: "Glyph of the Mounted King",
			description: "Mounting one of your Paladin class mounts automatically casts Blessing of Kings on you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_greaterblessingofkings.jpg",
		},
		[PaladinMinorGlyph.GlyphOfContemplation]: {
			name: "Glyph of Contemplation",
			description: "Teaches you the ability Contemplation.\u000D\u000A\u000D\u000A Allows you a moment of peace as you kneel in quiet contemplation to ponder the nature of the Light.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_relics_libramofhope.jpg",
		},
		[PaladinMinorGlyph.GlyphOfWingedVengeance]: {
			name: "Glyph of Winged Vengeance",
			description: "Your Avenging Wrath depicts 4 wings.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_revivechampion.jpg",
		},
		[PaladinMinorGlyph.GlyphOfSealOfBlood]: {
			name: "Glyph of Seal of Blood",
			description: "Your Seal of Truth now uses the Seal of Blood visual.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg",
		},
		[PaladinMinorGlyph.GlyphOfFireFromTheHeavens]: {
			name: "Glyph of Fire From the Heavens",
			description: "Your Judgment and Hammer of Wrath criticals call down fire from the sky.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_firestarter.jpg",
		},
		[PaladinMinorGlyph.GlyphOfFocusedWrath]: {
			name: "Glyph of Focused Wrath",
			description: "Holy Wrath only affects one target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_weaponmastery.jpg",
		},
		[PaladinMinorGlyph.GlyphOfTheFallingAvenger]: {
			name: "Glyph of the Falling Avenger",
			description: "You slow fall during Avenging Wrath.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_featherfall.jpg",
		},
		[PaladinMinorGlyph.GlyphOfTheRighteousRetreat]: {
			name: "Glyph of the Righteous Retreat",
			description: "During Divine Shield, you can invoke your Hearthstone 50% faster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_guildperk_hastyhearth.jpg",
		},
		[PaladinMinorGlyph.GlyphOfBladedJudgment]: {
			name: "Glyph of Bladed Judgment",
			description: "Your Judgment spell depicts an axe or sword instead of a hammer, if you have an axe or sword equipped.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_axe_1h_cataclysm_c_01.jpg",
		},
		[PaladinMinorGlyph.GlyphOfRighteousRetreat]: {
			name: "Glyph of Righteous Retreat",
			description: "During Divine Shield, you can invoke your Hearthstone 50% faster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_guildperk_hastyhearth.jpg",
		},
		[PaladinMinorGlyph.GlyphOfTheExorcist]: {
			name: "Glyph of the Exorcist",
			description: "Exorcism will now appear to remove the evil from its target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism.jpg",
		},
		[PaladinMinorGlyph.GlyphOfPillarOfLight]: {
			name: "Glyph of Pillar of Light",
			description: "Critical heals on other players display a small pillar of light at their location briefly.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_surgeoflight.jpg",
		},
	},
};
