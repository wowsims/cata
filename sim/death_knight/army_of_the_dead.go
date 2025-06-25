package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Summons an entire legion of Ghouls to fight for the Death Knight for 40 sec.
The Ghouls will swarm the area, taunting and fighting anything they can.
While channeling Army of the Dead, the Death Knight takes less damage equal to his Dodge plus Parry chance.
*/
func (dk *DeathKnight) registerArmyOfTheDead() {
	actionID := core.ActionID{SpellID: 42650}
	ghoulIndex := 0
	dmgReduction := 0.0

	aotdAura := dk.RegisterAura(core.Aura{
		Label:    "Army of the Dead",
		ActionID: actionID,
		Duration: time.Millisecond * 500 * 8,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			attackTable := dk.AttackTables[dk.CurrentTarget.UnitIndex]
			dmgReduction = 1.0 - (dk.GetTotalParryChanceAsDefender(attackTable) + dk.GetTotalDodgeChanceAsDefender(attackTable))
			dk.PseudoStats.DamageTakenMultiplier *= dmgReduction

			if sim.CurrentTime >= 0 {
				dk.AutoAttacks.CancelAutoSwing(sim)
			}
			dk.CancelGCDTimer(sim)

			ghoulIndex = 0
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 8,
				Period:   time.Millisecond * 500,
				OnAction: func(sim *core.Simulation) {
					ghoul := dk.ArmyGhoul[ghoulIndex]
					// Seems to always have two random ghouls spawn without a delay
					// Adding two RandomFloat calls here screws with tests though and it's minor enough that I don't care
					ghoul.EnableWithTimeout(sim, ghoul, time.Second*40)
					ghoulIndex++
				},
				CleanUp: func(sim *core.Simulation) {
					aura.Deactivate(sim)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageTakenMultiplier /= dmgReduction
			if sim.CurrentTime >= 0 {
				dk.AutoAttacks.EnableAutoSwing(sim)
			}
			dk.SetGCDTimer(sim, sim.CurrentTime)
		},
	})

	dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: DeathKnightSpellArmyOfTheDead,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDMin,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: aotdAura,
	})
}
