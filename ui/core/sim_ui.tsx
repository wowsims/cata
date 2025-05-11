import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { BaseModal } from './components/base_modal.jsx';
import { Component } from './components/component.js';
import { NoticeLocalSim } from './components/individual_sim_ui/notice_local_sim.jsx';
import { NumberPicker } from './components/pickers/number_picker.js';
import { ResultsViewer } from './components/results_viewer.jsx';
import { SimHeader } from './components/sim_header.jsx';
import { SimTab } from './components/sim_tab.js';
import { SimTitleDropdown } from './components/sim_title_dropdown.js';
import { SocialLinks } from './components/social_links.jsx';
import Toast from './components/toast';
import { REPO_NEW_ISSUE_URL } from './constants/other';
import { LaunchStatus, SimStatus } from './launched_sims.js';
import { PlayerSpec } from './player_spec.js';
import { ErrorOutcomeType } from './proto/api';
import { ActionId } from './proto_utils/action_id.js';
import { SimResult } from './proto_utils/sim_result';
import { Sim, SimError } from './sim.js';
import { RequestTypes } from './sim_signal_manager.js';
import { EventID, TypedEvent } from './typed_event.js';
import { WorkerProgressCallback } from './worker_pool';

const URLMAXLEN = 2048;
const globalKnownIssues: Array<string> = [];

// Config for displaying a warning to the user whenever a condition is met.
export interface SimWarning {
	updateOn: TypedEvent<any>;
	getContent: () => string | Array<string>;
}

export interface SimUIConfig {
	// Additional css class to add to the root element.
	cssClass: string;
	// Scheme used for themeing on a per-class Basis or for other sims
	cssScheme: string;
	// The spec, if an individual sim, or null if the raid sim.
	spec: PlayerSpec<any> | null;
	simStatus: SimStatus;
	knownIssues?: Array<string>;
	noticeText?: string;
}

// Shared UI for all individual sims and the raid sim.
export abstract class SimUI extends Component {
	readonly sim: Sim;
	readonly cssClass: string;
	readonly cssScheme: string;
	readonly isWithinRaidSim: boolean;

	// Emits when anything from the sim, raid, or encounter changes.
	readonly changeEmitter;

	readonly resultsViewer: ResultsViewer;
	readonly simHeader: SimHeader;

	readonly simContentContainer: HTMLElement;
	readonly simMain: HTMLElement;
	readonly simActionsContainer: HTMLElement;
	readonly iterationsPicker: HTMLElement;
	readonly simTabContentsContainer: HTMLElement;

