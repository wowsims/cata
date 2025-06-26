package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

/*
Tooltip:
Channels Jade lightning at the target, causing ${6*($m1+*$ap*0.386)} Nature damage over 6 sec. When dealing damage, you have a 30% chance to generate 1 Chi.

If the enemy attacks you within melee range while victim to Crackling Jade Lightning, they are knocked back a short distance. This effect has an 8 sec cooldown.

TODO: Check if it does a one-time hit check or per tick
TODO: Spell or melee hit / crit
TODO: Courageous Primal Diamond should make all ticks ignore mana cost
*/
func (monk *Monk) registerCracklingJadeLightning() {
	actionID := core.ActionID{SpellID: 117952}
	energyMetrics := monk.NewEnergyMetrics(actionID)
	manaMetrics := monk.NewManaMetrics(actionID)
	chiMetrics := monk.NewChiMetrics(core.ActionID{SpellID: 123333})
	avgScaling := monk.CalcScalingSpellDmg(0.1800000072)

	monk.RegisterSpell(core.SpellConfig{
		ActionID:       actionID,
		SpellSchool:    core.SpellSchoolNature,
		Flags:          core.SpellFlagChanneled | core.SpellFlagAPL,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: MonkSpellCracklingJadeLightning,
		MaxRange:       40,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: core.TernaryFloat64(monk.StanceMatches(WiseSerpent), 1.57, 0),
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Crackling Jade Lightning" + monk.Label,
			},
			NumberOfTicks:       6,
			TickLength:          time.Second,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				isWiseSerpent := monk.StanceMatches(WiseSerpent)
				currentResource := core.TernaryFloat64(isWiseSerpent, monk.CurrentMana(), monk.CurrentEnergy())
				cost := core.TernaryFloat64(isWiseSerpent, dot.Spell.Cost.GetCurrentCost(), 20.0)

				if currentResource >= cost {
					if isWiseSerpent {
						monk.SpendMana(sim, cost, manaMetrics)
					} else {
						monk.SpendEnergy(sim, cost, energyMetrics)
					}

					baseDamage := avgScaling + dot.Spell.MeleeAttackPower()*0.386
					dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicCrit)

					if sim.Proc(0.3, "Crackling Jade Lightning") {
						monk.AddChi(sim, dot.Spell, 1, chiMetrics)
					}
				} else {
					monk.AutoAttacks.EnableMeleeSwing(sim)
					monk.ExtendGCDUntil(sim, sim.CurrentTime+monk.ChannelClipDelay)

					// Deactivating within OnTick causes a panic since tickAction gets set to nil in the default OnExpire
					pa := sim.GetConsumedPendingActionFromPool()
					pa.NextActionAt = sim.CurrentTime

					pa.OnAction = func(sim *core.Simulation) {
						dot.Deactivate(sim)
					}

					sim.AddPendingAction(pa)
				}
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitNoHitCounter)

			if result.Landed() {
				dot := spell.Dot(target)
				dot.Apply(sim)
				expiresAt := dot.ExpiresAt()
				monk.AutoAttacks.StopMeleeUntil(sim, expiresAt)
				monk.ExtendGCDUntil(sim, expiresAt+monk.ReactionTime)
			}

			spell.DealOutcome(sim, result)
		},
	})
}
