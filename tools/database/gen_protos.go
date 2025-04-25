package database

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/tools/database/dbc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func convertTalentClassID(raw int) int {
	return 1 << (raw - 1)
}

var dbcClasses = []dbc.DbcClass{
	{ProtoClass: proto.Class_ClassWarrior, ID: 1},
	{ProtoClass: proto.Class_ClassPaladin, ID: 2},
	{ProtoClass: proto.Class_ClassHunter, ID: 3},
	{ProtoClass: proto.Class_ClassRogue, ID: 4},
	{ProtoClass: proto.Class_ClassPriest, ID: 5},
	{ProtoClass: proto.Class_ClassDeathKnight, ID: 6},
	{ProtoClass: proto.Class_ClassShaman, ID: 7},
	{ProtoClass: proto.Class_ClassMage, ID: 8},
	{ProtoClass: proto.Class_ClassWarlock, ID: 9},
	{ProtoClass: proto.Class_ClassDruid, ID: 11},
}

func classNameFromDBC(dbc dbc.DbcClass) string {
	switch dbc.ID {
	case 1:
		return "Warrior"
	case 2:
		return "Paladin"
	case 3:
		return "Hunter"
	case 4:
		return "Rogue"
	case 5:
		return "Priest"
	case 6:
		return "Death_Knight"
	case 7:
		return "Shaman"
	case 8:
		return "Mage"
	case 9:
		return "Warlock"
	case 11:
		return "Druid"
	default:
		return "Unknown"
	}
}

type TalentConfig struct {
	FieldName        string          `json:"fieldName"`
	FancyName        string          `json:"fancyName"`
	Location         TalentLocation  `json:"location"`
	SpellIds         []int           `json:"spellIds"`
	MaxPoints        int             `json:"maxPoints"`
	PrereqLocation   *TalentLocation `json:"prereqLocation,omitempty"`
	TabName          string          `json:"tabName"`
	ProtoFieldNumber int
}

type TalentLocation struct {
	RowIdx int `json:"rowIdx"`
	ColIdx int `json:"colIdx"`
}

type ClassData struct {
	ClassName          string
	LowerCaseClassName string
	FileName           string
	Talents            []TalentConfig
	TalentTabs         []TalentTabConfig
	GlyphsPrime        []Glyph
	GlyphsMajor        []Glyph
	GlyphsMinor        []Glyph
}

type Glyph struct {
	EnumName    string
	Name        string
	Description string
	IconUrl     string
	ID          int
}

const staticHeader = `syntax = "proto3";
package proto;
option go_package = "./proto";
import "common.proto";`

const protoTemplateStr = `
{{- $class := .ClassName -}}
// {{.ClassName}}Talents message.
message {{$class}}Talents {
{{- range $tab := .TalentTabs }}
    // {{$tab.Name}}
{{- range $talent := $tab.Talents }}
    {{- if eq $talent.MaxPoints 1 }}
    bool {{ final $talent.FancyName $class }} = {{ $talent.ProtoFieldNumber }};
    {{- else }}
    int32 {{ final $talent.FancyName $class }} = {{ $talent.ProtoFieldNumber }};
    {{- end }}
{{- end }}
{{- end }}
}

enum {{.ClassName}}MajorGlyph {
    {{.ClassName}}MajorGlyphNone = 0;
    {{- range .GlyphsMajor }}
    {{ protoOverride .EnumName $class }} = {{ .ID }};
    {{- end }}
}

enum {{.ClassName}}MinorGlyph {
    {{.ClassName}}MinorGlyphNone = 0;
    {{- range .GlyphsMinor }}
    {{ protoOverride .EnumName $class }} = {{ .ID }};
    {{- end }}
}
`

const tsTemplateStr = `import { {{.ClassName}}MajorGlyph, {{.ClassName}}MinorGlyph, {{.ClassName}}Talents } from '../proto/{{.FileName}}.js';
import { GlyphsConfig } from './glyphs_picker.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import {{.ClassName}}TalentJson from './trees/{{.FileName}}.json';
{{- $class := .ClassName -}}
export const {{.LowerCaseClassName}}TalentsConfig: TalentsConfig<{{.ClassName}}Talents> = newTalentsConfig({{.ClassName}}TalentJson);

export const {{.LowerCaseClassName}}GlyphsConfig: GlyphsConfig = {
	majorGlyphs: {
		{{- range .GlyphsMajor }}
		[{{$.ClassName}}MajorGlyph.{{protoOverride .EnumName $class}}]: {
			name: "{{.Name}}",
			description: "{{.Description}}",
			iconUrl: "{{.IconUrl}}",
		},
		{{- end }}
	},
	minorGlyphs: {
		{{- range .GlyphsMinor }}
		[{{$.ClassName}}MinorGlyph.{{protoOverride .EnumName $class}}]: {
			name: "{{.Name}}",
			description: "{{.Description}}",
			iconUrl: "{{.IconUrl}}",
		},
		{{- end }}
	},
};
`

