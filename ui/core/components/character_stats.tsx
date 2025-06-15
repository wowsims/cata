import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import * as Mechanics from '../constants/mechanics.js';
import { Player } from '../player.js';
import { Spec, Stat } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id';
import { getStatName, masterySpellIDs, masterySpellNames } from '../proto_utils/names.js';
import { Stats, UnitStat } from '../proto_utils/stats.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { Component } from './component.js';
import { NumberPicker } from './pickers/number_picker.js';

export type StatMods = { base?: Stats; gear?: Stats; talents?: Stats; buffs?: Stats; consumes?: Stats; final?: Stats; stats?: Array<Stat> };
export type StatWrites = { base: Stats; gear: Stats; talents: Stats; buffs: Stats; consumes: Stats; final: Stats; stats: Array<Stat> };

export class CharacterStats extends Component {
	readonly stats: Array<UnitStat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly meleeCritCapValueElem: HTMLTableCellElement | undefined;
	masteryElem: HTMLTableCellElement | undefined;

	private readonly player: Player<any>;
	private readonly modifyDisplayStats?: (player: Player<any>) => StatMods;
	private readonly overwriteDisplayStats?: (player: Player<any>) => StatWrites;

	constructor(
		parent: HTMLElement,
		player: Player<any>,
		statList: Array<UnitStat>,
		modifyDisplayStats?: (player: Player<any>) => StatMods,
		overwriteDisplayStats?: (player: Player<any>) => StatWrites,
	) {
		super(parent, 'character-stats-root');
		this.stats = statList;
		this.player = player;
		this.modifyDisplayStats = modifyDisplayStats;
		this.overwriteDisplayStats = overwriteDisplayStats;

		const label = document.createElement('label');
		label.classList.add('character-stats-label');
		label.textContent = 'Stats';
		this.rootElem.appendChild(label);

		const table = document.createElement('table');
		table.classList.add('character-stats-table');
		this.rootElem.appendChild(table);

		this.valueElems = [];
		this.stats.forEach(unitStat => {
			const statName = unitStat.getShortName(player.getClass());
			const valueRef = ref<HTMLTableCellElement>();
			const row = (
				<tr className="character-stats-table-row">
					<td className="character-stats-table-label">
						{statName}
						{unitStat.equalsStat(Stat.StatMasteryRating) && (
							<>
								<br />
								{masterySpellNames.get(this.player.getSpec())}
							</>
						)}
					</td>
					<td ref={valueRef} className="character-stats-table-value">
						{unitStat.hasRootStat() && this.bonusStatsLink(unitStat)}
					</td>
				</tr>
			);

			table.appendChild(row);
			this.valueElems.push(valueRef.value!);
		});

		if (this.shouldShowMeleeCritCap(player)) {
			const valueRef = ref<HTMLTableCellElement>();
			const row = (
				<tr className="character-stats-table-row">
					<td className="character-stats-table-label">Melee Crit Cap</td>
					<td ref={valueRef} className="character-stats-table-value"></td>
				</tr>
			);

			table.appendChild(row);
			this.meleeCritCapValueElem = valueRef.value!;
		} else {
			this.meleeCritCapValueElem = undefined;
		}

		this.updateStats(player);
		TypedEvent.onAny([player.currentStatsEmitter, player.sim.changeEmitter, player.talentsChangeEmitter]).on(() => {
			this.updateStats(player);
		});
	}

