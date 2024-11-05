package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

func (shaman *Shaman) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_TotemRemainingTime:
		return shaman.newValueTotemRemainingTime(rot, config.GetTotemRemainingTime())
	case *proto.APLValue_ShamanCanSnapshotStrongerFireElemental:
		return shaman.newValueCanSnapshotStrongerFireElemental(config.GetShamanCanSnapshotStrongerFireElemental())
	case *proto.APLValue_ShamanFireElementalDuration:
		return shaman.newValueFireElementalDuration(config.GetShamanFireElementalDuration())
	default:
		return nil
	}
}

type APLValueTotemRemainingTime struct {
	core.DefaultAPLValueImpl
	shaman    *Shaman
	totemType proto.ShamanTotems_TotemType
}

func (shaman *Shaman) newValueTotemRemainingTime(rot *core.APLRotation, config *proto.APLValueTotemRemainingTime) core.APLValue {
	if config.TotemType == proto.ShamanTotems_TypeUnknown {
		rot.ValidationWarning("Totem Type required.")
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

type APLValueShamanCanSnapshotStrongerFireElemental struct {
	core.DefaultAPLValueImpl
	shaman *Shaman
}

func (shaman *Shaman) newValueCanSnapshotStrongerFireElemental(_ *proto.APLValueShamanCanSnapshotStrongerFireElemental) core.APLValue {
	return &APLValueShamanCanSnapshotStrongerFireElemental{
		shaman: shaman,
	}
}
func (value *APLValueShamanCanSnapshotStrongerFireElemental) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeBool
}
func (value *APLValueShamanCanSnapshotStrongerFireElemental) GetBool(sim *core.Simulation) bool {
	shaman := value.shaman

	if shaman.FireElemental.IsEnabled() {
		simulatedStats := shaman.fireElementalStatInheritance(0, 0)(shaman.GetStats())
		potentialFireElementalSpellPower := simulatedStats[stats.SpellPower]
		currentFireElementalSpellPower := shaman.FireElemental.GetPet().GetInheritedStats()[stats.SpellPower]
		return potentialFireElementalSpellPower > currentFireElementalSpellPower
	}

	return true
}

func (value *APLValueShamanCanSnapshotStrongerFireElemental) String() string {
	return "Can Snapshot Stronger Fire Elemental"
}

type APLValueShamanFireElementalDuration struct {
	core.DefaultAPLValueImpl
	shaman   *Shaman
	duration time.Duration
}

func (shaman *Shaman) newValueFireElementalDuration(_ *proto.APLValueShamanFireElementalDuration) core.APLValue {
	return &APLValueShamanFireElementalDuration{
		shaman:   shaman,
		duration: time.Second * time.Duration(120*(1.0+0.20*float64(shaman.Talents.TotemicFocus))),
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
