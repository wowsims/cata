import { LaunchStatus, raidSimStatus, simLaunchStatuses } from '../launched_sims.js';
import { PlayerClass } from '../player_class';
import { PlayerClasses } from '../player_classes';
import { PlayerSpec } from '../player_spec';
import { Class, Spec } from '../proto/common';
import { raidSimIcon, raidSimLabel, raidSimSiteUrl, textCssClassForClass, textCssClassForSpec } from '../proto_utils/utils.js';
import { Component } from './component.js';

interface ClassOptions {
	type: 'Class';
	class: PlayerClass<Class>;
}

interface SpecOptions {
	type: 'Spec';
	spec: PlayerSpec<any>;
}

interface RaidOptions {
	type: 'Raid';
}

type SimTitleDropdownConfig = {
	noDropdown?: boolean;
};

// Dropdown menu for selecting a player.
export class SimTitleDropdown extends Component {
	private readonly dropdownMenu: HTMLElement | undefined;

	constructor(parent: HTMLElement, currentSpec: PlayerSpec<any> | null, config: SimTitleDropdownConfig = {}) {
		super(parent, 'sim-title-dropdown-root');

		const rootLinkArgs: SpecOptions | RaidOptions = currentSpec === null ? { type: 'Raid' } : { type: 'Spec', spec: currentSpec };
		const rootLink = this.buildRootSimLink(rootLinkArgs);

		if (config.noDropdown) {
			this.rootElem.innerHTML = rootLink.outerHTML;
			return;
		}

		this.rootElem.innerHTML = `
			<div class="dropdown sim-link-dropdown">
				${rootLink.outerHTML}
				<ul class="dropdown-menu"></ul>
			</div>
    	`;

		this.dropdownMenu = this.rootElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
		this.buildDropdown();

		// Prevent Bootstrap from closing the menu instead of opening class menus
		this.dropdownMenu.addEventListener('click', event => {
			const target = event.target as HTMLElement;
			const link = target.closest('a:not([href="javascript:void(0)"]');

			if (!link) {
				event.stopPropagation();
				event.preventDefault();
			}
		});
	}

	private buildDropdown() {
		if (raidSimStatus.status >= LaunchStatus.Alpha) {
			// Add the raid sim to the top of the dropdown
			const raidListItem = document.createElement('li');
			raidListItem.appendChild(this.buildRaidLink());
			this.dropdownMenu?.appendChild(raidListItem);
		}

		PlayerClasses.naturalOrder.forEach(klass => {
			const listItem = document.createElement('li');
			// Add the class to the dropdown with an additional spec dropdown
			listItem.appendChild(this.buildClassDropdown(klass));
			this.dropdownMenu?.appendChild(listItem);
		});
	}

	private buildClassDropdown(klass: PlayerClass<Class>) {
		const dropdownFragment = document.createElement('fragment');
		const dropdownMenu = document.createElement('ul');
		dropdownMenu.classList.add('dropdown-menu');

		// Generate the class link to act as a dropdown toggle for the spec dropdown
		const classLink = this.buildClassLink(klass);

		// Generate links for a class's specs
		Object.values(klass.specs).forEach(spec => {
			const listItem = document.createElement('li');
			const link = this.buildSpecLink(spec);

			listItem.appendChild(link);
			dropdownMenu.appendChild(listItem);
		});

		dropdownFragment.innerHTML = `
			<div class="dropend sim-link-dropdown">
				${classLink.outerHTML}
				${dropdownMenu.outerHTML}
			</div>
		`;

		return dropdownFragment.children[0] as HTMLElement;
	}

