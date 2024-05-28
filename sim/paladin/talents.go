package paladin

func (paladin *Paladin) ApplyTalents() {
	paladin.ApplyRetributionTalents()
	paladin.ApplyProtectionTalents()
	paladin.ApplyHolyTalents()
}
