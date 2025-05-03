package hunter

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (hunter *Hunter) registerAspectSpell(spellID int32, dependency *stats.StatDependency, label string) *core.Spell {
	actionID := core.ActionID{SpellID: spellID}

	aura := hunter.GetOrRegisterAura(core.Aura{
		Label:      label,
		ActionID:   actionID,
		Duration:   core.NeverExpires,
		BuildPhase: core.CharacterBuildPhaseBase,
		OnGain: func(a *core.Aura, sim *core.Simulation) {
			a.Unit.EnableDynamicStatDep(sim, dependency)
		},
		OnExpire: func(a *core.Aura, sim *core.Simulation) {
			a.Unit.DisableDynamicStatDep(sim, dependency)
		},
	})
	hunter.applySharedAspectConfig(aura)

	spell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})

	return spell
}
func (hunter *Hunter) registerHawkSpell() {
	hunter.registerAspectSpell(13165, hunter.NewDynamicMultiplyStat(stats.RangedAttackPower, 1.35), "Aspect of the Hawk")
}

func (hunter *Hunter) applySharedAspectConfig(aura *core.Aura) {
	aura.Duration = core.NeverExpires
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
}
