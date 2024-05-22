package core

import (
	"strconv"
	"time"
)

type OnSnapshot func(sim *Simulation, target *Unit, dot *Dot, isRollover bool)
type OnTick func(sim *Simulation, target *Unit, dot *Dot)

type DotConfig struct {
	IsAOE    bool // Set to true for AOE dots (Blizzard, Hurricane, Consecrate, etc)
	SelfOnly bool // Set to true to only create the self-hot.

	// Optional, will default to the corresponding spell.
	Spell *Spell

	Aura Aura

	NumberOfTicks int32         // number of ticks over the whole duration
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed  bool
	HasteAffectsDuration bool

	OnSnapshot OnSnapshot
	OnTick     OnTick

	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)
}

type Dot struct {
	Spell *Spell

	// Embed Aura, so we can use IsActive/Refresh/etc directly.
	*Aura

	NumberOfTicks int32         // number of ticks over the whole duration
	BaseTickCount int32         // base tick count without haste applied
	TickLength    time.Duration // time between each tick

	// If true, tick length will be shortened based on casting speed.
	AffectedByCastSpeed  bool
	HasteAffectsDuration bool

	OnSnapshot OnSnapshot
	OnTick     OnTick

	SnapshotBaseDamage         float64
	SnapshotCritChance         float64
	SnapshotAttackerMultiplier float64

	tickAction *PendingAction
	tickPeriod time.Duration

	// Number of ticks since last call to Apply().
	TickCount int32

	lastTickTime time.Duration
	isChanneled  bool

	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)
}

// TickPeriod is how fast the snapshot dot ticks.
func (dot *Dot) TickPeriod() time.Duration {
	return dot.tickPeriod
}

func (dot *Dot) NextTickAt() time.Duration {
	return dot.lastTickTime + dot.tickPeriod
}

func (dot *Dot) TimeUntilNextTick(sim *Simulation) time.Duration {
	return dot.NextTickAt() - sim.CurrentTime
}

func (dot *Dot) MaxTicksRemaining() int32 {
	return dot.NumberOfTicks - dot.TickCount
}

func (dot *Dot) NumTicksRemaining(sim *Simulation) int32 {
	maxTicksRemaining := dot.MaxTicksRemaining()
	finalTickAt := dot.lastTickTime + dot.tickPeriod*time.Duration(maxTicksRemaining)
	return max(0, int32((finalTickAt-sim.CurrentTime)/dot.tickPeriod)+1)
}

// Roll over = gets carried over with everlasting refresh and doesn't get applied if triggered when the spell is already up.
// - Example: critical strike rating, internal % damage modifiers: buffs or debuffs on player
// Nevermelting Ice, Shadow Mastery (ISB), Trick of the Trades, Deaths Embrace, Thaddius Polarity, Hera Spores, Crit on weapons from swapping

// Snapshot = calculation happens at refresh and application (stays up even if buff falls of, until new refresh or application)
// - Example: Spell power, Haste rating
// Blood Fury, Lightweave Embroid, Eradication, Bloodlust

// Dynamic = realtime update
// - Example: external % damage modifier debuffs on target
// Haunt, Curse of Shadow, Shadow Embrace

// Rollover is used to reset the duration of a dot from an external spell (not casting the dot itself)
// This keeps the snapshot crit and %dmg modifiers.
// However, sp and haste are recalculated.
func (dot *Dot) Rollover(sim *Simulation) {
	dot.TakeSnapshot(sim, true)

	dot.RecomputeAuraDuration() // recalculate haste
	dot.Aura.Refresh(sim)       // update aura's duration

	oldNextTick := dot.tickAction.NextActionAt
	dot.tickAction.Cancel(sim) // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	dot.tickAction.NextActionAt = oldNextTick
	sim.AddPendingAction(dot.tickAction)
}

func (dot *Dot) RescheduleNextTick(sim *Simulation) {
	dot.RecomputeAuraDuration()

	dot.tickAction.Cancel(sim) // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	dot.tickAction.NextActionAt = dot.lastTickTime + dot.tickPeriod
	sim.AddPendingAction(dot.tickAction)
}

// Snapshots and activates the Dot
// If the Dot is running it's duration will be refreshed and
// if there was a next Dot happening this will carry over to the new Dot
func (dot *Dot) Apply(sim *Simulation) {

	dot.TakeSnapshot(sim, false)

	dot.TickCount = 0

	// we a have running dot tick
	// the next tick never get's clipped and is added onto the dot's time for hasted dots
	// see: https://github.com/wowsims/cata/issues/50git
	if dot.tickAction != nil && !dot.tickAction.cancelled {

		// save next tick timer as timer is computed based on tick time
		// which we update in RecomputeAuraDuration
		nextTick := dot.TimeUntilNextTick(sim)
		dot.RecomputeAuraDuration()
		dot.Aura.Duration += nextTick

		// add extra tick
		dot.TickCount--

		// update tick action to work with new tick rate, but set next tick to still occur
		oldNextAction := dot.tickAction.NextActionAt
		dot.tickAction.Cancel(sim)
		periodicOptions := dot.basePeriodicOptions()
		periodicOptions.Period = dot.tickPeriod
		dot.tickAction = NewPeriodicAction(sim, periodicOptions)
		dot.tickAction.NextActionAt = oldNextAction
		sim.AddPendingAction(dot.tickAction)
	} else {
		dot.RecomputeAuraDuration()
	}

	dot.Aura.Activate(sim)
}

