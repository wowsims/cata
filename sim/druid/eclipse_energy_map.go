package druid

type EclipseEnergyValues struct {
	InEclipse float64
	NoEclipse float64
}

var EclipseEnergyMap = map[int32]EclipseEnergyValues{}

func (druid *Druid) SetSpellEclipseEnergy(spellID int32, inEclipseEnergy float64, noEclipseEnergy float64) {
	EclipseEnergyMap[spellID] = EclipseEnergyValues{
		InEclipse: inEclipseEnergy,
		NoEclipse: noEclipseEnergy,
	}
}

func (druid *Druid) SetSpellEclipseEnergyValues(spellID int32, values EclipseEnergyValues) {
	EclipseEnergyMap[spellID] = values
}

func (druid *Druid) GetSpellEclipseEnergy(spellID int32, inEclipse bool) float64 {
	energyValue := EclipseEnergyMap[spellID]

	if inEclipse {
		return energyValue.InEclipse
	}

	return energyValue.NoEclipse
}

func (druid *Druid) GetSpellEclipseEnergyValues(spellID int32) EclipseEnergyValues {
	energyValues := EclipseEnergyMap[spellID]

	return energyValues
}
