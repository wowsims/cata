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
	DemonologyWarlock_Options as WarlockOptions,
	WarlockMajorGlyph as MajorGlyph,
	WarlockMinorGlyph as MinorGlyph,
	WarlockOptions_Armor as Armor,
	WarlockOptions_Summon as Summon,
	WarlockOptions_WeaponImbue as WeaponImbue,
} from '../../core/proto/warlock';
import DemoApl from './apls/demo.apl.json';
import P1DemoDestroGear from './gear_sets/p1_demodestro.gear.json';
import P2DemoDestroGear from './gear_sets/p2_demodestro.gear.json';
import P3DemoAllianceGear from './gear_sets/p3_demo_alliance.gear.json';
import P3DemoHordeGear from './gear_sets/p3_demo_horde.gear.json';
import P4DemoGear from './gear_sets/p4_demo.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const BIS_TOOLTIP = "This gear preset is inspired from Zephan's Affliction guide: https://www.warcrafttavern.com/wotlk/guides/pve-affliction-warlock/";

export const P1_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P1 Demo/Destro', P1DemoDestroGear, { tooltip: BIS_TOOLTIP });
export const P2_DEMODESTRO_PRESET = PresetUtils.makePresetGear('P2 Demo/Destro', P2DemoDestroGear, { tooltip: BIS_TOOLTIP });
export const P3_DEMO_ALLIANCE_PRESET = PresetUtils.makePresetGear('P3 Demo [A]', P3DemoAllianceGear, {
	tooltip: BIS_TOOLTIP,
	faction: Faction.Alliance,
});
export const P3_DEMO_HORDE_PRESET = PresetUtils.makePresetGear('P3 Demo [H]', P3DemoHordeGear, { tooltip: BIS_TOOLTIP, faction: Faction.Horde });
export const P4_DEMO_PRESET = PresetUtils.makePresetGear('P4 Demo', P4DemoGear, { tooltip: BIS_TOOLTIP });

export const APL_Demo_Default = PresetUtils.makePresetAPLRotation('Demo', DemoApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/wotlk/talent-calc and copy the numbers in the url.

export const DemonologyTalents = {
	name: 'Demonology',
	data: SavedTalents.create({
		// talentsString: '-203203301035012530135201351-550000052',
		// glyphs: Glyphs.create({
		// 	major1: MajorGlyph.GlyphOfLifeTap,
		// 	major2: MajorGlyph.GlyphOfQuickDecay,
		// 	major3: MajorGlyph.GlyphOfFelguard,
		// 	minor1: MinorGlyph.GlyphOfSouls,
		// 	minor2: MinorGlyph.GlyphOfDrainSoul,
		// 	minor3: MinorGlyph.GlyphOfSubjugateDemon,
		// }),
	}),
};

export const DemonologyOptions = WarlockOptions.create({
	classOptions: {
		armor: Armor.FelArmor,
		summon: Summon.Felguard,
		weaponImbue: WeaponImbue.GrandSpellstone,
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
