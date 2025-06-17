import { hasTouch } from '../shared/bootstrap_overrides';
import { SimRequest } from '../worker/types';
import { getBrowserLanguageCode, setLanguageCode } from './constants/lang';
import { CURRENT_PHASE, LOCAL_STORAGE_PREFIX } from './constants/other';
import { Encounter } from './encounter';
import { getCurrentLang, setCurrentLang } from './locale_service';
import { Player, UnitMetadata } from './player';
import {
	BulkSettings,
	BulkSimCombosRequest,
	BulkSimCombosResult,
	BulkSimRequest,
	BulkSimResult,
	ComputeStatsRequest,
	ErrorOutcome,
	ErrorOutcomeType,
	Raid as RaidProto,
	RaidSimRequest,
	RaidSimResult,
	SimOptions,
	SimType,
	StatWeightsRequest,
	StatWeightsResult,
} from './proto/api.js';
import {
	ArmorType,
	Faction,
	Profession,
	PseudoStat,
	RangedWeaponType,
	Spec,
	Stat,
	UnitReference,
	UnitReference_Type as UnitType,
	WeaponType,
} from './proto/common.js';
import { Consumable, SimDatabase } from './proto/db';
import { SpellEffect } from './proto/spell';
import { DatabaseFilters, RaidFilterOption, SimSettings as SimSettingsProto, SourceFilterOption } from './proto/ui.js';
import { Database } from './proto_utils/database.js';
import { SimResult } from './proto_utils/sim_result.js';
import { Raid } from './raid.js';
import { runConcurrentSim, runConcurrentStatWeights } from './sim_concurrent';
import { RequestTypes, SimSignalManager } from './sim_signal_manager';
import { EventID, TypedEvent } from './typed_event.js';
import { getEnumValues, noop } from './utils.js';
import { generateRequestId, WorkerPool, WorkerProgressCallback } from './worker_pool.js';

export type RaidSimData = {
	request: RaidSimRequest;
	result: RaidSimResult;
};

export type StatWeightsData = {
	request: StatWeightsRequest;
	result: StatWeightsResult;
};

interface SimProps {
	// The type of sim. Default `SimType.SimTypeIndividual`
	type?: SimType;
}

const WASM_CONCURRENCY_STORAGE_KEY = `${LOCAL_STORAGE_PREFIX}_wasmconcurrency`;

// Core Sim module which deals only with api types, no UI-related stuff.
export class Sim {
	private readonly workerPool: WorkerPool;

	iterations = 12500;

	private phase: number = CURRENT_PHASE;
	private faction: Faction = Faction.Alliance;
	private fixedRngSeed = 0;
	private filters: DatabaseFilters = Sim.defaultFilters();
	private showDamageMetrics = true;
	private showThreatMetrics = false;
	private showHealingMetrics = false;
	private showExperimental = false;
	private wasmConcurrency = 0;
	private showQuickSwap = false;
	private showEPValues = false;
	private useCustomEPValues = false;
	private useSoftCapBreakpoints = true;
	private language = '';

	readonly type: SimType;
	readonly raid: Raid;
	readonly encounter: Encounter;

	private db_: Database | null = null;

	readonly iterationsChangeEmitter = new TypedEvent<void>();
	readonly phaseChangeEmitter = new TypedEvent<void>();
	readonly factionChangeEmitter = new TypedEvent<void>();
	readonly fixedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly lastUsedRngSeedChangeEmitter = new TypedEvent<void>();
	readonly filtersChangeEmitter = new TypedEvent<void>();
	readonly showDamageMetricsChangeEmitter = new TypedEvent<void>();
	readonly showThreatMetricsChangeEmitter = new TypedEvent<void>();
	readonly showHealingMetricsChangeEmitter = new TypedEvent<void>();
	readonly showExperimentalChangeEmitter = new TypedEvent<void>();
	readonly wasmConcurrencyChangeEmitter = new TypedEvent<void>();
	readonly showQuickSwapChangeEmitter = new TypedEvent<void>();
	readonly showEPValuesChangeEmitter = new TypedEvent<void>();
	readonly useCustomEPValuesChangeEmitter = new TypedEvent<void>();
	readonly useSoftCapBreakpointsChangeEmitter = new TypedEvent<void>();
	readonly languageChangeEmitter = new TypedEvent<void>();
	readonly crashEmitter = new TypedEvent<SimError>();

