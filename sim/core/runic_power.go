package core

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type RuneChangeType int32

const (
	None             RuneChangeType = 0
	SpendRune        RuneChangeType = 1
	GainRune         RuneChangeType = 2
	ConvertToDeath   RuneChangeType = 4
	ConvertFromDeath RuneChangeType = 8
)

func (r RuneChangeType) Matches(other RuneChangeType) bool {
	return (r & other) != 0
}

type OnRuneChange func(sim *Simulation, changeType RuneChangeType, runeRegen []int8)
type OnRunicPowerGain func(sim *Simulation)

type RuneMeta struct {
	regenMulti        float64
	regenAt           time.Duration // time at which the rune will no longer be spent.
	unscaledRegenLeft time.Duration // time which the rune spent in regen (Unscaled)

	revertAt time.Duration // time at which rune will no longer be kind death.
}

type runicPowerBar struct {
	unit *Unit

	maxRunicPower     float64
	currentRunicPower float64
	runeCD            time.Duration

	// These flags are used to simplify pending action checks
	// |DS|DS|DS|DS|DS|DS|
	runeStates int16
	runeMeta   [6]RuneMeta
	btSlot     int8

	bloodRuneGainMetrics  *ResourceMetrics
	frostRuneGainMetrics  *ResourceMetrics
	unholyRuneGainMetrics *ResourceMetrics
	deathRuneGainMetrics  *ResourceMetrics

	onRuneChange     OnRuneChange
	onRunicPowerGain OnRunicPowerGain

	pa *PendingAction

	runeRegenMultiplier  float64
	runicRegenMultiplier float64

	permanentDeaths []int8

	lastRegen []int8
}

// Constants for finding runes
// |DS|DS|DS|DS|DS|DS|
const (
	baseRuneState = 0 // unspent, no death

	allDeath = 0b101010101010
	allSpent = 0b010101010101

	anyBloodSpent  = 0b0101 << 0
	anyFrostSpent  = 0b0101 << 4
	anyUnholySpent = 0b0101 << 8
)

var (
	isDeaths     = [6]int16{0b10 << 0, 0b10 << 2, 0b10 << 4, 0b10 << 6, 0b10 << 8, 0b10 << 10}
	isSpents     = [6]int16{0b01 << 0, 0b01 << 2, 0b01 << 4, 0b01 << 6, 0b01 << 8, 0b01 << 10}
	isSpentDeath = [6]int16{0b11 << 0, 0b11 << 2, 0b11 << 4, 0b11 << 6, 0b11 << 8, 0b11 << 10}
)

func (rp *runicPowerBar) DebugString() string {
	ss := make([]string, len(rp.runeMeta))
	for i := range rp.runeMeta {
		ss[i] += fmt.Sprintf("Rune %d - D: %v S: %v\n\tRegenAt: %0.1f, RevertAt: %0.1f", i, rp.runeStates&isDeaths[i] != 0, rp.runeStates&isSpents[i] != 0, rp.runeMeta[i].regenAt.Seconds(), rp.runeMeta[i].revertAt.Seconds())
	}
	return strings.Join(ss, "\n")
}

func (rp *runicPowerBar) reset(sim *Simulation) {
	if rp.unit == nil {
		return
	}

	if rp.pa != nil {
		rp.pa.Cancel(sim)
	}

	for i := range rp.runeMeta {
		rp.runeMeta[i].regenAt = NeverExpires

		rp.runeMeta[i].revertAt = NeverExpires
	}

	rp.runeStates = baseRuneState
	for i := range rp.permanentDeaths {
		rp.runeStates |= isDeaths[i]
	}
}

func (unit *Unit) EnableRunicPowerBar(currentRunicPower float64, maxRunicPower float64, runeCD time.Duration,
	onRuneChange OnRuneChange, onRunicPowerGain OnRunicPowerGain) {
	unit.SetCurrentPowerBar(RunicPower)
	unit.runicPowerBar = runicPowerBar{
		unit: unit,

		maxRunicPower:        maxRunicPower,
		currentRunicPower:    currentRunicPower,
		runeCD:               runeCD,
		runeRegenMultiplier:  1.0,
		runicRegenMultiplier: 1.0,

		runeStates: baseRuneState,
		btSlot:     -1,

		onRuneChange:     onRuneChange,
		onRunicPowerGain: onRunicPowerGain,

		permanentDeaths: make([]int8, 0),
		lastRegen:       make([]int8, 0),
	}

	unit.bloodRuneGainMetrics = unit.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	unit.frostRuneGainMetrics = unit.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	unit.unholyRuneGainMetrics = unit.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	unit.deathRuneGainMetrics = unit.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})
}

func (unit *Unit) HasRunicPowerBar() bool {
	return unit.runicPowerBar.unit != nil
}

func (rp *runicPowerBar) SetPermanentDeathRunes(permanentDeaths []int8) {
	rp.permanentDeaths = permanentDeaths
}

func (rp *runicPowerBar) SetRuneCd(runeCd time.Duration) {
	rp.runeCD = runeCd
}

func (rp *runicPowerBar) CurrentRunicPower() float64 {
	return rp.currentRunicPower
}

