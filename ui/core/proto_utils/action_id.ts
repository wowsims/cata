import { getWowheadLanguagePrefix } from '../constants/lang';
import { CHARACTER_LEVEL } from '../constants/mechanics';
import { ActionID as ActionIdProto, ItemRandomSuffix, OtherAction, ReforgeStat } from '../proto/common';
import { ResourceType } from '../proto/spell';
import { IconData, UIItem as Item } from '../proto/ui';
import { buildWowheadTooltipDataset, WowheadTooltipItemParams, WowheadTooltipSpellParams } from '../wowhead';
import { Database } from './database';

// If true uses wotlkdb.com, else uses wowhead.com.
export const USE_WOTLK_DB = false;

// Uniquely identifies a specific item / spell / thing in WoW. This object is immutable.
export class ActionId {
	readonly itemId: number;
	readonly randomSuffixId: number;
	readonly reforgeId: number;
	readonly spellId: number;
	readonly otherId: OtherAction;
	readonly tag: number;

	readonly baseName: string; // The name without any tag additions.
	readonly name: string;
	readonly iconUrl: string;
	readonly spellIdTooltipOverride: number | null;

	private constructor(
		itemId: number,
		spellId: number,
		otherId: OtherAction,
		tag: number,
		baseName: string,
		name: string,
		iconUrl: string,
		randomSuffixId?: number,
		reforgeId?: number,
	) {
		this.itemId = itemId;
		this.randomSuffixId = randomSuffixId || 0;
		this.reforgeId = reforgeId || 0;
		this.spellId = spellId;
		this.otherId = otherId;
		this.tag = tag;

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
				if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				} else if (tag == 41570) {
					name += ' (Magmaw)';
				} else if (tag == 49416) {
					name += ' (Blazing Bone Construct)';
				} else if (tag == 56427) {
					name += ' (Warmaster Blackhorn)';
				} else if (tag == 56781) {
					name += ' (Goriona)';
				} else if (tag > 4191800) {
					name += ` (Animated Bone Warrior ${(tag - 4191800).toFixed(0)})`;
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
		this.baseName = baseName;
		this.name = name || baseName;
		this.iconUrl = iconUrl;
		this.spellIdTooltipOverride = this.spellTooltipOverride?.spellId || null;
	}

	anyId(): number {
		return this.itemId || this.spellId || this.otherId;
	}

	equals(other: ActionId): boolean {
		return this.equalsIgnoringTag(other) && this.tag == other.tag;
	}

	equalsIgnoringTag(other: ActionId): boolean {
		return this.itemId == other.itemId && this.randomSuffixId == other.randomSuffixId && this.spellId == other.spellId && this.otherId == other.otherId;
	}

	setBackground(elem: HTMLElement) {
		if (this.iconUrl) {
			elem.style.backgroundImage = `url('${this.iconUrl}')`;
		}
	}