	// Emits when any of the settings change (but not the raid / encounter).
	readonly settingsChangeEmitter: TypedEvent<void>;

	// Emits when any player, target, or pet has metadata changes (spells or auras).
	readonly unitMetadataEmitter = new TypedEvent<void>('UnitMetadata');

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	// Fires when a raid sim API call completes.
	readonly simResultEmitter = new TypedEvent<SimResult>();

	// Fires when a bulk sim API call starts.
	readonly bulkSimStartEmitter = new TypedEvent<BulkSimRequest>();
	// Fires when a bulk sim API call completes..
	readonly bulkSimResultEmitter = new TypedEvent<BulkSimResult>();

	private readonly _initPromise: Promise<any>;
	private lastUsedRngSeed = 0;

	// These callbacks are needed so we can apply BuffBot modifications automatically before sending requests.
	private modifyRaidProto: (raidProto: RaidProto) => void = noop;

	readonly signalManager: SimSignalManager;

	constructor({ type }: SimProps = {}) {
		this.type = type ?? SimType.SimTypeIndividual;

		this.workerPool = new WorkerPool(1);
		this.wasmConcurrencyChangeEmitter.on(async () => {
			// Prevent using worker concurrency when not running wasm. Local sim has native threading.
			if (await this.workerPool.isWasm()) {
				const nWorker = Math.max(1, Math.min(this.wasmConcurrency, navigator.hardwareConcurrency));
				this.workerPool.setNumWorkers(nWorker);
			}
		});

		let wasmConcurrencySetting = parseInt(window.localStorage.getItem(WASM_CONCURRENCY_STORAGE_KEY) ?? 'NaN');
		if (isNaN(wasmConcurrencySetting)) {
			wasmConcurrencySetting = 0;
			// Set a default worker count if env supports multiple threads. Should not be too high as to be safe for all situations.
			// TODO: Set based on browser/engine? E.g. Firefox has significant RAM and CPU usage per worker while Chrome can run many without a downside.
			if (navigator.hardwareConcurrency > 1) {
				wasmConcurrencySetting = Math.min(4, Math.floor(navigator.hardwareConcurrency / 2));
			}
		}
		this.setWasmConcurrency(TypedEvent.nextEventID(), wasmConcurrencySetting);

		this.signalManager = new SimSignalManager();

		this._initPromise = Database.get().then(db => {
			this.db_ = db;
		});

		this.raid = new Raid(this);
		this.encounter = new Encounter(this);

		this.settingsChangeEmitter = TypedEvent.onAny([
			this.iterationsChangeEmitter,
			this.phaseChangeEmitter,
			this.fixedRngSeedChangeEmitter,
			this.filtersChangeEmitter,
			this.showDamageMetricsChangeEmitter,
			this.showThreatMetricsChangeEmitter,
			this.showHealingMetricsChangeEmitter,
			this.showExperimentalChangeEmitter,
			this.wasmConcurrencyChangeEmitter,
			this.showQuickSwapChangeEmitter,
			this.showEPValuesChangeEmitter,
			this.useCustomEPValuesChangeEmitter,
			this.useSoftCapBreakpointsChangeEmitter,
			this.languageChangeEmitter,
		]);

		this.changeEmitter = TypedEvent.onAny([this.settingsChangeEmitter, this.raid.changeEmitter, this.encounter.changeEmitter]);

		TypedEvent.onAny([this.raid.changeEmitter, this.encounter.changeEmitter]).on(eventID => this.updateCharacterStats(eventID));

		this.language = getCurrentLang();
	}

	waitForInit(): Promise<void> {
		return this._initPromise;
	}

	/**
	 * Check if workers are running wasm.
	 * @returns true if workers are running wasm.
	 */
	isWasm() {
		return this.workerPool.isWasm();
	}

	/**
	 * Whether the current environment should use wasm/worker concurrency methods.
	 * @returns true if running wasm workers and concurrency setting is active.
	 */
	private async shouldUseWasmConcurrency() {
		return (await this.isWasm()) && this.getWasmConcurrency() >= 2 && this.workerPool.getNumWorkers() >= 2;
	}

	get db(): Database {
		return this.db_!;
	}