const talentJsonTemplate = `[
{{- $n := len . }}
{{- range $i, $tab := . }}
  {
    "name": "{{ $tab.Name }}",
    "backgroundUrl": "{{ $tab.BackgroundUrl }}",
    "talents": [
    {{- $m := len $tab.Talents }}
    {{- range $j, $talent := $tab.Talents }}
      {
        "fieldName": "{{ toCamelCase $talent.FancyName }}",
        "fancyName": "{{ $talent.FancyName }}",
        "location": {
          "rowIdx": {{ $talent.Location.RowIdx }},
          "colIdx": {{ $talent.Location.ColIdx }}
        },
        "spellIds": [{{- range $k, $id := $talent.SpellIds }}{{if $k}}, {{end}}{{ $id }}{{- end }}],
        "maxPoints": {{ $talent.MaxPoints }}{{ if $talent.PrereqLocation }},
        "prereqLocation": {
          "rowIdx": {{ $talent.PrereqLocation.RowIdx }},
          "colIdx": {{ $talent.PrereqLocation.ColIdx }}
        }{{ end }}
      }{{ if ne (add $j 1) $m }},{{ end }}
    {{- end }}
    ]
  }{{ if ne (add $i 1) $n }},{{ end }}
{{- end }}
]
`

func generateProtoFile(data ClassData) error {
	// Create the proto directory if it doesn't exist
	dirPath := "./proto"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dirPath, err)
	}
	fileName := fmt.Sprintf("%s/%s.proto", dirPath, strings.ToLower(data.ClassName))
	if err := createOrUpdateProtoFile(fileName, data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}
func generateTemplateContent(data ClassData) (string, error) {
	funcMap := template.FuncMap{
		"add":           func(a, b int) int { return a + b },
		"toCamelCase":   toCamelCase,
		"toSnakeCase":   toSnakeCase,
		"protoOverride": protoOverride,
		"final":         finalFieldName,
	}

	data.ClassName = strings.ReplaceAll(data.ClassName, "_", "")
	slices.SortFunc(data.GlyphsMajor, func(a, b Glyph) int {
		return a.ID - b.ID
	})
	slices.SortFunc(data.GlyphsMinor, func(a, b Glyph) int {
		return a.ID - b.ID
	})
	slices.SortFunc(data.GlyphsPrime, func(a, b Glyph) int {
		return a.ID - b.ID
	})

	tmpl, err := template.New("protoTemplate").Funcs(funcMap).Parse(protoTemplateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return b.String(), nil
}
func createOrUpdateProtoFile(filePath string, data ClassData) error {
	newGeneratedContent, err := generateTemplateContent(data)
	if err != nil {
		return fmt.Errorf("error generating template content: %w", err)
	}

	generatedBlock := fmt.Sprintf("// BEGIN GENERATED\n%s\n// END GENERATED", newGeneratedContent)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fullContent := staticHeader + "\n\n" + generatedBlock + "\n"
		err := os.WriteFile(filePath, []byte(fullContent), 0644)
		if err != nil {
			return fmt.Errorf("error creating file %s: %w", filePath, err)
		}
	} else {
		existingBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", filePath, err)
		}
		existingContent := string(existingBytes)

		updatedContent, err := updateGeneratedProtoSection(existingContent, newGeneratedContent)
		if err != nil {
			return fmt.Errorf("error updating generated section: %w", err)
		}
		err = os.WriteFile(filePath, []byte(updatedContent), 0644)
		if err != nil {
			return fmt.Errorf("error writing updated file %s: %w", filePath, err)
		}
	}

	return nil
}
func updateGeneratedProtoSection(fileContent, newContent string) (string, error) {
	const beginMarker = "// BEGIN GENERATED"
	const endMarker = "// END GENERATED"
	beginIdx := strings.Index(fileContent, beginMarker)
	if beginIdx == -1 {
		return "", fmt.Errorf("begin marker %q not found in the file", beginMarker)
	}
	endIdx := strings.LastIndex(fileContent, endMarker)
	if endIdx == -1 {
		return "", fmt.Errorf("end marker %q not found in the file", endMarker)
	}
	endIdx += len(endMarker)
	newBlock := fmt.Sprintf("%s\n%s\n%s", beginMarker, newContent, endMarker)
	updatedContent := fileContent[:beginIdx] + newBlock + fileContent[endIdx:]
	return updatedContent, nil
}

