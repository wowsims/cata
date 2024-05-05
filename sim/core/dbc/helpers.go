package dbc

// Utility function to check if an item exists in a slice of SpellEffectData pointers
func containsEffect(effects []*SpellEffectData, effect *SpellEffectData) bool {
	for _, e := range effects {
		if e.ID == effect.ID {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func effectCategorySubtypes() []uint {
	return []uint{
		A_MODIFY_CATEGORY_COOLDOWN,
		A_MOD_MAX_CHARGES,
		A_MOD_RECHARGE_TIME,
		A_MOD_RECHARGE_MULTIPLIER,
		A_HASTED_CATEGORY,
		A_MOD_RECHARGE_RATE_CATEGORY,
	}
}

func contains(subtype uint) bool {
	for _, s := range effectCategorySubtypes() {
		if s == subtype {
			return true
		}
	}
	return false
}
