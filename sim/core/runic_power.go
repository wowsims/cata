package core

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
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
	character *Character

	maxRunicPower      float64
	startingRunicPower float64
	currentRunicPower  float64
	runeCD             time.Duration

	// These flags are used to simplify pending action checks
	// |DS|DS|DS|DS|DS|DS|
	runeStates int16
	runeMeta   [6]RuneMeta

	bloodRuneGainMetrics  *ResourceMetrics
	frostRuneGainMetrics  *ResourceMetrics
	unholyRuneGainMetrics *ResourceMetrics
	deathRuneGainMetrics  *ResourceMetrics

	spellRunicPowerMetrics map[ActionID]*ResourceMetrics
	spellBloodRuneMetrics  map[ActionID]*ResourceMetrics
	spellFrostRuneMetrics  map[ActionID]*ResourceMetrics
	spellUnholyRuneMetrics map[ActionID]*ResourceMetrics
	spellDeathRuneMetrics  map[ActionID]*ResourceMetrics

	onRuneChange     OnRuneChange
	onRunicPowerGain OnRunicPowerGain

	pa *PendingAction

	runeRegenMultiplier  float64
	runicRegenMultiplier float64

	runicRegenMultiplierDisabled bool

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
	if rp.character == nil {
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

	rp.currentRunicPower = rp.startingRunicPower
}

func (character *Character) EnableRunicPowerBar(startingRunicPower float64, maxRunicPower float64, runeCD time.Duration,
	onRuneChange OnRuneChange, onRunicPowerGain OnRunicPowerGain) {
	character.SetCurrentPowerBar(RunicPower)
	character.runicPowerBar = runicPowerBar{
		character: character,

		maxRunicPower:        maxRunicPower,
		currentRunicPower:    startingRunicPower,
		startingRunicPower:   startingRunicPower,
		runeCD:               runeCD,
		runeRegenMultiplier:  1.0,
		runicRegenMultiplier: 1.0,

		runeStates: baseRuneState,

		onRuneChange:     onRuneChange,
		onRunicPowerGain: onRunicPowerGain,

		permanentDeaths: make([]int8, 0),
		lastRegen:       make([]int8, 0),

		spellRunicPowerMetrics: make(map[ActionID]*ResourceMetrics),
		spellBloodRuneMetrics:  make(map[ActionID]*ResourceMetrics),
		spellFrostRuneMetrics:  make(map[ActionID]*ResourceMetrics),
		spellUnholyRuneMetrics: make(map[ActionID]*ResourceMetrics),
		spellDeathRuneMetrics:  make(map[ActionID]*ResourceMetrics),
	}

	character.bloodRuneGainMetrics = character.NewBloodRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionBloodRuneGain, Tag: 1})
	character.frostRuneGainMetrics = character.NewFrostRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionFrostRuneGain, Tag: 1})
	character.unholyRuneGainMetrics = character.NewUnholyRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionUnholyRuneGain, Tag: 1})
	character.deathRuneGainMetrics = character.NewDeathRuneMetrics(ActionID{OtherID: proto.OtherAction_OtherActionDeathRuneGain, Tag: 1})
}

func (unit *Unit) HasRunicPowerBar() bool {
	return unit.runicPowerBar.character != nil
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

func (rp *runicPowerBar) MaximumRunicPower() float64 {
	return rp.maxRunicPower
}

func (rp *runicPowerBar) maybeFireChange(sim *Simulation, changeType RuneChangeType) {
	if changeType != None && rp.onRuneChange != nil {
		rp.onRuneChange(sim, changeType, rp.lastRegen)
		// Clear regen runes
		rp.lastRegen = make([]int8, 0)
	}
}

func (rp *runicPowerBar) addRunicPowerInternal(sim *Simulation, amount float64, metrics *ResourceMetrics, withMultiplier bool) {
	if amount < 0 {
		panic("Trying to add negative runic power!")
	}

	runicRegenMultiplier := rp.runicRegenMultiplier
	if !withMultiplier || rp.runicRegenMultiplierDisabled {
		runicRegenMultiplier = 1.0
	}
	amount *= runicRegenMultiplier
	newRunicPower := min(rp.currentRunicPower+amount, rp.maxRunicPower)

	if sim.CurrentTime > 0 {
		metrics.AddEvent(amount, newRunicPower-rp.currentRunicPower)
	}

	if sim.Log != nil {
		rp.character.Log(sim, "Gained %0.3f runic power from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower, rp.maxRunicPower)
	}

	rp.currentRunicPower = newRunicPower
}

func (rp *runicPowerBar) AddRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInternal(sim, amount, metrics, true)
	if rp.onRunicPowerGain != nil {
		rp.onRunicPowerGain(sim)
	}
}

func (rp *runicPowerBar) AddUnscaledRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	rp.addRunicPowerInternal(sim, amount, metrics, false)
	if rp.onRunicPowerGain != nil {
		rp.onRunicPowerGain(sim)
	}
}

