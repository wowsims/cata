package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterShatteringThrowCD() {
	shattDebuffs := warrior.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.ShatteringThrowAura(unit, warrior.UnitIndex)
	})

	ShatteringThrowSpell := warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 64382},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskShatteringThrow | SpellMaskSpecialAttack,
		MaxRange:       30,
		MissileSpeed:   50,

		RageCost: core.RageCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Minute * 5,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if warrior.AutoAttacks.MH().SwingSpeed == warrior.AutoAttacks.OH().SwingSpeed {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, true)
				} else {
					warrior.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+cast.CastTime, false)
				}
				if !warrior.StanceMatches(BattleStance) && warrior.BattleStance.IsReady(sim) {
					warrior.BattleStance.Cast(sim, nil)
				}
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance)
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   warrior.DefaultMeleeCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 12 + 0.5*spell.MeleeAttackPower()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
			if result.Landed() {
				shattDebuffs.Get(target).Activate(sim)
			}
		},

		RelatedAuraArrays: shattDebuffs.ToMap(),
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: ShatteringThrowSpell,
		Type:  core.CooldownTypeDPS,
	})
}
