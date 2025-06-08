import { PriestMajorGlyph, PriestMinorGlyph, PriestPrimeGlyph, PriestTalents } from '../proto/priest.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PriestTalentJson from './trees/priest.json';export const priestTalentsConfig: TalentsConfig<PriestTalents> = newTalentsConfig(PriestTalentJson);

export const priestGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[PriestPrimeGlyph.GlyphOfFlashHeal]: {
			name: "Glyph of Flash Heal",
			description: "Increases the critical effect chance of your Flash Heal on targets at or below $55679s2% health by $55679s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_flashheal.jpg",
		},
		[PriestPrimeGlyph.GlyphOfLightwell]: {
			name: "Glyph of Lightwell",
			description: "Increases the total amount of charges of your Lightwell by $55673s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_summonlightwell.jpg",
		},
		[PriestPrimeGlyph.GlyphOfShadowWordPain]: {
			name: "Glyph of Shadow Word: Pain",
			description: "Increases the periodic damage of your Shadow Word: Pain by $55681s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowwordpain.jpg",
		},
		[PriestPrimeGlyph.GlyphOfPowerWordBarrier]: {
			name: "Glyph of Power Word: Barrier",
			description: "Increases the healing received while under Power Word: Barrier by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_powerwordbarrier.jpg",
		},
		[PriestPrimeGlyph.GlyphOfPowerWordShield]: {
			name: "Glyph of Power Word: Shield",
			description: "Your Power Word: Shield also heals the target for $55672s1% of the absorption amount.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_powerwordshield.jpg",
		},
		[PriestPrimeGlyph.GlyphOfPrayerOfHealing]: {
			name: "Glyph of Prayer of Healing",
			description: "Your Prayer of Healing spell also heals an additional $55680s1% of its initial heal over 6 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_prayerofhealing02.jpg",
		},
		[PriestPrimeGlyph.GlyphOfRenew]: {
			name: "Glyph of Renew",
			description: "Increases the amount healed by your Renew by an additional $55674s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_renew.jpg",
		},
		[PriestPrimeGlyph.GlyphOfShadowWordDeath]: {
			name: "Glyph of Shadow Word: Death",
			description: "If your Shadow Word: Death fails to kill the target at or below $s1% health, your Shadow Word: Death's cooldown is instantly reset.  This effect cannot occur more often than once every $95652d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonicfortitude.jpg",
		},
		[PriestPrimeGlyph.GlyphOfMindFlay]: {
			name: "Glyph of Mind Flay",
			description: "Increases the damage done by your Mind Flay spell by $55687s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_siphonmana.jpg",
		},
		[PriestPrimeGlyph.GlyphOfDispersion]: {
			name: "Glyph of Dispersion",
			description: "Reduces the cooldown on Dispersion by ${$63229m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_dispersion.jpg",
		},
		[PriestPrimeGlyph.GlyphOfGuardianSpirit]: {
			name: "Glyph of Guardian Spirit",
			description: "Reduces the cooldown of your Guardian Spirit by $/1000;63231S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_guardianspirit.jpg",
		},
		[PriestPrimeGlyph.GlyphOfPenance]: {
			name: "Glyph of Penance",
			description: "Reduces the cooldown of Penance by ${$63235m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_penance.jpg",
		},
	},
	majorGlyphs: {
		[PriestMajorGlyph.GlyphOfCircleOfHealing]: {
			name: "Glyph of Circle of Healing",
			description: "Your Circle of Healing spell heals $55675s1 additional target, but its mana cost is increased by $55675s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_circleofrenewal.jpg",
		},
		[PriestMajorGlyph.GlyphOfDispelMagic]: {
			name: "Glyph of Dispel Magic",
			description: "Your Dispel Magic spell also heals your target for $55677s1% of maximum health when you successfully dispel a magical effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_dispelmagic.jpg",
		},
		[PriestMajorGlyph.GlyphOfFade]: {
			name: "Glyph of Fade",
			description: "Reduces the cooldown of your Fade spell by ${$55684m2/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
		},
		[PriestMajorGlyph.GlyphOfFearWard]: {
			name: "Glyph of Fear Ward",
			description: "Reduces the cooldown and duration of Fear Ward by ${$55678m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism.jpg",
		},
		[PriestMajorGlyph.GlyphOfHolyNova]: {
			name: "Glyph of Holy Nova",
			description: "Reduces the global cooldown of your Holy Nova by .5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_holynova.jpg",
		},
		[PriestMajorGlyph.GlyphOfInnerFire]: {
			name: "Glyph of Inner Fire",
			description: "Increases the armor gained from your Inner Fire spell by $55686m1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg",
		},
		[PriestMajorGlyph.GlyphOfMassDispel]: {
			name: "Glyph of Mass Dispel",
			description: "Reduces the cast time of Mass Dispel by $/1000;55691S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_massdispel.jpg",
		},
		[PriestMajorGlyph.GlyphOfPsychicHorror]: {
			name: "Glyph of Psychic Horror",
			description: "Reduces the cooldown of your Psychic Horror by $/1000;S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_psychichorrors.jpg",
		},
		[PriestMajorGlyph.GlyphOfPsychicScream]: {
			name: "Glyph of Psychic Scream",
			description: "Targets of your Psychic Scream spell now tremble in place instead of fleeing in fear, but the cooldown of Psychic Scream is increased by $/1000;S2 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_psychicscream.jpg",
		},
		[PriestMajorGlyph.GlyphOfScourgeImprisonment]: {
			name: "Glyph of Scourge Imprisonment",
			description: "Reduces the cast time of your Shackle Undead by ${$55690m1/-1000}.1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusade.jpg",
		},
		[PriestMajorGlyph.GlyphOfSmite]: {
			name: "Glyph of Smite",
			description: "Your Smite spell inflicts an additional $55692s1% damage against targets afflicted by Holy Fire.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_holysmite.jpg",
		},
		[PriestMajorGlyph.GlyphOfPrayerOfMending]: {
			name: "Glyph of Prayer of Mending",
			description: "The first charge of your Prayer of Mending heals for an additional $55685s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_prayerofmendingtga.jpg",
		},
		[PriestMajorGlyph.GlyphOfSpiritTap]: {
			name: "Glyph of Spirit Tap",
			description: "When you kill a target with your Shadow Word: Death that yields experience or honor, you receive $*4;81301s1% of your total mana over $81301d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_requiem.jpg",
		},
		[PriestMajorGlyph.GlyphOfDivineAccuracy]: {
			name: "Glyph of Divine Accuracy",
			description: "Increases your chance to hit with your Smite and Holy Fire spells by $63246s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_holyprotection.jpg",
		},
		[PriestMajorGlyph.GlyphOfDesperation]: {
			name: "Glyph of Desperation",
			description: "Allows Pain Suppression and Guardian Spirit to be cast while stunned.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_painsupression.jpg",
		},
	},
	minorGlyphs: {
		[PriestMinorGlyph.GlyphOfFading]: {
			name: "Glyph of Fading",
			description: "Reduces the mana cost of your Fade spell by $57985s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_invisibility.jpg",
		},
		[PriestMinorGlyph.GlyphOfLevitate]: {
			name: "Glyph of Levitate",
			description: "Your Levitate spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_layonhands.jpg",
		},
		[PriestMinorGlyph.GlyphOfFortitude]: {
			name: "Glyph of Fortitude",
			description: "Reduces the mana cost of your Power Word: Fortitude spell by $58009s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_wordfortitude.jpg",
		},
		[PriestMinorGlyph.GlyphOfShadowProtection]: {
			name: "Glyph of Shadow Protection",
			description: "Increases the duration of your Shadow Protection spell by ${$58015m1/60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_prayerofshadowprotection.jpg",
		},
		[PriestMinorGlyph.GlyphOfShackleUndead]: {
			name: "Glyph of Shackle Undead",
			description: "Increases the range of your Shackle Undead spell by $57986s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg",
		},
		[PriestMinorGlyph.GlyphOfShadowfiend]: {
			name: "Glyph of Shadowfiend",
			description: "Receive $58227s1% of your maximum mana if your Shadowfiend dies from damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowfiend.jpg",
		},
		[PriestMinorGlyph.GlyphOfShadow]: {
			name: "Glyph of Shadow",
			description: "Alters the appearance of your Shadowform.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowform.jpg",
		},
	},
};
