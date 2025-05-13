package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

// Currently Imbues are carried over on item swap
// func (shaman *Shaman) RegisterOnItemSwapWithImbue(effectID int32, procMask *core.ProcMask, aura *core.Aura) {
// 	shaman.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(sim *core.Simulation, slot proto.ItemSlot) {
// 		mask := core.ProcMaskUnknown
// 		if shaman.MainHand().TempEnchant == effectID {
// 			mask |= core.ProcMaskMeleeMH
// 		}
// 		if shaman.OffHand().TempEnchant == effectID {
// 			mask |= core.ProcMaskMeleeOH
// 		}
// 		*procMask = mask

// 		if mask == core.ProcMaskUnknown {
// 			aura.Deactivate(sim)
// 		} else {
// 			aura.Activate(sim)
// 		}
// 	})
// }

func (shaman *Shaman) newWindfuryImbueSpell(isMH bool) *core.Spell {
	apBonus := 4430.0

	tag := 1
	procMask := core.ProcMaskMeleeMHSpecial
	weaponDamageFunc := shaman.MHWeaponDamage
	if !isMH {
		tag = 2
		procMask = core.ProcMaskMeleeOHSpecial
		weaponDamageFunc = shaman.OHWeaponDamage
		apBonus *= 2 // applied after 50% offhand penalty
	}

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 8232, Tag: int32(tag)},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagPassiveSpell,

		DamageMultiplier: []float64{1, 1.20, 1.40}[shaman.Talents.ElementalWeapons],
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mAP := spell.MeleeAttackPower() + apBonus

			baseDamage1 := weaponDamageFunc(sim, mAP)
			baseDamage2 := weaponDamageFunc(sim, mAP)
			baseDamage3 := weaponDamageFunc(sim, mAP)
			result1 := spell.CalcDamage(sim, target, baseDamage1, spell.OutcomeMeleeSpecialHitAndCrit)
			result2 := spell.CalcDamage(sim, target, baseDamage2, spell.OutcomeMeleeSpecialHitAndCrit)
			result3 := spell.CalcDamage(sim, target, baseDamage3, spell.OutcomeMeleeSpecialHitAndCrit)
			spell.DealDamage(sim, result1)
			spell.DealDamage(sim, result2)
			spell.DealDamage(sim, result3)
		},
	}

	return shaman.RegisterSpell(spellConfig)
}

func (shaman *Shaman) RegisterWindfuryImbue(procMask core.ProcMask) {
	enchantID := int32(283)
	if procMask == core.ProcMaskUnknown {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = enchantID
		if shaman.ItemSwap.IsEnabled() {
			weapon := shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand)
			if weapon != nil {
				weapon.TempEnchant = enchantID
			}
		}
	}

	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = enchantID
		if shaman.ItemSwap.IsEnabled() {
			weapon := shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotOffHand)
			if weapon != nil {
				weapon.TempEnchant = enchantID
			}
		}
	}

	var proc = 0.2
	if procMask == core.ProcMaskMelee {
		proc = 0.36
	}
	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfWindfuryWeapon) {
		proc += 0.02
	}

	mhSpell := shaman.newWindfuryImbueSpell(true)
	ohSpell := shaman.newWindfuryImbueSpell(false)

	core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
		Name:       "Windfury Imbue",
		ProcMask:   procMask,
		ICD:        time.Second * 3,
		ProcChance: proc,
		Outcome:    core.OutcomeLanded,
		Callback:   core.CallbackOnSpellHitDealt,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.IsMH() {
				mhSpell.Cast(sim, result.Target)
			} else {
				ohSpell.Cast(sim, result.Target)
			}
		},
	})

	// Currently Imbues are carried over on item swap
	// shaman.RegisterOnItemSwapWithImbue(enchantID, &procMask, aura)
}