	private updateStats(player: Player<any>) {
		const playerStats = player.getCurrentStats();
		const statMods = this.modifyDisplayStats ? this.modifyDisplayStats(this.player) : {};

		const baseStats = Stats.fromProto(playerStats.baseStats);
		const gearStats = Stats.fromProto(playerStats.gearStats);
		const talentsStats = Stats.fromProto(playerStats.talentsStats);
		const buffsStats = Stats.fromProto(playerStats.buffsStats);
		const consumesStats = Stats.fromProto(playerStats.consumesStats);
		const bonusStats = player.getBonusStats();

		let baseDelta = baseStats.add(statMods.base || new Stats());
		let gearDelta = gearStats
			.subtract(baseStats)
			.subtract(bonusStats)
			.add(statMods.gear || new Stats());
		let talentsDelta = talentsStats.subtract(gearStats).add(statMods.talents || new Stats());
		let buffsDelta = buffsStats.subtract(talentsStats).add(statMods.buffs || new Stats());
		let consumesDelta = consumesStats.subtract(buffsStats).add(statMods.consumes || new Stats());

		let finalStats = Stats.fromProto(playerStats.finalStats)
			.add(statMods.base || new Stats())
			.add(statMods.gear || new Stats())
			.add(statMods.talents || new Stats())
			.add(statMods.buffs || new Stats())
			.add(statMods.consumes || new Stats())
			.add(statMods.final || new Stats());

		if (this.overwriteDisplayStats) {
			const statOverwrites = this.overwriteDisplayStats(this.player);
			if (statOverwrites.stats) {
				statOverwrites.stats.forEach((stat, _) => {
					baseDelta = baseDelta.withStat(stat, statOverwrites.base.getStat(stat));
					gearDelta = gearDelta.withStat(stat, statOverwrites.gear.getStat(stat));
					talentsDelta = talentsDelta.withStat(stat, statOverwrites.talents.getStat(stat));
					buffsDelta = buffsDelta.withStat(stat, statOverwrites.buffs.getStat(stat));
					consumesDelta = consumesDelta.withStat(stat, statOverwrites.consumes.getStat(stat));
					finalStats = finalStats.withStat(stat, statOverwrites.final.getStat(stat));
				});
			}
		}

		const masteryPoints =
			this.player.getBaseMastery() + (playerStats.finalStats?.stats[Stat.StatMasteryRating] || 0) / Mechanics.MASTERY_RATING_PER_MASTERY_POINT;

		this.stats.forEach((unitStat, idx) => {
			const bonusStatValue = unitStat.hasRootStat() ? bonusStats.getStat(unitStat.getRootStat()) : 0;
			let contextualClass: string;
			if (bonusStatValue == 0) {
				contextualClass = 'text-white';
			} else if (bonusStatValue > 0) {
				contextualClass = 'text-success';
			} else {
				contextualClass = 'text-danger';
			}

			const statLinkElemRef = ref<HTMLButtonElement>();

			// Custom "HACK" for Warlock/Protection Warrior..
			// they have two different mastery scalings
			// And a different base mastery value..
			let modifier = [this.player.getMasteryPerPointModifier()];
			let customBonus = [0];
			switch (player.getSpec()) {
				case Spec.SpecDestructionWarlock:
					customBonus = [1, 0];
					modifier = [1, ...modifier];
					break;
				case Spec.SpecDemonologyWarlock:
					customBonus = [0, 0];
					modifier = [1, ...modifier];
					break;
				case Spec.SpecProtectionWarrior:
					customBonus = [0, 0];
					modifier = [0.5, ...modifier];
					break;
				case Spec.SpecWindwalkerMonk:
					customBonus = [3.5, 0];
					break;
			}

			const valueElem = (
				<div className="stat-value-link-container">
					<button ref={statLinkElemRef} className={clsx('stat-value-link', contextualClass)}>
						{`${this.statDisplayString(finalStats, unitStat, true)} `}
					</button>
					{unitStat.equalsStat(Stat.StatMasteryRating) &&
						modifier.map((modifier, index) => (
							<a
								href={ActionId.makeSpellUrl(masterySpellIDs.get(this.player.getSpec()) || 0)}
								className={clsx('stat-value-link-mastery', contextualClass)}
								target="_blank">
								{`${(masteryPoints * modifier + customBonus[index]).toFixed(2)}%`}
							</a>
						))}
				</div>
			);

			const statLinkElem = statLinkElemRef.value!;

			this.valueElems[idx].querySelector('.stat-value-link-container')?.remove();
			this.valueElems[idx].prepend(valueElem);

			const tooltipContent = (
				<div>
					<div className="character-stats-tooltip-row">
						<span>Base:</span>
						<span>{this.statDisplayString(baseDelta, unitStat, true)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Gear:</span>
						<span>{this.statDisplayString(gearDelta, unitStat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Talents:</span>
						<span>{this.statDisplayString(talentsDelta, unitStat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Buffs:</span>
						<span>{this.statDisplayString(buffsDelta, unitStat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Consumes:</span>
						<span>{this.statDisplayString(consumesDelta, unitStat)}</span>
					</div>
					{bonusStatValue !== 0 && (
						<div className="character-stats-tooltip-row">
							<span>Bonus:</span>
							<span>{this.statDisplayString(bonusStats, unitStat)}</span>
						</div>
					)}
					<div className="character-stats-tooltip-row">
						<span>Total:</span>
						<span>{this.statDisplayString(finalStats, unitStat, true)}</span>
					</div>
				</div>
			);

			tippy(statLinkElem, {
				content: tooltipContent,
			});
		});

		if (this.meleeCritCapValueElem) {
			const meleeCritCapInfo = player.getMeleeCritCapInfo();

			const valueElem = <button className="stat-value-link">{this.meleeCritCapDisplayString(player, finalStats)} </button>;

			const capDelta = meleeCritCapInfo.playerCritCapDelta;
			if (capDelta == 0) {
				valueElem.classList.add('text-white');
			} else if (capDelta > 0) {
				valueElem.classList.add('text-danger');
			} else if (capDelta < 0) {
				valueElem.classList.add('text-success');
			}

			this.meleeCritCapValueElem.querySelector('.stat-value-link')?.remove();
			this.meleeCritCapValueElem.prepend(valueElem);

			const tooltipContent = (
				<div>
					<div className="character-stats-tooltip-row">
						<span>Glancing:</span>
						<span>{`${meleeCritCapInfo.glancing.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Suppression:</span>
						<span>{`${meleeCritCapInfo.suppression.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>To Hit Cap:</span>
						<span>{`${meleeCritCapInfo.remainingMeleeHitCap.toFixed(2)}%`}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>To Exp Cap:</span>
						<span>{`${meleeCritCapInfo.remainingExpertiseCap.toFixed(2)}%`}</span>
					</div>
					{meleeCritCapInfo.specSpecificOffset != 0 && (
						<div className="character-stats-tooltip-row">
							<span>Spec Offsets:</span>
							<span>{`${meleeCritCapInfo.specSpecificOffset.toFixed(2)}%`}</span>
						</div>
					)}
					<div className="character-stats-tooltip-row">
						<span>Final Crit Cap:</span>
						<span>{`${meleeCritCapInfo.baseCritCap.toFixed(2)}%`}</span>
					</div>
					<hr />
					<div className="character-stats-tooltip-row">
						<span>Can Raise By:</span>
						<span>{`${(meleeCritCapInfo.remainingExpertiseCap + meleeCritCapInfo.remainingMeleeHitCap).toFixed(2)}%`}</span>
					</div>
				</div>
			);

			tippy(valueElem, {
				content: tooltipContent,
			});
		}
	}

	private statDisplayString(deltaStats: Stats, unitStat: UnitStat, includeBase?: boolean): string {
		const rootRatingValue = unitStat.hasRootStat() ? deltaStats.getStat(unitStat.getRootStat()) : null;
		let derivedPercentOrPointsValue = unitStat.convertDefaultUnitsToPercent(deltaStats.getUnitStat(unitStat));

		if (unitStat.equalsStat(Stat.StatMasteryRating) && includeBase) {
			derivedPercentOrPointsValue = derivedPercentOrPointsValue! + this.player.getBaseMastery();
		}

		const hideRootRating = rootRatingValue === null || (rootRatingValue === 0 && derivedPercentOrPointsValue !== null);
		const rootRatingString = hideRootRating ? '' : String(Math.round(rootRatingValue));
		const percentOrPointsSuffix = unitStat.equalsStat(Stat.StatMasteryRating) ? ' Points' : '%';
		const percentOrPointsString = derivedPercentOrPointsValue === null ? '' : `${derivedPercentOrPointsValue.toFixed(2)}` + percentOrPointsSuffix;
		const wrappedPercentOrPointsString = hideRootRating || derivedPercentOrPointsValue === null ? percentOrPointsString : ` (${percentOrPointsString})`;
		return rootRatingString + wrappedPercentOrPointsString;
	}

	private bonusStatsLink(unitStat: UnitStat): HTMLElement {
		const rootStat = unitStat.getRootStat();
		const statName = getStatName(rootStat);
		const linkRef = ref<HTMLButtonElement>();
		const iconRef = ref<HTMLDivElement>();

		const link = (
			<button ref={linkRef} className="add-bonus-stats text-white ms-2" dataset={{ bsToggle: 'popover' }}>
				<i ref={iconRef} className="fas fa-plus-minus"></i>
			</button>
		);

		tippy(iconRef.value!, { content: `Bonus ${statName}` });
		tippy(linkRef.value!, {
			interactive: true,
			trigger: 'click',
			theme: 'bonus-stats-popover',
			placement: 'right',
			onShow: instance => {
				const picker = new NumberPicker(null, this.player, {
					id: `character-bonus-stat-${rootStat}`,
					label: `Bonus ${statName}`,
					extraCssClasses: ['mb-0'],
					changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
					getValue: (player: Player<any>) => player.getBonusStats().getStat(rootStat),
					setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
						const bonusStats = player.getBonusStats().withStat(rootStat, newValue);
						player.setBonusStats(eventID, bonusStats);
						instance?.hide();
					},
				});
				instance.setContent(picker.rootElem);
			},
		});

		return link as HTMLElement;
	}

	private shouldShowMeleeCritCap(player: Player<any>): boolean {
		return player.getPlayerSpec().isMeleeDpsSpec;
	}

	private meleeCritCapDisplayString(player: Player<any>, _finalStats: Stats): string {
		const playerCritCapDelta = player.getMeleeCritCap();

		if (playerCritCapDelta === 0.0) {
			return 'Exact';
		}

		const prefix = playerCritCapDelta > 0 ? 'Over by ' : 'Under by ';
		return `${prefix} ${Math.abs(playerCritCapDelta).toFixed(2)}%`;
	}
}
