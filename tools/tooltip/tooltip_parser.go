package tooltip

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/wowsims/mop/sim/core"
)

type TooltipDataProvider interface {
	GetAttackPower() float64
	GetDescriptionVariableString(spellId int64) string
	GetEffectAmplitude(spellId int64, effectIdx int64) float64
	GetEffectScaledValue(spellId int64, effectIdx int64) float64
	GetEffectBaseValue(spellId int64, effectIdx int64) float64 // basePoints + ?
	GetEffectChainAmplitude(spellId int64, effectidx int64) float64
	GetEffectMaxTargets(spellId int64, effectIdx int64) int64
	GetEffectPeriod(spellId int64, effectIdx int64) time.Duration
	GetEffectPointsPerResource(spellId int64, effectIdx int64) float64
	GetEffectRadius(spellId int64, effectIdx int64) float64
	GetEffectEnchantValue(enchantId int64, effectidx int64) float64
	GetMainHandWeapon() *core.Weapon
	GetOffHandWeapon() *core.Weapon
	GetPlayerLevel() float64
	GetSpecNum() int64 // The spec index for the class. Basically left to right in the talent window. i.E. Balance = 0, Guardian = 1, Feral = 2, Restoration = 4
	GetSpellDescription(spellId int64) string
	GetSpellDuration(spellId int64) time.Duration
	GetSpellIcon(spellId int64) string
	GetSpellMaxTargets(spellId int64) int64
	GetSpellName(spellid int64) string
	GetSpellPower() float64
	GetSpellPPM(spellId int64) float64
	GetSpellProcChance(spellId int64) float64
	GetSpellProcCooldown(spellId int64) time.Duration
	GetSpellRange(spellId int64) float64
	GetSpellStacks(spellId int64) int64 // Should return SpellAuraOptions ProcCharges or CumulativeAura
	HasAura(auraId int64) bool
	HasPassive(auraId int64) bool
	IsMaleGender() bool
	KnowsSpell(spellId int64) bool
}

// ****************************
// MATH Handling
//
// Base on participle sample math parser
//

type MathOperator int

const (
	OpMul MathOperator = iota
	OpDiv
	OpAdd
	OpSub
)

var operatorMap = map[string]MathOperator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

