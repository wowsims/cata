export type CacheHandlerOptions = {
	keysToKeep?: number;
};

export class CacheHandler<T> {
	keysToKeep: CacheHandlerOptions['keysToKeep'];
	private data = new Map<string, T>();

	constructor(options: CacheHandlerOptions = {}) {
		this.keysToKeep = options.keysToKeep;
	}

	has(id: string): boolean {
		return this.data.has(id);
	}

	get(id: string): T | undefined {
		return this.data.get(id);
	}

	set(id: string, result: T) {
		this.data.set(id, result);
		if (this.keysToKeep) this.keepMostRecent();
	}

	private keepMostRecent() {
		if (this.data.size > 2) {
			const keys = [...this.data.keys()];
			const keysToRemove = keys.slice(0, keys.length - 2);
			keysToRemove.forEach(key => this.data.delete(key));
		}
	}
}
