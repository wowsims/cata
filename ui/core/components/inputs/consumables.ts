import { Player } from '../../player';
import { Class, ConsumesSpec, Profession, Spec, Stat } from '../../proto/common';
import { Consumable } from '../../proto/db';
import { ActionId } from '../../proto_utils/action_id';
import { EventID, TypedEvent } from '../../typed_event';
import * as InputHelpers from '../input_helpers';
import { IconEnumValueConfig } from '../pickers/icon_enum_picker';
import { ActionInputConfig, ItemStatOption } from './stat_options';

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T;
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>;
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof ConsumesSpec;
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void;
	showWhen?: (player: Player<any>) => boolean;
}

function makeConsumeInputFactory<T extends number, SpecType extends Spec>(
	args: ConsumeInputFactoryArgs<T>,
): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<SpecType>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		const valueOptions = options.map(
			option =>
				({
					actionId: option.config.actionId,
					value: option.config.value,
					showWhen: (player: Player<SpecType>) =>
						(!option.config.showWhen || option.config.showWhen(player)) && (option.config.faction || player.getFaction()) == player.getFaction(),
				}) satisfies IconEnumValueConfig<Player<SpecType>, T>,
		);
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [{ value: 0, iconUrl: '', tooltip: 'None' } as unknown as IconEnumValueConfig<Player<SpecType>, T>].concat(valueOptions),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.gearChangeEmitter, player.professionChangeEmitter]),
			showWhen: (player: Player<any>) => (!args.showWhen || args.showWhen(player)) && valueOptions.some(option => option.showWhen?.(player)),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();
				if (newConsumes[args.consumesFieldName] === newValue) {
					return;
				}

				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}

///////////////////////////////////////////////////////////////////////////
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

export const ConjuredDarkRune = {
	actionId: ActionId.fromItemId(12662),
	value: 12662,
};
export const ConjuredHealthstone = {
	actionId: ActionId.fromItemId(5512),
	value: 5512,
};
export const ConjuredRogueThistleTea = {
	actionId: ActionId.fromItemId(7676),
	value: 7676,
	showWhen: <SpecType extends Spec>(player: Player<SpecType>) => player.getClass() == Class.ClassRogue,
};

export const CONJURED_CONFIG = [
	{ config: ConjuredRogueThistleTea, stats: [] },
	{ config: ConjuredHealthstone, stats: [Stat.StatStamina] },
	{ config: ConjuredDarkRune, stats: [Stat.StatIntellect] },
] as ConsumableStatOption<number>[];

export const makeConjuredInput = makeConsumeInputFactory({ consumesFieldName: 'conjuredId' });

export const ExplosiveBigDaddy = {
	actionId: ActionId.fromItemId(63396),
	value: 89637,
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
};

export const HighpoweredBoltGun = {
	actionId: ActionId.fromItemId(60223),
	value: 82207,
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
};

export const EXPLOSIVE_CONFIG = [
	{ config: ExplosiveBigDaddy, stats: [] },
	{ config: HighpoweredBoltGun, stats: [] },
] as ConsumableStatOption<number>[];
export const makeExplosivesInput = makeConsumeInputFactory({ consumesFieldName: 'explosiveId' });

///////////////////////////////////////////////////////////////////////////
//                                 Tinkers
///////////////////////////////////////////////////////////////////////////

export const TinkerHandsSynapseSprings = {
	actionId: ActionId.fromSpellId(82174),
	value: 82174,
};
export const TinkerHandsQuickflipDeflectionPlates = {
	actionId: ActionId.fromSpellId(82176),
	value: 82176,
};
export const TinkerHandsTazikShocker = {
	actionId: ActionId.fromSpellId(82179),
	value: 82179,
};
export const TinkerHandsSpinalHealingInjector = {
	actionId: ActionId.fromSpellId(82184),
	value: 82184,
};
export const TinkerHandsZ50ManaGulper = {
	actionId: ActionId.fromSpellId(82186),
	value: 82186,
};

export const TINKERS_HANDS_CONFIG = [
	{ config: TinkerHandsSynapseSprings, stats: [] },
	{ config: TinkerHandsQuickflipDeflectionPlates, stats: [] },
	{ config: TinkerHandsTazikShocker, stats: [] },
	{ config: TinkerHandsSpinalHealingInjector, stats: [] },
	{ config: TinkerHandsZ50ManaGulper, stats: [] },
] as ConsumableStatOption<number>[];

export const makeTinkerHandsInput = makeConsumeInputFactory({
	consumesFieldName: 'tinkerId',
	showWhen: (player: Player<any>) => player.hasProfession(Profession.Engineering),
});

export interface ConsumableInputOptions {
	consumesFieldName: keyof ConsumesSpec;
	setValue?: (eventID: EventID, player: Player<any>, newValue: number) => void;
}

export function makeConsumableInput(
	items: Consumable[],
	options: ConsumableInputOptions,
	tooltip?: string,
): InputHelpers.TypedIconEnumPickerConfig<Player<any>, number> {
	const valueOptions = items.map(item => ({
		value: item.id,
		iconUrl: item.icon,
		actionId: ActionId.fromItemId(item.id),
		tooltip: item.name,
	}));
	return {
		type: 'iconEnum',
		tooltip: tooltip,
		numColumns: items.length > 5 ? 2 : 1,
		values: [{ value: 0, iconUrl: '', tooltip: 'None' }].concat(valueOptions),
		equals: (a: number, b: number) => a === b,
		zeroValue: 0,
		changedEvent: (player: Player<any>) => player.consumesChangeEmitter,
		getValue: (player: Player<any>) => player.getConsumes()[options.consumesFieldName] as number,
		showWhen: (_: Player<any>) => !!valueOptions.length,
		setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
			if (options.setValue) {
				options.setValue(eventID, player, newValue);
			}

			const newConsumes = {
				...player.getConsumes(),
				[options.consumesFieldName]: newValue,
			};

			if (options.consumesFieldName === 'flaskId') {
				newConsumes.guardianElixirId = 0;
				newConsumes.battleElixirId = 0;
			}

			if (options.consumesFieldName === 'battleElixirId' || options.consumesFieldName === 'guardianElixirId') {
				newConsumes.flaskId = 0;
			}
			player.setConsumes(eventID, newConsumes);
		},
	};
}
