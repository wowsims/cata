package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (paladin *Paladin) registerInquisition() {
	if paladin.Talents.InquiryOfFaith == 0 {
		return
	}

	actionId := core.ActionID{SpellID: 84963}
	hpMetrics := paladin.NewHolyPowerMetrics(actionId)
	inquisitionDuration := 4 * time.Second * time.Duration([]float64{0, 1.66, 2.33, 3.0}[paladin.Talents.InquiryOfFaith])

	hasT11_4pc := paladin.HasSetBonus(ItemSetReinforcedSapphiriumBattleplate, 4)

	inquisitionMod := paladin.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.3,
		ClassMask:  SpellMaskModifiedByInquisition,
		School:     core.SpellSchoolHoly,
	})

	paladin.InquisitionAura = paladin.RegisterAura(core.Aura{
		Label:    "Inquisition",
		ActionID: actionId,
		Duration: inquisitionDuration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			inquisitionMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			inquisitionMod.Deactivate()
		},
	})

	// Inquisition self-buff.
	paladin.Inquisition = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		Flags:          core.SpellFlagAPL,
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

			if holyPower == 0 {
				return
			}

			if hasT11_4pc {
				holyPower += 1
			}

			paladin.InquisitionAura.Duration = inquisitionDuration * time.Duration(holyPower)
			paladin.InquisitionAura.Activate(sim)
			paladin.SpendHolyPower(sim, hpMetrics)
		},
	})
}
