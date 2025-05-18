package warlock

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warlock *Warlock) registerArchimondesDarkness() {
	if !warlock.Talents.ArchimondesDarkness {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_ModCharges_Flat,
		IntValue:  2,
		ClassMask: WarlockSpellDarkSoulInsanity,
	})

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: -time.Second * 100,
		ClassMask: WarlockSpellDarkSoulInsanity,
	})
}
