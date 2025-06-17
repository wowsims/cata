package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) applyOmenOfClarity() {
	var affectedSpells []*DruidSpell
	druid.ClearcastingAura = druid.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 16870},
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice([]*DruidSpell{
				// Balance
				druid.Starfire,
				druid.Wrath,
				druid.Hurricane,
				druid.WildMushrooms,
				druid.ForceOfNature,
				druid.Starsurge,

				// Feral
				druid.DemoralizingRoar,
				druid.FerociousBite,
				druid.Lacerate,
				druid.MangleBear,
				druid.MangleCat,
				druid.Maul,
				druid.Pulverize,
				druid.Rake,
				druid.Rip,
				druid.Shred,
				druid.SwipeBear,
				druid.SwipeCat,
				druid.Thrash,
			}, func(spell *DruidSpell) bool { return spell != nil })
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.PercentModifier *= -1
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range affectedSpells {
				spell.Cost.PercentModifier /= -1
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if aura.RemainingDuration(sim) == aura.Duration {
				// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate
				// if it was just activated.
				return
			}

			for _, as := range affectedSpells {
				// Mangle (Bear) handled separately in mangle.go in order to preferentially consume Berserk procs over Clearcasting procs
				if as.IsEqual(spell) && (as != druid.MangleBear) && (as != druid.WildMushrooms) {
					aura.Deactivate(sim)
					break
				}
			}
		},
	})

	druid.ProcOoc = func(sim *core.Simulation) {
		druid.ClearcastingAura.Activate(sim)
	}

	druid.RegisterAura(core.Aura{
		Label:    "Omen of Clarity",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
			if druid.HurricaneTickSpell.IsEqual(spell) {
				curCastTickSpeed := spell.CurDot().TickPeriod().Seconds() / 10
				hurricaneCoeff := 1.0 - (7.0 / 9.0)
				spellCoeff := hurricaneCoeff * curCastTickSpeed
				chanceToProc := ((1.5 / 60) * 3.5) * spellCoeff
				if sim.Proc(chanceToProc, "Clearcasting") {
					druid.ProcOoc(sim)
				}
			} else if druid.AutoAttacks.PPMProc(sim, 3.5, core.ProcMaskMeleeWhiteHit, "Omen of Clarity", spell) { // Melee
				druid.ProcOoc(sim)
			} else if spell.Flags.Matches(SpellFlagOmenTrigger) { // Spells
				// Heavily based on comment here
				// https://github.com/JamminL/wotlk-classic-bugs/issues/66#issuecomment-1182017571
				// Instants are treated as 1.5
				// Uses current cast time rather than default cast time (PPM is constant with haste)
				castTime := spell.CurCast.CastTime.Seconds()
				if castTime == 0 {
					castTime = 1.5
				}

				chanceToProc := (castTime / 60) * 3.5
				if druid.Typhoon.IsEqual(spell) { // Add Typhoon
					chanceToProc *= 0.25
				} else if druid.Moonfire.IsEqual(spell) { // Add Moonfire
					chanceToProc *= 0.076
				} else if druid.WildMushroomsDetonate.IsEqual(spell) {
					// Wild Mushroom: Detonate seems to have an 'almost' guaranteed chance to proc
					// setting to 0.5 to be safe
					chanceToProc = 0.5
				} else {
					chanceToProc *= 0.666
				}
				if sim.Proc(chanceToProc, "Clearcasting") {
					druid.ProcOoc(sim)
				}
			}
		},
	})
}
