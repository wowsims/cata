import { DeathKnightMajorGlyph, DeathKnightMinorGlyph, DeathKnightTalents } from '../proto/death_knight';
import { GlyphsConfig } from './glyphs_picker.jsx';
import { newTalentsConfig, TalentsConfig } from './talents_picker.jsx';
import DkTalentsJson from './trees/death_knight.json';

export const deathknightTalentsConfig: TalentsConfig<DeathKnightTalents> = newTalentsConfig(DkTalentsJson);

export const deathknightGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[DeathKnightMajorGlyph.GlyphOfAntiMagicShell]: {
			name: 'Glyph of Anti-Magic Shell',
			description: 'Increases the duration of your Anti-Magic Shell by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_antimagicshell.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfBloodStrike]: {
			name: 'Glyph of Blood Strike',
			description: 'Your Blood Strike causes an additional 20% damage to snared targets.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_deathstrike.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfBoneShield]: {
			name: 'Glyph of Bone Shield',
			description: 'Adds 1 additional charge to your Bone Shield.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_chest_leather_13.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfChainsOfIce]: {
			name: 'Glyph of Chains of Ice',
			description: 'Your Chains of Ice also causes 144 to 156 Frost damage, increased by your attack power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_chainsofice.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDancingRuneWeapon]: {
			name: 'Glyph of Dancing Rune Weapon',
			description: 'Increases the duration of Dancing Rune Weapon by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_07.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDarkCommand]: {
			name: 'Glyph of Dark Command',
			description: 'Increases the chance for your Dark Command ability to work successfully by 8%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_shamanrage.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDarkDeath]: {
			name: 'Glyph of Dark Death',
			description: 'Increases the damage or healing done by Death Coil by 15%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDeathAndDecay]: {
			name: 'Glyph of Death and Decay',
			description: 'Damage of your Death and Decay spell increased by 20%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathanddecay.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDeathGrip]: {
			name: 'Glyph of Death Grip',
			description: 'When you deal a killing blow that grants honor or experience, the cooldown of your Death Grip is refreshed.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_strangulate.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDeathStrike]: {
			name: 'Glyph of Death Strike',
			description:
				"Increases your Death Strike's damage by 1% for every 1 runic power you currently have (up to a maximum of 25%). The runic power is not consumed by this effect.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_butcher2.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfDisease]: {
			name: 'Glyph of Disease',
			description:
				'Your Pestilence ability now refreshes disease durations and secondary effects of diseases on your primary target back to their maximum duration.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfFrostStrike]: {
			name: 'Glyph of Frost Strike',
			description: 'Reduces the cost of your Frost Strike by 8 Runic Power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_empowerruneblade2.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfHeartStrike]: {
			name: 'Glyph of Heart Strike',
			description: 'Your Heart Strike also reduces the movement speed of your target by 50% for 10 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_weapon_shortblade_40.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfHowlingBlast]: {
			name: 'Glyph of Howling Blast',
			description: 'Your Howling Blast ability now infects your targets with Frost Fever.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_arcticwinds.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfHungeringCold]: {
			name: 'Glyph of Hungering Cold',
			description: 'Reduces the cost of Hungering Cold by 40 runic power.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_staff_15.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfIceboundFortitude]: {
			name: 'Glyph of Icebound Fortitude',
			description: 'Your Icebound Fortitude now always grants at least 40% damage reduction, regardless of your defense skill.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_iceboundfortitude.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfIcyTouch]: {
			name: 'Glyph of Icy Touch',
			description: 'Your Frost Fever disease deals 20% additional damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_icetouch.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfObliterate]: {
			name: 'Glyph of Obliterate',
			description: 'Increases the damage of your Obliterate ability by 25%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_classicon.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfPlagueStrike]: {
			name: 'Glyph of Plague Strike',
			description: 'Your Plague Strike does 20% additional damage.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_empowerruneblade.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfRuneStrike]: {
			name: 'Glyph of Rune Strike',
			description: 'Increases the critical strike chance of your Rune Strike by 10%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_darkconviction.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfRuneTap]: {
			name: 'Glyph of Rune Tap',
			description: 'Your Rune Tap now heals you for an additional 1% of your maximum health, and also heals your party for 10% of their maximum health.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_runetap.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfScourgeStrike]: {
			name: 'Glyph of Scourge Strike',
			description: 'Your Scourge Strike increases the duration of your diseases on the target by 3 sec, up to a maximum of 9 additional seconds.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_scourgestrike.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfStrangulate]: {
			name: 'Glyph of Strangulate',
			description: 'Reduces the cooldown of your Strangulate by 20 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_soulleech_3.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfTheGhoul]: {
			name: 'Glyph of the Ghoul',
			description: 'Your Ghoul receives an additional 40% of your Strength and 40% of your Stamina.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfUnbreakableArmor]: {
			name: 'Glyph of Unbreakable Armor',
			description: 'Increases the total armor granted by Unbreakable Armor to 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_armor_helm_plate_naxxramas_raidwarrior_c_01.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfUnholyBlight]: {
			name: 'Glyph of Unholy Blight',
			description: 'Increases the damage done by Unholy Blight by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_contagion.jpg',
		},
		[DeathKnightMajorGlyph.GlyphOfVampiricBlood]: {
			name: 'Glyph of Vampiric Blood',
			description: 'Increases the duration of your Vampiric Blood by 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_lifedrain.jpg',
		},
	},
	minorGlyphs: {
		[DeathKnightMinorGlyph.GlyphOfBloodTap]: {
			name: 'Glyph of Blood Tap',
			description: 'Your Blood Tap no longer causes damage to you.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_deathknight_bloodtap.jpg',
		},
		[DeathKnightMinorGlyph.GlyphOfCorpseExplosion]: {
			name: 'Glyph of Corpse Explosion',
			description: 'Increases the radius of effect on Corpse Explosion by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_creature_disease_02.jpg',
		},
		[DeathKnightMinorGlyph.GlyphOfDeathSEmbrace]: {
			name: "Glyph of Death's Embrace",
			description: 'Your Death Coil refunds 20 runic power when used to heal.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_deathcoil.jpg',
		},
		[DeathKnightMinorGlyph.GlyphOfHornOfWinter]: {
			name: 'Glyph of Horn of Winter',
			description: 'Increases the duration of your Horn of Winter ability by 1 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_horn_02.jpg',
		},
		[DeathKnightMinorGlyph.GlyphOfPestilence]: {
			name: 'Glyph of Pestilence',
			description: 'Increases the radius of your Pestilence effect by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_plaguecloud.jpg',
		},
		[DeathKnightMinorGlyph.GlyphOfRaiseDead]: {
			name: 'Glyph of Raise Dead',
			description: 'Your Raise Dead spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg',
		},
	},
};
