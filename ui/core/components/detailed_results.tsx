import { REPO_NAME } from '../constants/other';
import { IndividualSimUI } from '../individual_sim_ui';
import { DetailedResultsUpdate, SimRun, SimRunData } from '../proto/ui';
import { SimResult } from '../proto_utils/sim_result';
import { SimUI } from '../sim_ui';
import { TypedEvent } from '../typed_event';
import { Component } from './component';
import { AuraMetricsTable } from './detailed_results/aura_metrics';
import { CastMetricsTable } from './detailed_results/cast_metrics';
import { DamageMetricsTable } from './detailed_results/damage_metrics';
import { DpsHistogram } from './detailed_results/dps_histogram';
import { DtpsMetricsTable } from './detailed_results/dtps_metrics';
import { HealingMetricsTable } from './detailed_results/healing_metrics';
import { LogRunner } from './detailed_results/log_runner';
import { PlayerDamageMetricsTable } from './detailed_results/player_damage';
import { PlayerDamageTakenMetricsTable } from './detailed_results/player_damage_taken';
import { ResourceMetricsTable } from './detailed_results/resource_metrics';
import { SimResultData } from './detailed_results/result_component';
import { ResultsFilter } from './detailed_results/results_filter';
import { Timeline } from './detailed_results/timeline';
import { ToplineResults } from './detailed_results/topline_results';
import { RaidSimResultsManager } from './raid_sim_action';

type Tab = {
	isActive?: boolean;
	targetId: string;
	label: string;
	classes?: string[];
};

const tabs: Tab[] = [
	{
		isActive: true,
		targetId: 'damageTab',
		label: 'Damage',
		classes: ['damage-metrics-tab'],
	},
	{
		targetId: 'healingTab',
		label: 'Healing',
		classes: ['healing-metrics-tab'],
	},
	{
		targetId: 'damageTakenTab',
		label: 'Damage Taken',
		classes: ['threat-metrics-tab'],
	},
	{
		targetId: 'buffsTab',
		label: 'Buffs',
	},
	{
		targetId: 'debuffsTab',
		label: 'Debuffs',
	},
	{
		targetId: 'castsTab',
		label: 'Casts',
	},
	{
		targetId: 'resourcesTab',
		label: 'Resources',
	},
	{
		targetId: 'timelineTab',
		label: 'Timeline',
	},
	{
		targetId: 'logTab',
		label: 'Log',
	},
];

export abstract class DetailedResults extends Component {
	protected readonly simUI: SimUI | null;
	protected latestRun: SimRunData | null = null;

	private currentSimResult: SimResult | null = null;
	private resultsEmitter: TypedEvent<SimResultData | null> = new TypedEvent<SimResultData | null>();
	private resultsFilter: ResultsFilter;
	private rootDiv: HTMLElement;

