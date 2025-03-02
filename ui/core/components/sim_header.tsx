import clsx from 'clsx';
import tippy, { ReferenceElement as TippyReferenceElement } from 'tippy.js';
import { ref } from 'tsx-vanilla';

import { REPO_CHOOSE_NEW_ISSUE_URL, REPO_RELEASES_URL } from '../constants/other';
import { SimUI } from '../sim_ui';
import { isLocal, noop } from '../utils';
import { Component } from './component';
import { Exporter } from './exporter';
import { SimToolbarItem } from './header/sim_toolbar_item';
import { Importer } from './importer';
import { SettingsMenu } from './settings_menu';
import { SimTab } from './sim_tab';
import { SocialLinks } from './social_links';

interface ToolbarLinkArgs {
	parent: HTMLElement;
	href?: string;
	text?: string;
	icon?: string;
	tooltip?: string | HTMLElement;
	classes?: string;
	onclick?: () => void;
}

export class SimHeader extends Component {
	private simUI: SimUI;

	private simTabsContainer: HTMLElement;
	private simToolbar: HTMLElement;
	private knownIssuesLink: TippyReferenceElement<HTMLElement>;
	private knownIssuesContent: HTMLUListElement;

	constructor(parentElem: HTMLElement, simUI: SimUI) {
		super(parentElem, 'sim-header');
		this.simUI = simUI;
		this.simTabsContainer = this.rootElem.querySelector<HTMLElement>('.sim-tabs')!;
		this.simToolbar = this.rootElem.querySelector<HTMLElement>('.sim-toolbar')!;

		this.knownIssuesContent = (<ul className="text-start ps-3 mb-0"></ul>) as HTMLUListElement;
		this.knownIssuesLink = this.addKnownIssuesLink();
		this.addBugReportLink();
		this.addDownloadBinaryLink();
		this.addSimOptionsLink();
		this.addSocialLinks();

		// Allow styling the sticky header
		new IntersectionObserver(([e]) => e.target.classList.toggle('stuck', e.intersectionRatio < 1), { threshold: [1] }).observe(this.rootElem);
	}

	activateTab(className: string) {
		(this.simTabsContainer.getElementsByClassName(className)[0] as HTMLElement).click();
	}

	addTab(title: string, contentId: string) {
		const isFirstTab = this.simTabsContainer.children.length == 0;

		this.simTabsContainer.appendChild(
			<li
				className={`${contentId} nav-item`}
				attributes={{
					role: 'presentation',
					// @ts-expect-error
					'aria-controls': contentId,
				}}>
				<button
					className={`nav-link ${isFirstTab && 'active'}`}
					type="button"
					dataset={{
						bsToggle: 'tab',
						bsTarget: `#${contentId}`,
					}}
					attributes={{
						role: 'tab',
						'aria-selected': isFirstTab,
					}}>
					{title}
				</button>
			</li>,
		);
	}

	addSimTabLink(tab: SimTab) {
		const isFirstTab = this.simTabsContainer.children.length == 0;

		tab.navLink.setAttribute('aria-selected', isFirstTab.toString());

		if (isFirstTab) tab.navLink.classList.add('active', 'show');

		this.simTabsContainer.appendChild(tab.navItem);
	}

	addImportLink(label: string, importer: Importer, hideInRaidSim?: boolean, isUnsupported = false) {
		this.addImportExportLink('.import-dropdown', label, importer, hideInRaidSim, isUnsupported);
	}
	addExportLink(label: string, exporter: Exporter, hideInRaidSim?: boolean, isUnsupported = false) {
		this.addImportExportLink('.export-dropdown', label, exporter, hideInRaidSim, isUnsupported);
	}
	private addImportExportLink(cssClass: string, label: string, importerExporter: Importer | Exporter, _hideInRaidSim?: boolean, isUnsupported?: boolean) {
		const dropdownElem = this.rootElem.querySelector<HTMLElement>(cssClass)!;
		const menuElem = dropdownElem.querySelector<HTMLElement>('.dropdown-menu')!;
		const buttonRef = ref<HTMLButtonElement>();

		menuElem.appendChild(
			<li>
				<button ref={buttonRef} className={clsx('dropdown-item', isUnsupported && 'disabled')}>
					{label}
				</button>
			</li>,
		);
		if (buttonRef.value) {
			if (isUnsupported) {
				tippy(buttonRef.value, { content: 'Currently unsupported' });
				return;
			}
			buttonRef.value.addEventListener('click', () => importerExporter.open());
		}
	}

