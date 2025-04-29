import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Component } from '../core/components/component.js';
import { Player } from '../core/player.js';
import { PlayerClasses } from '../core/player_classes/index.js';
import { PlayerSpecs } from '../core/player_specs/index.js';
import { Class, RaidBuffs, Spec } from '../core/proto/common.js';
import { HunterOptions_PetType as HunterPetType } from '../core/proto/hunter.js';
import { PaladinAura } from '../core/proto/paladin.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem } from '../core/proto/shaman.js';
import { WarlockOptions_Summon as WarlockSummon } from '../core/proto/warlock.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { ClassSpecs, SpecTalents, textCssClassForClass } from '../core/proto_utils/utils.js';
import { Raid } from '../core/raid.js';
import { sum } from '../core/utils.js';
import { RaidSimUI } from './raid_sim_ui.js';

interface RaidStatsOptions {
	sections: Array<RaidStatsSectionOptions>;
}

interface RaidStatsSectionOptions {
	label: string;
	categories: Array<RaidStatsCategoryOptions>;
}

interface RaidStatsCategoryOptions {
	label: string;
	effects: Array<RaidStatsEffectOptions>;
}

type PlayerProvider = { class?: Class; condition: (player: Player<any>) => boolean };
type RaidProvider = (raid: Raid) => boolean;

interface RaidStatsEffectOptions {
	label: string;
	actionId?: ActionId;
	playerData?: PlayerProvider;
	raidData?: RaidProvider;
}

export class RaidStats extends Component {
	private readonly categories: Array<RaidStatsCategory>;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI) {
		super(parent, 'raid-stats');

		const categories: Array<RaidStatsCategory> = [];
		RAID_STATS_OPTIONS.sections.forEach(section => {
			const contentElemRef = ref<HTMLDivElement>();

			const sectionElem = (
				<div className="raid-stats-section">
					<div className="raid-stats-section-header">
						<label className="raid-stats-section-label form-label">{section.label}</label>
					</div>
					<div ref={contentElemRef} className="raid-stats-section-content"></div>
				</div>
			);
			this.rootElem.appendChild(sectionElem);

			const contentElem = contentElemRef.value!;
			section.categories.forEach(categoryOptions => {
				categories.push(new RaidStatsCategory(contentElem, raidSimUI, categoryOptions));
			});
		});
		this.categories = categories;

		raidSimUI.changeEmitter.on(_eventID => this.categories.forEach(c => c.update()));
	}
}

class RaidStatsCategory extends Component {
	readonly raidSimUI: RaidSimUI;
	private readonly options: RaidStatsCategoryOptions;
	private readonly effects: Array<RaidStatsEffect>;
	private readonly counterElem: HTMLElement;
	private readonly tooltipElem: HTMLElement;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI, options: RaidStatsCategoryOptions) {
		super(parent, 'raid-stats-category-root');
		this.raidSimUI = raidSimUI;
		this.options = options;

		const counterElemRef = ref<HTMLElement>();
		const categoryElemRef = ref<HTMLButtonElement>();
		this.rootElem.appendChild(
			<button ref={categoryElemRef} className="raid-stats-category">
				<span ref={counterElemRef} className="raid-stats-category-counter"></span>
				<span className="raid-stats-category-label">{options.label}</span>
			</button>,
		);

		this.counterElem = counterElemRef.value!;
		this.tooltipElem = (
			<div>
				<label className="raid-stats-category-label">{options.label}</label>
			</div>
		) as HTMLElement;

		this.effects = options.effects.map(opt => new RaidStatsEffect(this.tooltipElem, raidSimUI, opt));

		if (options.effects.length != 1 || options.effects[0].playerData?.class) {
			const statsLink = categoryElemRef.value!;

			// Using the title option here because outerHTML sanitizes and filters out the img src options
			tippy(statsLink, {
				theme: 'raid-stats-category-tooltip',
				placement: 'right',
				content: this.tooltipElem,
			});
		}
	}

	update() {
		this.effects.forEach(effect => effect.update());

		const total = sum(this.effects.map(effect => effect.count));
		this.counterElem.textContent = String(total);

		const statsLink = this.rootElem.querySelector<HTMLElement>('.raid-stats-category')!;
		statsLink?.classList[total === 0 ? 'remove' : 'add']('active');
	}
}

