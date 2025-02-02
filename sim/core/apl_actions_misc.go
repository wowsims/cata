package core

import (
	"fmt"
	"time"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type APLActionChangeTarget struct {
	defaultAPLActionImpl
	unit      *Unit
	newTarget UnitReference
}

func (rot *APLRotation) newActionChangeTarget(config *proto.APLActionChangeTarget) APLActionImpl {
	newTarget := rot.GetSourceUnit(config.NewTarget)
	if newTarget.Get() == nil {
		return nil
	}
	return &APLActionChangeTarget{
		newTarget: newTarget,
	}
}
func (action *APLActionChangeTarget) IsReady(sim *Simulation) bool {
	return action.unit.CurrentTarget != action.newTarget.Get()
}
func (action *APLActionChangeTarget) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.unit.Log(sim, "Changing target to %s", action.newTarget.Get().Label)
	}
	action.unit.CurrentTarget = action.newTarget.Get()
}
func (action *APLActionChangeTarget) String() string {
	return fmt.Sprintf("Change Target(%s)", action.newTarget.Get().Label)
}

type APLActionCancelAura struct {
	defaultAPLActionImpl
	aura *Aura
}

type APLActionActivateAura struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionCancelAura(config *proto.APLActionCancelAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionCancelAura{
		aura: aura.Get(),
	}
}

func (rot *APLRotation) newActionActivateAura(config *proto.APLActionActivateAura) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionActivateAura{
		aura: aura.Get(),
	}
}

func (action *APLActionCancelAura) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionCancelAura) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Cancelling aura %s", action.aura.ActionID)
	}
	action.aura.Deactivate(sim)
}
func (action *APLActionCancelAura) String() string {
	return fmt.Sprintf("Cancel Aura(%s)", action.aura.ActionID)
}

func (action *APLActionActivateAura) IsReady(sim *Simulation) bool {
	return (action.aura.Icd == nil) || action.aura.Icd.IsReady(sim)
}

func (action *APLActionActivateAura) Execute(sim *Simulation) {
	if !action.IsReady(sim) {
		if sim.Log != nil {
			action.aura.Unit.Log(sim, "Could not activate aura %s because it's not ready", action.aura.ActionID)
		}
		return
	}

	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Activating aura %s", action.aura.ActionID)
	}

	action.aura.Activate(sim)
	if action.aura.Icd != nil {
		action.aura.Icd.Use(sim)
	}
}

func (action *APLActionActivateAura) String() string {
	return fmt.Sprintf("Activate Aura(%s)", action.aura.ActionID)
}

type APLActionActivateAuraWithStacks struct {
	defaultAPLActionImpl
	aura      *Aura
	numStacks int32
}

func (rot *APLRotation) newActionActivateAuraWithStacks(config *proto.APLActionActivateAuraWithStacks) APLActionImpl {
	aura := rot.GetAPLAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	if aura.Get().MaxStacks == 0 {
		rot.ValidationMessage(proto.LogLevel_Warning, "%s is not a stackable aura", ProtoToActionID(config.AuraId))
		return nil
	}
	return &APLActionActivateAuraWithStacks{
		aura:      aura.Get(),
		numStacks: int32(min(config.NumStacks, aura.Get().MaxStacks)),
	}
}
func (action *APLActionActivateAuraWithStacks) IsReady(sim *Simulation) bool {
	return (action.aura.Icd == nil) || action.aura.Icd.IsReady(sim)
}
func (action *APLActionActivateAuraWithStacks) Execute(sim *Simulation) {
	if !action.IsReady(sim) {
		if sim.Log != nil {
			action.aura.Unit.Log(sim, "Could not activate aura %s (%d stacks) because it's not ready", action.aura.ActionID, action.numStacks)
		}
		return
	}
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Activating aura %s (%d stacks)", action.aura.ActionID, action.numStacks)
	}
	action.aura.Activate(sim)
	action.aura.SetStacks(sim, action.numStacks)
}
func (action *APLActionActivateAuraWithStacks) String() string {
	return fmt.Sprintf("Activate Aura(%s) Stacks(%d)", action.aura.ActionID, action.numStacks)
}

