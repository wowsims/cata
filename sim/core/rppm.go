// Implements the core functionality for the RPPM system
// RPPM values are either dynamically loaded from the spell data
// Or dan be provided manually if the values are not correct in the spell data
package core

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

type rppmMod interface {
	GetCoefficient(proc *RPPMProc) float64

	// Returns true if the mod is static in reference to the character or false; if the mod result might change
	IsStatic() bool
}

type rppmCritMod struct{}

func (r rppmCritMod) GetCoefficient(proc *RPPMProc) float64 {
	return max(1+proc.character.GetStat(stats.PhysicalCritPercent)/100.0,
		1+proc.character.GetStat(stats.SpellCritPercent)/100.0)
}

func (r rppmCritMod) IsStatic() bool {
	return false
}

type rppmHasteMod struct{}

func (r rppmHasteMod) GetCoefficient(proc *RPPMProc) float64 {
	// as of 5.2 this should no longer include non 'True haste mods' so only i.E. Lust
	return max(proc.character.TotalRealHasteMultiplier(), proc.character.TotalSpellHasteMultiplier())
}

func (r rppmHasteMod) IsStatic() bool {
	return false
}

type rppmSpecMod struct {
	spec        proto.Spec
	coefficient float64
}

func (r rppmSpecMod) GetCoefficient(proc *RPPMProc) float64 {
	if proc.character.Spec == r.spec {
		return 1 + r.coefficient
	}

	return 1.0
}

func (r rppmSpecMod) IsStatic() bool {
	return true
}

type rppmClassMod struct {
	classMask   int
	coefficient float64
}

func (r rppmClassMod) GetCoefficient(proc *RPPMProc) float64 {
	mask := 1 << (proc.character.Class - 1)
	if r.classMask&mask > 0 {
		return 1 + r.coefficient
	}

	return 1
}

func (r rppmClassMod) IsStatic() bool {
	return true
}

type rppmApproxIlvlMod struct {
	baseIlvl    int32
	coefficient float64
}

func (r rppmApproxIlvlMod) GetCoefficient(proc *RPPMProc) float64 {
	// We use an approximation here, or we'd need to load the complete random properties table into the sim
	// Just to calculate the difference in random prop points as not all points are available on the item that scales
	// The maximal relative error I observed when comparing real values and approximation was 0.07%
	// which was an increase of 0.0008% in proc chance. So for now I think we can neglect that

	if r.baseIlvl == proc.ilvl {
		return 1
	}

	// Each ilvl step is ~0.936% larger than the previous. They're rounded in the table
	// But over a range of >3 ilvl this approximation becomes very accurate
	return 1 + (math.Pow(1.00936, float64(proc.ilvl-r.baseIlvl))-1)*r.coefficient
}

func (r rppmApproxIlvlMod) IsStatic() bool {
	return true
}

type RPPMConfig struct {
	PPM         float64
	Coefficient float64
	Ilvl        int32
	Mods        []rppmMod
}

type RPPMProc struct {
	ppm         float64
	coefficient float64
	character   *Character
	lastProc    time.Duration
	lastCheck   time.Duration
	ilvl        int32
	mods        []rppmMod
}

// Attach a crit mod to the RPPM config
func (config RPPMConfig) WithCritMod() RPPMConfig {
	config.Mods = append(config.Mods, rppmCritMod{})

	return config
}

// Attach a haste mod to the RPPM config
// It uses the highest haste value that does not include effects like Slice and Dice
// It multiplies the actual proc chance by 1 + haste%
func (config RPPMConfig) WithHasteMod() RPPMConfig {
	config.Mods = append(config.Mods, rppmHasteMod{})

	return config
}

// Attach a class specific modifier to the RPPM config
// 1 - Warrior, 2 - Paladin, 4 - Hunter, 8 - Rogue, 16 - Priest, 32 - DK
// 64 - Shaman, 128 - Mage, 256 - Warlock, 512 - Monk, 1024 - Druid
// It multiplies the actual proc chance by 1 + coefficient
func (config RPPMConfig) WithClassMod(coefficient float64, classMask int) RPPMConfig {
	config.Mods = append(config.Mods, rppmClassMod{
		classMask:   classMask,
		coefficient: coefficient,
	})

	return config
}

// Attaches a spec mod to the RPPM config
// It multiplies the actual proc chance by 1 + coefficient
func (config RPPMConfig) WithSpecMod(coefficient float64, spec proto.Spec) RPPMConfig {
	config.Mods = append(config.Mods, rppmSpecMod{
		spec:        spec,
		coefficient: coefficient,
	})

	return config
}

