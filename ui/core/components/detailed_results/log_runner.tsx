// @ts-expect-error
import debounce from 'lodash/debounce';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, fragment } from 'tsx-vanilla';

import { SimLog } from '../../proto_utils/logs_parser';
import { TypedEvent } from '../../typed_event.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
export class LogRunner extends ResultComponent {
	private virtualScroll: CustomVirtualScroll | null = null;
	readonly showDebugChangeEmitter = new TypedEvent<void>('Show Debug');
	private showDebug = false;
	cacheOutput: {
		cacheKey: number | null;
		logs: SimLog[] | null;
		logsAsHTML: Element[] | null;
		logsAsText: string[] | null;
	} = { cacheKey: null, logs: null, logsAsHTML: null, logsAsText: null };

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'log-runner-root';
		super(config);

		// Existing setup code for the component...
		this.rootElem.innerHTML += `
			<div class="log-runner-actions">
				<input type="text" id="log-search-input" class="form-control" placeholder="Filter logs">
				<button id="log-runner-scroll-to-top-btn" class="btn btn-primary order-last">Top</button>
			</div>
			<div id="log-runner-logs-scroll" class="log-runner-scroll">
				<table class="metrics-table log-runner-table">
					<thead>
						<tr class="metrics-table-header-row">
							<th>Time</th>
							<th><div class="d-flex align-items-end">Event</div></th>
						</tr>
					</thead>
					<tbody id="log-runner-logs" class="log-runner-logs"></tbody>
				</table>
			</div>
		`;
		const searchInput = this.rootElem.querySelector('#log-search-input') as HTMLInputElement;

		// Use the 'input' event to trigger search as the user types
		const onSearchHandler = () => {
			this.searchLogs(searchInput.value);
		};
		searchInput.addEventListener('input', debounce(onSearchHandler, 150));
		const scrollToTopBtn = this.rootElem.querySelector('#log-runner-scroll-to-top-btn');
		scrollToTopBtn?.addEventListener('click', () => {
			this.virtualScroll?.scrollToTop();
		});
		new BooleanPicker<LogRunner>(this.rootElem.querySelector('.log-runner-actions')!, this, {
			extraCssClasses: ['show-debug-picker'],
			label: 'Show Debug Statements',
			inline: true,
			reverse: true,
			changedEvent: () => this.showDebugChangeEmitter,
			getValue: () => this.showDebug,
			setValue: (eventID, _logRunner, newValue) => {
				this.showDebug = newValue;
				this.showDebugChangeEmitter.emit(eventID);
			},
		});

		this.showDebugChangeEmitter.on(() => {
			const lastResults = this.getLastSimResult();
			this.onSimResult(lastResults);
		});
		this.initializeClusterize();
	}

	private initializeClusterize(): void {
		const scrollElem = this.rootElem.querySelector('#log-runner-logs-scroll') as HTMLElement;
		const contentElem = this.rootElem.querySelector('#log-runner-logs') as HTMLElement;

		this.virtualScroll = new CustomVirtualScroll({
			scrollContainer: scrollElem,
			contentContainer: contentElem,
			itemHeight: 30,
		});
	}

	searchLogs(searchQuery: string): void {
		// Regular expression to match quoted phrases or words
		const matchQuotesRegex = /"([^"]+)"|\S+/g;
		let match;
		const keywords: any[] = [];
		// Extract keywords and quoted phrases from the search query
		while ((match = matchQuotesRegex.exec(searchQuery))) {
			keywords.push(match[1] ? match[1].toLowerCase() : match[0].toLowerCase());
		}
		const filteredLogs = this.cacheOutput.logsAsHTML?.filter((_, index) => {
			const logText = this.cacheOutput.logsAsText![index];
			return keywords.every(keyword => {
				if (keyword.startsWith('"') && keyword.endsWith('"')) {
					// Remove quotes for exact phrase match
					return logText.includes(keyword.slice(1, -1));
				}
				return logText.includes(keyword);
			});
		});

		if (filteredLogs) {
			this.virtualScroll?.setItems(filteredLogs);
		}
	}

	onSimResult(resultData: SimResultData): void {
		this.virtualScroll?.setItems(this.getLogs(resultData) || []);
	}

	getLogs(resultData: SimResultData) {
		if (!resultData) return [];
		if (this.cacheOutput.cacheKey === resultData?.eventID) {
			return this.cacheOutput.logsAsHTML;
		}

		const validLogs = resultData.result.logs.filter(log => !log.isCastCompleted());
		this.cacheOutput.cacheKey = resultData?.eventID;
		this.cacheOutput.logs = validLogs;
		this.cacheOutput.logsAsHTML = validLogs.map(log => (
			<tr>
				<td className="log-timestamp">{log.formattedTimestamp()}</td>
				<td className="log-evdsfent">{log.toHTML(false)}</td>
			</tr>
		));
		this.cacheOutput.logsAsText = this.cacheOutput.logsAsHTML.map(element => fragmentToString(element).trim().toLowerCase());

		return this.cacheOutput.logsAsHTML;
	}
}

const fragmentToString = (element: Node | Element) => {
	const div = document.createElement('div');
	div.appendChild(element.cloneNode(true));
	return div.innerHTML;
};

class CustomVirtualScroll {
	private scrollContainer: HTMLElement;
	private contentContainer: HTMLElement;
	private items: Element[];
	private itemHeight: number;
	private visibleItemsCount: number;
	private startIndex: number;
	private placeholderTop: HTMLElement;
	private placeholderBottom: HTMLElement;

	constructor({ scrollContainer, contentContainer, itemHeight }: { scrollContainer: HTMLElement; contentContainer: HTMLElement; itemHeight: number }) {
		this.scrollContainer = scrollContainer;
		this.contentContainer = contentContainer;
		this.items = [];
		this.itemHeight = itemHeight;
		this.visibleItemsCount = 50; // +1 for buffer
		this.startIndex = 0;

		this.placeholderTop = document.createElement('div');
		this.placeholderBottom = document.createElement('div');
		contentContainer.prepend(this.placeholderTop);
		contentContainer.append(this.placeholderBottom);

		this.attachScrollListener();
	}

	scrollToTop(): void {
		this.scrollContainer.scrollTop = 0;
		this.startIndex = 0; // Reset startIndex to ensure items are updated correctly
		this.updateVisibleItems(); // Update the visible items after scrolling to top
	}

	setItems(newItems: CustomVirtualScroll['items']): void {
		// Adjust the type of newItems as needed
		this.items = newItems;
		this.scrollToTop();
	}

	private attachScrollListener(): void {
		this.scrollContainer.addEventListener('scroll', () => {
			const newIndex = Math.floor(this.scrollContainer.scrollTop / this.itemHeight);
			if (newIndex !== this.startIndex) {
				this.startIndex = newIndex;
				this.updateVisibleItems();
			}
		});
	}

	private updateVisibleItems(): void {
		const endIndex = this.startIndex + this.visibleItemsCount;
		const visibleItems = this.items.slice(this.startIndex, endIndex);

		// Reset content and adjust placeholders
		this.contentContainer.innerHTML = '';
		this.contentContainer.appendChild(this.placeholderTop);
		this.placeholderTop.style.height = `${this.startIndex * this.itemHeight}px`;
		const fragment = document.createDocumentFragment();
		visibleItems.forEach(item => fragment.appendChild(item));
		this.contentContainer.appendChild(fragment);

		this.contentContainer.appendChild(this.placeholderBottom);
		const remainingItems = this.items.length - endIndex;
		this.placeholderBottom.style.height = `${remainingItems * this.itemHeight}px`;
	}
}
