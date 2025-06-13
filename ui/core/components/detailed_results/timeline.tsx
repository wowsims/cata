import ApexCharts from 'apexcharts';
import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { CacheHandler } from '../../cache_handler';
import { OtherAction } from '../../proto/common';
import { ResourceType } from '../../proto/spell';
import { ActionId, buffAuraToSpellIdMap, resourceTypeToIcon } from '../../proto_utils/action_id';
import { AuraUptimeLog, CastLog, DpsLog, ResourceChangedLogGroup, SimLog, ThreatLogGroup } from '../../proto_utils/logs_parser';
import { resourceNames } from '../../proto_utils/names';
import SecondaryResource from '../../proto_utils/secondary_resource';
import { UnitMetrics } from '../../proto_utils/sim_result';
import { orderedResourceTypes } from '../../proto_utils/utils';
import { TypedEvent } from '../../typed_event';
import { bucket, distinct, fragmentToString, maxIndex, stringComparator } from '../../utils';
import { actionColors } from './color_settings';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component';

type TooltipHandler = (dataPointIndex: number) => Element;

const dpsColor = '#ed5653';
const manaColor = '#2E93fA';
const threatColor = '#b56d07';

const cachedSpellCastIcon = new CacheHandler<HTMLAnchorElement>();

interface TimelineConfig extends ResultComponentConfig {
	secondaryResource?: SecondaryResource | null;
}

export class Timeline extends ResultComponent {
	private readonly dpsResourcesPlotElem: HTMLElement;
	private dpsResourcesPlot: any;

	private readonly rotationPlotElem: HTMLElement;
	private readonly rotationLabels: HTMLElement;
	private readonly rotationTimeline: HTMLElement;
	private rotationTimelineTimeRulerElem: HTMLCanvasElement | null = null;
	private readonly rotationHiddenIdsContainer: HTMLElement;
	private readonly chartPicker: HTMLSelectElement;

	private prevResultData: SimResultData | null;
	private resultData: SimResultData | null;
	rendered: boolean;

	private hiddenIds: Array<ActionId>;
	private hiddenIdsChangeEmitter;
	private cacheHandler = new CacheHandler<{
		dpsResourcesPlotOptions: any;
		rotationLabels: Timeline['rotationLabels'];
		rotationTimeline: Timeline['rotationTimeline'];
		rotationHiddenIdsContainer: Timeline['rotationHiddenIdsContainer'];
		rotationTimelineTimeRulerElem: Timeline['rotationTimelineTimeRulerElem'];
		rotationTimelineTimeRulerImage: ImageData | undefined;
	}>({
		keysToKeep: 2,
	});

	private secondaryResource?: SecondaryResource | null;

	constructor(config: TimelineConfig) {
		config.rootCssClass = 'timeline-root';
		super(config);
		this.resultData = null;
		this.prevResultData = null;
		this.rendered = false;
		this.hiddenIds = [];
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();
		this.secondaryResource = config.secondaryResource;

		this.rootElem.appendChild(
			<div className="timeline-disclaimer">
				<div className="d-flex flex-column">
					<p>
						<i className="warning fa fa-exclamation-triangle fa-xl me-2"></i>
						Timeline data visualizes only 1 sim iteration.
					</p>
					<p>
						Note: You can move the timeline by holding <kbd>Shift</kbd> while scrolling, or by clicking and dragging.
					</p>
				</div>
				<select className="timeline-chart-picker form-select">
					<option className="rotation-option" value="rotation">
						Rotation
					</option>
					<option className="dps-option" value="dps">
						DPS
					</option>
					<option className="threat-option" value="threat">
						Threat
					</option>
				</select>
			</div>,
		);

		this.rootElem.appendChild(
			<div className="timeline-plots-container">
				<div className="timeline-plot dps-resources-plot hide"></div>
				<div className="timeline-plot rotation-plot">
					<div className="rotation-container">
						<div className="rotation-labels"></div>
						<div className="rotation-timeline" draggable={true}></div>
					</div>
					<div className="rotation-hidden-ids"></div>
				</div>
			</div>,
		);

		this.chartPicker = this.rootElem.querySelector('.timeline-chart-picker')!;
		this.chartPicker.addEventListener('change', () => this.onChartPickerSelectHandler());

		this.dpsResourcesPlotElem = this.rootElem.querySelector('.dps-resources-plot')!;
		this.dpsResourcesPlot = new ApexCharts(this.dpsResourcesPlotElem, {
			chart: {
				animations: {
					enabled: false,
				},
				background: 'transparent',
				foreColor: 'white',
				height: '100%',
				id: 'dpsResources',
				type: 'line',
				zoom: {
					enabled: true,
					allowMouseWheelZoom: false,
				},
			},
			series: [], // Set dynamically
			xaxis: {
				title: {
					text: 'Time (s)',
				},
			},
			noData: {
				text: 'Waiting for data...',
			},
			stroke: {
				width: 2,
				curve: 'straight',
			},
		});

		this.rotationPlotElem = this.rootElem.querySelector('.rotation-plot')!;
		this.rotationLabels = this.rootElem.querySelector('.rotation-labels')!;
		this.rotationTimeline = this.rootElem.querySelector('.rotation-timeline')!;
		this.rotationHiddenIdsContainer = this.rootElem.querySelector('.rotation-hidden-ids')!;

		let isMouseDown = false;
		let startX = 0;
		let scrollLeft = 0;
		this.rotationTimeline.addEventListener('dragstart', event => {
			event.preventDefault();
		});
		this.rotationTimeline.addEventListener('mousedown', event => {
			isMouseDown = true;
			startX = event.pageX - this.rotationTimeline.offsetLeft;
			scrollLeft = this.rotationTimeline.scrollLeft;
		});
		this.rotationTimeline.addEventListener('mouseleave', () => {
			isMouseDown = false;
			this.rotationTimeline.classList.remove('active');
		});
		this.rotationTimeline.addEventListener('mouseup', () => {
			isMouseDown = false;
			this.rotationTimeline.classList.remove('active');
		});
		this.rotationTimeline.addEventListener('mousemove', event => {
			if (!isMouseDown) return;
			event.preventDefault();
			const x = event.pageX - this.rotationTimeline.offsetLeft;
			const walk = (x - startX) * 3; //scroll-fast
			this.rotationTimeline.scrollLeft = scrollLeft - walk;
		});
	}

	onChartPickerSelectHandler() {
		if (this.chartPicker.value === 'rotation') {
			this.dpsResourcesPlotElem.classList.add('hide');
			this.rotationPlotElem.classList.remove('hide');
		} else {
			this.dpsResourcesPlotElem.classList.remove('hide');
			this.rotationPlotElem.classList.add('hide');
		}
	}

	onSimResult(resultData: SimResultData) {
		this.prevResultData = this.resultData;
		this.resultData = resultData;
		this.update();
	}