	constructor(parentElem: HTMLElement, sim: Sim, config: SimUIConfig) {
		super(parentElem, 'sim-ui');
		this.sim = sim;
		this.cssClass = config.cssClass;
		this.cssScheme = config.cssScheme;
		this.isWithinRaidSim = this.rootElem.closest('.within-raid-sim') != null;

		const container = (
			<>
				<div className="sim-root">
					<div className="sim-bg"></div>
					{config.noticeText ? (
						<div className="notices-banner alert border-bottom mb-0 text-center within-raid-sim-hide">{config.noticeText}</div>
					) : null}
					<div className="sim-container">
						<aside className="sim-sidebar">
							<div className="sim-title"></div>
							<div className="sim-sidebar-content">
								<div className="sim-sidebar-actions within-raid-sim-hide"></div>
								<div className="sim-sidebar-results within-raid-sim-hide"></div>
								<div className="sim-sidebar-stats"></div>
								<div className="sim-sidebar-socials"></div>
							</div>
						</aside>
						<div className="sim-content container-fluid"></div>
					</div>
				</div>
				<div className="sim-toast-container p-3 bottom-0 right-0" id="toastContainer"></div>
			</>
		);

		this.rootElem.appendChild(container);

		this.simContentContainer = this.rootElem.querySelector('.sim-content') as HTMLElement;
		this.simHeader = new SimHeader(this.simContentContainer, this);
		this.simMain = document.createElement('main');
		this.simMain.classList.add('sim-main', 'tab-content');
		this.simContentContainer.appendChild(this.simMain);

		this.rootElem.classList.add(this.cssClass);

		if (!this.isWithinRaidSim) {
			this.rootElem.classList.add('not-within-raid-sim');
		}

		this.changeEmitter = TypedEvent.onAny([this.sim.changeEmitter], 'SimUIChange');

		this.sim.crashEmitter.on((eventID: EventID, error: SimError) => this.handleCrash(error));

		const updateShowDamageMetrics = () => {
			if (this.sim.getShowDamageMetrics()) this.rootElem.classList.remove('hide-damage-metrics');
			else this.rootElem.classList.add('hide-damage-metrics');
		};
		updateShowDamageMetrics();
		this.sim.showDamageMetricsChangeEmitter.on(updateShowDamageMetrics);

		const updateShowThreatMetrics = () => {
			if (this.sim.getShowThreatMetrics()) this.rootElem.classList.remove('hide-threat-metrics');
			else this.rootElem.classList.add('hide-threat-metrics');
		};
		updateShowThreatMetrics();
		this.sim.showThreatMetricsChangeEmitter.on(updateShowThreatMetrics);

		const updateShowHealingMetrics = () => {
			if (this.sim.getShowHealingMetrics()) this.rootElem.classList.remove('hide-healing-metrics');
			else this.rootElem.classList.add('hide-healing-metrics');
		};
		updateShowHealingMetrics();
		this.sim.showHealingMetricsChangeEmitter.on(updateShowHealingMetrics);

		const updateShowEpRatios = () => {
			// Threat metrics *always* shows multiple columns, so
			// always show ratios when they are shown
			if (this.sim.getShowThreatMetrics()) {
				this.rootElem.classList.remove('hide-ep-ratios');
				// This case doesn't currently happen, but who knows
				// what the future holds...
			} else if (this.sim.getShowDamageMetrics() && this.sim.getShowHealingMetrics()) {
				this.rootElem.classList.remove('hide-ep-ratios');
			} else {
				this.rootElem.classList.add('hide-ep-ratios');
			}
		};

		updateShowEpRatios();
		this.sim.showDamageMetricsChangeEmitter.on(updateShowEpRatios);
		this.sim.showHealingMetricsChangeEmitter.on(updateShowEpRatios);
		this.sim.showThreatMetricsChangeEmitter.on(updateShowEpRatios);

		const updateShowExperimental = () => {
			if (this.sim.getShowExperimental()) this.rootElem.classList.remove('hide-experimental');
			else this.rootElem.classList.add('hide-experimental');
		};
		updateShowExperimental();
		this.sim.showExperimentalChangeEmitter.on(updateShowExperimental);

		this.addKnownIssues(config);

		// Sidebar Contents

		const titleElem = this.rootElem.querySelector('.sim-title') as HTMLElement;
		new SimTitleDropdown(titleElem, config.spec, { noDropdown: this.isWithinRaidSim });

		this.simActionsContainer = this.rootElem.querySelector('.sim-sidebar-actions') as HTMLElement;
		this.addNoticeForLocalSim();

		this.iterationsPicker = new NumberPicker(this.simActionsContainer, this.sim, {
			id: 'simui-iterations',
			label: 'Iterations',
			extraCssClasses: ['iterations-picker', 'within-raid-sim-hide'],
			changedEvent: (sim: Sim) => sim.iterationsChangeEmitter,
			getValue: (sim: Sim) => sim.getIterations(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setIterations(eventID, newValue);
			},
		}).rootElem;

		const resultsViewerElem = this.rootElem.querySelector('.sim-sidebar-results') as HTMLElement;
		this.resultsViewer = new ResultsViewer(resultsViewerElem);

		const socialsContainer = this.rootElem.querySelector('.sim-sidebar-socials') as HTMLElement;
		socialsContainer.appendChild(SocialLinks.buildDiscordLink());
		socialsContainer.appendChild(SocialLinks.buildGitHubLink());
		socialsContainer.appendChild(SocialLinks.buildPatreonLink());

		this.simTabContentsContainer = this.rootElem.querySelector('.sim-main.tab-content') as HTMLElement;

		if (!this.isWithinRaidSim) {
			window.addEventListener('message', async event => {
				if (event.data == 'runOnce') {
					this.runSimOnce();
				}
			});
		}
	}

