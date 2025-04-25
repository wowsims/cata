package druid

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
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
	eclipseTrigger   func(spell *core.Spell) bool // used to deactivate eclipse spells
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

func getEclipseMasteryBonus(masteryPoints float64) float64 {
	return (16 + masteryPoints*2) / 100
}

func (druid *Druid) RegisterEclipseAuras() {
	baselineEclipsePct := 0.25
	initialEclipseMasteryBonus := getEclipseMasteryBonus(druid.GetMasteryPoints())

	lunarSpellMod := druid.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: baselineEclipsePct + initialEclipseMasteryBonus,
	})

	solarSpellMod := druid.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolNature,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: baselineEclipsePct + initialEclipseMasteryBonus,
	})

	lunarEclipse := druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48518},
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			lunarSpellMod.UpdateFloatValue(baselineEclipsePct + getEclipseMasteryBonus(druid.GetMasteryPoints()))
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
			solarSpellMod.UpdateFloatValue(baselineEclipsePct + getEclipseMasteryBonus(druid.GetMasteryPoints()))
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
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			var eclipseEnergyMultiplier float64 = 1.0

			if druid.canEuphoriaProc(spell) && druid.hasEuphoriaProcced(sim) {
				eclipseEnergyMultiplier = 2
			}

			if energyGain := druid.GetSpellEclipseEnergy(spell.ClassSpellMask, druid.currentEclipse != NoEclipse); energyGain != 0 {
				switch spell.ClassSpellMask {
				case DruidSpellStarfire:
					druid.AddEclipseEnergy(energyGain*eclipseEnergyMultiplier, SolarEnergy, sim, solarMetric, spell)
				case DruidSpellWrath:
					druid.AddEclipseEnergy(energyGain*eclipseEnergyMultiplier, LunarEnergy, sim, lunarMetric, spell)
				case DruidSpellStarsurge:
					if druid.CanGainEnergy(SolarEnergy) {
						druid.AddEclipseEnergy(energyGain, SolarEnergy, sim, solarMetric, spell)
					} else {
						druid.AddEclipseEnergy(energyGain, LunarEnergy, sim, lunarMetric, spell)
					}
				case DruidSpellMoonfire: // Moonfire (under the effect of Lunar Shower)
					druid.AddEclipseEnergy(energyGain, SolarEnergy, sim, solarMetric, spell)
				case DruidSpellSunfire: // Sunfire (under the effect of Lunar Shower)
					druid.AddEclipseEnergy(energyGain, LunarEnergy, sim, lunarMetric, spell)
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// chekc if trigger is supposed to handle spell hit, then clear
			if druid.eclipseTrigger != nil && druid.eclipseTrigger(spell) {
				druid.eclipseTrigger = nil
			}
		},
	})
}

func (druid *Druid) hasEuphoriaProcced(sim *core.Simulation) bool {
	return sim.Proc(0.12*float64(druid.Talents.Euphoria), fmt.Sprintf("Euphoria %d/2", druid.Talents.Euphoria))
}

func (druid *Druid) canEuphoriaProc(spell *core.Spell) bool {
	if druid.Talents.Euphoria == 0 {
		return false
	}

	if druid.currentEclipse != NoEclipse {
		return false
	}

	if spell.ClassSpellMask != DruidSpellStarfire && spell.ClassSpellMask != DruidSpellWrath {
		return false
	}

	if druid.Talents.Euphoria == 1 {
		return true
	}

	if druid.Talents.Euphoria == 2 {
		if druid.CanGainEnergy(SolarEnergy) && druid.CurrentSolarEnergy() <= 35 && druid.CurrentLunarEnergy() == 0 {
			return true
		}

		if druid.CanGainEnergy(LunarEnergy) && druid.CurrentLunarEnergy() <= 35 && druid.CurrentSolarEnergy() == 0 {
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

func (eb *eclipseEnergyBar) AddEclipseEnergy(amount float64, kind EclipseEnergy, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) {
	if eb.druid == nil {
		return
	}

	// unit currently can not gain the specified energy
	if kind&eb.gainMask == 0 {
		return
	}

	if kind&SolarEnergy > 0 {
		remainder := eb.spendLunarEnergy(amount, sim, metrics, spell)
		eb.addSolarEnergy(remainder, sim, metrics, spell)
		return
	}

	remainder := eb.spendSolarEnergy(amount, sim, metrics, spell)
	eb.addLunarEnergy(remainder, sim, metrics, spell)
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
func (eb *eclipseEnergyBar) spendLunarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) float64 {
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
		eb.SetEclipse(NoEclipse, sim, spell)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addLunarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) {
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
		eb.SetEclipse(LunarEclipse, sim, spell)
	}

	metrics.AddEvent(amount, gain)
}

func (eb *eclipseEnergyBar) SetEclipse(eclipse Eclipse, sim *core.Simulation, spell *core.Spell) {
	if eb.currentEclipse == eclipse {
		return
	}

	if eclipse == LunarEclipse {
		eb.gainMask = SolarEnergy
		eb.invokeCallback(eclipse, true, sim)
		eb.currentEclipse = eclipse
	} else if eclipse == SolarEclipse {
		eb.gainMask = LunarEnergy
		eb.invokeCallback(eclipse, true, sim)
		eb.currentEclipse = eclipse
	} else {
		if spell.ClassSpellMask&DruidSpellWrath > 0 {
			eb.eclipseTrigger = func(triggerSpell *core.Spell) bool {
				if spell.ClassSpellMask == triggerSpell.ClassSpellMask && !(triggerSpell.ProcMask&core.ProcMaskSpellProc > 0) {
					eb.invokeCallback(eb.currentEclipse, false, sim)
					eb.currentEclipse = eclipse
					return true
				}

				return false
			}

			// eclipse state is only removed for non procced spells
		} else if !(spell.ProcMask&core.ProcMaskSpellProc > 0) {
			sim.AddPendingAction(&core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Millisecond*10,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(sim *core.Simulation) {
					// make sure we're not triggering twice for edge cases
					if eb.currentEclipse == eclipse {
						return
					}

					eb.invokeCallback(eb.currentEclipse, false, sim)
					eb.currentEclipse = eclipse
				},
			})
		}
	}
}

func (eb *eclipseEnergyBar) invokeCallback(eclipse Eclipse, gained bool, sim *core.Simulation) {
	for _, callback := range eb.eclipseCallbacks {
		callback(eclipse, gained, sim)
	}
}

func (eb *eclipseEnergyBar) spendSolarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) float64 {
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
		eb.SetEclipse(NoEclipse, sim, spell)
	}

	return remainder
}

func (eb *eclipseEnergyBar) addSolarEnergy(amount float64, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) {
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
		eb.SetEclipse(SolarEclipse, sim, spell)
	}

	metrics.AddEvent(amount, gain)
}

func (unit *Druid) NewSolarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeSolarEnergy)
}

func (unit *Druid) NewLunarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeLunarEnergy)
}
