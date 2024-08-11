import clsx from 'clsx';

import { SpellSchool } from '../../../proto/common';
import { formatToNumber, formatToPercent } from '../../../utils';
import { MetricsTotalBar } from './metrics_total_bar';

type MetricsCombinedTableEntry = {
	name: string;
	value: number;
	percentage: number;
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

	return (
		<table className="metrics-table">
			<thead className="metrics-table-header">
				<tr className="metrics-table-header-row">
					<th className="metrics-table-header-cell">Type</th>
					<th className="metrics-table-header-cell">Count</th>
				</tr>
			</thead>
			<tbody className="metrics-table-body">
				{displayedValues.map(({ name, value, percentage }) => (
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
					</tr>
				))}
			</tbody>
			{displayedValues.length > 1 && (
				<tfoot className="metrics-table-footer">
					<tr className="metrics-table-footer-row">
						<td>Total</td>
						<td className="text-end">{formatToNumber(total)}</td>
					</tr>
				</tfoot>
			)}
		</table>
	);
};
