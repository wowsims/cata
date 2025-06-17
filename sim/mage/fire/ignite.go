package fire

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/mage"
)

func (fire *FireMage) registerMastery() {

	fire.Ignite = RegisterIgniteEffect(&fire.Unit, IgniteConfig{
		ActionID:       core.ActionID{SpellID: 12846},
		ClassSpellMask: mage.MageSpellIgnite,
		DotAuraLabel:   "Ignite",
		DotAuraTag:     "IgniteDot",

		ProcTrigger: core.ProcTrigger{
			Name:     "Ignite Talent",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellDamage,
			Outcome:  core.OutcomeLanded,

			ExtraCondition: func(_ *core.Simulation, spell *core.Spell, _ *core.SpellResult) bool {
				return spell.Matches(FireSpellIgnitable)
			},
		},

		DamageCalculator: func(result *core.SpellResult) float64 {
			return result.Damage * fire.GetMasteryBonus()
		},
	})

	// This is needed because we want to listen for the spell "cast" event that refreshes the Dot
	fire.Ignite.Flags ^= core.SpellFlagNoOnCastComplete

}

func (fire *FireMage) GetMasteryBonus() float64 {
	return (.12 + 0.015*fire.GetMasteryPoints())
}

// Takes in the SpellResult for the triggering spell, and returns the damage per
// tick of a *fresh* Ignite triggered by that spell. Roll-over damage
// calculations for existing Ignites are handled internally.
type IgniteDamageCalculator func(result *core.SpellResult) float64

type IgniteConfig struct {
	ActionID           core.ActionID
	ClassSpellMask     int64
	DisableCastMetrics bool
	DotAuraLabel       string
	DotAuraTag         string
	ProcTrigger        core.ProcTrigger // Ignores the Handler field and creates a custom one, but uses all others.
	DamageCalculator   IgniteDamageCalculator
	IncludeAuraDelay   bool // "munching" and "free roll-over" interactions
	SpellSchool        core.SpellSchool
	NumberOfTicks      int32
	TickLength         time.Duration
	SetBonusAura       *core.Aura
}

func RegisterIgniteEffect(unit *core.Unit, config IgniteConfig) *core.Spell {
	spellFlags := core.SpellFlagIgnoreModifiers | core.SpellFlagNoSpellMods | core.SpellFlagNoOnCastComplete

	if config.DisableCastMetrics {
		spellFlags |= core.SpellFlagPassiveSpell
	}

	if config.SpellSchool == 0 {
		config.SpellSchool = core.SpellSchoolFire
	}

	if config.NumberOfTicks == 0 {
		config.NumberOfTicks = 2
	}

	if config.TickLength == 0 {
		config.TickLength = time.Second * 2
	}

	igniteSpell := unit.RegisterSpell(core.SpellConfig{
		ActionID:         config.ActionID,
		SpellSchool:      config.SpellSchool,
		ProcMask:         core.ProcMaskSpellProc,
		ClassSpellMask:   config.ClassSpellMask,
		Flags:            spellFlags,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     config.DotAuraLabel,
				Tag:       config.DotAuraTag,
				MaxStacks: math.MaxInt32,
			},

			NumberOfTicks:       config.NumberOfTicks,
			TickLength:          config.TickLength,
			AffectedByCastSpeed: false,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.Spell.CalcPeriodicDamage(sim, target, dot.SnapshotBaseDamage, dot.OutcomeTick)
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})

	refreshIgnite := func(sim *core.Simulation, target *core.Unit, damagePerTick float64) {
		// Cata Ignite
		// 1st ignite application = 4s, split into 2 ticks (2s, 0s)
		// Ignite refreshes: Duration = 4s + MODULO(remaining duration, 2), max 6s. Split damage over 3 ticks at 4s, 2s, 0s.
		dot := igniteSpell.Dot(target)
		dot.SnapshotBaseDamage = damagePerTick
		igniteSpell.Cast(sim, target)
		dot.Aura.SetStacks(sim, int32(dot.SnapshotBaseDamage))
	}

	var scheduledRefresh *core.PendingAction
	procTrigger := config.ProcTrigger
	procTrigger.Handler = func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
		target := result.Target
		dot := igniteSpell.Dot(target)
		outstandingDamage := dot.OutstandingDmg()
		newDamage := config.DamageCalculator(result)
		totalDamage := outstandingDamage + newDamage
		newTickCount := dot.BaseTickCount + core.TernaryInt32(dot.IsActive(), 1, 0)
		damagePerTick := totalDamage / float64(newTickCount)

		if config.IncludeAuraDelay {
			// Rough 2-bucket model for the aura update delay distribution based
			// on PTR measurements. Most updates occur on either the same or very
			// next spell batch after the proc, and can therefore be modeled by a
			// 0-10 ms random draw. But a reasonable minority fraction take ~10x
			// longer than this to fire. The origin of these longer delays is
			// likely not actually random in reality, but can be treated that way
			// in practice since the player cannot play around them.
			var delaySeconds float64

			if sim.Proc(0.75, "Aura Delay") {
				delaySeconds = 0.010 * sim.RandomFloat("Aura Delay")
			} else {
				delaySeconds = 0.090 + 0.020*sim.RandomFloat("Aura Delay")
			}

			applyDotAt := sim.CurrentTime + core.DurationFromSeconds(delaySeconds)

			// Cancel any prior aura updates already in the queue
			if (scheduledRefresh != nil) && (scheduledRefresh.NextActionAt > sim.CurrentTime) {
				scheduledRefresh.Cancel(sim)

				if sim.Log != nil {
					unit.Log(sim, "Previous %s proc was munched due to server aura delay", config.DotAuraLabel)
				}
			}

			// Schedule a delayed refresh of the DoT with cached damagePerTick value (allowing for "free roll-overs")
			if sim.Log != nil {
				unit.Log(sim, "Schedule travel (%0.1f ms) for %s", delaySeconds*1000, config.DotAuraLabel)

				if dot.IsActive() && (dot.NextTickAt() < applyDotAt) {
					unit.Log(sim, "%s rolled with %0.3f damage both ticking and rolled into next", config.DotAuraLabel, outstandingDamage)
				}
			}

			scheduledRefresh = core.StartDelayedAction(sim, core.DelayedActionOptions{
				DoAt:     applyDotAt,
				Priority: core.ActionPriorityDOT,

				OnAction: func(_ *core.Simulation) {
					refreshIgnite(sim, target, damagePerTick)
				},
			})
		} else {
			refreshIgnite(sim, target, damagePerTick)
		}
	}

	if config.SetBonusAura != nil {
		config.SetBonusAura.AttachProcTrigger(procTrigger)
	} else {
		core.MakeProcTriggerAura(unit, procTrigger)
	}

	return igniteSpell
}
