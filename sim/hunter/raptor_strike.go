package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (hunter *Hunter) registerRaptorStrikeSpell() {
	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48996},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		DamageMultiplier: 1,
		CritMultiplier:   hunter.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 373.575 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(_ *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	//if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveAutosOnly || !hunter.RaptorStrike.IsReady(sim) || hunter.CurrentMana() < hunter.RaptorStrike.DefaultCast.Cost {
	//	return nil
	//}
	return mhSwingSpell
}
