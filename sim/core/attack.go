package core

import (
	"slices"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// ReplaceMHSwing is called right before a main hand auto attack fires.
// It must never return nil, but either a replacement spell or the passed in regular mhSwingSpell.
// This allows for abilities that convert a white attack into a yellow attack.
type ReplaceMHSwing func(sim *Simulation, mhSwingSpell *Spell) *Spell

// Represents a generic weapon. Pets / unarmed / various other cases don't use
// actual weapon items so this is an abstraction of a Weapon.
type Weapon struct {
	BaseDamageMin        float64
	BaseDamageMax        float64
	AttackPowerPerDPS    float64
	SwingSpeed           float64
	NormalizedSwingSpeed float64
	CritMultiplier       float64
	SpellSchool          SpellSchool
	MinRange             float64
	MaxRange             float64
}

func (weapon *Weapon) DPS() float64 {
	if weapon.SwingSpeed == 0 {
		return 0
	}
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2.0 / weapon.SwingSpeed
}

func newWeaponFromUnarmed(critMultiplier float64) Weapon {
	// These numbers are probably wrong but nobody cares.
	return Weapon{
		BaseDamageMin:        0,
		BaseDamageMax:        0,
		SwingSpeed:           1,
		NormalizedSwingSpeed: 1,
		CritMultiplier:       critMultiplier,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
		MaxRange:             MaxMeleeRange,
	}
}

func getWeaponMaxRange(item *Item) float64 {
	switch item.RangedWeaponType {
	case proto.RangedWeaponType_RangedWeaponTypeUnknown:
		return MaxMeleeRange
	case proto.RangedWeaponType_RangedWeaponTypeWand:
	case proto.RangedWeaponType_RangedWeaponTypeThrown:
		return 30
	default:
		return 40
	}

	return 40
}

func newWeaponFromItem(item *Item, critMultiplier float64, bonusDps float64) Weapon {
	normalizedWeaponSpeed := 2.4
	if item.WeaponType == proto.WeaponType_WeaponTypeDagger {
		normalizedWeaponSpeed = 1.7
	} else if item.HandType == proto.HandType_HandTypeTwoHand {
		normalizedWeaponSpeed = 3.3
	} else if item.RangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
		normalizedWeaponSpeed = 2.8
	}

	return Weapon{
		BaseDamageMin:        item.WeaponDamageMin + bonusDps*item.SwingSpeed,
		BaseDamageMax:        item.WeaponDamageMax + bonusDps*item.SwingSpeed,
		SwingSpeed:           item.SwingSpeed,
		NormalizedSwingSpeed: normalizedWeaponSpeed,
		CritMultiplier:       critMultiplier,
		AttackPowerPerDPS:    DefaultAttackPowerPerDPS,
		MinRange:             0, // no more deadzone in MoP
		MaxRange:             getWeaponMaxRange(item),
	}
}

// Returns weapon stats using the main hand equipped weapon.
func (character *Character) WeaponFromMainHand(critMultiplier float64) Weapon {
	if weapon := character.GetMHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusMHDps)
	} else {
		return newWeaponFromUnarmed(critMultiplier)
	}
}

// Returns weapon stats using the off-hand equipped weapon.
func (character *Character) WeaponFromOffHand(critMultiplier float64) Weapon {
	if weapon := character.GetOHWeapon(); weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusOHDps)
	} else {
		return Weapon{}
	}
}

// Returns weapon stats using the ranged equipped weapon.
func (character *Character) WeaponFromRanged(critMultiplier float64) Weapon {
	weapon := character.Ranged()
	if weapon != nil {
		return newWeaponFromItem(weapon, critMultiplier, character.PseudoStats.BonusRangedDps)
	} else {
		return Weapon{}
	}
}

func (weapon *Weapon) GetSpellSchool() SpellSchool {
	if weapon.SpellSchool == SpellSchoolNone {
		return SpellSchoolPhysical
	} else {
		return weapon.SpellSchool
	}
}