class RaidStatsEffect extends Component {
	readonly raidSimUI: RaidSimUI;
	private readonly options: RaidStatsEffectOptions;
	private readonly counterElem: HTMLElement;

	curPlayers: Array<Player<any>>;
	count: number;

	constructor(parent: HTMLElement, raidSimUI: RaidSimUI, options: RaidStatsEffectOptions) {
		super(parent, 'raid-stats-effect');
		this.raidSimUI = raidSimUI;
		this.options = options;

		this.curPlayers = [];
		this.count = 0;

		const counterElemRef = ref<HTMLElement>();
		const labelElemRef = ref<HTMLElement>();
		const iconElemRef = ref<HTMLImageElement>();
		this.rootElem.appendChild(
			<>
				<span ref={counterElemRef} className="raid-stats-effect-counter"></span>
				<img ref={iconElemRef} className="raid-stats-effect-icon"></img>
				<span ref={labelElemRef} className="raid-stats-effect-label">
					{options.label}
				</span>
			</>,
		);

		this.counterElem = counterElemRef.value!;

		if (this.options.playerData?.class) {
			const playerCssClass = textCssClassForClass(PlayerClasses.fromProto(this.options.playerData.class));
			labelElemRef.value!.classList.add(playerCssClass);
		}

		if (options.actionId) {
			options.actionId.fill().then(actionId => (iconElemRef.value!.src = actionId.iconUrl));
		} else {
			iconElemRef.value!.remove();
		}
	}

	update() {
		if (this.options.playerData) {
			this.curPlayers = this.raidSimUI.getActivePlayers().filter(p => this.options.playerData!.condition(p));
		}

		const raidData = this.options.raidData && this.options.raidData(this.raidSimUI.sim.raid);

		this.count = this.curPlayers.length + (raidData ? 1 : 0);
		this.counterElem.textContent = String(this.count);
		if (this.count == 0) {
			this.rootElem.classList.remove('active');
		} else {
			this.rootElem.classList.add('active');
		}
	}
}

function negateIf(val: boolean, cond: boolean): boolean {
	return cond ? !val : val;
}

function playerClass<T extends Class>(clazz: T, extraCondition?: (player: Player<ClassSpecs<T>>) => boolean): PlayerProvider {
	return {
		class: clazz,
		condition: (player: Player<any>): boolean => {
			return player.isClass(clazz) && (!extraCondition || extraCondition(player));
		},
	};
}
function playerClassAndTalentInternal<T extends Class>(
	clazz: T,
	talentName: keyof SpecTalents<ClassSpecs<T>>,
	negateTalent: boolean,
	extraCondition?: (player: Player<ClassSpecs<T>>) => boolean,
): PlayerProvider {
	return {
		class: clazz,
		condition: (player: Player<any>): boolean => {
			return (
				player.isClass(clazz) &&
				negateIf(Boolean((player.getTalents() as any)[talentName]), negateTalent) &&
				(!extraCondition || extraCondition(player))
			);
		},
	};
}
function playerClassAndTalent<T extends Class>(
	clazz: T,
	talentName: keyof SpecTalents<ClassSpecs<T>>,
	extraCondition?: (player: Player<ClassSpecs<T>>) => boolean,
): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, false, extraCondition);
}
function playerClassAndMissingTalent<T extends Class>(
	clazz: T,
	talentName: keyof SpecTalents<ClassSpecs<T>>,
	extraCondition?: (player: Player<ClassSpecs<T>>) => boolean,
): PlayerProvider {
	return playerClassAndTalentInternal(clazz, talentName, true, extraCondition);
}
function playerSpecAndTalentInternal<T extends Spec>(
	spec: T,
	talentName: keyof SpecTalents<T>,
	negateTalent: boolean,
	extraCondition?: (player: Player<T>) => boolean,
): PlayerProvider {
	return {
		class: PlayerSpecs.fromProto(spec).classID,
		condition: (player: Player<any>): boolean => {
			return (
				player.isSpec(spec) && negateIf(Boolean((player.getTalents() as any)[talentName]), negateTalent) && (!extraCondition || extraCondition(player))
			);
		},
	};
}
function playerSpecAndTalent<T extends Spec>(spec: T, talentName: keyof SpecTalents<T>, extraCondition?: (player: Player<T>) => boolean): PlayerProvider {
	return playerSpecAndTalentInternal(spec, talentName, false, extraCondition);
}
function playerSpecAndMissingTalent<T extends Spec>(
	spec: T,
	talentName: keyof SpecTalents<T>,
	extraCondition?: (player: Player<T>) => boolean,
): PlayerProvider {
	return playerSpecAndTalentInternal(spec, talentName, true, extraCondition);
}

