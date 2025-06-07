// Implements the core functionality for the RPPM system
// RPPM values are either dynamically loaded from the spell data
// Or dan be provided manually if the values are not correct in the spell data
package core

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type rppmMod interface {
	GetCoefficient(unit *Unit) float64
}

type rppmCharMod interface {
	GetCoefficient(character *Character) float64
}

type rppmCritMod struct {
	coefficient float64
	kind        CritRPPMModKind
}

func (r rppmCritMod) GetCoefficient(unit *Unit) float64 {
	switch r.kind {
	case MeleeCrit:
		// We do not separate those two crit values apparentls
		// Treat them as the same
		fallthrough
	case RangedCrit:
		return 1 + unit.GetStat(stats.PhysicalCritPercent)/100.0
	case SpellCrit:
		return 1 + unit.GetStat(stats.SpellCritPercent)/100.0
	case LowestCrit:
		return math.Min(1+unit.GetStat(stats.PhysicalCritPercent)/100.0,
			1+unit.GetStat(stats.SpellCritPercent)/100.0)
	default:
		return 1.0
	}
}

type rppmHasteMod struct {
	coefficient float64
	kind        HasteRPPMModKind
}

// GetCoefficient implements RPPMMod.
func (r rppmHasteMod) GetCoefficient(unit *Unit) float64 {
	switch r.kind {
	case MeleeHaste:
		return unit.SwingSpeed()
	case RangedHaste:
		return unit.RangedSwingSpeed()
	case SpellHaste:
		return 1 / unit.CastSpeed
	case HighestHaste:
		// as of 5.2 this should no longer include non 'True haste mods' so only i.E. Lust
		// With that this all should be equal
		return unit.PseudoStats.AttackSpeedMultiplier * (1 + unit.GetStat(stats.HasteRating)/(HasteRatingPerHastePercent*100))
	case LowestHaste:
		return math.Min(unit.CastSpeed, math.Min(unit.SwingSpeed(), unit.RangedSwingSpeed()))
	default:
		return 1.0
	}
}

type rppmSpecMod struct {
	spec        proto.Spec
	coefficient float64
}

func (r rppmSpecMod) GetCoefficient(charater *Character) float64 {
	if charater.Spec == r.spec {
		return 1 + r.coefficient
	}

	return 1.0
}

type rppmClassMod struct {
	classMask   int
	coefficient float64
}

func toDBCClass(class proto.Class) int {
	switch class {
	case proto.Class_ClassWarrior:
		return 1
	case proto.Class_ClassPaladin:
		return 2
	case proto.Class_ClassHunter:
		return 3
	case proto.Class_ClassRogue:
		return 4
	case proto.Class_ClassPriest:
		return 5
	case proto.Class_ClassDeathKnight:
		return 6
	case proto.Class_ClassShaman:
		return 7
	case proto.Class_ClassMage:
		return 8
	case proto.Class_ClassWarlock:
		return 9
	case proto.Class_ClassMonk:
		return 10
	case proto.Class_ClassDruid:
		return 11
	default:
		return 1
	}
}

func (r rppmClassMod) GetCoefficient(character *Character) float64 {
	mask := 1 << (toDBCClass(character.Class) - 1)
	if r.classMask&mask > 0 {
		return 1 + r.coefficient
	}

	return 1
}

type rppmApproxIlvlMod struct {
	baseIlvl    int32
	realIlvl    int32
	coefficient float64
}

func (r rppmApproxIlvlMod) GetCoefficient(unit *Unit) float64 {

	// We use an approximation here, or we'd need to load the complete random properties table into the sim
	// Just to calculate the difference in random prop points as not all points are available on the item that scales
	// The maximal relativ error I observed when comparing real values and approximation was 0.07%
	// which was an increase of 0,0008% in proc chance. So for now I think we can neglect that

	if r.baseIlvl == r.realIlvl {
		return 1
	}

	// Each ilvl step is ~0.936% larger than the previous. They're rounded in the table
	// But over a range of >3 ilvl this approximation becomes very accurate
	return 1 + (math.Pow(1.00936, float64(r.realIlvl-r.baseIlvl))-1)*r.coefficient
}

type HasteRPPMModKind uint8

const (
	MeleeHaste HasteRPPMModKind = iota
	RangedHaste
	SpellHaste
	HighestHaste
	LowestHaste
)

type CritRPPMModKind uint8

const (
	MeleeCrit CritRPPMModKind = iota
	RangedCrit
	SpellCrit
	LowestCrit
)

type RPPMProc struct {
	ppm             float64
	coefficient     float64
	charCoefficient float64
	char            *Character
	lastProc        time.Duration
	lastCheck       time.Duration
	mod             []rppmMod
	charMods        []rppmCharMod
}