func (weapon *Weapon) EnemyWeaponDamage(sim *Simulation, attackPower float64, damageSpread float64) float64 {
	// Maximum damage range is 133% of minimum damage; AP contribution is % of minimum damage roll.
	// Patchwerk follows special damage range rules.
	// TODO: Scrape more logs to determine these values more accurately. AP defined in constants.go

	rand := 1 + damageSpread*sim.RandomFloat("Enemy Weapon Damage")

	return weapon.BaseDamageMin * (rand + attackPower*EnemyAutoAttackAPCoefficient)
}

func (weapon *Weapon) BaseDamage(sim *Simulation) float64 {
	return weapon.BaseDamageMin + (weapon.BaseDamageMax-weapon.BaseDamageMin)*sim.RandomFloat("Weapon Base Damage")
}

func (weapon *Weapon) AverageDamage() float64 {
	return (weapon.BaseDamageMin + weapon.BaseDamageMax) / 2
}

func (weapon *Weapon) CalculateWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (weapon *Weapon) CalculateAverageWeaponDamage(attackPower float64) float64 {
	return weapon.AverageDamage() + (weapon.SwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (weapon *Weapon) CalculateNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return weapon.BaseDamage(sim) + (weapon.NormalizedSwingSpeed*attackPower)/weapon.AttackPowerPerDPS
}

func (unit *Unit) MHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.mh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) MHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.mh.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) OHWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateWeaponDamage(sim, attackPower)
}
func (unit *Unit) OHNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return 0.5 * unit.AutoAttacks.oh.CalculateNormalizedWeaponDamage(sim, attackPower)
}
func (unit *Unit) RangedNormalizedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.ranged.CalculateNormalizedWeaponDamage(sim, attackPower)
}

func (unit *Unit) RangedWeaponDamage(sim *Simulation, attackPower float64) float64 {
	return unit.AutoAttacks.ranged.CalculateWeaponDamage(sim, attackPower)
}

type MeleeDamageCalculator func(attackPower float64, bonusWeaponDamage float64) float64

// Returns whether this hit effect is associated with the main-hand weapon.
func (spell *Spell) IsMH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeMH)
}

// Returns whether this hit effect is associated with the off-hand weapon.
func (spell *Spell) IsOH() bool {
	return spell.ProcMask.Matches(ProcMaskMeleeOH)
}

// Returns whether this hit effect is associated with either melee weapon.
func (spell *Spell) IsMelee() bool {
	return spell.ProcMask.Matches(ProcMaskMelee)
}

// Returns whether this hit effect is associated with a ranged weapon.
func (spell *Spell) IsRanged() bool {
	return spell.ProcMask.Matches(ProcMaskRanged)
}

func (aa *AutoAttacks) MH() *Weapon {
	return aa.mh.getWeapon()
}

func (aa *AutoAttacks) SetMH(weapon Weapon) {
	aa.mh.setWeapon(weapon)
}

func (aa *AutoAttacks) OH() *Weapon {
	return aa.oh.getWeapon()
}

func (aa *AutoAttacks) SetOH(weapon Weapon) {
	aa.oh.setWeapon(weapon)
}

func (aa *AutoAttacks) Ranged() *Weapon {
	return aa.ranged.getWeapon()
}

func (aa *AutoAttacks) SetRanged(weapon Weapon) {
	aa.ranged.setWeapon(weapon)
}

func (aa *AutoAttacks) MHAuto() *Spell {
	return aa.mh.spell
}

func (aa *AutoAttacks) OHAuto() *Spell {
	return aa.oh.spell
}

func (aa *AutoAttacks) RangedAuto() *Spell {
	return aa.ranged.spell
}

func (aa *AutoAttacks) OffhandSwingAt() time.Duration {
	return aa.oh.swingAt
}

func (aa *AutoAttacks) SetOffhandSwingAt(offhandSwingAt time.Duration) {
	aa.oh.swingAt = offhandSwingAt
}

func (aa *AutoAttacks) SetReplaceMHSwing(replaceSwing ReplaceMHSwing) {
	aa.mh.replaceSwing = replaceSwing
}

func (aa *AutoAttacks) MHConfig() *SpellConfig {
	return &aa.mh.config
}

func (aa *AutoAttacks) OHConfig() *SpellConfig {
	return &aa.oh.config
}

func (aa *AutoAttacks) RangedConfig() *SpellConfig {
	return &aa.ranged.config
}

