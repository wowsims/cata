package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) registerInquisition() {
	actionId := core.ActionID{SpellID: 84963}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)
	inquisitionDuration := time.Millisecond * time.Duration(4000*[]float64{1, 1.66, 2.33, 3.0}[paladin.Talents.InquiryOfFaith])

	paladin.InquisitionAura = paladin.RegisterAura(core.Aura{
		Label:    "Inquisition" + paladin.Label,
		ActionID: actionId,
		Duration: inquisitionDuration,
	})

	paladin.InquisitionAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.3,
		School:     core.SpellSchoolHoly,
	})

	// Inquisition self-buff.
	paladin.Inquisition = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		Flags:          core.SpellFlagAPL,
		ProcMask:       core.ProcMaskEmpty,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: SpellMaskInquisition,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.GetHolyPowerValue() > 0
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			holyPower := paladin.GetHolyPowerValue()

			if paladin.HasT11Ret4pc {
				holyPower += 1
			}

			paladin.InquisitionAura.Duration = inquisitionDuration * time.Duration(holyPower)
			paladin.SpendHolyPower(sim, hpMetrics)
			paladin.InquisitionAura.Activate(sim)
		},
	})
}
