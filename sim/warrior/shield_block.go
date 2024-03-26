package warrior

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func (warrior *Warrior) RegisterShieldBlockCD() {
	actionID := core.ActionID{SpellID: 2565}

	// TODO: check if the crit block chance adapts to the current target in-game
	// for now, use a most-likely attack table (level+3 boss)
	atkTableAttacker := &core.Unit{Level: warrior.Level + 3, Type: core.EnemyUnit}
	atkTable := core.NewAttackTable(atkTableAttacker, &warrior.Unit)
	extraAvoidance := 0.0
	warrior.ShieldBlockAura = warrior.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 25*core.BlockRatingPerBlockChance)

			// TODO: Wording of the second effect implies it'll use the higher of the two between total avoidance
			// or block, this should be tested though
			avoidance := warrior.GetTotalAvoidanceChance(atkTable)
			blockChance := warrior.GetTotalBlockChanceAsDefender(atkTable)
			highestChance := math.Max(avoidance, blockChance)
			if highestChance > core.CombatTableCoverageCap {
				extraAvoidance = highestChance - core.CombatTableCoverageCap
				warrior.CriticalBlockChance += extraAvoidance
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -25*core.BlockRatingPerBlockChance)
			if extraAvoidance > 0.0 {
				warrior.CriticalBlockChance -= extraAvoidance
			}
		},
	})

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second*60 - time.Second*10*time.Duration(warrior.Talents.ShieldMastery),
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
