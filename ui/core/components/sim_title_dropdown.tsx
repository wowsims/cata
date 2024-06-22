import clsx from 'clsx';
import { ref } from 'tsx-vanilla';

import { LaunchStatus, raidSimStatus, simLaunchStatuses } from '../launched_sims.js';
import { PlayerClass } from '../player_class.js';
import { PlayerClasses } from '../player_classes/index.js';
import { PlayerSpec } from '../player_spec.js';
import { PlayerSpecs } from '../player_specs/index.js';
import { Class, Spec } from '../proto/common.js';
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

		const rootLinkArgs: SpecOptions | RaidOptions = !!currentSpec ? { type: 'Spec', spec: currentSpec } : { type: 'Raid' };
		const rootLink = this.buildRootSimLink(rootLinkArgs);

		if (config.noDropdown) {
			this.rootElem.replaceChildren(rootLink);
			return;
		}

		const dropdownMenuRef = ref<HTMLUListElement>();
		this.rootElem.replaceChildren(
			<div className="dropdown sim-link-dropdown">
				{rootLink}
				<ul ref={dropdownMenuRef} className="dropdown-menu"></ul>
			</div>,
		);

		this.dropdownMenu = dropdownMenuRef.value!;
		this.buildDropdown();

		// Prevent Bootstrap from closing the menu instead of opening class menus
		this.dropdownMenu.addEventListener('click', event => {
			const target = event.target as HTMLElement;
			const link = target.closest('a');

			if (!link) {
				event.stopPropagation();
				event.preventDefault();
			}
		});
	}

	private buildDropdown() {
		if (raidSimStatus.status >= LaunchStatus.Alpha) {
			this.dropdownMenu?.appendChild(<li>{this.buildRaidLink()}</li>);
		}
		PlayerClasses.naturalOrder.forEach(klass => this.dropdownMenu?.appendChild(<li>{this.buildClassDropdown(klass)}</li>));
	}

	private buildClassDropdown(klass: PlayerClass<Class>) {
		return (
			<div className="dropend sim-link-dropdown">
				{this.buildClassLink(klass)}
				<ul className="dropdown-menu">
					{Object.values(klass.specs).map(spec => (
						<li>{this.buildSpecLink(spec)}</li>
					))}
				</ul>
			</div>
		);
	}

	private buildRootSimLink(data: SpecOptions | RaidOptions) {
		return (
			<button
				className={clsx('sim-link', this.getContextualKlass(data))}
				dataset={{ bsToggle: 'dropdown', bsTrigger: 'click' }}
				attributes={{ 'aria-expanded': 'false' }}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath(data)} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label text-white">WoWSims - Cataclysm</span>
						<span className="sim-link-title">
							{data.type === 'Raid' && raidSimLabel}
							{data.type === 'Spec' && PlayerSpecs.getFullSpecName(data.spec)}
						</span>
						{this.launchStatusLabel(data)}
					</div>
				</div>
			</button>
		);
	}

	private buildRaidLink() {
		return (
			<a href={raidSimSiteUrl} className={clsx('sim-link', this.getContextualKlass({ type: 'Raid' }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Raid' })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-title">{raidSimLabel}</span>
						{this.launchStatusLabel({ type: 'Raid' })}
					</div>
				</div>
			</a>
		);
	}

	private buildClassLink(klass: PlayerClass<Class>) {
		return (
			<button
				className={clsx('sim-link', this.getContextualKlass({ type: 'Class', class: klass }))}
				dataset={{ bsToggle: 'dropdown' }}
				attributes={{ 'aria-expanded': 'false' }}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Class', class: klass })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-title">{klass.friendlyName}</span>
					</div>
				</div>
			</button>
		);
	}

	private buildSpecLink(spec: PlayerSpec<any>) {
		return (
			<a href={spec.simLink} className={clsx('sim-link', this.getContextualKlass({ type: 'Spec', spec: spec }))}>
				<div className="sim-link-content">
					<img src={this.getSimIconPath({ type: 'Spec', spec: spec })} className="sim-link-icon" />
					<div className="d-flex flex-column">
						<span className="sim-link-label">{PlayerSpecs.getPlayerClass(spec).friendlyName}</span>
						<span className="sim-link-title">{spec.friendlyName}</span>
						{this.launchStatusLabel({ type: 'Spec', spec: spec })}
					</div>
				</div>
			</a>
		);
	}

	private launchStatusLabel(data: SpecOptions | RaidOptions) {
		if (
			(data.type === 'Raid' && raidSimStatus.status === LaunchStatus.Launched)
		)
			return null;

		const status = data.type === 'Raid' ? raidSimStatus.status : simLaunchStatuses[data.spec.specID as Spec].status;
		const phase = data.type === 'Raid' ? raidSimStatus.phase : simLaunchStatuses[data.spec.specID as Spec].phase;

		return (
			<span className="launch-status-label text-brand">
				{status === LaunchStatus.Unlaunched ? (
					<>Not Yet Supported</>
				) : (
					<>
						Phase {phase} - {LaunchStatus[status]}
					</>
				)}
			</span>
		);
	}

	private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
		let iconPath = '';
		switch (data.type) {
			case 'Raid':
				iconPath = raidSimIcon;
				break;
			case 'Class':
				iconPath = data.class.getIcon('large');
				break;
			case 'Spec':
				iconPath = data.spec.getIcon('large');
				break;
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
