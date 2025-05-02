import { MageMajorGlyph, MageMinorGlyph, MagePrimeGlyph, MageTalents } from '../proto/mage.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MageTalentJson from './trees/mage.json';export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

export const mageGlyphsConfig: GlyphsConfig = {
	primeGlyphs: {
		[MagePrimeGlyph.GlyphOfArcaneMissiles]: {
			name: "Glyph of Arcane Missiles",
			description: "Increases the critical strike chance of your Arcane Missiles spell by $56363s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_starfall.jpg",
		},
		[MagePrimeGlyph.GlyphOfFireball]: {
			name: "Glyph of Fireball",
			description: "Increases the critical strike chance of your Fireball spell by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_flamebolt.jpg",
		},
		[MagePrimeGlyph.GlyphOfFrostbolt]: {
			name: "Glyph of Frostbolt",
			description: "Increases the critical strike chance of your Frostbolt spell by $56370s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostbolt02.jpg",
		},
		[MagePrimeGlyph.GlyphOfPyroblast]: {
			name: "Glyph of Pyroblast",
			description: "Increases the critical strike chance of your Pyroblast spell by $56384s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_fireball02.jpg",
		},
		[MagePrimeGlyph.GlyphOfIceLance]: {
			name: "Glyph of Ice Lance",
			description: "Increases the damage of your Ice Lance spell by $56377s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostblast.jpg",
		},
		[MagePrimeGlyph.GlyphOfMageArmor]: {
			name: "Glyph of Mage Armor",
			description: "Your Mage Armor regenerates $s1% more mana.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magearmor.jpg",
		},
		[MagePrimeGlyph.GlyphOfMoltenArmor]: {
			name: "Glyph of Molten Armor",
			description: "Your Molten Armor grants an additional $s3% spell critical strike chance.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_moltenarmor.jpg",
		},
		[MagePrimeGlyph.GlyphOfConeOfCold]: {
			name: "Glyph of Cone of Cold",
			description: "Increases the damage of your Cone of Cold spell by $56364s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_glacier.jpg",
		},
		[MagePrimeGlyph.GlyphOfFrostfire]: {
			name: "Glyph of Frostfire",
			description: "Increases the damage done by your Frostfire Bolt by $s1% and your Frostfire Bolt now deals $s3% additional damage over 12 sec, stacking up to 3 times, but no longer reduces the victim's movement speed.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_frostfirebolt.jpg",
		},
		[MagePrimeGlyph.GlyphOfArcaneBlast]: {
			name: "Glyph of Arcane Blast",
			description: "Increases the damage from your Arcane Blast buff by $62210s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blast.jpg",
		},
		[MagePrimeGlyph.GlyphOfDeepFreeze]: {
			name: "Glyph of Deep Freeze",
			description: "Your Deep Freeze deals $s1% additional damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_deepfreeze.jpg",
		},
		[MagePrimeGlyph.GlyphOfArcaneBarrage]: {
			name: "Glyph of Arcane Barrage",
			description: "Increases the damage of your Arcane Barrage spell by $63092s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_arcanebarrage.jpg",
		},
		[MagePrimeGlyph.GlyphOfLivingBomb]: {
			name: "Glyph of Living Bomb",
			description: "Increases the damage of your Living Bomb spell by $s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_livingbomb.jpg",
		},
	},
	majorGlyphs: {
		[MageMajorGlyph.GlyphOfArcanePower]: {
			name: "Glyph of Arcane Power",
			description: "While Arcane Power is active the global cooldown of your Blink, Mana Shield, and Mirror Image is reduced to zero.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg",
		},
		[MageMajorGlyph.GlyphOfBlink]: {
			name: "Glyph of Blink",
			description: "Increases the distance you travel with the Blink spell by $56365s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg",
		},
		[MageMajorGlyph.GlyphOfEvocation]: {
			name: "Glyph of Evocation",
			description: "Your Evocation ability also causes you to regain ${$56380m1*4}% of your health over its duration.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_purge.jpg",
		},
		[MageMajorGlyph.GlyphOfFrostNova]: {
			name: "Glyph of Frost Nova",
			description: "Your Frost Nova targets can take an additional $56376s1% damage before the Frost Nova effect automatically breaks.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg",
		},
		[MageMajorGlyph.GlyphOfIceBlock]: {
			name: "Glyph of Ice Block",
			description: "Your Frost Nova cooldown is now reset every time you use Ice Block.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frost.jpg",
		},
		[MageMajorGlyph.GlyphOfIcyVeins]: {
			name: "Glyph of Icy Veins",
			description: "Your Icy Veins ability also removes all movement slowing and cast time slowing effects.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_coldhearted.jpg",
		},
		[MageMajorGlyph.GlyphOfInvisibility]: {
			name: "Glyph of Invisibility",
			description: "Increases your movement speed while Invisible by $87833s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_invisibility.jpg",
		},
		[MageMajorGlyph.GlyphOfPolymorph]: {
			name: "Glyph of Polymorph",
			description: "Your Polymorph spell also removes all damage over time effects from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_polymorph.jpg",
		},
		[MageMajorGlyph.GlyphOfDragonsBreath]: {
			name: "Glyph of Dragon's Breath",
			description: "Reduces the cooldown of your Dragon's Breath by ${$56373m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_head_dragon_01.jpg",
		},
		[MageMajorGlyph.GlyphOfBlastWave]: {
			name: "Glyph of Blast Wave",
			description: "Increases the duration of Blast Wave's slowing effect by ${$m1/1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg",
		},
		[MageMajorGlyph.GlyphOfSlow]: {
			name: "Glyph of Slow",
			description: "Increases the range of your Slow spell by $63091s1 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg",
		},
		[MageMajorGlyph.GlyphOfIceBarrier]: {
			name: "Glyph of Ice Barrier",
			description: "Increases the amount of damage absorbed by your Ice Barrier by $63095s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_ice_lament.jpg",
		},
		[MageMajorGlyph.GlyphOfManaShield]: {
			name: "Glyph of Mana Shield",
			description: "Reduces the cooldown of your Mana Shield by ${$70937m1/-1000} sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_shadow_detectlesserinvisibility.jpg",
		},
		[MageMajorGlyph.GlyphOfFrostArmor]: {
			name: "Glyph of Frost Armor",
			description: "Your Frost Armor also causes you to regenerate $s1% of your maximum mana every 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostarmor02.jpg",
		},
	},
	minorGlyphs: {
		[MageMinorGlyph.GlyphOfArcaneBrilliance]: {
			name: "Glyph of Arcane Brilliance",
			description: "Reduces the mana cost of your Arcane Brilliance spell by $57924s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_holy_magicalsentry.jpg",
		},
		[MageMinorGlyph.GlyphOfConjuring]: {
			name: "Glyph of Conjuring",
			description: "Reduces the mana cost of your Conjuring spells by $57928s1%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_conjurefoodrank9.jpg",
		},
		[MageMinorGlyph.GlyphOfTheMonkey]: {
			name: "Glyph of the Monkey",
			description: "Your Polymorph: Sheep spell polymorphs the target into a monkey instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_aspectofthemonkey.jpg",
		},
		[MageMinorGlyph.GlyphOfThePenguin]: {
			name: "Glyph of the Penguin",
			description: "Your Polymorph: Sheep spell polymorphs the target into a penguin instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_penguinpet.jpg",
		},
		[MageMinorGlyph.ZzoldglyphOfTheBearCub]: {
			name: "zzOLDGlyph of the Bear Cub",
			description: "Your Polymorph: Sheep spell polymorphs the target into a polar bear cub instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_pet_babyblizzardbear.jpg",
		},
		[MageMinorGlyph.GlyphOfSlowFall]: {
			name: "Glyph of Slow Fall",
			description: "Your Slow Fall spell no longer requires a reagent.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_featherfall.jpg",
		},
		[MageMinorGlyph.GlyphOfMirrorImage]: {
			name: "Glyph of Mirror Image",
			description: "Your Mirror Images cast Arcane Blast or Fireball instead of Frostbolt depending on your primary talent tree.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
		},
		[MageMinorGlyph.GlyphOfArmors]: {
			name: "Glyph of Armors",
			description: "Increases the duration of your Armor spells by ${$89749m1/60000} min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_chillingarmor.jpg",
		},
	},
};
