import tippy, { Props as TippyProps } from 'tippy.js';

import { SpellSchool } from '../../../proto/common';
import { formatToCompactNumber, sum } from '../../../utils';
import { MetricsTotalBar } from './metrics_total_bar';

type MetricsCombinedTableEntry = {
	name: string;
	value: number;
	percentage: number;
	min?: number;
	max?: number;
	average?: number;
};

type MetricsCombinedTooltipTableProps = {
	tooltipElement: HTMLElement;
	tooltipConfig?: Partial<TippyProps>;
	total: number;
	totalPercentage: number;
	values: MetricsCombinedTableEntry[];
	spellSchool: SpellSchool | undefined | null;
	hasMetricBars?: boolean;
	hasFooter?: boolean;
	headerValues?: (MetricsCombinedTableEntry['name'] | undefined)[];
};
export const MetricsCombinedTooltipTable = ({
	tooltipElement,
	tooltipConfig,
	headerValues,
	total,
	totalPercentage,
	values,
	spellSchool,
	hasMetricBars = true,
	hasFooter = true,
}: MetricsCombinedTooltipTableProps) => {
	const displayedValues = values.filter(v => v.value);
	const maxValue = Math.max(...displayedValues.map(a => a.value));
	const hasAverageColumn = displayedValues.some(d => typeof d.average === 'number');
	if (displayedValues.length) {
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
						{displayedValues
							.sort((a, b) => b.value - a.value)
							.map(({ name, value, percentage, average }) => (
								<tr>
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
							))}
					</tbody>
					{hasFooter && displayedValues.length > 1 && (
						<tfoot className="metrics-table-footer">
							<tr className="metrics-table-footer-row">
								<td>Total</td>
								<td className="text-end">{formatToCompactNumber(total)}</td>
								{hasAverageColumn ? (
									<td>{formatToCompactNumber(sum(displayedValues.map(v => v.average || 0)) / displayedValues.length)}</td>
								) : undefined}
							</tr>
						</tfoot>
					)}
				</table>
			),
		});
	}
	return null;
};
