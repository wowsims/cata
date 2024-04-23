package core

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

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
}

func (unit *Unit) EnableFocusBar(maxFocus float64, baseFocusPerSecond float64, isPlayer bool) {
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
	}
}

func (unit *Unit) HasFocusBar() bool {
	return unit.focusBar.unit != nil
}

func (fb *focusBar) CurrentFocus() float64 {
	return fb.currentFocus
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
	if fb.isPlayer {
		return fb.baseFocusPerSecond * fb.getTotalRegenMultiplier()
	} else {
		return fb.baseFocusPerSecond
	}
}

func (fb *focusBar) getTotalRegenMultiplier() float64 {
	return fb.hasteRatingMultiplier * fb.focusRegenMultiplier
}

func (fb *focusBar) AddFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}
	newFocus := min(fb.currentFocus+amount, fb.maxFocus)

	if fb.isPlayer {
		if sim.Log != nil {
			fb.unit.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, fb.currentFocus, newFocus, fb.maxFocus)
		}
		metrics.AddEvent(amount, newFocus-fb.currentFocus)
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

// Gives an immediate partial Focus tick and restarts the tick timer.
func (fb *focusBar) ResetFocusTick(sim *Simulation) {
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
	fb.hasteRatingMultiplier = 1.0 + fb.unit.GetStat(stats.MeleeHaste)/(100*HasteRatingPerHastePercent)
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
	fb.hasteRatingMultiplier = 1.0 + fb.unit.GetStat(stats.MeleeHaste)/(100*HasteRatingPerHastePercent)

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
	Cost float64

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.FocusRefundMetrics if not supplied.
}

type FocusCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newFocusCost(spell *Spell, options FocusCostOptions) *FocusCost {
	spell.DefaultCast.Cost = options.Cost
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.focusRefundMetrics
	}

	return &FocusCost{
		Refund:          options.Refund,
		RefundMetrics:   options.RefundMetrics,
		ResourceMetrics: spell.Unit.NewFocusMetrics(spell.ActionID),
	}
}

func (ec *FocusCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.ApplyCostModifiers(spell.CurCast.Cost)
	return spell.Unit.CurrentFocus() >= spell.CurCast.Cost
}

func (ec *FocusCost) CostFailureReason(_ *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough focus (Current Focus = %0.03f, Focus Cost = %0.03f)", spell.Unit.CurrentFocus(), spell.CurCast.Cost)
}
func (ec *FocusCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendFocus(sim, spell.CurCast.Cost, ec.ResourceMetrics)
}
func (ec *FocusCost) IssueRefund(sim *Simulation, spell *Spell) {
	if ec.Refund > 0 {
		spell.Unit.AddFocus(sim, ec.Refund*spell.CurCast.Cost, ec.RefundMetrics)
	}
}
