package dbc

type SpellPowerData struct {
	ID             uint
	SpellID        uint
	AuraID         uint
	PowerType      int
	Cost           int
	CostMax        int
	CostPerTick    int
	PctCost        float64
	PctCostMax     float64
	PctCostPerTick float64
}

func (sd *SpellPowerData) GetMaxCost() float64 {
	if sd.CostMax != 0 {
		return float64(sd.CostMax) / sd.costDivisor(!(sd.Cost != 0))
	}
	return float64(sd.PctCostMax) / sd.costDivisor(!(sd.Cost != 0))
}

func (sd *SpellPowerData) GetCostPerTick() float64 {
	return float64(sd.CostPerTick) / sd.costDivisor(!(sd.Cost != 0))
}

func (sd *SpellPowerData) GetCost() float64 {
	cost := 0.0
	if sd.Cost != 0 {
		cost = float64(sd.Cost)
	} else {
		cost = sd.PctCost
	}
	return cost / sd.costDivisor(!(sd.Cost != 0))
}
func (sd *SpellPowerData) costDivisor(percentage bool) float64 {
	switch sd.PowerType {
	case POWER_MANA:
		if percentage {
			return 100.0
		}
		return 1.0
	case POWER_RAGE, POWER_RUNIC_POWER, POWER_ASTRAL_POWER, POWER_SOUL_SHARDS:
		return 10.0
	default:
		return 1.0
	}

}
