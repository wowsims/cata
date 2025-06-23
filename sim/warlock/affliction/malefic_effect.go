package affliction

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/warlock"
)

func (affliction *AfflictionWarlock) registerMaleficEffect() {
	buildSpell := func(id int32) *core.Spell {
		return affliction.RegisterSpell(core.SpellConfig{
			ActionID:       core.ActionID{SpellID: id}.WithTag(1),
			Flags:          core.SpellFlagNoOnCastComplete | core.SpellFlagNoSpellMods | core.SpellFlagIgnoreAttackerModifiers,
			SpellSchool:    core.SpellSchoolShadow,
			ProcMask:       core.ProcMaskSpellDamage,
			ClassSpellMask: warlock.WarlockSpellMaleficGrasp,

			ThreatMultiplier: 1,
			DamageMultiplier: 1,
			CritMultiplier:   affliction.DefaultCritMultiplier(),
			BonusSpellPower:  0, // used to transmit base damage
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, spell.BonusSpellPower, spell.OutcomeMagicHit)
			},
		})
	}

	corruptionProc := buildSpell(172)
	agonyProc := buildSpell(980)
	uaProc := buildSpell(30108)

	procTable := map[*core.Spell]**core.Spell{
		corruptionProc: &affliction.Corruption,
		agonyProc:      &affliction.Agony,
		uaProc:         &affliction.UnstableAffliction,
	}

	// used to iterate over the map in constant order
	procKeys := []*core.Spell{corruptionProc, agonyProc, uaProc}
	affliction.ProcMaleficEffect = func(target *core.Unit, coeff float64, sim *core.Simulation) {

		// I don't like it but if sac is specced the damage replciation effect specifically is increased by 20%
		// Nothing we can do really properly with SpellMod's here nicely
		if affliction.Talents.GrimoireOfSacrifice {
			coeff *= 1.2
		}

		if affliction.T15_4pc.IsActive() {
			coeff *= 1.05
		}

		if affliction.T16_2pc_buff != nil && affliction.T16_2pc_buff.IsActive() {
			coeff *= 1.2
		}

		for _, proc := range procKeys {
			source := procTable[proc]
			dot := (*source).Dot(target)
			if !dot.IsActive() {
				continue
			}

			proc.BonusSpellPower = calculateDoTBaseTickDamage(dot) * coeff
			proc.Cast(sim, target)
		}
	}
}