// TODO: Not sure on the base damage here wowhead does not seem to be correct. in testing with 1.3 weapon and 129 sp it was 109 damage
func (shaman *Shaman) newFlametongueImbueSpell(weapon *core.Item) *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: int32(8024)},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamageProc,
		ClassSpellMask: SpellMaskFlametongueWeapon,
		Flags:          core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if weapon.SwingSpeed != 0 {
				damage := weapon.SwingSpeed * (68.5 + 0.08/2.6*spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (shaman *Shaman) ApplyFlametongueImbueToItem(item *core.Item) {
	enchantID := int32(5)
	if item == nil || item.TempEnchant == enchantID {
		return
	}

	//flametongue imbue does not stack
	if (shaman.HasMHWeapon() && shaman.GetMHWeapon().TempEnchant == enchantID) || (shaman.HasOHWeapon() && shaman.GetOHWeapon().TempEnchant == enchantID) {
		item.TempEnchant = enchantID
		return

	}
	if shaman.ItemSwap.IsEnabled() && (shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand).TempEnchant == int32(enchantID) || shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotOffHand).TempEnchant == int32(enchantID)) {
		item.TempEnchant = enchantID
		return
	}
	magicDamageBonus := 1.0 + (0.05 * (1 + 0.2*float64(shaman.Talents.ElementalWeapons)))

	shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= magicDamageBonus
	shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= magicDamageBonus
	shaman.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= magicDamageBonus

	if shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfFlametongueWeapon) {
		shaman.AddStat(stats.SpellCritPercent, 2)
	}

	item.TempEnchant = enchantID
}

func (shaman *Shaman) ApplyFlametongueImbue(procMask core.ProcMask) {
	if procMask.Matches(core.ProcMaskMeleeMH) && shaman.HasMHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.MainHand())
	}

	if procMask.Matches(core.ProcMaskMeleeOH) && shaman.HasOHWeapon() {
		shaman.ApplyFlametongueImbueToItem(shaman.OffHand())
	}
}

func (shaman *Shaman) ApplyFlametongueImbueSwap(procMask core.ProcMask) {
	if !shaman.ItemSwap.IsEnabled() {
		return
	}
	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.ApplyFlametongueImbueToItem(shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotMainHand))
	}

	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.ApplyFlametongueImbueToItem(shaman.ItemSwap.GetUnequippedItemBySlot(proto.ItemSlot_ItemSlotOffHand))
	}
}

func (shaman *Shaman) RegisterFlametongueImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown && !shaman.ItemSwap.IsEnabled() {
		return
	}

	for _, itemSlot := range core.MeleeWeaponSlots() {
		var weapon *core.Item
		var triggerProcMask core.ProcMask
		switch {
		case procMask.Matches(core.ProcMaskMeleeMH) && itemSlot == proto.ItemSlot_ItemSlotMainHand:
			weapon = shaman.MainHand()
			triggerProcMask = core.ProcMaskMeleeMH | core.ProcMaskMeleeProc
		case procMask.Matches(core.ProcMaskMeleeOH) && itemSlot == proto.ItemSlot_ItemSlotOffHand:
			weapon = shaman.OffHand()
			triggerProcMask = core.ProcMaskMeleeOH
		}

		if weapon == nil {
			continue
		}

		flameTongueSpell := shaman.newFlametongueImbueSpell(weapon)
		core.MakeProcTriggerAura(&shaman.Unit, core.ProcTrigger{
			Name:     fmt.Sprintf("Flametongue Imbue %s", itemSlot),
			ProcMask: triggerProcMask,
			Outcome:  core.OutcomeLanded,
			Callback: core.CallbackOnSpellHitDealt,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				flameTongueSpell.Cast(sim, result.Target)
			},
		})

	}

	// Currently Imbues are carried over on item swap
	// shaman.RegisterOnItemSwapWithImbue(5, &procMask, aura)
}

func (shaman *Shaman) frostbrandDDBCHandler(sim *core.Simulation, spell *core.Spell, attackTable *core.AttackTable) float64 {
	if spell.ClassSpellMask&(SpellMaskLightningBolt|SpellMaskChainLightning|SpellMaskLavaLash|SpellMaskEarthShock|SpellMaskFlameShock|SpellMaskFrostShock) > 0 {
		return 1 + 0.05*float64(shaman.Talents.FrozenPower)
	}
	return 1.0
}

func (shaman *Shaman) FrostbrandDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Frostbrand Attack-" + shaman.Label,
		ActionID: core.ActionID{SpellID: 8034},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.EnableDamageDoneByCaster(DDBC_FrostbrandWeapon, DDBC_Total, shaman.AttackTables[aura.Unit.UnitIndex], shaman.frostbrandDDBCHandler)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.DisableDamageDoneByCaster(DDBC_FrostbrandWeapon, shaman.AttackTables[aura.Unit.UnitIndex])
		},
	})
}

