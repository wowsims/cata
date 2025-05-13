package hunter

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (hunter *Hunter) registerSteadyShotSpell() {

	ssMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 56641})

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 56641},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: HunterSpellSteadyShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		MissileSpeed:   40,
		MinRange:       5,
		MaxRange:       40,
		FocusCost: core.FocusCostOptions{

			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Millisecond * 2000,
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		BonusCritPercent:         0,
		DamageMultiplierAdditive: 1 + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSteadyShot), 0.1, 0),
		DamageMultiplier:         1,
		CritMultiplier:           hunter.DefaultCritMultiplier(),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower()) + (280.182 + (spell.RangedAttackPower() * 0.021))
			intFocus := core.TernaryFloat64(hunter.T13_2pc.IsActive(), 9*2, 9)

			if hunter.Talents.Termination != 0 && sim.IsExecutePhase25() {
				intFocus += float64(hunter.Talents.Termination) * 3
			}

			if hunter.Talents.MasterMarksman != 0 {
				procChance := float64(hunter.Talents.MasterMarksman) * 0.2
				if sim.Proc(procChance, "Master Marksman Proc") && !hunter.MasterMarksmanCounterAura.IsActive() {
					hunter.MasterMarksmanCounterAura.Activate(sim)
				}
			}

			hunter.AddFocus(sim, intFocus, ssMetrics)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
