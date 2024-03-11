import { LOCAL_STORAGE_PREFIX } from '../constants/other';
import { PlayerSpec } from '../player_spec';
import { Spec } from '../proto/common';
import * as DeathKnightSpecs from './death_knight';
import * as DruidSpecs from './druid';
import * as HunterSpecs from './hunter';
import * as MageSpecs from './mage';
import * as PaladinSpecs from './paladin';
import * as PriestSpecs from './priest';
import * as RogueSpecs from './rogue';
import * as ShamanSpecs from './shaman';
import * as WarlockSpecs from './warlock';
import * as WarriorSpecs from './warrior';

const protoToPlayerSpec: Record<Spec, PlayerSpec<Spec> | undefined> = {
	[Spec.SpecUnknown]: undefined,
	// Death Knight
	[Spec.SpecBloodDeathKnight]: DeathKnightSpecs.BloodDeathKnight,
	[Spec.SpecFrostDeathKnight]: DeathKnightSpecs.FrostDeathKnight,
	[Spec.SpecUnholyDeathKnight]: DeathKnightSpecs.UnholyDeathKnight,
	// Druid
	[Spec.SpecBalanceDruid]: DruidSpecs.BalanceDruid,
	[Spec.SpecFeralDruid]: DruidSpecs.FeralDruid,
	[Spec.SpecRestorationDruid]: DruidSpecs.RestorationDruid,
	// Hunter
	[Spec.SpecBeastMasteryHunter]: HunterSpecs.BeastMasteryHunter,
	[Spec.SpecMarksmanshipHunter]: HunterSpecs.MarksmanshipHunter,
	[Spec.SpecSurvivalHunter]: HunterSpecs.SurvivalHunter,
	// Mage
	[Spec.SpecArcaneMage]: MageSpecs.ArcaneMage,
	[Spec.SpecFireMage]: MageSpecs.FireMage,
	[Spec.SpecFrostMage]: MageSpecs.FrostMage,
	// Paladin
	[Spec.SpecHolyPaladin]: PaladinSpecs.HolyPaladin,
	[Spec.SpecProtectionPaladin]: PaladinSpecs.ProtectionPaladin,
	[Spec.SpecRetributionPaladin]: PaladinSpecs.RetributionPaladin,
	// Priest
	[Spec.SpecDisciplinePriest]: PriestSpecs.DisciplinePriest,
	[Spec.SpecHolyPriest]: PriestSpecs.HolyPriest,
	[Spec.SpecShadowPriest]: PriestSpecs.ShadowPriest,
	// Rogue
	[Spec.SpecAssassinationRogue]: RogueSpecs.AssassinationRogue,
	[Spec.SpecCombatRogue]: RogueSpecs.CombatRogue,
	[Spec.SpecSubtletyRogue]: RogueSpecs.SubtletyRogue,
	// Shaman
	[Spec.SpecElementalShaman]: ShamanSpecs.ElementalShaman,
	[Spec.SpecEnhancementShaman]: ShamanSpecs.EnhancementShaman,
	[Spec.SpecRestorationShaman]: ShamanSpecs.RestorationShaman,
	// Warlock
	[Spec.SpecAfflictionWarlock]: WarlockSpecs.AfflictionWarlock,
	[Spec.SpecDemonologyWarlock]: WarlockSpecs.DemonologyWarlock,
	[Spec.SpecDestructionWarlock]: WarlockSpecs.DestructionWarlock,
	// Warrior
	[Spec.SpecArmsWarrior]: WarriorSpecs.ArmsWarrior,
	[Spec.SpecFuryWarrior]: WarriorSpecs.FuryWarrior,
	[Spec.SpecProtectionWarrior]: WarriorSpecs.ProtectionWarrior,
};

export const PlayerSpecs = {
	...DeathKnightSpecs,
	...DruidSpecs,
	...HunterSpecs,
	...MageSpecs,
	...PaladinSpecs,
	...PriestSpecs,
	...RogueSpecs,
	...ShamanSpecs,
	...WarlockSpecs,
	...WarriorSpecs,
	// Prefixes used for storing browser data for each site. Even if a Spec is
	// renamed, DO NOT change these values or people will lose their saved data.
	getLocalStorageKey: <SpecType extends Spec>(spec: PlayerSpec<SpecType>): string => {
		return `${LOCAL_STORAGE_PREFIX}_${spec.friendlyName.toLowerCase().replace(/\s/, '_')}_${spec.playerClass.friendlyName
			.toLowerCase()
			.replace(/\s/, '_')}`;
	},
	fromProto: <SpecType extends Spec>(protoId: SpecType): PlayerSpec<SpecType> => {
		if (protoId == Spec.SpecUnknown) {
			throw new Error('Invalid Spec');
		}

		return protoToPlayerSpec[protoId] as PlayerSpec<SpecType>;
	},
};
