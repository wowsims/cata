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
			ClassMask:  WarlockSpellCorruption | WarlockSpellSeedOfCorruption | WarlockSpellSeedOfCorruptionExposion | WarlockSpellUnstableAffliction,
			Kind:       core.SpellMod_GlobalCooldown_Flat,
			FloatValue: -1 * float64(warlock.Talents.EverlastingAffliction) * 0.25,
		})

		warlock.registerPandemic()
	}
}

func (warlock *Warlock) registerEradication() {
	if warlock.Talents.Eradication <= 0 {
		return
	}

	castSpeedMultiplier := []float64{1, 1.06, 1.12, 1.20}[warlock.Talents.Eradication]
	warlock.EradicationAura = warlock.RegisterAura(core.Aura{
		Label:    "Eradication",
		ActionID: core.ActionID{SpellID: 47197},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(castSpeedMultiplier)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.MultiplyCastSpeed(1 / castSpeedMultiplier)
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Eradication Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != warlock.Corruption {
				return
			}

			if sim.Proc(0.06, "Eradication") {
				warlock.EradicationAura.Activate(sim)
			}
		},
	})
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

	warlock.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
			//TODO: Does this need to deactivate somewhere?
			if isExecute == 25 {
				deathsEmbraceMod.Activate()
			}
		})
	})
}

func (warlock *Warlock) ShadowEmbraceDebuffAura(target *core.Unit) *core.Aura {
	shadowEmbraceBonus := []float64{0, 0.03, 0.04, 0.05}[warlock.Talents.ShadowEmbrace]

	return target.GetOrRegisterAura(core.Aura{
		Label:     "Shadow Embrace-" + warlock.Label,
		ActionID:  core.ActionID{SpellID: 32392},
		Duration:  time.Second * 12,
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

	warlock.RegisterAura(core.Aura{
		Label:    "Shadow Embrace Talent Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if (spell == warlock.ShadowBolt || spell == warlock.Haunt) && result.Landed() {
				aura := warlock.ShadowEmbraceAuras.Get(result.Target)
				aura.Activate(sim)
				aura.AddStack(sim)
			}
		},
	})
}

func (warlock *Warlock) registerEverlastingAffliction() {
	procChance := []float64{0, 0.33, 0.66, 1.0}[warlock.Talents.EverlastingAffliction]

	warlock.RegisterAura(core.Aura{
		Label:    "EverlastingAffliction Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			//TODO: Does this apply on haunt damage or just the dot?
			if spell != warlock.DrainSoul || spell != warlock.Haunt {
				return
			}

			if warlock.Corruption.Dot(aura.Unit).IsActive() && sim.Proc(procChance, "EverlastingAffliction") {
				//TODO: Should this Rollover or Apply like other dots in cata?
				warlock.Corruption.Dot(aura.Unit).Rollover(sim)
			}
		},
	})
}

func (warlock *Warlock) registerPandemic() {
	if warlock.Talents.Pandemic <= 0 {
		return
	}

	procChance := []float64{0, 0.5, 1.0}[warlock.Talents.Pandemic]

	pandemicAura := warlock.RegisterAura(core.Aura{
		Label:    "Pandemic Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			//TODO: Does this need to deactivate here?
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != warlock.DrainSoul {
				return
			}

			if warlock.UnstableAffliction.Dot(aura.Unit).IsActive() && sim.Proc(procChance, "Pandemic") {
				//TODO: Should this Rollover or Apply like other dots in cata?
				warlock.UnstableAffliction.Dot(aura.Unit).Rollover(sim)
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
		Duration: time.Second * 6,
	}

	warlock.NightfallProcAura = warlock.RegisterAura(core.Aura{
		Icd:      &icd,
		Label:    "Nightfall Shadow Trance",
		ActionID: core.ActionID{SpellID: 17941},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier -= 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.ShadowBolt.CastTimeMultiplier += 1
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// Check if the shadowbolt was instant cast and not a normal one
			if spell == warlock.ShadowBolt && spell.CurCast.CastTime == 0 {
				aura.Deactivate(sim)
			}
		},
	})

	warlock.RegisterAura(core.Aura{
		Label:    "Nightfall Hidden Aura",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == warlock.Corruption {
				if sim.Proc(nightfallProcChance, "Nightfall") {
					warlock.NightfallProcAura.Activate(sim)
				}
			}
		},
	})
}
