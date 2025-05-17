package destruction

import "github.com/wowsims/mop/sim/core"

func (destruction DestructionWarlock) ApplyEmbersHandler() {

	core.MakePermanent(destruction.RegisterAura(core.Aura{
		Label: "Burning Embers: Driver",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell.ClassSpellMask & SpellMaskCinderGenerator) == 0 {
				return
			}

			if result.DidCrit() {
				destruction.BurningEmbers.Gain(2, spell.ActionID, sim)
				return
			}

			destruction.BurningEmbers.Gain(1, spell.ActionID, sim)
		},
	}))
}
