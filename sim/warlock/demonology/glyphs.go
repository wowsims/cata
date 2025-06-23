package demonology

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (demo *DemonologyWarlock) registerImpSwarm() {
	if !demo.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImpSwarm) {
		return
	}

	demo.ImpSwarm = demo.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 104316},
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    demo.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for range 5 {
				demo.SpawnImp(sim)
			}

			spell.CD.Set(sim.CurrentTime + time.Duration(float64(time.Minute*2)/demo.TotalSpellHasteMultiplier()))
		},
	})

	demo.AddMajorCooldown(core.MajorCooldown{
		Spell:    demo.ImpSwarm,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeDPS,
	})
}
