package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (sinRogue *AssassinationRogue) registerVendetta() {
	if !sinRogue.Talents.Vendetta {
		return
	}

	actionID := core.ActionID{SpellID: 79140}

	vendettaAura := sinRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Vendetta",
			ActionID: actionID,
			Duration: 20 * time.Second,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				sinRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				sinRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= 1.2
			},
		})
	})

	sinRogue.Vendetta = sinRogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagAPL | core.SpellFlagMCD,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    sinRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura := vendettaAura.Get(target)
			aura.Activate(sim)
		},
	})

	sinRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    sinRogue.Vendetta,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return sinRogue.GCD.IsReady(sim) && sinRogue.ComboPoints() >= 4
		},
	})
}
