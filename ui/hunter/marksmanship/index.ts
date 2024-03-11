import { Player } from '../../core/player.js';
import { PlayerSpecs } from '../../core/player_specs';
import { Spec } from '../../core/proto/common.js';
import { Sim } from '../../core/sim.js';
import { TypedEvent } from '../../core/typed_event.js';
import { MarksmanshipHunterSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecMarksmanshipHunter>(PlayerSpecs.MarksmanshipHunter, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

new MarksmanshipHunterSimUI(document.body, player);
