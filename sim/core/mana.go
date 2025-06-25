package core

import (
	"fmt"
	"math"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

const ThreatPerManaGained = 0.5

type manaBar struct {
	unit     *Unit
	BaseMana float64

	currentMana float64

	manaRegenMultiplier float64
	hasteEffectsRegen   bool

	manaCombatMetrics    *ResourceMetrics
	manaNotCombatMetrics *ResourceMetrics
	JowManaMetrics       *ResourceMetrics
	VtManaMetrics        *ResourceMetrics
	JowiseManaMetrics    *ResourceMetrics

	ReplenishmentAura *Aura

	// For keeping track of OOM status.
	waitingForMana          float64
	waitingForManaStartTime time.Duration
}

// EnableManaBar will setup caster stat dependencies (int->mana and int->spellcrit)
// as well as enable the mana gain action to regenerate mana.
// It will then enable mana gain metrics for reporting.
func (character *Character) EnableManaBar() {
	character.EnableManaBarWithModifier(1.0)
	character.Unit.SetCurrentPowerBar(ManaBar)
}

func (character *Character) EnableManaBarWithModifier(modifier float64) {

	// Starting with cataclysm you get mp5 equal 5% of your base mana
	character.AddStat(stats.MP5, character.baseStats[stats.Mana]*0.05)

	if character.Unit.Type == PlayerUnit {
		// Pets might have different scaling so let them handle their scaling
		character.AddStatDependency(stats.Intellect, stats.SpellCritPercent, CritPerIntMaxLevel[character.Class])

		// Starting with cataclysm 1 intellect now provides 1 spell power
		character.AddStatDependency(stats.Intellect, stats.SpellPower, 1.0)

		// first 10 int should not count so remove them
		character.AddStat(stats.SpellPower, -10)
	}

	// Not a real spell, just holds metrics from mana gain threat.
	character.RegisterSpell(SpellConfig{
		ActionID: ActionID{OtherID: proto.OtherAction_OtherActionManaGain},
	})

	character.manaCombatMetrics = character.NewManaMetrics(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen, Tag: 1})
	character.manaNotCombatMetrics = character.NewManaMetrics(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen, Tag: 2})

	character.BaseMana = character.GetBaseStats()[stats.Mana]
	character.Unit.manaBar.unit = &character.Unit
	character.Unit.manaBar.manaRegenMultiplier = 1.0
}

func (unit *Unit) HasManaBar() bool {
	return unit.manaBar.unit != nil
}

// Gets the Maxiumum mana including bonus and temporary affects that would increase your mana pool.
func (unit *Unit) MaxMana() float64 {
	return unit.stats[stats.Mana]
}
func (unit *Unit) CurrentMana() float64 {
	return unit.currentMana
}
func (unit *Unit) CurrentManaPercent() float64 {
	return unit.CurrentMana() / unit.MaxMana()
}

func (unit *Unit) AddMana(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to add negative mana!")
	}

	oldMana := unit.CurrentMana()
	newMana := min(oldMana+amount, unit.MaxMana())

	metrics.AddEvent(amount, newMana-oldMana)

	if sim.Log != nil {
		unit.Log(sim, "Gained %0.3f mana from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, oldMana, newMana, unit.MaxMana())
	}

	unit.currentMana = newMana
	unit.Metrics.ManaGained += newMana - oldMana
}

func (unit *Unit) SpendMana(sim *Simulation, amount float64, metrics *ResourceMetrics) {
	if amount < 0 {
		panic("Trying to spend negative mana!")
	}

	newMana := unit.CurrentMana() - amount
	metrics.AddEvent(-amount, -amount)

	if sim.Log != nil {
		unit.Log(sim, "Spent %0.3f mana from %s (%0.3f --> %0.3f) of %0.0f total.", amount, metrics.ActionID, unit.CurrentMana(), newMana, unit.MaxMana())
	}

	unit.currentMana = newMana
	unit.Metrics.ManaSpent += amount
}

func (mb *manaBar) doneIteration(sim *Simulation) {
	if mb.unit == nil {
		return
	}

	if mb.waitingForMana != 0 {
		mb.unit.Metrics.AddOOMTime(sim, sim.CurrentTime-mb.waitingForManaStartTime)
	}

	manaGainSpell := mb.unit.GetSpell(ActionID{OtherID: proto.OtherAction_OtherActionManaGain})

	for _, resourceMetrics := range mb.unit.Metrics.resources {
		if resourceMetrics.Type != proto.ResourceType_ResourceTypeMana {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{OtherID: proto.OtherAction_OtherActionManaRegen}) {
			continue
		}
		if resourceMetrics.ActionID.SameActionIgnoreTag(ActionID{SpellID: 34917}) {
			// Vampiric Touch mana threat goes to the priest, so it's handled in the priest code.
			continue
		}
		if resourceMetrics.ActualGainForCurrentIteration() <= 0 {
			continue
		}

		manaGainSpell.SpellMetrics[0].Casts += resourceMetrics.EventsForCurrentIteration()
		manaGainSpell.ApplyAOEThreatIgnoreMultipliers(resourceMetrics.ActualGainForCurrentIteration() * ThreatPerManaGained)
	}
}

