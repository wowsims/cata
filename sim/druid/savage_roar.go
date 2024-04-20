package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) getSavageRoarMultiplier() float64 {
	return 1.8 + core.TernaryFloat64(druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfSavageRoar), 0.05, 0)
}

func (druid *Druid) registerSavageRoarSpell() {
	actionID := core.ActionID{SpellID: 52610}

	srm := druid.getSavageRoarMultiplier()
	durBonus := core.DurationFromSeconds(4.0 * float64(druid.Talents.EndlessCarnage))
	druid.SavageRoarDurationTable = [6]time.Duration{
		0,
		durBonus + time.Second*(9+5),
		durBonus + time.Second*(9+10),
		durBonus + time.Second*(9+15),
		durBonus + time.Second*(9+20),
		durBonus + time.Second*(9+25),
	}

	druid.SavageRoarAura = druid.RegisterAura(core.Aura{
		Label:    "Savage Roar Aura",
		ActionID: actionID,
		Duration: 9,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.MHAutoSpell.DamageMultiplier *= srm
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if druid.InForm(Cat) {
				druid.MHAutoSpell.DamageMultiplier /= srm
			}
		},
	})

	srSpell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		EnergyCost: core.EnergyCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.SavageRoarAura.Duration = druid.SavageRoarDurationTable[druid.ComboPoints()]
			druid.SavageRoarAura.Activate(sim)
			druid.SpendComboPoints(sim, spell.ComboPointMetrics())
		},
	})

	druid.SavageRoar = srSpell
}

func (druid *Druid) CurrentSavageRoarCost() float64 {
	return druid.SavageRoar.ApplyCostModifiers(druid.SavageRoar.DefaultCast.Cost)
}
