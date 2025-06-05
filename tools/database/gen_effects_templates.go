package database

const TmplStrOnUse = `package mop

import (
	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllOnUseCds() {
{{- range .Groups }}

	// {{ .Name }}
{{- range .Entries }}
  	{{- with index .Variants 0}}
	shared.NewSimpleStatActive({{ .ID }}) // {{ .Name }}
	{{- end}}
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
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	{{- end}}
	// {{.Tooltip}}
	{{- if len .Variants | eq 1}}
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		{{with index .Variants 0 -}}
		Name:     "{{ .Name }}",
		EnchantID: {{ .ID }},
		{{- end}}
		Callback: {{ .ProcInfo.Callback | asCoreCallback }},
		ProcMask: {{ .ProcInfo.ProcMask | asCoreProcMask }},
		Outcome:  {{ .ProcInfo.Outcome | asCoreOutcome }},
	})
	{{- else }}
	shared.NewProcStatBonusEffectWithVariants(shared.ProcStatBonusEffect{
		Callback: {{ .ProcInfo.Callback | asCoreCallback }},
		ProcMask: {{ .ProcInfo.ProcMask | asCoreProcMask }},
		Outcome:  {{ .ProcInfo.Outcome | asCoreOutcome }},
	}, []shared.ItemVariant{
		{{- range .Variants }}
		{ItemID: {{.ID}}, ItemName: "{{.Name}}"},
		{{- end}}
	})
	{{- end}}
{{- end }}

{{- end }}
}`

const TmplStrEnchant = `package mop

import (
	"github.com/wowsims/mop/sim/core"
 	"github.com/wowsims/mop/sim/common/shared"
)

func RegisterAllEnchants() {
{{- range .Groups }}

	// {{ .Name }}
{{- range .Entries }}
	{{if .ProcInfo.IsEmpty}}
	// TODO: Manual implementation required
	//       This can be ignored if the effect has already been implemented.
	//       With next db run the item will be removed if implemented.
	//
	{{- end}}
	{{- range (.Tooltip | formatStrings 100) }}
	// {{.}}
	{{- end}}
	shared.NewProcStatBonusEffect(shared.ProcStatBonusEffect{
		{{with index .Variants 0 -}}
		Name:     "{{ .Name }}",
		EnchantID: {{ .ID }},
		{{- end}}
		Callback:  {{ .ProcInfo.Callback | asCoreCallback }},
		ProcMask:  {{ .ProcInfo.ProcMask | asCoreProcMask }},
		Outcome:   {{ .ProcInfo.Outcome | asCoreOutcome }},
	})
{{- end }}

{{- end }}
}`

const TmplStrMissingEffects = `
// This file is auto generated
// Changes will be overwritten on next database generation

export const MISSING_ITEM_EFFECTS = [
{{- range .ItemEffects }}
    {{.}},
{{- end }}
]

export const MISSING_ENCHANT_EFFECTS = [
{{- range .EnchantEffects }}
    {{.}},
{{- end }}
]
`
