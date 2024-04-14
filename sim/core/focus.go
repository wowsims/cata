package core

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
)

// Time between focus ticks.
const FocusTickDuration = time.Millisecond * 200

type focusBar struct {
	unit *Unit

	maxFocus           float64
	currentFocus       float64
	baseFocusPerSecond float64

	// List of focus levels that might affect APL decisions. E.g:
	// [10, 15, 20, 30, 60, 85]
	focusDecisionThresholds []int

	// Slice with len == maxFocus+1 with each index corresponding to an amount of focus. Looks like this:
	// [0, 0, 0, 0, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, ...]
	// Increments by 1 at each value of focusDecisionThresholds.
	cumulativeFocusDecisionThresholds []int

	nextFocusTick time.Duration
	isPlayer      bool

	// Multiplies focus regen from ticks.
	FocusTickMultiplier float64

	regenMetrics       *ResourceMetrics
	focusRefundMetrics *ResourceMetrics
}

func (unit *Unit) EnableFocusBar(maxFocus float64, baseFocusPerSecond float64, isPlayer bool) {
	unit.SetCurrentPowerBar(FocusBar)

	unit.focusBar = focusBar{
		unit:                unit,
		maxFocus:            max(100, maxFocus),
		FocusTickMultiplier: 1,
		isPlayer:            isPlayer,
		baseFocusPerSecond:  baseFocusPerSecond,
		regenMetrics:        unit.NewFocusMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFocusRegen}),
		focusRefundMetrics:  unit.NewFocusMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

