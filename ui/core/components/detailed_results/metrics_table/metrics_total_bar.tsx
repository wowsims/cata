import clsx from 'clsx';

import { SpellSchool } from '../../../proto/common';
import { spellSchoolNames } from '../../../proto_utils/names';
import { formatToCompactNumber, formatToPercent } from '../../../utils';

export type MetricsTotalBarProps = {
	percentage: number | undefined | null;
	max: number | null;
	total: number;
	value: number;
	// Used for overlayed value display, such as shielding.
	// Will show as darkened bar on top of main bar.
	overlayValue?: number;
	spellSchool?: SpellSchool | undefined | null;
	classColor?: string | undefined | null;
};
export const MetricsTotalBar = ({ percentage, max, total, value, overlayValue, spellSchool, classColor }: MetricsTotalBarProps) => {
	const spellSchoolString = typeof spellSchool === 'number' ? spellSchoolNames.get(spellSchool) : undefined;
	return (
		<div className="metrics-total position-relative d-flex justify-content-between w-100">
			<div className="metrics-total-percentage">{formatToPercent(percentage || 0)}</div>
			<div className="metrics-total-bar ms-1 me-1">
				<div
					className={clsx(
						'metrics-total-bar-fill',
						spellSchoolString && `bg-spell-school-${spellSchoolString.toLowerCase()}`,
						classColor && `bg-${classColor.toLowerCase()}`,
					)}
					style={{ '--percentage': formatToPercent((value / (max ?? 1)) * 100) }}></div>
				{overlayValue ? (
					<div
						className="metrics-total-bar-fill bg-black bg-opacity-25"
						style={{ '--percentage': formatToPercent((overlayValue / (max ?? 1)) * 100) }}
					/>
				) : undefined}
			</div>
			<div className="metrics-total-amount">{formatToCompactNumber(total)}</div>
		</div>
	);
};
