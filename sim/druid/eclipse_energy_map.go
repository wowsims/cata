package druid

type EclipseEnergyValues struct {
	InEclipse float64
	NoEclipse float64
}

type EclipseEnergyMap = map[int64]EclipseEnergyValues

func (druid *Druid) SetSpellEclipseEnergy(spellMask int64, inEclipseEnergy float64, noEclipseEnergy float64) {
	druid.EclipseEnergyMap[spellMask] = EclipseEnergyValues{
		InEclipse: inEclipseEnergy,
		NoEclipse: noEclipseEnergy,
	}
}

func (druid *Druid) GetSpellEclipseEnergy(spellMask int64, inEclipse bool) float64 {
	energyValue := druid.EclipseEnergyMap[spellMask]

	if inEclipse {
		return energyValue.InEclipse
	}

	return energyValue.NoEclipse
}
