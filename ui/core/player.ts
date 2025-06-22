import { ItemSwapSettings } from './components/item_swap_picker';
import Toast from './components/toast';
import * as Mechanics from './constants/mechanics';
import { CURRENT_API_VERSION } from './constants/other';
import { SimSettingCategories } from './constants/sim_settings';
import { IndividualSimUIConfig } from './individual_sim_ui';
import { MAX_PARTY_SIZE, Party } from './party';
import { PlayerClass } from './player_class';
import { PlayerSpec } from './player_spec';
import { PlayerSpecs } from './player_specs';
import {
	AuraStats as AuraStatsProto,
	ErrorOutcomeType,
	Player as PlayerProto,
	PlayerStats,
	SpellStats as SpellStatsProto,
	StatWeightsResult,
	UnitMetadata as UnitMetadataProto,
} from './proto/api';
import { APLRotation, APLRotation_Type as APLRotationType, SimpleRotation } from './proto/apl';
import {
	Class,
	ConsumesSpec,
	Cooldowns,
	Faction,
	GemColor,
	Glyphs,
	HandType,
	HealingModel,
	IndividualBuffs,
	ItemLevelState,
	ItemRandomSuffix,
	ItemSlot,
	Profession,
	PseudoStat,
	Race,
	RangedWeaponType,
	ReforgeStat,
	Spec,
	Stat,
	UnitReference,
	UnitStats,
} from './proto/common';
import { SimDatabase } from './proto/db';
import {
	DungeonDifficulty,
	RaidFilterOption,
	SourceFilterOption,
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
	UIItem_FactionRestriction,
} from './proto/ui';
import { ActionId } from './proto_utils/action_id';
import { Database } from './proto_utils/database';
import { EquippedItem, ReforgeData } from './proto_utils/equipped_item';
import { Gear, ItemSwapGear } from './proto_utils/gear';
import { gemMatchesSocket, isUnrestrictedGem } from './proto_utils/gems';
import SecondaryResource from './proto_utils/secondary_resource';
import { StatCap, Stats } from './proto_utils/stats';
import {
	AL_CATEGORY_HARD_MODE,
	canEquipEnchant,
	canEquipItem,
	ClassOptions,
	ClassSpecs,
	emptyUnitReference,
	enchantAppliesToItem,
	getMetaGemEffectEP,
	getTalentTreePoints,
	isPVPItem,
	newUnitReference,
	raceToFaction,
	SpecClasses,
	SpecOptions,
	SpecRotation,
	SpecTalents,
	SpecTypeFunctions,
	specTypeFunctions,
	withSpec,
} from './proto_utils/utils';
import { Raid } from './raid';
import { Sim } from './sim';
import { playerTalentStringToProto } from './talents/factory';
import { EventID, TypedEvent } from './typed_event';
import { omitDeep, stringComparator, sum } from './utils';
import { WorkerProgressCallback } from './worker_pool';

export interface AuraStats {
	data: AuraStatsProto;
	id: ActionId;
}
export interface SpellStats {
	data: SpellStatsProto;
	id: ActionId;
}

export class UnitMetadata {
	private name: string;
	private auras: Array<AuraStats>;
	private spells: Array<SpellStats>;

	constructor() {
		this.name = '';
		this.auras = [];
		this.spells = [];
	}

	getName(): string {
		return this.name;
	}

	getAuras(): Array<AuraStats> {
		return this.auras.slice();
	}

	getSpells(): Array<SpellStats> {
		return this.spells.slice();
	}

	// Returns whether any updates were made.
	async update(metadata: UnitMetadataProto): Promise<boolean> {
		let newSpells = metadata!.spells.map(spell => {
			return {
				data: spell,
				id: ActionId.fromProto(spell.id!),
			};
		});
		let newAuras = metadata!.auras.map(aura => {
			return {
				data: aura,
				id: ActionId.fromProto(aura.id!),
			};
		});

		await Promise.all([...newSpells, ...newAuras].map(newSpell => newSpell.id.fill().then(newId => (newSpell.id = newId))));

		newSpells = newSpells.sort((a, b) => stringComparator(a.id.name, b.id.name));
		newAuras = newAuras.sort((a, b) => stringComparator(a.id.name, b.id.name));

		let anyUpdates = false;
		if (metadata.name != this.name) {
			this.name = metadata.name;
			anyUpdates = true;
		}
		if (newSpells.length != this.spells.length || newSpells.some((newSpell, i) => !newSpell.id.equals(this.spells[i].id))) {
			this.spells = newSpells;
			anyUpdates = true;
		}
		if (newAuras.length != this.auras.length || newAuras.some((newAura, i) => !newAura.id.equals(this.auras[i].id))) {
			this.auras = newAuras;
			anyUpdates = true;
		}

		return anyUpdates;
	}
}

export class UnitMetadataList {
	private metadatas: Array<UnitMetadata>;

	constructor() {
		this.metadatas = [];
	}

	async update(newMetadatas: Array<UnitMetadataProto>): Promise<boolean> {
		const oldLen = this.metadatas.length;

		if (newMetadatas.length > oldLen) {
			for (let i = oldLen; i < newMetadatas.length; i++) {
				this.metadatas.push(new UnitMetadata());
			}
		} else if (newMetadatas.length < oldLen) {
			this.metadatas = this.metadatas.slice(0, newMetadatas.length);
		}

		const anyUpdates = await Promise.all(newMetadatas.map((metadata, i) => this.metadatas[i].update(metadata)));

		return oldLen != this.metadatas.length || anyUpdates.some(v => v);
	}

	asList(): Array<UnitMetadata> {
		return this.metadatas.slice();
	}
}

export interface MeleeCritCapInfo {
	meleeCrit: number;
	meleeHit: number;
	expertise: number;
	suppression: number;
	glancing: number;
	hasOffhandWeapon: boolean;
	meleeHitCap: number;
	expertiseCap: number;
	remainingMeleeHitCap: number;
	remainingExpertiseCap: number;
	baseCritCap: number;
	specSpecificOffset: number;
	playerCritCapDelta: number;
}

export type AutoRotationGenerator<SpecType extends Spec> = (player: Player<SpecType>) => APLRotation;
export type SimpleRotationGenerator<SpecType extends Spec> = (
	player: Player<SpecType>,
	simpleRotation: SpecRotation<SpecType>,
	cooldowns: Cooldowns,
) => APLRotation;

export interface PlayerConfig<SpecType extends Spec> {
	autoRotation: AutoRotationGenerator<SpecType>;
	simpleRotation?: SimpleRotationGenerator<SpecType>;
	hiddenMCDs?: Array<number>; // spell IDs for any MCDs that should be omitted from the Simple Cooldowns UI
	secondaryResource?: SecondaryResource | null;
}

const SPEC_CONFIGS: Partial<Record<Spec, PlayerConfig<any>>> = {};

export function registerSpecConfig<SpecType extends Spec>(spec: SpecType, config: PlayerConfig<SpecType>) {
	SPEC_CONFIGS[spec] = config;
}

export function getSpecConfig<SpecType extends Spec>(spec: SpecType): PlayerConfig<SpecType> {
	const config = SPEC_CONFIGS[spec] as PlayerConfig<SpecType>;
	config.secondaryResource = SecondaryResource.create(spec);
	if (!config) {
		throw new Error('No config registered for Spec: ' + spec);
	}
	return config;
}

// Manages all the gear / consumes / other settings for a single Player.
export class Player<SpecType extends Spec> {
	readonly sim: Sim;
	private party: Party | null;
	private raid: Raid | null;

	readonly playerSpec: PlayerSpec<SpecType>;
	readonly playerClass: PlayerClass<SpecClasses<SpecType>>;
	readonly secondaryResource?: SecondaryResource | null;

