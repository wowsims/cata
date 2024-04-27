package balance

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/druid"
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

	moonkin.Character.PrimaryStat = stats.Intellect

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
	eclipseEnergyBar
	Options *proto.BalanceDruid_Options
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()
	moonkin.EnableEclipseBar()

	moonkin.RegisterTestSpells()
	//moonkin.RegisterBalanceSpells()

	// if moonkin.OwlkinFrenzyAura != nil && moonkin.Options.OkfUptime > 0 {
	// 	moonkin.Env.RegisterPreFinalizeEffect(func() {
	// 		core.ApplyFixedUptimeAura(moonkin.OwlkinFrenzyAura, float64(moonkin.Options.OkfUptime), time.Second*5, 0)
	// 	})
	// }
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.Druid.Reset(sim)
	moonkin.eclipseEnergyBar.reset()
	//moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")
}
