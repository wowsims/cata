package shaman

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func (shaman *Shaman) ApplyEnhancementTalents() {

	//Mental Quickness (AP -> SP in enhancement.go)
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskShock,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -90,
	})
	shaman.AddStaticMod(core.SpellModConfig{
		ClassMask: SpellMaskTotem | SpellMaskInstantSpell,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -75,
	})
	primalWisdomManaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: 63375})
	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:       "Mental Quickness",
		ProcMask:   core.ProcMaskMelee,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcChance: 0.4,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.AddMana(sim, 0.05*shaman.MaxMana(), primalWisdomManaMetrics)
		},
	})

	//Flurry
	flurryProcAura := shaman.RegisterAura(core.Aura{
		Label:     "Flurry Proc",
		ActionID:  core.ActionID{SpellID: 16278},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1.15)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1/1.15)
		},
	}).AttachStatDependency(shaman.NewDynamicMultiplyStat(stats.HasteRating, 1.5))

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:     "Flurry",
		ProcMask: core.ProcMaskMelee | core.ProcMaskMeleeProc,
		Callback: core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Outcome.Matches(core.OutcomeCrit) {
				flurryProcAura.Activate(sim)
				flurryProcAura.SetStacks(sim, 5)
				return
			}

			// Remove a stack.
			if flurryProcAura.IsActive() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				flurryProcAura.RemoveStack(sim)
			}
		},
	})

	//Searing Flames
	ftmod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskFlametongueWeapon,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.08,
	})
	llmod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask:  SpellMaskLavaLash,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.2,
	})

	searingFlameStackingAura := shaman.RegisterAura(core.Aura{
		Label:     "Searing Flames",
		ActionID:  core.ActionID{SpellID: 77661},
		Duration:  time.Second * 15,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			ftmod.UpdateFloatValue(float64(newStacks) * 0.08)
			ftmod.Activate()
			llmod.UpdateFloatValue(float64(newStacks) * 0.2)
			llmod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			ftmod.Deactivate()
			llmod.Deactivate()
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Matches(SpellMaskLavaLash) {
				return
			}
			aura.Deactivate(sim)
		},
	})

	core.MakeProcTriggerAura(&shaman.FireElemental.Unit, core.ProcTrigger{
		Name:           "Searing Flames Dummy Fire ele",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskFireElementalMelee,
		Outcome:        core.OutcomeLanded,
		ProcChance:     1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			searingFlameStackingAura.Activate(sim)
			searingFlameStackingAura.AddStack(sim)
		},
	})
	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Searing Flames Dummy Shaman",
		Callback:       core.CallbackOnSpellHitDealt,
		ClassSpellMask: SpellMaskSearingTotem,
		Outcome:        core.OutcomeLanded,
		ProcChance:     1,

		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			searingFlameStackingAura.Activate(sim)
			searingFlameStackingAura.AddStack(sim)
		},
	})

	//Static Shock
	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:           "Static Shock",
		Callback:       core.CallbackOnApplyEffects,
		ClassSpellMask: SpellMaskStormstrikeCast | SpellMaskLavaLash,
		ProcChance:     0.45,
		ExtraCondition: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return shaman.SelfBuffs.Shield == proto.ShamanShield_LightningShield
		},
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.LightningShieldDamage.Cast(sim, result.Target)
		},
	})

	//Maelstrom Weapon
	mwAffectedSpells := SpellMaskLightningBolt | SpellMaskChainLightning | SpellMaskEarthShock | SpellMaskElementalBlast
	mwCastTimemod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask:  mwAffectedSpells,
		Kind:       core.SpellMod_CastTime_Pct,
		FloatValue: -0.2,
	})
	mwManaCostmod := shaman.AddDynamicMod(core.SpellModConfig{
		ClassMask: mwAffectedSpells,
		Kind:      core.SpellMod_PowerCost_Pct,
		IntValue:  -20,
	})
	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:     "MaelstromWeapon Proc",
		ActionID:  core.ActionID{SpellID: 51530},
		Duration:  time.Second * 30,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			mwCastTimemod.UpdateFloatValue(float64(newStacks) * -0.2)
			mwCastTimemod.Activate()
			mwManaCostmod.UpdateIntValue(newStacks * -20)
			mwManaCostmod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mwCastTimemod.Deactivate()
			mwManaCostmod.Deactivate()
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.Matches(mwAffectedSpells) {
				return
			}
			//If AS is active and MW < 5 stacks, do not consume MW stacks
			//As i don't know which OnCastComplete is going to be executed first, check here if AS has not just been consumed/is active
			if aura.GetStacks() < 5 && shaman.Talents.AncestralSwiftness && shaman.AncestralSwiftnessInstantAura.TimeInactive(sim) == 0 {
				return
			}
			shaman.MaelstromWeaponAura.Deactivate(sim)
		},
	})

	dpm := shaman.NewLegacyPPMManager(10.0, core.ProcMaskMeleeOrMeleeProc)

	// This aura is hidden, just applies stacks of the proc aura.
	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:     "Maelstrom Weapon",
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
		DPM:      dpm,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			shaman.MaelstromWeaponAura.Activate(sim)
			shaman.MaelstromWeaponAura.AddStack(sim)
		},
	})
}