	addNoticeForLocalSim() {
		new NoticeLocalSim(this.simActionsContainer);
	}

	addAction(label: string, cssClass: string, onClick: (event: MouseEvent) => void): HTMLButtonElement {
		const buttonRef = ref<HTMLButtonElement>();
		this.simActionsContainer.appendChild(
			<button ref={buttonRef} className={clsx('sim-sidebar-action-button btn btn-primary w-100', cssClass)} onclick={onClick}>
				{label}
				<span className="sim-sidebar-action-button-loading-icon">
					<i className="fas fa-spinner fa-spin"></i>
				</span>
			</button>,
		);

		return buttonRef.value!;
	}

	addActionGroup(groups: ActionGroupItem[], groupOptions: { cssClass?: string } = {}) {
		const refs: HTMLButtonElement[] = [];
		const groupRef = ref<HTMLDivElement>();
		const { cssClass } = groupOptions;
		this.simActionsContainer.appendChild(
			<div ref={groupRef} className={clsx('d-flex btn-group w-100', cssClass)} attributes={{ role: 'group' }}>
				{groups.map(({ label, cssClass, children, onClick }) => (
					<button ref={ref => refs.push(ref)} onclick={onClick} className={clsx('sim-sidebar-action-button btn btn-primary', cssClass)}>
						{label}
						{children}
						<span className="sim-sidebar-action-button-loading-icon">
							<i className="fas fa-spinner fa-spin"></i>
						</span>
					</button>
				))}
			</div>,
		);

		return { group: groupRef.value!, children: refs };
	}

	addTab(title: string, cssClass: string, content: HTMLElement | Element) {
		const contentId = cssClass.replace(/\s+/g, '-') + '-tab';
		const isFirstTab = this.simTabContentsContainer.children.length == 0;

		this.simHeader.addTab(title, contentId);
		this.simTabContentsContainer.appendChild(
			<div id={contentId} className={clsx('tab-pane fade', isFirstTab && 'active show')}>
				{content}
			</div>,
		);
	}

	addSimTab(tab: SimTab) {
		this.simHeader.addSimTabLink(tab);
	}

	addWarning(warning: SimWarning) {
		this.resultsViewer.addWarning(warning);
	}

	private addKnownIssues(config: SimUIConfig) {
		let statusStr = '';
		switch (config.simStatus.status) {
			case LaunchStatus.Unlaunched:
				statusStr = 'This sim is a WORK IN PROGRESS. It is not fully developed and should not be used for general purposes.';
				break;
			case LaunchStatus.Alpha:
				statusStr = 'This sim is in ALPHA. Bugs are expected; please let us know if you find one!';
				break;
			case LaunchStatus.Beta:
				statusStr = 'This sim is in BETA. There may still be a few bugs; please let us know if you find one!';
				break;
		}

		if (statusStr) {
			config.knownIssues = [statusStr].concat(config.knownIssues || []);
		}
		config.knownIssues?.forEach(issue => this.simHeader.addKnownIssue(issue));
		globalKnownIssues?.forEach(issue => this.simHeader.addKnownIssue(issue));
	}

	// Returns a key suitable for the browser's localStorage feature.
	abstract getStorageKey(postfix: string): string;

	getSettingsStorageKey(): string {
		return this.getStorageKey('__currentSettings__');
	}

	getSavedEncounterStorageKey(): string {
		// By skipping the call to this.getStorageKey(), saved encounters will be
		// shared across all sims.
		return 'sharedData__savedEncounter__';
	}

	isIndividualSim(): boolean {
		return this.rootElem.classList.contains('individual-sim-ui');
	}