func (rp *runicPowerBar) spendRunicPower(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative runic power!")
	}

	amount = min(amount, rp.currentRunicPower)

	newRunicPower := rp.currentRunicPower - amount

	if sim.CurrentTime > 0 {
		metrics.AddEvent(-amount, -amount)
	}

	if sim.Log != nil {
		rp.character.Log(sim, "Spent %0.3f runic power from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, rp.currentRunicPower, newRunicPower, rp.maxRunicPower)
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

	rp.runeMeta[slot].revertAt = NeverExpires

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
		defaultCost := RuneCost(spell.Cost.GetCurrentCost())
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
		rp.character.Log(sim, "Gained %0.3f %s rune from %s (%d --> %d).", float64(gainAmount), name, metrics.ActionID, currRunes-gainAmount, currRunes)
	}
}

// spendRuneMetrics should be called after spending the rune
func (rp *runicPowerBar) spendRuneMetrics(sim *Simulation, metrics *ResourceMetrics, spendAmount int8) {
	metrics.AddEvent(-float64(spendAmount), -float64(spendAmount))

	if sim.Log != nil {
		name, currRunes := rp.typeAmount(metrics)
		rp.character.Log(sim, "Spent 1.000 %s rune from %s (%d --> %d).", name, metrics.ActionID, currRunes+spendAmount, currRunes)
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

func (rp *runicPowerBar) RegenAllFrostAndUnholyRunesAsDeath(sim *Simulation, deathRuneMetrics *ResourceMetrics) {
	changeType := None
	for i := 2; i < 6; i++ {
		if rp.runeStates&isSpents[i] > 0 {
			rp.regenRuneInternal(sim, sim.CurrentTime, int8(i))

			rp.gainRuneMetrics(sim, deathRuneMetrics, 1)

			changeType = GainRune
		}

		if rp.runeStates&isDeaths[i] == 0 {
			rp.ConvertToDeath(sim, int8(i), NeverExpires)
			changeType |= ConvertToDeath
		}
	}

	rp.maybeFireChange(sim, changeType)
}

func (rp *runicPowerBar) AnyDepletedRunes() bool {
	for slot := range rp.runeMeta {
		if rp.isDepleted(slot) {
			return true
		}
	}

	return false
}

func (rp *runicPowerBar) isDepleted(runeSlot int) bool {
	return rp.runeStates&isSpents[runeSlot] > 0 && rp.runeMeta[runeSlot].regenAt == NeverExpires
}

type depletedRune struct {
	runeSlot int8
	regenAt  time.Duration
}

func getHighestCDRune(sim *Simulation, possibleRunes []*depletedRune) int8 {
	maxRegenAt := sim.CurrentTime
	for _, rune := range possibleRunes {
		if rune.regenAt > maxRegenAt {
			maxRegenAt = rune.regenAt
		}
	}

	filteredRunes := make([]int8, 0)
	for _, rune := range possibleRunes {
		if rune.regenAt == maxRegenAt {
			filteredRunes = append(filteredRunes, rune.runeSlot)
		}
	}

	randomRuneIndex := int(math.Floor(sim.RandomFloat("Rune Regen") * float64(len(filteredRunes))))
	return filteredRunes[randomRuneIndex]
}

// Runic Empowerment prioritizes fully depleted runes with the highest cd
func (rp *runicPowerBar) RegenRunicEmpowermentRune(sim *Simulation, runeMetrics []*ResourceMetrics) {
	possibleRunes := make([]*depletedRune, 0)

	if rp.isDepleted(0) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 0, regenAt: rp.runeMeta[1].regenAt})
	} else if rp.isDepleted(1) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 1, regenAt: rp.runeMeta[0].regenAt})
	}
	if rp.isDepleted(2) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 2, regenAt: rp.runeMeta[3].regenAt})
	} else if rp.isDepleted(3) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 3, regenAt: rp.runeMeta[2].regenAt})
	}
	if rp.isDepleted(4) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 4, regenAt: rp.runeMeta[5].regenAt})
	} else if rp.isDepleted(5) {
		possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 5, regenAt: rp.runeMeta[4].regenAt})
	}

	if len(possibleRunes) == 0 {
		return
	}

	slot := getHighestCDRune(sim, possibleRunes)

	rp.regenRuneInternal(sim, sim.CurrentTime, slot)
	if rp.runeStates&isDeaths[slot] > 0 {
		rp.gainRuneMetrics(sim, runeMetrics[3], 1)
	} else {
		rp.gainRuneMetrics(sim, runeMetrics[slot/2], 1)
	}

	rp.maybeFireChange(sim, GainRune)
}