func (rp *runicPowerBar) maybeFireChange(sim *Simulation, changeType RuneChangeType) {
	if changeType != None && rp.onRuneChange != nil {
		rp.onRuneChange(sim, changeType, rp.lastRegen)
		// Clear regen runes
		rp.lastRegen = make([]int8, 0)
	}
}

func (rp *runicPowerBar) addRunicPowerInterval(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	newRunicPower := min(rp.currentRunicPower+(amount*rp.runicRegenMultiplier), rp.maxRunicPower)

	metrics.AddEvent(amount, newRunicPower-rp.currentRunicPower)

	if sim.Log != nil {
		rp.unit.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower, rp.maxRunicPower)
	}

	rp.currentRunicPower = newRunicPower
}

func (rp *runicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInterval(sim, amount, metrics)
	if rp.onRunicPowerGain != nil {
		rp.onRunicPowerGain(sim)
	}
}

func (rp *runicPowerBar) spendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	newRunicPower := rp.currentRunicPower - amount

	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		rp.unit.Log(sim, "Spent %0.3f runic power from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower, rp.maxRunicPower)
	}

	rp.currentRunicPower = newRunicPower
}

func (rp *runicPowerBar) SpendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.spendRunicPower(sim, amount, metrics)
}

// DeathRuneRegenAt returns the time the given death rune will regen at.
// If the rune is not death or not spent it returns NeverExpires.
func (rp *runicPowerBar) DeathRuneRegenAt(slot int32) time.Duration {
	// If not death or not spent, no regen time
	if rp.runeStates&isSpentDeath[slot] != isSpentDeath[slot] {
		return NeverExpires
	}
	return rp.runeMeta[slot].regenAt
}

// DeathRuneRevertAt returns the next time that a death rune will revert.
// If there is no death rune that needs to revert it returns NeverExpires.
func (rp *runicPowerBar) DeathRuneRevertAt() time.Duration {
	minRevertAt := rp.runeMeta[0].revertAt
	for _, rm := range rp.runeMeta[1:] {
		if rm.revertAt < minRevertAt {
			minRevertAt = rm.revertAt
		}
	}
	return minRevertAt
}

func (rp *runicPowerBar) normalSpentRuneReadyAt(slot int8) time.Duration {
	readyAt := NeverExpires
	if t := rp.runeMeta[slot].regenAt; t < readyAt && rp.runeStates&isSpentDeath[slot] == isSpents[slot] {
		readyAt = t
	}
	if t := rp.runeMeta[slot+1].regenAt; t < readyAt && rp.runeStates&isSpentDeath[slot+1] == isSpents[slot+1] {
		readyAt = t
	}
	return readyAt
}

// NormalSpentBloodRuneReadyAt returns the earliest time a spent non-death blood rune is ready.
func (rp *runicPowerBar) NormalSpentBloodRuneReadyAt(_ *Simulation) time.Duration {
	return rp.normalSpentRuneReadyAt(0)
}

func (rp *runicPowerBar) normalRuneReadyAt(sim *Simulation, slot int8) time.Duration {
	if rp.runeStates&isSpentDeath[slot] == 0 || rp.runeStates&isSpentDeath[slot+1] == 0 {
		return sim.CurrentTime
	}
	return rp.normalSpentRuneReadyAt(slot)
}

// NormalFrostRuneReadyAt returns the earliest time a non-death frost rune is ready.
func (rp *runicPowerBar) NormalFrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.normalRuneReadyAt(sim, 2)
}

func (rp *runicPowerBar) NormalUnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.normalRuneReadyAt(sim, 4)
}

func (rp *runicPowerBar) BloodDeathRuneBothReadyAt() time.Duration {
	if rp.runeStates&isDeaths[0] != 0 && rp.runeStates&isDeaths[1] != 0 {
		if max(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt) > 150000000*time.Minute {
			return min(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
		} else {
			return max(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
		}
	} else {
		return -1
	}
}

func (rp *runicPowerBar) RuneReadyAt(sim *Simulation, slot int8) time.Duration {
	if rp.runeStates&isSpents[slot] != isSpents[slot] {
		return sim.CurrentTime
	}
	return rp.runeMeta[slot].regenAt
}

func (rp *runicPowerBar) SpendRuneReadyAt(slot int8, spendAt time.Duration) time.Duration {
	return spendAt + rp.runeCD
}

// BloodRuneReadyAt returns the earliest time a (possibly death-converted) blood rune is ready.
func (rp *runicPowerBar) BloodRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyBloodSpent != anyBloodSpent { // if any are not spent
		return sim.CurrentTime
	}
	return min(rp.runeMeta[0].regenAt, rp.runeMeta[1].regenAt)
}

func (rp *runicPowerBar) FrostRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyFrostSpent != anyFrostSpent { // if any are not spent
		return sim.CurrentTime
	}
	return min(rp.runeMeta[2].regenAt, rp.runeMeta[3].regenAt)
}

func (rp *runicPowerBar) UnholyRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&anyUnholySpent != anyUnholySpent { // if any are not spent
		return sim.CurrentTime
	}
	return min(rp.runeMeta[4].regenAt, rp.runeMeta[5].regenAt)
}

