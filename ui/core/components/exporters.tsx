import { default as pako } from 'pako';
import { ref } from 'tsx-vanilla';

import * as Mechanics from '../constants/mechanics';
import { IndividualSimUI } from '../individual_sim_ui';
import { RaidSimRequest } from '../proto/api';
import { Consumes, ItemSlot, PseudoStat, Spec, Stat } from '../proto/common';
import { IndividualSimSettings } from '../proto/ui';
import { raceNames } from '../proto_utils/names';
import { UnitStat } from '../proto_utils/stats';
import { SimSettingCategories } from '../sim';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { arrayEquals, downloadString, getEnumValues, jsonStringifyWithFlattenedPaths } from '../utils';
import { BaseModal } from './base_modal';
import { CopyButton } from './copy_button';
import { IndividualLinkImporter, IndividualWowheadGearPlannerImporter } from './importers';
import { BooleanPicker } from './pickers/boolean_picker';
import { createWowheadGearPlannerLink, WowheadGearPlannerData, WowheadItemData } from './wowhead_helper';

interface ExporterOptions {
	title: string;
	header?: boolean;
	allowDownload?: boolean;
}

export abstract class Exporter extends BaseModal {
	private readonly textElem: HTMLElement;
	protected readonly changedEvent: TypedEvent<void> = new TypedEvent();

	constructor(parent: HTMLElement, simUI: SimUI, options: ExporterOptions) {
		super(parent, 'exporter', { title: options.title, header: true, footer: true });

		this.body.innerHTML = `
			<textarea spellCheck="false" class="exporter-textarea form-control"></textarea>
		`;
		this.textElem = this.rootElem.getElementsByClassName('exporter-textarea')[0] as HTMLElement;

		new CopyButton(this.footer!, {
			extraCssClasses: ['btn-primary', 'me-2'],
			getContent: () => this.textElem.innerHTML,
			text: 'Copy',
			tooltip: 'Copy to clipboard',
		});

		if (options.allowDownload) {
			const downloadBtnRef = ref<HTMLButtonElement>();
			this.footer!.appendChild(
				<button className="exporter-button btn btn-primary download-button" ref={downloadBtnRef}>
					<i className="fa fa-download me-1"></i>
					Download
				</button>,
			);

			const downloadButton = downloadBtnRef.value!;
			downloadButton.addEventListener('click', _event => {
				const data = this.textElem.textContent!;
				downloadString(data, 'wowsims.json');
			});
		}
	}

	protected init() {
		this.changedEvent.on(() => this.updateContent());
		this.updateContent();
	}

	private updateContent() {
		this.textElem.textContent = this.getData();
	}

	abstract getData(): string;
}

export class IndividualLinkExporter<SpecType extends Spec> extends Exporter {
	private static readonly exportPickerConfigs: Array<{
		category: SimSettingCategories;
		label: string;
		labelTooltip: string;
	}> = [
		{
			category: SimSettingCategories.Gear,
			label: 'Gear',
			labelTooltip: 'Also includes bonus stats and weapon swaps.',
		},
		{
			category: SimSettingCategories.Talents,
			label: 'Talents',
			labelTooltip: 'Talents and Glyphs.',
		},
		{
			category: SimSettingCategories.Rotation,
			label: 'Rotation',
			labelTooltip: 'Includes everything found in the Rotation tab.',
		},
		{
			category: SimSettingCategories.Consumes,
			label: 'Consumes',
			labelTooltip: 'Flask, pots, food, etc.',
		},
		{
			category: SimSettingCategories.External,
			label: 'Buffs & Debuffs',
			labelTooltip: 'All settings which are applied by other raid members.',
		},
		{
			category: SimSettingCategories.Miscellaneous,
			label: 'Misc',
			labelTooltip: 'Spec-specific settings, front/back of target, distance from target, etc.',
		},
		{
			category: SimSettingCategories.Encounter,
			label: 'Encounter',
			labelTooltip: 'Fight-related settings.',
		},
		// Intentionally exclude UISettings category here, because users almost
		// never intend to export them and it messes with other users' settings.
		// If they REALLY want to export UISettings, just use the JSON exporter.
	];

