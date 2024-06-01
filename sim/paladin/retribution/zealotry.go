package retribution

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/paladin"
	"time"
)

func (retPaladin *RetributionPaladin) RegisterZealotry() {
	if !retPaladin.Talents.Zealotry {
		return
	}

	actionId := core.ActionID{SpellID: 85696}

	retPaladin.ZealotryAura = retPaladin.RegisterAura(core.Aura{
		Label:    "Zealotry",
		ActionID: actionId,
		Duration: 20 * time.Second,
	})

	retPaladin.Zealotry = retPaladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskZealotry,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return retPaladin.GetHolyPowerValue() >= 3
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			holyPower := retPaladin.GetHolyPowerValue()

			if holyPower == 0 {
				return
			}

			retPaladin.ZealotryAura.Activate(sim)
		},
	})
}
