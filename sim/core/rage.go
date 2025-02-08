package core

import (
	"fmt"

	"github.com/wowsims/cata/sim/core/proto"
)

const MaxRage = 100.0
const ThreatPerRageGained = 5

type OnRageGainCB func(sim *Simulation, spell *Spell, result *SpellResult, rage float64) float64

type rageBar struct {
	unit *Unit

	startingRage float64
	currentRage  float64

	RageRefundMetrics *ResourceMetrics
}

type RageBarOptions struct {
	StartingRage   float64
	MHSwingSpeed   float64
	OHSwingSpeed   float64
	RageMultiplier float64

	// Called when rage is calculated from an OnSpellHitDealt event
	// but before it has been applied to the unit
	OnHitDealtRageGain OnRageGainCB

	// Called when rage is calculated from an OnSpellHitTaken event
	// but before it has been applied to the unit
	OnHitTakenRageGain OnRageGainCB
}

func (unit *Unit) EnableRageBar(options RageBarOptions) {
	rageFromDamageTakenMetrics := unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken})

	unit.SetCurrentPowerBar(RageBar)
	unit.RegisterAura(Aura{
		Label:    "RageBar",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}
			if result.Outcome.Matches(OutcomeMiss) {
				return
			}

			hitFactor := 6.5
			var speed float64
			if spell.ProcMask == ProcMaskMeleeMHAuto {
				speed = options.MHSwingSpeed
			} else if spell.ProcMask == ProcMaskMeleeOHAuto {
				// OH hits generate 50% of the rage they would if they were MH hits
				hitFactor /= 2
				speed = options.OHSwingSpeed
			} else {
				return
			}

			// Currently, rage does not get doubled for crits in cataclysm
			// Leaving the code here for reference with a note.
			//if result.Outcome.Matches(OutcomeCrit) {
			//	hitFactor *= 2
			//}

			// TODO: Cataclysm dodge/parry behavior
			// damage := result.Damage
			// if result.Outcome.Matches(OutcomeDodge | OutcomeParry) {
			// 	// Rage is still generated for dodges/parries, based on the damage it WOULD have done.
			// 	damage = result.PreOutcomeDamage
			// }

			// rage in cata is normalized so it only depends on weapon swing speed and some multipliers
			generatedRage := hitFactor * speed
			generatedRage *= options.RageMultiplier

			if options.OnHitDealtRageGain != nil {
				generatedRage = options.OnHitDealtRageGain(sim, spell, result, generatedRage)
			}

			var metrics *ResourceMetrics
			if spell.Cost != nil {
				metrics = spell.Cost.SpellCostFunctions.(*RageCost).ResourceMetrics
			} else {
				if spell.ResourceMetrics == nil {
					spell.ResourceMetrics = spell.Unit.NewRageMetrics(spell.ActionID)
				}
				metrics = spell.ResourceMetrics
			}
			unit.AddRage(sim, generatedRage, metrics)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if unit.GetCurrentPowerBar() != RageBar {
				return
			}

			// Formula taken from simc and verified using in-game
			// measurements. The in-game measurements agree closely for
			// incoming melee attacks but not for spells, so more work is
			// required to determine whether spells use a different formula or
			// whether Beta is simply bugged for spell Rage calculations. Note
			// that the below formula delivers full Rage even on missed
			// attacks! This is intentional and matches in-game behavior.
			generatedRage := result.PreOutcomeDamage / result.ResistanceMultiplier * 18.92 / unit.MaxHealth()

			if options.OnHitTakenRageGain != nil {
				generatedRage = options.OnHitTakenRageGain(sim, spell, result, generatedRage)
			}

			unit.AddRage(sim, generatedRage, rageFromDamageTakenMetrics)
		},
	})

	// Not a real spell, just holds metrics from rage gain threat.
	unit.RegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionRageGain},
	})

	unit.rageBar = rageBar{
		unit:              unit,
		startingRage:      max(0, min(options.StartingRage, MaxRage)),
		RageRefundMetrics: unit.NewRageMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
	}
}

func (unit *Unit) HasRageBar() bool {
	return unit.rageBar.unit != nil
}

func (rb *rageBar) CurrentRage() float64 {
	return rb.currentRage
}

func (rb *rageBar) AddRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative rage!")
	}

	newRage := min(rb.currentRage+amount, MaxRage)
	metrics.AddEvent(amount, newRage-rb.currentRage)

	if sim.Log != nil {
		rb.unit.Log(sim, "Gained %0.3f rage from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rb.currentRage, newRage, 100.0)
	}

	rb.currentRage = newRage
	if !sim.Options.Interactive {
		rb.unit.ReactToEvent(sim)
	}
}

func (rb *rageBar) SpendRage(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative rage!")
	}

	newRage := rb.currentRage - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rb.unit.Log(sim, "Spent %0.3f rage from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rb.currentRage, newRage, 100.0)
	}

	rb.currentRage = newRage
}

func (rb *rageBar) reset(_ *Simulation) {
	if rb.unit == nil {
		return
	}

	rb.currentRage = rb.startingRage
}

func (rb *rageBar) doneIteration() {
	if rb.unit == nil {
		return
	}

	rageGainSpell := rb.unit.GetSpell(ActionID{OtherID: proto.OtherAction_OtherActionRageGain})

	for _, resourceMetrics := range rb.unit.Metrics.resources {
		if resourceMetrics.Type != proto.ResourceType_ResourceTypeRage {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionDamageTaken}) {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionRefund}) {
			continue
		}
		if resourceMetrics.ActualGainForCurrentIteration() <= 0 {
			continue
		}

		// Need to exclude rage gained from white hits. Rather than have a manual list of all IDs that would
		// apply here (autos, WF attack, sword spec procs, etc), just check if the effect caused any damage.
		sourceSpell := rb.unit.GetSpell(resourceMetrics.ActionID)
		if sourceSpell != nil && sourceSpell.SpellMetrics[0].TotalDamage > 0 {
			continue
		}

		rageGainSpell.SpellMetrics[0].Casts += resourceMetrics.EventsForCurrentIteration()
		rageGainSpell.ApplyAOEThreatIgnoreMultipliers(resourceMetrics.ActualGainForCurrentIteration() * ThreatPerRageGained)
	}
}

type RageCostOptions struct {
	Cost float64

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.RageRefundMetrics if not supplied.
}
type RageCost struct {
	Refund          float64
	RefundMetrics   *ResourceMetrics
	ResourceMetrics *ResourceMetrics
}

func newRageCost(spell *Spell, options RageCostOptions) *SpellCost {
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.RageRefundMetrics
	}

	return &SpellCost{
		spell:      spell,
		BaseCost:   options.Cost,
		Multiplier: 100,
		SpellCostFunctions: &RageCost{
			Refund:          options.Refund,
			RefundMetrics:   options.RefundMetrics,
			ResourceMetrics: spell.Unit.NewRageMetrics(spell.ActionID),
		},
	}
}

func (rc *RageCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()
	return spell.Unit.CurrentRage() >= spell.CurCast.Cost
}
func (rc *RageCost) CostFailureReason(sim *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough rage (Current Rage = %0.03f, Rage Cost = %0.03f)", spell.Unit.CurrentRage(), spell.CurCast.Cost)
}
func (rc *RageCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendRage(sim, spell.CurCast.Cost, rc.ResourceMetrics)
	}
}
func (rc *RageCost) IssueRefund(sim *Simulation, spell *Spell) {
	if rc.Refund > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddRage(sim, rc.Refund*spell.CurCast.Cost, rc.RefundMetrics)
	}
}