	private readonly simUI: IndividualSimUI<SpecType>;
	private readonly exportCategories: Record<SimSettingCategories, boolean>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Sharable Link' });
		this.simUI = simUI;

		const exportCategories: Partial<Record<SimSettingCategories, boolean>> = {};
		(getEnumValues(SimSettingCategories) as Array<SimSettingCategories>).forEach(
			cat => (exportCategories[cat] = IndividualLinkImporter.DEFAULT_CATEGORIES.includes(cat)),
		);
		this.exportCategories = exportCategories as Record<SimSettingCategories, boolean>;

		const pickersContainer = document.createElement('div');
		pickersContainer.classList.add('link-exporter-pickers');
		this.body.prepend(pickersContainer);

		IndividualLinkExporter.exportPickerConfigs.forEach(exportConfig => {
			const category = exportConfig.category;
			new BooleanPicker(pickersContainer, this, {
				id: `link-exporter-${category}`,
				label: exportConfig.label,
				labelTooltip: exportConfig.labelTooltip,
				inline: true,
				getValue: () => this.exportCategories[category],
				setValue: (eventID: EventID, modObj: IndividualLinkExporter<SpecType>, newValue: boolean) => {
					this.exportCategories[category] = newValue;
					this.changedEvent.emit(eventID);
				},
				changedEvent: () => this.changedEvent,
			});
		});
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		return IndividualLinkExporter.createLink(
			this.simUI,
			(getEnumValues(SimSettingCategories) as Array<SimSettingCategories>).filter(c => this.exportCategories[c]),
		);
	}

	static createLink(simUI: IndividualSimUI<any>, exportCategories?: Array<SimSettingCategories>): string {
		if (!exportCategories) {
			exportCategories = IndividualLinkImporter.DEFAULT_CATEGORIES;
		}

		const proto = simUI.toProto(exportCategories);

		const protoBytes = IndividualSimSettings.toBinary(proto);
		// @ts-ignore Pako did some weird stuff between versions and the @types package doesn't correctly support this syntax for version 2.0.4 but it's completely valid
		// The syntax was removed in 2.1.0 and there were several complaints but the project seems to be largely abandoned now
		const deflated = pako.deflate(protoBytes, { to: 'string' });
		const encoded = btoa(String.fromCharCode(...deflated));

		const linkUrl = new URL(window.location.href);
		linkUrl.hash = encoded;
		if (arrayEquals(exportCategories, IndividualLinkImporter.DEFAULT_CATEGORIES)) {
			linkUrl.searchParams.delete(IndividualLinkImporter.CATEGORY_PARAM);
		} else {
			const categoryCharString = exportCategories.map(c => IndividualLinkImporter.CATEGORY_KEYS.get(c)).join('');
			linkUrl.searchParams.set(IndividualLinkImporter.CATEGORY_PARAM, categoryCharString);
		}
		return linkUrl.toString();
	}
}

