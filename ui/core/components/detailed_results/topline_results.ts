import { RaidSimResultsManager } from '../../components/raid_sim_action.js';
import { DeathKnight } from '../../player_classes/death_knight';
import { Hunter } from '../../player_classes/hunter';
import { Rogue } from '../../player_classes/rogue';
import { Warrior } from '../../player_classes/warrior';
import { PlayerSpec } from '../../player_spec';
import { PlayerSpecs } from '../../player_specs';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		let content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter);

		const noManaClasses = [DeathKnight, Rogue, Warrior, Hunter];

		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length == 1 && !!players[0].spec && !noManaClasses.some(klass => PlayerSpecs.getPlayerClass(players[0].spec as PlayerSpec<any>) == klass)) {
			const player = players[0];
			const secondsOOM = player.secondsOomAvg;
			const percentOOM = secondsOOM / resultData.result.encounterMetrics.durationSeconds;
			const dangerLevel = percentOOM < 0.01 ? 'safe' : percentOOM < 0.05 ? 'warning' : 'danger';

			content += `
				<div class="results-sim-percent-oom ${dangerLevel} damage-metrics">
					<span class="topline-result-avg">${secondsOOM.toFixed(1)}s</span>
				</div>
			`;
		}

		this.rootElem.innerHTML = content;
	}
}