	constructor(parent: HTMLElement, simUI: SimUI | null, cssScheme: string) {
		super(parent, 'detailed-results-manager-root');

		this.rootElem.appendChild(
			<div className="dr-root dr-no-results">
				<div className="dr-toolbar">
					<div className="results-filter"></div>
					<div className="tabs-filler"></div>
					<ul className="nav nav-tabs" attributes={{ role: 'tablist' }}>
						{tabs.map(({ label, targetId, isActive, classes }) => (
							<li className={`nav-item dr-tab-tab ${classes?.join(' ') || ''}`} attributes={{ role: 'presentation' }}>
								<button
									className={`nav-link${isActive ? ' active' : ''}`}
									type="button"
									attributes={{
										role: 'tab',
										// @ts-expect-error
										'aria-controls': targetId,
										'aria-selected': !!isActive,
									}}
									dataset={{
										bsToggle: 'tab',
										bsTarget: `#${targetId}`,
									}}>
									{label}
								</button>
							</li>
						))}
					</ul>
				</div>
				<div className="tab-content">
					<div id="noResultsTab" className="tab-pane dr-tab-content fade active show">
						Run a simulation to view results
					</div>
					<div id="damageTab" className="tab-pane dr-tab-content damage-content fade active show">
						<div className="dr-row topline-results" />
						<div className="dr-row all-players-only">
							<div className="player-damage-metrics" />
						</div>
						<div className="dr-row single-player-only">
							<div className="damage-metrics" />
						</div>
						{/* <div className="dr-row single-player-only">
							<div className="melee-metrics" />
						</div>
						<div className="dr-row single-player-only">
							<div className="spell-metrics" />
						</div> */}
						<div className="dr-row dps-histogram" />
					</div>
					<div id="healingTab" className="tab-pane dr-tab-content healing-content fade">
						<div className="dr-row topline-results" />
						<div className="dr-row single-player-only">
							<div className="healing-spell-metrics" />
						</div>
						<div className="dr-row hps-histogram" />
					</div>
					<div id="damageTakenTab" className="tab-pane dr-tab-content damage-taken-content fade">
						<div className="dr-row topline-results" />
						<div className="dr-row all-players-only">
							<div className="player-damage-taken-metrics" />
						</div>
						<div className="dr-row single-player-only">
							<div className="dtps-metrics" />
						</div>
						<div className="dr-row damage-taken-histogram single-player-only" />
					</div>
					<div id="buffsTab" className="tab-pane dr-tab-content buffs-content fade">
						<div className="dr-row">
							<div className="buff-aura-metrics" />
						</div>
					</div>
					<div id="debuffsTab" className="tab-pane dr-tab-content debuffs-content fade">
						<div className="dr-row">
							<div className="debuff-aura-metrics" />
						</div>
					</div>
					<div id="castsTab" className="tab-pane dr-tab-content casts-content fade">
						<div className="dr-row">
							<div className="cast-metrics" />
						</div>
					</div>
					<div id="resourcesTab" className="tab-pane dr-tab-content resources-content fade">
						<div className="dr-row">
							<div className="resource-metrics" />
						</div>
					</div>
					<div id="timelineTab" className="tab-pane dr-tab-content timeline-content fade">
						<div className="dr-row">
							<div className="timeline" />
						</div>
					</div>
					<div id="logTab" className="tab-pane dr-tab-content log-content fade">
						<div className="dr-row">
							<div className="log" />
						</div>
					</div>
				</div>
			</div>,
		);
		this.rootDiv = this.rootElem.querySelector('.dr-root')!;
		this.simUI = simUI;

		this.simUI?.sim.settingsChangeEmitter.on(async () => await this.updateSettings());

		// Allow styling the sticky toolbar
		const toolbar = document.querySelector<HTMLElement>('.dr-toolbar')!;
		new IntersectionObserver(
			([e]) => {
				e.target.classList.toggle('stuck', e.intersectionRatio < 1);
			},
			{
				// Intersect with the sim header or top of the separate tab
				rootMargin: this.simUI ? `-${this.simUI.simHeader.rootElem.offsetHeight + 1}px 0px 0px 0px` : '0px',
				threshold: [1],
			},
		).observe(toolbar);

		this.resultsFilter = new ResultsFilter({
			parent: this.rootElem.querySelector('.results-filter')!,
			resultsEmitter: this.resultsEmitter,
		});

		[...this.rootElem.querySelectorAll<HTMLElement>('.topline-results')]?.forEach(toplineResultsDiv => {
			new ToplineResults({ parent: toplineResultsDiv, resultsEmitter: this.resultsEmitter });
		});

		new CastMetricsTable({
			parent: this.rootElem.querySelector('.cast-metrics')!,
			resultsEmitter: this.resultsEmitter,
		});
		new DamageMetricsTable({
			parent: this.rootElem.querySelector('.damage-metrics')!,
			resultsEmitter: this.resultsEmitter,
		});

		new HealingMetricsTable({
			parent: this.rootElem.querySelector('.healing-spell-metrics')!,
			resultsEmitter: this.resultsEmitter,
		});
		new ResourceMetricsTable({
			parent: this.rootElem.querySelector('.resource-metrics')!,
			resultsEmitter: this.resultsEmitter,
			secondaryResource: (simUI as IndividualSimUI<any>)?.player?.secondaryResource,
		});
		new PlayerDamageMetricsTable(
			{ parent: this.rootElem.querySelector('.player-damage-metrics')!, resultsEmitter: this.resultsEmitter },
			this.resultsFilter,
		);
		new PlayerDamageTakenMetricsTable(
			{ parent: this.rootElem.querySelector('.player-damage-taken-metrics')!, resultsEmitter: this.resultsEmitter },
			this.resultsFilter,
		);
		new AuraMetricsTable(
			{
				parent: this.rootElem.querySelector('.buff-aura-metrics')!,
				resultsEmitter: this.resultsEmitter,
			},
			false,
		);
		new AuraMetricsTable(
			{
				parent: this.rootElem.querySelector('.debuff-aura-metrics')!,
				resultsEmitter: this.resultsEmitter,
			},
			true,
		);

		new DpsHistogram({
			parent: this.rootElem.querySelector('.dps-histogram')!,
			resultsEmitter: this.resultsEmitter,
		});

		new DtpsMetricsTable({
			parent: this.rootElem.querySelector('.dtps-metrics')!,
			resultsEmitter: this.resultsEmitter,
		});

		const timeline = new Timeline({
			parent: this.rootElem.querySelector('.timeline')!,
			cssScheme: cssScheme,
			resultsEmitter: this.resultsEmitter,
			secondaryResource: (simUI as IndividualSimUI<any>)?.player?.secondaryResource,
		});

		const tabEl = document.querySelector('button[data-bs-target="#timelineTab"]');
		tabEl?.addEventListener('shown.bs.tab', () => {
			timeline.render();
		});

		new LogRunner({
			parent: this.rootElem.querySelector('.log')!,
			cssScheme: cssScheme,
			resultsEmitter: this.resultsEmitter,
		});

		this.rootElem.classList.add('hide-threat-metrics', 'hide-threat-metrics');

		this.resultsFilter.changeEmitter.on(() => this.updateResults());

		this.resultsEmitter.on((_, resultData) => {
			if (resultData?.filter.player || resultData?.filter.player === 0) {
				this.rootDiv.classList.remove('all-players');
				this.rootDiv.classList.add('single-player');
			} else {
				this.rootDiv.classList.add('all-players');
				this.rootDiv.classList.remove('single-player');
			}
		});
	}

