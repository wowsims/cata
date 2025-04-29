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
		stats.Health:      43285,
		stats.Agility:     123,
		stats.Strength:    189,
		stats.Intellect:   37,
		stats.Spirit:      63,
		stats.Stamina:     173,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassPaladin: {
		stats.Health:      43285,
		stats.Agility:     97,
		stats.Strength:    164,
		stats.Intellect:   106,
		stats.Spirit:      117,
		stats.Stamina:     156,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassHunter: {
		stats.Health:            39037,
		stats.Agility:           198,
		stats.Strength:          80,
		stats.Intellect:         105,
		stats.Spirit:            97,
		stats.Stamina:           138,
		stats.AttackPower:       float64(CharacterLevel)*2.0 - 20,
		stats.RangedAttackPower: float64(CharacterLevel)*2.0 - 10,
	},
	proto.Class_ClassRogue: {
		stats.Health:      40529,
		stats.Agility:     206,
		stats.Strength:    122,
		stats.Intellect:   46,
		stats.Spirit:      71,
		stats.Stamina:     114,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 20,
	},
	proto.Class_ClassPriest: {
		stats.Health:    43285,
		stats.Agility:   54,
		stats.Strength:  46,
		stats.Intellect: 190,
		stats.Spirit:    198,
		stats.Stamina:   71,
	},
	proto.Class_ClassDeathKnight: {
		stats.Health:      43285,
		stats.Agility:     121,
		stats.Strength:    191,
		stats.Intellect:   36,
		stats.Spirit:      63,
		stats.Stamina:     174,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 20,
	},
	proto.Class_ClassShaman: {
		stats.Health:      37097,
		stats.Agility:     80,
		stats.Strength:    131,
		stats.Intellect:   139,
		stats.Spirit:      156,
		stats.Stamina:     158,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 30,
	},
	proto.Class_ClassMage: {
		stats.Health:    37113,
		stats.Agility:   46,
		stats.Strength:  37,
		stats.Intellect: 198,
		stats.Spirit:    190,
		stats.Stamina:   63,
	},
	proto.Class_ClassWarlock: {
		stats.Health:      38184,
		stats.Agility:     71,
		stats.Strength:    63,
		stats.Intellect:   173,
		stats.Spirit:      161,
		stats.Stamina:     96,
		stats.AttackPower: -10,
	},
	proto.Class_ClassDruid: {
		stats.Health:      39533,
		stats.Agility:     89,
		stats.Strength:    96,
		stats.Intellect:   156,
		stats.Spirit:      173,
		stats.Stamina:     106,
		stats.AttackPower: float64(CharacterLevel)*3.0 - 10,
	},
	proto.Class_ClassMonk: {
		stats.Health:      43285,
		stats.Agility:     85,
		stats.Strength:    68,
		stats.Intellect:   135,
		stats.Spirit:      154,
		stats.Stamina:     86,
		stats.AttackPower: float64(CharacterLevel)*2.0 - 20,
	},
}

var ClassBaseScaling = map[proto.Class]float64{
	proto.Class_ClassWarrior:     1125.227400,
	proto.Class_ClassPaladin:     1029.493400,
	proto.Class_ClassHunter:      1125.227400,
	proto.Class_ClassRogue:       1125.227400,
	proto.Class_ClassPriest:      945.188840,
	proto.Class_ClassDeathKnight: 1125.227400,
	proto.Class_ClassShaman:      1004.487900,
	proto.Class_ClassMage:        937.330080,
	proto.Class_ClassWarlock:     962.335630,
	proto.Class_ClassMonk:        986.626400,
	proto.Class_ClassDruid:       986.626400,
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
