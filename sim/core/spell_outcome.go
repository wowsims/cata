package core

import (
	"github.com/wowsims/mop/sim/core/stats"
)

// This function should do 3 things:
//  1. Set the Outcome of the hit effect.
//  2. Update spell outcome metrics.
//  3. Modify the damage if necessary.
type OutcomeApplier func(sim *Simulation, result *SpellResult, attackTable *AttackTable)

func (spell *Spell) OutcomeAlwaysHit(sim *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
	spell.SpellMetrics[result.Target.UnitIndex].Hits++
}

// Hit without Hits++ counter
func (spell *Spell) OutcomeAlwaysHitNoHitCounter(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
}

func (spell *Spell) OutcomeAlwaysMiss(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeMiss
	result.Damage = 0
	spell.SpellMetrics[result.Target.UnitIndex].Misses++
}

func (dot *Dot) OutcomeTick(_ *Simulation, result *SpellResult, _ *AttackTable) {
	result.Outcome = OutcomeHit
	dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
}
func (dot *Dot) OutcomeTickPhysicalHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.PhysicalHitCheck(sim, attackTable) {
		if dot.Spell.PhysicalCritCheck(sim, attackTable) {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritDamageMultiplier()
			dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		} else {
			result.Outcome = OutcomeHit
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		}
	} else {
		result.Outcome = OutcomeMiss
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}
func (dot *Dot) OutcomeTickPhysicalCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.PhysicalCritCheck(sim, attackTable) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritDamageMultiplier()
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
	}
}

func (dot *Dot) OutcomeTickMagicCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.MagicCritCheck(sim, result.Target) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritDamageMultiplier()
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
	}
}

func (dot *Dot) OutcomeTickMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.MagicHitCheck(sim, attackTable) {
		if dot.Spell.MagicCritCheck(sim, result.Target) {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritDamageMultiplier()
			dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		} else {
			result.Outcome = OutcomeHit
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (dot *Dot) OutcomeTickHealingCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.HealingCritCheck(sim) {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritDamageMultiplier()
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
	}
}

func (dot *Dot) OutcomeSnapshotCrit(sim *Simulation, result *SpellResult, _ *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritDamageMultiplier()
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
	} else {
		result.Outcome = OutcomeHit
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
	}
}

func (dot *Dot) OutcomeMagicHitAndSnapshotCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if dot.Spell.MagicHitCheck(sim, attackTable) {
		if sim.RandomFloat("Snapshot Crit Roll") < dot.SnapshotCritChance {
			result.Outcome = OutcomeCrit
			result.Damage *= dot.Spell.CritDamageMultiplier()
			dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		} else {
			result.Outcome = OutcomeHit
			dot.Spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		dot.Spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMagicHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicHitCheck(sim, attackTable) {
		if spell.MagicCritCheck(sim, result.Target) {
			result.Outcome = OutcomeCrit
			result.Damage *= spell.CritDamageMultiplier()
			if countHits {
				spell.SpellMetrics[result.Target.UnitIndex].Crits++
			}
		} else {
			result.Outcome = OutcomeHit
			if countHits {
				spell.SpellMetrics[result.Target.UnitIndex].Hits++
			}
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicCrit(sim *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeMagicCrit(sim, result, nil, true)
}
func (spell *Spell) OutcomeMagicCritNoHitCounter(sim *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeMagicCrit(sim, result, nil, false)
}
func (spell *Spell) outcomeMagicCrit(sim *Simulation, result *SpellResult, _ *AttackTable, countHits bool) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.MagicCritCheck(sim, result.Target) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritDamageMultiplier()
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		}
	} else {
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	}
}

func (spell *Spell) OutcomeHealing(_ *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealing(nil, result, nil, true)
}
func (spell *Spell) OutcomeHealingNoHitCounter(_ *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealing(nil, result, nil, false)
}
func (spell *Spell) outcomeHealing(_ *Simulation, result *SpellResult, _ *AttackTable, countHits bool) {
	result.Outcome = OutcomeHit
	if countHits {
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

func (spell *Spell) OutcomeHealingCrit(sim *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealingCrit(sim, result, nil, true)
}
func (spell *Spell) OutcomeHealingCritNoHitCounter(sim *Simulation, result *SpellResult, _ *AttackTable) {
	spell.outcomeHealingCrit(sim, result, nil, false)
}
func (spell *Spell) outcomeHealingCrit(sim *Simulation, result *SpellResult, _ *AttackTable, countHits bool) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.HealingCritCheck(sim) {
		result.Outcome = OutcomeCrit
		result.Damage *= spell.CritDamageMultiplier()
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		}
	} else {
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	}
}

func (spell *Spell) OutcomeTickMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
	}
}

