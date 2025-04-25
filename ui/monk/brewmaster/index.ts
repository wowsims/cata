import { Player } from '../../core/player';
import { PlayerSpecs } from '../../core/player_specs';
import { Spec } from '../../core/proto/common';
import { Sim } from '../../core/sim';
import { TypedEvent } from '../../core/typed_event';
import { BrewmasterMonkSimUI } from './sim';

const sim = new Sim();
const player = new Player<Spec.SpecBrewmasterMonk>(PlayerSpecs.BrewmasterMonk, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

new BrewmasterMonkSimUI(document.body, player);