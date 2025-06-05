package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type Stance uint8

const (
	StanceNone         = 0
	FierceTiger Stance = 1 << iota
	SturdyOx
	WiseSerpent
)

const stanceEffectCategory = "Stance"

func (monk *Monk) StanceMatches(other Stance) bool {
	return (monk.Stance & other) != 0
}

func (monk *Monk) makeStanceSpell(aura *core.Aura, stanceCD *core.Timer) *core.Spell {
	return monk.RegisterSpell(core.SpellConfig{
		ActionID: aura.ActionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    stanceCD,
				Duration: time.Second,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !aura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			aura.Activate(sim)
		},
	})
}

func (monk *Monk) registerStanceOfTheSturdyOx(stanceCD *core.Timer) {
	if monk.Spec != proto.Spec_SpecBrewmasterMonk {
		return
	}
	actionID := core.ActionID{SpellID: 115069}
	chiMetrics := monk.NewChiMetrics(actionID)
	stamDep := monk.NewDynamicMultiplyStat(stats.Stamina, 1.2)

	monk.StanceOfTheSturdyOxAura = monk.GetOrRegisterAura(core.Aura{
		Label:    "Stance of the Sturdy Ox" + monk.Label,
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			monk.Stance = SturdyOx
			monk.MultiplyEnergyRegenSpeed(sim, 1.1)
			monk.SetCurrentPowerBar(core.EnergyBar)

			currentChi := monk.GetChi()
			if currentChi > 1 {
				monk.SpendChi(sim, currentChi-1, chiMetrics)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.Stance = StanceNone
			monk.MultiplyEnergyRegenSpeed(sim, 1.0/1.1)
		},
	}).AttachMultiplicativePseudoStatBuff(
		&monk.PseudoStats.DamageTakenMultiplier, 0.75,
	).AttachAdditivePseudoStatBuff(
		&monk.PseudoStats.ReducedCritTakenChance, 0.06,
	).AttachStatDependency(
		stamDep,
	).AttachMultiplicativePseudoStatBuff(
		&monk.PseudoStats.ThreatMultiplier, 7.0,
	)

	monk.StanceOfTheSturdyOxAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
	monk.StanceOfTheSturdyOx = monk.makeStanceSpell(monk.StanceOfTheSturdyOxAura, stanceCD)
}

func (monk *Monk) registerStanceOfTheWiseSerpent(stanceCD *core.Timer) {
	if monk.Spec != proto.Spec_SpecMistweaverMonk {
		return
	}
	actionID := core.ActionID{SpellID: 117895}
	chiMetrics := monk.NewChiMetrics(actionID)

	hitDep := monk.NewDynamicStatDependency(stats.Spirit, stats.HitRating, 0.5)
	expDep := monk.NewDynamicStatDependency(stats.Spirit, stats.ExpertiseRating, 0.5)
	hasteDep := monk.NewDynamicMultiplyStat(stats.HasteRating, 1.5)
	// TODO: This should be a replacement not a dependency.
	// apDep := monk.NewDynamicStatDependency(stats.SpellPower, stats.AttackPower, 2)

	dmgDone := 0.0

	eminenceHeal := monk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

		DamageMultiplier: 0.25,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, dmgDone, spell.OutcomeHealing)
		},
	})

	// When the Monk deals non-autoattack damage, he/she will heal the lowest health nearby target within 20 yards equal to 25% of the damage done.
	eminenceAura := monk.RegisterAura(core.Aura{
		Label:    "Eminence" + monk.Label,
		ActionID: core.ActionID{SpellID: 126890},
		Duration: core.NeverExpires,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result == nil || !result.Landed() || result.Damage == 0 || spell.ProcMask.Matches(core.ProcMaskWhiteHit) {
				return
			}

			dmgDone = result.Damage
			// Should be a smart heal
			eminenceHeal.Cast(sim, &monk.Unit)
		},
	})

	monk.StanceOfTheWiseSerpentAura = monk.GetOrRegisterAura(core.Aura{
		Label:    "Stance of the Wise Serpent" + monk.Label,
		ActionID: core.ActionID{SpellID: 136336},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			monk.Stance = WiseSerpent
			// monk.EnableDynamicStatDep(sim, apDep)
			monk.SetCurrentPowerBar(core.ManaBar)
			eminenceAura.Activate(sim)

			currentChi := monk.GetChi()
			if currentChi > 1 {
				monk.SpendChi(sim, currentChi-1, chiMetrics)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			eminenceAura.Deactivate(sim)
			monk.SetCurrentPowerBar(core.EnergyBar)
			// monk.DisableDynamicStatDep(sim, apDep)
			monk.Stance = StanceNone
		},
	}).AttachMultiplicativePseudoStatBuff(
		&monk.PseudoStats.HealingDealtMultiplier, 1.2,
	).AttachStatDependency(
		hitDep,
	).AttachStatDependency(
		expDep,
	).AttachStatDependency(
		hasteDep,
	)

	monk.StanceOfTheWiseSerpentAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
	monk.StanceOfTheWiseSerpent = monk.makeStanceSpell(monk.StanceOfTheWiseSerpentAura, stanceCD)
}

/*
Increases your movement speed by 10%, increases damage done by 10% and increases the amount of Chi generated by your Jab and Expel Harm abilities by 1.
*/
func (monk *Monk) registerStanceOfTheFierceTiger(stanceCD *core.Timer) {
	actionID := core.ActionID{SpellID: 103985}
	chiMetrics := monk.NewChiMetrics(actionID)

	monk.StanceOfTheFierceTigerAura = monk.GetOrRegisterAura(core.Aura{
		Label:    "Stance of the Fierce Tiger" + monk.Label,
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			monk.Stance = FierceTiger
			monk.SetCurrentPowerBar(core.EnergyBar)

			currentChi := monk.GetChi()
			if currentChi > 1 {
				monk.SpendChi(sim, currentChi-1, chiMetrics)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			monk.Stance = StanceNone
		},
	}).AttachMultiplicativePseudoStatBuff(
		&monk.PseudoStats.DamageDealtMultiplier, 1.1,
	).AttachMultiplicativePseudoStatBuff(
		// This **does** stack with other movement speed buffs.
		&monk.PseudoStats.MovementSpeedMultiplier, 1.1,
	)

	monk.StanceOfTheFierceTigerAura.NewExclusiveEffect(stanceEffectCategory, true, core.ExclusiveEffect{})
	monk.StanceOfTheFierceTiger = monk.makeStanceSpell(monk.StanceOfTheFierceTigerAura, stanceCD)
}

func (monk *Monk) registerStances() {
	stanceCD := monk.NewTimer()
	monk.registerStanceOfTheSturdyOx(stanceCD)
	monk.registerStanceOfTheWiseSerpent(stanceCD)
	monk.registerStanceOfTheFierceTiger(stanceCD)
}
