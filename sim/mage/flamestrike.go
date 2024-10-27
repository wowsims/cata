package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) GetFlameStrikeConfig(spellId int32, isProc bool) core.SpellConfig {
	label := "Flamestrike - " + strconv.Itoa(int(spellId))
	if isProc {
		label += " - Proc"
	}

	config := core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellId},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlag(core.TernaryInt32(isProc, 0, int32(core.SpellFlagAPL))),
		ClassSpellMask: MageSpellFlamestrike,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
		},

		DamageMultiplierAdditive: 1,
		CritMultiplier:           mage.DefaultSpellCritMultiplier(),
		BonusCoefficient:         0.146,
		ThreatMultiplier:         1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: label,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := mage.CurrentTarget
				baseDamage := 0.103 * mage.ClassSpellScaling
				dot.Snapshot(target, baseDamage)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeSnapshotCrit)
				}
			},
			BonusCoefficient: 0.061,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := 0.662 * mage.ClassSpellScaling
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}

			if spell.AOEDot() != nil {
				spell.AOEDot().Apply(sim)
			}
		},
	}

	if !isProc {
		config.Cast = core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		}

		config.ManaCost = core.ManaCostOptions{
			BaseCost: 0.30,
		}
	} else {
		config.ProcMask = core.ProcMaskSpellProc
		config.ActionID = config.ActionID.WithTag(1)
	}

	return config
}
func (mage *Mage) registerFlamestrikeSpell() {
	mage.Flamestrike = mage.RegisterSpell(mage.GetFlameStrikeConfig(2120, false))
}
