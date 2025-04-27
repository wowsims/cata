package rogue

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (rogue *Rogue) registerShivSpell() {
	baseCost := int32(20)
	if ohWeapon := rogue.GetOHWeapon(); ohWeapon != nil {
		baseCost = baseCost + int32(10*ohWeapon.SwingSpeed)
	}

	rogue.Shiv = rogue.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 5938},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeOHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: RogueSpellShiv,

		EnergyCost: core.EnergyCostOptions{
			Cost: baseCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 * rogue.DWSMultiplier(),
		CritMultiplier:   rogue.CritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParryNoCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())

				switch rogue.Options.OhImbue {
				case proto.RogueOptions_DeadlyPoison:
					rogue.DeadlyPoison.Cast(sim, target)
				case proto.RogueOptions_InstantPoison:
					rogue.InstantPoison[ShivProc].Cast(sim, target)
				case proto.RogueOptions_WoundPoison:
					rogue.WoundPoison[ShivProc].Cast(sim, target)
				}
			}
		},
	})
}