	abstract postMessage(update: DetailedResultsUpdate): Promise<void>;

	protected async setSimRunData(simRunData: SimRunData) {
		this.latestRun = simRunData;
		await this.postMessage(
			DetailedResultsUpdate.create({
				data: {
					oneofKind: 'runData',
					runData: simRunData,
				},
			}),
		);
	}

	protected async updateSettings() {
		if (!this.simUI) return;
		await this.postMessage(
			DetailedResultsUpdate.create({
				data: {
					oneofKind: 'settings',
					settings: this.simUI.sim.toProto(),
				},
			}),
		);
	}

	private updateResults() {
		const eventID = TypedEvent.nextEventID();
		if (this.currentSimResult == null) {
			this.rootDiv.classList.add('dr-no-results');
			this.resultsEmitter.emit(eventID, null);
		} else {
			this.rootDiv.classList.remove('dr-no-results');
			this.resultsEmitter.emit(eventID, {
				eventID: eventID,
				result: this.currentSimResult,
				filter: this.resultsFilter.getFilter(),
			});
		}
	}

	protected async handleMessage(data: DetailedResultsUpdate) {
		switch (data.data.oneofKind) {
			case 'runData':
				const runData = data.data.runData;
				this.currentSimResult = await SimResult.fromProto(runData.run || SimRun.create());
				this.updateResults();
				break;
			case 'settings':
				const settings = data.data.settings;
				if (settings.showDamageMetrics) {
					this.rootElem.classList.remove('hide-damage-metrics');
				} else {
					this.rootElem.classList.add('hide-damage-metrics');
					const damageTabEl = document.getElementById('damageTab')!;
					const healingTabEl = document.getElementById('healingTab')!;
					if (damageTabEl.classList.contains('active')) {
						damageTabEl.classList.remove('active', 'show');
						healingTabEl.classList.add('active', 'show');

						const toolbar = document.getElementsByClassName('dr-toolbar')[0] as HTMLElement;
						toolbar.querySelector('.damage-metrics')?.children[0].classList.remove('active');
						toolbar.querySelector('.healing-metrics')?.children[0].classList.add('active');
					}
				}
				this.rootElem.classList[settings.showThreatMetrics ? 'remove' : 'add']('hide-threat-metrics');
				this.rootElem.classList[settings.showHealingMetrics ? 'remove' : 'add']('hide-healing-metrics');
				this.rootElem.classList[settings.showExperimental ? 'remove' : 'add']('hide-experimental');
				break;
		}
	}
}

