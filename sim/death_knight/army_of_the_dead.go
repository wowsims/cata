package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerArmyOfTheDeadSpell() {
	var ghoulIndex = 0
	aotdAura := dk.RegisterAura(core.Aura{
		Label:    "Army of the Dead",
		ActionID: core.ActionID{SpellID: 42650},
		Duration: time.Millisecond * 500 * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime >= 0 {
				dk.AutoAttacks.CancelAutoSwing(sim)
			}
			dk.CancelGCDTimer(sim)

			ghoulIndex = 0
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 8,
				Period:   time.Millisecond * 500,
				OnAction: func(sim *core.Simulation) {
					dk.ArmyGhoul[ghoulIndex].EnableWithTimeout(sim, dk.ArmyGhoul[ghoulIndex], time.Second*40)
					ghoulIndex++
				},
				CleanUp: func(sim *core.Simulation) {
					aura.Deactivate(sim)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if sim.CurrentTime >= 0 {
				dk.AutoAttacks.EnableAutoSwing(sim)
			}
			dk.SetGCDTimer(sim, sim.CurrentTime)
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 42650},
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellArmyOfTheDead,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 30,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 10,
			},
		},

		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			aotdAura.Activate(sim)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
			return dk.HasActiveAuraWithTag(core.UnholyFrenzyAuraTag)
		},
	})
}
