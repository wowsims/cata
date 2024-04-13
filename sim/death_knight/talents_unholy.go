package death_knight

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (dk *DeathKnight) ApplyUnholyTalents() {
	// Epidemic
	if dk.Talents.Epidemic > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:      core.SpellMod_DotNumberOfTicks_Flat,
			IntValue:  []int64{0, 1, 2, 4}[dk.Talents.Epidemic],
			ClassMask: DeathKnightSpellDisease,
		})
	}

	// Virulence
	if dk.Talents.Virulence > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDisease,
			FloatValue: 0.1 * float64(dk.Talents.Virulence),
		})
	}

	// Morbidity
	if dk.Talents.Morbidity > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathCoil,
			FloatValue: 0.05 * float64(dk.Talents.Morbidity),
		})

		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellDeathAndDecay,
			FloatValue: 0.1 * float64(dk.Talents.Morbidity),
		})
	}

	// Rage of Rivendare
	if dk.Talents.RageOfRivendare > 0 {
		dk.AddStaticMod(core.SpellModConfig{
			Kind:       core.SpellMod_DamageDone_Flat,
			ClassMask:  DeathKnightSpellPlagueStrike | DeathKnightSpellScourgeStrike | DeathKnightSpellFesteringStrike,
			FloatValue: 0.15 * float64(dk.Talents.RageOfRivendare),
		})
	}

	// Unholy Blight
	dk.applyUnholyBlight()

	// Runic Empowerement/Corruption
	dk.applyRunicEmpowerementCorruption()

	// Ebon Plaguebringer
	dk.applyEbonPlaguebringer()

	// Sudden Doom
	dk.applySuddenDoom()
}

func (dk *DeathKnight) applyRunicEmpowerementCorruption() {
	var handler func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult)
	if dk.Talents.RunicCorruption > 0 {
		multi := 1.0 + float64(dk.Talents.RunicCorruption)*0.5
		// Runic Corruption gives rune regen speed
		regenAura := dk.GetOrRegisterAura(core.Aura{
			Label:    "Runic Corruption",
			ActionID: core.ActionID{SpellID: 51460},
			Duration: time.Second * 3,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				dk.MultiplyRuneRegenSpeed(sim, multi)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				dk.MultiplyRuneRegenSpeed(sim, 1/multi)
			},
		})

		handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			regenAura.Activate(sim)
		}
	} else {
		// Runic Empowerement refreshes random runes on cd
		actionId := core.ActionID{SpellID: 81229}
		runeMetrics := []*core.ResourceMetrics{
			dk.NewBloodRuneMetrics(actionId),
			dk.NewFrostRuneMetrics(actionId),
			dk.NewUnholyRuneMetrics(actionId),
			dk.NewDeathRuneMetrics(actionId),
		}
		handler = func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dk.RegenRandomRune(sim, runeMetrics)
		}
	}

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Runic Empowerement",
		Callback:       core.CallbackOnSpellHitDealt,
		Outcome:        core.OutcomeLanded,
		ClassSpellMask: DeathKnightSpellDeathCoil | DeathKnightSpellRuneStrike | DeathKnightSpellFrostStrike,
		ProcChance:     0.45,
		Handler:        handler,
	})
}

func (dk *DeathKnight) applyUnholyBlight() {
	if !dk.Talents.UnholyBlight {
		return
	}

	unholyBlight := dk.Unit.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49194},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers | core.SpellFlagNoOnDamageDealt,

		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "UnholyBlight" + dk.Label,
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:           "Unholy Blight",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: DeathKnightSpellDeathCoil,
		Outcome:        core.OutcomeLanded,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			dot := unholyBlight.Dot(result.Target)

			newDamage := result.Damage * 0.10
			outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

			dot.SnapshotAttackerMultiplier = unholyBlight.DamageMultiplier
			dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(dot.NumberOfTicks)

			unholyBlight.Cast(sim, result.Target)
		},
	})
}

func (dk *DeathKnight) ebonPlaguebringerDiseaseMultiplier(spell *core.Spell, _ *core.AttackTable) float64 {
	return core.TernaryFloat64(spell.ClassSpellMask&DeathKnightSpellDisease > 0, 1.0+0.15*float64(dk.Talents.EbonPlaguebringer), 1.0)
}

func (dk *DeathKnight) applyEbonPlaguebringer() {
	if dk.Talents.EbonPlaguebringer == 0 {
		return
	}

	dk.EbonPlagueBringerAura = dk.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		aura := core.EbonPlaguebringerAura(dk.GetCharacter(), target, dk.Talents.Epidemic, dk.Talents.EbonPlaguebringer)
		aura.ApplyOnGain(func(aura *core.Aura, sim *core.Simulation) {
			dk.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = dk.ebonPlaguebringerDiseaseMultiplier
		})
		aura.ApplyOnExpire(func(aura *core.Aura, sim *core.Simulation) {
			dk.AttackTables[aura.Unit.UnitIndex].DamageDoneByCasterMultiplier = nil
		})
		return aura
	})
	dk.Env.RegisterPreFinalizeEffect(func() {
		dk.FrostFeverSpell.RelatedAuras = append(dk.FrostFeverSpell.RelatedAuras, dk.EbonPlagueBringerAura)
		dk.BloodPlagueSpell.RelatedAuras = append(dk.BloodPlagueSpell.RelatedAuras, dk.EbonPlagueBringerAura)
	})
}

func (dk *DeathKnight) applySuddenDoom() {
	if dk.Talents.SuddenDoom == 0 {
		return
	}

	mod := dk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_PowerCost_Pct,
		ClassMask:  DeathKnightSpellDeathCoil,
		FloatValue: -1,
	})

	aura := dk.GetOrRegisterAura(core.Aura{
		Label:    "Sudden Doom Proc",
		ActionID: core.ActionID{SpellID: 81340},
		Duration: time.Second * 10,

		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != DeathKnightSpellDeathCoil {
				return
			}

			if spell.CurCast.Cost > 0 {
				return
			}

			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&dk.Unit, core.ProcTrigger{
		Name:     "Sudden Doom",
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeMH,
		Outcome:  core.OutcomeLanded,
		PPM:      1.0 * float64(dk.Talents.SuddenDoom), // TODO: Find correct PPM

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			aura.Activate(sim)
		},
	})
}
