import { IndividualSimUI } from '../../individual_sim_ui';
import { Player } from '../../player';
import { SavedEPWeights } from '../../proto/ui';
import { Stats } from '../../proto_utils/stats';
import { TypedEvent } from '../../typed_event';
import { SavedDataManager, SavedDataManagerConfig } from '../saved_data_manager';

export const renderSavedEPWeights = (
	container: HTMLElement | null,
	simUI: IndividualSimUI<any>,
	options?: Partial<SavedDataManagerConfig<Player<any>, SavedEPWeights>>,
) => {
	const savedEPWeightsManager = new SavedDataManager<Player<any>, SavedEPWeights>(container, simUI.player, {
		label: 'EP Weights',
		header: { title: 'Saved EP weights' },
		storageKey: simUI.getSavedEPWeightsStorageKey(),
		getData: player =>
			SavedEPWeights.create({
				epWeights: player.getEpWeights().toProto(),
			}),
		setData: (eventID, player, newEPWeights) => {
			TypedEvent.freezeAllAndDo(() => {
				player.setEpWeights(eventID, Stats.fromProto(newEPWeights.epWeights));
			});
		},
		changeEmitters: [simUI.player.epWeightsChangeEmitter],
		equals: (a, b) => SavedEPWeights.equals(a, b),
		toJson: a => SavedEPWeights.toJson(a),
		fromJson: obj => SavedEPWeights.fromJson(obj),
		...options,
	});

	simUI.sim.waitForInit().then(() => {
		savedEPWeightsManager.loadUserData();
		simUI.individualConfig.presets.epWeights.forEach(({ name, epWeights, enableWhen, onLoad }) => {
			savedEPWeightsManager.addSavedData({
				name: name,
				isPreset: true,
				data: SavedEPWeights.create({
					epWeights: epWeights.toProto(),
				}),
				enableWhen,
				onLoad,
			});
		});
	});

	return savedEPWeightsManager;
};
