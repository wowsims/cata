import { SimRun, SimRunData } from '../proto/ui';
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

declare let Chart: any;

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

export class DetailedResults extends Component {
	protected readonly simUI: SimUI;
	protected latestRun: SimRunData | null = null;

	private currentSimResult: SimResult | null = null;
	private resultsEmitter: TypedEvent<SimResultData | null> = new TypedEvent<SimResultData | null>();
	private resultsFilter: ResultsFilter;
	private rootDiv: Element;

	constructor(parent: HTMLElement, simUI: SimUI, simResultsManager: RaidSimResultsManager) {
		super(parent, 'detailed-results-manager-root');

		const simOnceBtn = <button className="detailed-results-1-iteration-button btn btn-primary">Sim 1 Iteration</button>;
		this.rootDiv = (
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
			</div>
		);

		this.rootElem.appendChild(
			<>
				<div className="detailed-results-controls-div">{simOnceBtn}</div>
				{this.rootDiv}
			</>,
		);

		this.simUI = simUI;
		this.simUI.sim.settingsChangeEmitter.on(() => this.updateSettings());

		Chart.defaults.color = 'white';

		// Allow styling the sticky toolbar
		const toolbar = document.querySelector('.dr-toolbar') as HTMLElement;
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
			resultsEmitter: this.resultsEmitter,
		});

		const tabEl = document.querySelector('button[data-bs-target="#timelineTab"]');
		tabEl?.addEventListener('shown.bs.tab', () => {
			timeline.render();
		});

		new LogRunner({
			parent: this.rootElem.querySelector('.log')!,
			resultsEmitter: this.resultsEmitter,
		});

		this.rootElem.classList.add('hide-threat-metrics');

		this.resultsFilter.changeEmitter.on(async () => await this.updateResults(this.latestRun));

		this.resultsEmitter.on((_, resultData) => {
			if (resultData?.filter.player || resultData?.filter.player === 0) {
				this.rootDiv.classList.remove('all-players');
				this.rootDiv.classList.add('single-player');
			} else {
				this.rootDiv.classList.add('all-players');
				this.rootDiv.classList.remove('single-player');
			}
		});

		simOnceBtn.addEventListener('click', () => this.simUI.runSimOnce());

		simResultsManager.currentChangeEmitter.on(async () => {
			const runData = simResultsManager.getRunData();
			if (runData) {
				this.updateSettings();
				await this.updateResults(runData);
			}
		});
	}

	private updateSettings() {
		const settings = this.simUI.sim.toProto();
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
	}

	private async updateResults(simRunData: SimRunData | null) {
		this.latestRun = simRunData;
		this.currentSimResult = await SimResult.fromProto(simRunData?.run || SimRun.create());

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
}
