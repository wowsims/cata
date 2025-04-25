package feral

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
)

func (cat *FeralDruid) NewAPLValue(rot *core.APLRotation, config *proto.APLValue) core.APLValue {
	switch config.Value.(type) {
	case *proto.APLValue_CatExcessEnergy:
		return cat.newValueCatExcessEnergy(rot, config.GetCatExcessEnergy())
	case *proto.APLValue_CatNewSavageRoarDuration:
		return cat.newValueCatNewSavageRoarDuration(rot, config.GetCatNewSavageRoarDuration())
	default:
		return nil
	}
}

type APLValueCatExcessEnergy struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatExcessEnergy(_ *core.APLRotation, _ *proto.APLValueCatExcessEnergy) core.APLValue {
	return &APLValueCatExcessEnergy{
		cat: cat,
	}
}
func (value *APLValueCatExcessEnergy) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCatExcessEnergy) GetFloat(sim *core.Simulation) float64 {
	cat := value.cat
	pendingPool := PoolingActions{}
	pendingPool.create(4)

	simTimeRemain := sim.GetRemainingDuration()
	if ripDot := cat.Rip.CurDot(); ripDot.IsActive() && ripDot.RemainingDuration(sim) < simTimeRemain-time.Second*10 && cat.ComboPoints() == 5 {
		ripCost := core.Ternary(cat.berserkExpectedAt(sim, ripDot.ExpiresAt()), cat.Rip.DefaultCast.Cost*0.5, cat.Rip.DefaultCast.Cost)
		pendingPool.addAction(ripDot.ExpiresAt(), ripCost)
	}
	if rakeDot := cat.Rake.CurDot(); rakeDot.IsActive() && rakeDot.RemainingDuration(sim) < simTimeRemain-rakeDot.Duration {
		rakeCost := core.Ternary(cat.berserkExpectedAt(sim, rakeDot.ExpiresAt()), cat.Rake.DefaultCast.Cost*0.5, cat.Rake.DefaultCast.Cost)
		pendingPool.addAction(rakeDot.ExpiresAt(), rakeCost)
	}
	if cat.bleedAura.IsActive() && cat.bleedAura.RemainingDuration(sim) < simTimeRemain-time.Second {
		mangleCost := core.Ternary(cat.berserkExpectedAt(sim, cat.bleedAura.ExpiresAt()), cat.MangleCat.DefaultCast.Cost*0.5, cat.MangleCat.DefaultCast.Cost)
		pendingPool.addAction(cat.bleedAura.ExpiresAt(), mangleCost)
	}
	if cat.SavageRoarAura.IsActive() {
		roarCost := core.Ternary(cat.berserkExpectedAt(sim, cat.SavageRoarAura.ExpiresAt()), cat.SavageRoar.DefaultCast.Cost*0.5, cat.SavageRoar.DefaultCast.Cost)
		pendingPool.addAction(cat.SavageRoarAura.ExpiresAt(), roarCost)
	}

	pendingPool.sort()

	floatingEnergy := pendingPool.calcFloatingEnergy(cat, sim)
	return cat.CurrentEnergy() - floatingEnergy
}
func (value *APLValueCatExcessEnergy) String() string {
	return "Cat Excess Energy()"
}

type APLValueCatNewSavageRoarDuration struct {
	core.DefaultAPLValueImpl
	cat *FeralDruid
}

func (cat *FeralDruid) newValueCatNewSavageRoarDuration(_ *core.APLRotation, _ *proto.APLValueCatNewSavageRoarDuration) core.APLValue {
	return &APLValueCatNewSavageRoarDuration{
		cat: cat,
	}
}
func (value *APLValueCatNewSavageRoarDuration) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeDuration
}
func (value *APLValueCatNewSavageRoarDuration) GetDuration(_ *core.Simulation) time.Duration {
	cat := value.cat
	return cat.SavageRoarDurationTable[cat.ComboPoints()]
}
func (value *APLValueCatNewSavageRoarDuration) String() string {
	return "New Savage Roar Duration()"
}

func (cat *FeralDruid) NewAPLAction(rot *core.APLRotation, config *proto.APLAction) core.APLActionImpl {
	switch config.Action.(type) {
	case *proto.APLAction_CatOptimalRotationAction:
		return cat.newActionCatOptimalRotationAction(rot, config.GetCatOptimalRotationAction())
	default:
		return nil
	}
}

type APLActionCatOptimalRotationAction struct {
	cat        *FeralDruid
	lastAction time.Duration
}

func (impl *APLActionCatOptimalRotationAction) GetInnerActions() []*core.APLAction { return nil }
func (impl *APLActionCatOptimalRotationAction) GetAPLValues() []core.APLValue      { return nil }
func (impl *APLActionCatOptimalRotationAction) Finalize(*core.APLRotation)         {}
func (impl *APLActionCatOptimalRotationAction) PostFinalize(*core.APLRotation)     {}
func (impl *APLActionCatOptimalRotationAction) GetNextAction(*core.Simulation) *core.APLAction {
	return nil
}

