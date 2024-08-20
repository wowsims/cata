import { DeathKnight } from '../../player_classes/death_knight';
import { Hunter } from '../../player_classes/hunter';
import { Rogue } from '../../player_classes/rogue';
import { Warrior } from '../../player_classes/warrior';
import { PlayerSpecs } from '../../player_specs/index';
import { RaidSimResultsManager } from '../raid_sim_action';
import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component';

export class ToplineResults extends ResultComponent {
	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'topline-results-root';
		super(config);

		this.rootElem.classList.add('results-sim');
	}

	onSimResult(resultData: SimResultData) {
		const noManaClasses = [DeathKnight, Rogue, Warrior, Hunter];
		const players = resultData.result.getRaidIndexedPlayers(resultData.filter);
		const content = RaidSimResultsManager.makeToplineResultsContent(resultData.result, resultData.filter, {
			showOutOfMana: players.length === 1 && !!players[0].spec && !noManaClasses.some(klass => PlayerSpecs.getPlayerClass(players[0].spec!) === klass),
		});

		this.rootElem.replaceChildren(content);
	}
}
