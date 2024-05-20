package hunter

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (hunter *Hunter) registerAspectOfTheHawkSpell() {
	actionID := core.ActionID{SpellID: 13165}
	ap := 2700.0

	if hunter.Talents.OneWithNature > 0 {
		ap *= 1 + (float64(hunter.Talents.OneWithNature) * 0.1)
	}

	hunter.AspectOfTheHawkAura = hunter.NewTemporaryStatsAuraWrapped(
		"Aspect of the Hawk",
		actionID,
		stats.Stats{
			stats.RangedAttackPower: ap,
		},
		core.NeverExpires,
		func(aura *core.Aura) {
		})
	hunter.applySharedAspectConfig(true, hunter.AspectOfTheHawkAura)

	hunter.AspectOfTheHawk = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheHawkAura.Activate(sim)
		},
	})
}
func (hunter *Hunter) registerAspectOfTheFoxSpell() {
	actionID := core.ActionID{SpellID: 82661}
	restoreFocus := 2.0
	focusMetric := hunter.NewFocusMetrics(actionID)

	if hunter.Talents.OneWithNature > 0 {
		restoreFocus += 1 * float64(hunter.Talents.OneWithNature)
	}

	foxMod := hunter.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_AllowCastWhileMoving,
		ClassMask: HunterSpellCobraShot | HunterSpellSteadyShot,
	})
	hunter.AspectOfTheFoxAura = core.MakePermanent(hunter.GetOrRegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Aspect of the Fox",
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			foxMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			foxMod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if aura.IsActive() {
				hunter.AddFocus(sim, restoreFocus, focusMetric)
			}
		},
	}))

	hunter.applySharedAspectConfig(true, hunter.AspectOfTheFoxAura)

	hunter.AspectOfTheFox = hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.AspectOfTheFoxAura.Activate(sim)
		},
	})
}

func (hunter *Hunter) applySharedAspectConfig(isHawk bool, aura *core.Aura) {
	aura.Duration = core.NeverExpires
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})
}
