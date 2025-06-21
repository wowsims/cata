package combat

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerKillingSpreeCD() {
	mhWeaponSwing := comRogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51690, Tag: 1}, // actual spellID is 57841
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: rogue.RogueSpellKillingSpreeHit,

		DamageMultiplier: 1,
		CritMultiplier:   comRogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
	ohWeaponSwing := comRogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51690, Tag: 2}, // actual spellID is 57842
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics,
		ClassSpellMask: rogue.RogueSpellKillingSpreeHit,

		DamageMultiplier: 1.75, // Combat has a 1.75x OH damage multiplier
		CritMultiplier:   comRogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	auraDamageMult := 1.5
	comRogue.KillingSpreeAura = comRogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second*3 + 1, // +1 ensures the final hit is buffed
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.SetGCDTimer(sim, sim.CurrentTime+aura.Duration)
			comRogue.PseudoStats.DamageDealtMultiplier *= auraDamageMult
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        7,
				TickImmediately: true,
				OnAction: func(s *core.Simulation) {
					targetCount := sim.GetNumTargets()
					target := comRogue.CurrentTarget
					if targetCount > 1 && comRogue.HasActiveAura("Blade Flurry") {
						newUnitIndex := int32(math.Ceil(float64(targetCount)*sim.RandomFloat("Killing Spree"))) - 1
						target = sim.GetTargetUnit(newUnitIndex)
					}
					mhWeaponSwing.Cast(sim, target)
					ohWeaponSwing.Cast(sim, target)
					if comRogue.T16SpecMod != nil {
						if comRogue.T16SpecMod.IsActive {
							newMod := comRogue.T16SpecMod.GetFloatValue() * 1.1
							comRogue.T16SpecMod.UpdateFloatValue(newMod)
						}
						comRogue.T16SpecMod.Activate()
					}
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.PseudoStats.DamageDealtMultiplier /= auraDamageMult
			if comRogue.T16SpecMod != nil {
				comRogue.T16SpecMod.UpdateFloatValue(0.1)
				comRogue.T16SpecMod.Deactivate()
			}
		},
	})
	comRogue.KillingSpree = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51690},
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: rogue.RogueSpellKillingSpree,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    comRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			comRogue.BreakStealth(sim)
			comRogue.KillingSpreeAura.Activate(sim)
		},
	})

	comRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    comRogue.KillingSpree,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
	})
}
