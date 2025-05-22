package core

import (
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type DiminishingReturnsConstants struct {
	k, c_p, c_d, c_b float64
}

// https://github.com/raethkcj/MistsDiminishingReturns
var AvoidanceDRByClass = map[proto.Class]DiminishingReturnsConstants{
	proto.Class_ClassWarrior:     {0.956, 237.186, 90.6425, 150.376},
	proto.Class_ClassPaladin:     {0.886, 237.186, 66.5675, 150.376},
	proto.Class_ClassHunter:      {0.988, 0, 145.560, 0},
	proto.Class_ClassRogue:       {0.988, 145.560, 145.560, 0},
	proto.Class_ClassPriest:      {0.983, 0, 150.376, 0},
	proto.Class_ClassDeathKnight: {0.956, 237.186, 90.6425, 0},
	proto.Class_ClassShaman:      {0.988, 145.560, 145.560, 0},
	proto.Class_ClassMonk:        {1.422, 90.6425, 501.253, 0},
	proto.Class_ClassMage:        {0.983, 0, 150.376, 0},
	proto.Class_ClassWarlock:     {0.983, 0, 150.376, 0},
	proto.Class_ClassDruid:       {1.222, 0, 150.376, 0},
}

// Diminishing Returns for tank avoidance
// Non-diminishing sources are added separately in spell outcome funcs

func (unit *Unit) GetDiminishedDodgeChance() float64 {
	// undiminished Dodge % = D
	// diminished Dodge % = (D * Cd)/((k*Cd) + D)
	dodgeChance := unit.stats[stats.DodgeRating] / DodgeRatingPerDodgePercent / 100
	return (dodgeChance * unit.avoidanceParams.c_d) / (unit.avoidanceParams.k * unit.avoidanceParams.c_d + dodgeChance)
}

func (unit *Unit) GetDiminishedParryChance() float64 {
	// undiminished Parry % = P
	// diminished Parry % = (P * Cp)/((k*Cp) + P)
	parryChance := unit.stats[stats.ParryRating] / ParryRatingPerParryPercent / 100
	return (parryChance * unit.avoidanceParams.c_p) / (unit.avoidanceParams.k * unit.avoidanceParams.c_p + parryChance)
}

func (unit *Unit) GetDiminishedBlockChance() float64 {
	// undiminished Block % = B
	// diminished Block % = (B * Cb)/((k*Cb) + B)
	blockChance := unit.stats[stats.BlockPercent] / 100
	return (blockChance * unit.avoidanceParams.c_b) / (unit.avoidanceParams.k * unit.avoidanceParams.c_b + blockChance)
}
