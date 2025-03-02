/// <reference types="vite/client" />
import './shared/bootstrap_overrides';

import * as Popper from '@popperjs/core';
import { Dropdown, Modal, Tab } from 'bootstrap';
import { Chart, registerables } from 'chart.js';
import tippy from 'tippy.js';

declare global {
	interface Window {
		Popper: any;
		bootstrap: any;
	}
}

Chart.register(...registerables);
Chart.defaults.color = 'white';

tippy.setDefaultProps({ arrow: false, allowHTML: true });
window.Popper = Popper;
window.bootstrap = { Dropdown, Modal, Tab };

// Force scroll to top when refreshing
if (history.scrollRestoration) {
	history.scrollRestoration = 'manual';
} else {
	window.onbeforeunload = function () {
		window.scrollTo(0, 0);
	};
}

function docReady(fn: any) {
	// see if DOM is already available
	if (document.readyState === 'complete' || document.readyState === 'interactive') {
		// call on next available tick
		setTimeout(fn, 1);
	} else {
		document.addEventListener('DOMContentLoaded', fn);
	}
}

docReady(function () {
	document.body.classList.add('ready');
});