func (cat *FeralDruid) newActionCatOptimalRotationAction(_ *core.APLRotation, config *proto.APLActionCatOptimalRotationAction) core.APLActionImpl {
	rotationOptions := &proto.FeralDruid_Rotation{
		RotationType:        config.RotationType,
		MaintainFaerieFire:  config.MaintainFaerieFire,
		UseRake:             config.UseRake,
		UseBite:             config.UseBite,
		BiteTime:            config.BiteTime,
		BerserkBiteTime:     config.BerserkBiteTime,
		BiteDuringExecute:   config.BiteDuringExecute,
		MinRoarOffset:       config.MinRoarOffset,
		RipLeeway:           config.RipLeeway,
		ManualParams:        config.ManualParams,
		AllowAoeBerserk:     config.AllowAoeBerserk,
		MeleeWeave:          config.MeleeWeave,
		BearWeave:           config.BearWeave,
		SnekWeave:           config.SnekWeave,
		CancelPrimalMadness: config.CancelPrimalMadness,
	}

	cat.setupRotation(rotationOptions)

	// Pre-allocate PoolingActions
	cat.pendingPool = &PoolingActions{}
	cat.pendingPool.create(4)
	cat.pendingPoolWeaves = &PoolingActions{}
	cat.pendingPoolWeaves.create(2)

	return &APLActionCatOptimalRotationAction{
		cat: cat,
	}
}

func (action *APLActionCatOptimalRotationAction) IsReady(sim *core.Simulation) bool {
	return sim.CurrentTime > action.lastAction
}

func (action *APLActionCatOptimalRotationAction) Execute(sim *core.Simulation) {
	cat := action.cat

	// If a melee swing resulted in an Omen proc, then schedule the
	// next player decision based on latency.
	ccRefreshTime := cat.ClearcastingAura.ExpiresAt() - cat.ClearcastingAura.Duration

	if ccRefreshTime >= sim.CurrentTime-cat.ReactionTime {
		// Kick gcd loop, also need to account for any gcd 'left'
		// otherwise it breaks gcd logic
		kickTime := max(cat.NextGCDAt(), ccRefreshTime+cat.ReactionTime)
		cat.NextRotationAction(sim, kickTime)
	}

	action.lastAction = sim.CurrentTime

	// Keep up Sunder debuff if not provided externally. Do this here since FF can be
	// cast while moving.
	if cat.Rotation.MaintainFaerieFire {
		for _, aoeTarget := range sim.Encounter.TargetUnits {
			if cat.ShouldFaerieFire(sim, aoeTarget) {
				cat.FaerieFire.CastOrQueue(sim, aoeTarget)
			}
		}
	}

	// Off-GCD bear-weave checks.
	if cat.BearFormAura.IsActive() && !cat.ClearcastingAura.IsActive() {
		if cat.Enrage.IsReady(sim) && !cat.readyToShift {
			cat.Enrage.Cast(sim, nil)
		}

		if cat.Maul.CanCast(sim, cat.CurrentTarget) && ((cat.CurrentRage() >= cat.Maul.DefaultCast.Cost+cat.MangleBear.DefaultCast.Cost) || (cat.AutoAttacks.NextAttackAt() < cat.NextGCDAt())) {
			cat.Maul.Cast(sim, cat.CurrentTarget)
		}
	}

	// Handle movement before any rotation logic
	if cat.Moving || (cat.Hardcast.Expires > sim.CurrentTime) {
		return
	}
	if cat.DistanceFromTarget > core.MaxMeleeRange {
		// Try leaping first before defaulting to manual movement
		if cat.CatCharge.CanCast(sim, cat.CurrentTarget) {
			cat.CatCharge.Cast(sim, cat.CurrentTarget)
		} else {
			if sim.Log != nil {
				cat.Log(sim, "Out of melee range (%.6fy) and cannot Charge (remaining CD: %s), initiating manual run-in...", cat.DistanceFromTarget, cat.CatCharge.TimeToReady(sim))
			}

			cat.MoveTo(core.MaxMeleeRange-1, sim) // movement aura is discretized in 1 yard intervals, so need to overshoot to guarantee melee range
			return
		}
	}

	if !cat.GCD.IsReady(sim) {
		cat.WaitUntil(sim, cat.NextGCDAt())
		return
	}

	cat.TryTigersFury(sim)
	cat.TryBerserk(sim)

	if sim.CurrentTime >= cat.nextActionAt {
		cat.OnGCDReady(sim)
	} else {
		cat.WaitUntil(sim, cat.nextActionAt)
	}
}

func (action *APLActionCatOptimalRotationAction) Reset(*core.Simulation) {
	action.cat.usingHardcodedAPL = true
	action.cat.cachedRipEndThresh = time.Second * 10 // placeholder until first calc
	action.lastAction = core.DurationFromSeconds(-100)
}

func (action *APLActionCatOptimalRotationAction) String() string {
	return "Execute Optimal Cat Action()"
}
