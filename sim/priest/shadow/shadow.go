package shadow

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/priest"
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

	// add spirit -> spell hit conversion for Twisted Faith talent
	if spriest.Talents.TwistedFaith > 0 {
		spriest.AddStatDependency(stats.Spirit, stats.SpellHit, 0.5*float64(spriest.Talents.TwistedFaith))
	}

	spriest.SelfBuffs.PowerInfusionTarget = &proto.UnitReference{}
	if spriest.Talents.PowerInfusion && shadowOptions.Options.PowerInfusionTarget != nil {
		spriest.SelfBuffs.PowerInfusionTarget = shadowOptions.Options.PowerInfusionTarget
	}

	return spriest
}

type ShadowPriest struct {
	*priest.Priest
	options *proto.ShadowPriest_Options

	shadowOrbsAura      *core.Aura
	empoweredShadowAura *core.Aura
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

func (spriest *ShadowPriest) ApplyTalents() {

	// apply shadow spec specific auras
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
			OnGain: func(_ *core.Aura, _ *core.Simulation) {
				// only shadow damage here
				spriest.Priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] += 0.15
			},
			OnExpire: func(_ *core.Aura, _ *core.Simulation) {
				spriest.Priest.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] -= 0.15
			},
		},
	)

	spriest.Priest.ShadowCritMultiplier = 1.0

	// mastery aura
	spriest.shadowOrbsAura = spriest.RegisterAura(core.Aura{
		Label:     "Shadow Orb",
		ActionID:  core.ActionID{SpellID: 77487},
		Duration:  time.Minute,
		MaxStacks: 3,
		OnStacksChange: func(_ *core.Aura, _ *core.Simulation, _ int32, newStacks int32) {
			priest.AddOrReplaceMod(&spriest.DamageDonePercentAddMods, &priest.PriestAuraMod[float64]{
				SpellID:    77487,
				ClassSpell: priest.PriestSpellMindBlast | priest.PriestSpellMindSpike,

				// 10% + 11.6 base value from our mastery
				BaseValue: 0.216,

				// add 1.45% dmg per mastery point
				DynamicValue: func(p *priest.Priest) float64 {
					return p.GetMasteryPoints() * 0.0145
				},
				Stacks: newStacks,
			})
		},

		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
			if !spriest.MindBlast.IsEqual(spell) && !spriest.MindSpike.IsEqual(spell) {
				return
			}

			spriest.empoweredShadowAura.Activate(sim)
			aura.Deactivate(sim)
		},
	})

	spriest.empoweredShadowAura = spriest.RegisterAura(core.Aura{
		Label:    "Empowered Shadow",
		ActionID: core.ActionID{SpellID: 95799},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			priest.AddOrReplaceMod(&spriest.DamageDonePercentAddMods, &priest.PriestAuraMod[float64]{
				SpellID:    95799,
				ClassSpell: priest.PriestSpellDoT,
				School:     core.SpellSchoolShadow,

				// since we're simming 85 players 10 + 11.6 mastery base value
				BaseValue: 0.216,
				DynamicValue: func(p *priest.Priest) float64 {
					return p.GetMasteryPoints() * 0.0145
				},
			})
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			priest.RemoveMod(&spriest.DamageDonePercentAddMods, 95799)
		},
	})

	spriest.RegisterAura(core.Aura{
		Label:    "Shadow Orb Power",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},

		OnSpellHitDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			handleShadowOrbPower(spriest, sim, spell, result)
		},

		OnPeriodicDamageDealt: func(_ *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			handleShadowOrbPower(spriest, sim, spell, result)
		},
	})

	spriest.Priest.ApplyTalents()
}

func handleShadowOrbPower(spriest *ShadowPriest, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
	if !result.Landed() {
		return
	}

	if spell == spriest.ShadowWordPain.Spell || spell.SpellID == spriest.MindFlayAPL.SpellID {
		procChance := spriest.GetClassSpellProcChance(0.1, priest.PriestSpellShadowOrbPassive, core.SpellSchoolShadow)
		if sim.RandomFloat("Shadow Orb Power") < procChance {
			spriest.shadowOrbsAura.Activate(sim)
			spriest.shadowOrbsAura.AddStack(sim)
		}
	}
}
