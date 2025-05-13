package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Tooltip:
You spin while kicking in the air, dealing ${1.59*(1.75/1.59)*$<low>} to ${1.59*(1.75/1.59)*$<high>} damage to all nearby enemies

-- Teachings of the Monastery --

	and $117640m1 healing to nearby allies

-- Teachings of the Monastery --

	every 0.75 sec, within 8 yards.

-- NOT Glyph of Spinning Crane Kick --
Movement speed is reduced by 30%.
-- NOT Glyph of Spinning Crane Kick --

Generates 1 Chi, if it hits at least 3 targets. Lasts 2 sec.
During Spinning Crane Kick, you can continue to dodge and parry.
*/
func (monk *Monk) registerSpinningCraneKick() {
	// Rushing Jade Wind replaces Spinning Crane Kick
	if monk.Talents.RushingJadeWind && monk.Level >= 90 {
		return
	}

	actionID := core.ActionID{SpellID: 101546}
	debuffActionID := core.ActionID{SpellID: 107270}
	chiMetrics := monk.NewChiMetrics(actionID)
	numTargets := monk.Env.GetNumTargets()

	spinningCraneKickTickSpell := monk.RegisterSpell(core.SpellConfig{
		ActionID:       debuffActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,
		ClassSpellMask: MonkSpellSpinningCraneKick,
		MaxRange:       8,

		DamageMultiplier: 1.75, // 1.59 * (1.75 / 1.59),
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			results := make([]*core.SpellResult, numTargets)

			for idx := int32(0); idx < numTargets; idx++ {
				currentTarget := sim.Environment.GetTargetUnit(idx)
				baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)
				results[idx] = spell.CalcDamage(sim, currentTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			}

			for idx := int32(0); idx < numTargets; idx++ {
				spell.DealDamage(sim, results[idx])
			}
		},
	})

	glyphOfSpinningCraneKick := monk.HasMajorGlyph(proto.MonkMajorGlyph_GlyphOfSpinningCraneKick)
	spinningCraneKickAura := monk.RegisterAura(core.Aura{
		Label:    "Spinning Crane Kick" + monk.Label,
		ActionID: actionID,
		Duration: time.Millisecond * 750 * 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if !glyphOfSpinningCraneKick {
				monk.MultiplyMovementSpeed(sim, 0.7)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if !glyphOfSpinningCraneKick {
				monk.MultiplyMovementSpeed(sim, 1.0/0.7)
			}
		},
	})

	isWiseSerpent := monk.StanceMatches(WiseSerpent)
	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagChanneled | SpellFlagBuilder | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellSpinningCraneKick,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryInt32(isWiseSerpent, 0, 40),
		},
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(isWiseSerpent, 7.15, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label:    "Spinning Crane Kick" + monk.Label,
				ActionID: debuffActionID,
			},
			NumberOfTicks:        3,
			TickLength:           time.Millisecond * 750,
			AffectedByCastSpeed:  true,
			HasteReducesDuration: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				spinningCraneKickTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dot := spell.AOEDot()
			dot.Apply(sim)
			dot.TickOnce(sim)

			expiresAt := dot.ExpiresAt()
			monk.AutoAttacks.DelayMeleeBy(sim, expiresAt-sim.CurrentTime)
			monk.ExtendGCDUntil(sim, expiresAt+monk.ReactionTime)

			remainingDuration := dot.RemainingDuration(sim)
			spinningCraneKickAura.Duration = remainingDuration
			spinningCraneKickAura.Activate(sim)

			if numTargets >= 3 {
				monk.AddChi(sim, spell, 1, chiMetrics)
			}
		},
	})
}