type WeaponAttack struct {
	Weapon

	agent Agent
	unit  *Unit

	config SpellConfig
	spell  *Spell

	replaceSwing ReplaceMHSwing

	swingAt time.Duration

	curSwingSpeed    float64
	curSwingDuration time.Duration
	enabled          bool
}

func (wa *WeaponAttack) getWeapon() *Weapon {
	return &wa.Weapon
}

func (wa *WeaponAttack) setWeapon(weapon Weapon) {
	wa.Weapon = weapon
	wa.spell.CritMultiplier = weapon.CritMultiplier
	wa.updateSwingDuration(wa.curSwingSpeed)
}

// inlineable stub for swing
func (wa *WeaponAttack) trySwing(sim *Simulation) time.Duration {
	if sim.CurrentTime < wa.swingAt {
		return wa.swingAt
	}
	return wa.swing(sim)
}

func (wa *WeaponAttack) swing(sim *Simulation) time.Duration {
	attackSpell := wa.spell

	if wa.replaceSwing != nil {
		// Need to check APL here to allow last-moment HS queue casts.
		wa.unit.ReactToEvent(sim)

		// Need to check this again in case the DoNextAction call swapped items.
		if wa.replaceSwing != nil {
			// Allow MH swing to be overridden for abilities like Heroic Strike.
			attackSpell = wa.replaceSwing(sim, attackSpell)
		}
	}

	// Update swing timer BEFORE the cast, so that APL checks for TimeToNextAuto behave correctly
	// if the attack causes APL evaluations (e.g. from rage gain).
	wa.swingAt = sim.CurrentTime + wa.curSwingDuration
	attackSpell.Cast(sim, wa.unit.CurrentTarget)

	if !sim.Options.Interactive && wa.unit.Rotation != nil {
		wa.unit.ReactToEvent(sim)
	}

	return wa.swingAt
}

func (wa *WeaponAttack) updateSwingDuration(curSwingSpeed float64) {
	wa.curSwingSpeed = curSwingSpeed
	wa.curSwingDuration = DurationFromSeconds(wa.SwingSpeed / wa.curSwingSpeed)
}

func (wa *WeaponAttack) addWeaponAttack(sim *Simulation, swingSpeed float64) {
	if !wa.enabled {
		return
	}

	wa.updateSwingDuration(swingSpeed)
	sim.addWeaponAttack(wa)
	sim.rescheduleWeaponAttack(wa.swingAt)
}

type AutoAttacks struct {
	AutoSwingMelee  bool
	AutoSwingRanged bool

	IsDualWielding bool

	character *Character

	mh     WeaponAttack
	oh     WeaponAttack
	ranged WeaponAttack
}

// Options for initializing auto attacks.
type AutoAttackOptions struct {
	MainHand        Weapon
	OffHand         Weapon
	Ranged          Weapon
	AutoSwingMelee  bool // If true, core engine will handle calling SwingMelee() for you.
	AutoSwingRanged bool // If true, core engine will handle calling SwingRanged() for you.
	ReplaceMHSwing  ReplaceMHSwing
}

