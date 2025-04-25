package feral

import (
	"math"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/druid"
)

func (cat *FeralDruid) OnGCDReady(sim *core.Simulation) {
	if !cat.usingHardcodedAPL {
		return
	}

	if !cat.GCD.IsReady(sim) {
		return
	}

	cat.bleedAura = cat.CurrentTarget.GetExclusiveEffectCategory(core.BleedEffectCategory).GetActiveAura()

	if cat.preRotationCleanup(sim) {
		valid := false
		nextAction := time.Duration(0)
		if cat.Rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget {
			valid, nextAction = cat.doRotation(sim)
		} else {
			valid, nextAction = cat.doAoeRotation(sim)
		}
		if valid {
			cat.postRotation(sim, nextAction)
		}
	}

	// Check for an opportunity to cancel Primal Madness if we just casted a spell.
	if !cat.GCD.IsReady(sim) && cat.PrimalMadnessAura.IsActive() && cat.Rotation.CancelPrimalMadness {
		// Determine cancellation threshold based on the expected Energy
		// loss when Primal Madness will naturally expire.
		energyThresh := cat.primalMadnessBonus

		// Apply a conservative correction to account for the cost of losing one final buffed Shred at the very
		// end of the TF or Zerk window due to an early cancellation.
		energyThresh -= core.TernaryFloat64(cat.BerserkAura.IsActive(), 0.5, 0.15) * cat.Shred.DefaultCast.Cost

		// Apply input delay realism to Energy measurement for a real player.
		energyThresh += cat.EnergyRegenPerSecond() * cat.ReactionTime.Seconds()

		if cat.CurrentEnergy() < energyThresh {
			cat.PrimalMadnessAura.Deactivate(sim)
		}
	}
}

func (cat *FeralDruid) NextRotationAction(sim *core.Simulation, kickAt time.Duration) {
	cat.nextActionAt = kickAt
	cat.WaitUntil(sim, kickAt)
}

func (cat *FeralDruid) shiftBearCat(sim *core.Simulation, powershift bool) bool {
	cat.waitingForTick = false

	// If we have just now decided to shift, then we do not execute the
	// shift immediately, but instead trigger an input delay for realism.
	if !cat.readyToShift {
		cat.readyToShift = true
		return false
	}
	cat.readyToShift = false

	toCat := !cat.InForm(druid.Cat)
	if powershift {
		toCat = !toCat
	}

	cat.lastShift = sim.CurrentTime
	if toCat {
		return cat.CatForm.Cast(sim, nil)
	} else {
		cat.BearForm.Cast(sim, nil)
		// Bundle Enrage if available
		if cat.Enrage.IsReady(sim) {
			cat.Enrage.Cast(sim, nil)
		}
		return true
	}
}

func (cat *FeralDruid) canBite(sim *core.Simulation, isExecutePhase bool) bool {
	if cat.tempSnapshotAura.IsActive() && isExecutePhase {
		return true
	}

	biteTime := core.TernaryDuration(cat.BerserkAura.IsActive(), cat.Rotation.BerserkBiteTime, cat.Rotation.BiteTime)

	if cat.SavageRoarAura.RemainingDuration(sim) < biteTime {
		return false
	}

	if isExecutePhase {
		return (cat.Rip.NewSnapshotPower > cat.Rip.CurrentSnapshotPower-0.001) || cat.BerserkAura.IsActive()
	}

	return cat.Rip.CurDot().RemainingDuration(sim) >= biteTime
}

func (cat *FeralDruid) berserkExpectedAt(sim *core.Simulation, futureTime time.Duration) bool {
	if cat.BerserkAura.IsActive() {
		return futureTime < cat.BerserkAura.ExpiresAt()
	}

	if !cat.Talents.Berserk {
		return false
	}

	if cat.Berserk.IsReady(sim) {
		return cat.TigersFuryAura.IsActive() || cat.tfExpectedBefore(sim, futureTime)
	}

	return futureTime > cat.Berserk.ReadyAt()
}

func (cat *FeralDruid) calcBuilderDpe(sim *core.Simulation) (float64, float64) {
	// Calculate current damage-per-Energy of Rake vs. Shred. Used to
	// determine whether Rake is worth casting when player stats change upon a
	// dynamic proc occurring
	shredDpc := cat.Shred.ExpectedInitialDamage(sim, cat.CurrentTarget)
	potentialRakeTicks := min(cat.Rake.CurDot().BaseTickCount, int32(sim.GetRemainingDuration()/time.Second*3))
	rakeDpc := cat.Rake.ExpectedInitialDamage(sim, cat.CurrentTarget) + cat.Rake.ExpectedTickDamage(sim, cat.CurrentTarget)*float64(potentialRakeTicks)
	return rakeDpc / cat.Rake.DefaultCast.Cost, shredDpc / cat.Shred.DefaultCast.Cost
}