type APLActionActivateAllItemSwapStatBuffAuras struct {
	defaultAPLActionImpl
	character *Character

	statTypesToMatch []stats.Stat

	allSubactions []*APLActionActivateAura
}

func (rot *APLRotation) newActionActivateAllItemSwapStatBuffAuras(config *proto.APLActionActivateAllItemSwapStatBuffAuras) APLActionImpl {

	unit := rot.unit
	character := unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter()
	statTypesToMatch := stats.IntTupleToStatsList(config.StatType1, config.StatType2, config.StatType3)

	allSubactions := MapSlice(rot.GetAPLItemProcAuras(statTypesToMatch, 0, false, true, &proto.UUID{Value: ""}), func(statBuffAura *StatBuffAura) *APLActionActivateAura {
		return &APLActionActivateAura{
			aura: statBuffAura.Aura,
		}
	})

	return &APLActionActivateAllItemSwapStatBuffAuras{
		character:        character,
		statTypesToMatch: statTypesToMatch,
		allSubactions:    allSubactions,
	}
}

func (action *APLActionActivateAllItemSwapStatBuffAuras) IsReady(sim *Simulation) bool {
	return len(action.allSubactions) > 0
}

func (action *APLActionActivateAllItemSwapStatBuffAuras) Execute(sim *Simulation) {
	for _, subaction := range action.allSubactions {
		subaction.Execute(sim)
	}
}

func (action *APLActionActivateAllItemSwapStatBuffAuras) String() string {
	return fmt.Sprintf("ActivateAllItemSwapStatBuffAurasFor(%s)", StringFromStatTypes(action.statTypesToMatch))
}

func (action *APLActionActivateAllItemSwapStatBuffAuras) PostFinalize(rot *APLRotation) {
	if len(action.allSubactions) == 0 {
		rot.ValidationMessage(proto.LogLevel_Warning, "%s will not activate any Auras! There are no proc items buffing the specified stat type(s).", action)
	} else {
		actionIDs := MapSlice(action.allSubactions, func(subaction *APLActionActivateAura) ActionID {
			return subaction.aura.ActionID
		})

		rot.ValidationMessage(proto.LogLevel_Information, "%s will activate the following Auras: %s", action, StringFromActionIDs(actionIDs))
	}
}

type APLActionTriggerICD struct {
	defaultAPLActionImpl
	aura *Aura
}

func (rot *APLRotation) newActionTriggerICD(config *proto.APLActionTriggerICD) APLActionImpl {
	aura := rot.GetAPLICDAura(rot.GetSourceUnit(&proto.UnitReference{Type: proto.UnitReference_Self}), config.AuraId)
	if aura.Get() == nil {
		return nil
	}
	return &APLActionTriggerICD{
		aura: aura.Get(),
	}
}
func (action *APLActionTriggerICD) IsReady(sim *Simulation) bool {
	return action.aura.IsActive()
}
func (action *APLActionTriggerICD) Execute(sim *Simulation) {
	if sim.Log != nil {
		action.aura.Unit.Log(sim, "Triggering ICD %s", action.aura.ActionID)
	}
	action.aura.Icd.Use(sim)
}
func (action *APLActionTriggerICD) String() string {
	return fmt.Sprintf("Trigger ICD(%s)", action.aura.ActionID)
}

type APLActionItemSwap struct {
	defaultAPLActionImpl
	character *Character
	swapSet   proto.APLActionItemSwap_SwapSet
}

