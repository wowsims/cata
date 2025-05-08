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
		ClassMask: SpellMaskChainLightning | SpellMaskLightningBolt | SpellMaskLavaBeam,
		Kind:      core.SpellMod_CastTime_Flat,
		TimeValue: time.Millisecond * -500,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask:  SpellMaskChainLightning | SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskChainLightningOverload | SpellMaskLavaBeam | SpellMaskLavaBeamOverload,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: 0.7,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskChainLightning | SpellMaskLavaBeam,
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

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Rolling Thunder",
		ActionID:       actionID,
		ClassSpellMask: SpellMaskChainLightning | SpellMaskChainLightningOverload | SpellMaskLightningBolt | SpellMaskLightningBoltOverload | SpellMaskLavaBeam | SpellMaskLavaBeamOverload,
		Callback:       core.CallbackOnSpellHitDealt,
		ProcChance:     0.6,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			nStack := core.TernaryInt32(shaman.T14Ele4pc.IsActive(), 2, 1)
			shaman.AddMana(sim, 0.02*shaman.MaxMana()*float64(nStack), manaMetrics)
			shaman.LightningShieldAura.Activate(sim)
			shaman.LightningShieldAura.SetStacks(sim, shaman.LightningShieldAura.GetStacks()+nStack)
		},
	})

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
	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Lava Surge",
		Duration: time.Second * 6,
		ActionID: core.ActionID{SpellID: 77762},
	}).AttachSpellMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaBurst,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -1.0,
	})

	core.MakePermanent(shaman.RegisterAura(core.Aura{
		Label:           "Lava Surge Proc Aura",
		ActionIDForProc: core.ActionID{SpellID: 77762},
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
					procAura.Activate(sim)
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
			if spell.ClassSpellMask != SpellMaskLavaBurst || !procAura.IsActive() {
				return
			}
			//If lava surge procs during LvB cast time, it is not consumed and lvb does not go on cd
			if spell.CurCast.CastTime > 0 {
				spell.CD.Reset()
				return
			}
			procAura.Deactivate(sim)
		},
	}))
}