func (spell *Spell) OutcomeTickMagicHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.MagicHitCheck(sim, attackTable) {
		if spell.MagicCritCheck(sim, result.Target) {
			result.Outcome = OutcomeCrit
			result.Damage *= spell.CritDamageMultiplier()
			spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		} else {
			result.Outcome = OutcomeHit
			spell.SpellMetrics[result.Target.UnitIndex].Ticks++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMagicHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMagicHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMagicHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.MagicHitCheck(sim, attackTable) {
		result.Outcome = OutcomeHit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Hits++
		}
	} else {
		result.Outcome = OutcomeMiss
		result.Damage = 0
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
	}
}

func (spell *Spell) OutcomeMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWhite(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWhiteNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWhite(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0
	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) {
			if result.applyAttackTableGlance(spell, attackTable, roll, &chance) ||
				result.applyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
				result.applyAttackTableBlock(sim, spell, attackTable)
			} else if !result.applyAttackTableBlock(sim, spell, attackTable) {
				result.applyAttackTableHit(spell, countHits)
			}
		}
	} else {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableGlance(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeWhiteNoGlance(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0
	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) {
			if result.applyAttackTableCrit(spell, attackTable, roll, &chance, true) {
				result.applyAttackTableBlock(sim, spell, attackTable)
			} else if !result.applyAttackTableBlock(sim, spell, attackTable) {
				result.applyAttackTableHit(spell, true)
			}
		}
	} else {
		if !result.applyAttackTableMiss(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCrit(spell, attackTable, roll, &chance, true) {
			result.applyAttackTableHit(spell, true)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(sim, spell, attackTable) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) {
			if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
				result.applyAttackTableBlock(sim, spell, attackTable)
			} else if !result.applyAttackTableBlock(sim, spell, attackTable) {
				result.applyAttackTableHit(spell, countHits)
			}
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialNoParry(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoParry(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWeaponSpecialNoParryNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoParry(sim, result, attackTable, false)
}

func (spell *Spell) outcomeMeleeWeaponSpecialNoParry(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(sim, spell, attackTable) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, countHits)
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWeaponSpecialHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialHitAndCrit(sim, result, attackTable, false)
}

// Like OutcomeMeleeSpecialHitAndCrit, but blocks prevent crits (all weapon damage based attacks).
func (spell *Spell) outcomeMeleeWeaponSpecialHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if spell.Unit.PseudoStats.InFrontOfTarget {
		roll := sim.RandomFloat("White Hit Table")
		chance := 0.0

		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(sim, spell, attackTable) &&
			!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		spell.outcomeMeleeSpecialHitAndCrit(sim, result, attackTable, countHits)
	}
}

func (spell *Spell) OutcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeWeaponSpecialNoCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeWeaponSpecialNoCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeWeaponSpecialNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	unit := spell.Unit
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if unit.PseudoStats.InFrontOfTarget {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableParry(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableBlock(sim, spell, attackTable) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
			!result.applyAttackTableDodge(spell, attackTable, roll, &chance) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParry(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParry(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialNoBlockDodgeParry(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialNoBlockDodgeParryNoCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialNoBlockDodgeParryNoCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialCritOnly(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialCritOnlyNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialCritOnly(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeMeleeSpecialBlockAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialBlockAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeMeleeSpecialBlockAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeMeleeSpecialBlockAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeMeleeSpecialBlockAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {

	if spell.Unit.PseudoStats.InFrontOfTarget {
		if result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableBlock(sim, spell, attackTable)
		} else if !result.applyAttackTableBlock(sim, spell, attackTable) {
			result.applyAttackTableHit(spell, countHits)
		}
	} else {
		if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) OutcomeRangedHit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableDodge(spell, attackTable, roll, &chance) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCrit(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitAndCritNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCrit(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHitAndCrit(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableDodge(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}

}

func (dot *Dot) OutcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	dot.outcomeRangedHitAndCritSnapshot(sim, result, attackTable, true)
}
func (dot *Dot) OutcomeRangedHitAndCritSnapshotNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	dot.outcomeRangedHitAndCritSnapshot(sim, result, attackTable, false)
}
func (dot *Dot) outcomeRangedHitAndCritSnapshot(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(dot.Spell, attackTable, roll, &chance) &&
		!result.applyAttackTableDodge(dot.Spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRollSnapshot(sim, dot) {
		result.applyAttackTableHit(dot.Spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedHitAndCritNoBlock(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCritNoBlock(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedHitAndCritNoBlockNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedHitAndCritNoBlock(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedHitAndCritNoBlock(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("White Hit Table")
	chance := 0.0

	if !result.applyAttackTableMissNoDWPenalty(spell, attackTable, roll, &chance) &&
		!result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeRangedCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedCritOnly(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeRangedCritOnlyNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeRangedCritOnly(sim, result, attackTable, false)
}
func (spell *Spell) outcomeRangedCritOnly(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {

	if !result.applyAttackTableCritSeparateRoll(sim, spell, attackTable, countHits) {
		result.applyAttackTableHit(spell, countHits)
	}
}

func (spell *Spell) OutcomeEnemyMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeEnemyMeleeWhite(sim, result, attackTable, true)
}
func (spell *Spell) OutcomeEnemyMeleeWhiteNoHitCounter(sim *Simulation, result *SpellResult, attackTable *AttackTable) {
	spell.outcomeEnemyMeleeWhite(sim, result, attackTable, false)
}
func (spell *Spell) outcomeEnemyMeleeWhite(sim *Simulation, result *SpellResult, attackTable *AttackTable, countHits bool) {
	roll := sim.RandomFloat("Enemy White Hit Table")
	chance := 0.0

	if !result.applyEnemyAttackTableMiss(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableDodge(spell, attackTable, roll, &chance) &&
		!result.applyEnemyAttackTableParry(spell, attackTable, roll, &chance) {
		if result.applyEnemyAttackTableCrit(spell, attackTable, roll, &chance, countHits) {
			result.applyEnemyAttackTableBlock(sim, spell, attackTable)
		} else if !result.applyEnemyAttackTableBlock(sim, spell, attackTable) {
			result.applyAttackTableHit(spell, countHits)
		}
	}
}

func (spell *Spell) fixedCritCheck(sim *Simulation, critChance float64) bool {
	return sim.RandomFloat("Fixed Crit Roll") < critChance
}

func (spell *Spell) GetPhysicalMissChance(attackTable *AttackTable) float64 {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(attackTable)

	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}

	return max(0, missChance)
}

func (result *SpellResult) applyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance = spell.GetPhysicalMissChance(attackTable)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableMissNoDWPenalty(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := attackTable.BaseMissChance - spell.PhysicalHitChance(attackTable)
	*chance = max(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableBlock(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	chance := attackTable.BaseBlockChance

	if sim.RandomFloat("Block Roll") < chance {
		result.Outcome |= OutcomeBlock
		if result.DidCrit() {
			// Subtract Crits because they happen before Blocks
			spell.SpellMetrics[result.Target.UnitIndex].Crits--
			spell.SpellMetrics[result.Target.UnitIndex].CritBlocks++
		} else if result.DidGlance() {
			// Subtract Glances because they happen before Blocks
			spell.SpellMetrics[result.Target.UnitIndex].Glances--
			spell.SpellMetrics[result.Target.UnitIndex].GlanceBlocks++

		} else {
			spell.SpellMetrics[result.Target.UnitIndex].Blocks++
		}
		damageReduced := result.Damage * (1 - result.Target.BlockDamageReduction())
		result.Damage = max(0, damageReduced)

		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if spell.Flags.Matches(SpellFlagCannotBeDodged) {
		return false
	}

	*chance += max(0, attackTable.BaseDodgeChance-spell.DodgeSuppression())

	if roll < *chance {
		result.Outcome = OutcomeDodge
		spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += max(0, attackTable.BaseParryChance-spell.ParrySuppression(attackTable))

	if roll < *chance {
		result.Outcome = OutcomeParry
		spell.SpellMetrics[result.Target.UnitIndex].Parries++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableGlance(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	*chance += attackTable.BaseGlanceChance

	if roll < *chance {
		result.Outcome = OutcomeGlance
		spell.SpellMetrics[result.Target.UnitIndex].Glances++
		// TODO glancing blow damage reduction is actually a range ([65%, 85%] vs. +3, [80%, 90%] vs. +2, [91%, 99%] vs. +1 and +0)
		result.Damage *= attackTable.GlanceMultiplier
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCrit(spell *Spell, attackTable *AttackTable, roll float64, chance *float64, countHits bool) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	*chance += spell.PhysicalCritChance(attackTable)

	if roll < *chance {
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		}
		result.Damage *= spell.CritDamageMultiplier()
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableCritSeparateRoll(sim *Simulation, spell *Spell, attackTable *AttackTable, countHits bool) bool {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}
	if spell.PhysicalCritCheck(sim, attackTable) {
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		}
		result.Damage *= spell.CritDamageMultiplier()
		return true
	}
	return false
}
func (result *SpellResult) applyAttackTableCritSeparateRollSnapshot(sim *Simulation, dot *Dot) bool {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}
	if sim.RandomFloat("Physical Crit Roll") < dot.SnapshotCritChance {
		result.Outcome = OutcomeCrit
		result.Damage *= dot.Spell.CritDamageMultiplier()
		dot.Spell.SpellMetrics[result.Target.UnitIndex].CritTicks++
		return true
	}
	return false
}

func (result *SpellResult) applyAttackTableHit(spell *Spell, countHits bool) {
	result.Outcome = OutcomeHit
	if countHits {
		spell.SpellMetrics[result.Target.UnitIndex].Hits++
	}
}

func (result *SpellResult) applyEnemyAttackTableMiss(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	missChance := result.Target.GetTotalChanceToBeMissedAsDefender(attackTable) + spell.Unit.PseudoStats.IncreasedMissChance
	if spell.Unit.AutoAttacks.IsDualWielding && !spell.Unit.PseudoStats.DisableDWMissPenalty {
		missChance += 0.19
	}
	*chance += max(0, missChance)

	if roll < *chance {
		result.Outcome = OutcomeMiss
		spell.SpellMetrics[result.Target.UnitIndex].Misses++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableBlock(sim *Simulation, spell *Spell, attackTable *AttackTable) bool {
	if !result.Target.PseudoStats.CanBlock || result.Target.PseudoStats.Stunned {
		return false
	}

	chance := result.Target.GetTotalBlockChanceAsDefender(attackTable)

	if sim.RandomFloat("Player Block") < chance {
		result.Outcome |= OutcomeBlock
		if result.DidCrit() {
			// Subtract Crits because they happen before Blocks
			spell.SpellMetrics[result.Target.UnitIndex].Crits--
			spell.SpellMetrics[result.Target.UnitIndex].CritBlocks++
		} else if result.DidGlance() {
			// Subtract Glances because they happen before Blocks
			spell.SpellMetrics[result.Target.UnitIndex].Glances--
			spell.SpellMetrics[result.Target.UnitIndex].GlanceBlocks++
		} else {
			spell.SpellMetrics[result.Target.UnitIndex].Blocks++
		}

		if result.Target.Blockhandler != nil {
			result.Target.Blockhandler(sim, spell, result)
			return true
		}

		result.Damage = max(0, result.Damage*(1-result.Target.BlockDamageReduction()))

		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableDodge(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if result.Target.PseudoStats.Stunned {
		return false
	}

	*chance += max(result.Target.GetTotalDodgeChanceAsDefender(attackTable), 0.0)

	if roll < *chance {
		result.Outcome = OutcomeDodge
		spell.SpellMetrics[result.Target.UnitIndex].Dodges++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableParry(spell *Spell, attackTable *AttackTable, roll float64, chance *float64) bool {
	if !result.Target.PseudoStats.CanParry || result.Target.PseudoStats.Stunned {
		return false
	}

	*chance += result.Target.GetTotalParryChanceAsDefender(attackTable)

	if roll < *chance {
		result.Outcome = OutcomeParry
		spell.SpellMetrics[result.Target.UnitIndex].Parries++
		result.Damage = 0
		return true
	}
	return false
}

func (result *SpellResult) applyEnemyAttackTableCrit(spell *Spell, _ *AttackTable, roll float64, chance *float64, countHits bool) bool {

	critPercent := spell.Unit.stats[stats.PhysicalCritPercent] + spell.BonusCritPercent
	critChance := critPercent / 100
	critChance -= result.Target.PseudoStats.ReducedCritTakenChance
	*chance += max(0, critChance)

	if roll < *chance {
		result.Outcome = OutcomeCrit
		if countHits {
			spell.SpellMetrics[result.Target.UnitIndex].Crits++
		}
		result.Damage *= 2
		return true
	}
	return false
}

func (spell *Spell) OutcomeExpectedTick(_ *Simulation, _ *SpellResult, _ *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicAlwaysHit(_ *Simulation, _ *SpellResult, _ *AttackTable) {
	// result.Damage *= 1
}
func (spell *Spell) OutcomeExpectedMagicHit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicCrit(_ *Simulation, result *SpellResult, _ *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier += spell.SpellCritChance(result.Target) * (spell.CritDamageMultiplier() - 1)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMagicHitAndCrit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	if spell.CritMultiplier == 0 {
		panic("Spell " + spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier -= spell.SpellChanceToMiss(attackTable)
	averageMultiplier += averageMultiplier * spell.SpellCritChance(result.Target) * (spell.CritDamageMultiplier() - 1)

	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMeleeWhite(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	missChance := spell.GetPhysicalMissChance(attackTable)
	dodgeChance := TernaryFloat64(spell.Flags.Matches(SpellFlagCannotBeDodged), 0, max(0, attackTable.BaseDodgeChance-spell.DodgeSuppression()))
	parryChance := TernaryFloat64(spell.Unit.PseudoStats.InFrontOfTarget, max(0, attackTable.BaseParryChance-spell.ParrySuppression(attackTable)), 0)
	glanceChance := attackTable.BaseGlanceChance
	blockChance := TernaryFloat64(spell.Unit.PseudoStats.InFrontOfTarget, attackTable.BaseBlockChance, 0)
	whiteCritCap := 1.0 - missChance - dodgeChance - parryChance - glanceChance
	critChance := min(spell.PhysicalCritChance(attackTable), whiteCritCap)
	averageMultiplier := (1.0 - missChance - dodgeChance - parryChance + (spell.CritDamageMultiplier()-1)*critChance - glanceChance*(1.0-attackTable.GlanceMultiplier)) * (1.0 - blockChance*result.Target.BlockDamageReduction())
	result.Damage *= averageMultiplier
}

func (spell *Spell) OutcomeExpectedMeleeWeaponSpecialHitAndCrit(_ *Simulation, result *SpellResult, attackTable *AttackTable) {
	missChance := max(0, attackTable.BaseMissChance-spell.PhysicalHitChance(attackTable))
	dodgeChance := TernaryFloat64(spell.Flags.Matches(SpellFlagCannotBeDodged), 0, max(0, attackTable.BaseDodgeChance-spell.DodgeSuppression()))
	parryChance := TernaryFloat64(spell.Unit.PseudoStats.InFrontOfTarget, max(0, attackTable.BaseParryChance-spell.ParrySuppression(attackTable)), 0)
	blockChance := TernaryFloat64(spell.Unit.PseudoStats.InFrontOfTarget, attackTable.BaseBlockChance, 0)
	critChance := spell.PhysicalCritChance(attackTable)
	critFactor := (spell.CritDamageMultiplier() - 1) * critChance
	averageMultiplier := (1.0 - missChance - dodgeChance - parryChance) * (1.0 + critFactor - blockChance*(critFactor+result.Target.BlockDamageReduction()))
	result.Damage *= averageMultiplier
}

func (dot *Dot) OutcomeExpectedMagicSnapshotCrit(_ *Simulation, result *SpellResult, _ *AttackTable) {
	if dot.Spell.CritMultiplier == 0 {
		panic("Spell " + dot.Spell.ActionID.String() + " missing CritMultiplier")
	}

	averageMultiplier := 1.0
	averageMultiplier += dot.SnapshotCritChance * (dot.Spell.CritDamageMultiplier() - 1)

	result.Damage *= averageMultiplier
}
