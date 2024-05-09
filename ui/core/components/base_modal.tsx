import { Modal } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, ref } from 'tsx-vanilla';

import { Component } from './component';

type ModalSize = 'sm' | 'md' | 'lg' | 'xl';

type BaseModalConfig = {
	closeButton?: {
		// When true, the button will be rendered in a fixed position on the screen.
		// Primarily used for the raid sim's embedded player editors
		fixed?: boolean;
	};
	// Whether or not to add a modal-footer element
	footer?: boolean;
	// Whether or not to add a modal-header element
	header?: boolean;
	// Whether or not to allow modal contents to extend past the screen height.
	// When true, the modal is fixed to the screen height and body contents will scroll.
	scrollContents?: boolean;
	// Specify the size of the modal
	size?: ModalSize;
	// A title for the modal
	title?: string | null;
	// Should the modal be disposed on close?
	disposeOnClose?: boolean;
};

const DEFAULT_CONFIG = {
	footer: false,
	header: true,
	scrollContents: false,
	size: 'lg' as ModalSize,
	title: null,
};

export class BaseModal extends Component {
	readonly modalConfig: BaseModalConfig;

	isOpen = false;
	onHideCallbacks: Array<() => void> = [];

	readonly modal: Modal;
	readonly dialog: HTMLElement;
	readonly header: HTMLElement | undefined;
	readonly body: HTMLElement;
	readonly footer: HTMLElement | undefined;

	constructor(parent: HTMLElement, cssClass: string, config: BaseModalConfig = { disposeOnClose: true }) {
		super(parent, 'modal');
		this.modalConfig = { ...DEFAULT_CONFIG, ...config };

		const dialogRef = ref<HTMLDivElement>();
		const headerRef = ref<HTMLDivElement>();
		const bodyRef = ref<HTMLDivElement>();
		const footerRef = ref<HTMLDivElement>();

		const modalSizeKlass = this.modalConfig.size && this.modalConfig.size != 'md' ? `modal-${this.modalConfig.size}` : '';

		this.rootElem.classList.add('fade');
		this.rootElem.appendChild(
			<div className={`modal-dialog ${cssClass} ${modalSizeKlass} ${this.modalConfig.scrollContents ? 'modal-overflow-scroll' : ''}`} ref={dialogRef}>
				<div className="modal-content">
					<div className={`modal-header ${this.modalConfig.header || this.modalConfig.title ? '' : 'p-0 border-0'}`} ref={headerRef}>
						{this.modalConfig.title && <h5 className="modal-title">{this.modalConfig.title}</h5>}
						<button
							type="button"
							className={`btn-close ${this.modalConfig.closeButton?.fixed ? 'position-fixed' : ''}`}
							onclick={() => this.close()}
							attributes={{ 'aria-label': 'Close' }}>
							<i className="fas fa-times fa-2xl"></i>
						</button>
					</div>
					<div className="modal-body" ref={bodyRef} />
					{this.modalConfig.footer && <div className="modal-footer" ref={footerRef} />}
				</div>
			</div>,
		);

		this.dialog = dialogRef.value!;
		this.header = headerRef.value!;
		this.body = bodyRef.value!;
		this.footer = footerRef.value!;

		this.modal = new Modal(this.rootElem);

		if (this.modalConfig.disposeOnClose) {
			this.rootElem.addEventListener('hidden.bs.modal', _ => {
				this.rootElem.remove();
				this.dispose();
			});
		}
	}

	open() {
		const closeModalOnEscKey = this.closeModalOnEscKey.bind(this);
		const showBSFn = this.showBSFn.bind(this);
		const hideBSFn = this.hideBSFn.bind(this);
		const hiddenBSFn = this.hiddenBSFn.bind(this);

		document.addEventListener('keydown', closeModalOnEscKey);
		this.rootElem.addEventListener('show.bs.modal', showBSFn);
		this.rootElem.addEventListener('hide.bs.modal', hideBSFn);
		this.rootElem.addEventListener('hidden.bs.modal', hiddenBSFn);

		this.addOnHideCallback(() => document.removeEventListener('keydown', closeModalOnEscKey));
		this.addOnHideCallback(() => this.rootElem.removeEventListener('show.bs.modal', showBSFn));
		this.addOnHideCallback(() => this.rootElem.removeEventListener('hide.bs.modal', hideBSFn));
		this.addOnHideCallback(() => this.rootElem.removeEventListener('hidden.bs.modal', hiddenBSFn));

		this.modal.show();
		this.isOpen = true;
	}

	close() {
		this.modal.hide();
		this.isOpen = false;
	}

	// Allows you to add a callback that will be run when the modal is hidden. Primarily used for disposing of event listeners on hide
	addOnHideCallback(fn: () => void) {
		this.onHideCallbacks.push(fn);
	}

	// Callbacks for on show and on hide
	protected onShow(_e: Event) {
		return;
	}

	protected onHide(_e: Event) {
		this.onHideCallbacks.forEach(callback => callback());
		return;
	}

	// Hacks for better looking multi modals
	private async showBSFn(event: Event) {
		// Prevent the event from bubbling up to parent modals
		event.stopImmediatePropagation();

		// Wait for the backdrop to be injected into the DOM
		const backdrop = (await new Promise(resolve => {
			setTimeout(() => {
				// @ts-ignore
				if (this.modal._backdrop._element)
					// @ts-ignore
					resolve(this.modal._backdrop._element);
			}, 100);
		})) as HTMLElement;
		// Then move it from <body> to the parent element
		this.rootElem.insertAdjacentElement('afterend', backdrop);
		this.onShow(event);
	}

	private hideBSFn(event: Event) {
		// Prevent the event from bubbling up to parent modals
		event.stopImmediatePropagation();
		this.onHide(event);
	}

	private hiddenBSFn(event: Event) {
		// Prevent the event from bubbling up to parent modals
		// Do not use stopImmediatePropagation here. It prevents Bootstrap from removing the modal,
		// leading to other issues
		event.stopPropagation();
	}

	private closeModalOnEscKey(event: KeyboardEvent) {
		if (event.key == 'Escape') {
			this.close();
		}
	}
}
