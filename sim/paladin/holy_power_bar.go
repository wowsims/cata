package paladin

import (
	"github.com/wowsims/mop/sim/core"
)

type HolyPowerBar struct {
	*core.DefaultSecondaryResourceBarImpl
	paladin *Paladin
}

// Spend implements core.SecondaryResourceBar.
func (h HolyPowerBar) Spend(sim *core.Simulation, amount int32, action core.ActionID) {
	if h.paladin.DivinePurposeAura.IsActive() {
		return
	}

	h.DefaultSecondaryResourceBarImpl.Spend(sim, amount, action)
}

// SpendUpTo implements core.SecondaryResourceBar.
func (h HolyPowerBar) SpendUpTo(sim *core.Simulation, limit int32, action core.ActionID) int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return h.DefaultSecondaryResourceBarImpl.SpendUpTo(sim, limit, action)
}

// Value implements core.SecondaryResourceBar.
func (h HolyPowerBar) Value() int32 {
	if h.paladin.DivinePurposeAura.IsActive() {
		return 5
	}

	return h.DefaultSecondaryResourceBarImpl.Value()
}

func (h HolyPowerBar) CanSpend(amount int32) bool {
	if h.paladin.DivinePurposeAura.IsActive() {
		return true
	}

	return h.DefaultSecondaryResourceBarImpl.CanSpend(amount)
}
