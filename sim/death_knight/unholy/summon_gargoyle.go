package unholy

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/death_knight"
)

func (uhdk *UnholyDeathKnight) registerSummonGargoyleSpell() {
	spell := uhdk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 49206},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: death_knight.DeathKnightSpellSummonGargoyle,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    uhdk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)

			uhdk.Gargoyle.SetExpireTime(sim.CurrentTime + time.Second*30)
			uhdk.Gargoyle.EnableWithTimeout(sim, uhdk.Gargoyle, time.Second*30)
			// Start casting after a 2.5s delay to simulate the summon animation
			uhdk.Gargoyle.SetGCDTimer(sim, sim.CurrentTime+time.Millisecond*2500)
		},

		RelatedSelfBuff: uhdk.RegisterAura(core.Aura{
			Label:    "Summon Gargoyle",
			ActionID: core.ActionID{SpellID: 49206},
			Duration: time.Second * 30,
		}),
	})

	uhdk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
