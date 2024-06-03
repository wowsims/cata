package paladin

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (paladin *Paladin) registerConsecrationSpell() {
	numTicks := int32(10)
	if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfConsecration) {
		numTicks += 2
	}

	consAvgDamage := core.CalcScalingSpellAverageEffect(proto.Class_ClassPaladin, 0.079)

	paladin.Consecration = paladin.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 26573},
		SpellSchool:    core.SpellSchoolHoly,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: SpellMaskConsecration,

		MaxRange: 8,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.55,
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
		CritMultiplier:   paladin.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				ActionID: core.ActionID{SpellID: 81297},
				Label:    "Consecration",
			},
			NumberOfTicks: numTicks,
			TickLength:    time.Second * 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				// Consecration recalculates everything on each tick
				baseDamage := consAvgDamage +
					.027*dot.Spell.MeleeAttackPower() +
					.027*dot.Spell.SpellPower()
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
