package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type OnFocusGain func(*Simulation, float64)

type focusBar struct {
	unit *Unit

	maxFocus           float64
	currentFocus       float64
	baseFocusPerSecond float64
	focusTickDuration  time.Duration
	nextFocusTick      time.Duration
	isPlayer           bool

	// These two terms are multiplied together to scale the total Focus regen from ticks.
	focusRegenMultiplier  float64
	hasteRatingMultiplier float64

	regenMetrics       *ResourceMetrics
	focusRefundMetrics *ResourceMetrics

	OnFocusGain OnFocusGain
}

func (unit *Unit) EnableFocusBar(maxFocus float64, baseFocusPerSecond float64, isPlayer bool, onFocusGain OnFocusGain) {
	unit.SetCurrentPowerBar(FocusBar)

	unit.focusBar = focusBar{
		unit:                  unit,
		maxFocus:              max(100, maxFocus),
		focusTickDuration:     unit.ReactionTime,
		focusRegenMultiplier:  1,
		hasteRatingMultiplier: 1,
		isPlayer:              isPlayer,
		baseFocusPerSecond:    baseFocusPerSecond,
		regenMetrics:          unit.NewFocusMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFocusRegen}),
		focusRefundMetrics:    unit.NewFocusMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
		OnFocusGain:           onFocusGain,
	}
}

func (unit *Unit) HasFocusBar() bool {
	return unit.focusBar.unit != nil
}

func (fb *focusBar) CurrentFocus() float64 {
	return fb.currentFocus
}

func (fb *focusBar) MaximumFocus() float64 {
	return fb.maxFocus
}

func (fb *focusBar) NextFocusTickAt() time.Duration {
	return fb.nextFocusTick
}

func (fb *focusBar) MultiplyFocusRegenSpeed(sim *Simulation, multiplier float64) {
	fb.ResetFocusTick(sim)
	fb.focusRegenMultiplier *= multiplier
}

func (fb *focusBar) FocusRegenPerTick() float64 {
	ticksPerSecond := float64(time.Second) / float64(fb.focusTickDuration)
	return fb.FocusRegenPerSecond() / ticksPerSecond
}

func (fb *focusBar) FocusRegenPerSecond() float64 {
	return fb.baseFocusPerSecond * fb.getTotalRegenMultiplier()
}

func (fb *focusBar) TimeToTargetFocus(targetFocus float64) time.Duration {
	if fb.currentFocus >= targetFocus {
		return time.Duration(0)
	}

	return DurationFromSeconds((targetFocus - fb.currentFocus) / fb.FocusRegenPerSecond())
}

func (fb *focusBar) getTotalRegenMultiplier() float64 {
	return fb.hasteRatingMultiplier * fb.focusRegenMultiplier
}

func (fb *focusBar) AddFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}
	newFocus := min(fb.currentFocus+amount, fb.maxFocus)
	if sim.Log != nil {
		fb.unit.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, fb.currentFocus, newFocus, fb.maxFocus)
	}
	if fb.isPlayer {
		metrics.AddEvent(amount, newFocus-fb.currentFocus)
	}

	if fb.OnFocusGain != nil {
		fb.OnFocusGain(sim, newFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) SpendFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative focus!")
	}

	newFocus := fb.currentFocus - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		fb.unit.Log(sim, "Spent %0.3f focus from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, fb.currentFocus, newFocus, fb.maxFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) IsTicking(sim *Simulation) bool {
	return (fb.nextFocusTick != 0) && (sim.CurrentTime <= fb.nextFocusTick) && (fb.nextFocusTick-sim.CurrentTime <= fb.focusTickDuration)
}

// Gives an immediate partial Focus tick and restarts the tick timer.
func (fb *focusBar) ResetFocusTick(sim *Simulation) {
	if !fb.IsTicking(sim) {
		return
	}

	timeSinceLastTick := max(sim.CurrentTime-(fb.NextFocusTickAt()-fb.focusTickDuration), 0)
	partialTickAmount := fb.FocusRegenPerSecond() * timeSinceLastTick.Seconds()
	fb.AddFocus(sim, partialTickAmount, fb.regenMetrics)
	fb.nextFocusTick = sim.CurrentTime + fb.focusTickDuration
	sim.RescheduleTask(fb.nextFocusTick)
}

func (fb *focusBar) processDynamicHasteRatingChange(sim *Simulation) {
	if fb.unit == nil {
		return
	}

	fb.ResetFocusTick(sim)
	fb.hasteRatingMultiplier = 1.0 + fb.unit.GetStat(stats.HasteRating)/(100*HasteRatingPerHastePercent)
}

func (fb *focusBar) RunTask(sim *Simulation) time.Duration {
	if sim.CurrentTime < fb.nextFocusTick {
		return fb.nextFocusTick
	}
	fb.AddFocus(sim, fb.FocusRegenPerTick(), fb.regenMetrics)
	fb.nextFocusTick = sim.CurrentTime + fb.focusTickDuration
	return fb.nextFocusTick
}

func (fb *focusBar) reset(sim *Simulation) {
	if fb.unit == nil {
		return
	}

	fb.currentFocus = fb.maxFocus
	fb.hasteRatingMultiplier = 1.0 + fb.unit.GetStat(stats.HasteRating)/(100*HasteRatingPerHastePercent)
	fb.focusRegenMultiplier = 1.0

	if fb.unit.Type != PetUnit {
		fb.enable(sim, sim.Environment.PrepullStartTime())
	}
}

func (fb *focusBar) enable(sim *Simulation, startAt time.Duration) {
	sim.AddTask(fb)
	fb.nextFocusTick = startAt + time.Duration(sim.RandomFloat("Focus Tick")*float64(fb.focusTickDuration))
	sim.RescheduleTask(fb.nextFocusTick)
}

func (fb *focusBar) disable(sim *Simulation) {
	fb.nextFocusTick = NeverExpires
	sim.RemoveTask(fb)
}

type FocusCostOptions struct {
	Cost int32

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.FocusRefundMetrics if not supplied.
}

type FocusCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newFocusCost(spell *Spell, options FocusCostOptions) *SpellCost {
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.focusRefundMetrics
	}

	return &SpellCost{
		spell:           spell,
		BaseCost:        options.Cost,
		PercentModifier: 1,
		ResourceCostImpl: &FocusCost{
			Refund:          options.Refund,
			RefundMetrics:   options.RefundMetrics,
			ResourceMetrics: spell.Unit.NewFocusMetrics(spell.ActionID),
		},
	}
}

func (ec *FocusCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()
	return spell.Unit.CurrentFocus() >= spell.CurCast.Cost
}

func (ec *FocusCost) CostFailureReason(_ *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough focus (Current Focus = %0.03f, Focus Cost = %0.03f)", spell.Unit.CurrentFocus(), spell.CurCast.Cost)
}
func (ec *FocusCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendFocus(sim, spell.CurCast.Cost, ec.ResourceMetrics)
}
func (ec *FocusCost) IssueRefund(sim *Simulation, spell *Spell) {
	if ec.Refund > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddFocus(sim, ec.Refund*spell.CurCast.Cost, ec.RefundMetrics)
	}
}
