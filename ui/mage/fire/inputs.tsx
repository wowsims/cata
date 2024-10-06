import * as InputHelpers from '../../core/components/input_helpers';
import { Spec } from '../../core/proto/common';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                       FIRE INPUTS
		// ********************************************************
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'igniteCombustThreshold',
			label: 'Ignite Combust Threshold',
			labelTooltip: 'The Ignite damage threshold to use Combustion during Bloodlust',
			description: (
				<>
					<p>Should be set to the Ignite damage threshold at which you want to use Combustion during Bloodlust.</p>
					<p>You can check the Sim or your logs to find a value that is feasible to hit.</p>
					<p>Furthermore a % of this value will be used for other Combust usages outside of the Bloodlust window:</p>
					<p>Example: Setting the Ignite Combust Threshold to 30.000 will:</p>
					<ul>
						<li>Cast Combust whilst Bloodlust is active when Ignite exceeds 30.000 (100%) damage.</li>
						<li>Cast Combust at the last moment when Bloodlust (+ Berserking) is running out when Ignite exceeds 10.000 (33%) damage</li>
						<li>Cast Combust outside of Bloodlust when Ignite exceeds 15.000 (50%) damage</li>
						<li>Cast Combust when encounter is ending in 15 seconds when Ignite exceeds 10.000 (33%) damage</li>
					</ul>
				</>
			),
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().igniteCombustThreshold,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'igniteLastMomentLustPercentage',
			label: 'Max Ignite % - Last moment during Bloodlust',
			labelTooltip: 'Using a % of the "Ignite Combust Threshold" to use Combustion when Bloodlust (+ Berserking) is about to run out.',
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().igniteLastMomentLustPercentage,
			percent: true,
			max: 1,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFireMage>({
			fieldName: 'igniteNoLustPercentage',
			label: 'Max Ignite % - Outside of Bloodlust',
			labelTooltip: 'Using a % of the "Ignite Combust Threshold" to use Combustion when Bloodlust is not up.',
			changeEmitter: player => player.rotationChangeEmitter,
			getValue: player => player.getSimpleRotation().igniteNoLustPercentage,
			percent: true,
			max: 1,
			positive: true,
		}),
	],
};