func (rp *runicPowerBar) bothRunesReadyAt(sim *Simulation, slot int8) time.Duration {
	switch (rp.runeStates >> (2 * slot)) & 0b0101 {
	case 0b0000:
		return sim.CurrentTime
	case 0b0001:
		return rp.runeMeta[slot].regenAt
	case 0b0100:
		return rp.runeMeta[slot+1].regenAt
	default:
		return max(rp.runeMeta[slot].regenAt, rp.runeMeta[slot+1].regenAt)
	}
}

func (rp *runicPowerBar) NextBloodRuneReadyAt(sim *Simulation) time.Duration {
	return rp.bothRunesReadyAt(sim, 0)
}

func (rp *runicPowerBar) NextFrostRuneReadyAt(sim *Simulation) time.Duration {
	return rp.bothRunesReadyAt(sim, 2)
}

func (rp *runicPowerBar) NextUnholyRuneReadyAt(sim *Simulation) time.Duration {
	return rp.bothRunesReadyAt(sim, 4)
}

// AnySpentRuneReadyAt returns the next time that a rune will regenerate.
// It will be NeverExpires if there is no rune pending regeneration.
func (rp *runicPowerBar) AnySpentRuneReadyAt() time.Duration {
	minRegenAt := rp.runeMeta[0].regenAt
	for _, rm := range rp.runeMeta[1:] {
		if rm.regenAt < minRegenAt {
			minRegenAt = rm.regenAt
		}
	}
	return minRegenAt
}

func (rp *runicPowerBar) AnyRuneReadyAt(sim *Simulation) time.Duration {
	if rp.runeStates&allSpent != allSpent {
		return sim.CurrentTime
	}
	return rp.AnySpentRuneReadyAt()
}

// ConvertFromDeath reverts the rune to its original type.
func (rp *runicPowerBar) ConvertFromDeath(sim *Simulation, slot int8) {
	if slices.Contains(rp.permanentDeaths, slot) {
		return
	}

	rp.runeStates ^= isDeaths[slot]
	rp.runeMeta[slot].revertAt = NeverExpires

	if rp.runeStates&isSpents[slot] == 0 {
		metrics := rp.bloodRuneGainMetrics
		if slot == 2 || slot == 3 {
			metrics = rp.frostRuneGainMetrics
		} else if slot == 4 || slot == 5 {
			metrics = rp.unholyRuneGainMetrics
		}
		rp.spendRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
		rp.gainRuneMetrics(sim, metrics, 1)
	}
}

// ConvertToDeath converts the given slot to death and sets up the reversion conditions
func (rp *runicPowerBar) ConvertToDeath(sim *Simulation, slot int8, revertAt time.Duration) {
	if slices.Contains(rp.permanentDeaths, slot) && rp.runeStates&isDeaths[slot] > 0 {
		return
	}

	rp.runeStates |= isDeaths[slot]

	if rp.btSlot != slot {
		rp.runeMeta[slot].revertAt = NeverExpires
	} else {
		if rp.runeMeta[slot].revertAt != NeverExpires {
			rp.runeMeta[slot].revertAt = max(rp.runeMeta[slot].revertAt, revertAt)
		} else {
			rp.runeMeta[slot].revertAt = revertAt
		}
	}

	// Note we gained
	metrics := rp.bloodRuneGainMetrics
	if slot == 2 || slot == 3 {
		metrics = rp.frostRuneGainMetrics
	} else if slot == 4 || slot == 5 {
		metrics = rp.unholyRuneGainMetrics
	}
	if rp.runeStates&isSpents[slot] == 0 {
		// Only lose/gain if it wasn't spent (which it should be at this point)
		rp.spendRuneMetrics(sim, metrics, 1)
		rp.gainRuneMetrics(sim, rp.deathRuneGainMetrics, 1)
	}
}

func (rp *runicPowerBar) CancelBloodTap(sim *Simulation) {
	if rp.btSlot == -1 {
		return
	}
	rp.ConvertFromDeath(sim, rp.btSlot)
	bloodTapAura := rp.unit.GetAura("Blood Tap")
	bloodTapAura.Deactivate(sim)
	rp.btSlot = -1

	rp.maybeFireChange(sim, ConvertFromDeath)
}

func (rp *runicPowerBar) BloodTapConversion(sim *Simulation, bloodMetrics *ResourceMetrics, deathMetrics *ResourceMetrics) {
	changeType := None

	// 1. converts a blood rune -> death rune
	// 2. then convert one inactive blood or death rune -> active
	if rp.runeStates&isDeaths[0] == 0 {
		rp.btSlot = 0
		rp.ConvertToDeath(sim, 0, sim.CurrentTime+time.Second*20)
		changeType |= ConvertToDeath
	} else if rp.runeStates&isDeaths[1] == 0 {
		rp.btSlot = 1
		rp.ConvertToDeath(sim, 1, sim.CurrentTime+time.Second*20)
		changeType |= ConvertToDeath
	}

	if rp.runeStates&isSpents[0] > 0 {
		rp.regenRuneInternal(sim, sim.CurrentTime, 0)
		if rp.runeStates&isDeaths[0] > 0 {
			rp.gainRuneMetrics(sim, deathMetrics, 1)
		} else {
			rp.gainRuneMetrics(sim, bloodMetrics, 1)
		}
		changeType |= GainRune
	} else if rp.runeStates&isSpents[1] > 0 {
		rp.regenRuneInternal(sim, sim.CurrentTime, 1)
		if rp.runeStates&isDeaths[1] > 0 {
			rp.gainRuneMetrics(sim, deathMetrics, 1)
		} else {
			rp.gainRuneMetrics(sim, bloodMetrics, 1)
		}
		changeType |= GainRune
	}

	// if PA isn't running, make it run 20s from now to disable BT
	rp.launchPA(sim, sim.CurrentTime+20.0*time.Second)

	rp.maybeFireChange(sim, changeType)
}