func (shaman *Shaman) newFrostbrandImbueSpell() *core.Spell {
	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 8033},
		SpellSchool: core.SpellSchoolFrost,

		// TODO: Is this correct?
		ProcMask: core.ProcMaskEmpty,
		Flags:    core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		CritMultiplier:   shaman.DefaultCritMultiplier(),
		ThreatMultiplier: 1,
		BonusCoefficient: 0.1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := shaman.ClassSpellScaling * 0.60900002718 //spell id 8034
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (shaman *Shaman) RegisterFrostbrandImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskUnknown {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.MainHand().TempEnchant = 2
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.OffHand().TempEnchant = 2
	}

	dpm := shaman.AutoAttacks.NewPPMManager(9.0, procMask)

	mhSpell := shaman.newFrostbrandImbueSpell()
	ohSpell := shaman.newFrostbrandImbueSpell()

	fbDebuffAuras := shaman.NewEnemyAuraArray(shaman.FrostbrandDebuffAura)

	aura := shaman.RegisterAura(core.Aura{
		Label:    "Frostbrand Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if dpm.Proc(sim, spell.ProcMask, "Frostbrand Weapon") {
				if spell.IsMH() {
					mhSpell.Cast(sim, result.Target)
				} else {
					ohSpell.Cast(sim, result.Target)
				}
				fbDebuffAuras.Get(result.Target).Activate(sim)
			}
		},
	})

	shaman.ItemSwap.RegisterEnchantProc(2, aura)
}

func (shaman *Shaman) newEarthlivingImbueSpell() *core.Spell {
	glyphBonus := core.Ternary(shaman.HasPrimeGlyph(proto.ShamanPrimeGlyph_GlyphOfEarthlivingWeapon), 1.2, 1.0)

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51730},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagPassiveSpell,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Hot: core.DotConfig{
			Aura: core.Aura{
				Label:    "Earthliving",
				ActionID: core.ActionID{SpellID: 51945},
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = (shaman.ClassSpellScaling*0.57400000095 + (0.038 * dot.Spell.HealingPower(target))) * glyphBonus
				dot.SnapshotAttackerMultiplier = dot.Spell.CasterHealingMultiplier()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Hot(target).Apply(sim)
		},
	})
}

func (shaman *Shaman) ApplyEarthlivingImbueToItem(item *core.Item) {
	enchantId := int32(3345)

	if item == nil || item.TempEnchant == enchantId {
		return
	}

	spBonus := 532.0 * (1.0 + float64(shaman.Talents.ElementalWeapons)*0.20)

	newStats := stats.Stats{stats.SpellPower: spBonus}
	item.Stats = item.Stats.Add(newStats)
	item.TempEnchant = enchantId
}

func (shaman *Shaman) RegisterEarthlivingImbue(procMask core.ProcMask) {
	if procMask == core.ProcMaskEmpty && !shaman.ItemSwap.IsEnabled() {
		return
	}

	if procMask.Matches(core.ProcMaskMeleeMH) {
		shaman.ApplyEarthlivingImbueToItem(shaman.MainHand())
	}
	if procMask.Matches(core.ProcMaskMeleeOH) {
		shaman.ApplyEarthlivingImbueToItem(shaman.OffHand())
	}

	imbueSpell := shaman.newEarthlivingImbueSpell()

	shaman.RegisterAura(core.Aura{
		Label:    "Earthliving Imbue",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnHealDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell != shaman.ChainHeal && spell != shaman.HealingSurge && spell != shaman.HealingWave && spell != shaman.Riptide {
				return
			}

			if procMask.Matches(core.ProcMaskMeleeMH) && sim.RandomFloat("earthliving") < 0.2 {
				imbueSpell.Cast(sim, result.Target)
			}

			if procMask.Matches(core.ProcMaskMeleeOH) && sim.RandomFloat("earthliving") < 0.2 {
				imbueSpell.Cast(sim, result.Target)
			}
		},
	})

	// Currently Imbues are carried over on item swap
	// shaman.RegisterOnItemSwapWithImbue(3350, &procMask, aura)
}
