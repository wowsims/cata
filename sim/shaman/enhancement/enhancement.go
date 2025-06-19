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

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, options.TalentsString, selfBuffs, true, enhOptions.ClassOptions.FeleAutocast),
	}

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultCritMultiplier()),
		AutoSwingMelee: true,
	})

	enh.ApplySyncType(enhOptions.SyncType)
	enh.ApplyFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}

	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}

	enh.SpiritWolves = &SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

	enh.PseudoStats.CanParry = true

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

	SpiritWolves *SpiritWolves

	StormStrikeDebuffAuras core.AuraArray
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.UnleashedRage = true
	enh.Shaman.AddRaidBuffs(raidBuffs)
}

func (enh *EnhancementShaman) ApplyTalents() {
	enh.ApplyEnhancementTalents()
	enh.Shaman.ApplyTalents()
	enh.ApplyArmorSpecializationEffect(stats.Agility, proto.ArmorType_ArmorTypeMail, 86529)
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	// In the Initialize due to frost brand adding the aura to the enemy
	enh.RegisterFrostbrandImbue(enh.getImbueProcMask(proto.ShamanImbue_FrostbrandWeapon))
	enh.RegisterFlametongueImbue(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))
	enh.RegisterWindfuryImbue(enh.getImbueProcMask(proto.ShamanImbue_WindfuryWeapon))

	if enh.ItemSwap.IsEnabled() {
		enh.ApplyFlametongueImbueSwap(enh.getImbueProcMask(proto.ShamanImbue_FlametongueWeapon))
		enh.RegisterItemSwapCallback(core.AllWeaponSlots(), func(_ *core.Simulation, slot proto.ItemSlot) {
			enh.ApplySyncType(proto.ShamanSyncType_Auto)
		})
	}

	//Mental Quickness
	enh.GetSpellPowerValue = func(spell *core.Spell) float64 {
		return spell.MeleeAttackPower() * 0.65
	}

	// Mastery: Enhanced Elements
	masteryMod := enh.AddDynamicMod(core.SpellModConfig{
		Kind:              core.SpellMod_DamageDone_Pct,
		School:            core.SpellSchoolFire | core.SpellSchoolFrost | core.SpellSchoolNature,
		ShouldApplyToPets: true,
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

	enh.registerLavaLashSpell()
	enh.registerFireNovaSpell()
	enh.registerStormstrikeSpell()
	enh.registerStormblastSpell()
	enh.registerFeralSpirit()
}

func (enh EnhancementShaman) getMasteryBonus() float64 {
	return 0.16 + 0.02*enh.GetMasteryPoints()
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