// ApplyOrReset is used for rolling dots that reset the tick timer on reapplication.
// This is more efficient than Apply(), and works around tickAction.CleanUp() wrongly generating
// an extra ticks if (re-)application and tick happen at the same time.
func (dot *Dot) ApplyOrReset(sim *Simulation) {
	if !dot.IsActive() {
		dot.Apply(sim)
		return
	}

	dot.TakeSnapshot(sim, true)

	dot.RecomputeAuraDuration() // recalculate haste
	dot.Aura.Refresh(sim)       // update aura's duration

	dot.TickCount = 0

	oldTickAction := dot.tickAction
	dot.tickAction = nil      // prevent tickAction.CleanUp() from adding an extra tick
	oldTickAction.Cancel(sim) // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	sim.AddPendingAction(dot.tickAction)
}

func (dot *Dot) Cancel(sim *Simulation) {
	if dot.Aura.IsActive() {
		dot.Aura.Deactivate(sim)
	}
}

// Call this after manually changing NumberOfTicks or TickLength.
func (dot *Dot) RecomputeAuraDuration() {
	if dot.AffectedByCastSpeed {
		dot.tickPeriod = dot.Spell.Unit.ApplyCastSpeedForSpell(dot.TickLength, dot.Spell)

		// cata haste logic here for dots
		// channels seem not to be affected by the same logic
		// see: https://youtu.be/Rr4YyKaU7Ik?si=Isuce7Z1bQWMWpMi&t=53
		if !dot.isChanneled && !dot.HasteAffectsDuration {
			dot.NumberOfTicks = int32(round(float64(dot.GetBaseDuration()) / float64(dot.tickPeriod)))
		}

		dot.Aura.Duration = dot.tickPeriod * time.Duration(dot.NumberOfTicks)
	} else {
		dot.tickPeriod = dot.TickLength
		dot.Aura.Duration = dot.tickPeriod * time.Duration(dot.NumberOfTicks)
	}
}

func (dot *Dot) AddTicks(num int32) {
	dot.BaseTickCount += num
	dot.NumberOfTicks += num
}

func (dot *Dot) GetBaseDuration() time.Duration {
	return time.Duration(dot.BaseTickCount) * dot.TickLength
}

// Takes a new snapshot of this Dot's effects.
//
// In most cases this will be called automatically, and should only be called
// to force a new snapshot to be taken.
//
//	doRollover will apply previously snapshotted crit/%dmg instead of recalculating.
func (dot *Dot) TakeSnapshot(sim *Simulation, doRollover bool) {
	if dot.OnSnapshot != nil {
		dot.OnSnapshot(sim, dot.Unit, dot, doRollover)
	}
}

// Forces an instant tick. Does not reset the tick timer or aura duration,
// the tick is simply an extra tick.
func (dot *Dot) TickOnce(sim *Simulation) {
	dot.lastTickTime = sim.CurrentTime
	dot.OnTick(sim, dot.Unit, dot)

	if dot.isChanneled {
		// Note: even if the clip delay is 0ms, need a WaitUntil so that APL is called after the channel aura fully fades.
		if dot.MaxTicksRemaining() == 0 {
			if dot.Spell.Unit.GCD.IsReady(sim) {
				dot.Spell.Unit.WaitUntil(sim, sim.CurrentTime+dot.Spell.Unit.ChannelClipDelay)
			}
		} else if dot.Spell.Unit.Rotation.shouldInterruptChannel(sim) {
			dot.Cancel(sim)
			if dot.Spell.Unit.GCD.IsReady(sim) {
				dot.Spell.Unit.WaitUntil(sim, sim.CurrentTime+dot.Spell.Unit.ChannelClipDelay)
			}
		}
	}
}

// ManualTick forces the dot forward one tick
// Will cancel the dot if it is out of ticks.
func (dot *Dot) ManualTick(sim *Simulation) {
	if dot.lastTickTime != sim.CurrentTime {
		dot.TickCount++
		if dot.NumTicksRemaining(sim) <= 0 {
			dot.Cancel(sim)
		} else {
			dot.TickOnce(sim)
		}
	}
}

