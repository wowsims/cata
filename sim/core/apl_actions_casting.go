package core

import (
	"fmt"

	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

type APLActionCastSpell struct {
	defaultAPLActionImpl
	spell  *Spell
	target UnitReference
}

func (rot *APLRotation) newActionCastSpell(config *proto.APLActionCastSpell) APLActionImpl {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	target := rot.GetTargetUnit(config.Target)
	if target.Get() == nil {
		return nil
	}
	return &APLActionCastSpell{
		spell:  spell,
		target: target,
	}
}
func (action *APLActionCastSpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCastOrQueue(sim, action.target.Get()) && (!action.spell.Flags.Matches(SpellFlagMCD) || action.spell.Flags.Matches(SpellFlagReactive) || action.spell.Unit.GCD.IsReady(sim) || action.spell.Unit.Rotation.inSequence)
}
func (action *APLActionCastSpell) Execute(sim *Simulation) {
	action.spell.CastOrQueue(sim, action.target.Get())
}
func (action *APLActionCastSpell) String() string {
	return fmt.Sprintf("Cast Spell(%s)", action.spell.ActionID)
}

type APLActionCastFriendlySpell struct {
	defaultAPLActionImpl
	spell  *Spell
	target UnitReference
}

func (rot *APLRotation) newActionCastFriendlySpell(config *proto.APLActionCastFriendlySpell) APLActionImpl {
	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	target := rot.GetTargetUnit(config.Target)
	if target.Get() == nil {
		return nil
	}
	return &APLActionCastFriendlySpell{
		spell:  spell,
		target: target,
	}
}
func (action *APLActionCastFriendlySpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCastOrQueue(sim, action.target.Get()) && (!action.spell.Flags.Matches(SpellFlagMCD) || action.spell.Flags.Matches(SpellFlagReactive) || action.spell.Unit.GCD.IsReady(sim) || action.spell.Unit.Rotation.inSequence)
}
func (action *APLActionCastFriendlySpell) Execute(sim *Simulation) {
	action.spell.CastOrQueue(sim, action.target.Get())
}
func (action *APLActionCastFriendlySpell) String() string {
	return fmt.Sprintf("Cast Friendly Spell(%s)", action.spell.ActionID)
}

type APLActionChannelSpell struct {
	defaultAPLActionImpl
	spell       *Spell
	target      UnitReference
	interruptIf APLValue
	allowRecast bool
}

func (rot *APLRotation) newActionChannelSpell(config *proto.APLActionChannelSpell) APLActionImpl {
	interruptIf := rot.coerceTo(rot.newAPLValue(config.InterruptIf), proto.APLValueType_ValueTypeBool)
	if interruptIf == nil {
		return rot.newActionCastSpell(&proto.APLActionCastSpell{
			SpellId: config.SpellId,
			Target:  config.Target,
		})
	}

	spell := rot.GetAPLSpell(config.SpellId)
	if spell == nil {
		return nil
	}
	if !spell.Flags.Matches(SpellFlagChanneled) {
		return nil
	}

	target := rot.GetTargetUnit(config.Target)
	if target.Get() == nil {
		return nil
	}

	return &APLActionChannelSpell{
		spell:       spell,
		target:      target,
		interruptIf: interruptIf,
		allowRecast: config.AllowRecast,
	}
}
func (action *APLActionChannelSpell) GetAPLValues() []APLValue {
	return []APLValue{action.interruptIf}
}
func (action *APLActionChannelSpell) IsReady(sim *Simulation) bool {
	return action.spell.CanCastOrQueue(sim, action.target.Get())
}
func (action *APLActionChannelSpell) Execute(sim *Simulation) {
	action.spell.CastOrQueue(sim, action.target.Get())
	action.spell.Unit.Rotation.interruptChannelIf = action.interruptIf
	action.spell.Unit.Rotation.allowChannelRecastOnInterrupt = action.allowRecast
}
func (action *APLActionChannelSpell) String() string {
	return fmt.Sprintf("Channel Spell(%s, interruptIf=%s)", action.spell.ActionID, action.interruptIf)
}

type APLActionMultidot struct {
	defaultAPLActionImpl
	spell      *Spell
	maxDots    int32
	maxOverlap APLValue

	nextTarget *Unit
}

func (rot *APLRotation) newActionMultidot(config *proto.APLActionMultidot) APLActionImpl {
	unit := rot.unit

	spell := rot.GetAPLMultidotSpell(config.SpellId)
	if spell == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"}, &proto.UUID{Value: ""})
	}

	maxDots := config.MaxDots
	numTargets := unit.Env.GetNumTargets()
	if spell.Flags.Matches(SpellFlagHelpful) {
		numTargets = int32(len(unit.Env.Raid.AllPlayerUnits))
	}
	if numTargets < maxDots {
		rot.ValidationMessage(proto.LogLevel_Warning, "Encounter only has %d targets. Using that for Max Dots instead of %d", numTargets, maxDots)
		maxDots = numTargets
	}

	return &APLActionMultidot{
		spell:      spell,
		maxDots:    maxDots,
		maxOverlap: maxOverlap,
	}
}
func (action *APLActionMultidot) GetAPLValues() []APLValue {
	return []APLValue{action.maxOverlap}
}
func (action *APLActionMultidot) Reset(*Simulation) {
	action.nextTarget = nil
}
func (action *APLActionMultidot) IsReady(sim *Simulation) bool {
	maxOverlap := action.maxOverlap.GetDuration(sim)

	if action.spell.Flags.Matches(SpellFlagHelpful) {
		for i := int32(0); i < action.maxDots; i++ {
			target := sim.Raid.AllPlayerUnits[i]
			dot := action.spell.Dot(target)
			if (!dot.IsActive() || dot.RemainingDuration(sim) < maxOverlap) && action.spell.CanCastOrQueue(sim, target) {
				action.nextTarget = target
				return true
			}
		}
	} else {
		for i := int32(0); i < action.maxDots; i++ {
			target := sim.Encounter.TargetUnits[i]
			dot := action.spell.Dot(target)
			if (!dot.IsActive() || dot.RemainingDuration(sim) < maxOverlap) && action.spell.CanCastOrQueue(sim, target) {
				action.nextTarget = target
				return true
			}
		}
	}
	return false
}
func (action *APLActionMultidot) Execute(sim *Simulation) {
	action.spell.CastOrQueue(sim, action.nextTarget)
}
func (action *APLActionMultidot) String() string {
	return fmt.Sprintf("Multidot(%s)", action.spell.ActionID)
}