// Returns the rate of mana regen per second from mp5.
func (unit *Unit) MP5ManaRegenPerSecond() float64 {
	return unit.stats[stats.MP5] / 5.0
}

// Returns the rate of mana regen per second from spirit.
func (unit *Unit) SpiritManaRegenPerSecond() float64 {
	return 0.001 + unit.stats[stats.Spirit]*math.Sqrt(unit.stats[stats.Intellect])*0.003345
}

// Returns the rate of mana regen per second, assuming this unit is
// considered to be casting.
func (unit *Unit) ManaRegenPerSecondWhileCombat() float64 {
	regenRate := unit.MP5ManaRegenPerSecond()

	if unit.manaBar.hasteEffectsRegen {
		regenRate *= unit.TotalSpellHasteMultiplier()
	}

	spiritRegenRate := 0.0
	if unit.PseudoStats.SpiritRegenRateCombat != 0 || unit.PseudoStats.ForceFullSpiritRegen {
		spiritRegenRate = unit.SpiritManaRegenPerSecond() * unit.PseudoStats.SpiritRegenMultiplier
		if !unit.PseudoStats.ForceFullSpiritRegen {
			spiritRegenRate *= unit.PseudoStats.SpiritRegenRateCombat
		}
	}
	regenRate += spiritRegenRate

	regenRate *= unit.manaRegenMultiplier

	return regenRate
}

// Returns the rate of mana regen per second, assuming this unit is
// considered to be not casting.
func (unit *Unit) ManaRegenPerSecondWhileNotCombat() float64 {
	regenRate := unit.MP5ManaRegenPerSecond()

	if unit.manaBar.hasteEffectsRegen {
		regenRate *= unit.TotalSpellHasteMultiplier()
	}

	regenRate += unit.SpiritManaRegenPerSecond() * unit.PseudoStats.SpiritRegenMultiplier

	regenRate *= unit.manaRegenMultiplier

	return regenRate
}

func (unit *Unit) UpdateManaRegenRates() {
	unit.manaTickWhileCombat = unit.ManaRegenPerSecondWhileCombat() * 2
	unit.manaTickWhileNotCombat = unit.ManaRegenPerSecondWhileNotCombat() * 2
}

func (unit *Unit) MultiplyManaRegenSpeed(sim *Simulation, multiplier float64) {
	unit.manaRegenMultiplier *= multiplier
	unit.UpdateManaRegenRates()
}

func (unit *Unit) HasteEffectsManaRegen() {
	unit.manaBar.hasteEffectsRegen = true
}

// Applies 1 'tick' of mana regen, which worth 2s of regeneration based on mp5/int/spirit/etc.
func (unit *Unit) ManaTick(sim *Simulation) {
	if sim.CurrentTime > 0 {
		regen := unit.manaTickWhileCombat
		unit.AddMana(sim, max(0, regen), unit.manaCombatMetrics)
	} else {
		regen := unit.manaTickWhileNotCombat
		unit.AddMana(sim, max(0, regen), unit.manaNotCombatMetrics)
	}
}

// Returns the amount of time this Unit would need to wait in order to reach
// the desired amount of mana, via mana regen.
//
// Assumes that desiredMana > currentMana. Calculation assumes the Unit
// will not take any actions during this period that would reset the 5-second rule.
func (unit *Unit) TimeUntilManaRegen(desiredMana float64) time.Duration {
	// +1 at the end is to deal with floating point math rounding errors.
	manaNeeded := desiredMana - unit.CurrentMana()
	regenTime := NeverExpires

	regenWhileCasting := unit.ManaRegenPerSecondWhileCombat()
	if regenWhileCasting != 0 {
		regenTime = DurationFromSeconds(manaNeeded/regenWhileCasting) + 1
	}

	// TODO: this needs to have access to the sim to see current time vs unit.PseudoStats.FiveSecondRule.
	//  it is possible that we have been waiting.
	//  In practice this function is always used right after a previous cast so no big deal for now.
	if regenTime > time.Second*5 {
		regenTime = time.Second * 5
		manaNeeded -= regenWhileCasting * 5
		// now we move into spirit based regen.
		regenTime += DurationFromSeconds(manaNeeded / unit.ManaRegenPerSecondWhileNotCombat())
	}

	return regenTime
}

