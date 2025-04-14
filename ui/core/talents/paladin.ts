import { PaladinMajorGlyph, PaladinMinorGlyph, PaladinPrimeGlyph, PaladinTalents } from '../proto/paladin.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PaladinTalentJson from './trees/paladin.json';export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

export const paladinGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[PaladinPrimeGlyph.GlyphOfJudgement]: {
			name: "Glyph of Judgement",
			description: "Your Judgement deals $54922s1% more damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousfury.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfCrusaderStrike]: {
			name: "Glyph of Crusader Strike",
			description: "Increases the critical strike chance of Crusader Strike by $54927s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusaderstrike.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfExorcism]: {
			name: "Glyph of Exorcism",
			description: "Your Exorcism causes an additional $s1% of its damage over $879d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfWordOfGlory]: {
			name: "Glyph of Word of Glory",
			description: "Increases the healing done by Word of Glory by $54936s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_helmet_96.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfDivineFavor]: {
			name: "Glyph of Divine Favor",
			description: "Increases the duration of Divine Favor by $/1000;s1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_divineillumination.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfSealOfInsight]: {
			name: "Glyph of Seal of Insight",
			description: "While Seal of Insight is active, the effect of your healing spells is increased by $54943s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingaura.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfSealOfTruth]: {
			name: "Glyph of Seal of Truth",
			description: "Your Seal of Truth and Seal of Righteousness also grant $56416s1 expertise while active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfHammerOfTheRighteous]: {
			name: "Glyph of Hammer of the Righteous",
			description: "Increases the damage of both the physical and Holy components of Hammer of the Righteous by $63219s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_hammeroftherighteous.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfTemplarsVerdict]: {
			name: "Glyph of Templar's Verdict",
			description: "Increases the damage of Templar's Verdict by $63220s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_paladin_templarsverdict.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfShieldOfTheRighteous]: {
			name: "Glyph of Shield of the Righteous",
			description: "Increases the damage of Shield of the Righteous by $63222s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_shieldofvengeance.jpg",
		},
		[PaladinPrimeGlyph.GlyphOfHolyShock]: {
			name: "Glyph of Holy Shock",
			description: "Increases the critical effect chance of Holy Shock by $63224s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_searinglight.jpg",
		},
	},
	majorGlyphs: {
		[PaladinMajorGlyph.GlyphOfRebuke]: {
			name: "Glyph of Rebuke",
			description: "Reduces the mana cost of Rebuke by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_rebuke.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHammerOfJustice]: {
			name: "Glyph of Hammer of Justice",
			description: "Increases your Hammer of Justice range by $s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofmight.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivineProtection]: {
			name: "Glyph of Divine Protection",
			description: "Removes the physical damage reduction of your Divine Protection, but increases the magical damage reduction by $s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_divineprotection.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHammerOfWrath]: {
			name: "Glyph of Hammer of Wrath",
			description: "Reduces the mana cost of Hammer of Wrath by ${$54926m1/-1}%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_hammer_04.jpg",
		},
		[PaladinMajorGlyph.GlyphOfConsecration]: {
			name: "Glyph of Consecration",
			description: "Increases the duration and cooldown of Consecration by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg",
		},
		[PaladinMajorGlyph.GlyphOfFocusedShield]: {
			name: "Glyph of Focused Shield",
			description: "Your Avenger's Shield hits $s1 fewer targets, but for $s2% more damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTurnEvil]: {
			name: "Glyph of Turn Evil",
			description: "Reduces the casting time of your Turn Evil spell by $54931s1%, but increases the cooldown by ${$54931m2/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_turnundead.jpg",
		},
		[PaladinMajorGlyph.GlyphOfCleansing]: {
			name: "Glyph of Cleansing",
			description: "Reduces the mana cost of your Cleanse by $54935s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_purify.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTheAsceticCrusader]: {
			name: "Glyph of the Ascetic Crusader",
			description: "Reduces the mana cost of your Crusader Strike by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusaderstrike.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivinity]: {
			name: "Glyph of Divinity",
			description: "When you use Lay on Hands, you also gain $s1% of your maximum mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg",
		},
		[PaladinMajorGlyph.GlyphOfLightOfDawn]: {
			name: "Glyph of Light of Dawn",
			description: "Light of Dawn affects $54940s1 fewer targets, but heals each target for $54940s2% more.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_paladin_lightofdawn.jpg",
		},
		[PaladinMajorGlyph.GlyphOfLayOnHands]: {
			name: "Glyph of Lay on Hands",
			description: "Reduces the cooldown of your Lay on Hands spell by ${$57955m1/-60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg",
		},
		[PaladinMajorGlyph.GlyphOfHolyWrath]: {
			name: "Glyph of Holy Wrath",
			description: "Your Holy Wrath now also stuns Elementals and Dragonkin.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_purifyingpower.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDazingShield]: {
			name: "Glyph of Dazing Shield",
			description: "Your Avenger's Shield now also dazes targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg",
		},
		[PaladinMajorGlyph.GlyphOfBeaconOfLight]: {
			name: "Glyph of Beacon of Light",
			description: "Your Beacon of Light costs no mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_paladin_beaconoflight.jpg",
		},
		[PaladinMajorGlyph.GlyphOfDivinePlea]: {
			name: "Glyph of Divine Plea",
			description: "Your Divine Plea provides an additional 6% of your total mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_aspiration.jpg",
		},
		[PaladinMajorGlyph.GlyphOfSalvation]: {
			name: "Glyph of Salvation",
			description: "Hand of Salvation no longer permanently reduces threat over time but instead reduces all threat as long as Hand of Salvation lasts.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofsalvation.jpg",
		},
		[PaladinMajorGlyph.GlyphOfTheLongWord]: {
			name: "Glyph of the Long Word",
			description: "Your Word of Glory heals for $s1% less up front, but provides an additional $s2% healing over $85673d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_helmet_96.jpg",
		},
	},
	minorGlyphs: {
		[PaladinMinorGlyph.GlyphOfRighteousness]: {
			name: "Glyph of Righteousness",
			description: "Reduces the mana cost of Seal of Righteousness by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousnessaura.jpg",
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfMight]: {
			name: "Glyph of Blessing of Might",
			description: "Reduces the mana cost of your Blessing of Might by $57958s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_greaterblessingofkings.jpg",
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfKings]: {
			name: "Glyph of Blessing of Kings",
			description: "Reduces the mana cost of Blessing of Kings by $57937s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_magearmor.jpg",
		},
		[PaladinMinorGlyph.GlyphOfInsight]: {
			name: "Glyph of Insight",
			description: "Reduces the mana cost of Seal of Insight by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingaura.jpg",
		},
		[PaladinMinorGlyph.GlyphOfTruth]: {
			name: "Glyph of Truth",
			description: "Reduces the mana cost of Seal of Truth by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg",
		},
		[PaladinMinorGlyph.GlyphOfJustice]: {
			name: "Glyph of Justice",
			description: "Reduces the mana cost of Seal of Justice by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofwrath.jpg",
		},
	},
};