func (o *MathOperator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

type MathValue struct {
	Number        *float64       `parser:"  @('-'? (Float|Int))"`
	Variable      *string        `parser:"| @Ident"`
	ComputedValue *ComputedValue `parser:"| @@"`
	Subexpression *Expression    `parser:"| '(' @@ ')'"`
}

type MathFactorTerm struct {
	Base     *MathValue `parser:"@@"`
	Exponent *MathValue `parser:"( '^' @@ )?"`
}

type MathMultTerm struct {
	Operator MathOperator    `parser:"@('*' | '/')"`
	Factor   *MathFactorTerm `parser:"@@"`
}

type MathSimpleTerm struct {
	Left  *MathFactorTerm `parser:"@@"`
	Right []*MathMultTerm `parser:"@@*"`
}

type MathAddTerm struct {
	Operator MathOperator    `parser:"@('+' | '-')"`
	Term     *MathSimpleTerm `parser:"@@"`
}

type Expression struct {
	Left  *MathSimpleTerm `parser:"@@"`
	Right []*MathAddTerm  `parser:"@@*"`
}

type MathExpression struct {
	Expression *Expression `parser:"'${'@@ ')'*'}'"`
	Round      *float64    `parser:"@Float?"`
}

func (m MathExpression) GetDecimalPlace() int64 {
	if m.Round != nil {
		return int64(*m.Round * 10)
	}

	return 0
}

func (m MathExpression) Eval(ctx *TooltipContext) float64 {
	val := m.Expression.Eval(ctx)
	if m.Round == nil {
		return val
	}

	factor := math.Pow(10, *m.Round*10)
	return math.Round(val*factor) / factor
}

func (o MathOperator) String() string {
	switch o {
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpSub:
		return "-"
	case OpAdd:
		return "+"
	}
	panic("unsupported operator")
}

func (o MathOperator) Eval(l, r float64) float64 {
	switch o {
	case OpMul:
		return l * r
	case OpDiv:
		return l / r
	case OpAdd:
		return l + r
	case OpSub:
		return l - r
	}
	panic("unsupported operator")
}

func (v *MathValue) Eval(ctx *TooltipContext) float64 {
	switch {
	case v.Number != nil:
		return *v.Number
	case v.ComputedValue != nil:
		return v.ComputedValue.Eval(ctx)
	default:
		return v.Subexpression.Eval(ctx)
	}
}

func (f *MathFactorTerm) Eval(ctx *TooltipContext) float64 {
	b := f.Base.Eval(ctx)
	if f.Exponent != nil {
		return math.Pow(b, f.Exponent.Eval(ctx))
	}
	return b
}

func (t *MathSimpleTerm) Eval(ctx *TooltipContext) float64 {
	n := t.Left.Eval(ctx)
	for _, r := range t.Right {
		n = r.Operator.Eval(n, r.Factor.Eval(ctx))
	}
	return n
}

func (e *Expression) Eval(ctx *TooltipContext) float64 {
	l := e.Left.Eval(ctx)
	for _, r := range e.Right {
		l = r.Operator.Eval(l, r.Term.Eval(ctx))
	}
	return l
}

// ****************************
// BOOL HANDLING
//
// Handles spell conditions
// Those are usually in the form of (a|p|s)SpellID
// a = HasAura
// p = HasPassive
// s = KnowsSpell
//
// There exists one exception of such conditions that is (c)[1-4]
// c = Current specialization index. It can be 1 to 4, top to bottom in the talent window
type SpellCondition struct {
	Op      string
	SpellId int64
}

func (s *SpellCondition) Capture(values []string) error {
	s.Op = values[0][:1]
	s.SpellId, _ = strconv.ParseInt(values[0][1:], 10, 64)
	return nil
}

type BoolTerminalValue struct {
	Negated        *string           `parser:"@NOT?"`
	SpellCondition *SpellCondition   `parser:"@SpellCond"`
	SpellRef       *SimpleSpellValue `parser:"|@@"`
	Number         *float64          `parser:"| @(Float|Int)"`
	Subexpression  *BoolExpression   `parser:"| '(' @@ ')'"`
}

type BoolCompareBranch struct {
	Operator string             `parser:"@BOP"`
	Value    *BoolTerminalValue `parser:"@@"`
}

type BoolGeneralTerm struct {
	Left  *BoolTerminalValue   `parser:"@@"`
	Right []*BoolCompareBranch `parser:"@@*"`
}

// For boolean chaining with & and |
type BoolChainTerm struct {
	Operator string           `parser:"@BOC"`
	Term     *BoolGeneralTerm `parser:"@@"`
}

// Root term for boolean expressions
type BoolExpression struct {
	Left  *BoolGeneralTerm `parser:"@@"`
	Right []*BoolChainTerm `parser:"@@*"`
}

func (b BoolTerminalValue) Eval(ctx *TooltipContext) float64 {
	if b.Number != nil {
		return *b.Number
	}

	if b.SpellRef != nil {
		return b.SpellRef.Eval(ctx)
	}

	if b.EvalBool(ctx) {
		return 1
	}

	return 0
}

func (s SpellCondition) EvalBool(ctx *TooltipContext) bool {
	op := s.Op
	if op[:1] == "?" {
		op = op[1:]
	}

	switch op {
	case "c": // class
		return ctx.DataProvider.GetSpecNum() == s.SpellId
	case "s":
		return ctx.DataProvider.KnowsSpell(s.SpellId)
	case "p":
		fallthrough
	case "a":
		return ctx.DataProvider.HasAura(s.SpellId)
	default:
		panic("Unsupported spell condition")
	}
}

func (b BoolTerminalValue) EvalBool(ctx *TooltipContext) bool {
	if b.SpellCondition != nil {
		result := b.SpellCondition.EvalBool(ctx)
		if b.Negated != nil {
			return !result
		}

		return result
	}

	if b.Subexpression != nil {
		return b.Subexpression.Eval(ctx)
	}

	return false
}

func (b BoolGeneralTerm) Eval(ctx *TooltipContext) bool {
	if len(b.Right) == 0 {
		return b.Left.EvalBool(ctx)
	}

	left := b.Left.Eval(ctx)
	right := b.Right[0].Value.Eval(ctx)
	switch b.Right[0].Operator {
	case ">":
		return left > right
	case "<":
		return left < right
	case "!=":
		return left != right
	case "=":
		return left == right
	}

	return false
}

func (b BoolExpression) Eval(ctx *TooltipContext) bool {
	if len(b.Right) == 0 {
		return b.Left.Eval(ctx)
	}

	left := b.Left.Eval(ctx)
	for _, right := range b.Right {
		rVal := right.Term.Eval(ctx)
		switch right.Operator {
		case "|":
			left = left || rVal
		case "&":
			left = left && rVal
		}
	}

	return left
}

// Parses values that reference spell params like $m1 or $5565s1
// If a spellID is given the lookup is performed on another spell's values
type SimpleSpellValue struct {
	SpellId  *int64        `parser:"'$'?@Int?"`
	Selector SpellEntryRef `parser:"@SpMod"`
}

type SpellEntryRef struct {
	EffectColumn string
	EffectIndex  int64
}

// Due to some complications the the lexing context and many inconsistencies on blizzards side
// We need to manually capture the effect column and index
func (s *SpellEntryRef) Capture(values []string) error {
	val := values[0]
	(*s).EffectColumn = val[:1]
	if len(val) > 1 {
		(*s).EffectIndex, _ = strconv.ParseInt(val[1:], 10, 64)
	}

	return nil
}

// External lookups and references
// Usually $<varname> variables are dynamically computed through SpellXSpellDescription references
// While $[a-zA-Z]{2,} variables are typically statically provided to the tooltip rendering context
// They reference character stats and static spell values like proc chance and RPPM
type VariableRef struct {
	DynamicVarName string `parser:"'$<'@VarRefName'>'"`
	StaticVarName  string `parser:"|'$'@VarName (?!'(')"`
}

type DescriptionRef struct {
	SpellId int64 `parser:"@Int"`
}

type SpellNameRef struct {
	SpellId int64 `parser:"@Int"`
}

type SpellIconRef struct {
	SpellId int64 `parser:"@Int"`
}

type TernaryRightSide struct {
	SecondValue *TooltipAST `parser:"'['@@']'"`
	Chained     *Ternary    `parser:"|@@"`
}
type Ternary struct {
	BoolExpr    *BoolExpression   `parser:"@@"`
	FirstValue  *TooltipAST       `parser:"'['@@(']'|']?')"`
	SecondValue *TernaryRightSide `parser:"@@"`
}

type VariableAssignment struct {
	VariableName string        `parser:"'$'@VarName'='"`
	Value        *ComplexValue `parser:"@@"`
}

type SimpleCompute struct {
	Op    string            `parser:"'$'@('/' | '*')"`
	Num   int64             `parser:"@Int"`
	Value *SimpleSpellValue `parser:"';'@@"`
}

// Tooltips support functions like clamp / min / max
type Function struct {
	Name string        `parser:"'$'@VarName"`
	Args *[]Expression `parser:"'(' @@ (',' @@)* ')'"`
}

func (f *Function) Eval(ctx *TooltipContext) float64 {
	switch strings.ToLower(f.Name) {
	case "max":
		left := (*f.Args)[0].Eval(ctx)
		right := (*f.Args)[1].Eval(ctx)
		if left > right {
			return left
		}

		return right
	case "floor":
		arg := (*f.Args)[0].Eval(ctx)
		return math.Floor(arg)
	case "cond":
		cond := (*f.Args)[0].Eval(ctx)
		if cond > 0 {
			return (*f.Args)[1].Eval(ctx)
		}

		return (*f.Args)[2].Eval(ctx)
	case "gt": // implement java like compare
		left := (*f.Args)[0].Eval(ctx)
		right := (*f.Args)[1].Eval(ctx)
		if left > right {
			return 1
		}

		return 0
	case "gte":
		left := (*f.Args)[0].Eval(ctx)
		right := (*f.Args)[1].Eval(ctx)
		if left >= right {
			return 1
		}

		return 0
	case "lt":
		left := (*f.Args)[0].Eval(ctx)
		right := (*f.Args)[1].Eval(ctx)
		if left < right {
			return 1
		}

		return 0
	case "clamp":
		arg := (*f.Args)[0].Eval(ctx)
		min := (*f.Args)[1].Eval(ctx)
		max := (*f.Args)[2].Eval(ctx)
		return math.Min(math.Max(arg, min), max)
	default:
		panic(f.Name + " not implmemented")
	}
}

// Short ternaryies like $lowner:owners are hard to lex vers $low=
// So some magic here..
//
// $gmal:female
// $lsingular:plural <- seems to be based on last compuated element
type ShortTernary struct {
	Type  string `parser:"@ShortTern"`
	Right string `parser:"@Option';'"`
}

type NegativeComputedValue struct {
	Negative *string        `parser:"'-'"`
	Value    *ComputedValue `parser:"@@"`
}

type ComputedValue struct {
	Negative      *NegativeComputedValue `parser:"@@"`
	SimpleCompute *SimpleCompute         `parser:"|@@"`
	VariableRef   *VariableRef           `parser:"|@@"`
	Function      *Function              `parser:"|@@"`
	SpellValue    *SimpleSpellValue      `parser:"|@@"`
}

type ComplexValue struct {
	MathExpression     *MathExpression     `parser:"@@"`
	VariableAssignment *VariableAssignment `parser:"|@@"`
	Terniary           *Ternary            `parser:"|'$?'@@"`
	ShortTernary       *ShortTernary       `parser:"|@@"`
	ComputedValue      *ComputedValue      `parser:"|@@"`
	Word               *string             `parser:"|@(Ident|Int|Op|Float)"`
	Punctuation        *string             `parser:"|@Punct"`
	DescriptionRef     *DescriptionRef     `parser:"|DescLookup @@"`
	SpellNameRef       *SpellNameRef       `parser:"|'$@spellname'@@"`
	SpellIconRef       *SpellIconRef       `parser:"|'$@spellicon'@@"`
}

func (c ComplexValue) isPunctuation() bool {
	return c.Punctuation != nil
}

type Word struct {
	Word string `parser:"@Ident"`
}

type TooltipAST struct {
	Values *[]ComplexValue `parser:"@@*"`
}

var ternParser = regexp.MustCompile(`^\$(\d*)([glLG])(.+):$`)

func (d DescriptionRef) String(ctx *TooltipContext) string {
	desc := ctx.DataProvider.GetSpellDescription(d.SpellId)
	if len(desc) == 0 {
		return desc
	}

	result, error := ParseTooltip(desc, ctx.DataProvider, d.SpellId)
	if error != nil {
		return ""
	}

	return result.String()
}

func (s SpellNameRef) String(ctx *TooltipContext) string {
	return ctx.DataProvider.GetSpellName(s.SpellId)
}

func (s SpellIconRef) String(ctx *TooltipContext) string {
	iconPath := ctx.DataProvider.GetSpellIcon(s.SpellId)
	if len(iconPath) == 0 {
		return ""
	}

	return fmt.Sprintf("|T%s:24|t", iconPath)
}

func (s ShortTernary) String(ctx *TooltipContext) string {

	// Some spells have short hands as $89808lspell:spells
	// The only way I could make sense of the 'l' operator is a ref to the last computed value
	// so spells that are not this do not make sense. We ignore them
	// We will get this as $\d*[l|g]\w+: cleanup a bit.

	match := ternParser.FindStringSubmatch(s.Type)
	t := match[2]
	left := match[3]
	switch t {
	case "G":
		fallthrough
	case "g":
		if ctx.DataProvider.IsMaleGender() {
			return left
		}

		return s.Right
	case "L":
		fallthrough
	case "l":
		if ctx.LastEval <= 1 {
			return left
		}

		return s.Right
	default:
		panic(s.Type + " is unsupported shorthand ternary.")
	}
}

func (c ComputedValue) Eval(ctx *TooltipContext) float64 {
	switch {
	case c.Negative != nil:
		return -1 * c.Negative.Value.Eval(ctx)
	case c.SimpleCompute != nil:
		return c.SimpleCompute.Eval(ctx)
	case c.SpellValue != nil:
		return c.SpellValue.Eval(ctx)
	case c.VariableRef != nil:
		return c.VariableRef.getValue(ctx)
	case c.Function != nil:
		return c.Function.Eval(ctx)
	default:
		panic("Invalid computed value")
	}
}

func (s SimpleSpellValue) getSpellId(ctx *TooltipContext) int64 {
	if s.SpellId != nil {
		return *s.SpellId
	}

	return ctx.SpellId
}

func (v VariableRef) getValue(ctx *TooltipContext) float64 {
	varName := v.DynamicVarName
	if len(varName) == 0 {
		varName = v.StaticVarName
	}

	value, ok := ctx.Variables[varName]
	if ok {
		ctx.LastEval = value
		return value
	}

	fmt.Printf("Warning: Variable (%s) not registered. Using default value (0)\n", varName)
	return 0.0
}

func (v VariableRef) String(ctx *TooltipContext) string {
	return fmt.Sprintf("%0.f", v.getValue(ctx))
}

func (s SimpleCompute) Eval(ctx *TooltipContext) float64 {
	switch s.Op {
	case "/":
		return s.Value.Eval(ctx) / float64(s.Num)
	case "*":
		return s.Value.Eval(ctx) / float64(s.Num)
	default:
		panic("OP not implemented")
	}
}

func (s SimpleCompute) String(ctx *TooltipContext) string {
	right := s.Value.Eval(ctx)
	ctx.LastEval = right
	if s.Op == "/" {
		return fmt.Sprintf("%.1f", right/float64(s.Num))
	} else {
		return fmt.Sprintf("%.1f", right*float64(s.Num))
	}
}

func (s SimpleSpellValue) Eval(ctx *TooltipContext) float64 {
	switch strings.ToLower(s.Selector.EffectColumn) {
	case "e":
		return ctx.DataProvider.GetEffectAmplitude(s.getSpellId(ctx), s.Selector.EffectIndex-1)
	case "h":
		return ctx.DataProvider.GetSpellProcChance(s.getSpellId(ctx))
	case "d":
		return float64(ctx.DataProvider.GetSpellDuration(s.getSpellId(ctx))) / float64(time.Second)
	case "w":
		// This does not properly evaluate in client for Spell Descriptions. In theory it seems to refer to the specific extra values of a buff
		// i.E. the actual stamina buffed by a priester to display it correctly client side
		// So for Buff Tooltip rendering we want to treat at probably the same as scaled effect value
		fallthrough
	case "s":
		return ctx.DataProvider.GetEffectScaledValue(s.getSpellId(ctx), s.Selector.EffectIndex-1)
	case "M":
		fallthrough
	case "m":
		return ctx.DataProvider.GetEffectBaseValue(s.getSpellId(ctx), s.Selector.EffectIndex-1)
	case "T":
		fallthrough
	case "t":
		return float64(ctx.DataProvider.GetEffectPeriod(s.getSpellId(ctx), s.Selector.EffectIndex-1))
	case "x":
		return float64(ctx.DataProvider.GetEffectMaxTargets(s.getSpellId(ctx), s.Selector.EffectIndex-1))
	case "i":
		return float64(ctx.DataProvider.GetSpellMaxTargets(s.getSpellId(ctx)))
	case "F":
		fallthrough
	case "f":
		return float64(ctx.DataProvider.GetEffectChainAmplitude(s.getSpellId(ctx), s.Selector.EffectIndex-1))
	case "o":
		spellId := s.getSpellId(ctx)
		baseDamage := ctx.DataProvider.GetEffectScaledValue(spellId, s.Selector.EffectIndex-1)
		period := ctx.DataProvider.GetEffectPeriod(spellId, s.Selector.EffectIndex-1)
		if period == 0 {
			return 0
		}

		duration := ctx.DataProvider.GetSpellDuration(spellId)
		return float64(duration) / float64(period) * baseDamage

		// maybe R = Max Range, r = min range?
	case "R":
		fallthrough
	case "r":
		if s.Selector.EffectIndex == 0 {
			return ctx.DataProvider.GetSpellRange(s.getSpellId(ctx))
		}
		return ctx.DataProvider.GetEffectRadius(s.getSpellId(ctx), s.Selector.EffectIndex-1)
	case "A":
		fallthrough
	case "a":
		radius := ctx.DataProvider.GetEffectRadius(s.getSpellId(ctx), s.Selector.EffectIndex-1)
		if radius == 0 {
			radius = float64(ctx.DataProvider.GetSpellRange(s.getSpellId(ctx)))
		}

		return radius

		// DUnno what's the difference
	case "n":
		fallthrough
	case "u":
		return float64(ctx.DataProvider.GetSpellStacks(s.getSpellId(ctx)))
	case "b":
		return float64(ctx.DataProvider.GetEffectPointsPerResource(s.getSpellId(ctx), s.Selector.EffectIndex-1))
	case "V":
		fallthrough // Max Target Level SpelLTargetRestrictions.dbc
	case "v":
		return 0
	case "k":
		return ctx.DataProvider.GetEffectEnchantValue(s.getSpellId(ctx), s.Selector.EffectIndex-1)
	default:
		return 0.0
	}
}

func (s SimpleSpellValue) String(ctx *TooltipContext) string {
	value := s.Eval(ctx)
	ctx.LastEval = value
	switch s.Selector.EffectColumn {
	case "D":
		fallthrough
	case "d":
		// value is returned in seconds, normalize to time value
		value *= float64(time.Second)
		duration := time.Duration(value)
		if value > float64(time.Hour*2) {
			return fmt.Sprintf("%dhrs", duration/time.Hour)
		}
		if value >= float64(time.Hour) {
			return fmt.Sprintf("%dhr", duration/time.Hour)
		}
		if value >= float64(time.Minute) {
			return fmt.Sprintf("%dmin", duration/time.Minute)
		}
		if value < float64(time.Second) {
			return fmt.Sprintf("%dms", duration/time.Millisecond)
		} else {
			return fmt.Sprintf("%ds", duration/time.Second)
		}
	case "e":
		fallthrough
	case "T":
		fallthrough
	case "t":
		return fmt.Sprintf("%.1f", value/float64(time.Second))
	case "s":
		fallthrough
	case "A":
		fallthrough
	case "a":
		fallthrough
	case "f":
		fallthrough
	case "F":
		fallthrough
	case "m":
		fallthrough
	case "h":
		fallthrough
	case "o":
		fallthrough
	case "x":
		fallthrough
	case "u":
		fallthrough
	case "R":
		fallthrough
	case "r":
		fallthrough
	case "i":
		fallthrough
	case "n":
		fallthrough
	case "M":
		fallthrough
	case "b":
		fallthrough
	case "k":

		// Apparently spell ref values are always positive and explicitly prefixed by '-' in the tooltip
		return fmt.Sprintf("%.0f", math.Abs(value))
	default:
		return "{UNK: " + s.Selector.EffectColumn + "}"
	}
}

func (c ComputedValue) String(ctx *TooltipContext) string {
	if c.SpellValue != nil {
		return c.SpellValue.String(ctx)
	}

	if c.SimpleCompute != nil {
		return c.SimpleCompute.String(ctx)
	}

	if c.VariableRef != nil {
		return c.VariableRef.String(ctx)
	}

	return ""
}

func (t TooltipAST) String(ctx *TooltipContext) string {
	var tooltip string
	if t.Values == nil {
		return ""
	}

	for _, val := range *t.Values {
		value := val.String(ctx)
		if len(value) > 0 {
			if len(tooltip) > 0 && !val.isPunctuation() {
				lastChar := tooltip[len(tooltip)-1:]

				// for now do not render a space after + or -
				// might want to change lexing behaviour in root context
				// later on if we find unaccaptable inconsistencies
				if lastChar != "+" && lastChar != "-" {
					tooltip += " "
				}

			}

			tooltip += value
		}
	}

	return tooltip
}

func (c ComplexValue) String(ctx *TooltipContext) string {
	if c.VariableAssignment != nil {
		return ""
	}

	if c.Terniary != nil {
		return (*c.Terniary).String(ctx)
	}

	if c.Word != nil {
		return *c.Word
	}

	if c.ComputedValue != nil {
		return (*c.ComputedValue).String(ctx)
	}

	if c.isPunctuation() {
		return *c.Punctuation
	}

	if c.MathExpression != nil {
		val := c.MathExpression.Eval(ctx)
		ctx.LastEval = val
		decimal := c.MathExpression.GetDecimalPlace()
		format := "%." + strconv.FormatInt(decimal, 10) + "f"
		return fmt.Sprintf(format, val)
	}

	if c.ShortTernary != nil {
		return (*c.ShortTernary).String(ctx)
	}

	if c.SpellIconRef != nil {
		return (*c.SpellIconRef).String(ctx)
	}

	if c.DescriptionRef != nil {
		return (*c.DescriptionRef).String(ctx)
	}

	if c.SpellNameRef != nil {
		return (*c.SpellNameRef).String(ctx)
	}

	return ""
}

func (t Ternary) String(ctx *TooltipContext) string {
	if t.BoolExpr.Eval(ctx) {
		return t.FirstValue.String(ctx)
	} else {
		if t.SecondValue.Chained != nil {
			return t.SecondValue.Chained.String(ctx)
		}

		return t.SecondValue.SecondValue.String(ctx)
	}
}

// define some tooltip fixes
type tooltipFix struct {
	Regex   *regexp.Regexp
	Replace string
}

var fixes = []tooltipFix{
	{Regex: regexp.MustCompile(`\(\<\$`), Replace: "($<"},
	{Regex: regexp.MustCompile(`,\<\$`), Replace: ",$<"},
	{Regex: regexp.MustCompile(`\]\]`), Replace: "]"},
	{Regex: regexp.MustCompile(`(.)\$[bB]([^\d])`), Replace: "$1\n$2"},
	{Regex: regexp.MustCompile(`\)\r\n\[`), Replace: ")["},
	{Regex: regexp.MustCompile(`\]\$\[`), Replace: "]["},
	{Regex: regexp.MustCompile(`\{(\d+[a-zA-Z]\d)`), Replace: "{$$$1"},
	{Regex: regexp.MustCompile(`(\$\?[^\[$]+)\$\?`), Replace: "$1"},
	{Regex: regexp.MustCompile(`([\(|&?][ap]\d+)[a-z]\d`), Replace: "$1"},
}

func applyFixes(tooltip string) string {
	for _, f := range fixes {
		tooltip = f.Regex.ReplaceAllString(tooltip, f.Replace)
	}

	return tooltip
}

func getLexer() *lexer.StatefulDefinition {
	return lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: "CommentStart", Pattern: `--`, Action: lexer.Push("Comment")},
			{Name: "TernStart", Pattern: `(\$|\])\?`, Action: lexer.Push("Ternary")},
			{Name: "DescLookup", Pattern: `\$@(spelldesc|spelltooltip)`, Action: nil},
			{Name: "SpellLookup", Pattern: `\$@spellname`, Action: nil},
			{Name: "IconLookup", Pattern: `\$@spellicon`, Action: nil},
			lexer.Include("Shared"),
			{Name: "Ident", Pattern: `[#a-zA-Z'0-9|()"&][a-zA-Z'0-9"#|()-\\_&]*`, Action: nil},
			{Name: "Tok", Pattern: `[\[\]\{\}=?\<\>]`, Action: nil},
			{Name: "Punct", Pattern: `[.,:\!?%;\]\r\n]`},
			{Name: "SpellCond2", Pattern: `\?[aspc]\d{2,}`, Action: nil},
		},
		"Comment": {
			{Name: "Comment", Pattern: ".+?\n", Action: lexer.Pop()},
		},
		"Boolean": {
			{Name: "BEND", Pattern: `\)`, Action: nil},
			{Name: "BSTART", Pattern: `\(`, Action: nil},
			{Name: "BOC", Pattern: `[\|\&]`, Action: nil},
			{Name: "BOP", Pattern: `([\<\>]|!=|=)`, Action: nil},
			{Name: "NOT", Pattern: `!`, Action: nil},
		},
		"Math": {
			{Name: "MathEnd", Pattern: `\}`, Action: lexer.Pop()},
			{Name: "Brackets", Pattern: `[()]`, Action: nil},
			{Name: "ArgSep", Pattern: `,`, Action: nil},
			{Name: `Op`, Pattern: `[-+/*=]`, Action: nil},
			lexer.Include("Shared"),
		},
		"Shared": {
			{Name: `Whitespace`, Pattern: `[ \t]+`, Action: nil},
			{Name: "MathStart", Pattern: `\$\{`, Action: lexer.Push("Math")},
			{Name: "ShortTern", Pattern: `\$\d*[lgLG][a-zA-Z0-9 ]+:`, Action: lexer.Push("ShortTern")},
			{Name: "Float", Pattern: `-?(\d+)?\.\d+`, Action: nil},
			{Name: "Int", Pattern: `-?\d+`, Action: nil},
			{Name: "VarRef", Pattern: `\$\<`, Action: lexer.Push("VarRef")},
			{Name: `Op`, Pattern: `[-+/*=]`, Action: nil},
			{Name: "Var", Pattern: `\$`, Action: lexer.Push("Variable")},
		},
		"ShortTern": {
			{Name: `Option`, Pattern: `[^;:]+`, Action: nil},
			{Name: `ShortTernEnd`, Pattern: `;`, Action: lexer.Pop()},
			{Name: `ShortTernDiv`, Pattern: `:`, Action: nil},
		},
		"VarRef": {
			{Name: "VarRefEnd", Pattern: `\>`, Action: lexer.Pop()},
			{Name: "VarRefName", Pattern: `[a-zA-Z0-9]+`, Action: nil},
		},
		"Ternary": {
			lexer.Include("Boolean"),
			lexer.Include("Shared"),
			{Name: "SpellCond", Pattern: `[aspc]\d{1,}`, Action: nil},
			{Name: "TernEnd", Pattern: `(\[|=)`, Action: lexer.Pop()},
		},
		"Variable": {
			lexer.Include("Shared"),
			{Name: "SimpleComp", Pattern: `[/*]`, Action: nil},
			{Name: "SimpleTok", Pattern: `[;]`, Action: nil},
			{Name: "VarName", Pattern: `[a-zA-Z]{2,}\d*`, Action: lexer.Pop()},
			{Name: "SpMod", Pattern: "[a-zA-Z][0-9]?", Action: lexer.Pop()},
			{Name: "SpellCond2", Pattern: `\?[aspc]\d{2,}`, Action: lexer.Pop()},
		},
	})
}