func (cat *FeralDruid) calcRipEndThresh(sim *core.Simulation) time.Duration {
	// Use cached value when below 5 CP
	if cat.ComboPoints() < 5 {
		return cat.cachedRipEndThresh
	}

	// Calculate the minimum DoT duration at which a Rip cast will provide higher DPE than a Bite cast
	expectedBiteDPE := cat.FerociousBite.ExpectedInitialDamage(sim, cat.CurrentTarget) / cat.FerociousBite.DefaultCast.Cost
	expectedRipTickDPE := cat.Rip.ExpectedTickDamage(sim, cat.CurrentTarget) / cat.Rip.DefaultCast.Cost
	numTicksToBreakEven := 1 + int32(expectedBiteDPE/expectedRipTickDPE)

	if sim.Log != nil {
		cat.Log(sim, "Bite Break-Even Point = %d Rip ticks", numTicksToBreakEven)
	}

	ripDot := cat.Rip.CurDot()
	endThresh := time.Duration(numTicksToBreakEven) * ripDot.BaseTickLength

	// Store the result so we can keep using it even when not at 5 CP
	cat.cachedRipEndThresh = endThresh

	return endThresh
}

func (cat *FeralDruid) clipRoar(sim *core.Simulation, isExecutePhase bool) bool {
	ripDot := cat.Rip.CurDot()
	ripdotRemaining := ripDot.RemainingDuration(sim)
	simTimeRemaining := sim.GetRemainingDuration()

	if !ripDot.IsActive() || (simTimeRemaining-ripdotRemaining < cat.cachedRipEndThresh) {
		return false
	}

	// Project Rip end time assuming full Glyph of Shred extensions
	remainingExtensions := cat.maxRipTicks - ripDot.BaseTickCount
	ripDur := ripdotRemaining + time.Duration(remainingExtensions)*ripDot.BaseTickLength
	roarDur := cat.SavageRoarAura.RemainingDuration(sim)

	if roarDur > (ripDur + cat.Rotation.RipLeeway) {
		return false
	}

	if roarDur >= simTimeRemaining {
		return false
	}

	// Calculate when roar would end if casted now
	newRoarDur := cat.SavageRoarDurationTable[cat.ComboPoints()]

	// If a fresh Roar cast now would cover us to end of fight, then clip now for maximum CP efficiency.
	if newRoarDur >= simTimeRemaining {
		return true
	}

	// If waiting another GCD to build an additional CP would lower our total Roar casts for the fight, then force a wait.
	if newRoarDur+time.Second+core.TernaryDuration(cat.ComboPoints() < 5, time.Second*5, 0) >= simTimeRemaining {
		return false
	}

	// Clip as soon as we have enough CPs for the new roar to expire well
	// after the current rip
	if !isExecutePhase {
		return newRoarDur >= (ripDur + cat.Rotation.MinRoarOffset)
	}

	// Under Execute conditions, ignore the offset rule and instead optimize for as few Roar casts as possible.
	if cat.ComboPoints() < 5 {
		return false
	}

	minRoarsPossible := (simTimeRemaining - roarDur) / newRoarDur
	projectedRoarCasts := simTimeRemaining / newRoarDur
	return projectedRoarCasts == minRoarsPossible
}

func (cat *FeralDruid) tfExpectedBefore(sim *core.Simulation, futureTime time.Duration) bool {
	if !cat.TigersFury.IsReady(sim) {
		return cat.TigersFury.ReadyAt() < futureTime
	}
	if cat.BerserkAura.IsActive() {
		return cat.BerserkAura.ExpiresAt() < futureTime
	}
	return true
}

func (cat *FeralDruid) calcTfEnergyThresh(leewayTime time.Duration) float64 {
	delayTime := leewayTime + core.TernaryDuration(cat.ClearcastingAura.IsActive(), time.Second, 0) + core.TernaryDuration(cat.StampedeCatAura.IsActive() && (cat.Rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget), time.Second, 0)
	return 40.0 - delayTime.Seconds()*cat.EnergyRegenPerSecond()
}

func (cat *FeralDruid) TryTigersFury(sim *core.Simulation) {
	// Handle tigers fury
	if !cat.TigersFury.IsReady(sim) {
		return
	}

	gcdTimeToRdy := cat.GCD.TimeToReady(sim)
	leewayTime := max(gcdTimeToRdy, cat.ReactionTime)
	tfEnergyThresh := cat.calcTfEnergyThresh(leewayTime)
	tfNow := (cat.CurrentEnergy() < tfEnergyThresh) && !cat.BerserkAura.IsActive() && (!cat.T13Feral4pBonus.IsActive() || !cat.StampedeCatAura.IsActive() || (cat.Rotation.RotationType == proto.FeralDruid_Rotation_Aoe))

	if tfNow {
		cat.TigersFury.Cast(sim, nil)
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		cat.NextRotationAction(sim, sim.CurrentTime+leewayTime)
	}
}

func (cat *FeralDruid) TryBerserk(sim *core.Simulation) {
	// Berserk algorithm: time Berserk for just after a Tiger's Fury
	// *unless* we'll lose Berserk uptime by waiting for Tiger's Fury to
	// come off cooldown. The latter exception is necessary for
	// Lacerateweave rotation since TF timings can drift over time.
	simTimeRemain := sim.GetRemainingDuration()
	tfCdRemain := cat.TigersFury.TimeToReady(sim)
	waitForTf := cat.Talents.Berserk && (tfCdRemain <= cat.BerserkAura.Duration) && (tfCdRemain+cat.ReactionTime < simTimeRemain-cat.BerserkAura.Duration)
	isClearcast := cat.ClearcastingAura.IsActive()
	berserkNow := cat.Rotation.UseBerserk && cat.Berserk.IsReady(sim) && !waitForTf && !isClearcast

	if berserkNow {
		cat.Berserk.Cast(sim, nil)
		cat.UpdateMajorCooldowns()

		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		gcdTimeToRdy := cat.GCD.TimeToReady(sim)
		leewayTime := max(gcdTimeToRdy, cat.ReactionTime)
		cat.NextRotationAction(sim, sim.CurrentTime+leewayTime)
	}
}

