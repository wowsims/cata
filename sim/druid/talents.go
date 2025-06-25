package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/stats"
)

func (druid *Druid) ApplyTalents() {
	druid.registerYserasGift()
	druid.registerRenewal()
	druid.registerCenarionWard()

	druid.registerForceOfNature()

	druid.registerHeartOfTheWild()
	druid.registerNaturesVigil()
}

func (druid *Druid) registerHeartOfTheWild() {
	if !druid.Talents.HeartOfTheWild {
		return
	}

	// Apply 6% increase to Stamina, Agility, and Intellect
	statMultiplier := 1.06
	druid.MultiplyStat(stats.Stamina, statMultiplier)
	druid.MultiplyStat(stats.Agility, statMultiplier)
	druid.MultiplyStat(stats.Intellect, statMultiplier)

	// The activation spec specific effects are not implemented - most likely irrelevant for the sim unless proven otherwise
}

func (druid *Druid) registerNaturesVigil() {
	if !druid.Talents.NaturesVigil {
		return
	}

	actionID := core.ActionID{SpellID: 124974}

	naturesVigilAura := druid.RegisterAura(core.Aura{
		Label:    "Nature's Vigil",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(_ *core.Aura, _ *core.Simulation) {
			druid.PseudoStats.DamageDealtMultiplier *= 1.12
		},
		OnExpire: func(_ *core.Aura, _ *core.Simulation) {
			druid.PseudoStats.DamageDealtMultiplier /= 1.12
		},
	})

	naturesVigilSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:        actionID,
		Flags:           core.SpellFlagAPL,
		RelatedSelfBuff: naturesVigilAura,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: 0,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.RelatedSelfBuff.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: naturesVigilSpell.Spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (druid *Druid) registerYserasGift() {
	if !druid.Talents.YserasGift {
		return
	}

	healingSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 145108},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, target, 0.05*spell.Unit.MaxHealth(), spell.OutcomeHealing)
		},
	})

	druid.RegisterResetEffect(func(sim *core.Simulation) {
		core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period:   time.Second * 5,
			Priority: core.ActionPriorityDOT,

			OnAction: func(sim *core.Simulation) {
				healingSpell.Cast(sim, &druid.Unit)
			},
		})
	})
}

func (druid *Druid) registerRenewal() {
	if !druid.Talents.Renewal {
		return
	}

	renewalSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 108238},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful | core.SpellFlagAPL | core.SpellFlagIgnoreModifiers,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.CalcAndDealHealing(sim, spell.Unit, 0.3*spell.Unit.MaxHealth(), spell.OutcomeHealing)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: renewalSpell.Spell,
		Type:  core.CooldownTypeSurvival,
	})
}

func (druid *Druid) registerCenarionWard() {
	if !druid.Talents.CenarionWard {
		return
	}

	// First register the HoT spell that gets triggered when the target takes damage.
	baseTickDamage := 11.27999973297 * druid.ClassSpellScaling // ~12349

	// SP is snapshot at the time of the original buff cast according to simc
	var spSnapshot float64

	cenarionWardHot := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 102352},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskSpellHealing,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,
		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label: "Cenarion Ward (HoT)",
			},

			NumberOfTicks: 3,
			TickLength:    time.Second * 2,

			OnSnapshot: func(_ *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseTickDamage + spSnapshot*1.04
				dot.SnapshotAttackerMultiplier = dot.CasterPeriodicHealingMultiplier()
				dot.SnapshotCritChance = dot.Spell.HealingCritChance()
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(target).Apply(sim)
		},
	})

	// Then register the buff that triggers the HoT upon taking damage.
	buffActionID := core.ActionID{SpellID: 102351}

	buffConfig := core.Aura{
		Label:    "Cenarion Ward",
		ActionID: buffActionID,
		Duration: time.Second * 30,

		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
			if result.Damage > 0 {
				cenarionWardHot.Cast(sim, aura.Unit)
				aura.Deactivate(sim)
			}
		},
	}

	cenarionWardBuffs := druid.NewAllyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(buffConfig)
	})

	// Finally, register the spell that applies the buff.
	druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: buffActionID,
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 14.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},

			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			spSnapshot = cenarionWardHot.HealingPower(target)
			cenarionWardBuffs.Get(target).Activate(sim)
		},
	})
}

func (druid *Druid) registerForceOfNature() {
	if !druid.Talents.ForceOfNature {
		return
	}

	druid.ForceOfNature = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 106737},
		Flags:        core.SpellFlagAPL,
		Charges:      3,
		RechargeTime: time.Second * 20,

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			druid.Treants[spell.GetNumCharges()].Enable(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: druid.ForceOfNature.Spell,
		Type:  core.CooldownTypeDPS,
	})
}