// Plague leech prioritizes runes based on spec
// Unholy prefers to regen a pair of B/F runes first, then U if no other runes are available.
// Blood and Frost prefers to regen a pair of F/U runes first, then B if no other runes are available.
func (rp *runicPowerBar) ConvertAndRegenPlagueLeechRunes(sim *Simulation, spell *Spell, runeMetrics []*ResourceMetrics) {
	runesToRegen := make([]int8, 0)

	if rp.character.Spec == proto.Spec_SpecUnholyDeathKnight {
		if rp.isDepleted(0) {
			runesToRegen = append(runesToRegen, 0)
		} else if rp.isDepleted(1) {
			runesToRegen = append(runesToRegen, 1)
		}

		if rp.isDepleted(2) {
			runesToRegen = append(runesToRegen, 2)
		} else if rp.isDepleted(3) {
			runesToRegen = append(runesToRegen, 3)
		}

		if len(runesToRegen) < 2 {
			if rp.isDepleted(4) {
				runesToRegen = append(runesToRegen, 4)
			} else if rp.isDepleted(5) {
				runesToRegen = append(runesToRegen, 5)
			}
		}
	} else {
		if rp.isDepleted(2) {
			runesToRegen = append(runesToRegen, 2)
		} else if rp.isDepleted(3) {
			runesToRegen = append(runesToRegen, 3)
		}

		if rp.isDepleted(4) {
			runesToRegen = append(runesToRegen, 4)
		} else if rp.isDepleted(5) {
			runesToRegen = append(runesToRegen, 5)
		}

		if len(runesToRegen) < 2 {
			if rp.isDepleted(0) {
				runesToRegen = append(runesToRegen, 0)
			} else if rp.isDepleted(1) {
				runesToRegen = append(runesToRegen, 1)
			}
		}
	}

	if len(runesToRegen) == 0 {
		return
	}

	for _, slot := range runesToRegen {
		spell.Unit.ConvertToDeath(sim, slot, NeverExpires)

		rp.regenRuneInternal(sim, sim.CurrentTime, slot)
		if rp.runeStates&isDeaths[slot] > 0 {
			rp.gainRuneMetrics(sim, runeMetrics[3], 1)
		} else {
			rp.gainRuneMetrics(sim, runeMetrics[slot/2], 1)
		}
	}

	rp.maybeFireChange(sim, ConvertToDeath|GainRune)
}