func (cat *FeralDruid) preRotationCleanup(sim *core.Simulation) bool {
	if cat.BerserkAura.IsActive() {
		cat.berserkUsed = true
	}

	// If we previously decided to shift, then execute the shift now once
	// the input delay is over.
	if cat.readyToShift {
		cat.shiftBearCat(sim, false)

		// Reset swing timer from snek (or idol/weapon swap) when going into cat
		if cat.InForm(druid.Cat) && cat.Rotation.SnekWeave {
			if cat.AutoAttacks.NextAttackAt()-sim.CurrentTime > cat.AutoAttacks.MainhandSwingSpeed() {
				cat.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
			}
		}

		// Bundle a leave-weave with the Cat Form GCD if possible
		if cat.InForm(druid.Cat) && cat.Rotation.MeleeWeave {
			timeToMove := core.DurationFromSeconds((cat.CatCharge.MinRange+1-cat.DistanceFromTarget)/cat.GetMovementSpeed()) + cat.ReactionTime

			if cat.CatCharge.TimeToReady(sim) < timeToMove {
				cat.MoveTo(cat.CatCharge.MinRange+1, sim)
				cat.NextRotationAction(sim, sim.CurrentTime+cat.ReactionTime)
			}
		}

		// To prep for the above, pre-position to max melee range during Bear Form GCD
		if cat.InForm(druid.Bear) {
			cat.MoveTo(core.MaxMeleeRange-1, sim)
			cat.NextRotationAction(sim, sim.CurrentTime+cat.ReactionTime)
		}

		return false
	}

	return true
}

func (cat *FeralDruid) postRotation(sim *core.Simulation, nextAction time.Duration) {
	// Also schedule an action right at Energy cap to make sure we never
	// accidentally over-cap while waiting on other timers.
	timeToCap := core.DurationFromSeconds((cat.MaximumEnergy() - cat.CurrentEnergy()) / cat.EnergyRegenPerSecond())
	nextAction = min(nextAction, sim.CurrentTime+timeToCap)

	nextAction += cat.ReactionTime

	if nextAction <= sim.CurrentTime {
		panic("nextaction in the past")
	} else {
		cat.NextRotationAction(sim, nextAction)
	}
}

func (cat *FeralDruid) calcBleedRefreshTime(sim *core.Simulation, bleedSpell *druid.DruidSpell, bleedDot *core.Dot, isExecutePhase bool, isRip bool) time.Duration {
	if !bleedDot.IsActive() {
		return sim.CurrentTime - cat.ReactionTime
	}

	// If we're not gaining a stronger snapshot, then use the standard 1
	// tick refresh window.
	bleedEnd := bleedDot.ExpiresAt()
	standardRefreshTime := bleedEnd - bleedDot.BaseTickLength

	if !cat.tempSnapshotAura.IsActive() {
		return standardRefreshTime
	}

	// For Rip specifically, also bypass clipping calculations during Execute phase or
	// if CP count is too low for the calculation to be relevant.
	if isRip && (isExecutePhase || (cat.ComboPoints() < cat.Rotation.MinCombosForRip)) {
		return standardRefreshTime
	}

	// Likewise, if the existing buff will still be up at the start of the normal
	// window, then don't clip unnecessarily. For long buffs that cover a full bleed
	// duration, project "buffEnd" forward in time such that we block clips if we are
	// already maxing out the number of full durations we can snapshot.
	buffRemains := cat.tempSnapshotAura.RemainingDuration(sim)
	maxTickCount := core.TernaryInt32(isRip, cat.maxRipTicks, bleedDot.BaseTickCount)
	maxBleedDur := bleedDot.BaseTickLength * time.Duration(maxTickCount)
	numCastsCovered := buffRemains / maxBleedDur
	buffEnd := cat.tempSnapshotAura.ExpiresAt() - numCastsCovered*maxBleedDur

	if buffEnd > standardRefreshTime+cat.ReactionTime {
		return standardRefreshTime
	}

	// Potential clips for a buff snapshot should be done as late as possible
	latestPossibleSnapshot := buffEnd - cat.ReactionTime*time.Duration(2)
	numClippedTicks := (bleedEnd - latestPossibleSnapshot) / bleedDot.BaseTickLength
	targetClipTime := standardRefreshTime - numClippedTicks*bleedDot.BaseTickLength

	// Since the clip can cost us 30-35 Energy, we need to determine whether the damage gain is worth the
	// spend. First calculate the maximum number of buffed bleed ticks we can get out before the fight
	// ends.
	buffedTickCount := min(maxTickCount, int32((sim.Duration-targetClipTime)/bleedDot.BaseTickLength))

	// Perform a DPE comparison vs. Shred
	expectedDamageGain := (bleedSpell.NewSnapshotPower - bleedSpell.CurrentSnapshotPower) * float64(buffedTickCount)

	// For Rake specifically, we get 1 free "tick" immediately upon cast.
	if !isRip {
		expectedDamageGain += bleedSpell.NewSnapshotPower
	}

	energyEquivalent := expectedDamageGain / cat.Shred.ExpectedInitialDamage(sim, cat.CurrentTarget) * cat.Shred.DefaultCast.Cost

	// Finally, discount the effective Energy cost of the clip based on the number of clipped ticks.
	discountedRefreshCost := float64(numClippedTicks) / float64(maxTickCount) * bleedSpell.DefaultCast.Cost

	if sim.Log != nil {
		cat.Log(sim, "%s buff snapshot is worth %.1f Energy, discounted refresh cost is %.1f Energy.", bleedSpell.ShortName, energyEquivalent, discountedRefreshCost)
	}

	return core.TernaryDuration(energyEquivalent > discountedRefreshCost, targetClipTime, standardRefreshTime)
}

