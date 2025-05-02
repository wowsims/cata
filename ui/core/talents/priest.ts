import { PriestMajorGlyph, PriestMinorGlyph, PriestTalents } from '../proto/priest.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PriestTalentJson from './trees/priest.json';export const priestTalentsConfig: TalentsConfig<PriestTalents> = newTalentsConfig(PriestTalentJson);

export const priestGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[PriestMajorGlyph.GlyphOfCircleOfHealing]: {
			name: "Glyph of Circle of Healing",
			description: "Your Circle of Healing spell heals $55675s1 additional target, but its mana cost is increased by $55675s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_circleofrenewal.jpg",
		},
		[PriestMajorGlyph.GlyphOfPurify]: {
			name: "Glyph of Purify",
			description: "Your Purify spell also heals your target for $55677s1% of maximum health when you successfully dispel a magical effect or disease.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_dispelmagic.jpg",
		},
		[PriestMajorGlyph.GlyphOfFade]: {
			name: "Glyph of Fade",
			description: "Your Fade ability now also reduces all damage taken by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
		},
		[PriestMajorGlyph.GlyphOfFearWard]: {
			name: "Glyph of Fear Ward",
			description: "Reduces the cooldown of Fear Ward by ${$55678m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism.jpg",
		},
		[PriestMajorGlyph.GlyphOfInnerSanctum]: {
			name: "Glyph of Inner Sanctum",
			description: "Spell damage taken is reduced by $s1% while within Inner Fire, and the movement speed bonus of your Inner Will is increased by $s2%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/priest_icon_innewill.jpg",
		},
		[PriestMajorGlyph.GlyphOfHolyNova]: {
			name: "Glyph of Holy Nova",
			description: "Teaches you the ability Holy Nova.\u000D\u000A\u000D\u000ACauses an explosion of holy light around the caster, causing $132157s1 Holy damage to all enemy targets within $132157A1 yards and healing up to 5 targets within $23455A1 yards for $23455s1. \u000D\u000A\u000D\u000AHealing is divided among the number of targets healed.  These effects cause no threat.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_holynova.jpg",
		},
		[PriestMajorGlyph.GlyphOfInnerFire]: {
			name: "Glyph of Inner Fire",
			description: "Increases the armor gained from your Inner Fire spell by $55686m1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_innerfire.jpg",
		},
		[PriestMajorGlyph.GlyphOfDeepWells]: {
			name: "Glyph of Deep Wells",
			description: "Increases the total amount of charges of your Lightwell by $55673s1.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_summonlightwell.jpg",
		},
		[PriestMajorGlyph.GlyphOfMassDispel]: {
			name: "Glyph of Mass Dispel",
			description: "Causes your Mass Dispel to be potent enough to remove Magic effects that are normally undispellable.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_massdispel.jpg",
		},
		[PriestMajorGlyph.GlyphOfPsychicHorror]: {
			name: "Glyph of Psychic Horror",
			description: "Reduces the cooldown of your Psychic Horror by $/1000;S1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_psychichorrors.jpg",
		},
		[PriestMajorGlyph.GlyphOfHolyFire]: {
			name: "Glyph of Holy Fire",
			description: "Increases the range of your Holy Fire, Smite, and Power Word: Solace spells by $s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_searinglight.jpg",
		},
		[PriestMajorGlyph.GlyphOfWeakenedSoul]: {
			name: "Glyph of Weakened Soul",
			description: "Reduces the duration of the Weakened Soul effect caused by Power Word: Shield by ${$m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_ashestoashes.jpg",
		},
		[PriestMajorGlyph.GlyphOfPowerWordShield]: {
			name: "Glyph of Power Word: Shield",
			description: "$55672s1% of the absorb from your Power Word: Shield spell is converted into healing.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_powerwordshield.jpg",
		},
		[PriestMajorGlyph.GlyphOfSpiritOfRedemption]: {
			name: "Glyph of Spirit of Redemption",
			description: "Increases the duration of Spirit of Redemption by ${$m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_enchant_essenceeternallarge.jpg",
		},
		[PriestMajorGlyph.GlyphOfPsychicScream]: {
			name: "Glyph of Psychic Scream",
			description: "Targets of your Psychic Scream and your Psyfiend\'s Psychic Terror now tremble in place instead of fleeing in fear.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_psychicscream.jpg",
		},
		[PriestMajorGlyph.GlyphOfRenew]: {
			name: "Glyph of Renew",
			description: "Your Renew heals for $s1% more each time it heals, but its duration is reduced by ${$m2/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_renew.jpg",
		},
		[PriestMajorGlyph.GlyphOfScourgeImprisonment]: {
			name: "Glyph of Scourge Imprisonment",
			description: "Reduces the cast time of your Shackle Undead by ${$55690m1/-1000}.1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_crusade.jpg",
		},
		[PriestMajorGlyph.GlyphOfMindBlast]: {
			name: "Glyph of Mind Blast",
			description: "When you critically hit with your Mind Blast, you cause the target to be unable to move for $87194d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unholyfrenzy.jpg",
		},
		[PriestMajorGlyph.GlyphOfDispelMagic]: {
			name: "Glyph of Dispel Magic",
			description: "Your Dispel Magic spell also damages your target for $119856s1 Holy damage when you successfully dispel a magical effect.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_nullifydisease.jpg",
		},
		[PriestMajorGlyph.GlyphOfSmite]: {
			name: "Glyph of Smite",
			description: "Your Smite spell inflicts an additional $55692s1% damage against targets afflicted by $?s139139[Power Word: Solace][Holy Fire], but that additional damage does not get transferred by Atonement.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_holysmite.jpg",
		},
		[PriestMajorGlyph.GlyphOfPrayerOfMending]: {
			name: "Glyph of Prayer of Mending",
			description: "The first charge of your Prayer of Mending heals for an additional $55685s1% but your Prayer of Mending has $s2 fewer $Lcharge:charges;.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_prayerofmendingtga.jpg",
		},
		[PriestMajorGlyph.GlyphOfLevitate]: {
			name: "Glyph of Levitate",
			description: "Increases your movement speed while Levitating and for 10 sec afterward by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_priest_pathofdevout.jpg",
		},
		[PriestMajorGlyph.GlyphOfReflectiveShield]: {
			name: "Glyph of Reflective Shield",
			description: "Causes $s1% of the damage you absorb with Power Word: Shield to reflect back at the attacker.  This damage causes no threat.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_priest_reflectiveshield.jpg",
		},
		[PriestMajorGlyph.GlyphOfDispersion]: {
			name: "Glyph of Dispersion",
			description: "Reduces the cooldown on Dispersion by ${$63229m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_dispersion.jpg",
		},
		[PriestMajorGlyph.GlyphOfLeapOfFaith]: {
			name: "Glyph of Leap of Faith",
			description: "Your Leap of Faith spell now also clears all movement impairing effects from your target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/priest_spell_leapoffaith_a.jpg",
		},
		[PriestMajorGlyph.GlyphOfPenance]: {
			name: "Glyph of Penance",
			description: "Increases the mana cost of Penance by $s1% but allows Penance to be cast while moving.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_penance.jpg",
		},
		[PriestMajorGlyph.GlyphOfFocusedMending]: {
			name: "Glyph of Focused Mending",
			description: "Causes your Prayer of Mending to only bounce between the target and the caster.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_healingfocus.jpg",
		},
		[PriestMajorGlyph.GlyphOfMindSpike]: {
			name: "Glyph of Mind Spike",
			description: "Your successful non-instant Mind Spikes, reduce the cast time of your next Mind Blast within $81292d by $81292s1%. This effect can stack up to $81292u times.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_priest_mindspike.jpg",
		},
		[PriestMajorGlyph.GlyphOfBindingHeal]: {
			name: "Glyph of Binding Heal",
			description: "Your Binding Heal spell now heals a third friendly target within 20 yards, but costs $s1% more mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_blindingheal.jpg",
		},
		[PriestMajorGlyph.GlyphOfMindFlay]: {
			name: "Glyph of Mind Flay",
			description: "Your Mind Flay spell no longer slows your victim\'s movement speed.  Instead, each time Mind Flay deals damage you will be granted $120587s1% increased movement speed for $120587d, stacking up to $120587u times.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_siphonmana.jpg",
		},
		[PriestMajorGlyph.GlyphOfShadowWordDeath]: {
			name: "Glyph of Shadow Word: Death",
			description: "Your Shadow Word: Death can now be cast at any time, but deals 25% damage against targets above 20% health$?s15407[ and does not generate a Shadow Orb when used against them][].  \u000D\u000A\u000D\u000ACasting Shadow Word: Death now also does damage to you equivalent to the damage it would do to an enemy above 20% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonicfortitude.jpg",
		},
		[PriestMajorGlyph.GlyphOfVampiricEmbrace]: {
			name: "Glyph of Vampiric Embrace",
			description: "Your Vampiric Embrace converts an additional $s2% of the damage you deal into healing, but the duration is reduced by ${$m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unsummonbuilding.jpg",
		},
		[PriestMajorGlyph.ZzoldGlyphOfVampiricTouch]: {
			name: "zzOLD Glyph of Vampiric Touch",
			description: "Your Shadow Word: Death can now be cast at any time, but deals 25% damage against targets above 20% health$?s15407[ and does not generate a Shadow Orb when used against them][].  \u000D\u000A\u000D\u000ACasting Shadow Word: Death now also does damage to you equivalent to the damage it would do to an enemy above 20% health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_demonicfortitude.jpg",
		},
		[PriestMajorGlyph.GlyphOfLightspring]: {
			name: "Glyph of Lightspring",
			description: "Your Lightwell no longer automatically heals nearby targets, but can be clicked by players to deal $s2% more healing than normal.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_summonlightwell.jpg",
		},
		[PriestMajorGlyph.GlyphOfLightwell]: {
			name: "Glyph of Lightwell",
			description: "Your Lightwell no longer automatically heals nearby targets, but can be clicked by players to deal $s2% more healing than normal.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_summonlightwell.jpg",
		},
	},
	minorGlyphs: {
		[PriestMinorGlyph.GlyphOfShadowRavens]: {
			name: "Glyph of Shadow Ravens",
			description: "Your Shadow Orbs now appear as Shadow Ravens.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_priest_shadoworbs.jpg",
		},
		[PriestMinorGlyph.GlyphOfBorrowedTime]: {
			name: "Glyph of Borrowed Time",
			description: "Your Borrowed Time is now displayed visually.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_borrowedtime.jpg",
		},
		[PriestMinorGlyph.GlyphOfShackleUndead]: {
			name: "Glyph of Shackle Undead",
			description: "Changes the appearance of your Shackle Undead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg",
		},
		[PriestMinorGlyph.GlyphOfDarkArchangel]: {
			name: "Glyph of Dark Archangel",
			description: "When you apply Devouring Plague to a target, you take on the form of a Dark Archangel for $129208d.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonvoidwalker.jpg",
		},
		[PriestMinorGlyph.GlyphOfShadow]: {
			name: "Glyph of Shadow",
			description: "Alters the appearance of your Shadowform to be less transparent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_shadowform.jpg",
		},
		[PriestMinorGlyph.GlyphOfTheHeavens]: {
			name: "Glyph of the Heavens",
			description: "Your Levitate targets will appear to be riding on a cloud for the duration of the spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_elemental_primal_air.jpg",
		},
		[PriestMinorGlyph.GlyphOfConfession]: {
			name: "Glyph of Confession",
			description: "Teaches you the ability Confession.\u000D\u000A\u000D\u000ACompels a friendly target to confess a secret.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soothingkiss.jpg",
		},
		[PriestMinorGlyph.GlyphOfHolyResurrection]: {
			name: "Glyph of Holy Resurrection",
			description: "Your resurrection target appears bathed in holy light for the duration of the Resurrection cast.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_surgeoflight.jpg",
		},
		[PriestMinorGlyph.GlyphOfTheValkyr]: {
			name: "Glyph of the Val'kyr",
			description: "While Spirit of Redemption is active, you now appear as a Val\'kyr.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/achievement_boss_svalasorrowgrave.jpg",
		},
		[PriestMinorGlyph.GlyphOfShadowyFriends]: {
			name: "Glyph of Shadowy Friends",
			description: "Your Shadowform extends to your non-combat pets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antishadow.jpg",
		},
		[PriestMinorGlyph.GlyphOfAngels]: {
			name: "Glyph of Angels",
			description: "Your heal spells momentarily grant you angelic wings.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_priest_archangel.jpg",
		},
		[PriestMinorGlyph.GlyphOfTheSha]: {
			name: "Glyph of the Sha",
			description: "Transforms your Shadowfiend and Mindbender into a Sha Beast.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_hand_1h_shaclaw.jpg",
		},
		[PriestMinorGlyph.GlyphOfShiftedAppearances]: {
			name: "Glyph of Shifted Appearances",
			description: "Void Shift causes you and your target to exchange appearances for several seconds.\u000D\u000A\u000D\u000ADoes not affect mounted players.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_priest_voidshift.jpg",
		},
		[PriestMinorGlyph.GlyphOfInspiredHymns]: {
			name: "Glyph of Inspired Hymns",
			description: "While channeling Hymns, a spirit appears above you.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_divinehymn.jpg",
		},
	},
};
