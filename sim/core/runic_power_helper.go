package core

import (
	"fmt"
)

// RuneCost's bit layout is: <16r.4d.4u.4f.4b>. Each part is just a count now (0..15 for runes).
type RuneCost int32

func NewRuneCost(rp int16, blood, frost, unholy, death int8) RuneCost {
	return RuneCost(rp)<<16 |
		RuneCost(death&0xf)<<12 |
		RuneCost(unholy&0xf)<<8 |
		RuneCost(frost&0xf)<<4 |
		RuneCost(blood&0xf)
}

func (rc RuneCost) String() string {
	return fmt.Sprintf("RP: %d, Blood: %d, Frost: %d, Unholy: %d, Death: %d", rc.RunicPower(), rc.Blood(), rc.Frost(), rc.Unholy(), rc.Death())
}

// HasRune returns if this cost includes a rune portion.
func (rc RuneCost) HasRune() bool {
	return rc&0xffff > 0
}

func (rc RuneCost) RunicPower() int16 {
	return int16(rc >> 16)
}

func (rc RuneCost) Blood() int8 {
	return int8(rc & 0xf)
}

func (rc RuneCost) Frost() int8 {
	return int8((rc >> 4) & 0xf)
}

func (rc RuneCost) Unholy() int8 {
	return int8((rc >> 8) & 0xf)
}

func (rc RuneCost) Death() int8 {
	return int8((rc >> 12) & 0xf)
}