// Blood tap prioritizes runes based on spec
// Blood prefers to regen B runes first, then the F/U rune with the highest CD if no blood runes are available.
// Frost prefers to regen U runes first, then the B/F rune with the highest CD if no other runes are available.
// Unholy prefers to regen a B/F rune with the highest CD first, then an U rune if no other runes are available.
func (rp *runicPowerBar) ConvertAndRegenBloodTapRune(sim *Simulation, spell *Spell, runeMetrics []*ResourceMetrics) bool {
	slot := int8(-1)
	possibleRunes := make([]*depletedRune, 0)

	if rp.character.Spec == proto.Spec_SpecBloodDeathKnight {
		if rp.isDepleted(0) {
			slot = 0
		} else if rp.isDepleted(1) {
			slot = 1
		} else {
			if rp.isDepleted(2) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 2, regenAt: rp.runeMeta[3].regenAt})
			} else if rp.isDepleted(3) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 3, regenAt: rp.runeMeta[2].regenAt})
			}
			if rp.isDepleted(4) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 4, regenAt: rp.runeMeta[5].regenAt})
			} else if rp.isDepleted(5) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 5, regenAt: rp.runeMeta[4].regenAt})
			}
		}
	} else if rp.character.Spec == proto.Spec_SpecFrostDeathKnight {
		if rp.isDepleted(4) {
			slot = 4
		} else if rp.isDepleted(5) {
			slot = 5
		} else {
			if rp.isDepleted(0) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 0, regenAt: rp.runeMeta[1].regenAt})
			} else if rp.isDepleted(1) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 1, regenAt: rp.runeMeta[0].regenAt})
			}
			if rp.isDepleted(2) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 2, regenAt: rp.runeMeta[3].regenAt})
			} else if rp.isDepleted(3) {
				possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 3, regenAt: rp.runeMeta[2].regenAt})
			}
		}
	} else {
		if rp.isDepleted(0) {
			possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 0, regenAt: rp.runeMeta[1].regenAt})
		} else if rp.isDepleted(1) {
			possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 1, regenAt: rp.runeMeta[0].regenAt})
		}
		if rp.isDepleted(2) {
			possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 2, regenAt: rp.runeMeta[3].regenAt})
		} else if rp.isDepleted(3) {
			possibleRunes = append(possibleRunes, &depletedRune{runeSlot: 3, regenAt: rp.runeMeta[2].regenAt})
		}

		if len(possibleRunes) == 0 {
			if rp.isDepleted(4) {
				slot = 4
			} else if rp.isDepleted(5) {
				slot = 5
			}
		}
	}

	if len(possibleRunes) > 0 {
		slot = getHighestCDRune(sim, possibleRunes)
	}

	if slot == -1 {
		return false
	}

	spell.Unit.ConvertToDeath(sim, slot, NeverExpires)

	rp.regenRuneInternal(sim, sim.CurrentTime, slot)
	if rp.runeStates&isDeaths[slot] > 0 {
		rp.gainRuneMetrics(sim, runeMetrics[3], 1)
	} else {
		rp.gainRuneMetrics(sim, runeMetrics[slot/2], 1)
	}
	rp.maybeFireChange(sim, ConvertToDeath|GainRune)

	return true
}

func (rp *runicPowerBar) MultiplyRuneRegenSpeed(sim *Simulation, multiplier float64) {
	rp.runeRegenMultiplier *= multiplier
	rp.updateRegenTimes(sim)
}

func (rp *runicPowerBar) MultiplyRunicRegen(multiply float64) {
	rp.runicRegenMultiplier *= multiply
}

func (rp *runicPowerBar) GetRuneRegenMultiplier() float64 {
	return rp.runeRegenMultiplier
}

func (rp *runicPowerBar) getTotalRegenMultiplier() float64 {
	hasteMultiplier := 1.0 + rp.character.GetStat(stats.HasteRating)/(100*HasteRatingPerHastePercent)
	totalMultiplier := 1 / (hasteMultiplier * rp.runeRegenMultiplier)
	return totalMultiplier
}

