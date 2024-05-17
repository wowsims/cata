import { DeathKnight } from '../../player_classes/death_knight.js';
import { Hunter } from '../../player_classes/hunter.js';
import { Rogue } from '../../player_classes/rogue.js';
import { Warrior } from '../../player_classes/warrior.js';
import { PlayerSpecs } from '../../player_specs/index.js';
import { RaidSimResultsManager } from '../raid_sim_action.jsx';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		const content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter);

		const noManaClasses = [DeathKnight, Rogue, Warrior, Hunter];

		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		if (players.length === 1 && !!players[0].spec && !noManaClasses.some(klass => PlayerSpecs.getPlayerClass(players[0].spec!) === klass)) {
			const player = players[0];
			const secondsOOM = player.secondsOomAvg;
			const percentOOM = secondsOOM / resultData.result.encounterMetrics.durationSeconds;
			const dangerLevel = percentOOM < 0.01 ? 'safe' : percentOOM < 0.05 ? 'warning' : 'danger';

			content.appendChild(
				<div className={`results-sim-percent-oom ${dangerLevel} damage-metrics`}>
					<span className="topline-result-avg">{secondsOOM.toFixed(1)}s</span>
				</div>,
			);
		}

		this.rootElem.innerHTML = '';
		this.rootElem.appendChild(content);
	}
}
