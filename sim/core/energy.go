package core

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type energyBar struct {
	unit *Unit

	maxEnergy           float64
	currentEnergy       float64
	startingComboPoints int32
	maxComboPoints      int32
	comboPoints         int32
	nextEnergyTick      time.Duration

	// Time between Energy ticks.
	EnergyTickDuration time.Duration
	EnergyPerTick      float64

	// These two terms are multiplied together to scale the total Energy regen from ticks.
	energyRegenMultiplier float64
	hasteRatingMultiplier float64

	regenMetrics        *ResourceMetrics
	EnergyRefundMetrics *ResourceMetrics

	ownerClass              proto.Class
	comboPointsResourceName string // "chi" or "combo points"
	hasNoRegen              bool   // some units have an energy bar but do not require regen ticks
}
type EnergyBarOptions struct {
	StartingComboPoints int32
	MaxComboPoints      int32
	MaxEnergy           float64
	UnitClass           proto.Class
	HasNoRegen          bool
}

func (unit *Unit) EnableEnergyBar(options EnergyBarOptions) {
	unit.SetCurrentPowerBar(EnergyBar)

	unit.energyBar = energyBar{
		unit:                    unit,
		maxEnergy:               max(10, options.MaxEnergy),
		maxComboPoints:          options.MaxComboPoints,
		EnergyTickDuration:      unit.ReactionTime,
		EnergyPerTick:           10.0 * unit.ReactionTime.Seconds(),
		energyRegenMultiplier:   1,
		hasteRatingMultiplier:   1,
		regenMetrics:            unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionEnergyRegen}),
		EnergyRefundMetrics:     unit.NewEnergyMetrics(ActionID{OtherID: proto.OtherAction_OtherActionRefund}),
		startingComboPoints:     max(0, min(int32(options.StartingComboPoints), 5)),
		ownerClass:              options.UnitClass,
		comboPointsResourceName: Ternary(options.UnitClass == proto.Class_ClassMonk, "chi", "combo points"),
		hasNoRegen:              options.HasNoRegen,
	}
}

func (unit *Unit) HasEnergyBar() bool {
	return unit.energyBar.unit != nil
}

func (eb *energyBar) CurrentEnergy() float64 {
	return eb.currentEnergy
}

func (eb *energyBar) MaximumEnergy() float64 {
	return eb.maxEnergy
}

func (eb *energyBar) NextEnergyTickAt() time.Duration {
	return eb.nextEnergyTick
}

func (eb *energyBar) MultiplyEnergyRegenSpeed(sim *Simulation, multiplier float64) {
	eb.ResetEnergyTick(sim)
	eb.energyRegenMultiplier *= multiplier
}

func (eb *energyBar) EnergyRegenPerSecond() float64 {
	return 10.0 * eb.hasteRatingMultiplier * eb.energyRegenMultiplier
}

func (eb *energyBar) TimeToTargetEnergy(targetEnergy float64) time.Duration {
	if eb.currentEnergy >= targetEnergy {
		return time.Duration(0)
	}

	return DurationFromSeconds((targetEnergy - eb.currentEnergy) / eb.EnergyRegenPerSecond())
}

func (eb *energyBar) CurrentEnergyRegenMultiplier() float64 {
	return eb.energyRegenMultiplier
}

func (eb *energyBar) AddEnergy(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative energy!")
	}

	newEnergy := min(eb.currentEnergy+amount, eb.maxEnergy)
	metrics.AddEvent(amount, newEnergy-eb.currentEnergy)

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %0.3f energy from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, eb.currentEnergy, newEnergy, eb.maxEnergy)
	}

	eb.currentEnergy = newEnergy
}

func (eb *energyBar) SpendEnergy(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative energy!")
	}

	newEnergy := eb.currentEnergy - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %0.3f energy from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, eb.currentEnergy, newEnergy, eb.maxEnergy)
	}

	eb.currentEnergy = newEnergy
}

func (eb *energyBar) ComboPoints() int32 {
	return eb.comboPoints
}