func (rp *runicPowerBar) LeftBloodRuneReady() bool {
	return rp.runeStates&isSpents[0] == 0
}

func (rp *runicPowerBar) RuneIsActive(slot int8) bool {
	return rp.runeStates&isSpents[slot] == 0
}

func (rp *runicPowerBar) RuneIsDeath(slot int8) bool {
	return rp.runeStates&isDeaths[slot] != 0
}

// rune state to count of non-death, non-spent runes (0b00)
var rs2c = []int8{
	0b0000: 2, 0b0001: 1, 0b0010: 1, 0b0011: 1, 0b0100: 1, 0b0101: 0, 0b0110: 0, 0b0111: 0,
	0b1000: 1, 0b1001: 0, 0b1010: 0, 0b1011: 0, 0b1100: 1, 0b1101: 0, 0b1110: 0, 0b1111: 0,
}

func (rp *runicPowerBar) CurrentBloodRunes() int8 {
	return rs2c[(rp.runeStates>>0)&0b1111]
}

func (rp *runicPowerBar) CurrentFrostRunes() int8 {
	return rs2c[(rp.runeStates>>4)&0b1111]
}

func (rp *runicPowerBar) CurrentUnholyRunes() int8 {
	return rs2c[(rp.runeStates>>8)&0b1111]
}

// rune state to count of death, non-spent runes (0b10)
var rs2d = []int8{
	0b0000: 0, 0b0001: 0, 0b0010: 1, 0b0011: 0, 0b0100: 0, 0b0101: 0, 0b0110: 1, 0b0111: 0,
	0b1000: 1, 0b1001: 1, 0b1010: 2, 0b1011: 1, 0b1100: 0, 0b1101: 0, 0b1110: 1, 0b1111: 0,
}

func (rp *runicPowerBar) CurrentDeathRunes() int8 {
	return rs2d[(rp.runeStates>>0)&0b1111] + rs2d[(rp.runeStates>>4)&0b1111] + rs2d[(rp.runeStates>>8)&0b1111]
}

func (rp *runicPowerBar) DeathRunesInFU() int8 {
	return rs2d[(rp.runeStates>>4)&0b1111] + rs2d[(rp.runeStates>>8)&0b1111]
}

// rune state to count of non-spent runes (0bx0), masking death runes
var rs2cd = [16]int8{
	0b0000: 2, 0b0001: 1, 0b0100: 1,
}

func (rp *runicPowerBar) CurrentBloodOrDeathRunes() int8 {
	return rs2cd[(rp.runeStates>>0)&0b0101]
}

func (rp *runicPowerBar) CurrentFrostOrDeathRunes() int8 {
	return rs2cd[(rp.runeStates>>4)&0b0101]
}

func (rp *runicPowerBar) CurrentUnholyOrDeathRunes() int8 {
	return rs2cd[(rp.runeStates>>8)&0b0101]
}

func (rp *runicPowerBar) AllRunesSpent() bool {
	return rp.runeStates&allSpent == allSpent
}

func (rp *runicPowerBar) OptimalRuneCost(cost RuneCost) RuneCost {
	var b, f, u, d int8

	if b = cost.Blood(); b > 0 {
		if cb := rp.CurrentBloodRunes(); cb < b {
			d += b - cb
			b = cb
		}
	}

	if f = cost.Frost(); f > 0 {
		if cf := rp.CurrentFrostRunes(); cf < f {
			d += f - cf
			f = cf
		}
	}

	if u = cost.Unholy(); u > 0 {
		if cu := rp.CurrentUnholyRunes(); cu < u {
			d += u - cu
			u = cu
		}
	}

	if d == 0 {
		return cost
	}

	d += cost.Death()

	if cd := rp.CurrentDeathRunes(); cd >= d {
		return NewRuneCost(cost.RunicPower(), b, f, u, d)
	}

	return 0
}

