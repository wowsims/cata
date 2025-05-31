package core

func (result *SpellResult) applyArmor(spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	armorMitigationMultiplier := spell.armorMultiplier(isPeriodic, attackTable)

	result.Damage *= armorMitigationMultiplier

	result.ArmorMultiplier = armorMitigationMultiplier
	result.PreOutcomeDamage = result.Damage
}

// Returns Armor mitigation fraction for the spell
func (spell *Spell) armorMultiplier(isPeriodic bool, attackTable *AttackTable) float64 {
	if spell.Flags.Matches(SpellFlagIgnoreArmor) {
		return 1
	}

	// There are no (partial) resists in MoP
	if !spell.SpellSchool.Matches(SpellSchoolPhysical) {
		return 1
	}

	// All physical dots (Bleeds) ignore armor.
	if isPeriodic && !spell.Flags.Matches(SpellFlagApplyArmorReduction) {
		return 1
	}

	// return armor mitigation fraction
	return attackTable.getArmorDamageModifier()
}

// https://web.archive.org/web/20130511200023/http://elitistjerks.com/f15/t29453-combat_ratings_level_85_cataclysm/p40/#post2171306
func (at *AttackTable) getArmorDamageModifier() float64 {
	if at.IgnoreArmor {
		return 1.0
	}

	ignoreArmorFactor := Clamp(at.ArmorIgnoreFactor, 0.0, 1.0)

	// Assume target > 80
	armorConstant := float64(at.Attacker.Level)*4037.5 - 317117.5
	defenderArmor := at.Defender.Armor() * (1.0 - ignoreArmorFactor)
	return 1 - defenderArmor/(defenderArmor+armorConstant)
}
