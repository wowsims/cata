import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { CacheHandler } from '../../../cache_handler';
import { TOOLTIP_METRIC_LABELS } from '../../../constants/tooltips';
import { ActionId } from '../../../proto_utils/action_id';
import { ActionMetrics, AuraMetrics, ResourceMetrics, UnitMetrics } from '../../../proto_utils/sim_result';
import { TypedEvent } from '../../../typed_event';
import { ResultComponent, ResultComponentConfig, SimResultData } from '../result_component';

declare let $: any;

export enum ColumnSortType {
	None,
	Ascending,
	Descending,
}

export interface MetricsColumnConfig<T> {
	name: string;
	tooltip?: string;
	headerCellClass?: string;
	columnClass?: string;
	sort?: ColumnSortType;

	getValue?: (metric: T, isChildRow?: boolean) => number;

	// Either getDisplayString or fillCell must be specified.
	getDisplayString?: (metric: T, isChildRow?: boolean) => string;
	fillCell?: (metric: T, cellElem: HTMLElement, rowElem: HTMLElement, isChildRow?: boolean) => void;
}

const cachedMetricsTableIcon = new CacheHandler<HTMLAnchorElement>();

export abstract class MetricsTable<T extends ActionMetrics | AuraMetrics | UnitMetrics | ResourceMetrics> extends ResultComponent {
	private readonly columnConfigs: Array<MetricsColumnConfig<T>>;

	protected readonly tableElem: HTMLElement;
	protected readonly bodyElem: HTMLElement;

	readonly onUpdate = new TypedEvent<void>('MetricsTableUpdate');

	constructor(config: ResultComponentConfig, columnConfigs: Array<MetricsColumnConfig<T>>) {
		super(config);
		this.columnConfigs = columnConfigs;

		this.rootElem.appendChild(
			<table className="metrics-table tablesorter">
				<thead className="metrics-table-header">
					<tr className="metrics-table-header-row"></tr>
				</thead>
				<tbody className="metrics-table-body"></tbody>
			</table>,
		);

		this.tableElem = this.rootElem.querySelector<HTMLTableSectionElement>('.metrics-table')!;
		this.bodyElem = this.rootElem.querySelector<HTMLTableSectionElement>('.metrics-table-body')!;

		const headerRowElem = this.rootElem.querySelector<HTMLTableRowElement>('.metrics-table-header-row')!;
		this.columnConfigs.forEach(columnConfig => {
			const headerCell = document.createElement('th');
			const tooltipContent = columnConfig.tooltip || TOOLTIP_METRIC_LABELS[columnConfig.name as keyof typeof TOOLTIP_METRIC_LABELS];
			headerCell.classList.add('metrics-table-header-cell');
			if (columnConfig.columnClass) {
				headerCell.classList.add(...columnConfig.columnClass.split(' '));
			}
			if (columnConfig.headerCellClass) {
				headerCell.classList.add(...columnConfig.headerCellClass.split(' '));
			}
			headerCell.appendChild(<span>{columnConfig.name}</span>);
			if (tooltipContent) {
				const tooltip = tippy(headerCell, {
					content: tooltipContent,
					ignoreAttributes: true,
				});
				this.addOnResetCallback(() => tooltip.destroy());
			}
			headerRowElem.appendChild(headerCell);
		});

		const sortList = this.columnConfigs
			.map((config, i) => [i, config.sort == ColumnSortType.Ascending ? 0 : 1])
			.filter(sortData => this.columnConfigs[sortData[0]].sort);

		$(this.tableElem).tablesorter({
			sortList: sortList,
			cssChildRow: 'child-metric',
		});
	}

	protected sortMetrics(metrics: Array<T>) {
		this.columnConfigs
			.filter(config => config.sort)
			.forEach(config => {
				if (!config.getValue) {
					throw new Error("Can' apply group sorting without getValue");
				}
				if (config.sort == ColumnSortType.Ascending) {
					metrics.sort((a, b) => config.getValue!(a) - config.getValue!(b));
				} else {
					metrics.sort((a, b) => config.getValue!(b) - config.getValue!(a));
				}
			});
	}

