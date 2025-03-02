import { ref } from 'tsx-vanilla';

import { SimUI } from '../sim_ui';
import { BaseModal } from './base_modal';
import Toast from './toast';

export interface ImporterOptions {
	title: string;
	allowFileUpload?: boolean;
}

export abstract class Importer extends BaseModal {
	protected abstract readonly simUI: SimUI;

	protected readonly textElem: HTMLTextAreaElement;
	protected readonly descriptionElem: HTMLElement;
	protected readonly importButton: HTMLButtonElement;
	private readonly allowFileUpload: boolean;

	constructor(parent: HTMLElement, options: ImporterOptions) {
		super(parent, 'importer', { title: options.title, footer: true, disposeOnClose: false });

		this.allowFileUpload = options.allowFileUpload || false;
		const uploadInputId = 'upload-input-' + options.title.toLowerCase().replaceAll(' ', '-');

		const descriptionElemRef = ref<HTMLDivElement>();
		const textElemRef = ref<HTMLTextAreaElement>();
		const importButtonRef = ref<HTMLButtonElement>();
		const uploadInputRef = ref<HTMLInputElement>();

		this.body.replaceChildren(
			<>
				<div ref={descriptionElemRef} className="import-description"></div>
				<textarea ref={textElemRef} className="importer-textarea form-control" attributes={{ spellcheck: false }}></textarea>
			</>,
		);

		this.footer!.appendChild(
			<>
				{this.allowFileUpload && (
					<>
						<label htmlFor={uploadInputId} className="importer-button btn btn-primary upload-button me-2">
							<i className="fas fa-file-arrow-up me-1"></i>
							Upload File
						</label>
						<input ref={uploadInputRef} type="file" id={uploadInputId} className="importer-upload-input d-none" hidden />
					</>
				)}
				<button ref={importButtonRef} className="importer-button btn btn-primary import-button">
					<i className="fa fa-download me-1"></i>
					Import
				</button>
			</>,
		);

		this.descriptionElem = descriptionElemRef.value!;
		this.textElem = textElemRef.value!;

		if (this.allowFileUpload && uploadInputRef.value) {
			uploadInputRef.value.addEventListener('change', async event => {
				this.textElem.textContent = await (event as any).target.files[0].text();
			});
		}

		this.importButton = importButtonRef.value!;
		this.importButton.addEventListener('click', async _event => {
			try {
				await this.onImport(this.textElem.value || '');
			} catch (error: any) {
				new Toast({ variant: 'error', body: `Import error: ${error?.message || error}` });
			}
		});
	}

	abstract onImport(data: string): Promise<void>;
}