// Attach an approximate Ilvl scaling to the RPPM config
// The proc chance will be multiplied by 1.00936^(ilvlDiff)
func (config RPPMConfig) WithApproximateIlvlMod(coefficient float64, baseIlvl int32) RPPMConfig {
	config.Mods = append(config.Mods, rppmApproxIlvlMod{
		coefficient: coefficient,
		baseIlvl:    baseIlvl,
	})

	return config
}

// Create a new RPPM Proc with the given ppm (usually from the ProcsPerMinute record)
func NewRPPMProc(character *Character, config RPPMConfig) DynamicProc {
	proc := &RPPMProc{
		character:   character,
		ppm:         config.PPM,
		coefficient: TernaryFloat64(config.Coefficient > 0, config.Coefficient, 1),
		ilvl:        config.Ilvl,
		lastProc:    -time.Second * 120,
		lastCheck:   -time.Second * 10,
		mods:        []rppmMod{},
	}

	if config.Mods != nil {
		for _, mod := range config.Mods {
			if mod.IsStatic() {
				proc.coefficient *= mod.GetCoefficient(proc)
			} else {
				proc.mods = append(proc.mods, mod)
			}
		}
	}

	return proc
}

// Does not change the state of the RPPMProc.
// Only calculates the proc chance.
//
// To actually modify the state correctly call:
//
//	Proc(sim, string)
func (proc *RPPMProc) getProcChance(sim *Simulation) float64 {
	basePpm := proc.ppm

	baseCoeff := proc.coefficient
	for _, mod := range proc.mods {
		baseCoeff *= mod.GetCoefficient(proc)
	}

	lastCheck := math.Min(10.0, (sim.CurrentTime - proc.lastCheck).Seconds())
	lastProc := math.Min(1000.0, (sim.CurrentTime - proc.lastProc).Seconds())

	// TODO: Adjust implementation if needed
	// Temporary implementation, targeting the 'intended' MOP proc behavior
	// https://github.com/ClassicWoWCommunity/cata-classic-bugs/issues/1774

	realPPM := basePpm * baseCoeff
	procCoefficient := math.Max(1, 1+((lastProc/(60/realPPM))-1.5)*3) // Bad luck protection
	baseProcChance := realPPM * (lastCheck / 60.0)
	return baseProcChance * procCoefficient
}

func (proc *RPPMProc) Proc(sim *Simulation, label string) bool {
	result := sim.Proc(proc.getProcChance(sim), label)
	proc.lastCheck = sim.CurrentTime
	if result {
		proc.lastProc = sim.CurrentTime
	}

	return result
}

func (proc *RPPMProc) Chance(sim *Simulation) float64 {
	return proc.getProcChance(sim)
}

func (proc *RPPMProc) Reset() {
	proc.lastCheck = time.Second * -10
	proc.lastProc = time.Second * -120
}

func RppmModFromProto(config *proto.RppmMod) (rppmMod, error) {
	switch modType := config.GetModType().(type) {
	case *proto.RppmMod_ClassMask:
		return rppmClassMod{
			classMask:   int(config.GetClassMask()),
			coefficient: config.GetCoefficient(),
		}, nil
	case *proto.RppmMod_Crit:
		return rppmCritMod{}, nil
	case *proto.RppmMod_Haste:
		return rppmHasteMod{}, nil
	case *proto.RppmMod_Spec:
		return rppmSpecMod{
			coefficient: config.GetCoefficient(),
			spec:        config.GetSpec(),
		}, nil
	case *proto.RppmMod_Ilvl:
		return rppmApproxIlvlMod{
			coefficient: config.GetCoefficient(),
			baseIlvl:    config.GetIlvl(),
		}, nil
	case nil:
		return nil, fmt.Errorf("rppmMod: ModType is not set")
	default:
		return nil, fmt.Errorf("unknown ModType: %T", modType)
	}
}

func RppmConfigFromProcEffectProto(effect *proto.ProcEffect) RPPMConfig {
	config := RPPMConfig{
		PPM: effect.GetRppm().GetRate(),
	}

	for _, protoMod := range effect.GetRppm().GetMods() {
		mod, error := RppmModFromProto(protoMod)
		if error != nil {
			panic("Could not parse rrpm mod from proto")
		}

		config.Mods = append(config.Mods, mod)
	}

	return config
}
