package dbc

import (
	"github.com/wowsims/cata/sim/core/proto"
)

type ItemArmorQuality struct {
	ItemLevel int
	Quality   []float64 // len 7
}

func (iat *ItemArmorQuality) ToProto() *proto.ItemQualityValue {
	if len(iat.Quality) != 7 {
		return &proto.ItemQualityValue{
			ItemLevel: int32(iat.ItemLevel),
		}
	}

	return &proto.ItemQualityValue{
		ItemLevel: int32(iat.ItemLevel),
		Quality: &proto.QualityValues{
			Common:    iat.Quality[0],
			Uncommon:  iat.Quality[1],
			Rare:      iat.Quality[2],
			Epic:      iat.Quality[3],
			Legendary: iat.Quality[4],
			Artifact:  iat.Quality[5],
			Heirloom:  iat.Quality[6],
		},
	}
}

type ItemArmorShield struct {
	ItemLevel int
	Quality   []float64 // len 7
}

func (iat *ItemArmorShield) ToProto() *proto.ItemQualityValue {
	if len(iat.Quality) != 7 {
		return &proto.ItemQualityValue{
			ItemLevel: int32(iat.ItemLevel),
		}
	}

	return &proto.ItemQualityValue{
		ItemLevel: int32(iat.ItemLevel),
		Quality: &proto.QualityValues{
			Common:    iat.Quality[0],
			Uncommon:  iat.Quality[1],
			Rare:      iat.Quality[2],
			Epic:      iat.Quality[3],
			Legendary: iat.Quality[4],
			Artifact:  iat.Quality[5],
			Heirloom:  iat.Quality[6],
		},
	}
}

type ItemArmorTotal struct {
	ItemLevel int
	Cloth     float64
	Leather   float64
	Mail      float64
	Plate     float64
}

func (iat *ItemArmorTotal) ToProto() *proto.ItemArmorTotal {
	return &proto.ItemArmorTotal{
		ItemLevel: int32(iat.ItemLevel),
		Cloth:     iat.Cloth,
		Mail:      iat.Mail,
		Leather:   iat.Leather,
		Plate:     iat.Plate,
	}
}

type ArmorLocation struct {
	Id       int
	Modifier [5]float64
}
