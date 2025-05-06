package paladin

import "github.com/wowsims/mop/sim/core"

type HolyPowerBar struct {
	*core.DefaultSecondaryResourceBarImpl
	paladin *Paladin
}

// Spend implements core.SecondaryResourceBar.
func (h HolyPowerBar) Spend(amount int32, action core.ActionID, sim *core.Simulation) {
	if h.paladin.DivinePurposeAura.IsActive() {
		return
	}

	h.DefaultSecondaryResourceBarImpl.Spend(amount, action, sim)
}

// SpendUpTo implements core.SecondaryResourceBar.
func (h HolyPowerBar) SpendUpTo(limit int32, action core.ActionID, sim *core.Simulation) int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.DefaultSecondaryResourceBarImpl.SpendUpTo(limit, action, sim)
}

// Value implements core.SecondaryResourceBar.
func (h HolyPowerBar) Value() int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.DefaultSecondaryResourceBarImpl.Value()
}
