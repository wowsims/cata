package affliction

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerSoulburn() {

	castTimeMod := affliction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  warlock.WarlockSpellSummonImp | warlock.WarlockSpellSummonSuccubus | warlock.WarlockSpellSummonFelhunter | warlock.WarlockSpellSoulFire,
		FloatValue: -1.0,
	})

	drainLifeCastMod := affliction.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  warlock.WarlockSpellDrainLife,
		FloatValue: -0.5,
	})

	affliction.SoulBurnAura = affliction.RegisterAura(core.Aura{
		Label:    "Soulburn",
		ActionID: core.ActionID{SpellID: 74434},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Activate()
			drainLifeCastMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			castTimeMod.Deactivate()
			drainLifeCastMod.Deactivate()
		},
	})

	affliction.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 74434},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellSoulBurn,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    affliction.NewTimer(),
				Duration: time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			affliction.SoulBurnAura.Activate(sim)

			// if we cast this >= 20 sec pre pull we will regenerate the shard
			if sim.CurrentTime > -20*time.Second {
				affliction.SoulShards.Spend(sim, 1, spell.ActionID)
			}
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return affliction.SoulShards.CanSpend(1)
		},
	})
}
