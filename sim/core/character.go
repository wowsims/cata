package core

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/core/stats"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CharacterBuildPhase uint8

func (cbp CharacterBuildPhase) Matches(other CharacterBuildPhase) bool {
	return (cbp & other) != 0
}

const (
	CharacterBuildPhaseNone CharacterBuildPhase = 0
	CharacterBuildPhaseBase CharacterBuildPhase = 1 << iota
	CharacterBuildPhaseGear
	CharacterBuildPhaseTalents
	CharacterBuildPhaseBuffs
	CharacterBuildPhaseConsumes
)

const CharacterBuildPhaseAll = CharacterBuildPhaseBase | CharacterBuildPhaseGear | CharacterBuildPhaseTalents | CharacterBuildPhaseBuffs | CharacterBuildPhaseConsumes

// Character is a data structure to hold all the shared values that all
// class logic shares.
// All players have stats, equipment, auras, etc
type Character struct {
	Unit

	Name  string // Different from Label, needed for returned results.
	Race  proto.Race
	Class proto.Class
	Spec  proto.Spec

	// Current gear.
	Equipment

	// Stat buff auras associated with any proc effects in the Character's equippable items
	ItemProcBuffs []*StatBuffAura

	//Item Swap Handler
	ItemSwap ItemSwap

	// Consumables this Character will be using.
	Consumables *proto.ConsumesSpec

	// Base stats for this Character.
	baseStats stats.Stats

	// Bonus stats for this Character, specified in the UI and/or EP
	// calculator
	bonusStats     stats.Stats
	bonusMHDps     float64
	bonusOHDps     float64
	bonusRangedDps float64

	professions [2]proto.Profession

	glyphs [6]int32

	// Used for effects like "Increased Armor Value from Items"
	*EquipScalingManager

	// Provides major cooldown management behavior.
	majorCooldownManager

	// Up reference to this Character's Party.
	Party *Party

	// This character's index within its party [0-4].
	PartyIndex int

	// This stores a timer on spell category ID so that we can track on use effects.
	spellCategoryTimers map[int32]*Timer

	Pets []*Pet // cached in AddPet, for advance()
}

func NewCharacter(party *Party, partyIndex int, player *proto.Player) Character {
	if player.Database != nil {
		addToDatabase(player.Database)
	}

	character := Character{
		Unit: Unit{
			Type:        PlayerUnit,
			Index:       int32(party.Index*5 + partyIndex),
			Level:       CharacterLevel,
			auraTracker: newAuraTracker(),
			PseudoStats: stats.NewPseudoStats(),
			Metrics:     NewUnitMetrics(),

			StatDependencyManager: stats.NewStatDependencyManager(),
			avoidanceParams:       AvoidanceDRByClass[player.Class],

			ReactionTime:            time.Duration(max(player.ReactionTimeMs, 10)) * time.Millisecond,
			ChannelClipDelay:        max(0, time.Duration(player.ChannelClipDelayMs)*time.Millisecond),
			StartDistanceFromTarget: player.DistanceFromTarget,
		},

		Name:  player.Name,
		Race:  player.Race,
		Class: player.Class,
		Spec:  PlayerProtoToSpec(player),

		Equipment: ProtoToEquipment(player.Equipment),

		professions: [2]proto.Profession{
			player.Profession1,
			player.Profession2,
		},

		Party:      party,
		PartyIndex: partyIndex,

		majorCooldownManager: newMajorCooldownManager(player.Cooldowns),
	}
	character.GCD = character.NewTimer()
	character.RotationTimer = character.NewTimer()

	character.Label = fmt.Sprintf("%s (#%d)", character.Name, character.Index+1)

	if player.Glyphs != nil {
		character.glyphs = [6]int32{
			player.Glyphs.Major1,
			player.Glyphs.Major2,
			player.Glyphs.Major3,
			player.Glyphs.Minor1,
			player.Glyphs.Minor2,
			player.Glyphs.Minor3,
		}
	}

	character.Consumables = &proto.ConsumesSpec{}
	if player.Consumables != nil {
		character.Consumables = player.Consumables
	}

	character.baseStats = BaseStats[BaseStatsKey{Race: character.Race, Class: character.Class}]

	character.AddStats(character.baseStats)
	character.addUniversalStatDependencies()

	if player.BonusStats != nil {
		if player.BonusStats.Stats != nil {
			character.bonusStats = stats.FromUnitStatsProto(player.BonusStats)
		}
		if player.BonusStats.PseudoStats != nil {
			ps := player.BonusStats.PseudoStats
			character.bonusMHDps = ps[proto.PseudoStat_PseudoStatMainHandDps]
			character.bonusOHDps = ps[proto.PseudoStat_PseudoStatOffHandDps]
			character.bonusRangedDps = ps[proto.PseudoStat_PseudoStatRangedDps]
			character.PseudoStats.BonusMHDps += character.bonusMHDps
			character.PseudoStats.BonusOHDps += character.bonusOHDps
			character.PseudoStats.BonusRangedDps += character.bonusRangedDps
		}
	}

	if weapon := character.OffHand(); weapon.ID != 0 {
		if weapon.WeaponType == proto.WeaponType_WeaponTypeShield {
			character.PseudoStats.CanBlock = true
		}
	}
	character.PseudoStats.InFrontOfTarget = player.InFrontOfTarget

	if player.EnableItemSwap && player.ItemSwap != nil {
		character.enableItemSwap(player.ItemSwap, character.DefaultCritMultiplier(), character.DefaultCritMultiplier(), character.DefaultCritMultiplier())
	}

	character.EquipScalingManager = character.NewEquipScalingManager()

	return character
}