	private addToolbarLink({ parent, tooltip, classes, onclick, text, ...itemArgs }: ToolbarLinkArgs): HTMLElement {
		const itemRef = ref<HTMLAnchorElement>();
		parent.appendChild(
			<SimToolbarItem linkRef={itemRef} buttonClassName={classes} {...itemArgs}>
				{text}
			</SimToolbarItem>,
		);

		if (onclick) itemRef.value!.addEventListener('click', onclick);
		if (tooltip)
			tippy(itemRef.value!, {
				content: tooltip,
				placement: 'bottom',
			});
		return itemRef.value!;
	}

	private addKnownIssuesLink() {
		return this.addToolbarLink({
			parent: this.simToolbar,
			text: 'Known Issues',
			tooltip: this.knownIssuesContent,
			classes: 'known-issues link-danger hide',
		});
	}

	addKnownIssue(issue: string) {
		const listItem = (<li></li>) as HTMLLIElement;
		// Using innerHTML here because the issue text can contain stringified HTML
		listItem.innerHTML = issue;
		this.knownIssuesContent.appendChild(listItem);

		this.knownIssuesLink.classList.remove('hide');
		this.knownIssuesLink._tippy?.setContent(this.knownIssuesContent);
	}

	private addBugReportLink() {
		this.addToolbarLink({
			href: REPO_CHOOSE_NEW_ISSUE_URL,
			parent: this.simToolbar,
			icon: 'fas fa-bug fa-lg',
			tooltip: 'Report a bug or<br>Request a feature',
		});
	}

	private addDownloadBinaryLink() {
		const href = REPO_RELEASES_URL;
		const icon = 'fas fa-gauge-high fa-lg';
		const parent = this.simToolbar;

		if (isLocal()) {
			fetch('/version')
				.then(resp => {
					resp.json()
						.then(versionInfo => {
							if (versionInfo.outdated == 2) {
								this.addToolbarLink({
									href: href,
									parent: parent,
									icon: icon,
									tooltip: 'Newer version of simulator available for download',
									classes: 'downbin link-danger',
								});
							}
						})
						.catch(_error => {
							console.warn('No version info found!');
						});
				})
				.catch(noop);
		} else {
			this.addToolbarLink({
				href: href,
				parent: parent,
				icon: icon,
				tooltip: 'Download simulator for faster simulating',
				classes: 'downbin',
			});
		}
	}

	private addSimOptionsLink() {
		const settingsMenu = new SettingsMenu(this.simUI.rootElem, this.simUI);
		this.addToolbarLink({
			parent: this.simToolbar,
			icon: 'fas fa-cog fa-lg',
			tooltip: 'Show Sim Options',
			classes: 'sim-options',
			onclick: () => settingsMenu.open(),
		});
	}

	private addSocialLinks() {
		const container = (<div className="sim-toolbar-socials" />) as HTMLElement;
		this.simToolbar.appendChild(container);

		this.addDiscordLink(container);
		this.addGitHubLink(container);
		this.addPatreonLink(container);
	}

	private addDiscordLink(container: HTMLElement) {
		container.appendChild(<SimToolbarItem>{SocialLinks.buildDiscordLink()}</SimToolbarItem>);
	}

	private addGitHubLink(container: HTMLElement) {
		container.appendChild(<SimToolbarItem>{SocialLinks.buildGitHubLink()}</SimToolbarItem>);
	}

	private addPatreonLink(container: HTMLElement) {
		container.appendChild(<SimToolbarItem>{SocialLinks.buildPatreonLink()}</SimToolbarItem>);
	}

	protected customRootElement(): HTMLElement {
		return (
			<header className="sim-header">
				<div className="sim-header-container">
					<ul className="sim-tabs nav nav-tabs" attributes={{ role: 'tablist' }}></ul>
					<div className="import-export nav within-raid-sim-hide">
						<div className="dropdown sim-dropdown-menu import-dropdown">
							<button className="import-link" attributes={{ 'aria-expanded': 'false' }} dataset={{ bsToggle: 'dropdown', bsDisplay: 'dynamic' }}>
								<i className="fa fa-download"></i> Import
							</button>
							<ul className="dropdown-menu"></ul>
						</div>
						<div className="dropdown sim-dropdown-menu export-dropdown">
							<button className="export-link" attributes={{ 'aria-expanded': 'false' }} dataset={{ bsToggle: 'dropdown', bsDisplay: 'dynamic' }}>
								<i className="fa fa-right-from-bracket"></i> Export
							</button>
							<ul className="dropdown-menu"></ul>
						</div>
					</div>
					<div className="sim-toolbar nav"></div>
				</div>
			</header>
		) as HTMLElement;
	}
}
