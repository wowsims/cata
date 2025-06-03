package database

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools/database/dbc"
	"github.com/wowsims/mop/tools/tooltip"
)

type ProcInfo struct {
	Outcome  core.HitOutcome
	Callback core.AuraCallback
	ProcMask core.ProcMask
	IsEmpty  bool
}

// Entry represents a single effect with its ID and display name.
type Entry struct {
	ID       int
	Name     string
	Tooltip  string
	ProcInfo ProcInfo
}

// Group holds a category of effects.
type Group struct {
	Name    string
	Entries []Entry
}

// Define your groups and effects here.
// The map key is the group name, and the inner map is ID -> display name.

const TmplStrOnUse = `package mop

import (
	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllOnUseCds() {
{{- range .Groups }}

	// {{ .Name }}
{{- range .Entries }}
	shared.NewSimpleStatActive({{ .ID }}) // {{ .Name }}
{{- end }}

{{- end }}
}`
const TmplStrProc = `package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllOnUseProcs() {
{{- range .Groups }}

	// {{ .Name }}
{{- range .Entries }}
	{{if .ProcInfo.IsEmpty}}
	// TODO: Overwrite me
	{{- end}}
	// {{.Tooltip}}
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		Name:     "{{ .Name }}",
		ItemID:   {{ .ID }},
		Callback: {{ .ProcInfo.Callback | asCoreCallback }},
		ProcMask: {{ .ProcInfo.ProcMask | asCoreProcMask }},
		Outcome:  {{ .ProcInfo.Outcome | asCoreOutcome }},
	})
{{- end }}

{{- end }}
}`

