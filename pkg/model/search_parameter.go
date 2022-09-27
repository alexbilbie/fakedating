package model

type SearchParameter struct {
	Latitude, Longitude float64
	Radius              uint
	AgeLower, AgeUpper  uint
	Offset              uint
}

type SearchParameterOpt = func(*SearchParameter)

func SearchWithLocationConstraint(latitude, longitude float64, radius uint) SearchParameterOpt {
	return func(params *SearchParameter) {
		if radius < 25 {
			radius = 25
		}
		if radius > 500 {
			radius = 500
		}
		params.Latitude = latitude
		params.Longitude = longitude
	}
}

func SearchWithAgeConstraint(lower uint, upper uint) SearchParameterOpt {
	return func(params *SearchParameter) {
		params.AgeLower = lower
		params.AgeUpper = upper
	}
}