// Computes the focus thresholds.
// Computes the focus thresholds.
func (fb *focusBar) setupFocusThresholds() {
	if fb.unit == nil {
		return
	}
	var focusThresholds []int

	// Focus thresholds from spell costs.
	for _, action := range fb.unit.Rotation.allAPLActions() {
		for _, spell := range action.GetAllSpells() {
			if _, ok := spell.Cost.(*FocusCost); ok {
				focusThresholds = append(focusThresholds, int(math.Ceil(spell.DefaultCast.Cost)))
			}
		}
	}

	// Focus thresholds from conditional comparisons.
	for _, action := range fb.unit.Rotation.allAPLActions() {
		for _, value := range action.GetAllAPLValues() {
			if cmpValue, ok := value.(*APLValueCompare); ok {
				_, lhsIsFocus := cmpValue.lhs.(*APLValueCurrentFocus)
				_, rhsIsFocus := cmpValue.rhs.(*APLValueCurrentFocus)
				if !lhsIsFocus && !rhsIsFocus {
					continue
				}

				lhsConstVal := getConstAPLFloatValue(cmpValue.lhs)
				rhsConstVal := getConstAPLFloatValue(cmpValue.rhs)

				if lhsIsFocus && rhsConstVal != -1 {
					focusThresholds = append(focusThresholds, int(math.Ceil(rhsConstVal)))
				} else if rhsIsFocus && lhsConstVal != -1 {
					focusThresholds = append(focusThresholds, int(math.Ceil(lhsConstVal)))
				}
			}
		}
	}

	slices.SortStableFunc(focusThresholds, func(t1, t2 int) int {
		return t1 - t2
	})

	// Add each unique value to the final thresholds list.
	curVal := 0
	for _, threshold := range focusThresholds {
		if threshold > curVal {
			fb.focusDecisionThresholds = append(fb.focusDecisionThresholds, threshold)
			curVal = threshold
		}
	}

	curFocus := 0
	cumulativeVal := 0
	fb.cumulativeFocusDecisionThresholds = make([]int, int(fb.maxFocus)+1)
	for _, threshold := range fb.focusDecisionThresholds {
		for curFocus < threshold {
			fb.cumulativeFocusDecisionThresholds[curFocus] = cumulativeVal
			curFocus++
		}
		cumulativeVal++
	}
	for curFocus < len(fb.cumulativeFocusDecisionThresholds) {
		fb.cumulativeFocusDecisionThresholds[curFocus] = cumulativeVal
		curFocus++
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

// Returns the rate of focus regen per second from melee haste.
// Todo: Verify that this is actually how it works. Check simc code below
// player_t::focus_regen_per_second =========================================
// double player_t::focus_regen_per_second() const
//
//	{
//	  double r = base_focus_regen_per_second * ( 1.0 / composite_attack_haste() );
//	  return r;
//	}
func (fb *focusBar) FocusRegenPerTick() float64 {
	ticksPerSecond := float64(time.Second) / float64(FocusTickDuration)
	if fb.isPlayer {
		hastePercent := fb.unit.RangedSwingSpeed()
		tick := fb.baseFocusPerSecond * hastePercent / ticksPerSecond
		return tick
	} else {
		tick := fb.baseFocusPerSecond / ticksPerSecond
		return tick
	}
}

func (fb *focusBar) FocusRegenPerSecond() float64 {
	if fb.isPlayer {
		hastePercent := fb.unit.RangedSwingSpeed()
		return fb.baseFocusPerSecond * hastePercent
	} else {
		return fb.baseFocusPerSecond
	}
}

func (fb *focusBar) onFocusGain(sim *Simulation, crossedThreshold bool) {
	if sim.CurrentTime < 0 {
		return
	}

	if !sim.Options.Interactive && crossedThreshold {
		fb.unit.Rotation.DoNextAction(sim)
	}
}

func (fb *focusBar) addFocusInternal(sim *Simulation, amount float64, metrics *ResourceMetrics) bool {
	if amount < 0 {
		panic("Trying to add negative focus!")
	}
	newFocus := min(fb.currentFocus+amount, fb.maxFocus)

	if fb.isPlayer {
		if sim.Log != nil {
			fb.unit.Log(sim, "Gained %0.3f focus from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, fb.currentFocus, newFocus)
		}
		metrics.AddEvent(amount, newFocus-fb.currentFocus)
	}

	crossedThreshold := fb.cumulativeFocusDecisionThresholds == nil || fb.cumulativeFocusDecisionThresholds[int(fb.currentFocus)] != fb.cumulativeFocusDecisionThresholds[int(newFocus)]
	fb.currentFocus = newFocus

	return crossedThreshold
}
func (fb *focusBar) AddFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	crossedThreshold := fb.addFocusInternal(sim, amount, metrics)
	fb.onFocusGain(sim, crossedThreshold)
}

func (fb *focusBar) SpendFocus(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative focus!")
	}

	newFocus := fb.currentFocus - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		fb.unit.Log(sim, "Spent %0.3f focus from %s (%0.3f --> %0.3f).", amount, metrics.ActionID, fb.currentFocus, newFocus)
	}

	fb.currentFocus = newFocus
}

func (fb *focusBar) RunTask(sim *Simulation) time.Duration {
	if sim.CurrentTime < fb.nextFocusTick {
		return fb.nextFocusTick
	}
	crossedThreshold := fb.addFocusInternal(sim, fb.FocusRegenPerTick(), fb.regenMetrics)
	fb.onFocusGain(sim, crossedThreshold)

	fb.nextFocusTick = sim.CurrentTime + FocusTickDuration
	return fb.nextFocusTick
}

func (fb *focusBar) reset(sim *Simulation) {
	if fb.unit == nil {
		return
	}

	fb.currentFocus = fb.maxFocus
	if fb.unit.Type != PetUnit {
		fb.enable(sim, sim.Environment.PrepullStartTime())
	}
}

func (fb *focusBar) enable(sim *Simulation, startAt time.Duration) {
	sim.AddTask(fb)
	fb.nextFocusTick = startAt + time.Duration(sim.RandomFloat("Focus Tick")*float64(FocusTickDuration))
	sim.RescheduleTask(fb.nextFocusTick)

	if fb.cumulativeFocusDecisionThresholds != nil && sim.Log != nil {
		fb.unit.Log(sim, "[DEBUG] APL Focus decision thresholds: %v", fb.focusDecisionThresholds)
	}
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
