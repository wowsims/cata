package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerArcaneMissilesSpell() {

	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)

	mage.ArcaneMissilesTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 7268},
		SpellSchool: core.SpellSchoolArcane,
		// unlike Mind Flay, this CAN proc JoW. It can also proc trinkets without the "can proc from proc" flag
		// such as illustration of the dragon soul
		// however, it cannot proc Nibelung so we add the ProcMaskNotInSpellbook flag
		ProcMask:     core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:        SpellFlagMage | core.SpellFlagNoLogs,
		MissileSpeed: 20,

		BonusCritRating: 0 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5*core.CritRatingPerCritChance, 0) +
			core.TernaryFloat64(mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfArcaneMissiles), 5*core.CritRatingPerCritChance, 0),
		DamageMultiplier: 1 + .02*float64(mage.Talents.TormentTheWeak),
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := 0.432*mage.ScalingBaseDamage + 0.278*spell.SpellPower()
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 7268},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.ArcaneMissilesAura.IsActive()
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles!",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.ArcaneMissilesAura.IsActive() {
						if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
							mage.MissileBarrageAura.Deactivate(sim)
						}
					}

					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.
					dot := mage.ArcaneMissiles.Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
					mage.ArcaneBlastAura.Deactivate(sim)
				},
			},
			NumberOfTicks:       3 + int32(mage.Talents.ImprovedArcaneMissiles),
			TickLength:          time.Millisecond*700 - 100*time.Duration(mage.Talents.MissileBarrage),
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.ArcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
