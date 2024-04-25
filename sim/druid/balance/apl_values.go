package balance

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (druid *BalanceDruid) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CurrentSolarEnergy:
		return druid.newValueCurrentSolarEnergy(config.GetCurrentSolarEnergy())
	case *proto.APLValue_CurrentLunarEnergy:
		return druid.newValueCurrentLunarEnergy(config.GetCurrentLunarEnergy())
	case *proto.APLValue_DruidCurrentEclipsePhase:
		return druid.newValueCurrentEclipsePhase(config.GetDruidCurrentEclipsePhase())
	default:
		return nil
	}
}

type APLValueCurrentSolarEnergy struct {
	core.DefaultAPLValueImpl
	druid *BalanceDruid
}

func (druid *BalanceDruid) newValueCurrentSolarEnergy(_ *proto.APLValueCurrentSolarEnergy) core.APLValue {
	return &APLValueCurrentSolarEnergy{
		druid: druid,
	}
}

func (value *APLValueCurrentSolarEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}

func (value *APLValueCurrentSolarEnergy) GetInt(sim *core.Simulation) int32 {
	return int32(value.druid.CurrentSolarEnergy())
}

func (value *APLValueCurrentSolarEnergy) String() string {
	return "Current Solar Energy"
}

type APLValueCurrentLunarEnergy struct {
	core.DefaultAPLValueImpl
	druid *BalanceDruid
}

func (druid *BalanceDruid) newValueCurrentLunarEnergy(_ *proto.APLValueCurrentLunarEnergy) core.APLValue {
	return &APLValueCurrentLunarEnergy{
		druid: druid,
	}
}

func (value *APLValueCurrentLunarEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeInt
}

func (value *APLValueCurrentLunarEnergy) GetInt(sim *core.Simulation) int32 {
	return int32(value.druid.CurrentLunarEnergy())
}

func (value *APLValueCurrentLunarEnergy) String() string {
	return "Current Solar Energy"
}

type APLValueCurrentEclipsePhase struct {
	core.DefaultAPLValueImpl
	phase proto.APLValueEclipsePhase
	druid *BalanceDruid
}

func (value *APLValueCurrentEclipsePhase) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}

func (value *APLValueCurrentEclipsePhase) GetBool(sim *core.Simulation) bool {
	if value.druid.gainMask&SolarAndLunarEnergy == SolarAndLunarEnergy {
		return value.phase == proto.APLValueEclipsePhase_NeutralPhase

		// if we can only gain lunar energy we're in solar eclipse phase
	} else if value.druid.gainMask&LunarEnergy > 0 {
		return value.phase == proto.APLValueEclipsePhase_SolarPhase
	}

	return value.phase == proto.APLValueEclipsePhase_LunarPhase
}

func (value *APLValueCurrentEclipsePhase) String() string {
	return "Current Eclipse Phase"
}

func (druid *BalanceDruid) newValueCurrentEclipsePhase(config *proto.APLValueCurrentEclipsePhase) core.APLValue {
	return &APLValueCurrentEclipsePhase{
		druid: druid,
		phase: config.EclipsePhase,
	}
}
