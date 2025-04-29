package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (priest *Priest) registerDispersionSpell() {
	if !priest.Talents.Dispersion {
		return
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 47585})
	var pa *core.PendingAction
	dispersionAura := priest.GetOrRegisterAura(core.Aura{
		Label:    "Dispersion",
		ActionID: core.ActionID{SpellID: 47585},
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second,
				NumTicks: 6,
				OnAction: func(sim *core.Simulation) {
					manaGain := priest.MaxMana() * 0.06
					priest.AddMana(sim, manaGain, manaMetric)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
		},
	})

	spell := priest.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 47585},
		ProcMask:       core.ProcMaskEmpty,
		SpellSchool:    core.SpellSchoolShadow,
		ClassSpellMask: PriestSpellDispersion,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 120,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dispersionAura.Activate(sim)
			priest.WaitUntil(sim, dispersionAura.ExpiresAt())
		},
	})

	priest.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: 1,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.01
		},
	})
}