	async runSim(onProgress: WorkerProgressCallback) {
		this.resultsViewer.setPending();
		try {
			await this.sim.signalManager.abortType(RequestTypes.All);
			const result = await this.sim.runRaidSim(TypedEvent.nextEventID(), onProgress);
			if (!(result instanceof SimResult) && result.type == ErrorOutcomeType.ErrorOutcomeAborted) {
				new Toast({
					variant: 'info',
					body: 'Raid sim cancelled.',
				});
			}
		} catch (e) {
			this.resultsViewer.hideAll();
			this.handleCrash(e);
		}
	}

	async runSimOnce() {
		this.resultsViewer.setPending();
		try {
			await this.sim.runRaidSimWithLogs(TypedEvent.nextEventID());
		} catch (e) {
			this.resultsViewer.hideAll();
			this.handleCrash(e);
		}
	}

	async handleCrash(error: any): Promise<void> {
		if (!(error instanceof SimError)) {
			if (error.message) {
				new Toast({
					variant: 'error',
					body: error.message,
				});
			} else {
				alert(error);
			}
			return;
		}

		new Toast({
			variant: 'error',
			body: 'The Simulation failed. Generated an error report.',
		});

		const errorStr = (error as SimError).errorStr;
		if (errorStr.startsWith('[USER_ERROR] ')) {
			let alertStr = errorStr.substring('[USER_ERROR] '.length);
			alertStr = await ActionId.replaceAllInString(alertStr);
			alert(alertStr);
			return;
		}

		if (window.confirm('Simulation Failure:\n' + errorStr + '\nPress Ok to file crash report')) {
			// Splice out just the line numbers
			const hash = this.hashCode(errorStr);
			const link = this.toLink();
			const rngSeed = this.sim.getLastUsedRngSeed();
			fetch('https://api.github.com/search/issues?q=is:issue+is:open+repo:wowsims/mop+' + hash)
				.then(resp => {
					resp.json().then(issues => {
						if (issues.total_count > 0) {
							window.open(issues.items[0].html_url, '_blank');
						} else {
							const url = new URL(REPO_NEW_ISSUE_URL);
							url.searchParams.append('title', `Crash Report ${hash}`);
							url.searchParams.append('assignees', '');
							url.searchParams.append('labels', '');

							const maxBodyLength = URLMAXLEN - url.toString().length;
							let issueBody = `Link:\n${link}\n\nRNG Seed: ${rngSeed}\n\n${errorStr}`;
							if (link.includes('/raid/')) {
								// Move the actual error before the link, as it will likely get truncated.
								issueBody = `${errorStr}\nRNG Seed: ${rngSeed}\nLink:\n${link}`;
							}
							let truncated = false;
							while (issueBody.length > maxBodyLength - (truncated ? 3 : 0)) {
								issueBody = issueBody.slice(0, issueBody.lastIndexOf('%')); // Avoid truncating in the middle of a URLencoded segment.
								truncated = true;
							}
							if (truncated) {
								issueBody += '...';
								// The raid links are too large and will always cause truncation.
								// Prompt the user to add more information to the issue.
								new CrashModal(this.rootElem, link).open();
							}
							url.searchParams.append('body', issueBody);

							window.open(url.toString(), '_blank');
						}
					});
				})
				.catch(fetchErr => {
					alert('Failed to file report... try again another time:' + fetchErr);
				});
		}
	}

	hashCode(str: string): number {
		let hash = 0;
		for (let i = 0, len = str.length; i < len; i++) {
			const chr = str.charCodeAt(i);
			hash = (hash << 5) - hash + chr;
			hash |= 0; // Convert to 32bit integer
		}
		return hash;
	}

	abstract applyDefaults(eventID: EventID): void;
	abstract toLink(): string;
}

class CrashModal extends BaseModal {
	constructor(parent: HTMLElement, link: string) {
		super(parent, 'crash', { title: 'Extra Crash Information' });
		this.body.appendChild(
			<div className="sim-crash-report">
				<h3 className="sim-crash-report-header">
					Please append the following complete link to the issue you just created. This will simplify debugging the issue.
				</h3>
				<textarea className="sim-crash-report-text form-control">{link}</textarea>
			</div>,
		);
	}
}

export type ActionGroupItem = { label?: string; children?: Element; cssClass?: string; onClick?: (event: MouseEvent) => void };
