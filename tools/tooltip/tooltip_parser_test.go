package tooltip

import (
	"testing"

	"github.com/wowsims/mop/tools/database/dbc"
)

var db = dbc.GetDBC()

func Test_WhenInvalidTernaryGiven_ThenProperlyApplyFixes(t *testing.T) {
	tp, error := ParseTooltip("$<dam> damage every ${$16914d3/10}.2 seconds$?$w1!=0[ and movement slowed by $w1%][].",
		NewTestDataProvider(CharacterConfig{SpellPower: 1000}),
		16914,
	)

	if error != nil {
		t.Fatal()
	}

	if _, ok := tp.Context.Variables["dam"]; !ok {
		t.Fail()
	}

	if tp.String() != "6624 damage every 1.00 seconds." {
		t.Fail()
	}
}

func Test_WhenLConditionGiven_ThenProperlyEvaluate(t *testing.T) {
	SimpleTooltipTest(55685, "The first charge of your Prayer of Mending heals for an additional 60% but your Prayer of Mending has 1 fewer charge.", t)
}

func Test_WhenGivenDurationShorter60_ThenRenderSeconds(t *testing.T) {
	SimpleTooltipTest(54810,
		"For 6s after activating Frenzied Regeneration, healing effects on you are 40% more powerful. However, your Frenzied Regeneration now always costs 50 Rage and no longer converts Rage into health.",
		t,
	)
}

func Test_WhenGivenDurationLongerThan2Hours_ThenRnderHrs(t *testing.T) {
	SimpleTooltipTest(56382, "When cast on critters, your Polymorph spells now last 24hrs and can be cast on multiple targets.", t)
}

func SimpleTooltipTest(spellId int, expectedDescription string, t *testing.T) {
	spell := db.Spells[spellId]
	tp, error := ParseTooltip(spell.Description,
		NewTestDataProvider(CharacterConfig{}),
		int64(spellId),
	)

	if error != nil {
		t.Fatal()
	}

	if tp.String() != expectedDescription {
		t.Fail()
	}
}

func NewTestDataProvider(config CharacterConfig) *TestDataProvider {
	return &TestDataProvider{
		DBCTooltipDataProvider: &DBCTooltipDataProvider{
			DBC: db,
		},
		Character: &config,
	}
}

// add here over time to overwrite fixed values for tests
type CharacterConfig struct {
	SpellPower float64
}

type TestDataProvider struct {
	*DBCTooltipDataProvider
	Character *CharacterConfig
}

func (t TestDataProvider) GetSpellPower() float64 {
	return t.Character.SpellPower
}
