package balance

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/druid"
)

/*
  unit specific balance energy bar
*/

type EclipseEnergy byte

const (
	NoEnergy            EclipseEnergy = 0
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
	moonkin          *BalanceDruid
	lunarEnergy      float64
	solarEnergy      float64
	currentEclipse   Eclipse
	gainMask         EclipseEnergy // which energy the unit is currently allowed to accumulate
	previousGainMask EclipseEnergy // used to restore gain mask after CA
	eclipseCallbacks []EclipseCallback
	eclipseTrigger   func(spell *core.Spell) bool // used to deactivate eclipse spells
}

func (eb *eclipseEnergyBar) reset() {
	if eb.moonkin == nil {
		return
	}

	eb.lunarEnergy = 0
	eb.solarEnergy = 0

	// in neutral state we can gain both
	eb.gainMask = SolarAndLunarEnergy
	eb.currentEclipse = NoEclipse
}

func (eb *eclipseEnergyBar) resetWithMask(gainMask EclipseEnergy) {
	eb.reset()
	eb.gainMask = gainMask
}

func (moonkin *BalanceDruid) EnableEclipseBar() {
	moonkin.eclipseEnergyBar = eclipseEnergyBar{
		moonkin:          moonkin,
		gainMask:         SolarAndLunarEnergy,
		eclipseCallbacks: moonkin.eclipseEnergyBar.eclipseCallbacks,
	}
}

func (moonkin *BalanceDruid) SuspendEclipseBar() {
	moonkin.eclipseEnergyBar.previousGainMask = moonkin.eclipseEnergyBar.gainMask
	moonkin.eclipseEnergyBar.resetWithMask(NoEnergy)
}

func (moonkin *BalanceDruid) RestoreEclipseBar() {
	moonkin.eclipseEnergyBar.resetWithMask(moonkin.eclipseEnergyBar.previousGainMask)
}

func (moonkin *BalanceDruid) ActivateEclipse(eclipse Eclipse, sim *core.Simulation) {
	moonkin.eclipseEnergyBar.invokeCallback(eclipse, true, sim)
}

func (moonkin *BalanceDruid) DeactivateEclipse(eclipse Eclipse, sim *core.Simulation) {
	moonkin.eclipseEnergyBar.invokeCallback(eclipse, false, sim)
}

func getEclipseMasteryBonus(masteryPoints float64) float64 {
	return (15.0 + (8.0 * 1.875) + (masteryPoints * 1.875)) / 100
}

