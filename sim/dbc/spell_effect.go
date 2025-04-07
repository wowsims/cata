package dbc

type SpellEffect struct {
	ID                             int
	DifficultyID                   int
	EffectIndex                    int
	EffectType                     SpellEffectType
	EffectAmplitude                float64
	EffectAttributes               int
	EffectAura                     EffectAuraType
	EffectAuraPeriod               int
	EffectBasePoints               int
	EffectBonusCoefficient         float64
	EffectChainAmplitude           float64
	EffectChainTargets             int
	EffectDieSides                 int
	EffectItemType                 int
	EffectMechanic                 int
	EffectPointsPerResource        float64
	EffectPosFacing                float64
	EffectRealPointsPerLevel       float64
	EffectTriggerSpell             int
	BonusCoefficientFromAP         float64
	PvpMultiplier                  float64
	Coefficient                    float64
	Variance                       float64
	ResourceCoefficient            float64
	GroupSizeBasePointsCoefficient float64
	// Grouped properties parsed from JSON strings:
	EffectMiscValues      []int // from EffectMiscValue, EffectMiscValue_0, EffectMiscValue_1
	EffectRadiusIndices   []int // from EffectRadiusIndex, EffectRadiusIndex_0, EffectRadiusIndex_1
	EffectSpellClassMasks []int // from EffectSpellClassMask, EffectSpellClassMask_0, EffectSpellClassMask_1, EffectSpellClassMask_2, EffectSpellClassMask_3
	ImplicitTargets       []int // from ImplicitTarget, ImplicitTarget_0, ImplicitTarget_1
	SpellID               int
}
