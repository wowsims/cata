package protection

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/paladin"
)

// Consecrates the land beneath you, causing 8222 Holy damage over 9 sec to enemies who enter the area.
func (prot *ProtectionPaladin) registerConsecrationSpell() {
	numTargets := prot.Env.GetNumTargets()
	prot.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 26573},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL | core.SpellFlagAoE,
		ClassSpellMask: paladin.SpellMaskConsecration,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 7,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    prot.NewTimer(),
				Duration: 9 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   prot.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 26573},
				Label:    "Consecration" + prot.Label,
			},
			NumberOfTicks: 9,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				results := make([]*core.SpellResult, numTargets)

				// Consecration recalculates everything on each tick
				baseDamage := prot.CalcScalingSpellDmg(0.80000001192) + 0.07999999821*dot.Spell.MeleeAttackPower()

				for idx := range numTargets {
					currentTarget := sim.Environment.GetTargetUnit(idx)
					results[idx] = dot.Spell.CalcPeriodicDamage(sim, currentTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}

				for idx := range numTargets {
					dot.Spell.DealPeriodicDamage(sim, results[idx])
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
