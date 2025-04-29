import { MageMajorGlyph, MageMinorGlyph, MageTalents } from '../proto/mage.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MageTalentJson from './trees/mage.json';

export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

export const mageGlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		[MageMajorGlyph.GlyphOfArcanePower]: {
			name: 'Glyph of Arcane Power',
			description: 'While Arcane Power is active the global cooldown of your Blink, Mana Shield, and Mirror Image is reduced to zero.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_lightning.jpg',
		},
		[MageMajorGlyph.GlyphOfBlastWave]: {
			name: 'Glyph of Blast Wave',
			description: "Increases the duration of Blast Wave's slowing effect by 1 sec.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_excorcism_02.jpg',
		},
		[MageMajorGlyph.GlyphOfBlink]: {
			name: 'Glyph of Blink',
			description: 'Increases the distance you travel with the Blink spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_arcane_blink.jpg',
		},
		[MageMajorGlyph.GlyphOfDragonSBreath]: {
			name: "Glyph of Dragon's Breath",
			description: "Reduces the cooldown of your Dragon's Breath by 3 sec.",
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_head_dragon_01.jpg',
		},
		[MageMajorGlyph.GlyphOfEvocation]: {
			name: 'Glyph of Evocation',
			description: 'Your Evocation ability also causes you to regain 40% of your health over its duration.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_purge.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostArmor]: {
			name: 'Glyph of Frost Armor',
			description: 'Your Frost Armor also causes you to regenerate 2% of your maximum mana every 5 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostarmor02.jpg',
		},
		[MageMajorGlyph.GlyphOfFrostNova]: {
			name: 'Glyph of Frost Nova',
			description: 'Your Frost Nova targets can take an additional 20% damage before the Frost Nova effect automatically breaks.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frostnova.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBarrier]: {
			name: 'Glyph of Ice Barrier',
			description: 'Increases the amount of damage absorbed by your Ice Barrier by 30%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_ice_lament.jpg',
		},
		[MageMajorGlyph.GlyphOfIceBlock]: {
			name: 'Glyph of Ice Block',
			description: 'Your Frost Nova cooldown is now reset every time you use Ice Block.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_frost.jpg',
		},
		[MageMajorGlyph.GlyphOfIcyVeins]: {
			name: 'Glyph of Icy Veins',
			description: 'Your Icy Veins ability also removes all movement slowing and cast time slowing effects.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_coldhearted.jpg',
		},
		[MageMajorGlyph.GlyphOfInvisibility]: {
			name: 'Glyph of Invisibility',
			description: 'Increases your movement speed while Invisible by 40%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_invisibility.jpg',
		},
		[MageMajorGlyph.GlyphOfManaShield]: {
			name: 'Glyph of Mana Shield',
			description: 'Reduces the cooldown of your Mana Shield by 2 sec.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_detectlesserinvisibility.jpg',
		},
		[MageMajorGlyph.GlyphOfPolymorph]: {
			name: 'Glyph of Polymorph',
			description: 'Your Polymorph spell also removes all damage over time effects from the target.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_polymorph.jpg',
		},
		[MageMajorGlyph.GlyphOfSlow]: {
			name: 'Glyph of Slow',
			description: 'Increases the range of your Slow spell by 5 yards.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_slow.jpg',
		},
	},
	minorGlyphs: {
		[MageMinorGlyph.GlyphOfArcaneBrilliance]: {
			name: 'Glyph of Arcane Brilliance',
			description: 'Reduces the mana cost of your Arcane Brilliance spell by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_magicalsentry.jpg',
		},
		[MageMinorGlyph.GlyphOfArmors]: {
			name: 'Glyph of Armors',
			description: 'Increases the duration of your Armor spells by 30 min.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_frost_chillingarmor.jpg',
		},
		[MageMinorGlyph.GlyphOfConjuring]: {
			name: 'Glyph of Conjuring',
			description: 'Reduces the mana cost of your Conjuring spells by 50%.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_mage_conjurefoodrank9.jpg',
		},
		[MageMinorGlyph.GlyphOfMirrorImage]: {
			name: 'Glyph of Mirror Image',
			description: 'Your Mirror Images cast Arcane Blast or Fireball instead of Frostbolt depending on your primary talent tree.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_lesserinvisibilty.jpg',
		},
		[MageMinorGlyph.GlyphOfSlowFall]: {
			name: 'Glyph of Slow Fall',
			description: 'Your Slow Fall spell no longer requires a reagent.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/spell_magic_featherfall.jpg',
		},
		[MageMinorGlyph.GlyphOfTheMonkey]: {
			name: 'Glyph of the Monkey',
			description: 'Your Polymorph: Sheep spell polymorphs the target into a monkey instead.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/ability_hunter_aspectofthemonkey.jpg',
		},
		[MageMinorGlyph.GlyphOfThePenguin]: {
			name: 'Glyph of the Penguin',
			description: 'Your Polymorph: Sheep spell polymorphs the target into a penguin instead.',
			iconUrl: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_penguinpet.jpg',
		},
	},
};
