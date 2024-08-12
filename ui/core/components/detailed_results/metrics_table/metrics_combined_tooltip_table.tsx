import clsx from 'clsx';

import { SpellSchool } from '../../../proto/common';
import { formatToCompactNumber, formatToNumber, formatToPercent, sum } from '../../../utils';
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
	total: number;
	totalPercentage: number;
	values: MetricsCombinedTableEntry[];
	spellSchool: SpellSchool | undefined | null;
};
export const MetricsCombinedTooltipTable = ({ total, totalPercentage, values, spellSchool }: MetricsCombinedTooltipTableProps) => {
	const displayedValues = values.filter(v => v.value);
	const maxValue = Math.max(...displayedValues.map(a => a.value));
	const hasAverageColumn = displayedValues.some(d => typeof d.average === 'number');

	return (
		<table className="metrics-table">
			<thead className="metrics-table-header">
				<tr className="metrics-table-header-row">
					<th className="metrics-table-header-cell">Type</th>
					<th className="metrics-table-header-cell">Count</th>
					{hasAverageColumn ? <th className="metrics-table-header-cell">Average</th> : undefined}
				</tr>
			</thead>
			<tbody className="metrics-table-body">
				{displayedValues
					.sort((a, b) => b.value - a.value)
					.map(({ name, value, percentage, average }) => (
						<tr>
							<td>{name}</td>
							<td>
								<MetricsTotalBar
									spellSchool={spellSchool}
									percentage={(percentage / totalPercentage) * 100}
									max={maxValue}
									total={value}
									value={value}
								/>
							</td>
							{typeof average === 'number' ? <td>{formatToCompactNumber(average)}</td> : undefined}
						</tr>
					))}
			</tbody>
			{displayedValues.length > 1 && (
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
	);
};
