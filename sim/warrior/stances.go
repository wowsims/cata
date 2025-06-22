package warrior

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

type Stance uint8

const (
	StanceNone          = 0
	BattleStance Stance = 1 << iota
	DefensiveStance
	BerserkerStance
)

const stanceEffectCategory = "Stance"

func (warrior *Warrior) StanceMatches(other Stance) bool {
	return (warrior.Stance & other) != 0
}

func (warrior *Warrior) makeStanceSpell(stance Stance, aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	actionID := aura.ActionID

	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Millisecond * 1500,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.Stance != stance
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			// TODO: see if this is fixed in 4.4.0
			if warrior.WarriorInputs.StanceSnapshot {
				// Delayed, so same-GCD casts are affected by the current aura.
				//  Alternatively, those casts could just (artificially) happen before the stance change.
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt:     sim.CurrentTime + 10*time.Millisecond,
					OnAction: aura.Activate,
				})
			} else {
				aura.Activate(sim)
			}

			warrior.Stance = stance
		},
	})
}

func (warrior *Warrior) registerBattleStanceAura() {
	actionID := core.ActionID{SpellID: 2457}

	warrior.BattleStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Battle Stance",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAutoAttackRageGen(2.0)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyAutoAttackRageGen(1.0 / 2.0)
		},
	})
	warrior.BattleStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerDefensiveStanceAura() {
	actionID := core.ActionID{SpellID: 71}
	rageMetrics := warrior.NewRageMetrics(actionID)

	var pa *core.PendingAction
	warrior.DefensiveStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Defensive Stance",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.GetAura("RageBar").Deactivate(sim)
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 3,
				OnAction: func(sim *core.Simulation) {
					if sim.CurrentTime > 0 {
						warrior.AddRage(sim, 1, rageMetrics)
					}
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.GetAura("RageBar").Activate(sim)
			if pa != nil {
				pa.Cancel(sim)
				pa = nil
			}
		},
	}).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.ThreatMultiplier, 7,
	).AttachMultiplicativePseudoStatBuff(
		&warrior.PseudoStats.DamageTakenMultiplier, 0.75,
	)

	warrior.DefensiveStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerBerserkerStanceAura() {
	actionId := core.ActionID{SpellID: 2458}
	rageMetrics := warrior.NewRageMetrics(actionId)
	warrior.BerserkerStanceAura = warrior.GetOrRegisterAura(core.Aura{
		Label:    "Berserker Stance",
		ActionID: actionId,
		Duration: core.NeverExpires,
	}).AttachProcTrigger(core.ProcTrigger{
		Name:     "Berserker Stance - Rage Gain",
		ActionID: actionId,
		Callback: core.CallbackOnSpellHitTaken,
		Outcome:  core.OutcomeLanded,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			warrior.AddRage(sim, math.Floor(result.Damage/warrior.MaxHealth()*100), rageMetrics)
		},
	})
	warrior.BerserkerStanceAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
}

func (warrior *Warrior) registerStances() {
	stanceCD := warrior.NewTimer()
	warrior.registerBattleStanceAura()
	warrior.registerDefensiveStanceAura()
	warrior.registerBerserkerStanceAura()
	warrior.BattleStance = warrior.makeStanceSpell(BattleStance, warrior.BattleStanceAura, stanceCD)
	warrior.DefensiveStance = warrior.makeStanceSpell(DefensiveStance, warrior.DefensiveStanceAura, stanceCD)
	warrior.BerserkerStance = warrior.makeStanceSpell(BerserkerStance, warrior.BerserkerStanceAura, stanceCD)
}