func protoOverride(name string, className string) string {
	if name == "GlyphOfDeathCoil" && className == "Warlock" {
		return "GlyphOfDeathCoilWarlock"
	}
	if name == "GlyphOfStampede" && className == "Hunter" {
		return "GlyphOfStampedeHunter"
	}
	if name == "Tnt" || name == "tnt" {
		return "TNT"
	}
	return name
}
func generateTsFile(data ClassData) error {
	dirPath := "./ui/core/talents"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dirPath, err)
	}

	fileName := fmt.Sprintf("%s/%s.ts", dirPath, strings.ToLower(data.ClassName))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", fileName, err)
	}
	defer file.Close()

	funcMap := template.FuncMap{
		"add":           func(a, b int) int { return a + b },
		"toCamelCase":   toCamelCase,
		"toSnakeCase":   toSnakeCase,
		"protoOverride": protoOverride,
	}
	tmpl, err := template.New("tsTemplate").Funcs(funcMap).Parse(tsTemplateStr)
	if err != nil {
		return err
	}
	data.ClassName = strings.ReplaceAll(data.ClassName, "_", "")
	data.FileName = data.LowerCaseClassName
	data.LowerCaseClassName = strings.ReplaceAll(data.LowerCaseClassName, "_", "")
	if data.LowerCaseClassName == "deathknight" {
		data.LowerCaseClassName = "deathKnight"
	}
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template for %s: %w", data.ClassName, err)
	}

	return nil
}
func finalFieldName(fancyName, className string) string {
	snake := toSnakeCase(fancyName)
	return protoOverride(snake, className)
}

type TalentTabConfig struct {
	Name          string         `json:"name"`
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func generateTalentJson(tabs []TalentTabConfig, className string) error {
	// Create the directory if it doesn't exist
	dirPath := "ui/core/talents/trees"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dirPath, err)
	}

	// Create the file with the class name
	filePath := fmt.Sprintf("%s/%s.json", dirPath, strings.ToLower(className))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	defer file.Close()

	funcMap := template.FuncMap{
		"add":         func(a, b int) int { return a + b },
		"toCamelCase": toCamelCase,
		"toSnakeCase": toSnakeCase,
	}

	tmpl, err := template.New("talentJson").Funcs(funcMap).Parse(talentJsonTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(file, tabs); err != nil {
		return fmt.Errorf("error executing template for %s: %w", className, err)
	}

	fmt.Printf("Generated %s.json\n", strings.ToLower(className))
	return nil
}

