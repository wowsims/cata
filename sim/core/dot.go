package core

import (
	"math"
	"strconv"
	"time"
)

type OnSnapshot func(sim *Simulation, target *Unit, dot *Dot, isRollover bool)
type OnTick func(sim *Simulation, target *Unit, dot *Dot)

type DotConfig struct {
	// Optional, will default to the corresponding spell.
	Spell *Spell

	OnSnapshot OnSnapshot
	OnTick     OnTick

	Aura Aura

	TickLength    time.Duration // time between each tick
	NumberOfTicks int32         // number of ticks over the whole duration

	IsAOE                bool // Set to true for AOE dots (Blizzard, Hurricane, Consecrate, etc)
	SelfOnly             bool // Set to true to only create the self-hot.
	AffectedByCastSpeed  bool // tick length are shortened based on casting speed
	HasteReducesDuration bool // does not gain additional ticks after a certain haste threshold

	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)

	PeriodicDamageMultiplier float64 // Multiplier for periodic damage on top of the spell's damage multiplier
}

type Dot struct {
	Spell *Spell

	*Aura // Embed Aura, so we can use IsActive/Refresh/etc directly.

	onSnapshot OnSnapshot
	onTick     OnTick
	tickAction *PendingAction

	tickPeriod     time.Duration // hasted time between each tick, rounded to full ms
	BaseTickLength time.Duration // time between each tick

	SnapshotBaseDamage         float64
	SnapshotCritChance         float64
	SnapshotAttackerMultiplier float64

	BaseTickCount  int32 // base tick count without haste applied
	remainingTicks int32
	tmpExtraTicks  int32 // extra ticks that are added during the runtime of the dot

	BonusCoefficient float64 // EffectBonusCoefficient in SpellEffect client DB table, "SP mod" on Wowhead (not necessarily shown there even if > 0)

	PeriodicDamageMultiplier float64 // Multiplier for periodic damage on top of the spell's damage multiplier

	affectedByCastSpeed  bool // tick length are shortened based on casting speed
	hasteReducesDuration bool // does not gain additional ticks after a haste threshold, HasteAffectsDuration in dbc
	isChanneled          bool
}

// Takes a new snapshot of this Dot's effects.
//
// In most cases this will be called automatically, and should only be called
// to force a new snapshot to be taken.
//
// doRollover will apply previously snapshotted crit/%dmg instead of recalculating.
func (dot *Dot) TakeSnapshot(sim *Simulation, doRollover bool) {
	if dot.onSnapshot != nil {
		dot.onSnapshot(sim, dot.Unit, dot, doRollover)
	}
}

// Snapshots and activates the Dot
// If the Dot is already active it's duration will be refreshed and the last tick from the previous application will be
// transfered to the new one
func (dot *Dot) Apply(sim *Simulation) {
	if dot.Spell.Flags&SpellFlagSupressDoTApply > 0 {
		return
	}

	dot.TakeSnapshot(sim, false)
	dot.recomputeAuraDuration(sim)
	dot.Activate(sim)
}

// Rolls over and activates the Dot
// If the Dot is already active it's duration will be refreshed and the last tick from the previous application will be
// transfered to the new one
func (dot *Dot) ApplyRollover(sim *Simulation) {
	if dot.Spell.Flags&SpellFlagSupressDoTApply > 0 {
		return
	}

	dot.TakeSnapshot(sim, true)
	dot.recomputeAuraDuration(sim)
	dot.Activate(sim)
}

func (dot *Dot) recomputeAuraDuration(sim *Simulation) {
	nextTick := dot.TimeUntilNextTick(sim)

	dot.remainingTicks = dot.BaseTickCount
	dot.tmpExtraTicks = 0
	if dot.affectedByCastSpeed {
		// round the tickPeriod to the nearest full ms, same as ingame. This can best be seen ingame in how haste caps
		// work. For example shadowflame should take 1009 haste rating with the 5%/3% haste buffs without rounding, but
		// because of the rounding it already applies at 1007 haste rating.
		dot.tickPeriod = dot.Spell.Unit.ApplyCastSpeedForSpell(dot.BaseTickLength, dot.Spell).Round(time.Millisecond)

		if !dot.hasteReducesDuration {
			dot.remainingTicks = dot.HastedTickCount()
		}
	} else {
		dot.tickPeriod = dot.BaseTickLength
	}
	dot.Duration = dot.tickPeriod * time.Duration(dot.remainingTicks)

	// we a have running dot tick
	// the next tick never gets clipped and is added onto the dot's time for hasted dots
	// see: https://github.com/wowsims/mop/issues/50
	if dot.IsActive() {
		dot.Duration += nextTick
		dot.remainingTicks++
	}
}

