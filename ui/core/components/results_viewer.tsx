import { Icon, Link, LinkIcon } from '@wowsims/ui';
import tippy, { inlinePositioning, Instance as TippyInstance, Props as TippyProps } from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { Component } from '../components/component.js';
import { TypedEvent } from '../typed_event.js';

// Config for displaying a warning to the user whenever a condition is met.
interface SimWarning {
	updateOn: TypedEvent<any>;
	getContent: () => string | Array<string>;
}

interface WarningLinkArgs {
	parent: HTMLElement;
	href?: string;
	text?: string;
	icon?: LinkIcon;
	tooltip?: HTMLElement | Element;
	classes?: string;
}

const TOOLTIP_HTML_BASE = <ul className="text-start ps-3 mb-0"></ul>;

export class ResultsViewer extends Component {
	readonly pendingElem: HTMLDivElement;
	readonly contentElem: HTMLDivElement;
	readonly warningElem: HTMLDivElement;
	private warningsLink: HTMLElement;

	private warnings: Array<SimWarning> = [];
	private warningsTooltip: TippyInstance | null = null;

	constructor(parentElem: HTMLElement) {
		super(parentElem, 'results-viewer');

		const pendingElemRef = ref<HTMLDivElement>();
		const contentElemRef = ref<HTMLDivElement>();
		const warningElemRef = ref<HTMLDivElement>();

		this.rootElem.appendChild(
			<>
				<div ref={pendingElemRef} className="results-pending">
					<div className="loader"></div>
				</div>
				<div ref={contentElemRef} className="results-content"></div>
				<div ref={warningElemRef} className="warning-zone text-center"></div>
			</>,
		);
		this.pendingElem = pendingElemRef.value!;
		this.contentElem = contentElemRef.value!;
		this.warningElem = warningElemRef.value!;

		this.warningsLink = this.addWarningsLink();
		this.updateWarnings();

		this.hideAll();
	}

	private addWarningLink(args: WarningLinkArgs): HTMLElement {
		const item = (
			<div className="sim-toolbar-item">
				<Link
					as={args.href ? undefined : 'button'}
					href={args.href ? args.href : undefined}
					target={args.href ? '_blank' : '_self'}
					className={args.classes}
					iconRight={args.icon}>
					{args.text || ''}
				</Link>
			</div>
		) as HTMLElement;

		args.parent.appendChild(item);

		if (args.tooltip) {
			this.warningsTooltip = tippy(item, {
				appendTo: 'parent',
				content: args.tooltip,
				placement: 'bottom',
				inlinePositioning: true,
				plugins: [inlinePositioning],
			});
		}

		return item;
	}

	private addWarningsLink() {
		return this.addWarningLink({
			parent: this.warningElem,
			icon: <Icon icon="exclamation-triangle" size="3x" />,
			tooltip: TOOLTIP_HTML_BASE,
			classes: 'warning link-warning',
		}) as HTMLElement;
	}

	addWarning(warning: SimWarning) {
		this.warnings.push(warning);
		warning.updateOn.on(() => this.updateWarnings());
		this.updateWarnings();
	}

	private updateWarnings() {
		const activeWarnings = this.warnings
			.map(warning => warning.getContent())
			.flat()
			.filter(content => content !== '');

		const list = ((this.warningsTooltip?.props.content as Element)?.cloneNode(true) || <></>) as HTMLElement;
		if (list) list.innerHTML = '';
		list.appendChild(<>{activeWarnings?.map(warning => <li>{warning}</li>)}</>);

		this.warningsLink.parentElement?.classList?.[activeWarnings.length ? 'remove' : 'add']('hide');
		this.warningsTooltip?.setContent(list);
	}

	hideAll() {
		this.contentElem.classList.add('hide');
		this.pendingElem.classList.add('hide');
	}

	setPending() {
		this.contentElem.classList.add('hide');
		this.pendingElem.classList.remove('hide');
	}

	setContent(html: Element | HTMLElement | string) {
		if (typeof html === 'string') {
			this.contentElem.innerHTML = html;
		} else {
			this.contentElem.innerHTML = '';
			this.contentElem.appendChild(html);
		}
		this.contentElem.classList.remove('hide');
		this.pendingElem.classList.add('hide');
	}
}