func (unit *Unit) EnableAutoAttacks(agent Agent, options AutoAttackOptions) {
	if options.MainHand.AttackPowerPerDPS == 0 {
		options.MainHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}
	if options.OffHand.AttackPowerPerDPS == 0 {
		options.OffHand.AttackPowerPerDPS = DefaultAttackPowerPerDPS
	}

	unit.AutoAttacks = AutoAttacks{
		AutoSwingMelee:  options.AutoSwingMelee,
		AutoSwingRanged: options.AutoSwingRanged,

		IsDualWielding: options.OffHand.SwingSpeed != 0,

		character: agent.GetCharacter(),

		mh: WeaponAttack{
			agent:        agent,
			unit:         unit,
			Weapon:       options.MainHand,
			replaceSwing: options.ReplaceMHSwing,
		},
		oh: WeaponAttack{
			agent:  agent,
			unit:   unit,
			Weapon: options.OffHand,
		},
		ranged: WeaponAttack{
			agent:  agent,
			unit:   unit,
			Weapon: options.Ranged,
		},
	}

	unit.AutoAttacks.mh.config = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 1},
		SpellSchool: options.MainHand.GetSpellSchool(),
		ProcMask:    ProcMaskMeleeMHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           options.MainHand.CritMultiplier,
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite)
		},

		ExpectedInitialDamage: func(sim *Simulation, target *Unit, spell *Spell, _ bool) *SpellResult {
			baseDamage := spell.Unit.AutoAttacks.MH().CalculateAverageWeaponDamage(spell.MeleeAttackPower())
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMeleeWhite)
		},
	}

	unit.AutoAttacks.oh.config = SpellConfig{
		ActionID:    ActionID{OtherID: proto.OtherAction_OtherActionAttack, Tag: 2},
		SpellSchool: options.OffHand.GetSpellSchool(),
		ProcMask:    ProcMaskMeleeOHAuto,
		Flags:       SpellFlagMeleeMetrics | SpellFlagNoOnCastComplete,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           options.OffHand.CritMultiplier,
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower())

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWhite)
		},
	}

	unit.AutoAttacks.ranged.config = SpellConfig{
		ActionID:     ActionID{OtherID: proto.OtherAction_OtherActionShoot},
		SpellSchool:  options.Ranged.GetSpellSchool(),
		ProcMask:     ProcMaskRangedAuto,
		Flags:        SpellFlagMeleeMetrics,
		MissileSpeed: 40,

		DamageMultiplier:         1,
		DamageMultiplierAdditive: 1,
		CritMultiplier:           options.Ranged.CritMultiplier,
		ThreatMultiplier:         1,

		BonusCoefficient: 1,

		ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
			baseDamage := spell.Unit.RangedWeaponDamage(sim, spell.RangedAttackPower())
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	}

	if unit.Type == EnemyUnit {
		unit.AutoAttacks.mh.config.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
		unit.AutoAttacks.oh.config.ApplyEffects = func(sim *Simulation, target *Unit, spell *Spell) {
			ap := max(0, spell.Unit.stats[stats.AttackPower])
			baseDamage := spell.Unit.AutoAttacks.mh.EnemyWeaponDamage(sim, ap, spell.Unit.PseudoStats.DamageSpread) * 0.5

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)
		}
	}
}

func (aa *AutoAttacks) finalize() {
	if aa.AutoSwingMelee {
		aa.mh.spell = aa.mh.unit.GetOrRegisterSpell(aa.mh.config)

		// Will keep keep the OH spell registered for Item swapping
		aa.oh.spell = aa.oh.unit.GetOrRegisterSpell(aa.oh.config)
	}
	if aa.AutoSwingRanged {
		aa.ranged.spell = aa.ranged.unit.GetOrRegisterSpell(aa.ranged.config)
	}
}

func (aa *AutoAttacks) anyEnabled() bool {
	return aa.mh.enabled || aa.oh.enabled || aa.ranged.enabled
}

func (aa *AutoAttacks) reset(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	aa.mh.enabled = false
	aa.oh.enabled = false
	aa.ranged.enabled = false

	aa.mh.swingAt = NeverExpires
	aa.oh.swingAt = NeverExpires

	if aa.AutoSwingMelee {
		aa.mh.updateSwingDuration(aa.mh.unit.SwingSpeed())
		aa.mh.swingAt = 0

		if aa.IsDualWielding {
			aa.oh.updateSwingDuration(aa.mh.curSwingSpeed)
			aa.oh.swingAt = 0

			// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
			if aa.oh.unit.Type == EnemyUnit {
				aa.oh.swingAt = DurationFromSeconds(aa.mh.SwingSpeed / 2)
			} else {
				if sim.RandomFloat("SwingResetWeapon") < 0.5 {
					aa.mh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				} else {
					aa.oh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
				}
			}
		}

	}

	aa.ranged.swingAt = NeverExpires

	if aa.AutoSwingRanged {
		aa.ranged.updateSwingDuration(aa.ranged.unit.RangedSwingSpeed())
		aa.ranged.swingAt = 0
	}
}

