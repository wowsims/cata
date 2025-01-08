package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

type HolyPowerBar struct {
	paladin *Paladin

	holyPower int32
}

// CurrentHolyPower returns the actual amount of holy power the paladin has, not counting the Divine Purpose proc.
func (paladin *Paladin) CurrentHolyPower() int32 {
	return paladin.holyPower
}

// GetHolyPowerValue returns the amount of holy power used for calculating the damage done by Templar's Verdict and duration of Inquisition.
func (paladin *Paladin) GetHolyPowerValue() int32 {
	if paladin.DivinePurposeAura.IsActive() {
		return 3
	}

	return paladin.CurrentHolyPower()
}

func (paladin *Paladin) initializeHolyPowerBar() {
	paladin.HolyPowerBar = HolyPowerBar{
		paladin:   paladin,
		holyPower: paladin.StartingHolyPower,
	}
}

func (pb *HolyPowerBar) Reset() {
	if pb.paladin == nil {
		return
	}

	pb.holyPower = pb.paladin.StartingHolyPower
}

func (paladin *Paladin) HasHolyPowerBar() bool {
	return paladin.HolyPowerBar.paladin != nil
}

func (pb *HolyPowerBar) GainHolyPower(sim *core.Simulation, amountToAdd int32, metrics *core.ResourceMetrics) {
	if pb.paladin == nil {
		return
	}

	newHolyPower := min(pb.holyPower+amountToAdd, 3)
	metrics.AddEvent(float64(amountToAdd), float64(newHolyPower-pb.holyPower))

	if sim.Log != nil {
		pb.paladin.Log(sim, "Gained %d holy power from %s (%d --> %d) of %0.0f total.", amountToAdd, metrics.ActionID, pb.holyPower, newHolyPower, 3.0)
	}

	pb.holyPower = newHolyPower
}

func (pb *HolyPowerBar) SpendHolyPower(sim *core.Simulation, metrics *core.ResourceMetrics) {
	if pb.paladin == nil {
		return
	}

	if pb.paladin.DivinePurposeAura.IsActive() {
		// Aura deactivation handled in talents_retribution.go:applyDivinePurpose()
		return
	}

	if sim.Log != nil {
		pb.paladin.Log(sim, "Spent %d holy power from %s (%d --> %d) of %0.0f total.", pb.holyPower, metrics.ActionID, pb.holyPower, 0, 3.0)
	}

	metrics.AddEvent(float64(-pb.holyPower), float64(-pb.holyPower))
	pb.holyPower = 0
}

func (unit *Paladin) NewHolyPowerMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeHolyPower)
}
