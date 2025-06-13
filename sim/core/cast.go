package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
)

// A cast corresponds to any action which causes the in-game castbar to be
// shown, and activates the GCD. Note that a cast can also be instant, i.e.
// the effects are applied immediately even though the GCD is still activated.

// Callback for when a cast is finished, i.e. when the in-game castbar reaches full.
type OnCastComplete func(aura *Aura, sim *Simulation, spell *Spell)

type Hardcast struct {
	Expires    time.Duration
	ActionID   ActionID
	OnComplete func(*Simulation, *Unit)
	Target     *Unit
	CanMove    bool
}

// Input for constructing the CastSpell function for a spell.
type CastConfig struct {
	// Default cast values with all static effects applied.
	DefaultCast Cast

	// Dynamic modifications for each cast.
	ModifyCast func(*Simulation, *Spell, *Cast)

	// Ignores haste when calculating the GCD and cast time for this cast.
	// Automatically set if GCD and cast times are all 0, e.g. for empty casts.
	IgnoreHaste bool

	CD       Cooldown
	SharedCD Cooldown

	CastTime func(spell *Spell) time.Duration
}

type Cast struct {
	// Amount of resource that will be consumed by this cast.
	Cost float64

	// The length of time the GCD will be on CD as a result of this cast.
	GCD time.Duration

	// The minimum length of time for the GCD. Can be left out to use the default of 1s
	GCDMin time.Duration

	// The amount of time between the call to spell.Cast() and when the spell
	// effects are invoked.
	CastTime time.Duration

	// A dummy value to give 0gcd casts a non empty struct
	NonEmpty bool
}

func (cast *Cast) EffectiveTime() time.Duration {
	gcd := max(0, cast.GCD)
	if cast.GCD > 0 {
		if cast.GCDMin != 0 {
			gcd = max(cast.GCDMin, gcd)
		} else {
			gcd = max(GCDMin, gcd)
		}
	}
	return max(gcd, cast.CastTime)
}

type CastFunc func(*Simulation, *Unit)
type CastSuccessFunc func(*Simulation, *Unit) bool

func (spell *Spell) castFailureHelper(sim *Simulation, message string, vals ...any) bool {
	if sim.CurrentTime < 0 && spell.Unit.Rotation != nil {
		spell.Unit.Rotation.ValidationMessage(proto.LogLevel_Warning, fmt.Sprintf(spell.ActionID.String()+" failed to cast: "+message, vals...))
	} else {
		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, fmt.Sprintf(spell.ActionID.String()+" failed to cast: "+message, vals...))
		}
	}
	return false
}

func (spell *Spell) makeCastFunc(config CastConfig) CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		spell.CurCast = spell.DefaultCast

		if config.ModifyCast != nil {
			config.ModifyCast(sim, spell, &spell.CurCast)
			if spell.CurCast.Cost != spell.DefaultCast.Cost {
				// Costs need to be modified using the unit and spell multipliers, so that
				// their affects are also visible in the spell.CanCast() function, which
				// does not invoke ModifyCast.
				panic("May not modify cost in ModifyCast!")
			}
		}

		if spell.Flags.Matches(SpellFlagSwapped) {
			return spell.castFailureHelper(sim, "spell attached to an un-equipped item")
		}

		if spell.ExtraCastCondition != nil {
			if !spell.ExtraCastCondition(sim, target) {
				return spell.castFailureHelper(sim, "extra spell condition")
			}
		}

		if spell.Cost != nil {
			if !spell.Cost.MeetsRequirement(sim, spell) {
				return spell.castFailureHelper(sim, spell.Cost.CostFailureReason(sim, spell))
			}
		}

		if spell.MaxCharges > 0 && spell.charges == 0 {
			return spell.castFailureHelper(sim, "not enough charges")
		}

		if !config.IgnoreHaste {
			spell.CurCast.GCD = max(0, spell.Unit.ApplyCastSpeed(spell.CurCast.GCD)).Round(time.Millisecond)
			spell.CurCast.CastTime = spell.Unit.ApplyCastSpeedForSpell(spell.CurCast.CastTime, spell).Round(time.Millisecond)
		}

		if config.CD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.CD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on cooldown for %s, curTime = %s", spell.CD.TimeToReady(sim), sim.CurrentTime)
			}
		}

		if config.SharedCD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.SharedCD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on shared cooldown for %s, curTime = %s", spell.SharedCD.TimeToReady(sim), sim.CurrentTime)
			}
		}

		// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
		if spell.CurCast.GCD > 0 && !spell.Unit.GCD.IsReady(sim) {
			return spell.castFailureHelper(sim, "GCD on cooldown for %s, curTime = %s", spell.Unit.GCD.TimeToReady(sim), sim.CurrentTime)
		}

		if hc := spell.Unit.Hardcast; hc.Expires > sim.CurrentTime {
			return spell.castFailureHelper(sim, "casting/channeling %v for %s, curTime = %s", hc.ActionID, hc.Expires-sim.CurrentTime, sim.CurrentTime)
		}

		if effectiveTime := spell.CurCast.EffectiveTime(); effectiveTime != 0 {

			// do not add channeled time here as they have variable cast length
			// cast time for channels is handled in dot.OnExpire
			if !spell.Flags.Matches(SpellFlagChanneled) {
				spell.SpellMetrics[target.UnitIndex].TotalCastTime += effectiveTime
			}

			spell.Unit.SetGCDTimer(sim, max(sim.CurrentTime+effectiveTime, spell.Unit.NextGCDAt()))
		}

		if (spell.Flags&SpellFlagCanCastWhileMoving == 0) && (spell.CurCast.CastTime > 0) && spell.Unit.Moving {
			return spell.castFailureHelper(sim, "casting/channeling while moving not allowed!")
		}

		// Hardcasts
		if spell.CurCast.CastTime > 0 {
			if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
				spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
					spell.ActionID, max(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
			}

			spell.Unit.Hardcast = Hardcast{
				Expires:  sim.CurrentTime + spell.CurCast.CastTime,
				ActionID: spell.ActionID,
				OnComplete: func(sim *Simulation, target *Unit) {
					if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
						spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
					}

					if spell.Cost != nil {
						spell.Cost.SpendCost(sim, spell)
					}

					if spell.MaxCharges > 0 {
						spell.ConsumeCharge(sim)
					}

					if config.CD.Timer != nil || spell.rechargeTimer != nil {
						spell.triggerCooldown(sim)
					}

					if config.SharedCD.Timer != nil {
						spell.SharedCD.Set(sim.CurrentTime + time.Duration(float64(spell.SharedCD.Duration)*spell.CdMultiplier))
					}

					spell.applyEffects(sim, target)

					if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
						spell.Unit.OnCastComplete(sim, spell)
					}
				},
				Target:  target,
				CanMove: spell.Flags&SpellFlagCanCastWhileMoving > 0,
			}

			spell.Unit.newHardcastAction(sim)
			return true
		}

		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, max(0, spell.CurCast.Cost), spell.CurCast.CastTime, spell.CurCast.EffectiveTime())
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		if spell.Cost != nil {
			spell.Cost.SpendCost(sim, spell)
		}

		if spell.MaxCharges > 0 {
			spell.ConsumeCharge(sim)
		}

		if config.CD.Timer != nil || spell.rechargeTimer != nil {
			spell.triggerCooldown(sim)
		}

		if config.SharedCD.Timer != nil {
			spell.SharedCD.Set(sim.CurrentTime + time.Duration(float64(spell.SharedCD.Duration)*spell.CdMultiplier))
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}

