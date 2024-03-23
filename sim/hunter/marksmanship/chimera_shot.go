package marksmanship

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *MarksmanshipHunter) registerChimeraShotSpell() {
	if !hunter.Talents.ChimeraShot {
		return
	}

	hunter.ChimeraShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53209},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 50 - (float64(hunter.Talents.Efficiency) * 2),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfChimeraShot), time.Second*1, 0),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			wepDmg := hunter.AutoAttacks.Ranged().CalculateNormalizedWeaponDamage(sim, spell.RangedAttackPower(target))
			baseDamage := 0.732*spell.RangedAttackPower(target) + 1620.33

			result := spell.CalcDamage(sim, target, wepDmg+baseDamage, spell.OutcomeRangedHitAndCrit)
			if result.Landed() {
				if hunter.SerpentSting.Dot(target).IsActive() {
					hunter.SerpentSting.Dot(target).Rollover(sim)
				}
			}
			spell.DealDamage(sim, result)
		},
	})
}
