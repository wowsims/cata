import * as InputHelpers from '../core/components/input_helpers';
import { HunterSpecs } from '../core/proto_utils/utils';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet';
// import { makePetTypeInputConfig } from '../core/talents/hunter_pet';

// // Configuration for class-specific UI elements on the settings tab.
// // These don't need to be in a separate file but it keeps things cleaner.

// export const PetTypeInput = <SpecType extends HunterSpecs>() => makePetTypeInputConfig<SpecType>();
export const PetTypeInput = <SpecType extends HunterSpecs>() => makePetTypeInputConfig<SpecType>();

export const PetUptime = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'petUptime',
		label: 'Pet Uptime (%)',
		labelTooltip: 'Percent of the fight duration for which your pet will be alive.',
		percent: true,
	});

export const AQTierPrepull = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'useAqTier',
		label: 'Prepull T2.5',
		labelTooltip: 'Uses T2.5 during prepull',
	});
export const NaxxTierPrepull = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'useNaxxTier',
		label: 'Prepull T3',
		labelTooltip: 'Use T3 during prepull.',
	});

export const GlaiveTossChance = <SpecType extends HunterSpecs>() =>
	InputHelpers.makeClassOptionsNumberInput<SpecType>({
		fieldName: 'glaiveTossSuccess',
		label: 'Glaive Toss Success %',
		labelTooltip: 'The chance that Glaive Toss hits secondary targets in percentages.',
	});
