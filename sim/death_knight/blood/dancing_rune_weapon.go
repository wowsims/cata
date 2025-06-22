package blood

import (
	"time"

	"github.com/wowsims/mop/sim/core"
	"github.com/wowsims/mop/sim/core/proto"
	"github.com/wowsims/mop/sim/death_knight"
)

/*
Summons a second rune weapon that fights on its own for 12 sec, mirroring the Death Knight's attacks.
The rune weapon also assists in defense of its master, granting an additional 20% parry chance

-- Glyph of Dancing Rune Weapon --

and increasing threat generation by 100%

-- /Glyph of Dancing Rune Weapon --

while active.
*/
func (dk *BloodDeathKnight) registerDancingRuneWeapon() {
	duration := time.Second * 12

	hasGlyph := dk.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfDancingRuneWeapon)

	t124PAura := dk.RegisterAura(core.Aura{
		Label:    "Flaming Rune Weapon" + dk.Label,
		ActionID: core.ActionID{SpellID: 101162},
		Duration: duration,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance += 0.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.PseudoStats.BaseParryChance -= 0.15
		},
	})

	dancingRuneWeaponAura := dk.RegisterAura(core.Aura{
		Label:    "Dancing Rune Weapon" + dk.Label,
		ActionID: core.ActionID{SpellID: 81256},
		Duration: duration,
		// Casts copy
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			copySpell := dk.RuneWeapon.RuneWeaponSpells[spell.ActionID]
			if copySpell == nil {
				return
			}

			death_knight.CopySpellMultipliers(spell, copySpell, dk.CurrentTarget)

			copySpell.Cast(sim, dk.CurrentTarget)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if dk.T12Tank4pc.IsActive() {
				t124PAura.Activate(sim)
			}
		},
	}).AttachAdditivePseudoStatBuff(&dk.PseudoStats.BaseParryChance, 0.2)

	if hasGlyph {
		dancingRuneWeaponAura.AttachMultiplicativePseudoStatBuff(&dk.PseudoStats.ThreatMultiplier, 2.0)
	}

	spell := dk.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 49028},
		Flags:          core.SpellFlagAPL | core.SpellFlagReadinessTrinket,
		ClassSpellMask: death_knight.DeathKnightSpellDancingRuneWeapon,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				NonEmpty: true,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 90,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.RuneWeapon.EnableWithTimeout(sim, dk.RuneWeapon, duration)
			dk.RuneWeapon.CancelGCDTimer(sim)
			spell.RelatedSelfBuff.Activate(sim)
		},

		RelatedSelfBuff: dancingRuneWeaponAura,
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
