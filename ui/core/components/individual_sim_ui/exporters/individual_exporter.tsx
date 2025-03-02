import { SimSettingCategories } from '../../../constants/sim_settings';
import { IndividualSimUI } from '../../../individual_sim_ui';
import { Spec } from '../../../proto/common';
import { EventID } from '../../../typed_event';
import { getEnumValues } from '../../../utils';
import { Exporter, ExporterOptions } from '../../exporter';
import { BooleanPicker } from '../../pickers/boolean_picker';
import { IndividualImporter } from '../importers/individual_importer';

interface IndividualExporterOptions extends ExporterOptions {
	selectCategories?: boolean;
}

export abstract class IndividualExporter<SpecType extends Spec> extends Exporter {
	protected static readonly exportPickerConfigs: Array<{
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
	protected readonly simUI: IndividualSimUI<SpecType>;
	protected readonly exportCategories: Record<SimSettingCategories, boolean>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>, options: IndividualExporterOptions) {
		super(parent, options as ExporterOptions);

		this.simUI = simUI;
		const exportCategories: Partial<Record<SimSettingCategories, boolean>> = {};
		(getEnumValues(SimSettingCategories) as Array<SimSettingCategories>).forEach(
			cat => (exportCategories[cat] = IndividualImporter.DEFAULT_CATEGORIES.includes(cat)),
		);
		this.exportCategories = exportCategories as Record<SimSettingCategories, boolean>;

		if (options.selectCategories) {
			const categoryPickerContainer = (<div className="exporter-category-pickers" />) as HTMLElement;
			this.body.prepend(categoryPickerContainer);

			IndividualExporter.exportPickerConfigs.forEach(exportConfig => {
				const category = exportConfig.category;
				new BooleanPicker(categoryPickerContainer, this, {
					id: `link-exporter-${category}`,
					label: exportConfig.label,
					labelTooltip: exportConfig.labelTooltip,
					inline: true,
					getValue: () => this.exportCategories[category],
					setValue: (eventID: EventID, _modObj: IndividualExporter<SpecType>, newValue: boolean) => {
						this.exportCategories[category] = newValue;
						this.changedEvent.emit(eventID);
					},
					changedEvent: () => this.changedEvent,
				});
			});
		}
	}
}
