package druid

import (
	"fmt"

	"github.com/wowsims/cata/sim/core"
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

type EclipseCallback func(eclipse Eclipse, gained bool, sim *core.Simulation)
type eclipseEnergyBar struct {
	druid            *Druid
	lunarEnergy      float64
	solarEnergy      float64
	currentEclipse   Eclipse
	gainMask         EclipseEnergy // which energy the unit is currently allowed to accumulate
	eclipseCallbacks []EclipseCallback
}

func (eb *eclipseEnergyBar) reset() {
	if eb.druid == nil {
		return
	}

	eb.lunarEnergy = 0
	eb.solarEnergy = 0

	// in neutral state we can gain both
	eb.gainMask = SolarEnergy | LunarEnergy
	eb.currentEclipse = NoEclipse
}

func (druid *Druid) EnableEclipseBar() {
	druid.eclipseEnergyBar = eclipseEnergyBar{
		druid:            druid,
		gainMask:         SolarEnergy | LunarEnergy,
		eclipseCallbacks: druid.eclipseEnergyBar.eclipseCallbacks,
	}
}

func (druid *Druid) RegisterEclipseAuras() {
	baselineEclipsePct := 0.25

	lunarSpellMod := druid.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: baselineEclipsePct + (0.16 + druid.GetMasteryPoints()*0.02),
	})

	solarSpellMod := druid.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolNature,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: baselineEclipsePct + (0.16 + druid.GetMasteryPoints()*0.02),
	})

	lunarEclipse := druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48518},
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lunarSpellMod.UpdateFloatValue(baselineEclipsePct + (0.16 + druid.GetMasteryPoints()*0.02))
			lunarSpellMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			lunarSpellMod.Deactivate()
		},
	})

	solarEclipse := druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48517},
		Label:    "Eclipse (Solar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			solarSpellMod.UpdateFloatValue(baselineEclipsePct + (0.16 + druid.GetMasteryPoints()*0.02))
			solarSpellMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			solarSpellMod.Deactivate()
		},
	})

	druid.AddEclipseCallback(func(eclipse Eclipse, gained bool, sim *core.Simulation) {
		if eclipse == LunarEclipse {
			if gained {
				lunarEclipse.Activate(sim)
			} else {
				lunarEclipse.Deactivate(sim)
			}
		} else {
			if gained {
				solarEclipse.Activate(sim)
			} else {
				solarEclipse.Deactivate(sim)
			}
		}
	})
}

func (druid *Druid) RegisterEclipseEnergyGainAura() {
	solarMetric := druid.NewSolarEnergyMetrics(core.ActionID{SpellID: 89265})
	lunarMetric := druid.NewLunarEnergyMetrics(core.ActionID{SpellID: 89265})

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 89265},
		Label:    "Eclipse Energy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			var eclipseEnergyMultiplier float64 = 1.0

			if druid.canEuphoriaProc() && druid.hasEuphoriaProcced(sim) {
				eclipseEnergyMultiplier = 2
			}

			var energyGain = druid.GetSpellEclipseEnergy(spell.ActionID.SpellID, druid.currentEclipse != NoEclipse)

			switch spell.ActionID.SpellID {
			case 2912: // Starfire
				druid.AddEclipseEnergy(energyGain*eclipseEnergyMultiplier, SolarEnergy, sim, solarMetric)
			case 5176: // Wrath
				druid.AddEclipseEnergy(energyGain*eclipseEnergyMultiplier, LunarEnergy, sim, lunarMetric)
			case 78674: // Starsurge
				if druid.CanGainEnergy(SolarEnergy) {
					druid.AddEclipseEnergy(energyGain, SolarEnergy, sim, solarMetric)
				} else {
					druid.AddEclipseEnergy(energyGain, LunarEnergy, sim, lunarMetric)
				}
			}
		},
	})
}

