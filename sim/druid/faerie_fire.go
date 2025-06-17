package druid

import (
	"time"

	"github.com/wowsims/mop/sim/core"
)

func (druid *Druid) registerFaerieFireSpell() {
	actionID := core.ActionID{SpellID: 770}
	manaCostOptions := core.ManaCostOptions{
		BaseCostPercent: 7.5,
	}
	gcd := core.GCDDefault
	ignoreHaste := false
	cd := core.Cooldown{}
	flatThreatBonus := 48.
	flags := core.SpellFlagAPL
	formMask := Humanoid | Moonkin

	if druid.InForm(Cat | Bear) {
		manaCostOptions = core.ManaCostOptions{}
		formMask = Cat | Bear
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
	}

	if druid.InForm(Cat) {
		gcd = time.Second
		ignoreHaste = true
	}

	druid.FaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WeakenedArmorAura(target)
	})

	druid.FaerieFire = druid.RegisterSpell(formMask, core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       flags,

		ManaCost: manaCostOptions,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: gcd,
			},
			IgnoreHaste: ignoreHaste,
			CD:          cd,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  flatThreatBonus,
		DamageMultiplier: 1,
		CritMultiplier:   druid.DefaultCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.0
			outcome := spell.OutcomeMagicHit

			if druid.InForm(Bear) {
				baseDamage = 10.0 + 0.302*spell.MeleeAttackPower()
				outcome = spell.OutcomeMagicHitAndCrit
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, outcome)

			if result.Landed() {
				druid.TryApplyFaerieFireEffect(sim, target)

				if druid.InForm(Bear) && sim.Proc(0.25, "Mangle CD Reset") {
					druid.MangleBear.CD.Reset()
				}
			}
		},

		RelatedAuraArrays: druid.FaerieFireAuras.ToMap(),
	})
}

func (druid *Druid) CanApplyFaerieFireDebuff(target *core.Unit) bool {
	return druid.FaerieFireAuras.Get(target).IsActive() || !druid.FaerieFireAuras.Get(target).ExclusiveEffects[0].Category.AnyActive()
}

func (druid *Druid) TryApplyFaerieFireEffect(sim *core.Simulation, target *core.Unit) {
	if druid.CanApplyFaerieFireDebuff(target) {
		aura := druid.FaerieFireAuras.Get(target)
		aura.Activate(sim)

		if aura.IsActive() {
			aura.SetStacks(sim, 3)
		}
	}
}

func (druid *Druid) ShouldFaerieFire(sim *core.Simulation, target *core.Unit) bool {
	if druid.FaerieFire == nil {
		return false
	}

	if !druid.FaerieFire.CanCastOrQueue(sim, target) {
		return false
	}

	return druid.FaerieFireAuras.Get(target).ShouldRefreshExclusiveEffects(sim, time.Second*6)
}
