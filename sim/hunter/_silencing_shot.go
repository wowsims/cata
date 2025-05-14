package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerSilencingShotSpell() {
	if !hunter.Talents.SilencingShot {
		return
	}

	hunter.SilencingShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34490},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MinRange:    5,
		MaxRange:    40,
		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Silencing Shot does nothing in wotlk for damage except maybe restore 10 focus
			if hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfSilencingShot) {
				focusMetics := hunter.NewFocusMetrics(core.ActionID{SpellID: 34490})
				hunter.AddFocus(sim, 10, focusMetics)
			}
		},
	})
}
