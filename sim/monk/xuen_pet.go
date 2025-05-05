package monk

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type Xuen struct {
	core.Pet

	Monk                    *Monk
	CracklingTigerLightning *core.Spell
}

var baseStats = stats.Stats{
	stats.Strength:    0,
	stats.Agility:     0,
	stats.Stamina:     0,
	stats.Intellect:   0,
	stats.AttackPower: 1141,
	stats.Mana:        0,
}

func (monk *Monk) NewXuen() *Xuen {
	xuen := &Xuen{
		Pet:  core.NewPet("Xuen, The White Tiger", &monk.Character, baseStats, monk.xuenStatInheritance(), false, false),
		Monk: monk,
	}

	xuen.OnPetEnable = func(sim *core.Simulation) {
		xuen.AutoAttacks.PauseMeleeBy(sim, time.Duration(1))
	}

	xuen.DelayInitialInheritance(time.Millisecond * 500)

	actionID := core.ActionID{SpellID: 123996}
	xuen.CracklingTigerLightning = xuen.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		DamageMultiplier: 1,
		CritMultiplier:   monk.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Millisecond * 500,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    xuen.NewTimer(),
				Duration: time.Second * 1,
			},
		},

		BonusCoefficient: 0.505,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := monk.CalcScalingSpellDmg(0.293) + xuen.GetStat(stats.AttackPower)*spell.BonusCoefficient

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	xuen.PseudoStats.DamageTakenMultiplier *= 0.1
	xuen.PseudoStats.DamageDealtMultiplier = monk.PseudoStats.DamageDealtMultiplier

	xuen.EnableAutoAttacks(xuen, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        monk.CalcScalingSpellDmg(1), // Currently does 1240 in bugged state
			BaseDamageMax:        monk.CalcScalingSpellDmg(1), // Currently does 1241 in bugged state
			SwingSpeed:           1,
			NormalizedSwingSpeed: 1,
			CritMultiplier:       monk.DefaultCritMultiplier(),
			SpellSchool:          core.SpellSchoolNature,
		},
		AutoSwingMelee: true,
	})

	xuen.AutoAttacks.MHConfig().BonusCoefficient = 0
	xuen.AutoAttacks.MHConfig().Flags |= core.SpellFlagIgnoreTargetModifiers

	monk.RegisterOnStanceChanged(func(sim *core.Simulation, _ Stance) {
		xuen.PseudoStats.DamageDealtMultiplier = monk.PseudoStats.DamageDealtMultiplier
	})

	monk.AddPet(xuen)

	return xuen
}

func (monk *Monk) xuenStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.PhysicalCritPercent: ownerStats[stats.PhysicalCritPercent],
			stats.Stamina:             ownerStats[stats.Stamina],
			stats.Intellect:           ownerStats[stats.Intellect] * 0.3,
			stats.AttackPower:         ownerStats[stats.AttackPower],
			stats.PhysicalHitPercent:  ownerStats[stats.PhysicalHitPercent],
			stats.HasteRating:         ownerStats[stats.HasteRating],
		}
	}
}

func (xuen *Xuen) Initialize() {
}

func (xuen *Xuen) ExecuteCustomRotation(sim *core.Simulation) {
	if xuen.CracklingTigerLightning.CanCast(sim, xuen.CurrentTarget) {
		xuen.CracklingTigerLightning.Cast(sim, xuen.CurrentTarget)
	}
}

func (xuen *Xuen) Reset(sim *core.Simulation) {
	xuen.Disable(sim)
}

func (xuen *Xuen) OnPetDisable(sim *core.Simulation) {
}

func (xuen *Xuen) GetPet() *core.Pet {
	return &xuen.Pet
}