func (eb *energyBar) SetMaxComboPoints(maxComboPoints int32) {
	eb.maxComboPoints = maxComboPoints
	if eb.maxComboPoints < eb.comboPoints {
		eb.comboPoints = eb.maxComboPoints
	}
}

func (eb *energyBar) MaxComboPoints() int32 {
	return eb.maxComboPoints
}

func (eb *energyBar) IsReset(sim *Simulation) bool {
	return (eb.nextEnergyTick != 0) && (eb.nextEnergyTick-sim.CurrentTime <= eb.EnergyTickDuration)
}

func (eb *energyBar) IsTicking(sim *Simulation) bool {
	return eb.IsReset(sim) && (sim.CurrentTime <= eb.nextEnergyTick) && !eb.hasNoRegen
}

// Gives an immediate partial energy tick and restarts the tick timer.
func (eb *energyBar) ResetEnergyTick(sim *Simulation) {
	if !eb.IsTicking(sim) {
		return
	}

	timeSinceLastTick := max(sim.CurrentTime-(eb.NextEnergyTickAt()-eb.EnergyTickDuration), 0)
	partialTickAmount := (eb.EnergyPerTick * eb.hasteRatingMultiplier * eb.energyRegenMultiplier) * (float64(timeSinceLastTick) / float64(eb.EnergyTickDuration))
	eb.AddEnergy(sim, partialTickAmount, eb.regenMetrics)
	eb.nextEnergyTick = sim.CurrentTime + eb.EnergyTickDuration
	sim.RescheduleTask(eb.nextEnergyTick)
}

func (eb *energyBar) processDynamicHasteRatingChange(sim *Simulation) {
	if eb.unit == nil {
		return
	}

	eb.ResetEnergyTick(sim)
	eb.hasteRatingMultiplier = 1.0 + eb.unit.GetStat(stats.HasteRating)/(100*HasteRatingPerHastePercent)
}

// Used for dynamic updates to maximum Energy, such as from the Druid Primal Madness talent
func (eb *energyBar) UpdateMaxEnergy(sim *Simulation, bonusEnergy float64, metrics *ResourceMetrics) {
	if !eb.IsReset(sim) {
		eb.maxEnergy += bonusEnergy
	} else {
		eb.updateMaxEnergyInternal(sim, bonusEnergy, metrics)
	}
}

func (eb *energyBar) updateMaxEnergyInternal(sim *Simulation, bonusEnergy float64, metrics *ResourceMetrics) {
	// Reset tick timer first so that Energy is properly zeroed out when
	// bonusEnergy < -currentEnergy.
	eb.ResetEnergyTick(sim)

	eb.maxEnergy += bonusEnergy

	if bonusEnergy >= 0 {
		eb.AddEnergy(sim, bonusEnergy, metrics)
	} else {
		eb.SpendEnergy(sim, min(-bonusEnergy, eb.currentEnergy), metrics)
	}
}

func (eb *energyBar) AddComboPoints(sim *Simulation, pointsToAdd int32, metrics *ResourceMetrics) {
	newComboPoints := min(eb.comboPoints+pointsToAdd, eb.maxComboPoints)
	metrics.AddEvent(float64(pointsToAdd), float64(newComboPoints-eb.comboPoints))

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %d %s from %s (%d --> %d) of %0.0f total.", pointsToAdd, eb.comboPointsResourceName, metrics.ActionID, eb.comboPoints, newComboPoints, eb.maxComboPoints)
	}

	eb.comboPoints = newComboPoints
}

func (eb *energyBar) SpendPartialComboPoints(sim *Simulation, pointsToSpend int32, metrics *ResourceMetrics) {
	eb.spendComboPointsInternal(sim, pointsToSpend, metrics)
}

func (eb *energyBar) SpendComboPoints(sim *Simulation, metrics *ResourceMetrics) {
	eb.spendComboPointsInternal(sim, eb.comboPoints, metrics)
}

func (eb *energyBar) spendComboPointsInternal(sim *Simulation, pointsToSpend int32, metrics *ResourceMetrics) {
	pointsToSpend = min(pointsToSpend, eb.comboPoints)
	newComboPoints := eb.comboPoints - pointsToSpend
	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %d %s from %s (%d --> %d) of %0.0f total.", pointsToSpend, eb.comboPointsResourceName, metrics.ActionID, eb.comboPoints, newComboPoints, eb.maxComboPoints)
	}
	metrics.AddEvent(float64(-pointsToSpend), float64(-pointsToSpend))
	eb.comboPoints = newComboPoints
}