func (spell *Spell) triggerCooldown(sim *Simulation) {
	cd := time.Duration(float64(spell.CD.Duration) * spell.CdMultiplier)

	// if recharge timer is higher than the actual cooldown of the spell we use
	if spell.MaxCharges > 0 && spell.charges == 0 {
		// spell.CdMultiplier would be considered within the the recharge time if we ever need that
		cd = TernaryDuration(cd > spell.NextChargeIn(sim), cd, spell.NextChargeIn(sim))
	}

	if cd > 0 {
		spell.CD.Set(sim.CurrentTime + cd)
	}
}

func (spell *Spell) makeCastFuncSimple() CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		if spell.Flags.Matches(SpellFlagSwapped) {
			return spell.castFailureHelper(sim, "spell attached to an un-equipped item")
		}

		if spell.ExtraCastCondition != nil {
			if !spell.ExtraCastCondition(sim, target) {
				return spell.castFailureHelper(sim, "extra spell condition")
			}
		}

		if spell.CD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.CD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on cooldown for %s, curTime = %s", spell.CD.TimeToReady(sim), sim.CurrentTime)
			}
		}

		if spell.SharedCD.Timer != nil {
			// By panicking if spell is on CD, we force each sim to properly check for their own CDs.
			if !spell.SharedCD.IsReady(sim) {
				return spell.castFailureHelper(sim, "still on shared cooldown for %s, curTime = %s", spell.SharedCD.TimeToReady(sim), sim.CurrentTime)
			}
		}

		if spell.MaxCharges > 0 && spell.charges == 0 {
			return spell.castFailureHelper(sim, "not enough charges")
		}

		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, 0.0, "0s", "0s")
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		if spell.MaxCharges > 0 {
			spell.ConsumeCharge(sim)
		}

		if spell.CD.Timer != nil || spell.rechargeTimer != nil {
			spell.triggerCooldown(sim)
		}

		if spell.SharedCD.Timer != nil {
			spell.SharedCD.Set(sim.CurrentTime + time.Duration(float64(spell.SharedCD.Duration)*spell.CdMultiplier))
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}

func (spell *Spell) makeCastFuncAutosOrProcs() CastSuccessFunc {
	return func(sim *Simulation, target *Unit) bool {
		if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
			spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
				spell.ActionID, 0.0, "0s", "0s")
			spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
		}

		spell.applyEffects(sim, target)

		if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
			spell.Unit.OnCastComplete(sim, spell)
		}

		return true
	}
}

// Procs a spell, circumventing all checks, cooldowns, gcd's and so on
// Simply logging the cast and applying the effect
// Can be used for spells that proc off other spells and are the same spell id
func (spell *Spell) Proc(sim *Simulation, target *Unit) {
	if sim.Log != nil && !spell.Flags.Matches(SpellFlagNoLogs) {
		spell.Unit.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s, Effective Time = %s)",
			spell.ActionID, 0.0, "0s", "0s")
		spell.Unit.Log(sim, "Completed cast %s", spell.ActionID)
	}

	spell.applyEffects(sim, target)

	if !spell.Flags.Matches(SpellFlagNoOnCastComplete) {
		spell.Unit.OnCastComplete(sim, spell)
	}
}