func (rp *runicPowerBar) updateRegenTimes(sim *Simulation) {
	if rp.character == nil {
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

	// Query APL if a change occurred
	if changeType != None {
		rp.character.ReactToEvent(sim)
	}
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
	if !slices.Contains(rp.permanentDeaths, slot) {
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

type RuneCostOptions struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	DeathRuneCost  int8
	RunicPowerCost float64
	RunicPowerGain float64
	Refundable     bool
	RefundCost     float64
}

type RuneCostImpl struct {
	BloodRuneCost  int8
	FrostRuneCost  int8
	UnholyRuneCost int8
	DeathRuneCost  int8
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

func newRuneCost(spell *Spell, options RuneCostOptions) *SpellCost {
	return &SpellCost{
		spell:           spell,
		BaseCost:        int32(NewRuneCost(int16(options.RunicPowerCost), options.BloodRuneCost, options.FrostRuneCost, options.UnholyRuneCost, options.DeathRuneCost)),
		PercentModifier: 1,
		ResourceCostImpl: &RuneCostImpl{
			BloodRuneCost:  options.BloodRuneCost,
			FrostRuneCost:  options.FrostRuneCost,
			UnholyRuneCost: options.UnholyRuneCost,
			DeathRuneCost:  options.DeathRuneCost,
			RunicPowerCost: options.RunicPowerCost,
			RunicPowerGain: options.RunicPowerGain,
			Refundable:     options.Refundable,
			RefundCost:     options.RefundCost,

			runicPowerMetrics: Ternary(options.RunicPowerCost > 0 || options.RunicPowerGain > 0, spell.Unit.NewRunicPowerMetrics(spell.ActionID), nil),
			bloodRuneMetrics:  Ternary(options.BloodRuneCost > 0, spell.Unit.NewBloodRuneMetrics(spell.ActionID), nil),
			frostRuneMetrics:  Ternary(options.FrostRuneCost > 0, spell.Unit.NewFrostRuneMetrics(spell.ActionID), nil),
			unholyRuneMetrics: Ternary(options.UnholyRuneCost > 0, spell.Unit.NewUnholyRuneMetrics(spell.ActionID), nil),
			deathRuneMetrics:  spell.Unit.NewDeathRuneMetrics(spell.ActionID),
		},
	}
}

func (rc *RuneCostImpl) GetConfig() RuneCostOptions {
	return RuneCostOptions{
		BloodRuneCost:  rc.BloodRuneCost,
		FrostRuneCost:  rc.FrostRuneCost,
		UnholyRuneCost: rc.UnholyRuneCost,
		DeathRuneCost:  rc.DeathRuneCost,
		RunicPowerCost: rc.RunicPowerCost,
		RunicPowerGain: rc.RunicPowerGain,
		Refundable:     rc.Refundable,
	}
}

func (rc *RuneCostImpl) MeetsRequirement(_ *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()

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
	} else {
		rc.spendRefundableRunicPowerCost(sim, spell)
	}
}

func (rc *RuneCostImpl) spendRefundableRunicPowerCost(sim *Simulation, spell *Spell) {
	if rc.Refundable && rc.RunicPowerCost > 0 {
		refundCost := TernaryFloat64(rc.RefundCost > 0, rc.RefundCost, rc.RunicPowerCost*0.1)
		spell.Unit.spendRunicPower(sim, refundCost, spell.RunicPowerMetrics())
	}
}

func (spell *Spell) SpendRefundableCost(sim *Simulation, result *SpellResult) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendRefundableCost(sim, spell, result)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertBloodRune(sim *Simulation, spell *Spell, landed bool) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !landed {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		rc.spendRefundableRunicPowerCost(sim, spell)
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	for _, slot := range slots {
		if slot == 0 || slot == 1 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
			changeType |= ConvertToDeath
		}
	}

	if rc.RunicPowerGain > 0 {
		spell.Unit.AddRunicPower(sim, rc.RunicPowerGain, spell.RunicPowerMetrics())
	}

	spell.Unit.maybeFireChange(sim, changeType)
}

func (spell *Spell) SpendRefundableCostAndConvertBloodRune(sim *Simulation, landed bool) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendRefundableCostAndConvertBloodRune(sim, spell, landed)
}

func (rc *RuneCostImpl) spendCostAndConvertFrostOrUnholyRune(sim *Simulation, spell *Spell, landed bool, refundable bool) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if refundable && !landed {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		rc.spendRefundableRunicPowerCost(sim, spell)
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
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

func (spell *Spell) SpendRefundableCostAndConvertFrostOrUnholyRune(sim *Simulation, landed bool) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendCostAndConvertFrostOrUnholyRune(sim, spell, landed, true)
}

