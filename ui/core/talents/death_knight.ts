import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightTalents } from '../proto/death_knight.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DeathKnightTalentJson from './trees/death_knight.json';export const deathKnightTalentsConfig: TalentsConfig<DeathKnightTalents> = newTalentsConfig(DeathKnightTalentJson);

export const deathKnightGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DeathKnightMajorGlyph.GlyphOfAntiMagicShell]: {
			name: "Glyph of Anti-Magic Shell",
			description: "Causes your Anti-Magic Shell to absorb all incoming magical damage, up to the absorption limit.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antimagicshell.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfUnholyFrenzy]: {
			name: "Glyph of Unholy Frenzy",
			description: "Causes your Unholy Frenzy to no longer deal damage to the affected target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_unholyfrenzy.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfIceboundFortitude]: {
			name: "Glyph of Icebound Fortitude",
			description: "Reduces the cooldown of your Icebound Fortitude by 50%, but also reduces its duration by 75%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_iceboundfortitude.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfChainsOfIce]: {
			name: "Glyph of Chains of Ice",
			description: "Your Chains of Ice also causes 143 Frost damage, with additional damage depending on your attack power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDeathGrip]: {
			name: "Glyph of Death Grip",
			description: "Increases the range of your Death Grip ability by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_strangulate.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDeathAndDecay]: {
			name: "Glyph of Death and Decay",
			description: "Your Death and Decay also reduces the movement speed of enemies within its radius by 50%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathanddecay.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfShiftingPresences]: {
			name: "Glyph of Shifting Presences",
			description: "You retain 70% of your Runic Power when switching Presences.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_bloodpresence.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfIcyTouch]: {
			name: "Glyph of Icy Touch",
			description: "Your Icy Touch dispels one helpful Magic effect from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_icetouch.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfEnduringInfection]: {
			name: "Glyph of Enduring Infection",
			description: "Your diseases are undispellable, but their damage dealt is reduced by 15%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_creature_disease_05.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfPestilence]: {
			name: "Glyph of Pestilence",
			description: "Increases the radius of your Pestilence effect by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_plaguecloud.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfMindFreeze]: {
			name: "Glyph of Mind Freeze",
			description: "Reduces the cooldown of your Mind Freeze ability by 1 sec, but also raises its cost by 10 Runic Power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_mindfreeze.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfStrangulate]: {
			name: "Glyph of Strangulate",
			description: "Increases the Silence duration of your Strangulate ability by 2 sec when used on a target who is casting a spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soulleech_3.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfPillarOfFrost]: {
			name: "Glyph of Pillar of Frost",
			description: "Empowers your Pillar of Frost, making you immune to all effects that cause loss of control of your character, but also reduces your movement speed by 70% while the ability is active.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_deathknight_pillaroffrost.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfVampiricBlood]: {
			name: "Glyph of Vampiric Blood",
			description: "Increases the bonus healing received while your Vampiric Blood is active by an additional 15%, but your Vampiric Blood no longer grants you health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfUnholyCommand]: {
			name: "Glyph of Unholy Command",
			description: "Immediately finishes the cooldown of your Death Grip upon dealing a killing blow to a target that grants experience or honor.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_skull.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfOutbreak]: {
			name: "Glyph of Outbreak",
			description: "Your Outbreak spell no longer has a cooldown, but now costs 30 Runic Power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathvortex.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon]: {
			name: "Glyph of Dancing Rune Weapon",
			description: "Increases your threat generation by 100% while your Dancing Rune Weapon is active, but reduces its damage dealt by 25%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_sword_07.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDarkSimulacrum]: {
			name: "Glyph of Dark Simulacrum",
			description: "Reduces the cooldown of Dark Simulacrum by 30 sec and increases its duration by 4 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_consumemagic.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDeathCoil]: {
			name: "Glyph of Death Coil",
			description: "Your Death Coil spell is now usable on all allies. When cast on a non-undead ally, Death Coil shrouds them with a protective barrier that absorbs up to 168 damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfDarkSuccor]: {
			name: "Glyph of Dark Succor",
			description: "When you kill an enemy that yields experience or honor, while in Frost or Unholy Presence, your next Death Strike within 15s is free and will restore at least 20% of your maximum health.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_butcher2.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfSwiftDeath]: {
			name: "Glyph of Swift Death",
			description: "The haste effect granted by Soul Reaper now also increases your movement speed by 30% for the duration.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_deathknight_soulreaper.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfLoudHorn]: {
			name: "Glyph of Loud Horn",
			description: "Your Horn of Winter now generates an additional 10 Runic Power, but the cooldown is increased by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_horn_04.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfRegenerativeMagic]: {
			name: "Glyph of Regenerative Magic",
			description: "If Anti-Magic Shell expires after its full duration, the cooldown is reduced by up to 50%, based on the amount of damage absorbtion remaining.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antimagicshell.jpg",
		},
		[DeathKnightMajorGlyph.GlyphOfFesteringBlood]: {
			name: "Glyph of Festering Blood",
			description: "Blood Boil will now treat all targets as though they have Blood Plague or Frost Fever applied.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_yorsahj_bloodboil_green.jpg",
		},
	},
	minorGlyphs: {
		[DeathKnightMinorGlyph.GlyphOfTheGeist]: {
			name: "Glyph of the Geist",
			description: "Your Raise Dead spell summons a geist instead of a ghoul.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfDeathsEmbrace]: {
			name: "Glyph of Death's Embrace",
			description: "Your Death Coil refunds 20 Runic Power when used to heal an allied minion, but will no longer trigger Blood Tap when used this way.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfHornOfWinter]: {
			name: "Glyph of Horn of Winter",
			description: "When used outside of combat, your Horn of Winter ability causes a brief, localized snow flurry.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_horn_02.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfArmyOfTheDead]: {
			name: "Glyph of Army of the Dead",
			description: "The ghouls summoned by your Army of the Dead no longer taunt their target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_armyofthedead.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfFoulMenagerie]: {
			name: "Glyph of Foul Menagerie",
			description: "Causes your Army of the Dead spell to summon an assortment of undead minions.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_armyofthedead.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfPathOfFrost]: {
			name: "Glyph of Path of Frost",
			description: "Your Path of Frost ability allows you to fall from a greater distance without suffering damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_pathoffrost.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfResilientGrip]: {
			name: "Glyph of Resilient Grip",
			description: "When your Death Grip ability fails because its target is immune, its cooldown is reset.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/warlock_curse_shadow.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfDeathGate]: {
			name: "Glyph of Death Gate",
			description: "Reduces the cast time of your Death Gate spell by 60%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_teleportundercity.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfCorpseExplosion]: {
			name: "Glyph of Corpse Explosion",
			description: "Teaches you the ability Corpse Explosion.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_corpseexplode.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfTranquilGrip]: {
			name: "Glyph of Tranquil Grip",
			description: "Your Death Grip spell no longer taunts the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_rogue_envelopingshadows.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfTheSkeleton]: {
			name: "Glyph of the Skeleton",
			description: "Your Raise Dead spell summons a skeleton instead of a ghoul.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_bone_humanskull_01.jpg",
		},
		[DeathKnightMinorGlyph.GlyphOfTheLongWinter]: {
			name: "Glyph of the Long Winter",
			description: "The effect of your Horn of Winter now lasts for 1 hour.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_deathknight_remorselesswinters.jpg",
		},
	},
};
