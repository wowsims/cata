package core

import (
	"sync"

	"github.com/wowsims/cata/sim/core/dbc"
)

var (
	once     sync.Once
	instance *SpellGen
)

type SpellGen struct {
	dbc *dbc.DBC
}

func GetInstance() *SpellGen {
	once.Do(func() {
		instance = &SpellGen{}
	})
	return instance
}

func (s *SpellGen) GetDBC() *dbc.DBC {
	if s.dbc == nil {
		s.dbc = dbc.NewDBC()
	}
	return s.dbc
}

func (sg *SpellGen) ParseSpellData(spellId uint, unit *Unit) *SpellConfig {
	dbc := sg.GetDBC()
	dbcSpell := dbc.FetchSpell(spellId)
	s := SpellConfig{}
	if !dbcSpell.Ok() {
		return nil
	}
	s.SpellID = int32(dbcSpell.ID)
	s.MissileSpeed = dbcSpell.PrjSpeed
	s.SpellSchool = SpellSchool(dbcSpell.School) // Todo? Does this match 1 to 1?
	//s.TicksCanCrit = dbcSpell.Flags()
	s.Cast = CastConfig{
		DefaultCast: Cast{
			GCD:      dbcSpell.GCD,
			CastTime: dbcSpell.CastTime,
		},
		CD: Cooldown{
			Duration: dbcSpell.Cooldown, //Todo: There's also CategoryCooldown and Category
			Timer:    unit.NewTimer(),   // ??
		},
	}
	if dbcSpell.HasDirectDamageEffect() {
		//
	}
	if dbcSpell.HasPeriodicDamageEffect() {
		s.Dot = DotConfig{}
	}
	return &s
}
func (sg *SpellGen) ParseEffects(dbcSpell *dbc.SpellData, spellConfig *SpellConfig, unit *Unit) {
	effects := dbcSpell.Effects
	for i := 0; i < int(dbcSpell.EffectsCount); i++ {
		effect := effects[i]
		switch effect.Type {
		case dbc.E_SCHOOL_DAMAGE:
		case dbc.E_HEALTH_LEECH:
			//parse direct dmg mod
			break
		case dbc.E_NORMALIZED_WEAPON_DMG:
			//set normalised wpn dmg
		case dbc.E_WEAPON_DAMAGE:
			// normal wpn dmg
		case dbc.E_WEAPON_PERCENT_DAMAGE:
			// wpn prct dmg
		case dbc.E_PERSISTENT_AREA_AURA:
			//
		case dbc.E_APPLY_AURA:
			//
			switch effect.Subtype {
			case dbc.A_PERIODIC_DAMAGE:
			case dbc.A_PERIODIC_LEECH:
				//parse effect periodic mods
				//keep going
			}

		}
	}
}
func (sg *SpellGen) ParseDirectEffect(effect *dbc.SpellEffectData, spellConfig *SpellConfig, unit *Unit) {

	// spCoeff := effect.SPCoeff
	// apCoeff := effect.APCoeff
	// delta := effect.MDelta

	// min := effect.Min(sg.GetDBC(), 85, 85)
	// max := effect.Max(sg.GetDBC(), 85, 85)
	// maxRadius := effect.GetRadiusMax()
	//Todo: Need a way to apply these to the spell

}

func (sg *SpellGen) ParsePeriodicEffect(effect *dbc.SpellEffectData, spellConfig *SpellConfig, unit *Unit) {
	// spCoeff := effect.SPCoeff
	// apCoeff := effect.APCoeff
	// tickDamage := effect.Average(sg.GetDBC(), 85, 85)

	// if effect.Amplitude > 0 {
	// 	tickTime := effect.Amplitude
	// 	dotDuration := effect.GetSpell(sg.GetDBC()).Duration
	// }
}
func (sg *SpellGen) GetOutCome(dbcSspellpell *dbc.SpellData) HitOutcome {
	return OutcomeHit
}

const (
	SX_RANGED_ABILITY                 uint = 1
	SX_ABILITY                        uint = 4
	SX_TRADESKILL_ABILITY             uint = 5
	SX_PASSIVE                        uint = 6
	SX_HIDDEN                         uint = 7
	SX_REQ_STEALTH                    uint = 17
	SX_CANCEL_AUTO_ATTACK             uint = 20
	SX_NO_D_P_B                       uint = 21
	SX_NO_COMBAT                      uint = 22
	SX_NO_CANCEL                      uint = 31
	SX_CHANNELED                      uint = 34
	SX_NO_STEALTH_BREAK               uint = 37
	SX_CHANNELED_2                    uint = 38
	SX_MELEE_COMBAT_START             uint = 41
	SX_NO_THREAT                      uint = 42
	SX_DISCOUNT_ON_MISS               uint = 59
	SX_DONT_DISPLAY_IN_AURA_BAR       uint = 60
	SX_CANNOT_CRIT                    uint = 93
	SX_FOOD_AURA                      uint = 95
	SX_NOT_A_PROC                     uint = 105
	SX_REQ_MAIN_HAND                  uint = 106
	SX_DISABLE_PLAYER_PROCS           uint = 112
	SX_DISABLE_TARGET_PROCS           uint = 113
	SX_ALWAYS_HIT                     uint = 114
	SX_REQ_OFF_HAND                   uint = 120
	SX_TREAT_AS_PERIODIC              uint = 121
	SX_CAN_PROC_FROM_PROCS            uint = 122
	SX_DISABLE_TARGET_MULT            uint = 136
	SX_DISABLE_WEAPON_PROCS           uint = 151
	SX_TICK_ON_APPLICATION            uint = 169
	SX_DOT_HASTED                     uint = 173
	SX_NO_PARTIAL_RESISTS             uint = 183
	SX_REQ_LINE_OF_SIGHT              uint = 186
	SX_IGNORE_FOR_MOD_TIME_RATE       uint = 196
	SX_DISABLE_PLAYER_MULT            uint = 221
	SX_NO_DODGE                       uint = 247
	SX_NO_PARRY                       uint = 248
	SX_NO_MISS                        uint = 249
	SX_NO_BLOCK                       uint = 256
	SX_TICK_MAY_CRIT                  uint = 265
	SX_DURATION_HASTED                uint = 273
	SX_DOT_HASTED_MELEE               uint = 278
	SX_MASTERY_AFFECTS_POINTS         uint = 285
	SX_FIXED_TRAVEL_TIME              uint = 292
	SX_DISABLE_PLAYER_HEALING_MULT    uint = 312
	SX_DISABLE_TARGET_POSITIVE_MULT   uint = 321
	SX_TARGET_SPECIFIC_COOLDOWN       uint = 330
	SX_ROLLING_PERIODIC               uint = 334
	SX_SCALE_ILEVEL                   uint = 354
	SX_ONLY_PROC_FROM_CLASS_ABILITIES uint = 415
	SX_ALLOW_CLASS_ABILITY_PROCS      uint = 416
	SX_REFRESH_EXTENDS_DURATION       uint = 436
)
