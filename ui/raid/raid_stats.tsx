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
					label: 'Major Haste',
					effects: [
						{
							label: 'Bloodlust',
							actionId: ActionId.fromSpellId(2825),
							playerData: playerClass(Class.ClassShaman),
						},
						{
							label: 'Heroism',
							actionId: ActionId.fromSpellId(32182),
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
					label: 'Stats',
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
					],
				},
				{
					label: 'Attack Power',
					effects: [
						{
							label: 'Trueshot Aura',
							actionId: ActionId.fromSpellId(19506),
							playerData: playerClass(Class.ClassHunter),
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
					label: 'Attack Speed',
					effects: [
						{
							label: 'Unholy Aura',
							actionId: ActionId.fromSpellId(55610),
							playerData: playerClass(Class.ClassDeathKnight),
						},
						{
							label: "Serpent's Swiftness",
							actionId: ActionId.fromSpellId(128433),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: "Swiftblade's Cunning",
							actionId: ActionId.fromSpellId(113742),
							playerData: playerClass(Class.ClassRogue),
						},
						{
							label: 'Unleashed Rage',
							actionId: ActionId.fromSpellId(30809),
							playerData: playerClass(Class.ClassShaman),
						},
					],
				},
				{
					label: 'Spell Power',
					effects: [
						{
							label: 'Arcane Brilliance',
							actionId: ActionId.fromSpellId(1459),
							playerData: playerClass(Class.ClassMage),
						},
						{
							label: 'Still Water',
							actionId: ActionId.fromSpellId(126309),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Burning Wrath',
							actionId: ActionId.fromSpellId(77747),
							playerData: playerClass(Class.ClassShaman),
						},
						{
							label: 'Dark Intent',
							actionId: ActionId.fromSpellId(109773),
							playerData: playerClass(Class.ClassWarlock),
						},
					],
				},
				{
					label: 'Spell Haste',
					effects: [
						{
							label: 'Shadow Form',
							actionId: ActionId.fromSpellId(15473),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Moonkin Aura',
							actionId: ActionId.fromSpellId(24907),
							playerData: playerClass(Class.ClassDruid),
						},
						{
							label: 'Mind Quickening',
							actionId: ActionId.fromSpellId(49868),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Elemental Oath',
							actionId: ActionId.fromSpellId(51470),
							playerData: playerClass(Class.ClassShaman),
						},
					],
				},
				{
					label: 'Crit %',
					effects: [
						{
							label: 'Leader of the Pack',
							actionId: ActionId.fromSpellId(17007),
							playerData: playerClass(Class.ClassDruid),
						},
						{
							label: 'Furious Howl',
							actionId: ActionId.fromSpellId(24604),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Terrifying Roar',
							actionId: ActionId.fromSpellId(90309),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Legacy of the White Tiger',
							actionId: ActionId.fromSpellId(116781),
							playerData: playerClass(Class.ClassMonk),
						},
					],
				},
				{
					label: 'Mastery',
					effects: [
						{
							label: 'Blessing of Might',
							actionId: ActionId.fromSpellId(19740),
							playerData: playerClass(Class.ClassPaladin),
						},
						{
							label: 'Roar of Courage',
							actionId: ActionId.fromSpellId(93435),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Spirit Beast Blessing',
							actionId: ActionId.fromSpellId(128997),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Grace of Air',
							actionId: ActionId.fromSpellId(116956),
							playerData: playerClass(Class.ClassShaman),
						},
					],
				},
				{
					label: 'Stamina',
					effects: [
						{
							label: 'Power Word: Fortitude',
							actionId: ActionId.fromSpellId(21562),
							playerData: playerClass(Class.ClassPriest),
						},
						{
							label: 'Qiraji Fortitude',
							actionId: ActionId.fromSpellId(90364),
							playerData: playerClass(Class.ClassHunter),
						},
						{
							label: 'Commanding Shout',
							actionId: ActionId.fromSpellId(469),
							playerData: playerClass(Class.ClassWarrior),
						},
					],
				},
				{
					label: 'Mana Regen',
					effects: [
						{
							label: 'Mana Tide Totem',
							actionId: ActionId.fromSpellId(5675),
							playerData: playerClass(Class.ClassShaman),
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
