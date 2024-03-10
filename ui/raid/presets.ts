import { naturalSpecOrder } from 'ui/core/spec.js';

import { IndividualSimUI, IndividualSimUIConfig, RaidSimPreset } from '../core/individual_sim_ui.js';
import { getSpecConfig,Player } from '../core/player.js';
import {
	Spec
} from '../core/proto/common.js';
import { TankDeathknightSimUI } from '../death_knight/blood/sim.js';
import { DeathknightSimUI } from '../death_knight/sim.js';
import { FeralTankDruidSimUI } from '../druid/_feral_tank/sim.js';
import { BalanceDruidSimUI } from '../druid/balance/sim.js';
import { FeralDruidSimUI } from '../druid/feral/sim.js';
import { RestorationDruidSimUI } from '../druid/restoration_druid/sim.js';
import { HunterSimUI } from '../hunter/sim.js';
import { MageSimUI } from '../mage/sim.js';
import { HolyPaladinSimUI } from '../paladin/holy/sim.js';
import { ProtectionPaladinSimUI } from '../paladin/protection/sim.js';
import { RetributionPaladinSimUI } from '../paladin/retribution/sim.js';
import { ShadowPriestSimUI } from '../priest/shadow/sim.js';
import { RogueSimUI } from '../rogue/sim.js';
import { ElementalShamanSimUI } from '../shaman/elemental/sim.js';
import { EnhancementShamanSimUI } from '../shaman/enhancement/sim.js';
import { RestorationShamanSimUI } from '../shaman/restoration/sim.js';
import { WarlockSimUI } from '../warlock/sim.js';
import { ProtectionWarriorSimUI } from '../warrior/protection/sim.js';
import { WarriorSimUI } from '../warrior/sim.js';

export const specSimFactories: Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>> = {
	[Spec.SpecTankDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new TankDeathknightSimUI(parentElem, player),
	[Spec.SpecDeathknight]: (parentElem: HTMLElement, player: Player<any>) => new DeathknightSimUI(parentElem, player),
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecFeralTankDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralTankDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	[Spec.SpecHunter]: (parentElem: HTMLElement, player: Player<any>) => new HunterSimUI(parentElem, player),
	[Spec.SpecMage]: (parentElem: HTMLElement, player: Player<any>) => new MageSimUI(parentElem, player),
	[Spec.SpecRogue]: (parentElem: HTMLElement, player: Player<any>) => new RogueSimUI(parentElem, player),
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	[Spec.SpecWarrior]: (parentElem: HTMLElement, player: Player<any>) => new WarriorSimUI(parentElem, player),
	[Spec.SpecProtectionWarrior]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionWarriorSimUI(parentElem, player),
	[Spec.SpecWarlock]: (parentElem: HTMLElement, player: Player<any>) => new WarlockSimUI(parentElem, player),
};

export const playerPresets: Array<RaidSimPreset<any>> = naturalSpecOrder
	.map(spec => getSpecConfig(spec.protoID))
	.map(config => {
		const indSimUiConfig = config as IndividualSimUIConfig<any>;
		return indSimUiConfig.raidSimPresets;
	})
	.flat();

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];
