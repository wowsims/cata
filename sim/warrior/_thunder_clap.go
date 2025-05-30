package warrior

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (warrior *Warrior) RegisterThunderClapSpell() {
	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.ThunderClapAura(target)
	})

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6343},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagAoE | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskThunderClap | SpellMaskSpecialAttack,

		RageCost: core.RageCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		DamageMultiplier: 1.0,
		ThreatMultiplier: 1.85,
		CritMultiplier:   warrior.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 303.0 + 0.228*spell.MeleeAttackPower()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
				if result.Landed() {
					warrior.ThunderClapAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuraArrays: warrior.ThunderClapAuras.ToMap(),
	})
}
