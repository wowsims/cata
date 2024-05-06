package druid

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *Druid) registerMaulSpell() {
	flatBaseDamage := 578.0
	if druid.Ranged().ID == 23198 { // Idol of Brutality
		flatBaseDamage += 50
	} else if druid.Ranged().ID == 38365 { // Idol of Perspicacious Attacks
		flatBaseDamage += 120
	}

	numHits := core.TernaryInt32(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMaul) && druid.Env.GetNumTargets() > 1, 2, 1)
	rendAndTearMod := []float64{1.0, 1.07, 1.13, 1.2}[druid.Talents.RendAndTear]

	druid.Maul = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48480},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},

		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,
		FlatThreatBonus:  424,
		BonusCoefficient: 1,
		MaxRange:         core.MaxMeleeRange,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Need to specially deactivate CC here in case maul is cast simultaneously with another spell.
			if druid.ClearcastingAura != nil {
				druid.ClearcastingAura.Deactivate(sim)
			}

			baseDamage := flatBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				modifier := 1.0
				if druid.BleedCategories.Get(curTarget).AnyActive() {
					modifier += .3
				}
				if druid.AssumeBleedActive || druid.Rip.Dot(curTarget).IsActive() || druid.Rake.Dot(curTarget).IsActive() || druid.Lacerate.Dot(curTarget).IsActive() {
					modifier *= rendAndTearMod
				}
				if hitIndex > 0 {
					modifier *= 0.5
				}

				result := spell.CalcAndDealDamage(sim, curTarget, baseDamage * modifier, spell.OutcomeMeleeSpecialHitAndCrit)

				if !result.Landed() {
					spell.IssueRefund(sim)
				}

				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}

			druid.MaulQueueAura.Deactivate(sim)
		},
	})

	druid.MaulQueueAura = druid.RegisterAura(core.Aura{
		Label:    "Maul Queue Aura",
		ActionID: druid.Maul.ActionID,
		Duration: core.NeverExpires,
	})

	druid.MaulQueueSpell = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    druid.Maul.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.MaulQueueAura.IsActive() &&
				druid.CurrentRage() >= druid.Maul.DefaultCast.Cost &&
				sim.CurrentTime >= druid.Hardcast.Expires
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			druid.MaulQueueAura.Activate(sim)
		},
	})
}

func (druid *Druid) QueueMaul(sim *core.Simulation) {
	if druid.MaulQueueSpell.CanCast(sim, druid.CurrentTarget) {
		druid.MaulQueueSpell.Cast(sim, druid.CurrentTarget)
	}
}

// Returns true if the regular melee swing should be used, false otherwise.
func (druid *Druid) MaulReplaceMH(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	if !druid.MaulQueueAura.IsActive() {
		return mhSwingSpell
	}

	if !druid.Maul.Spell.CanCast(sim, druid.CurrentTarget) {
		druid.MaulQueueAura.Deactivate(sim)
		return mhSwingSpell
	}

	return druid.Maul.Spell
}