func (cat *FeralDruid) canMeleeWeave(sim *core.Simulation, regenRate float64, currentEnergy float64, isClearcast bool, upcomingTimers *PoolingActions) bool {
	if !cat.Rotation.MeleeWeave || !cat.CatCharge.IsReady(sim) || isClearcast || cat.BerserkAura.IsActive() {
		return false
	}

	// Estimate time to run out and charge back in
	runOutTime := core.DurationFromSeconds((cat.CatCharge.MinRange+1-cat.DistanceFromTarget)/cat.GetMovementSpeed()) + cat.ReactionTime
	chargeInTime := core.DurationFromSeconds((cat.CatCharge.MinRange+1)/80) + cat.ReactionTime
	weaveDuration := runOutTime + chargeInTime
	weaveEnergy := 100.0 - weaveDuration.Seconds()*regenRate

	if currentEnergy > weaveEnergy {
		return false
	}

	// Prioritize all timers over weaving
	weaveEnd := sim.CurrentTime + weaveDuration
	isPooling, nextRefresh := upcomingTimers.nextRefreshTime()

	if (isPooling && (nextRefresh < weaveEnd)) || cat.tfExpectedBefore(sim, weaveEnd) {
		return false
	}

	// Also add an end-of-fight condition to make sure we can spend down our Energy
	// post-weave before the encounter ends.
	energyToDump := currentEnergy + weaveDuration.Seconds()*regenRate
	timeToDump := core.DurationFromSeconds(math.Floor(energyToDump / cat.Shred.DefaultCast.Cost))
	return weaveEnd+timeToDump < sim.Duration
}

func (cat *FeralDruid) canBearWeave(sim *core.Simulation, furorCap float64, regenRate float64, currentEnergy float64, excessEnergy float64, upcomingTimers *PoolingActions, shiftCost float64) bool {
	if !cat.Rotation.BearWeave || cat.ClearcastingAura.IsActive() || cat.BerserkAura.IsActive() || (cat.StampedeCatAura.IsActive() && (cat.Rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget)) {
		return false
	}

	// If we can Shred now and then weave on the next GCD, prefer that.
	if excessEnergy > cat.Shred.DefaultCast.Cost {
		return false
	}

	// Calculate effective Energy cap for out-of-form pooling
	targetWeaveDuration := core.GCDDefault*2 + cat.ReactionTime*2

	if (cat.Talents.Furor == 3) && (!cat.Rotation.MeleeWeave || (cat.CatCharge.TimeToReady(sim) > targetWeaveDuration)) {
		targetWeaveDuration += core.GCDDefault
	}

	weaveEnergy := furorCap - targetWeaveDuration.Seconds()*regenRate

	if currentEnergy > weaveEnergy {
		return false
	}

	// Prioritize all timers over weaving
	earliestWeaveEnd := sim.CurrentTime + core.GCDDefault*3 + cat.ReactionTime*2
	isPooling, nextRefresh := upcomingTimers.nextRefreshTime()

	if isPooling && (nextRefresh < earliestWeaveEnd) {
		return false
	}

	// Mana check
	if cat.CurrentMana() < shiftCost*2 {
		cat.Metrics.MarkOOM(sim)
		return false
	}

	// Also add a condition to make sure we can spend down our Energy post-weave before
	// the encounter ends or TF is ready.
	energyToDump := currentEnergy + (earliestWeaveEnd-sim.CurrentTime).Seconds()*regenRate
	timeToDump := earliestWeaveEnd + core.DurationFromSeconds(math.Floor(energyToDump/cat.Shred.DefaultCast.Cost))
	return (timeToDump < sim.Duration) && !cat.tfExpectedBefore(sim, timeToDump)
}

func (cat *FeralDruid) terminateBearWeave(sim *core.Simulation, isClearcast bool, currentEnergy float64, furorCap float64, regenRate float64, upcomingTimers *PoolingActions) bool {
	// Shift back early if a bear auto resulted in an Omen proc
	if isClearcast {
		return true
	}

	// Also terminate early if Feral Charge is off cooldown to avoid accumulating delays for Ravage opportunities
	if cat.Rotation.MeleeWeave && cat.CatCharge.IsReady(sim) && (sim.CurrentTime-cat.lastShift > core.GCDDefault) {
		return true
	}

	// Check Energy pooling leeway
	nextGCDLength := core.TernaryDuration(cat.Lacerate.CanCast(sim, cat.CurrentTarget) || cat.MangleBear.CanCast(sim, cat.CurrentTarget), core.GCDDefault, core.GCDMin)
	smallestWeaveExtension := nextGCDLength + cat.ReactionTime
	finalEnergy := currentEnergy + smallestWeaveExtension.Seconds()*regenRate

	if finalEnergy > furorCap {
		return true
	}

	// Check timer leeway
	earliestWeaveEnd := sim.CurrentTime + smallestWeaveExtension + core.GCDDefault
	isPooling, nextRefresh := upcomingTimers.nextRefreshTime()

	if isPooling && (nextRefresh < earliestWeaveEnd) {
		return true
	}

	// Also add a condition to prevent extending a weave if we don't have enough time
	// to spend the pooled Energy thus far.
	energyToDump := finalEnergy + 1.5*regenRate // need to include Cat Form GCD here
	timeToDump := earliestWeaveEnd + core.DurationFromSeconds(math.Floor(energyToDump/cat.Shred.DefaultCast.Cost))
	return (timeToDump >= sim.Duration) || cat.tfExpectedBefore(sim, timeToDump)
}

