package warlock

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (warlock *Warlock) ApplyAfflictionTalents() {

	if warlock.Talents.ImprovedCorruption > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellCorruption,
			Kind:       core.SpellMod_DamageDone_Flat,
			FloatValue: 0.04 * float64(warlock.Talents.ImprovedCorruption),
		})
	}

	if warlock.Talents.DoomAndGloom > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellBaneOfAgony | WarlockSpellBaneOfDoom,
			Kind:       core.SpellMod_BonusCrit_Rating,
			FloatValue: 4.0 * float64(warlock.Talents.DoomAndGloom) * core.CritRatingPerCritChance,
		})
	}

	warlock.registerEradication()
	warlock.registerShadowEmbrace()
	warlock.registerDeathsEmbrace()
	warlock.registerNightfall()

	if warlock.Talents.EverlastingAffliction > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask:  WarlockSpellCorruption | WarlockSpellSeedOfCorruption | WarlockSpellSeedOfCorruptionExposion | WarlockSpellUnstableAffliction,
			Kind:       core.SpellMod_BonusCrit_Rating,
			FloatValue: 5.0 * float64(warlock.Talents.EverlastingAffliction) * core.CritRatingPerCritChance,
		})

		warlock.registerEverlastingAffliction()
	}

	if warlock.Talents.Pandemic > 0 {
		warlock.AddStaticMod(core.SpellModConfig{
			ClassMask: WarlockSpellBaneOfAgony | WarlockSpellBaneOfDoom | WarlockSpellCurseOfElements | WarlockSpellCurseOfWeakness | WarlockSpellCurseOfTongues,
			Kind:      core.SpellMod_GlobalCooldown_Flat,
			TimeValue: time.Duration(-250*warlock.Talents.Pandemic) * time.Millisecond,
		})

		warlock.registerPandemic()
	}
}

func (warlock *Warlock) registerEradication() {
	if warlock.Talents.Eradication <= 0 {
		return
	}

	castSpeedMultiplier := []float64{1, 1.06, 1.12, 1.20}[warlock.Talents.Eradication]
	eradicationAura := warlock.RegisterAura(core.Aura{
		Label:    "Eradication",
		ActionID: core.ActionID{SpellID: 47197},
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(castSpeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / castSpeedMultiplier)
		},
	})

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label: "Eradication Talent",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask&(WarlockSpellCorruption) > 0 && sim.Proc(0.06, "Eradication") {
					eradicationAura.Activate(sim)
				}
			},
		}))
}

func (warlock *Warlock) registerDeathsEmbrace() {
	if warlock.Talents.DeathsEmbrace <= 0 {
		return
	}

	deathsEmbraceMod := warlock.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  WarlockShadowDamage,
		FloatValue: 0.04 * float64(warlock.Talents.DeathsEmbrace),
	})

	deathsEmbraceAura := warlock.RegisterAura(core.Aura{
		Label:    "Deaths Embrace",
		ActionID: core.ActionID{SpellID: 1120},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			deathsEmbraceMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			deathsEmbraceMod.Deactivate()
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				deathsEmbraceAura.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) ShadowEmbraceDebuffAura(target *core.Unit) *core.Aura {
	shadowEmbraceBonus := []float64{0, 0.03, 0.04, 0.05}[warlock.Talents.ShadowEmbrace]

	return target.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Embrace-" + warlock.Label,
		ActionID:  core.ActionID{SpellID: 32392},
		Duration:  12 * time.Second,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier /= 1.0 + shadowEmbraceBonus*float64(oldStacks)
			warlock.AttackTables[aura.Unit.UnitIndex].HauntSEDamageTakenMultiplier *= 1.0 + shadowEmbraceBonus*float64(newStacks)
		},
	})
}

func (warlock *Warlock) registerShadowEmbrace() {
	if warlock.Talents.ShadowEmbrace <= 0 {
		return
	}

	warlock.ShadowEmbraceAuras = warlock.NewEnemyAuraArray(warlock.ShadowEmbraceDebuffAura)

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label: "Shadow Embrace Talent Hidden Aura",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if (spell == warlock.ShadowBolt || spell == warlock.Haunt) && result.Landed() {
					aura := warlock.ShadowEmbraceAuras.Get(result.Target)
					aura.Activate(sim)
					aura.AddStack(sim)
				}
			},
		}))
}

func (warlock *Warlock) registerEverlastingAffliction() {
	procChance := []float64{0, 0.33, 0.66, 1.0}[warlock.Talents.EverlastingAffliction]

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label: "EverlastingAffliction Talent",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask != WarlockSpellDrainSoul {
					return
				}

				if warlock.Corruption.Dot(result.Target).IsActive() && sim.Proc(procChance, "EverlastingAffliction") {
					warlock.Corruption.Dot(result.Target).Apply(sim)
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ClassSpellMask == WarlockSpellHaunt && result.Landed() {
					if warlock.Corruption.Dot(result.Target).IsActive() && sim.Proc(procChance, "EverlastingAffliction") {
						warlock.Corruption.Dot(result.Target).Apply(sim)
					}
				}
			},
		}))
}

func (warlock *Warlock) registerPandemic() {
	if warlock.Talents.Pandemic <= 0 {
		return
	}

	procChance := []float64{0, 0.5, 1.0}[warlock.Talents.Pandemic]

	pandemicAura := warlock.RegisterAura(core.Aura{
		Label:    "Pandemic Talent",
		Duration: core.NeverExpires,
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != warlock.DrainSoul {
				return
			}

			if warlock.UnstableAffliction.Dot(result.Target).IsActive() && sim.Proc(procChance, "Pandemic") {
				warlock.UnstableAffliction.Dot(result.Target).Apply(sim)
			}
		},
	})

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			if isExecute == 25 {
				pandemicAura.Activate(sim)
			}
		})
	})
}

func (warlock *Warlock) registerNightfall() {
	if warlock.Talents.Nightfall <= 0 && !warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfCorruption) {
		return
	}

	nightfallProcChance := 0.02*float64(warlock.Talents.Nightfall) +
		0.04*core.TernaryFloat64(warlock.HasPrimeGlyph(proto.WarlockPrimeGlyph_GlyphOfCorruption), 1, 0)

	icd := core.Cooldown{
		Timer:    warlock.NewTimer(),
		Duration: 6 * time.Second,
	}

	nightfallProcAura := warlock.RegisterAura(core.Aura{
		Icd:      &icd,
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: 10 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if spell.ClassSpellMask&(WarlockSpellShadowBolt) > 0 && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	core.MakePermanent(
		warlock.RegisterAura(core.Aura{
			Label: "Nightfall Hidden Aura",
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == warlock.Corruption {
					if sim.Proc(nightfallProcChance, "Nightfall") {
						nightfallProcAura.Activate(sim)
					}
				}
			},
		}))
}
