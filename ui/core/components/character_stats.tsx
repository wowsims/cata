import clsx from 'clsx';
import tippy from 'tippy.js';
import { ref } from 'tsx-vanilla';

import * as Mechanics from '../constants/mechanics.js';
import { Player } from '../player.js';
import { Class, PseudoStat, Stat } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id';
import { getClassStatName, masterySpellIDs, masterySpellNames, statOrder } from '../proto_utils/names.js';
import { Stats, statToPercentageOrPoints } from '../proto_utils/stats.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { Component } from './component.js';
import { NumberPicker } from './pickers/number_picker.js';

export type StatMods = { base?: Stats; gear?: Stats; talents?: Stats; buffs?: Stats; consumes?: Stats; final?: Stats; stats?: Array<Stat> };
export type StatWrites = { base: Stats; gear: Stats; talents: Stats; buffs: Stats; consumes: Stats; final: Stats; stats: Array<Stat> };

export class CharacterStats extends Component {
	readonly stats: Array<Stat>;
	readonly valueElems: Array<HTMLTableCellElement>;
	readonly meleeCritCapValueElem: HTMLTableCellElement | undefined;
	masteryElem: HTMLTableCellElement | undefined;

	private readonly player: Player<any>;
	private readonly modifyDisplayStats?: (player: Player<any>) => StatMods;
	private readonly overwriteDisplayStats?: (player: Player<any>) => StatWrites;