	static makeItemUrl(id: number, randomSuffixId?: number, reforgeId?: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		const url = new URL(`https://wowhead.com/cata/${langPrefix}item=${id}`);
		url.searchParams.set('level', String(CHARACTER_LEVEL));
		url.searchParams.set('rand', String(randomSuffixId || 0));
		if (reforgeId) url.searchParams.set('forg', String(reforgeId));
		return url.toString();
	}
	static makeSpellUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return `https://wotlkdb.com/?spell=${id}`;
		} else {
			return `https://wowhead.com/cata/${langPrefix}spell=${id}`;
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
			return `https://wowhead.com/cata/${langPrefix}quest=${id}`;
		}
	}
	static makeNpcUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return 'https://wotlkdb.com/?npc=' + id;
		} else {
			return `https://wowhead.com/cata/${langPrefix}npc=${id}`;
		}
	}
	static makeZoneUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		if (USE_WOTLK_DB) {
			return 'https://wotlkdb.com/?zone=' + id;
		} else {
			return `https://wowhead.com/cata/${langPrefix}zone=${id}`;
		}
	}

	setWowheadHref(elem: HTMLAnchorElement) {
		if (this.itemId) {
			elem.href = ActionId.makeItemUrl(this.itemId, this.randomSuffixId, this.reforgeId);
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

	async fillAndSet(elem: HTMLAnchorElement, setHref: boolean, setBackground: boolean): Promise<ActionId> {
		const filled = await this.fill();
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

		// handle DRT
		let tag = this.tag;
		if (tag >= 71086) {
			name = 'Dragonwrath - ' + name;
			tag -= 71086;
		}

		switch (baseName) {
			case 'Minor Speed':
				name = 'Minor Run Speed (8%)';
				break;
			case 'Explosive Shot':
				if (this.spellId == 53301) {
					name += ' (First)';
				} else if (this.spellId == 1215485) {
					name += ' (Second)';
				}
				break;
			case 'Explosive Trap':
				if (tag == 1) {
					name += ' (Weaving)';
				}
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
				}
				break;
			case 'Shattering Throw':
				if (tag === -1) {
					name += ' (raid)';
				} else {
					name += ` (self)`;
				}
				break;
			case 'Envenom':
			case 'Eviscerate':
			case 'Expose Armor':
			case 'Rupture':
			case 'Slice and Dice':
			case 'Recuperate':
				if (tag) name += ` (${tag} CP)`;
				break;
			case 'Instant Poison':
			case 'Wound Poison':
				if (tag == 1) {
					name += ' (Deadly)';
				} else if (tag == 2) {
					name += ' (Shiv)';
				} else if (tag == 3) {
					name += ' (Fan of Knives)';
				}
				break;
			case 'Fan of Knives':
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
				if (this.spellId == 89775) {
					name += ' (DoT)';
				}
				break;
			case 'Stormstrike':
				if (tag == 0) {
					name += ' (Cast)';
				} else if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Chain Lightning':
			case 'Lightning Bolt':
			case 'Lava Burst':
				if (tag == 6) {
					name += ' (Overload)';
				} else if (tag) {
					name += ` (${tag} MW)`;
				}
				break;
			case 'Flame Shock':
				if (tag == 1) {
					name += ' (DoT)';
				}
				break;
			case 'Fulmination':
				name += ` (${tag + 3})`;
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
			case 'Holy Shield':
				if (tag == 1) {
					name += ' (Proc)';
				}
				break;
			case 'Censure':
				if (tag == 2) {
					name += ' (DoT)';
				}
				break;
			case 'Exorcism':
				if (tag === 3) {
					name = 'Glyph of Exorcism (DoT)';
				}
				break;
			case 'Seal of Righteousness':
				if (tag === 2) {
					name += ' (DS)';
				}
				break;
			// For targetted buffs, tag is the source player's raid index or -1 if none.
			case 'Bloodlust':
			case 'Ferocious Inspiration':
			case 'Innervate':
			case 'Focus Magic':
			case 'Mana Tide Totem':
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
			case 'Rune Strike':
				if (tag == 0) {
					name += ' (Queue)';
				} else if (tag == 1) {
					name += ' (Main Hand)';
				} else if (tag == 2) {
					name += ' (Off Hand)';
				}
				break;
			case 'Raging Blow':
			case 'Whirlwind':
			case 'Slam':
			case 'Frost Strike':
			case 'Plague Strike':
			case 'Blood Strike':
			case 'Obliterate':
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
				// Warrior - T12 4P proc
				if (baseName === 'Raging Blow' && tag === 3) {
					name = 'Fiery attack';
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
				} else if (tag == 3) {
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
					name += ' (Improved)';
					break;
				}
			case 'Improved Steady Shot':
				if (tag == 2) {
					name += ' (pre)';
				}
				break;
			case 'Immolate':
				if (tag == 1) {
					name += ' (DoT)';
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
			case 'Virtuous Empowerment':
				name = 'Battleplate of Radiant Glory - T13 2pc';
				break;
			case 'Hurricane':
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
			default:
				if (tag) {
					name += ' (??)';
				}
				break;
		}

		const iconOverrideId = this.spellTooltipOverride || this.spellIconOverride;
		let iconUrl = ActionId.makeIconUrl(tooltipData['icon']);
		if (iconOverrideId) {
			const overrideTooltipData = await ActionId.getTooltipData(iconOverrideId, { signal: options?.signal });
			iconUrl = ActionId.makeIconUrl(overrideTooltipData['icon']);
		}

		return new ActionId(this.itemId, this.spellId, this.otherId, this.tag, baseName, name, iconUrl, this.randomSuffixId, this.reforgeId);
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
			throw new Error('Empty action id!');
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
		return new ActionId(this.itemId, this.spellId, this.otherId, 0, this.baseName, this.baseName, this.iconUrl, this.randomSuffixId, this.reforgeId);
	}

	static fromEmpty(): ActionId {
		return new ActionId(0, 0, OtherAction.OtherActionNone, 0, '', '', '');
	}

	static fromItemId(itemId: number, tag?: number, randomSuffixId?: number, reforgeId?: number): ActionId {
		return new ActionId(itemId, 0, OtherAction.OtherActionNone, tag || 0, '', '', '', randomSuffixId, reforgeId);
	}

	static fromSpellId(spellId: number, tag?: number): ActionId {
		return new ActionId(0, spellId, OtherAction.OtherActionNone, tag || 0, '', '', '');
	}

	static fromOtherId(otherId: OtherAction, tag?: number): ActionId {
		return new ActionId(0, 0, otherId, tag || 0, '', '', '');
	}

	static fromPetName(petName: string): ActionId {
		return petNameToActionId[petName] || new ActionId(0, 0, OtherAction.OtherActionPet, 0, petName, petName, petNameToIcon[petName] || '');
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
		return new ActionId(
			idType == 'ItemID' ? id : 0,
			idType == 'SpellID' ? id : 0,
			idType == 'OtherID' ? id : 0,
			match[7] ? parseInt(match[7]) : 0,
			'',
			'',
			'',
		);
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
	[JSON.stringify({ spellId: 13889 }), { spellId: 109709 }], // Minor Run Speed
]);

const spellIdTooltipOverrides: Map<string, ActionIdOverride> = new Map([
	[JSON.stringify({ spellId: 47897, tag: 1 }), { spellId: 47960 }], // Shadowflame Dot
	[JSON.stringify({ spellId: 55090, tag: 1 }), { spellId: 70890 }], // Shadowflame Dot
	[JSON.stringify({ spellId: 12294, tag: 3 }), { spellId: 99237 }], // Warrior - T12 4P Fiery Attack - Mortal Strike
	[JSON.stringify({ spellId: 85288, tag: 3 }), { spellId: 99237 }], // Warrior - T12 4P Fiery Attack - Raging Blow
	[JSON.stringify({ spellId: 23922, tag: 3 }), { spellId: 99240 }], // Warrior - T12 2P Combust - Shield Slam
	[JSON.stringify({ spellId: 77767, tag: 3 }), { spellId: 99058 }], // Hunter - T12 2P Flaming Arrow - Cobra shot
	[JSON.stringify({ spellId: 56641, tag: 3 }), { spellId: 99058 }], // Hunter - T12 2P Flaming Arrow - Steady shot
	[JSON.stringify({ spellId: 35395, tag: 3 }), { spellId: 99092 }], // Paladin - T12 2P Flames of the Faithful
	[JSON.stringify({ spellId: 53600, tag: 3 }), { spellId: 99075 }], // Paladin - T12 2P Righteous Flames
	[JSON.stringify({ spellId: 879, tag: 3 }), { spellId: 54934 }], // Paladin - Glyph of Exorcism
	[JSON.stringify({ spellId: 49020, tag: 3 }), { spellId: 99000 }], // Death Knight - T12 4P Flaming Torment
	[JSON.stringify({ spellId: 55090, tag: 3 }), { spellId: 99000 }], // Death Knight - T12 4P Flaming Torment

	// Off-Hand attacks
	[JSON.stringify({ spellId: 45902, tag: 2 }), { spellId: 66215 }], // Death Knight - Blood Strike Off-Hand
	[JSON.stringify({ spellId: 45462, tag: 2 }), { spellId: 49998 }], // Death Knight - Death Strike Off-Hand
	[JSON.stringify({ spellId: 85948, tag: 2 }), { spellId: 86061 }], // Death Knight - Festering Strike Off-Hand
	[JSON.stringify({ spellId: 49143, tag: 2 }), { spellId: 66196 }], // Death Knight - Frost Strike Off-Hand
	[JSON.stringify({ spellId: 49020, tag: 2 }), { spellId: 66198 }], // Death Knight - Obliterate Off-Hand
	[JSON.stringify({ spellId: 45462, tag: 2 }), { spellId: 66216 }], // Death Knight - Plague Strike Off-Hand
	[JSON.stringify({ spellId: 56815, tag: 2 }), { spellId: 66217 }], // Death Knight - Rune Strike Off-Hand
	[JSON.stringify({ spellId: 1329, tag: 2 }), { spellId: 27576 }], // Rogue - Mutilate Off-Hand
	[JSON.stringify({ spellId: 17364, tag: 2 }), { spellId: 32176 }], // Shaman - Stormstrike Off-Hand
	[JSON.stringify({ spellId: 85288, tag: 2 }), { spellId: 85384 }], // Warrior - Raging Blow Off-Hand
	[JSON.stringify({ spellId: 1464, tag: 2 }), { spellId: 97992 }], // Warrior - Slam Off-Hand
	[JSON.stringify({ spellId: 1680, tag: 2 }), { spellId: 44949 }], // Warrior - Whirlwind Off-Hand
]);

export const defaultTargetIcon = 'https://wow.zamimg.com/images/wow/icons/large/spell_shadow_metamorphosis.jpg';

const petNameToActionId: Record<string, ActionId> = {
	'Ancient Guardian': ActionId.fromSpellId(86150),
	'Army of the Dead': ActionId.fromSpellId(42650),
	Bloodworm: ActionId.fromSpellId(50452),
	'Flame Orb': ActionId.fromSpellId(82731),
	Gargoyle: ActionId.fromSpellId(49206),
	Ghoul: ActionId.fromSpellId(46584),
	'Gnomish Flame Turret': ActionId.fromItemId(23841),
	'Greater Earth Elemental': ActionId.fromSpellId(2062),
	'Greater Fire Elemental': ActionId.fromSpellId(2894),
	'Mirror Image': ActionId.fromSpellId(55342),
	'Mirror Image T12 2pc': ActionId.fromSpellId(55342),
	'Rune Weapon': ActionId.fromSpellId(49028),
	Shadowfiend: ActionId.fromSpellId(34433),
	'Spirit Wolf 1': ActionId.fromSpellId(51533),
	'Spirit Wolf 2': ActionId.fromSpellId(51533),
	Valkyr: ActionId.fromSpellId(71844),
	'Tentacle of the Old Ones': ActionId.fromSpellId(107818),
	Treant: ActionId.fromSpellId(33831),
	'Water Elemental': ActionId.fromSpellId(31687),
};

// https://wowhead.com/cata/hunter-pets
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
	[ResourceType.ResourceTypeComboPoints]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_mace_2h_pvp410_c_01.jpg',
	[ResourceType.ResourceTypeFocus]: 'https://wow.zamimg.com/images/wow/icons/medium/ability_hunter_focusfire.jpg',
	[ResourceType.ResourceTypeRunicPower]: 'https://wow.zamimg.com/images/wow/icons/medium/inv_sword_62.jpg',
	[ResourceType.ResourceTypeBloodRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_bloodpresence.jpg',
	[ResourceType.ResourceTypeFrostRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_frostpresence.jpg',
	[ResourceType.ResourceTypeUnholyRune]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_deathknight_unholypresence.jpg',
	[ResourceType.ResourceTypeDeathRune]: '/cata/assets/img/death_rune.png',
	[ResourceType.ResourceTypeSolarEnergy]: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_eclipseorange.jpg',
	[ResourceType.ResourceTypeLunarEnergy]: 'https://wow.zamimg.com/images/wow/icons/large/ability_druid_eclipse.jpg',
	[ResourceType.ResourceTypeHolyPower]: 'https://wow.zamimg.com/images/wow/icons/medium/spell_holy_holybolt.jpg',
};

// Use this to connect a buff row to a cast row in the timeline view
export const buffAuraToSpellIdMap: Record<number, ActionId> = {
	96228: ActionId.fromSpellId(82174), // Synapse Springs - Agi
	96229: ActionId.fromSpellId(82174), // Synapse Springs - Str
	96230: ActionId.fromSpellId(82174), // Synapse Springs - Int
};