type APLActionMultishield struct {
	defaultAPLActionImpl
	spell      *Spell
	maxShields int32
	maxOverlap APLValue

	nextTarget *Unit
}

func (rot *APLRotation) newActionMultishield(config *proto.APLActionMultishield) APLActionImpl {
	unit := rot.unit

	spell := rot.GetAPLMultishieldSpell(config.SpellId)
	if spell == nil {
		return nil
	}

	maxOverlap := rot.coerceTo(rot.newAPLValue(config.MaxOverlap), proto.APLValueType_ValueTypeDuration)
	if maxOverlap == nil {
		maxOverlap = rot.newValueConst(&proto.APLValueConst{Val: "0ms"}, &proto.UUID{Value: ""})
	}

	maxShields := config.MaxShields
	numTargets := int32(len(unit.Env.Raid.AllPlayerUnits))
	if numTargets < maxShields {
		rot.ValidationMessage(proto.LogLevel_Warning, "Encounter only has %d targets. Using that for Max Shields instead of %d", numTargets, maxShields)
		maxShields = numTargets
	}

	return &APLActionMultishield{
		spell:      spell,
		maxShields: maxShields,
		maxOverlap: maxOverlap,
	}
}
func (action *APLActionMultishield) GetAPLValues() []APLValue {
	return []APLValue{action.maxOverlap}
}
func (action *APLActionMultishield) Reset(*Simulation) {
	action.nextTarget = nil
}
func (action *APLActionMultishield) IsReady(sim *Simulation) bool {
	maxOverlap := action.maxOverlap.GetDuration(sim)

	for i := int32(0); i < action.maxShields; i++ {
		target := sim.Raid.AllPlayerUnits[i]
		shield := action.spell.Shield(target)
		if (!shield.IsActive() || shield.RemainingDuration(sim) < maxOverlap) && action.spell.CanCastOrQueue(sim, target) {
			action.nextTarget = target
			return true
		}
	}
	return false
}
func (action *APLActionMultishield) Execute(sim *Simulation) {
	action.spell.CastOrQueue(sim, action.nextTarget)
}
func (action *APLActionMultishield) String() string {
	return fmt.Sprintf("Multishield(%s)", action.spell.ActionID)
}

type APLActionCastAllStatBuffCooldowns struct {
	defaultAPLActionImpl
	character *Character

	statTypesToMatch []stats.Stat

	allSubactions   []*APLActionCastSpell
	readySubactions []*APLActionCastSpell
}

