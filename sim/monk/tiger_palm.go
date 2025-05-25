package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

/*
Tooltip:
Attack with the palm of your hand, dealing

-- Teachings of the Monastery --

	${6*$<low>} to ${6*$<high>}

-- else --

	${3*$<low>} to ${3*$<high>}

--

	damage.

Also grants you Tiger Power, causing your attacks to ignore 30% of enemies' armor for 20 sec.
*/
var tigerPalmActionID = core.ActionID{SpellID: 100787}
var tigerPowerActionID = core.ActionID{SpellID: 125359}

func tigerPowerBuffConfig(monk *Monk, isSEFClone bool) core.Aura {
	config := core.Aura{
		Label:    "Tiger Power" + monk.Label,
		ActionID: tigerPowerActionID,
		Duration: time.Second * 20,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range sim.Encounter.TargetUnits {
				aura.Unit.AttackTables[target.UnitIndex].ArmorIgnoreFactor += 0.3
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, target := range sim.Encounter.TargetUnits {
				aura.Unit.AttackTables[target.UnitIndex].ArmorIgnoreFactor -= 0.3
			}
		},
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
	}

	return config
}

func tigerPalmSpellConfig(monk *Monk, isSEFClone bool, overrides core.SpellConfig) core.SpellConfig {
	config := core.SpellConfig{
		ActionID:       tigerPalmActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellTigerPalm,
		MaxRange:       core.MaxMeleeRange,

		Cast: overrides.Cast,

		DamageMultiplier: 3.0,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ExtraCastCondition: overrides.ExtraCastCondition,

		ApplyEffects: overrides.ApplyEffects,
	}

	if isSEFClone {
		config.ActionID = config.ActionID.WithTag(SEFSpellID)
		config.Flags &= ^(core.SpellFlagAPL | SpellFlagSpender)
	}

	return config
}

func (monk *Monk) registerTigerPalm() {
	chiMetrics := monk.NewChiMetrics(tigerPalmActionID)
	isBrewmaster := monk.Spec == proto.Spec_SpecBrewmasterMonk

	tigerPowerBuff := monk.RegisterAura(tigerPowerBuffConfig(monk, false))

	monk.RegisterSpell(tigerPalmSpellConfig(monk, false, core.SpellConfig{
		ActionID:       tigerPalmActionID,
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		Flags:          core.SpellFlagMeleeMetrics | SpellFlagSpender | core.SpellFlagAPL,
		ClassSpellMask: MonkSpellTigerPalm,
		MaxRange:       core.MaxMeleeRange,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 3.0,
		ThreatMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return isBrewmaster || monk.GetChi() >= 1 || monk.ComboBreakerTigerPalmAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				tigerPowerBuff.Activate(sim)

				if monk.ComboBreakerTigerPalmAura.IsActive() || isBrewmaster {
					monk.SpendChi(sim, 0, chiMetrics)
				} else {
					monk.SpendChi(sim, 1, chiMetrics)
				}
			}

			spell.DealOutcome(sim, result)
		},
	}))
}

func (pet *StormEarthAndFirePet) registerSEFTigerPalm() {
	tigerPowerBuff := pet.RegisterAura(tigerPowerBuffConfig(pet.owner, true))

	pet.RegisterSpell(tigerPalmSpellConfig(pet.owner, true, core.SpellConfig{
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := pet.owner.CalculateMonkStrikeDamage(sim, spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				tigerPowerBuff.Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},
	}))
}
