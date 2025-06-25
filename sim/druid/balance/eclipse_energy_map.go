package balance

type EclipseEnergyValues struct {
	InEclipse float64
	NoEclipse float64
}

type EclipseEnergyMap = map[int64]EclipseEnergyValues

func (moonkin *BalanceDruid) SetSpellEclipseEnergy(spellMask int64, inEclipseEnergy float64, noEclipseEnergy float64) {
	moonkin.EclipseEnergyMap[spellMask] = EclipseEnergyValues{
		InEclipse: inEclipseEnergy,
		NoEclipse: noEclipseEnergy,
	}
}

func (moonkin *BalanceDruid) GetSpellEclipseEnergy(spellMask int64, inEclipse bool) float64 {
	energyValue := moonkin.EclipseEnergyMap[spellMask]

	if inEclipse {
		return energyValue.InEclipse
	}

	return energyValue.NoEclipse
}