func (spell *Spell) SpendCostAndConvertFrostOrUnholyRune(sim *Simulation, landed bool) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendCostAndConvertFrostOrUnholyRune(sim, spell, landed, false)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertBloodOrFrostRune(sim *Simulation, spell *Spell, landed bool) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !landed {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		rc.spendRefundableRunicPowerCost(sim, spell)
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	for _, slot := range slots {
		if slot == 0 || slot == 1 {
			spell.Unit.ConvertToDeath(sim, slot, NeverExpires)
			changeType |= ConvertToDeath
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

func (spell *Spell) SpendRefundableCostAndConvertBloodOrFrostRune(sim *Simulation, landed bool) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendRefundableCostAndConvertBloodOrFrostRune(sim, spell, landed)
}

func (rc *RuneCostImpl) spendRefundableCostAndConvertFrostRune(sim *Simulation, spell *Spell, landed bool) {
	cost := RuneCost(spell.CurCast.Cost) // cost was already optimized in RuneSpell.Cast
	if cost == 0 {
		return // it was free this time. we don't care
	}
	if !landed {
		// misses just don't get spent as a way to avoid having to cancel regeneration PAs
		// only spend RP
		rc.spendRefundableRunicPowerCost(sim, spell)
		return
	}

	changeType, slots := spell.Unit.spendRuneCost(sim, spell, cost)
	for _, slot := range slots {
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

func (spell *Spell) SpendRefundableCostAndConvertFrostRune(sim *Simulation, landed bool) {
	spell.Cost.ResourceCostImpl.(*RuneCostImpl).spendRefundableCostAndConvertFrostRune(sim, spell, landed)
}

func (rc *RuneCostImpl) IssueRefund(_ *Simulation, _ *Spell) {
	// Instead of issuing refunds we just don't charge the cost of spells which
	// miss; this is better for perf since we'd have to cancel the regen actions.
}

func (spell *Spell) RuneCostImpl() *RuneCostImpl {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl)
}

func (spell *Spell) RunicPowerMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl).runicPowerMetrics
}

func (spell *Spell) BloodRuneMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl).bloodRuneMetrics
}

func (spell *Spell) FrostRuneMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl).frostRuneMetrics
}

func (spell *Spell) UnholyRuneMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl).unholyRuneMetrics
}

func (spell *Spell) DeathRuneMetrics() *ResourceMetrics {
	return spell.Cost.ResourceCostImpl.(*RuneCostImpl).deathRuneMetrics
}

func (rp *runicPowerBar) NewRunicPowerMetrics(action ActionID) *ResourceMetrics {
	metric, ok := rp.spellRunicPowerMetrics[action]
	if !ok {
		metric = rp.character.newRunicPowerMetrics(action)
		rp.spellRunicPowerMetrics[action] = metric
	}

	return metric
}

func (rp *runicPowerBar) NewBloodRuneMetrics(action ActionID) *ResourceMetrics {
	metric, ok := rp.spellBloodRuneMetrics[action]
	if !ok {
		metric = rp.character.newBloodRuneMetrics(action)
		rp.spellBloodRuneMetrics[action] = metric
	}

	return metric
}

func (rp *runicPowerBar) NewFrostRuneMetrics(action ActionID) *ResourceMetrics {
	metric, ok := rp.spellFrostRuneMetrics[action]
	if !ok {
		metric = rp.character.newFrostRuneMetrics(action)
		rp.spellFrostRuneMetrics[action] = metric
	}

	return metric
}

func (rp *runicPowerBar) NewUnholyRuneMetrics(action ActionID) *ResourceMetrics {
	metric, ok := rp.spellUnholyRuneMetrics[action]
	if !ok {
		metric = rp.character.newUnholyRuneMetrics(action)
		rp.spellUnholyRuneMetrics[action] = metric
	}

	return metric
}

func (rp *runicPowerBar) NewDeathRuneMetrics(action ActionID) *ResourceMetrics {
	metric, ok := rp.spellDeathRuneMetrics[action]
	if !ok {
		metric = rp.character.newDeathRuneMetrics(action)
		rp.spellDeathRuneMetrics[action] = metric
	}

	return metric
}
