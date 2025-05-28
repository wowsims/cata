package database

import (
	"fmt"
	"os"
	"sort"
	"text/template"
)

// Entry represents a single effect with its ID and display name.
type Entry struct {
	ID      int
	Name    string
	Tooltip string
}

// Group holds a category of effects.
type Group struct {
	Name    string
	Entries []Entry
}

// Define your groups and effects here.
// The map key is the group name, and the inner map is ID -> display name.

const TmplStrOnUse = `package cata

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
const TmplStrProc = `package cata

import (
	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllOnUseProcs() {
{{- range .Groups }}
	// {{ .Name }}
{{- range .Entries }}

	//{{.Tooltip}}
	//shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
	//	Name:     "{{ .Name }}",
	//	ItemID:   {{ .ID }},
	//	Callback: core.CallbackOnSpellHitDealt, // Fill this
	//	ProcMask: core.ProcMaskMeleeOrRanged, // Fill this
	//	Outcome:  core.OutcomeCrit, // Fill this
	//})
{{- end }}

{{- end }}
}`

func GenerateEffectsFile(groups []Group, outFile string, templateString string) error {
	// Check if file already exists
	if _, err := os.Stat(outFile); err == nil {
		return fmt.Errorf("file %s already exists", outFile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("unable to check file %s: %w", outFile, err)
	}

	// Ensure groups and entries are sorted
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	for _, grp := range groups {
		sort.Slice(grp.Entries, func(i, j int) bool { return grp.Entries[i].ID < grp.Entries[j].ID })
	}

	tmpl := template.Must(template.New("effects").Parse(templateString))
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
