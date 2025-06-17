package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Word of Glory:

Consumes up to 3 Holy Power to heal a friendly target for (<5239-5838> + 0.49 * <SP>)

-- Glyph of Harsh Words --
or harm an enemy target for (<4029-4490> + 0.377 * <SP>)
-- /Glyph of Harsh Words --

per charge of Holy Power.

-- Glyph of Protector of the Innocent --
If used to heal another target, you will be healed for 20% of the amount healed
-- /Glyph of Protector of the Innocent --

-- Glyph of Word of Glory --
Your damage is increased by 3% per Holy Power spent for 6 sec after you cast Word of Glory
-- /Glyph of Word of Glory --

Eternal Flame:

Consumes up to 3 Holy Power to place a protective Holy flame on a friendly target,
which heals them for (<5239-5838> + 0.49 * <SP>) and an additional (712 + 0.819 * <SP>) every 3 sec for 30 sec.
Healing increased per charge of Holy Power.
The heal over time is increased by 50% if used on the Paladin.
Replaces Word of Glory.
*/
func (paladin *Paladin) registerWordOfGlory() {
	isProt := paladin.Spec == proto.Spec_SpecProtectionPaladin
	eternalFlame := paladin.Talents.EternalFlame

	actionID := core.ActionID{SpellID: core.TernaryInt32(eternalFlame, 114163, 85673)}

	config := core.SpellConfig{
		ActionID:       actionID,
		Flags:          core.SpellFlagAPL | core.SpellFlagHelpful,
		ProcMask:       core.ProcMaskSpellHealing,
		SpellSchool:    core.SpellSchoolHoly,
		ClassSpellMask: SpellMaskWordOfGlory,
		MetricSplits:   4,

		MaxRange: 40,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.TernaryDuration(isProt, 0, core.GCDDefault),
				NonEmpty: isProt,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Millisecond * 1500,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				paladin.DynamicHolyPowerSpent = core.TernaryInt32(paladin.BastionOfPowerAura.IsActive(), 3, paladin.SpendableHolyPower())
				spell.SetMetricsSplit(paladin.DynamicHolyPowerSpent)
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return paladin.BastionOfPowerAura.IsActive() || paladin.HolyPower.CanSpend(1)
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		BonusCoefficient: 0.49000000954,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if target.IsOpponent(&paladin.Unit) {
				target = &paladin.Unit
			}

			baseHealing := paladin.CalcAndRollDamageRange(sim, 4.84999990463, 0.1080000028)

			damageMultiplier := spell.DamageMultiplier
			spell.DamageMultiplier *= float64(paladin.DynamicHolyPowerSpent)
			if target == &paladin.Unit {
				spell.DamageMultiplier *= 1.0 + paladin.BastionOfGloryMultiplier
			}

			result := spell.CalcHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)

			spell.DamageMultiplier = damageMultiplier

			if !paladin.BastionOfPowerAura.IsActive() {
				paladin.HolyPower.SpendUpTo(sim, paladin.DynamicHolyPowerSpent, actionID)
			}

			spell.DealHealing(sim, result)

			if eternalFlame {
				spell.Hot(target).Apply(sim)
			}
		},
	}

	if eternalFlame {
		config.Hot = paladin.eternalFlameHotConfig()
	}

	paladin.RegisterSpell(config)
}

func (paladin *Paladin) eternalFlameHotConfig() core.DotConfig {
	return core.DotConfig{
		Aura: core.Aura{
			Label:    "Eternal Flame",
			Duration: time.Second * 30,
		},

		TickLength:          time.Second * 3,
		NumberOfTicks:       10,
		AffectedByCastSpeed: true,

		BonusCoefficient: 0.08190000057,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			dot.SnapshotHeal(target, paladin.CalcScalingSpellDmg(0.62300002575))
			dot.SnapshotAttackerMultiplier *= float64(paladin.DynamicHolyPowerSpent)
			if target == &paladin.Unit {
				dot.SnapshotAttackerMultiplier *= 1.0 + paladin.BastionOfGloryMultiplier
				dot.SnapshotAttackerMultiplier *= 1.5
			}
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeSnapshotCrit)
		},
	}
}
