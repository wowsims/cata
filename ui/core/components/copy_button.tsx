import { Button, Icon } from '@wowsims/ui';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { cloneChildren } from '../utils';
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
			<Button ref={btnRef} className={`btn ${config.extraCssClasses?.join(' ') ?? ''}`} iconLeft="copy">
				{config.text ?? 'Copy to Clipboard'}
			</Button>
		);

		super(parent, 'copy-button', buttonElem as HTMLElement);
		this.config = config;

		const button = btnRef.value!;
		button.addEventListener('click', () => {
			if (button.disabled) return;

			const data = this.config.getContent();
			if (!navigator.clipboard) {
				alert(data);
			} else {
				navigator.clipboard.writeText(data);
				const defaultState = cloneChildren(button);
				button.disabled = true;
				button.replaceChildren(
					<>
						<Icon icon="check" className="me-1" />
						Copied
					</>,
				);
				setTimeout(() => {
					button.replaceChildren(...defaultState);
					button.disabled = false;
				}, 1500);
			}
		});

		if (config.tooltip) {
			tippy(button, { content: config.tooltip });
		}
	}
}