func (dot *Dot) basePeriodicOptions() PeriodicActionOptions {
	return PeriodicActionOptions{
		//Priority: ActionPriorityDOT,
		OnAction: func(sim *Simulation) {
			if dot.lastTickTime != sim.CurrentTime {
				dot.TickCount++
				dot.TickOnce(sim)
			}
		},
		CleanUp: func(sim *Simulation) {
			// In certain cases, the last tick and the dot aura expiration can happen in
			// different orders, so we might need to apply the last tick.
			if dot.tickAction != nil && dot.tickAction.NextActionAt == sim.CurrentTime {
				if dot.lastTickTime != sim.CurrentTime {
					dot.TickCount++
					dot.TickOnce(sim)
				}
			}
		},
	}
}

func newDot(config Dot) *Dot {
	dot := &Dot{}
	*dot = config

	dot.tickPeriod = dot.TickLength
	dot.Aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)

	dot.Aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		dot.lastTickTime = sim.CurrentTime
		periodicOptions := dot.basePeriodicOptions()
		periodicOptions.Period = dot.tickPeriod
		dot.tickAction = NewPeriodicAction(sim, periodicOptions)
		sim.AddPendingAction(dot.tickAction)
		if dot.isChanneled {
			dot.Spell.Unit.ChanneledDot = dot
		}
	})
	dot.Aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		if dot.tickAction != nil {
			dot.tickAction.Cancel(sim)
			dot.tickAction = nil
		}
		if dot.isChanneled {
			dot.Spell.Unit.ChanneledDot = nil
			dot.Spell.Unit.Rotation.interruptChannelIf = nil
			dot.Spell.Unit.Rotation.allowChannelRecastOnInterrupt = false
		}
	})

	return dot
}

type DotArray []*Dot

func (dots DotArray) Get(target *Unit) *Dot {
	return dots[target.UnitIndex]
}

func (spell *Spell) createDots(config DotConfig, isHot bool) {
	if config.NumberOfTicks == 0 && config.TickLength == 0 {
		return
	}

	if config.Spell == nil {
		config.Spell = spell
	}
	dot := Dot{
		Spell: config.Spell,

		NumberOfTicks:        config.NumberOfTicks,
		BaseTickCount:        config.NumberOfTicks,
		TickLength:           config.TickLength,
		AffectedByCastSpeed:  config.AffectedByCastSpeed,
		HasteAffectsDuration: config.HasteAffectsDuration,

		OnSnapshot: config.OnSnapshot,
		OnTick:     config.OnTick,

		isChanneled: config.Spell.Flags.Matches(SpellFlagChanneled),

		BonusCoefficient: config.BonusCoefficient,
	}

	auraConfig := config.Aura
	if auraConfig.ActionID.IsEmptyAction() {
		auraConfig.ActionID = dot.Spell.ActionID
	}

	caster := dot.Spell.Unit
	if config.IsAOE || config.SelfOnly {
		dot.Aura = caster.GetOrRegisterAura(auraConfig)
		spell.aoeDot = newDot(dot)
	} else {
		auraConfig.Label += "-" + strconv.Itoa(int(caster.UnitIndex))
		if spell.dots == nil {
			spell.dots = make([]*Dot, len(caster.Env.AllUnits))
		}
		for _, target := range caster.Env.AllUnits {
			if isHot != caster.IsOpponent(target) {
				dot.Aura = target.GetOrRegisterAura(auraConfig)
				spell.dots[target.UnitIndex] = newDot(dot)
			}
		}
	}
}

func (dot *Dot) CopyDotAndApply(sim *Simulation, originaldot *Dot) {
	dot.TakeSnapshot(sim, false)
	dot.SnapshotBaseDamage = originaldot.SnapshotBaseDamage

	dot.tickPeriod = originaldot.tickPeriod
	//dot.NumberOfTicks = originaldot.NumTicksRemaining(sim) // originaldot.NumberOfTicks
	dot.TickCount = 0 // originaldot.TickCount

	// must be set before Activate
	dot.Aura.Duration = originaldot.ExpiresAt() - sim.CurrentTime // originaldot.Aura.Duration
	dot.UpdateExpires(originaldot.ExpiresAt())

	dot.Aura.Activate(sim)

	dot.lastTickTime = originaldot.lastTickTime // must be set after Activate

	// Copied from RescheduleNextTick without the RecomputeAuraDuration
	dot.tickAction.Cancel(sim) // remove old PA ticker

	// recreate with new period, resetting the next tick.
	periodicOptions := dot.basePeriodicOptions()
	periodicOptions.Period = dot.tickPeriod
	dot.tickAction = NewPeriodicAction(sim, periodicOptions)
	dot.tickAction.NextActionAt = dot.lastTickTime + dot.tickPeriod
	sim.AddPendingAction(dot.tickAction)
}
