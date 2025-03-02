export enum SimSettingCategories {
	Gear = 0,
	Talents,
	Rotation,
	Consumes,
	Miscellaneous, // Spec-specific settings, Distance from target, tank status, etc
	External, // Buffs and debuffs
	Encounter,
	UISettings, // # iterations, EP weights, filters, etc
}

export const SIM_CATEGORY_KEYS: Map<SimSettingCategories, string> = (() => {
	const map = new Map();
	// Use single-letter abbreviations since these will be included in sim links.
	map.set(SimSettingCategories.Gear, 'g');
	map.set(SimSettingCategories.Talents, 't');
	map.set(SimSettingCategories.Rotation, 'r');
	map.set(SimSettingCategories.Consumes, 'c');
	map.set(SimSettingCategories.Miscellaneous, 'm');
	map.set(SimSettingCategories.External, 'x');
	map.set(SimSettingCategories.Encounter, 'e');
	map.set(SimSettingCategories.UISettings, 'u');
	return map;
})();
