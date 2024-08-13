export const BUFFS_SECTION = 'Buffs provided by other party/raid members.';
export const DEBUFFS_SECTION = 'Debuffs applied by other raid members.';
export const COOLDOWNS_SECTION =
	'Specify cooldown timings, in seconds. Cooldowns will be used as soon as possible after their specified timings. When not specified, cooldowns will be used when ready and it is sensible to do so.<br><br>Multiple timings can be provided by separating with commas. Any cooldown usages after the last provided timing will use the default logic.';
export const BLESSINGS_SECTION =
	'Specify Paladin Blessings for each role, in order of priority. Blessings in the 1st column will be used if there is at least 1 Paladin in the raid, 2nd column if at least 2, etc.';

export const BASIC_BIS_DISCLAIMER =
	"<p>Preset gear lists are intended as rough approximations of BIS, and will often not be the absolute highest-DPS setup for you. Your optimal gear setup will depend on many factors; that's why we have a sim!</p><p>Items may also be omitted from the presets if they are highly contested and clearly better utilized on other classes, to encourage equitable gearing for the raid as a whole.</p>";

export const HEALING_SIM_DISCLAIMER =
	'*** WARNING - USE AT YOUR OWN RISK ***\n\nThe entire concept of a healing sim is EXPERIMENTAL. All results should be taken with an EXTREMELY large grain of salt.\n\nThis tool is currently more similar to a spreadsheet than a true sim. Options for specifying incoming damage profiles in order to have proper reactive rotations have not yet been added.';
export const EP_TOOLTIP = `
	EP (Equivalence Points) is way of comparing items by multiplying the raw stats of an item with your current stat weights.
	More EP does not necessarily mean more DPS, as EP doesn't take into account stat caps and non-linear stat calculations.
`;

export const TOOLTIP_METRIC_LABELS = {
	// Damage metrics
	Damage: 'Total Damage done',
	DPS: 'Damage / Encounter Duration',
	TPS: 'Threat / Encounter Duration',
	DPET: 'Damage / Avg Cast Time',
	'Damage Avg Cast': 'Damage / Casts',
	'Avg Hit': 'Damage / (Hits + Crits + Glances + Blocks)',
	// Healing metrics
	Healing: 'Total Healing done',
	'Healing Avg Cast': 'Healing / Casts',
	HPM: 'Healing / Mana',
	HPET: 'Healing / Avg Cast Time',
	HPS: 'Healing / Encounter Duration',
	// Damage taken metrics
	'Damage Taken': 'Total Damage taken',
	DTPS: 'Damage Taken / Encounter Duration',
	COD: 'Chance of Death',
	// Cast metrics
	Casts: 'Casts',
	CPM: 'Casts / (Encounter Duration / 60 Seconds)',
	'Cast Time': 'Average cast time in seconds',
	// Hit metrics
	Hits: 'Hits + Crits + Glances + Blocks',
	'Crit %': 'Crits / Hits',
	'Hit Miss %': 'Misses / (Hits + Crits + Glances + Blocks)',
	'Cast Miss %': 'Misses / Casts',
	// Encounter
	DUR: 'Encounter Duration',
	OOM: 'Spent Out of Mana',
	TTO: 'Time to Out of Mana in seconds',
	// Aura metrcis
	Procs: 'Procs',
	PPM: 'Procs Per Minute',
	Uptime: 'Uptime / Encounter Duration',
	// Resource Metrics
	Gain: 'Gain',
	'Gain / s': 'Gain / Second',
	'Avg Gain': 'Gain / Event',
	'Wasted Gain': 'Gain that was wasted because of resource cap.',
} as const;
