package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) registerAscendanceSpell() {

	ascendanceAura := shaman.GetOrRegisterAura(core.Aura{
		Label:    "Ascendance",
		ActionID: core.ActionID{SpellID: 114049},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa := &core.PendingAction{
				NextActionAt: aura.ExpiresAt(),
				Priority:     core.ActionPriorityGCD,
				OnAction:     func(sim *core.Simulation) {},
			}
			sim.AddPendingAction(pa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			//Lava Beam cast gets cancelled if ascendance fades during it
			if (shaman.Hardcast.ActionID.SpellID == 114074) && shaman.Hardcast.Expires > sim.CurrentTime {
				shaman.CancelHardcast(sim)
			}
			if shaman.Spec == proto.Spec_SpecEnhancementShaman {
				shaman.Stormstrike.CD.Set(shaman.Stormblast.CD.ReadyAt())
			}
		},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaBurst,
		Kind:       core.SpellMod_Cooldown_Multiplier,
		FloatValue: -1,
	})

	shaman.Ascendance = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 114049},
		SpellSchool:    core.SpellSchoolPhysical,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskAscendance,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 5.2,
			PercentModifier: 100,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 3,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ascendanceAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: shaman.Ascendance,
		Type:  core.CooldownTypeDPS,
	})
}