	setModifyRaidProto(newModFn: (raidProto: RaidProto) => void) {
		this.modifyRaidProto = newModFn;
	}

	getModifiedRaidProto(): RaidProto {
		const raidProto = this.raid.toProto(false, true);
		this.modifyRaidProto(raidProto);

		// Remove any inactive meta gems, since the backend doesn't have its own validation.
		raidProto.parties.forEach(party => {
			party.players.forEach(player => {
				if (!player.equipment) {
					return;
				}

				let gear = this.db.lookupEquipmentSpec(player.equipment);
				let gearChanged = false;

				const isBlacksmith = [player.profession1, player.profession2].includes(Profession.Blacksmithing);

				// Disable meta gem if inactive.
				if (gear.hasInactiveMetaGem(isBlacksmith)) {
					gear = gear.withoutMetaGem();
					gearChanged = true;
				}

				// Remove bonus sockets if not blacksmith.
				if (!isBlacksmith) {
					gear = gear.withoutBlacksmithSockets();
					gearChanged = true;
				}

				if (gearChanged) {
					player.equipment = gear.asSpec();
				}

				// Include consumables in the player db
				const pdb = player.database!;

				const newConsumables: Consumable[] = [];
				const newSpellEffects: SpellEffect[] = [];
				const seenConsumableIds = new Set<number>();
				const seenEffectIds = new Set<number>();
				Object.values(player.consumables ?? []).forEach((cid: number) => {
					if (!cid || seenConsumableIds.has(cid)) return;
					const consume = this.db.getConsumable(cid);
					if (!consume) return;
					seenConsumableIds.add(consume.id);
					newConsumables.push(consume);
					for (const eid of consume.effectIds) {
						if (seenEffectIds.has(eid)) continue;
						const effect = this.db.getSpellEffect(eid);
						if (!effect) continue;

						seenEffectIds.add(effect.id);
						newSpellEffects.push(effect);
					}
				});

				// swap in the fresh arrays
				pdb.consumables = newConsumables;
				pdb.spellEffects = newSpellEffects;
				player.database = pdb;
			});
		});

		return raidProto;
	}

	makeRaidSimRequest(debug: boolean): RaidSimRequest {
		const raid = this.getModifiedRaidProto();
		const encounter = this.encounter.toProto();

		// TODO: remove any replenishment from sim request here? probably makes more sense to do it inside the sim to protect against accidents

		return RaidSimRequest.create({
			requestId: generateRequestId(SimRequest.raidSimAsync),
			type: this.type,
			raid: raid,
			encounter: encounter,
			simOptions: SimOptions.create({
				iterations: debug ? 1 : this.getIterations(),
				randomSeed: BigInt(this.nextRngSeed()),
				debugFirstIteration: true,
			}),
		});
	}

	async runBulkSim(bulkSettings: BulkSettings, bulkItemsDb: SimDatabase, onProgress: WorkerProgressCallback): Promise<BulkSimResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.targets.length < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = BulkSimRequest.create({
			baseSettings: this.makeRaidSimRequest(false),
			bulkSettings: bulkSettings,
		});

		if (request.baseSettings != null && request.baseSettings.simOptions != null) {
			request.baseSettings.simOptions.debugFirstIteration = false;
		}

		if (!request.baseSettings?.raid || request.baseSettings?.raid?.parties.length == 0 || request.baseSettings?.raid?.parties[0].players.length == 0) {
			throw new Error('Raid must contain exactly 1 player for bulk sim.');
		}

		// Attach the extra database to the player.
		const playerDatabase = request.baseSettings.raid.parties[0].players[0].database!;
		playerDatabase.items.push(...bulkItemsDb.items);
		playerDatabase.enchants.push(...bulkItemsDb.enchants);
		playerDatabase.gems.push(...bulkItemsDb.gems);
		playerDatabase.reforgeStats.push(...bulkItemsDb.reforgeStats);
		playerDatabase.itemEffectRandPropPoints.push(...bulkItemsDb.itemEffectRandPropPoints);
		playerDatabase.randomSuffixes.push(...bulkItemsDb.randomSuffixes);
		playerDatabase.consumables.push(...bulkItemsDb.consumables);
		playerDatabase.spellEffects.push(...bulkItemsDb.spellEffects);

