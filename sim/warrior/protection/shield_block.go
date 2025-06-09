package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/warrior"
)

func (war *ProtectionWarrior) registerShieldBlock() {
	actionId := core.ActionID{SpellID: 2565}

	// extra avoidance to crit block effect seems to be based on basic level+3 target
	atkTableAttacker := &core.Unit{Level: war.Level + 3, Type: core.EnemyUnit}
	atkTable := core.NewAttackTable(atkTableAttacker, &war.Unit)

	extraAvoidance := 0.0
	war.ShieldBlockAura = war.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionId,
		Duration: time.Second * 6,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			avoidance := war.GetTotalAvoidanceChance(atkTable)
			if avoidance > core.CombatTableCoverageCap {
				extraAvoidance = avoidance - core.CombatTableCoverageCap
				war.CriticalBlockChance[1] += extraAvoidance
			} else {
				extraAvoidance = 0.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if extraAvoidance > 0.0 {
				war.CriticalBlockChance[1] -= extraAvoidance
			}
		},
	}).AttachStatBuff(stats.BlockPercent, 100)

	war.RegisterSpell(core.SpellConfig{
		ActionID:       actionId,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: warrior.SpellMaskShieldBlock,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,

		Charges:      2,
		RechargeTime: 9 * time.Second,

		RageCost: core.RageCostOptions{
			Cost: 60,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    war.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			war.ShieldBlockAura.Activate(sim)
		},
	})
}