func (rp *runicPowerBar) spendRuneCost(sim *Simulation, spell *Spell, cost RuneCost) (RuneChangeType, []int8) {
	if !cost.HasRune() {
		if rpc := cost.RunicPower(); rpc > 0 {
			rp.spendRunicPower(sim, float64(cost.RunicPower()), spell.RunicPowerMetrics())
		}
		return None, nil
	}

	b, f, u, d := cost.Blood(), cost.Frost(), cost.Unholy(), cost.Death()
	slots := make([]int8, 0, b+f+u+d)
	for i := int8(0); i < b; i++ {
		slots = append(slots, rp.spendRune(sim, 0, spell.BloodRuneMetrics()))
	}
	for i := int8(0); i < f; i++ {
		slots = append(slots, rp.spendRune(sim, 2, spell.FrostRuneMetrics()))
	}
	for i := int8(0); i < u; i++ {
		slots = append(slots, rp.spendRune(sim, 4, spell.UnholyRuneMetrics()))
	}
	if d > 0 {
		defaultCost := RuneCost(spell.DefaultCast.Cost)
		for i, mu := int8(0), defaultCost.Unholy()-u; i < mu; i++ {
			slots = append(slots, rp.spendDeathRune(sim, []int8{4, 5, 2, 3, 0, 1}, spell.DeathRuneMetrics()))
		}
		for i, mf := int8(0), defaultCost.Frost()-f; i < mf; i++ {
			slots = append(slots, rp.spendDeathRune(sim, []int8{2, 3, 4, 5, 0, 1}, spell.DeathRuneMetrics()))
		}
		for i, mb := int8(0), defaultCost.Blood()-b; i < mb; i++ {
			slots = append(slots, rp.spendDeathRune(sim, []int8{0, 1, 4, 5, 2, 3}, spell.DeathRuneMetrics()))
		}
	}

	if rpc := cost.RunicPower(); rpc > 0 {
		rp.AddRunicPower(sim, float64(rpc), spell.RunicPowerMetrics())
	}

	return SpendRune, slots
}

func (rp *runicPowerBar) typeAmount(metrics *ResourceMetrics) (string, int8) {
	switch metrics.Type {
	case proto.ResourceType_ResourceTypeDeathRune:
		return "death", rp.CurrentDeathRunes()
	case proto.ResourceType_ResourceTypeBloodRune:
		return "blood", rp.CurrentBloodRunes()
	case proto.ResourceType_ResourceTypeFrostRune:
		return "frost", rp.CurrentFrostRunes()
	case proto.ResourceType_ResourceTypeUnholyRune:
		return "unholy", rp.CurrentUnholyRunes()
	default:
		panic("invalid metrics for rune gaining")
	}
}

// gainRuneMetrics should be called after gaining the rune
func (rp *runicPowerBar) gainRuneMetrics(sim *Simulation, metrics *ResourceMetrics, gainAmount int8) {
	metrics.AddEvent(float64(gainAmount), float64(gainAmount))

	if sim.Log != nil {
		name, currRunes := rp.typeAmount(metrics)
		rp.unit.Log(sim, "Gained %0.3f %s rune from %s (%d --> %d).", float64(gainAmount), name, metrics.ActionID, currRunes-gainAmount, currRunes)
	}
}

// spendRuneMetrics should be called after spending the rune
func (rp *runicPowerBar) spendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, spendAmount int8) {
	metrics.AddEvent(-float64(spendAmount), -float64(spendAmount))

	if sim.Log != nil {
		name, currRunes := rp.typeAmount(metrics)
		rp.unit.Log(sim, "Spent 1.000 %s rune from %s (%d --> %d).", name, metrics.ActionID, currRunes+spendAmount, currRunes)
	}
}

func (rp *runicPowerBar) regenRune(sim *Simulation, regenAt time.Duration, slot int8) {
	rp.regenRuneInternal(sim, regenAt, slot)

	metrics := rp.bloodRuneGainMetrics
	if rp.runeStates&isDeaths[slot] > 0 {
		metrics = rp.deathRuneGainMetrics
	} else if slot == 2 || slot == 3 {
		metrics = rp.frostRuneGainMetrics
	} else if slot == 4 || slot == 5 {
		metrics = rp.unholyRuneGainMetrics
	}

	rp.gainRuneMetrics(sim, metrics, 1)
}

func (rp *runicPowerBar) regenRuneInternal(sim *Simulation, regenAt time.Duration, slot int8) {
	rp.lastRegen = append(rp.lastRegen, slot)
	rp.runeStates ^= isSpents[slot] // unset spent flag for this rune.
	rp.runeMeta[slot].regenAt = NeverExpires

	// if other slot rune is spent start and not regening start regen
	otherSlot := (slot/2)*2 + (slot+1)%2
	if rp.runeStates&isSpents[otherSlot] > 0 && rp.runeMeta[otherSlot].regenAt == NeverExpires {
		rp.launchRuneRegen(sim, otherSlot)
	}
}

func (rp *runicPowerBar) RegenAllRunes(sim *Simulation, metrics []*ResourceMetrics) {
	changeType := None
	for i := range rp.runeMeta {
		if rp.runeStates&isSpents[i] > 0 {
			rp.regenRuneInternal(sim, sim.CurrentTime, int8(i))

			metric := metrics[0]
			if rp.runeStates&isDeaths[i] > 0 {
				metric = metrics[3]
			} else if i == 2 || i == 3 {
				metric = metrics[1]
			} else if i == 4 || i == 5 {
				metric = metrics[2]
			}

			rp.gainRuneMetrics(sim, metric, 1)

			changeType = GainRune
		}
	}

	rp.maybeFireChange(sim, changeType)
}