func (aa *AutoAttacks) startPull(sim *Simulation) {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return
	}

	if aa.mh.unit.CurrentTarget == nil {
		return
	}

	if aa.anyEnabled() {
		return
	}

	if aa.AutoSwingMelee {
		if aa.mh.swingAt == NeverExpires {
			aa.mh.swingAt = 0
		}

		if aa.IsDualWielding {
			if aa.oh.swingAt == NeverExpires {
				aa.oh.swingAt = 0

				// Apply random delay of 0 - 50% swing time, to one of the weapons if dual wielding
				if aa.oh.unit.Type == EnemyUnit {
					aa.oh.swingAt = DurationFromSeconds(aa.mh.SwingSpeed / 2)
				} else {
					if sim.RandomFloat("SwingResetWeapon") < 0.5 {
						aa.mh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
					} else {
						aa.oh.swingAt = DurationFromSeconds(sim.RandomFloat("SwingResetDelay") * aa.mh.SwingSpeed / 2)
					}
				}
			}
			if aa.oh.IsInRange() {
				aa.oh.enabled = true
				aa.oh.addWeaponAttack(sim, aa.mh.curSwingSpeed)
			}

		}

		if aa.mh.IsInRange() {
			aa.mh.enabled = true
			aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		}
	}

	if aa.AutoSwingRanged {
		if aa.ranged.swingAt == NeverExpires {
			aa.ranged.swingAt = 0
		}
		if aa.ranged.IsInRange() {
			aa.ranged.enabled = true
			aa.ranged.addWeaponAttack(sim, aa.ranged.unit.RangedSwingSpeed())
		}

	}
}

func (wa *WeaponAttack) IsInRange() bool {
	return (wa.MinRange == 0. || wa.MinRange < wa.unit.DistanceFromTarget) && (wa.MaxRange == 0. || wa.MaxRange >= wa.unit.DistanceFromTarget)
}

// Stops the auto swing action for the rest of the iteration. Used for pets
// after being disabled.
func (aa *AutoAttacks) CancelAutoSwing(sim *Simulation) {
	aa.CancelMeleeSwing(sim)
	aa.CancelRangedSwing(sim)
}

// Re-enables the auto swing action for the iteration
func (aa *AutoAttacks) EnableAutoSwing(sim *Simulation) {
	aa.EnableMeleeSwing(sim)
	aa.EnableRangedSwing(sim)
}

func (aa *AutoAttacks) EnableMeleeSwing(sim *Simulation) {
	if !aa.AutoSwingMelee {
		return
	}

	if aa.mh.unit.CurrentTarget == nil {
		return
	}

	aa.mh.swingAt = max(aa.mh.swingAt, sim.CurrentTime, 0)
	if aa.mh.IsInRange() && !aa.mh.enabled {
		aa.mh.enabled = true
		aa.mh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
	}

	if aa.IsDualWielding && !aa.oh.enabled {
		aa.oh.swingAt = max(aa.oh.swingAt, sim.CurrentTime, 0)
		if aa.oh.IsInRange() {
			aa.oh.enabled = true
			aa.oh.addWeaponAttack(sim, aa.mh.unit.SwingSpeed())
		}
	}

	if !aa.IsDualWielding && aa.oh.enabled {
		sim.removeWeaponAttack(&aa.oh)
		aa.oh.enabled = false
	}
}

func (aa *AutoAttacks) EnableRangedSwing(sim *Simulation) {
	if !aa.AutoSwingRanged || aa.ranged.enabled {
		return
	}

	if aa.ranged.unit.CurrentTarget == nil {
		return
	}

	aa.ranged.swingAt = max(aa.ranged.swingAt, sim.CurrentTime, 0)
	if aa.ranged.IsInRange() {
		aa.ranged.enabled = true
		aa.ranged.addWeaponAttack(sim, aa.ranged.unit.RangedSwingSpeed())
	}
}

func (aa *AutoAttacks) CancelMeleeSwing(sim *Simulation) {
	if !aa.AutoSwingMelee {
		return
	}

	if aa.mh.enabled {
		sim.removeWeaponAttack(&aa.mh)
		aa.mh.enabled = false
	}

	if aa.IsDualWielding && aa.oh.enabled {
		aa.oh.enabled = false
		sim.removeWeaponAttack(&aa.oh)
	}
}

