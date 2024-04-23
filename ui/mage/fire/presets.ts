import * as PresetUtils from '../../core/preset_utils';
import { Conjured, Consumes, Faction, Flask, Food, Glyphs, Potions, Profession, Spec } from '../../core/proto/common';
import {
	FireMage_Options as MageOptions,
	FireMage_Rotation as MageRotation,
	FireMage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
	MageMajorGlyph,
	MageMinorGlyph,
	MageOptions_ArmorType as ArmorType,
} from '../../core/proto/mage';
import { SavedTalents } from '../../core/proto/ui';
import FireApl from './apls/fire.apl.json';
import FireAoeApl from './apls/fire_aoe.apl.json';
import P1FireGear from './gear_sets/p1_fire.gear.json';
import P2FireGear from './gear_sets/p2_fire.gear.json';
import P3FireAllianceGear from './gear_sets/p3_fire_alliance.gear.json';
import P3FireHordeGear from './gear_sets/p3_fire_horde.gear.json';
import P4FireAllianceGear from './gear_sets/p4_fire_alliance.gear.json';
import P4FireHordeGear from './gear_sets/p4_fire_horde.gear.json';
import PreraidFireGear from './gear_sets/preraid_fire.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const FIRE_PRERAID_PRESET = PresetUtils.makePresetGear('Fire Preraid Preset', PreraidFireGear, { talentTree: 1 });
export const FIRE_P1_PRESET = PresetUtils.makePresetGear('Fire P1 Preset', P1FireGear, { talentTree: 1 });
export const FIRE_P2_PRESET = PresetUtils.makePresetGear('Fire P2 Preset', P2FireGear, {
	talentTree: 1,
});
export const FIRE_P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('Fire P3 Preset [A]', P3FireAllianceGear, {
	talentTree: 1,
	faction: Faction.Alliance,
});
export const FIRE_P3_PRESET_HORDE = PresetUtils.makePresetGear('Fire P3 Preset [H]', P3FireHordeGear, {
	talentTree: 1,
	faction: Faction.Horde,
});
export const FIRE_P4_PRESET_ALLIANCE = PresetUtils.makePresetGear('Fire P4 Preset [A]', P4FireAllianceGear, {
	talentTree: 1,
	faction: Faction.Alliance,
});
export const FIRE_P4_PRESET_HORDE = PresetUtils.makePresetGear('Fire P4 Preset [H]', P4FireHordeGear, {
	talentTree: 1,
	faction: Faction.Horde,
});

export const DefaultSimpleRotation = MageRotation.create({
	primaryFireSpell: PrimaryFireSpell.Fireball,
	maintainImprovedScorch: false,
});

export const ROTATION_PRESET_SIMPLE = PresetUtils.makePresetSimpleRotation('Simple Default', Spec.SpecFireMage, DefaultSimpleRotation);
export const FIRE_ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Fire', FireApl, { talentTree: 1 });
export const FIRE_ROTATION_PRESET_AOE = PresetUtils.makePresetAPLRotation('Fire AOE', FireAoeApl, { talentTree: 1 });

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.
export const FireTalents = {
	name: 'Fire',
	data: SavedTalents.create({
		// talentsString: '23000503110003-0055030012303331053120301351',
		// glyphs: Glyphs.create({
		// 	major1: MageMajorGlyph.GlyphOfFireball,
		// 	major2: MageMajorGlyph.GlyphOfMoltenArmor,
		// 	major3: MageMajorGlyph.GlyphOfLivingBomb,
		// 	minor1: MageMinorGlyph.GlyphOfSlowFall,
		// 	minor2: MageMinorGlyph.GlyphOfFrostWard,
		// 	minor3: MageMinorGlyph.GlyphOfBlastWave,
		// }),
	}),
};

export const DefaultFireOptions = MageOptions.create({
	classOptions: {
		armor: ArmorType.MoltenArmor,
	},
});

export const DefaultFireConsumes = Consumes.create({
	flask: Flask.FlaskOfTheDraconicMind,
	food: Food.FoodSeafoodFeast,
	defaultPotion: Potions.VolcanicPotion,
	defaultConjured: Conjured.ConjuredFlameCap,
	prepopPotion: Potions.VolcanicPotion,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	ebonPlaguebringer: true,
	mangle: true,
	criticalMass: true,
	demoralizingShout: true,
	frostFever: true,
	judgement: true,
});
	
export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
};
