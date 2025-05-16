package core

import (
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
)

func applyRaceEffects(agent Agent) {
	character := agent.GetCharacter()

	switch character.Race {
	case proto.Race_RaceBloodElf:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.PseudoStats.ReducedFireHitTakenChance += 0.02
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02

		var actionID ActionID

		var resourceMetrics *ResourceMetrics = nil
		if resourceMetrics == nil {
			if character.HasRunicPowerBar() {
				actionID = ActionID{SpellID: 50613}
				resourceMetrics = character.NewRunicPowerMetrics(actionID)
			} else if character.HasEnergyBar() {
				actionID = ActionID{SpellID: 25046}
				resourceMetrics = character.NewEnergyMetrics(actionID)
			} else if character.HasManaBar() {
				actionID = ActionID{SpellID: 28730}
				resourceMetrics = character.NewManaMetrics(actionID)
			} else if character.HasRageBar() {
				actionID = ActionID{SpellID: 69179}
				resourceMetrics = character.NewRageMetrics(actionID)
			} else if character.HasFocusBar() {
				actionID = ActionID{SpellID: 80483}
				resourceMetrics = character.NewFocusMetrics(actionID)
			}
		}

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, spell *Spell) {
				if spell.Unit.HasRunicPowerBar() {
					spell.Unit.AddRunicPower(sim, 15.0, resourceMetrics)
				} else if spell.Unit.HasEnergyBar() {
					spell.Unit.AddEnergy(sim, 15.0, resourceMetrics)
				} else if spell.Unit.HasManaBar() {
					spell.Unit.AddMana(sim, spell.Unit.MaxMana()*0.06, resourceMetrics)
				} else if spell.Unit.HasRageBar() {
					spell.Unit.AddRage(sim, 15.0, resourceMetrics)
				} else if spell.Unit.HasFocusBar() {
					spell.Unit.AddFocus(sim, 15.0, resourceMetrics)
				}
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell:    spell,
			Type:     CooldownTypeDPS,
			Priority: CooldownPriorityLow,
			ShouldActivate: func(sim *Simulation, character *Character) bool {
				if spell.Unit.HasRunicPowerBar() {
					return character.CurrentRunicPower() <= character.maxRunicPower-15
				} else if spell.Unit.HasEnergyBar() {
					return character.CurrentEnergy() <= character.maxEnergy-15
				} else if spell.Unit.HasRageBar() {
					return character.CurrentRage() <= MaxRage-15
				} else if spell.Unit.HasFocusBar() {
					return character.CurrentFocus() <= character.maxFocus-15
				}
				return true
			},
		})
	case proto.Race_RaceDraenei:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
		character.AddStats(stats.Stats{
			stats.PhysicalHitPercent: 1,
			stats.SpellHitPercent:    1,
		})
		// TODO: Gift of the naaru for healers
	case proto.Race_RaceDwarf:
		character.PseudoStats.ReducedFrostHitTakenChance += 0.02

		// Crack Shot: 1% Expertise with Ranged Weapons
		ranged := character.Ranged()
		if ranged != nil && (ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow ||
			ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun ||
			ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeCrossbow) {
			character.AddStat(stats.ExpertiseRating, ExpertisePerQuarterPercentReduction*4)
		}

		applyWeaponSpecialization(character, 3*ExpertisePerQuarterPercentReduction,
			proto.WeaponType_WeaponTypeMace)

		actionID := ActionID{SpellID: 20594}

		statDep := character.NewDynamicMultiplyStat(stats.Armor, 1.1)
		stoneFormAura := character.NewTemporaryStatsAuraWrapped("Stoneform", actionID, stats.Stats{}, time.Second*8, func(aura *Aura) {
			aura.ApplyOnGain(func(aura *Aura, sim *Simulation) {
				aura.Unit.EnableDynamicStatDep(sim, statDep)
			})
			aura.ApplyOnExpire(func(aura *Aura, sim *Simulation) {
				aura.Unit.DisableDynamicStatDep(sim, statDep)
			})
		})

		spell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,
			Flags:    SpellFlagNoOnCastComplete,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				stoneFormAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: spell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceGnome:
		character.PseudoStats.ReducedArcaneHitTakenChance += 0.02
		character.MultiplyStat(stats.Mana, 1.05)
		applyOneHandWeaponSpecialization(character, 3*ExpertisePerQuarterPercentReduction,
			proto.WeaponType_WeaponTypeSword, proto.WeaponType_WeaponTypeDagger)
	case proto.Race_RaceHuman:
		character.MultiplyStat(stats.Spirit, 1.03)
		applyWeaponSpecialization(character, 3*ExpertisePerQuarterPercentReduction,
			proto.WeaponType_WeaponTypeMace, proto.WeaponType_WeaponTypeSword)
	case proto.Race_RaceNightElf:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.PseudoStats.ReducedPhysicalHitTakenChance += 0.02

		// Shadowmeld
		actionID := ActionID{SpellID: 58984}

		shmeldAura := character.RegisterAura(Aura{
			Label:    "Shadowmeld",
			ActionID: actionID,
			Duration: NeverExpires,
			// Shadowmeld counts as a stealth stance for (most?) spell requirements, but does not enable any additional bonuses/auras tied to stealth.
			// Implementation vaguely mirrors Rogue Vanish
			OnCastComplete: func(aura *Aura, sim *Simulation, spell *Spell) {
				if spell.ActionID != actionID {
					aura.Deactivate(sim)
				}
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				if character.AutoAttacks.MHConfig() != nil {
					character.AutoAttacks.EnableAutoSwing(sim)
				}
			},
		})

		shmeldSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},

			ApplyEffects: func(sim *Simulation, target *Unit, spell *Spell) {
				shmeldAura.Activate(sim)
				character.AutoAttacks.CancelAutoSwing(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: shmeldSpell,
			Type:  CooldownTypeUnknown,
			ShouldActivate: func(s *Simulation, c *Character) bool {
				// No reason to auto-cast this
				return false
			},
		})

	case proto.Race_RaceOrc:
		// Command (Pet damage +5%)
		for _, pet := range character.Pets {
			pet.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Blood Fury
		actionID := ActionID{SpellID: 33697}
		apBonus := 0.0
		spBonus := 0.0

		switch character.Class {
		case proto.Class_ClassMage:
			spBonus = 584.5
		case proto.Class_ClassWarlock:
			spBonus = 584.5
		case proto.Class_ClassShaman:
			spBonus = 584.5
			apBonus = 1169.0
		default:
			apBonus = 1169.0
		}

		buffStats := stats.Stats{stats.AttackPower: apBonus, stats.RangedAttackPower: apBonus, stats.SpellPower: spBonus}
		RegisterTemporaryStatsOnUseCD(character, "Blood Fury", buffStats, time.Second*15, SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 2,
				},
			},
		})

		// Axe specialization
		applyWeaponSpecialization(character, 3*ExpertisePerQuarterPercentReduction,
			proto.WeaponType_WeaponTypeAxe, proto.WeaponType_WeaponTypeFist)
	case proto.Race_RaceTauren:
		character.PseudoStats.ReducedNatureHitTakenChance += 0.02
		character.AddStat(stats.Health, character.GetBaseStats()[stats.Health]*0.05)
	case proto.Race_RaceTroll:
		// Dead Eye: 1% Expertise with Guns, Bows or Crossbows.
		ranged := character.Ranged()
		if ranged != nil && (ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeBow ||
			ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeGun ||
			ranged.RangedWeaponType == proto.RangedWeaponType_RangedWeaponTypeCrossbow) {
			character.AddStat(stats.ExpertiseRating, ExpertisePerQuarterPercentReduction*4)
		}

		// Beast Slaying (+5% damage to beasts)
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.DamageDealtMultiplier *= 1.05
		}

		// Berserking
		actionID := ActionID{SpellID: 26297}

		berserkingAura := character.RegisterAura(Aura{
			Label:    "Berserking (Troll)",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1.2)
				character.MultiplyAttackSpeed(sim, 1.2)
				character.MultiplyResourceRegenSpeed(sim, 1.2)
			},
			OnExpire: func(aura *Aura, sim *Simulation) {
				character.MultiplyCastSpeed(1 / 1.2)
				character.MultiplyAttackSpeed(sim, 1/1.2)
				character.MultiplyResourceRegenSpeed(sim, 1/1.2)
			},
		})

		berserkingSpell := character.RegisterSpell(SpellConfig{
			ActionID: actionID,

			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				berserkingAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(MajorCooldown{
			Spell: berserkingSpell,
			Type:  CooldownTypeDPS,
		})
	case proto.Race_RaceUndead:
		character.PseudoStats.ReducedShadowHitTakenChance += 0.02
	case proto.Race_RaceWorgen:
		character.AddStat(stats.PhysicalCritPercent, 1)
		character.AddStat(stats.SpellCritPercent, 1)
	case proto.Race_RaceGoblin:
		character.PseudoStats.MeleeSpeedMultiplier *= 1.01
		character.PseudoStats.RangedSpeedMultiplier *= 1.01
		character.PseudoStats.CastSpeedMultiplier *= 1.01
	}
}

func applyWeaponSpecialization(character *Character, expertiseBonus float64, weaponTypes ...proto.WeaponType) {
	mask := character.GetProcMaskForTypes(weaponTypes...)

	if mask == ProcMaskMelee || (mask == ProcMaskMeleeMH && !character.HasOHWeapon()) {
		character.AddStat(stats.ExpertiseRating, expertiseBonus)
	} else {
		character.OnSpellRegistered(func(spell *Spell) {
			if spell.ProcMask.Matches(mask) {
				spell.BonusExpertiseRating += expertiseBonus
			}
		})
	}
}

func applyOneHandWeaponSpecialization(character *Character, expertiseBonus float64, weaponTypes ...proto.WeaponType) {
	mask := character.GetProcMaskForTypesAndHand(false, weaponTypes...)

	if mask == ProcMaskMelee || (mask == ProcMaskMeleeMH && !character.HasOHWeapon()) {
		character.AddStat(stats.ExpertiseRating, expertiseBonus)
	} else {
		character.OnSpellRegistered(func(spell *Spell) {
			if spell.ProcMask.Matches(mask) {
				spell.BonusExpertiseRating += expertiseBonus
			}
		})
	}
}