func (rot *APLRotation) newActionItemSwap(config *proto.APLActionItemSwap) APLActionImpl {
	if config.SwapSet == proto.APLActionItemSwap_Unknown {
		rot.ValidationMessage(proto.LogLevel_Warning, "Unknown item swap set")
		return nil
	}

	character := rot.unit.Env.Raid.GetPlayerFromUnit(rot.unit).GetCharacter()
	if !character.ItemSwap.IsEnabled() {
		if config.SwapSet != proto.APLActionItemSwap_Main {
			rot.ValidationMessage(proto.LogLevel_Warning, "No swap set configured in Settings.")
		}
		return nil
	}

	return &APLActionItemSwap{
		character: character,
		swapSet:   config.SwapSet,
	}
}
func (action *APLActionItemSwap) IsReady(sim *Simulation) bool {
	return (action.swapSet == proto.APLActionItemSwap_Main) == action.character.ItemSwap.IsSwapped()
}
func (action *APLActionItemSwap) Execute(sim *Simulation) {
	if action.character.ItemSwap.swapSet == action.swapSet {
		if sim.Log != nil {
			action.character.Log(sim, "Item Swap already set to %s", action.swapSet)
		}
	} else {
		if sim.Log != nil {
			action.character.Log(sim, "Item Swap to set %s", action.swapSet)
		}
	}

	action.character.ItemSwap.SwapItems(sim, action.swapSet, false)
}
func (action *APLActionItemSwap) String() string {
	return fmt.Sprintf("Item Swap(%s)", action.swapSet)
}

type APLActionMove struct {
	defaultAPLActionImpl
	unit      *Unit
	moveRange APLValue
}

func (rot *APLRotation) newActionMove(config *proto.APLActionMove) APLActionImpl {
	return &APLActionMove{
		unit:      rot.unit,
		moveRange: rot.newAPLValue(config.RangeFromTarget),
	}
}
func (action *APLActionMove) IsReady(sim *Simulation) bool {
	isPrepull := sim.CurrentTime < 0
	return !action.unit.Moving && (action.moveRange.GetFloat(sim) != action.unit.DistanceFromTarget || isPrepull) && action.unit.Hardcast.Expires < sim.CurrentTime
}
func (action *APLActionMove) Execute(sim *Simulation) {
	moveRange := action.moveRange.GetFloat(sim)
	if sim.Log != nil {
		action.unit.Log(sim, "[DEBUG] Moving to %.1f yards", moveRange)
	}

	action.unit.MoveTo(moveRange, sim)
}
func (action *APLActionMove) String() string {
	return fmt.Sprintf("Move(%s)", action.moveRange)
}

type APLActionCustomRotation struct {
	defaultAPLActionImpl
	unit  *Unit
	agent Agent

	lastExecutedAt time.Duration
}

func (rot *APLRotation) newActionCustomRotation(config *proto.APLActionCustomRotation) APLActionImpl {
	agent := rot.unit.Env.GetAgentFromUnit(rot.unit)
	if agent == nil {
		panic("Agent not found for custom rotation")
	}

	return &APLActionCustomRotation{
		unit:  rot.unit,
		agent: agent,
	}
}
func (action *APLActionCustomRotation) Reset(sim *Simulation) {
	action.lastExecutedAt = -1
}
func (action *APLActionCustomRotation) IsReady(sim *Simulation) bool {
	// Prevent infinite loops by only allowing this action to be performed once at each timestamp.
	return action.lastExecutedAt != sim.CurrentTime
}
func (action *APLActionCustomRotation) Execute(sim *Simulation) {
	action.lastExecutedAt = sim.CurrentTime
	action.agent.ExecuteCustomRotation(sim)
}
func (action *APLActionCustomRotation) String() string {
	return "Custom Rotation()"
}

type APLActionMoveDuration struct {
	defaultAPLActionImpl
	unit         *Unit
	moveDuration APLValue
}

func (rot *APLRotation) newActionMoveDuration(config *proto.APLActionMoveDuration) APLActionImpl {
	return &APLActionMoveDuration{
		unit:         rot.unit,
		moveDuration: rot.newAPLValue(config.Duration),
	}
}

func (action *APLActionMoveDuration) Execute(sim *Simulation) {
	action.unit.MoveDuration(action.moveDuration.GetDuration(sim), sim)
}

func (action *APLActionMoveDuration) IsReady(sim *Simulation) bool {

	// only alow us to move if we're not already moving or movement action is running out this step
	if action.unit.Moving && action.unit.movementAction.NextActionAt != sim.CurrentTime {
		return false
	}

	if action.moveDuration.GetDuration(sim) == time.Duration(0) {
		return false
	}

	// check current casting state
	return (action.unit.Hardcast.Expires < sim.CurrentTime || action.unit.Hardcast.CanMove) && action.unit.ChanneledDot == nil
}

func (action *APLActionMoveDuration) String() string {
	return "MoveDuration()"
}