	private updatePlot() {
		if (this.resultData == null) {
			return;
		}

		const cachedData = this.cacheHandler.get(this.resultData.result.request.requestId);
		if (cachedData) {
			const { dpsResourcesPlotOptions, rotationLabels, rotationTimeline, rotationHiddenIdsContainer, rotationTimelineTimeRulerImage } = cachedData;
			this.rotationLabels.replaceChildren(...rotationLabels.cloneNode(true).childNodes);
			this.rotationTimeline.replaceChildren(...rotationTimeline.cloneNode(true).childNodes);

			this.rotationHiddenIdsContainer.replaceChildren(...rotationHiddenIdsContainer.cloneNode(true).childNodes);
			this.dpsResourcesPlot.updateOptions(dpsResourcesPlotOptions);

			if (rotationTimelineTimeRulerImage)
				this.rotationTimeline
					.querySelector<HTMLCanvasElement>('.rotation-timeline-canvas')
					?.getContext('2d')
					?.putImageData(rotationTimelineTimeRulerImage, 0, 0);

			this.onChartPickerSelectHandler();
			return;
		}

		const duration = this.resultData!.result.result.firstIterationDuration || 1;
		const options: any = {
			theme: {
				mode: 'dark',
			},
			series: [],
			colors: [],
			xaxis: {
				min: 0,
				max: duration,
				tickAmount: 10,
				decimalsInFloat: 1,
				labels: {
					show: true,
				},
				title: {
					text: 'Time (s)',
				},
			},
			yaxis: [],
			chart: {
				events: {
					beforeResetZoom: () => {
						return {
							xaxis: {
								min: 0,
								max: duration,
							},
						};
					},
				},
				toolbar: {
					show: false,
				},
			},
		};

		let tooltipHandlers: Array<TooltipHandler | null> = [];
		options.tooltip = {
			enabled: true,
			custom: (data: { series: any; seriesIndex: number; dataPointIndex: number; w: any }) => {
				if (tooltipHandlers[data.seriesIndex]) {
					return fragmentToString(tooltipHandlers[data.seriesIndex]!(data.dataPointIndex));
				} else {
					throw new Error('No tooltip handler for series ' + data.seriesIndex);
				}
			},
		};

		const players = this.resultData!.result.getRaidIndexedPlayers(this.resultData!.filter);
		if (players.length == 1) {
			const player = players[0];

			const rotationOption = this.rootElem.querySelector('.rotation-option')!;
			rotationOption.classList.remove('hide');
			const threatOption = this.rootElem.querySelector('.threat-option')!;
			threatOption.classList.add('hide');

			try {
				this.updateRotationChart(player, duration);
			} catch (e) {
				console.log('Failed to update rotation chart: ', e);
			}

			const dpsData = this.addDpsSeries(player, options, '');
			this.addDpsYAxis(dpsData.maxDps, options);
			tooltipHandlers.push(dpsData.tooltipHandler);
			tooltipHandlers.push(this.addManaSeries(player, options));
			tooltipHandlers.push(this.addThreatSeries(player, options, ''));
			tooltipHandlers = tooltipHandlers.filter(handler => !!handler);

			this.addMajorCooldownAnnotations(player, options);
		} else {
			if (this.chartPicker.value == 'rotation') {
				this.chartPicker.value = 'dps';
				return;
			}
			const rotationOption = this.rootElem.querySelector('.rotation-option')!;
			rotationOption.classList.add('hide');
			const threatOption = this.rootElem.querySelector('.threat-option')!;
			threatOption.classList.remove('hide');

			this.clearRotationChart();

			if (this.chartPicker.value == 'dps') {
				let maxDps = 0;
				players.forEach(player => {
					const dpsData = this.addDpsSeries(player, options, `var(--bs-${player.classColor}`);
					maxDps = Math.max(maxDps, dpsData.maxDps);
					tooltipHandlers.push(dpsData.tooltipHandler);
				});
				this.addDpsYAxis(maxDps, options);
			} else {
				// threat
				let maxThreat = 0;
				players.forEach(player => {
					tooltipHandlers.push(this.addThreatSeries(player, options, player.classColor));
					maxThreat = Math.max(maxThreat, player.maxThreat);
				});
				this.addThreatYAxis(maxThreat, options);
			}
		}

		this.dpsResourcesPlot.updateOptions(options);

		this.rotationTimelineTimeRulerElem?.toBlob(() => {
			this.cacheHandler.set(this.resultData!.result.request.requestId, {
				dpsResourcesPlotOptions: options,
				rotationLabels: this.rotationLabels.cloneNode(true) as HTMLElement,
				rotationTimeline: this.rotationTimeline.cloneNode(true) as HTMLElement,
				rotationHiddenIdsContainer: this.rotationHiddenIdsContainer.cloneNode(true) as HTMLElement,
				rotationTimelineTimeRulerElem: this.rotationTimelineTimeRulerElem?.cloneNode(true) as HTMLCanvasElement,
				rotationTimelineTimeRulerImage: this.rotationTimelineTimeRulerElem
					?.getContext('2d')
					?.getImageData(0, 0, this.rotationTimelineTimeRulerElem.width, this.rotationTimelineTimeRulerElem.height),
			});
		});
	}

