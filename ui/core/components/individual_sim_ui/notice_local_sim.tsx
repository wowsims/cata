import { LOCAL_STORAGE_PREFIX, REPO_RELEASES_URL } from '../../constants/other';
import { isDevMode, isLocal } from '../../utils';
import { Component } from '../component';
import Toast from '../toast';

export class NoticeLocalSim extends Component {
	container: HTMLElement;
	toast: Toast | null = null;
	constructor(parent: HTMLElement) {
		super(null);
		this.container = parent;

		if (this.hasSeenNotice || isLocal() || isDevMode()) return;

		this.render();

		this.toast?.element.addEventListener(
			'hide.bs.toast',
			() => {
				this.setHasSeenNotice();
			},
			{ once: true },
		);
	}

	private get settingsKey(): string {
		return `${LOCAL_STORAGE_PREFIX}_notice-local-sim.v1`;
	}

	private get hasSeenNotice() {
		return window.localStorage.getItem(this.settingsKey);
	}

	private setHasSeenNotice() {
		window.localStorage.setItem(this.settingsKey, 'true');
	}

	render() {
		this.toast = new Toast({
			additionalClasses: ['toast-notice-local-download'],
			container: this.container,
			variant: 'info',
			title: 'Did you know?',
			autohide: false,
			body: (
				<>
					<p>You can download our local sim to speed up your simulations significantly?</p>
					<a href={REPO_RELEASES_URL} className="btn btn-outline-light" target="_blank" onclick={() => this.setHasSeenNotice()}>
						Download
					</a>
				</>
			),
		});
	}
}
