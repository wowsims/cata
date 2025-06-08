package dbc

type ItemArmorQuality struct {
	ItemLevel int
	Quality   []float64 // len 7
}

type ItemArmorShield struct {
	ItemLevel int
	Quality   []float64 // len 7
}

type ItemArmorTotal struct {
	ItemLevel int
	Cloth     float64
	Leather   float64
	Mail      float64
	Plate     float64
}

type ArmorLocation struct {
	Id       int
	Modifier [5]float64
}