	private buildRootSimLink(data: SpecOptions | RaidOptions): HTMLElement {
		const iconPath = this.getSimIconPath(data);
		const textKlass = this.getContextualKlass(data);

		let label = '';
		if (data.type == 'Raid') {
			label = raidSimLabel;
		} else if (data.type == 'Spec') {
			label = data.spec.friendlyName;
		}

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="javascript:void(0)" class="sim-link ${textKlass}" role="button" data-bs-toggle="dropdown" aria-expanded="false">
				<div class="sim-link-content">
				<img src="${iconPath}" class="sim-link-icon">
				<div class="d-flex flex-column">
					<span class="sim-link-label text-white">WoWSims - Cataclysm</span>
					<span class="sim-link-title">${label}</span>
					${this.launchStatusLabel(data)}
				</div>
				</div>
			</a>
		`;

		return fragment.children[0] as HTMLElement;
	}

	private buildRaidLink(): HTMLElement {
		const textKlass = this.getContextualKlass({ type: 'Raid' });
		const iconPath = this.getSimIconPath({ type: 'Raid' });
		const label = raidSimLabel;

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="${raidSimSiteUrl}" class="sim-link ${textKlass}">
				<div class="sim-link-content">
				<img src="${iconPath}" class="sim-link-icon">
				<div class="d-flex flex-column">
					<span class="sim-link-title">${label}</span>
					${this.launchStatusLabel({ type: 'Raid' })}
				</div>
				</div>
			</a>
		`;

		return fragment.children[0] as HTMLElement;
	}

	private buildClassLink(klass: PlayerClass<Class>): HTMLElement {
		const textKlass = this.getContextualKlass({ type: 'Class', class: klass });
		const iconPath = this.getSimIconPath({ type: 'Class', class: klass });
		const label = klass.friendlyName;

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="javascript:void(0)" class="sim-link ${textKlass}" role="button" data-bs-toggle="dropdown" aria-expanded="false">
				<div class="sim-link-content">
				<img src="${iconPath}" class="sim-link-icon">
				<div class="d-flex flex-column">
					<span class="sim-link-title">${label}</span>
				</div>
				</div>
			</a>
		`;

		return fragment.children[0] as HTMLElement;
	}

	private buildSpecLink(spec: PlayerSpec<any>): HTMLElement {
		const textKlass = this.getContextualKlass({ type: 'Spec', spec: spec });
		const iconPath = this.getSimIconPath({ type: 'Spec', spec: spec });

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="${spec.simLink}" class="sim-link ${textKlass}" role="button">
				<div class="sim-link-content">
				<img src="${iconPath}" class="sim-link-icon">
				<div class="d-flex flex-column">
					<span class="sim-link-label">${spec.playerClass.friendlyName}</span>
					<span class="sim-link-title">${spec.friendlyName}</span>
					${this.launchStatusLabel({ type: 'Spec', spec: spec })}
				</div>
				</div>
			</a>
		`;

		return fragment.children[0] as HTMLElement;
	}

	private launchStatusLabel(data: SpecOptions | RaidOptions): string {
		if (
			(data.type == 'Raid' && raidSimStatus.status == LaunchStatus.Launched) ||
			(data.type == 'Spec' && simLaunchStatuses[data.spec.protoID as Spec].status == LaunchStatus.Launched)
		)
			return '';

		const status = data.type == 'Raid' ? raidSimStatus.status : simLaunchStatuses[data.spec.protoID as Spec].status;
		const phase = data.type == 'Raid' ? raidSimStatus.phase : simLaunchStatuses[data.spec.protoID as Spec].phase;

		const elem = document.createElement('span');
		elem.classList.add('launch-status-label', 'text-brand');
		elem.textContent = status == LaunchStatus.Unlaunched ? 'Not Yet Supported' : `Phase ${phase} - ${LaunchStatus[status]}`;

		return elem.outerHTML;
	}

	private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
		let iconPath = '';
		if (data.type == 'Raid') {
			iconPath = raidSimIcon;
		} else if (data.type == 'Class') {
			iconPath = data.class.getIcon('large');
		} else if (data.type == 'Spec') {
			iconPath = data.spec.getIcon('large');
		}

		return iconPath;
	}

	private getContextualKlass(data: ClassOptions | SpecOptions | RaidOptions): string {
		switch (data.type) {
			case 'Class':
				return textCssClassForClass(data.class);
			case 'Spec':
				return textCssClassForSpec(data.spec);
			default:
				return 'text-white';
		}
	}
}
