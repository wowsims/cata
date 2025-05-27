import tippy from 'tippy.js';

import { Player } from '../../player';
import { PetSpec } from '../../proto/hunter';
import { HunterSpecs } from '../../proto_utils/utils';
import { TypedEvent } from '../../typed_event';
import { Component } from '../component';

export class PetSpecPicker<SpecType extends HunterSpecs> extends Component {
	private readonly player: Player<SpecType>;
	private readonly container: HTMLDivElement;

	constructor(parent: HTMLElement, player: Player<SpecType>) {
		super(parent, 'pet-spec-picker');
		this.player = player;

		// Header
		const header = document.createElement('div');
		header.className = 'talent-tree-header';
		const title = document.createElement('span');
		title.className = 'talent-tree-title';
		title.textContent = 'Pet Spec';
		header.appendChild(title);
		this.rootElem.appendChild(header);

		// Container for spec items (vertical block)
		this.container = document.createElement('div');
		this.container.className = 'talent-tree-main pet-spec-list';
		this.rootElem.appendChild(this.container);

		// Specs with corresponding icon keys
		const specs: Array<{ spec: PetSpec; label: string; iconKey: string }> = [
			{ spec: PetSpec.Ferocity, label: 'Ferocity', iconKey: 'ability_druid_kingofthejungle' },
			{ spec: PetSpec.Tenacity, label: 'Tenacity', iconKey: 'ability_druid_demoralizingroar' },
			{ spec: PetSpec.Cunning, label: 'Cunning', iconKey: 'ability_eyeoftheowl' },
		];

		specs.forEach(({ spec, label, iconKey }) => {
			// Wrap each as a talent-picker style item
			const item = document.createElement('div');
			item.className = 'talent-picker-root pet-spec-item';
			item.addEventListener('click', () => this.onClickSpec(spec));

			// Icon div with background image
			const icon = document.createElement('div');
			icon.className = 'talent-picker-icon';
			icon.style.backgroundImage = `url('https://wow.zamimg.com/images/wow/icons/large/${iconKey}.jpg')`;
			item.appendChild(icon);

			// Label text
			const lbl = document.createElement('div');
			lbl.className = 'talent-picker-label';
			lbl.textContent = label;
			item.appendChild(lbl);

			// Tooltip
			tippy(item, { content: label });

			this.container.appendChild(item);
		});

		// Listen and render selection
		player.specOptionsChangeEmitter.on(() => this.renderActive());
		this.renderActive();
	}

	private onClickSpec(newSpec: PetSpec) {
		const current = this.player.getClassOptions().petSpec;
		if (newSpec === current) return;
		const opts = this.player.getClassOptions();
		opts.petSpec = newSpec;
		this.player.setClassOptions(TypedEvent.nextEventID(), opts);
	}

	private renderActive() {
		const active = this.player.getClassOptions().petSpec;
		const order = [PetSpec.Ferocity, PetSpec.Tenacity, PetSpec.Cunning];
		Array.from(this.container.children).forEach((el, idx) => {
			const spec = order[idx];
			el.classList.toggle('selected', spec === active);
		});
	}
}