	private addDpsYAxis(maxDps: number, options: any) {
		const dpsAxisMax = Math.ceil(maxDps / 100) * 100;
		options.yaxis.push({
			color: dpsColor,
			seriesName: 'DPS',
			min: 0,
			max: dpsAxisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'DPS',
				style: {
					color: dpsColor,
				},
			},
			axisBorder: {
				show: true,
				color: dpsColor,
			},
			axisTicks: {
				color: dpsColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [dpsColor],
				},
			},
		});
	}

	private addThreatYAxis(maxThreat: number, options: any) {
		const axisMax = Math.ceil(maxThreat / 10000) * 10000;
		options.yaxis.push({
			color: threatColor,
			seriesName: 'Threat',
			min: 0,
			max: axisMax,
			tickAmount: 10,
			decimalsInFloat: 0,
			title: {
				text: 'Threat',
				style: {
					color: threatColor,
				},
			},
			axisBorder: {
				show: true,
				color: threatColor,
			},
			axisTicks: {
				color: threatColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [threatColor],
				},
			},
		});
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addDpsSeries(unit: UnitMetrics, options: any, colorOverride: string): { maxDps: number; tooltipHandler: TooltipHandler } {
		const dpsLogs = unit.dpsLogs.filter(log => log.timestamp >= 0);

		options.colors.push(colorOverride || dpsColor);
		options.series.push({
			name: 'DPS',
			type: 'line',
			data: dpsLogs.map(log => {
				return {
					x: log.timestamp,
					y: log.dps,
				};
			}),
		});

		return {
			maxDps: dpsLogs[maxIndex(dpsLogs.map(l => l.dps))!]?.dps,
			tooltipHandler: (dataPointIndex: number) => {
				const log = dpsLogs[dataPointIndex];
				return this.dpsTooltip(log, true, unit, colorOverride);
			},
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addManaSeries(unit: UnitMetrics, options: any): TooltipHandler | null {
		const manaLogs = unit.groupedResourceLogs[ResourceType.ResourceTypeMana].filter(log => log.timestamp >= 0);
		if (manaLogs.length == 0) {
			return null;
		}
		const maxMana = manaLogs[0].valueBefore;

		options.colors.push(manaColor);
		options.series.push({
			name: 'Mana',
			type: 'line',
			data: manaLogs.map(log => {
				return {
					x: log.timestamp,
					y: log.valueAfter,
				};
			}),
		});
		options.yaxis.push({
			seriesName: 'Mana',
			opposite: true, // Appear on right side
			min: 0,
			max: maxMana,
			tickAmount: 10,
			title: {
				text: 'Mana',
				style: {
					color: manaColor,
				},
			},
			axisBorder: {
				show: true,
				color: manaColor,
			},
			axisTicks: {
				color: manaColor,
			},
			labels: {
				minWidth: 30,
				style: {
					colors: [manaColor],
				},
				formatter: (val: string) => {
					const v = parseFloat(val);
					return `${v.toFixed(0)} (${((v / maxMana) * 100).toFixed(0)}%)`;
				},
			},
		} as any);

		return (dataPointIndex: number) => {
			const log = manaLogs[dataPointIndex];
			return this.resourceTooltip(log, maxMana, true);
		};
	}

	// Returns a function for drawing the tooltip, or null if no series was added.
	private addThreatSeries(unit: UnitMetrics, options: any, colorOverride: string): TooltipHandler | null {
		options.colors.push(colorOverride || threatColor);
		options.series.push({
			name: 'Threat',
			type: 'line',
			data: unit.threatLogs
				.filter(log => log.timestamp >= 0)
				.map(log => {
					return {
						x: log.timestamp,
						y: log.threatAfter,
					};
				}),
		});

		return (dataPointIndex: number) => {
			const log = unit.threatLogs[dataPointIndex];
			return this.threatTooltip(log, true, unit, colorOverride);
		};
	}

	private addMajorCooldownAnnotations(unit: UnitMetrics, options: any) {
		const mcdLogs = unit.majorCooldownLogs;
		const mcdAuraLogs = unit.majorCooldownAuraUptimeLogs;

		// Figure out how much to vertically offset cooldown icons, for cooldowns
		// used very close to each other. This is so the icons don't overlap.
		const MAX_ALLOWED_DIST = 10;
		const cooldownIconOffsets = mcdLogs.map(
			(mcdLog, mcdIdx) => mcdLogs.filter((cdLog, cdIdx) => cdIdx < mcdIdx && cdLog.timestamp > mcdLog.timestamp - MAX_ALLOWED_DIST).length,
		);

		const distinctMcdAuras = distinct(mcdAuraLogs, (a, b) => a.actionId!.equalsIgnoringTag(b.actionId!));
		// Sort by name so auras keep their same colors even if timings change.
		distinctMcdAuras.sort((a, b) => stringComparator(a.actionId!.name, b.actionId!.name));
		const mcdAuraColors = mcdAuraLogs.map(
			mcdAuraLog => actionColors[distinctMcdAuras.findIndex(dAura => dAura.actionId!.equalsIgnoringTag(mcdAuraLog.actionId!))],
		);

		options.annotations = {
			position: 'back',
			xaxis: mcdAuraLogs.map((log, i) => {
				return {
					x: log.gainedAt,
					x2: log.fadedAt,
					fillColor: mcdAuraColors[i],
				};
			}),
			points: mcdLogs.map((log, i) => {
				return {
					x: log.timestamp,
					y: 0,
					image: {
						path: log.actionId!.iconUrl,
						width: 20,
						height: 20,
						offsetY: cooldownIconOffsets[i] * -25,
					},
				};
			}),
		};
	}

	private clearRotationChart() {
		this.rotationLabels.replaceChildren(<div className="rotation-label-header"></div>);
		const canvasRef = ref<HTMLCanvasElement>();
		this.rotationTimeline.replaceChildren(
			<div className="rotation-timeline-header">
				<canvas ref={canvasRef} className="rotation-timeline-canvas" />
			</div>,
		);
		this.rotationTimelineTimeRulerElem = canvasRef.value || null;
		this.rotationHiddenIdsContainer.replaceChildren();
		this.hiddenIdsChangeEmitter = new TypedEvent<void>();
	}

	private updateRotationChart(player: UnitMetrics, duration: number) {
		const targets = this.resultData!.result.getTargets(this.resultData!.filter);
		if (targets.length == 0) {
			return;
		}
		const target = targets[0];

		this.clearRotationChart();

		try {
			this.drawRotationTimeRuler(this.rotationTimeline.querySelector('.rotation-timeline-canvas')!, duration);
		} catch (e) {
			console.log('Failed to draw rotation: ', e);
		}

		orderedResourceTypes.forEach(resourceType => this.addResourceRow(resourceType, player.groupedResourceLogs[resourceType], duration));

		const buffsById = Object.values(bucket(player.auraUptimeLogs, log => log.actionId!.toString()));
		buffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const debuffsById = Object.values(bucket(target.auraUptimeLogs, log => log.actionId!.toString()));
		debuffsById.sort((a, b) => stringComparator(a[0].actionId!.name, b[0].actionId!.name));
		const buffsAndDebuffsById = buffsById.concat(debuffsById);

		auraAsResource.forEach(auraId => {
			const auraIndex = buffsById.findIndex(auraUptimeLogs => auraUptimeLogs?.[0].actionId!.spellId === auraId);
			if (auraIndex !== -1) {
				this.addAuraRow(buffsById[auraIndex], duration);
			}
		});

		const playerCastsByAbility = this.getSortedCastsByAbility(player);
		playerCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));

		if (player.pets.length > 0) {
			const playerPets = new Map<string, UnitMetrics>();
			player.pets.forEach(petsLog => {
				const petCastsByAbility = this.getSortedCastsByAbility(petsLog);
				if (petCastsByAbility.length > 0) {
					// Because multiple pets can have the same name and we parse cast logs
					// by pet name each individual pet ends up with all the casts of pets
					// with the same name. Because of this we can just grab the first pet
					// of each name and visualize only that.
					if (!playerPets.has(petsLog.name)) {
						playerPets.set(petsLog.name, petsLog);
					}
				}
			});

			playerPets.forEach(pet => {
				this.addSeparatorRow(duration);
				this.addPetRow(pet.name, duration);
				orderedResourceTypes.forEach(resourceType => this.addResourceRow(resourceType, pet.groupedResourceLogs[resourceType], duration));
				const petCastsByAbility = this.getSortedCastsByAbility(pet);
				petCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));
			});
		}

		// Don't add a row for buffs that were already visualized in a cast row or are prioritized.
		const buffsToShow = buffsById.filter(auraUptimeLogs =>
			playerCastsByAbility.findIndex(
				casts => auraUptimeLogs[0].actionId && (casts[0].actionId!.equalsIgnoringTag(auraUptimeLogs[0].actionId) || auraAsResource.includes(auraUptimeLogs[0].actionId.anyId())),
			),
		);
		if (buffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			buffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}

		const targetCastsByAbility = this.getSortedCastsByAbility(target);
		if (targetCastsByAbility.length > 0) {
			this.addSeparatorRow(duration);
			targetCastsByAbility.forEach(castLogs => this.addCastRow(castLogs, buffsAndDebuffsById, duration));
		}

		// Add a row for all debuffs, even those which have already been visualized in a cast row.
		const debuffsToShow = debuffsById;
		if (debuffsToShow.length > 0) {
			this.addSeparatorRow(duration);
			debuffsToShow.forEach(auraUptimeLogs => this.addAuraRow(auraUptimeLogs, duration));
		}
	}

	private getSortedCastsByAbility(player: UnitMetrics): Array<Array<CastLog>> {
		const meleeActionIds = player.getMeleeActions().map(action => action.actionId);
		const spellActionIds = player.getSpellActions().map(action => action.actionId);
		const getActionCategory = (actionId: ActionId): number => {
			const fixedCategory = idToCategoryMap[actionId.anyId()];
			if (fixedCategory != null) {
				return fixedCategory;
			} else if (meleeActionIds.find(meleeActionId => meleeActionId.equals(actionId))) {
				return MELEE_ACTION_CATEGORY;
			} else if (spellActionIds.find(spellActionId => spellActionId.equals(actionId))) {
				return SPELL_ACTION_CATEGORY;
			} else {
				return DEFAULT_ACTION_CATEGORY;
			}
		};

		const castsByAbility = Object.values(
			bucket(player.castLogs, log => {
				if (idsToGroupForRotation.includes(log.actionId!.spellId)) {
					return log.actionId!.toStringIgnoringTag();
				} else {
					return log.actionId!.toString();
				}
			}),
		);

		castsByAbility.sort((a, b) => {
			const categoryA = getActionCategory(a[0].actionId!);
			const categoryB = getActionCategory(b[0].actionId!);
			if (categoryA != categoryB) {
				return categoryA - categoryB;
			} else if (a[0].actionId!.anyId() == b[0].actionId!.anyId()) {
				return a[0].actionId!.tag - b[0].actionId!.tag;
			} else {
				return stringComparator(a[0].actionId!.name, b[0].actionId!.name);
			}
		});

		return castsByAbility;
	}

	private makeLabelElem(actionId: ActionId, isHiddenLabel: boolean, isAura?: boolean): JSX.Element {
		const labelText = idsToGroupForRotation.includes(actionId.spellId) ? actionId.baseName : actionId.name;
		const labelIcon = ref<HTMLAnchorElement>();
		const hideElem = ref<HTMLElement>();
		const labelElem = (
			<div className={clsx('rotation-label rotation-row', isHiddenLabel && 'rotation-label-hidden')}>
				<span ref={hideElem} className={clsx('fas', isHiddenLabel ? 'fa-eye' : 'fa-eye-slash')}></span>
				<a ref={labelIcon} className="rotation-label-icon"></a>
				<span className="rotation-label-text">{labelText}</span>
			</div>
		);
		const onClickHandler = () => {
			if (isHiddenLabel) {
				const index = this.hiddenIds.findIndex(hiddenId => hiddenId.equals(actionId));
				if (index != -1) {
					this.hiddenIds.splice(index, 1);
				}
			} else {
				this.hiddenIds.push(actionId);
			}
			this.hiddenIdsChangeEmitter.emit(TypedEvent.nextEventID());
		};
		hideElem.value!.addEventListener('click', onClickHandler);
		const tooltip = tippy(hideElem.value!, {
			theme: 'timeline-tooltip',
			placement: 'bottom',
			content: isHiddenLabel ? 'Show Row' : 'Hide Row',
		});

		const updateHidden = () => {
			if (isHiddenLabel == Boolean(this.hiddenIds.find(hiddenId => hiddenId.equals(actionId)))) {
				labelElem.classList.remove('hide');
			} else {
				labelElem.classList.add('hide');
			}
		};
		const event = this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		actionId.setBackgroundAndHref(labelIcon.value!);
		actionId.setWowheadDataset(labelIcon.value!, { useBuffAura: isAura });

		this.addOnResetCallback(() => {
			hideElem.value?.removeEventListener('click', onClickHandler);
			tooltip.destroy();
			event.dispose();
		});

		return labelElem;
	}

	private makeRowElem(actionId: ActionId, duration: number): JSX.Element {
		const rowElem = (
			<div
				className="rotation-timeline-row rotation-row"
				style={{
					width: this.timeToPx(duration),
				}}></div>
		);

		const updateHidden = () => {
			if (this.hiddenIds.find(hiddenId => hiddenId.equals(actionId))) {
				rowElem.classList.add('hide');
			} else {
				rowElem.classList.remove('hide');
			}
		};
		const event = this.hiddenIdsChangeEmitter.on(updateHidden);
		updateHidden();
		this.addOnResetCallback(() => event.dispose());
		return rowElem;
	}

	private addPetRow(petName: string, duration: number) {
		const actionId = ActionId.fromPetName(petName);
		const rowElem = this.makeRowElem(actionId, duration);

		const iconElem = document.createElement('div');
		this.rotationLabels.appendChild(iconElem);

		actionId.fill().then(filledActionId => {
			const labelText = idsToGroupForRotation.includes(filledActionId.spellId) ? filledActionId.baseName : filledActionId.name;
			const labelIcon = ref<HTMLAnchorElement>();
			const labelElem = (
				<div className="rotation-label rotation-row">
					<a ref={labelIcon} className="rotation-label-icon"></a>
					<span className="rotation-label-text">{labelText}</span>
				</div>
			);
			filledActionId.setBackgroundAndHref(labelIcon.value!);
			iconElem.appendChild(labelElem);
		});

		this.rotationTimeline.appendChild(rowElem);
	}

	private addSeparatorRow(duration: number) {
		const separatorElem = <div className="rotation-timeline-separator"></div>;
		this.rotationLabels.appendChild(separatorElem.cloneNode());
		separatorElem.style.width = this.timeToPx(duration);
		this.rotationTimeline.appendChild(separatorElem);
	}

	private addResourceRow(resourceType: ResourceType, resourceLogs: Array<ResourceChangedLogGroup>, duration: number) {
		if (resourceLogs.length == 0) {
			return;
		}
		const startValue = function (group: ResourceChangedLogGroup): number {
			if (group.maxValue == null) {
				return resourceLogs[0].valueBefore;
			}

			return group.maxValue;
		};

		let resourceName = resourceNames.get(resourceType);
		let resourceIcon = resourceTypeToIcon[resourceType];
		if (resourceType == ResourceType.ResourceTypeGenericResource && !!this.secondaryResource) {
			resourceName = this.secondaryResource.name;
			resourceIcon = this.secondaryResource.icon || '';
		}

		const labelElem = (
			<div className="rotation-label rotation-row">
				<a
					className="rotation-label-icon"
					style={{
						backgroundImage: `url('${resourceIcon}')`,
					}}></a>
				<span className="rotation-label-text">{resourceName}</span>
			</div>
		);

		this.rotationLabels.appendChild(labelElem);

		const rowElem = (
			<div
				className="rotation-timeline-row rotation-row"
				style={{
					width: this.timeToPx(duration),
				}}></div>
		);

		resourceLogs.forEach((resourceLogGroup, i) => {
			const cNames = resourceNames.get(resourceType)!.toLowerCase().replaceAll(' ', '-');
			const resourceElem = (
				<div
					className={`rotation-timeline-resource series-color ${cNames}`}
					style={{
						left: this.timeToPx(resourceLogGroup.timestamp),
						width: this.timeToPx((resourceLogs[i + 1]?.timestamp || duration) - resourceLogGroup.timestamp),
					}}></div>
			);

			if (percentageResources.includes(resourceType)) {
				resourceElem.textContent = ((resourceLogGroup.valueAfter / startValue(resourceLogGroup)) * 100).toFixed(0) + '%';
			} else {
				if (
					resourceType == ResourceType.ResourceTypeEnergy ||
					resourceType == ResourceType.ResourceTypeFocus ||
					resourceType == ResourceType.ResourceTypeSolarEnergy ||
					resourceType == ResourceType.ResourceTypeLunarEnergy
				) {
					const bgElem = document.createElement('div');
					bgElem.classList.add('rotation-timeline-resource-fill');
					bgElem.classList.add(cNames);
					bgElem.style.height = ((resourceLogGroup.valueAfter / startValue(resourceLogGroup)) * 100).toFixed(0) + '%';
					resourceElem.appendChild(bgElem);
				} else {
					resourceElem.textContent = Math.floor(resourceLogGroup.valueAfter).toFixed(0);
				}
			}
			rowElem.appendChild(resourceElem);

			const tooltip = tippy(resourceElem, {
				placement: 'bottom',
				content: this.resourceTooltipElem(resourceLogGroup, startValue(resourceLogGroup), false),
			});
			this.addOnResetCallback(() => tooltip.destroy());
		});
		this.rotationTimeline.appendChild(rowElem);
	}

	private addCastRow(castLogs: Array<CastLog>, aurasById: Array<Array<AuraUptimeLog>>, duration: number) {
		const actionId = castLogs[0].actionId!;

		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true));

		const rowElem = this.makeRowElem(actionId, duration);
		castLogs.forEach(castLog => {
			const castElem = (
				<div
					className="rotation-timeline-cast"
					style={{
						left: this.timeToPx(castLog.timestamp),
						minWidth: this.timeToPx(castLog.castTime + castLog.travelTime),
					}}
				/>
			);
			rowElem.appendChild(castElem);

			if (castLog.travelTime != 0) {
				const travelTimeElem = (
					<div
						className="rotation-timeline-travel-time"
						style={{
							left: this.timeToPx(castLog.castTime),
							minWidth: this.timeToPx(castLog.travelTime),
						}}
					/>
				);
				castElem.appendChild(travelTimeElem);
			}

			if (castLog.damageDealtLogs.length > 0) {
				const ddl = castLog.damageDealtLogs[0];
				if (ddl.miss || ddl.dodge || ddl.parry) {
					castElem.classList.add('outcome-miss');
				} else if (ddl.glance || ddl.block || ddl.partialResist1_4 || ddl.partialResist2_4 || ddl.partialResist3_4) {
					castElem.classList.add('outcome-partial');
				} else if (ddl.crit) {
					castElem.classList.add('outcome-crit');
				} else {
					castElem.classList.add('outcome-hit');
				}
			}

			const actionIdAsString = actionId.toString();
			const cachedIconElem = cachedSpellCastIcon.get(actionIdAsString)?.cloneNode() as HTMLAnchorElement | undefined;
			let iconElem = cachedIconElem;
			if (!iconElem) {
				iconElem = (<a className="rotation-timeline-cast-icon" />) as HTMLAnchorElement;
				actionId.setBackground(iconElem);
				cachedSpellCastIcon.set(actionIdAsString, iconElem);
			}
			castElem.appendChild(iconElem);

			const travelTimeStr = castLog.travelTime == 0 ? '' : ` + ${castLog.travelTime.toFixed(2)}s travel time`;
			const totalDamage = castLog.totalDamage();

			const tt = (
				<div className="timeline-tooltip">
					<span>
						{castLog.actionId!.name} from {castLog.timestamp.toFixed(2)}s to {(castLog.timestamp + castLog.castTime).toFixed(2)}s (
						{castLog.castTime > 0 && `${castLog.castTime.toFixed(2)}s, `} {castLog.effectiveTime.toFixed(2)}s GCD Time)
						{travelTimeStr.length > 0 && travelTimeStr}
					</span>
					{castLog.damageDealtLogs.length > 0 && (
						<ul className="rotation-timeline-cast-damage-list">
							{castLog.damageDealtLogs.map(ddl => (
								<li>
									<span>
										{ddl.timestamp.toFixed(2)}s - {ddl.result()}
									</span>
									{ddl.source?.isTarget && <span className="threat-metrics"> ({ddl.threat.toFixed(1)} Threat)</span>}
								</li>
							))}
						</ul>
					)}
					{totalDamage > 0 && (
						<span>
							Total: {totalDamage.toFixed(2)} ({(totalDamage / (castLog.effectiveTime || 1)).toFixed(2)} DPET)
						</span>
					)}
				</div>
			);

			const tooltip = tippy(castElem, {
				placement: 'bottom',
				content: tt,
			});
			this.addOnResetCallback(() => tooltip.destroy());

			castLog.damageDealtLogs
				.filter(ddl => ddl.tick)
				.forEach(ddl => {
					const tickElem = (
						<div
							className="rotation-timeline-tick"
							style={{
								left: this.timeToPx(ddl.timestamp),
							}}
						/>
					);
					rowElem.appendChild(tickElem);

					const tt = (
						<div className="timeline-tooltip">
							<span>
								{ddl.timestamp.toFixed(2)}s - {ddl.actionId!.name} {ddl.result()}
							</span>
							{ddl.source?.isTarget && <span className="threat-metrics"> ({ddl.threat.toFixed(1)} Threat)</span>}
						</div>
					);

					const tooltip = tippy(tickElem, {
						placement: 'bottom',
						content: tt,
					});
					this.addOnResetCallback(() => tooltip.destroy());
				});
		});

		// If there are any auras that correspond to this cast, visualize them in the same row.
		aurasById
			.filter(auraUptimeLogs => actionId.equals(buffAuraToSpellIdMap[auraUptimeLogs[0].actionId!.spellId] ?? auraUptimeLogs[0].actionId!))
			.forEach(auraUptimeLogs => this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem, true));

		this.rotationTimeline.appendChild(rowElem);
	}

	private addAuraRow(auraUptimeLogs: Array<AuraUptimeLog>, duration: number) {
		const actionId = auraUptimeLogs[0].actionId!;

		const rowElem = this.makeRowElem(actionId, duration);
		this.rotationLabels.appendChild(this.makeLabelElem(actionId, false, true));
		this.rotationHiddenIdsContainer.appendChild(this.makeLabelElem(actionId, true, true));
		this.rotationTimeline.appendChild(rowElem);

		this.applyAuraUptimeLogsToRow(auraUptimeLogs, rowElem, false);
	}

	private applyAuraUptimeLogsToRow(auraUptimeLogs: Array<AuraUptimeLog>, rowElem: JSX.Element, hasCast: boolean) {
		auraUptimeLogs.forEach(aul => {
			const auraElem = (
				<div
					className="rotation-timeline-aura"
					style={{
						left: this.timeToPx(aul.gainedAt),
						minWidth: this.timeToPx(aul.fadedAt === aul.gainedAt ? 0.001 : aul.fadedAt - aul.gainedAt),
					}}
				/>
			);
			rowElem.appendChild(auraElem);

			const tt = (
				<div className="timeline-tooltip">
					<span>
						{aul.actionId!.name}: {aul.gainedAt.toFixed(2)}s - {aul.fadedAt.toFixed(2)}s
					</span>
				</div>
			);

			const tooltip = tippy(auraElem, {
				placement: 'bottom',
				content: tt,
			});
			this.addOnResetCallback(() => tooltip.destroy());

			aul.stacksChange.forEach((scl, i) => {
				if (scl.timestamp == aul.fadedAt) {
					return;
				}

				const stacksChangeElem = (
					<div
						className="rotation-timeline-stacks-change"
						style={{
							left: this.timeToPx(scl.timestamp - aul.timestamp),
							width: this.timeToPx(aul.stacksChange[i + 1] ? aul.stacksChange[i + 1].timestamp - scl.timestamp : aul.fadedAt - scl.timestamp),
							textIndent: hasCast ? '30px' : undefined,
						}}>
						{String(scl.newStacks)}
					</div>
				);
				auraElem.appendChild(stacksChangeElem);
			});
		});
	}

	private timeToPxValue(time: number): number {
		return time * 100;
	}
	private timeToPx(time: number): string {
		return this.timeToPxValue(time) + 'px';
	}

	private drawRotationTimeRuler(canvas: HTMLCanvasElement, duration: number) {
		const height = 30;
		canvas.width = this.timeToPxValue(duration);
		canvas.height = height;

		const ctx = canvas.getContext('2d')!;
		ctx.strokeStyle = 'white';

		ctx.font = 'bold 14px SimDefaultFont';
		ctx.fillStyle = 'white';
		ctx.lineWidth = 2;
		ctx.beginPath();

		// Bottom border line
		ctx.moveTo(0, height);
		ctx.lineTo(canvas.width, height);

		// Tick lines
		const numTicks = 1 + Math.floor(duration * 10);
		for (let i = 0; i <= numTicks; i++) {
			const time = i * 0.1;
			let x = this.timeToPxValue(time);
			if (i == 0) {
				ctx.textAlign = 'left';
				x++;
			} else if (i % 10 == 0 && time + 1 > duration) {
				ctx.textAlign = 'right';
				x--;
			} else {
				ctx.textAlign = 'center';
			}

			let lineHeight = 0;
			if (i % 10 == 0) {
				lineHeight = height * 0.5;
				ctx.fillText(time + 's', x, height - height * 0.6);
			} else if (i % 5 == 0) {
				lineHeight = height * 0.25;
			} else {
				lineHeight = height * 0.125;
			}
			ctx.moveTo(x, height);
			ctx.lineTo(x, height - lineHeight);
		}
		ctx.stroke();
	}

	private dpsTooltip(log: DpsLog, _includeAuras: boolean, player: UnitMetrics, colorOverride: string) {
		const showPlayerLabel = colorOverride != '';
		return (
			<div className="timeline-tooltip dps">
				<div className="timeline-tooltip-header">
					{showPlayerLabel ? (
						<>
							<img className="timeline-tooltip-icon" src="${player.iconUrl}" />
							<span className="" style="color: ${colorOverride}">
								{player.label}
							</span>
							<span> - </span>
						</>
					) : null}
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<ul className="timeline-dps-events">{log.damageLogs.map(damageLog => this.tooltipLogItem(damageLog, damageLog.result()))}</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">DPS: {log.dps.toFixed(2)}</span>
					</div>
				</div>
				{this.tooltipAurasSection(log)}
			</div>
		);
	}

	private threatTooltip(log: ThreatLogGroup, includeAuras: boolean, player: UnitMetrics, colorOverride: string) {
		const showPlayerLabel = colorOverride != '';
		return (
			<div className="timeline-tooltip threat">
				<div className="timeline-tooltip-header">
					{showPlayerLabel ? (
						<>
							<img className="timeline-tooltip-icon" src={player.iconUrl} />
							<span className="" style={{ color: colorOverride }}>
								{player.label}
							</span>
							<span> - </span>
						</>
					) : null}
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<div className="timeline-tooltip-body-row">
						<span className="series-color">Before: {log.threatBefore.toFixed(1)}</span>
					</div>
					<ul className="timeline-threat-events">{log.logs.map(log => this.tooltipLogItem(log, <>{log.threat.toFixed(1)} Threat</>))}</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">After: {log.threatAfter.toFixed(1)}</span>
					</div>
				</div>
				{includeAuras ? this.tooltipAurasSection(log) : null}
			</div>
		);
	}

	private resourceTooltipElem(log: ResourceChangedLogGroup, maxValue: number, includeAuras: boolean) {
		const valToDisplayString = percentageResources.includes(log.resourceType)
			? (val: number) => `${val.toFixed(1)} (${((val / maxValue) * 100).toFixed(0)}%)`
			: (val: number) => `${val.toFixed(1)}`;

		return (
			<div className={`timeline-tooltip ${resourceNames.get(log.resourceType)!.toLowerCase().replaceAll(' ', '-')}`}>
				<div className="timeline-tooltip-header">
					<span className="bold">{log.timestamp.toFixed(2)}s</span>
				</div>
				<div className="timeline-tooltip-body">
					<div className="timeline-tooltip-body-row">
						<span className="series-color">Before: {valToDisplayString(log.valueBefore)}</span>
					</div>
					<ul className="timeline-mana-events">
						{log.logs.map(manaChangedLog => this.tooltipLogItemElem(manaChangedLog, <>{manaChangedLog.resultString()}</>))}
					</ul>
					<div className="timeline-tooltip-body-row">
						<span className="series-color">After: {valToDisplayString(log.valueAfter)}</span>
					</div>
				</div>
				{includeAuras && this.tooltipAurasSectionElem(log)}
			</div>
		);
	}

	private resourceTooltip(log: ResourceChangedLogGroup, maxValue: number, includeAuras: boolean) {
		return this.resourceTooltipElem(log, maxValue, includeAuras);
	}

	private tooltipLogItem(log: SimLog, value: Element) {
		return this.tooltipLogItemElem(log, value);
	}

	private tooltipLogItemElem(log: SimLog, value: Element): JSX.Element {
		return (
			<li>
				{log.actionId && log.actionId.iconUrl && <img className="timeline-tooltip-icon" src={log.actionId.iconUrl}></img>}
				{log.actionId && <span>{log.actionId.name}</span>}
				<span className="series-color">{value}</span>
			</li>
		);
	}

	private tooltipAurasSection(log: SimLog) {
		if (log.activeAuras.length == 0) {
			return '';
		}
		return this.tooltipAurasSectionElem(log);
	}

	private tooltipAurasSectionElem(log: SimLog): JSX.Element {
		if (log.activeAuras.length == 0) {
			return <></>;
		}

		return (
			<div className="timeline-tooltip-auras">
				<div className="timeline-tooltip-body-row">
					<span className="bold">Active Auras</span>
				</div>
				<ul className="timeline-active-auras">
					{log.activeAuras.map(auraLog => (
						<li>
							{auraLog.actionId!.iconUrl && <img className="timeline-tooltip-icon" src={auraLog.actionId!.iconUrl}></img>}
							<span>{auraLog.actionId!.name}</span>
						</li>
					))}
				</ul>
			</div>
		);
	}

	update() {
		this.reset();
		if (!this.rendered) this.dpsResourcesPlot.render();
		this.updatePlot();
		this.rendered = true;
	}

	render() {
		if (this.rendered) return;
		this.update();
	}

	reset() {
		const previousResultRequestId = this.prevResultData?.result.request.requestId;
		if (previousResultRequestId && !this.cacheHandler.get(previousResultRequestId)) return;
		super.reset();
	}
}

