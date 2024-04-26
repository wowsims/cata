package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (druid *Druid) registerMoonfireSpell() {
	// TODO: Genesis increase to Moonfire duration
	// TODO: Shooting stars proc on periodic damage
	// TODO: Glyph of Moonfire increase to periodic damage
	// TODO: Calculate ticks based on haste
	numTicks := druid.moonfireTicks()
	//hasMoonfireGlyph := druid.HasPrimeGlyph(proto.DruidPrimeGlyph_GlyphOfMoonfire)
	//bonusPeriodicDamageMultiplier := core.TernaryFloat64(hasMoonfireGlyph, 0.2, 0)

	druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 48463},
		SpellSchool:    core.SpellSchoolArcane,
		ProcMask:       core.ProcMaskSpellDamage,
		ClassSpellMask: DruidSpellMoonfire | DruidArcaneSpells | DruidSpellDoT,
		Flags:          SpellFlagNaturesGrace | SpellFlagOmenTrigger | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.21,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,

		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Moonfire",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {

				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {

				},
			},
			NumberOfTicks: druid.moonfireTicks(),
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// dot.Spell.DamageMultiplier = baseDamageMultiplier + bonusPeriodicDamageMultiplier
				// dot.SnapshotBaseDamage = 200 + 0.13*dot.Spell.SpellPower()
				// attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				// dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				// dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable, true)
				// dot.Spell.DamageMultiplier = baseDamageMultiplier - malusInitialDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(406, 476) + 0.15*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				druid.ExtendingMoonfireStacks = 3
				dot := spell.Dot(target)
				dot.NumberOfTicks = numTicks
				dot.Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (druid *Druid) moonfireTicks() int32 {
	return 4
}
