import { Dropdown } from 'bootstrap';
import { hideAll as hideAllTippyInstances } from 'tippy.js';

import { isDescendant } from './utils';

Dropdown.Default.offset = [0, -1];
//Dropdown.Default.display = "static";

const body = document.querySelector('body') as HTMLElement;

function hasTouch() {
	return 'ontouchstart' in window || navigator.maxTouchPoints > 0;
}

function hasHover() {
	return window.matchMedia('(any-hover: hover)').matches;
}

// Disable 'mouseover' to avoid needed to double click on mobile
// Leaving 'mouseleave', however still allows dropdown to close when clicking new box
if (!hasTouch() || hasHover()) {
	// Custom dropdown event handlers for mouseover dropdowns
	body.addEventListener(
		'mouseover',
		event => {
			const target = event.target as HTMLElement;
			const toggle = target.closest('[data-bs-toggle=dropdown]:not([data-bs-trigger=click])');
			if (toggle && !toggle.classList.contains('open-on-click')) {
				const dropdown = Dropdown.getOrCreateInstance(toggle);
				dropdown.show();
			}
		},
		true,
	);
}

body.addEventListener(
	'mouseleave',
	event => {
		const e = event as MouseEvent;
		const target = event.target as HTMLElement;
		const toggle = target.closest('[data-bs-toggle=dropdown]:not([data-bs-trigger=click])') as HTMLElement | null;
		// Hide dropdowns when hovering off of the toggle, so long as the new target is not part of the dropdown as well
		if (toggle) {
			const dropdown = Dropdown.getOrCreateInstance(toggle);
			const dropdownMenu = toggle.nextElementSibling as HTMLElement;
			const relatedTarget = e.relatedTarget as HTMLElement;
			if (relatedTarget == null || (!isDescendant(relatedTarget, dropdownMenu) && !isDescendant(relatedTarget, toggle))) dropdown.hide();
		}

		const dropdownMenu = target.closest('.dropdown-menu') as HTMLElement;
		// Hide dropdowns when hovering off of the menu, so long as the new target is not part of the dropdown as well
		if (dropdownMenu) {
			const toggle = dropdownMenu.previousElementSibling as HTMLElement;
			const dropdown = Dropdown.getOrCreateInstance(toggle);
			const relatedTarget = e.relatedTarget as HTMLElement;
			if (relatedTarget == null || (!isDescendant(relatedTarget, dropdownMenu) && e.relatedTarget != toggle)) dropdown.hide();
		}
	},
	true,
);

document.addEventListener('keydown', event => {
	if (event.key == 'Escape') {
		hideAllTippyInstances();
	}
});