func (aa *AutoAttacks) CancelRangedSwing(sim *Simulation) {
	if !aa.AutoSwingRanged || !aa.ranged.enabled {
		return
	}

	aa.ranged.enabled = false
	sim.removeWeaponAttack(&aa.ranged)
}

// The amount of time between two MH swings.
func (aa *AutoAttacks) MainhandSwingSpeed() time.Duration {
	return aa.mh.curSwingDuration
}

// The amount of time between two OH swings.
func (aa *AutoAttacks) OffhandSwingSpeed() time.Duration {
	return aa.oh.curSwingDuration
}

// Optionally replaces the given swing spell with an Agent-specified MH Swing replacer.
// This is for effects like Heroic Strike or Raptor Strike.
func (aa *AutoAttacks) MaybeReplaceMHSwing(sim *Simulation, mhSwingSpell *Spell) *Spell {
	if aa.mh.replaceSwing == nil {
		return mhSwingSpell
	}

	// Allow MH swing to be overridden for abilities like Heroic Strike.
	return aa.mh.replaceSwing(sim, mhSwingSpell)
}

func (aa *AutoAttacks) UpdateSwingTimers(sim *Simulation) {
	if !aa.anyEnabled() {
		return
	}

	if aa.AutoSwingRanged && aa.ranged.enabled {
		aa.ranged.updateSwingDuration(aa.ranged.unit.RangedSwingSpeed())
		// ranged attack speed changes aren't applied mid-"swing"
	}

	if aa.AutoSwingMelee && aa.mh.enabled {
		oldSwingSpeed := aa.mh.curSwingSpeed
		aa.mh.updateSwingDuration(aa.mh.unit.SwingSpeed())
		f := oldSwingSpeed / aa.mh.curSwingSpeed

		if remainingSwingTime := aa.mh.swingAt - sim.CurrentTime; remainingSwingTime > 0 {
			aa.mh.swingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
		}

		sim.rescheduleWeaponAttack(aa.mh.swingAt)

		if aa.IsDualWielding && aa.oh.enabled {
			aa.oh.updateSwingDuration(aa.mh.curSwingSpeed)

			if remainingSwingTime := aa.oh.swingAt - sim.CurrentTime; remainingSwingTime > 0 {
				aa.oh.swingAt = sim.CurrentTime + time.Duration(float64(remainingSwingTime)*f)
			}

			sim.rescheduleWeaponAttack(aa.oh.swingAt)
		}
	}
}

// StopMeleeUntil should be used whenever a non-melee spell is cast. It stops melee, then restarts it
// at end of cast, but with a reset swing timer (as if swings had just landed).
func (aa *AutoAttacks) StopMeleeUntil(sim *Simulation, readyAt time.Duration, desyncOH bool) {
	if !aa.AutoSwingMelee { // if not auto swinging, don't auto restart.
		return
	}

	aa.mh.swingAt = readyAt + aa.mh.curSwingDuration
	sim.rescheduleWeaponAttack(aa.mh.swingAt)

	if aa.IsDualWielding {
		aa.oh.swingAt = readyAt + aa.oh.curSwingDuration
		if desyncOH {
			// Used by warrior to desync offhand after unglyphed Shattering Throw.
			aa.oh.swingAt += aa.oh.curSwingDuration / 2
		}
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
	}
}
func (aa *AutoAttacks) StopRangedUntil(sim *Simulation, readyAt time.Duration) {
	if !aa.AutoSwingRanged { // if not auto swinging, don't auto restart.
		return
	}

	aa.ranged.swingAt = readyAt + aa.ranged.curSwingDuration
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}

// Delays all swing timers for the specified amount.
func (aa *AutoAttacks) DelayMeleeBy(sim *Simulation, delay time.Duration) {
	if delay <= 0 {
		return
	}

	aa.mh.swingAt += delay
	sim.rescheduleWeaponAttack(aa.mh.swingAt)

	if aa.IsDualWielding {
		aa.oh.swingAt += delay
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
	}
}