func GenerateEffectsFile(groups []Group, outFile string, templateString string) error {
	// Check if file already exists
	if _, err := os.Stat(outFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to check file %s: %w", outFile, err)
	}

	// Ensure groups and entries are sorted
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	for _, grp := range groups {
		sort.Slice(grp.Entries, func(i, j int) bool {
			if grp.Entries[i].ProcInfo.IsEmpty != grp.Entries[j].ProcInfo.IsEmpty {
				return grp.Entries[i].ProcInfo.IsEmpty
			}

			return grp.Entries[i].ID < grp.Entries[j].ID
		})
	}

	funcMap := map[string]any{
		"asCoreCallback": asCoreCallback,
		"asCoreProcMask": asCoreProcMask,
		"asCoreOutcome":  asCoreOutcome,
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

func GenerateEffects(instance *dbc.DBC, iconsMap map[int]string, db *WowDatabase) {
	groupMap := map[string]Group{}
	groupMapProc := map[string]Group{}

	// Example loop over your items
	for _, item := range instance.Items {
		parsed := item.ToUIItem()

		if parsed.Icon == "" {
			parsed.Icon = strings.ToLower(GetIconName(iconsMap, item.FDID))
		}

		parsed.ItemEffect = dbc.MergeItemEffectsForAllStates(parsed)
		db.MergeItem(parsed)

		if TryParseOnUseEffect(parsed, item, groupMap) {
			continue
		}

		TryParseProcEffect(parsed, item, instance, groupMapProc)
	}

	// Sorting done in GenerateEffectsFile
	var groups []Group
	for _, grp := range groupMap {
		groups = append(groups, grp)
	}

	var procGroups []Group
	for _, grp := range groupMapProc {
		procGroups = append(procGroups, grp)
	}

	GenerateEffectsFile(groups, "sim/common/mop/stat_bonus_cds.go", TmplStrOnUse)
	GenerateEffectsFile(procGroups, "sim/common/mop/stat_bonus_procs.go", TmplStrProc)
}

func TryParseProcEffect(parsed *proto.UIItem, item dbc.Item, instance *dbc.DBC, groupMapProc map[string]Group) {
	if parsed.ItemEffect.GetProc() != nil && item.ItemLevel > 416 {
		tooltipString, id := dbc.GetItemEffectSpellTooltip(item.Id)
		tooltip, _ := tooltip.ParseTooltip(tooltipString, tooltip.DBCTooltipDataProvider{DBC: instance}, int64(id))

		grp, exists := groupMapProc["Procs"]
		if !exists {
			grp = Group{Name: "Procs"}
		}

		entry := Entry{Tooltip: tooltip.String(), ID: int(parsed.Id), Name: parsed.Name}
		entry.ProcInfo = BuildProcInfo(parsed, instance, entry.Tooltip)
		grp.Entries = append(grp.Entries, entry)
		groupMapProc["Procs"] = grp
	}
}

func TryParseOnUseEffect(parsed *proto.UIItem, item dbc.Item, groupMap map[string]Group) bool {
	if parsed.ItemEffect.GetOnUse() != nil && item.ItemLevel > 416 { // MoP constraints
		stats := parsed.ItemEffect.ScalingOptions[int32(proto.ItemLevelState_Base)].Stats
		var firstStat proto.Stat = proto.Stat_StatStrength
		found := false
		for k := range stats {
			stat := proto.Stat(k)
			if !found || stat < firstStat {
				firstStat = stat
				found = true
			}
		}

		groupName := firstStat.String()
		grp, exists := groupMap[groupName]
		if !exists {
			grp = Group{Name: groupName}
		}
		grp.Entries = append(grp.Entries, Entry{ID: int(parsed.Id), Name: parsed.Name})
		groupMap[groupName] = grp
		return true
	}

	return false
}

var critMatcher = regexp.MustCompile(`critical [^\s]+ [^f]`)
var pureHealMatcher = regexp.MustCompile(`healing spells`)

func BuildProcInfo(parsed *proto.UIItem, instance *dbc.DBC, tooltip string) ProcInfo {
	var info = ProcInfo{
		IsEmpty: true,
	}

	itemEffectInfo, ok := instance.ItemEffectsByParentID[int(parsed.Id)]
	if !ok || len(itemEffectInfo) > 1 {
		fmt.Printf("WARN: Can not generate proc info for Item: %d, not supported.\n", parsed.Id)
		return info
	}

	procId := itemEffectInfo[0].SpellID
	procSpell, ok := instance.Spells[int(procId)]
	if !ok {
		panic(fmt.Sprintf("Could not find proc aura %d spell for item effect %d.\n", procId, parsed.Id))
	}

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

	if info.ProcMask.Matches(core.ProcMaskMelee) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskMeleeProc
	}

	if info.ProcMask.Matches(core.ProcMaskRanged) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskRangedProc
	}

	if info.ProcMask.Matches(core.ProcMaskSpellDamage) && procSpell.Attributes[3]&dbc.ATTR_EX_3_CAN_PROC_FROM_PROCS > 0 {
		info.ProcMask |= core.ProcMaskSpellDamageProc
	}

	findOutcome := func() core.HitOutcome {
		if critMatcher.MatchString(tooltip) {
			return core.OutcomeCrit
		}

		return core.OutcomeLanded
	}

	info.Outcome = findOutcome()

	if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_ANY_DIRECT_TAKEN > 0 {
		info.Callback |= core.CallbackOnSpellHitTaken

		// For now we do not support self damage procs as they usually have custom extra proc conditions
		// like On dodge or on On parry or x amount of damage taken
		return ProcInfo{IsEmpty: true}
	}

	if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_ANY_DIRECT_DEALT > 0 {
		info.Callback |= core.CallbackOnSpellHitDealt
	}

	if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HARMFUL_PERIODIC > 0 {
		info.Callback |= core.CallbackOnPeriodicDamageDealt
	}

	if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HELPFUL_SPELL > 0 {
		info.Callback |= core.CallbackOnHealDealt

		// handle HoTs onyl with direct heals for now, there are some odd cases with HoT / DoT overlaps
		if procSpell.ProcTypeMask[0]&dbc.PROC_FLAG_DEAL_HELPFUL_PERIODIC > 0 {
			info.Callback |= core.CallbackOnPeriodicHealDealt
		}
	}

	// On hit proc
	if itemEffectInfo[0].TriggerType == 2 {
		info.Callback |= core.CallbackOnSpellHitDealt
	}

	// check for pure healing spell
	if pureHealMatcher.MatchString(tooltip) {
		info.Callback &= ^core.CallbackOnSpellHitDealt
		info.Callback &= ^core.CallbackOnPeriodicDamageDealt
	}

	info.IsEmpty = info.Callback == core.CallbackEmpty &&
		info.Outcome == core.OutcomeEmpty &&
		info.ProcMask == core.ProcMaskEmpty

	return info
}

