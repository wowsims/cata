package database

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"github.com/wowsims/mop/tools/database/dbc"
	"github.com/wowsims/mop/tools/tooltip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func convertTalentClassID(raw int) int {
	return 1 << (raw - 1)
}

type TalentConfig struct {
	FieldName        string         `json:"fieldName"`
	FancyName        string         `json:"fancyName"`
	Location         TalentLocation `json:"location"`
	SpellId          int            `json:"spellId"`
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
	TalentTab          TalentTabConfig
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
{{- range $talent := .TalentTab.Talents }}
    bool {{ final $talent.FancyName $class }} = {{ $talent.ProtoFieldNumber }};
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

const talentJsonTemplate = `
{
	"backgroundUrl": "{{ .BackgroundUrl }}",
	"talents": [
	{{- $m := len .Talents }}
	{{- range $j, $talent := .Talents }}
		{
			"fieldName": "{{ toCamelCase $talent.FancyName }}",
			"fancyName": "{{ $talent.FancyName }}",
			"location": {
				"rowIdx": {{ $talent.Location.RowIdx }},
				"colIdx": {{ $talent.Location.ColIdx }}
			},
			"spellId": {{ $talent.SpellId }}
		}{{ if ne (add $j 1) $m }},{{ end }}
	{{- end }}
	]
}
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
	BackgroundUrl string         `json:"backgroundUrl"`
	Talents       []TalentConfig `json:"talents"`
}

func generateTalentJson(tab TalentTabConfig, className string) error {
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

	if err := tmpl.Execute(file, tab); err != nil {
		return fmt.Errorf("error executing template for %s: %w", className, err)
	}

	fmt.Printf("Generated %s.json\n", strings.ToLower(className))
	return nil
}

func transformRawTalentsToTab(rawTalents []RawTalent) (TalentTabConfig, error) {
	tab := TalentTabConfig{
		BackgroundUrl: fmt.Sprintf("https://wow.zamimg.com/images/wow/talents/backgrounds/cata/%s.jpg", "TODO"),
		Talents:       []TalentConfig{},
	}

	for _, rt := range rawTalents {
		fieldName := strings.ToLower(rt.TalentName[:1]) + rt.TalentName[1:]
		talent := TalentConfig{
			FieldName: fieldName,
			FancyName: rt.TalentName,
			Location: TalentLocation{
				RowIdx: rt.TierID,
				ColIdx: rt.ColumnIndex,
			},
			SpellId: rt.SpellID,
		}

		tab.Talents = append(tab.Talents, talent)
		slices.SortFunc(tab.Talents, func(a, b TalentConfig) int {
			return cmp.Or(
				cmp.Compare(a.Location.RowIdx, b.Location.RowIdx),
				cmp.Compare(a.Location.ColIdx, b.Location.ColIdx),
			)
		})
	}

	fieldNum := 1
	for i := range tab.Talents {
		tab.Talents[i].ProtoFieldNumber = fieldNum
		fieldNum++
	}

	return tab, nil
}

func transformRawTalentsToConfigsForClass(rawTalents []RawTalent, classID int) ([]TalentConfig, error) {
	var talents []TalentConfig

	for _, rt := range rawTalents {
		if classID == rt.ClassMask {
			fieldName := strings.ToLower(rt.TalentName[:1]) + rt.TalentName[1:]
			talent := TalentConfig{
				FieldName: fieldName,
				FancyName: rt.TalentName,
				Location: TalentLocation{
					RowIdx: rt.TierID,
					ColIdx: rt.ColumnIndex,
				},
				SpellId: rt.SpellID,
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

	for _, dbcClass := range dbc.Classes {
		className := dbc.ClassNameFromDBC(dbcClass)

		classTalents := []RawTalent{}
		for _, rt := range rawTalents {
			if dbcClass.ID == rt.ClassMask {
				classTalents = append(classTalents, rt)
			}
		}

		tab, err := transformRawTalentsToTab(classTalents)
		if err != nil {
			fmt.Printf("Error transforming talents for %s: %v\n", className, err)
			continue
		}

		if err := generateTalentJson(tab, className); err != nil {
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
func convertRawGlyphToGlyph(r RawGlyph, dbc *dbc.DBC) Glyph {
	tooltip, _ := tooltip.ParseTooltip(r.Description, tooltip.DBCTooltipDataProvider{DBC: dbc}, int64(r.SpellId))
	return Glyph{
		EnumName: strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(properTitle(r.Name), ":", ""),
					"'", ""),
				"-", ""),
			" ", ""),
		Name:        r.Name,
		Description: template.JSEscapeString(tooltip.String()),
		IconUrl:     "",
		ID:          int(r.ItemId),
	}
}

func GenerateProtos(dbcData *dbc.DBC) {
	helper, err := NewDBHelper()
	if err != nil {
		fmt.Printf("Error creating DB helper: %v\n", err)
		return
	}
	defer helper.Close()

	var ignoredGlyphs = []int32{102153, 104054}
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
	iconsMap, _ := LoadArtTexturePaths("./tools/DB2ToSqlite/listfile.csv")
	for _, dbcClass := range dbc.Classes {
		className := dbc.ClassNameFromDBC(dbcClass)
		data := ClassData{
			ClassName:          className,
			LowerCaseClassName: strings.ToLower(className),
			Talents:            []TalentConfig{},
			TalentTab:          TalentTabConfig{},
			GlyphsMajor:        []Glyph{},
			GlyphsMinor:        []Glyph{},
		}

		// Process glyphs
		for _, raw := range rawGlyphs {
			if slices.Contains(ignoredGlyphs, raw.ItemId) || strings.Contains(raw.Name, "Deprecated") || strings.Contains(raw.Name, "zzz") || (len(raw.Name) > 2 && raw.Name[:2] == "zz") {
				continue
			}
			if glyphBelongsToClass(raw, dbcClass) {
				g := convertRawGlyphToGlyph(raw, dbcData)
				g.IconUrl = "https://wow.zamimg.com/images/wow/icons/large/" + strings.ToLower(GetIconName(iconsMap, int(raw.FDID))) + ".jpg"
				switch raw.GlyphType {
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
			if dbcClass.ID == rt.ClassMask {
				classTalents = append(classTalents, rt)
			}
		}
		talents, err := transformRawTalentsToConfigsForClass(rawTalents, dbcClass.ID)
		if err != nil {
			fmt.Printf("Error processing talents for %s: %v\n", className, err)
		}
		talentTab, err := transformRawTalentsToTab(classTalents)

		if err != nil {
			fmt.Printf("Error grouping talents for %s: %v\n", className, err)
		}

		data.Talents = talents
		data.TalentTab = talentTab

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
