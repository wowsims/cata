import { Spec } from '../proto/common';
import { SecondaryResourceType } from '../proto/spell';

export interface SecondaryResourceConfig {
	name: string;
	icon: string;
}

export const SECONDARY_RESOURCES = new Map<SecondaryResourceType, SecondaryResourceConfig>([
	[
		SecondaryResourceType.SecondaryResourceTypeArcaneCharges,
		{
			name: 'Arcane Charges',
			icon: 'https://wow.zamimg.com/images/wow/icons/medium/spell_arcane_arcane01.jpg',
		},
	],
	[
		SecondaryResourceType.SecondaryResourceTypeShadowOrbs,
		{
			name: 'Shadow Orbs',
			icon: 'https://wow.zamimg.com/images/wow/icons/medium/spell_priest_shadoworbs.jpg',
		},
	],
	[
		SecondaryResourceType.SecondaryResourceTypeDemonicFury,
		{
			name: 'Demonic Fury',
			icon: 'https://wow.zamimg.com/images/wow/icons/medium/ability_warlock_eradication.jpg',
		},
	],
	[
		SecondaryResourceType.SecondaryResourceTypeHolyPower,
		{
			name: 'Holy Power',
			icon: 'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_holybolt.jpg',
		},
	],
	[
		SecondaryResourceType.SecondaryResourceTypeBurningEmbers,
		{
			name: 'Burning Embers',
			icon: 'https://wow.zamimg.com/images/wow/icons/medium/inv_mace_2h_pvp410_c_01.jpg',
		},
	],
	[
		SecondaryResourceType.SecondaryResourceTypeSoulShards,
		{
			name: 'Soul Shards',
			icon: 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_gem_amethyst_02.jpg',
		},
	],
]);

const RESOURCE_TYPE_PER_SPEC = new Map<Spec, SecondaryResourceType>([
	// Paladin
	[Spec.SpecRetributionPaladin, SecondaryResourceType.SecondaryResourceTypeHolyPower],
	[Spec.SpecProtectionPaladin, SecondaryResourceType.SecondaryResourceTypeHolyPower],
	[Spec.SpecHolyPaladin, SecondaryResourceType.SecondaryResourceTypeHolyPower],
	// Warlock
	[Spec.SpecAfflictionWarlock, SecondaryResourceType.SecondaryResourceTypeSoulShards],
	[Spec.SpecDemonologyWarlock, SecondaryResourceType.SecondaryResourceTypeDemonicFury],
	[Spec.SpecDestructionWarlock, SecondaryResourceType.SecondaryResourceTypeBurningEmbers],
	// Priest
	[Spec.SpecShadowPriest, SecondaryResourceType.SecondaryResourceTypeShadowOrbs],
	// Mage
	[Spec.SpecArcaneMage, SecondaryResourceType.SecondaryResourceTypeArcaneCharges],
]);

class SecondaryResource {
	private readonly config: SecondaryResourceConfig | null;
	constructor(spec: Spec) {
		const type = SecondaryResource.getGenericResourcesForSpec(spec);
		this.config = type || null;
	}

	get name() {
		return this.config?.name;
	}
	get icon() {
		return this.config?.icon;
	}

	static hasSecondaryResource(spec: Spec): boolean {
		return RESOURCE_TYPE_PER_SPEC.has(spec);
	}

	static getGenericResourcesForSpec(spec: Spec) {
		const type = RESOURCE_TYPE_PER_SPEC.get(spec);
		if (!type) return null;
		return SECONDARY_RESOURCES.get(type);
	}

	static create(spec: Spec): SecondaryResource | null {
		if (!SecondaryResource.hasSecondaryResource(spec)) return null;
		return new SecondaryResource(spec);
	}

	replaceResourceName(value: string) {
		return value.replaceAll(/{GENERIC_RESOURCE}/g, this.name || '');
	}
}

export default SecondaryResource;
