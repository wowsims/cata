import { Button } from '@wowsims/ui';
// @ts-expect-error
import debounce from 'lodash/debounce';
import { ref } from 'tsx-vanilla';

import { SimLog } from '../../proto_utils/logs_parser';
import { TypedEvent } from '../../typed_event.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
export class LogRunner extends ResultComponent {
	private virtualScroll: CustomVirtualScroll | null = null;
	readonly showDebugChangeEmitter = new TypedEvent<void>('Show Debug');
	private showDebug = false;
	private ui: {
		search: HTMLInputElement;
		actions: HTMLDivElement;
		buttonToTop: HTMLButtonElement;
		scrollContainer: HTMLDivElement;
		contentContainer: HTMLTableSectionElement;
	};
	cacheOutput: {
		cacheKey: number | null;
		logs: SimLog[] | null;
		logsAsHTML: Element[] | null;
		logsAsText: string[] | null;
	} = { cacheKey: null, logs: null, logsAsHTML: null, logsAsText: null };

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'log-runner-root';
		super(config);

		const searchRef = ref<HTMLInputElement>();
		const actionsRef = ref<HTMLDivElement>();
		const buttonToTopRef = ref<HTMLButtonElement>();
		const scrollContainerRef = ref<HTMLDivElement>();
		const contentContainerRef = ref<HTMLTableSectionElement>();

		this.rootElem.appendChild(
			<>
				<div ref={actionsRef} className="log-runner-actions">
					<input ref={searchRef} type="text" className="form-control log-search-input" placeholder="Filter logs" />
					<Button variant="primary" ref={buttonToTopRef} className="order-last log-runner-scroll-to-top-btn">
						Top
					</Button>
				</div>
				<div ref={scrollContainerRef} className="log-runner-scroll">
					<table className="metrics-table log-runner-table">
						<thead>
							<tr className="metrics-table-header-row">
								<th>Time</th>
								<th>
									<div className="d-flex align-items-end">Event</div>
								</th>
							</tr>
						</thead>
						<tbody ref={contentContainerRef} className="log-runner-logs"></tbody>
					</table>
				</div>
			</>,
		);

		this.ui = {
			search: searchRef.value!,
			actions: actionsRef.value!,
			buttonToTop: buttonToTopRef.value!,
			scrollContainer: scrollContainerRef.value!,
			contentContainer: contentContainerRef.value!,
		};

		// Use the 'input' event to trigger search as the user types
		const onSearchHandler = () => {
			this.searchLogs(this.ui.search.value);
		};
		this.ui.search?.addEventListener('input', debounce(onSearchHandler, 150));
		this.ui.buttonToTop?.addEventListener('click', () => {
			this.virtualScroll?.scrollToTop();
		});
		new BooleanPicker<LogRunner>(this.ui.actions, this, {
			id: 'log-runner-show-debug',
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
		this.virtualScroll = new CustomVirtualScroll({
			scrollContainer: this.ui.scrollContainer,
			contentContainer: this.ui.contentContainer,
			itemHeight: 32,
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
		this.cacheOutput.logsAsHTML = validLogs.map(log => this.renderItem(log));
		this.cacheOutput.logsAsText = this.cacheOutput.logsAsHTML.map(element => fragmentToString(element).trim().toLowerCase());

		return this.cacheOutput.logsAsHTML;
	}

	renderItem(log: SimLog) {
		return (
			<tr>
				<td className="log-timestamp">{log.formattedTimestamp()}</td>
				<td className="log-evdsfent">{log.toHTML(false)}</td>
			</tr>
		) as HTMLTableRowElement;
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
		const remainingItems = this.items.length - endIndex;

		// Update the height of the placeholders before it's placed in the dom to prevent rerender
		this.placeholderTop.style.height = `${this.startIndex * this.itemHeight}px`;
		this.placeholderBottom.style.height = `${remainingItems * this.itemHeight}px`;
		this.contentContainer.replaceChildren(
			<>
				{this.placeholderTop}
				{visibleItems.map(item => item)}
				{this.placeholderBottom}
			</>,
		);
	}
}