type EquipScalingManager struct {
	itemStatMultipliers map[stats.Stat]float64
	cachedEquipStats    stats.Stats
	equipStatsApplied   bool
	equipCacheValid     bool
}

func (character *Character) NewEquipScalingManager() *EquipScalingManager {
	return &EquipScalingManager{
		itemStatMultipliers: make(map[stats.Stat]float64),
		cachedEquipStats:    character.Equipment.Stats().Add(character.bonusStats),
		equipCacheValid:     true,
	}
}

func (character *Character) AddDynamicEquipStats(sim *Simulation, equipStats stats.Stats) {
	character.AddStatsDynamic(sim, equipStats.ApplyMultipliers(character.itemStatMultipliers))
	character.equipCacheValid = false
}

func (character *Character) applyEquipScalingInternal(stat stats.Stat, multiplier float64) float64 {
	character.updateCachedEquipStats()
	oldMultiplier, exists := character.itemStatMultipliers[stat]

	if !exists {
		oldMultiplier = 1.0
	}

	newMultiplier := oldMultiplier * multiplier
	character.itemStatMultipliers[stat] = newMultiplier

	return character.cachedEquipStats[stat] * (newMultiplier - oldMultiplier)
}

func (character *Character) ApplyEquipScaling(stat stats.Stat, multiplier float64) {
	statDiff := character.applyEquipScalingInternal(stat, multiplier)
	// Equipment stats already applied, so need to manually at the bonus to
	// the character now to ensure correct values
	if character.equipStatsApplied {
		character.AddStat(stat, statDiff)
	}
}

func (character *Character) ApplyDynamicEquipScaling(sim *Simulation, stat stats.Stat, multiplier float64) {
	if character.Env.MeasuringStats && (character.Env.State != Finalized) {
		character.ApplyEquipScaling(stat, multiplier)
	} else {
		statDiff := character.applyEquipScalingInternal(stat, multiplier)
		character.AddStatDynamic(sim, stat, statDiff)
	}
}

func (character *Character) RemoveEquipScaling(stat stats.Stat, multiplier float64) {
	character.ApplyEquipScaling(stat, 1/multiplier)
}

func (character *Character) RemoveDynamicEquipScaling(sim *Simulation, stat stats.Stat, multiplier float64) {
	character.ApplyDynamicEquipScaling(sim, stat, 1/multiplier)
}

func (character *Character) updateCachedEquipStats() {
	if !character.equipCacheValid {
		character.cachedEquipStats = character.Equipment.Stats().Add(character.bonusStats)
		character.equipCacheValid = true
	}
}