func (cat *FeralDruid) doRotation(sim *core.Simulation) (bool, time.Duration) {
	// Store state variables for re-use
	rotation := &cat.Rotation
	curEnergy := cat.CurrentEnergy()
	curCp := cat.ComboPoints()
	isClearcast := cat.ClearcastingAura.IsActive()
	simTimeRemain := sim.GetRemainingDuration()
	shiftCost := cat.CatForm.DefaultCast.Cost
	rakeDot := cat.Rake.CurDot()
	ripDot := cat.Rip.CurDot()
	lacerateDot := cat.Lacerate.CurDot()
	isBleedActive := cat.AssumeBleedActive || ripDot.IsActive() || rakeDot.IsActive() || lacerateDot.IsActive()
	regenRate := cat.EnergyRegenPerSecond()
	isExecutePhase := rotation.BiteDuringExecute && sim.IsExecutePhase25()
	tfActive := cat.TigersFuryAura.IsActive()
	berserkActive := cat.BerserkAura.IsActive()
	t11Active := cat.StrengthOfThePantherAura.IsActive()

	// Prioritize using Rip with omen procs if bleed isnt active
	ripCcCheck := core.Ternary(isBleedActive, !isClearcast, true)

	// Allow Clearcast Rakes if we will lose Rake uptime by Shredding first
	rakeCcCheck := !isClearcast || !rakeDot.IsActive() || (rakeDot.RemainingDuration(sim) < time.Second)

	// Use DPE calculation for deciding the end-of-fight breakpoint for Rip vs. Bite usage
	baseEndThresh := cat.calcRipEndThresh(sim)
	finalTickLeeway := core.TernaryDuration(ripDot.IsActive(), ripDot.TimeUntilNextTick(sim), 0)
	endThreshForClip := baseEndThresh + finalTickLeeway
	ripRefreshTime := cat.calcBleedRefreshTime(sim, cat.Rip, ripDot, isExecutePhase, true)
	ripNow := (curCp >= rotation.MinCombosForRip) && (!ripDot.IsActive() || ((sim.CurrentTime > ripRefreshTime) && !isExecutePhase)) && (simTimeRemain >= endThreshForClip) && ripCcCheck
	biteAtEnd := (curCp >= rotation.MinCombosForBite) && ((simTimeRemain < endThreshForClip) || (ripDot.IsActive() && (simTimeRemain-ripDot.RemainingDuration(sim) < baseEndThresh)))

	// Delay Rip refreshes if Tiger's Fury will be usable soon enough for the snapshot to outweigh the lost Rip ticks from waiting
	if ripNow && !tfActive && !berserkActive {
		buffedTickCount := min(cat.maxRipTicks, int32((simTimeRemain-finalTickLeeway)/ripDot.BaseTickLength))
		delayBreakpoint := finalTickLeeway + core.DurationFromSeconds(0.15*float64(buffedTickCount)*ripDot.BaseTickLength.Seconds())

		if cat.tfExpectedBefore(sim, sim.CurrentTime+delayBreakpoint) {
			delaySeconds := delayBreakpoint.Seconds()
			energyToDump := curEnergy + delaySeconds*regenRate - cat.calcTfEnergyThresh(cat.ReactionTime)
			secondsToDump := math.Ceil(energyToDump / cat.Shred.DefaultCast.Cost)

			if (secondsToDump < delaySeconds) && (!cat.tempSnapshotAura.IsActive() || (cat.tempSnapshotAura.RemainingDuration(sim) > delayBreakpoint)) {
				ripNow = false
			}
		}
	}

	// Clip Mangle if it won't change the total number of Mangles we have to
	// cast before the fight ends.
	t11BuildNow := (cat.StrengthOfThePantherAura != nil) && (cat.StrengthOfThePantherAura.GetStacks() < 3) && !rotation.BearWeave
	t11RefreshNow := t11Active && (cat.StrengthOfThePantherAura.RemainingDuration(sim) < cat.ReactionTime+max(time.Second, core.DurationFromSeconds((cat.MangleCat.DefaultCast.Cost-(curEnergy-cat.Shred.DefaultCast.Cost-core.TernaryFloat64(cat.PrimalMadnessAura.IsActive() && (cat.PrimalMadnessAura.ExpiresAt() < cat.StrengthOfThePantherAura.ExpiresAt()), cat.primalMadnessBonus, 0)))/regenRate))) && (simTimeRemain > time.Second)
	t11RefreshNext := t11Active && (cat.StrengthOfThePantherAura.RemainingDuration(sim) < time.Second*2+cat.ReactionTime) && (simTimeRemain > time.Second*2)
	mangleRefreshNow := !cat.bleedAura.IsActive() && (simTimeRemain > time.Second)
	mangleRefreshPending := (!t11RefreshNow && !mangleRefreshNow) && ((cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < (simTimeRemain-time.Second)) || (t11Active && (cat.StrengthOfThePantherAura.GetStacks() == 3) && (cat.StrengthOfThePantherAura.RemainingDuration(sim) < simTimeRemain-time.Second)))
	clipMangle := false

	if mangleRefreshPending && !t11Active && !cat.Rotation.BearWeave {
		numManglesRemaining := 1 + int32((sim.Duration-time.Second-cat.bleedAura.ExpiresAt())/time.Minute)
		earliestMangle := sim.Duration - time.Duration(numManglesRemaining)*time.Minute
		clipMangle = (sim.CurrentTime >= earliestMangle) && !isClearcast
	}

	mangleNow := cat.MangleCat != nil && (mangleRefreshNow || clipMangle)

	biteBeforeRip := (curCp >= rotation.MinCombosForBite) && ripDot.IsActive() && cat.SavageRoarAura.IsActive() && (rotation.UseBite || isExecutePhase) && cat.canBite(sim, isExecutePhase)
	biteNow := (biteBeforeRip || biteAtEnd) && !isClearcast

	// Ignore minimum CP enforcement during Execute phase if Rip is about to fall off
	emergencyBiteNow := isExecutePhase && ripDot.IsActive() && (ripDot.RemainingDuration(sim) < ripDot.BaseTickLength) && (curCp >= 1)
	biteNow = (biteNow || emergencyBiteNow) && !t11RefreshNext

	// Rake calcs
	rakeRefreshTime := cat.calcBleedRefreshTime(sim, cat.Rake, rakeDot, isExecutePhase, false)
	rakeNow := rotation.UseRake && (!rakeDot.IsActive() || (sim.CurrentTime > rakeRefreshTime)) && (simTimeRemain > rakeDot.BaseTickLength) && rakeCcCheck

	// Additionally, don't Rake if the current Shred DPE is higher due to
	// trinket procs etc.
	if rotation.RakeDpeCheck && rakeNow {
		rakeDpe, shredDpe := cat.calcBuilderDpe(sim)
		rakeNow = (rakeDpe > shredDpe)
	}

	// Additionally, don't Rake if there is insufficient time to max out
	// our available glyph of shred extensions before rip falls off
	if rakeNow && ripDot.IsActive() {
		remainingExt := cat.maxRipTicks - ripDot.BaseTickCount
		remainingRipDur := ripDot.RemainingDuration(sim) + time.Duration(remainingExt)*ripDot.BaseTickLength
		energyForShreds := curEnergy - cat.CurrentRakeCost() - cat.Rip.DefaultCast.Cost + remainingRipDur.Seconds()*regenRate + core.Ternary(cat.tfExpectedBefore(sim, sim.CurrentTime+remainingRipDur), 60.0, 0.0)
		maxShredsPossible := min(energyForShreds/cat.Shred.DefaultCast.Cost, (ripDot.ExpiresAt() - (sim.CurrentTime + time.Second)).Seconds())
		rakeNow = remainingExt == 0 || (maxShredsPossible > float64(remainingExt))
	}

	// Apply same TF Rip delay logic to Rake as well
	if rakeNow && !tfActive && !berserkActive {
		finalRakeTickLeeway := core.TernaryDuration(rakeDot.IsActive(), rakeDot.TimeUntilNextTick(sim), 0)
		buffedTickCount := min(rakeDot.BaseTickCount, int32((simTimeRemain-finalRakeTickLeeway)/rakeDot.BaseTickLength))
		delayBreakpoint := finalRakeTickLeeway + core.DurationFromSeconds(0.15*float64(buffedTickCount)*rakeDot.BaseTickLength.Seconds())

		if cat.tfExpectedBefore(sim, sim.CurrentTime+delayBreakpoint) {
			delaySeconds := delayBreakpoint.Seconds()
			energyToDump := curEnergy + delaySeconds*regenRate - cat.calcTfEnergyThresh(cat.ReactionTime)
			secondsToDump := math.Ceil(energyToDump / cat.Shred.DefaultCast.Cost)

			if secondsToDump < delaySeconds {
				rakeNow = false
			}
		}
	}

	// Roar calcs
	roarNow := (curCp >= 1) && (!cat.SavageRoarAura.IsActive() || cat.clipRoar(sim, isExecutePhase)) && (ripDot.IsActive() || (curCp < 3) || (simTimeRemain < baseEndThresh))

	// Ravage calc
	ravageNow := cat.Ravage.CanCast(sim, cat.CurrentTarget) && !isClearcast && ((curEnergy+2*regenRate < cat.MaximumEnergy()) || (cat.StampedeCatAura.RemainingDuration(sim) < time.Second*2))

	// Pooling calcs
	ripRefreshPending := ripDot.IsActive() && (ripDot.RemainingDuration(sim) < simTimeRemain-baseEndThresh) && (curCp >= core.TernaryInt32(isExecutePhase, 1, rotation.MinCombosForRip))
	rakeRefreshPending := rakeDot.IsActive() && (rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.BaseTickLength)
	roarRefreshPending := cat.SavageRoarAura.IsActive() && (cat.SavageRoarAura.RemainingDuration(sim) < simTimeRemain-cat.ReactionTime) && (curCp >= 1)
	cat.pendingPool.reset()
	cat.pendingPoolWeaves.reset()

	if ripRefreshPending && (sim.CurrentTime < ripRefreshTime) {
		baseCost := core.Ternary(isExecutePhase, cat.FerociousBite.DefaultCast.Cost, cat.Rip.DefaultCast.Cost)
		refreshCost := core.Ternary(cat.berserkExpectedAt(sim, ripRefreshTime), baseCost*0.5, baseCost)
		cat.pendingPool.addAction(ripRefreshTime, refreshCost)
		cat.pendingPoolWeaves.addAction(ripRefreshTime, refreshCost)
	}
	if rakeRefreshPending && (sim.CurrentTime < rakeRefreshTime) {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeRefreshTime), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		cat.pendingPool.addAction(rakeRefreshTime, rakeCost)
		cat.pendingPoolWeaves.addAction(rakeRefreshTime, rakeCost)
	}
	if mangleRefreshPending {
		mangleRefreshTime := cat.bleedAura.ExpiresAt()
		if t11Active {
			mangleRefreshTime = cat.StrengthOfThePantherAura.ExpiresAt() - cat.ReactionTime*2
		}
		mangleCost := core.Ternary(cat.berserkExpectedAt(sim, mangleRefreshTime), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		cat.pendingPool.addAction(mangleRefreshTime, mangleCost)
	}
	if roarRefreshPending {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		cat.pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)
	}

	cat.pendingPool.sort()
	cat.pendingPoolWeaves.sort()
	floatingEnergy := cat.pendingPool.calcFloatingEnergy(cat, sim)
	excessE := curEnergy - floatingEnergy
	latencySecs := cat.ReactionTime.Seconds()

	// Check melee-weaving conditions
	meleeWeaveNow := cat.canMeleeWeave(sim, regenRate, curEnergy, isClearcast, cat.pendingPool)

	// Check bear-weaving conditions
	furorCap := min(float64(100*cat.Talents.Furor)/3.0, 100.0-1.5*regenRate)
	bearWeaveNow := cat.canBearWeave(sim, furorCap, regenRate, curEnergy, excessE, cat.pendingPoolWeaves, shiftCost)
	// Main  decision tree starts here
	timeToNextAction := time.Duration(0)

	if !cat.CatFormAura.IsActive() {
		// First determine what we want to do with the next GCD.
		if cat.terminateBearWeave(sim, isClearcast, curEnergy, furorCap, regenRate, cat.pendingPoolWeaves) {
			cat.readyToShift = true
		} else if cat.MangleBear.CanCast(sim, cat.CurrentTarget) {
			cat.MangleBear.Cast(sim, cat.CurrentTarget)
		} else if cat.Lacerate.CanCast(sim, cat.CurrentTarget) {
			cat.Lacerate.Cast(sim, cat.CurrentTarget)
		} else if cat.FaerieFire.CanCast(sim, cat.CurrentTarget) {
			cat.FaerieFire.Cast(sim, cat.CurrentTarget)
		} else {
			cat.readyToShift = true
		}

		// Last second Maul check if we are about to shift back..
		if cat.readyToShift && cat.Maul.CanCast(sim, cat.CurrentTarget) && !isClearcast {
			cat.Maul.Cast(sim, cat.CurrentTarget)
		}

		if !cat.readyToShift {
			timeToNextAction = cat.ReactionTime
		}
	} else if t11RefreshNow {
		if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentMangleCatCost() - curEnergy) / regenRate)
	} else if ripNow {
		if cat.Rip.CanCast(sim, cat.CurrentTarget) {
			cat.Rip.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentRipCost() - curEnergy) / regenRate)
	} else if roarNow {
		if cat.SavageRoar.CanCast(sim, cat.CurrentTarget) {
			cat.SavageRoar.Cast(sim, nil)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentSavageRoarCost() - curEnergy) / regenRate)
	} else if biteNow && ((curEnergy >= cat.CurrentFerociousBiteCost()) || !bearWeaveNow) {
		if cat.FerociousBite.CanCast(sim, cat.CurrentTarget) {
			cat.FerociousBite.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentFerociousBiteCost() - curEnergy) / regenRate)
	} else if mangleNow {
		if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentMangleCatCost() - curEnergy) / regenRate)
	} else if rakeNow {
		if cat.Rake.CanCast(sim, cat.CurrentTarget) {
			cat.Rake.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentRakeCost() - curEnergy) / regenRate)
	} else if t11BuildNow {
		if cat.MangleCat.CanCast(sim, cat.CurrentTarget) {
			cat.MangleCat.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		timeToNextAction = core.DurationFromSeconds((cat.CurrentMangleCatCost() - curEnergy) / regenRate)
	} else if bearWeaveNow {
		cat.readyToShift = true
	} else if meleeWeaveNow {
		// Perform a final check to make sure we will have enough Energy to actually cast Feral Charge once we are in range, and delay the run-out slightly if not.
		minRunOutSeconds := (cat.CatCharge.MinRange - cat.DistanceFromTarget) / cat.GetMovementSpeed() // intentionally under-estimate for a safe buffer
		projectedEnergy := curEnergy + minRunOutSeconds*regenRate

		if tfActive && cat.PrimalMadnessAura.IsActive() {
			latestCharge := sim.CurrentTime + core.DurationFromSeconds((cat.CatCharge.MinRange+1-cat.DistanceFromTarget)/cat.GetMovementSpeed()) + cat.ReactionTime*2

			if cat.TigersFuryAura.ExpiresAt() < latestCharge {
				projectedEnergy -= cat.primalMadnessBonus
			}
		}

		if projectedEnergy >= cat.CatCharge.DefaultCast.Cost {
			cat.MoveTo(cat.CatCharge.MinRange+1, sim)
		} else {
			timeToNextAction = core.DurationFromSeconds((cat.CatCharge.DefaultCast.Cost - projectedEnergy) / regenRate)
		}
	} else if ravageNow {
		cat.Ravage.Cast(sim, cat.CurrentTarget)
		return false, 0
	} else if !t11RefreshNext && (isClearcast || !ripRefreshPending || !cat.tempSnapshotAura.IsActive() || (ripRefreshTime+cat.ReactionTime-sim.CurrentTime > core.GCDMin)) {
		fillerSpell := core.Ternary(rotation.MangleFiller, cat.MangleCat, cat.Shred)
		fillerCost := core.TernaryFloat64(rotation.MangleFiller, cat.CurrentMangleCatCost(), cat.CurrentShredCost())

		if excessE >= fillerCost || isClearcast {
			fillerSpell.Cast(sim, cat.CurrentTarget)
			return false, 0
		}
		// Also Shred if we're about to cap on Energy. Catches some edge
		// cases where floating_energy > 100 due to too many synced timers.
		if curEnergy > cat.MaximumEnergy()-regenRate*latencySecs {
			fillerSpell.Cast(sim, cat.CurrentTarget)
			return false, 0
		}

		timeToNextAction = core.DurationFromSeconds((fillerCost - excessE) / regenRate)

		if berserkActive {
			if curEnergy >= fillerCost {
				fillerSpell.Cast(sim, cat.CurrentTarget)
				return false, 0
			}
			timeToNextAction = core.DurationFromSeconds((fillerCost - curEnergy) / regenRate)
		}
	}

	// Model in latency when waiting on Energy for our next action
	nextAction := sim.CurrentTime + timeToNextAction
	paValid, rt := cat.pendingPool.nextRefreshTime()
	if paValid {
		nextAction = min(nextAction, rt)
	}

	return true, nextAction
}

