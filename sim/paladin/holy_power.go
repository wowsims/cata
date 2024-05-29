package paladin

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

type HolyPowerBar struct {
	paladin *Paladin

	DivinePurpose bool
	holyPower     int32
}

func (paladin *Paladin) CurrentHolyPower() int32 {
	return paladin.holyPower
}

func (paladin *Paladin) GetHolyPowerValue() int32 {
	if paladin.DivinePurpose {
		return 3
	}

	return paladin.CurrentHolyPower()
}

func (paladin *Paladin) InitializeHolyPowerBar() {
	paladin.HolyPowerBar = HolyPowerBar{
		paladin:   paladin,
		holyPower: 0,
	}
}

func (pb *HolyPowerBar) Reset() {
	if pb.paladin == nil {
		return
	}

	pb.holyPower = 0
}

func (paladin *Paladin) HasHolyPowerBar() bool {
	return paladin.HolyPowerBar.paladin != nil
}

func (pb *HolyPowerBar) GainHolyPower(sim *core.Simulation, amountToAdd int32, metrics *core.ResourceMetrics) {
	if pb.paladin == nil {
		return
	}

	newHolyPower := min(pb.holyPower+amountToAdd, 3)
	metrics.AddEvent(float64(newHolyPower), float64(newHolyPower-pb.holyPower))

	if sim.Log != nil {
		pb.paladin.Log(sim, "Gained %d holy power from %s (%d --> %d) of %0.0f total.", newHolyPower, metrics.ActionID, pb.holyPower, newHolyPower, 3.0)
	}

	pb.holyPower = newHolyPower
}

func (pb *HolyPowerBar) SpendHolyPower(sim *core.Simulation, metrics *core.ResourceMetrics) {
	if pb.paladin == nil {
		return
	}

	if pb.DivinePurpose {
		pb.paladin.Log(sim, "Consumed Divine Purpose")
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