// PauseMeleeBy will prevent any swing from completing for the specified time.
// This replicates a /stopattack and /startattack with a brief "pause" in the middle.
// It's possible that no swing time is lost if the pauseTime is less than the remaining swing time.
// Used by Rogue Gouge
func (aa *AutoAttacks) PauseMeleeBy(sim *Simulation, pauseTime time.Duration) {
	if !aa.AutoSwingMelee {
		return
	}

	timeToResume := sim.CurrentTime + pauseTime
	if aa.mh.swingAt < timeToResume {
		aa.mh.swingAt = timeToResume
		sim.rescheduleWeaponAttack(aa.mh.swingAt)
	}
	if aa.IsDualWielding && aa.oh.swingAt < timeToResume {
		aa.oh.swingAt = timeToResume
		sim.rescheduleWeaponAttack(aa.oh.swingAt)
	}
}

func (aa *AutoAttacks) DelayRangedUntil(sim *Simulation, readyAt time.Duration) {
	if readyAt <= aa.ranged.swingAt {
		return
	}

	aa.ranged.swingAt = readyAt
	sim.rescheduleWeaponAttack(aa.ranged.swingAt)
}

// Returns the time at which the next attack will occur.
func (aa *AutoAttacks) NextAttackAt() time.Duration {
	return min(aa.mh.swingAt, aa.oh.swingAt)
}

// Used to prevent artificial Haste breakpoints arising from APL evaluations after autos occurring at
// locally optimal timings.
func (aa *AutoAttacks) RandomizeMeleeTiming(sim *Simulation) {
	swingDur := aa.MainhandSwingSpeed()
	randomAutoOffset := DurationFromSeconds(sim.RandomFloat("Melee Timing") * swingDur.Seconds() / 2)
	aa.StopMeleeUntil(sim, sim.CurrentTime-swingDur+randomAutoOffset, true)
}

type DynamicProcManager struct {
	procMasks   []ProcMask
	procChances []float64
}

// Returns whether the effect procced.
func (dpm *DynamicProcManager) Proc(sim *Simulation, procMask ProcMask, label string) bool {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return sim.RandomFloat(label) < dpm.procChances[i]
		}
	}
	return false
}

func (dpm *DynamicProcManager) Chance(procMask ProcMask) float64 {
	for i, m := range dpm.procMasks {
		if m.Matches(procMask) {
			return dpm.procChances[i]
		}
	}
	return 0
}

// PPMManager for static ProcMasks
func (aa *AutoAttacks) NewPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := aa.newDynamicProcManager(ppm, 0, procMask)

	if aa.character != nil {
		aa.character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
			dpm = aa.character.AutoAttacks.newDynamicProcManager(ppm, 0, procMask)
		})
	}

	return &dpm
}

// PPMManager for static ProcMasks and no item swap callback
func (aa *AutoAttacks) NewStaticPPMManager(ppm float64, procMask ProcMask) *DynamicProcManager {
	dpm := aa.newDynamicProcManager(ppm, 0, procMask)

	return &dpm
}

// Dynamic Proc Manager for static ProcMasks and no item swap callback
func (aa *AutoAttacks) NewStaticDynamicProcManager(fixedProcChance float64, procMask ProcMask) *DynamicProcManager {
	dpm := aa.newDynamicProcManager(0, fixedProcChance, procMask)

	return &dpm
}

// Dynamic Proc Manager for dynamic ProcMasks on weapon enchants
func (aa *AutoAttacks) NewDynamicProcManagerForEnchant(effectID int32, ppm float64, fixedProcChance float64) *DynamicProcManager {
	return aa.newDynamicProcManagerWithDynamicProcMask(ppm, fixedProcChance, func() ProcMask {
		return aa.character.getCurrentProcMaskForWeaponEnchant(effectID)
	})
}

// Dynamic Proc Manager for dynamic ProcMasks on weapon effects
func (aa *AutoAttacks) NewDynamicProcManagerForWeaponEffect(itemID int32, ppm float64, fixedProcChance float64) *DynamicProcManager {
	return aa.newDynamicProcManagerWithDynamicProcMask(ppm, fixedProcChance, func() ProcMask {
		return aa.character.getCurrentProcMaskForWeaponEffect(itemID)
	})
}

