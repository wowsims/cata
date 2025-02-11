package druid

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (druid *Druid) registerFaerieFireSpell() {
	actionID := core.ActionID{SpellID: 770}
	manaCostOptions := core.ManaCostOptions{
		BaseCostPercent: 8,
	}
	gcd := core.GCDDefault
	ignoreHaste := false
	cd := core.Cooldown{}
	flatThreatBonus := 48.
	flags := SpellFlagOmenTrigger
	formMask := Humanoid | Moonkin

	if druid.InForm(Cat | Bear) {
		actionID = core.ActionID{SpellID: 16857}
		manaCostOptions = core.ManaCostOptions{}
		gcd = time.Second
		ignoreHaste = true
		flags = core.SpellFlagNone
		formMask = Cat | Bear
		cd = core.Cooldown{
			Timer:    druid.NewTimer(),
			Duration: time.Second * 6,
		}
	}
	flags |= core.SpellFlagAPL

	druid.FaerieFireAuras = druid.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.FaerieFireAura(target)
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
		CritMultiplier:   druid.DefaultSpellCritMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.0
			outcome := spell.OutcomeMagicHit
			if druid.InForm(Bear) {
				baseDamage = 2950 + 0.108*spell.MeleeAttackPower()
				outcome = spell.OutcomeMagicHitAndCrit
			}

			result := spell.CalcAndDealDamage(sim, target, baseDamage, outcome)
			if result.Landed() {
				druid.TryApplyFaerieFireEffect(sim, target)
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
			aura.SetStacks(sim, aura.GetStacks()+1+druid.Talents.FeralAggression)
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
