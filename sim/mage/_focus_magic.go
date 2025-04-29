package mage

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (mage *Mage) applyFocusMagic() {
	if !mage.Talents.FocusMagic {
		return
	}

	// This is used only for the individual sim.
	if mage.Party.Raid.Size() == 1 {
		if mage.ArcaneOptions.FocusMagicPercentUptime > 0 {
			selfAura, _ := core.FocusMagicAura(&mage.Unit, nil)
			core.ApplyFixedUptimeAura(selfAura, float64(mage.ArcaneOptions.FocusMagicPercentUptime)/100, time.Second*10, 1)
		}
		return
	}

	focusMagicTarget := mage.GetUnit(mage.ArcaneOptions.FocusMagicTarget)
	if focusMagicTarget == nil {
		return
	} else if focusMagicTarget == &mage.Unit {
		// When self is selected, give permanent self buff.
		selfAura, _ := core.FocusMagicAura(&mage.Unit, nil)
		core.MakePermanent(selfAura)
		return
	}

	core.FocusMagicAura(&mage.Unit, focusMagicTarget)
}