type TooltipContext struct {
	DataProvider TooltipDataProvider
	Variables    map[string]float64
	LastEval     float64
	SpellId      int64
}

type Tooltip struct {
	Context *TooltipContext
	AST     *TooltipAST
}

func (t Tooltip) String() string {
	return t.AST.String(t.Context)
}

func ParseTooltip(tooltip string, dataProvider TooltipDataProvider, spellId int64) (*Tooltip, error) {
	def := getLexer()
	parser, error := participle.Build[TooltipAST](participle.Lexer(def), participle.Elide("Whitespace", "Comment", "CommentStart"), participle.UseLookahead(-1))
	if error != nil {
		panic(error)
	}

	// register variables
	variables := map[string]float64{
		"STR":          1,
		"INT":          1,
		"AP":           dataProvider.GetAttackPower(),
		"RAP":          dataProvider.GetAttackPower(),
		"SP":           dataProvider.GetSpellPower(),
		"SPS":          dataProvider.GetSpellPower(),
		"SPFR":         dataProvider.GetSpellPower(),
		"SPN":          dataProvider.GetSpellPower(),
		"SPH":          dataProvider.GetSpellPower(),
		"pctH":         1,
		"PL":           dataProvider.GetPlayerLevel(),
		"pl":           dataProvider.GetPlayerLevel(),
		"proccooldown": dataProvider.GetSpellProcCooldown(spellId).Seconds(),
		"procrppm":     dataProvider.GetSpellPPM(spellId),
	}

	mhWeapon := dataProvider.GetMainHandWeapon()
	if mhWeapon != nil {
		variables["mwb"] = mhWeapon.BaseDamageMin
		variables["mws"] = mhWeapon.SwingSpeed
		variables["MWB"] = mhWeapon.BaseDamageMax
		variables["MWS"] = mhWeapon.SwingSpeed
	}

	ohWeapon := dataProvider.GetOffHandWeapon()
	if ohWeapon != nil {
		variables["owb"] = ohWeapon.BaseDamageMin
		variables["ows"] = ohWeapon.SwingSpeed
		variables["OWB"] = ohWeapon.BaseDamageMax
		variables["OWS"] = ohWeapon.SwingSpeed
	}

	tooltip = applyFixes(tooltip)
	ctx := &TooltipContext{DataProvider: dataProvider, Variables: variables, SpellId: spellId, LastEval: 0}
	varString := applyFixes(dataProvider.GetDescriptionVariableString(spellId))
	vars, error := parser.ParseString("", varString)
	if error != nil {
		return nil, error
	}

	if vars.Values != nil {
		for _, c := range *vars.Values {
			if c.VariableAssignment != nil {
				// Complex values can evaluate to string or float - use float and parse
				result := c.VariableAssignment.Value.String(ctx)
				f, e := strconv.ParseFloat(result, 64)
				if e != nil {
					panic("[" + strconv.FormatInt(spellId, 10) + "] Variable description result does not evaluate to float.")
				}
				variables[c.VariableAssignment.VariableName] = f
			}
		}
	}

	value, error := parser.ParseString("", tooltip)
	if error != nil {
		return nil, error

	}

	return &Tooltip{AST: value, Context: ctx}, nil
}