	private name = '';
	private buffs: IndividualBuffs = IndividualBuffs.create();
	private consumables: ConsumesSpec = ConsumesSpec.create();
	private bonusStats: Stats = new Stats();
	private gear: Gear = new Gear({});
	//private bulkEquipmentSpec: BulkEquipmentSpec = BulkEquipmentSpec.create();
	itemSwapSettings: ItemSwapSettings;
	private race: Race;
	private profession1: Profession = 0;
	private profession2: Profession = 0;
	aplRotation: APLRotation = APLRotation.create();
	private talentsString = '';
	private glyphs: Glyphs = Glyphs.create();
	private specOptions: SpecOptions<SpecType>;
	private reactionTime = 0;
	private channelClipDelay = 0;
	private inFrontOfTarget = false;
	private distanceFromTarget = 0;
	private healingModel: HealingModel = HealingModel.create();
	private healingEnabled = false;
	private challengeModeEnabled = false;

	private readonly autoRotationGenerator: AutoRotationGenerator<SpecType> | null = null;
	private readonly simpleRotationGenerator: SimpleRotationGenerator<SpecType> | null = null;
	readonly hiddenMCDs: Array<number>;

	private itemEPCache = new Array<Map<string, number>>();
	private gemEPCache = new Map<number, number>();
	private randomSuffixEPCache = new Map<number, number>();
	private enchantEPCache = new Map<number, number>();
	private upgradeEPCache = new Map<string, number>();
	private talents: SpecTalents<SpecType> | null = null;
	private specConfig: IndividualSimUIConfig<SpecType>;

	readonly specTypeFunctions: SpecTypeFunctions<SpecType>;

	private static readonly numEpRatios = 6;
	private epRatios: Array<number> = new Array<number>(Player.numEpRatios).fill(0);
	private epWeights: Stats = new Stats();
	private statCaps: Stats = new Stats();
	private softCapBreakpoints: StatCap[] = [];
	private breakpointLimits: Stats = new Stats();
	private currentStats: PlayerStats = PlayerStats.create();
	private metadata: UnitMetadata = new UnitMetadata();
	private petMetadatas: UnitMetadataList = new UnitMetadataList();

	readonly nameChangeEmitter = new TypedEvent<void>('PlayerName');
	readonly buffsChangeEmitter = new TypedEvent<void>('PlayerBuffs');
	readonly consumesChangeEmitter = new TypedEvent<void>('PlayerConsumes');
	readonly bonusStatsChangeEmitter = new TypedEvent<void>('PlayerBonusStats');
	readonly gearChangeEmitter = new TypedEvent<void>('PlayerGear');
	readonly professionChangeEmitter = new TypedEvent<void>('PlayerProfession');
	readonly raceChangeEmitter = new TypedEvent<void>('PlayerRace');
	readonly rotationChangeEmitter = new TypedEvent<void>('PlayerRotation');
	readonly talentsChangeEmitter = new TypedEvent<void>('PlayerTalents');
	readonly glyphsChangeEmitter = new TypedEvent<void>('PlayerGlyphs');
	readonly specOptionsChangeEmitter = new TypedEvent<void>('PlayerSpecOptions');
	readonly inFrontOfTargetChangeEmitter = new TypedEvent<void>('PlayerInFrontOfTarget');
	readonly distanceFromTargetChangeEmitter = new TypedEvent<void>('PlayerDistanceFromTarget');
	readonly healingModelChangeEmitter = new TypedEvent<void>('PlayerHealingModel');
	readonly epWeightsChangeEmitter = new TypedEvent<void>('PlayerEpWeights');
	readonly statCapsChangeEmitter = new TypedEvent<void>('StatCaps');
	readonly softCapBreakpointsChangeEmitter = new TypedEvent<void>('SoftCapBreakpoints');
	readonly breakpointLimitsChangeEmitter = new TypedEvent<void>('BreakpointLimits');
	readonly miscOptionsChangeEmitter = new TypedEvent<void>('PlayerMiscOptions');
	readonly challengeModeChangeEmitter = new TypedEvent<void>('ChallengeMode');

	readonly currentStatsEmitter = new TypedEvent<void>('PlayerCurrentStats');
	readonly epRatiosChangeEmitter = new TypedEvent<void>('PlayerEpRatios');
	readonly epRefStatChangeEmitter = new TypedEvent<void>('PlayerEpRefStat');

	// Emits when any of the above emitters emit.
	readonly changeEmitter: TypedEvent<void>;

	constructor(spec: PlayerSpec<SpecType>, sim: Sim) {
		this.sim = sim;
		this.party = null;
		this.raid = null;

		this.playerSpec = spec;
		this.playerClass = PlayerSpecs.getPlayerClass(spec);

		this.race = this.playerClass.races[0];
		this.specTypeFunctions = specTypeFunctions[this.getSpec()] as SpecTypeFunctions<SpecType>;
		this.specOptions = this.specTypeFunctions.optionsCreate();

		this.specConfig = getSpecConfig<SpecType>(this.getSpec()) as IndividualSimUIConfig<SpecType>;
		this.secondaryResource = this.specConfig.secondaryResource;

		this.autoRotationGenerator = this.specConfig.autoRotation;
		if (this.specConfig.simpleRotation) {
			this.simpleRotationGenerator = this.specConfig.simpleRotation;
		} else {
			this.simpleRotationGenerator = null;
		}
		this.hiddenMCDs = this.specConfig.hiddenMCDs || new Array<number>();

		for (let i = 0; i < ItemSlot.ItemSlotOffHand + 1; ++i) {
			this.itemEPCache[i] = new Map();
		}

		this.itemSwapSettings = new ItemSwapSettings(this);

		this.bindChallengeModeChange();

		this.changeEmitter = TypedEvent.onAny(
			[
				this.nameChangeEmitter,
				this.buffsChangeEmitter,
				this.consumesChangeEmitter,
				this.bonusStatsChangeEmitter,
				this.gearChangeEmitter,
				this.professionChangeEmitter,
				this.raceChangeEmitter,
				this.rotationChangeEmitter,
				this.talentsChangeEmitter,
				this.glyphsChangeEmitter,
				this.specOptionsChangeEmitter,
				this.miscOptionsChangeEmitter,
				this.inFrontOfTargetChangeEmitter,
				this.distanceFromTargetChangeEmitter,
				this.healingModelChangeEmitter,
				this.epWeightsChangeEmitter,
				this.epRatiosChangeEmitter,
				this.epRefStatChangeEmitter,
				this.statCapsChangeEmitter,
				this.breakpointLimitsChangeEmitter,
				this.challengeModeChangeEmitter,
			],
			'PlayerChange',
		);
	}

	bindChallengeModeChange() {
		this.challengeModeChangeEmitter.on(() => {
			this.setGear(TypedEvent.nextEventID(), this.gear, true);
		});
	}

	getSpecIcon(): string {
		return this.playerSpec.getIcon('medium');
	}

	getPlayerSpec(): PlayerSpec<SpecType> {
		return this.playerSpec;
	}

	getSpec(): SpecType {
		return this.getPlayerSpec().specID;
	}

	getPlayerClass(): PlayerClass<SpecClasses<SpecType>> {
		return this.playerClass;
	}

	getClass(): SpecClasses<SpecType> {
		return this.playerSpec.classID;
	}

	getClassColor(): string {
		return this.playerClass.hexColor;
	}

	canEnableTargetDummies(): boolean {
		const healingSpellClasses: Class[] = [Class.ClassDruid, Class.ClassPaladin, Class.ClassPriest, Class.ClassShaman, Class.ClassMonk];
		return healingSpellClasses.includes(this.getClass());
	}

