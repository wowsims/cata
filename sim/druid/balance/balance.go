package balance

import (
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/druid"
)

const (
	WrathBaseEnergyGain     float64 = 15
	StarsurgeBaseEnergyGain float64 = 20
	StarfireBaseEnergyGain  float64 = 20
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
		Druid:            druid.New(character, druid.Moonkin, selfBuffs, options.TalentsString),
		Options:          balanceOptions.Options,
		EclipseEnergyMap: make(EclipseEnergyMap),
	}

	moonkin.registerTreants()

	moonkin.SelfBuffs.InnervateTarget = &proto.UnitReference{}
	if balanceOptions.Options.ClassOptions.InnervateTarget != nil {
		moonkin.SelfBuffs.InnervateTarget = balanceOptions.Options.ClassOptions.InnervateTarget
	}

	return moonkin
}

type BalanceDruid struct {
	*druid.Druid
	eclipseEnergyBar
	Options *proto.BalanceDruid_Options

	EclipseEnergyMap EclipseEnergyMap

	AstralCommunion      *druid.DruidSpell
	AstralStorm          *druid.DruidSpell
	AstralStormTickSpell *druid.DruidSpell
	CelestialAlignment   *druid.DruidSpell
	ChosenOfElune        *druid.DruidSpell
	Starfall             *druid.DruidSpell
	Starfire             *druid.DruidSpell
	Sunfire              *druid.DruidSpell
	Starsurge            *druid.DruidSpell

	AstralInsight   *core.Aura // Soul of the Forest
	DreamOfCenarius *core.Aura
	NaturesGrace    *core.Aura
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()

	moonkin.EnableEclipseBar()
	moonkin.RegisterEclipseAuras()
	moonkin.RegisterEclipseEnergyGainAura()

	moonkin.RegisterBalancePassives()
	moonkin.RegisterBalanceSpells()
}

func (moonkin *BalanceDruid) ApplyTalents() {
	moonkin.Druid.ApplyTalents()

	moonkin.ApplyBalanceTalents()
}

func (moonkin *BalanceDruid) RegisterBalanceSpells() {
	moonkin.registerSunfireSpell()
	moonkin.registerStarfireSpell()
	moonkin.registerStarsurgeSpell()
	moonkin.registerStarfallSpell()
	moonkin.registerAstralCommunionSpell()
	moonkin.registerCelestialAlignmentSpell()
	moonkin.registerAstralStormSpell()
	moonkin.registerWildMushrooms()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.eclipseEnergyBar.reset()
	moonkin.Druid.Reset(sim)
}