func (character *Character) EquipStats() stats.Stats {
	character.updateCachedEquipStats()
	return character.cachedEquipStats.ApplyMultipliers(character.itemStatMultipliers)
}

func (character *Character) applyEquipment() {
	if character.EquipScalingManager == nil {
		return
	}

	if character.equipStatsApplied {
		panic("Equipment stats already applied to character!")
	}
	character.AddStats(character.EquipStats())
	character.equipStatsApplied = true
}

func (character *Character) addUniversalStatDependencies() {
	character.Unit.addUniversalStatDependencies()
	character.AddStat(stats.Health, 20-14*20)
	character.AddStatDependency(stats.Stamina, stats.Health, 14)
}

// Returns a partially-filled PlayerStats proto for use in the CharacterStats api call.
func (character *Character) applyAllEffects(agent Agent, raidBuffs *proto.RaidBuffs, partyBuffs *proto.PartyBuffs, individualBuffs *proto.IndividualBuffs) *proto.PlayerStats {
	playerStats := &proto.PlayerStats{}

	measureStats := func() *proto.UnitStats {
		baseStats := character.GetStats()
		character.stats = character.SortAndApplyStatDependencies(character.stats)
		measuredStatsProto := &proto.UnitStats{
			Stats:       character.GetStats().ToProtoArray(),
			PseudoStats: character.GetPseudoStatsProto(),
			ApiVersion:  GetCurrentProtoVersion(),
		}
		character.stats = baseStats
		return measuredStatsProto
	}

	applyRaceEffects(agent)
	character.applyProfessionEffects()
	character.applyBuildPhaseAuras(CharacterBuildPhaseBase)
	playerStats.BaseStats = measureStats()

	character.applyEquipment()
	character.applyItemEffects(agent)
	character.applyItemSetBonusEffects(agent)
	character.applyBuildPhaseAuras(CharacterBuildPhaseGear)
	playerStats.GearStats = measureStats()

	agent.ApplyTalents()
	character.applyBuildPhaseAuras(CharacterBuildPhaseTalents)
	playerStats.TalentsStats = measureStats()

	applyBuffEffects(agent, raidBuffs, partyBuffs, individualBuffs)
	character.applyBuildPhaseAuras(CharacterBuildPhaseBuffs)
	playerStats.BuffsStats = measureStats()

	applyConsumeEffects(agent)
	character.applyBuildPhaseAuras(CharacterBuildPhaseConsumes)
	playerStats.ConsumesStats = measureStats()
	character.clearBuildPhaseAuras(CharacterBuildPhaseAll)

	for _, petAgent := range character.PetAgents {
		applyPetBuffEffects(petAgent, raidBuffs, partyBuffs, individualBuffs)
	}

	return playerStats
}
func (character *Character) applyBuildPhaseAuras(phase CharacterBuildPhase) {
	sim := Simulation{}
	character.Env.MeasuringStats = true
	for _, aura := range character.auras {
		if aura.BuildPhase.Matches(phase) {
			aura.Activate(&sim)
		}
	}
	character.Env.MeasuringStats = false
}
func (character *Character) clearBuildPhaseAuras(phase CharacterBuildPhase) {
	sim := Simulation{}
	character.Env.MeasuringStats = true
	for _, aura := range character.auras {
		if aura.BuildPhase.Matches(phase) {
			aura.Deactivate(&sim)
		}
	}
	character.Env.MeasuringStats = false
}
func (character *Character) CalculateMasteryPoints() float64 {
	return character.GetStat(stats.MasteryRating) / MasteryRatingPerMasteryPoint
}

// Apply effects from all equipped core.
func (character *Character) applyItemEffects(agent Agent) {
	registeredItemEffects := make(map[int32]bool)
	registeredItemEnchantEffects := make(map[int32]bool)

	character.Equipment.applyItemEffects(agent, registeredItemEffects, registeredItemEnchantEffects, true)

	if character.ItemSwap.IsEnabled() {
		character.ItemSwap.unEquippedItems.applyItemEffects(agent, registeredItemEffects, registeredItemEnchantEffects, false)
	}
}