func (rp *runicPowerBar) RegenRandomDepletedRune(sim *Simulation, runeMetrics []*ResourceMetrics) {
	changeType := None
	possibleRunes := make([]int, 0)
	for i := range rp.runeMeta {
		if rp.runeStates&isSpents[i] > 0 && rp.runeMeta[i].regenAt == NeverExpires {
			possibleRunes = append(possibleRunes, i)
		}
	}

	if len(possibleRunes) == 0 {
		return
	}

	randomRuneIndex := int(math.Floor(sim.RandomFloat("Rune Regen") * float64(len(possibleRunes))))
	randomRune := int8(possibleRunes[randomRuneIndex])

	rp.regenRuneInternal(sim, sim.CurrentTime, randomRune)
	changeType = GainRune
	if rp.runeStates&isDeaths[randomRune] > 0 {
		rp.gainRuneMetrics(sim, runeMetrics[3], 1)
	} else {
		rp.gainRuneMetrics(sim, runeMetrics[randomRune/2], 1)
	}
	rp.maybeFireChange(sim, changeType)
}

func (rp *runicPowerBar) MultiplyRuneRegenSpeed(sim *Simulation, multiplier float64) {
	rp.runeRegenMultiplier *= multiplier
	rp.updateRegenTimes(sim)
}

func (rp *runicPowerBar) MultiplyRunicRegen(multiply float64) {
	rp.runicRegenMultiplier *= multiply
}

func (rp *runicPowerBar) getTotalRegenMultiplier() float64 {
	hasteMultiplier := 1.0 + rp.unit.GetStat(stats.MeleeHaste)/(100*HasteRatingPerHastePercent)
	totalMultiplier := 1 / (hasteMultiplier * rp.runeRegenMultiplier)
	return totalMultiplier
}

func (rp *runicPowerBar) updateRegenTimes(sim *Simulation) {
	if rp.unit == nil {
		return
	}
	totalMultiplier := rp.getTotalRegenMultiplier()

	for slot := int8(0); slot < 6; slot++ {
		if rp.runeStates&isSpents[slot] > 0 && rp.runeMeta[slot].regenAt != NeverExpires {
			// Is spent so we save current progress and then rescale pa
			regenTime := DurationFromSeconds(rp.runeMeta[slot].unscaledRegenLeft.Seconds() * rp.runeMeta[slot].regenMulti)
			startTime := rp.runeMeta[slot].regenAt - regenTime
			unscaledRegenDone := (sim.CurrentTime - startTime).Seconds() / rp.runeMeta[slot].regenMulti
			regenLeft := (rp.runeMeta[slot].unscaledRegenLeft.Seconds() - unscaledRegenDone)

			rp.runeMeta[slot].regenMulti = totalMultiplier
			rp.runeMeta[slot].regenAt = sim.CurrentTime + DurationFromSeconds(regenLeft*totalMultiplier)
			rp.runeMeta[slot].unscaledRegenLeft = DurationFromSeconds(regenLeft)

			rp.launchPA(sim, rp.runeMeta[slot].regenAt)
		}
	}
}

func (rp *runicPowerBar) launchRuneRegen(sim *Simulation, slot int8) {
	totalMultiplier := rp.getTotalRegenMultiplier()

	rp.runeMeta[slot].regenMulti = totalMultiplier
	rp.runeMeta[slot].regenAt = sim.CurrentTime + DurationFromSeconds(rp.runeCD.Seconds()*totalMultiplier)
	rp.runeMeta[slot].unscaledRegenLeft = rp.runeCD

	rp.launchPA(sim, rp.runeMeta[slot].regenAt)
}

func (rp *runicPowerBar) launchPA(sim *Simulation, at time.Duration) {
	if rp.pa != nil {
		if at >= rp.pa.NextActionAt {
			return
		}
		// If this new regen is before currently scheduled one, we must cancel old regen and start a new one.
		rp.pa.Cancel(sim)
	}

	pa := &PendingAction{
		NextActionAt: at,
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		if !pa.cancelled {
			// regenerate and revert
			rp.Advance(sim, sim.CurrentTime)

			// Check when we need next check
			pa.NextActionAt = min(rp.AnySpentRuneReadyAt(), rp.DeathRuneRevertAt())
			if pa.NextActionAt < NeverExpires {
				sim.AddPendingAction(pa)
			}
		}
	}
	rp.pa = pa
	sim.AddPendingAction(pa)

}

func (rp *runicPowerBar) Advance(sim *Simulation, newTime time.Duration) {
	changeType := None
	if rp.runeStates&allDeath > 0 {
		for i, rm := range rp.runeMeta {
			if rm.revertAt <= newTime {
				if rp.btSlot == int8(i) {
					rp.btSlot = -1 // this was the BT slot
				}
				rp.ConvertFromDeath(sim, int8(i))
				changeType |= ConvertFromDeath
			}
		}
	}

	if rp.runeStates&allSpent > 0 {
		for i, rm := range rp.runeMeta {
			if rm.regenAt <= newTime && rp.runeStates&isSpents[i] > 0 {
				rp.regenRune(sim, newTime, int8(i))
				changeType |= GainRune
			}
		}
	}

	rp.maybeFireChange(sim, changeType)
}

func (rp *runicPowerBar) spendRune(sim *Simulation, firstSlot int8, metrics *ResourceMetrics) int8 {
	slot := rp.findReadyRune(firstSlot)
	rp.runeStates |= isSpents[slot]

	rp.spendRuneMetrics(sim, metrics, 1)

	// if other rune is not spent start regen
	otherSlot := (slot/2)*2 + (slot+1)%2
	if rp.runeStates&isSpents[otherSlot] == 0 {
		rp.launchRuneRegen(sim, slot)
	}
	return slot
}

