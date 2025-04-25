package warrior

type SpellEffectWeaponDmgPctConfig struct {
	BaseWeapon_Pct    float64
	Coefficient       float64
	EffectPerLevel    float64
	BaseSpellLevel    int32
	MaxSpellLevel     int32
	ClassSpellScaling float64
}

func (config *SpellEffectWeaponDmgPctConfig) CalcSpellDamagePct() float64 {
	return config.BaseWeapon_Pct + (config.EffectPerLevel * (float64(config.MaxSpellLevel - config.BaseSpellLevel)) / 100)
}

func (config *SpellEffectWeaponDmgPctConfig) CalcAddedSpellDamage() float64 {
	return config.Coefficient * config.ClassSpellScaling
}
