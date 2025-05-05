package paladin

import "github.com/wowsims/mop/sim/core"

type HolyPowerBar struct {
	resourceBar core.SecondaryResourceBar
	paladin     *Paladin
}

// CanSpend implements core.SecondaryResourceBar.
func (h HolyPowerBar) CanSpend(limit float64) bool {
	return h.Value() >= limit
}

// Gain implements core.SecondaryResourceBar.
func (h HolyPowerBar) Gain(amount float64, action core.ActionID, sim *core.Simulation) {
	h.resourceBar.Gain(amount, action, sim)
}

// Reset implements core.SecondaryResourceBar.
func (h HolyPowerBar) Reset(sim *core.Simulation) {
	h.resourceBar.Reset(sim)
}

// Spend implements core.SecondaryResourceBar.
func (h HolyPowerBar) Spend(amount float64, action core.ActionID, sim *core.Simulation) {
	if h.paladin.DivinePurposeAura.IsActive() {
		return
	}

	h.resourceBar.Spend(amount, action, sim)
}

// SpendUpTo implements core.SecondaryResourceBar.
func (h HolyPowerBar) SpendUpTo(limit float64, action core.ActionID, sim *core.Simulation) float64 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.resourceBar.SpendUpTo(limit, action, sim)
}

// Value implements core.SecondaryResourceBar.
func (h HolyPowerBar) Value() float64 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.resourceBar.Value()
}
