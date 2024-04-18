package talent_trees

import _ "embed"

//go:embed death_knight.json
var DeathKnightTalentsConfig string

//go:embed druid.json
var DruidTalentsConfig string

//go:embed hunter.json
var HunterTalentsConfig string

//go:embed mage.json
var MageTalentsConfig string

//go:embed paladin.json
var PaladinTalentsConfig string

//go:embed priest.json
var PriestTalentsConfig string

//go:embed rogue.json
var RogueTalentsConfig string

//go:embed shaman.json
var ShamanTalentsConfig string

//go:embed warlock.json
var WarlockTalentsConfig string

//go:embed warrior.json
var WarriorTalentsConfig string
