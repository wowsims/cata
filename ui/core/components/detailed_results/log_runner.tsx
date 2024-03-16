// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { TypedEvent } from '../../typed_event.js';
import { BooleanPicker } from '../boolean_picker.js';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';
export class LogRunner extends ResultComponent {
    private virtualScroll: CustomVirtualScroll | null = null;
	readonly showDebugChangeEmitter = new TypedEvent<void>('Show Debug');
	private showDebug = false;
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'log-runner-root';
		super(config);

		// Existing setup code for the component...
		this.rootElem.innerHTML += `
			<div class="show-debug-container"></div>
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
		

		new BooleanPicker<LogRunner>(this.rootElem.querySelector('.show-debug-container')!, this, {
			extraCssClasses: ['show-debug-picker'],
			label: 'Show Debug Statements',
			inline: true,
			reverse: true,
			changedEvent: () => this.showDebugChangeEmitter,
			getValue: () => this.showDebug,
			setValue: (eventID, _logRunner, newValue) => {
				this.showDebug = newValue;
				this.showDebugChangeEmitter.emit(eventID);
			}
		});

		this.showDebugChangeEmitter.on(() => {
			// Refresh the logs display based on the new showDebug state.
			if (this.getLastSimResult()) {
				this.onSimResult(this.getLastSimResult());
			}
		});
		this.initializeClusterize();
	}

    private initializeClusterize(): void {
		const scrollElem = this.rootElem.querySelector('#log-runner-logs-scroll') as HTMLElement;
		const contentElem = this.rootElem.querySelector('#log-runner-logs') as HTMLElement;

		this.virtualScroll = new CustomVirtualScroll({
			scrollContainer: scrollElem,
			contentContainer: contentElem,
			itemHeight: 30 s
		});
    }

    onSimResult(resultData: SimResultData): void {
        let logs = resultData.result.logs.filter(log => !log.isCastCompleted()).map(log => `<tr><td class="log-timestamp">${log.formattedTimestamp()}</td><td class="log-event">${log.toString(false).trim()}</td></tr>`);
       	this.virtualScroll?.setItems(logs);
    }
}

class CustomVirtualScroll {
  private scrollContainer: HTMLElement;
  private contentContainer: HTMLElement;
  private items: string[]; 
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

  setItems(newItems: string[]): void { // Adjust the type of newItems as needed
    this.items = newItems;
    this.updateVisibleItems();
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

    visibleItems.forEach(item => {
      const itemElement = document.createElement('div'); // Or any other element that suits your items
      itemElement.innerHTML = item;
      this.contentContainer.appendChild(itemElement);
    });

    this.contentContainer.appendChild(this.placeholderBottom);
    const remainingItems = this.items.length - endIndex;
    this.placeholderBottom.style.height = `${remainingItems * this.itemHeight}px`;
  }
}
