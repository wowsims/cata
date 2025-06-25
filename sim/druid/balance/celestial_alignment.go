package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/druid"
)

func (moonkin *BalanceDruid) registerCelestialAlignmentSpell() {
	actionID := core.ActionID{SpellID: 112071}

	celestialAlignmentAura := moonkin.RegisterAura(core.Aura{
		Label:    "Celestial Alignment",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(_ *core.Aura, sim *core.Simulation) {
			moonkin.SuspendEclipseBar()

			// Activate both eclipse damage bonuses
			moonkin.ActivateEclipse(LunarEclipse, sim)
			moonkin.ActivateEclipse(SolarEclipse, sim)
		},
		OnExpire: func(_ *core.Aura, sim *core.Simulation) {
			moonkin.DeactivateEclipse(LunarEclipse, sim)
			moonkin.DeactivateEclipse(SolarEclipse, sim)

			// Restore previous eclipse gain mask
			moonkin.RestoreEclipseBar()
		},
		OnCastComplete: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask == druid.DruidSpellMoonfire {
				moonkin.Sunfire.Dot(spell.Unit.CurrentTarget).Apply(sim)
			}

			if spell.ClassSpellMask == druid.DruidSpellSunfire {
				moonkin.Moonfire.Dot(spell.Unit.CurrentTarget).Apply(sim)
			}
		},
	})

	moonkin.CelestialAlignment = moonkin.RegisterSpell(druid.Humanoid|druid.Moonkin, core.SpellConfig{
		ActionID:        actionID,
		SpellSchool:     core.SpellSchoolArcane,
		Flags:           core.SpellFlagAPL,
		RelatedSelfBuff: celestialAlignmentAura,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    moonkin.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
	})

	moonkin.AddMajorCooldown(core.MajorCooldown{
		Spell: moonkin.CelestialAlignment.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