function raidBuff(buffName: keyof RaidBuffs): RaidProvider {
	return (raid: Raid): boolean => {
		return Boolean(raid.getBuffs()[buffName]);
	};
}

const RAID_STATS_OPTIONS: RaidStatsOptions = {
	sections: [
		{
			label: 'Roles',
			categories: [
				{
					label: 'Tanks',
					effects: [
						{
							label: 'Tanks',
							playerData: {
								condition: player => player.getPlayerSpec().isTankSpec,
							},
						},
					],
				},
				{
					label: 'Healers',
					effects: [
						{
							label: 'Healers',
							playerData: { condition: player => player.getPlayerSpec().isHealingSpec },
						},
					],
				},
				{
					label: 'Melee',
					effects: [
						{
							label: 'Melee',
							playerData: { condition: player => player.getPlayerSpec().isMeleeDpsSpec },
						},
					],
				},
				{
					label: 'Ranged',
					effects: [
						{
							label: 'Ranged',
							playerData: { condition: player => player.getPlayerSpec().isRangedDpsSpec },
						},
					],
				},
			],
		},
		{
			label: 'Buffs',
			categories: [
				{
					label: 'Bloodlust',
					effects: [
						{
							label: 'Bloodlust',
							actionId: ActionId.fromSpellId(2825),
							playerData: playerClass(Class.ClassShaman),
						},
						{
							label: 'Time Warp',
							actionId: ActionId.fromSpellId(80353),
							playerData: playerClass(Class.ClassMage),
						},
					],
				},
				{
					label: 'Stats %',
					effects: [
						{
							label: 'Mark of the Wild',
							actionId: ActionId.fromSpellId(1126),
							playerData: playerClass(Class.ClassDruid),
						},
						{
							label: 'Blessing of Kings',
							actionId: ActionId.fromSpellId(20217),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Drums of the Burning Wild',
							actionId: ActionId.fromItemId(63140),
							raidData: raidBuff('drumsOfTheBurningWild'),
						},
					],
				},
				{
					label: 'Strength/Agility',
					effects: [
						{
							label: 'Strength of Earth Totem',
							actionId: ActionId.fromSpellId(8075),
							playerData: playerClass(
								Class.ClassShaman,
								player => player.getSpecOptions().classOptions?.totems?.earth == EarthTotem.StrengthOfEarthTotem,
							),
						},
						{
							label: 'Horn of Winter',
							actionId: ActionId.fromSpellId(57330),
							playerData: playerClass(Class.ClassDeathKnight),
						},
						{
							label: 'Battle Shout',
							actionId: ActionId.fromSpellId(6673),
							playerData: playerClass(Class.ClassWarrior),
						},
					],
				},
				{
					label: 'Armor',
					effects: [
						{
							label: 'Devotion Aura',
							actionId: ActionId.fromSpellId(465),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Stoneskin Totem',
							actionId: ActionId.fromSpellId(8071),
							playerData: playerClass(
								Class.ClassShaman,
								player => player.getSpecOptions().classOptions?.totems?.earth == EarthTotem.StoneskinTotem,
							),
						},
					],
				},
				{
					label: 'Attack Power %',
					effects: [
						{
							label: 'Blessing of Might',
							actionId: ActionId.fromSpellId(19740),
							playerData: playerClass(Class.ClassPaladin),
						},
						// {
						// 	label: 'Abominations Might',
						// 	actionId: ActionId.fromSpellId(53138),
						// 	playerData: playerClassAndTalent(Class.ClassDeathKnight, 'abominationsMight'),
						// },
						// {
						// 	label: 'Unleashed Rage',
						// 	actionId: ActionId.fromSpellId(30808),
						// 	playerData: playerClassAndTalent(Class.ClassShaman, 'unleashedRage'),
						// },
						{
							label: 'Trueshot Aura',
							actionId: ActionId.fromSpellId(19506),
							playerData: playerClass(Class.ClassHunter),
						},
					],
				},
				{
					label: 'Spell Power',
					effects: [
						// {
						// 	label: 'Demonic Pact',
						// 	actionId: ActionId.fromSpellId(47236),
						// 	playerData: playerClassAndTalent(Class.ClassWarlock, 'demonicPact'),
						// },
						// {
						// 	label: 'Totemic Wrath',
						// 	actionId: ActionId.fromSpellId(77746),
						// 	playerData: playerClassAndTalent(Class.ClassShaman, 'totemicWrath'),
						// },
						{
							label: 'Arcane Brilliance',
							actionId: ActionId.fromSpellId(1459),
							playerData: playerClass(Class.ClassMage),
						},
						{
							label: 'Flametongue Totem',
							actionId: ActionId.fromSpellId(8227),
							playerData: playerClass(
								Class.ClassShaman,
								player => player.getSpecOptions().classOptions?.totems?.fire == FireTotem.FlametongueTotem,
							),
						},
					],
				},
				{
					label: '+3% Damage',
					effects: [
						// {
						// 	label: 'Communion',
						// 	actionId: ActionId.fromSpellId(31876),
						// 	playerData: playerClassAndTalent(Class.ClassPaladin, 'communion'),
						// },
						// {
						// 	label: 'Arcane Tactics',
						// 	actionId: ActionId.fromSpellId(82930),
						// 	playerData: playerClassAndTalent(Class.ClassMage, 'arcaneTactics'),
						// },
						// {
						// 	label: 'Ferocious Inspiration',
						// 	actionId: ActionId.fromSpellId(34460),
						// 	playerData: playerClassAndTalent(Class.ClassHunter, 'ferociousInspiration'),
						// },
					],
				},
				{
					label: 'Melee Haste',
					effects: [
						// {
						// 	label: 'Icy Talons',
						// 	actionId: ActionId.fromSpellId(55610),
						// 	playerData: playerClassAndTalent(Class.ClassDeathKnight, 'improvedIcyTalons'),
						// },
						// {
						// 	label: 'Hunting Party',
						// 	actionId: ActionId.fromSpellId(53290),
						// 	playerData: playerClassAndTalent(Class.ClassHunter, 'huntingParty'),
						// },
						{
							label: 'Windfury Totem',
							actionId: ActionId.fromSpellId(8512),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().classOptions?.totems?.air == AirTotem.WindfuryTotem),
						},
					],
				},
				{
					label: 'Spell Haste',
					effects: [
						// {
						// 	label: 'Shadow Form',
						// 	actionId: ActionId.fromSpellId(15473),
						// 	playerData: playerClassAndTalent(Class.ClassPriest, 'shadowform'),
						// },
						// {
						// 	label: 'Moonkin Form',
						// 	actionId: ActionId.fromSpellId(24858),
						// 	playerData: playerClassAndTalent(Class.ClassDruid, 'moonkinForm'),
						// },
						{
							label: 'Wrath of Air Totem',
							actionId: ActionId.fromSpellId(3738),
							playerData: playerClass(Class.ClassShaman, player => player.getSpecOptions().classOptions?.totems?.air == AirTotem.WrathOfAirTotem),
						},
					],
				},
				{
					label: '+5% Crit',
					effects: [
						// {
						// 	label: 'Leader of the Pack',
						// 	actionId: ActionId.fromSpellId(17007),
						// 	playerData: playerClassAndTalent(Class.ClassDruid, 'leaderOfThePack'),
						// },
						// {
						// 	label: 'Elemental Oath',
						// 	actionId: ActionId.fromSpellId(51470),
						// 	playerData: playerClassAndTalent(Class.ClassShaman, 'elementalOath'),
						// },
						// {
						// 	label: 'Honor Among Thieves',
						// 	actionId: ActionId.fromSpellId(51701),
						// 	playerData: playerClassAndTalent(Class.ClassRogue, 'honorAmongThieves'),
						// },
						// {
						// 	label: 'Rampage',
						// 	actionId: ActionId.fromSpellId(29801),
						// 	playerData: playerClassAndTalent(Class.ClassWarrior, 'rampage'),
						// },
					],
				},
				{
					label: 'Mana',
					effects: [
						{
							label: 'Arcane Brilliance',
							actionId: ActionId.fromSpellId(1459),
							playerData: playerClass(Class.ClassMage),
						},
						{
							label: 'Fel Intelligence',
							actionId: ActionId.fromSpellId(54424),
							playerData: playerClass(Class.ClassWarlock),
						},
					],
				},
				{
					label: 'MP5',
					effects: [
						{
							label: 'Blessing of Might',
							actionId: ActionId.fromSpellId(19740),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Fel Intelligence',
							actionId: ActionId.fromSpellId(54424),
							playerData: playerClass(Class.ClassWarlock),
						},
						{
							label: 'Mana Spring Totem',
							actionId: ActionId.fromSpellId(5675),
							playerData: playerClass(
								Class.ClassShaman,
								player => player.getSpecOptions().classOptions?.totems?.water == WaterTotem.ManaSpringTotem,
							),
						},
					],
				},
				{
					label: 'Replenishment',
					effects: [
						// {
						// 	label: 'Vampiric Touch',
						// 	actionId: ActionId.fromSpellId(34914),
						// 	playerData: playerClassAndTalent(Class.ClassPriest, 'vampiricTouch'),
						// },
						// {
						// 	label: 'Communion',
						// 	actionId: ActionId.fromSpellId(31876),
						// 	playerData: playerClassAndTalent(Class.ClassPaladin, 'communion'),
						// },
						// {
						// 	label: 'Revitalize',
						// 	actionId: ActionId.fromSpellId(48544),
						// 	playerData: playerClassAndTalent(Class.ClassDruid, 'revitalize'),
						// },
						// {
						// 	label: 'Soul Leach',
						// 	actionId: ActionId.fromSpellId(30295),
						// 	playerData: playerClassAndTalent(Class.ClassWarlock, 'soulLeech'),
						// },
						// {
						// 	label: 'Enduring Winter',
						// 	actionId: ActionId.fromSpellId(86508),
						// 	playerData: playerClassAndTalent(Class.ClassMage, 'enduringWinter'),
						// },
					],
				},
				{
					label: 'Stamina',
					effects: [
						{
							label: 'Power Word Fortitude',
							actionId: ActionId.fromSpellId(21562),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Blood Pact',
							actionId: ActionId.fromSpellId(6307),
							playerData: playerClass(Class.ClassWarlock),
						},
						{
							label: 'Commanding Shout',
							actionId: ActionId.fromSpellId(469),
							playerData: playerClass(Class.ClassWarrior),
						},
					],
				},
				{
					label: 'Resistances',
					effects: [
						{
							label: 'Resistance Aura',
							actionId: ActionId.fromSpellId(19891),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Elemental Resistance Totem',
							actionId: ActionId.fromSpellId(8184),
							playerData: playerClass(
								Class.ClassShaman,
								player => player.getSpecOptions().classOptions?.totems?.water == WaterTotem.ElementalResistanceTotem,
							),
						},
						{
							label: 'Aspect of the Wild',
							actionId: ActionId.fromSpellId(20043),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Shadow Protection',
							actionId: ActionId.fromSpellId(27683),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Blessing of Kings',
							actionId: ActionId.fromSpellId(20217),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Mark of the Wild',
							actionId: ActionId.fromSpellId(1126),
							playerData: playerClass(Class.ClassDruid),
						},
					],
				},
			],
		},
		{
			label: 'External Buffs',
			categories: [
				{
					label: 'Innervate',
					effects: [
						{
							label: 'Innervate',
							actionId: ActionId.fromSpellId(29166),
							playerData: playerClass(Class.ClassDruid),
						},
					],
				},
				{
					label: 'Power Infusion',
					effects: [
						{
							label: 'Power Infusion',
							actionId: ActionId.fromSpellId(10060),
							playerData: playerClassAndTalent(Class.ClassPriest, 'powerInfusion'),
						},
					],
				},
				{
					label: 'Focus Magic',
					effects: [
						// {
						// 	label: 'Focus Magic',
						// 	actionId: ActionId.fromSpellId(54648),
						// 	playerData: playerClassAndTalent(Class.ClassMage, 'focusMagic'),
						// },
					],
				},
				{
					label: 'Tricks of the Trade',
					effects: [
						{
							label: 'Tricks of the Trade',
							actionId: ActionId.fromSpellId(57933),
							playerData: playerClass(Class.ClassRogue),
						},
					],
				},
				{
					label: 'Dark Intent',
					effects: [
						{
							label: 'Dark Intent',
							actionId: ActionId.fromSpellId(85759),
							playerData: playerClass(Class.ClassWarlock),
						},
					],
				},
				// {
				// 	label: 'Unholy Frenzy',
				// 	effects: [
				// 		{
				// 			label: 'Unholy Frenzy',
				// 			actionId: ActionId.fromSpellId(49016),
				// 			playerData: playerClassAndTalent(Class.ClassDeathKnight, 'unholyFrenzy'),
				// 		},
				// 	],
				// },
				// {
				// 	label: 'Pain Suppression',
				// 	effects: [
				// 		{
				// 			label: 'Pain Suppression',
				// 			actionId: ActionId.fromSpellId(33206),
				// 			playerData: playerClassAndTalent(Class.ClassPriest, 'painSuppression'),
				// 		},
				// 	],
				// },
				// {
				// 	label: 'Divine Guardian',
				// 	effects: [
				// 		{
				// 			label: 'Divine Guardian',
				// 			actionId: ActionId.fromSpellId(70940),
				// 			playerData: playerClassAndTalent(Class.ClassPaladin, 'divineGuardian'),
				// 		},
				// 	],
				// },
				// {
				// 	label: 'Mana Tide',
				// 	effects: [
				// 		{
				// 			label: 'Mana Tide Totem',
				// 			actionId: ActionId.fromSpellId(16190),
				// 			playerData: playerClassAndTalent(Class.ClassShaman, 'manaTideTotem'),
				// 		},
				// 	],
				// },
			],
		},
		{
			label: 'DPS Debuffs',
			categories: [
				{
					label: '-Armor %',
					effects: [
						{
							label: 'Sunder Armor',
							actionId: ActionId.fromSpellId(7386),
							playerData: playerClass(Class.ClassWarrior),
						},
						{
							label: 'Expose Armor',
							actionId: ActionId.fromSpellId(8647),
							playerData: playerClass(Class.ClassRogue),
						},
						{
							label: 'Faerie Fire',
							actionId: ActionId.fromSpellId(770),
							playerData: playerClass(Class.ClassDruid),
						},
						{
							label: 'Corosive Spit',
							actionId: ActionId.fromSpellId(35387),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.Serpent),
						},
					],
				},
				{
					label: 'Phys Vuln',
					effects: [
						// {
						// 	label: 'Blood Frenzy',
						// 	actionId: ActionId.fromSpellId(29859),
						// 	playerData: playerClassAndTalent(Class.ClassWarrior, 'bloodFrenzy'),
						// },
						// {
						// 	label: 'Savage Combat',
						// 	actionId: ActionId.fromSpellId(58413),
						// 	playerData: playerClassAndTalent(Class.ClassRogue, 'savageCombat'),
						// },
						// {
						// 	label: 'Brittle Bones',
						// 	actionId: ActionId.fromSpellId(81328),
						// 	playerData: playerClassAndTalent(Class.ClassDeathKnight, 'brittleBones'),
						// },
						{
							label: 'Acid Spit',
							actionId: ActionId.fromSpellId(55749),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.Worm),
						},
					],
				},
				{
					label: '+Bleed %',
					effects: [
						// {
						// 	label: 'Blood Frenzy',
						// 	actionId: ActionId.fromSpellId(29859),
						// 	playerData: playerClassAndTalent(Class.ClassWarrior, 'bloodFrenzy'),
						// },
						{
							label: 'Mangle',
							actionId: ActionId.fromSpellId(33878),
							playerData: playerClass(Class.ClassDruid, player => player.isSpec(Spec.SpecFeralDruid)),
						},
						// {
						// 	label: 'Hemorrhage',
						// 	actionId: ActionId.fromSpellId(16511),
						// 	playerData: playerClassAndTalent(Class.ClassRogue, 'hemorrhage'),
						// },
						{
							label: 'Stampede',
							actionId: ActionId.fromSpellId(57386),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.Rhino),
						},
					],
				},
				{
					label: 'Spell Crit',
					effects: [
						// {
						// 	label: 'Critical Mass',
						// 	actionId: ActionId.fromSpellId(12873),
						// 	playerData: playerClassAndTalent(Class.ClassMage, 'criticalMass'),
						// },
						// {
						// 	label: 'Shadow and Flame',
						// 	actionId: ActionId.fromSpellId(17801),
						// 	playerData: playerClassAndTalent(Class.ClassWarlock, 'shadowAndFlame'),
						// },
					],
				},
				{
					label: 'Spell Dmg',
					effects: [
						// {
						// 	label: 'Ebon Plaguebringer',
						// 	actionId: ActionId.fromSpellId(51160),
						// 	playerData: playerClassAndTalent(Class.ClassDeathKnight, 'ebonPlaguebringer'),
						// },
						// {
						// 	label: 'Earth and Moon',
						// 	actionId: ActionId.fromSpellId(60433),
						// 	playerData: playerSpecAndTalent(Spec.SpecBalanceDruid, 'earthAndMoon'),
						// },
						{
							label: 'Curse of Elements',
							actionId: ActionId.fromSpellId(1490),
							playerData: playerClass(Class.ClassWarlock),
						},
						// {
						// 	label: 'Master Poisoner',
						// 	actionId: ActionId.fromSpellId(58410),
						// 	playerData: playerClassAndTalent(Class.ClassRogue, 'masterPoisoner'),
						// },
						{
							label: 'Fire Breath',
							actionId: ActionId.fromSpellId(34889),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.Dragonhawk),
						},
						{
							label: 'Lightning Breath',
							actionId: ActionId.fromSpellId(24844),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.WindSerpent),
						},
					],
				},
			],
		},
		{
			label: 'Mitigation Debuffs',
			categories: [
				{
					label: '-Dmg %',
					effects: [
						// {
						// 	label: 'Vindication',
						// 	actionId: ActionId.fromSpellId(26016),
						// 	playerData: playerClassAndTalent(Class.ClassPaladin, 'vindication', player =>
						// 		[Spec.SpecRetributionPaladin, Spec.SpecProtectionPaladin].includes(player.getSpec()),
						// 	),
						// },
						{
							label: 'Curse of Weakness',
							actionId: ActionId.fromSpellId(702),
							playerData: playerClass(Class.ClassWarlock),
						},
						{
							label: 'Demoralizing Roar',
							actionId: ActionId.fromSpellId(99),
							playerData: playerClass(Class.ClassDruid, player => player.isSpec(Spec.SpecFeralDruid)),
						},
						// {
						// 	label: 'Scarlet Fever',
						// 	actionId: ActionId.fromSpellId(81130),
						// 	playerData: playerClassAndTalent(Class.ClassDeathKnight, 'scarletFever'),
						// },
						{
							label: 'Demoralizing Shout',
							actionId: ActionId.fromSpellId(1160),
							playerData: playerClass(Class.ClassWarrior),
						},
					],
				},
				{
					label: 'Atk Speed',
					effects: [
						{
							label: 'Thunder Clap',
							actionId: ActionId.fromSpellId(6343),
							playerData: playerClass(Class.ClassWarrior),
						},
						{
							label: 'Frost Fever',
							actionId: ActionId.fromSpellId(59921),
							playerData: playerClass(Class.ClassDeathKnight),
						},
						// {
						// 	label: 'Judgements of the Just',
						// 	actionId: ActionId.fromSpellId(53696),
						// 	playerData: playerClassAndTalent(Class.ClassPaladin, 'judgementsOfTheJust'),
						// },
						// {
						// 	label: 'Infected Wounds',
						// 	actionId: ActionId.fromSpellId(48484),
						// 	playerData: playerClassAndTalent(Class.ClassDruid, 'infectedWounds', player =>
						// 		[Spec.SpecFeralDruid, Spec.SpecFeralDruid].includes(player.getSpec()),
						// 	),
						// },
						{
							label: 'Earth Shock',
							actionId: ActionId.fromSpellId(8042),
							playerData: playerClass(Class.ClassShaman),
						},
						{
							label: 'Dust Cloud',
							actionId: ActionId.fromSpellId(50285),
							playerData: playerClass(Class.ClassHunter, player => player.getSpecOptions().classOptions?.petType == HunterPetType.Tallstrider),
						},
					],
				},
			],
		},
	],
};
