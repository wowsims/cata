import { Icon, Link } from '@wowsims/ui';
import tippy from 'tippy.js';

import { Component } from './component';

export class SocialLinks extends Component {
	static buildDiscordLink() {
		const anchor = (
			<Link
				variant="alt"
				href="https://discord.gg/p3DgvmnDCS"
				target="_blank"
				className="discord-link"
				dataset={{ tippyContent: 'Join us on Discord' }}
				iconLeft={<Icon icon="discord" type="fab" size="lg" />}
			/>
		);
		tippy(anchor);
		return anchor;
	}

	static buildGitHubLink() {
		const anchor = (
			<Link
				variant="alt"
				href="https://github.com/wowsims/cata"
				target="_blank"
				className="github-link"
				dataset={{ tippyContent: 'Contribute on GitHub' }}
				iconLeft={<Icon icon="github" type="fab" size="lg" />}
			/>
		);
		tippy(anchor);
		return anchor;
	}

	static buildPatreonLink() {
		const anchor = (
			<Link
				variant="alt"
				href="https://patreon.com/wowsims"
				target="_blank"
				className="patreon-link"
				dataset={{ tippyContent: 'Support us on Patreon' }}
				iconLeft={<Icon icon="patreon" type="fab" size="lg" />}>
				Patreon
			</Link>
		);
		tippy(anchor);
		return anchor;
	}
}
