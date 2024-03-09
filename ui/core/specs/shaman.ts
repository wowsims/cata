import { IconSize } from "../class";
import { Shaman } from "../classes/shaman";
import { Spec as SpecProto } from '../proto/common.js';
import { getSpecSiteUrl } from "../proto_utils/utils";
import { Spec } from "../spec";

export class Elemental extends Spec {
	static protoID = SpecProto.SpecElementalShaman;
	static class = Shaman;
	static friendlyName = 'Elemental';
	static simLink = getSpecSiteUrl('elemental_shaman');

	readonly protoID = Elemental.protoID;
	readonly class = Elemental.class;
	readonly friendlyName = Elemental.friendlyName;
	readonly simLink = Enhancement.simLink;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightning.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Shaman.getIcon(size)
	}
}

export class Enhancement extends Spec {
	static protoID = SpecProto.SpecEnhancementShaman;
	static class = Shaman;
	static friendlyName = 'Enhancement';
	static simLink = getSpecSiteUrl('enhancement_shaman');

	readonly protoID = Enhancement.protoID;
	readonly class = Enhancement.class;
	readonly friendlyName = Enhancement.friendlyName;
	readonly simLink = Enhancement.simLink;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_lightningshield.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Shaman.getIcon(size)
	}
}

export class Restoration extends Spec {
	static protoID = SpecProto.SpecRestorationShaman;
	static class = Shaman;
	static friendlyName = 'Restoration';
	static simLink = getSpecSiteUrl('restoration_shaman');

	readonly protoID = Restoration.protoID;
	readonly class = Restoration.class;
	readonly friendlyName = Restoration.friendlyName;
	readonly simLink = Restoration.simLink;

	static getIcon = (size: IconSize): string => {
		return `https://wow.zamimg.com/images/wow/icons/${size}/spell_nature_magicimmunity.jpg`;
	}

	getIcon = (size: IconSize): string => {
		return Shaman.getIcon(size)
	}
}