type FeralDruidRotation struct {
	RotationType proto.FeralDruid_Rotation_AplType

	BearWeave           bool
	MaintainFaerieFire  bool
	MinCombosForRip     int32
	UseRake             bool
	UseBite             bool
	BiteTime            time.Duration
	BerserkBiteTime     time.Duration
	BiteDuringExecute   bool
	MinCombosForBite    int32
	MangleFiller        bool
	MinRoarOffset       time.Duration
	RipLeeway           time.Duration
	SnekWeave           bool
	RakeDpeCheck        bool
	UseBerserk          bool
	MeleeWeave          bool
	CancelPrimalMadness bool
}

func (cat *FeralDruid) setupRotation(rotation *proto.FeralDruid_Rotation) {
	cat.Rotation = FeralDruidRotation{
		RotationType:        rotation.RotationType,
		BearWeave:           rotation.BearWeave,
		MaintainFaerieFire:  rotation.MaintainFaerieFire,
		MinCombosForRip:     5,
		UseRake:             rotation.UseRake,
		UseBite:             rotation.UseBite,
		BiteTime:            time.Duration(float64(rotation.BiteTime) * float64(time.Second)),
		BerserkBiteTime:     time.Duration(float64(rotation.BerserkBiteTime) * float64(time.Second)),
		BiteDuringExecute:   core.Ternary(cat.Talents.BloodInTheWater > 0, rotation.BiteDuringExecute, false),
		MinCombosForBite:    5,
		MangleFiller:        cat.PseudoStats.InFrontOfTarget || cat.CannotShredTarget,
		MinRoarOffset:       time.Duration(float64(rotation.MinRoarOffset) * float64(time.Second)),
		RipLeeway:           time.Duration(float64(rotation.RipLeeway) * float64(time.Second)),
		SnekWeave:           rotation.SnekWeave,
		RakeDpeCheck:        true,
		UseBerserk:          cat.Talents.Berserk && ((rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget) || rotation.AllowAoeBerserk),
		MeleeWeave:          rotation.MeleeWeave && (cat.Talents.Stampede > 0) && (rotation.RotationType == proto.FeralDruid_Rotation_SingleTarget) && !cat.CannotShredTarget && !cat.PseudoStats.InFrontOfTarget,
		CancelPrimalMadness: rotation.CancelPrimalMadness,
	}

	// Use automatic values unless specified
	if rotation.ManualParams {
		return
	}

	cat.Rotation.UseRake = true
	cat.Rotation.UseBite = true
	cat.Rotation.BiteDuringExecute = (cat.Talents.BloodInTheWater == 2)
	cat.Rotation.CancelPrimalMadness = rotation.CancelPrimalMadness && (rotation.RotationType == proto.FeralDruid_Rotation_Aoe)

	cat.Rotation.RipLeeway = 1 * time.Second
	cat.Rotation.MinRoarOffset = 31 * time.Second
	cat.Rotation.BiteTime = 11 * time.Second
	cat.Rotation.BerserkBiteTime = 6 * time.Second
}
