import { MageMajorGlyph, MageMinorGlyph, MageTalents } from '../proto/mage.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MageTalentJson from './trees/mage.json';export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

export const mageGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[MageMajorGlyph.GlyphOfArcaneExplosion]: {
			name: "Glyph of Arcane Explosion",
			description: "Increases the radius of your Arcane Explosion by 5 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_wispsplode.jpg",
		},
		[MageMajorGlyph.GlyphOfBlink]: {
			name: "Glyph of Blink",
			description: "Increases the distance you travel with the Blink spell by 8 yards.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg",
		},
		[MageMajorGlyph.GlyphOfEvocation]: {
			name: "Glyph of Evocation",
			description: "Your Evocation ability also causes you to regain 60% of your health over its duration.\u000D\u000A\u000D\u000A With the Invocation talent, you instead gain 10% of your health upon completing an Evocation.\u000D\u000A\u000D\u000A With the Rune of Power talent, you gain 1% of your health per second while standing in your own Rune of Power.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_purge.jpg",
		},
		[MageMajorGlyph.GlyphOfCombustion]: {
			name: "Glyph of Combustion",
			description: "Increases the direct damage, the duration of the damage over time effect and the cooldown of Combustion by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_fire_sealoffire.jpg",
		},
		[MageMajorGlyph.GlyphOfFrostNova]: {
			name: "Glyph of Frost Nova",
			description: "Reduces the cooldown of Frost Nova by 5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg",
		},
		[MageMajorGlyph.GlyphOfIceBlock]: {
			name: "Glyph of Ice Block",
			description: "When Ice Block terminates, it triggers an instant free Frost Nova and makes you immune to all spells for 3s.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frost.jpg",
		},
		[MageMajorGlyph.GlyphOfSplittingIce]: {
			name: "Glyph of Splitting Ice",
			description: "Your Ice Lance and Icicles now hit 1 additional target for 50% damage.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostblast.jpg",
		},
		[MageMajorGlyph.GlyphOfConeOfCold]: {
			name: "Glyph of Cone of Cold",
			description: "Increases the damage done by Cone of Cold by 200%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_glacier.jpg",
		},
		[MageMajorGlyph.GlyphOfRapidDisplacement]: {
			name: "Glyph of Rapid Displacement",
			description: "Blink now has 2 charges, gaining a charge every 15 sec, but no longer frees the caster from stuns and bonds.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg",
		},
		[MageMajorGlyph.GlyphOfManaGem]: {
			name: "Glyph of Mana Gem",
			description: "Your Conjure Mana Gem spell now creates a Brilliant Mana Gem, which holds up to 10 charges.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_jewelcrafting_gem_05.jpg",
		},
		[MageMajorGlyph.GlyphOfPolymorph]: {
			name: "Glyph of Polymorph",
			description: "Your Polymorph spell also removes all damage over time effects from the target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_polymorph.jpg",
		},
		[MageMajorGlyph.GlyphOfIcyVeins]: {
			name: "Glyph of Icy Veins",
			description: "Your Icy Veins causes your Frostbolt, Frostfire Bolt, Ice Lance, and your Water Elemental\'s Waterbolt spells to split into 3 smaller bolts that each do 40% damage, instead of increasing spell casting speed.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_coldhearted.jpg",
		},
		[MageMajorGlyph.GlyphOfSpellsteal]: {
			name: "Glyph of Spellsteal",
			description: "Spellsteal now also heals you for 5% of your maximum health when it successfully steals a spell.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_arcane02.jpg",
		},
		[MageMajorGlyph.GlyphOfFrostfireBolt]: {
			name: "Glyph of Frostfire Bolt",
			description: "Reduces the cast time of Frostfire Bolt by 0.5 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_frostfirebolt.jpg",
		},
		[MageMajorGlyph.GlyphOfRemoveCurse]: {
			name: "Glyph of Remove Curse",
			description: "Increases the damage you deal by 5% for 0ms after you successfully remove a curse.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_removecurse.jpg",
		},
		[MageMajorGlyph.GlyphOfArcanePower]: {
			name: "Glyph of Arcane Power",
			description: "Increases the duration and cooldown of Arcane Power by 100%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg",
		},
		[MageMajorGlyph.GlyphOfWaterElemental]: {
			name: "Glyph of Water Elemental",
			description: "Increases the health of your Water Elemental by 40%, and allows it to cast while moving. When in Assist mode and in combat, commanding your Water Elemental to Follow will cause it to stay near you and autocast Waterbolt when your target is in range.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_summonwaterelemental_2.jpg",
		},
		[MageMajorGlyph.GlyphOfSlow]: {
			name: "Glyph of Slow",
			description: "Your Arcane Blast spell applies the Slow spell to any target it damages if no target is currently affected by your Slow.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg",
		},
		[MageMajorGlyph.GlyphOfDeepFreeze]: {
			name: "Glyph of Deep Freeze",
			description: "Your Deep Freeze spell is no longer on the global cooldown, but its duration is reduced by 1 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_mage_deepfreeze.jpg",
		},
		[MageMajorGlyph.GlyphOfCounterspell]: {
			name: "Glyph of Counterspell",
			description: "Your Counterspell can now be cast while casting or channeling other spells, but its cooldown is increased by 4 sec.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_iceshock.jpg",
		},
		[MageMajorGlyph.GlyphOfInfernoBlast]: {
			name: "Glyph of Inferno Blast",
			description: "Your Inferno Blast spell spreads Pyroblast, Ignite, and Combustion to 1 additional target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_mage_infernoblast.jpg",
		},
		[MageMajorGlyph.GlyphOfArmors]: {
			name: "Glyph of Armors",
			description: "Reduces the cast time of your Frost Armor, Mage Armor, and Molten Armor spells by 1.5 sec, and increases the defensive effect of each Armor by an additional 10%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostarmor02.jpg",
		},
	},
	minorGlyphs: {
		[MageMinorGlyph.GlyphOfLooseMana]: {
			name: "Glyph of Loose Mana",
			description: "Your Mana Gem now restores mana over 6 sec, rather than instantly.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_gem_sapphire_02.jpg",
		},
		[MageMinorGlyph.GlyphOfMomentum]: {
			name: "Glyph of Momentum",
			description: "Your Blink spell teleports you in the direction you are moving instead of the direction you are facing.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg",
		},
		[MageMinorGlyph.GlyphOfCrittermorph]: {
			name: "Glyph of Crittermorph",
			description: "When cast on critters, your Polymorph spells now last 1440min and can be cast on multiple targets.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_doublepolymorph2.jpg",
		},
		[MageMinorGlyph.GlyphOfThePorcupine]: {
			name: "Glyph of the Porcupine",
			description: "Your Polymorph spell polymorphs the target into a porcupine instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_polymorphpig.jpg",
		},
		[MageMinorGlyph.GlyphOfConjureFamiliar]: {
			name: "Glyph of Conjure Familiar",
			description: "Teaches you the ability Conjure Familiar.\u000D\u000A\u000D\u000A Conjures a familiar stone, containing either an Arcane, Fiery, or Icy Familiar.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_nature_elementalabsorption.jpg",
		},
		[MageMinorGlyph.GlyphOfTheMonkey]: {
			name: "Glyph of the Monkey",
			description: "Your Polymorph spell polymorphs the target into a monkey instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/ability_hunter_aspectofthemonkey.jpg",
		},
		[MageMinorGlyph.GlyphOfThePenguin]: {
			name: "Glyph of the Penguin",
			description: "Your Polymorph spell polymorphs the target into a penguin instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_penguinpet.jpg",
		},
		[MageMinorGlyph.GlyphOfTheBearCub]: {
			name: "Glyph of the Bear Cub",
			description: "Your Polymorph spell polymorphs the target into a polar bear cub instead.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_pet_babyblizzardbear.jpg",
		},
		[MageMinorGlyph.GlyphOfArcaneLanguage]: {
			name: "Glyph of Arcane Language",
			description: "Your Arcane Brilliance spell allows you to comprehend your allies\' racial languages.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_fish_68.jpg",
		},
		[MageMinorGlyph.GlyphOfIllusion]: {
			name: "Glyph of Illusion",
			description: "Teaches you the ability Illusion.\u000D\u000A\u000D\u000A Transforms the Mage to look like someone else for 0ms.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_mask_01.jpg",
		},
		[MageMinorGlyph.GlyphOfMirrorImage]: {
			name: "Glyph of Mirror Image",
			description: "Your Mirror Images cast Arcane Blast or Fireball instead of Frostbolt depending on your primary talent tree.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg",
		},
		[MageMinorGlyph.GlyphOfRapidTeleportation]: {
			name: "Glyph of Rapid Teleportation",
			description: "After casting a Mage Teleport spell, or entering a Mage Portal, your movement speed is increased by 70% for 1min.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_arcane_portaldalaran.jpg",
		},
		[MageMinorGlyph.GlyphOfDiscreetMagic]: {
			name: "Glyph of Discreet Magic",
			description: "Your Nether Tempest, Living Bomb, Frost Bomb, Arcane Barrage, and Inferno Blast no longer affect targets more than 5 yds away from their primary target.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_enchant_essencemagicsmall.jpg",
		},
		[MageMinorGlyph.GlyphOfTheUnboundElemental]: {
			name: "Glyph of the Unbound Elemental",
			description: "Your Water Elemental is replaced by an Unbound Water Elemental.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/spell_frost_summonwaterelemental.jpg",
		},
		[MageMinorGlyph.GlyphOfEvaporation]: {
			name: "Glyph of Evaporation",
			description: "Reduces the size of your Water Elemental by 40%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_elemental_mote_water01.jpg",
		},
		[MageMinorGlyph.GlyphOfCondensation]: {
			name: "Glyph of Condensation",
			description: "Increases the size of your Water Elemental by 40%.",
			iconUrl: "https://wow.zamimg.com/images/wow/icons/large/inv_misc_volatilewater.jpg",
		},
	},
};