// TickPeriod is how fast the snapshotted dot ticks.
func (dot *Dot) TickPeriod() time.Duration {
	return dot.tickPeriod
}

func (dot *Dot) NextTickAt() time.Duration {
	if !dot.IsActive() {
		return 0
	}
	return dot.tickAction.NextActionAt
}

func (dot *Dot) TimeUntilNextTick(sim *Simulation) time.Duration {
	return dot.NextTickAt() - sim.CurrentTime
}

func (dot *Dot) calculateHastedTickCount(baseDuration time.Duration, tickPeriod time.Duration) int32 {
	return int32(math.RoundToEven(float64(baseDuration) / float64(tickPeriod)))
}

// Returns the total amount of ticks with the snapshotted haste
func (dot *Dot) HastedTickCount() int32 {
	return dot.calculateHastedTickCount(dot.BaseDuration(), dot.tickPeriod)
}

func (dot *Dot) ExpectedTickCount() int32 {
	tickCount := dot.BaseTickCount
	if dot.affectedByCastSpeed && !dot.hasteReducesDuration {
		tickPeriod := dot.Spell.Unit.ApplyCastSpeedForSpell(dot.BaseTickLength, dot.Spell).Round(time.Millisecond)
		tickCount = dot.calculateHastedTickCount(dot.BaseDuration(), tickPeriod)
	}
	return tickCount
}

func (dot *Dot) RemainingTicks() int32 {
	return dot.remainingTicks
}

func (dot *Dot) TickCount() int32 {
	return dot.HastedTickCount() + dot.tmpExtraTicks - dot.remainingTicks
}

func (dot *Dot) OutstandingDmg() float64 {
	return TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.remainingTicks), 0)
}

func (dot *Dot) BaseDuration() time.Duration {
	return time.Duration(dot.BaseTickCount) * dot.BaseTickLength
}

// Adds a tick to the current active dot and extends it's duration
func (dot *Dot) AddTick() {
	if !dot.active {
		return
	}

	dot.tmpExtraTicks++
	dot.remainingTicks++
	dot.UpdateExpires(dot.expires + dot.TickPeriod())
}

// Copy's the original DoT's period and duration to the current DoT.
// This is only currently used for Mage's Impact DoT spreading and Enhancement's ImprovedLava Lash.
func (dot *Dot) CopyDotAndApply(sim *Simulation, originaldot *Dot) {
	dot.TakeSnapshot(sim, false)
	dot.SnapshotBaseDamage = originaldot.SnapshotBaseDamage

	dot.tickPeriod = originaldot.tickPeriod
	dot.remainingTicks = originaldot.remainingTicks
	dot.tmpExtraTicks = 0

	// must be set before Activate
	dot.Duration = originaldot.ExpiresAt() - sim.CurrentTime // originaldot.Duration
	dot.UpdateExpires(originaldot.ExpiresAt())

	dot.Activate(sim)

	// recreate with new period, resetting the next tick.
	if dot.tickAction != nil {
		dot.tickAction.Cancel(sim)
	}
	pa := &PendingAction{
		NextActionAt: originaldot.tickAction.NextActionAt,
		OnAction:     dot.periodicTick,
	}
	dot.tickAction = pa
	sim.AddPendingAction(dot.tickAction)
}

// This is the incredibly cursed way fel flame uses to increase dot duration, don't use unless you know what you're
// doing. It extends the duration, immediately recalculates the next tick and then fits as many ticks into the rest of
// the aura duration as it can. This will cause aura duration and dot ticks to desync ingame, so the aura will fall off
// prematurely to what is shown.
//
// Sometimes the game also decides to tick one last time anyway, even though the time since the last tick is absurdly
// low, though this isn't implemented until someone figures out the conditions.
func (dot *Dot) DurationExtendSnapshot(sim *Simulation, extendBy time.Duration) {
	if !dot.IsActive() {
		panic("Can't extend a non-active dot")
	}
	dot.TakeSnapshot(sim, false)

	previousTick := dot.tickAction.NextActionAt - dot.tickPeriod
	dot.tickPeriod = dot.Spell.Unit.ApplyCastSpeedForSpell(dot.BaseTickLength, dot.Spell).Round(time.Millisecond)

	// ensure the tick is at least scheduled for the future ..
	nextTick := max(previousTick+dot.tickPeriod, sim.CurrentTime+1*time.Millisecond)

	dot.tickAction.Cancel(sim)
	dot.tickAction = &PendingAction{
		NextActionAt: nextTick,
		// Priority:     ActionPriorityDOT,
		OnAction: dot.periodicTick,
	}

	// cap the total duration to the amount of hasted ticks a new dot would have
	extendDuration := min(dot.RemainingDuration(sim)+extendBy,
		dot.tickPeriod*time.Duration(dot.HastedTickCount()-1)+(nextTick-sim.CurrentTime))
	dot.remainingTicks = int32((extendDuration-(nextTick-sim.CurrentTime))/dot.tickPeriod) + 1

	dot.Duration = nextTick - sim.CurrentTime + time.Duration(dot.remainingTicks-1)*dot.tickPeriod
	sim.AddPendingAction(dot.tickAction)
	dot.Refresh(sim)
}

