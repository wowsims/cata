package arms

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ArmsWarrior) RegisterSweepingStrikes() {
	if !war.Talents.SweepingStrikes {
		return
	}

	actionID := core.ActionID{SpellID: 12328}

	var curDmg float64
	ssHit := war.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty, // No proc mask, so it won't proc itself.
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreAttackerModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, curDmg, spell.OutcomeAlwaysHit)
		},
	})

	ssAura := war.RegisterAura(core.Aura{
		Label:    "Sweeping Strikes",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Damage <= 0 || !spell.ProcMask.Matches(core.ProcMaskMelee) || war.Env.GetNumTargets() < 2 {
				return
			}

			if (spell == war.Execute && !sim.IsExecutePhase20()) || spell == war.Whirlwind {
				curDmg = spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				curDmg = result.Damage
			}

			// Undo armor reduction to get the raw damage value.
			curDmg /= result.ResistanceMultiplier

			ssHit.Cast(sim, war.Env.NextTargetUnit(result.Target))
			ssHit.SpellMetrics[result.Target.UnitIndex].Casts--
		},
	})

	ssCD := war.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: warrior.SpellMaskSweepingStrikes,

		RageCost: core.RageCostOptions{
			Cost: 30,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return war.StanceMatches(warrior.BattleStance | warrior.BerserkerStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			ssAura.Activate(sim)
		},
	})

	war.AddMajorCooldown(core.MajorCooldown{
		Spell: ssCD,
		Type:  core.CooldownTypeDPS,
	})
}
