// Implements a generic resource bar that can be used to implement secondary resources
// TODO: Check whether pre-pull OOC resource loss needs to be supported for DemonicFury
package core

import (
	"github.com/wowsims/mop/sim/core/proto"
)

type OnGainCallback func(sim *Simulation, gain int32, realGain int32, actionID ActionID)
type OnSpendCallback func(sim *Simulation, amount int32, actionID ActionID)

type SecondaryResourceBar interface {
	CanSpend(limit int32) bool                                     // Check whether the current resource is available or not
	Spend(sim *Simulation, amount int32, action ActionID)          // Spend the specified amount of resource
	SpendUpTo(sim *Simulation, limit int32, action ActionID) int32 // Spends as much resource as possible up to the speciefied limit; Returns the amount of resource spent
	Gain(sim *Simulation, amount int32, action ActionID)           // Gain the amount specified from the action
	Reset(sim *Simulation)                                         // Resets the current resource bar
	Value() int32                                                  // Returns the current amount of resource
	RegisterOnGain(callback OnGainCallback)                        // Registers a callback that will be called. Gain = amount gained, realGain = actual amount gained due to caps
	RegisterOnSpend(callback OnSpendCallback)                      // Registers a callback that will be called when the resource was spend
}

type SecondaryResourceConfig struct {
	Type    proto.SecondaryResourceType // The type of resource the bar tracks
	Max     int32                       // The maximum amount the bar tracks
	Default int32                       // The default value this bar should be initialized with
}

// Default implementation of SecondaryResourceBar
// Use RegisterSecondaryResourceBar to intantiate the resource bar
type DefaultSecondaryResourceBarImpl struct {
	config  SecondaryResourceConfig
	value   int32
	unit    *Unit
	metrics map[ActionID]*ResourceMetrics
	onGain  []OnGainCallback
	onSpend []OnSpendCallback
}

// CanSpend implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) CanSpend(limit int32) bool {
	return bar.value >= limit
}

// Gain implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Gain(sim *Simulation, amount int32, action ActionID) {
	if amount < 0 {
		panic("Can not gain negative amount")
	}

	oldValue := bar.value
	bar.value = min(bar.value+amount, bar.config.Max)
	amountGained := bar.value - oldValue
	metrics := bar.GetMetric(action)
	metrics.AddEvent(float64(amount), float64(amountGained))
	if sim.Log != nil {
		bar.unit.Log(
			sim,
			"Gained %d %s from %s (%d --> %d) of %d total.",
			amountGained,
			proto.SecondaryResourceType_name[int32(bar.config.Type)],
			action,
			oldValue,
			bar.value,
			bar.config.Max,
		)
	}

	bar.invokeOnGain(sim, amount, amountGained, action)
}

// Reset implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Reset(sim *Simulation) {
	bar.value = 0
	if bar.config.Default > 0 {
		bar.Gain(sim, bar.config.Default, ActionID{SpellID: int32(bar.config.Type)})
	}
}

// Spend implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Spend(sim *Simulation, amount int32, action ActionID) {
	if amount > bar.value {
		panic("Trying to spend more resource than is available.")
	}

	if amount < 0 {
		panic("Trying to spend negative amount.")
	}

	metrics := bar.GetMetric(action)
	if sim.Log != nil {
		bar.unit.Log(
			sim,
			"Spent %d %s from %s (%d --> %d) of %d total.",
			amount,
			proto.SecondaryResourceType_name[int32(bar.config.Type)],
			metrics.ActionID,
			bar.value,
			bar.value-amount,
			bar.config.Max,
		)
	}

	metrics.AddEvent(float64(-amount), float64(-amount))
	bar.invokeOnSpend(sim, amount, action)
	bar.value -= amount
}

// SpendUpTo implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) SpendUpTo(sim *Simulation, limit int32, action ActionID) int32 {
	if bar.value > limit {
		bar.Spend(sim, limit, action)
		return limit
	}

	value := bar.value
	bar.Spend(sim, value, action)
	return value
}

// Value implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Value() int32 {
	return bar.value
}

func (bar *DefaultSecondaryResourceBarImpl) Max() int32 {
	return bar.config.Max
}

func (bar *DefaultSecondaryResourceBarImpl) GetMetric(action ActionID) *ResourceMetrics {
	metric, ok := bar.metrics[action]
	if !ok {
		metric = bar.unit.NewGenericMetric(action)
		bar.metrics[action] = metric
	}

	return metric
}

func (bar *DefaultSecondaryResourceBarImpl) RegisterOnGain(callback OnGainCallback) {
	if callback == nil {
		panic("Can not register nil callback")
	}

	bar.onGain = append(bar.onGain, callback)
}

func (bar *DefaultSecondaryResourceBarImpl) RegisterOnSpend(callback OnSpendCallback) {
	if callback == nil {
		panic("Can not register nil callback")
	}

	bar.onSpend = append(bar.onSpend, callback)
}

func (bar *DefaultSecondaryResourceBarImpl) invokeOnGain(sim *Simulation, gain int32, realGain int32, actionID ActionID) {
	for _, callback := range bar.onGain {
		callback(sim, gain, realGain, actionID)
	}
}

func (bar *DefaultSecondaryResourceBarImpl) invokeOnSpend(sim *Simulation, amount int32, actionID ActionID) {
	for _, callback := range bar.onSpend {
		callback(sim, amount, actionID)
	}
}

func (unit *Unit) NewDefaultSecondaryResourceBar(config SecondaryResourceConfig) *DefaultSecondaryResourceBarImpl {
	if config.Type <= 0 {
		panic("Invalid SecondaryResourceType given.")
	}

	if config.Max <= 0 {
		panic("Invalid maximum resource value given.")
	}

	if config.Default < 0 || config.Default > config.Max {
		panic("Invalid default value given for resource bar")
	}

	return &DefaultSecondaryResourceBarImpl{
		config:  config,
		unit:    unit,
		metrics: make(map[ActionID]*ResourceMetrics),
		onGain:  []OnGainCallback{},
		onSpend: []OnSpendCallback{},
	}
}

func (unit *Unit) RegisterSecondaryResourceBar(config SecondaryResourceBar) {
	if unit.secondaryResourceBar != nil {
		panic("A secondary resource bar has already been registered.")
	}

	unit.secondaryResourceBar = config
}

func (unit *Unit) RegisterNewDefaultSecondaryResourceBar(config SecondaryResourceConfig) SecondaryResourceBar {
	bar := unit.NewDefaultSecondaryResourceBar(config)
	unit.RegisterSecondaryResourceBar(bar)
	return bar
}

func (unit *Unit) GetSecondaryResourceBar() SecondaryResourceBar {
	return unit.secondaryResourceBar
}
