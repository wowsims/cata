import { IndividualSimUI, IndividualSimUIConfig, RaidSimPreset } from '../core/individual_sim_ui.js';
import { getSpecConfig, Player } from '../core/player.js';
import { PlayerClasses } from '../core/player_classes';
import { Spec } from '../core/proto/common.js';
import { BloodDeathKnightSimUI } from '../death_knight/blood/sim';
import { FrostDeathKnightSimUI } from '../death_knight/frost/sim';
import { UnholyDeathKnightSimUI } from '../death_knight/unholy/sim';
import { BalanceDruidSimUI } from '../druid/balance/sim.js';
import { FeralDruidSimUI } from '../druid/feral/sim.js';
import { GuardianDruidSimUI } from '../druid/guardian/sim';
import { RestorationDruidSimUI } from '../druid/restoration/sim.js';
import { BeastMasteryHunterSimUI } from '../hunter/beast_mastery/sim';
import { MarksmanshipHunterSimUI } from '../hunter/marksmanship/sim';
import { SurvivalHunterSimUI } from '../hunter/survival/sim';
import { ArcaneMageSimUI } from '../mage/arcane/sim';
import { FireMageSimUI } from '../mage/fire/sim';
import { FrostMageSimUI } from '../mage/frost/sim';
import { HolyPaladinSimUI } from '../paladin/holy/sim.js';
import { ProtectionPaladinSimUI } from '../paladin/protection/sim.js';
import { RetributionPaladinSimUI } from '../paladin/retribution/sim.js';
import { DisciplinePriestSimUI } from '../priest/discipline/sim';
import { HolyPriestSimUI } from '../priest/holy/sim';
import { ShadowPriestSimUI } from '../priest/shadow/sim.js';
import { AssassinationRogueSimUI } from '../rogue/assassination/sim';
import { CombatRogueSimUI } from '../rogue/combat/sim';
import { SubtletyRogueSimUI } from '../rogue/subtlety/sim';
import { ElementalShamanSimUI } from '../shaman/elemental/sim.js';
import { EnhancementShamanSimUI } from '../shaman/enhancement/sim.js';
import { RestorationShamanSimUI } from '../shaman/restoration/sim.js';
import { AfflictionWarlockSimUI } from '../warlock/affliction/sim';
import { DemonologyWarlockSimUI } from '../warlock/demonology/sim';
import { DestructionWarlockSimUI } from '../warlock/destruction/sim';
import { ArmsWarriorSimUI } from '../warrior/arms/sim';
import { FuryWarriorSimUI } from '../warrior/fury/sim';
import { ProtectionWarriorSimUI } from '../warrior/protection/sim';

