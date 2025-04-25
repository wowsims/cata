package balance

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/druid"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewBalanceDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BalanceDruid)
			if !ok {
				panic("Invalid spec value for Balance Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBalanceDruid(character *core.Character, options *proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	moonkin := &BalanceDruid{
		Druid:   druid.New(character, druid.Moonkin, selfBuffs, options.TalentsString),
		Options: balanceOptions.Options,
	}

	moonkin.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if balanceOptions.Options.ClassOptions.InnervateTarget != nil {
		moonkin.SelfBuffs.InnervateTarget = balanceOptions.Options.ClassOptions.InnervateTarget
	}

	return moonkin
}

type BalanceOnUseTrinket struct {
	Cooldown *core.MajorCooldown
	Stat     stats.Stat
}

type BalanceDruid struct {
	*druid.Druid
	Options *proto.BalanceDruid_Options
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()

	moonkin.RegisterBalanceSpells()

	// if moonkin.OwlkinFrenzyAura != nil && moonkin.Options.OkfUptime > 0 {
	// 	moonkin.Env.RegisterPreFinalizeEffect(func() {
	// 		core.ApplyFixedUptimeAura(moonkin.OwlkinFrenzyAura, float64(moonkin.Options.OkfUptime), time.Second*5, 0)
	// 	})
	// }
}

func (moonkin *BalanceDruid) ApplyTalents() {

	// moonkin.EnableEclipseBar()
	// moonkin.RegisterEclipseAuras()
	// moonkin.RegisterEclipseEnergyGainAura()

	// Moonfury passive
	moonkin.RegisterAura(
		core.Aura{
			Label:    "Moonfury",
			Duration: core.NeverExpires,
			ActionID: core.ActionID{
				SpellID: 16913,
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
		},
	)

	moonkin.AddStaticMod(core.SpellModConfig{
		ClassMask:  druid.DruidSpellWrath | druid.DruidSpellStarfire | druid.DruidSpellStarsurge | druid.DruidSpellStarfall | druid.DruidSpellDoT,
		Kind:       core.SpellMod_CritMultiplier_Flat,
		FloatValue: 1.0,
	})

	moonkin.AddStaticMod(core.SpellModConfig{
		School:     core.SpellSchoolArcane | core.SpellSchoolNature,
		ClassMask:  druid.DruidSpellsAll,
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: 0.1,
	})

	// Apply druid talents
	// moonkin.Druid.ApplyTalents()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.Druid.Reset(sim)
	//moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")
}