func transformRawTalentsToTabs(rawTalents []RawTalent) ([]TalentTabConfig, error) {
	tabsMap := make(map[string]*TalentTabConfig)
	for _, rt := range rawTalents {
		tab, exists := tabsMap[rt.TabName]
		if !exists {
			tab = &TalentTabConfig{
				Name:          rt.TabName,
				BackgroundUrl: fmt.Sprintf("https://wow.zamimg.com/images/wow/talents/backgrounds/cata/%s.jpg", rt.BackgroundFile),
				Talents:       []TalentConfig{},
			}
			tabsMap[rt.TabName] = tab
		}

		var spellIds []int
		if err := json.Unmarshal([]byte(rt.SpellRank), &spellIds); err != nil {
			return nil, fmt.Errorf("parsing SpellRank for talent %s: %w", rt.TalentName, err)
		}

		filtered := []int{}
		for _, id := range spellIds {
			if id != 0 {
				filtered = append(filtered, id)
			}
		}

		maxPoints := len(filtered)
		fieldName := strings.ToLower(rt.TalentName[:1]) + rt.TalentName[1:]
		talent := TalentConfig{
			FieldName: fieldName,
			FancyName: rt.TalentName,
			Location: TalentLocation{
				RowIdx: rt.TierID,
				ColIdx: rt.ColumnIndex,
			},
			SpellIds:  filtered,
			MaxPoints: maxPoints,
		}

		if (rt.PrereqRow.Valid && rt.PrereqRow.Int64 != 0) || (rt.PrereqCol.Valid && rt.PrereqCol.Int64 != 0) {
			talent.PrereqLocation = &TalentLocation{
				RowIdx: int(rt.PrereqRow.Int64),
				ColIdx: int(rt.PrereqCol.Int64),
			}
		}

		tab.Talents = append(tab.Talents, talent)
		slices.SortFunc(tab.Talents, func(a, b TalentConfig) int {
			return cmp.Or(
				cmp.Compare(a.Location.RowIdx, b.Location.RowIdx),
				cmp.Compare(a.Location.ColIdx, b.Location.ColIdx),
			)
		})
	}

	var tabs []TalentTabConfig

	for _, t := range tabsMap {
		tabs = append(tabs, *t)
	}
	slices.SortFunc(tabs, func(a, b TalentTabConfig) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fieldNum := 1
	for i := range tabs {
		for j := range tabs[i].Talents {
			tabs[i].Talents[j].ProtoFieldNumber = fieldNum
			fieldNum++
		}
	}
	return tabs, nil
}

func transformRawTalentsToConfigsForClass(rawTalents []RawTalent, classID int) ([]TalentConfig, error) {
	var talents []TalentConfig
	for _, rt := range rawTalents {

		converted := convertTalentClassID(classID)
		if converted == rt.ClassMask {
			var spellIds []int
			if err := json.Unmarshal([]byte(rt.SpellRank), &spellIds); err != nil {
				return nil, fmt.Errorf("parsing SpellRank for talent %s: %w", rt.TalentName, err)
			}

			filtered := []int{}
			for _, id := range spellIds {
				if id != 0 {
					filtered = append(filtered, id)
				}
			}

			maxPoints := len(filtered)
			fieldName := strings.ToLower(rt.TalentName[:1]) + rt.TalentName[1:]
			talent := TalentConfig{
				FieldName: fieldName,
				FancyName: rt.TalentName,
				Location: TalentLocation{
					RowIdx: rt.TierID,
					ColIdx: rt.ColumnIndex,
				},
				SpellIds:  filtered,
				MaxPoints: maxPoints,
			}

			if (rt.PrereqRow.Valid && rt.PrereqRow.Int64 != 0) || (rt.PrereqCol.Valid && rt.PrereqCol.Int64 != 0) {
				talent.PrereqLocation = &TalentLocation{
					RowIdx: int(rt.PrereqRow.Int64),
					ColIdx: int(rt.PrereqCol.Int64),
				}
			}

			talents = append(talents, talent)
		}
	}
	slices.SortFunc(talents, func(a, b TalentConfig) int {
		return cmp.Or(
			cmp.Compare(a.Location.RowIdx, b.Location.RowIdx),
			cmp.Compare(a.Location.ColIdx, b.Location.ColIdx),
		)
	})

	return talents, nil
}

func GenerateTalentJsonFromDB(dbHelper *DBHelper) error {
	rawTalents, err := LoadTalents(dbHelper)
	if err != nil {
		return fmt.Errorf("error loading talents: %w", err)
	}

	for _, dbc := range dbcClasses {
		className := classNameFromDBC(dbc)

		classTalents := []RawTalent{}
		for _, rt := range rawTalents {
			converted := convertTalentClassID(dbc.ID)
			if converted == rt.ClassMask {
				classTalents = append(classTalents, rt)
			}
		}

		tabs, err := transformRawTalentsToTabs(classTalents)
		if err != nil {
			fmt.Printf("Error transforming talents for %s: %v\n", className, err)
			continue
		}

		if err := generateTalentJson(tabs, className); err != nil {
			fmt.Printf("Error generating talent JSON for %s: %v\n", className, err)
		}
	}

	return nil
}

func glyphBelongsToClass(r RawGlyph, dbc dbc.DbcClass) bool {
	return r.ClassMask == int32(dbc.ID)
}
func properTitle(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}
func convertRawGlyphToGlyph(r RawGlyph) Glyph {
	return Glyph{
		EnumName: strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(properTitle(r.Name), ":", ""),
					"'", ""),
				"-", ""),
			" ", ""),
		Name:        r.Name,
		Description: template.JSEscapeString(r.Description),
		IconUrl:     "",
		ID:          int(r.ItemId),
	}
}

