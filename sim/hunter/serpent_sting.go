package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerSerpentStingSpell() {
	noxiousStingsMultiplier := 1 + 0.05*float64(hunter.Talents.NoxiousStings)

	impSSCritChance := float64(hunter.Talents.ImprovedSerpentSting) * 5
	impSSCritChance += core.TernaryFloat64(hunter.HasSetBonus(ItemSetLightningChargedBattleGear, 2), 5, 0)

	hunter.SerpentSting = hunter.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 1978},
		SpellSchool:  core.SpellSchoolNature,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagAPL,
		MissileSpeed: 40,
		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		BonusCritRating: impSSCritChance + core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfSerpentSting), 6, 0)*core.CritRatingPerCritChance,

		DamageMultiplierAdditive: 1,
		// SS uses Spell Crit which is multiplied by toxicology
		CritMultiplier:   hunter.SpellCritMultiplier(1, float64(hunter.Talents.Toxicology)*0.5),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SerpentStingDot",
				Tag:   "SerpentSting",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= noxiousStingsMultiplier

				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					hunter.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= noxiousStingsMultiplier
				},
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 3,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDmg := 460 + 0.08*dot.Spell.RangedAttackPower(target)
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDmg, dot.OutcomeTickPhysicalCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var result *core.SpellResult

			if hunter.Talents.ImprovedSerpentSting != 0 {
				baseDamage := (460 * 5) + 0.40*spell.RangedAttackPower(target)
				result = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			} else {
				result = spell.CalcOutcome(sim, target, spell.OutcomeRangedHitAndCrit)
			}

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					//spell.SpellMetrics[target.UnitIndex].Hits--
					spell.Dot(target).Apply(sim)

					spell.DealOutcome(sim, result)
				}
			})
		},
	})
}