		this.bulkSimStartEmitter.emit(TypedEvent.nextEventID(), request);

		const signals = this.signalManager.registerRunning(RequestTypes.BulkSim);
		try {
			const result = await this.workerPool.bulkSimAsync(request, onProgress, signals);

			if (result.error) {
				if (result.error.type != ErrorOutcomeType.ErrorOutcomeError) return result;
				throw new SimError(result.error.message);
			}

			this.bulkSimResultEmitter.emit(TypedEvent.nextEventID(), result);
			return result;
		} catch (error) {
			if (error instanceof SimError) throw error;
			console.log(error);
			throw new Error('Something went wrong running your raid sim. Reload the page and try again.');
		} finally {
			this.signalManager.unregisterRunning(signals);
		}
	}

	async calculateBulkCombinations(bulkSettings: BulkSettings, bulkItemsDb: SimDatabase): Promise<BulkSimCombosResult | null> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.targets.length < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		const request = BulkSimCombosRequest.create({
			baseSettings: this.makeRaidSimRequest(false),
			bulkSettings: bulkSettings,
		});

		if (request.baseSettings != null && request.baseSettings.simOptions != null) {
			request.baseSettings.simOptions.debugFirstIteration = false;
		}

		if (!request.baseSettings?.raid || request.baseSettings?.raid?.parties.length == 0 || request.baseSettings?.raid?.parties[0].players.length == 0) {
			throw new Error('Raid must contain exactly 1 player for bulk sim.');
		}

		// Attach the extra database to the player.
		const playerDatabase = request.baseSettings.raid.parties[0].players[0].database!;
		playerDatabase.items.push(...bulkItemsDb.items);
		playerDatabase.enchants.push(...bulkItemsDb.enchants);
		playerDatabase.gems.push(...bulkItemsDb.gems);
		playerDatabase.reforgeStats.push(...bulkItemsDb.reforgeStats);
		playerDatabase.itemEffectRandPropPoints.push(...bulkItemsDb.itemEffectRandPropPoints);
		playerDatabase.randomSuffixes.push(...bulkItemsDb.randomSuffixes);

		this.bulkSimStartEmitter.emit(TypedEvent.nextEventID(), request);

		try {
			const result = await this.workerPool.bulkSimCombosAsync(request);
			if (result.errorResult != '') {
				throw new SimError(result.errorResult);
			}
			return result;
		} catch (error) {
			if (error instanceof SimError) throw error;
			console.log(error);
			throw new Error('Something went wrong running your raid sim. Reload the page and try again.');
		}
	}

	async runRaidSim(eventID: EventID, onProgress: WorkerProgressCallback): Promise<SimResult | ErrorOutcome> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.targets.length < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		const signals = this.signalManager.registerRunning(RequestTypes.RaidSim);
		try {
			await this.waitForInit();

			const request = this.makeRaidSimRequest(false);

			let result;
			// Only use worker base concurrency when running wasm. Local sim has native threading.
			if (await this.shouldUseWasmConcurrency()) {
				result = await runConcurrentSim(request, this.workerPool, onProgress, signals);
			} else {
				result = await this.workerPool.raidSimAsync(request, onProgress, signals);
			}

			if (result.error) {
				if (result.error.type != ErrorOutcomeType.ErrorOutcomeError) return result.error;
				throw new SimError(result.error.message);
			}
			const simResult = await SimResult.makeNew(request, result);
			this.simResultEmitter.emit(eventID, simResult);
			return simResult;
		} catch (error) {
			if (error instanceof SimError) throw error;
			console.error(error);
			throw new Error('Something went wrong running your raid sim. Reload the page and try again.');
		} finally {
			this.signalManager.unregisterRunning(signals);
		}
	}

	async runRaidSimWithLogs(eventID: EventID): Promise<SimResult | null> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.targets.length < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		const signals = this.signalManager.registerRunning(RequestTypes.RaidSim);
		try {
			await this.waitForInit();

			const request = this.makeRaidSimRequest(true);
			const result = await this.workerPool.raidSimAsync(request, noop, signals);
			if (result.error) {
				throw new SimError(result.error.message);
			}
			const simResult = await SimResult.makeNew(request, result);
			this.simResultEmitter.emit(eventID, simResult);
			return simResult;
		} catch (error) {
			if (error instanceof SimError) throw error;
			console.error(error);
			throw new Error('Something went wrong running your raid sim. Reload the page and try again.');
		} finally {
			this.signalManager.unregisterRunning(signals);
		}
	}

	// This should be invoked internally whenever stats might have changed.
	async updateCharacterStats(eventID: EventID) {
		if (eventID == 0) {
			// Skip the first event ID because it interferes with the loaded stats.
			return;
		}
		eventID = TypedEvent.nextEventID();

		await this.waitForInit();
		// Capture the current players so we avoid issues if something changes while
		// request is in-flight.

		const players = this.raid.getPlayers();
		const req = ComputeStatsRequest.create({
			raid: this.getModifiedRaidProto(),
			encounter: this.encounter.toProto(),
		});
		const result = await this.workerPool.computeStats(req);
		if (result.errorResult != '') {
			this.crashEmitter.emit(eventID, new SimError(result.errorResult));
			return;
		}

		TypedEvent.freezeAllAndDo(async () => {
			const playerUpdatePromises = result
				.raidStats!.parties.map((partyStats, partyIndex) =>
					partyStats.players.map((playerStats, playerIndex) => {
						const player = players[partyIndex * 5 + playerIndex];
						if (player) {
							player.setCurrentStats(eventID, playerStats);
							return player.updateMetadata();
						} else {
							return null;
						}
					}),
				)
				.flat()
				.filter(p => p != null) as Array<Promise<boolean>>;

			const targetUpdatePromise = this.encounter.targetsMetadata.update(result.encounterStats!.targets.map(t => t.metadata!));

			const anyUpdates = await Promise.all(playerUpdatePromises.concat([targetUpdatePromise]));
			if (anyUpdates.some(v => v)) {
				this.unitMetadataEmitter.emit(eventID);
			}
		});
	}

	async statWeights(
		player: Player<any>,
		epStats: Array<Stat>,
		epPseudoStats: Array<PseudoStat>,
		epReferenceStat: Stat,
		onProgress: WorkerProgressCallback,
	): Promise<StatWeightsResult> {
		if (this.raid.isEmpty()) {
			throw new Error('Raid is empty! Try adding some players first.');
		} else if (this.encounter.targets.length < 1) {
			throw new Error('Encounter has no targets! Try adding some targets first.');
		}

		await this.waitForInit();

		if (player.getParty() == null) {
			console.warn('Trying to get stat weights without a party!');
			return StatWeightsResult.create();
		} else {
			const tanks = this.raid
				.getTanks()
				.map(tank => tank.index)
				.includes(player.getRaidIndex())
				? [UnitReference.create({ type: UnitType.Player, index: 0 })]
				: [];
			const request = StatWeightsRequest.create({
				player: player.toProto(false, true),
				raidBuffs: this.raid.getBuffs(),
				partyBuffs: player.getParty()!.getBuffs(),
				debuffs: this.raid.getDebuffs(),
				encounter: this.encounter.toProto(),
				simOptions: SimOptions.create({
					iterations: this.getIterations(),
					randomSeed: BigInt(this.nextRngSeed()),
					debug: false,
				}),
				tanks: tanks,

				statsToWeigh: epStats,
				pseudoStatsToWeigh: epPseudoStats,
				epReferenceStat: epReferenceStat,
			});

			const signals = this.signalManager.registerRunning(RequestTypes.StatWeights);
			try {
				let result: StatWeightsResult;
				// Only use worker based concurrency when running wasm.
				if (await this.shouldUseWasmConcurrency()) {
					result = await runConcurrentStatWeights(request, this.workerPool, onProgress, signals);
				} else {
					result = await this.workerPool.statWeightsAsync(request, onProgress, signals);
				}
				if (result.error) {
					if (result.error.type != ErrorOutcomeType.ErrorOutcomeError) return result;
					throw new SimError(result.error.message);
				}
				return result;
			} catch (error) {
				if (error instanceof SimError) throw error;
				console.error(error);
				throw new Error('Something went wrong calculating your stat weights. Reload the page and try again.');
			} finally {
				this.signalManager.unregisterRunning(signals);
			}
		}
	}

	getUnitMetadata(ref: UnitReference | undefined, contextPlayer: Player<any> | null, defaultRef: UnitReference): UnitMetadata | undefined {
		if (!ref || ref.type == UnitType.Unknown) {
			return this.getUnitMetadata(defaultRef, contextPlayer, defaultRef);
		} else if (ref.type == UnitType.Player) {
			return this.raid.getPlayerFromUnitReference(ref)?.getMetadata();
		} else if (ref.type == UnitType.Target) {
			return this.encounter.targetsMetadata.asList()[ref.index];
		} else if (ref.type == UnitType.Pet) {
			const owner = this.raid.getPlayerFromUnitReference(ref.owner, contextPlayer);
			if (owner) {
				return owner.getPetMetadatas().asList()[ref.index];
			} else {
				return undefined;
			}
		} else if (ref.type == UnitType.Self) {
			return contextPlayer?.getMetadata();
		} else if (ref.type == UnitType.CurrentTarget) {
			return this.encounter.targetsMetadata.asList()[0];
		}
		return undefined;
	}

	getPhase(): number {
		return this.phase;
	}
	setPhase(eventID: EventID, newPhase: number) {
		if (newPhase != this.phase) {
			this.phase = newPhase;
			this.phaseChangeEmitter.emit(eventID);
		}
	}

	getFaction(): Faction {
		return this.faction;
	}
	setFaction(eventID: EventID, newFaction: Faction) {
		if (newFaction != this.faction && !!newFaction) {
			this.faction = newFaction;
			this.factionChangeEmitter.emit(eventID);
		}
	}

	getFixedRngSeed(): number {
		return this.fixedRngSeed;
	}
	setFixedRngSeed(eventID: EventID, newFixedRngSeed: number) {
		if (newFixedRngSeed != this.fixedRngSeed) {
			this.fixedRngSeed = newFixedRngSeed;
			this.fixedRngSeedChangeEmitter.emit(eventID);
		}
	}

	static MAX_RNG_SEED = Math.pow(2, 32) - 1;
	private nextRngSeed(): number {
		let rngSeed = 0;
		if (this.fixedRngSeed) {
			rngSeed = this.fixedRngSeed;
		} else {
			rngSeed = Math.floor(Math.random() * Sim.MAX_RNG_SEED);
		}

		this.lastUsedRngSeed = rngSeed;
		this.lastUsedRngSeedChangeEmitter.emit(TypedEvent.nextEventID());
		return rngSeed;
	}
	getLastUsedRngSeed(): number {
		return this.lastUsedRngSeed;
	}

	getFilters(): DatabaseFilters {
		// Make a defensive copy
		return DatabaseFilters.clone(this.filters);
	}
	setFilters(eventID: EventID, newFilters: DatabaseFilters) {
		if (DatabaseFilters.equals(newFilters, this.filters)) {
			return;
		}

		// Make a defensive copy
		this.filters = DatabaseFilters.clone(newFilters);
		this.filtersChangeEmitter.emit(eventID);
	}

	getShowDamageMetrics(): boolean {
		return this.showDamageMetrics;
	}
	setShowDamageMetrics(eventID: EventID, newShowDamageMetrics: boolean) {
		if (newShowDamageMetrics != this.showDamageMetrics) {
			this.showDamageMetrics = newShowDamageMetrics;
			this.showDamageMetricsChangeEmitter.emit(eventID);
		}
	}

	getShowThreatMetrics(): boolean {
		return this.showThreatMetrics;
	}
	setShowThreatMetrics(eventID: EventID, newShowThreatMetrics: boolean) {
		if (newShowThreatMetrics != this.showThreatMetrics) {
			this.showThreatMetrics = newShowThreatMetrics;
			this.showThreatMetricsChangeEmitter.emit(eventID);
		}
	}

	getShowHealingMetrics(): boolean {
		return (
			this.showHealingMetrics ||
			(this.showThreatMetrics &&
				[Spec.SpecBloodDeathKnight, Spec.SpecGuardianDruid, Spec.SpecBrewmasterMonk, Spec.SpecProtectionPaladin].includes(
					this.raid.getPlayer(0)?.playerSpec.specID,
				))
		);
	}
	setShowHealingMetrics(eventID: EventID, newShowHealingMetrics: boolean) {
		if (newShowHealingMetrics != this.showHealingMetrics) {
			this.showHealingMetrics = newShowHealingMetrics;
			this.showHealingMetricsChangeEmitter.emit(eventID);
		}
	}

	getShowExperimental(): boolean {
		return this.showExperimental;
	}
	setShowExperimental(eventID: EventID, newShowExperimental: boolean) {
		if (newShowExperimental != this.showExperimental) {
			this.showExperimental = newShowExperimental;
			this.showExperimentalChangeEmitter.emit(eventID);
		}
	}

	getWasmConcurrency(): number {
		return this.wasmConcurrency;
	}
	setWasmConcurrency(eventID: EventID, newWasmConcurrency: number) {
		if (newWasmConcurrency != this.wasmConcurrency) {
			this.wasmConcurrency = newWasmConcurrency;
			window.localStorage.setItem(WASM_CONCURRENCY_STORAGE_KEY, newWasmConcurrency.toString());
			this.wasmConcurrencyChangeEmitter.emit(eventID);
		}
	}

	getShowQuickSwap(): boolean {
		return !hasTouch() && this.showQuickSwap;
	}
	setShowQuickSwap(eventID: EventID, newShowQuickSwap: boolean) {
		if (newShowQuickSwap != this.showQuickSwap) {
			this.showQuickSwap = newShowQuickSwap;
			this.showQuickSwapChangeEmitter.emit(eventID);
		}
	}

	getShowEPValues(): boolean {
		return this.showEPValues;
	}
	setShowEPValues(eventID: EventID, newShowEPValues: boolean) {
		if (newShowEPValues != this.showEPValues) {
			this.showEPValues = newShowEPValues;
			this.showEPValuesChangeEmitter.emit(eventID);
		}
	}

	getUseCustomEPValues(): boolean {
		return this.useCustomEPValues;
	}
	setUseCustomEPValues(eventID: EventID, newUseCustomEPValues: boolean) {
		if (newUseCustomEPValues !== this.useCustomEPValues) {
			this.useCustomEPValues = newUseCustomEPValues;
			this.useCustomEPValuesChangeEmitter.emit(eventID);
		}
	}

	getUseSoftCapBreakpoints(): boolean {
		return this.useSoftCapBreakpoints;
	}
	setUseSoftCapBreakpoints(eventID: EventID, newUseSoftCapBreakpoints: boolean) {
		if (newUseSoftCapBreakpoints !== this.useSoftCapBreakpoints) {
			this.useSoftCapBreakpoints = newUseSoftCapBreakpoints;
			this.useSoftCapBreakpointsChangeEmitter.emit(eventID);
		}
	}

	getLanguage(): string {
		return this.language;
	}
	setLanguage(eventID: EventID, newLanguage: string) {
		newLanguage = newLanguage || getBrowserLanguageCode();
		if (newLanguage != this.language) {
			this.language = newLanguage;
			setCurrentLang(this.language);
			setLanguageCode(this.language);
			this.languageChangeEmitter.emit(eventID);
		}
	}

	getIterations(): number {
		return this.iterations;
	}
	setIterations(eventID: EventID, newIterations: number) {
		if (newIterations != this.iterations) {
			this.iterations = newIterations;
			this.iterationsChangeEmitter.emit(eventID);
		}
	}

	static readonly ALL_ARMOR_TYPES = (getEnumValues(ArmorType) as Array<ArmorType>).filter(v => v != 0);
	static readonly ALL_WEAPON_TYPES = (getEnumValues(WeaponType) as Array<WeaponType>).filter(v => v != 0);
	static readonly ALL_RANGED_WEAPON_TYPES = (getEnumValues(RangedWeaponType) as Array<RangedWeaponType>).filter(v => v != 0);
	static readonly ALL_SOURCES = (getEnumValues(SourceFilterOption) as Array<SourceFilterOption>).filter(v => v != 0);
	static readonly ALL_RAIDS = (getEnumValues(RaidFilterOption) as Array<RaidFilterOption>).filter(v => v != 0);

	toProto(): SimSettingsProto {
		const filters = this.getFilters();
		if (filters.armorTypes.length == Sim.ALL_ARMOR_TYPES.length) {
			filters.armorTypes = [];
		}
		if (filters.weaponTypes.length == Sim.ALL_WEAPON_TYPES.length) {
			filters.weaponTypes = [];
		}
		if (filters.rangedWeaponTypes.length == Sim.ALL_RANGED_WEAPON_TYPES.length) {
			filters.rangedWeaponTypes = [];
		}
		if (filters.sources.length == Sim.ALL_SOURCES.length) {
			filters.sources = [];
		}
		if (filters.raids.length == Sim.ALL_RAIDS.length) {
			filters.raids = [];
		}

		return SimSettingsProto.create({
			iterations: this.getIterations(),
			phase: this.getPhase(),
			fixedRngSeed: BigInt(this.getFixedRngSeed()),
			showDamageMetrics: this.getShowDamageMetrics(),
			showThreatMetrics: this.getShowThreatMetrics(),
			showHealingMetrics: this.getShowHealingMetrics(),
			showExperimental: this.getShowExperimental(),
			showQuickSwap: this.getShowQuickSwap(),
			showEpValues: this.getShowEPValues(),
			useCustomEpValues: this.getUseCustomEPValues(),
			useSoftCapBreakpoints: this.getUseSoftCapBreakpoints(),
			language: this.getLanguage(),
			faction: this.getFaction(),
			filters: filters,
		});
	}

	fromProto(eventID: EventID, proto: SimSettingsProto) {
		TypedEvent.freezeAllAndDo(() => {
			this.setIterations(eventID, proto.iterations || 12500);
			this.setPhase(eventID, proto.phase || CURRENT_PHASE);
			this.setFixedRngSeed(eventID, Number(proto.fixedRngSeed));
			this.setShowDamageMetrics(eventID, proto.showDamageMetrics);
			this.setShowThreatMetrics(eventID, proto.showThreatMetrics);
			this.setShowHealingMetrics(eventID, proto.showHealingMetrics);
			this.setShowExperimental(eventID, proto.showExperimental);
			this.setShowQuickSwap(eventID, proto.showQuickSwap);
			this.setShowEPValues(eventID, proto.showEpValues);
			this.setUseCustomEPValues(eventID, proto.useCustomEpValues);
			this.setUseSoftCapBreakpoints(eventID, proto.useSoftCapBreakpoints);
			this.setLanguage(eventID, proto.language);
			this.setFaction(eventID, proto.faction || Faction.Alliance);

			const filters = proto.filters || Sim.defaultFilters();
			if (filters.armorTypes.length == 0) {
				if (this.type == SimType.SimTypeIndividual) {
					// For Individual sims, by default only show the class's default armor type because of armor specialization
					filters.armorTypes = [this.raid.getActivePlayers()[0].getPlayerClass().armorTypes[0]];
				} else {
					filters.armorTypes = Sim.ALL_ARMOR_TYPES.slice();
				}
			}
			if (filters.weaponTypes.length == 0) {
				filters.weaponTypes = Sim.ALL_WEAPON_TYPES.slice();
			}
			if (filters.rangedWeaponTypes.length == 0) {
				filters.rangedWeaponTypes = Sim.ALL_RANGED_WEAPON_TYPES.slice();
			}
			if (filters.sources.length == 0) {
				filters.sources = Sim.ALL_SOURCES.slice();
			}
			if (filters.raids.length == 0) {
				filters.raids = Sim.ALL_RAIDS.slice();
			}
			this.setFilters(eventID, filters);
		});
	}

	applyDefaults(eventID: EventID, isTankSim: boolean, isHealingSim: boolean) {
		this.fromProto(
			eventID,
			SimSettingsProto.create({
				iterations: 12500,
				phase: CURRENT_PHASE,
				faction: Faction.Alliance,
				showDamageMetrics: !isHealingSim,
				showThreatMetrics: isTankSim,
				showHealingMetrics: isHealingSim,
				language: this.getLanguage(), // Don't change language.
				filters: Sim.defaultFilters(),
				showEpValues: false,
				useSoftCapBreakpoints: true,
			}),
		);
	}

	static defaultFilters(): DatabaseFilters {
		return DatabaseFilters.create({
			oneHandedWeapons: true,
			twoHandedWeapons: true,
		});
	}
}

export class SimError extends Error {
	readonly errorStr: string;

	constructor(errorStr: string) {
		super(errorStr);
		this.errorStr = errorStr;
	}
}
