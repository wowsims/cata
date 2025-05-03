package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerHuntersMarkSpell() {
	actionID := core.ActionID{SpellID: 1130}
	rangedMult := 1.05
	hunter.HuntersMarkAura = hunter.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "HuntersMark-" + hunter.Label,
			ActionID: actionID,
			Duration: 5 * time.Minute,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].RangedDamageTakenMulitplier *= rangedMult
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				hunter.AttackTables[aura.Unit.UnitIndex].RangedDamageTakenMulitplier /= rangedMult
			},
		})
	})

	config := core.SpellConfig{
		ActionID: actionID,
		ProcMask: core.ProcMaskEmpty,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aura := range hunter.HuntersMarkAura {
				if aura.IsActive() {
					aura.Deactivate(sim)
				}
			}
			// Activating Hunters Mark for the new target
			hunter.HuntersMarkAura.Get(target).Activate(sim)
		},
	}

	hunter.HuntersMarkSpell = hunter.RegisterSpell(config)

	config.Cast = core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
		IgnoreHaste: true,
	}
	config.Flags = core.SpellFlagAPL

	hunter.RegisterSpell(config)
}