const MELEE_ACTION_CATEGORY = 1;
const SPELL_ACTION_CATEGORY = 2;
const DEFAULT_ACTION_CATEGORY = 3;

const auraAsResource = [
	// Vengeance
	84840, // Druid
	84839, // Paladin
	93098, // Warrior
	93099, // Death Knight
	120267, // Monk

	// Monk
	124255, // Stagger
];

// Hard-coded spell categories for controlling rotation ordering.
const idToCategoryMap: Record<number, number> = {
	[OtherAction.OtherActionMove]: 0,
	[OtherAction.OtherActionAttack]: 0.01,
	[OtherAction.OtherActionShoot]: 0.5,

	// Druid
	[48480]: 0.1, // Maul
	[48564]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Bear)
	[48568]: MELEE_ACTION_CATEGORY + 0.2, // Lacerate
	[48562]: MELEE_ACTION_CATEGORY + 0.3, // Swipe (Bear)

	[48566]: MELEE_ACTION_CATEGORY + 0.1, // Mangle (Cat)
	[48572]: MELEE_ACTION_CATEGORY + 0.2, // Shred
	[49800]: MELEE_ACTION_CATEGORY + 0.51, // Rip
	[52610]: MELEE_ACTION_CATEGORY + 0.52, // Savage Roar
	[48577]: MELEE_ACTION_CATEGORY + 0.53, // Ferocious Bite

	[48465]: SPELL_ACTION_CATEGORY + 0.1, // Starfire
	[48461]: SPELL_ACTION_CATEGORY + 0.2, // Wrath
	[53201]: SPELL_ACTION_CATEGORY + 0.3, // Starfall
	[48468]: SPELL_ACTION_CATEGORY + 0.4, // Insect Swarm
	[48463]: SPELL_ACTION_CATEGORY + 0.5, // Moonfire

	// Hunter
	[53217]: 0.6, // Wild Quiver
	[53209]: MELEE_ACTION_CATEGORY + 0.1, // Chimera Shot
	[53353]: MELEE_ACTION_CATEGORY + 0.11, // Chimera Shot Serpent
	[53301]: MELEE_ACTION_CATEGORY + 0.1, // Explosive Shot
	[1215485]: MELEE_ACTION_CATEGORY + 0.12, // Explosive Shot
	[49050]: MELEE_ACTION_CATEGORY + 0.2, // Aimed Shot
	[49048]: MELEE_ACTION_CATEGORY + 0.21, // Multi Shot
	[3044]: MELEE_ACTION_CATEGORY + 0.22, // Arcane Shot
	[56641]: MELEE_ACTION_CATEGORY + 0.27, // Steady Shot
	[53351]: MELEE_ACTION_CATEGORY + 0.28, // Kill Shot
	[34490]: MELEE_ACTION_CATEGORY + 0.29, // Silencing Shot
	[1978]: MELEE_ACTION_CATEGORY + 0.3, // Serpent Sting
	[53238]: MELEE_ACTION_CATEGORY + 0.31, // Piercing Shots
	[63672]: MELEE_ACTION_CATEGORY + 0.32, // Black Arrow
	[49067]: MELEE_ACTION_CATEGORY + 0.33, // Explosive Trap
	[77767]: MELEE_ACTION_CATEGORY + 0.34, // Cobra Shot

	// Paladin
	[76672]: MELEE_ACTION_CATEGORY + 0.01, // Hand of Light (mastery)
	[35395]: MELEE_ACTION_CATEGORY + 0.02, // Crusader Strike
	[99092]: MELEE_ACTION_CATEGORY + 0.03, // Flames of the Faithful (ret T12 2pc)
	[53595]: MELEE_ACTION_CATEGORY + 0.04, // Hammer of the Righteous (Physical)
	[88263]: MELEE_ACTION_CATEGORY + 0.05, // Hammer of the Righteous (Holy)
	[53385]: MELEE_ACTION_CATEGORY + 0.06, // Divine Storm
	[85256]: MELEE_ACTION_CATEGORY + 0.07, // Templar's Verdict
	[20271]: MELEE_ACTION_CATEGORY + 0.08, // Judgment
	[42463]: MELEE_ACTION_CATEGORY + 0.09, // Seal of Truth (on-hit)
	[31803]: MELEE_ACTION_CATEGORY + 0.1, // Censure (Seal of Truth)
	[101423]: MELEE_ACTION_CATEGORY + 0.11, // Seal of Righteousness
	[53600]: MELEE_ACTION_CATEGORY + 0.12, // Shield of the Righteous
	[99075]: MELEE_ACTION_CATEGORY + 0.13, // Righteous Flames (prot T12 2pc)
	[879]: MELEE_ACTION_CATEGORY + 0.15, // Exorcism
	[26573]: MELEE_ACTION_CATEGORY + 0.16, // Consecration
	[119072]: MELEE_ACTION_CATEGORY + 0.17, // Holy Wrath
	[24275]: MELEE_ACTION_CATEGORY + 0.18, // Hammer of Wrath
	[114852]: MELEE_ACTION_CATEGORY + 0.19, // Holy Prism (Damage)
	[114919]: MELEE_ACTION_CATEGORY + 0.19, // Arcing Light (Damage)
	[114916]: MELEE_ACTION_CATEGORY + 0.19, // Execution Sentence
	[114871]: MELEE_ACTION_CATEGORY + 0.2, // Holy Prism (Heal)
	[119952]: MELEE_ACTION_CATEGORY + 0.2, // Arcing Light (Heal)
	[146586]: MELEE_ACTION_CATEGORY + 0.2, // Stay of Execution
	[84963]: SPELL_ACTION_CATEGORY + 0.01, // Inquisition
	[54428]: SPELL_ACTION_CATEGORY + 0.02, // Divine Plea
	[498]: SPELL_ACTION_CATEGORY + 0.03, // Divine Protection
	[99090]: SPELL_ACTION_CATEGORY + 0.04, // Flaming Aegis (Prot T12 4pc)
	[66233]: SPELL_ACTION_CATEGORY + 0.05, // Ardent Defender
	[31884]: SPELL_ACTION_CATEGORY + 0.06, // Avenging Wrath
	[114232]: SPELL_ACTION_CATEGORY + 0.07, // Sanctified Wrath
	[105809]: SPELL_ACTION_CATEGORY + 0.08, // Holy Avenger,
	[86698]: SPELL_ACTION_CATEGORY + 0.09, // Guardian of Ancient Kings
	[86704]: SPELL_ACTION_CATEGORY + 0.1, // Ancient Fury
	[20925]: SPELL_ACTION_CATEGORY + 0.11, // Sacred Shield (Ret / Prot)
	[148039]: SPELL_ACTION_CATEGORY + 0.11, // Sacred Shield (Holy)
	[65148]: SPELL_ACTION_CATEGORY + 0.12, // Sacred Shield (Absorb)
	[114039]: SPELL_ACTION_CATEGORY + 0.13, // Hand of Purity

	// Priest
	[48300]: SPELL_ACTION_CATEGORY + 0.11, // Devouring Plague
	[48125]: SPELL_ACTION_CATEGORY + 0.12, // Shadow Word: Pain
	[48160]: SPELL_ACTION_CATEGORY + 0.13, // Vampiric Touch
	[48135]: SPELL_ACTION_CATEGORY + 0.14, // Holy Fire
	[48123]: SPELL_ACTION_CATEGORY + 0.19, // Smite
	[48127]: SPELL_ACTION_CATEGORY + 0.2, // Mind Blast
	[48158]: SPELL_ACTION_CATEGORY + 0.3, // Shadow Word: Death
	[48156]: SPELL_ACTION_CATEGORY + 0.4, // Mind Flay

	// Rogue
	[6774]: MELEE_ACTION_CATEGORY + 0.1, // Slice and Dice
	[8647]: MELEE_ACTION_CATEGORY + 0.2, // Expose Armor
	[48672]: MELEE_ACTION_CATEGORY + 0.3, // Rupture
	[57993]: MELEE_ACTION_CATEGORY + 0.3, // Envenom
	[48668]: MELEE_ACTION_CATEGORY + 0.4, // Eviscerate
	[48666]: MELEE_ACTION_CATEGORY + 0.5, // Mutilate
	[48665]: MELEE_ACTION_CATEGORY + 0.6, // Mutilate (MH)
	[48664]: MELEE_ACTION_CATEGORY + 0.7, // Mutilate (OH)
	[48638]: MELEE_ACTION_CATEGORY + 0.5, // Sinister Strike
	[51723]: MELEE_ACTION_CATEGORY + 0.8, // Fan of Knives
	[57973]: SPELL_ACTION_CATEGORY + 0.1, // Deadly Poison
	[57968]: SPELL_ACTION_CATEGORY + 0.2, // Instant Poison

	// Shaman
	[8232]: 0.11, // Windfury Weapon
	[8024]: 0.12, // Flametongue Weapon
	[8033]: 0.12, // Frostbrand Weapon
	[17364]: MELEE_ACTION_CATEGORY + 0.1, // Stormstrike
	[60103]: MELEE_ACTION_CATEGORY + 0.2, // Lava Lash
	[49233]: SPELL_ACTION_CATEGORY + 0.21, // Flame Shock
	[49231]: SPELL_ACTION_CATEGORY + 0.22, // Earth Shock
	[49236]: SPELL_ACTION_CATEGORY + 0.23, // Frost Shock
	[60043]: SPELL_ACTION_CATEGORY + 0.31, // Lava Burst
	[49238]: SPELL_ACTION_CATEGORY + 0.32, // Lightning Bolt
	[49271]: SPELL_ACTION_CATEGORY + 0.33, // Chain Lightning
	[61657]: SPELL_ACTION_CATEGORY + 0.41, // Fire Nova
	[58734]: SPELL_ACTION_CATEGORY + 0.42, // Magma Totem
	[58704]: SPELL_ACTION_CATEGORY + 0.43, // Searing Totem
	[49281]: SPELL_ACTION_CATEGORY + 0.51, // Lightning Shield
	[49279]: SPELL_ACTION_CATEGORY + 0.52, // Lightning Shield (Proc)
	[2825]: DEFAULT_ACTION_CATEGORY + 0.1, // Bloodlust

	// Warlock
	[603]: SPELL_ACTION_CATEGORY + 0.01, // Curse of Doom
	[980]: SPELL_ACTION_CATEGORY + 0.02, // Curse of Agony
	[172]: SPELL_ACTION_CATEGORY + 0.1, // Corruption
	[48181]: SPELL_ACTION_CATEGORY + 0.2, // Haunt
	[30108]: SPELL_ACTION_CATEGORY + 0.3, // Unstable Affliction
	[348]: SPELL_ACTION_CATEGORY + 0.31, // Immolate
	[17962]: SPELL_ACTION_CATEGORY + 0.32, // Conflagrate
	[50796]: SPELL_ACTION_CATEGORY + 0.49, // Chaos Bolt
	[686]: SPELL_ACTION_CATEGORY + 0.5, // Shadow Bolt
	[29722]: SPELL_ACTION_CATEGORY + 0.51, // Incinerate
	[6353]: SPELL_ACTION_CATEGORY + 0.52, // Soul Fire
	[1120]: SPELL_ACTION_CATEGORY + 0.6, // Drain Soul
	[1454]: SPELL_ACTION_CATEGORY + 0.7, // Life Tap
	[59672]: SPELL_ACTION_CATEGORY + 0.8, // Metamorphosis
	[50589]: SPELL_ACTION_CATEGORY + 0.81, // Immolation Aura
	[47193]: SPELL_ACTION_CATEGORY + 0.82, // Demonic Empowerment

	// Mage
	[42842]: SPELL_ACTION_CATEGORY + 0.01, // Frostbolt
	[47610]: SPELL_ACTION_CATEGORY + 0.02, // Frostfire Bolt
	[42897]: SPELL_ACTION_CATEGORY + 0.02, // Arcane Blast
	[42833]: SPELL_ACTION_CATEGORY + 0.02, // Fireball
	[42859]: SPELL_ACTION_CATEGORY + 0.03, // Scorch
	[42891]: SPELL_ACTION_CATEGORY + 0.1, // Pyroblast
	[42846]: SPELL_ACTION_CATEGORY + 0.1, // Arcane Missiles
	[44572]: SPELL_ACTION_CATEGORY + 0.1, // Deep Freeze
	[44781]: SPELL_ACTION_CATEGORY + 0.2, // Arcane Barrage
	[42914]: SPELL_ACTION_CATEGORY + 0.2, // Ice Lance
	[55360]: SPELL_ACTION_CATEGORY + 0.2, // Living Bomb
	[55362]: SPELL_ACTION_CATEGORY + 0.21, // Living Bomb (Explosion)
	[12654]: SPELL_ACTION_CATEGORY + 0.3, // Ignite
	[12472]: SPELL_ACTION_CATEGORY + 0.4, // Icy Veins
	[11129]: SPELL_ACTION_CATEGORY + 0.4, // Combustion
	[12042]: SPELL_ACTION_CATEGORY + 0.4, // Arcane Power
	[11958]: SPELL_ACTION_CATEGORY + 0.41, // Cold Snap
	[12043]: SPELL_ACTION_CATEGORY + 0.41, // Presence of Mind
	[31687]: SPELL_ACTION_CATEGORY + 0.41, // Water Elemental
	[55342]: SPELL_ACTION_CATEGORY + 0.5, // Mirror Image
	[33312]: SPELL_ACTION_CATEGORY + 0.51, // Mana Gems
	[12051]: SPELL_ACTION_CATEGORY + 0.52, // Evocate
	[44401]: SPELL_ACTION_CATEGORY + 0.6, // Missile Barrage
	[44448]: SPELL_ACTION_CATEGORY + 0.6, // Hot Streak
	[44545]: SPELL_ACTION_CATEGORY + 0.6, // Fingers of Frost
	[44549]: SPELL_ACTION_CATEGORY + 0.61, // Brain Freeze
	[12536]: SPELL_ACTION_CATEGORY + 0.61, // Clearcasting

	// Warrior
	[47520]: 0.1, // Cleave
	[47450]: 0.1, // Heroic Strike
	[47475]: MELEE_ACTION_CATEGORY + 0.05, // Slam
	[23881]: MELEE_ACTION_CATEGORY + 0.1, // Bloodthirst
	[47486]: MELEE_ACTION_CATEGORY + 0.1, // Mortal Strike
	[30356]: MELEE_ACTION_CATEGORY + 0.1, // Shield Slam
	[47498]: MELEE_ACTION_CATEGORY + 0.21, // Devastate
	[47467]: MELEE_ACTION_CATEGORY + 0.22, // Sunder Armor
	[57823]: MELEE_ACTION_CATEGORY + 0.23, // Revenge
	[1680]: MELEE_ACTION_CATEGORY + 0.24, // Whirlwind
	[7384]: MELEE_ACTION_CATEGORY + 0.25, // Overpower
	[47471]: MELEE_ACTION_CATEGORY + 0.42, // Execute
	[12867]: SPELL_ACTION_CATEGORY + 0.51, // Deep Wounds
	[58874]: SPELL_ACTION_CATEGORY + 0.52, // Damage Shield
	[47296]: SPELL_ACTION_CATEGORY + 0.53, // Critical Block
	[46924]: SPELL_ACTION_CATEGORY + 0.61, // Bladestorm
	[2565]: SPELL_ACTION_CATEGORY + 0.62, // Shield Block
	[64382]: SPELL_ACTION_CATEGORY + 0.65, // Shattering Throw
	[71]: DEFAULT_ACTION_CATEGORY + 0.1, // Defensive Stance
	[2457]: DEFAULT_ACTION_CATEGORY + 0.1, // Battle Stance
	[2458]: DEFAULT_ACTION_CATEGORY + 0.1, // Berserker Stance

	// Death Knight
	[51425]: MELEE_ACTION_CATEGORY + 0.05, // Obliterate
	[55268]: MELEE_ACTION_CATEGORY + 0.1, // Frost strike
	[49930]: MELEE_ACTION_CATEGORY + 0.15, // Blood strike
	[50842]: MELEE_ACTION_CATEGORY + 0.2, // Pestilence
	[51411]: MELEE_ACTION_CATEGORY + 0.25, // Howling Blast
	[49895]: MELEE_ACTION_CATEGORY + 0.25, // Death Coil
	[49938]: MELEE_ACTION_CATEGORY + 0.25, // Death and Decay
	[63560]: MELEE_ACTION_CATEGORY + 0.25, // Ghoul Frenzy
	[50536]: MELEE_ACTION_CATEGORY + 0.25, // Unholy Blight
	[57623]: MELEE_ACTION_CATEGORY + 0.25, // HoW
	[59131]: MELEE_ACTION_CATEGORY + 0.3, // Icy touch
	[49921]: MELEE_ACTION_CATEGORY + 0.3, // Plague strike
	[51271]: MELEE_ACTION_CATEGORY + 0.35, // UA
	[45529]: MELEE_ACTION_CATEGORY + 0.35, // BT
	[47568]: MELEE_ACTION_CATEGORY + 0.35, // ERW
	[49206]: MELEE_ACTION_CATEGORY + 0.35, // Summon Gargoyle
	[46584]: MELEE_ACTION_CATEGORY + 0.35, // Raise Dead
	[55095]: MELEE_ACTION_CATEGORY + 0.4, // Frost Fever
	[55078]: MELEE_ACTION_CATEGORY + 0.4, // Blood Plague
	[49655]: MELEE_ACTION_CATEGORY + 0.4, // Wandering Plague
	[50401]: MELEE_ACTION_CATEGORY + 0.5, // Razor Frost
	[51460]: MELEE_ACTION_CATEGORY + 0.5, // Necrosis
	[50463]: MELEE_ACTION_CATEGORY + 0.5, // BCB
	[50689]: DEFAULT_ACTION_CATEGORY + 0.1, // Blood Presence
	[48263]: DEFAULT_ACTION_CATEGORY + 0.1, // Frost Presence
	[48265]: DEFAULT_ACTION_CATEGORY + 0.1, // Unholy Presence

	// Monk
	[120274]: 0.02, // Tiger Strikes (Main Hand)
	[120278]: 0.03, // Tiger Strikes (Off Hand)
	[100780]: MELEE_ACTION_CATEGORY + 0.01, // Jab
	[100787]: MELEE_ACTION_CATEGORY + 0.02, // Tiger Palm
	[100784]: MELEE_ACTION_CATEGORY + 0.03, // Blackout Kick
	[130320]: MELEE_ACTION_CATEGORY + 0.04, // Rising Sun Kick
	[113656]: MELEE_ACTION_CATEGORY + 0.05, // Fists of Fury (Cast)
	[117418]: MELEE_ACTION_CATEGORY + 0.06, // Fists of Fury (Hit)
	[101546]: MELEE_ACTION_CATEGORY + 0.07, // Spinning Crane Kick (Cast)
	[107270]: MELEE_ACTION_CATEGORY + 0.08, // Spinning Crane Kick (Hit)
	[116847]: MELEE_ACTION_CATEGORY + 0.07, // Rushing Jade Wind (Cast)
	[148187]: MELEE_ACTION_CATEGORY + 0.08, // Rushing Jade Wind (Hit)
	[115098]: SPELL_ACTION_CATEGORY + 0.01, // Chi Wave
	[132467]: SPELL_ACTION_CATEGORY + 0.011, // Chi Wave (Damage)
	[132463]: SPELL_ACTION_CATEGORY + 0.012, // Chi Wave (Heal)
	[124098]: SPELL_ACTION_CATEGORY + 0.01, // Zen Sphere (Damage)
	[124081]: SPELL_ACTION_CATEGORY + 0.011, // Zen Sphere (Heal)
	[125033]: SPELL_ACTION_CATEGORY + 0.011, // Zen Sphere: Detonate (Damage)
	[124101]: SPELL_ACTION_CATEGORY + 0.011, // Zen Sphere: Detonate (Heal)
	[123986]: SPELL_ACTION_CATEGORY + 0.01, // Chi Burst
	[148135]: SPELL_ACTION_CATEGORY + 0.011, // Chi Burst (Damage)
	[130654]: SPELL_ACTION_CATEGORY + 0.012, // Chi Burst (Heal)
	[116740]: SPELL_ACTION_CATEGORY + 0.02, // Tigereye Brew
	[115399]: SPELL_ACTION_CATEGORY + 0.03, // Chi Brew
	[115288]: SPELL_ACTION_CATEGORY + 0.04, // Energizing Brew
	[126456]: SPELL_ACTION_CATEGORY + 0.05, // Fortifying Brew
	[123904]: SPELL_ACTION_CATEGORY + 0.06, // Invoke Xuen, the White Tiger
	[115008]: SPELL_ACTION_CATEGORY + 0.06, // Chi Torpedo

	// Generic
	[53307]: SPELL_ACTION_CATEGORY + 0.931, // Thorns
	[54043]: SPELL_ACTION_CATEGORY + 0.932, // Retribution Aura
	[54758]: SPELL_ACTION_CATEGORY + 0.933, // Hyperspeed Acceleration
	[42641]: SPELL_ACTION_CATEGORY + 0.941, // Sapper
	[40536]: SPELL_ACTION_CATEGORY + 0.942, // Explosive Decoy
	[41119]: SPELL_ACTION_CATEGORY + 0.943, // Saronite Bomb
	[40771]: SPELL_ACTION_CATEGORY + 0.944, // Cobalt Frag Bomb

	// Souldrinker - to pair up the damage part with the healing
	[109828]: SPELL_ACTION_CATEGORY + 0.945, // Drain Life - LFR
	[108022]: SPELL_ACTION_CATEGORY + 0.946, // Drain Life - Normal
	[109831]: SPELL_ACTION_CATEGORY + 0.947, // Drain Life - Heroic

	// No'Kaled - to pair up the different spells it can proc
	[109871]: SPELL_ACTION_CATEGORY + 0.948, // Flameblast - LFR
	[109869]: SPELL_ACTION_CATEGORY + 0.949, // Iceblast - LFR
	[109867]: SPELL_ACTION_CATEGORY + 0.95, // Shadowblast - LFR
	[107785]: SPELL_ACTION_CATEGORY + 0.951, // Flameblast - Normal
	[107789]: SPELL_ACTION_CATEGORY + 0.952, // Iceblast - Normal
	[107787]: SPELL_ACTION_CATEGORY + 0.953, // Shadowblast - Normal
	[109872]: SPELL_ACTION_CATEGORY + 0.954, // Flameblast - Heroic
	[109870]: SPELL_ACTION_CATEGORY + 0.955, // Iceblast - Heroic
	[109868]: SPELL_ACTION_CATEGORY + 0.956, // Shadowblast - Heroic
};

const idsToGroupForRotation: Array<number> = [
	6774, // Slice and Dice
	8647, // Expose Armor
	48668, // Eviscerate
	48672, // Rupture
	51690, // Killing Spree
	57993, // Envenom
];

const percentageResources: Array<ResourceType> = [ResourceType.ResourceTypeHealth, ResourceType.ResourceTypeMana];
