import { getWowheadLanguagePrefix } from '../constants/lang';
import { CHARACTER_LEVEL, MAX_CHALLENGE_MODE_ILVL } from '../constants/mechanics';
import { ActionID as ActionIdProto, ItemLevelState, ItemRandomSuffix, OtherAction, ReforgeStat } from '../proto/common';
import { ResourceType } from '../proto/spell';
import { IconData, UIItem as Item } from '../proto/ui';
import { buildWowheadTooltipDataset, WowheadTooltipItemParams, WowheadTooltipSpellParams } from '../wowhead';
import { Database } from './database';

// If true uses wotlkdb.com, else uses wowhead.com.
export const USE_WOTLK_DB = false;

type ActionIdOptions = {
	itemId?: number;
	spellId?: number;
	otherId?: OtherAction;
	tag?: number;
	baseName?: string;
	name?: string;
	iconUrl?: string;
	randomSuffixId?: number;
	reforgeId?: number;
	upgradeStep?: number;
};

// Uniquely identifies a specific item / spell / thing in WoW. This object is immutable.
export class ActionId {
	readonly itemId: number;
	readonly randomSuffixId: number;
	readonly reforgeId: number;
	readonly upgradeStep: ItemLevelState;
	readonly spellId: number;
	readonly otherId: OtherAction;
	readonly tag: number;

	readonly baseName: string; // The name without any tag additions.
	readonly name: string;
	readonly iconUrl: string;
	readonly spellIdTooltipOverride: number | null;