func (rot *APLRotation) newActionCastAllStatBuffCooldowns(config *proto.APLActionCastAllStatBuffCooldowns) APLActionImpl {
	unit := rot.unit
	actionImpl := &APLActionCastAllStatBuffCooldowns{
		character:        unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter(),
		statTypesToMatch: stats.IntTupleToStatsList(config.StatType1, config.StatType2, config.StatType3),
	}

	unit.Env.RegisterPostFinalizeEffect(func() {
		// This needs to happen after the rotation is finalized so that
		// all manually casted MCDs are removed from the cooldown list
		// before it is filtered for the desired buff types.
		actionImpl.processMajorCooldowns()
	})

	return actionImpl
}
func (action *APLActionCastAllStatBuffCooldowns) processMajorCooldowns() {
	matchingSpells := action.character.GetMatchingStatBuffSpells(action.statTypesToMatch)
	action.allSubactions = MapSlice(matchingSpells, func(buffSpell *Spell) *APLActionCastSpell {
		return &APLActionCastSpell{
			spell:  buffSpell,
			target: action.character.Rotation.GetTargetUnit(nil),
		}
	})

	action.character.Env.RegisterPostFinalizeEffect(func() {
		// We again need a delayed evaluation here so that other instances of
		// this action within the rotation can assemble their spell lists
		// before the filtered spells are irreversibly removed from the MCD
		// manager.
		for _, buffSpell := range matchingSpells {
			action.character.removeInitialMajorCooldown(buffSpell.ActionID)
		}
	})
}
func (action *APLActionCastAllStatBuffCooldowns) IsReady(sim *Simulation) bool {
	action.readySubactions = FilterSlice(action.allSubactions, func(subaction *APLActionCastSpell) bool {
		return subaction.IsReady(sim)
	})

	return Ternary(action.character.Rotation.inSequence, len(action.readySubactions) == len(action.allSubactions), len(action.readySubactions) > 0)
}
func (action *APLActionCastAllStatBuffCooldowns) Execute(sim *Simulation) {
	actionSetToUse := Ternary(sim.CurrentTime < 0, action.allSubactions, action.readySubactions)

	for _, subaction := range actionSetToUse {
		subaction.Execute(sim)
	}
}
func (action *APLActionCastAllStatBuffCooldowns) String() string {
	return fmt.Sprintf("CastAllBuffCooldownsFor(%s)", StringFromStatTypes(action.statTypesToMatch))
}
func (action *APLActionCastAllStatBuffCooldowns) PostFinalize(rot *APLRotation) {
	if len(action.allSubactions) == 0 {
		rot.ValidationMessage(proto.LogLevel_Warning, "%s will not cast any spells! There are either no major cooldowns buffing the specified stat type(s), or all of them are manually cast in the APL.", action)
	} else {
		actionIDs := MapSlice(action.allSubactions, func(subaction *APLActionCastSpell) ActionID {
			return subaction.spell.ActionID
		})

		rot.ValidationMessage(proto.LogLevel_Warning, "%s will cast the following spells: %s", action, StringFromActionIDs(actionIDs))
	}
}

type APLActionAutocastOtherCooldowns struct {
	defaultAPLActionImpl
	character *Character

	nextReadyMCD *MajorCooldown
}

func (rot *APLRotation) newActionAutocastOtherCooldowns(config *proto.APLActionAutocastOtherCooldowns) APLActionImpl {
	unit := rot.unit
	return &APLActionAutocastOtherCooldowns{
		character: unit.Env.Raid.GetPlayerFromUnit(unit).GetCharacter(),
	}
}
func (action *APLActionAutocastOtherCooldowns) Reset(*Simulation) {
	action.nextReadyMCD = nil
}
func (action *APLActionAutocastOtherCooldowns) IsReady(sim *Simulation) bool {
	action.nextReadyMCD = action.character.getFirstReadyMCD(sim)

	// Explicitly check for GCD because MCDs are usually desired to be cast immediately
	// before the next spell, rather than immediately after the previous spell. This is
	// true even for MCDs which do not require the GCD. The one exception to this rule
	// is Engineering explosives, which should instead be cast right *after* incurring
	// a GCD, since they are off-GCD but have small cast times.
	return (action.nextReadyMCD != nil) && ((action.character.GCD.IsReady(sim) != action.nextReadyMCD.Type.Matches(CooldownTypeExplosive)) || action.nextReadyMCD.Spell.Flags.Matches(SpellFlagReactive))
}
func (action *APLActionAutocastOtherCooldowns) Execute(sim *Simulation) {
	action.nextReadyMCD.tryActivateHelper(sim, action.character)
	action.character.UpdateMajorCooldowns()
}
func (action *APLActionAutocastOtherCooldowns) String() string {
	return "Autocast Other Cooldowns"
}
func (action *APLActionAutocastOtherCooldowns) PostFinalize(rot *APLRotation) {
	if len(action.character.initialMajorCooldowns) == 0 {
		rot.ValidationMessage(proto.LogLevel_Warning, "%s will not cast any spells! There are either no major cooldowns configured for this character, or all of them are manually cast in the APL.", action)
	} else {
		actionIDs := MapSlice(action.character.initialMajorCooldowns, func(mcd MajorCooldown) ActionID {
			return mcd.Spell.ActionID
		})

		rot.ValidationMessage(proto.LogLevel_Warning, "%s will cast the following spells: %s", action, StringFromActionIDs(actionIDs))
	}
}
