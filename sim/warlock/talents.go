package warlock

import "github.com/wowsims/mop/sim/core"

func (warlock *Warlock) registerArchimondesDarkness() {
	if !warlock.Talents.ArchimondesDarkness {
		return
	}

	warlock.AddStaticMod(core.SpellModConfig{
		Kind:      core.SpellMod_ModCharges_Flat,
		IntValue:  1,
		ClassMask: WarlockSpellDarkSoulInsanity,
	})
}
