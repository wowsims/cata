package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (shaman *Shaman) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_TotemRemainingTime:
		return shaman.newValueTotemRemainingTime(rot, config.GetTotemRemainingTime(), config.Uuid)
	case *proto.APLValue_ShamanFireElementalDuration:
		return shaman.newValueFireElementalDuration(config.GetShamanFireElementalDuration(), config.Uuid)
	default:
		return nil
	}
}

type APLValueTotemRemainingTime struct {
	core.DefaultAPLValueImpl
	shaman    *Shaman
	totemType proto.ShamanTotems_TotemType
}

func (shaman *Shaman) newValueTotemRemainingTime(rot *core.APLRotation, config *proto.APLValueTotemRemainingTime, uuid *proto.UUID) core.APLValue {
	if config.TotemType == proto.ShamanTotems_TypeUnknownTotem {
		rot.ValidationMessageByUUID(uuid, proto.LogLevel_Warning, "Totem Type required.")
		return nil
	}
	return &APLValueTotemRemainingTime{
		shaman:    shaman,
		totemType: config.TotemType,
	}
}
func (value *APLValueTotemRemainingTime) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueTotemRemainingTime) GetDuration(sim *core.Simulation) time.Duration {
	if value.totemType == proto.ShamanTotems_Earth {
		return max(0, value.shaman.TotemExpirations[EarthTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Air {
		return max(0, value.shaman.TotemExpirations[AirTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Fire {
		return max(0, value.shaman.TotemExpirations[FireTotem]-sim.CurrentTime)
	} else if value.totemType == proto.ShamanTotems_Water {
		return max(0, value.shaman.TotemExpirations[WaterTotem]-sim.CurrentTime)
	} else {
		return 0
	}
}
func (value *APLValueTotemRemainingTime) String() string {
	return fmt.Sprintf("Totem Remaining Time(%s)", value.totemType.String())
}

type APLValueShamanFireElementalDuration struct {
	core.DefaultAPLValueImpl
	shaman   *Shaman
	duration time.Duration
}

func (shaman *Shaman) newValueFireElementalDuration(_ *proto.APLValueShamanFireElementalDuration, _ *proto.UUID) core.APLValue {
	return &APLValueShamanFireElementalDuration{
		shaman:   shaman,
		duration: time.Second * 60,
	}
}

func (value *APLValueShamanFireElementalDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}

func (value *APLValueShamanFireElementalDuration) GetDuration(sim *core.Simulation) time.Duration {
	return value.duration
}

func (value *APLValueShamanFireElementalDuration) String() string {
	return "Fire Elemental Total Duration"
}
