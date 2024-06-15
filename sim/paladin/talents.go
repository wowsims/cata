package paladin

func (paladin *Paladin) ApplyTalents() {
	paladin.applyRetributionTalents()
	paladin.applyProtectionTalents()
	paladin.applyHolyTalents()
}
