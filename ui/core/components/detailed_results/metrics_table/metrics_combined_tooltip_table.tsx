import clsx from 'clsx';
import tippy, { Props as TippyProps } from 'tippy.js';

import { SpellSchool } from '../../../proto/common';
import { formatToCompactNumber } from '../../../utils';
import { MetricsTotalBar, MetricsTotalBarProps } from './metrics_total_bar';

type MetricsCombinedGroup = {
	name?: string;
	cssClass?: string;
	total?: number;
	totalPercentage: number;
	spellSchool: SpellSchool | undefined | null;
	data: MetricsCombinedTableEntry[];
};
type MetricsCombinedTableEntry = {
	name: string;
	min?: number;
	max?: number;
	average?: number;
	percentage: number;
} & Pick<MetricsTotalBarProps, 'value'>;

type MetricsCombinedTooltipTableProps = {
	tooltipElement: HTMLElement;
	tooltipConfig?: Partial<TippyProps>;
	groups: MetricsCombinedGroup[];
	hasMetricBars?: boolean;
	headerValues?: (MetricsCombinedGroup['name'] | undefined)[];
};
export const MetricsCombinedTooltipTable = ({
	tooltipElement,
	tooltipConfig,
	headerValues,
	groups,
	hasMetricBars = true,
}: MetricsCombinedTooltipTableProps) => {
	const displayGroups = groups.filter(group => group.data.some(d => d.value)).map(group => ({ ...group, data: group.data.filter(v => v.value) }));
	const hasAverageColumn = displayGroups.some(group => group.data.find(d => typeof d.average === 'number'));

	if (groups.length) {
		tippy(tooltipElement, {
			placement: 'auto',
			theme: 'metrics-table',
			maxWidth: 'none',
			...tooltipConfig,
			content: (
				<table className="metrics-table">
					<thead className="metrics-table-header">
						<tr className="metrics-table-header-row">
							<th className="metrics-table-header-cell">{headerValues?.[0] || 'Type'}</th>
							<th className="metrics-table-header-cell">{headerValues?.[1] || 'Count'}</th>
							{hasAverageColumn ? <th className="metrics-table-header-cell">{headerValues?.[2] || 'Average'}</th> : undefined}
						</tr>
					</thead>
					<tbody className="metrics-table-body">
						{displayGroups.map(({ name: groupName, cssClass, data, spellSchool, totalPercentage }) => {
							const maxValue = Math.max(...data.map(a => a.value));
							const columnCount = data.some(d => typeof d.average === 'number') ? 3 : 2;
							return (
								<>
									{groupName && displayGroups.length > 1 ? (
										<tr className={clsx('metrics-table-group-header', cssClass)}>
											<th className="text-start fw-normal" colSpan={columnCount}>
												{groupName}
											</th>
										</tr>
									) : undefined}
									{data
										.sort((a, b) => b.value - a.value)
										.map(({ name, value, percentage, average }) => (
											<>
												<tr className={clsx(cssClass)}>
													<td>{name}</td>
													<td>
														{hasMetricBars ? (
															<MetricsTotalBar
																spellSchool={spellSchool}
																percentage={(percentage / totalPercentage) * 100}
																max={maxValue}
																total={value}
																value={value}
															/>
														) : (
															formatToCompactNumber(value)
														)}
													</td>
													{typeof average === 'number' ? <td>{formatToCompactNumber(average)}</td> : undefined}
												</tr>
											</>
										))}
								</>
							);
						})}
					</tbody>
				</table>
			),
		});
	}
	return null;
};
