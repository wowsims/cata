import * as InputHelpers from "../../core/components/input_helpers";
import { Spec } from "../../core/proto/common";

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingHolyPower = () =>
	InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRetributionPaladin>({
		fieldName: 'startingHolyPower',
		label: 'Starting Holy Power',
		labelTooltip: "Initial Holy Power at the start of each iteration.",
	});