// Attach a crit mot to the RPPM Proc
// The most common value used is LowestCrit
func (proc *RPPMProc) WithCritMod(coefficient float64, kind CritRPPMModKind) *RPPMProc {
	proc.mod = append(proc.mod, rppmCritMod{
		kind:        kind,
		coefficient: coefficient,
	})

	return proc
}

// Attach a haste mod to the RPPM Proc
// The most common used kind is HighestHaste
// It uses the highest haste value that does not include effects like Slice and Dice
// It multiplies the actual proc chance by 1 + haste%
func (proc *RPPMProc) WithHasteMod(coeffienct float64, kind HasteRPPMModKind) *RPPMProc {
	proc.mod = append(proc.mod, rppmHasteMod{
		kind:        kind,
		coefficient: coeffienct,
	})

	return proc
}

// Attach a class specific modifier to the RPPM
// 1 - Warrior, 2 - Paladin, 4 - Hunter, 8 - Rogue, 16 - Priest, 32 - DK
// 64 - Shaman, 128 - Mage, 256 - Warlock, 512 - Monk, 1024 - Druid
// It multiplies the actual proc chance by 1 + coefficient
func (proc *RPPMProc) WithClassMod(coefficient float64, classMask int) *RPPMProc {
	mod := rppmClassMod{
		classMask:   classMask,
		coefficient: coefficient,
	}

	if proc.char != nil {
		proc.charCoefficient *= mod.GetCoefficient(proc.char)
	} else {
		proc.charMods = append(proc.charMods, mod)
	}

	return proc
}

// Attaches a spec mod to the RPPM
// It multiplies the actual proc chance by 1 + coefficient
func (proc *RPPMProc) WithSpecMod(coefficient float64, spec proto.Spec) *RPPMProc {
	mod := rppmSpecMod{
		spec:        spec,
		coefficient: coefficient,
	}

	if proc.char != nil {
		proc.charCoefficient *= mod.GetCoefficient(proc.char)
	} else {
		proc.charMods = append(proc.charMods, mod)
	}

	return proc
}

// Set the base coefficient for the RPPM proc chance
// It's a multiplicator for the RPPM based proc chance
func (proc *RPPMProc) WithCoefficient(coefficient float64) *RPPMProc {
	proc.coefficient = coefficient
	return proc
}

// Attach an approximate Ilvl scaling to the Mod
// The proc chance will be multiplied by 1.00936^(ilvlDiff)
func (proc *RPPMProc) WithApproximateIlvlMod(coefficient float64, baseIlvl int32, realIlvl int32) *RPPMProc {
	proc.mod = append(proc.mod, rppmApproxIlvlMod{
		coefficient: coefficient,
		baseIlvl:    int32(baseIlvl),
		realIlvl:    realIlvl,
	})

	return proc
}

func (proc *RPPMProc) ForCharacter(character *Character) *RPPMProc {
	proc.char = character

	for _, mod := range proc.charMods {
		proc.charCoefficient *= mod.GetCoefficient(character)
	}

	proc.charMods = []rppmCharMod{}
	return proc
}

// Create a new RPPM Proc with the given ppm (usually from the ProcsPerMinute record)
func NewRPPMProc(ppm float64) *RPPMProc {
	proc := &RPPMProc{
		ppm:             ppm,
		coefficient:     1,
		charCoefficient: 1,
		lastProc:        -time.Second * 120,
		lastCheck:       -time.Second * 10,
		mod:             []rppmMod{},
	}

	return proc
}

// Does not change the state of the RPPMProc.
// Only calculates the proc chance.
//
// To actually modify the state correctly call:
//
//	Proc(character, sim)
func (proc *RPPMProc) getProcChance(unit *Unit, sim *Simulation) float64 {
	basePpm := proc.ppm

	baseCoeff := proc.coefficient
	for _, mod := range proc.mod {
		baseCoeff *= mod.GetCoefficient(unit)
	}

	lastCheck := math.Min(10.0, float64((sim.CurrentTime-proc.lastCheck)/time.Second))
	lastProc := math.Min(1000.0, float64((sim.CurrentTime-proc.lastProc)/time.Second))

	procCoefficient := math.Max(1, 1+((lastProc/(60/proc.ppm))-1.5)*3) // Bad luck protection
	baseProcChance := (basePpm * (lastCheck / 60.0)) * baseCoeff
	return baseProcChance * procCoefficient * proc.charCoefficient
}

func (proc *RPPMProc) Proc(unit *Unit, sim *Simulation, label string) bool {
	result := sim.Proc(proc.getProcChance(unit, sim), label)
	proc.lastCheck = sim.CurrentTime
	if result {
		proc.lastProc = sim.CurrentTime
	}

	return result
}

func (proc *RPPMProc) Chance(unit *Unit, sim *Simulation) float64 {
	return proc.getProcChance(unit, sim)
}

func (proc *RPPMProc) Reset() {
	proc.lastCheck = time.Second * -10
	proc.lastProc = time.Second * -120
}
