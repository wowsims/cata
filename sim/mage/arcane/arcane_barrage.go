package arcane

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/mage"
)

func (Mage *ArcaneMage) registerArcaneBarrageSpell() {

	Mage.ArcaneBarrage = Mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44425},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          mage.SpellFlagMage | mage.ArcaneMissileSpells | core.SpellFlagAPL,
		ClassSpellMask: mage.MageSpellArcaneBarrage,
		MissileSpeed:   24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    Mage.NewTimer(),
				Duration: time.Second * 4,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   Mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.907,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.413 * Mage.ScalingBaseDamage
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
			Mage.ArcaneBlastAura.Deactivate(sim)
		},
	})
}