func (druid *Druid) hasEuphoriaProcced(sim *core.Simulation) bool {
	return sim.Proc(0.12*float64(druid.Talents.Euphoria), fmt.Sprintf("Euphoria %d/2", druid.Talents.Euphoria))
}

func (druid *Druid) canEuphoriaProc() bool {
	if druid.Talents.Euphoria == 0 {
		return false
	}

	if druid.currentEclipse != NoEclipse {
		return false
	}

	if druid.Talents.Euphoria == 1 {
		return true
	}

	if druid.Talents.Euphoria == 2 {
		if druid.CanGainEnergy(SolarEnergy) && druid.CurrentSolarEnergy() <= 35 {
			return true
		}

		if druid.CanGainEnergy(LunarEnergy) && druid.CurrentLunarEnergy() <= 35 {
			return true
		}
	}

	return false
}

func (druid *Druid) HasEclipseBar() bool {
	return druid.eclipseEnergyBar.druid != nil
}

func (eb *eclipseEnergyBar) AddEclipseCallback(callback EclipseCallback) {
	eb.eclipseCallbacks = append(eb.eclipseCallbacks, callback)
}

func (eb *eclipseEnergyBar) AddEclipseEnergy(amount float64, kind EclipseEnergy, sim *core.Simulation, metrics *core.ResourceMetrics) {
	if eb.druid == nil {
		return
	}

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

func (eb *eclipseEnergyBar) CanGainEnergy(kind EclipseEnergy) bool {
	return eb.gainMask&kind > 0
}

// spends the given amount of energy and returns how much energy remains
// this might be added to the solar energy
func (eb *eclipseEnergyBar) spendLunarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics) float64 {
	if amount == 0 || eb.lunarEnergy == 0 {
		return amount
	}

	spend := min(amount, eb.lunarEnergy)
	remainder := amount - spend
	old := eb.lunarEnergy
	eb.lunarEnergy -= spend

	if sim.Log != nil {
		eb.druid.Log(sim, "Spent %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.lunarEnergy, 100.0)
	}

	if eb.lunarEnergy == 0 {
		eb.SetEclipse(NoEclipse, sim)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addLunarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics) {
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
		eb.druid.Log(sim, "Gained %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.lunarEnergy, 100.0)
	}

	if eb.lunarEnergy == 100 {
		eb.SetEclipse(LunarEclipse, sim)
	}

	metrics.AddEvent(amount, gain)
}

func (eb *eclipseEnergyBar) ForceEclipse(eclipse Eclipse, sim *core.Simulation) {
	eb.reset()

	if eclipse == SolarEclipse {
		eb.solarEnergy = 100
	} else if eclipse == LunarEclipse {
		eb.lunarEnergy = 100
	}

	eb.SetEclipse(eclipse, sim)
}

func (eb *eclipseEnergyBar) SetEclipse(eclipse Eclipse, sim *core.Simulation) {
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

func (eb *eclipseEnergyBar) invokeCallback(eclipse Eclipse, gained bool, sim *core.Simulation) {
	for _, callback := range eb.eclipseCallbacks {
		callback(eclipse, gained, sim)
	}
}

func (eb *eclipseEnergyBar) spendSolarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics) float64 {
	if amount == 0 || eb.solarEnergy == 0 {
		return amount
	}

	spend := min(amount, eb.solarEnergy)
	remainder := amount - spend
	old := eb.solarEnergy
	eb.solarEnergy -= spend

	if sim.Log != nil {
		eb.druid.Log(sim, "Spent %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.solarEnergy, 100.0)
	}

	if eb.solarEnergy == 0 {
		eb.SetEclipse(NoEclipse, sim)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addSolarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics) {
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
		eb.druid.Log(sim, "Gained %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.solarEnergy, 100.0)
	}

	if eb.solarEnergy == 100 {
		eb.SetEclipse(SolarEclipse, sim)
	}

	metrics.AddEvent(amount, gain)
}

func (unit *Druid) NewSolarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeSolarEnergy)
}

func (unit *Druid) NewLunarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeLunarEnergy)
}
