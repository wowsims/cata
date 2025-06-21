import * as Mechanics from './constants/mechanics';
import { CURRENT_API_VERSION } from './constants/other';
import { UnitMetadataList } from './player';
import { Encounter as EncounterProto, MobType, PresetEncounter, PresetTarget, SpellSchool, Stat, Target as TargetProto, TargetInput } from './proto/common';
import { Stats } from './proto_utils/stats';
import { Sim } from './sim';
import { EventID, TypedEvent } from './typed_event';

// Manages all the settings for an Encounter.
export class Encounter {
	readonly sim: Sim;

	private duration = 300;
	private durationVariation = 60;
	private executeProportion20 = 0.2;
	private executeProportion25 = 0.25;
	private executeProportion35 = 0.35;
	private executeProportion45 = 0.45;
	private executeProportion90 = 0.9;
	private useHealth = false;
	targets: Array<TargetProto>;
	targetsMetadata: UnitMetadataList;

	readonly targetsChangeEmitter = new TypedEvent<void>();
	readonly durationChangeEmitter = new TypedEvent<void>();
	readonly executeProportionChangeEmitter = new TypedEvent<void>();

	// Emits when any of the above emitters emit.
	readonly changeEmitter = new TypedEvent<void>();

	constructor(sim: Sim) {
		this.sim = sim;
		this.targets = [Encounter.defaultTargetProto()];
		this.targetsMetadata = new UnitMetadataList();

		[this.targetsChangeEmitter, this.durationChangeEmitter, this.executeProportionChangeEmitter].forEach(emitter =>
			emitter.on(eventID => this.changeEmitter.emit(eventID)),
		);
	}

	get primaryTarget(): TargetProto {
		return this.targets[0];
	}

	getDurationVariation(): number {
		return this.durationVariation;
	}
	setDurationVariation(eventID: EventID, newDuration: number) {
		if (newDuration == this.durationVariation) return;

		this.durationVariation = newDuration;
		this.durationChangeEmitter.emit(eventID);
	}

	getDuration(): number {
		return this.duration;
	}
	setDuration(eventID: EventID, newDuration: number) {
		if (newDuration == this.duration) return;

		this.duration = newDuration;
		this.durationChangeEmitter.emit(eventID);
	}

	getExecuteProportion20(): number {
		return this.executeProportion20;
	}
	setExecuteProportion20(eventID: EventID, newExecuteProportion20: number) {
		if (newExecuteProportion20 == this.executeProportion20) return;

		this.executeProportion20 = newExecuteProportion20;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion25(): number {
		return this.executeProportion25;
	}
	setExecuteProportion25(eventID: EventID, newExecuteProportion25: number) {
		if (newExecuteProportion25 == this.executeProportion25) return;

		this.executeProportion25 = newExecuteProportion25;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion35(): number {
		return this.executeProportion35;
	}
	setExecuteProportion35(eventID: EventID, newExecuteProportion35: number) {
		if (newExecuteProportion35 == this.executeProportion35) return;

		this.executeProportion35 = newExecuteProportion35;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion45(): number {
		return this.executeProportion45;
	}
	setExecuteProportion45(eventID: EventID, newExecuteProportion45: number) {
		if (newExecuteProportion45 == this.executeProportion45) return;

		this.executeProportion45 = newExecuteProportion45;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getExecuteProportion90(): number {
		return this.executeProportion90;
	}
	setExecuteProportion90(eventID: EventID, newExecuteProportion90: number) {
		if (newExecuteProportion90 == this.executeProportion90) return;

		this.executeProportion90 = newExecuteProportion90;
		this.executeProportionChangeEmitter.emit(eventID);
	}
	getUseHealth(): boolean {
		return this.useHealth;
	}
	setUseHealth(eventID: EventID, newUseHealth: boolean) {
		if (newUseHealth == this.useHealth) return;

		this.useHealth = newUseHealth;
		this.durationChangeEmitter.emit(eventID);
		this.executeProportionChangeEmitter.emit(eventID);
	}

	matchesPreset(preset: PresetEncounter): boolean {
		return preset.targets.length == this.targets.length && this.targets.every((t, i) => TargetProto.equals(t, preset.targets[i].target));
	}

	applyPreset(eventID: EventID, preset: PresetEncounter) {
		this.targets = preset.targets.map(presetTarget => presetTarget.target || TargetProto.create());
		this.targetsChangeEmitter.emit(eventID);
	}

	applyPresetTarget(eventID: EventID, preset: PresetTarget, index: number) {
		this.targets[index] = preset.target || TargetProto.create();
		this.targetsChangeEmitter.emit(eventID);
	}

	toProto(): EncounterProto {
		return EncounterProto.create({
			duration: this.duration,
			durationVariation: this.durationVariation,
			executeProportion20: this.executeProportion20,
			executeProportion25: this.executeProportion25,
			executeProportion35: this.executeProportion35,
			executeProportion45: this.executeProportion45,
			executeProportion90: this.executeProportion90,
			useHealth: this.useHealth,
			targets: this.targets,
			apiVersion: CURRENT_API_VERSION,
		});
	}

	fromProto(eventID: EventID, proto: EncounterProto) {
		// Fix out-of-date protos before importing
		Encounter.updateProtoVersion(proto);

		TypedEvent.freezeAllAndDo(() => {
			this.setDuration(eventID, proto.duration);
			this.setDurationVariation(eventID, proto.durationVariation);
			this.setExecuteProportion20(eventID, proto.executeProportion20);
			this.setExecuteProportion25(eventID, proto.executeProportion25);
			this.setExecuteProportion35(eventID, proto.executeProportion35);
			this.setExecuteProportion45(eventID, proto.executeProportion45);
			this.setExecuteProportion90(eventID, proto.executeProportion90);
			this.setUseHealth(eventID, proto.useHealth);
			this.targets = proto.targets;
			this.targetsChangeEmitter.emit(eventID);
		});
	}

	applyDefaults(eventID: EventID) {
		this.fromProto(
			eventID,
			EncounterProto.create({
				duration: 300,
				durationVariation: 60,
				executeProportion20: 0.2,
				executeProportion25: 0.25,
				executeProportion35: 0.35,
				executeProportion45: 0.45,
				executeProportion90: 0.9,
				targets: [Encounter.defaultTargetProto()],
				apiVersion: CURRENT_API_VERSION,
			}),
		);
	}

	static defaultTargetProto(): TargetProto {
		// Copy default raid target used as fallback for missing DB.
		// https://github.com/wowsims/mop/blob/3570c4fcf1a4e2cd81926019d4a1b3182f613de1/sim/encounters/register_all.go#L24
		return TargetProto.create({
			id: 31146,
			name: 'Raid Target',
			level: Mechanics.BOSS_LEVEL,
			mobType: MobType.MobTypeMechanical,
			stats: Stats.fromMap({
				[Stat.StatArmor]: 24835,
				[Stat.StatAttackPower]: 650,
				[Stat.StatHealth]: 120016403,
			}).asProtoArray(),
			minBaseDamage: 550000,
			damageSpread: 0.4,
			tankIndex: 0,
			swingSpeed: 2,
			suppressDodge: false,
			parryHaste: false,
			dualWield: false,
			dualWieldPenalty: false,
			spellSchool: SpellSchool.SpellSchoolPhysical,
			targetInputs: new Array<TargetInput>(0),
		});
	}

	static updateProtoVersion(proto: EncounterProto) {
		if (!(proto.apiVersion < CURRENT_API_VERSION)) {
			return;
		}
	}
}
