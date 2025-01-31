package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerSteadyShotSpell() {

	ssMetrics := hunter.NewFocusMetrics(core.ActionID{SpellID: 56641})

	hunter.SteadyShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 56641},
		SpellSchool:    core.SpellSchoolPhysical,
		ClassSpellMask: HunterSpellSteadyShot,
		ProcMask:       core.ProcMaskRangedSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,
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
		DamageMultiplierAdditive: core.TernaryInt64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSteadyShot), 10, 0),
		DamageMultiplier:         1,
		CritMultiplier:           hunter.CritMultiplier(true, false, false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) + (280.182 + (spell.RangedAttackPower(target) * 0.021))
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