func (sim *Simulation) initManaTickAction() {
	var unitsWithManaBars []*Unit

	for _, party := range sim.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			if character.HasManaBar() {
				unitsWithManaBars = append(unitsWithManaBars, &player.GetCharacter().Unit)
			}

			for _, petAgent := range character.PetAgents {
				if petAgent.GetPet().HasManaBar() {
					unitsWithManaBars = append(unitsWithManaBars, &petAgent.GetCharacter().Unit)
				}
			}
		}
	}

	if len(unitsWithManaBars) == 0 {
		return
	}

	interval := time.Second * 2
	pa := &PendingAction{
		NextActionAt: sim.Environment.PrepullStartTime() + interval,
		Priority:     ActionPriorityRegen,
	}
	pa.OnAction = func(sim *Simulation) {
		for _, unit := range unitsWithManaBars {
			if unit.IsEnabled() {
				unit.ManaTick(sim)
			}
		}

		pa.NextActionAt = sim.CurrentTime + interval
		sim.AddPendingAction(pa)
	}
	sim.AddPendingAction(pa)
}

func (mb *manaBar) reset() {
	if mb.unit == nil {
		return
	}

	mb.currentMana = mb.unit.MaxMana()

	mb.manaRegenMultiplier = 1.0

	mb.waitingForMana = 0
	mb.waitingForManaStartTime = 0
}

func (mb *manaBar) IsOOM() bool {
	return mb.waitingForMana != 0
}
func (mb *manaBar) StartOOMEvent(sim *Simulation, requiredMana float64) {
	mb.waitingForManaStartTime = sim.CurrentTime
	mb.waitingForMana = requiredMana
	mb.unit.Metrics.MarkOOM(sim)
}
func (mb *manaBar) EndOOMEvent(sim *Simulation) {
	eventDuration := sim.CurrentTime - mb.waitingForManaStartTime
	mb.unit.Metrics.AddOOMTime(sim, eventDuration)
	mb.waitingForManaStartTime = 0
	mb.waitingForMana = 0
}

func (unit *Unit) HasteEffectsRegen() {
	unit.manaBar.hasteEffectsRegen = true
}

type ManaCostOptions struct {
	BaseCostPercent float64 // The cost of the spell as a percentage (0-100) of the unit's base mana.
	FlatCost        int32   // Alternative to BaseCostPercent for giving a flat value.
	PercentModifier float64 // Will default to 1. PercentModifier stored as a float i.e. 40% reduction is (0.6 multiplier) to the base cost
}
type ManaCost struct {
	ResourceMetrics *ResourceMetrics
}

func newManaCost(spell *Spell, options ManaCostOptions) *SpellCost {
	return &SpellCost{
		spell:           spell,
		BaseCost:        TernaryInt32(options.FlatCost > 0, options.FlatCost, int32(options.BaseCostPercent*spell.Unit.BaseMana)/100),
		PercentModifier: TernaryFloat64(options.PercentModifier == 0, 1, options.PercentModifier),
		ResourceCostImpl: &ManaCost{
			ResourceMetrics: spell.Unit.NewManaMetrics(spell.ActionID),
		},
	}
}

func (mc *ManaCost) MeetsRequirement(sim *Simulation, spell *Spell) bool {
	spell.CurCast.Cost = spell.Cost.GetCurrentCost()
	meetsRequirement := spell.Unit.CurrentMana() >= spell.CurCast.Cost

	if spell.CurCast.Cost > 0 {
		if meetsRequirement {
			if spell.Unit.IsOOM() {
				spell.Unit.EndOOMEvent(sim)
			}
		} else {
			if spell.Unit.IsOOM() {
				// Continuation of OOM event.
				spell.Unit.waitingForMana = min(spell.Unit.waitingForMana, spell.CurCast.Cost)
			} else {
				spell.Unit.StartOOMEvent(sim, spell.CurCast.Cost)
			}
		}
	}

	return meetsRequirement
}
func (mc *ManaCost) CostFailureReason(sim *Simulation, spell *Spell) string {
	return fmt.Sprintf("not enough mana (Current Mana = %0.03f, Mana Cost = %0.03f)", spell.Unit.CurrentMana(), spell.CurCast.Cost)
}
func (mc *ManaCost) SpendCost(sim *Simulation, spell *Spell) {
	if spell.CurCast.Cost > 0 {
		spell.Unit.SpendMana(sim, spell.CurCast.Cost, mc.ResourceMetrics)
	}
}
func (mc *ManaCost) IssueRefund(_ *Simulation, _ *Spell) {}
