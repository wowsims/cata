import { IndividualSimUI, IndividualSimUIConfig, RaidSimPreset } from '../core/individual_sim_ui.js';
import { getSpecConfig, Player } from '../core/player.js';
import { naturalPlayerClassOrder } from '../core/player_class.js';
import { Spec } from '../core/proto/common.js';
import { BloodDeathknightSimUI } from '../death_knight/blood/sim.js';
import { FrostDeathKnightSimUI } from '../death_knight/frost/sim';
import { UnholyDeathknightSimUI } from '../death_knight/unholy/sim';
import { BalanceDruidSimUI } from '../druid/balance/sim.js';
import { FeralDruidSimUI } from '../druid/feral/sim.js';
import { RestorationDruidSimUI } from '../druid/restoration/sim.js';
import { BeastMasteryHunterSimUI } from '../hunter/beast_mastery/sim';
import { MarksmanshipHunterSimUI } from '../hunter/marksmanship/sim';
import { SurvivalHunterSimUI } from '../hunter/survival/sim';
import { HolyPaladinSimUI } from '../paladin/holy/sim.js';
import { ProtectionPaladinSimUI } from '../paladin/protection/sim.js';
import { RetributionPaladinSimUI } from '../paladin/retribution/sim.js';
import { ShadowPriestSimUI } from '../priest/shadow/sim.js';
import { ElementalShamanSimUI } from '../shaman/elemental/sim.js';
import { EnhancementShamanSimUI } from '../shaman/enhancement/sim.js';
import { RestorationShamanSimUI } from '../shaman/restoration/sim.js';

export const specSimFactories: Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>> = {
	// Death Knight
	[Spec.SpecBloodDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new BloodDeathknightSimUI(parentElem, player),
	[Spec.SpecFrostDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new FrostDeathKnightSimUI(parentElem, player),
	[Spec.SpecUnholyDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new UnholyDeathknightSimUI(parentElem, player),
	// Druid
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	// Hunter
	[Spec.SpecBeastMasteryHunter]: (parentElem: HTMLElement, player: Player<any>) => new BeastMasteryHunterSimUI(parentElem, player),
	[Spec.SpecMarksmanshipHunter]: (parentElem: HTMLElement, player: Player<any>) => new MarksmanshipHunterSimUI(parentElem, player),
	[Spec.SpecSurvivalHunter]: (parentElem: HTMLElement, player: Player<any>) => new SurvivalHunterSimUI(parentElem, player),
	// Mage
	// Paladin
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	// Priest
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	// Rogue
	// Shaman
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	// Warlock
	// Warrior
};

export const playerPresets: Array<RaidSimPreset<any>> = naturalPlayerClassOrder
	.map(playerClass => Object.values(playerClass.specs))
	.flat()
	.map(spec => getSpecConfig(spec.protoID))
	.map(config => {
		const indSimUiConfig = config as IndividualSimUIConfig<any>;
		return indSimUiConfig.raidSimPresets;
	})
	.flat();

export const implementedSpecs: Array<Spec> = [...new Set(playerPresets.map(preset => preset.spec))];
