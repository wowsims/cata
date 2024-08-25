package druid

import (
	"sync"
)

type EclipseEnergyValues struct {
	InEclipse float64
	NoEclipse float64
}

var EclipseEnergyMap = map[int64]EclipseEnergyValues{}
var rwMutex sync.RWMutex

func (druid *Druid) SetSpellEclipseEnergy(spellMask int64, inEclipseEnergy float64, noEclipseEnergy float64) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	EclipseEnergyMap[spellMask] = EclipseEnergyValues{
		InEclipse: inEclipseEnergy,
		NoEclipse: noEclipseEnergy,
	}
}

func (druid *Druid) GetSpellEclipseEnergy(spellMask int64, inEclipse bool) float64 {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	energyValue := EclipseEnergyMap[spellMask]

	if inEclipse {
		return energyValue.InEclipse
	}

	return energyValue.NoEclipse
}