func (aa *AutoAttacks) newDynamicProcManagerWithDynamicProcMask(ppm float64, fixedProcChance float64, procMaskFn func() ProcMask) *DynamicProcManager {
	dpm := aa.newDynamicProcManager(ppm, fixedProcChance, procMaskFn())

	if aa.character != nil {
		aa.character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
			dpm = aa.character.AutoAttacks.newDynamicProcManager(ppm, fixedProcChance, procMaskFn())
		})
	}

	return &dpm

}

func (aa *AutoAttacks) newDynamicProcManager(ppm float64, fixedProcChance float64, procMask ProcMask) DynamicProcManager {
	if (ppm != 0) && (fixedProcChance != 0) {
		panic("Cannot simultaneously specify both a ppm and a fixed proc chance!")
	}

	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return DynamicProcManager{}
	}

	dpm := DynamicProcManager{procMasks: make([]ProcMask, 0, 2), procChances: make([]float64, 0, 2)}

	mergeOrAppend := func(speed float64, mask ProcMask) {
		if speed == 0 || mask == 0 {
			return
		}

		if i := slices.Index(dpm.procChances, speed); i != -1 {
			dpm.procMasks[i] |= mask
			return
		}

		dpm.procMasks = append(dpm.procMasks, mask)
		dpm.procChances = append(dpm.procChances, speed)
	}

	mergeOrAppend(aa.mh.SwingSpeed, procMask&^ProcMaskRanged&^ProcMaskMeleeOH) // "everything else", even if not explicitly flagged MH
	mergeOrAppend(aa.oh.SwingSpeed, procMask&ProcMaskMeleeOH)
	mergeOrAppend(aa.ranged.SwingSpeed, procMask&ProcMaskRanged)

	for i := range dpm.procChances {
		if fixedProcChance != 0 {
			dpm.procChances[i] = fixedProcChance
		} else {
			dpm.procChances[i] *= ppm / 60
		}
	}

	return dpm
}

// Returns whether a PPM-based effect procced.
// Using NewPPMManager() is preferred; this function should only be used when
// the attacker is not known at initialization time.
func (aa *AutoAttacks) PPMProc(sim *Simulation, ppm float64, procMask ProcMask, label string, spell *Spell) bool {
	if !aa.AutoSwingMelee && !aa.AutoSwingRanged {
		return false
	}

	switch {
	case spell.ProcMask.Matches(procMask &^ ProcMaskMeleeOH &^ ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.mh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskMeleeOH):
		return sim.RandomFloat(label) < ppm*aa.oh.SwingSpeed/60.0
	case spell.ProcMask.Matches(procMask & ProcMaskRanged):
		return sim.RandomFloat(label) < ppm*aa.ranged.SwingSpeed/60.0
	}
	return false
}

func (unit *Unit) applyParryHaste() {
	if !unit.PseudoStats.ParryHaste || !unit.AutoAttacks.AutoSwingMelee {
		return
	}

	unit.RegisterAura(Aura{
		Label:    "Parry Haste",
		Duration: NeverExpires,
		OnReset: func(aura *Aura, sim *Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitTaken: func(aura *Aura, sim *Simulation, spell *Spell, result *SpellResult) {
			if !result.Outcome.Matches(OutcomeParry) {
				return
			}

			remainingTime := aura.Unit.AutoAttacks.mh.swingAt - sim.CurrentTime
			swingSpeed := aura.Unit.AutoAttacks.mh.curSwingDuration
			minRemainingTime := time.Duration(float64(swingSpeed) * 0.2) // 20% of Swing Speed
			defaultReduction := minRemainingTime * 2                     // 40% of Swing Speed

			if remainingTime <= minRemainingTime {
				return
			}

			parryHasteReduction := min(defaultReduction, remainingTime-minRemainingTime)
			newReadyAt := aura.Unit.AutoAttacks.mh.swingAt - parryHasteReduction
			if sim.Log != nil {
				aura.Unit.Log(sim, "MH Swing reduced by %s due to parry haste, will now occur at %s", parryHasteReduction, newReadyAt)
			}

			aura.Unit.AutoAttacks.mh.swingAt = newReadyAt
			sim.rescheduleWeaponAttack(newReadyAt)
		},
	})
}