	shouldEnableTargetDummies(): boolean {
		if (this.getPlayerSpec().isHealingSpec || this.getPlayerSpec().isTankSpec) {
			return true;
		}

		if (!this.itemSwapSettings.getEnableItemSwap()) {
			return false;
		}

		// Not comprehensive, add other relevant IDs here as needed.
		const healingProcTrinkets: number[] = [72898, 77969, 77204, 77989];

		return this.itemSwapSettings.getGear().hasTrinketFromOptions(healingProcTrinkets);
	}

	// TODO: Cata - Check this
	isSpec<T extends Spec>(specId: T): this is Player<T> {
		return (this.getSpec() as unknown) == specId;
	}

	isClass<T extends Class>(classId: T): this is Player<ClassSpecs<T>> {
		return (this.getClass() as unknown) == classId;
	}

	getParty(): Party | null {
		return this.party;
	}

	getRaid(): Raid | null {
		return this.raid;
	}

	// Returns this player's index within its party [0-4].
	getPartyIndex(): number {
		if (this.party == null) {
			throw new Error("Can't get party index for player without a party!");
		}

		return this.party.getPlayers().indexOf(this);
	}

	// Returns this player's index within its raid [0-24].
	getRaidIndex(): number {
		if (this.party == null) {
			throw new Error("Can't get raid index for player without a party!");
		}

		return this.party.getIndex() * MAX_PARTY_SIZE + this.getPartyIndex();
	}

	// This should only ever be called from party.
	setParty(newParty: Party | null) {
		if (newParty == null) {
			this.party = null;
			this.raid = null;
		} else {
			this.party = newParty;
			this.raid = newParty.raid;
		}
	}

	getOtherPartyMembers(): Array<Player<any>> {
		if (this.party == null) {
			return [];
		}

		return this.party.getPlayers().filter(player => player != null && player != this) as Array<Player<any>>;
	}

	// Returns all items that this player can wear in the given slot.
	getItems(slot: ItemSlot): Array<Item> {
		return this.sim.db.getItems(slot).filter(item => canEquipItem(item, this.playerSpec, slot));
	}

	// Returns all random suffixes that this player would be interested in for the given base item.
	getRandomSuffixes(item: Item): Array<ItemRandomSuffix> {
		return item.randomSuffixOptions
			.map(id => this.sim.db.getRandomSuffixById(id))
			.filter((suffix): suffix is ItemRandomSuffix => !!suffix && this.computeRandomSuffixEP(suffix) > 0);
	}

	// Returns all reforgings that are valid with a given item
	getAvailableReforgings(equippedItem: EquippedItem): Array<ReforgeData> {
		return this.sim.db.getAvailableReforges(equippedItem.item).map(reforge => equippedItem.getReforgeData(reforge)!);
	}

	// Returns reforge given an id
	getReforge(id: number): ReforgeStat | undefined {
		return this.sim.db.getReforgeById(id);
	}

	// Returns all enchants that this player can wear in the given slot.
	getEnchants(slot: ItemSlot): Array<Enchant> {
		return this.sim.db.getEnchants(slot).filter(enchant => canEquipEnchant(enchant, this.playerSpec));
	}

	// Returns all tinkers that this player can wear in the given slot.
	// For the purpose of this function, they are all enchants still, however we split them since you can have both on the same item.
	getTinkers(slot: ItemSlot): Array<Enchant> {
		return this.sim.db.getEnchants(slot).filter(enchant => enchant.requiredProfession == Profession.Engineering);
	}

	// Returns all gems that this player can wear of the given color.
	getGems(socketColor?: GemColor): Array<Gem> {
		return this.sim.db.getGems(socketColor);
	}

	getEpWeights(): Stats {
		return this.epWeights;
	}

	setEpWeights(eventID: EventID, newEpWeights: Stats) {
		this.epWeights = newEpWeights;
		this.epWeightsChangeEmitter.emit(eventID);

		this.gemEPCache = new Map();
		this.enchantEPCache = new Map();
		this.randomSuffixEPCache = new Map();
		this.upgradeEPCache = new Map();
		for (let i = 0; i < ItemSlot.ItemSlotOffHand + 1; ++i) {
			this.itemEPCache[i] = new Map();
		}
	}

	getStatCaps(): Stats {
		return this.statCaps;
	}

	setStatCaps(eventID: EventID, newStatCaps: Stats) {
		this.statCaps = newStatCaps;
		this.statCapsChangeEmitter.emit(eventID);
	}

	getSoftCapBreakpoints(): StatCap[] {
		return this.softCapBreakpoints;
	}

	setSoftCapBreakpoints(eventID: EventID, newSoftCapBreakpoints: StatCap[]) {
		this.softCapBreakpoints = newSoftCapBreakpoints;
		this.softCapBreakpointsChangeEmitter.emit(eventID);
	}
	getBreakpointLimits(): Stats {
		return this.breakpointLimits;
	}

	setBreakpointLimits(eventID: EventID, newLimits: Stats) {
		this.breakpointLimits = newLimits;
		this.breakpointLimitsChangeEmitter.emit(eventID);
	}

	getDefaultEpRatios(isTankSpec: boolean, isHealingSpec: boolean): Array<number> {
		const defaultRatios = new Array(Player.numEpRatios).fill(0);
		if (isHealingSpec) {
			// By default only value HPS EP for healing spec
			defaultRatios[1] = 1;
		} else if (isTankSpec) {
			// By default value TPS and DTPS EP equally for tanking spec
			defaultRatios[2] = 1;
			defaultRatios[3] = 1;
			if (this.getSpec() == Spec.SpecBloodDeathKnight) {
				// Add healing EPs for BDKs
				defaultRatios[1] = 1;
			}
		} else {
			// By default only value DPS EP
			defaultRatios[0] = 1;
		}

		return defaultRatios;
	}

	getEpRatios() {
		return this.epRatios.slice();
	}

	setEpRatios(eventID: EventID, newRatios: Array<number>) {
		this.epRatios = newRatios;
		this.epRatiosChangeEmitter.emit(eventID);
	}

	async computeStatWeights(
		_eventID: EventID,
		epStats: Array<Stat>,
		epPseudoStats: Array<PseudoStat>,
		epReferenceStat: Stat,
		onProgress: WorkerProgressCallback,
	): Promise<StatWeightsResult | null> {
		try {
			const result = await this.sim.statWeights(this, epStats, epPseudoStats, epReferenceStat, onProgress);
			if (result.error) {
				if (result.error.type == ErrorOutcomeType.ErrorOutcomeAborted) {
					new Toast({
						variant: 'info',
						body: 'Statweight sim cancelled.',
					});
				}
				return null;
			}
			return result;
		} catch (error: any) {
			// TODO: Show crash report like for raid sim?
			console.error(error);
			new Toast({
				variant: 'error',
				body: error?.message || 'Something went wrong calculating your stat weights. Reload the page and try again.',
			});
			return null;
		}
	}

	getCurrentStats(): PlayerStats {
		return PlayerStats.clone(this.currentStats);
	}

	setCurrentStats(eventID: EventID, newStats: PlayerStats) {
		this.currentStats = newStats;
		this.currentStatsEmitter.emit(eventID);
	}

	getMetadata(): UnitMetadata {
		return this.metadata;
	}

	getPetMetadatas(): UnitMetadataList {
		return this.petMetadatas;
	}

	async updateMetadata(): Promise<boolean> {
		const playerPromise = this.metadata.update(this.currentStats.metadata!);
		const petsPromise = this.petMetadatas.update(this.currentStats.pets.map(p => p.metadata!));
		const playerUpdated = await playerPromise;
		const petsUpdated = await petsPromise;
		return playerUpdated || petsUpdated;
	}

