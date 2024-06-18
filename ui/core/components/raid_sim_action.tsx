import tippy from 'tippy.js';

import { DistributionMetrics as DistributionMetricsProto, ProgressMetrics, Raid as RaidProto } from '../proto/api.js';
import { Encounter as EncounterProto, Spec } from '../proto/common.js';
import { SimRunData } from '../proto/ui.js';
import { ActionMetrics, SimResult, SimResultFilter } from '../proto_utils/sim_result.js';
import { RequestTypes } from '../sim_signal_manager';
import { SimUI } from '../sim_ui.jsx';
import { EventID, TypedEvent } from '../typed_event.js';
import { formatDeltaTextElem, sum } from '../utils.js';

export function addRaidSimAction(simUI: SimUI): RaidSimResultsManager {
	const resultsViewer = simUI.resultsViewer
	let isRunning = false;
	let waitAbort = false;
	simUI.addAction('Simulate', 'dps-action', async ev => {
		const button = ev.target as HTMLButtonElement;
		button.disabled = true;
		if (!isRunning) {
			isRunning = true;

			resultsViewer.addAbortButton(async () => {
				if (waitAbort) return;
				try {
					waitAbort = true;
					await simUI.sim.signalManager.abortType(RequestTypes.RaidSim);
				} catch (error) {
					console.error('Error on sim abort!');
					console.error(error);
				} finally {
					waitAbort = false;
					if (!isRunning) button.disabled = false;
				}
			});

			await simUI.runSim((progress: ProgressMetrics) => {
				resultsManager.setSimProgress(progress);
			});

			resultsViewer.removeAbortButton();
			if (!waitAbort) button.disabled = false;
			isRunning = false;
		}
	});

	const resultsManager = new RaidSimResultsManager(simUI);
	simUI.sim.simResultEmitter.on((eventID, simResult) => {
		resultsManager.setSimResult(eventID, simResult);
	});
	return resultsManager;
}

export type ReferenceData = {
	simResult: SimResult;
	settings: any;
	raidProto: RaidProto;
	encounterProto: EncounterProto;
};

export interface ResultMetrics {
	cod: string;
	dps: string;
	dpasp: string;
	dtps: string;
	tmi: string;
	dur: string;
	hps: string;
	tps: string;
	tto: string;
}

export interface ResultMetricCategories {
	damage: string;
	demo: string;
	healing: string;
	threat: string;
}

export interface ResultsLineArgs {
	average: number;
	stdev?: number;
	classes?: string;
}

export class RaidSimResultsManager {
	static resultMetricCategories: { [ResultMetrics: string]: keyof ResultMetricCategories } = {
		dps: 'damage',
		dpasp: 'demo',
		tps: 'threat',
		dtps: 'threat',
		tmi: 'threat',
		cod: 'threat',
		tto: 'healing',
		hps: 'healing',
	};

	static resultMetricClasses: { [ResultMetrics: string]: string } = {
		cod: 'results-sim-cod',
		dps: 'results-sim-dps',
		dpasp: 'results-sim-dpasp',
		dtps: 'results-sim-dtps',
		tmi: 'results-sim-tmi',
		dur: 'results-sim-dur',
		hps: 'results-sim-hps',
		tps: 'results-sim-tps',
		tto: 'results-sim-tto',
	};

	static metricsClasses: { [ResultMetricCategories: string]: string } = {
		damage: 'damage-metrics',
		demo: 'demo-metrics',
		healing: 'healing-metrics',
		threat: 'threat-metrics',
	};

	readonly currentChangeEmitter: TypedEvent<void> = new TypedEvent<void>();
	readonly referenceChangeEmitter: TypedEvent<void> = new TypedEvent<void>();

	readonly changeEmitter: TypedEvent<void> = new TypedEvent<void>();

	private readonly simUI: SimUI;

	private currentData: ReferenceData | null = null;
	private referenceData: ReferenceData | null = null;

	constructor(simUI: SimUI) {
		this.simUI = simUI;

		[this.currentChangeEmitter, this.referenceChangeEmitter].forEach(emitter => emitter.on(eventID => this.changeEmitter.emit(eventID)));
	}

