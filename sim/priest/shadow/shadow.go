package shadow

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/priest"
)

func RegisterShadowPriest() {
	core.RegisterAgentFactory(
		proto.Player_ShadowPriest{},
		proto.Spec_SpecShadowPriest,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewShadowPriest(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_ShadowPriest)
			if !ok {
				panic("Invalid spec value for Shadow Priest!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewShadowPriest(character *core.Character, options *proto.Player) *ShadowPriest {
	shadowOptions := options.GetShadowPriest()

	selfBuffs := priest.SelfBuffs{
		UseShadowfiend: true,
		UseInnerFire:   shadowOptions.Options.ClassOptions.Armor == proto.PriestOptions_InnerFire,
	}

	basePriest := priest.New(character, selfBuffs, options.TalentsString)
	basePriest.Latency = float64(basePriest.ChannelClipDelay.Milliseconds())
	spriest := &ShadowPriest{
		Priest:  basePriest,
		options: shadowOptions.Options,
	}

	// TODO: Fix this to work with the new talent system.
	// spriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	// if spriest.Talents.PowerInfusion && shadowOptions.Options.PowerInfusionTarget != nil {
	// 	spriest.SelfBuffs.PowerInfusionTarget = shadowOptions.Options.PowerInfusionTarget
	// }

	return spriest
}

type ShadowPriest struct {
	*priest.Priest
	options *proto.ShadowPriest_Options
}

func (spriest *ShadowPriest) GetPriest() *priest.Priest {
	return spriest.Priest
}

func (spriest *ShadowPriest) Initialize() {
	spriest.Priest.Initialize()
}

func (spriest *ShadowPriest) Reset(sim *core.Simulation) {
	spriest.Priest.Reset(sim)
}

func getMasteryBonus(masteryPoints float64) float64 {
	return (21.6 + masteryPoints*1.45) / 100
}

func (spriest *ShadowPriest) ApplyTalents() {
	// spriest.Priest.ApplyTalents()

	// apply shadow spec specific auras
	// make it an aura so it's visible that it's used in the timeline
	spriest.AddStaticMod(core.SpellModConfig{
		FloatValue: 0.15,
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	spriest.RegisterAura(
		core.Aura{
			Label:    "ShadowPower",
			Duration: core.NeverExpires,
			ActionID: core.ActionID{
				SpellID: 87327,
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
		},
	)

	// Shadow Power
	spriest.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 1.0,
		School:     core.SpellSchoolShadow,
		ClassMask:  int64(priest.PriestShadowSpells),
	})

	shadowOrbMod := spriest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  int64(priest.PriestSpellMindBlast) | int64(priest.PriestSpellMindSpike),
		FloatValue: getMasteryBonus(spriest.GetMasteryPoints()),
		Kind:       core.SpellMod_DamageDone_Pct,
	})

	// mastery aura
	spriest.ShadowOrbsAura = spriest.RegisterAura(core.Aura{
		Label:     "Shadow Orb",
		ActionID:  core.ActionID{SpellID: 77487},
		Duration:  time.Minute,
		MaxStacks: 3,
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, oldStacks int32, newStacks int32) {
			shadowOrbMod.UpdateFloatValue(getMasteryBonus(spriest.GetMasteryPoints()) * float64(newStacks))
			shadowOrbMod.Activate()
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if spell.ClassSpellMask&(priest.PriestSpellMindBlast|priest.PriestSpellMindSpike) == 0 {
				return
			}

			spriest.EmpoweredShadowAura.Deactivate(sim)
			spriest.EmpoweredShadowAura.Activate(sim)
			aura.Deactivate(sim)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shadowOrbMod.Deactivate()
		},
	})

	spriest.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		shadowOrbMod.UpdateFloatValue(getMasteryBonus(core.MasteryRatingToMasteryPoints(newMastery)) * float64(spriest.ShadowOrbsAura.GetStacks()))
	})

	empoweredShadowMod := spriest.AddDynamicMod(core.SpellModConfig{
		ClassMask:  priest.PriestSpellDoT | priest.PriestSpellMindSear,
		Kind:       core.SpellMod_DamageDone_Flat,
		FloatValue: getMasteryBonus(spriest.GetMasteryPoints()),
	})

	spriest.EmpoweredShadowAura = spriest.RegisterAura(core.Aura{
		Label:    "Empowered Shadow",
		ActionID: core.ActionID{SpellID: 95799},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			empoweredShadowMod.UpdateFloatValue(getMasteryBonus(aura.Unit.GetMasteryPoints()))
			empoweredShadowMod.Activate()
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			empoweredShadowMod.Deactivate()
		},
	})

	spriest.RegisterAura(core.Aura{
		Label:    "Shadow Orb Power",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			handleShadowOrbPower(spriest, sim, spell, result)
		},
	})
}

func handleShadowOrbPower(spriest *ShadowPriest, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	if !result.Landed() {
		return
	}

	if spell.ClassSpellMask&(priest.PriestSpellShadowWordPain|priest.PriestSpellMindFlay) > 0 {
		// TODO: Fix this to work with the new talent system.
		// procChance := 0.1 + float64(spriest.Talents.HarnessedShadows)*0.04
		procChance := 0.1
		if sim.Proc(procChance, "Shadow Orb Power") {
			spriest.ShadowOrbsAura.Activate(sim)
			spriest.ShadowOrbsAura.AddStack(sim)
		}
	}
}
