import tippy from 'tippy.js';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment, ref } from 'tsx-vanilla';

import { wowheadSupportedLanguages } from '../constants/lang.js';
import { Sim } from '../sim.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { BaseModal } from './base_modal.jsx';
import { BooleanPicker } from './boolean_picker.js';
import { EnumPicker } from './enum_picker.js';
import { NumberPicker } from './number_picker.js';
import Toast from './toast';

export class SettingsMenu extends BaseModal {
	private readonly simUI: SimUI;

	constructor(parent: HTMLElement, simUI: SimUI) {
		super(parent, 'settings-menu', { title: 'Options', footer: true });
		this.simUI = simUI;

		const restoreDefaultsButton = ref<HTMLButtonElement>();
		const fixedRngSeed = ref<HTMLDivElement>();
		const lastUsedRngSeed = ref<HTMLDivElement>();
		const language = ref<HTMLDivElement>();
		const showThreatMetrics = ref<HTMLDivElement>();
		const showExperimental = ref<HTMLDivElement>();
		const showQuickSwap = ref<HTMLDivElement>();

		const body = (
			<div>
				<div className="picker-group">
					<div className="fixed-rng-seed-container">
						<div ref={fixedRngSeed} className="fixed-rng-seed"></div>
						<div className="form-text">
							<span>Last used RNG seed:</span>&nbsp;
							<span ref={lastUsedRngSeed} className="last-used-rng-seed">
								0
							</span>
						</div>
					</div>
					<div ref={language} className="language-picker within-raid-sim-hide"></div>
				</div>
				<div ref={showThreatMetrics} className="show-threat-metrics-picker w-50 pe-2"></div>
				<div ref={showExperimental} className="show-experimental-picker w-50 pe-2"></div>
				<div ref={showQuickSwap} className="show-quick-swap-picker w-50 pe-2"></div>
			</div>
		);

		this.body.innerHTML = '';
		this.body.appendChild(body);

		const footer = (
			<button ref={restoreDefaultsButton} className="restore-defaults-button btn btn-primary">
				Restore Defaults
			</button>
		);
		if (this.footer) {
			this.footer.innerHTML = '';
			this.footer.appendChild(footer);
		}

		if (restoreDefaultsButton.value) {
			tippy(restoreDefaultsButton.value, {
				content: 'Restores all default settings (gear, consumes, buffs, talents, EP weights, etc). Saved settings are preserved.',
			});
			restoreDefaultsButton.value.addEventListener('click', () => {
				this.simUI.applyDefaults(TypedEvent.nextEventID());
				new Toast({
					variant: 'success',
					body: 'Restored to default settings.',
				});
			});
		}

		if (fixedRngSeed.value)
			new NumberPicker(fixedRngSeed.value, this.simUI.sim, {
				label: 'Fixed RNG Seed',
				labelTooltip:
					'Seed value for the random number generator used during sims, or 0 to use different randomness each run. Use this to share exact sim results or for debugging.',
				extraCssClasses: ['mb-0'],
				changedEvent: (sim: Sim) => sim.fixedRngSeedChangeEmitter,
				getValue: (sim: Sim) => sim.getFixedRngSeed(),
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					sim.setFixedRngSeed(eventID, newValue);
				},
			});

		if (lastUsedRngSeed.value) {
			lastUsedRngSeed.value.textContent = String(this.simUI.sim.getLastUsedRngSeed());
			this.simUI.sim.lastUsedRngSeedChangeEmitter.on(() => {
				if (lastUsedRngSeed.value) lastUsedRngSeed.value.textContent = String(this.simUI.sim.getLastUsedRngSeed());
			});
		}

		if (language.value) {
			const langs = Object.keys(wowheadSupportedLanguages);
			const defaultLang = langs.indexOf('en');
			const languagePicker = new EnumPicker(language.value, this.simUI.sim, {
				label: 'Language',
				labelTooltip: 'Controls the language for Wowhead tooltips.',
				values: langs.map((lang, i) => {
					return {
						name: wowheadSupportedLanguages[lang],
						value: i,
					};
				}),
				changedEvent: (sim: Sim) => sim.languageChangeEmitter,
				getValue: (sim: Sim) => {
					const idx = langs.indexOf(sim.getLanguage());
					return idx == -1 ? defaultLang : idx;
				},
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					sim.setLanguage(eventID, langs[newValue] || 'en');
				},
			});
			// Refresh page after language change, to apply the changes.
			languagePicker.changeEmitter.on(() => setTimeout(() => location.reload(), 100));
		}

		if (showThreatMetrics.value)
			new BooleanPicker(showThreatMetrics.value, this.simUI.sim, {
				label: 'Show Threat/Tank Options',
				labelTooltip: 'Shows all options and metrics relevant to tanks, like TPS/DTPS.',
				inline: true,
				changedEvent: (sim: Sim) => sim.showThreatMetricsChangeEmitter,
				getValue: (sim: Sim) => sim.getShowThreatMetrics(),
				setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
					sim.setShowThreatMetrics(eventID, newValue);
				},
			});

		if (showExperimental.value)
			new BooleanPicker(showExperimental.value, this.simUI.sim, {
				label: 'Show Experimental',
				labelTooltip: 'Shows experimental options, if there are any active experiments.',
				inline: true,
				changedEvent: (sim: Sim) => sim.showExperimentalChangeEmitter,
				getValue: (sim: Sim) => sim.getShowExperimental(),
				setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
					sim.setShowExperimental(eventID, newValue);
				},
			});
		if (showQuickSwap.value)
			new BooleanPicker(showQuickSwap.value, this.simUI.sim, {
				label: 'Show quick swap interface',
				labelTooltip: 'Allows you to quickly swap between Gems/Enchants through your favorites. (Disabled on touch devices)',
				inline: true,
				changedEvent: (sim: Sim) => sim.showQuickSwapChangeEmitter,
				getValue: (sim: Sim) => sim.getShowQuickSwap(),
				setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
					sim.setShowQuickSwap(eventID, newValue);
				},
			});
	}
}
