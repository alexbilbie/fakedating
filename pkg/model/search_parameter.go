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
		if radius < 1 {
			radius = 1
		}
		if radius > 500 {
			radius = 500
		}
		params.Latitude = latitude
		params.Longitude = longitude
		params.Radius = radius
	}
}

func SearchWithAgeConstraint(lower uint, upper uint) SearchParameterOpt {
	return func(params *SearchParameter) {
		params.AgeLower = lower
		params.AgeUpper = upper
	}
}

func SearchWithOffset(offset uint) SearchParameterOpt {
	return func(params *SearchParameter) {
		params.Offset = offset
	}
}
