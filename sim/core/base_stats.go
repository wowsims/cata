package core

import (
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked level 70 of the race/class you want, ideally without any talents to mess up base stats.
//  Basic stats are as-shown (str/agi/stm/int/spirit)

// Base Spell Crit is calculated by
//   1. Take as-shown value (troll shaman have 3.5%)
//   2. Calculate the bonus from int (for troll shaman that would be 104/78.1=1.331% crit)
//   3. Subtract as-shown from int bouns (3.5-1.331=2.169)
//   4. 2.169*22.08 (rating per crit percent) = 47.89 crit rating.

// Base mana can be looked up here: https://wowwiki-archive.fandom.com/wiki/Base_mana

// These are also scattered in various dbc/casc files,
// `octbasempbyclass.txt`, `combatratings.txt`, `chancetospellcritbase.txt`, etc.

var RaceOffsets = map[proto.Race]stats.Stats{
	proto.Race_RaceUnknown: stats.Stats{},
	proto.Race_RaceHuman:   stats.Stats{},
	proto.Race_RaceOrc: {
		stats.Agility:   -3,
		stats.Strength:  3,
		stats.Intellect: -3,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   1,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   4,
		stats.Strength:  -4,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   0,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -4,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   2,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   0,
	},
	proto.Race_RaceBloodElf: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceDraenei: {
		stats.Agility:   -3,
		stats.Strength:  1,
		stats.Intellect: 0,
		stats.Spirit:    2,
		stats.Stamina:   0,
	},
	proto.Race_RaceGoblin: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceWorgen: {
		stats.Agility:   2,
		stats.Strength:  3,
		stats.Intellect: -4,
		stats.Spirit:    -1,
		stats.Stamina:   0,
	},
	proto.Race_RaceAlliancePandaren: {
		stats.Agility:   -2,
		stats.Strength:  0,
		stats.Intellect: -1,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceHordePandaren: {
		stats.Agility:   -2,
		stats.Strength:  0,
		stats.Intellect: -1,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
}

var ClassBaseStats = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown: {},
	proto.Class_ClassWarrior: {
		stats.Health:      146663,
		stats.Agility:     133,
		stats.Strength:    206,
		stats.Intellect:   39,
		stats.Spirit:      67,
		stats.Stamina:     188,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassPaladin: {
		stats.Health:      146663,
		stats.Agility:     105,
		stats.Strength:    178,
		stats.Intellect:   114,
		stats.Spirit:      123,
		stats.Stamina:     169,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassHunter: {
		stats.Health:            146663,
		stats.Agility:           216,
		stats.Strength:          86,
		stats.Intellect:         105,
		stats.Spirit:            113,
		stats.Stamina:           151,
		stats.AttackPower:       float64(CharacterLevel)*2.0 - 20,
		stats.RangedAttackPower: float64(CharacterLevel)*2.0 - 10,
	},
	proto.Class_ClassRogue: {
		stats.Health:      146663,
		stats.Agility:     225,
		stats.Strength:    132,
		stats.Intellect:   48,
		stats.Spirit:      77,
		stats.Stamina:     123,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 20,
	},
	proto.Class_ClassPriest: {
		stats.Health:    146663,
		stats.Agility:   58,
		stats.Strength:  48,
		stats.Intellect: 207,
		stats.Spirit:    216,
		stats.Stamina:   77,
	},
	proto.Class_ClassDeathKnight: {
		stats.Health:      146663,
		stats.Agility:     131,
		stats.Strength:    209,
		stats.Intellect:   38,
		stats.Spirit:      69,
		stats.Stamina:     190,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassShaman: {
		stats.Health:      146663,
		stats.Agility:     86,
		stats.Strength:    142,
		stats.Intellect:   151,
		stats.Spirit:      169,
		stats.Stamina:     161,
		stats.AttackPower: float64(CharacterLevel) * 2.0,
	},
	proto.Class_ClassMage: {
		stats.Health:    146663,
		stats.Agility:   48,
		stats.Strength:  39,
		stats.Intellect: 215,
		stats.Spirit:    207,
		stats.Stamina:   67,
	},
	proto.Class_ClassWarlock: {
		stats.Health:      146663,
		stats.Agility:     77,
		stats.Strength:    67,
		stats.Intellect:   188,
		stats.Spirit:      198,
		stats.Stamina:     104,
		stats.AttackPower: -10,
	},
	proto.Class_ClassDruid: {
		stats.Health:      146663,
		stats.Agility:     95,
		stats.Strength:    104,
		stats.Intellect:   169,
		stats.Spirit:      188,
		stats.Stamina:     114,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 10,
	},
	proto.Class_ClassMonk: {
		stats.Health:      146663,
		stats.Agility:     113,
		stats.Strength:    94,
		stats.Intellect:   169,
		stats.Spirit:      190,
		stats.Stamina:     113,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 30,
	},
}

var ClassBaseScaling = map[proto.Class]float64{
	proto.Class_ClassUnknown:     1710.000000,
	proto.Class_ClassWarrior:     1246.298600,
	proto.Class_ClassPaladin:     1141.926000,
	proto.Class_ClassHunter:      1246.298600,
	proto.Class_ClassRogue:       1246.298600,
	proto.Class_ClassPriest:      1049.328400,
	proto.Class_ClassDeathKnight: 1246.298600,
	proto.Class_ClassShaman:      1114.501700,
	proto.Class_ClassMage:        1040.778600,
	proto.Class_ClassWarlock:     1068.202900,
	proto.Class_ClassMonk:        1094.739700,
	proto.Class_ClassDruid:       1094.739700,
}

func AddBaseStatsCombo(r proto.Race, c proto.Class) {
	BaseStats[BaseStatsKey{Race: r, Class: c}] = ClassBaseStats[c].Add(RaceOffsets[r]).Add(ExtraClassBaseStats[c])
}

func init() {
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDruid)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDruid)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassDruid)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassDruid)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassDeathKnight)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassDeathKnight)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassHunter)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassMage)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassPaladin)

	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassPriest)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassRogue)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceAlliancePandaren, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceHordePandaren, proto.Class_ClassShaman)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassWarlock)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassWarrior)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceWorgen, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceAlliancePandaren, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceGoblin, proto.Class_ClassMonk)
	AddBaseStatsCombo(proto.Race_RaceHordePandaren, proto.Class_ClassMonk)
}