export class IndividualJsonExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'JSON Export', allowDownload: true });
		this.simUI = simUI;
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		return jsonStringifyWithFlattenedPaths(IndividualSimSettings.toJson(this.simUI.toProto()), 2, (value, path) => {
			if (['stats', 'pseudoStats'].includes(path[path.length - 1])) {
				return true;
			}

			if (['player', 'equipment', 'items'].every((v, i) => path[i] == v)) {
				return path.length > 3;
			}

			if (path[0] == 'player' && path[1] == 'rotation' && ['prepullActions', 'priorityList'].includes(path[2])) {
				return path.length > 3;
			}

			return false;
		});
	}
}
export class IndividualWowheadGearPlannerExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Wowhead Export', allowDownload: true });
		this.simUI = simUI;
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;

		const classStr = player.getPlayerClass().friendlyName.replaceAll(/\s/g, '-').toLowerCase();
		const raceStr = raceNames.get(player.getRace())!.replaceAll(/\s/g, '-').toLowerCase();
		const url = `https://www.wowhead.com/cata/gear-planner/${classStr}/${raceStr}/`;

		const addGlyph = (glyphItemId: number): number => {
			const spellId = this.simUI.sim.db.glyphItemToSpellId(glyphItemId);
			if (!spellId) {
				return 0;
			}
			return spellId;
		};

		const glyphs = player.getGlyphs();

		const data = {
			level: Mechanics.CHARACTER_LEVEL,
			talents: player.getTalentsString().split("-"),
			glyphs: [
				addGlyph(glyphs.prime1),
				addGlyph(glyphs.prime2),
				addGlyph(glyphs.prime3),
				addGlyph(glyphs.major1),
				addGlyph(glyphs.major2),
				addGlyph(glyphs.major3),
				addGlyph(glyphs.minor1),
				addGlyph(glyphs.minor2),
				addGlyph(glyphs.minor3),
			],
			items: []
		} as WowheadGearPlannerData

		const gear = player.getGear();

		gear.getItemSlots()
			.sort((slot1, slot2) => IndividualWowheadGearPlannerImporter.slotIDs[slot1] - IndividualWowheadGearPlannerImporter.slotIDs[slot2])
			.forEach(itemSlot => {
				const item = gear.getEquippedItem(itemSlot);
				if (!item) {
					return;
				}

				const slotId = IndividualWowheadGearPlannerImporter.slotIDs[itemSlot];
				const itemData = {
					slotId: slotId,
					itemId: item.id,

				} as WowheadItemData
				if(item._randomSuffix?.id) {
					itemData.randomEnchantId = item._randomSuffix.id
				}
				if(item._enchant) {
					itemData.enchantIds = [item._enchant.spellId]
				}

				if (ItemSlot.ItemSlotHands == itemSlot) {
					//Todo: IF Hands we want to append any tinkers if existing
				}

				if(item._gems) {
					itemData.gemItemIds = item._gems.map(gem => {return gem?.id ?? 0})
				}
				if(item._reforge) {
					itemData.reforge = item._reforge.id
				}
				data.items.push(itemData)
			});

		const hash = createWowheadGearPlannerLink(data)

		return url + hash;
	}
}

export class Individual60UEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: '60Upgrades Cataclysm EP Export', allowDownload: true });
		this.simUI = simUI;
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats.forEach(stat => {
			const statName = Individual60UEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return (
			`https://sixtyupgrades.com/cata/ep/import?name=${encodeURIComponent(`${player.getPlayerSpec().friendlyName} WoWSims Weights`)}` +
			Object.keys(namesToWeights)
				.map(statName => `&${statName}=${namesToWeights[statName].toFixed(3)}`)
				.join('')
		);
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return Individual60UEPExporter.statNames[stat.getStat()];
		} else {
			return Individual60UEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'strength',
		[Stat.StatAgility]: 'agility',
		[Stat.StatStamina]: 'stamina',
		[Stat.StatIntellect]: 'intellect',
		[Stat.StatSpirit]: 'spirit',
		[Stat.StatSpellPower]: 'spellDamage',
		[Stat.StatMP5]: 'mp5',
		[Stat.StatHitRating]: 'hitRating',
		[Stat.StatCritRating]: 'critRating',
		[Stat.StatHasteRating]: 'hasteRating',
		[Stat.StatSpellPenetration]: 'spellPen',
		[Stat.StatAttackPower]: 'attackPower',
		[Stat.StatMasteryRating]: 'masteryRating',
		[Stat.StatExpertiseRating]: 'expertiseRating',
		[Stat.StatMana]: 'mana',
		[Stat.StatArmor]: 'armor',
		[Stat.StatRangedAttackPower]: 'attackPower',
		[Stat.StatDodgeRating]: 'dodgeRating',
		[Stat.StatParryRating]: 'parryRating',
		[Stat.StatResilienceRating]: 'resilienceRating',
		[Stat.StatHealth]: 'health',
		[Stat.StatArcaneResistance]: 'arcaneResistance',
		[Stat.StatFireResistance]: 'fireResistance',
		[Stat.StatFrostResistance]: 'frostResistance',
		[Stat.StatNatureResistance]: 'natureResistance',
		[Stat.StatShadowResistance]: 'shadowResistance',
		[Stat.StatBonusArmor]: 'armorBonus',
	};
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'dps',
		[PseudoStat.PseudoStatRangedDps]: 'rangedDps',
	};
}

