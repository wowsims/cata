package core

func (result *SpellResult) applyResistances(sim *Simulation, spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	resistanceMultiplier := spell.resistanceMultiplier(sim, isPeriodic, attackTable)

	result.Damage *= resistanceMultiplier

	result.ArmorMultiplier = resistanceMultiplier
	result.PreOutcomeDamage = result.Damage
}

// Modifies damage based on Armor
func (spell *Spell) resistanceMultiplier(sim *Simulation, isPeriodic bool, attackTable *AttackTable) float64 {
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

	// Physical resistance (armor).
	return attackTable.getArmorDamageModifier(spell)
}

// https://web.archive.org/web/20130511200023/http://elitistjerks.com/f15/t29453-combat_ratings_level_85_cataclysm/p40/#post2171306
func (at *AttackTable) getArmorDamageModifier(spell *Spell) float64 {
	if at.IgnoreArmor {
		return 1.0
	}

	ignoreArmorFactor := Clamp(at.ArmorIgnoreFactor, 0.0, 1.0)

	// Assume target > 80
	armorConstant := float64(at.Attacker.Level)*4037.5 - 317117.5
	defenderArmor := at.Defender.Armor() - (at.Defender.Armor() * ignoreArmorFactor)
	return 1 - defenderArmor/(defenderArmor+armorConstant)
}
