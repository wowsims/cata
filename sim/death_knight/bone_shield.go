package death_knight

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (dk *DeathKnight) registerBoneShieldSpell() {
	if !dk.Talents.BoneShield {
		return
	}

	actionID := core.ActionID{SpellID: 49222}
	stackRemovalCd := 0 * time.Second

	aura := dk.RegisterAura(core.Aura{
		Label:     "Bone Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 5,
		MaxStacks: 6,
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				if sim.CurrentTime > stackRemovalCd+2*time.Second {
					stackRemovalCd = sim.CurrentTime

					aura.RemoveStack(sim)
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier *= 1.02
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
			stackRemovalCd = sim.CurrentTime
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.DamageDealtMultiplier /= 1.02
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,
		ClassSpellMask: DeathKnightSpellBoneShield,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			aura.Activate(sim)
			aura.SetStacks(sim, aura.MaxStacks)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return dk.CurrentHealthPercent() < 0.6
			},
		})
	}
}