	private addRow(metric: T, isChildRow = false): HTMLElement {
		const rowElem = document.createElement('tr');
		this.bodyElem.appendChild(rowElem);

		this.columnConfigs.forEach(columnConfig => {
			const cellElem = document.createElement('td');
			if (columnConfig.getValue) {
				cellElem.dataset.text = String(columnConfig.getValue(metric, isChildRow));
			}
			if (columnConfig.columnClass) {
				cellElem.classList.add(columnConfig.columnClass);
			}
			if (columnConfig.fillCell) {
				columnConfig.fillCell(metric, cellElem, rowElem, isChildRow);
			} else if (columnConfig.getDisplayString) {
				cellElem.textContent = columnConfig.getDisplayString(metric, isChildRow);
			} else {
				throw new Error('Metrics column config does not provide content function: ' + columnConfig.name);
			}
			rowElem.appendChild(cellElem);
		});

		this.customizeRowElem(metric, rowElem, isChildRow);
		return rowElem;
	}

	private addGroup(metrics: Array<T>) {
		if (metrics.length == 0) {
			return;
		}

		if (metrics.length == 1 && this.shouldCollapse(metrics[0])) {
			this.addRow(metrics[0]);
			return;
		}

		// Manually sort because tablesorter doesn't let us apply sorting to child rows.
		this.sortMetrics(metrics);

		const mergedMetrics = this.mergeMetrics(metrics);
		const parentRow = this.addRow(mergedMetrics);
		const childRows = metrics.map(metric => this.addRow(metric, true));
		childRows.forEach(childRow => childRow.classList.add('child-metric'));

		let expand = true;
		parentRow.classList.add('parent-metric', 'expand');
		parentRow.addEventListener('click', () => {
			expand = !expand;
			if (expand) {
				childRows.forEach(row => row.classList.remove('hide'));
				parentRow.classList.add('expand');
			} else {
				childRows.forEach(row => row.classList.add('hide'));
				parentRow.classList.remove('expand');
			}
		});
	}

	onSimResult(resultData: SimResultData) {
		this.reset();
		const groupedMetrics = this.getGroupedMetrics(resultData).filter(group => group.length > 0);
		if (groupedMetrics.length) {
			this.rootElem.classList.remove('hide');
		} else {
			this.rootElem.classList.add('hide');
			this.onUpdate.emit(resultData.eventID);
			return;
		}

		groupedMetrics.forEach(group => this.addGroup(group));
		$(this.tableElem).trigger('update');
		this.onUpdate.emit(resultData.eventID);
	}

	reset() {
		super.reset();
		this.bodyElem.replaceChildren();
	}

	// Whether a single-element group should have its parent row removed.
	// Override this to add custom behavior.
	protected shouldCollapse(metric: T): boolean {
		return true;
	}

	// Override this to customize rowElem after it has been populated.
	protected customizeRowElem(metric: T, rowElem: HTMLElement, isChildRow = false) {
		return;
	}

	// Override this to provide custom merge behavior.
	protected mergeMetrics(metrics: Array<T>): T {
		return metrics[0];
	}

	// Returns grouped metrics to display.
	abstract getGroupedMetrics(resultData: SimResultData): Array<Array<T>>;

	static nameCellConfig<T>(
		getData: (metric: T) => {
			name: string;
			actionId: ActionId;
			metricType: string;
		} & Pick<MetricsColumnConfig<T>, 'columnClass' | 'headerCellClass'>,
	): MetricsColumnConfig<T> {
		return {
			name: 'Name',
			fillCell: (metric: T, cellElem: HTMLElement, rowElem: HTMLElement) => {
				const data = getData(metric);
				const actionIdAsString = data.actionId.toString();
				const iconElemRef = ref<HTMLAnchorElement>();
				const iconElem = cachedMetricsTableIcon.get(actionIdAsString);
				cellElem.appendChild(
					<div className="metrics-action">
						{iconElem?.cloneNode() || <a ref={iconElemRef} className="metrics-action-icon"></a>}
						<span className="metrics-action-name text-truncate">{data.name}</span>
						<span className="expand-toggle fa fa-caret-right"></span>
						<span className="expand-toggle fa fa-caret-down"></span>
					</div>,
				);
				if (!iconElem && iconElemRef.value) {
					data.actionId.setBackgroundAndHref(iconElemRef.value);
					data.actionId.setWowheadDataset(iconElemRef.value, {
						useBuffAura: data.metricType === 'AuraMetrics',
					});
					cachedMetricsTableIcon.set(actionIdAsString, iconElemRef.value);
				}
			},
		};
	}

	static playerNameCellConfig(): MetricsColumnConfig<UnitMetrics> {
		return {
			name: 'Name',
			columnClass: 'name-cell',
			fillCell: (player: UnitMetrics, cellElem: HTMLElement, rowElem: HTMLElement) => {
				cellElem.appendChild(
					<>
						<img className="metrics-action-icon" src={player.iconUrl}></img>
						<span className={`metrics-action-name text-${player.classColor}`}>{player.label}</span>
					</>,
				);
			},
		};
	}
}