// Forces an instant tick. Does not reset the tick timer or aura duration,
// the tick is simply an extra tick.
func (dot *Dot) TickOnce(sim *Simulation) {
	dot.onTick(sim, dot.Unit, dot)
}

func (dot *Dot) periodicTick(sim *Simulation) {
	dot.remainingTicks--
	dot.TickOnce(sim)
	if dot.isChanneled {
		// Note: even if the clip delay is 0ms, need a WaitUntil so that APL is called after the channel aura fades.
		if dot.remainingTicks == 0 && dot.Spell.Unit.GCD.IsReady(sim) {
			dot.Spell.Unit.WaitUntil(sim, sim.CurrentTime+dot.Spell.Unit.ChannelClipDelay)
		} else if dot.Spell.Unit.Rotation.shouldInterruptChannel(sim) {
			dot.tickAction.NextActionAt = NeverExpires // don't tick again in ApplyOnExpire
			dot.Deactivate(sim)
			if dot.Spell.Unit.GCD.IsReady(sim) {
				dot.Spell.Unit.WaitUntil(sim, sim.CurrentTime+dot.Spell.Unit.ChannelClipDelay)
			}
			return // don't schedule another tick
		}
	}

	dot.tickAction.NextActionAt = sim.CurrentTime + dot.tickPeriod
	sim.AddPendingAction(dot.tickAction)
}

func newDot(config Dot) *Dot {
	dot := &config

	dot.tickPeriod = dot.BaseTickLength
	dot.Duration = dot.tickPeriod * time.Duration(dot.BaseTickCount)

	dot.ApplyOnGain(func(aura *Aura, sim *Simulation) {
		dot.tickAction = &PendingAction{
			NextActionAt: sim.CurrentTime + dot.tickPeriod,
			// Priority:     ActionPriorityDOT,
			OnAction: dot.periodicTick,
		}
		sim.AddPendingAction(dot.tickAction)
		if dot.isChanneled {
			dot.Spell.Unit.ChanneledDot = dot
		}
	})
	dot.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
		// the core scheduling fails to process ticks first so we need to apply the last tick
		if dot.tickAction.NextActionAt == sim.CurrentTime {
			dot.remainingTicks--
			dot.TickOnce(sim)
			// Note: even if the clip delay is 0ms, need a WaitUntil so that APL is called after the channel aura fades.
			if dot.isChanneled && dot.Spell.Unit.GCD.IsReady(sim) {
				dot.Spell.Unit.WaitUntil(sim, sim.CurrentTime+dot.Spell.Unit.ChannelClipDelay)
			}
		}

		dot.tickAction.Cancel(sim)
		dot.tickAction = nil
		if dot.isChanneled {
			dot.Spell.Unit.ChanneledDot = nil
			dot.Spell.Unit.Rotation.interruptChannelIf = nil
			dot.Spell.Unit.Rotation.allowChannelRecastOnInterrupt = false
			// track time metrics for channels
			dot.Spell.SpellMetrics[aura.Unit.UnitIndex].TotalCastTime += dot.fadeTime - dot.StartedAt()
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

	if config.PeriodicDamageMultiplier == 0 {
		config.PeriodicDamageMultiplier = 1
	}

	if config.Spell == nil {
		config.Spell = spell
	}
	dot := Dot{
		Spell: config.Spell,

		remainingTicks:       config.NumberOfTicks,
		BaseTickCount:        config.NumberOfTicks,
		BaseTickLength:       config.TickLength,
		onSnapshot:           config.OnSnapshot,
		onTick:               config.OnTick,
		affectedByCastSpeed:  config.AffectedByCastSpeed,
		hasteReducesDuration: config.HasteReducesDuration,
		isChanneled:          config.Spell.Flags.Matches(SpellFlagChanneled),

		BonusCoefficient: config.BonusCoefficient,

		PeriodicDamageMultiplier: config.PeriodicDamageMultiplier,
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