	private constructor({ itemId, spellId, otherId, tag, baseName, name, iconUrl, randomSuffixId, reforgeId, upgradeStep }: ActionIdOptions = {}) {
		this.itemId = itemId ?? 0;
		this.randomSuffixId = randomSuffixId ?? 0;
		this.reforgeId = reforgeId ?? 0;
		this.upgradeStep = upgradeStep ?? 0;
		this.spellId = spellId ?? 0;
		this.otherId = otherId ?? OtherAction.OtherActionNone;
		this.tag = tag ?? 0;

		switch (otherId) {
			case OtherAction.OtherActionNone:
				break;
			case OtherAction.OtherActionWait:
				baseName = 'Wait';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_pocketwatch_01.jpg';
				break;
			case OtherAction.OtherActionManaRegen:
				name = 'Mana Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				if (tag == 1) {
					name += ' (In Combat)';
				} else if (tag == 2) {
					name += ' (Out of Combat)';
				}
				break;
			case OtherAction.OtherActionEnergyRegen:
				baseName = 'Energy Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeEnergy];
				break;
			case OtherAction.OtherActionFocusRegen:
				baseName = 'Focus Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeFocus];
				break;
			case OtherAction.OtherActionManaGain:
				baseName = 'Mana Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				break;
			case OtherAction.OtherActionRageGain:
				baseName = 'Rage Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeRage];
				break;
			case OtherAction.OtherActionAttack:
				name = 'Attack';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_04.jpg';
				if (this.tag == 1) {
					name += ' (Main Hand)';
				} else if (this.tag == 2) {
					name += ' (Off Hand)';
				} else if (this.tag == 41570) {
					name += ' (Magmaw)';
				} else if (this.tag == 49416) {
					name += ' (Blazing Bone Construct)';
				} else if (this.tag == 56427) {
					name += ' (Warmaster Blackhorn)';
				} else if (this.tag == 56781) {
					name += ' (Goriona)';
				} else if (this.tag > 4191800) {
					name += ` (Animated Bone Warrior ${(this.tag - 4191800).toFixed(0)})`;
				}
				break;
			case OtherAction.OtherActionShoot:
				name = 'Shoot';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/ability_marksmanship.jpg';
				break;
			case OtherAction.OtherActionPet:
				break;
			case OtherAction.OtherActionRefund:
				baseName = 'Refund';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_misc_coin_01.jpg';
				break;
			case OtherAction.OtherActionDamageTaken:
				baseName = 'Damage Taken';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_sword_04.jpg';
				break;
			case OtherAction.OtherActionHealingModel:
				baseName = 'Incoming HPS';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/spell_holy_renew.jpg';
				break;
			case OtherAction.OtherActionBloodRuneGain:
				baseName = 'Blood Rune Gain';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_deathstrike.jpg';
				break;
			case OtherAction.OtherActionFrostRuneGain:
				baseName = 'Frost Rune Gain';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_deathstrike2.jpg';
				break;
			case OtherAction.OtherActionUnholyRuneGain:
				baseName = 'Unholy Rune Gain';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_empowerruneblade.jpg';
				break;
			case OtherAction.OtherActionDeathRuneGain:
				baseName = 'Death Rune Gain';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_empowerruneblade.jpg';
				break;
			case OtherAction.OtherActionPotion:
				baseName = 'Potion';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/large/inv_alchemy_elixir_04.jpg';
				break;
			case OtherAction.OtherActionMove:
				baseName = 'Moving';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/inv_boots_cloth_03.jpg';
				break;
			case OtherAction.OtherActionPrepull:
				baseName = 'Prepull';
				iconUrl = 'https://wow.zamimg.com/images/wow/icons/medium/inv_misc_pocketwatch_02.jpg';
				break;
		}
		this.baseName = baseName ?? '';
		this.name = (name || baseName) ?? '';
		this.iconUrl = iconUrl ?? '';
		this.spellIdTooltipOverride = this.spellTooltipOverride?.spellId || null;
	}

	anyId(): number {
		return this.itemId || this.spellId || this.otherId;
	}

	equals(other: ActionId): boolean {
		return this.equalsIgnoringTag(other) && this.tag == other.tag;
	}

	equalsIgnoringTag(other: ActionId): boolean {
		return (
			this.itemId == other.itemId &&
			this.randomSuffixId == other.randomSuffixId &&
			this.spellId == other.spellId &&
			this.otherId == other.otherId &&
			this.upgradeStep === other.upgradeStep
		);
	}

	setBackground(elem: HTMLElement) {
		if (this.iconUrl) {
			elem.style.backgroundImage = `url('${this.iconUrl}')`;
		}
	}

	static makeItemUrl(id: number, randomSuffixId?: number, reforgeId?: number, upgradeStep?: ItemLevelState): string {
		const langPrefix = getWowheadLanguagePrefix();
		const url = new URL(`https://wowhead.com/mop-classic/${langPrefix}item=${id}`);
		url.searchParams.set('level', String(CHARACTER_LEVEL));
		url.searchParams.set('rand', String(randomSuffixId || 0));
		if (reforgeId) url.searchParams.set('forg', String(reforgeId));
		if (upgradeStep === ItemLevelState.ChallengeMode) {
			url.searchParams.set('ilvl', String(MAX_CHALLENGE_MODE_ILVL));
		} else if (upgradeStep) {
			url.searchParams.set('upgd', String(upgradeStep));
		}
		return url.toString();
	}
	static makeSpellUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return `https://wotlkdb.com/?spell=${id}`;
		} else {
			return `https://wowhead.com/mop-classic/${langPrefix}spell=${id}`;
		}
	}
	static async makeItemTooltipData(id: number, params?: Omit<WowheadTooltipItemParams, 'itemId'>) {
		return buildWowheadTooltipDataset({ itemId: id, ...params });
	}
	static async makeSpellTooltipData(id: number, params?: Omit<WowheadTooltipSpellParams, 'spellId'>) {
		return buildWowheadTooltipDataset({ spellId: id, ...params });
	}
	static makeQuestUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return 'https://wotlkdb.com/?quest=' + id;
		} else {
			return `https://wowhead.com/mop-classic/${langPrefix}quest=${id}`;
		}
	}
	static makeNpcUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return 'https://wotlkdb.com/?npc=' + id;
		} else {
			return `https://wowhead.com/mop-classic/${langPrefix}npc=${id}`;
		}
	}
	static makeZoneUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return 'https://wotlkdb.com/?zone=' + id;
		} else {
			return `https://wowhead.com/mop-classic/${langPrefix}zone=${id}`;
		}
	}

	setWowheadHref(elem: HTMLAnchorElement) {
		if (this.itemId) {
			elem.href = ActionId.makeItemUrl(this.itemId, this.randomSuffixId, this.reforgeId, this.upgradeStep);
		} else if (this.spellId) {
			elem.href = ActionId.makeSpellUrl(this.spellIdTooltipOverride || this.spellId);
		}
	}

	async setWowheadDataset(elem: HTMLElement, params?: Omit<WowheadTooltipItemParams, 'itemId'> | Omit<WowheadTooltipSpellParams, 'spellId'>) {
		(this.itemId
			? ActionId.makeItemTooltipData(this.itemId, params)
			: ActionId.makeSpellTooltipData(this.spellIdTooltipOverride || this.spellId, params)
		).then(url => {
			if (elem) elem.dataset.wowhead = url;
		});
	}

	setBackgroundAndHref(elem: HTMLAnchorElement) {
		this.setBackground(elem);
		this.setWowheadHref(elem);
	}

	async fillAndSet(elem: HTMLAnchorElement, setHref: boolean, setBackground: boolean, options: { signal?: AbortSignal } = {}): Promise<ActionId> {
		const filled = await this.fill(undefined, options);
		if (setHref) {
			filled.setWowheadHref(elem);
		}
		if (setBackground) {
			filled.setBackground(elem);
		}
		return filled;
	}

	// Returns an ActionId with the name and iconUrl fields filled.
	// playerIndex is the optional index of the player to whom this ID corresponds.
	async fill(playerIndex?: number, options: { signal?: AbortSignal } = {}): Promise<ActionId> {
		if (this.name || this.iconUrl) {
			return this;
		}

		if (this.otherId) {
			return this;
		}
		const tooltipData = await ActionId.getTooltipData(this, { signal: options?.signal });

		const baseName = tooltipData['name'];
		let name = baseName;

		let tag = this.tag;
		// Handle Monk's Storm, Earth and Fire
		if (tag >= 138228) {
			tag -= 138228;
		}
		// Handle DTR
		else if (tag >= 71086 && tag <= 71096) {
			name = 'Dragonwrath - ' + name;
			tag -= 71086;
		}

		switch (baseName) {
			case 'Minor Speed':
				name = 'Minor Run Speed (8%)';
				break;
			case 'Arcane Missiles':
				break;
			case 'Arcane Blast':
				if (tag == 1) {
					name += ' (No Stacks)';
				} else if (tag == 2) {
					name += ` (1 Stack)`;
				} else if (tag > 2) {
					name += ` (${tag - 1} Stacks)`;
				}
				break;
			case 'Hot Streak':
				if (tag) name += ' (Crits)';
				break;
			case 'Fireball':
			case 'Flamestrike':
				if (tag == 1) name += ' (Blast Wave)';
				break;
			case 'Pyroblast':
			case 'Combustion':
				if (tag) name += ' (DoT)';
				break;
			case 'Frost Bomb':
			case 'Living Bomb':
				if (tag == 1) name += ' (DoT)';
				else if (tag == 2) name += ' (Explosion)';
				break;
			case 'Evocation':
				if (tag == 1) {
					name += ' (1 Tick)';
				} else if (tag == 2) {
					name += ' (2 Tick)';
				} else if (tag == 3) {
					name += ' (3 Tick)';
				} else if (tag == 4) {
					name += ' (4 Tick)';
				} else if (tag == 5) {
					name += ' (5 Tick)';
				}
				break;
			case 'Mind Flay':
				if (this.spellId === 15407) {
					if (tag == 1) {
						name += ' (1 Tick)';
					} else if (tag == 2) {
						name += ' (2 Tick)';
					} else if (tag == 3) {
						name += ' (3 Tick)';
					} else if (tag == 77486) {
						name += ' (Mastery)';
					}
				} else {
					// Gurthalak, Voice of the Deeps
					if (tag === 0) {
						name += ' (LFR)';
					} else if (tag === 1) {
						name += ' (Normal)';
					} else if (tag === 2) {
						name += ' (Heroic)';
					}
				}
				break;
			case 'Mind Sear':
				if (tag == 1) {
					name += ' (1 Tick)';
				} else if (tag == 2) {
					name += ' (2 Tick)';
				} else if (tag == 3) {
					name += ' (3 Tick)';
				} else if (tag == 77486) {
					name += ' (Mastery)';
				}

				break;
			case 'Devotion Aura':
			case 'Shattering Throw':
			case 'Skull Banner':
			case 'Stormlash':
				if (tag === -1) {
					name += ' (raid)';
				} else {
					name += ` (self)`;
				}
				break;
			case 'Envenom':
			case 'Eviscerate':
			case 'Rupture':
			case 'Slice and Dice':
				if (tag) name += ` (${tag} CP)`;
				break;
			case 'Crimson Tempest':
				if (tag == 7) {
					name += ' (DoT)';
				} else if (tag) {
					name += ` (${tag} CP)`;
				}
				break;
			case 'Deadly Poison':
				if (tag == 1) {
					name += ' (DoT)';
				} else {
					name += ' (Hit)';
				}
				break;
			case 'Shadow Blades':
				if (tag == 1) {
					name = 'Shadow Blade';
				} else if (tag == 2) {
					name = 'Shadow Blade Off-hand';
				}
				break;
			case 'Bladestorm':
			case 'Killing Spree':
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Tricks of the Trade':
				if (tag == 1) {
					name += ' (Not Self)';
				}
				break;
			case 'Mutilate':
				if (tag == 0) {
					name += ' (Cast)';
				} else if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Hemorrhage':
				if (tag == 1) {
					name += ' (Hit)';
				} else {
					name += ' (DoT)';
				}
				break;
			case 'Wind Lash':
			case 'Stormstrike':
			case 'Stormblast':
				if (tag == 0) {
					name += ' (Cast)';
				} else if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Elemental Blast':
			case 'Chain Lightning':
			case 'Lightning Bolt':
			case 'Lava Beam':
			case 'Lava Burst':
				if (tag == 6) {
					name += ' (Overload)';
				} else if (tag == 7) {
					name += ' (Echo)';
				} else if (tag == 8) {
					name += ' (Echo Overload)';
				} else if (tag == 9) {
					name += ' (Mastery)';
				} else if (tag == 10) {
					name += ' (Haste)';
				} else if (tag == 11) {
					name += ' (Crit)';
				} else if (tag == 12) {
					name += ' (Agi)';
				} else if (tag) {
					name += ` (${tag} MW)`;
				}
				break;
			case 'Flame Shock':
			case 'Earth Shock':
			case 'Frost Shock':
			case 'Fire Nova':
				if (tag == 1) {
					name += ' (DoT)';
				} else if (tag == 7) {
					name += ' (Echo)';
				}
				break;
			case 'Fulmination':
				name += ` (${tag + 2})`;
				break;
			case 'Lightning Shield':
				if (tag == 1) {
					name += ' (Wasted)';
				}
				break;
			case 'Crescendo of Suffering':
				if (tag == 1) {
					name += ' (Pre-Pull)';
				}
				break;
			case 'Shadowflame':
			case 'Moonfire':
			case 'Sunfire':
				if (tag == 1) {
					name += ' (DoT)';
				}
				break;
			case 'Sacred Shield':
				if (this.spellId === 65148) {
					name += ' (Absorb)';
				}
				break;
			case 'Censure':
				if (tag == 2) {
					name += ' (DoT)';
				}
				break;
			case 'Hammer of the Righteous':
				if (this.spellId === 88263) {
					name += ' (Holy)';
				} else {
					name += ' (Physical)';
				}
				break;
			case 'Seal of Righteousness':
				if (tag === 1) {
					name += ' (DS)';
				}
				break;
			case 'Harsh Word':
				name = 'Word of Glory (Damage';
				if (tag > 0) {
					name += `, ${tag} HP`;
				}
				name += ')';
				break;
			case 'Word of Glory':
				if (tag > 0) {
					name += `(${tag} HP)`;
				}
				break;
			case 'Eternal Flame':
				if (tag > 0) {
					name += `(${tag} HP)`;
				}
				break;
			case 'Exorcism':
				if (tag === 2) {
					name = 'Glyph of Mass Exorcism';
				}
				break;
			case "Templar's Verdict":
				if (tag === 2) {
					name += ' (T15 4P)';
				}
				break;
			// For targetted buffs, tag is the source player's raid index or -1 if none.
			case 'Bloodlust':
			case 'Ferocious Inspiration':
			case 'Innervate':
			case 'Focus Magic':
			case 'Mana Tide Totem':
			case 'Stormlash Totem':
			case 'Unholy Frenzy':
			case 'Power Infusion':
				if (tag != -1) {
					if (tag === playerIndex || playerIndex == undefined) {
						name += ` (self)`;
					} else {
						name += ` (from #${tag + 1})`;
					}
				} else {
					name += ' (raid)';
				}
				break;
			case 'Elemental Mastery':
				if (this.spellId === 64701) {
					name = `${name} (Buff)`;
				} else {
					name = `${name} (Instant)`;
				}
				break;
			case 'Heart Strike':
				if (tag == 2) {
					name += ' (Off-target)';
				}
				break;
			case 'Raging Blow':
				if (tag == 2) {
					name += ' (Main Hand)';
				} else if (tag == 3) {
					name += ' (Off Hand)';
				}
				break;
			case 'Soul Reaper':
				if (tag === 0) {
					name += ' (Tick)';
				}
				break;
			case 'Whirlwind':
			case 'Storm Bolt':
			case 'Frost Strike':
			case 'Plague Strike':
			case 'Blood Strike':
			case 'Obliterate':
			case 'Rune of Razorice':
			case 'Blood-Caked Strike':
			case 'Festering Strike':
			case 'Razor Frost':
			case 'Lightning Speed':
			case 'Windfury Weapon':
			case 'Berserk':
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				// Death Knight - T12 4P proc
				if (baseName === 'Obliterate' && tag === 3) {
					name = 'Flaming Torment (T12 4P)';
				}
				break;
			case 'Death Strike':
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				} else if (this.spellId === 45470) {
					name += ' (Heal)';
				}
				break;
			case 'Battle Shout':
				if (tag == 1) {
					name += ' (Snapshot)';
				}
				break;
			case 'Heroic Strike':
			case 'Cleave':
			case 'Maul':
				if (tag == 1) {
					name += ' (Queue)';
				}
				break;
			case 'Seed of Corruption':
				if (tag == 0) {
					name += ' (DoT)';
				} else if (tag == 1) {
					name += ' (Explosion)';
				}
				break;
			case 'Thunderfury':
				if (tag == 1) {
					name += ' (ST)';
				} else if (tag == 2) {
					name += ' (MT)';
				}
				break;
			case 'Devouring Plague':
				if (tag == 1) {
					name += ' (DoT)';
					break;
				}
				if (tag == 77486) {
					name += ' (Mastery)';
					break;
				}
			case 'Shadow Word: Death':
				if (tag == 1) {
					name += ' (No Orb)';
				}
			case 'Steady Focus':
				if (tag == 2) {
					name += ' (pre)';
				}
				break;
			case 'Immolate':
				if (tag == 1) {
					name += ' (DoT)';
				}
				break;
			case 'Immolation Aura':
				if (tag == 2) {
					name += ' (Tick)';
				}
				break;
			case 'Frozen Blows':
			case 'Opportunity Strike':
				break;
			// Warrior - T12 2P proc
			case 'Shield Slam':
				if (tag === 3) {
					name = 'Combust (T12 2P)';
				}
				break;
			// Warrior - T13 4P proc
			case 'Colossus Smash':
				if (this.spellId === 108126) {
					name += ' (T13 4P)';
				}
				break;
			// Warrior - T12 4P proc
			case 'Mortal Strike':
				if (tag === 3) {
					name = 'Fiery attack (T12 4P)';
				}
				break;
			case 'Slam':
				if (tag == 1) {
					name += ' (Sweeping Strikes)';
				}
				break;
			// Hunter - T12 2P proc
			case 'Steady Shot':
			case 'Cobra Shot':
				if (tag === 3) {
					name = 'Flaming Arrow (T12 2P)';
				}
				break;
			// Paladin - T12 4P proc
			case 'Shield of the Righteous':
				if (tag === 3) {
					name = 'Righteous Flames (T12 2P)';
				}
				break;
			// Paladin - T12 4P proc
			case 'Crusader Strike':
				if (tag === 3) {
					name = 'Flames of the Faithful (T12 2P)';
				}
				break;
			// Death Knight - T12 4P proc
			case 'Scourge Strike':
				if (tag === 3) {
					name = 'Flaming Torment (T12 4P)';
				}
				break;
			// Death Knight - T12 2P proc
			case 'Burning Blood':
				name += ' (T12 2P)';
				break;
			// Death Knight - T12 4P proc
			case 'Flaming Rune Weapon':
				name += ' (T12 4P)';
				break;
			// Souldrinker - Drain Life
			case 'Drain Life':
				if (this.spellId === 109828) {
					name += ' 1.3%';
				} else if (this.spellId === 108022) {
					name += ' 1.5%';
				} else if (this.spellId === 109831) {
					name += ' 1.7%';
				}

				if (tag === 2) {
					name += ' (Heal)';
				}
				break;
			// No'Kaled
			case 'Flameblast':
			case 'Iceblast':
			case 'Shadowblast':
				if (this.spellId === 109871 || this.spellId === 109869 || this.spellId === 109867) {
					name += ' (LFR)';
				} else if (this.spellId === 107785 || this.spellId === 107789 || this.spellId === 107787) {
					name += ' (Normal)';
				} else if (this.spellId === 109872 || this.spellId === 109870 || this.spellId === 109868) {
					name += ' (Heroic)';
				}
				break;
			case 'Death Coil':
				if (tag === 2) {
					name += ' (Heal)';
				}
				break;

			case 'Item - Paladin T11 Retribution 4P Bonus':
				name = 'Reinforced Sapphirium Battleplate - T11 4pc';
				break;
			case 'Item - Paladin T12 Retribution 4P Bonus':
				name = 'Battleplate of Immolation - T12 4pc';
				break;
			case 'Item - Paladin T14 Retribution 2P Bonus':
				name = 'White Tiger Battlegear - T14 2pc';
				break;
			case 'Item - Paladin T10 Retribution 4P Bonus':
				name = 'White Tiger Battlegear - T14 4pc';
				break;
			case 'Item - Paladin T14 Protection 2P Bonus':
				name = 'White Tiger Plate - T14 2pc';
				break;
			case 'Item - Paladin T14 Protection 4P Bonus':
				name = 'White Tiger Plate - T14 4pc';
				break;
			case 'Item - Paladin T15 Retribution 2P Bonus':
				name = 'Battlegear of the Lightning Emperor - T15 2pc';
				break;
			case 'Item - Paladin T15 Retribution 4P Bonus':
				name = 'Battlegear of the Lightning Emperor - T15 4pc';
				break;
			case 'Item - Paladin T15 Protection 2P Bonus':
				name = 'Plate of the Lightning Emperor - T15 2pc';
				break;
			case 'Item - Paladin T15 Protection 4P Bonus':
				name = 'Plate of the Lightning Emperor - T15 4pc';
				break;
			case 'Item - Paladin T16 Retribution 2P Bonus':
				name = 'Battlegear of Winged Triumph - T16 2pc';
				break;
			case 'Item - Paladin T16 Retribution 4P Bonus':
				name = 'Battlegear of Winged Triumph - T16 4pc';
				break;
			case 'Item - Paladin T16 Protection 2P Bonus':
				name = 'Plate of Winged Triumph - T16 2pc';
				break;
			case 'Item - Paladin T16 Protection 4P Bonus':
				name = 'Plate of Winged Triumph - T16 4pc';
				break;
			case 'Hurricane':
			case 'Dancing Steel':
			case 'Bloody Dancing Steel':
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				} else if (tag == 3) {
					name += ' (Spell)';
				}
				break;
			case 'Landslide':
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Jade Spirit':
				if (tag == 1) {
					name += ' (Intellect)';
				} else if (tag == 2) {
					name += ' (Spirit)';
				}
				break;

			// Monk
			case 'Tiger Strikes':
				if (this.spellId === 120274) {
					name += ' (Main Hand)';
				} else {
					name += ' (Off Hand)';
				}
				break;
			case 'Blackout Kick':
				if (tag === 2) {
					name += ' (DoT)';
				}
				break;
			case 'Expel Harm':
				if (this.spellId === 115072) {
					name += ' (Heal)';
				} else {
					name += ' (Damage)';
				}
				break;
			case 'Chi Wave':
				if (this.spellId === 132463) {
					name += ' (Heal)';
				} else if (this.spellId === 132467) {
					name += ' (Damage)';
				}
				break;
			case 'Zen Sphere':
			case 'Zen Sphere: Detonate':
				if (this.spellId === 124081) {
					if ([4, 5].includes(this.tag)) {
						name += ': Detonate';
					}
					if (this.tag === 1) {
						name += ' (# of Spheres)';
					} else if ([0, 2, 4].includes(this.tag)) {
						name += ' (Heal)';
					} else if ([3, 5].includes(this.tag)) {
						name += ' (Damage)';
					}
				}
				break;
			case 'Chi Burst':
				if (this.spellId === 130654) {
					name += ' (Heal)';
				} else if (this.spellId === 148135) {
					name += ' (Damage)';
				}
				break;
			case 'Stagger':
				if (this.tag === 1) {
					name += ' (Dot)';
				}
				break;
			case 'Dampen Harm':
				break;
			case 'Healing Sphere':
				if (this.spellId === 115460) {
					if (tag === 1) {
						name += ' (Stacks)';
					} else {
						name += ' (Cast)';
					}
				}
				if (this.spellId === 115464) {
					name += ' (Heal)';
				}
				break;
			case 'Vampiric Touch':
			case 'Shadow Word: Pain':
				if (tag == 77486) {
					name += ' (Mastery)';
				}

				break;
			case 'Cascade':
				if (tag == 1) {
					name += ' (Bounce)';
				}

				break;
			case 'Agony':
			case 'Unstable Affliction':
			case 'Corruption':
				if (tag == 1) {
					name += ' (Malefic)';
				}
				break;
			case 'Holy Prism':
				if (this.spellId === 114852) {
					if (tag === 1) {
						name += ' (Damage)';
					} else if (tag === 2) {
						name += ' (Aoe heal)';
					}
				} else if (this.spellId === 114871) {
					if (tag === 1) {
						name += ' (Heal)';
					} else if (tag === 2) {
						name += ' (Aoe damage)';
					}
				}
				break;
			case 'Alter Time':
				if (tag == 1) {
					name += " (Restore)"
				}
				break;
			default:
				if (tag) {
					name += ' (??)';
				}
				break;
		}

		let iconUrl = ActionId.makeIconUrl(tooltipData['icon']);

		const iconOverrideId = this.spellTooltipOverride || this.spellIconOverride;
		if (iconOverrideId) {
			const overrideTooltipData = await ActionId.getTooltipData(iconOverrideId, { signal: options?.signal });
			iconUrl = ActionId.makeIconUrl(overrideTooltipData['icon']);
		}

		return new ActionId({
			itemId: this.itemId,
			spellId: this.spellId,
			otherId: this.otherId,
			tag: this.tag,
			baseName,
			name,
			iconUrl,
			randomSuffixId: this.randomSuffixId,
			reforgeId: this.reforgeId,
			upgradeStep: this.upgradeStep,
		});
	}

	toString(): string {
		return this.toStringIgnoringTag() + (this.tag ? '-' + this.tag : '');
	}

	toStringIgnoringTag(): string {
		if (this.itemId) {
			return 'item-' + this.itemId;
		} else if (this.spellId) {
			return 'spell-' + this.spellId;
		} else if (this.otherId) {
			return 'other-' + this.otherId;
		} else {
			console.error('Empty action id!');
			return '';
		}
	}

	toProto(): ActionIdProto {
		const protoId = ActionIdProto.create({
			tag: this.tag,
		});

		if (this.itemId) {
			protoId.rawId = {
				oneofKind: 'itemId',
				itemId: this.itemId,
			};
		} else if (this.spellId) {
			protoId.rawId = {
				oneofKind: 'spellId',
				spellId: this.spellId,
			};
		} else if (this.otherId) {
			protoId.rawId = {
				oneofKind: 'otherId',
				otherId: this.otherId,
			};
		}

		return protoId;
	}

	toProtoString(): string {
		return ActionIdProto.toJsonString(this.toProto());
	}

	withoutTag(): ActionId {
		return new ActionId({
			itemId: this.itemId,
			spellId: this.spellId,
			otherId: this.otherId,
			baseName: this.baseName,
			iconUrl: this.iconUrl,
			randomSuffixId: this.randomSuffixId,
			reforgeId: this.reforgeId,
			upgradeStep: this.upgradeStep,
		});
	}

	static fromEmpty(): ActionId {
		return new ActionId();
	}

	static fromItemId(itemId: number, tag?: number, randomSuffixId?: number, reforgeId?: number, upgradeStep?: ItemLevelState): ActionId {
		return new ActionId({
			itemId,
			tag,
			randomSuffixId,
			reforgeId,
			upgradeStep,
		});
	}

	static fromSpellId(spellId: number, tag?: number): ActionId {
		return new ActionId({ spellId, tag });
	}

	static fromOtherId(otherId: OtherAction, tag?: number): ActionId {
		return new ActionId({ otherId, tag });
	}

	static fromPetName(petName: string): ActionId {
		return (
			petNameToActionId[petName] ||
			new ActionId({
				baseName: petName,
				iconUrl: petNameToIcon[petName],
			})
		);
	}

	static fromItem(item: Item): ActionId {
		return ActionId.fromItemId(item.id);
	}

	static fromRandomSuffix(item: Item, randomSuffix: ItemRandomSuffix): ActionId {
		return ActionId.fromItemId(item.id, 0, randomSuffix.id);
	}

	static fromReforge(item: Item, reforge: ReforgeStat): ActionId {
		return ActionId.fromItemId(item.id, 0, 0, reforge.id);
	}

	static fromProto(protoId: ActionIdProto): ActionId {
		if (protoId.rawId.oneofKind == 'spellId') {
			return ActionId.fromSpellId(protoId.rawId.spellId, protoId.tag);
		} else if (protoId.rawId.oneofKind == 'itemId') {
			return ActionId.fromItemId(protoId.rawId.itemId, protoId.tag);
		} else if (protoId.rawId.oneofKind == 'otherId') {
			return ActionId.fromOtherId(protoId.rawId.otherId, protoId.tag);
		} else {
			return ActionId.fromEmpty();
		}
	}

	private static readonly logRegex = /{((SpellID)|(ItemID)|(OtherID)): (\d+)(, Tag: (-?\d+))?}/;
	private static readonly logRegexGlobal = new RegExp(ActionId.logRegex, 'g');
	private static fromMatch(match: RegExpMatchArray): ActionId {
		const idType = match[1];
		const id = parseInt(match[5]);
		return new ActionId({
			itemId: idType == 'ItemID' ? id : undefined,
			spellId: idType == 'SpellID' ? id : undefined,
			otherId: idType == 'OtherID' ? id : undefined,
			tag: match[7] ? parseInt(match[7]) : undefined,
		});
	}
	static fromLogString(str: string): ActionId {
		const match = str.match(ActionId.logRegex);
		if (match) {
			return ActionId.fromMatch(match);
		} else {
			console.warn('Failed to parse action id from log: ' + str);
			return ActionId.fromEmpty();
		}
	}

	static async replaceAllInString(str: string): Promise<string> {
		const matches = [...str.matchAll(ActionId.logRegexGlobal)];

		const replaceData = await Promise.all(
			matches.map(async match => {
				const actionId = ActionId.fromMatch(match);
				const filledId = await actionId.fill();
				return {
					firstIndex: match.index || 0,
					len: match[0].length,
					actionId: filledId,
				};
			}),
		);

		// Loop in reverse order so we can greedily apply the string replacements.
		for (let i = replaceData.length - 1; i >= 0; i--) {
			const data = replaceData[i];
			str = str.substring(0, data.firstIndex) + data.actionId.name + str.substring(data.firstIndex + data.len);
		}

		return str;
	}

	private static makeIconUrl(iconLabel: string): string {
		if (USE_WOTLK_DB) {
			return `https://wotlkdb.com/static/images/wow/icons/large/${iconLabel}.jpg`;
		} else {
			return `https://wow.zamimg.com/images/wow/icons/large/${iconLabel}.jpg`;
		}
	}

	static async getTooltipData(actionId: ActionId, options: { signal?: AbortSignal } = {}): Promise<IconData> {
		if (actionId.itemId) {
			return await Database.getItemIconData(actionId.itemId, { signal: options?.signal });
		} else {
			return await Database.getSpellIconData(actionId.spellId, { signal: options?.signal });
		}
	}

	get spellIconOverride(): ActionId | null {
		const override = spellIdIconOverrides.get(JSON.stringify({ spellId: this.spellId }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromSpellId(override.spellId!);
	}

	get spellTooltipOverride(): ActionId | null {
		const override = spellIdTooltipOverrides.get(JSON.stringify({ spellId: this.spellId, tag: this.tag }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromSpellId(override.spellId!);
	}
}

type ActionIdOverride = { itemId?: number; spellId?: number };

// Some items/spells have weird icons, so use this to show a different icon instead.
const spellIdIconOverrides: Map<string, ActionIdOverride> = new Map([
	[JSON.stringify({ spellId: 37212 }), { itemId: 29035 }], // Improved Wrath of Air Totem
	[JSON.stringify({ spellId: 37223 }), { itemId: 29040 }], // Improved Strength of Earth Totem
	[JSON.stringify({ spellId: 37447 }), { itemId: 30720 }], // Serpent-Coil Braid
	[JSON.stringify({ spellId: 37443 }), { itemId: 30196 }], // Robes of Tirisfal (4pc bonus)
	[JSON.stringify({ spellId: 90299 }), { itemId: 65214 }], // Reinforced Sapphirium Battleplate (4pc bonus)
	[JSON.stringify({ spellId: 99116 }), { itemId: 71512 }], // Battleplate of Immolation (4pc bonus)
	[JSON.stringify({ spellId: 105767 }), { itemId: 78727 }], // Battleplate of Radiant Glory (2pc bonus)
	[JSON.stringify({ spellId: 123180 }), { itemId: 85343 }], // White Tiger Battlegear (2pc bonus)
	[JSON.stringify({ spellId: 70762 }), { itemId: 85343 }], // White Tiger Battlegear  (4pc bonus)
	[JSON.stringify({ spellId: 123104 }), { itemId: 95290 }], // White Tiger Plate (2pc bonus)
	[JSON.stringify({ spellId: 123107 }), { itemId: 95290 }], // White Tiger Plate (4pc bonus)
	[JSON.stringify({ spellId: 138159 }), { itemId: 95280 }], // Battlegear of the Lightning Emperor (2pc bonus)
	[JSON.stringify({ spellId: 138164 }), { itemId: 95280 }], // Battlegear of the Lightning Emperor (4pc bonus)
	[JSON.stringify({ spellId: 138238 }), { itemId: 95290 }], // Plate of the Lightning Emperor (2pc bonus)
	[JSON.stringify({ spellId: 138244 }), { itemId: 95290 }], // Plate of the Lightning Emperor (4pc bonus)
	[JSON.stringify({ spellId: 144586 }), { itemId: 99136 }], // Battlegear of Winged Triumph (2pc bonus)
	[JSON.stringify({ spellId: 144593 }), { itemId: 99136 }], // Battlegear of Winged Triumph (4pc bonus)
	[JSON.stringify({ spellId: 144580 }), { itemId: 99126 }], // Plate of Winged Triumph (2pc bonus)
	[JSON.stringify({ spellId: 144566 }), { itemId: 99126 }], // Plate of Winged Triumph (4pc bonus)
	[JSON.stringify({ spellId: 13889 }), { spellId: 109709 }], // Minor Run Speed
	[JSON.stringify({ spellId: 65658 }), { spellId: 48721 }], // Blood Boil RP regen
]);

const spellIdTooltipOverrides: Map<string, ActionIdOverride> = new Map([
	[JSON.stringify({ spellId: 47897, tag: 1 }), { spellId: 47960 }], // Shadowflame Dot
	[JSON.stringify({ spellId: 55090, tag: 2 }), { spellId: 70890 }], // Death Knight - Scourge Strike (Shadow)
	[JSON.stringify({ spellId: 114867, tag: 1 }), { spellId: 114866 }], // Death Knight - Soul Reaper (Blood)
	[JSON.stringify({ spellId: 114867, tag: 2 }), { spellId: 130735 }], // Death Knight - Soul Reaper (Frost)
	[JSON.stringify({ spellId: 114867, tag: 3 }), { spellId: 130736 }], // Death Knight - Soul Reaper (Unholy)
	[JSON.stringify({ spellId: 12294, tag: 3 }), { spellId: 99237 }], // Warrior - T12 4P Fiery Attack - Mortal Strike
	[JSON.stringify({ spellId: 23922, tag: 3 }), { spellId: 99240 }], // Warrior - T12 2P Combust - Shield Slam
	[JSON.stringify({ spellId: 77767, tag: 3 }), { spellId: 99058 }], // Hunter - T12 2P Flaming Arrow - Cobra shot
	[JSON.stringify({ spellId: 56641, tag: 3 }), { spellId: 99058 }], // Hunter - T12 2P Flaming Arrow - Steady shot
	[JSON.stringify({ spellId: 35395, tag: 3 }), { spellId: 99092 }], // Paladin - T12 2P Ret Flames of the Faithful
	[JSON.stringify({ spellId: 53600, tag: 3 }), { spellId: 99075 }], // Paladin - T12 2P Prot Righteous Flames
	[JSON.stringify({ spellId: 85256, tag: 2 }), { spellId: 138165 }], // Paladin - T15 4P Ret Templar's Verdict
	[JSON.stringify({ spellId: 879, tag: 2 }), { spellId: 122032 }], // Paladin - Glyph of Mass Exorcism
	[JSON.stringify({ spellId: 49020, tag: 3 }), { spellId: 99000 }], // Death Knight - T12 4P Flaming Torment
	[JSON.stringify({ spellId: 55090, tag: 3 }), { spellId: 99000 }], // Death Knight - T12 4P Flaming Torment
	[JSON.stringify({ spellId: 2818, tag: 2 }), { spellId: 113780 }], // Rogue - Deadly Poison - Hit
	[JSON.stringify({ spellId: 121411, tag: 7 }), { spellId: 122233 }], // Rogue - Crimson Tempest - DoT
	[JSON.stringify({ spellId: 121471, tag: 1 }), { spellId: 121473 }], // Rogue - Shadow Blade

	// Off-Hand attacks
	[JSON.stringify({ spellId: 49998, tag: 2 }), { spellId: 66188 }], // Death Knight - Death Strike Off-Hand
	[JSON.stringify({ spellId: 49143, tag: 2 }), { spellId: 66196 }], // Death Knight - Frost Strike Off-Hand
	[JSON.stringify({ spellId: 49020, tag: 2 }), { spellId: 66198 }], // Death Knight - Obliterate Off-Hand
	[JSON.stringify({ spellId: 45462, tag: 2 }), { spellId: 66216 }], // Death Knight - Plague Strike Off-Hand
	[JSON.stringify({ spellId: 1329, tag: 2 }), { spellId: 27576 }], // Rogue - Mutilate Off-Hand
	[JSON.stringify({ spellId: 121471, tag: 2 }), { spellId: 121474 }], // Rogue - Shadow Blade Off-Hand
	[JSON.stringify({ spellId: 17364, tag: 2 }), { spellId: 32176 }], // Shaman - Stormstrike Off-Hand
	[JSON.stringify({ spellId: 85288, tag: 2 }), { spellId: 96103 }], // Warrior - Raging Blow Main-Hand
	[JSON.stringify({ spellId: 85288, tag: 3 }), { spellId: 85384 }], // Warrior - Raging Blow Off-Hand
	[JSON.stringify({ spellId: 1680, tag: 2 }), { spellId: 44949 }], // Warrior - Whirlwind Off-Hand
	[JSON.stringify({ spellId: 107570, tag: 2 }), { spellId: 145585 }], // Warrior - Storm Bolt Off-Hand

	// Monk - Zen Sphere
	[JSON.stringify({ spellId: 124081, tag: 3 }), { spellId: 124098 }],
	[JSON.stringify({ spellId: 124081, tag: 4 }), { spellId: 124101 }],
	[JSON.stringify({ spellId: 124081, tag: 5 }), { spellId: 125033 }],

	// Mage - Living Bomb
	[JSON.stringify({ spellId: 44457, tag: 2 }), { spellId: 44461 }], // Living Bomb Explosion

	// Warlock - Immolation Aura
	[JSON.stringify({ spellId: 104025, tag: 2 }), { spellId: 129476 }],
]);

export const defaultTargetIcon = 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_metamorphosis.jpg';

const petNameToActionId: Record<string, ActionId> = {
	'Ancient Guardian': ActionId.fromSpellId(86698),
	'Army of the Dead': ActionId.fromSpellId(42650),
	Bloodworm: ActionId.fromSpellId(50452),
	'Fallen Zandalari': ActionId.fromSpellId(138342),
	'Flame Orb': ActionId.fromSpellId(82731),
	'Frozen Orb': ActionId.fromSpellId(84721),
	Gargoyle: ActionId.fromSpellId(49206),
	Ghoul: ActionId.fromSpellId(46584),
	'Gnomish Flame Turret': ActionId.fromItemId(23841),
	'Greater Earth Elemental': ActionId.fromSpellId(2062),
	'Greater Fire Elemental': ActionId.fromSpellId(2894),
	'Primal Earth Elemental': ActionId.fromSpellId(2062),
	'Primal Fire Elemental': ActionId.fromSpellId(2894),
	'Mirror Image': ActionId.fromSpellId(55342),
	'Mirror Image T12 2pc': ActionId.fromSpellId(55342),
	'Rune Weapon': ActionId.fromSpellId(49028),
	Shadowfiend: ActionId.fromSpellId(34433),
	Mindbender: ActionId.fromSpellId(123040),
	'Spirit Wolf 1': ActionId.fromSpellId(51533),
	'Spirit Wolf 2': ActionId.fromSpellId(51533),
	Valkyr: ActionId.fromSpellId(71844),
	'Tentacle of the Old Ones': ActionId.fromSpellId(107818),
	Treant: ActionId.fromSpellId(33831),
	'Water Elemental': ActionId.fromSpellId(31687),
	Felhunter: ActionId.fromSpellId(691),
	Imp: ActionId.fromSpellId(688),
	Succubus: ActionId.fromSpellId(712),
	Voidwalker: ActionId.fromSpellId(697),
	Doomguard: ActionId.fromSpellId(18540),
	Infernal: ActionId.fromSpellId(1122),
	'Fel Imp': ActionId.fromSpellId(112866),
	Shivarra: ActionId.fromSpellId(112868),
	Observer: ActionId.fromSpellId(112869),
	Voidlord: ActionId.fromSpellId(112867),
	Terrorguard: ActionId.fromSpellId(112927),
	Abyssal: ActionId.fromSpellId(112921),
	'Grimoire: Imp': ActionId.fromSpellId(111859),
	'Grimoire: Voidwalker': ActionId.fromSpellId(111895),
	'Grimoire: Felhunter': ActionId.fromSpellId(111897),
	'Grimoire: Succubus': ActionId.fromSpellId(111896),
	Felguard: ActionId.fromSpellId(30146),
	'Wild Imp': ActionId.fromSpellId(114592),
	'Grimoire: Felguard': ActionId.fromSpellId(111898),
	Wrathguard: ActionId.fromSpellId(112870),
	'Xuen, The White Tiger': ActionId.fromSpellId(123904),
	'Earth Spirit': ActionId.fromSpellId(138121),
	'Storm Spirit': ActionId.fromSpellId(138122),
	'Fire Spirit': ActionId.fromSpellId(138123),
};

// https://wowhead.com/mop-classic/hunter-pets
const petNameToIcon: Record<string, string> = {
	Bat: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_bat.jpg',
	Bear: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_bear.jpg',
	'Bird of Prey': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	Boar: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_boar.jpg',
	'Burning Treant': 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_forceofnature.jpg',
	'Carrion Bird': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_vulture.jpg',
	Cat: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_cat.jpg',
	Chimaera: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_chimera.jpg',
	'Core Hound': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_corehound.jpg',
	Crab: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_crab.jpg',
	Crocolisk: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_crocolisk.jpg',
	Devilsaur: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_devilsaur.jpg',
	Dragonhawk: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_dragonhawk.jpg',
	'Fallen Zandalari': 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_animatedead.jpg',
	Felguard: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelguard.jpg',
	Felhunter: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonfelhunter.jpg',
	Infernal: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summoninfernal.jpg',
	Doomguard: 'https://wow.zamimg.com/images/wow/icons/large/warlock_summon_doomguard.jpg',
	'Ebon Imp': 'https://wow.zamimg.com/images/wow/icons/large/spell_nature_removecurse.jpg',
	'Fiery Imp': 'https://wow.zamimg.com/images/wow/icons/large/ability_warlock_empoweredimp.jpg',
	Gorilla: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_gorilla.jpg',
	Hyena: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_hyena.jpg',
	Imp: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonimp.jpg',
	Moth: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_moth.jpg',
	'Nether Ray': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_netherray.jpg',
	Owl: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_owl.jpg',
	Raptor: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_raptor.jpg',
	Ravager: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_ravager.jpg',
	Rhino: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_rhino.jpg',
	Scorpid: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_scorpid.jpg',
	Serpent: 'https://wow.zamimg.com/images/wow/icons/medium/spell_nature_guardianward.jpg',
	Silithid: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_silithid.jpg',
	Spider: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_spider.jpg',
	'Shale Spider': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_spider.jpg',
	'Spirit Beast': 'https://wow.zamimg.com/images/wow/icons/medium/ability_druid_primalprecision.jpg',
	'Spore Bat': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_sporebat.jpg',
	Succubus: 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_summonsuccubus.jpg',
	Tallstrider: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_tallstrider.jpg',
	Turtle: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_turtle.jpg',
	'Warp Stalker': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_warpstalker.jpg',
	Wasp: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_wasp.jpg',
	'Wind Serpent': 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_windserpent.jpg',
	Wolf: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_wolf.jpg',
	Worm: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_pet_worm.jpg',
	Fox: 'https://wow.zamimg.com/images/wow/icons/medium/inv_misc_monstertail_07.jpg',
};

export function getPetIconFromName(name: string): string | ActionId | undefined {
	return petNameToActionId[name] || petNameToIcon[name];
}

export const resourceTypeToIcon: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: '',
	[ResourceType.ResourceTypeHealth]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_elemental_mote_life01.jpg',
	[ResourceType.ResourceTypeMana]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_elemental_mote_mana.jpg',
	[ResourceType.ResourceTypeEnergy]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_shadow_shadowworddominate.jpg',
	[ResourceType.ResourceTypeRage]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_misc_emotionangry.jpg',
	[ResourceType.ResourceTypeChi]: 'https://wow.zamimg.com/images/wow/icons/medium/ability_monk_healthsphere.jpg',
	[ResourceType.ResourceTypeComboPoints]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_mace_2h_pvp410_c_01.jpg',
	[ResourceType.ResourceTypeFocus]: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_focusfire.jpg',
	[ResourceType.ResourceTypeRunicPower]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_sword_62.jpg',
	[ResourceType.ResourceTypeBloodRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_bloodpresence.jpg',
	[ResourceType.ResourceTypeFrostRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_frostpresence.jpg',
	[ResourceType.ResourceTypeUnholyRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_unholypresence.jpg',
	[ResourceType.ResourceTypeDeathRune]: '/mop/assets/img/death_rune.png',
	[ResourceType.ResourceTypeSolarEnergy]: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_eclipseorange.jpg',
	[ResourceType.ResourceTypeLunarEnergy]: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_eclipse.jpg',
	[ResourceType.ResourceTypeGenericResource]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_holybolt.jpg',
};

// Use this to connect a buff row to a cast row in the timeline view
export const buffAuraToSpellIdMap: Record<number, ActionId> = {
	96228: ActionId.fromSpellId(82174), // Synapse Springs - Agi
	96229: ActionId.fromSpellId(82174), // Synapse Springs - Str
	96230: ActionId.fromSpellId(82174), // Synapse Springs - Int

	132403: ActionId.fromSpellId(53600), // Shield of the Righteous
	138169: ActionId.fromSpellId(85256), // Paladin T15 Ret 4P Templar's Verdict
};
