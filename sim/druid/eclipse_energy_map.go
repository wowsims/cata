package druid

import (
	"sync"
)

type EclipseEnergyValues struct {
	InEclipse float64
	NoEclipse float64
}

var EclipseEnergyMap = map[int32]EclipseEnergyValues{}
var rwMutex sync.RWMutex

func (druid *Druid) SetSpellEclipseEnergy(spellID int32, inEclipseEnergy float64, noEclipseEnergy float64) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	EclipseEnergyMap[spellID] = EclipseEnergyValues{
		InEclipse: inEclipseEnergy,
		NoEclipse: noEclipseEnergy,
	}
}

func (druid *Druid) GetSpellEclipseEnergy(spellID int32, inEclipse bool) float64 {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	energyValue := EclipseEnergyMap[spellID]

	if inEclipse {
		return energyValue.InEclipse
	}

	return energyValue.NoEclipse
}
