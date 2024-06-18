import * as InputHelpers from "../../core/components/input_helpers";
import {PaladinSpecs} from "../../core/proto_utils/utils";

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SnapshotGuardian = <SpecType extends PaladinSpecs>() =>
	InputHelpers.makeClassOptionsBooleanInput<SpecType>({
		fieldName: 'snapshotGuardian',
		label: 'Snapshot T11 Protection 4pc set bonus',
		labelTooltip: "Enable this to make the first Guardian of Ancient Kings cast during pre-pull snapshot the T11 Protection 4pc set bonus (50% increased duration).",
	});
