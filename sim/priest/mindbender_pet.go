package priest

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

type MindBender struct {
	core.Pet

	Priest          *Priest
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

func (priest *Priest) NewMindBender() *MindBender {
	mindbender := &MindBender{
		Pet: core.NewPet(core.PetConfig{
			Name:                            "Mindbender",
			Owner:                           &priest.Character,
			BaseStats:                       baseStats,
			NonHitExpStatInheritance:        priest.mindbenderStatInheritance(),
			IsGuardian:                      false,
			EnabledOnStart:                  false,
			HasDynamicMeleeSpeedInheritance: true,
		}),
		Priest: priest,
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	core.MakePermanent(mindbender.GetOrRegisterAura(core.Aura{
		Label: "Autoattack mana regen",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			restoreMana := priest.MaxMana() * 0.0175
			priest.AddMana(sim, restoreMana, manaMetric)
		},
	}))

	actionID := core.ActionID{SpellID: 63619}
	mindbender.ShadowcrawlAura = mindbender.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
	}).AttachMultiplicativePseudoStatBuff(&mindbender.PseudoStats.DamageDealtMultiplier, 1.15)

	mindbender.Shadowcrawl = mindbender.RegisterSpell(core.SpellConfig{
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
			mindbender.ShadowcrawlAura.Activate(sim)
		},
	})

	mindbender.PseudoStats.DamageTakenMultiplier *= 0.1

	// never misses
	mindbender.AddStats(stats.Stats{
		stats.HitRating:       8 * core.PhysicalHitRatingPerHitPercent,
		stats.ExpertiseRating: 14 * core.ExpertisePerQuarterPercentReduction * 4,
	})

	mindbender.EnableAutoAttacks(mindbender, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        priest.CalcScalingSpellDmg(1.76),
			BaseDamageMax:        priest.CalcScalingSpellDmg(1.76),
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
			AttackPowerPerDPS:    core.DefaultAttackPowerPerDPS,
		},
		AutoSwingMelee: true,
	})

	mindbender.AutoAttacks.MHConfig().BonusCoefficient = 1

	mindbender.EnableManaBar()
	priest.AddPet(mindbender)

	return mindbender
}

func (priest *Priest) mindbenderStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{
			stats.PhysicalCritPercent: ownerStats[stats.SpellCritPercent],
			stats.Intellect:           (ownerStats[stats.Intellect] - 10) * 0.3,
			stats.Stamina:             ownerStats[stats.Stamina] * 0.75,
			stats.SpellPower:          0.88 * ownerStats[stats.SpellPower],
			stats.HasteRating:         ownerStats[stats.HasteRating],
		}
	}
}

func (mindbender *MindBender) Initialize() {
}

func (mindbender *MindBender) ExecuteCustomRotation(sim *core.Simulation) {
	mindbender.Shadowcrawl.Cast(sim, nil)
}

func (mindbender *MindBender) Reset(sim *core.Simulation) {
	mindbender.ShadowcrawlAura.Deactivate(sim)
	mindbender.Disable(sim)
}

func (mindbender *MindBender) OnPetDisable(sim *core.Simulation) {
	mindbender.ShadowcrawlAura.Deactivate(sim)
}

func (mindbender *MindBender) GetPet() *core.Pet {
	return &mindbender.Pet
}