	constructor(
		parent: HTMLElement,
		player: Player<any>,
		stats: Array<Stat>,
		modifyDisplayStats?: (player: Player<any>) => StatMods,
		overwriteDisplayStats?: (player: Player<any>) => StatWrites,
	) {
		super(parent, 'character-stats-root');
		this.stats = statOrder.filter(stat => stats.includes(stat));
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
		this.stats.forEach(stat => {
			const statName = getClassStatName(stat, player.getClass());
			const valueRef = ref<HTMLTableCellElement>();
			const hasteTooltipRef = ref<HTMLButtonElement>();
			const row = (
				<tr className="character-stats-table-row">
					<td className="character-stats-table-label">
						{statName}
						{[Stat.StatMeleeHaste, Stat.StatSpellHaste].includes(stat) && (
							<button ref={hasteTooltipRef} className="d-inline ms-1">
								<i className="fa-regular fa-circle-question" />
							</button>
						)}
						{stat === Stat.StatMastery && (
							<>
								<br />
								{masterySpellNames.get(this.player.getSpec())}
							</>
						)}
					</td>
					<td ref={valueRef} className="character-stats-table-value">
						{this.bonusStatsLink(stat)}
					</td>
				</tr>
			);

			if (hasteTooltipRef.value)
				tippy(hasteTooltipRef.value, {
					content: (
						<>
							<p className="mb-1">
								<strong>Why is buffs showing 0?</strong>
							</p>
							<p>Haste is a multiplicate stat, this makes it hard to properly calculate ratings on the fly (for example when reforging).</p>
							<p className="mb-0">
								The <strong>total</strong> value still shows the correct haste value including all (raid) buffs.
							</p>
						</>
					),
				});

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
		const debuffStats = this.getDebuffStats();
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
			.add(statMods.final || new Stats())
			.add(debuffStats);

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

		// Apply multiplicative Haste buffs to the final displayed value
		const baseMeleeHasteMultiplier = 1 + finalStats.getStat(Stat.StatMeleeHaste) / (Mechanics.HASTE_RATING_PER_HASTE_PERCENT * 100);
		const meleeHasteBuffsMultiplier =
			this.player.getClass() == Class.ClassHunter
				? finalStats.getPseudoStat(PseudoStat.PseudoStatRangedSpeedMultiplier)
				: finalStats.getPseudoStat(PseudoStat.PseudoStatMeleeSpeedMultiplier);
		finalStats = finalStats.withStat(
			Stat.StatMeleeHaste,
			(baseMeleeHasteMultiplier * meleeHasteBuffsMultiplier - 1) * 100 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT,
		);
		const baseSpellHasteMultiplier = 1 + finalStats.getStat(Stat.StatSpellHaste) / (Mechanics.HASTE_RATING_PER_HASTE_PERCENT * 100);
		const spellHasteBuffsMultiplier = finalStats.getPseudoStat(PseudoStat.PseudoStatCastSpeedMultiplier);
		finalStats = finalStats.withStat(
			Stat.StatSpellHaste,
			(baseSpellHasteMultiplier * spellHasteBuffsMultiplier - 1) * 100 * Mechanics.HASTE_RATING_PER_HASTE_PERCENT,
		);

		const masteryPoints =
			this.player.getBaseMastery() + (playerStats.finalStats?.stats[Stat.StatMastery] || 0) / Mechanics.MASTERY_RATING_PER_MASTERY_POINT;

		this.stats.forEach((stat, idx) => {
			const bonusStatValue = bonusStats.getStat(stat);
			let contextualClass: string;
			if (bonusStatValue == 0) {
				contextualClass = 'text-white';
			} else if (bonusStatValue > 0) {
				contextualClass = 'text-success';
			} else {
				contextualClass = 'text-danger';
			}

			const statLinkElemRef = ref<HTMLButtonElement>();

			const valueElem = (
				<div className="stat-value-link-container">
					<button ref={statLinkElemRef} className={clsx('stat-value-link', contextualClass)}>
						{`${this.statDisplayString(finalStats, finalStats, stat, true)} `}
					</button>
					{stat === Stat.StatMastery && (
						<a
							href={ActionId.makeSpellUrl(masterySpellIDs.get(this.player.getSpec()) || 0)}
							className={clsx('stat-value-link-mastery', contextualClass)}
							target="_blank">
							{`${(masteryPoints * this.player.getMasteryPerPointModifier()).toFixed(2)}%`}
						</a>
					)}
				</div>
			);

			const statLinkElem = statLinkElemRef.value!;

			this.valueElems[idx].querySelector('.stat-value-link-container')?.remove();
			this.valueElems[idx].prepend(valueElem);

			const tooltipContent = (
				<div>
					<div className="character-stats-tooltip-row">
						<span>Base:</span>
						<span>{this.statDisplayString(baseStats, baseDelta, stat, true)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Gear:</span>
						<span>{this.statDisplayString(gearStats, gearDelta, stat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Talents:</span>
						<span>{this.statDisplayString(talentsStats, talentsDelta, stat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Buffs:</span>
						<span>{this.statDisplayString(buffsStats, buffsDelta, stat)}</span>
					</div>
					<div className="character-stats-tooltip-row">
						<span>Consumes:</span>
						<span>{this.statDisplayString(consumesStats, consumesDelta, stat)}</span>
					</div>
					{debuffStats.getStat(stat) != 0 && (
						<div className="character-stats-tooltip-row">
							<span>Debuffs:</span>
							<span>{this.statDisplayString(debuffStats, debuffStats, stat)}</span>
						</div>
					)}
					{bonusStatValue !== 0 && (
						<div className="character-stats-tooltip-row">
							<span>Bonus:</span>
							<span>{this.statDisplayString(bonusStats, bonusStats, stat)}</span>
						</div>
					)}
					<div className="character-stats-tooltip-row">
						<span>Total:</span>
						<span>{this.statDisplayString(finalStats, finalStats, stat, true)}</span>
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

	private statDisplayString(stats: Stats, deltaStats: Stats, stat: Stat, includeBase?: boolean): string {
		let rawValue = deltaStats.getStat(stat);

		rawValue *= 1;

		let displayStr = String(Math.round(rawValue));
		const statAsPercentageOrPoint = statToPercentageOrPoints(stat, rawValue, stats);

		if (stat == Stat.StatMeleeHit) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatSpellHit) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeCrit || stat == Stat.StatSpellCrit) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatMeleeHaste) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatSpellHaste) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatExpertise) {
			// As of 06/20, Blizzard has changed Expertise to no longer truncate at quarter percent intervals. Note that
			// in-game character sheet tooltips will still display the truncated values, but it has been tested to behave
			// continuously in reality since the patch.
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatBlock) {
			// TODO: Figure out how to display these differently for the components than the final value
			//displayStr += ` (${(rawValue / Mechanics.BLOCK_RATING_PER_BLOCK_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatDodge) {
			//displayStr += ` (${(rawValue / Mechanics.DODGE_RATING_PER_DODGE_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatParry) {
			//displayStr += ` (${(rawValue / Mechanics.PARRY_RATING_PER_PARRY_CHANCE).toFixed(2)}%)`;
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatResilience) {
			displayStr += ` (${statAsPercentageOrPoint.toFixed(2)}%)`;
		} else if (stat == Stat.StatMastery) {
			displayStr += ` (${(statAsPercentageOrPoint + (includeBase ? this.player.getBaseMastery() : 0)).toFixed(2)} Points)`;
		}

		return displayStr;
	}

	private getDebuffStats(): Stats {
		let debuffStats = new Stats();

		const debuffs = this.player.sim.raid.getDebuffs();
		if (debuffs.criticalMass || debuffs.shadowAndFlame) {
			debuffStats = debuffStats.addStat(Stat.StatSpellCrit, 5 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);
		}

		return debuffStats;
	}

	private bonusStatsLink(stat: Stat): HTMLElement {
		const statName = getClassStatName(stat, this.player.getClass());
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
					id: `character-bonus-stat-${stat}`,
					label: `Bonus ${statName}`,
					extraCssClasses: ['mb-0'],
					changedEvent: (player: Player<any>) => player.bonusStatsChangeEmitter,
					getValue: (player: Player<any>) => player.getBonusStats().getStat(stat),
					setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
						const bonusStats = player.getBonusStats().withStat(stat, newValue);
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