func (moonkin *BalanceDruid) RegisterEclipseAuras() {
	manaMetrics := moonkin.NewManaMetrics(core.ActionID{SpellID: 79577 /* Eclipse */})

	docEclipseMasteryBonus := 0.0
	eclipseMasteryBonus := getEclipseMasteryBonus(moonkin.GetMasteryPoints())

	lunarSpellMod := moonkin.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane,
		ProcMask:   core.ProcMaskSpellDamage,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: docEclipseMasteryBonus + eclipseMasteryBonus,
	})

	solarSpellMod := moonkin.AddDynamicMod(core.SpellModConfig{
		School:     core.SpellSchoolNature,
		ProcMask:   core.ProcMaskSpellDamage,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: docEclipseMasteryBonus + eclipseMasteryBonus,
	})

	moonkin.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		if !moonkin.IsInEclipse() {
			return
		}

		masteryBonusDiff := core.MasteryRatingToMasteryPoints(newMastery) - core.MasteryRatingToMasteryPoints(oldMastery)

		if lunarSpellMod.IsActive {
			lunarSpellMod.UpdateFloatValue(lunarSpellMod.GetFloatValue() + getEclipseMasteryBonus(masteryBonusDiff))
		} else if solarSpellMod.IsActive {
			solarSpellMod.UpdateFloatValue(solarSpellMod.GetFloatValue() + getEclipseMasteryBonus(masteryBonusDiff))
		}
	})

	lunarEclipse := moonkin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48518},
		Label:    "Eclipse (Lunar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if moonkin.DreamOfCenarius.IsActive() {
				docEclipseMasteryBonus = 0.25
				moonkin.DreamOfCenarius.Deactivate(sim)
			}

			lunarSpellMod.UpdateFloatValue(docEclipseMasteryBonus + getEclipseMasteryBonus(moonkin.GetMasteryPoints()))
			lunarSpellMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			docEclipseMasteryBonus = 0.0
			lunarSpellMod.Deactivate()
		},
	})

	solarEclipse := moonkin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 48517},
		Label:    "Eclipse (Solar)",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if moonkin.DreamOfCenarius.IsActive() {
				docEclipseMasteryBonus = 0.25
				moonkin.DreamOfCenarius.Deactivate(sim)
			}

			solarSpellMod.UpdateFloatValue(docEclipseMasteryBonus + getEclipseMasteryBonus(moonkin.GetMasteryPoints()))
			solarSpellMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			docEclipseMasteryBonus = 0.0
			solarSpellMod.Deactivate()
		},
	})

	moonkin.AddEclipseCallback(func(eclipse Eclipse, gained bool, sim *core.Simulation) {
		if gained {
			// Moonkins are energized for 50% maximum mana every time they enter eclipse.
			moonkin.AddMana(sim, moonkin.MaxMana()*0.5, manaMetrics)
		}

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

func (moonkin *BalanceDruid) RegisterEclipseEnergyGainAura() {
	solarMetric := moonkin.NewSolarEnergyMetrics(core.ActionID{SpellID: 89265})
	lunarMetric := moonkin.NewLunarEnergyMetrics(core.ActionID{SpellID: 89265})

	moonkin.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 89265},
		Label:    "Eclipse Energy",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {

			if energyGain := moonkin.GetSpellEclipseEnergy(spell.ClassSpellMask, moonkin.currentEclipse != NoEclipse); energyGain != 0 {
				switch spell.ClassSpellMask {
				case druid.DruidSpellStarfire:
					moonkin.AddEclipseEnergy(energyGain, SolarEnergy, sim, solarMetric, spell)
				case druid.DruidSpellWrath:
					moonkin.AddEclipseEnergy(energyGain, LunarEnergy, sim, lunarMetric, spell)
				case druid.DruidSpellStarsurge:
					if moonkin.CanGainEnergy(SolarAndLunarEnergy) {
						moonkin.AddEclipseEnergy(energyGain, LunarEnergy, sim, solarMetric, spell)
					} else {
						moonkin.AddEclipseEnergy(energyGain, SolarEnergy, sim, lunarMetric, spell)
					}
				}
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			// chekc if trigger is supposed to handle spell hit, then clear
			if moonkin.eclipseTrigger != nil && moonkin.eclipseTrigger(spell) {
				moonkin.eclipseTrigger = nil
			}
		},
	})
}

func (moonkin *BalanceDruid) HasEclipseBar() bool {
	return moonkin.eclipseEnergyBar.moonkin != nil
}

func (moonkin *BalanceDruid) IsInEclipse() bool {
	return moonkin.currentEclipse != NoEclipse
}

func (eb *eclipseEnergyBar) AddEclipseCallback(callback EclipseCallback) {
	eb.eclipseCallbacks = append(eb.eclipseCallbacks, callback)
}

func (eb *eclipseEnergyBar) AddEclipseEnergy(amount float64, kind EclipseEnergy, sim *core.Simulation, metrics *core.ResourceMetrics, spell *core.Spell) {
	if eb.moonkin == nil {
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

func (eb *eclipseEnergyBar) StoreGainMaskAndSuspend() {
	eb.previousGainMask = eb.gainMask
	eb.gainMask = NoEnergy
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
		eb.moonkin.Log(sim, "Spent %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.lunarEnergy, 100.0)
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
		eb.moonkin.Log(sim, "Gained %0.0f lunar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.lunarEnergy, 100.0)
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
		if spell.Matches(druid.DruidSpellWrath) {
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
		eb.moonkin.Log(sim, "Spent %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", spend, metrics.ActionID, old, eb.solarEnergy, 100.0)
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
		eb.moonkin.Log(sim, "Gained %0.0f solar energy from %s (%0.0f --> %0.0f) of %0.0f total.", gain, metrics.ActionID, old, eb.solarEnergy, 100.0)
	}

	if eb.solarEnergy == 100 {
		eb.SetEclipse(SolarEclipse, sim, spell)
	}

	metrics.AddEvent(amount, gain)
}

func (unit *BalanceDruid) NewSolarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeSolarEnergy)
}

func (unit *BalanceDruid) NewLunarEnergyMetrics(actionID core.ActionID) *core.ResourceMetrics {
	return unit.Metrics.NewResourceMetrics(actionID, proto.ResourceType_ResourceTypeLunarEnergy)
}