func (eb *energyBar) RunTask(sim *Simulation) time.Duration {
	if sim.CurrentTime < eb.nextEnergyTick {
		return eb.nextEnergyTick
	}

	eb.AddEnergy(sim, eb.EnergyPerTick*eb.hasteRatingMultiplier*eb.energyRegenMultiplier, eb.regenMetrics)
	eb.nextEnergyTick = sim.CurrentTime + eb.EnergyTickDuration
	return eb.nextEnergyTick
}

func (eb *energyBar) reset(sim *Simulation) {
	if eb.unit == nil {
		return
	}

	eb.currentEnergy = eb.maxEnergy
	eb.comboPoints = eb.startingComboPoints

	eb.hasteRatingMultiplier = 1.0 + eb.unit.GetStat(stats.HasteRating)/(100*HasteRatingPerHastePercent)
	eb.energyRegenMultiplier = 1.0

	if eb.unit.Type != PetUnit {
		eb.enable(sim, sim.Environment.PrepullStartTime())
	}
}

func (eb *energyBar) enable(sim *Simulation, startAt time.Duration) {
	if eb.hasNoRegen {
		return
	}

	sim.AddTask(eb)
	eb.nextEnergyTick = startAt + time.Duration(sim.RandomFloat("Energy Tick")*float64(eb.EnergyTickDuration))
	sim.RescheduleTask(eb.nextEnergyTick)
}

func (eb *energyBar) disable(sim *Simulation) {
	eb.nextEnergyTick = NeverExpires
	sim.RemoveTask(eb)
}

type EnergyCostOptions struct {
	Cost int32

	Refund        float64
	RefundMetrics *ResourceMetrics // Optional, will default to unit.EnergyRefundMetrics if not supplied.
}
type EnergyCost struct {
	Refund            float64
	RefundMetrics     *ResourceMetrics
	ResourceMetrics   *ResourceMetrics
	ComboPointMetrics *ResourceMetrics
}

func newEnergyCost(spell *Spell, options EnergyCostOptions, energyBar *energyBar) *SpellCost {
	if options.Refund > 0 && options.RefundMetrics == nil {
		options.RefundMetrics = spell.Unit.EnergyRefundMetrics
	}

	return &SpellCost{
		spell:           spell,
		BaseCost:        options.Cost,
		PercentModifier: 1,
		ResourceCostImpl: &EnergyCost{
			Refund:            options.Refund,
			RefundMetrics:     options.RefundMetrics,
			ResourceMetrics:   spell.Unit.NewEnergyMetrics(spell.ActionID),
			ComboPointMetrics: Ternary(energyBar.ownerClass == proto.Class_ClassMonk, spell.Unit.NewChiMetrics(spell.ActionID), spell.Unit.NewComboPointMetrics(spell.ActionID)),
		},
	}
}

func (ec *EnergyCost) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()
	return spell.Unit.CurrentEnergy() >= spell.CurCast.Cost
}
func (ec *EnergyCost) CostFailureReason(_ *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough energy (Current Energy = %0.03f, Energy Cost = %0.03f)", spell.Unit.CurrentEnergy(), spell.CurCast.Cost)
}
func (ec *EnergyCost) SpendCost(sim *Simulation, spell *Spell) {
	spell.Unit.SpendEnergy(sim, spell.CurCast.Cost, ec.ResourceMetrics)
}
func (ec *EnergyCost) IssueRefund(sim *Simulation, spell *Spell) {
	if ec.Refund > 0 && spell.CurCast.Cost > 0 {
		spell.Unit.AddEnergy(sim, ec.Refund*spell.CurCast.Cost, ec.RefundMetrics)
	}
}

func (spell *Spell) EnergyMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*EnergyCost).ResourceMetrics
}

func (spell *Spell) ComboPointMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*EnergyCost).ComboPointMetrics
}