func (rp *runicPowerBar) findReadyRune(slot int8) int8 {
	if rp.runeStates&isSpentDeath[slot] == 0 {
		return slot
	}
	if rp.runeStates&isSpentDeath[slot+1] == 0 {
		return slot + 1
	}
	panic(fmt.Sprintf("findReadyRune(%d) - no slot found (runeStates = %12b)", slot, rp.runeStates))
}

func (rp *runicPowerBar) spendDeathRune(sim *Simulation, order []int8, metrics *ResourceMetrics) int8 {
	slot := rp.findReadyDeathRune(order)
	if rp.btSlot != slot && !slices.Contains(rp.permanentDeaths, slot) {
		rp.runeMeta[slot].revertAt = NeverExpires // disable revert at
		rp.runeStates ^= isDeaths[slot]           // clear death bit to revert.
	}

	// mark spent bit to spend
	rp.runeStates |= isSpents[slot]

	rp.spendRuneMetrics(sim, metrics, 1)

	// if other rune is not spent start regen
	otherSlot := (slot/2)*2 + (slot+1)%2
	if rp.runeStates&isSpents[otherSlot] == 0 {
		rp.launchRuneRegen(sim, slot)
	}
	return slot
}

// findReadyDeathRune returns the slot of first available death rune in the order given.
func (rp *runicPowerBar) findReadyDeathRune(order []int8) int8 {
	for _, slot := range order {
		if rp.runeStates&isSpentDeath[slot] == isDeaths[slot] {
			return slot
		}
	}
	panic(fmt.Sprintf("findReadyDeathRune() - no slot found (runeStates = %12b)", rp.runeStates))
}

func (rp *runicPowerBar) IsBloodTappedRune(slot int8) bool {
	return slot == rp.btSlot
}

type RuneCostOptions struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	RunicPowerCost float64
	RunicPowerGain float64
	Refundable     bool
	RefundCost     float64
}

type RuneCostImpl struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	RunicPowerCost float64
	RunicPowerGain float64
	Refundable     bool
	RefundCost     float64

	runicPowerMetrics *ResourceMetrics
	bloodRuneMetrics  *ResourceMetrics
	frostRuneMetrics  *ResourceMetrics
	unholyRuneMetrics *ResourceMetrics
	deathRuneMetrics  *ResourceMetrics
}

func newRuneCost(spell *Spell, options RuneCostOptions) *RuneCostImpl {
	baseCost := float64(NewRuneCost(int16(options.RunicPowerCost), options.BloodRuneCost, options.FrostRuneCost, options.UnholyRuneCost, 0))
	spell.DefaultCast.Cost = baseCost
	spell.CurCast.Cost = baseCost

	return &RuneCostImpl{
		BloodRuneCost:  options.BloodRuneCost,
		FrostRuneCost:  options.FrostRuneCost,
		UnholyRuneCost: options.UnholyRuneCost,
		RunicPowerCost: options.RunicPowerCost,
		RunicPowerGain: options.RunicPowerGain,
		Refundable:     options.Refundable,
		RefundCost:     options.RefundCost,

		runicPowerMetrics: Ternary(options.RunicPowerCost > 0 || options.RunicPowerGain > 0, spell.Unit.NewRunicPowerMetrics(spell.ActionID), nil),
		bloodRuneMetrics:  Ternary(options.BloodRuneCost > 0, spell.Unit.NewBloodRuneMetrics(spell.ActionID), nil),
		frostRuneMetrics:  Ternary(options.FrostRuneCost > 0, spell.Unit.NewFrostRuneMetrics(spell.ActionID), nil),
		unholyRuneMetrics: Ternary(options.UnholyRuneCost > 0, spell.Unit.NewUnholyRuneMetrics(spell.ActionID), nil),
		deathRuneMetrics:  spell.Unit.NewDeathRuneMetrics(spell.ActionID),
	}
}

func (rc *RuneCostImpl) GetConfig() RuneCostOptions {
	return RuneCostOptions{
		BloodRuneCost:  rc.BloodRuneCost,
		FrostRuneCost:  rc.FrostRuneCost,
		UnholyRuneCost: rc.UnholyRuneCost,
		RunicPowerCost: rc.RunicPowerCost,
		RunicPowerGain: rc.RunicPowerGain,
		Refundable:     rc.Refundable,
	}
}

func (rc *RuneCostImpl) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost *= spell.CostMultiplier // TODO this looks fishy - multiplying and rune costs don't go well together

	cost := RuneCost(spell.CurCast.Cost)
	if cost == 0 {
		return true
	}

	if !cost.HasRune() {
		if float64(cost.RunicPower()) > spell.Unit.CurrentRunicPower() {
			return false
		}
	}

	optCost := spell.Unit.OptimalRuneCost(cost)
	if optCost == 0 { // no combo of runes to fulfill cost
		return false
	}
	spell.CurCast.Cost = float64(optCost) // assign chosen runes to the cost
	return true
}

