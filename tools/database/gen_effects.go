package database

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/wowsims/mop/sim/common"
	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools/database/dbc"
	"github.com/wowsims/mop/tools/tooltip"
)

// Sets the minimum itemlevel that should be considered for this expansions
const MIN_EFFECT_ILVL = 416

type ProcInfo struct {
	Outcome  core.HitOutcome
	Callback core.AuraCallback
	ProcMask core.ProcMask
}

// Entry represents a single effect with its ID and display name.
type Variant struct {
	ID   int
	Name string
}

type Entry struct {
	Variants  []*Variant
	Tooltip   []string
	ProcInfo  ProcInfo
	Supported bool
	Harmful   bool
}

// Group holds a category of effects.
type Group struct {
	Name    string
	Entries []*Entry
}

var missingEffectsMap = map[string][]Variant{
	"EnchantEffects": {},
	"ItemEffects":    {},
}

type EffectParseResult byte

const (
	EffectParseResultInvalid     EffectParseResult = iota // Returned when the effect is invalid for the current parameters
	EffectParseResultUnsupported                          // Returned when the effect could be parsed but is not supported for effect generation
	EffectParseResultSuccess                              // Returned when the effect was parsed successfuly
)

func GenerateEffectsFile(groups []*Group, outFile string, templateString string) error {
	if _, err := os.Stat(outFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to check file %s: %w", outFile, err)
	}

	// Ensure groups and entries are sorted
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	for _, grp := range groups {
		sort.Slice(grp.Entries, func(i, j int) bool {
			if grp.Entries[i].Supported != grp.Entries[j].Supported {
				return !grp.Entries[i].Supported
			}

			return grp.Entries[i].Variants[0].ID < grp.Entries[j].Variants[0].ID
		})
	}

	funcMap := map[string]any{
		"asCoreCallback": asCoreCallback,
		"asCoreProcMask": asCoreProcMask,
		"asCoreOutcome":  asCoreOutcome,
		"formatStrings":  formatStrings,
	}
	tmpl := template.Must(template.New("effects").Funcs(funcMap).Parse(templateString))
	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outFile, err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, map[string]interface{}{"Groups": groups}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

const missingEffectsFileName = "ui/core/constants/missing_effects_auto_gen.ts"

func GenerateMissingEffectsFile() error {
	if _, err := os.Stat(missingEffectsFileName); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to check file %s: %w", missingEffectsFileName, err)
	}

	tmpl := template.Must(template.New("missingEffects").Parse(TmplStrMissingEffects))
	f, err := os.Create(missingEffectsFileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", missingEffectsFileName, err)
	}
	defer f.Close()

	sort.Slice(missingEffectsMap["EnchantEffects"], func(i, j int) bool {
		return missingEffectsMap["EnchantEffects"][i].ID < missingEffectsMap["EnchantEffects"][j].ID
	})
	sort.Slice(missingEffectsMap["ItemEffects"], func(i, j int) bool {
		return missingEffectsMap["ItemEffects"][i].ID < missingEffectsMap["ItemEffects"][j].ID
	})

	if err := tmpl.Execute(f, missingEffectsMap); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func GenerateEnchantEffects(instance *dbc.DBC, db *WowDatabase) {
	groupMapProc := map[string]Group{}
	enchantSpellEffects := map[int]*dbc.SpellEffect{}

	for _, effect := range instance.SpellEffectsById {
		if effect.EffectType == dbc.E_ENCHANT_ITEM {
			enchantSpellEffects[effect.EffectMiscValues[0]] = &effect
		}
	}

	for _, enchant := range instance.Enchants {
		parsed := enchant.ToProto()
		if _, ok := db.Enchants[EnchantToDBKey(parsed)]; !ok {
			continue
		}

		if TryParseEnchantEffect(parsed, groupMapProc, instance, enchantSpellEffects) == EffectParseResultUnsupported {
			missingEffectsMap["EnchantEffects"] = append(missingEffectsMap["EnchantEffects"], Variant{ID: enchant.EffectId, Name: enchant.Name})
		}
	}

	var procGroups []*Group
	for _, grp := range groupMapProc {
		procGroups = append(procGroups, &grp)
	}
	GenerateEffectsFile(procGroups, "sim/common/mop/enchants_auto_gen.go", TmplStrEnchant)
}

func GenerateItemEffects(instance *dbc.DBC, db *WowDatabase, itemSources map[int][]*proto.DropSource) {
	groupMapOnUse := map[string]Group{}
	groupMapProc := map[string]Group{}

	// Example loop over your items
	for _, parsed := range db.Items {
		parsed.ItemEffect = dbc.MergeItemEffectsForAllStates(parsed)

		result := TryParseOnUseEffect(parsed, groupMapOnUse)
		if result == EffectParseResultSuccess {
			continue
		}

		if (result == EffectParseResultUnsupported) ||
			(TryParseProcEffect(parsed, instance, groupMapProc) == EffectParseResultUnsupported) {
			missingEffectsMap["ItemEffects"] = append(missingEffectsMap["ItemEffects"],
				Variant{
					ID:   int(parsed.Id),
					Name: parsed.Name + BuildItemDifficultyPostfix(itemSources, int(parsed.Id), instance),
				})
		}
	}

	// Sorting done in GenerateEffectsFile
	var onUseGroups []*Group
	for _, grp := range groupMapOnUse {
		onUseGroups = append(onUseGroups, &grp)
	}

	// Merge variants
	var procGroups []*Group
	needsStatPostfix := map[string]bool{}
	for _, grp := range groupMapProc {
		newEntries := []*Entry{}
		entryGroupings := map[string]*Entry{}

		// sort entries first to make tooltip generation consistent for variants
		sort.Slice(grp.Entries, func(i, j int) bool {
			return grp.Entries[i].Variants[0].ID < grp.Entries[j].Variants[0].ID
		})

		for _, entry := range grp.Entries {
			var idx int64 = 0
			added := false

			// Make sure to only group by name and proc mask, each proc mask will create it's own sub group
			for _, group := range entryGroupings {
				if group.Variants[0].Name == entry.Variants[0].Name {
					idx++
					if group.ProcInfo.ProcMask == entry.ProcInfo.ProcMask {
						group.AddVariant(entry.Variants[0])
						added = true
						break
					}
				}
			}

			if !added {
				groupName := entry.Variants[0].Name
				if idx > 0 {
					needsStatPostfix[groupName] = true
					groupName += "(" + strconv.FormatInt(idx, 10) + ")"
				}

				newEntries = append(newEntries, entry)
				entryGroupings[entry.Variants[0].Name] = entry
			}
		}

		grp.Entries = newEntries
		procGroups = append(procGroups, &grp)
	}

	updateNames := func(entries []*Entry) {
		for _, entry := range entries {
			for _, variant := range entry.Variants {
				if _, ok := needsStatPostfix[variant.Name]; ok {
					item := db.Items[int32(variant.ID)]
					variant.Name += " - " + GetEffectStatString(item)
				}

				variant.Name += BuildItemDifficultyPostfix(itemSources, variant.ID, instance)
			}
		}
	}

	// Update Item names
	for _, grp := range onUseGroups {
		updateNames(grp.Entries)
	}

	for _, grp := range procGroups {
		updateNames(grp.Entries)
	}

	GenerateEffectsFile(onUseGroups, "sim/common/mop/stat_bonus_cds_auto_gen.go", TmplStrOnUse)
	GenerateEffectsFile(procGroups, "sim/common/mop/stat_bonus_procs_auto_gen.go", TmplStrProc)
}

func GenerateItemEffectRandomPropPoints(instance *dbc.DBC, db *WowDatabase) {
	for id, allocMap := range instance.RandomPropertiesByIlvl {
		ilvl := int32(id)
		if ilvl < core.MinIlvl || ilvl > core.MaxIlvl {
			continue
		}
		db.ItemEffectRandPropPoints[ilvl] = &proto.ItemEffectRandPropPoints{
			Ilvl:           ilvl,
			RandPropPoints: allocMap[proto.ItemQuality_ItemQualityEpic][0],
		}
	}
}

func BuildItemDifficultyPostfix(itemSources map[int][]*proto.DropSource, itemId int, instance *dbc.DBC) string {
	difficultyPostfix := ""
	if sources, ok := itemSources[itemId]; ok {
		name := DifficultyToShortName(sources[0].Difficulty)
		if len(name) > 0 {
			difficultyPostfix += " " + name
		}
	}

	if item, ok := instance.Items[itemId]; ok {
		if len(item.NameDescription) > 0 && item.NameDescription != "Heroic" {
			difficultyPostfix += " (" + item.NameDescription + ")"
		}

		if item.Flags1.Has(dbc.HORDE_SPECIFIC) {
			difficultyPostfix += " (Horde)"
		}

		if item.Flags1.Has(dbc.ALLIANCE_SPECIFIC) {
			difficultyPostfix += " (Alliance)"
		}
	}

	return difficultyPostfix
}

func TryParseProcEffect(parsed *proto.UIItem, instance *dbc.DBC, groupMapProc map[string]Group) EffectParseResult {
	if parsed.ItemEffect.GetProc() != nil && parsed.ScalingOptions[0].Ilvl > MIN_EFFECT_ILVL {
		// Effect was already manually implemented
		if core.HasItemEffect(parsed.Id) {
			return EffectParseResultSuccess
		}

		tooltipString, id := dbc.GetItemEffectSpellTooltip(int(parsed.Id))
		tooltip, _ := tooltip.ParseTooltip(tooltipString, tooltip.DBCTooltipDataProvider{DBC: instance}, int64(id))

		grp, exists := groupMapProc["Procs"]
		if !exists {
			grp = Group{Name: "Procs"}
		}

		renderedTooltip := tooltip.String()
		entry := Entry{Tooltip: strings.Split(renderedTooltip, "\n"), Variants: []*Variant{{ID: int(parsed.Id), Name: parsed.Name}}}
		entry.ProcInfo, entry.Supported = BuildProcInfo(parsed, instance, renderedTooltip)
		entry.Harmful = true
		grp.Entries = append(grp.Entries, &entry)
		groupMapProc["Procs"] = grp

		if !entry.Supported {
			return EffectParseResultUnsupported
		}

		return EffectParseResultSuccess
	}

	// check if the item has any kind of proc as we only support stat proc parsing right now
	if effects, ok := instance.ItemEffectsByParentID[int(parsed.Id)]; ok && parsed.ScalingOptions[0].Ilvl > MIN_EFFECT_ILVL {
		for _, effect := range effects {
			if SpellHasTriggerEffect(effect.SpellID, instance) {
				return EffectParseResultUnsupported
			}
		}
	}

	return EffectParseResultInvalid
}

func TryParseOnUseEffect(parsed *proto.UIItem, groupMap map[string]Group) EffectParseResult {
	// Effect was already manually implemented
	if core.HasItemEffect(parsed.Id) {
		return EffectParseResultSuccess
	}

	if parsed.ItemEffect.GetOnUse() != nil && parsed.ScalingOptions[0].Ilvl > MIN_EFFECT_ILVL { // MoP constraints

		if parsed.ItemEffect.GetOnUse().CooldownMs < 0 && parsed.ItemEffect.GetOnUse().CategoryCooldownMs < 0 {
			return EffectParseResultUnsupported
		}

		groupName := GetEffectStatString(parsed)
		grp, exists := groupMap[groupName]
		if !exists {
			grp = Group{Name: groupName}
		}
		grp.Entries = append(grp.Entries, &Entry{Variants: []*Variant{{ID: int(parsed.Id), Name: parsed.Name}}})
		groupMap[groupName] = grp
		return EffectParseResultSuccess
	}

	return EffectParseResultInvalid
}

func TryParseEnchantEffect(enchant *proto.UIEnchant, groupMapProc map[string]Group, instance *dbc.DBC, enchantSpellEffects map[int]*dbc.SpellEffect) EffectParseResult {
	if (enchant.EnchantEffect.GetProc() != nil || EnchantHasDummyEffect(enchant, instance)) && enchant.EffectId > 4267 {

		// Effect was already manually implemented
		if core.HasEnchantEffect(enchant.EffectId) {
			return EffectParseResultSuccess
		}

		if enchantingSpell, ok := enchantSpellEffects[int(enchant.EffectId)]; ok {
			tooltipString := instance.Spells[enchantingSpell.SpellID].Description
			tooltip, _ := tooltip.ParseTooltip(tooltipString, tooltip.DBCTooltipDataProvider{DBC: instance}, int64(enchantingSpell.SpellID))

			grp, exists := groupMapProc["Enchants"]
			if !exists {
				grp = Group{Name: "Enchants"}
			}

			renderedTooltip := tooltip.String()
			entry := Entry{Tooltip: strings.Split(renderedTooltip, "\n"), Variants: []*Variant{{ID: int(enchant.EffectId), Name: enchant.Name}}}
			entry.ProcInfo, entry.Supported = BuildEnchantProcInfo(enchant, instance, renderedTooltip)
			entry.Harmful = true
			grp.Entries = append(grp.Entries, &entry)
			groupMapProc["Enchants"] = grp
			if !entry.Supported {
				return EffectParseResultUnsupported
			}

			return EffectParseResultSuccess
		}
	}

	return EffectParseResultInvalid
}

var critMatcher = regexp.MustCompile(`critical ([^\s]+|damage,?)( chance)? [^fbc]`)
var pureHealMatcher = regexp.MustCompile(`healing spells`)
var hasHealMatcher = regexp.MustCompile(`heal(ing)?[^,]`)
var hasGenericMatcher = regexp.MustCompile(`a spell`)

func BuildProcInfo(parsed *proto.UIItem, instance *dbc.DBC, tooltip string) (ProcInfo, bool) {
	itemEffectInfo, ok := instance.ItemEffectsByParentID[int(parsed.Id)]
	if !ok {
		fmt.Printf("WARN: Can not generate proc info for Item: %d, not found.\n", parsed.Id)
	}

	// if we have multiple spells find the first that has a proc aura assigned
	for _, effectInfo := range itemEffectInfo {
		procId := effectInfo.SpellID
		procSpell, ok := instance.Spells[int(procId)]
		if !ok {
			panic(fmt.Sprintf("Could not find proc aura %d spell for item effect %d.\n", procId, parsed.Id))
		}

		if len(procSpell.ProcTypeMask) == 0 || procSpell.ProcTypeMask[0] == 0 {
			continue
		}

		itemType := proto.ItemType_ItemTypeUnknown
		if itemEffectInfo[0].TriggerType == 2 {
			itemType = proto.ItemType_ItemTypeWeapon
		}

		procInfo, supported := BuildSpellProcInfo(&procSpell, tooltip, itemType)

		// we do not support generation of more than one proc effect right now
		if len(itemEffectInfo) > 1 {
			return procInfo, false
		}

		if SpellHasDummyEffect(int(procId), instance) {
			return procInfo, false
		}

		if SpellUsesStacks(int(procId), instance) {
			return procInfo, false
		}

		return procInfo, supported
	}

	return ProcInfo{}, false
}

func BuildEnchantProcInfo(enchant *proto.UIEnchant, instance *dbc.DBC, tooltip string) (ProcInfo, bool) {
	procSpellID := enchant.SpellId
	if procSpellID == 0 {
		fmt.Printf("WARN: Enchant %d with no spell id", enchant.EffectId)
		return ProcInfo{}, false
	}

	procSpell, ok := instance.Spells[int(procSpellID)]
	if !ok {
		panic(fmt.Sprintf("Could not find proc aura %d spell for item effect %d.\n", procSpellID, enchant.EffectId))
	}

	procInfo, supported := BuildSpellProcInfo(&procSpell, tooltip, enchant.Type)
	if SpellHasDummyEffect(int(procSpellID), instance) {
		return procInfo, false
	}

	if SpellUsesStacks(int(procSpellID), instance) {
		return procInfo, false
	}

	return procInfo, supported
}

func BuildSpellProcInfo(procSpell *dbc.Spell, tooltip string, itemType proto.ItemType) (ProcInfo, bool) {
	var info = ProcInfo{}

	// On hit proc
	if itemType == proto.ItemType_ItemTypeWeapon {
		info.Callback |= core.CallbackOnSpellHitDealt
		info.ProcMask |= core.ProcMaskMelee
	}

	if itemType == proto.ItemType_ItemTypeRanged {
		info.Callback |= core.CallbackOnSpellHitDealt
		info.ProcMask |= core.ProcMaskRanged
	}

	if len(procSpell.ProcTypeMask) > 0 {
		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_MELEE_SWING > 0 {
			info.ProcMask |= core.ProcMaskMeleeWhiteHit
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_MELEE_ABILITY > 0 {
			info.ProcMask |= core.ProcMaskMeleeSpecial
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_RANGED_ATTACK > 0 {
			info.ProcMask |= core.ProcMaskRangedAuto
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_RANGED_ABILITY > 0 {
			info.ProcMask |= core.ProcMaskRangedSpecial
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HARMFUL_PERIODIC > 0 {
			info.ProcMask |= core.ProcMaskSpellDamage
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HARMFUL_SPELL > 0 {
			info.ProcMask |= core.ProcMaskSpellDamage
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_ANY_DIRECT_TAKEN > 0 {
			info.Callback |= core.CallbackOnSpellHitTaken
			info.Outcome = core.OutcomeLanded

			if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_TAKE_MELEE_SWING > 0 {
				info.ProcMask |= core.ProcMaskMeleeWhiteHit
			}

			if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_TAKE_MELEE_ABILITY > 0 {
				info.ProcMask |= core.ProcMaskMeleeSpecial
			}

			if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_TAKE_HARMFUL_SPELL > 0 {
				info.ProcMask |= core.ProcMaskSpellDamage
			}

			// For now we do not support self damage procs as they usually have custom extra proc conditions
			// like On dodge or on On parry or x amount of damage taken
			return info, false
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_ANY_DIRECT_DEALT > 0 {
			info.Callback |= core.CallbackOnSpellHitDealt
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HARMFUL_PERIODIC > 0 {
			info.Callback |= core.CallbackOnPeriodicDamageDealt
		}

		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HELPFUL_SPELL > 0 &&
			(hasHealMatcher.MatchString(tooltip) || hasGenericMatcher.MatchString(tooltip)) {
			info.Callback |= core.CallbackOnHealDealt
			info.ProcMask |= core.ProcMaskSpellHealing

			// handle HoTs onyl with direct heals for now, there are some odd cases with HoT / DoT overlaps
			if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HELPFUL_PERIODIC > 0 {
				info.Callback |= core.CallbackOnPeriodicHealDealt
			}

			// Check if we have periodic damage flag but only heal paired with it
			// This usually indicates a pure heal proc mask
			if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_ANY_DIRECT_DEALT == 0 {
				info.Callback &= ^core.CallbackOnPeriodicDamageDealt
				info.Callback &= ^core.CallbackOnSpellHitDealt
				info.ProcMask &= ^core.ProcMaskSpellDamage
			}
		}
	}

	if info.ProcMask.Matches(core.ProcMaskMelee) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskMeleeProc
	}

	if info.ProcMask.Matches(core.ProcMaskRanged) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskRangedProc
	}

	if info.ProcMask.Matches(core.ProcMaskSpellDamage) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskSpellDamageProc
	}

	if critMatcher.MatchString(tooltip) {
		info.Outcome = core.OutcomeCrit
	} else {
		info.Outcome = core.OutcomeLanded
	}

	// check for pure healing spell
	if pureHealMatcher.MatchString(tooltip) {
		info.Callback &= ^core.CallbackOnSpellHitDealt
		info.Callback &= ^core.CallbackOnPeriodicDamageDealt
	}

	unsupported := info.Callback == core.CallbackEmpty &&
		info.Outcome == core.OutcomeEmpty &&
		info.ProcMask == core.ProcMaskEmpty

	return info, !unsupported
}

func asCoreCallback(callback core.AuraCallback) string {
	callbacks := []string{}
	for i := range 32 {
		callbackFlag := core.AuraCallback(1 << i)
		if callbackFlag >= core.CallbackLast {
			break
		}

		if callback.Matches(callbackFlag) {
			callbacks = append(callbacks, "core."+callbackFlag.String())
		}
	}

	if len(callbacks) == 0 {
		return "core.CallbackEmpty"
	}

	return strings.Join(callbacks, " | ")
}

func asCoreProcMask(procMask core.ProcMask) string {
	procs := []string{}
	for i := range 32 {
		procFlag := core.ProcMask(1 << i)
		if procFlag >= core.ProcMaskLast {
			break
		}

		if procMask.Matches(procFlag) {
			procs = append(procs, "core."+procFlag.String())
		}
	}

	if len(procs) == 0 {
		return "core.ProcMaskEmpty"
	}
	return strings.Join(procs, " | ")
}

func asCoreOutcome(outcome core.HitOutcome) string {
	if outcome == core.OutcomeCrit {
		return "core.OutcomeCrit"
	}

	if outcome.Matches(core.OutcomeLanded) {
		return "core.OutcomeLanded"
	}

	return "core.OutcomeEmpty"
}

func (entry *Entry) AddVariant(variant *Variant) {
	entry.Variants = append(entry.Variants, variant)
	sort.Slice(entry.Variants, func(i, j int) bool {
		return entry.Variants[i].ID < entry.Variants[j].ID
	})
}