	getName(): string {
		return this.name;
	}
	setName(eventID: EventID, newName: string) {
		if (newName != this.name) {
			this.name = newName;
			this.nameChangeEmitter.emit(eventID);
		}
	}

	getLabel(): string {
		if (this.party) {
			return `${this.name} (#${this.getRaidIndex() + 1})`;
		} else {
			return this.name;
		}
	}

	getRace(): Race {
		return this.race;
	}
	setRace(eventID: EventID, newRace: Race) {
		if (newRace != this.race) {
			this.race = newRace;
			this.raceChangeEmitter.emit(eventID);
		}
	}

	getProfession1(): Profession {
		return this.profession1;
	}
	setProfession1(eventID: EventID, newProfession: Profession) {
		if (newProfession != this.profession1) {
			this.profession1 = newProfession;
			this.professionChangeEmitter.emit(eventID);
		}
	}
	getProfession2(): Profession {
		return this.profession2;
	}
	setProfession2(eventID: EventID, newProfession: Profession) {
		if (newProfession != this.profession2) {
			this.profession2 = newProfession;
			this.professionChangeEmitter.emit(eventID);
		}
	}
	getProfessions(): Array<Profession> {
		return [this.profession1, this.profession2].filter(p => p != Profession.ProfessionUnknown);
	}
	setProfessions(eventID: EventID, newProfessions: Array<Profession>) {
		TypedEvent.freezeAllAndDo(() => {
			this.setProfession1(eventID, newProfessions[0] || Profession.ProfessionUnknown);
			this.setProfession2(eventID, newProfessions[1] || Profession.ProfessionUnknown);
		});
	}
	hasProfession(prof: Profession): boolean {
		return this.getProfessions().includes(prof);
	}
	isBlacksmithing(): boolean {
		return this.hasProfession(Profession.Blacksmithing);
	}

	getFaction(): Faction {
		return raceToFaction[this.getRace()];
	}

	getBuffs(): IndividualBuffs {
		// Make a defensive copy
		return IndividualBuffs.clone(this.buffs);
	}

	setBuffs(eventID: EventID, newBuffs: IndividualBuffs) {
		if (IndividualBuffs.equals(this.buffs, newBuffs)) return;

		// Make a defensive copy
		this.buffs = IndividualBuffs.clone(newBuffs);
		this.buffsChangeEmitter.emit(eventID);
	}

	getConsumes(): ConsumesSpec {
		// Make a defensive copy
		return ConsumesSpec.clone(this.consumables);
	}

	setConsumes(eventID: EventID, newConsumes: ConsumesSpec) {
		if (ConsumesSpec.equals(this.consumables, newConsumes)) return;

		// Make a defensive copy
		this.consumables = ConsumesSpec.clone(newConsumes);
		this.consumesChangeEmitter.emit(eventID);
	}

	canDualWield2H(): boolean {
		return this.getSpec() == Spec.SpecFuryWarrior;
	}

	equipItem(eventID: EventID, slot: ItemSlot, newItem: EquippedItem | null) {
		this.setGear(eventID, this.gear.withEquippedItem(slot, newItem, this.canDualWield2H()));
	}

	getEquippedItem(slot: ItemSlot): EquippedItem | null {
		return this.gear.getEquippedItem(slot);
	}

	getEquippedItems(): Array<EquippedItem | null> {
		return this.gear.getEquippedItems();
	}

	getGear(): Gear {
		return this.gear;
	}

	setGear(eventID: EventID, newGear: Gear, forceUpdate?: boolean) {
		if (newGear.equals(this.gear) && !forceUpdate) return;
		this.gear = newGear.withChallengeMode(this.challengeModeEnabled);
		this.gearChangeEmitter.emit(eventID);
	}

	/*
	setBulkEquipmentSpec(eventID: EventID, newBulkEquipmentSpec: BulkEquipmentSpec) {
		if (BulkEquipmentSpec.equals(this.bulkEquipmentSpec, newBulkEquipmentSpec))
			return;

		TypedEvent.freezeAllAndDo(() => {
			this.bulkEquipmentSpec = newBulkEquipmentSpec;
			this.bulkGearChangeEmitter.emit(eventID);
		});
	}

	getBulkEquipmentSpec(): BulkEquipmentSpec {
		return BulkEquipmentSpec.clone(this.bulkEquipmentSpec);
	}
	*/

	getBonusStats(): Stats {
		return this.bonusStats;
	}

	setBonusStats(eventID: EventID, newBonusStats: Stats) {
		if (newBonusStats.equals(this.bonusStats)) return;

		this.bonusStats = newBonusStats;
		this.bonusStatsChangeEmitter.emit(eventID);
	}

	getMeleeCritCapInfo(): MeleeCritCapInfo {
		const meleeCrit = this.currentStats.finalStats?.pseudoStats[PseudoStat.PseudoStatPhysicalCritPercent] || 0.0;
		const meleeHit = this.currentStats.finalStats?.pseudoStats[PseudoStat.PseudoStatPhysicalHitPercent] || 0.0;
		const expertise = (this.currentStats.finalStats?.stats[Stat.StatExpertiseRating] || 0.0) / Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION / 4;
		//const agility = (this.currentStats.finalStats?.stats[Stat.StatAgility] || 0.0) / this.getClass();
		const suppression = 4.8;
		const glancing = 24.0;

		const hasOffhandWeapon = this.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand)?.item.weaponSpeed !== undefined;
		// Due to warrior HS bug, hit cap for crit cap calculation should be 8% instead of 27%
		const meleeHitCap = hasOffhandWeapon && this.getClass() != Class.ClassWarrior ? 27.0 : 8.0;
		const dodgeCap = 6.5;
		const parryCap = this.getInFrontOfTarget() ? 14.0 : 0;
		const expertiseCap = dodgeCap + parryCap;

		const remainingMeleeHitCap = Math.max(meleeHitCap - meleeHit, 0.0);
		const remainingDodgeCap = Math.max(dodgeCap - expertise, 0.0);
		const remainingParryCap = Math.max(parryCap - expertise, 0.0);
		const remainingExpertiseCap = remainingDodgeCap + remainingParryCap;

		const specSpecificOffset = 0.0;

		// if (this.getSpec() === Spec.SpecEnhancementShaman) {
		// 	// Elemental Devastation uptime is near 100%
		// 	// TODO: Cata - Check this
		// 	const ranks = (this as unknown as Player<Spec.SpecEnhancementShaman>).getTalents().elementalDevastation;
		// 	specSpecificOffset = 3.0 * ranks;
		// }

		const baseCritCap = 100.0 - glancing + suppression - remainingMeleeHitCap - remainingExpertiseCap - specSpecificOffset;
		const playerCritCapDelta = meleeCrit - baseCritCap;

