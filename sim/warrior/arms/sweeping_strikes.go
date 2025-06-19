package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) registerSweepingStrikes() {
	actionID := core.ActionID{SpellID: 12328}

	var copyDamage float64
	hitSpell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		ClassSpellMask: warrior.SpellMaskSweepingStrikesHit,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

		DamageMultiplier: 0.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, copyDamage, spell.OutcomeAlwaysHit)
		},
	})

	war.SweepingStrikesAura = core.MakeProcTriggerAura(&war.Unit, core.ProcTrigger{
		Name:     "Sweeping Strikes",
		ActionID: actionID,
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeHit,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if war.Env.GetNumTargets() < 2 || result.PreOutcomeDamage <= 0 || spell.Matches(warrior.SpellMaskSweepingStrikesHit) {
				return
			}

			copyDamage = result.PreOutcomeDamage

			hitSpell.Cast(sim, war.Env.NextTargetUnit(result.Target))
		},
	})

	spell := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: warrior.SpellMaskSweepingStrikes,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			war.SweepingStrikesAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.Env.GetNumTargets() > 1
		},
	})
}