export const specSimFactories: Partial<Record<Spec, (parentElem: HTMLElement, player: Player<any>) => IndividualSimUI<any>>> = {
	// Death Knight
	[Spec.SpecBloodDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new BloodDeathKnightSimUI(parentElem, player),
	[Spec.SpecFrostDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new FrostDeathKnightSimUI(parentElem, player),
	[Spec.SpecUnholyDeathKnight]: (parentElem: HTMLElement, player: Player<any>) => new UnholyDeathKnightSimUI(parentElem, player),
	// Druid
	[Spec.SpecBalanceDruid]: (parentElem: HTMLElement, player: Player<any>) => new BalanceDruidSimUI(parentElem, player),
	[Spec.SpecFeralDruid]: (parentElem: HTMLElement, player: Player<any>) => new FeralDruidSimUI(parentElem, player),
	[Spec.SpecRestorationDruid]: (parentElem: HTMLElement, player: Player<any>) => new RestorationDruidSimUI(parentElem, player),
	[Spec.SpecGuardianDruid]: (parentElem: HTMLElement, player: Player<any>) => new GuardianDruidSimUI(parentElem, player),
	// Hunter
	[Spec.SpecBeastMasteryHunter]: (parentElem: HTMLElement, player: Player<any>) => new BeastMasteryHunterSimUI(parentElem, player),
	[Spec.SpecMarksmanshipHunter]: (parentElem: HTMLElement, player: Player<any>) => new MarksmanshipHunterSimUI(parentElem, player),
	[Spec.SpecSurvivalHunter]: (parentElem: HTMLElement, player: Player<any>) => new SurvivalHunterSimUI(parentElem, player),
	// Mage
	[Spec.SpecArcaneMage]: (parentElem: HTMLElement, player: Player<any>) => new ArcaneMageSimUI(parentElem, player),
	[Spec.SpecFireMage]: (parentElem: HTMLElement, player: Player<any>) => new FireMageSimUI(parentElem, player),
	[Spec.SpecFrostMage]: (parentElem: HTMLElement, player: Player<any>) => new FrostMageSimUI(parentElem, player),
	// Paladin
	[Spec.SpecHolyPaladin]: (parentElem: HTMLElement, player: Player<any>) => new HolyPaladinSimUI(parentElem, player),
	[Spec.SpecProtectionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionPaladinSimUI(parentElem, player),
	[Spec.SpecRetributionPaladin]: (parentElem: HTMLElement, player: Player<any>) => new RetributionPaladinSimUI(parentElem, player),
	// Priest
	[Spec.SpecDisciplinePriest]: (parentElem: HTMLElement, player: Player<any>) => new DisciplinePriestSimUI(parentElem, player),
	[Spec.SpecHolyPriest]: (parentElem: HTMLElement, player: Player<any>) => new HolyPriestSimUI(parentElem, player),
	[Spec.SpecShadowPriest]: (parentElem: HTMLElement, player: Player<any>) => new ShadowPriestSimUI(parentElem, player),
	// Rogue
	[Spec.SpecAssassinationRogue]: (parentElem: HTMLElement, player: Player<any>) => new AssassinationRogueSimUI(parentElem, player),
	[Spec.SpecCombatRogue]: (parentElem: HTMLElement, player: Player<any>) => new CombatRogueSimUI(parentElem, player),
	[Spec.SpecSubtletyRogue]: (parentElem: HTMLElement, player: Player<any>) => new SubtletyRogueSimUI(parentElem, player),
	// Shaman
	[Spec.SpecElementalShaman]: (parentElem: HTMLElement, player: Player<any>) => new ElementalShamanSimUI(parentElem, player),
	[Spec.SpecEnhancementShaman]: (parentElem: HTMLElement, player: Player<any>) => new EnhancementShamanSimUI(parentElem, player),
	[Spec.SpecRestorationShaman]: (parentElem: HTMLElement, player: Player<any>) => new RestorationShamanSimUI(parentElem, player),
	// Warlock
	[Spec.SpecAfflictionWarlock]: (parentElem: HTMLElement, player: Player<any>) => new AfflictionWarlockSimUI(parentElem, player),
	[Spec.SpecDemonologyWarlock]: (parentElem: HTMLElement, player: Player<any>) => new DemonologyWarlockSimUI(parentElem, player),
	[Spec.SpecDestructionWarlock]: (parentElem: HTMLElement, player: Player<any>) => new DestructionWarlockSimUI(parentElem, player),
	// Warrior
	[Spec.SpecArmsWarrior]: (parentElem: HTMLElement, player: Player<any>) => new ArmsWarriorSimUI(parentElem, player),
	[Spec.SpecFuryWarrior]: (parentElem: HTMLElement, player: Player<any>) => new FuryWarriorSimUI(parentElem, player),
	[Spec.SpecProtectionWarrior]: (parentElem: HTMLElement, player: Player<any>) => new ProtectionWarriorSimUI(parentElem, player),
};

export const playerPresets: Array<RaidSimPreset<any>> = PlayerClasses.naturalOrder
	.map(playerClass => Object.values(playerClass.specs))
	.flat()
	.map(playerSpec => getSpecConfig(playerSpec.specID))
	.map(config => {
		const indSimUiConfig = config as IndividualSimUIConfig<any>;
		return indSimUiConfig.raidSimPresets;
	})
	.flat();

export const implementedSpecs: Array<any> = [...new Set(playerPresets.map(preset => preset.spec))];