export class WindowedDetailedResults extends DetailedResults {
	constructor(parent: HTMLElement) {
		super(parent, null, new URLSearchParams(window.location.search).get('cssScheme') ?? '');

		window.addEventListener('message', async event => await this.handleMessage(DetailedResultsUpdate.fromJson(event.data)));

		this.rootElem.insertAdjacentHTML('beforeend', `<div class="sim-bg"></div>`);
	}

	async postMessage(update: DetailedResultsUpdate): Promise<void> {
		await this.handleMessage(update);
	}
}

export class EmbeddedDetailedResults extends DetailedResults {
	private tabWindow: Window | null = null;

	constructor(parent: HTMLElement, simUI: SimUI, simResultsManager: RaidSimResultsManager) {
		super(parent, simUI, simUI.cssScheme);

		const newTabBtn = (
			<div className="detailed-results-controls-div">
				<button className="detailed-results-new-tab-button btn btn-primary">View in Separate Tab</button>
				<button className="detailed-results-1-iteration-button btn btn-primary">Sim 1 Iteration</button>
			</div>
		);

		this.rootElem.prepend(newTabBtn);

		const url = new URL(`${window.location.protocol}//${window.location.host}/${REPO_NAME}/detailed_results/index.html`);
		url.searchParams.append('cssClass', simUI.cssClass);

		if (simUI.isIndividualSim()) {
			url.searchParams.append('isIndividualSim', '');
			this.rootElem.classList.add('individual-sim');
		}

		const newTabButton = this.rootElem.querySelector('.detailed-results-new-tab-button');
		newTabButton?.addEventListener('click', () => {
			if (this.tabWindow == null || this.tabWindow.closed) {
				this.tabWindow = window.open(url.href, 'Detailed Results');
				this.tabWindow!.addEventListener('load', async () => {
					if (this.latestRun) {
						await Promise.all([this.updateSettings(), this.setSimRunData(this.latestRun)]);
					}
				});
			} else {
				this.tabWindow.focus();
			}
		});

		const simButton = this.rootElem.querySelector('.detailed-results-1-iteration-button');
		simButton?.addEventListener('click', () => {
			(window.opener || window.parent)!.postMessage('runOnce', '*');
		});

		simResultsManager.currentChangeEmitter.on(async () => {
			const runData = simResultsManager.getRunData();
			if (runData) {
				await Promise.all([this.updateSettings(), this.setSimRunData(runData)]);
			}
		});
	}

	async postMessage(update: DetailedResultsUpdate) {
		if (this.tabWindow) {
			this.tabWindow.postMessage(DetailedResultsUpdate.toJson(update), '*');
		}
		await this.handleMessage(update);
	}
}
