package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (warlock *Warlock) registerSoulburn() {

	warlock.SoulBurnAura = warlock.RegisterAura(core.Aura{
		Label:    "Soulburn",
		ActionID: core.ActionID{SpellID: 74434},
		Duration: core.NeverExpires,
	})

	warlock.SoulBurnAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellSummonFelguard | WarlockSpellSummonImp | WarlockSpellSummonSuccubus | WarlockSpellSummonFelhunter | WarlockSpellSoulFire,
		FloatValue: -1.0,
	})

	warlock.SoulBurnAura.AttachSpellMod(core.SpellModConfig{
		Kind:       core.SpellMod_CastTime_Pct,
		ClassMask:  WarlockSpellDrainLife,
		FloatValue: -0.5,
	})

	warlock.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 74434},
		SpellSchool:    core.SpellSchoolShadow,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: WarlockSpellSoulBurn,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 45 * time.Second,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.SoulBurnAura.Activate(sim)
			warlock.RemoveSoulShard()
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.SoulShards > 0
		},
	})
}
