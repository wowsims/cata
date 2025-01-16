package priest

// Pre-cache values instead of checking gear during sim runs
func (priest *Priest) registerSetBonuses() {
	priest.T7TwoSetBonus = priest.CouldHaveSetBonus(ItemSetValorous, 2)
	priest.T7FourSetBonus = priest.CouldHaveSetBonus(ItemSetValorous, 4)
	priest.T8TwoSetBonus = priest.CouldHaveSetBonus(ItemSetConquerorSanct, 2)
	priest.T8FourSetBonus = priest.CouldHaveSetBonus(ItemSetConquerorSanct, 4)
	priest.T9TwoSetBonus = priest.CouldHaveSetBonus(ItemSetZabras, 2)
	priest.T9FourSetBonus = priest.CouldHaveSetBonus(ItemSetZabras, 4)
	priest.T10TwoSetBonus = priest.CouldHaveSetBonus(ItemSetCrimsonAcolyte, 2)
	priest.T10FourSetBonus = priest.CouldHaveSetBonus(ItemSetCrimsonAcolyte, 4)
}
