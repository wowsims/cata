package paladin

import "github.com/wowsims/mop/sim/core"

type HolyPowerBar struct {
	resourceBar core.SecondaryResourceBar
	paladin     *Paladin
}

// RegisterOnGain implements core.SecondaryResourceBar.
func (h HolyPowerBar) RegisterOnGain(callback core.OnGainCallback) {
	h.resourceBar.RegisterOnGain(callback)
}

// RegisterOnSpend implements core.SecondaryResourceBar.
func (h HolyPowerBar) RegisterOnSpend(callback core.OnSpendCallback) {
	h.resourceBar.RegisterOnSpend(callback)
}

// CanSpend implements core.SecondaryResourceBar.
func (h HolyPowerBar) CanSpend(limit int32) bool {
	return h.Value() >= limit
}

// Gain implements core.SecondaryResourceBar.
func (h HolyPowerBar) Gain(amount int32, action core.ActionID, sim *core.Simulation) {
	h.resourceBar.Gain(amount, action, sim)
}

// Reset implements core.SecondaryResourceBar.
func (h HolyPowerBar) Reset(sim *core.Simulation) {
	h.resourceBar.Reset(sim)
}

// Spend implements core.SecondaryResourceBar.
func (h HolyPowerBar) Spend(amount int32, action core.ActionID, sim *core.Simulation) {
	if h.paladin.DivinePurposeAura.IsActive() {
		return
	}

	h.resourceBar.Spend(amount, action, sim)
}

// SpendUpTo implements core.SecondaryResourceBar.
func (h HolyPowerBar) SpendUpTo(limit int32, action core.ActionID, sim *core.Simulation) int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.resourceBar.SpendUpTo(limit, action, sim)
}

// Value implements core.SecondaryResourceBar.
func (h HolyPowerBar) Value() int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.resourceBar.Value()
}
