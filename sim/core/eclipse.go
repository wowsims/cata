package core

import (
	"github.com/wowsims/cata/sim/core/proto"
)

/*
  unit specific balance energy bar
*/

type EclipseEnergy byte

const (
	SolarEnergy         EclipseEnergy = 1
	LunarEnergy         EclipseEnergy = 2
	SolarAndLunarEnergy               = SolarEnergy | LunarEnergy
)

type Eclipse byte

const (
	NoEclipse    Eclipse = 0
	SolarEclipse Eclipse = 1
	LunarEclipse Eclipse = 2
)

type EclipseCallback func(eclipse Eclipse, gained bool, sim *Simulation)
type eclipseEnergyBar struct {
	unit             *Unit
	lunarEnergy      float64
	solarEnergy      float64
	currentEclipse   Eclipse
	gainMask         EclipseEnergy // which energy the unit is currently allowed to accumulate
	eclipseCallbacks []EclipseCallback
}

func (eb *eclipseEnergyBar) reset() {
	if eb.unit == nil {
		return
	}

	eb.lunarEnergy = 0
	eb.solarEnergy = 0

	// in neutral state we can gain both
	eb.gainMask = SolarEnergy | LunarEnergy
	eb.currentEclipse = NoEclipse
}

func (unit *Unit) EnableEclipseBar() {
	unit.eclipseEnergyBar = eclipseEnergyBar{
		unit:     unit,
		gainMask: SolarEnergy | LunarEnergy,
	}
}

func (unit *Unit) HasEclipseBar() bool {
	return unit.eclipseEnergyBar.unit != nil
}

func (eb *eclipseEnergyBar) AddEclipseCallback(callback EclipseCallback) {
	eb.eclipseCallbacks = append(eb.eclipseCallbacks, callback)
}

func (eb *eclipseEnergyBar) AddEclipseEnergy(amount float64, kind EclipseEnergy, sim *Simulation, metrics *ResourceMetrics) {

	// unit currently can not gain the specified energy
	if kind&eb.gainMask == 0 {
		return
	}

	if kind&SolarEnergy > 0 {
		remainder := eb.spendLunarEnergy(amount, sim, metrics)
		eb.addSolarEnergy(remainder, sim, metrics)
		return
	}

	remainder := eb.spendSolarEnergy(amount, sim, metrics)
	eb.addLunarEnergy(remainder, sim, metrics)
}

func (eb *eclipseEnergyBar) CurrentSolarEnergy() int32 {
	return int32(eb.solarEnergy)
}

func (eb *eclipseEnergyBar) CurrentLunarEnergy() int32 {
	return int32(eb.lunarEnergy)
}

// spends the given amount of energy and returns how much energy remains
// this might be added to the solar energy
func (eb *eclipseEnergyBar) spendLunarEnergy(amount float64, sim *Simulation, metrics *ResourceMetrics) float64 {
	if amount == 0 || eb.lunarEnergy == 0 {
		return amount
	}

	spend := min(amount, eb.lunarEnergy)
	remainder := amount - spend
	old := eb.lunarEnergy
	eb.lunarEnergy -= spend

	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.lunarEnergy, 100.0)
	}

	if eb.lunarEnergy == 0 {
		eb.SetEclipse(NoEclipse, sim)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addLunarEnergy(amount float64, sim *Simulation, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Tried to add negative amount of lunar energy.")
	}

	if amount == 0 {
		return
	}

	gain := min(eb.lunarEnergy+amount, 100.0) - eb.lunarEnergy

	old := eb.lunarEnergy
	eb.lunarEnergy += gain

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.lunarEnergy, 100.0)
	}

	if eb.lunarEnergy == 100 {
		eb.SetEclipse(LunarEclipse, sim)
	}

	metrics.AddEvent(amount, gain)
}

func (eb *eclipseEnergyBar) SetEclipse(eclipse Eclipse, sim *Simulation) {
	if eb.currentEclipse == eclipse {
		return
	}

	if eclipse == LunarEclipse {
		eb.gainMask = SolarEnergy
		eb.invokeCallback(eclipse, true, sim)
	} else if eclipse == SolarEclipse {
		eb.gainMask = LunarEnergy
		eb.invokeCallback(eclipse, true, sim)
	} else {
		eb.invokeCallback(eb.currentEclipse, false, sim)
	}

	eb.currentEclipse = eclipse
}

func (eb *eclipseEnergyBar) invokeCallback(eclipse Eclipse, gained bool, sim *Simulation) {
	for _, callback := range eb.eclipseCallbacks {
		callback(eclipse, gained, sim)
	}
}

func (eb *eclipseEnergyBar) spendSolarEnergy(amount float64, sim *Simulation, metrics *ResourceMetrics) float64 {
	if amount == 0 || eb.solarEnergy == 0 {
		return amount
	}

	spend := min(amount, eb.solarEnergy)
	remainder := amount - spend
	old := eb.solarEnergy
	eb.solarEnergy -= spend

	if sim.Log != nil {
		eb.unit.Log(sim, "Spent %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.solarEnergy, 100.0)
	}

	if eb.solarEnergy == 0 {
		eb.SetEclipse(NoEclipse, sim)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addSolarEnergy(amount float64, sim *Simulation, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Tried to add negative amount of solar energy.")
	}

	if amount == 0 {
		return
	}

	gain := min(eb.solarEnergy+amount, 100.0) - eb.solarEnergy

	old := eb.solarEnergy
	eb.solarEnergy += gain

	if sim.Log != nil {
		eb.unit.Log(sim, "Gained %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.solarEnergy, 100.0)
	}

	if eb.solarEnergy == 100 {
		eb.SetEclipse(SolarEclipse, sim)
	}

	metrics.AddEvent(amount, gain)
}

func (unit *Unit) NewSolarEnergyMetric(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeSolarEnergy)
}

func (unit *Unit) NewLunarEnergyMetrics(actionID ActionID) *ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeLunarEnergy)
}
