import * as PresetUtils from '../../core/preset_utils';
import {
	Consumes,
	Debuffs,
	Faction,
	Flask,
	Food,
	Glyphs,
	IndividualBuffs,
	PetFood,
	Potions,
	Profession,
	RaidBuffs,
	TristateEffect,
} from '../../core/proto/common';
import { SavedTalents } from '../../core/proto/ui';
import {
	DestructionWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WeaponImbue,
} from '../../core/proto/warlock';
import DestroApl from './apls/destro.apl.json';
import P1DemoDestroGear from './gear_sets/p1_demodestro.gear.json';
import P2DemoDestroGear from './gear_sets/p2_demodestro.gear.json';
import P3DestroAllianceGear from './gear_sets/p3_destro_alliance.gear.json';
import P3DestroHordeGear from './gear_sets/p3_destro_horde.gear.json';
import P4DestroGear from './gear_sets/p4_destro.gear.json';
import PreraidDemoDestroGear from './gear_sets/preraid_demodestro.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const PRERAID_DEMODESTRO_PRESET = PresetUtils.makePresetGear('Preraid Demo/Destro', PreraidDemoDestroGear, {
	tooltip: BIS_TOOLTIP,
});
export const P1_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P1 Demo/Destro', P1DemoDestroGear, { tooltip: BIS_TOOLTIP });
export const P2_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P2 Demo/Destro', P2DemoDestroGear, { tooltip: BIS_TOOLTIP });
export const P3_DESTRO_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3 Destro [A]', P3DestroAllianceGear, {
	tooltip: BIS_TOOLTIP,
	faction: Faction.Alliance,
});
export const P3_DESTRO_HORDE_PRESET = PresetUtils.makePresetGear('P3 Destro [H]', P3DestroHordeGear, {
	tooltip: BIS_TOOLTIP,
	faction: Faction.Horde,
});
export const P4_DESTRO_PRESET = PresetUtils.makePresetGear('P4 Destro', P4DestroGear, { tooltip: BIS_TOOLTIP });

export const APL_Destro_Default = PresetUtils.makePresetAPLRotation('Destro', DestroApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const DestructionTalents = {
	name: 'Destruction',
	data: SavedTalents.create({
		// talentsString: '-03310030003-05203205210331051335230351',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfConflagrate,
		// 	major2: MajorGlyph.GlyphOfLifeTap,
		// 	major3: MajorGlyph.GlyphOfIncinerate,
		// 	minor1: MinorGlyph.GlyphOfSouls,
		// 	minor2: MinorGlyph.GlyphOfDrainSoul,
		// 	minor3: MinorGlyph.GlyphOfSubjugateDemon,
		// }),
	}),
};

export const DestructionOptions = WarlockOptions.create({
	classOptions: {
		armor: Armor.FelArmor,
		summon: Summon.Imp,
		weaponImbue: WeaponImbue.GrandFirestone,
		detonateSeed: true,
	},
});

export const DefaultConsumes = Consumes.create({
	flask: Flask.FlaskOfTheFrostWyrm,
	food: Food.FoodFishFeast,
	petFood: PetFood.PetFoodSpicedMammothTreats,
	defaultPotion: Potions.VolcanicPotion,
	prepopPotion: Potions.VolcanicPotion,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	arcaneBrilliance: true,
	bloodlust: true,
	markOfTheWild: true,
	icyTalons: true,
	moonkinForm: true,
	leaderOfThePack: true,
	powerWordFortitude: true,
	strengthOfEarthTotem: true,
	trueshotAura: true,
	wrathOfAirTotem: true,
	demonicPact: true,
	blessingOfKings: true,
	blessingOfMight: true,
	communion: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	vampiricTouch: true,
});

export const DestroIndividualBuffs = IndividualBuffs.create({
	vampiricTouch: true,
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

export const DestroDebuffs = Debuffs.create({
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
	distanceFromTarget: 25,
	profession1: Profession.Engineering,
	profession2: Profession.Tailoring,
	channelClipDelay: 150,
};
