package warrior

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warrior *Warrior) RegisterShieldBlockCD() {
	actionID := core.ActionID{SpellID: 2565}

	T12_4P_bonus := warrior.RegisterAura(core.Aura{
		Label:    "T12 4P Bonus",
		ActionID: core.ActionID{SpellID: 99242},
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Parry, 6*core.ParryRatingPerParryChance)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Parry, -6*core.ParryRatingPerParryChance)
		},
	})

	// extra avoidance to crit block effect seems to be based on basic level+3 target
	atkTableAttacker := &core.Unit{Level: warrior.Level + 3, Type: core.EnemyUnit}
	atkTable := core.NewAttackTable(atkTableAttacker, &warrior.Unit)

	hasT12_4P_bonus := warrior.HasSetBonus(ItemSetMoltenGiantBattleplate, 4)

	extraAvoidance := 0.0
	warrior.ShieldBlockAura = warrior.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 25*core.BlockRatingPerBlockChance)

			avoidance := warrior.GetTotalAvoidanceChance(atkTable)
			if avoidance > core.CombatTableCoverageCap {
				extraAvoidance = avoidance - core.CombatTableCoverageCap
				warrior.CriticalBlockChance += extraAvoidance
			} else {
				extraAvoidance = 0.0
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -25*core.BlockRatingPerBlockChance)

			if extraAvoidance > 0.0 {
				warrior.CriticalBlockChance -= extraAvoidance
			}

			if hasT12_4P_bonus && sim.CurrentTime != sim.Duration {
				T12_4P_bonus.Activate(sim)
			}
		},
	})

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: SpellMaskShieldBlock,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 60,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock && warrior.StanceMatches(DefensiveStance)
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.ShieldBlockAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.ShieldBlock,
		Type:  core.CooldownTypeDPS,
	})

}
