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
		FocusCost: core.FocusCostOptions{

			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      time.Second,
				CastTime: time.Millisecond*2000 - core.TernaryDuration(hunter.HasSetBonus(ItemSetLightningChargedBattleGear, 4), time.Millisecond*200, 0),
			},
			IgnoreHaste: true,
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.CastTime = spell.CastTime()
			},

			CastTime: func(spell *core.Spell) time.Duration {
				return time.Duration(float64(spell.DefaultCast.CastTime) / hunter.RangedSwingSpeed())
			},
		},
		BonusCritRating:          0,
		DamageMultiplierAdditive: 1 + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSteadyShot), 0.1, 0),
		DamageMultiplier:         1,
		CritMultiplier:           hunter.CritMultiplier(true, false, false),
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target)) * 0.62) + (280.182 + spell.RangedAttackPower(target)*0.021)
			focus := 9.0
			if hunter.Talents.Termination != 0 && sim.IsExecutePhase25() {
				focus = float64(hunter.Talents.Termination) * 3
			}

			if hunter.Talents.MasterMarksman != 0 {
				procChance := float64(hunter.Talents.MasterMarksman) * 0.2
				if sim.Proc(procChance, "Master Marksman Proc") && !hunter.MasterMarksmanCounterAura.IsActive() {
					hunter.MasterMarksmanCounterAura.Activate(sim)
					//hunter.MasterMarksmanCounterAura.AddStack(sim)
				}
			}

			hunter.AddFocus(sim, focus, ssMetrics)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
