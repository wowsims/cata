package guardian

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (bear *GuardianDruid) registerSavageDefenseSpell() {
	bear.SavageDefenseAura = bear.RegisterAura(core.Aura{
		Label:    "Savage Defense",
		ActionID: core.ActionID{SpellID: 132402},
		Duration: time.Second * 6,

		OnGain: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.BaseDodgeChance += 0.45
		},

		OnExpire: func(aura *core.Aura, _ *core.Simulation) {
			aura.Unit.PseudoStats.BaseDodgeChance -= 0.45
		},
	})

	bear.SavageDefense = bear.RegisterSpell(druid.Bear, core.SpellConfig{
		ActionID:        core.ActionID{SpellID: 62606},
		SpellSchool:     core.SpellSchoolNature,
		ProcMask:        core.ProcMaskEmpty,
		Flags:           core.SpellFlagAPL,
		Charges:         3,
		RechargeTime:    time.Second * 9,
		RelatedSelfBuff: bear.SavageDefenseAura,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    bear.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},

		RageCost: core.RageCostOptions{
			Cost: 60,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			bear.SavageDefenseAura.Activate(sim)
		},
	})
}