export class IndividualPawnEPExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'Pawn EP Export', allowDownload: true });
		this.simUI = simUI;
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		const player = this.simUI.player;
		const epValues = player.getEpWeights();
		const allUnitStats = UnitStat.getAll();

		const namesToWeights: Record<string, number> = {};
		allUnitStats.forEach(stat => {
			const statName = IndividualPawnEPExporter.getName(stat);
			const weight = epValues.getUnitStat(stat);
			if (weight == 0 || statName == '') {
				return;
			}

			// Need to add together stats with the same name (e.g. hit/crit/haste).
			if (namesToWeights[statName]) {
				namesToWeights[statName] += weight;
			} else {
				namesToWeights[statName] = weight;
			}
		});

		return (
			`( Pawn: v1: "${player.getPlayerSpec().friendlyName} WoWSims Weights": Class=${player.getPlayerClass().friendlyName},` +
			Object.keys(namesToWeights)
				.map(statName => `${statName}=${namesToWeights[statName].toFixed(3)}`)
				.join(',') +
			' )'
		);
	}

	static getName(stat: UnitStat): string {
		if (stat.isStat()) {
			return IndividualPawnEPExporter.statNames[stat.getStat()];
		} else {
			return IndividualPawnEPExporter.pseudoStatNames[stat.getPseudoStat()] || '';
		}
	}

	static statNames: Record<Stat, string> = {
		[Stat.StatStrength]: 'Strength',
		[Stat.StatAgility]: 'Agility',
		[Stat.StatStamina]: 'Stamina',
		[Stat.StatIntellect]: 'Intellect',
		[Stat.StatSpirit]: 'Spirit',
		[Stat.StatSpellPower]: 'SpellDamage',
		[Stat.StatMP5]: 'Mp5',
		[Stat.StatHitRating]: 'HitRating',
		[Stat.StatCritRating]: 'CritRating',
		[Stat.StatHasteRating]: 'HasteRating',
		[Stat.StatSpellPenetration]: 'SpellPen',
		[Stat.StatAttackPower]: 'Ap',
		[Stat.StatMasteryRating]: 'MasteryRating',
		[Stat.StatExpertiseRating]: 'ExpertiseRating',
		[Stat.StatMana]: 'Mana',
		[Stat.StatArmor]: 'Armor',
		[Stat.StatRangedAttackPower]: 'Ap',
		[Stat.StatDodgeRating]: 'DodgeRating',
		[Stat.StatParryRating]: 'ParryRating',
		[Stat.StatResilienceRating]: 'ResilienceRating',
		[Stat.StatHealth]: 'Health',
		[Stat.StatArcaneResistance]: 'ArcaneResistance',
		[Stat.StatFireResistance]: 'FireResistance',
		[Stat.StatFrostResistance]: 'FrostResistance',
		[Stat.StatNatureResistance]: 'NatureResistance',
		[Stat.StatShadowResistance]: 'ShadowResistance',
		[Stat.StatBonusArmor]: 'Armor2',
	};
	static pseudoStatNames: Partial<Record<PseudoStat, string>> = {
		[PseudoStat.PseudoStatMainHandDps]: 'MeleeDps',
		[PseudoStat.PseudoStatRangedDps]: 'RangedDps',
	};
}

export class IndividualCLIExporter<SpecType extends Spec> extends Exporter {
	private readonly simUI: IndividualSimUI<SpecType>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
		super(parent, simUI, { title: 'CLI Export', allowDownload: true });
		this.simUI = simUI;
	}

	open() {
		super.open();
		this.init();
	}

	getData(): string {
		const raidSimJson: any = RaidSimRequest.toJson(this.simUI.sim.makeRaidSimRequest(false));
		delete raidSimJson.raid?.parties[0]?.players[0]?.database;
		return JSON.stringify(raidSimJson, null, 2);
	}
}