func (character *Character) AddPet(pet PetAgent) {
	if character.Env != nil {
		panic("Pets must be added during construction!")
	}

	character.PetAgents = append(character.PetAgents, pet)
	character.Pets = append(character.Pets, pet.GetPet())
}

func (character *Character) GetBaseStats() stats.Stats {
	return character.baseStats
}

// Returns the crit multiplier for a spell.
// https://web.archive.org/web/20081014064638/http://elitistjerks.com/f31/t12595-relentless_earthstorm_diamond_-_melee_only/p4/
// https://github.com/TheGroxEmpire/TBC_DPS_Warrior_Sim/issues/30
// TODO "primaryModifiers" could be modelled as a PseudoStat, since they're unit-specific. "secondaryModifiers" apply to a specific set of spells.
func (character *Character) calculateCritMultiplier(normalCritDamage float64, primaryModifiers float64, secondaryModifiers float64) float64 {
	return 1.0 + (normalCritDamage*primaryModifiers-1.0)*(1.0+secondaryModifiers)
}
func (character *Character) CritMultiplier(primaryModifiers float64, secondaryModifiers float64) float64 {
	return character.calculateCritMultiplier(2.0, primaryModifiers, secondaryModifiers)
}
func (character *Character) DefaultCritMultiplier() float64 {
	return character.CritMultiplier(1, 0)
}

func (character *Character) AddRaidBuffs(_ *proto.RaidBuffs) {
}
func (character *Character) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (character *Character) initialize(agent Agent) {
	character.majorCooldownManager.initialize(character)
	character.ItemSwap.initialize(character)

	character.rotationAction = &PendingAction{
		Priority: ActionPriorityGCD,
		OnAction: func(sim *Simulation) {
			if hc := &character.Hardcast; hc.Expires != startingCDTime && hc.Expires <= sim.CurrentTime {
				hc.Expires = startingCDTime
				if hc.OnComplete != nil {
					hc.OnComplete(sim, hc.Target)
				}
			}

			if sim.CurrentTime < 0 {
				return
			}

			if sim.Options.Interactive {
				if character.GCD.IsReady(sim) {
					sim.NeedsInput = true
				}
				return
			}

			if character.Rotation != nil {
				character.Rotation.DoNextAction(sim)
				return
			}
		},
	}
}

func (character *Character) Finalize() {
	if character.Env.IsFinalized() {
		return
	}

	character.PseudoStats.ParryHaste = character.PseudoStats.CanParry

	character.Unit.finalize()

	character.majorCooldownManager.finalize()
}

func (character *Character) FillPlayerStats(playerStats *proto.PlayerStats) {
	if playerStats == nil {
		return
	}

	character.applyBuildPhaseAuras(CharacterBuildPhaseAll)
	playerStats.FinalStats = &proto.UnitStats{
		Stats:       character.GetStats().ToProtoArray(),
		PseudoStats: character.GetPseudoStatsProto(),
		ApiVersion:  GetCurrentProtoVersion(),
	}

	character.clearBuildPhaseAuras(CharacterBuildPhaseAll)
	playerStats.Sets = character.GetActiveSetBonusNames()

	playerStats.Metadata = character.GetMetadata()
	for _, pet := range character.Pets {
		playerStats.Pets = append(playerStats.Pets, &proto.PetStats{
			Metadata: pet.GetMetadata(),
		})
	}

	if character.Rotation != nil {
		playerStats.RotationStats = character.Rotation.getStats()
	}
}

func (character *Character) reset(sim *Simulation, agent Agent) {
	character.Unit.reset(sim, agent)
	character.majorCooldownManager.reset(sim)
	character.CurrentTarget = character.defaultTarget

	agent.Reset(sim)

	character.ItemSwap.reset(sim)

	for _, petAgent := range character.PetAgents {
		petAgent.GetPet().reset(sim, petAgent)
	}
}

func (character *Character) HasProfession(prof proto.Profession) bool {
	return prof == character.professions[0] || prof == character.professions[1]
}

func (character *Character) HasGlyph(glyphID int32) bool {
	for _, g := range character.glyphs {
		if g == glyphID {
			return true
		}
	}
	return false
}

