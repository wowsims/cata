package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (priest *Priest) newMindFlaySpell(numTicksIdx int32) *PriestSpell {
	numTicks := numTicksIdx
	flags := core.SpellFlagChanneled
	if numTicksIdx == 0 {
		numTicks = 3
		flags |= core.SpellFlagAPL
	}

	return priest.RegisterSpell(PriestSpellMindFlay, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 15407}.WithTag(numTicksIdx),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "MindFlay-" + strconv.Itoa(int(numTicksIdx)),
			},
			NumberOfTicks:       numTicks,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 187.147 + 0.288*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				spell.DealOutcome(sim, result)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := 187.147 + 0.288*spell.SpellPower()
			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
		},
	})
}
