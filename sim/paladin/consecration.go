package paladin

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (paladin *Paladin) registerConsecrationSpell() {
	numTicks := int32(10)
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration) {
		numTicks += 2
	}

	consAvgDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 0.07900000364)

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 26573},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskConsecration,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 55,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: 30 * time.Second,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 26573},
				Label:    "Consecration" + paladin.Label,
			},
			NumberOfTicks: numTicks,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Consecration recalculates everything on each tick
				baseDamage := consAvgDamage +
					0.0270000007*dot.Spell.MeleeAttackPower() +
					0.0270000007*dot.Spell.SpellPower()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
