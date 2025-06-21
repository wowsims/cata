package protection

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/paladin"
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
func (prot *ProtectionPaladin) registerShieldOfTheRighteous() {
	prot.BastionOfGloryAura = prot.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 114637},
		Label:     "Bastion of Glory" + prot.Label,
		Duration:  time.Second * 20,
		MaxStacks: 5,

		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			if newStacks == 0 {
				prot.BastionOfGloryMultiplier = 0.0
				return
			}

			prot.BastionOfGloryMultiplier = 0.1*float64(newStacks) + prot.ShieldOfTheRighteousAdditiveMultiplier
		},
	}).AttachProcTrigger(core.ProcTrigger{
		Callback:       core.CallbackOnCastComplete,
		ClassSpellMask: paladin.SpellMaskWordOfGlory,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			prot.BastionOfGloryAura.Deactivate(sim)
		},
	})

	var snapshotDmgReduction float64
	shieldOfTheRighteousAura := prot.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 132403},
		Label:     "Shield of the Righteous" + prot.Label,
		Duration:  time.Second * 3,
		MaxStacks: 100,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			snapshotDmgReduction = 1.0 +
				(-0.25-prot.ShieldOfTheRighteousAdditiveMultiplier)*(1.0+prot.ShieldOfTheRighteousMultiplicativeMultiplier)

			snapshotDmgReduction = max(0.2, snapshotDmgReduction)

			prot.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= snapshotDmgReduction

			percent := int32(math.Round((1.0 - snapshotDmgReduction) * 100))
			if percent > 0 {
				aura.SetStacks(sim, percent)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			prot.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= snapshotDmgReduction
		},
	})

	prot.AddDefensiveCooldownAura(shieldOfTheRighteousAura)

	actionID := core.ActionID{SpellID: 53600}

	prot.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		ClassSpellMask: paladin.SpellMaskShieldOfTheRighteous,

		MaxRange: core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return prot.PseudoStats.CanBlock && prot.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 1,
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := prot.CalcScalingSpellDmg(0.73199999332) + 0.61699998379*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				prot.HolyPower.Spend(sim, 3, actionID)
			}

			// Buff should apply even if the spell misses/dodges/parries
			// It also extends on refresh and only recomputes the damage taken mod on application, not on refresh
			if spell.RelatedSelfBuff.IsActive() {
				spell.RelatedSelfBuff.UpdateExpires(spell.RelatedSelfBuff.ExpiresAt() + spell.RelatedSelfBuff.Duration)
			} else {
				spell.RelatedSelfBuff.Activate(sim)
			}

			prot.BastionOfGloryAura.Activate(sim)
			prot.BastionOfGloryAura.AddStack(sim)

			spell.DealDamage(sim, result)
		},

		RelatedSelfBuff: shieldOfTheRighteousAura,
	})
}
