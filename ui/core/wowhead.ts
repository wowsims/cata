import { getLanguageCode } from './constants/lang';
import { CHARACTER_LEVEL } from './constants/mechanics';
import { Database } from './proto_utils/database';

export type WowheadTooltipItemParams = {
	/**
	 * @description Item ID
	 * @see item - mapped value from wowhead
	 * */
	itemId: number;
	/**
	 * @description Item level
	 * @see ilvl - mapped value from wowhead
	 * */
	itemLevel?: number;
	/**
	 * @description Level
	 * @see lvl - mapped value from wowhead
	 * */
	level?: number;
	/**
	 * @description Enchant
	 * @see ench - mapped value from wowhead
	 * */
	enchantId?: number;
	/**
	 * @description Gems
	 * @see gems - mapped value from wowhead
	 * */
	gemIds?: number[];
	/**
	 * @description Extra Socket
	 * @see sock - mapped value from wowhead
	 * */
	hasExtraSocket?: boolean;
	/**
	 * @description Item Set Pieces
	 * @see pcs - mapped value from wowhead
	 * */
	setPieceIds?: number[];
	/**
	 * @description Random Enchantment
	 * @see rand - mapped value from wowhead
	 * */
	randomEnchantmentId?: number;
	/**
	 * @description Reforges
	 * @see forg - mapped value from wowhead
	 * */
	reforgeId?: number;
	/**
	 * @description Upgrades
	 * @see upgd - mapped value from wowhead
	 * */
	upgradeStep?: number;
	/**
	 * @description Transmogrified to
	 * @see transmog - mapped value from wowhead
	 * */
	transmogId?: number;
};

export type WowheadTooltipSpellParams = {
	/**
	 * @description Spell ID
	 * @see spell - mapped value from wowhead
	 * */
	spellId: number;
	/**
	 * @description Level
	 * @see lvl - mapped value from wowhead
	 * */
	level?: number;
	/**
	 * @description Buff
	 * @see buff - mapped value from wowhead
	 * */
	useBuffAura?: boolean;
	/**
	 * @description Difficulty
	 * @see dd - mapped value from wowhead
	 * */
	difficultyId?: 14 | 15 | 16;
};

export const WOWHEAD_EXPANSION_ENV = 15;

export const buildWowheadTooltipDataset = async (options: WowheadTooltipItemParams | WowheadTooltipSpellParams) => {
	const lang = getLanguageCode();
	const params = new URLSearchParams();
	const langPrefix = lang && lang != 'en' ? lang + '.' : '';
	params.set('domain', `${langPrefix}mop-classic`);
	params.set('dataEnv', String(WOWHEAD_EXPANSION_ENV));

	params.set('lvl', String(options.level || CHARACTER_LEVEL));

	if ('spellId' in options) {
		if (options.spellId) {
			params.set('spell', String(options.spellId));
		}
		if (options.useBuffAura) {
			const data = await Database.getSpellIconData(options.spellId);
			if (data.hasBuff) params.set('buff', '1');
		}
	}

	if ('itemId' in options) {
		params.set('item', String(options.itemId));
		if (options.itemLevel) {
			params.set('ilvl', String(options.itemLevel));
		}
		if (options.gemIds?.length) {
			params.set('gems', options.gemIds.join(':'));
		}
		if (options.enchantId) {
			params.set('ench', String(options.enchantId));
		}
		if (options.reforgeId) {
			params.set('forg', String(options.reforgeId));
		}
		if (options.randomEnchantmentId) {
			params.set('rand', String(options.randomEnchantmentId));
		}
		if (typeof options.upgradeStep === 'number') {
			params.set('upgd', String(options.upgradeStep));
		}
		if (options.setPieceIds?.length) {
			params.set('pcs', options.setPieceIds.join(':'));
		}
		if (options.hasExtraSocket) {
			params.set('sock', '');
		}
		if (options.transmogId) {
			params.set('transmog', String(options.transmogId));
		}
	}

	return decodeURIComponent(params.toString());
};
