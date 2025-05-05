// Implements a generic resource bar that can be used to implement secondary resources
// TODO: Check whether pre-pull OOC resource loss needs to be supported for DemonicFury
package core

import "math"

type SecondaryResourceType int32

const (
	SoulShards    SecondaryResourceType = 117198
	HolyPower     SecondaryResourceType = 138248
	Maelstrom     SecondaryResourceType = 53817
	Chi           SecondaryResourceType = 97272
	ArcaneCharges SecondaryResourceType = 36032
	ShadowOrbs    SecondaryResourceType = 95740
	BurningEmbers SecondaryResourceType = 108647
	DemonicFury   SecondaryResourceType = 104315
)

type SecondaryResourceBar interface {
	CanSpend(limit float64) bool                                       // Check whether the current resource is available or not
	Spend(amount float64, action ActionID, sim *Simulation)            // Spend the specified amount of resource
	SpendUpTo(limit float64, action ActionID, sim *Simulation) float64 // Spends as much resource as possible up to the speciefied limit; Returns the amount of resource spent
	Gain(amount float64, action ActionID, sim *Simulation)             // Gain the amount specified from the action
	Reset(sim *Simulation)                                             // Resets the current resource bar
	Value() float64                                                    // Returns the current amount of resource
}

type SecondaryResourceConfig struct {
	Type    SecondaryResourceType // The type of resource the bar tracks
	Max     float64               // The maximum amount the bar tracks
	Default float64               // The default value this bar should be initialized with
}

// Default implementation of SecondaryResourceBar
// Use RegisterSecondaryResourceBar to intantiate the resource bar
type DefaultSecondaryResourceBarImpl struct {
	config  SecondaryResourceConfig
	value   float64
	unit    *Unit
	metrics map[ActionID]*ResourceMetrics
}

// CanSpend implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) CanSpend(limit float64) bool {
	return bar.value >= limit
}

// Gain implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Gain(amount float64, action ActionID, sim *Simulation) {
	oldValue := bar.value
	bar.value = min(bar.value+amount, bar.config.Max)
	amountGained := bar.value - oldValue
	metrics := bar.GetMetric(action)
	metrics.AddEvent(amount, amountGained)
	if sim.Log != nil {
		bar.unit.Log(
			sim,
			"Gained %0.0f generic resource from %s (%0.0f --> %0.0f) of %0.0f total.",
			amountGained,
			action,
			oldValue,
			bar.value,
			bar.config.Max,
		)
	}
}

// Reset implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Reset(sim *Simulation) {
	bar.value = 0
	if bar.config.Default > 0 {
		bar.Gain(bar.config.Default, ActionID{SpellID: int32(bar.config.Type)}, sim)
	}
}

// Spend implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Spend(amount float64, action ActionID, sim *Simulation) {
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
			"Spent %0.0f generic resource from %s (%0.0f --> %0.0f) of %0.0f total.",
			amount,
			metrics.ActionID,
			bar.value,
			bar.value-amount,
			bar.config.Max,
		)
	}

	metrics.AddEvent(-amount, -amount)
	bar.value -= amount
}

// SpendUpTo implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) SpendUpTo(limit float64, action ActionID, sim *Simulation) float64 {
	if bar.value > limit {
		bar.Spend(limit, action, sim)
		return limit
	}

	max := math.Floor(bar.value)
	bar.Spend(max, action, sim)
	return max
}

// Value implements SecondaryResourceBar.
func (bar *DefaultSecondaryResourceBarImpl) Value() float64 {
	return bar.value
}

func (bar *DefaultSecondaryResourceBarImpl) GetMetric(action ActionID) *ResourceMetrics {
	metric, ok := bar.metrics[action]
	if !ok {
		metric = bar.unit.NewGenericMetric(action)
		bar.metrics[action] = metric
	}

	return metric
}

func (unit *Unit) RegisterSecondaryResourceBar(config SecondaryResourceConfig) SecondaryResourceBar {
	if config.Type <= 0 {
		panic("Invalid SecondaryResourceType given.")
	}

	if config.Max <= 0 {
		panic("Invalid maximum resource value given.")
	}

	if config.Default < 0 || config.Default > config.Max {
		panic("Invalid default value given for resource bar")
	}

	if unit.SecondaryResourceBar != nil {
		panic("A secondary resource bar has already been registered.")
	}

	if unit.Env.State == Finalized {
		panic("Can not add secondary resource bar after unit has been finalized")
	}

	unit.SecondaryResourceBar = &DefaultSecondaryResourceBarImpl{config: config, unit: unit, metrics: make(map[ActionID]*ResourceMetrics)}
	return unit.SecondaryResourceBar
}
