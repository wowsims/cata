package enhancement

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"github.com/wowsims/mop/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		proto.Spec_SpecEnhancementShaman,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewEnhancementShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_EnhancementShaman)
			if !ok {
				panic("Invalid spec value for Enhancement Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewEnhancementShaman(character *core.Character, options *proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman().Options

	selfBuffs := shaman.SelfBuffs{
		Shield:  enhOptions.ClassOptions.Shield,
		ImbueMH: enhOptions.ClassOptions.ImbueMh,
		ImbueOH: enhOptions.ImbueOh,
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.ClassOptions.Totems != nil {
		totems = enhOptions.ClassOptions.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, totems, selfBuffs, true),
	}

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	enh.ApplySyncType(enhOptions.SyncType)
	// enh.ApplyFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	// enh.SpiritWolves = &shaman.SpiritWolves{
	// 	SpiritWolf1: enh.NewSpiritWolf(1),
	// 	SpiritWolf2: enh.NewSpiritWolf(2),
	// }

	return enh
}

func (enh *EnhancementShaman) getImbueProcMask(imbue proto.ShamanImbue) core.ProcMask {
	var mask core.ProcMask
	if enh.SelfBuffs.ImbueMH == imbue {
		mask |= core.ProcMaskMeleeMH
	}
	if enh.SelfBuffs.ImbueOH == imbue {
		mask |= core.ProcMaskMeleeOH
	}
	return mask
}

type EnhancementShaman struct {
	*shaman.Shaman
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) ApplyTalents() {
	enh.Shaman.ApplyTalents()
	enh.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeMail, 86529)
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	// In the Initialize due to frost brand adding the aura to the enemy
	// enh.RegisterFrostbrandImbue(enh.getImbueProcMask(proto.ShamanImbue_FrostbrandWeapon))
	// enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))
	// enh.RegisterWindfuryImbue(enh.getImbueProcMask(proto.ShamanImbue_WindfuryWeapon))

	if enh.ItemSwap.IsEnabled() {
		// enh.ApplyFlametongueImbueSwap(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))
		enh.RegisterItemSwapCallback(core.MeleeWeaponSlots(), func(_ *core.Simulation, slot proto.ItemSlot) {
			enh.ApplySyncType(proto.ShamanSyncType_Auto)
		})
	}

	enh.GetSpellPowerValue = func(spell *core.Spell) float64 {
		return spell.MeleeAttackPower() * 0.55
	}

	// Mastery: Enhanced Elements
	masteryMod := enh.AddDynamicMod(core.SpellModConfig{
		Kind:   core.SpellMod_DamageDone_Pct,
		School: core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
	})

	enh.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(enh.getMasteryBonus())
	})

	core.MakePermanent(enh.GetOrRegisterAura(core.Aura{
		Label:    "Mastery: Enhanced Elements",
		ActionID: core.ActionID{SpellID: 77223},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.UpdateFloatValue(enh.getMasteryBonus())
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	enh.applyPrimalWisdom()
	// enh.registerLavaLashSpell()
}

func (enh EnhancementShaman) getMasteryBonus() float64 {
	return 0.2 + 0.025*enh.GetMasteryPoints()
}

func (enh *EnhancementShaman) applyPrimalWisdom() {
	manaMetrics := enh.NewManaMetrics(core.ActionID{SpellID: 63375})

	enh.RegisterAura(core.Aura{
		Label:    "Primal Wisdom",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if sim.RandomFloat("Primal Wisdom") < 0.4 {
				enh.AddMana(sim, 0.05*enh.BaseMana, manaMetrics)
			}
		},
	})
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)
}

func (enh *EnhancementShaman) AutoSyncWeapons() proto.ShamanSyncType {
	if mh, oh := enh.MainHand(), enh.OffHand(); mh.SwingSpeed != oh.SwingSpeed {
		return proto.ShamanSyncType_NoSync
	}
	return proto.ShamanSyncType_SyncMainhandOffhandSwings
}

func (enh *EnhancementShaman) ApplySyncType(syncType proto.ShamanSyncType) {
	const FlurryICD = time.Millisecond * 500

	if syncType == proto.ShamanSyncType_Auto {
		syncType = enh.AutoSyncWeapons()
	}

	switch syncType {
	case proto.ShamanSyncType_SyncMainhandOffhandSwings:
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed(); nextMHSwingAt != aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	case proto.ShamanSyncType_DelayOffhandSwings:
		enh.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
			if aa := &enh.AutoAttacks; aa.OffhandSwingAt()-sim.CurrentTime > FlurryICD {
				if nextMHSwingAt := sim.CurrentTime + aa.MainhandSwingSpeed() + 100*time.Millisecond; nextMHSwingAt > aa.OffhandSwingAt() {
					aa.SetOffhandSwingAt(nextMHSwingAt)
				}
			}
			return mhSwingSpell
		})
	default:
		enh.AutoAttacks.SetReplaceMHSwing(nil)
	}
}