func (character *Character) HasTrinketEquipped(itemID int32) bool {
	return character.Trinket1().ID == itemID ||
		character.Trinket2().ID == itemID
}

func (character *Character) HasRingEquipped(itemID int32) bool {
	return character.Finger1().ID == itemID || character.Finger2().ID == itemID
}

func (character *Character) HasMetaGemEquipped(gemID int32) bool {
	for _, gem := range character.Head().Gems {
		if gem.ID == gemID {
			return true
		}
	}
	return false
}

// Returns the MH weapon if one is equipped, and null otherwise.
func (character *Character) GetMHWeapon() *Item {
	weapon := character.MainHand()
	if weapon.ID == 0 {
		return nil
	}
	return weapon
}
func (character *Character) HasMHWeapon() bool {
	return character.GetMHWeapon() != nil
}

// Returns the OH weapon if one is equipped, and null otherwise. Note that
// shields / Held-in-off-hand items are NOT counted as weapons in this function.
func (character *Character) GetOHWeapon() *Item {
	weapon := character.OffHand()
	if weapon.ID == 0 ||
		weapon.WeaponType == proto.WeaponType_WeaponTypeShield ||
		weapon.WeaponType == proto.WeaponType_WeaponTypeOffHand {
		return nil
	} else {
		return weapon
	}
}
func (character *Character) HasOHWeapon() bool {
	return character.GetOHWeapon() != nil
}

func (character *Character) HasRangedWeapon() bool {
	return character.Ranged() != nil
}

func (character *Character) GetDynamicProcMaskForWeaponEnchant(effectID int32) *ProcMask {
	return character.getDynamicProcMaskPointer(func() ProcMask {
		return character.getCurrentProcMaskForWeaponEnchant(effectID)
	})
}

func (character *Character) getDynamicProcMaskPointer(procMaskFn func() ProcMask) *ProcMask {
	procMask := procMaskFn()

	character.RegisterItemSwapCallback(AllWeaponSlots(), func(sim *Simulation, slot proto.ItemSlot) {
		procMask = procMaskFn()
	})

	return &procMask
}

func (character *Character) getCurrentProcMaskForWeaponEnchant(effectID int32) ProcMask {
	return character.getCurrentProcMaskFor(func(weapon *Item) bool {
		return weapon.Enchant.EffectID == effectID
	})
}

func (character *Character) GetDynamicProcMaskForWeaponEffect(itemID int32) *ProcMask {
	return character.getDynamicProcMaskPointer(func() ProcMask {
		return character.getCurrentProcMaskForWeaponEffect(itemID)
	})
}

func (character *Character) getCurrentProcMaskForWeaponEffect(itemID int32) ProcMask {
	return character.getCurrentProcMaskFor(func(weapon *Item) bool {
		return weapon.ID == itemID
	})
}

func (character *Character) GetProcMaskForTypes(weaponTypes ...proto.WeaponType) ProcMask {
	return character.getCurrentProcMaskFor(func(weapon *Item) bool {
		return weapon != nil && slices.Contains(weaponTypes, weapon.WeaponType)
	})
}

func (character *Character) GetProcMaskForTypesAndHand(twohand bool, weaponTypes ...proto.WeaponType) ProcMask {
	return character.getCurrentProcMaskFor(func(weapon *Item) bool {
		return weapon != nil && (weapon.HandType == proto.HandType_HandTypeTwoHand) == twohand && slices.Contains(weaponTypes, weapon.WeaponType)
	})
}

func (character *Character) getCurrentProcMaskFor(pred func(item *Item) bool) ProcMask {
	mask := ProcMaskUnknown

	if character == nil {
		return mask
	}

	if pred(character.MainHand()) {
		if character.MainHand().RangedWeaponType > 0 {
			mask |= ProcMaskRanged
		} else {
			mask |= ProcMaskMeleeMH
		}
	}
	if pred(character.OffHand()) {
		mask |= ProcMaskMeleeOH
	}
	return mask
}