func GenerateProtos() {
	helper, err := NewDBHelper()
	if err != nil {
		fmt.Printf("Error creating DB helper: %v\n", err)
		return
	}
	defer helper.Close()

	rawGlyphs, err := LoadGlyphs(helper)
	if err != nil {
		fmt.Printf("Error loading glyphs: %v\n", err)
		return
	}

	rawTalents, err := LoadTalents(helper)
	if err != nil {
		fmt.Printf("Error loading talents: %v\n", err)
		return
	}

	var classesData []ClassData
	iconsMap, _ := LoadArtTexturePaths("./assets/db_inputs/ArtTextureID.lua")
	for _, dbc := range dbcClasses {
		className := classNameFromDBC(dbc)
		data := ClassData{
			ClassName:          className,
			LowerCaseClassName: strings.ToLower(className),
			Talents:            []TalentConfig{},
			TalentTabs:         []TalentTabConfig{},
			GlyphsPrime:        []Glyph{},
			GlyphsMajor:        []Glyph{},
			GlyphsMinor:        []Glyph{},
		}

		// Process glyphs
		for _, raw := range rawGlyphs {
			if strings.Contains(raw.Name, "Deprecated") || strings.Contains(raw.Name, "zzz") {
				continue
			}
			if glyphBelongsToClass(raw, dbc) {
				g := convertRawGlyphToGlyph(raw)
				g.IconUrl = "https://wow.zamimg.com/images/wow/icons/large/" + strings.ToLower(GetIconName(iconsMap, int(raw.FDID))) + ".jpg"
				switch raw.GlyphType {
				case 2: // prime
					data.GlyphsPrime = append(data.GlyphsPrime, g)
				case 0: // major
					data.GlyphsMajor = append(data.GlyphsMajor, g)
				case 1: // minor
					data.GlyphsMinor = append(data.GlyphsMinor, g)
				default:
					fmt.Printf("Unknown glyph type %d in raw glyph %+v\n", raw.GlyphType, raw)
				}
			}
		}
		classTalents := []RawTalent{}
		for _, rt := range rawTalents {
			converted := convertTalentClassID(dbc.ID)
			if converted == rt.ClassMask {
				classTalents = append(classTalents, rt)
			}
		}
		talents, err := transformRawTalentsToConfigsForClass(rawTalents, dbc.ID)
		if err != nil {
			fmt.Printf("Error processing talents for %s: %v\n", className, err)
		}
		talentTabs, err := transformRawTalentsToTabs(classTalents)
		slices.SortFunc(talentTabs, func(a, b TalentTabConfig) int {
			return cmp.Compare(a.Name, b.Name)
		})
		if err != nil {
			fmt.Printf("Error grouping talents for %s: %v\n", className, err)
		}
		var filteredTabs []TalentTabConfig
		for _, tab := range talentTabs {
			var filteredTalents []TalentConfig
			for _, t := range tab.Talents {
				if convertTalentClassID(t.MaxPoints) == convertTalentClassID(dbc.ID) {
					filteredTalents = append(filteredTalents, t)
				}
			}
			if len(filteredTalents) > 0 {
				tab.Talents = filteredTalents
				filteredTabs = append(filteredTabs, tab)
			}
		}
		data.Talents = talents
		data.TalentTabs = talentTabs

		classesData = append(classesData, data)
	}

	for _, classData := range classesData {
		if err := generateProtoFile(classData); err != nil {
			fmt.Printf("Error generating proto file for %s: %v\n", classData.ClassName, err)
		} else {
			fmt.Printf("Generated proto/%s.generated.proto\n", strings.ToLower(classData.ClassName))
		}

		if err := generateTsFile(classData); err != nil {
			fmt.Printf("Error generating TS file for %s: %v\n", classData.ClassName, err)
		} else {
			fmt.Printf("Generated %s.ts\n", strings.ToLower(classData.ClassName))
		}
	}

	if err := GenerateTalentJsonFromDB(helper); err != nil {
		fmt.Printf("Error generating talent json files: %v\n", err)
	}
}

func toSnakeCase(s string) string {
	var words []string
	// Split by whitespace.
	for _, w := range strings.Fields(s) {
		var b strings.Builder
		for _, r := range w {
			switch {
			case r == '-':
				b.WriteRune('_')
			case r == '\'':
				continue
			case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
				b.WriteRune(unicode.ToLower(r))
			default:
				continue
			}
		}
		cleaned := b.String()
		if cleaned != "" {
			words = append(words, cleaned)
		}
	}
	return strings.Join(words, "_")
}
func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	var result strings.Builder
	for i, w := range words {
		var b strings.Builder
		for _, r := range w {
			switch {
			case r == '-' || r == '\'':
				continue
			case unicode.IsLetter(r) || unicode.IsDigit(r):
				b.WriteRune(unicode.ToLower(r))
			default:
				continue
			}
		}
		cleaned := b.String()
		if cleaned == "" {
			continue
		}
		if i == 0 {
			result.WriteString(cleaned)
		} else {
			result.WriteString(strings.Title(cleaned))
		}
	}
	return result.String()
}
