package demonology

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func (demonology *DemonologyWarlock) registerMetamorphosis() {
	if !demonology.Talents.Metamorphosis {
		return
	}

	var immolationAura *core.Spell
	metaDmgMod := 0.0
	glyphBonus := core.TernaryDuration(demonology.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfMetamorphosis), 6, 0)

	metamorphosisAura := demonology.RegisterAura(core.Aura{
		Label:    "Metamorphosis Aura",
		ActionID: core.ActionID{SpellID: 59672},
		Duration: (30 + glyphBonus) * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			metaDmgMod = 1.2 + 0.184 + 0.023*demonology.GetMasteryPoints()
			if sim.CurrentTime <= 0 && demonology.prepullMastery > 0 {
				masteryPoints := core.MasteryRatingToMasteryPoints(float64(demonology.prepullMastery))
				metaDmgMod = 1.2 + 0.184 + 0.023*masteryPoints

				if demonology.prepullPostSnapshotMana > 0 {
					prepullPostSnapshotManaMetric := demonology.NewManaMetrics(core.ActionID{OtherID: proto.OtherAction_OtherActionPrepull})
					maxMana := demonology.MaxMana()
					currentMana := demonology.CurrentMana()
					postSnapshotMana := currentMana - float64(demonology.prepullPostSnapshotMana)
					if postSnapshotMana < 0 {
						postSnapshotMana = currentMana
					} else if postSnapshotMana > maxMana {
						postSnapshotMana = maxMana
					}
					demonology.SpendMana(sim, postSnapshotMana, prepullPostSnapshotManaMetric)
				}
			}
			aura.Unit.PseudoStats.DamageDealtMultiplier *= metaDmgMod

			if sim.Log != nil {
				demonology.Log(sim, "[DEBUG]: meta damage mod: %v", metaDmgMod)
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= metaDmgMod
			immolationAura.AOEDot().Deactivate(sim)
		},
	})

	demonology.Metamorphosis = demonology.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 59672},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: 180 * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			metamorphosisAura.Activate(sim)
		},
	})

	demonology.AddMajorCooldown(core.MajorCooldown{
		Spell: demonology.Metamorphosis,
		Type:  core.CooldownTypeDPS,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			//TODO: This will probably need tuning for Cata and the new Impending Doom talent
			// Changed the execute phase to 25 and I think the demonic pact can be removed.
			if !demonology.GetAura("Demonic Pact").IsActive() {
				return false
			}
			MetamorphosisNumber := (float64(sim.Duration) + float64(metamorphosisAura.Duration)) / float64(demonology.Metamorphosis.CD.Duration)
			if MetamorphosisNumber < 1 {
				return demonology.HasActiveAura("Bloodlust-"+core.BloodlustActionID.WithTag(-1).String()) || sim.IsExecutePhase25()
			}

			return true
		},
	})

	immolationAura = demonology.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 50589},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskEmpty,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: warlock.WarlockSpellImmolationAura,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.64,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    demonology.NewTimer(),
				Duration: 30 * time.Second,
			},
		},
		ExtraCastCondition: func(_ *core.Simulation, _ *core.Unit) bool {
			return metamorphosisAura.IsActive()
		},

		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,
		BonusCoefficient:         0.10000000149,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Immolation Aura",
			},
			NumberOfTicks:       15,
			TickLength:          1 * time.Second,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDmg := demonology.CalcScalingSpellDmg(0.58899998665) * sim.Encounter.AOECapMultiplier()

				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, dot.Spell.OutcomeMagicHit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