func (character *Character) GetProcMaskForWeaponSlot(slot proto.ItemSlot) ProcMask {
	item := character.GetItemBySlot(slot)
	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		if item.RangedWeaponType > 0 {
			return ProcMaskRanged
		} else {
			return ProcMaskMeleeMH
		}
	case proto.ItemSlot_ItemSlotOffHand:
		return ProcMaskMeleeOH
	}

	return ProcMaskEmpty
}

func (character *Character) doneIteration(sim *Simulation) {
	character.ItemSwap.doneIteration(sim)

	// Need to do pets first, so we can add their results to the owners.
	for _, pet := range character.Pets {
		pet.doneIteration(sim)
		character.Metrics.AddFinalPetMetrics(&pet.Metrics)
	}

	character.Unit.doneIteration(sim)
}

func (character *Character) GetPseudoStatsProto() []float64 {
	return []float64{
		proto.PseudoStat_PseudoStatMainHandDps: character.AutoAttacks.MH().DPS(),
		proto.PseudoStat_PseudoStatOffHandDps:  character.AutoAttacks.OH().DPS(),
		proto.PseudoStat_PseudoStatRangedDps:   character.AutoAttacks.Ranged().DPS(),

		// Base values are modified by Enemy attackTables, but we display for LVL 90 enemy as paperdoll default
		proto.PseudoStat_PseudoStatDodgePercent: (character.PseudoStats.BaseDodgeChance + character.GetDiminishedDodgeChance()) * 100,
		proto.PseudoStat_PseudoStatParryPercent: Ternary(character.PseudoStats.CanParry, (character.PseudoStats.BaseParryChance+character.GetDiminishedParryChance())*100, 0),
		proto.PseudoStat_PseudoStatBlockPercent: Ternary(character.PseudoStats.CanBlock, (character.PseudoStats.BaseBlockChance+character.GetDiminishedBlockChance())*100, 0),

		// Used by UI to incorporate multiplicative Haste buffs into final character stats display.
		proto.PseudoStat_PseudoStatRangedSpeedMultiplier: character.PseudoStats.RangedSpeedMultiplier * character.PseudoStats.AttackSpeedMultiplier,
		proto.PseudoStat_PseudoStatMeleeSpeedMultiplier:  character.PseudoStats.MeleeSpeedMultiplier * character.PseudoStats.AttackSpeedMultiplier,
		proto.PseudoStat_PseudoStatCastSpeedMultiplier:   character.PseudoStats.CastSpeedMultiplier,
		proto.PseudoStat_PseudoStatMeleeHastePercent:     (character.TotalMeleeHasteMultiplier() - 1) * 100,
		proto.PseudoStat_PseudoStatRangedHastePercent:    (character.TotalRangedHasteMultiplier() - 1) * 100,
		proto.PseudoStat_PseudoStatSpellHastePercent:     (character.TotalSpellHasteMultiplier() - 1) * 100,

		// School-specific fully buffed Hit/Crit are represented as proper Stats in the back-end so
		// that stat dependencies will work correctly, but are stored as PseudoStats in proto
		// messages. This is done so that the stats arrays embedded in database files and saved
		// Encounter settings can omit these extraneous fields.
		proto.PseudoStat_PseudoStatPhysicalHitPercent:  character.GetStat(stats.PhysicalHitPercent),
		proto.PseudoStat_PseudoStatSpellHitPercent:     character.GetStat(stats.SpellHitPercent),
		proto.PseudoStat_PseudoStatPhysicalCritPercent: character.GetStat(stats.PhysicalCritPercent),
		proto.PseudoStat_PseudoStatSpellCritPercent:    character.GetStat(stats.SpellCritPercent),
	}
}

func (character *Character) GetMetricsProto() *proto.UnitMetrics {
	metrics := character.Metrics.ToProto()
	metrics.Name = character.Name
	metrics.UnitIndex = character.UnitIndex
	metrics.Auras = character.auraTracker.GetMetricsProto()

	metrics.Pets = make([]*proto.UnitMetrics, len(character.Pets))
	for i, pet := range character.Pets {
		metrics.Pets[i] = pet.GetMetricsProto()
	}

	return metrics
}

