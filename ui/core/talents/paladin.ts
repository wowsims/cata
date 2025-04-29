import { PaladinMajorGlyph, PaladinMinorGlyph, PaladinTalents } from '../proto/paladin.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PaladinTalentJson from './trees/paladin.json';

export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

export const paladinGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PaladinMajorGlyph.GlyphOfBeaconOfLight]: {
			name: 'Glyph of Beacon of Light',
			description: 'Your Beacon of Light costs no mana.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_paladin_beaconoflight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfCleansing]: {
			name: 'Glyph of Cleansing',
			description: 'Reduces the mana cost of your Cleanse by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_purify.jpg',
		},
		[PaladinMajorGlyph.GlyphOfConsecration]: {
			name: 'Glyph of Consecration',
			description: 'Increases the duration and cooldown of Consecration by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDazingShield]: {
			name: 'Glyph of Dazing Shield',
			description: "Your Avenger's Shield now also dazes targets.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinePlea]: {
			name: 'Glyph of Divine Plea',
			description: 'Your Divine Plea provides an additional 6% of your total mana.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_aspiration.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivineProtection]: {
			name: 'Glyph of Divine Protection',
			description: 'Removes the physical damage reduction of your Divine Protection, but increases the magical damage reduction by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_divineprotection.jpg',
		},
		[PaladinMajorGlyph.GlyphOfDivinity]: {
			name: 'Glyph of Divinity',
			description: 'When you use Lay on Hands, you also gain 10% of your maximum mana.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMajorGlyph.GlyphOfFocusedShield]: {
			name: 'Glyph of Focused Shield',
			description: "Your Avenger's Shield hits 2 fewer targets, but for 30% more damage.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_avengersshield.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfJustice]: {
			name: 'Glyph of Hammer of Justice',
			description: 'Increases your Hammer of Justice range by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofmight.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHammerOfWrath]: {
			name: 'Glyph of Hammer of Wrath',
			description: 'Reduces the mana cost of Hammer of Wrath by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_hammer_04.jpg',
		},
		[PaladinMajorGlyph.GlyphOfHolyWrath]: {
			name: 'Glyph of Holy Wrath',
			description: 'Your Holy Wrath now also stuns Elementals and Dragonkin.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_purifyingpower.jpg',
		},
		[PaladinMajorGlyph.GlyphOfLayOnHands]: {
			name: 'Glyph of Lay on Hands',
			description: 'Reduces the cooldown of your Lay on Hands spell by 3 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg',
		},
		[PaladinMajorGlyph.GlyphOfLightOfDawn]: {
			name: 'Glyph of Light of Dawn',
			description: 'Light of Dawn affects 2 fewer targets, but heals each target for 25% more.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_paladin_lightofdawn.jpg',
		},
		[PaladinMajorGlyph.GlyphOfRebuke]: {
			name: 'Glyph of Rebuke',
			description: 'Reduces the mana cost of Rebuke by 100%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_rebuke.jpg',
		},
		[PaladinMajorGlyph.GlyphOfReckoning]: {
			name: 'Glyph of Reckoning',
			description: 'Your Hand of Reckoning spell no longer taunts the target and can deal damage to untauntable targets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_questionmark.jpg',
		},
		[PaladinMajorGlyph.GlyphOfSalvation]: {
			name: 'Glyph of Salvation',
			description: 'Hand of Salvation no longer permanently reduces threat over time but instead reduces all threat as long as Hand of Salvation lasts.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofsalvation.jpg',
		},
		[PaladinMajorGlyph.GlyphOfTheAsceticCrusader]: {
			name: 'Glyph of the Ascetic Crusader',
			description: 'Reduces the mana cost of your Crusader Strike by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusaderstrike.jpg',
		},
		[PaladinMajorGlyph.GlyphOfTheLongWord]: {
			name: 'Glyph of the Long Word',
			description: 'Your Word of Glory heals for 50% less up front, but provides an additional 50% healing over 6 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_helmet_96.jpg',
		},
		[PaladinMajorGlyph.GlyphOfTurnEvil]: {
			name: 'Glyph of Turn Evil',
			description: 'Reduces the casting time of your Turn Evil spell by 100%, but increases the cooldown by 8 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_turnundead.jpg',
		},
	},
	minorGlyphs: {
		[PaladinMinorGlyph.GlyphOfBlessingOfKings]: {
			name: 'Glyph of Blessing of Kings',
			description: 'Reduces the mana cost of Blessing of Kings by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_magearmor.jpg',
		},
		[PaladinMinorGlyph.GlyphOfBlessingOfMight]: {
			name: 'Glyph of Blessing of Might',
			description: 'Reduces the mana cost of your Blessing of Might by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_greaterblessingofkings.jpg',
		},
		[PaladinMinorGlyph.GlyphOfInsight]: {
			name: 'Glyph of Insight',
			description: 'Reduces the mana cost of Seal of Insight by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingaura.jpg',
		},
		[PaladinMinorGlyph.GlyphOfJustice]: {
			name: 'Glyph of Justice',
			description: 'Reduces the mana cost of Seal of Justice by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofwrath.jpg',
		},
		[PaladinMinorGlyph.GlyphOfRighteousness]: {
			name: 'Glyph of Righteousness',
			description: 'Reduces the mana cost of Seal of Righteousness by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_righteousnessaura.jpg',
		},
		[PaladinMinorGlyph.GlyphOfTruth]: {
			name: 'Glyph of Truth',
			description: 'Reduces the mana cost of Seal of Truth by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_sealofvengeance.jpg',
		},
	},
};