func asCoreCallback(callback core.AuraCallback) string {
	callbacks := []string{}
	if callback.Matches(core.CallbackOnApplyEffects) {
		callbacks = append(callbacks, "core.CallabckOnApplyEffects")
	}

	if callback.Matches(core.CallbackOnCastComplete) {
		callbacks = append(callbacks, "core.CallbackOnCastComplete")
	}

	if callback.Matches(core.CallbackOnHealDealt) {
		callbacks = append(callbacks, "core.CallbackOnHealDealt")
	}

	if callback.Matches(core.CallbackOnPeriodicDamageDealt) {
		callbacks = append(callbacks, "core.CallbackOnPeriodicDamageDealt")
	}

	if callback.Matches(core.CallbackOnPeriodicHealDealt) {
		callbacks = append(callbacks, "core.CallbackOnPeriodicHealDealt")
	}

	if callback.Matches(core.CallbackOnSpellHitDealt) {
		callbacks = append(callbacks, "core.CallbackOnSpellHitDealt")
	}

	if callback.Matches(core.CallbackOnSpellHitTaken) {
		callbacks = append(callbacks, "core.CallbackOnSpellHitTaken")
	}

	if len(callbacks) == 0 {
		return "core.CallbackEmpty"
	}

	return strings.Join(callbacks, " | ")
}

func asCoreProcMask(procMask core.ProcMask) string {
	procs := []string{}

	if procMask.Matches(core.ProcMaskMeleeMHAuto) {
		procs = append(procs, "core.ProcMaskMeleeMHAuto")
	}

	if procMask.Matches(core.ProcMaskMeleeOHAuto) {
		procs = append(procs, "core.ProcMaskMeleeOHAuto")
	}

	if procMask.Matches(core.ProcMaskMeleeMHSpecial) {
		procs = append(procs, "core.ProcMaskMeleeMHSpecial")
	}

	if procMask.Matches(core.ProcMaskMeleeOHSpecial) {
		procs = append(procs, "core.ProcMaskMeleeOHSpecial")
	}

	if procMask.Matches(core.ProcMaskRangedAuto) {
		procs = append(procs, "core.ProcMaskRangedAuto")
	}

	if procMask.Matches(core.ProcMaskRangedSpecial) {
		procs = append(procs, "core.ProcMaskRangedSpecial")
	}

	if procMask.Matches(core.ProcMaskSpellDamage) {
		procs = append(procs, "core.ProcMaskSpellDamage")
	}

	if procMask.Matches(core.ProcMaskSpellHealing) {
		procs = append(procs, "core.ProcMaskSpellHealing")
	}

	if procMask.Matches(core.ProcMaskSpellProc) {
		procs = append(procs, "core.ProcMaskSpellProc")
	}

	if procMask.Matches(core.ProcMaskMeleeProc) {
		procs = append(procs, "core.ProcMaskMeleeProc")
	}

	if procMask.Matches(core.ProcMaskRangedProc) {
		procs = append(procs, "core.ProcMaskRangedProc")
	}

	if procMask.Matches(core.ProcMaskSpellDamageProc) {
		procs = append(procs, "core.ProcMaskSpellDamageProc")
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
