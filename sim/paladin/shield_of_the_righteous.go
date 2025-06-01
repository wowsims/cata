package paladin

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Instantly slam the target with your shield, causing (836 + 0.617 * <AP>) Holy damage, reducing the physical damage you take by (25 + <Divine Bulwark>)% for 3 sec, and causing Bastion of Glory.

Bastion of Glory
Increases the strength of

-- Eternal Flame --
your Eternal Flame
-- else --
your Word of Glory
----------

when used to heal yourself by 10%.

-- Selfless Healer --
Selfless Healer also increases healing from Flash of Light on yourself by 20% per stack
-- /Selfless Healer --

Stacks up to 5 times.
*/
func (paladin *Paladin) registerShieldOfTheRighteous() {
	paladin.BastionOfGloryAura = paladin.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 114637},
		Label:     "Bastion of Glory" + paladin.Label,
		Duration:  time.Second * 20,
		MaxStacks: 5,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks == 0 {
				paladin.BastionOfGloryMultiplier = 0.0
				return
			}

			paladin.BastionOfGloryMultiplier = 0.1*float64(newStacks) + paladin.ShieldOfTheRighteousAdditiveMultiplier
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: SpellMaskWordOfGlory,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			paladin.BastionOfGloryAura.Deactivate(sim)
		},
	})

	var snapshotDmgReduction float64
	shieldOfTheRighteousAura := paladin.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 132403},
		Label:     "Shield of the Righteous" + paladin.Label,
		Duration:  time.Second * 3,
		MaxStacks: 100,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			snapshotDmgReduction = 1.0 +
				(-0.25-paladin.ShieldOfTheRighteousAdditiveMultiplier)*(1.0+paladin.ShieldOfTheRighteousMultiplicativeMultiplier)

			snapshotDmgReduction = max(0.2, snapshotDmgReduction)

			paladin.PseudoStats.SchoolDamageTakenMultiplier[core.SpellSchoolPhysical] *= snapshotDmgReduction

			percent := int32(math.Round((1.0 - snapshotDmgReduction) * 100))
			if percent > 0 {
				aura.SetStacks(sim, percent)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.SchoolDamageTakenMultiplier[core.SpellSchoolPhysical] /= snapshotDmgReduction
		},
	})

	paladin.AddDefensiveCooldownAura(shieldOfTheRighteousAura)

	actionID := core.ActionID{SpellID: 53600}

	paladin.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: SpellMaskShieldOfTheRighteous,

		MaxRange: core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.PseudoStats.CanBlock && paladin.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := paladin.CalcScalingSpellDmg(0.73199999332) + 0.61699998379*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Spend(sim, 3, actionID)
			}

			// Buff should apply even if the spell misses/dodges/parries
			// It also extends on refresh and only recomputes the damage taken mod on application, not on refresh
			if spell.RelatedSelfBuff.IsActive() {
				spell.RelatedSelfBuff.UpdateExpires(spell.RelatedSelfBuff.ExpiresAt() + spell.RelatedSelfBuff.Duration)
			} else {
				spell.RelatedSelfBuff.Activate(sim)
			}

			paladin.BastionOfGloryAura.Activate(sim)
			paladin.BastionOfGloryAura.AddStack(sim)

			spell.DealDamage(sim, result)
		},

		RelatedSelfBuff: shieldOfTheRighteousAura,
	})
}