	setSimProgress(progress: ProgressMetrics) {
		if (progress.finalRaidResult && progress.finalRaidResult.errorResult) {
			this.simUI.resultsViewer.hideAll();
			return;
		}
		this.simUI.resultsViewer.setContent(
			<div className="results-sim">
				<div className="results-sim-dps damage-metrics">
					<span className="topline-result-avg">{progress.dps.toFixed(2)}</span>
				</div>
				{this.simUI.isIndividualSim() && (
					<div className="results-sim-hps healing-metrics">
						<span className="topline-result-avg">{progress.hps.toFixed(2)}</span>
					</div>
				)}
				<div>
					{progress.presimRunning ? 'presimulations running' : `${progress.completedIterations} / ${progress.totalIterations}`}
					<br />
					iterations complete
				</div>
			</div>,
		);
	}

	setSimResult(eventID: EventID, simResult: SimResult) {
		this.currentData = {
			simResult: simResult,
			settings: {
				raid: RaidProto.toJson(this.simUI.sim.raid.toProto()),
				encounter: EncounterProto.toJson(this.simUI.sim.encounter.toProto()),
			},
			raidProto: RaidProto.clone(simResult.request.raid || RaidProto.create()),
			encounterProto: EncounterProto.clone(simResult.request.encounter || EncounterProto.create()),
		};
		this.currentChangeEmitter.emit(eventID);

		this.simUI.resultsViewer.setContent(
			<div className="results-sim">
				{RaidSimResultsManager.makeToplineResultsContent(simResult)}
				<div className="results-sim-reference">
					<button className="results-sim-set-reference">
						<i className={`fa fa-map-pin fa-lg text-${this.simUI.cssScheme} me-2`} />
						Save as Reference
					</button>
					<div className="results-sim-reference-bar">
						<button className="results-sim-reference-swap me-3">
							<i className="fas fa-arrows-rotate me-1" />
							Swap
						</button>
						<button className="results-sim-reference-delete">
							<i className="fa fa-times fa-lg me-1" />
							Cancel
						</button>
					</div>
				</div>
			</div>,
		);

		const setResultTooltip = (selector: string, content: Element | HTMLElement | string) => {
			const resultDivElem = this.simUI.resultsViewer.contentElem.querySelector<HTMLElement>(selector);
			if (resultDivElem) {
				tippy(resultDivElem, { content, placement: 'right' });
			}
		};
		setResultTooltip('.results-sim-dps', 'Damage Per Second');
		setResultTooltip('.results-sim-dpasp', 'Demonic Pact Average Spell Power');
		setResultTooltip('.results-sim-tto', 'Time To OOM');
		setResultTooltip('.results-sim-hps', 'Healing+Shielding Per Second, including overhealing.');
		setResultTooltip('.results-sim-tps', 'Threat Per Second');
		setResultTooltip('.results-sim-dtps', 'Damage Taken Per Second');
		setResultTooltip('.results-sim-dur', 'Average Fight Duration');
		setResultTooltip(
			'.results-sim-tmi',
			<>
				<p>Theck-Meloree Index (TMI)</p>
				<p>A measure of incoming damage smoothness which combines the benefits of avoidance with effective health.</p>
				<p>
					<b>Lower is better.</b> This represents the % of your HP to expect in a 6-second burst window based on the encounter settings.
				</p>
			</>,
		);
		setResultTooltip(
			'.results-sim-cod',
			<>
				<p>Chance of Death</p>
				<p>
					The percentage of iterations in which the player died, based on incoming damage from the enemies and incoming healing (see the{' '}
					<b>Incoming HPS</b> and <b>Healing Cadence</b> options).
				</p>
				<p>
					DTPS alone is not a good measure of tankiness because it is not affected by health and ignores damage spikes. Chance of Death attempts to
					capture overall tankiness.
				</p>
			</>,
		);

		if (!this.simUI.isIndividualSim()) {
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-reference-diff-separator')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-dpasp')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-tto')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-hps')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-tps')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-dtps')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-tmi')].forEach(e => e.remove());
			[...this.simUI.resultsViewer.contentElem.querySelectorAll('results-sim-cod')].forEach(e => e.remove());
		}

		const simReferenceSetButton = this.simUI.resultsViewer.contentElem.querySelector<HTMLSpanElement>('.results-sim-set-reference');
		if (simReferenceSetButton) {
			simReferenceSetButton.addEventListener('click', () => {
				this.referenceData = this.currentData;
				this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
				this.updateReference();
			});
			tippy(simReferenceSetButton, { content: 'Use as reference' });
		}

		const simReferenceSwapButton = this.simUI.resultsViewer.contentElem.querySelector<HTMLSpanElement>('.results-sim-reference-swap');
		if (simReferenceSwapButton) {
			simReferenceSwapButton.addEventListener('click', () => {
				TypedEvent.freezeAllAndDo(() => {
					if (this.currentData && this.referenceData) {
						const swapEventID = TypedEvent.nextEventID();
						const tmpData = this.currentData;
						this.currentData = this.referenceData;
						this.referenceData = tmpData;

						this.simUI.sim.raid.fromProto(swapEventID, this.currentData.raidProto);
						this.simUI.sim.encounter.fromProto(swapEventID, this.currentData.encounterProto);
						this.setSimResult(swapEventID, this.currentData.simResult);

						this.referenceChangeEmitter.emit(swapEventID);
						this.updateReference();
					}
				});
			});
			tippy(simReferenceSwapButton, {
				content: 'Swap reference with current',
				ignoreAttributes: true,
			});
		}
		const simReferenceDeleteButton = this.simUI.resultsViewer.contentElem.querySelector<HTMLSpanElement>('.results-sim-reference-delete');
		if (simReferenceDeleteButton) {
			simReferenceDeleteButton.addEventListener('click', () => {
				this.referenceData = null;
				this.referenceChangeEmitter.emit(TypedEvent.nextEventID());
				this.updateReference();
			});
			tippy(simReferenceDeleteButton, {
				content: 'Remove reference',
				ignoreAttributes: true,
			});
		}

		this.updateReference();
	}

	private updateReference() {
		if (!this.referenceData || !this.currentData) {
			// Remove references
			this.simUI.resultsViewer.contentElem.querySelector('.results-sim-reference')?.classList.remove('has-reference');
			this.simUI.resultsViewer.contentElem.querySelectorAll('.results-reference').forEach(e => e.classList.add('hide'));
			return;
		} else {
			// Add references references
			this.simUI.resultsViewer.contentElem.querySelector('.results-sim-reference')?.classList.add('has-reference');
			this.simUI.resultsViewer.contentElem.querySelectorAll('.results-reference').forEach(e => e.classList.remove('hide'));
		}

		this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['dps']} .results-reference-diff`, res => res.raidMetrics.dps, 2);
		if (this.simUI.isIndividualSim()) {
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['hps']} .results-reference-diff`, res => res.raidMetrics.hps, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['dpasp']} .results-reference-diff`, res => res.getFirstPlayer()!.dpasp, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tto']} .results-reference-diff`, res => res.getFirstPlayer()!.tto, 2);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tps']} .results-reference-diff`, res => res.getFirstPlayer()!.tps, 2);
			this.formatToplineResult(
				`.${RaidSimResultsManager.resultMetricClasses['dtps']} .results-reference-diff`,
				res => res.getFirstPlayer()!.dtps,
				2,
				true,
			);
			this.formatToplineResult(`.${RaidSimResultsManager.resultMetricClasses['tmi']} .results-reference-diff`, res => res.getFirstPlayer()!.tmi, 2, true);
			this.formatToplineResult(
				`.${RaidSimResultsManager.resultMetricClasses['cod']} .results-reference-diff`,
				res => res.getFirstPlayer()!.chanceOfDeath,
				1,
				true,
			);
		} else {
			this.formatToplineResult(
				`.${RaidSimResultsManager.resultMetricClasses['dtps']} .results-reference-diff`,
				res => sum(res.getPlayers()!.map(player => player.dtps.avg)) / res.getPlayers().length,
				2,
				true,
			);
		}
	}

	private formatToplineResult(
		querySelector: string,
		getMetrics: (result: SimResult) => DistributionMetricsProto | number,
		precision: number,
		lowerIsBetter?: boolean,
	) {
		const elem = this.simUI.resultsViewer.contentElem.querySelector<HTMLSpanElement>(querySelector);
		if (!elem) {
			return;
		}

		const cur = this.currentData!.simResult;
		const ref = this.referenceData!.simResult;
		const curMetricsTemp = getMetrics(cur);
		const refMetricsTemp = getMetrics(ref);
		if (typeof curMetricsTemp === 'number') {
			const curMetrics = curMetricsTemp as number;
			const refMetrics = refMetricsTemp as number;
			formatDeltaTextElem(elem, refMetrics, curMetrics, precision, lowerIsBetter);
		} else {
			const curMetrics = curMetricsTemp as DistributionMetricsProto;
			const refMetrics = refMetricsTemp as DistributionMetricsProto;
			const isDiff = this.applyZTestTooltip(elem, ref.iterations, refMetrics.avg, refMetrics.stdev, cur.iterations, curMetrics.avg, curMetrics.stdev);
			formatDeltaTextElem(elem, refMetrics.avg, curMetrics.avg, precision, lowerIsBetter, !isDiff);
		}
	}

	private applyZTestTooltip(elem: HTMLElement, n1: number, avg1: number, stdev1: number, n2: number, avg2: number, stdev2: number): boolean {
		const delta = avg1 - avg2;
		const err1 = stdev1 / Math.sqrt(n1);
		const err2 = stdev2 / Math.sqrt(n2);
		const denom = Math.sqrt(Math.pow(err1, 2) + Math.pow(err2, 2));
		const z = Math.abs(delta / denom);
		const isDiff = z > 1.96;

		let significance_str = '';
		if (isDiff) {
			significance_str = `Difference is significantly different (Z = ${z.toFixed(3)}).`;
		} else {
			significance_str = `Difference is not significantly different (Z = ${z.toFixed(3)}).`;
		}
		tippy(elem, {
			content: significance_str,
			ignoreAttributes: true,
		});

		return isDiff;
	}

	getRunData(): SimRunData | null {
		if (!this.currentData) {
			return null;
		}

		return SimRunData.create({
			run: this.currentData.simResult.toProto(),
			referenceRun: this.referenceData?.simResult.toProto(),
		});
	}

	getCurrentData(): ReferenceData | null {
		if (!this.currentData) {
			return null;
		}

		// Defensive copy.
		return {
			simResult: this.currentData.simResult,
			settings: JSON.parse(JSON.stringify(this.currentData.settings)),
			raidProto: this.currentData.raidProto,
			encounterProto: this.currentData.encounterProto,
		};
	}

	getReferenceData(): ReferenceData | null {
		if (!this.referenceData) {
			return null;
		}

		// Defensive copy.
		return {
			simResult: this.referenceData.simResult,
			settings: JSON.parse(JSON.stringify(this.referenceData.settings)),
			raidProto: this.referenceData.raidProto,
			encounterProto: this.referenceData.encounterProto,
		};
	}

	static makeToplineResultsContent(simResult: SimResult, filter?: SimResultFilter): HTMLElement {
		const players = simResult.getRaidIndexedPlayers(filter);
		const content = (<></>) as HTMLElement;

		if (players.length === 1) {
			const playerMetrics = players[0];
			if (playerMetrics.getTargetIndex(filter) === null) {
				const dpsMetrics = playerMetrics.dps;
				const dpaspMetrics = playerMetrics.dpasp;
				const tpsMetrics = playerMetrics.tps;
				const dtpsMetrics = playerMetrics.dtps;
				const tmiMetrics = playerMetrics.tmi;
				content.appendChild(
					this.buildResultsLine({
						average: dpsMetrics.avg,
						stdev: dpsMetrics.stdev,
						classes: this.getResultsLineClasses('dps'),
					}),
				);

				// Hide dpasp if it's zero.
				const dpaspContent = this.buildResultsLine({
					average: dpaspMetrics.avg,
					stdev: dpaspMetrics.stdev,
					classes: this.getResultsLineClasses('dpasp'),
				});
				if (dpaspMetrics.avg === 0) {
					dpaspContent.classList.add('hide');
				}
				content.appendChild(dpaspContent);

				content.appendChild(
					this.buildResultsLine({
						average: tpsMetrics.avg,
						stdev: tpsMetrics.stdev,
						classes: this.getResultsLineClasses('tps'),
					}),
				);
				content.appendChild(
					this.buildResultsLine({
						average: dtpsMetrics.avg,
						stdev: dtpsMetrics.stdev,
						classes: this.getResultsLineClasses('dtps'),
					}),
				);

				if (players[0].spec?.specID == Spec.SpecBloodDeathKnight) {
					content.appendChild(
						this.buildResultsLine({
							average: playerMetrics.hps.avg,
							stdev: playerMetrics.hps.stdev,
							classes: this.getResultsLineClasses('hps'),
						}),
					);
				}

				content.appendChild(
					this.buildResultsLine({
						average: tmiMetrics.avg,
						stdev: tmiMetrics.stdev,
						classes: this.getResultsLineClasses('tmi'),
					}),
				);

				content.appendChild(
					this.buildResultsLine({
						average: playerMetrics.chanceOfDeath,
						classes: this.getResultsLineClasses('cod'),
					}),
				);
			} else {
				const actions = simResult.getRaidIndexedActionMetrics(filter);
				if (!!actions.length) {
					const mergedActions = ActionMetrics.merge(actions);
					content.appendChild(
						this.buildResultsLine({
							average: mergedActions.dps,
							classes: this.getResultsLineClasses('dps'),
						}),
					);
					content.appendChild(
						this.buildResultsLine({
							average: mergedActions.tps,
							classes: this.getResultsLineClasses('tps'),
						}),
					);
				}

				const targetActions = simResult
					.getTargets(filter)
					.map(target => target.actions)
					.flat()
					.map(action => action.forTarget({ player: playerMetrics.unitIndex }));
				if (!!targetActions.length) {
					const mergedTargetActions = ActionMetrics.merge(targetActions);
					content.appendChild(
						this.buildResultsLine({
							average: mergedTargetActions.dps,
							classes: this.getResultsLineClasses('dtps'),
						}),
					);
				}

				if (players[0].spec?.specID == Spec.SpecBloodDeathKnight) {
					content.appendChild(
						this.buildResultsLine({
							average: playerMetrics.hps.avg,
							stdev: playerMetrics.hps.stdev,
							classes: this.getResultsLineClasses('hps'),
						}),
					);
				}
			}

			if (players[0].spec?.specID != Spec.SpecBloodDeathKnight) {
				content.appendChild(
					this.buildResultsLine({
						average: playerMetrics.tto.avg,
						stdev: playerMetrics.tto.stdev,
						classes: this.getResultsLineClasses('tto'),
					}),
				);
				content.appendChild(
					this.buildResultsLine({
						average: playerMetrics.hps.avg,
						stdev: playerMetrics.hps.stdev,
						classes: this.getResultsLineClasses('hps'),
					}),
				);
			}
		} else {
			const dpsMetrics = simResult.raidMetrics.dps;
			content.appendChild(
				this.buildResultsLine({
					average: dpsMetrics.avg,
					stdev: dpsMetrics.stdev,
					classes: this.getResultsLineClasses('dps'),
				}),
			);

			const targetActions = simResult
				.getTargets(filter)
				.map(target => target.actions)
				.flat()
				.map(action => action.forTarget(filter));
			if (!!targetActions.length) {
				const mergedTargetActions = ActionMetrics.merge(targetActions);
				content.appendChild(
					this.buildResultsLine({
						average: mergedTargetActions.dps,
						classes: this.getResultsLineClasses('dtps'),
					}),
				);
			}

			const hpsMetrics = simResult.raidMetrics.hps;
			content.appendChild(
				this.buildResultsLine({
					average: hpsMetrics.avg,
					stdev: hpsMetrics.stdev,
					classes: this.getResultsLineClasses('hps'),
				}),
			);
		}

		if (simResult.request.encounter?.useHealth) {
			content.appendChild(
				this.buildResultsLine({
					average: simResult.result.avgIterationDuration,
					classes: this.getResultsLineClasses('dur'),
				}),
			);
		}

		return content;
	}

	private static getResultsLineClasses(metric: keyof ResultMetrics): string {
		const classes = [this.resultMetricClasses[metric]];
		if (this.resultMetricCategories[metric]) classes.push(this.metricsClasses[this.resultMetricCategories[metric]]);

		return classes.join(' ');
	}

	private static buildResultsLine(args: ResultsLineArgs): Element {
		return (
			<div className={`results-metric ${args.classes}`}>
				<span className="topline-result-avg">{args.average.toFixed(2)}</span>
				{args.stdev && (
					<span className="topline-result-stdev">
						(<i className="fas fa-plus-minus fa-xs"></i>
						{args.stdev.toFixed()})
					</span>
				)}
				<div className="results-reference hide">
					<span className="results-reference-diff"></span> vs reference
				</div>
			</div>
		);
	}
}
