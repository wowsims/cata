package core

import (
	"github.com/wowsims/mop/sim/core/stats"
)

// Could be in constants.go, but they won't be used anywhere else
// C values are divided by 100 so that we are working with 1% = 0.01
// TODO: UPDATE FOR MOP
// Reference for Cata values: https://web.archive.org/web/20130127084642/http://elitistjerks.com/f15/t29453-combat_ratings_level_85_cataclysm/
const Diminish_k_Druid = 0.972
const Diminish_k_Nondruid = 0.956
const Diminish_Cd_Druid = 116.890707 / 100
const Diminish_Cd_Nondruid = 65.631440 / 100
const Diminish_Cp = 65.631440 / 100
const Diminish_kCd_Druid = (Diminish_k_Druid * Diminish_Cd_Druid)
const Diminish_kCd_Nondruid = (Diminish_k_Nondruid * Diminish_Cd_Nondruid)
const Diminish_kCp = (Diminish_k_Nondruid * Diminish_Cp)

// Diminishing Returns for tank avoidance
// Non-diminishing sources are added separately in spell outcome funcs

func (unit *Unit) GetDiminishedDodgeChance() float64 {
	// undiminished Dodge % = D
	// diminished Dodge % = (D * Cd)/((k*Cd) + D)
	dodgeChance := unit.stats[stats.DodgeRating] / DodgeRatingPerDodgePercent / 100

	if unit.PseudoStats.CanParry {
		return (dodgeChance * Diminish_Cd_Nondruid) / (Diminish_kCd_Nondruid + dodgeChance)
	} else {
		return (dodgeChance * Diminish_Cd_Druid) / (Diminish_kCd_Druid + dodgeChance)
	}
}

func (unit *Unit) GetDiminishedParryChance() float64 {
	// undiminished Parry % = P
	// diminished Parry % = (P * Cp)/((k*Cp) + P)
	parryChance := unit.stats[stats.ParryRating] / ParryRatingPerParryPercent / 100
	return (parryChance * Diminish_Cp) / (Diminish_kCp + parryChance)
}

func (unit *Unit) GetDiminishedMissChance() float64 {
	// Defense Rating is gone in Cata, so there are no diminished sources of Miss
	return 0
}