func (rc *RuneCostImpl) CostFailureReason(_ *Simulation, _ *Spell) string {
	return "not enough RP or runes"
}

func (rc *RuneCostImpl) SpendCost(sim *Simulation, spell *Spell) {
	// Spend now if there is no way to refund the spell
	if !rc.Refundable {
		changeType, _ := spell.Unit.spendRuneCost(sim, spell, RuneCost(spell.CurCast.Cost))
		spell.Unit.maybeFireChange(sim, changeType)

		if rc.RunicPowerGain > 0 && spell.CurCast.Cost > 0 {
			spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
		}
	}
}

func (rc *RuneCostImpl) spendRefundableCost(sim *Simulation, spell *Spell, result *SpellResult) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if result.Landed() {
		changeType, _ := spell.Unit.spendRuneCost(sim, spell, cost)
		spell.Unit.maybeFireChange(sim, changeType)

		if rc.RunicPowerGain > 0 {
			spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
		}
	} else if rc.RefundCost > 0 {
		spell.Unit.spendRunicPower(sim, rc.RefundCost, spell.RunicPowerMetrics())
	}
}

func (spell *Spell) SpendRefundableCost(sim *Simulation, result *SpellResult) {
	spell.Cost.(*RuneCostImpl).spendRefundableCost(sim, spell, result)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertBloodRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !result.Landed() {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		if rc.RefundCost > 0 {
			spell.Unit.spendRunicPower(sim, rc.RefundCost, spell.RunicPowerMetrics())
		}
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	if !sim.Proc(convertChance, "Reaping") {
		spell.Unit.maybeFireChange(sim, changeType)
		return
	}

	for _, slot := range slots {
		if slot == 0 || slot == 1 {
			// If the slot to be converted is already blood-tapped, then we convert the other blood rune
			if spell.Unit.IsBloodTappedRune(slot) {
				otherRune := (slot + 1) % 2
				spell.Unit.ConvertToDeath(sim, otherRune, NeverExpires)
				changeType |= ConvertToDeath
			} else {
				spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
				changeType |= ConvertToDeath
			}
		}
	}

	if rc.RunicPowerGain > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}

	spell.Unit.maybeFireChange(sim, changeType)
}

func (spell *Spell) SpendRefundableCostAndConvertBloodRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).spendRefundableCostAndConvertBloodRune(sim, spell, result, convertChance)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !result.Landed() {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		if rc.RefundCost > 0 {
			spell.Unit.spendRunicPower(sim, rc.RefundCost, spell.RunicPowerMetrics())
		}
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	if !sim.Proc(convertChance, "Blood Rites") {
		spell.Unit.maybeFireChange(sim, changeType)
		return
	}

	for _, slot := range slots {
		if slot == 2 || slot == 3 || slot == 4 || slot == 5 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
			changeType |= ConvertToDeath
		}
	}

	if rc.RunicPowerGain > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}

	spell.Unit.maybeFireChange(sim, changeType)
}

func (spell *Spell) SpendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).spendRefundableCostAndConvertFrostOrUnholyRune(sim, spell, result, convertChance)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertBloodOrFrostRune(sim *Simulation, spell *Spell, result *SpellResult, convertChance float64) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !result.Landed() {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		if rc.RefundCost > 0 {
			spell.Unit.spendRunicPower(sim, rc.RefundCost, spell.RunicPowerMetrics())
		}
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	if !sim.Proc(convertChance, "Reaping") {
		spell.Unit.maybeFireChange(sim, changeType)
		return
	}

	for _, slot := range slots {
		if slot == 0 || slot == 1 {
			// If the slot to be converted is already blood-tapped, then we convert the other blood rune
			if spell.Unit.IsBloodTappedRune(slot) {
				otherRune := (slot + 1) % 2
				spell.Unit.ConvertToDeath(sim, otherRune, NeverExpires)
				changeType |= ConvertToDeath
			} else {
				spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
				changeType |= ConvertToDeath
			}
		}
		if slot == 2 || slot == 3 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
			changeType |= ConvertToDeath
		}
	}

	if rc.RunicPowerGain > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}

	spell.Unit.maybeFireChange(sim, changeType)
}

func (spell *Spell) SpendRefundableCostAndConvertBloodOrFrostRune(sim *Simulation, result *SpellResult, convertChance float64) {
	spell.Cost.(*RuneCostImpl).spendRefundableCostAndConvertBloodOrFrostRune(sim, spell, result, convertChance)
}

func (rc *RuneCostImpl) IssueRefund(_ *Simulation, _ *Spell) {
	// Instead of issuing refunds we just don't charge the cost of spells which
	// miss; this is better for perf since we'd have to cancel the regen actions.
}

func (spell *Spell) RuneCostImpl() *RuneCostImpl {
	return spell.Cost.(*RuneCostImpl)
}

func (spell *Spell) RunicPowerMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).runicPowerMetrics
}

func (spell *Spell) BloodRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).bloodRuneMetrics
}

func (spell *Spell) FrostRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).frostRuneMetrics
}

func (spell *Spell) UnholyRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).unholyRuneMetrics
}

func (spell *Spell) DeathRuneMetrics() *ResourceMetrics {
	return spell.Cost.(*RuneCostImpl).deathRuneMetrics
}
