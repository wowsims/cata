package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (shaman *Shaman) ApplyElementalTalents() {

	//Elemental Precision
	shaman.AddStat(stats.SpellHitPercent, -shaman.GetBaseStats()[stats.Spirit]/core.SpellHitRatingPerHitPercent)
	shaman.AddStatDependency(stats.Spirit, stats.SpellHitPercent, 1.0/core.SpellHitRatingPerHitPercent)

	//Shamanism
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskChainLightning | SpellMaskLightningBolt,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * -500,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskChainLightning | SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightningOverload,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.7,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskChainLightning,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Second * -3,
	})

	// Elemental Fury
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskFire | SpellMaskNature |
			SpellMaskFrost | SpellMaskMagmaTotem | SpellMaskSearingTotem | SpellMaskEarthquake,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 0.5,
	})

	//Spiritual Insight
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskEarthShock | SpellMaskFlameShock,
		Kind:      core.SpellMod_Cooldown_Flat,
		TimeValue: time.Second * -1,
	})
	shaman.MultiplyStat(stats.Mana, 5)

	//Fulmination
	shaman.Fulmination = shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 88767},
		SpellSchool:    core.SpellSchoolNature,
		ProcMask:       core.ProcMaskSpellProc,
		Flags:          core.SpellFlagPassiveSpell,
		ClassSpellMask: SpellMaskFulmination,
		ManaCost: core.ManaCostOptions{
			BaseCostPercent: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			ModifyCast: func(s1 *core.Simulation, spell *core.Spell, c *core.Cast) {
				spell.SetMetricsSplit(shaman.LightningShieldAura.GetStacks() - 1)
			},
		},
		MetricSplits: 6,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			totalDamage := (shaman.ClassSpellScaling*0.56499999762 + 0.38800001144*spell.SpellPower()) * (float64(shaman.LightningShieldAura.GetStacks()) - 1)
			result := spell.CalcDamage(sim, target, totalDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Fulmination Proc",
		ProcChance:     1.0,
		ClassSpellMask: SpellMaskEarthShock,
		Callback:       core.CallbackOnCastComplete,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.Fulmination.Cast(sim, result.Target)
			shaman.LightningShieldAura.SetStacks(sim, 1)
		},
	})

	//Rolling Thunder
	actionID := core.ActionID{SpellID: 88765}
	manaMetrics := shaman.NewManaMetrics(actionID)

	wastedLSChargeAura := shaman.RegisterAura(core.Aura{
		Label:    "Wasted Lightning Shield Charge",
		Duration: core.NeverExpires,
		ActionID: core.ActionID{
			SpellID: 324,
			Tag:     1,
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Deactivate(sim)
		},
	})

	shaman.RegisterAura(*core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Rolling Thunder",
		ActionID:       actionID,
		ClassSpellMask: SpellMaskChainLightning | SpellMaskChainLightningOverload | SpellMaskLightningBolt | SpellMaskLightningBoltOverload,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcChance:     0.6,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.AddMana(sim, 0.02*shaman.MaxMana(), manaMetrics)
			if shaman.LightningShieldAura.GetStacks() == 7 {
				wastedLSChargeAura.Activate(sim)
			}
			shaman.LightningShieldAura.Activate(sim)
			shaman.LightningShieldAura.AddStack(sim)
		},
	}))

	//Elemental Focus
	var triggeringSpell *core.Spell
	var triggerTime time.Duration

	canTriggerSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskLavaBurst | SpellMaskFireNova | SpellMaskEarthShock | SpellMaskFlameShock | SpellMaskFrostShock
	canConsumeSpells := canTriggerSpells
	costReductionMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_PowerCost_Pct,
		ClassMask: canConsumeSpells,
		IntValue:  -25,
	})

	damageMod := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		School:     core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		FloatValue: 0.2,
	})
	damageModEarthquake := shaman.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Flat,
		ClassMask:  SpellMaskEarthquake,
		FloatValue: 0.2,
	})

	maxStacks := int32(2)

	// TODO: verify what procs and what consumes in game
	clearcastingAura := shaman.RegisterAura(core.Aura{
		Label:     "Clearcasting",
		ActionID:  core.ActionID{SpellID: 16246},
		Duration:  time.Second * 15,
		MaxStacks: maxStacks,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			costReductionMod.Activate()
			damageMod.Activate()
			damageModEarthquake.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			costReductionMod.Deactivate()
			damageMod.Deactivate()
			damageModEarthquake.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Flags.Matches(SpellFlagShock|SpellFlagFocusable) || (spell.ClassSpellMask&(SpellMaskOverload|SpellMaskThunderstorm) != 0) {
				return
			}
			if spell == triggeringSpell && sim.CurrentTime == triggerTime {
				return
			}
			aura.RemoveStack(sim)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Elemental Focus",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagShock|SpellFlagFocusable) || (spell.ClassSpellMask&(SpellMaskOverload|SpellMaskUnleashFlame|SpellMaskEarthquake) != 0) {
				return
			}
			if !result.Outcome.Matches(core.OutcomeCrit) {
				return
			}
			triggeringSpell = spell
			triggerTime = sim.CurrentTime
			clearcastingAura.Activate(sim)
			clearcastingAura.SetStacks(sim, maxStacks)
		},
	})

	//Lava Surge
	//TODO If it procs during lava burst cast time, lava burst won't get on cd so the proc can be used and not wasted
	instantLavaBurstMod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaBurst,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1.0,
	})
	shaman.RegisterAura(core.Aura{
		Label:    "Lava Surge",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			//TODO verify proc chance in game
			if spell.ClassSpellMask != SpellMaskFlameShockDot || !sim.Proc(0.2, "LavaSurge") {
				return
			}

			// Set up a PendingAction to reset the CD just after this
			// timestep rather than immediately. This guarantees that
			// an existing Lava Burst cast that is set to finish on
			// this timestep will apply the cooldown *before* it gets
			// reset by the Lava Surge proc.
			pa := &core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Duration(1),
				Priority:     core.ActionPriorityDOT,

				OnAction: func(sim *core.Simulation) {
					shaman.LavaBurst.CD.Reset()
					instantLavaBurstMod.Activate()
				},
			}
			sim.AddPendingAction(pa)

			// Additionally, trigger a rotational wait so that the agent has an
			// opportunity to cast another Lava Burst after the reset, rather
			// than defaulting to a lower priority spell. Since this Lava Burst
			// cannot be spell queued (the CD was only just now reset), apply
			// input delay to the rotation call.
			if shaman.RotationTimer.IsReady(sim) {
				shaman.WaitUntil(sim, sim.CurrentTime+shaman.ReactionTime)
			}
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.ClassSpellMask != SpellMaskLavaBurst || !instantLavaBurstMod.IsActive {
				return
			}
			//If lava surge procs during LvB cast time, it is not consumed
			if spell.CurCast.CastTime > 0 {
				return
			}
			instantLavaBurstMod.Deactivate()
		},
	})
}
