package hunter

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (hp *HunterPet) ApplySpikedCollar() {
	if hp.hunterOwner.Options.PetSpec != proto.PetSpec_Ferocity {
		return
	}
	critDep := hp.NewDynamicMultiplyStat(stats.Strength, 1.2)
	core.MakePermanent(hp.RegisterAura(core.Aura{
		Label:    "Spiked Collar",
		ActionID: core.ActionID{SpellID: 53184},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.EnableDynamicStatDep(sim, critDep)
			hp.PseudoStats.DamageDealtMultiplier *= 1.1
			hp.MultiplyAttackSpeed(sim, 1.1)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.DisableDynamicStatDep(sim, critDep)
			hp.PseudoStats.DamageDealtMultiplier /= 1.1
			hp.MultiplyAttackSpeed(sim, 1/1.1)
		},
	}))
}

func (hp *HunterPet) ApplyCombatExperience() {
	core.MakePermanent(hp.RegisterAura(core.Aura{
		Label:    "Combat Experience",
		ActionID: core.ActionID{SpellID: 20782},
		Duration: core.NeverExpires,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.DamageDealtMultiplier *= 1.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.DamageDealtMultiplier /= 1.5
		},
	}))
}