		return {
			meleeCrit,
			meleeHit,
			expertise,
			suppression,
			glancing,
			hasOffhandWeapon,
			meleeHitCap,
			expertiseCap,
			remainingMeleeHitCap,
			remainingExpertiseCap,
			baseCritCap,
			specSpecificOffset,
			playerCritCapDelta,
		};
	}

	getMeleeCritCap() {
		return this.getMeleeCritCapInfo().playerCritCapDelta;
	}

	setAplRotation(eventID: EventID, newRotation: APLRotation) {
		if (APLRotation.equals(newRotation, this.aplRotation)) return;

		this.aplRotation = APLRotation.clone(newRotation);
		this.rotationChangeEmitter.emit(eventID);
	}

	getSimpleRotation(): SpecRotation<SpecType> {
		const jsonStr = this.aplRotation.simple?.specRotationJson || '';
		if (!jsonStr) {
			return this.specTypeFunctions.rotationCreate();
		}

		try {
			const json = JSON.parse(jsonStr);
			return this.specTypeFunctions.rotationFromJson(json);
		} catch (e) {
			console.warn(`Error parsing rotation spec options: ${e}\n\nSpec options: '${jsonStr}'`);
			return this.specTypeFunctions.rotationCreate();
		}
	}

	setSimpleRotation(eventID: EventID, newRotation: SpecRotation<SpecType>) {
		if (this.specTypeFunctions.rotationEquals(newRotation, this.getSimpleRotation())) return;

		if (!this.aplRotation.simple) {
			this.aplRotation.simple = SimpleRotation.create();
		}
		this.aplRotation.simple.specRotationJson = JSON.stringify(this.specTypeFunctions.rotationToJson(newRotation));

		this.rotationChangeEmitter.emit(eventID);
	}

	getSimpleCooldowns(): Cooldowns {
		// Make a defensive copy
		return Cooldowns.clone(this.aplRotation.simple?.cooldowns || Cooldowns.create());
	}

	setSimpleCooldowns(eventID: EventID, newCooldowns: Cooldowns) {
		if (Cooldowns.equals(this.getSimpleCooldowns(), newCooldowns)) return;

		if (!this.aplRotation.simple) {
			this.aplRotation.simple = SimpleRotation.create();
		}
		this.aplRotation.simple.cooldowns = newCooldowns;
		this.rotationChangeEmitter.emit(eventID);
	}

	getRotationType(): APLRotationType {
		if (this.aplRotation.type == APLRotationType.TypeUnknown) {
			return APLRotationType.TypeAPL;
		} else {
			return this.aplRotation.type;
		}
	}

	hasSimpleRotationGenerator(): boolean {
		return this.simpleRotationGenerator != null;
	}

	getResolvedAplRotation(): APLRotation {
		const type = this.getRotationType();
		if (type == APLRotationType.TypeAuto && this.autoRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			const rot = APLRotation.clone(this.autoRotationGenerator(this));
			rot.type = APLRotationType.TypeAuto;
			return rot;
		} else if (type == APLRotationType.TypeSimple && this.simpleRotationGenerator) {
			// Clone to avoid modifying preset rotations, which are often returned directly.
			const simpleRot = this.getSimpleRotation();
			const rot = APLRotation.clone(this.simpleRotationGenerator(this, simpleRot, this.getSimpleCooldowns()));
			rot.simple = this.aplRotation.simple;
			rot.type = APLRotationType.TypeSimple;
			return rot;
		} else {
			return this.aplRotation;
		}
	}

	getTalents(): SpecTalents<SpecType> {
		if (this.talents == null) {
			this.talents = playerTalentStringToProto(this.playerSpec, this.talentsString) as SpecTalents<SpecType>;
		}
		return this.talents!;
	}

	getTalentsString(): string {
		return this.talentsString;
	}

	setTalentsString(eventID: EventID, newTalentsString: string) {
		if (newTalentsString == this.talentsString) return;

		this.talentsString = newTalentsString;
		this.talents = null;
		this.talentsChangeEmitter.emit(eventID);
	}

	getTalentTreePoints(): Array<number> {
		return getTalentTreePoints(this.getTalentsString());
	}

	getTalentTreeIcon(): string {
		return this.playerSpec.getIcon('medium');
	}

	getGlyphs(): Glyphs {
		// Make a defensive copy
		return Glyphs.clone(this.glyphs);
	}

	setGlyphs(eventID: EventID, newGlyphs: Glyphs) {
		if (Glyphs.equals(this.glyphs, newGlyphs)) return;

		// Make a defensive copy
		this.glyphs = Glyphs.clone(newGlyphs);
		this.glyphsChangeEmitter.emit(eventID);
	}

	getMajorGlyphs(): Array<number> {
		return [this.glyphs.major1, this.glyphs.major2, this.glyphs.major3].filter(glyph => glyph != 0);
	}

	getMinorGlyphs(): Array<number> {
		return [this.glyphs.minor1, this.glyphs.minor2, this.glyphs.minor3].filter(glyph => glyph != 0);
	}

	getAllGlyphs(): Array<number> {
		return this.getMajorGlyphs().concat(this.getMinorGlyphs());
	}

	getClassOptions(): ClassOptions<SpecType> {
		return this.getSpecOptions().classOptions as ClassOptions<SpecType>;
	}

	setClassOptions(eventID: EventID, newClassOptions: ClassOptions<SpecType>) {
		const newSpecOptions = this.getSpecOptions();
		newSpecOptions.classOptions = newClassOptions;
		if (this.specTypeFunctions.optionsEquals(newSpecOptions, this.specOptions)) return;

		this.specOptions = this.specTypeFunctions.optionsCopy(newSpecOptions);
		this.specOptionsChangeEmitter.emit(eventID);
	}

	getSpecOptions(): SpecOptions<SpecType> {
		return this.specTypeFunctions.optionsCopy(this.specOptions);
	}

	setSpecOptions(eventID: EventID, newSpecOptions: SpecOptions<SpecType>) {
		if (this.specTypeFunctions.optionsEquals(newSpecOptions, this.specOptions)) return;

		this.specOptions = this.specTypeFunctions.optionsCopy(newSpecOptions);
		this.specOptionsChangeEmitter.emit(eventID);
	}

	getReactionTime(): number {
		return this.reactionTime;
	}

	setReactionTime(eventID: EventID, newReactionTime: number) {
		if (newReactionTime == this.reactionTime) return;

		this.reactionTime = newReactionTime;
		this.miscOptionsChangeEmitter.emit(eventID);
	}

	getChannelClipDelay(): number {
		return this.channelClipDelay;
	}

	setChannelClipDelay(eventID: EventID, newChannelClipDelay: number) {
		if (newChannelClipDelay == this.channelClipDelay) return;

		this.channelClipDelay = newChannelClipDelay;
		this.miscOptionsChangeEmitter.emit(eventID);
	}

	getChallengeModeEnabled(): boolean {
		return this.challengeModeEnabled;
	}

	setChallengeModeEnabled(eventID: EventID, value: boolean) {
		if (value === this.challengeModeEnabled) return;

		this.challengeModeEnabled = value;
		this.challengeModeChangeEmitter.emit(eventID);
	}

	getInFrontOfTarget(): boolean {
		return this.inFrontOfTarget;
	}

	setInFrontOfTarget(eventID: EventID, newInFrontOfTarget: boolean) {
		if (newInFrontOfTarget == this.inFrontOfTarget) return;

		this.inFrontOfTarget = newInFrontOfTarget;
		this.inFrontOfTargetChangeEmitter.emit(eventID);
	}

	getDistanceFromTarget(): number {
		return this.distanceFromTarget;
	}

	setDistanceFromTarget(eventID: EventID, newDistanceFromTarget: number) {
		if (newDistanceFromTarget == this.distanceFromTarget) return;

		this.distanceFromTarget = newDistanceFromTarget;
		this.distanceFromTargetChangeEmitter.emit(eventID);
	}

	setDefaultHealingParams(hm: HealingModel) {
		const boss = this.sim.encounter.primaryTarget;
		const dualWield = boss.dualWield;
		if (hm.cadenceSeconds == 0) {
			let maxCadence = 1.5 * boss.swingSpeed;
			if (dualWield) {
				maxCadence /= 2;
			}
			hm.cadenceSeconds = 0.4;
			hm.cadenceVariation = maxCadence - hm.cadenceSeconds;
		}
		if (hm.hps == 0) {
			hm.hps = (0.25 * boss.minBaseDamage) / boss.swingSpeed;
			if (dualWield) {
				hm.hps *= 1.5;
			}
		}
	}

	enableHealing() {
		this.healingEnabled = true;
		const hm = this.getHealingModel();
		if (hm.cadenceSeconds == 0 || hm.hps == 0) {
			this.setDefaultHealingParams(hm);
			this.setHealingModel(0, hm);
		}
	}

	getHealingModel(): HealingModel {
		// Make a defensive copy
		return HealingModel.clone(this.healingModel);
	}

	setHealingModel(eventID: EventID, newHealingModel: HealingModel) {
		if (HealingModel.equals(this.healingModel, newHealingModel)) return;

		// Make a defensive copy
		this.healingModel = HealingModel.clone(newHealingModel);
		// If we have enabled healing model and try to set 0s cadence or 0 incoming HPS, then set intelligent defaults instead based on boss parameters.
		if (this.healingEnabled) {
			this.setDefaultHealingParams(this.healingModel);
		}
		this.healingModelChangeEmitter.emit(eventID);
	}

	computeStatsEP(stats?: Stats): number {
		if (stats == undefined) {
			return 0;
		}
		return stats.computeEP(this.epWeights);
	}

	computeGemEP(gem: Gem): number {
		if (this.gemEPCache.has(gem.id)) {
			return this.gemEPCache.get(gem.id)!;
		}

		const epFromStats = this.computeStatsEP(new Stats(gem.stats));
		const epFromEffect = getMetaGemEffectEP(this.playerSpec, gem, Stats.fromProto(this.currentStats.finalStats));
		let bonusEP = 0;
		// unique items are slightly worse than non-unique because you can have only one.
		if (gem.unique) {
			bonusEP -= 0.01;
		}

		const ep = epFromStats + epFromEffect + bonusEP;
		this.gemEPCache.set(gem.id, ep);
		return ep;
	}

	computeEnchantEP(enchant: Enchant): number {
		if (this.enchantEPCache.has(enchant.effectId)) {
			return this.enchantEPCache.get(enchant.effectId)!;
		}

		const ep = this.computeStatsEP(new Stats(enchant.stats));
		this.enchantEPCache.set(enchant.effectId, ep);
		return ep;
	}

	computeRandomSuffixEP(randomSuffix: ItemRandomSuffix): number {
		if (this.randomSuffixEPCache.has(randomSuffix.id)) {
			return this.randomSuffixEPCache.get(randomSuffix.id)!;
		}

		const ep = this.computeStatsEP(new Stats(randomSuffix.stats));
		this.randomSuffixEPCache.set(randomSuffix.id, ep);
		return ep;
	}

	computeReforgingEP(reforging: ReforgeData): number {
		let stats = new Stats([]);
		stats = stats.addStat(reforging.fromStat, reforging.fromAmount);
		stats = stats.addStat(reforging.toStat, reforging.toAmount);

		return this.computeStatsEP(stats);
	}

	computeUpgradeEP(equippedItem: EquippedItem, upgradeLevel: ItemLevelState, slot: ItemSlot): number {
		const cacheKey = `${equippedItem.id}-${slot}-${equippedItem.randomSuffix?.id}-${upgradeLevel}`;
		if (this.upgradeEPCache.has(cacheKey)) {
			return this.upgradeEPCache.get(cacheKey)!;
		}

		const stats = equippedItem.withUpgrade(upgradeLevel).calcStats(slot);
		const ep = this.computeStatsEP(stats);
		this.upgradeEPCache.set(cacheKey, ep);

		return ep;
	}

	computeItemEP(item: Item, slot: ItemSlot): number {
		if (item == null) return 0;
		const cacheKey = `${item.id}-${this.challengeModeEnabled}`;

		const cached = this.itemEPCache[slot].get(cacheKey);
		if (cached !== undefined) return cached;

		const equippedItem = new EquippedItem({
			item,
			challengeMode: this.challengeModeEnabled,
		}).withDynamicStats();
		const itemStats = equippedItem.calcStats(slot);

		// For random suffix items, use the suffix option with the highest EP for the purposes of ranking items in the picker.
		let maxSuffixEP = 0;
		if (item.randomSuffixOptions.length > 0) {
			const suffixEPs = equippedItem.item.randomSuffixOptions.map(id => this.computeRandomSuffixEP(this.sim.db.getRandomSuffixById(id)! || 0));
			maxSuffixEP = (Math.max(...suffixEPs) * equippedItem.item.randPropPoints) / 10000;
		}

		let ep = itemStats.computeEP(this.epWeights) + maxSuffixEP;

		// unique items are slightly worse than non-unique because you can have only one.
		if (item.unique) {
			ep -= 0.01;
		}

		// Compare whether its better to match sockets + get socket bonus, or just use best gems.
		const bestGemEPNotMatchingSockets = sum(
			item.gemSockets.map(socketColor => {
				const gems = this.sim.db.getGems(socketColor).filter(gem => isUnrestrictedGem(gem, this.sim.getPhase()));
				if (gems.length > 0) {
					return Math.max(...gems.map(gem => this.computeGemEP(gem)));
				} else {
					return 0;
				}
			}),
		);

		const bestGemEPMatchingSockets =
			sum(
				item.gemSockets.map(socketColor => {
					const gems = this.sim.db
						.getGems(socketColor)
						.filter(gem => isUnrestrictedGem(gem, this.sim.getPhase()) && gemMatchesSocket(gem, socketColor));
					if (gems.length > 0) {
						return Math.max(...gems.map(gem => this.computeGemEP(gem)));
					} else {
						return 0;
					}
				}),
			) + this.computeStatsEP(new Stats(item.socketBonus));

		ep += Math.max(bestGemEPMatchingSockets, bestGemEPNotMatchingSockets);

		this.itemEPCache[slot].set(cacheKey, ep);
		return ep;
	}

	async setWowheadData(equippedItem: EquippedItem, elem: HTMLElement) {
		const isBlacksmithing = this.hasProfession(Profession.Blacksmithing);
		const gemIds = equippedItem.gems.length ? equippedItem.curGems(isBlacksmithing).map(gem => (gem ? gem.id : 0)) : [];
		const enchantIds = [equippedItem.enchant?.effectId, equippedItem.tinker?.effectId].filter((id): id is number => id !== undefined);
		equippedItem.asActionId().setWowheadDataset(elem, {
			gemIds,
			itemLevel: Number(equippedItem.ilvl),
			enchantIds: enchantIds,
			reforgeId: equippedItem.reforge?.id,
			randomEnchantmentId: equippedItem.randomSuffix?.id,
			setPieceIds: this.gear
				.asArray()
				.filter(ei => ei != null)
				.map(ei => ei!.item.id),
			hasExtraSocket: equippedItem.hasExtraSocket(isBlacksmithing),
			upgradeStep: equippedItem.upgrade,
		});

		elem.dataset.whtticon = 'false';
	}

	static ARMOR_SLOTS: Array<ItemSlot> = [
		ItemSlot.ItemSlotHead,
		ItemSlot.ItemSlotShoulder,
		ItemSlot.ItemSlotChest,
		ItemSlot.ItemSlotWrist,
		ItemSlot.ItemSlotHands,
		ItemSlot.ItemSlotLegs,
		ItemSlot.ItemSlotWaist,
		ItemSlot.ItemSlotFeet,
	];

	static WEAPON_SLOTS: Array<ItemSlot> = [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand];

	static readonly DIFFICULTY_SRCS: Partial<Record<SourceFilterOption, DungeonDifficulty>> = {
		[SourceFilterOption.SourceDungeon]: DungeonDifficulty.DifficultyNormal,
		[SourceFilterOption.SourceDungeonH]: DungeonDifficulty.DifficultyHeroic,
		[SourceFilterOption.SourceRaidRF]: DungeonDifficulty.DifficultyRaid25RF,
		[SourceFilterOption.SourceRaid]: DungeonDifficulty.DifficultyRaid25,
		[SourceFilterOption.SourceRaidH]: DungeonDifficulty.DifficultyRaid25H,
	};

	static readonly HEROIC_TO_NORMAL: Partial<Record<DungeonDifficulty, DungeonDifficulty>> = {
		[DungeonDifficulty.DifficultyHeroic]: DungeonDifficulty.DifficultyNormal,
		[DungeonDifficulty.DifficultyRaid10H]: DungeonDifficulty.DifficultyRaid10,
		[DungeonDifficulty.DifficultyRaid25H]: DungeonDifficulty.DifficultyRaid25,
	};

	static readonly RAID_IDS: Partial<Record<RaidFilterOption, number>> = {
		[RaidFilterOption.RaidIcecrownCitadel]: 4812,
		[RaidFilterOption.RaidRubySanctum]: 4987,
		[RaidFilterOption.RaidBlackwingDescent]: 5094,
		[RaidFilterOption.RaidTheBastionOfTwilight]: 5334,
		[RaidFilterOption.RaidBaradinHold]: 5600,
		[RaidFilterOption.RaidThroneOfTheFourWinds]: 5638,
		[RaidFilterOption.RaidFirelands]: 5723,
		[RaidFilterOption.RaidDragonSoul]: 5892,
	};

	get armorSpecializationArmorType() {
		// We always pick the first entry since this is always the preffered armor type
		return this.playerClass.armorTypes[0];
	}

	hasArmorSpecializationBonus() {
		return [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
		].some(itemSlot => {
			const item = this.getEquippedItem(itemSlot)?.item;
			if (!item) return false;
			const armorType = item.armorType;
			return armorType !== this.armorSpecializationArmorType;
		});
	}

	filterItemData<T>(itemData: Array<T>, getItemFunc: (val: T) => Item, slot: ItemSlot): Array<T> {
		const filters = this.sim.getFilters();

		const filterItems = (itemData: Array<T>, filterFunc: (item: Item) => boolean) => {
			return itemData.filter(itemElem => filterFunc(getItemFunc(itemElem)));
		};

		if (filters.minIlvl != 0) {
			itemData = filterItems(itemData, item => (item.scalingOptions?.[ItemLevelState.Base].ilvl || item.ilvl) >= filters.minIlvl);
		}
		if (filters.maxIlvl != 0) {
			itemData = filterItems(itemData, item => (item.scalingOptions?.[ItemLevelState.Base].ilvl || item.ilvl) <= filters.maxIlvl);
		}

		if (filters.factionRestriction != UIItem_FactionRestriction.UNSPECIFIED) {
			itemData = filterItems(
				itemData,
				item => item.factionRestriction == filters.factionRestriction || item.factionRestriction == UIItem_FactionRestriction.UNSPECIFIED,
			);
		}

		if (!filters.sources.includes(SourceFilterOption.SourceCrafting)) {
			itemData = filterItems(itemData, item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'crafted'));
		}
		if (!filters.sources.includes(SourceFilterOption.SourceQuest)) {
			itemData = filterItems(itemData, item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'quest'));
		}
		if (!filters.sources.includes(SourceFilterOption.SourceReputation)) {
			itemData = filterItems(itemData, item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'rep'));
		}
		if (!filters.sources.includes(SourceFilterOption.SourcePvp)) {
			itemData = filterItems(itemData, item => !isPVPItem(item));
		}

		for (const [srcOptionStr, difficulty] of Object.entries(Player.DIFFICULTY_SRCS)) {
			const srcOption = parseInt(srcOptionStr) as SourceFilterOption;
			if (!filters.sources.includes(srcOption)) {
				itemData = filterItems(
					itemData,
					item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'drop' && itemSrc.source.drop.difficulty == difficulty),
				);

				if (difficulty == DungeonDifficulty.DifficultyRaid10H || difficulty == DungeonDifficulty.DifficultyRaid25H) {
					const normalDifficulty = Player.HEROIC_TO_NORMAL[difficulty];
					itemData = filterItems(
						itemData,
						item =>
							!item.sources.some(
								itemSrc =>
									itemSrc.source.oneofKind == 'drop' &&
									itemSrc.source.drop.difficulty == normalDifficulty &&
									itemSrc.source.drop.category == AL_CATEGORY_HARD_MODE,
							),
					);
				}
			}
		}

		for (const [raidOptionStr, zoneId] of Object.entries(Player.RAID_IDS)) {
			const raidOption = parseInt(raidOptionStr) as RaidFilterOption;
			if (!filters.raids.includes(raidOption)) {
				itemData = filterItems(
					itemData,
					item => !item.sources.some(itemSrc => itemSrc.source.oneofKind == 'drop' && itemSrc.source.drop.zoneId == zoneId),
				);
			}
		}

		if (Player.ARMOR_SLOTS.includes(slot)) {
			itemData = filterItems(itemData, item => {
				if (!filters.armorTypes.includes(item.armorType)) {
					return false;
				}

				return true;
			});
		} else if (Player.WEAPON_SLOTS.includes(slot)) {
			itemData = filterItems(itemData, item => {
				if (item.handType == HandType.HandTypeUnknown && item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeUnknown) {
					return false;
				}

				if (!filters.weaponTypes.includes(item.weaponType) && item.handType > HandType.HandTypeUnknown) {
					return false;
				}

				if (!filters.oneHandedWeapons && item.handType != HandType.HandTypeTwoHand) {
					return false;
				}
				if (!filters.twoHandedWeapons && item.handType == HandType.HandTypeTwoHand) {
					return false;
				}

				// Ranged weapons are equiped in MH slot from MoP onwards
				if (!filters.rangedWeaponTypes.includes(item.rangedWeaponType) && item.rangedWeaponType > RangedWeaponType.RangedWeaponTypeUnknown) {
					return false;
				}

				let minSpeed = slot == ItemSlot.ItemSlotMainHand ? filters.minMhWeaponSpeed : filters.minOhWeaponSpeed;
				let maxSpeed = slot == ItemSlot.ItemSlotMainHand ? filters.maxMhWeaponSpeed : filters.maxOhWeaponSpeed;
				if (item.rangedWeaponType > 0) {
					minSpeed = filters.minRangedWeaponSpeed;
					maxSpeed = filters.maxRangedWeaponSpeed;
				}

				if (minSpeed > 0 && item.weaponSpeed < minSpeed) {
					return false;
				}
				if (maxSpeed > 0 && item.weaponSpeed > maxSpeed) {
					return false;
				}

				return true;
			});
		}

		return itemData;
	}

	filterEnchantData<T>(enchantData: Array<T>, getEnchantFunc: (val: T) => Enchant, slot: ItemSlot, currentEquippedItem: EquippedItem | null): Array<T> {
		if (!currentEquippedItem) {
			return enchantData;
		}

		//const filters = this.sim.getFilters();

		return enchantData.filter(enchantElem => {
			const enchant = getEnchantFunc(enchantElem);

			if (!enchantAppliesToItem(enchant, currentEquippedItem.item)) {
				return false;
			}

			return true;
		});
	}

	filterGemData<T>(gemData: Array<T>, getGemFunc: (val: T) => Gem, slot: ItemSlot, socketColor: GemColor): Array<T> {
		const filters = this.sim.getFilters();

		const isJewelcrafting = this.hasProfession(Profession.Jewelcrafting);
		return gemData.filter(gemElem => {
			const gem = getGemFunc(gemElem);
			if (!isJewelcrafting && gem.requiredProfession == Profession.Jewelcrafting) {
				return false;
			}

			if (filters.matchingGemsOnly && !gemMatchesSocket(gem, socketColor)) {
				return false;
			}

			// This is not exactly a player selected filter, just a general filter to remove any gems with stats that is not in use for the player.
			// i.e dead gems.
			const statsFilter = this.specConfig.gemStats ?? this.specConfig.epStats;
			const positiveStatIds = gem.stats.map((value, statId) => (value > 0 ? statId : -1)).filter(statId => statId >= 0);
			if (!positiveStatIds.length) {
				return false;
			}
			return !positiveStatIds.some(statId => !statsFilter.includes(statId));
		});
	}

	makeUnitReference(): UnitReference {
		if (this.party == null) {
			return emptyUnitReference();
		} else {
			return newUnitReference(this.getRaidIndex());
		}
	}

	private toDatabase(): SimDatabase {
		const dbGear = this.getGear().toDatabase(this.sim.db);
		const dbItemSwapGear = this.itemSwapSettings.getGear().toDatabase(this.sim.db);
		return Database.mergeSimDatabases(dbGear, dbItemSwapGear);
	}

	toProto(forExport?: boolean, forSimming?: boolean, exportCategories?: Array<SimSettingCategories>): PlayerProto {
		const exportCategory = (cat: SimSettingCategories) => !exportCategories || exportCategories.length == 0 || exportCategories.includes(cat);

		const gear = this.getGear();
		const aplRotation = forSimming
			? this.getResolvedAplRotation()
			: // When exporting we want to omit the uuid field to prevent bloat
			  omitDeep(this.aplRotation, ['uuid']);

		let player = PlayerProto.create({
			class: this.getClass(),
			database: forExport ? undefined : this.toDatabase(),
		});
		if (exportCategory(SimSettingCategories.Gear)) {
			PlayerProto.mergePartial(player, {
				equipment: gear.asSpec(),
				bonusStats: this.getBonusStats().toProto(),
				enableItemSwap: this.itemSwapSettings.getEnableItemSwap(),
				itemSwap: this.itemSwapSettings.toProto(),
			});
		}
		if (exportCategory(SimSettingCategories.Talents)) {
			PlayerProto.mergePartial(player, {
				talentsString: this.getTalentsString(),
				glyphs: this.getGlyphs(),
			});
		}
		if (exportCategory(SimSettingCategories.Rotation)) {
			PlayerProto.mergePartial(player, {
				cooldowns: Cooldowns.create({
					hpPercentForDefensives: this.getSimpleCooldowns().hpPercentForDefensives,
				}),
				rotation: aplRotation,
			});
		}
		if (exportCategory(SimSettingCategories.Consumes)) {
			PlayerProto.mergePartial(player, {
				consumables: this.getConsumes(),
			});
		}
		if (exportCategory(SimSettingCategories.Miscellaneous)) {
			PlayerProto.mergePartial(player, {
				name: this.getName(),
				race: this.getRace(),
				profession1: this.getProfession1(),
				profession2: this.getProfession2(),
				reactionTimeMs: this.getReactionTime(),
				channelClipDelayMs: this.getChannelClipDelay(),
				inFrontOfTarget: this.getInFrontOfTarget(),
				distanceFromTarget: this.getDistanceFromTarget(),
				healingModel: this.getHealingModel(),
				challengeMode: this.getChallengeModeEnabled(),
			});
			player = withSpec(this.getSpec(), player, this.getSpecOptions());
		}
		if (exportCategory(SimSettingCategories.External)) {
			PlayerProto.mergePartial(player, {
				buffs: this.getBuffs(),
			});
		}
		return player;
	}

	fromProto(eventID: EventID, proto: PlayerProto, includeCategories?: Array<SimSettingCategories>) {
		// Fix potential out-of-date protos before importing
		TypedEvent.freezeAllAndDo(() => {
			Player.updateProtoVersion(proto);
			const loadCategory = (cat: SimSettingCategories) => !includeCategories || includeCategories.length == 0 || includeCategories.includes(cat);
			eventID = TypedEvent.nextEventID();
			if (loadCategory(SimSettingCategories.Gear)) {
				this.setGear(eventID, proto.equipment ? this.sim.db.lookupEquipmentSpec(proto.equipment) : new Gear({}));
				this.itemSwapSettings.setItemSwapSettings(
					eventID,
					proto.enableItemSwap,
					proto.itemSwap ? this.sim.db.lookupItemSwap(proto.itemSwap) : new ItemSwapGear({}),
					Stats.fromProto(proto.itemSwap?.prepullBonusStats),
				);
				this.setBonusStats(eventID, Stats.fromProto(proto.bonusStats || UnitStats.create()));
				//this.setBulkEquipmentSpec(eventID, BulkEquipmentSpec.create()); // Do not persist the bulk equipment settings.
			}
			if (loadCategory(SimSettingCategories.Talents)) {
				this.setTalentsString(eventID, proto.talentsString);
				this.setGlyphs(eventID, proto.glyphs || Glyphs.create());
			}
			if (loadCategory(SimSettingCategories.Rotation)) {
				if (proto.rotation?.type == APLRotationType.TypeUnknown) {
					if (!proto.rotation) {
						proto.rotation = APLRotation.create();
					}
					proto.rotation.type = APLRotationType.TypeAuto;
				}
				this.setAplRotation(eventID, proto.rotation || APLRotation.create());
			}
			if (loadCategory(SimSettingCategories.Consumes)) {
				this.setConsumes(eventID, proto.consumables || ConsumesSpec.create());
			}
			if (loadCategory(SimSettingCategories.Miscellaneous)) {
				this.setSpecOptions(eventID, this.specTypeFunctions.optionsFromPlayer(proto));
				this.setName(eventID, proto.name);
				this.setRace(eventID, proto.race);
				this.setProfession1(eventID, proto.profession1);
				this.setProfession2(eventID, proto.profession2);
				this.setReactionTime(eventID, proto.reactionTimeMs);
				this.setChannelClipDelay(eventID, proto.channelClipDelayMs);
				this.setInFrontOfTarget(eventID, proto.inFrontOfTarget);
				this.setDistanceFromTarget(eventID, proto.distanceFromTarget);
				this.setHealingModel(eventID, proto.healingModel || HealingModel.create());
				this.setChallengeModeEnabled(eventID, proto.challengeMode);
			}
			if (loadCategory(SimSettingCategories.External)) {
				this.setBuffs(eventID, proto.buffs || IndividualBuffs.create());
			}
		});
	}

	clone(eventID: EventID): Player<SpecType> {
		const newPlayer = new Player<SpecType>(this.playerSpec, this.sim);
		newPlayer.fromProto(eventID, this.toProto());
		return newPlayer;
	}

	applySharedDefaults(eventID: EventID) {
		TypedEvent.freezeAllAndDo(() => {
			this.setReactionTime(eventID, 100);
			this.setInFrontOfTarget(eventID, this.playerSpec.isTankSpec);
			this.setHealingModel(
				eventID,
				HealingModel.create({
					burstWindow: this.playerSpec.isTankSpec ? 6 : 0,
				}),
			);
			this.setSimpleCooldowns(
				eventID,
				Cooldowns.create({
					hpPercentForDefensives: this.playerSpec.isTankSpec ? 0.4 : 0,
				}),
			);
			this.setBonusStats(eventID, new Stats());
		});
	}

	getBaseMastery(): number {
		return 8;
	}

	getMasteryPerPointModifier(): number {
		return Mechanics.masteryPercentPerPoint.get(this.getSpec()) || 0;
	}
	static updateProtoVersion(proto: PlayerProto) {
		if (!(proto.apiVersion < CURRENT_API_VERSION)) {
			return;
		}
	}
}
