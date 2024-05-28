import { Player } from '../../core/player';
import { PlayerSpecs } from '../../core/player_specs';
import { Spec } from '../../core/proto/common';
import { Sim } from '../../core/sim';
import { TypedEvent } from '../../core/typed_event';
import { GuardianDruidSimUI } from './sim';

const sim = new Sim();
const player = new Player<Spec.SpecGuardianDruid>(PlayerSpecs.GuardianDruid, sim);
player.enableHealing();

sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

const simUI = new GuardianDruidSimUI(document.body, player);
