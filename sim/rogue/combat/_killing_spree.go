package combat

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/rogue"
)

func (comRogue *CombatRogue) registerKillingSpreeCD() {
	if !comRogue.Talents.KillingSpree {
		return
	}

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

		DamageMultiplier: 1 * comRogue.DWSMultiplier(),
		CritMultiplier:   comRogue.CritMultiplier(false),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})

	hasGlyph := comRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfKillingSpree)
	auraDamageMult := core.TernaryFloat64(hasGlyph, 1.3, 1.2)
	comRogue.KillingSpreeAura = comRogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second*2 + 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.SetGCDTimer(sim, sim.CurrentTime+aura.Duration)
			comRogue.PseudoStats.DamageDealtMultiplier *= auraDamageMult
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        5,
				TickImmediately: true,
				OnAction: func(s *core.Simulation) {
					targetCount := sim.GetNumTargets()
					target := comRogue.CurrentTarget
					if targetCount > 1 {
						newUnitIndex := int32(math.Ceil(float64(targetCount)*sim.RandomFloat("Killing Spree"))) - 1
						target = sim.GetTargetUnit(newUnitIndex)
					}
					mhWeaponSwing.Cast(sim, target)
					ohWeaponSwing.Cast(sim, target)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.PseudoStats.DamageDealtMultiplier /= auraDamageMult
		},
	})
	comRogue.KillingSpree = comRogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 51690},
		Flags:          core.SpellFlagAPL,
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
