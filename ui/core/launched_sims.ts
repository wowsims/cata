import { Phase } from './constants/other';
import { Player } from './player';
import { Spec } from './proto/common';

// This file is for anything related to launching a new sim. DO NOT touch this
// file until your sim is ready to launch!

export enum LaunchStatus {
	Unlaunched,
	Alpha,
	Beta,
	Launched,
}

export type SimStatus = {
	phase: Phase;
	status: LaunchStatus;
};

export const raidSimStatus: SimStatus = {
	phase: Phase.Phase1,
	status: LaunchStatus.Unlaunched,
};

// This list controls which links are shown in the top-left dropdown menu.
export const simLaunchStatuses: Record<Spec, SimStatus> = {
	[Spec.SpecUnknown]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	// Death Knight
	[Spec.SpecBloodDeathKnight]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFrostDeathKnight]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecUnholyDeathKnight]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Druid
	[Spec.SpecBalanceDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFeralDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecGuardianDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecRestorationDruid]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	// Hunter
	[Spec.SpecBeastMasteryHunter]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecMarksmanshipHunter]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecSurvivalHunter]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Mage
	[Spec.SpecArcaneMage]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFireMage]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFrostMage]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	// Monk
	[Spec.SpecBrewmasterMonk]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecMistweaverMonk]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecWindwalkerMonk]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Paladin
	[Spec.SpecHolyPaladin]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecProtectionPaladin]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecRetributionPaladin]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Priest
	[Spec.SpecDisciplinePriest]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecHolyPriest]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	[Spec.SpecShadowPriest]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Rogue
	[Spec.SpecAssassinationRogue]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecCombatRogue]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	[Spec.SpecSubtletyRogue]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Launched,
	},
	// Shaman
	[Spec.SpecElementalShaman]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecEnhancementShaman]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecRestorationShaman]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Unlaunched,
	},
	// Warlock
	[Spec.SpecAfflictionWarlock]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecDemonologyWarlock]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecDestructionWarlock]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	// Warrior
	[Spec.SpecArmsWarrior]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecFuryWarrior]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
	[Spec.SpecProtectionWarrior]: {
		phase: Phase.Phase1,
		status: LaunchStatus.Alpha,
	},
};

export const getSpecLaunchStatus = (player: Player<any>) => simLaunchStatuses[player.getSpec() as Spec].status;
