import { Button } from '@wowsims/ui';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Component } from './component';

export interface CopyButtonConfig {
	getContent: () => string;
	extraCssClasses?: string[];
	text?: string;
	tooltip?: string;
}

export class CopyButton extends Component {
	private readonly config: CopyButtonConfig;

	constructor(parent: HTMLElement, config: CopyButtonConfig) {
		const btnRef = ref<HTMLButtonElement>();
		const buttonElem = (
			<Button ref={btnRef} className={`btn ${config.extraCssClasses?.join(' ') ?? ''}`} iconLeft={{ icon: 'copy', className: 'me-1' }}>
				{config.text ?? 'Copy to Clipboard'}
			</Button>
		);

		super(parent, 'copy-button', buttonElem as HTMLElement);
		this.config = config;

		const button = btnRef.value!;
		let clicked = false;
		button.addEventListener('click', _event => {
			if (clicked) return;

			const data = this.config.getContent();
			if (navigator.clipboard == undefined) {
				alert(data);
			} else {
				clicked = true;
				navigator.clipboard.writeText(data);
				const originalContent = button.innerHTML;
				button.innerHTML = '<i class="fas fa-check me-1"></i>Copied';
				setTimeout(() => {
					button.innerHTML = originalContent;
					clicked = false;
				}, 1500);
			}
		});

		if (config.tooltip) {
			tippy(button, { content: config.tooltip });
		}
	}
}