func (character *Character) GetDefensiveTrinketCD() *Timer {
	return character.GetOrInitSpellCategoryTimer(1190)
}
func (character *Character) GetOffensiveTrinketCD() *Timer {
	return character.GetOrInitSpellCategoryTimer(1141)
}
func (character *Character) GetConjuredCD() *Timer {
	return character.GetOrInitSpellCategoryTimer(30)
}
func (character *Character) GetPotionCD() *Timer {
	return character.GetOrInitSpellCategoryTimer(4)
}

func (character *Character) AddStatProcBuff(effectID int32, procAura *StatBuffAura, isEnchant bool, eligibleSlots []proto.ItemSlot) {
	hasEquippedCheck := Ternary(isEnchant, character.Equipment.containsEnchantInSlots, character.Equipment.containsItemInSlots)

	procAura.IsSwapped = !hasEquippedCheck(effectID, eligibleSlots)
	character.ItemProcBuffs = append(character.ItemProcBuffs, procAura)

	character.RegisterItemSwapCallback(eligibleSlots, func(sim *Simulation, slot proto.ItemSlot) {
		procAura.IsSwapped = !hasEquippedCheck(effectID, eligibleSlots)
	})

}

func (character *Character) GetMatchingItemProcAuras(statTypesToMatch []stats.Stat, minIcd time.Duration) []*StatBuffAura {
	includeIcdFilter := (minIcd > 0)
	return FilterSlice(character.ItemProcBuffs, func(aura *StatBuffAura) bool {
		return aura.BuffsMatchingStat(statTypesToMatch) && (!includeIcdFilter || ((aura.Icd != nil) && (aura.Icd.Duration > minIcd)))
	})
}

// Uses proto reflection to set fields in a talents proto (e.g. MageTalents,
// WarriorTalents) based on a talentsStr.
func FillTalentsProto(data protoreflect.Message, talentsStr string) {
	fieldDescriptors := data.Descriptor().Fields()

	for talentIdx, talentValStr := range talentsStr {
		talentVal, _ := strconv.Atoi(string(talentValStr))
		if talentVal != 0 {
			talentOffset := talentIdx*3 + talentVal
			fd := fieldDescriptors.ByNumber(protowire.Number(talentOffset))
			if fd == nil {
				panic(fmt.Sprintf("Couldn't find proto field for talent #%d, full string: %s", talentOffset, talentsStr))
			}
			data.Set(fd, protoreflect.ValueOfBool(true))
		}
	}
}

func (character *Character) MeetsArmorSpecializationRequirement(armorType proto.ArmorType) bool {
	for _, itemSlot := range ArmorSpecializationSlots() {
		item := character.Equipment[itemSlot]
		if item.ArmorType == proto.ArmorType_ArmorTypeUnknown {
			continue
		}
		if item.ArmorType != armorType {
			return false
		}
	}

	return true
}

func (character *Character) RegisterArmorSpecializationTracker(armorType proto.ArmorType, spellID int32) *Aura {
	isEnabled := character.MeetsArmorSpecializationRequirement(armorType)

	aura := character.RegisterAura(Aura{
		Label:      "Armor Specialization",
		ActionID:   ActionID{SpellID: spellID},
		BuildPhase: Ternary(isEnabled, CharacterBuildPhaseTalents, CharacterBuildPhaseNone),
		Duration:   NeverExpires,
	})

	if isEnabled {
		aura = MakePermanent(aura)
	}

	character.RegisterItemSwapCallback(ArmorSpecializationSlots(),
		func(sim *Simulation, _ proto.ItemSlot) {
			if character.MeetsArmorSpecializationRequirement(armorType) {
				if !aura.IsActive() {
					aura.Activate(sim)
				}
			} else {
				aura.Deactivate(sim)
			}
		})

	return aura
}

func (character *Character) ApplyArmorSpecializationEffect(primaryStat stats.Stat, armorType proto.ArmorType, spellID int32) *Aura {
	armorSpecializationDependency := character.NewDynamicMultiplyStat(primaryStat, 1.05)
	trackerAura := character.RegisterArmorSpecializationTracker(armorType, spellID)
	trackerAura.AttachStatDependency(armorSpecializationDependency)
	return trackerAura
}
