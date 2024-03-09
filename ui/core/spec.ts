import { Class, IconSize } from './class.js';
import { Spec as SpecProto } from './proto/common.js';

export abstract class Spec {
	abstract readonly protoID: SpecProto;
	abstract readonly class: Class;
	abstract readonly friendlyName: string;
	abstract readonly simLink: string;

	abstract getIcon(size: IconSize): string;
}
