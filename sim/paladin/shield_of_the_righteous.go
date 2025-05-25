package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) registerShieldOfTheRighteous() {
	actionID := core.ActionID{SpellID: 53600}
	baseDamage := paladin.CalcScalingSpellDmg(0.73199999332)
	apCoef := 0.61699998379

	var bastionOfGloryAura *core.Aura
	bastionOfGloryAura = paladin.RegisterAura(core.Aura{
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
			bastionOfGloryAura.Deactivate(sim)
		},
	})

	var snapshotDmgReduction float64
	shieldOfTheRighteousAura := paladin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 132403},
		Label:    "Shield of the Righteous" + paladin.Label,
		Duration: time.Second * 3,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			snapshotDmgReduction = 1.0 +
				(-0.25-paladin.ShieldOfTheRighteousAdditiveMultiplier)*(1.0+paladin.ShieldOfTheRighteousMultiplicativeMultiplier)
			paladin.PseudoStats.SchoolDamageTakenMultiplier[core.SpellSchoolPhysical] *= snapshotDmgReduction
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			paladin.PseudoStats.SchoolDamageTakenMultiplier[core.SpellSchoolPhysical] /= snapshotDmgReduction
		},
	})

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
			return paladin.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield && paladin.HolyPower.CanSpend(3)
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDamage + apCoef*spell.MeleeAttackPower()

			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				paladin.HolyPower.Spend(3, actionID, sim)
			}

			// Buff should apply even if the spell misses/dodges/parries
			// It also extends on refresh and only recomputes the damage taken mod on application, not on refresh
			if spell.RelatedSelfBuff.IsActive() {
				spell.RelatedSelfBuff.UpdateExpires(spell.RelatedSelfBuff.ExpiresAt() + spell.RelatedSelfBuff.Duration)
			} else {
				spell.RelatedSelfBuff.Activate(sim)
			}

			bastionOfGloryAura.Activate(sim)
			bastionOfGloryAura.AddStack(sim)

			spell.DealOutcome(sim, result)
		},

		RelatedSelfBuff: shieldOfTheRighteousAura,
	})
}
