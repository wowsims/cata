package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

// DISCLAIMER: Shadowfiend need some extensive research on Level 85
// Proper Spell Scaling? Wiki says 37.5%, patch notes state 30%
// WoW Sims implemented priest crit scaling but we do not
// Right now Stats are inherited statically on spawn, but testing
// indicates shadow fiend scales per hit based on owner spell power
type Shadowfiend struct {
	core.Pet

	Priest          *Priest
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

var baseStats = stats.Stats{
	stats.Strength:  0,
	stats.Agility:   0,
	stats.Stamina:   348,
	stats.Intellect: 350,
	// stats.AttackPower: 896, // Level 85
	stats.AttackPower: 1077, // Level 90
	stats.Mana:        12295,
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	shadowfiend := &Shadowfiend{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Shadowfiend",
			Owner:                           &priest.Character,
			BaseStats:                       baseStats,
			StatInheritance:                 priest.shadowfiendStatInheritance(),
			EnabledOnStart:                  false,
			IsGuardian:                      false,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		Priest: priest,
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	core.MakePermanent(shadowfiend.GetOrRegisterAura(core.Aura{
		Label: "Autoattack mana regen",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			restoreMana := priest.MaxMana() * 0.03
			priest.AddMana(sim, restoreMana, manaMetric)
		},
	}))

	actionID := core.ActionID{SpellID: 63619}
	shadowfiend.ShadowcrawlAura = shadowfiend.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
	}).AttachMultiplicativePseudoStatBuff(&shadowfiend.PseudoStats.DamageDealtMultiplier, 1.15)

	shadowfiend.Shadowcrawl = shadowfiend.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowfiend.ShadowcrawlAura.Activate(sim)
		},
	})

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	// never misses
	shadowfiend.AddStats(stats.Stats{
		stats.HitRating:       8 * core.PhysicalHitRatingPerHitPercent,
		stats.ExpertiseRating: 14 * core.ExpertisePerQuarterPercentReduction * 4,
	})

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        priest.CalcScalingSpellDmg(2.0),
			BaseDamageMax:        priest.CalcScalingSpellDmg(2.0),
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
			AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.AutoAttacks.MHConfig().BonusCoefficient = 1
	shadowfiend.EnableManaBar()
	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.PhysicalCritPercent: ownerStats[stats.SpellCritPercent],
			stats.Intellect:           (ownerStats[stats.Intellect] - 10) * 0.3,
			stats.Stamina:             ownerStats[stats.Stamina] * 0.75,
			stats.SpellPower:          ownerStats[stats.SpellPower],
			stats.HasteRating:         ownerStats[stats.HasteRating],
		}
	}
}

func (shadowfiend *Shadowfiend) Initialize() {
}

func (shadowfiend *Shadowfiend) ExecuteCustomRotation(sim *core.Simulation) {
	shadowfiend.Shadowcrawl.Cast(sim, nil)
}

func (shadowfiend *Shadowfiend) Reset(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
	shadowfiend.Disable(sim)
}

func (shadowfiend *Shadowfiend) OnPetDisable(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
}

func (shadowfiend *Shadowfiend) GetPet() *core.Pet {
	return &shadowfiend.Pet
}
